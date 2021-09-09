package iface

import (
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type LocalStoreCarrierDB interface {
	GetYarnName() (string, error)
	SetSeedNode(seed *pb.SeedPeer) (pb.ConnState, error)
	DeleteSeedNode(id string) error
	GetSeedNode(id string) (*pb.SeedPeer, error)
	GetSeedNodeList() ([]*pb.SeedPeer, error)
	SetRegisterNode(typ pb.RegisteredNodeType, node *pb.YarnRegisteredPeerDetail) (pb.ConnState, error)
	DeleteRegisterNode(typ pb.RegisteredNodeType, id string) error
	GetRegisterNode(typ pb.RegisteredNodeType, id string) (*pb.YarnRegisteredPeerDetail, error)
	GetRegisterNodeList(typ pb.RegisteredNodeType) ([]*pb.YarnRegisteredPeerDetail, error)

	InsertLocalResource(resource *types.LocalResource) error
	RemoveLocalResource(jobNodeId string) error
	GetLocalResource(jobNodeId string) (*types.LocalResource, error)
	GetLocalResourceList() (types.LocalResourceArray, error)
	// powerId -> jobNodeId
	StoreLocalResourceIdByPowerId(powerId, jobNodeId string) error
	RemoveLocalResourceIdByPowerId(powerId string) error
	QueryLocalResourceIdByPowerId(powerId string) (string, error)
	//// metaDataId -> dataNodeId
	//StoreLocalResourceIdByMetadataId(metaDataId, dataNodeId string) error
	//RemoveLocalResourceIdByMetadataId(metaDataId string) error
	//QueryLocalResourceIdByMetadataId(metaDataId string) (string, error)

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
	// about task exec status (taskId -> "yes")
	StoreLocalTaskExecuteStatus(taskId string) error
	RemoveLocalTaskExecuteStatus(taskId string) error
	HasLocalTaskExecute(taskId string) (bool, error)
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
	SyncPowerUsed(resource *types.LocalResource) error
}

type IdentityCarrierDB interface {
	InsertIdentity(identity *types.Identity) error
	StoreIdentity(identity *apipb.Organization) error
	RemoveIdentity() error
	GetIdentityId() (string, error)
	GetIdentity() (*apipb.Organization, error)
	RevokeIdentity(identity *types.Identity) error
	GetIdentityList() (types.IdentityArray, error)
	//GetIdentityListByIds(identityIds []string) (types.IdentityArray, error)
	HasIdentity(identity *apipb.Organization) (bool, error)

	// v2.0
	GetMetadataAuthorityList(identityId string, lastUpdate uint64) (types.MetadataAuthArray, error)
}

type TaskCarrierDB interface {
	StoreTaskEvent(event *libTypes.TaskEvent) error
	GetTaskEventList(taskId string) ([]*libTypes.TaskEvent, error)
	RemoveTaskEventList(taskId string) error
	StoreLocalTask(task *types.Task) error
	RemoveLocalTask(taskId string) error
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
	GetTaskEventListByTaskId(taskId string) ([]*libTypes.TaskEvent, error)
	GetTaskEventListByTaskIds(taskIds []string) ([]*libTypes.TaskEvent, error)
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
