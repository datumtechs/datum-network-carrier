package iface

import (
	"github.com/RosettaFlow/Carrier-Go/lib/center/api"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type LocalStoreCarrierDB interface {
	GetYarnName() (string, error)
	SetSeedNode(seed *types.SeedNodeInfo) (types.NodeConnStatus, error)
	DeleteSeedNode(id string) error
	GetSeedNode(id string) (*types.SeedNodeInfo, error)
	GetSeedNodeList() ([]*types.SeedNodeInfo, error)
	SetRegisterNode(typ types.RegisteredNodeType, node *types.RegisteredNodeInfo) (types.NodeConnStatus, error)
	DeleteRegisterNode(typ types.RegisteredNodeType, id string) error
	GetRegisterNode(typ types.RegisteredNodeType, id string) (*types.RegisteredNodeInfo, error)
	GetRegisterNodeList(typ types.RegisteredNodeType) ([]*types.RegisteredNodeInfo, error)

	InsertLocalResource(resource *types.LocalResource) error
	RemoveLocalResource(jobNodeId string) error
	GetLocalResource(jobNodeId string) (*types.LocalResource, error)
	GetLocalResourceList() (types.LocalResourceArray, error)
	// powerId -> jobNodeId
	StoreLocalResourceIdByPowerId(powerId, jobNodeId string) error
	RemoveLocalResourceIdByPowerId(powerId string) error
	QueryLocalResourceIdByPowerId(powerId string) (string, error)
	//// metaDataId -> dataNodeId
    //StoreLocalResourceIdByMetaDataId(metaDataId, dataNodeId string) error
    //RemoveLocalResourceIdByMetaDataId(metaDataId string) error
    //QueryLocalResourceIdByMetaDataId(metaDataId string) (string, error)

	// about jobRerource   (jobNodeId -> {jobNodeId, powerId, resource, slotTotal, slotUsed})
	StoreLocalResourceTable(resource *types.LocalResourceTable) error
	RemoveLocalResourceTable(resourceId string) error
	StoreLocalResourceTables(resources []*types.LocalResourceTable) error
	QueryLocalResourceTable(resourceId string) (*types.LocalResourceTable, error)
	QueryLocalResourceTables() ([]*types.LocalResourceTable, error)
	// about Org power resource (identityId -> {identityId, resourceTotal, resourceUsed})
	StoreOrgResourceTable(resource *types.RemoteResourceTable) error
	StoreOrgResourceTables(resources []*types.RemoteResourceTable) error
	RemoveOrgResourceTable(identityId string) error
	QueryOrgResourceTable(identityId string) (*types.RemoteResourceTable, error)
	QueryOrgResourceTables() ([]*types.RemoteResourceTable, error)
	// about slotUnit (key -> slotUnit)
	StoreNodeResourceSlotUnit(slot *types.Slot) error
	RemoveNodeResourceSlotUnit() error
	QueryNodeResourceSlotUnit() (*types.Slot, error)
	// about TaskPowerUsed  (taskId -> {taskId, jobNodeId, slotCount})
	StoreLocalTaskPowerUsed(taskPowerUsed *types.LocalTaskPowerUsed) error
	StoreLocalTaskPowerUseds(taskPowerUseds []*types.LocalTaskPowerUsed) error
	RemoveLocalTaskPowerUsed(taskId string) error
	QueryLocalTaskPowerUsed(taskId string) (*types.LocalTaskPowerUsed, error)
	QueryLocalTaskPowerUseds() ([]*types.LocalTaskPowerUsed, error)
	// resourceTaskIds Mapping (jobNodeId -> [taskId, taskId, ..., taskId])
	StoreJobNodeRunningTaskId(jobNodeId, taskId string) error
	RemoveJobNodeRunningTaskId(jobNodeId, taskId string) error
	GetRunningTaskCountOnJobNode(jobNodeId string) (uint32, error)
	GetJobNodeRunningTaskIdList(jobNodeId string) ([]string, error)
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
}

type MetadataCarrierDB interface {
	InsertMetadata(metadata *types.Metadata) error
	RevokeMetadata(metadata *types.Metadata) error
	GetMetadataByDataId(dataId string) (*types.Metadata, error)
	GetMetadataList() (types.MetadataArray, error)
}

type ResourceCarrierDB interface {
	InsertResource(resource *types.Resource) error
	RevokeResource(resource *types.Resource) error
	GetResourceList() (types.ResourceArray, error)
	SyncPowerUsed (resource *types.LocalResource) error
}

type IdentityCarrierDB interface {
	InsertIdentity(identity *types.Identity) error
	StoreIdentity(identity *types.NodeAlias) error
	RemoveIdentity() error
	GetIdentityId() (string, error)
	GetIdentity() (*types.NodeAlias, error)
	RevokeIdentity(identity *types.Identity) error
	GetIdentityList() (types.IdentityArray, error)
	//GetIdentityListByIds(identityIds []string) (types.IdentityArray, error)
	HasIdentity(identity *types.NodeAlias) (bool, error)
}

type TaskCarrierDB interface {
	StoreTaskEvent(event *types.TaskEventInfo) error
	GetTaskEventList(taskId string) ([]*types.TaskEventInfo, error)
	RemoveTaskEventList(taskId string) error
	StoreLocalTask(task *types.Task) error
	RemoveLocalTask(taskId string) error
	UpdateLocalTaskState(taskId, state string) error  // 任务的状态 (pending: 等在中; running: 计算中; failed: 失败; success: 成功)
	GetLocalTask(taskId string) (*types.Task, error)
	GetLocalTaskListByIds(taskIds []string) (types.TaskDataArray, error)
	GetLocalTaskList() (types.TaskDataArray, error)
	GetLocalTaskAndEvents(taskId string) (*types.Task, error)
	GetLocalTaskAndEventsListByIds(taskIds []string) (types.TaskDataArray, error)
	GetLocalTaskAndEventsList() (types.TaskDataArray, error)

	// about task on datacenter
	InsertTask(task *types.Task) error
	GetTaskListByIdentityId(identityId string) (types.TaskDataArray, error)
	GetRunningTaskCountOnOrg() uint32
	GetTaskEventListByTaskId(taskId string) ([]*api.TaskEvent, error)
	GetTaskEventListByTaskIds(taskIds []string) ([]*api.TaskEvent, error)
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