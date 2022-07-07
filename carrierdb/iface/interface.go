package iface

import (
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/types"
)

type LocalStoreCarrierDB interface {
	// about datacenter config
	SetConfig(config *types.CarrierChainConfig) error
	// about carrier
	QueryYarnName() (string, error)
	SetSeedNode(seed *carrierapipb.SeedPeer) error
	RemoveSeedNode(addr string) error
	QuerySeedNodeList() ([]*carrierapipb.SeedPeer, error)
	SetRegisterNode(typ carrierapipb.RegisteredNodeType, node *carrierapipb.YarnRegisteredPeerDetail) error
	DeleteRegisterNode(typ carrierapipb.RegisteredNodeType, id string) error
	QueryRegisterNode(typ carrierapipb.RegisteredNodeType, id string) (*carrierapipb.YarnRegisteredPeerDetail, error)
	QueryRegisterNodeList(typ carrierapipb.RegisteredNodeType) ([]*carrierapipb.YarnRegisteredPeerDetail, error)
	// about powerId -> jobNodeId
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
	// about DataResourceDataUpload (originId -> {originId, dataNodeId, metaDataId, filePath})
	StoreDataResourceDataUpload(dataResourceDataUpload *types.DataResourceDataUpload) error
	StoreDataResourceDataUploads(dataResourceDataUploads []*types.DataResourceDataUpload) error
	RemoveDataResourceDataUpload(originId string) error
	QueryDataResourceDataUpload(originId string) (*types.DataResourceDataUpload, error)
	QueryDataResourceDataUploads() ([]*types.DataResourceDataUpload, error)
	// about DataResourceDiskUsed (metaDataId -> {metaDataId, dataNodeId, diskUsed})
	StoreDataResourceDiskUsed(dataResourceDiskUsed *types.DataResourceDiskUsed) error
	RemoveDataResourceDiskUsed(metaDataId string) error
	QueryDataResourceDiskUsed(metaDataId string) (*types.DataResourceDiskUsed, error)
	// v 0.2.0  about user metadataAuthUsed by metadataId (userType + user + metadataId -> metadataAuthId)
	StoreUserMetadataAuthIdByMetadataId(userType commonconstantpb.UserType, user, metadataId, metadataAuthId string) error
	QueryUserMetadataAuthIdByMetadataId(userType commonconstantpb.UserType, user, metadataId string) (string, error)
	HasUserMetadataAuthIdByMetadataId(userType commonconstantpb.UserType, user, metadataId string) (bool, error)
	RemoveUserMetadataAuthIdByMetadataId(userType commonconstantpb.UserType, user, metadataId string) error
	// v 0.2.0 about metadata used taskId    (metadataId -> [taskId, taskId, ..., taskId])
	StoreMetadataHistoryTaskId(metadataId, taskId string) error
	HasMetadataHistoryTaskId(metadataId, taskId string) (bool, error)
	QueryMetadataHistoryTaskIdCount(metadataId string) (uint32, error)
	QueryMetadataHistoryTaskIds(metadataId string) ([]string, error)
	// v 0.2.0 about Message Cache
	StoreMessageCache(value interface{}) error
	RemovePowerMsg(powerId string) error
	RemoveAllPowerMsg() error
	RemoveMetadataMsg(metadataId string) error
	RemoveMetadataUpdateMsg(metadataId string) error
	RemoveAllMetadataMsg() error
	RemoveMetadataAuthMsg(metadataAuthId string) error
	RemoveAllMetadataAuthMsg() error
	RemoveTaskMsg(taskId string) error
	RemoveAllTaskMsg() error
	QueryPowerMsgArr() (types.PowerMsgArr, error)
	QueryMetadataMsgArr() (types.MetadataMsgArr, error)
	QueryMetadataUpdateMsgArr() (types.MetadataUpdateMsgArr, error)
	QueryMetadataAuthorityMsgArr() (types.MetadataAuthorityMsgArr, error)
	QueryTaskMsgArr() (types.TaskMsgArr, error)

	StoreOrgWallet(orgWallet *types.OrgWallet) error
	// QueryOrgWallet does not return ErrNotFound if the organization wallet not found.
	QueryOrgWallet() (*types.OrgWallet, error)
}

