package core

import (
	"github.com/RosettaFlow/Carrier-Go/lib/center/api"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type CarrierDB interface {
	GetIdentityId() (string, error)
	InsertData(blocks types.Blocks) (int, error)
	InsertMetadata(metadata *types.Metadata) error
	StoreIdentity(identity *types.NodeAlias) error
	DelIdentity() error
	GetYarnName() (string, error)
	GetIdentity() (*types.NodeAlias, error)
	GetMetadataByDataId(dataId string) (*types.Metadata, error)
	GetMetadataListByNodeId(nodeId string) (types.MetadataArray, error)
	GetMetadataList() (types.MetadataArray, error)
	HasIdentityId(identityId string) (bool, error)
	HasIdentity(identity *types.NodeAlias) (bool, error)
	InsertResource(resource *types.Resource) error
	GetResourceByDataId(powerId string) (*types.Resource, error)
	GetResourceListByNodeId(nodeId string) (types.ResourceArray, error)
	GetResourceList() (types.ResourceArray, error)
	// InsertIdentity saves new identity info to the center of data.
	InsertIdentity(identity *types.Identity) error
	// RevokeIdentity revokes the identity info to the center of data.
	RevokeIdentity(identity *types.Identity) error
	GetIdentityList() (types.IdentityArray, error)
	GetIdentityByNodeId(nodeId string) (*types.Identity, error)
	InsertTask(task *types.Task) error
	GetTaskList() (types.TaskDataArray, error)
	GetTaskDataListByNodeId(nodeId string) (types.TaskDataArray, error)
	GetTaskEventListByTaskId(taskId string) ([]*api.TaskEvent, error)
	SetSeedNode(seed *types.SeedNodeInfo) (types.NodeConnStatus, error)
	DeleteSeedNode(id string) error
	GetSeedNode(id string) (*types.SeedNodeInfo, error)
	GetSeedNodeList() ([]*types.SeedNodeInfo, error)
	SetRegisterNode(typ types.RegisteredNodeType, node *types.RegisteredNodeInfo) (types.NodeConnStatus, error)
	DeleteRegisterNode(typ types.RegisteredNodeType, id string) error
	GetRegisterNode(typ types.RegisteredNodeType, id string) (*types.RegisteredNodeInfo, error)
	GetRegisterNodeList(typ types.RegisteredNodeType) ([]*types.RegisteredNodeInfo, error)
	StoreRunningTask(task *types.Task) error
	StoreJobNodeRunningTaskId(jobNodeId, taskId string) error
	IncreaseRunningTaskCountOnOrg() uint32
	IncreaseRunningTaskCountOnJobNode(jobNodeId string) uint32
	GetRunningTaskCountOnOrg() uint32
	GetRunningTaskCountOnJobNode(jobNodeId string) uint32
	GetJobNodeRunningTaskIdList(jobNodeId string) []string
}
