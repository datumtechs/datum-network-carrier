package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/db"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	"github.com/RosettaFlow/Carrier-Go/lib/center/api"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"sync/atomic"
)

// DataCenter is mainly responsible for communicating with the data center service
type DataCenter struct {
	ctx       context.Context
	config    *params.DataCenterConfig
	client    *grpclient.GrpcClient
	mu        sync.RWMutex // global mutex for locking data center operations.
	serviceMu sync.RWMutex // data processor lock

	db db.Database // Low level persistent database to store final content.

	processor     Processor      // block processor interface
	running       int32          // running must be called atomically
	procInterrupt int32          // interrupt signaler for block processing
	wg            sync.WaitGroup // chain processing wait group for shutting down
}

// NewDataCenter returns a fully initialised data center using information available in the database.
func NewDataCenter(ctx context.Context, db db.Database, config *params.DataCenterConfig) (*DataCenter, error) {
	if config.GrpcUrl == "" || config.Port == 0 {
		panic("Invalid Grpc Config.")
	}
	client, err := grpclient.NewGrpcClient(ctx, fmt.Sprintf("%v:%v", config.GrpcUrl, config.Port))
	if err != nil {
		log.WithError(err).Error("dial grpc server failed")
		return nil, err
	}
	dc := &DataCenter{
		ctx:    ctx,
		config: config,
		client: client,
		db:     db,
	}
	return dc, nil
}

func (dc *DataCenter) getProcInterrupt() bool {
	return atomic.LoadInt32(&dc.procInterrupt) == 1
}

func (dc *DataCenter) SetProcessor(processor Processor) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	dc.processor = processor
}

func (dc *DataCenter) GrpcClient() *grpclient.GrpcClient {
	return dc.client
}

// ************************************* public api (datachain) *******************************************

func (dc *DataCenter) GetYarnName() (string, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	name, err :=  rawdb.ReadYarnName(dc.db)
	if nil != err {
		return "", err
	}
	return name, nil
}

// InsertChain saves the data of block to the database.
func (dc *DataCenter) InsertData(blocks types.Blocks) (int, error) {
	if len(blocks) == 0 {
		return 0, nil
	}
	// check.
	for i := 0; i < len(blocks); i++ {
		if blocks[i].NumberU64() != blocks[i-1].NumberU64()+1 || blocks[i].ParentHash() != blocks[i-1].Hash() {
			log.WithFields(logrus.Fields{
				"number": blocks[i].NumberU64(),
				"hash":   blocks[i].Hash(),
			}).Error("Non contiguous block insert")
			return 0, fmt.Errorf("non contiguous insert: item %d is #%d", i-1, blocks[i-1].NumberU64())
		}
	}
	// pre-checks passed, start the full block imports
	dc.wg.Add(1)
	defer dc.wg.Done()

	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()

	headers := make([]*types.Header, len(blocks))
	seals := make([]bool, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
		seals[i] = true
	}
	for i, block := range blocks {
		if atomic.LoadInt32(&dc.procInterrupt) == 1 {
			log.Debug("Premature abort during blocks processing")
			break
		}
		err := dc.processor.Process(block, dc.config)
		if err != nil {
			// for err, how to deal with????
			return i, err
		}
	}
	return len(blocks), nil
}

// on yarn node api
func (dc *DataCenter) SetSeedNode(seed *types.SeedNodeInfo) (types.NodeConnStatus, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.WriteSeedNodes(dc.db, seed)
	return types.NONCONNECTED, nil
}

func (dc *DataCenter) DeleteSeedNode(id string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.DeleteSeedNode(dc.db, id)
	return nil
}

func (dc *DataCenter) GetSeedNode(id string) (*types.SeedNodeInfo, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.ReadSeedNode(dc.db, id)
}

func (dc *DataCenter) GetSeedNodeList() ([]*types.SeedNodeInfo, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.ReadAllSeedNodes(dc.db)
}

func (dc *DataCenter) SetRegisterNode(typ types.RegisteredNodeType, node *types.RegisteredNodeInfo) (types.NodeConnStatus, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.WriteRegisterNodes(dc.db, typ, node)
	return node.ConnState, nil
}

