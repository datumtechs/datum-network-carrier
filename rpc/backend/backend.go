package backend

import (
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"math/big"
)

type Backend interface {
	SendMsg(msg types.Msg) error

	// system (the yarn node self info)
	GetNodeInfo() (*pb.YarnNodeInfo, error)

	// local node resource api

	SetSeedNode(seed *pb.SeedPeer) (pb.ConnState, error)
	DeleteSeedNode(addr string) error
	GetSeedNodeList() ([]*pb.SeedPeer, error)
	SetRegisterNode(typ pb.RegisteredNodeType, node *pb.YarnRegisteredPeerDetail) (pb.ConnState, error)
	UpdateRegisterNode(typ pb.RegisteredNodeType, node *pb.YarnRegisteredPeerDetail) (pb.ConnState, error)
	DeleteRegisterNode(typ pb.RegisteredNodeType, id string) error
	GetRegisterNode(typ pb.RegisteredNodeType, id string) (*pb.YarnRegisteredPeerDetail, error)
	GetRegisterNodeList(typ pb.RegisteredNodeType) ([]*pb.YarnRegisteredPeerDetail, error)

	SendTaskEvent(event *libtypes.TaskEvent) error

	// v 2.0
	ReportTaskResourceUsage(nodeType pb.NodeType, ip, port string, usage *types.TaskResuorceUsage) error

	// v 0.4.0
	GenerateObServerProxyWalletAddress() (string, error)

	// metadata api
	IsInternalMetadata(metadataId string) (bool, error)
	GetInternalMetadataDetail(metadataId string) (*types.Metadata, error)
	GetMetadataDetail(metadataId string) (*types.Metadata, error)
	GetGlobalMetadataDetailList(lastUpdate, pageSize uint64) ([]*pb.GetGlobalMetadataDetail, error)
	GetGlobalMetadataDetailListByIdentityId(identityId string, lastUpdate, pageSize uint64) ([]*pb.GetGlobalMetadataDetail, error)
	GetLocalMetadataDetailList(lastUpdate, pageSize uint64) ([]*pb.GetLocalMetadataDetail, error)
	GetLocalInternalMetadataDetailList() ([]*pb.GetLocalMetadataDetail, error) // add by v 0.3.0
	GetMetadataUsedTaskIdList(identityId, metadataId string) ([]string, error)
	UpdateGlobalMetadata(metadata *types.Metadata) error // add by v 0.4.0

	// metadataAuthority api
	AuditMetadataAuthority(audit *types.MetadataAuthAudit) (libtypes.AuditMetadataOption, error)
	GetLocalMetadataAuthorityList(lastUpdate, pageSize uint64) (types.MetadataAuthArray, error)
	GetGlobalMetadataAuthorityList(lastUpdate, pageSize uint64) (types.MetadataAuthArray, error)
	HasValidMetadataAuth(userType libtypes.UserType, user, identityId, metadataId string) (bool, error)

	// power api
	GetGlobalPowerSummaryList() ([]*pb.GetGlobalPowerSummary, error)
	GetGlobalPowerDetailList(lastUpdate, pageSize uint64) ([]*pb.GetGlobalPowerDetail, error)
	GetLocalPowerDetailList() ([]*pb.GetLocalPowerDetail, error)

	// identity api

	GetNodeIdentity() (*types.Identity, error)
	GetIdentityList(lastUpdate, pageSize uint64) ([]*types.Identity, error)

	// task api
	GetLocalTask(taskId string) (*pb.TaskDetailShow, error)
	GetLocalTaskDetailList(lastUpdate, pageSize uint64) ([]*pb.TaskDetailShow, error)
	GetGlobalTaskDetailList(lastUpdate, pageSize uint64) ([]*pb.TaskDetailShow, error)
	GetTaskDetailListByTaskIds(taskIds []string) ([]*pb.TaskDetailShow, error) // v3.0
	GetTaskEventList(taskId string) ([]*pb.TaskEventShow, error)
	GetTaskEventListByTaskIds(taskIds []string) ([]*pb.TaskEventShow, error)
	HasLocalTask() (bool, error)

	// about jobResource
	QueryPowerRunningTaskList(powerId string) ([]string, error)

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
	StoreTaskResultFileSummary(taskId, originId, dataHash, filePath, dataNodeId, extra string) error
	QueryTaskResultFileSummary(taskId string) (*types.TaskResultFileSummary, error)
	QueryTaskResultFileSummaryList() (types.TaskResultFileSummaryArr, error)

	// v 0.4.0
	EstimateTaskGas(dataTokenTransferList []*pb.DataTokenTransferItem) (gasLimit uint64, gasPrice *big.Int, err error)
}
