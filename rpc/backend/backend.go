package backend

import (
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type Backend interface {
	SendMsg(msg types.Msg) error

	// system (the yarn node self info)
	GetNodeInfo() (*pb.YarnNodeInfo, error)
	GetRegisteredPeers() ([]*pb.YarnRegisteredPeer, error)

	// local node resource api
	SetSeedNode(seed *pb.SeedPeer) (types.NodeConnStatus, error)
	DeleteSeedNode(id string) error
	GetSeedNode(id string) (*pb.SeedPeer, error)
	GetSeedNodeList() ([]*pb.SeedPeer, error)
	SetRegisterNode(typ pb.RegisteredNodeType, node *pb.YarnRegisteredPeerDetail) (types.NodeConnStatus, error)
	UpdateRegisterNode(typ pb.RegisteredNodeType, node *pb.YarnRegisteredPeerDetail) (types.NodeConnStatus, error)
	DeleteRegisterNode(typ pb.RegisteredNodeType, id string) error
	GetRegisterNode(typ pb.RegisteredNodeType, id string) (*pb.YarnRegisteredPeerDetail, error)
	GetRegisterNodeList(typ pb.RegisteredNodeType) ([]*pb.YarnRegisteredPeerDetail, error)

	SendTaskEvent(event *libTypes.TaskEvent) error

	// metadata api
	GetMetaDataDetail(identityId, metaDataId string) (*pb.GetMetaDataDetailResponse, error)
	GetMetaDataDetailList() ([]*pb.GetMetaDataDetailResponse, error)
	GetMetaDataDetailListByOwner(identityId string) ([]*pb.GetMetaDataDetailResponse, error)

	// power api
	GetPowerTotalDetailList() ([]*types.OrgPowerDetail, error)
	GetPowerSingleDetailList() ([]*types.NodePowerDetail, error)

	// identity api
	GetNodeIdentity() (*types.Identity, error)
	GetIdentityList() ([]*types.Identity, error)

	// task api
	GetTaskDetailList() ([]*types.TaskDetailShow, error)
	GetTaskEventList(taskId string) ([]*types.TaskEvent, error)
	GetTaskEventListByTaskIds(taskIds []string) ([]*types.TaskEvent, error)

	// about DataResourceTable
	//StoreDataResourceTable(dataResourceTable *types.DataResourceTable) error
	//StoreDataResourceTables(dataResourceTables []*types.DataResourceTable) error
	//RemoveDataResourceTable(nodeId string) error
	//QueryDataResourceTable(nodeId string) (*types.DataResourceTable, error)
	QueryDataResourceTables() ([]*types.DataResourceTable, error)

	// about DataResourceFileUpload
	StoreDataResourceFileUpload(dataResourceDataUsed *types.DataResourceFileUpload) error
	StoreDataResourceFileUploads(dataResourceDataUseds []*types.DataResourceFileUpload) error
	//RemoveDataResourceFileUpload(originId string) error
	QueryDataResourceFileUpload(originId string) (*types.DataResourceFileUpload, error)
	//QueryDataResourceFileUploads() ([]*types.DataResourceFileUpload, error)
}
