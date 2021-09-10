package backend

import (
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type Backend interface {
	SendMsg(msg types.Msg) error

	// system (the yarn node self info)
	GetNodeInfo() (*pb.YarnNodeInfo, error)
	GetRegisteredPeers() ([]*pb.YarnRegisteredPeer, error)

	// local node resource api
	SetSeedNode(seed *pb.SeedPeer) (pb.ConnState, error)
	DeleteSeedNode(id string) error
	GetSeedNode(id string) (*pb.SeedPeer, error)
	GetSeedNodeList() ([]*pb.SeedPeer, error)
	SetRegisterNode(typ pb.RegisteredNodeType, node *pb.YarnRegisteredPeerDetail) (pb.ConnState, error)
	UpdateRegisterNode(typ pb.RegisteredNodeType, node *pb.YarnRegisteredPeerDetail) (pb.ConnState, error)
	DeleteRegisterNode(typ pb.RegisteredNodeType, id string) error
	GetRegisterNode(typ pb.RegisteredNodeType, id string) (*pb.YarnRegisteredPeerDetail, error)
	GetRegisterNodeList(typ pb.RegisteredNodeType) ([]*pb.YarnRegisteredPeerDetail, error)

	SendTaskEvent(event *types.ReportTaskEvent) error

	// metadata api
	GetMetadataDetail(identityId, metaDataId string) (*pb.GetMetadataDetailResponse, error)
	GetMetadataDetailList() ([]*pb.GetMetadataDetailResponse, error)
	GetMetadataDetailListByOwner(identityId string) ([]*pb.GetMetadataDetailResponse, error)

	// power api
	GetPowerTotalDetailList() ([]*pb.GetPowerTotalDetailResponse, error)
	GetPowerSingleDetailList() ([]*pb.GetPowerSingleDetailResponse, error)

	// identity api
	GetNodeIdentity() (*types.Identity, error)
	GetIdentityList() ([]*types.Identity, error)
	GetMetadataAuthorityList(identityId string, lastUpdate uint64) (types.MetadataAuthArray, error)


	// task api
	GetTaskDetailList() ([]*pb.TaskDetailShow, error)
	GetTaskEventList(taskId string) ([]*pb.TaskEventShow, error)
	GetTaskEventListByTaskIds(taskIds []string) ([]*pb.TaskEventShow, error)

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