func (dc *DataCenter) DeleteRegisterNode(typ types.RegisteredNodeType, id string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.DeleteRegisterNode(dc.db, typ, id)
	return nil
}

func (dc *DataCenter) GetRegisterNode(typ types.RegisteredNodeType, id string) (*types.RegisteredNodeInfo, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.ReadRegisterNode(dc.db, typ, id)
}

func (dc *DataCenter) GetRegisterNodeList(typ types.RegisteredNodeType) ([]*types.RegisteredNodeInfo, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.ReadAllRegisterNodes(dc.db, typ)
}

// about metaData
// on datecenter
func (dc *DataCenter) InsertMetadata(metadata *types.Metadata) error {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	response, err := dc.client.SaveMetaData(dc.ctx, types.NewMetaDataSaveRequest(metadata))
	if err != nil {
		log.WithError(err).WithField("hash", metadata.Hash()).Errorf("InsertMetadata failed")
		return err
	}
	if response.Status != 0 {
		return fmt.Errorf("insert metadata error: %s", response.Msg)
	}
	return nil
}

func (dc *DataCenter) RevokeMetadata(metadata *types.Metadata) error {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	response, err := dc.client.RevokeMetaData(dc.ctx, types.NewMetaDataRevokeRequest(metadata))
	if err != nil {
		log.WithError(err).WithField("hash", metadata.Hash()).Errorf("RevokeMetadata failed")
		return err
	}
	if response.Status != 0 {
		return fmt.Errorf("revoke metadata error: %s", response.Msg)
	}
	return nil
}

func (dc *DataCenter) GetMetadataByDataId(dataId string) (*types.Metadata, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	metadataByIdResponse, err := dc.client.GetMetadataById(dc.ctx, &api.MetadataByIdRequest{
		MetadataId: dataId,
	})
	return types.NewMetadataFromResponse(metadataByIdResponse), err
}

func (dc *DataCenter) GetMetadataList() (types.MetadataArray, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	metaDataListResponse, err := dc.client.GetMetadataList(dc.ctx, &api.MetadataListRequest{
		LastUpdateTime:      uint64( timeutils.Now().Unix()),
	})
	return types.NewMetadataArrayFromDetailListResponse(metaDataListResponse), err
}

// about power on local
func (dc *DataCenter) InsertLocalResource(resource *types.LocalResource) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.WriteLocalResource(dc.db, resource)
	return nil
}

func (dc *DataCenter) RemoveLocalResource(jobNodeId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.DeleteLocalResource(dc.db, jobNodeId)
	return nil
}

func (dc *DataCenter) GetLocalResource(jobNodeId string) (*types.LocalResource, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.ReadLocalResource(dc.db, jobNodeId)
}

func (dc *DataCenter) GetLocalResourceList() (types.LocalResourceArray, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.ReadAllLocalResource(dc.db)
}

func (dc *DataCenter) StoreLocalResourceIdByPowerId(powerId, jobNodeId string) error {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.StoreLocalResourceIdByPowerId(dc.db, powerId, jobNodeId)
}

func (dc *DataCenter) RemoveLocalResourceIdByPowerId(powerId string) error  {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.RemoveLocalResourceIdByPowerId(dc.db, powerId)
}

func (dc *DataCenter) QueryLocalResourceIdByPowerId(powerId string) (string, error)  {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryLocalResourceIdByPowerId(dc.db, powerId)
}

//func (dc *DataCenter)  StoreLocalResourceIdByMetaDataId(metaDataId, dataNodeId string) error {
//	dc.mu.RLock()
//	defer dc.mu.RUnlock()
//	return rawdb.StoreLocalResourceIdByMetaDataId(dc.db, metaDataId, dataNodeId)
//}
//func (dc *DataCenter)  RemoveLocalResourceIdByMetaDataId(metaDataId string) error {
//	dc.mu.RLock()
//	defer dc.mu.RUnlock()
//	return rawdb.RemoveLocalResourceIdByMetaDataId(dc.db, metaDataId)
//}
//func (dc *DataCenter)  QueryLocalResourceIdByMetaDataId(metaDataId string) (string, error) {
//	dc.mu.RLock()
//	defer dc.mu.RUnlock()
//	return rawdb.QueryLocalResourceIdByMetaDataId(dc.db, metaDataId)
//}

