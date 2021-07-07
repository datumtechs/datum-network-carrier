package backend

import (
	"github.com/RosettaFlow/Carrier-Go/types"
)

type Backend interface {
	SendMsg(msg types.Msg) error

	// system (the yarn node self info)
	GetNodeInfo() (*types.YarnNodeInfo, error)
	GetRegisteredPeers() (*types.YarnRegisteredNodeDetail, error)

	// local node resource api
	SetSeedNode(seed *types.SeedNodeInfo) (types.NodeConnStatus, error)
	DeleteSeedNode(id string) error
	GetSeedNode(id string) (*types.SeedNodeInfo, error)
	GetSeedNodeList() ([]*types.SeedNodeInfo, error)
	SetRegisterNode(typ types.RegisteredNodeType, node *types.RegisteredNodeInfo) (types.NodeConnStatus, error)
	DeleteRegisterNode(typ types.RegisteredNodeType, id string) error
	GetRegisterNode(typ types.RegisteredNodeType, id string) (*types.RegisteredNodeInfo, error)
	GetRegisterNodeList(typ types.RegisteredNodeType) ([]*types.RegisteredNodeInfo, error)

	SendTaskEvent(event *types.TaskEventInfo) error

	// metadata api
	GetMetaDataDetail(identityId, metaDataId string) (*types.OrgMetaDataInfo, error)
	GetMetaDataDetailList() ([]*types.OrgMetaDataInfo, error)
	GetMetaDataDetailListByOwner(identityId string) ([]*types.OrgMetaDataInfo, error)


	// power api
	GetPowerTotalDetailList() ([]*types.OrgPowerDetail, error)
	GetPowerSingleDetailList() ([]*types.NodePowerDetail, error)


	// identity api
	GetNodeIdentity() (*types.Identity, error)
	GetIdentityList() ([]*types.Identity, error)

	// task api
	GetTaskDetailList() ([]*types.TaskDetailShow, error)
	GetTaskEventList(taskId string) ([]*types.TaskEvent, error)
}
