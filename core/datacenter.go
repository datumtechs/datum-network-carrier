package core

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/db"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	"github.com/RosettaFlow/Carrier-Go/lib/center/api"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/sirupsen/logrus"
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
func NewDataCenter(config *params.DataCenterConfig) (*DataCenter, error) {
	// todo: When to call Close??
	client, err := grpclient.Dial(fmt.Sprintf("%v:%v", config.GrpcUrl, config.Port))
	if err != nil {
		log.WithError(err).Error("dial grpc server failed")
		return nil, err
	}
	dc := &DataCenter{
		config: config,
		client: client,
	}
	return dc, nil
}

func (dc *DataCenter) getProcInterrupt() bool {
	return atomic.LoadInt32(&dc.procInterrupt) == 1
}

func (dc *DataCenter) SetProcessor(processor Processor) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	// do setting...
	dc.processor = processor
}

func (dc *DataCenter) GrpcClient() *grpclient.GrpcClient {
	return dc.client
}

// ************************************* public api (datachain) *******************************************

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

// InsertMetadata saves new metadata info to the center of data.
func (dc *DataCenter) InsertMetadata(metadata *types.Metadata) error {
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

func (dc *DataCenter) StoreYarnName(name string) error {
	rawdb.WriteYarnName(dc.db, name)
	return nil
}

func (dc *DataCenter) GetYarnName() (string, error) {
	return rawdb.ReadYarnName(dc.db), nil
}

func (dc *DataCenter) DelYarnName() error {
	rawdb.DeleteYarnName(dc.db)
	return nil
}

func (dc *DataCenter) StoreIdentityId(identity string) error {
	rawdb.WriteIdentityStr(dc.db, identity)
	return nil
}

func (dc *DataCenter) DelIdentityId() error {
	 rawdb.DeleteIdentityStr(dc.db)
	 return nil
}

func (dc *DataCenter) GetIdentityID() (string, error) {
	return rawdb.ReadIdentityStr(dc.db), nil
}

func (dc *DataCenter) GetMetadataByDataId(dataId string) (*types.Metadata, error) {
	metadataByIdResponse, err := dc.client.GetMetadataById(dc.ctx, &api.MetadataByIdRequest{
		MetadataId:           dataId,
	})
	return types.NewMetadataFromResponse(metadataByIdResponse), err
}

func (dc *DataCenter) GetMetadataListByNodeId(nodeId string) (types.MetadataArray, error) {
	// todo: not need to coding, temporarily.
	return nil, nil
}

func (dc *DataCenter) GetMetadataList() (types.MetadataArray, error) {
	metaDataSummaryListResponse, err := dc.client.GetMetaDataSummaryList(dc.ctx)
	return types.NewMetadataArrayFromResponse(metaDataSummaryListResponse), err
}

// InsertResource saves new resource info to the center of data.
func (dc *DataCenter) InsertResource(resource *types.Resource) error {
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

func (dc *DataCenter) GetResourceByDataId(powerId string) (*types.Resource, error) {
	// todo: not need to coding, temporarily.
	return nil, nil
}

func (dc *DataCenter) GetResourceListByNodeId(nodeId string) (types.ResourceArray, error) {
	powerTotalSummaryResponse, err := dc.client.GetPowerSummaryByNodeId(dc.ctx, &api.PowerSummaryByNodeIdRequest{
		NodeId:               nodeId,
	})
	return types.NewResourceFromResponse(powerTotalSummaryResponse), err
}

func (dc *DataCenter) GetResourceList() (types.ResourceArray, error) {
	powerListRequest, err := dc.client.GetPowerList(dc.ctx, &api.PowerListRequest{})
	return types.NewResourceArrayFromPowerListResponse(powerListRequest), err
}

// InsertIdentity saves new identity info to the center of data.
func (dc *DataCenter) InsertIdentity(identity *types.Identity) error {
	response, err := dc.client.SaveIdentity(dc.ctx, types.NewSaveIdentityRequest(identity))
	if err != nil {
		log.WithError(err).WithField("hash", identity.Hash()).Errorf("InsertIdentity failed")
		return err
	}
	if response.Status != 0 {
		return fmt.Errorf("insert indentity error: %s", response.Msg)
	}
	return nil
}

// RevokeIdentity revokes the identity info to the center of data.
func (dc *DataCenter) RevokeIdentity(identity *types.Identity) error {
	response, err := dc.client.RevokeIdentityJoin(dc.ctx, &api.RevokeIdentityJoinRequest{
		Member:               &api.Organization{
			Name:                 identity.Name(),
			NodeId:               identity.NodeId(),
			IdentityId:           identity.IdentityId(),
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
	identityListResponse, err := dc.client.GetIdentityList(dc.ctx, &api.IdentityListRequest{})
	return types.NewIdentityArrayFromIdentityListResponse(identityListResponse), err
}

func (dc *DataCenter) GetIdentityByNodeId(nodeId string) (*types.Identity, error) {
	// todo: 读取本地节点，然后进行远程查询。
	return nil, nil
}

// InsertTask saves new task info to the center of data.
func (dc *DataCenter) InsertTask(task *types.Task) error {
	response, err := dc.client.SaveTask(dc.ctx, types.NewTaskDetail(task))
	if err != nil {
		log.WithError(err).WithField("hash", task.Hash()).Errorf("InsertTask failed")
		return err
	}
	if response.Status != 0 {
		return fmt.Errorf("insert task error: %s", response.Msg)
	}
	return nil
}

func (dc *DataCenter) GetTaskDataByTaskId(taskId string) (types.TaskDataArray, error) {
	taskListResponse, err := dc.client.ListTask(dc.ctx, &api.TaskListRequest{
		TaskId:               taskId,
	})
	return types.NewTaskArrayFromResponse(taskListResponse), err
}

func (dc *DataCenter) GetTaskDataListByNodeId(nodeId string) (types.TaskDataArray, error) {
	// todo: not to coding, temporary.
	return nil, nil
}

func (dc *DataCenter) GetTaskEventListByTaskId(taskId string) ([]*api.TaskEvent, error) {
	taskEventResponse, err := dc.client.ListTaskEvent(dc.ctx, &api.TaskEventRequest{
		TaskId:               taskId,
	})
	return taskEventResponse.TaskEventList, err
}

func (dc *DataCenter) SetSeedNode(seed *types.SeedNodeInfo) (types.NodeConnStatus, error) {
	rawdb.WriteSeedNodes(dc.db, seed)
	// todo: need to coding more logic...
	return types.NONCONNECTED, nil
}

func (dc *DataCenter) DeleteSeedNode(id string) error {
	rawdb.DeleteSeedNode(dc.db, id)
	return nil
}

func (dc *DataCenter) GetSeedNode(id string) (*types.SeedNodeInfo, error) {
	return rawdb.ReadSeedNode(dc.db, id), nil
}

func (dc *DataCenter) GetSeedNodeList() ([]*types.SeedNodeInfo, error) {
	return rawdb.ReadAllSeedNodes(dc.db), nil
}

func (dc *DataCenter) SetRegisterNode(typ types.RegisteredNodeType, node *types.RegisteredNodeInfo) (types.NodeConnStatus, error) {
	rawdb.WriteRegisterNodes(dc.db, typ, node)
	// todo: need to establish conn to registered node. heartbeat detection
	return types.NONCONNECTED, nil
}

func (dc *DataCenter) DeleteRegisterNode(typ types.RegisteredNodeType, id string) error {
	rawdb.DeleteRegisterNode(dc.db, typ, id)
	return nil
}

func (dc *DataCenter) GetRegisterNode(typ types.RegisteredNodeType, id string) (*types.RegisteredNodeInfo, error) {
	return rawdb.ReadRegisterNode(dc.db, typ, id), nil
}

func (dc *DataCenter) GetRegisterNodeList(typ types.RegisteredNodeType) ([]*types.RegisteredNodeInfo, error) {
	return rawdb.ReadAllRegisterNodes(dc.db, typ), nil
}

func (dc *DataCenter) StoreRunningTask(task *types.Task) error {
	rawdb.WriteRunningTask(dc.db, task)
	return nil
}

func (dc *DataCenter) StoreJobNodeRunningTaskId(jobNodeId, taskId string) error {
	rawdb.WriteRunningTaskIDList(dc.db, jobNodeId, taskId)
	return nil
}

func (dc *DataCenter) IncreaseRunningTaskCountOnOrg() uint32 {
	return rawdb.IncreaseRunningTaskCountForOrg(dc.db)
}

func (dc *DataCenter) IncreaseRunningTaskCountOnJobNode(jobNodeId string) uint32 {
	return rawdb.IncreaseRunningTaskCountForOrg(dc.db)
}

func (dc *DataCenter) GetRunningTaskCountOnOrg () uint32 {
	return rawdb.ReadRunningTaskCountForOrg(dc.db)
}

func (dc *DataCenter) GetRunningTaskCountOnJobNode (jobNodeId string) uint32 {
	return rawdb.ReadRunningTaskCountForJobNode(dc.db, jobNodeId)
}

func (dc *DataCenter) GetJobNodeRunningTaskIdList (jobNodeId string) []string {
	return rawdb.ReadRunningTaskIDList(dc.db, jobNodeId)
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