// about power on datacenter
func (dc *DataCenter) InsertResource(resource *types.Resource) error {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	response, err := dc.client.SaveResource(dc.ctx, types.NewPublishPowerRequest(resource))
	if err != nil {
		log.WithError(err).WithField("hash", resource.Hash()).Errorf("InsertResource failed")
		return err
	}
	if response.Status != 0 {
		return fmt.Errorf("insert resource error: %s", response.Msg)
	}
	return nil
}


func (dc *DataCenter) RevokeResource(resource *types.Resource) error {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	response, err := dc.client.RevokeResource(dc.ctx, types.RevokePowerRequest(resource))
	if err != nil {
		log.WithError(err).WithField("hash", resource.Hash()).Errorf("RevokeResource failed")
		return err
	}
	if response.Status != 0 {
		return fmt.Errorf("revoke resource error: %s", response.Msg)
	}
	return nil
}


func (dc *DataCenter) SyncPowerUsed (resource *types.LocalResource) error {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	response, err := dc.client.SyncPower(dc.ctx, types.NewSyncPowerRequest(resource))
	if err != nil {
		log.WithError(err).WithField("hash", resource.Hash()).Errorf("SyncPowerUsed failed")
		return err
	}
	if response.Status != 0 {
		return fmt.Errorf("sync resource used error: %s", response.Msg)
	}
	return nil
}


func (dc *DataCenter) GetResourceListByIdentityId(identityId string) (types.ResourceArray, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	powerTotalSummaryResponse, err := dc.client.GetPowerSummaryByIdentityId(dc.ctx, &api.PowerSummaryByIdentityRequest{
		IdentityId: identityId,
	})
	return types.NewResourceFromResponse(powerTotalSummaryResponse), err
}

func (dc *DataCenter) GetResourceList() (types.ResourceArray, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	powerListRequest, err := dc.client.GetPowerTotalSummaryList(dc.ctx)
	return types.NewResourceArrayFromPowerTotalSummaryListResponse(powerListRequest), err
}

// about identity on local
func (dc *DataCenter) StoreIdentity(identity *types.NodeAlias) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.WriteLocalIdentity(dc.db, identity)
	return nil
}

func (dc *DataCenter) RemoveIdentity() error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.DeleteLocalIdentity(dc.db)
	return nil
}

func (dc *DataCenter) GetIdentityId() (string, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	identity, err := rawdb.ReadLocalIdentity(dc.db)
	if nil != err {
		return "", err
	}
	return identity.GetNodeIdentityId(), nil
}

func (dc *DataCenter) GetIdentity() (*types.NodeAlias, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	identity, err := rawdb.ReadLocalIdentity(dc.db)
	if nil != err {
		return nil, err
	}
	return identity, nil
}

// about identity on datacenter
func (dc *DataCenter) HasIdentity(identity *types.NodeAlias) (bool, error) {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()
	responses, err := dc.client.GetIdentityList(dc.ctx, &api.IdentityListRequest{
		LastUpdateTime: uint64(timeutils.Now().Second()),
	})
	if err != nil {
		return false, err
	}
	for _, organization := range responses.IdentityList {
		if strings.EqualFold(organization.IdentityId, identity.IdentityId) {
			return true, nil
		}
	}
	return false, nil
}

func (dc *DataCenter) InsertIdentity(identity *types.Identity) error {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	response, err := dc.client.SaveIdentity(dc.ctx, types.NewSaveIdentityRequest(identity))
	if err != nil {
		//log.WithError(err).WithField("hash", identity.Hash()).Errorf("InsertIdentity failed")
		return err
	}
	if response.Status != 0 {
		return fmt.Errorf("insert indentity error: %s", response.Msg)
	}
	return nil
}

