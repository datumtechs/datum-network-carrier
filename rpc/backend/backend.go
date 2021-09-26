package backend

import (
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
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

	// v 2.0
	ReportResourceExpense(nodeType pb.NodeType, taskId, ip, port string, usage *libtypes.ResourceUsageOverview) error


	// metadata api

	GetMetadataDetail(identityId, metadataId string) (*types.Metadata, error)
	GetGlobalMetadataDetailList() ([]*pb.GetGlobalMetadataDetailResponse, error)
	GetLocalMetadataDetailList() ([]*pb.GetLocalMetadataDetailResponse, error)
	GetMetadataUsedTaskIdList(identityId, metadataId string) ([]string, error)

	// metadataAuthority api

	AuditMetadataAuthority(audit *types.MetadataAuthAudit) (apicommonpb.AuditMetadataOption, error)
	GetMetadataAuthorityList() (types.MetadataAuthArray, error)
	GetMetadataAuthorityListByUser(userType apicommonpb.UserType, user string) (types.MetadataAuthArray, error)
	HasValidUserMetadataAuth(userType apicommonpb.UserType, user, metadataId string) (bool, error)

	// power api

	GetGlobalPowerDetailList() ([]*pb.GetGlobalPowerDetailResponse, error)
	GetLocalPowerDetailList() ([]*pb.GetLocalPowerDetailResponse, error)

	// identity api

	GetNodeIdentity() (*types.Identity, error)
	GetIdentityList() ([]*types.Identity, error)

	// task api

	GetTaskDetailList() ([]*types.TaskEventShowAndRole, error)
	GetTaskEventList(taskId string) ([]*pb.TaskEventShow, error)
	GetTaskEventListByTaskIds(taskIds []string) ([]*pb.TaskEventShow, error)

	// about DataResourceTable

	StoreDataResourceTable(dataResourceTable *types.DataResourceTable) error
	StoreDataResourceTables(dataResourceTables []*types.DataResourceTable) error
	RemoveDataResourceTable(nodeId string) error
	QueryDataResourceTable(nodeId string) (*types.DataResourceTable, error)
	QueryDataResourceTables() ([]*types.DataResourceTable, error)

	// about DataResourceFileUpload

	StoreDataResourceFileUpload(dataResourceDataUsed *types.DataResourceFileUpload) error
	StoreDataResourceFileUploads(dataResourceDataUseds []*types.DataResourceFileUpload) error
	RemoveDataResourceFileUpload(originId string) error
	QueryDataResourceFileUpload(originId string) (*types.DataResourceFileUpload, error)
	QueryDataResourceFileUploads() ([]*types.DataResourceFileUpload, error)

	// about task result file

	StoreTaskUpResultFile(turf *types.TaskUpResultFile) error
	QueryTaskUpResultFile(taskId string) (*types.TaskUpResultFile, error)
	RemoveTaskUpResultFile(taskId string) error
	StoreTaskResultFileSummary(taskId, originId, filePath, dataNodeId string) error
	QueryTaskResultFileSummary (taskId string) (*types.TaskResultFileSummary, error)
	QueryTaskResultFileSummaryList () (types.TaskResultFileSummaryArr, error)
}
