package rpc

import (
	"github.com/RosettaFlow/Carrier-Go/event"
	"github.com/RosettaFlow/Carrier-Go/types"
	pbtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)

type Backend interface {

	SendMsg (msg types.Msg) error

	// system (the yarn node self info)
	GetNodeInfo() (*types.YarnNodeInfo, error)
	GetRegisteredPeers() (*types.YarnRegisteredNodeDetail, error)

	// local node resource api
 	SetSeedNode (seed *types.SeedNodeInfo) (types.NodeConnStatus,error)
	DeleteSeedNode(id string) error
	GetSeedNode (id string) (*types.SeedNodeInfo, error)
	GetSeedNodeList () ([]*types.SeedNodeInfo, error)
	SetRegisterNode (typ types.RegisteredNodeType, node *types.RegisteredNodeInfo) (types.NodeConnStatus,error)
	DeleteRegisterNode (typ types.RegisteredNodeType, id string) error
	GetRegisterNode (typ types.RegisteredNodeType, id string) (*types.RegisteredNodeInfo, error)
	GetRegisterNodeList (typ types.RegisteredNodeType) ([]*types.RegisteredNodeInfo, error)

	SendTaskEvent(event *event.TaskEvent) error

	// identity api
	ApplyIdentityJoin(identity *types.Identity) error
	RevokeIdentityJoin(identity *types.Identity) error

	// power api
	GetPowerTotalSummaryList() ([]*types.OrgResourcePowerAndTaskCount, error)
	GetPowerSingleSummaryList() ([]*types.NodeResourceUsagePowerRes, error)
	GetPowerTotalSummaryByState(state string) ([]*types.OrgResourcePowerAndTaskCount, error)
	GetPowerSingleSummaryByState(state string) ([]*types.NodeResourceUsagePowerRes, error)
	GetPowerTotalSummaryByOwner(identityId string) (*types.OrgResourcePowerAndTaskCount, error)
	GetPowerSingleSummaryByOwner(identityId string) ([]*types.NodeResourceUsagePowerRes, error)
	GetPowerSingleDetail(identityId, powerId string) (*types.OrgPowerTaskDetail, error)
	// metadata api
	GetMetaDataSummaryList() ([]*types.OrgMetaDataSummary, error)
	GetMetaDataSummaryByState(state string) ([]*types.OrgMetaDataSummary, error)
	GetMetaDataSummaryByOwner(identityId string) ([]*types.OrgMetaDataSummary, error)
	GetMetaDataDetail(identityId, metaDataId string) ([]types.OrgMetaDataInfo, error)

	// task api
	GetTaskSummaryList() ([]*types.Task, error)
	GetTaskJoinSummaryList() ([]*types.Task, error)
	GetTaskDetail(taskId string) (*types.Task, error)
	GetTaskEventList(taskId string) ([]*pbtypes.EventData, error)

}