func (dc *DataCenter) RevokeIdentity(identity *types.Identity) error {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	response, err := dc.client.RevokeIdentityJoin(dc.ctx, &api.RevokeIdentityJoinRequest{
		Member: &api.Organization{
			Name:       identity.Name(),
			NodeId:     identity.NodeId(),
			IdentityId: identity.IdentityId(),
		},
	})
	if err != nil {
		return err
	}
	if response.GetStatus() != 0 {
		return fmt.Errorf("revokeIdeneity err: %s", response.GetMsg())
	}
	return nil
}

func (dc *DataCenter) GetIdentityList() (types.IdentityArray, error) {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()
	identityListResponse, err := dc.client.GetIdentityList(dc.ctx, &api.IdentityListRequest{LastUpdateTime: uint64(timeutils.UnixMsec())})
	return types.NewIdentityArrayFromIdentityListResponse(identityListResponse), err
}

//func (dc *DataCenter) GetIdentityListByIds(identityIds []string) (types.IdentityArray, error) {
//	dc.serviceMu.RLock()
//	defer dc.mu.RUnlock()
//
//	return nil, nil
//}

// about task on local
// local task
func (dc *DataCenter) StoreLocalTask(task *types.Task) error {
	if task == nil {
		return errors.New("invalid params for task")
	}
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.WriteLocalTask(dc.db, task)
	return nil
}

func (dc *DataCenter) RemoveLocalTask(taskId string) error {
	if taskId == "" {
		return errors.New("invalid params for taskId to DelLocalTask")
	}
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.DeleteLocalTask(dc.db, taskId)
	return nil
}

func (dc *DataCenter) UpdateLocalTaskState(taskId, state string) error {
	if taskId == "" || state == "" {
		return errors.New("invalid params taskId or state for UpdateLocalTaskState")
	}
	dc.mu.Lock()
	defer dc.mu.Unlock()
	task, err := rawdb.ReadLocalTask(dc.db, taskId)
	if nil != err {
		return err
	}
	task.TaskData().State = state
	rawdb.DeleteLocalTask(dc.db, taskId)
	rawdb.WriteLocalTask(dc.db, task)
	return nil
}

func (dc *DataCenter) GetLocalTask(taskId string) (*types.Task, error) {
	if taskId == "" {
		return nil, errors.New("invalid params taskId for GetLocalTask")
	}
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.ReadLocalTask(dc.db, taskId)
}

func (dc *DataCenter) GetLocalTaskListByIds(taskIds []string) (types.TaskDataArray, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.ReadLocalTaskByIds(dc.db, taskIds)
}

func (dc *DataCenter) GetLocalTaskList() (types.TaskDataArray, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.ReadAllLocalTasks(dc.db)
}

func (dc *DataCenter) StoreJobNodeRunningTaskId(jobNodeId, taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreResourceTaskId(dc.db, jobNodeId, taskId)
}

func (dc *DataCenter) RemoveJobNodeRunningTaskId(jobNodeId, taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveResourceTaskId(dc.db, jobNodeId, taskId)
}

func (dc *DataCenter) GetRunningTaskCountOnJobNode(jobNodeId string) (uint32, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	taskIds, err := rawdb.QueryResourceTaskIds(dc.db, jobNodeId)
	if nil != err {
		return 0, err
	}
	return uint32(len(taskIds)), nil
}

func (dc *DataCenter) GetJobNodeRunningTaskIdList(jobNodeId string) ([]string, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryResourceTaskIds(dc.db, jobNodeId)
}

// about task on datacenter
func (dc *DataCenter) InsertTask(task *types.Task) error {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	response, err := dc.client.SaveTask(dc.ctx, types.NewTaskDetail(task))
	if err != nil {
		log.WithError(err).WithField("taskId", task.TaskId()).Errorf("InsertTask failed")
		return err
	}
	if response.Status != 0 {
		return fmt.Errorf("insert task, taskId: {%s},  error: %s", task.TaskId(), response.Msg)
	}
	return nil
}

