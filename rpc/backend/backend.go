package backend

import (
	"github.com/datumtechs/datum-network-carrier/core/policy"
	"github.com/datumtechs/datum-network-carrier/grpclient"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/types"
	"math/big"
)

type Backend interface {
	// add by v 0.4.0
	GetCarrierChainConfig() *types.CarrierChainConfig
	GetPolicyEngine() *policy.PolicyEngine
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
	GetMetadataAuthority(metadataAuthId string) (*types.MetadataAuthority, error)
	GetLocalMetadataAuthorityList(lastUpdate, pageSize uint64) (types.MetadataAuthArray, error)
	GetGlobalMetadataAuthorityList(lastUpdate, pageSize uint64) (types.MetadataAuthArray, error)
	VerifyMetadataAuthInfo(auth *types.MetadataAuthority) (bool, error) // add by v0.5.0
	VerifyMetadataAuthWithMetadataOption(auth *types.MetadataAuthority) (bool, error)

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

	// about DataResourceDataUpload
	StoreDataResourceDataUpload(dataResourceDataUsed *types.DataResourceDataUpload) error
	StoreDataResourceDataUploads(dataResourceDataUseds []*types.DataResourceDataUpload) error
	RemoveDataResourceDataUpload(originId string) error
	QueryDataResourceDataUpload(originId string) (*types.DataResourceDataUpload, error)
	QueryDataResourceDataUploads() ([]*types.DataResourceDataUpload, error)

	// about task result data
	StoreTaskResultDataSummary(taskId, originId, dataHash, metadataOption, dataNodeId, extra string, dataType commonconstantpb.OrigindataType) error
	QueryTaskResultDataSummary(taskId string) (*types.TaskResultDataSummary, error)
	QueryTaskResultDataSummaryList() (types.TaskResultDataSummaryArr, error)

	// v 0.4.0
	EstimateTaskGas(taskSponsorAddress string, tkItemList []*carrierapipb.TkItem) (gasLimit uint64, gasPrice *big.Int, err error)
	// v 0.5.0
	GetQueryDataNodeClientByNodeId(nodeId string) (*grpclient.DataNodeClient, bool)

	CreateDID() (string, *carrierapipb.TxInfo, error)
	// CreateVC... if context is empty means default value
	// CreateVC... if expirationDate is empty means default value

	ApplyVCLocal(issuerDid, issuerUrl, applicantDid string, pctId uint64, claim, expirationDate, vccontext, extInfo string) error
	ApplyVCRemote(issuerDid, applicantDid string, pctId uint64, claim, expirationDate, vccontext, extInfo, reqDigest, reqSignature string) error
	DownloadVCLocal(issuerDid, issuerUrl, applicantDid string) *carrierapipb.DownloadVCResponse
	DownloadVCRemote(issuerDid, applicantDid string, reqDigest, reqSignature string) *carrierapipb.DownloadVCResponse

	CreateVC(did string, context string, pctId uint64, claim string, expirationDate string) (string, *carrierapipb.TxInfo, error)
	SubmitProposal(proposalType int, proposalUrl string, candidateAddress string, candidateServiceUrl string) (string, *carrierapipb.TxInfo, error)
	WithdrawProposal(proposalId *big.Int) (bool, *carrierapipb.TxInfo, error)
	VoteProposal(proposalId *big.Int) (bool, *carrierapipb.TxInfo, error)
	EffectProposal(proposalId *big.Int) (bool, *carrierapipb.TxInfo, error)
}