type MetadataCarrierDB interface {
	// about global metadata
	InsertMetadata(metadata *types.Metadata) error
	RevokeMetadata(metadata *types.Metadata) error
	QueryMetadataById(metadataId string) (*types.Metadata, error)
	QueryMetadataByIds(metadataIds []string) ([]*types.Metadata, error) // add by v 0.4.0
	QueryMetadataList(lastUpdate, pageSize uint64) (types.MetadataArray, error)
	QueryMetadataListByIdentity(identityId string, lastUpdate, pageSize uint64) (types.MetadataArray, error)
	UpdateGlobalMetadata(metadata *types.Metadata) error // add by v 0.4.0
	// v 0.3.0 about internal metadata (generate from a task result file)
	StoreInternalMetadata(metadata *types.Metadata) error
	IsInternalMetadataById(metadataId string) (bool, error)
	QueryInternalMetadataById(metadataId string) (*types.Metadata, error)
	QueryInternalMetadataList() (types.MetadataArray, error)
}

type ResourceCarrierDB interface {
	// about global power resource
	InsertResource(resource *types.Resource) error
	RevokeResource(resource *types.Resource) error
	QueryGlobalResourceSummaryList() (types.ResourceArray, error)
	QueryGlobalResourceDetailList(lastUpdate, pageSize uint64) (types.ResourceArray, error)
	SyncPowerUsed(resource *types.LocalResource) error
	// about local resource (local jobNode resource)
	StoreLocalResource(resource *types.LocalResource) error
	RemoveLocalResource(jobNodeId string) error
	QueryLocalResource(jobNodeId string) (*types.LocalResource, error)
	QueryLocalResourceList() (types.LocalResourceArray, error)
}

type IdentityCarrierDB interface {
	// about global identity
	InsertIdentity(identity *types.Identity) error
	RevokeIdentity(identity *types.Identity) error
	QueryIdentityList(lastUpdate, pageSize uint64) (types.IdentityArray, error)
	//QueryIdentityListByIds(identityIds []string) (types.IdentityArray, error)
	// about local identity
	HasIdentity(identity *carriertypespb.Organization) (bool, error)
	StoreIdentity(identity *carriertypespb.Organization) error
	RemoveIdentity() error
	QueryIdentityId() (string, error)
	QueryIdentity() (*carriertypespb.Organization, error)
}