func (dc *DataCenter) GetTaskListByIdentityId(identityId string) (types.TaskDataArray, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	//taskListResponse, err := dc.client.ListTask(dc.ctx, &api.TaskListRequest{LastUpdateTime: uint64(timeutils.UnixMsec())})
	taskListResponse, err := dc.client.ListTaskByIdentity(dc.ctx, &api.TaskListByIdentityRequest{
		LastUpdateTime: uint64(timeutils.UnixMsec()),
		IdentityId: identityId,
	})
	return types.NewTaskArrayFromResponse(taskListResponse), err
}

func (dc *DataCenter) GetRunningTaskCountOnOrg() uint32 {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	taskList, err := rawdb.ReadAllLocalTasks(dc.db)
	if nil != err {
		return 0
	}
	if taskList != nil {
		return uint32(taskList.Len())
	}
	return 0
}

func (dc *DataCenter) GetTaskEventListByTaskId(taskId string) ([]*api.TaskEvent, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	taskEventResponse, err := dc.client.ListTaskEvent(dc.ctx, &api.TaskEventRequest{
		TaskId: taskId,
	})
	return taskEventResponse.TaskEventList, err
}

func (dc *DataCenter) GetTaskEventListByTaskIds(taskIds []string) ([]*api.TaskEvent, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()

	eventList := make([]*api.TaskEvent, 0)
	for _, taskId := range taskIds {
		taskEventResponse, err := dc.client.ListTaskEvent(dc.ctx, &api.TaskEventRequest{
			TaskId: taskId,
		})
		if nil != err {
			return nil, err
		}
		eventList = append(eventList, taskEventResponse.TaskEventList...)
	}
	return eventList, nil
}

// For ResourceManager
// about jobRerource
func (dc *DataCenter) StoreLocalResourceTable(resource *types.LocalResourceTable) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreNodeResource(dc.db, resource)
}

func (dc *DataCenter) StoreLocalResourceTables(resources []*types.LocalResourceTable) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreNodeResources(dc.db, resources)
}

func (dc *DataCenter) RemoveLocalResourceTable(resourceId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveNodeResource(dc.db, resourceId)
}

func (dc *DataCenter) QueryLocalResourceTable(resourceId string) (*types.LocalResourceTable, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryNodeResource(dc.db, resourceId)
}

func (dc *DataCenter) QueryLocalResourceTables() ([]*types.LocalResourceTable, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryNodeResources(dc.db)
}

// about Org power resource
func (dc *DataCenter) StoreOrgResourceTable(resource *types.RemoteResourceTable) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreOrgResource(dc.db, resource)
}

func (dc *DataCenter) StoreOrgResourceTables(resources []*types.RemoteResourceTable) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreOrgResources(dc.db, resources)
}

func (dc *DataCenter) RemoveOrgResourceTable(resourceId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveOrgResource(dc.db, resourceId)
}

func (dc *DataCenter) QueryOrgResourceTable(resourceId string) (*types.RemoteResourceTable, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryOrgResource(dc.db, resourceId)
}

func (dc *DataCenter) QueryOrgResourceTables() ([]*types.RemoteResourceTable, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryOrgResources(dc.db)
}

// about slotUnit
func (dc *DataCenter) StoreNodeResourceSlotUnit(slot *types.Slot) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreNodeResourceSlotUnit(dc.db, slot)
}

func (dc *DataCenter) RemoveNodeResourceSlotUnit() error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveNodeResourceSlotUnit(dc.db)
}

func (dc *DataCenter) QueryNodeResourceSlotUnit() (*types.Slot, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryNodeResourceSlotUnit(dc.db)
}

// about TaskPowerUsed
func (dc *DataCenter) StoreLocalTaskPowerUsed(taskPowerUsed *types.LocalTaskPowerUsed) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreLocalTaskPowerUsed(dc.db, taskPowerUsed)
}

func (dc *DataCenter) StoreLocalTaskPowerUseds(taskPowerUseds []*types.LocalTaskPowerUsed) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreLocalTaskPowerUseds(dc.db, taskPowerUseds)
}

