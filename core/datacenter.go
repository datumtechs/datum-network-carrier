package core

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/db"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
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
	name, err := rawdb.ReadYarnName(dc.db)
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
func (dc *DataCenter) SetSeedNode(seed *pb.SeedPeer) (pb.ConnState, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.WriteSeedNodes(dc.db, seed)
	return pb.ConnState_ConnState_UnConnected, nil
}

func (dc *DataCenter) DeleteSeedNode(id string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.DeleteSeedNode(dc.db, id)
	return nil
}

func (dc *DataCenter) GetSeedNode(id string) (*pb.SeedPeer, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.ReadSeedNode(dc.db, id)
}

func (dc *DataCenter) GetSeedNodeList() ([]*pb.SeedPeer, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.ReadAllSeedNodes(dc.db)
}

func (dc *DataCenter) SetRegisterNode(typ pb.RegisteredNodeType, node *pb.YarnRegisteredPeerDetail) (pb.ConnState, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.WriteRegisterNodes(dc.db, typ, node)
	return node.ConnState, nil
}

func (dc *DataCenter) DeleteRegisterNode(typ pb.RegisteredNodeType, id string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.DeleteRegisterNode(dc.db, typ, id)
	return nil
}

func (dc *DataCenter) GetRegisterNode(typ pb.RegisteredNodeType, id string) (*pb.YarnRegisteredPeerDetail, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.ReadRegisterNode(dc.db, typ, id)
}

func (dc *DataCenter) GetRegisterNodeList(typ pb.RegisteredNodeType) ([]*pb.YarnRegisteredPeerDetail, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.ReadAllRegisterNodes(dc.db, typ)
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

func (dc *DataCenter) HasLocalTaskPowerUsed(taskId, partyId string) (bool, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.HasLocalTaskPowerUsed(dc.db, taskId, partyId)
}


func (dc *DataCenter) RemoveLocalTaskPowerUsed(taskId, partyId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveLocalTaskPowerUsed(dc.db, taskId, partyId)
}

func (dc *DataCenter) RemoveLocalTaskPowerUsedByTaskId(taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveLocalTaskPowerUsedByTaskId(dc.db, taskId)
}

func (dc *DataCenter) QueryLocalTaskPowerUsed(taskId, partyId string) (*types.LocalTaskPowerUsed, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryLocalTaskPowerUsed(dc.db, taskId, partyId)
}

func (dc *DataCenter) QueryLocalTaskPowerUsedsByTaskId(taskId string) ([]*types.LocalTaskPowerUsed, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	res, _ := rawdb.QueryLocalTaskPowerUsedsByTaskId(dc.db, taskId)
	return res, nil
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
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreDataResourceDiskUsed(dc.db, dataResourceDiskUsed)
}

func (dc *DataCenter) RemoveDataResourceDiskUsed(metaDataId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveDataResourceDiskUsed(dc.db, metaDataId)
}

func (dc *DataCenter) QueryDataResourceDiskUsed(metaDataId string) (*types.DataResourceDiskUsed, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryDataResourceDiskUsed(dc.db, metaDataId)
}

// about LocalTaskExecuteStatus
func (dc *DataCenter) StoreLocalTaskExecuteStatus(taskId, partyId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreLocalTaskExecuteStatus(dc.db, taskId, partyId)
}

func (dc *DataCenter) RemoveLocalTaskExecuteStatus(taskId, partyId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveLocalTaskExecuteStatus(dc.db, taskId, partyId)
}

func (dc *DataCenter) HasLocalTaskExecute(taskId, partyId string) (bool, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.HasLocalTaskExecute(dc.db, taskId, partyId)
}

// about UserMetadataAuthUsed
func (dc *DataCenter) StoreUserMetadataAuthUsed (userType apicommonpb.UserType, user, metadataAuthId string)  error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreUserMetadataAauthUsed(dc.db, userType, user, metadataAuthId)
}

