package iface

import (
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type LocalStoreCarrierDB interface {
	// about datacenter config
	SetConfig (config *params.CarrierChainConfig) error
	// about carrier
	QueryYarnName() (string, error)
	SetSeedNode(seed *pb.SeedPeer) error
	RemoveSeedNode(addr string) error
	QuerySeedNodeList() ([]*pb.SeedPeer, error)
	SetRegisterNode(typ pb.RegisteredNodeType, node *pb.YarnRegisteredPeerDetail) error
	DeleteRegisterNode(typ pb.RegisteredNodeType, id string) error
	QueryRegisterNode(typ pb.RegisteredNodeType, id string) (*pb.YarnRegisteredPeerDetail, error)
	QueryRegisterNodeList(typ pb.RegisteredNodeType) ([]*pb.YarnRegisteredPeerDetail, error)
	// about local resource (local jobNode resource)
	InsertLocalResource(resource *types.LocalResource) error
	RemoveLocalResource(jobNodeId string) error
	QueryLocalResource(jobNodeId string) (*types.LocalResource, error)
	QueryLocalResourceList() (types.LocalResourceArray, error)
	// powerId -> jobNodeId
	StoreJobNodeIdIdByPowerId(powerId, jobNodeId string) error
	RemoveJobNodeIdByPowerId(powerId string) error
	QueryJobNodeIdByPowerId(powerId string) (string, error)

	// about jobRerource   (jobNodeId -> {jobNodeId, powerId, resource, slotTotal, slotUsed})
	StoreLocalResourceTable(resource *types.LocalResourceTable) error
	RemoveLocalResourceTable(resourceId string) error
	StoreLocalResourceTables(resources []*types.LocalResourceTable) error
	QueryLocalResourceTable(resourceId string) (*types.LocalResourceTable, error)
	QueryLocalResourceTables() ([]*types.LocalResourceTable, error)
	// about DataResourceTable (dataNodeId -> {dataNodeId, totalDisk, usedDisk})
	StoreDataResourceTable(StoreDataResourceTables *types.DataResourceTable) error
	StoreDataResourceTables(dataResourceTables []*types.DataResourceTable) error
	RemoveDataResourceTable(nodeId string) error
	QueryDataResourceTable(nodeId string) (*types.DataResourceTable, error)
	QueryDataResourceTables() ([]*types.DataResourceTable, error)
	// about DataResourceFileUpload (originId -> {originId, dataNodeId, metaDataId, filePath})
	StoreDataResourceFileUpload(dataResourceFileUpload *types.DataResourceFileUpload) error
	StoreDataResourceFileUploads(dataResourceFileUploads []*types.DataResourceFileUpload) error
	RemoveDataResourceFileUpload(originId string) error
	QueryDataResourceFileUpload(originId string) (*types.DataResourceFileUpload, error)
	QueryDataResourceFileUploads() ([]*types.DataResourceFileUpload, error)
	// about DataResourceDiskUsed (metaDataId -> {metaDataId, dataNodeId, diskUsed})
	StoreDataResourceDiskUsed(dataResourceDiskUsed *types.DataResourceDiskUsed) error
	RemoveDataResourceDiskUsed(metaDataId string) error
	QueryDataResourceDiskUsed(metaDataId string) (*types.DataResourceDiskUsed, error)
	// v2.0  about user metadataAuthUsed (userType + user -> metadataAuthId ...)
	//StoreUserMetadataAuthUsed(userType apicommonpb.UserType, user, metadataAuthId string) error
	//QueryUserMetadataAuthUsedCount(userType apicommonpb.UserType, user string) (uint32, error)
	//QueryUserMetadataAuthUseds(userType apicommonpb.UserType, user string) ([]string, error)
	//RemoveAllUserMetadataAuthUsed(userType apicommonpb.UserType, user string) error
	// v 2.0  about user metadataAuthUsed by metadataId (userType + user + metadataId -> metadataAuthId)
	StoreUserMetadataAuthIdByMetadataId(userType apicommonpb.UserType, user, metadataId, metadataAuthId string) error
	QueryUserMetadataAuthIdByMetadataId(userType apicommonpb.UserType, user, metadataId string) (string, error)
	HasUserMetadataAuthIdByMetadataId(userType apicommonpb.UserType, user, metadataId string) (bool, error)
	RemoveUserMetadataAuthIdByMetadataId(userType apicommonpb.UserType, user, metadataId string) error
	// v 2.0 about metadata used taskId    (metadataId -> [taskId, taskId, ..., taskId])
	StoreMetadataHistoryTaskId(metadataId, taskId string) error
	HasMetadataHistoryTaskId(metadataId, taskId string) (bool, error)
	QueryMetadataHistoryTaskIdCount(metadataId string) (uint32, error)
	QueryMetadataHistoryTaskIds(metadataId string) ([]string, error)
	// v 2.0 about Message Cache
	StoreMessageCache(value interface{}) error
	RemovePowerMsg(powerId string) error
	RemoveAllPowerMsg() error
	RemoveMetadataMsg(metadataId string) error
	RemoveAllMetadataMsg() error
	RemoveMetadataAuthMsg(metadataAuthId string) error
	RemoveAllMetadataAuthMsg() error
	RemoveTaskMsg(taskId string) error
	RemoveAllTaskMsg() error
	QueryPowerMsgArr() (types.PowerMsgArr, error)
	QueryMetadataMsgArr() (types.MetadataMsgArr, error)
	QueryMetadataAuthorityMsgArr() (types.MetadataAuthorityMsgArr, error)
	QueryTaskMsgArr() (types.TaskMsgArr, error)
}

type MetadataCarrierDB interface {
	StoreInternalMetadata(metadata *types.Metadata) error
	IsInternalMetadataByDataId(metadataId string) (bool, error)
	QueryInternalMetadataByDataId(metadataId string) (*types.Metadata, error)
	QueryInternalMetadataList() (types.MetadataArray, error)
	InsertMetadata(metadata *types.Metadata) error
	RevokeMetadata(metadata *types.Metadata) error
	QueryMetadataByDataId(dataId string) (*types.Metadata, error)
	QueryMetadataList(lastUpdate uint64, pageSize uint64) (types.MetadataArray, error)
}

type ResourceCarrierDB interface {
	InsertResource(resource *types.Resource) error
	RevokeResource(resource *types.Resource) error
	QueryGlobalResourceSummaryList(lastUpdate uint64, pageSize uint64) (types.ResourceArray, error)
	QueryGlobalResourceDetailList(lastUpdate uint64, pageSize uint64) (types.ResourceArray, error)
	SyncPowerUsed(resource *types.LocalResource) error
}

type IdentityCarrierDB interface {
	InsertIdentity(identity *types.Identity) error
	StoreIdentity(identity *apicommonpb.Organization) error
	RemoveIdentity() error
	QueryIdentityId() (string, error)
	QueryIdentity() (*apicommonpb.Organization, error)
	RevokeIdentity(identity *types.Identity) error
	QueryIdentityList(lastUpdate uint64, pageSize uint64) (types.IdentityArray, error)
	//QueryIdentityListByIds(identityIds []string) (types.IdentityArray, error)
	HasIdentity(identity *apicommonpb.Organization) (bool, error)
	// v2.0
	InsertMetadataAuthority(metadataAuth *types.MetadataAuthority) error
	UpdateMetadataAuthority(metadataAuth *types.MetadataAuthority) error
	//RevokeMetadataAuthority(metadataAuth *types.MetadataAuthority) error
	QueryMetadataAuthority(metadataAuthId string) (*types.MetadataAuthority, error)
	QueryMetadataAuthorityListByIds(metadataAuthIds []string) (types.MetadataAuthArray, error)
	QueryMetadataAuthorityListByIdentityId(identityId string, lastUpdate uint64, pageSize uint64) (types.MetadataAuthArray, error)
	QueryMetadataAuthorityList(lastUpdate uint64, pageSize uint64) (types.MetadataAuthArray, error)
}

type TaskCarrierDB interface {
	StoreLocalTask(task *types.Task) error
	RemoveLocalTask(taskId string) error
	HasLocalTask(taskId string) (bool, error)
	QueryLocalTask(taskId string) (*types.Task, error)
	QueryLocalTaskListByIds(taskIds []string) (types.TaskDataArray, error)
	QueryLocalTaskList() (types.TaskDataArray, error)
	QueryLocalTaskAndEvents(taskId string) (*types.Task, error)
	QueryLocalTaskAndEventsListByIds(taskIds []string) (types.TaskDataArray, error)
	QueryLocalTaskAndEventsList() (types.TaskDataArray, error)
	// about local task event
	StoreTaskEvent(event *libtypes.TaskEvent) error
	QueryTaskEventList(taskId string) ([]*libtypes.TaskEvent, error)
	QueryTaskEventListByPartyId (taskId, partyId string) ([]*libtypes.TaskEvent, error)
	RemoveTaskEventList(taskId string) error
	RemoveTaskEventListByPartyId(taskId, partyId string) error
	// about task on datacenter
	InsertTask(task *types.Task) error
	QueryTaskListByIdentityId(identityId string) (types.TaskDataArray, error)
	QueryTaskEventListByTaskId(taskId string) ([]*libtypes.TaskEvent, error)
	QueryTaskEventListByTaskIds(taskIds []string) ([]*libtypes.TaskEvent, error)
	// v 1.0 about TaskPowerUsed  (prefix + taskId + partyId -> {taskId, partId, jobNodeId, slotCount})
	StoreLocalTaskPowerUsed(taskPowerUsed *types.LocalTaskPowerUsed) error
	StoreLocalTaskPowerUseds(taskPowerUseds []*types.LocalTaskPowerUsed) error
	HasLocalTaskPowerUsed(taskId, partyId string) (bool, error)
	RemoveLocalTaskPowerUsed(taskId, partyId string) error
	RemoveLocalTaskPowerUsedByTaskId(taskId string) error
	QueryLocalTaskPowerUsed(taskId, partyId string) (*types.LocalTaskPowerUsed, error)
	QueryLocalTaskPowerUsedsByTaskId(taskId string) ([]*types.LocalTaskPowerUsed, error)
	QueryLocalTaskPowerUseds() ([]*types.LocalTaskPowerUsed, error)
	// about JobNodeTaskPartyId (prefix + jobNodeId + taskId -> [partyId, ..., partyId])
	StoreJobNodeTaskPartyId(jobNodeId, taskId, partyId string) error
	RemoveJobNodeTaskPartyId(jobNodeId, taskId, partyId string) error
	RemoveJobNodeTaskIdAllPartyIds(jobNodeId, taskId string) error
	QueryJobNodeRunningTaskIdList(jobNodeId string) ([]string, error)
	QueryJobNodeRunningTaskCount(jobNodeId string) (uint32, error)
	QueryJobNodeRunningTaskAllPartyIdList(jobNodeId, taskId string) ([]string, error)
	HasJobNodeRunningTaskId(jobNodeId, taskId string) (bool, error)
	HasJobNodeTaskPartyId(jobNodeId, taskId, partyId string) (bool, error)
	QueryJobNodeTaskPartyIdCount(jobNodeId, taskId string) (uint32, error)
	// v 2.0 about jobNode history task count (prefix + jobNodeId -> history task count AND prefix + jobNodeId + taskId -> index)
	StoreJobNodeHistoryTaskId (jobNodeId, taskId string) error
	HasJobNodeHistoryTaskId (jobNodeId, taskId string) (bool, error)
	QueryJobNodeHistoryTaskCount (jobNodeId string) (uint32, error)
	// v 2.0  about TaskResultFileMetadataId  (taskId -> {taskId, originId, metadataId})
	StoreTaskUpResultFile(turf *types.TaskUpResultFile) error
	QueryTaskUpResultFile(taskId string) (*types.TaskUpResultFile, error)
	QueryTaskUpResultFileList() ([]*types.TaskUpResultFile, error)
	RemoveTaskUpResultFile(taskId string) error
	// v 2.0 about task powerPartyIds (prefix + taskId -> powerPartyIds)
	StoreTaskPowerPartyIds(taskId string, powerPartyIds []string) error
	QueryTaskPowerPartyIds(taskId string) ([]string, error)
	RemoveTaskPowerPartyIds (taskId string) error
	// v 2.0 about task partyIds of all partners (prefix + taskId -> [partyId, ..., partyId]  for task sender)
	StoreTaskPartnerPartyIds(taskId string, partyIds []string) error
	HasTaskPartnerPartyIds(taskId string) (bool, error)
	QueryTaskPartnerPartyIds(taskId string) ([]string, error)
	RemoveTaskPartnerPartyId (taskId, partyId string) error
	RemoveTaskPartnerPartyIds (taskId string) error
	// v 1.0 -> v 2.0 about task exec status (prefix + taskId + partyId -> "cons"|"exec")
	StoreLocalTaskExecuteStatusValConsByPartyId(taskId, partyId string) error
	StoreLocalTaskExecuteStatusValExecByPartyId(taskId, partyId string) error
	RemoveLocalTaskExecuteStatusByPartyId(taskId, partyId string) error
	HasLocalTaskExecuteStatusParty(taskId string) (bool, error)
	HasLocalTaskExecuteStatusByPartyId(taskId, partyId string) (bool, error)
	HasLocalTaskExecuteStatusValConsByPartyId(taskId, partyId string) (bool, error)
	HasLocalTaskExecuteStatusValExecByPartyId(taskId, partyId string) (bool, error)
	// v 2.0 about NeedExecuteTask
	StoreNeedExecuteTask(task *types.NeedExecuteTask) error
	RemoveNeedExecuteTaskByPartyId(taskId, partyId string) error
	RemoveNeedExecuteTask(taskId string) error
	ForEachNeedExecuteTaskWwithPrefix (prifix []byte, f func(key, value []byte) error) error
	ForEachNeedExecuteTask (f func(key, value []byte) error) error
	// v 2.0 about taskbullet
	StoreTaskBullet(bullet *types.TaskBullet) error
	RemoveTaskBullet(taskId string) error
	ForEachTaskBullets(f func(key, value []byte) error) error

}

type ForConsensusDB interface {
	IdentityCarrierDB
	TaskCarrierDB
}

type ForHandleDB interface {
	LocalStoreCarrierDB
	IdentityCarrierDB
	ResourceCarrierDB
	MetadataCarrierDB
	TaskCarrierDB
}

type ForResourceDB interface {
	LocalStoreCarrierDB
	IdentityCarrierDB
	ResourceCarrierDB
	TaskCarrierDB
}

type ForScheduleDB interface {
	IdentityCarrierDB
	LocalStoreCarrierDB
}