func (dc *DataCenter) RemoveLocalTaskPowerUsed(taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveLocalTaskPowerUsed(dc.db, taskId)
}

func (dc *DataCenter) QueryLocalTaskPowerUsed(taskId string) (*types.LocalTaskPowerUsed, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryLocalTaskPowerUsed(dc.db, taskId)
}

func (dc *DataCenter) QueryLocalTaskPowerUseds() ([]*types.LocalTaskPowerUsed, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	res, _ := rawdb.QueryLocalTaskPowerUseds(dc.db)
	return res, nil
}

// about DataResourceTable
func (dc *DataCenter) StoreDataResourceTable(dataResourceTable *types.DataResourceTable) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreDataResourceTable(dc.db, dataResourceTable)
}

func (dc *DataCenter) StoreDataResourceTables(dataResourceTables []*types.DataResourceTable) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreDataResourceTables(dc.db, dataResourceTables)
}

func (dc *DataCenter) RemoveDataResourceTable(nodeId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveDataResourceTable(dc.db, nodeId)
}

func (dc *DataCenter) QueryDataResourceTable(nodeId string) (*types.DataResourceTable, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryDataResourceTable(dc.db, nodeId)
}

func (dc *DataCenter) QueryDataResourceTables() ([]*types.DataResourceTable, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryDataResourceTables(dc.db)
}

// about DataResourceFileUpload
func (dc *DataCenter) StoreDataResourceFileUpload(dataResourceDataUsed *types.DataResourceFileUpload) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreDataResourceFileUpload(dc.db, dataResourceDataUsed)
}

func (dc *DataCenter) StoreDataResourceFileUploads(dataResourceDataUseds []*types.DataResourceFileUpload) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreDataResourceFileUploads(dc.db, dataResourceDataUseds)
}

func (dc *DataCenter) RemoveDataResourceFileUpload(originId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveDataResourceFileUpload(dc.db, originId)
}

func (dc *DataCenter) QueryDataResourceFileUpload(originId string) (*types.DataResourceFileUpload, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryDataResourceFileUpload(dc.db, originId)
}

func (dc *DataCenter) QueryDataResourceFileUploads() ([]*types.DataResourceFileUpload, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryDataResourceFileUploads(dc.db)
}

// about DataResourceDiskUsed
func (dc *DataCenter) StoreDataResourceDiskUsed(dataResourceDiskUsed *types.DataResourceDiskUsed) error {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.StoreDataResourceDiskUsed(dc.db, dataResourceDiskUsed)
}
func (dc *DataCenter) RemoveDataResourceDiskUsed(metaDataId string) error {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.RemoveDataResourceDiskUsed(dc.db, metaDataId)
}
func (dc *DataCenter)  QueryDataResourceDiskUsed(metaDataId string) (*types.DataResourceDiskUsed, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryDataResourceDiskUsed(dc.db, metaDataId)
}


func (dc *DataCenter) StoreTaskEvent(event *types.TaskEventInfo) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.WriteTaskEvent(dc.db, event)
	log.Debugf("Store task eventList, event: %s", event.String())
	return nil
}

func (dc *DataCenter) GetTaskEventList(taskId string) ([]*types.TaskEventInfo, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	list, err := rawdb.ReadTaskEvent(dc.db, taskId)
	if nil != err {
		return nil, err
	}
	log.Debugf("Query task eventList, taskId: {%s}, eventList Len: {%d}", taskId, len(list))
	return list, nil
}


func (dc *DataCenter) RemoveTaskEventList(taskId string) error  {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.DeleteTaskEvent(dc.db, taskId)
	log.Debugf("Remove task eventList, taskId: {%s}", taskId)
	return nil
}

// ****************************************************************************************************************

func (dc *DataCenter) Stop() {
	if !atomic.CompareAndSwapInt32(&dc.running, 0, 1) {
		return
	}
	atomic.StoreInt32(&dc.procInterrupt, 1)
	dc.wg.Wait()
	dc.client.Close()
	log.Info("Datacenter manager stopped")
}