func (dc *DataCenter) QueryUserMetadataAuthUsedCount (userType apicommonpb.UserType, user string) (uint32, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryUserMetadataAuthUsedCount(dc.db, userType, user)
}

func (dc *DataCenter) QueryUserMetadataAuthUseds (userType apicommonpb.UserType, user string) ([]string, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	count, err := rawdb.QueryUserMetadataAuthUsedCount(dc.db, userType, user)
	if nil != err {
		return nil, err
	}

	if 0 == count {
		return nil, rawdb.ErrNotFound
	}

	metadataAuthIds := make([]string, 0)

	for index := 1; index <= int(count); index++ {

		metadataAuthId, err := rawdb.QueryUserMetadataAuthUsedByIndex(dc.db, userType, user, uint32(index))
		switch {
		case rawdb.IsNoDBNotFoundErr(err):
			return nil, err
		case rawdb.IsDBNotFoundErr(err):
			continue
		}
		metadataAuthIds = append(metadataAuthIds, metadataAuthId)
	}
	if len(metadataAuthIds) == 0 {
		return nil, rawdb.ErrNotFound
	}
	return metadataAuthIds, nil
}

func (dc *DataCenter) RemoveAllUserMetadataAuthUsed (userType apicommonpb.UserType, user string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	count, err := rawdb.QueryUserMetadataAuthUsedCount(dc.db, userType, user)
	switch {
	case rawdb.IsNoDBNotFoundErr(err):
		return err
	case rawdb.IsDBNotFoundErr(err) || 0 == count:
		return nil
	}

	for index := 1; index <= int(count); index++ {
		err := rawdb.RemoveUserMetadataAuthUsedByIndex(dc.db, userType, user, uint32(index))
		switch {
		case rawdb.IsNoDBNotFoundErr(err):
			return err
		case rawdb.IsDBNotFoundErr(err):
			continue
		}
	}
	return rawdb.RemoveUserMetadataAuthUsedCount(dc.db, userType, user)
}

func (dc *DataCenter) StoreUserMetadataAuthIdByMetadataId (userType apicommonpb.UserType, user, metadataId, metadataAuthId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreUserMetadataAuthIdByMetadataId(dc.db, userType, user, metadataId, metadataAuthId)
}

func (dc *DataCenter) QueryUserMetadataAuthIdByMetadataId (userType apicommonpb.UserType, user, metadataId string) (string, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryUserMetadataAuthIdByMetadataId(dc.db, userType, user, metadataId)
}

func (dc *DataCenter) HasUserMetadataAuthIdByMetadataId (userType apicommonpb.UserType, user, metadataId string) (bool, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.HasUserMetadataAuthIdByMetadataId(dc.db, userType, user, metadataId)
}

func (dc *DataCenter) RemoveUserMetadataAuthIdByMetadataId (userType apicommonpb.UserType, user, metadataId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveUserMetadataAuthIdByMetadataId(dc.db, userType, user, metadataId)
}


// about metadata used taskId
func (dc *DataCenter) StoreMetadataUsedTaskId (metadataId, taskId string)  error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreMetadataUsedTaskId(dc.db, metadataId, taskId)
}

func (dc *DataCenter) QueryMetadataUsedTaskIdCount (metadataId string) (uint32, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryMetadataUsedTaskIdCount(dc.db, metadataId)
}

func (dc *DataCenter) QueryMetadataUsedTaskIds (metadataId string) ([]string, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	count, err := rawdb.QueryMetadataUsedTaskIdCount(dc.db, metadataId)
	if nil != err {
		return nil, err
	}

	if 0 == count {
		return nil, rawdb.ErrNotFound
	}

	taskIds := make([]string, 0)

	for index := 1; index <= int(count); index++ {

		taskId, err := rawdb.QueryMetadataUsedTaskIdByIndex(dc.db, metadataId, uint32(index))
		switch {
		case rawdb.IsNoDBNotFoundErr(err):
			return nil, err
		case rawdb.IsDBNotFoundErr(err):
			continue
		}
		taskIds = append(taskIds, taskId)
	}
	if len(taskIds) == 0 {
		return nil, rawdb.ErrNotFound
	}
	return taskIds, nil
}

