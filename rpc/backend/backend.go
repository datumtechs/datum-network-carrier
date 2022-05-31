package backend

import (
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/types"
	"math/big"
)

type Backend interface {
	// add by v 0.4.0
	GetCarrierChainConfig() *types.CarrierChainConfig

	SendMsg(msg types.Msg) error

	// system (the yarn node self info)
	GetNodeInfo() (*carrierapipb.YarnNodeInfo, error)

	// local node resource api

	SetSeedNode(seed *carrierapipb.SeedPeer) (carrierapipb.ConnState, error)
	DeleteSeedNode(addr string) error
	GetSeedNodeList() ([]*carrierapipb.SeedPeer, error)
	SetRegisterNode(typ carrierapipb.RegisteredNodeType, node *carrierapipb.YarnRegisteredPeerDetail) (carrierapipb.ConnState, error)
	UpdateRegisterNode(typ carrierapipb.RegisteredNodeType, node *carrierapipb.YarnRegisteredPeerDetail) (carrierapipb.ConnState, error)
	DeleteRegisterNode(typ carrierapipb.RegisteredNodeType, id string) error
	GetRegisterNode(typ carrierapipb.RegisteredNodeType, id string) (*carrierapipb.YarnRegisteredPeerDetail, error)
	GetRegisterNodeList(typ carrierapipb.RegisteredNodeType) ([]*carrierapipb.YarnRegisteredPeerDetail, error)

	SendTaskEvent(event *carriertypespb.TaskEvent) error

	// v 2.0
	ReportTaskResourceUsage(nodeType carrierapipb.NodeType, ip, port string, usage *types.TaskResuorceUsage) error

	// v 0.4.0
	GenerateObServerProxyWalletAddress() (string, error)

	// metadata api
	IsInternalMetadata(metadataId string) (bool, error)
	GetInternalMetadataDetail(metadataId string) (*types.Metadata, error)
	GetMetadataDetail(metadataId string) (*types.Metadata, error)
	GetGlobalMetadataDetailList(lastUpdate, pageSize uint64) ([]*carrierapipb.GetGlobalMetadataDetail, error)
	GetGlobalMetadataDetailListByIdentityId(identityId string, lastUpdate, pageSize uint64) ([]*carrierapipb.GetGlobalMetadataDetail, error)
	GetLocalMetadataDetailList(lastUpdate, pageSize uint64) ([]*carrierapipb.GetLocalMetadataDetail, error)
	GetLocalInternalMetadataDetailList() ([]*carrierapipb.GetLocalMetadataDetail, error) // add by v 0.3.0
	GetMetadataUsedTaskIdList(identityId, metadataId string) ([]string, error)
	UpdateGlobalMetadata(metadata *types.Metadata) error // add by v 0.4.0

	// metadataAuthority api
	AuditMetadataAuthority(audit *types.MetadataAuthAudit) (commonconstantpb.AuditMetadataOption, error)
	GetLocalMetadataAuthorityList(lastUpdate, pageSize uint64) (types.MetadataAuthArray, error)
	GetGlobalMetadataAuthorityList(lastUpdate, pageSize uint64) (types.MetadataAuthArray, error)
	HasValidMetadataAuth(userType commonconstantpb.UserType, user, identityId, metadataId string) (bool, error)

	// power api
	GetGlobalPowerSummaryList() ([]*carrierapipb.GetGlobalPowerSummary, error)
	GetGlobalPowerDetailList(lastUpdate, pageSize uint64) ([]*carrierapipb.GetGlobalPowerDetail, error)
	GetLocalPowerDetailList() ([]*carrierapipb.GetLocalPowerDetail, error)

	// identity api

	GetNodeIdentity() (*types.Identity, error)
	GetIdentityList(lastUpdate, pageSize uint64) ([]*types.Identity, error)

	// task api
	GetLocalTask(taskId string) (*carriertypespb.TaskDetail, error)
	GetLocalTaskDetailList(lastUpdate, pageSize uint64) ([]*carriertypespb.TaskDetail, error)
	GetGlobalTaskDetailList(lastUpdate, pageSize uint64) ([]*carriertypespb.TaskDetail, error)
	GetTaskDetailListByTaskIds(taskIds []string) ([]*carriertypespb.TaskDetail, error) // v3.0
	GetTaskEventList(taskId string) ([]*carriertypespb.TaskEvent, error)
	GetTaskEventListByTaskIds(taskIds []string) ([]*carriertypespb.TaskEvent, error)
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
	StoreTaskResultFileSummary(taskId, originId, dataHash, metadataOption, dataNodeId, extra string, dataType commonconstantpb.OrigindataType) error
	QueryTaskResultFileSummary(taskId string) (*types.TaskResultFileSummary, error)
	QueryTaskResultFileSummaryList() (types.TaskResultFileSummaryArr, error)

	// v 0.4.0
	EstimateTaskGas(taskSponsorAddress string, dataTokenTransferList []string) (gasLimit uint64, gasPrice *big.Int, err error)
}
