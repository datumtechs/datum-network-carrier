package iface

import (
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type LocalStoreCarrierDB interface {
	QueryYarnName() (string, error)
	SetSeedNode(seed *pb.SeedPeer) (pb.ConnState, error)
	DeleteSeedNode(id string) error
	QuerySeedNode(id string) (*pb.SeedPeer, error)
	QuerySeedNodeList() ([]*pb.SeedPeer, error)
	SetRegisterNode(typ pb.RegisteredNodeType, node *pb.YarnRegisteredPeerDetail) (pb.ConnState, error)
	DeleteRegisterNode(typ pb.RegisteredNodeType, id string) error
	QueryRegisterNode(typ pb.RegisteredNodeType, id string) (*pb.YarnRegisteredPeerDetail, error)
	QueryRegisterNodeList(typ pb.RegisteredNodeType) ([]*pb.YarnRegisteredPeerDetail, error)

	InsertLocalResource(resource *types.LocalResource) error
	RemoveLocalResource(jobNodeId string) error
	QueryLocalResource(jobNodeId string) (*types.LocalResource, error)
	QueryLocalResourceList() (types.LocalResourceArray, error)
	// powerId -> jobNodeId
	StoreLocalResourceIdByPowerId(powerId, jobNodeId string) error
	RemoveLocalResourceIdByPowerId(powerId string) error
	QueryLocalResourceIdByPowerId(powerId string) (string, error)

	// about jobRerource   (jobNodeId -> {jobNodeId, powerId, resource, slotTotal, slotUsed})
	StoreLocalResourceTable(resource *types.LocalResourceTable) error
	RemoveLocalResourceTable(resourceId string) error
	StoreLocalResourceTables(resources []*types.LocalResourceTable) error
	QueryLocalResourceTable(resourceId string) (*types.LocalResourceTable, error)
	QueryLocalResourceTables() ([]*types.LocalResourceTable, error)
	// about slotUnit (key -> slotUnit)
	StoreNodeResourceSlotUnit(slot *types.Slot) error
	RemoveNodeResourceSlotUnit() error
	QueryNodeResourceSlotUnit() (*types.Slot, error)
	// about TaskPowerUsed  (prefix + taskId + partyId -> {taskId, partId, jobNodeId, slotCount})
	StoreLocalTaskPowerUsed(taskPowerUsed *types.LocalTaskPowerUsed) error
	StoreLocalTaskPowerUseds(taskPowerUseds []*types.LocalTaskPowerUsed) error
	HasLocalTaskPowerUsed(taskId, partyId string) (bool, error)
	RemoveLocalTaskPowerUsed(taskId, partyId string) error
	RemoveLocalTaskPowerUsedByTaskId(taskId string) error
	QueryLocalTaskPowerUsed(taskId, partyId string) (*types.LocalTaskPowerUsed, error)
	QueryLocalTaskPowerUsedsByTaskId(taskId string) ([]*types.LocalTaskPowerUsed, error)
	QueryLocalTaskPowerUseds() ([]*types.LocalTaskPowerUsed, error)
	// about resourceTaskIds Mapping (jobNodeId -> [taskId, taskId, ..., taskId])
	StoreJobNodeRunningTaskId(jobNodeId, taskId string) error
	RemoveJobNodeRunningTaskId(jobNodeId, taskId string) error
	QueryRunningTaskCountOnJobNode(jobNodeId string) (uint32, error)
	QueryJobNodeRunningTaskIdList(jobNodeId string) ([]string, error)
	// about resource task party count (prefix + taskId -> n (n: partyId count))
	IncreaseResourceTaskPartyIdCount(jobNodeId, taskId string) error
	DecreaseResourceTaskPartyIdCount(jobNodeId, taskId string) error
	QueryResourceTaskPartyIdCount(jobNodeId, taskId string) (uint32, error)
	// about task totalCount on jobNode ever (prefix + jobNodeId -> taskTotalCount)
	IncreaseResourceTaskTotalCount (jobNodeId string) error
	DecreaseResourceTaskTotalCount (jobNodeId string) error
	RemoveResourceTaskTotalCount (jobNodeId string) error
	QueryResourceTaskTotalCount (jobNodeId string) (uint32, error)

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
	// about task exec status (prefix + taskId + partyId -> "yes")
	StoreLocalTaskExecuteStatus(taskId, partyId string) error
	RemoveLocalTaskExecuteStatus(taskId, partyId string) error
	HasLocalTaskExecute(taskId, partyId string) (bool, error)
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
	StoreMetadataUsedTaskId(metadataId, taskId string) error
	QueryMetadataUsedTaskIdCount(metadataId string) (uint32, error)
	QueryMetadataUsedTaskIds(metadataId string) ([]string, error)
	RemoveAllMetadataUsedTaskId(metadataId string) error
	// v 2.0  about TaskResultFileMetadataId  (taskId -> {taskId, originId, metadataId})
	StoreTaskUpResultFile(turf *types.TaskUpResultFile) error
	QueryTaskUpResultFile(taskId string) (*types.TaskUpResultFile, error)
	QueryTaskUpResultFileList() ([]*types.TaskUpResultFile, error)
	RemoveTaskUpResultFile(taskId string) error
	// V 2.0 about task used resource  (taskId -> resourceUsed)
	StoreTaskResuorceUsage(usage *types.TaskResuorceUsage) error
	QueryTaskResuorceUsage(taskId, partyId string) (*types.TaskResuorceUsage, error)
	RemoveTaskResuorceUsage(taskId, partyId string) error
	RemoveTaskResuorceUsageByTaskId(taskId string) error
	// v 2.0 about task powerPartyIds (prefix + taskId -> powerPartyIds)
	StoreTaskPowerPartyIds(taskId string, powerPartyIds []string) error
	QueryTaskPowerPartyIds(taskId string) ([]string, error)
	RemoveTaskPowerPartyIds (taskId string) error
	// v 2.0 about Message Cache
	StoreMessageCache(value interface{})
	QueryPowerMsgArr() (types.PowerMsgArr, error)
	QueryMetadataMsgArr() (types.MetadataMsgArr, error)
	QueryMetadataAuthorityMsgArr() (types.MetadataAuthorityMsgArr, error)
	QueryTaskMsgArr() (types.TaskMsgArr, error)
}

type MetadataCarrierDB interface {
	StoreLocalMetadata(metadata *types.Metadata) error
	QueryLocalMetadataByDataId(metadataId string) (*types.Metadata, error)
	QueryLocalMetadataList() (types.MetadataArray, error)
	InsertMetadata(metadata *types.Metadata) error
	RevokeMetadata(metadata *types.Metadata) error
	QueryMetadataByDataId(dataId string) (*types.Metadata, error)
	QueryMetadataList() (types.MetadataArray, error)
}

type ResourceCarrierDB interface {
	InsertResource(resource *types.Resource) error
	RevokeResource(resource *types.Resource) error
	QueryResourceList() (types.ResourceArray, error)
	SyncPowerUsed(resource *types.LocalResource) error
}

type IdentityCarrierDB interface {
	InsertIdentity(identity *types.Identity) error
	StoreIdentity(identity *apicommonpb.Organization) error
	RemoveIdentity() error
	QueryIdentityId() (string, error)
	QueryIdentity() (*apicommonpb.Organization, error)
	RevokeIdentity(identity *types.Identity) error
	QueryIdentityList() (types.IdentityArray, error)
	//QueryIdentityListByIds(identityIds []string) (types.IdentityArray, error)
	HasIdentity(identity *apicommonpb.Organization) (bool, error)

	// v2.0
	InsertMetadataAuthority(metadataAuth *types.MetadataAuthority) error
	UpdateMetadataAuthority(metadataAuth *types.MetadataAuthority) error
	//RevokeMetadataAuthority(metadataAuth *types.MetadataAuthority) error
	QueryMetadataAuthority(metadataAuthId string) (*types.MetadataAuthority, error)
	QueryMetadataAuthorityListByIds(metadataAuthIds []string) (types.MetadataAuthArray, error)
	QueryMetadataAuthorityListByIdentityId(identityId string, lastUpdate uint64) (types.MetadataAuthArray, error)
	QueryMetadataAuthorityList(lastUpdate uint64) (types.MetadataAuthArray, error)
}

type TaskCarrierDB interface {
	StoreTaskEvent(event *libtypes.TaskEvent) error
	QueryTaskEventList(taskId string) ([]*libtypes.TaskEvent, error)
	RemoveTaskEventList(taskId string) error
	StoreLocalTask(task *types.Task) error
	RemoveLocalTask(taskId string) error
	QueryLocalTask(taskId string) (*types.Task, error)
	QueryLocalTaskListByIds(taskIds []string) (types.TaskDataArray, error)
	QueryLocalTaskList() (types.TaskDataArray, error)
	QueryLocalTaskAndEvents(taskId string) (*types.Task, error)
	QueryLocalTaskAndEventsListByIds(taskIds []string) (types.TaskDataArray, error)
	QueryLocalTaskAndEventsList() (types.TaskDataArray, error)

	// about task on datacenter
	InsertTask(task *types.Task) error
	QueryTaskListByIdentityId(identityId string) (types.TaskDataArray, error)
	QueryRunningTaskCountOnOrg() uint32
	QueryTaskEventListByTaskId(taskId string) ([]*libtypes.TaskEvent, error)
	QueryTaskEventListByTaskIds(taskIds []string) ([]*libtypes.TaskEvent, error)

	// about scheduling
	StoreScheduling(bullet *types.TaskBullet)
	DeleteScheduling(bullet *types.TaskBullet)
	RecoveryScheduling() (*types.TaskBullets, *types.TaskBullets,map[string]*types.TaskBullet)
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