func (dc *DataCenter) RemoveAllMetadataUsedTaskId (metadataId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	count, err := rawdb.QueryMetadataUsedTaskIdCount(dc.db, metadataId)
	switch {
	case rawdb.IsNoDBNotFoundErr(err):
		return err
	case rawdb.IsDBNotFoundErr(err) || 0 == count:
		return nil
	}

	for index := 1; index <= int(count); index++ {
		err := rawdb.RemoveMetadataUsedTaskIdByIndex(dc.db, metadataId, uint32(index))
		switch {
		case rawdb.IsNoDBNotFoundErr(err):
			return err
		case rawdb.IsDBNotFoundErr(err):
			continue
		}
	}
	return rawdb.RemoveMetadataUsedTaskIdCount(dc.db, metadataId)
}

// about TaskResultFileMetadataId
func (dc *DataCenter) StoreTaskUpResultFile(turf *types.TaskUpResultFile)  error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreTaskUpResultFile(dc.db, turf)
}

func (dc *DataCenter) QueryTaskUpResultFile(taskId string)  (*types.TaskUpResultFile, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryTaskUpResultFile(dc.db, taskId)
}

func (dc *DataCenter) QueryTaskUpResultFileList () ([]*types.TaskUpResultFile, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryTaskUpResultFileList(dc.db)
}

func (dc *DataCenter) RemoveTaskUpResultFile(taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveTaskUpResultFile(dc.db, taskId)
}

func (dc *DataCenter) StoreTaskResuorceUsage(usage *types.TaskResuorceUsage) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreTaskResuorceUsage(dc.db, usage)
}

func (dc *DataCenter) QueryTaskResuorceUsage(taskId, partyId string) (*types.TaskResuorceUsage, error)  {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryTaskResuorceUsage(dc.db, taskId, partyId)
}

func (dc *DataCenter) RemoveTaskResuorceUsage(taskId, partyId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveTaskResuorceUsage(dc.db, taskId, partyId)
}

func (dc *DataCenter) RemoveTaskResuorceUsageByTaskId (taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveTaskResuorceUsageByTaskId(dc.db, taskId)
}

func (dc *DataCenter) StoreTaskEvent(event *libtypes.TaskEvent) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.WriteTaskEvent(dc.db, event)
	log.Debugf("Store task eventList, event: %s", event.String())
	return nil
}

func (dc *DataCenter) GetTaskEventList(taskId string) ([]*libtypes.TaskEvent, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	list, err := rawdb.ReadTaskEvent(dc.db, taskId)
	if nil != err {
		return nil, err
	}
	log.Debugf("Query local task eventList, taskId: {%s}, local eventList Len: {%d}", taskId, len(list))
	return list, nil
}

func (dc *DataCenter) RemoveTaskEventList(taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.DeleteTaskEvent(dc.db, taskId)
	log.Debugf("Remove task eventList, taskId: {%s}", taskId)
	return nil
}

// about Message Cache
func (dc *DataCenter) StoreMessageCache(value interface{}) {
	rawdb.StoreMessageCache(dc.db, value)
}

func (dc *DataCenter) QueryPowerMsgArr() (types.PowerMsgArr,error) {
	return rawdb.QueryPowerMsgArr(dc.db)
}

func (dc *DataCenter) QueryMetadataMsgArr() (types.MetadataMsgArr,error) {
	return rawdb.QueryMetadataMsgArr(dc.db)
}

func (dc *DataCenter) QueryMetadataAuthorityMsgArr() (types.MetadataAuthorityMsgArr,error) {
	return rawdb.QueryMetadataAuthorityMsgArr(dc.db)
}

func (dc *DataCenter) QueryTaskMsgArr()(types.TaskMsgArr,error)  {
	return rawdb.QueryTaskMsgArr(dc.db)
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
