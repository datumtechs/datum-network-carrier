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
	StoreLocalResourceIdByPowerId(powerId, jobNodeId string) error
	RemoveLocalResourceIdByPowerId(powerId string) error
	QueryLocalResourceIdByPowerId(powerId string) (string, error)

	// about jobRerource
	StoreLocalResourceTable(resource *types.LocalResourceTable) error
	RemoveLocalResourceTable(resourceId string) error
	StoreLocalResourceTables(resources []*types.LocalResourceTable) error
	QueryLocalResourceTable(resourceId string) (*types.LocalResourceTable, error)
	QueryLocalResourceTables() ([]*types.LocalResourceTable, error)
	// about Org power resource
	StoreOrgResourceTable(resource *types.RemoteResourceTable) error
	StoreOrgResourceTables(resources []*types.RemoteResourceTable) error
	RemoveOrgResourceTable(identityId string) error
	QueryOrgResourceTable(identityId string) (*types.RemoteResourceTable, error)
	QueryOrgResourceTables() ([]*types.RemoteResourceTable, error)
	// about slotUnit
	StoreNodeResourceSlotUnit(slot *types.Slot) error
	RemoveNodeResourceSlotUnit() error
	QueryNodeResourceSlotUnit() (*types.Slot, error)
	// about TaskPowerUsed
	StoreLocalTaskPowerUsed(taskPowerUsed *types.LocalTaskPowerUsed) error
	StoreLocalTaskPowerUseds(taskPowerUseds []*types.LocalTaskPowerUsed) error
	RemoveLocalTaskPowerUsed(taskId string) error
	QueryLocalTaskPowerUsed(taskId string) (*types.LocalTaskPowerUsed, error)
	QueryLocalTaskPowerUseds() ([]*types.LocalTaskPowerUsed, error)
	// about DataRereouceTable
	StoreDataRereouceTable(dataRereouceTable *types.DataRereouceTable) error
	StoreDataRereouceTables(dataRereouceTables []*types.DataRereouceTable) error
	RemoveDataRereouceTable(nodeId string) error
	QueryDataRereouceTable(nodeId string) (*types.DataRereouceTable, error)
	QueryDataRereouceTables() ([]*types.DataRereouceTable, error)
	// about DataResourceDataUsed
	StoreDataResourceDataUsed(dataResourceDataUsed *types.DataResourceDataUsed) error
	StoreDataResourceDataUseds(dataResourceDataUseds []*types.DataResourceDataUsed) error
	RemoveDataResourceDataUsed(originId string) error
	QueryDataResourceDataUsed(originId string) (*types.DataResourceDataUsed, error)
	QueryDataResourceDataUseds() ([]*types.DataResourceDataUsed, error)
}

type MetadataCarrierDB interface {
	InsertMetadata(metadata *types.Metadata) error
	GetMetadataByDataId(dataId string) (*types.Metadata, error)
	GetMetadataList() (types.MetadataArray, error)
}

type ResourceCarrierDB interface {
	InsertResource(resource *types.Resource) error
	GetResourceList() (types.ResourceArray, error)
}

type IdentityCarrierDB interface {
	InsertIdentity(identity *types.Identity) error
	StoreIdentity(identity *types.NodeAlias) error
	RemoveIdentity() error
	GetIdentityId() (string, error)
	GetIdentity() (*types.NodeAlias, error)
	RevokeIdentity(identity *types.Identity) error
	GetIdentityList() (types.IdentityArray, error)
	HasIdentity(identity *types.NodeAlias) (bool, error)
}

type TaskCarrierDB interface {
	StoreTaskEvent(event *types.TaskEventInfo) error
	GetTaskEventList(taskId string) ([]*types.TaskEventInfo, error)
	CleanTaskEventList(taskId string) error
	StoreLocalTask(task *types.Task) error
	RemoveLocalTask(taskId string) error
	UpdateLocalTaskState(taskId, state string) error
	GetLocalTask(taskId string) (*types.Task, error)
	GetLocalTaskListByIds(taskIds []string) (types.TaskDataArray, error)
	GetLocalTaskList() (types.TaskDataArray, error)
	StoreJobNodeRunningTaskId(jobNodeId, taskId string) error
	RemoveJobNodeRunningTaskId(jobNodeId, taskId string) error
	GetRunningTaskCountOnJobNode(jobNodeId string) (uint32, error)
	GetJobNodeRunningTaskIdList(jobNodeId string) ([]string, error)
	// about task on datacenter
	InsertTask(task *types.Task) error
	GetTaskList() (types.TaskDataArray, error)
	GetRunningTaskCountOnOrg() uint32
	GetTaskEventListByTaskId(taskId string) ([]*api.TaskEvent, error)
}

type ConsensusDB interface {
	IdentityCarrierDB
	TaskCarrierDB
}

type HandleDB interface {
	LocalStoreCarrierDB
	IdentityCarrierDB
	ResourceCarrierDB
	MetadataCarrierDB
}

type ResourceDB interface {
	LocalStoreCarrierDB
	IdentityCarrierDB
	ResourceCarrierDB
}

type ScheduleDB interface {
	IdentityCarrierDB
	LocalStoreCarrierDB
}