type MetadataAuthorityCarrierDB interface {
	// v2.0
	InsertMetadataAuthority(metadataAuth *types.MetadataAuthority) error
	UpdateMetadataAuthority(metadataAuth *types.MetadataAuthority) error
	//RevokeMetadataAuthority(metadataAuth *types.MetadataAuthority) error
	QueryMetadataAuthority(metadataAuthId string) (*types.MetadataAuthority, error)
	QueryMetadataAuthorityListByIds(metadataAuthIds []string) (types.MetadataAuthArray, error)
	QueryMetadataAuthorityListByIdentityId(identityId string, lastUpdate, pageSize uint64) (types.MetadataAuthArray, error)
	QueryMetadataAuthorityList(lastUpdate, pageSize uint64) (types.MetadataAuthArray, error)
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
	StoreTaskEvent(event *carriertypespb.TaskEvent) error
	QueryTaskEventList(taskId string) ([]*carriertypespb.TaskEvent, error)
	QueryTaskEventListByPartyId(taskId, partyId string) ([]*carriertypespb.TaskEvent, error)
	RemoveTaskEventList(taskId string) error
	RemoveTaskEventListByPartyId(taskId, partyId string) error
	// about global task on datacenter
	InsertTask(task *types.Task) error
	QueryGlobalTaskList(lastUpdate, pageSize uint64) (types.TaskDataArray, error)
	QueryTaskListByIdentityId(identityId string, lastUpdate, pageSize uint64) (types.TaskDataArray, error)
	QueryTaskListByTaskIds(taskIds []string) (types.TaskDataArray, error)
	QueryTaskEventListByTaskId(taskId string) ([]*carriertypespb.TaskEvent, error)
	QueryTaskEventListByTaskIds(taskIds []string) ([]*carriertypespb.TaskEvent, error)
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
	QueryJobNodeRunningTaskIdsAndPartyIdsPairs(jobNodeId string) (map[string][]string, error)
	QueryJobNodeRunningTaskAllPartyIdList(jobNodeId, taskId string) ([]string, error)
	HasJobNodeRunningTaskId(jobNodeId, taskId string) (bool, error)
	HasJobNodeTaskPartyId(jobNodeId, taskId, partyId string) (bool, error)
	QueryJobNodeTaskPartyIdCount(jobNodeId, taskId string) (uint32, error)
	// v 2.0 about jobNode history task count (prefix + jobNodeId -> history task count AND prefix + jobNodeId + taskId -> index)
	StoreJobNodeHistoryTaskId(jobNodeId, taskId string) error
	HasJobNodeHistoryTaskId(jobNodeId, taskId string) (bool, error)
	QueryJobNodeHistoryTaskCount(jobNodeId string) (uint32, error)
	// v 2.0  about TaskResultData  (taskId -> {taskId, originId, metadataId})
	StoreTaskUpResultData(turf *types.TaskUpResultData) error
	QueryTaskUpResulData(taskId string) (*types.TaskUpResultData, error)
	QueryTaskUpResultDataList() ([]*types.TaskUpResultData, error)
	RemoveTaskUpResultData(taskId string) error
	// v 2.0 about task partyIds of all partners (prefix + taskId -> [partyId, ..., partyId]  for task sender)
	StoreTaskPartnerPartyIds(taskId string, partyIds []string) error
	HasTaskPartnerPartyIds(taskId string) (bool, error)
	QueryTaskPartnerPartyIds(taskId string) ([]string, error)
	RemoveTaskPartnerPartyId(taskId, partyId string) error
	RemoveTaskPartnerPartyIds(taskId string) error
	// v 1.0 -> v 2.0 about task exec status (prefix + taskId + partyId -> "consensusing"|"running|terminating")
	StoreLocalTaskExecuteStatusValConsensusByPartyId(taskId, partyId string) error
	StoreLocalTaskExecuteStatusValRunningByPartyId(taskId, partyId string) error
	StoreLocalTaskExecuteStatusValTerminateByPartyId(taskId, partyId string) error // add by v 0.3.0
	RemoveLocalTaskExecuteStatusByPartyId(taskId, partyId string) error
	HasLocalTaskExecuteStatusParty(taskId string) (bool, error)
	HasLocalTaskExecuteStatusByPartyId(taskId, partyId string) (bool, error)
	HasLocalTaskExecuteStatusConsensusByPartyId(taskId, partyId string) (bool, error)
	HasLocalTaskExecuteStatusRunningByPartyId(taskId, partyId string) (bool, error)
	HasLocalTaskExecuteStatusTerminateByPartyId(taskId, partyId string) (bool, error) // add by v 0.3.0
	// v 2.0 about NeedExecuteTask
	StoreNeedExecuteTask(task *types.NeedExecuteTask) error
	RemoveNeedExecuteTaskByPartyId(taskId, partyId string) error
	RemoveNeedExecuteTask(taskId string) error
	ForEachNeedExecuteTaskWithPrefix(prifix []byte, f func(key, value []byte) error) error
	ForEachNeedExecuteTask(f func(key, value []byte) error) error
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
	MetadataAuthorityCarrierDB
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
