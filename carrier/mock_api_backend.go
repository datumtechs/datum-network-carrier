package carrier

import (
	"context"
	"github.com/datumtechs/datum-network-carrier/core/policy"
	"github.com/datumtechs/datum-network-carrier/grpclient"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/types"
	"math/big"
)

type MockApiBackend struct {
}

func (receiver *MockApiBackend) CheckRequestIpIsPrivate(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetCarrierChainConfig() *types.CarrierChainConfig {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetPolicyEngine() *policy.PolicyEngine {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) SendMsg(msg types.Msg) error {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetNodeInfo() (*carrierapipb.YarnNodeInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) SetSeedNode(seed *carrierapipb.SeedPeer) (carrierapipb.ConnState, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetSeedNodeList() ([]*carrierapipb.SeedPeer, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) SetRegisterNode(typ carrierapipb.RegisteredNodeType, node *carrierapipb.YarnRegisteredPeerDetail) (carrierapipb.ConnState, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) UpdateRegisterNode(typ carrierapipb.RegisteredNodeType, node *carrierapipb.YarnRegisteredPeerDetail) (carrierapipb.ConnState, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) DeleteRegisterNode(typ carrierapipb.RegisteredNodeType, id string) error {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetRegisterNode(typ carrierapipb.RegisteredNodeType, id string) (*carrierapipb.YarnRegisteredPeerDetail, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetRegisterNodeList(typ carrierapipb.RegisteredNodeType) ([]*carrierapipb.YarnRegisteredPeerDetail, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) SendTaskEvent(event *carriertypespb.TaskEvent) error {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) ReportTaskResourceUsage(nodeType carrierapipb.NodeType, ip, port string, usage *types.TaskResuorceUsage) error {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GenerateObServerProxyWalletAddress() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) IsInternalMetadata(metadataId string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetInternalMetadataDetail(metadataId string) (*types.Metadata, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetMetadataDetail(metadataId string) (*types.Metadata, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetGlobalMetadataDetailList(lastUpdate, pageSize uint64) ([]*carrierapipb.GetGlobalMetadataDetail, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetGlobalMetadataDetailListByIdentityId(identityId string, lastUpdate, pageSize uint64) ([]*carrierapipb.GetGlobalMetadataDetail, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetLocalMetadataDetailList(lastUpdate, pageSize uint64) ([]*carrierapipb.GetLocalMetadataDetail, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetLocalInternalMetadataDetailList() ([]*carrierapipb.GetLocalMetadataDetail, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetMetadataUsedTaskIdList(identityId, metadataId string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) UpdateGlobalMetadata(metadata *types.Metadata) error {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) AuditMetadataAuthority(audit *types.MetadataAuthAudit) (commonconstantpb.AuditMetadataOption, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetMetadataAuthority(metadataAuthId string) (*types.MetadataAuthority, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetLocalMetadataAuthorityList(lastUpdate, pageSize uint64) (types.MetadataAuthArray, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetGlobalMetadataAuthorityList(lastUpdate, pageSize uint64) (types.MetadataAuthArray, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) VerifyMetadataAuthInfo(auth *types.MetadataAuthority) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) VerifyMetadataAuthWithMetadataOption(auth *types.MetadataAuthority) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetGlobalPowerSummaryList() ([]*carrierapipb.GetGlobalPowerSummary, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetGlobalPowerDetailList(lastUpdate, pageSize uint64) ([]*carrierapipb.GetGlobalPowerDetail, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetLocalPowerDetailList() ([]*carrierapipb.GetLocalPowerDetail, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetNodeIdentity() (*types.Identity, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetIdentityList(lastUpdate, pageSize uint64) ([]*types.Identity, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetLocalTask(taskId string) (*carriertypespb.TaskDetail, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetLocalTaskDetailList(lastUpdate, pageSize uint64) ([]*carriertypespb.TaskDetail, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetGlobalTaskDetailList(lastUpdate, pageSize uint64) ([]*carriertypespb.TaskDetail, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetTaskDetailListByTaskIds(taskIds []string) ([]*carriertypespb.TaskDetail, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetTaskEventList(taskId string) ([]*carriertypespb.TaskEvent, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetTaskEventListByTaskIds(taskIds []string) ([]*carriertypespb.TaskEvent, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) HasLocalTask() (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) QueryPowerRunningTaskList(powerId string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) StoreDataResourceTable(dataResourceTable *types.DataResourceTable) error {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) StoreDataResourceTables(dataResourceTables []*types.DataResourceTable) error {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) RemoveDataResourceTable(nodeId string) error {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) QueryDataResourceTable(nodeId string) (*types.DataResourceTable, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) QueryDataResourceTables() ([]*types.DataResourceTable, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) StoreDataResourceDataUpload(dataResourceDataUsed *types.DataResourceDataUpload) error {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) StoreDataResourceDataUploads(dataResourceDataUseds []*types.DataResourceDataUpload) error {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) RemoveDataResourceDataUpload(originId string) error {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) QueryDataResourceDataUpload(originId string) (*types.DataResourceDataUpload, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) QueryDataResourceDataUploads() ([]*types.DataResourceDataUpload, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) StoreTaskResultDataSummary(taskId, originId, dataHash, metadataOption, dataNodeId, extra string, dataType commonconstantpb.OrigindataType) error {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) QueryTaskResultDataSummary(taskId string) (*types.TaskResultDataSummary, error) {
	if taskId == "task:0x26915db614c60bc03bf0f2239e5528813bc22c0d3e6848fd45a4fa8681a0192e" {
		return &types.TaskResultDataSummary{OriginId: "OriginId_001", DataType: 1, MetadataOption: `{"dataPath": "./jobnode_service.go"}`}, nil
	}
	return nil, nil
}

func (receiver *MockApiBackend) QueryTaskResultDataSummaryList() (types.TaskResultDataSummaryArr, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) EstimateTaskGas(taskSponsorAddress string, tkItemList []*carrierapipb.TkItem) (gasLimit uint64, gasPrice *big.Int, err error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) GetQueryDataNodeClientByNodeId(nodeId string) (*grpclient.DataNodeClient, bool) {
	if c, err := grpclient.NewDataNodeClient(context.Background(), "127.0.0.1:3344", ""); err == nil {
		return c, true
	}
	return nil, false
}

func (receiver *MockApiBackend) CreateDID() (string, *carrierapipb.TxInfo, error, int) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) ApplyVCLocal(issuerDid, issuerUrl, applicantDid string, pctId uint64, claim, expirationDate, vccontext, extInfo string) error {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) ApplyVCRemote(issuerDid, applicantDid string, pctId uint64, claim, expirationDate, vccontext, extInfo, reqDigest, reqSignature string) error {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) DownloadVCLocal(issuerDid, issuerUrl, applicantDid string) *carrierapipb.DownloadVCResponse {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) DownloadVCRemote(issuerDid, applicantDid string, reqDigest, reqSignature string) *carrierapipb.DownloadVCResponse {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) CreateVC(applicantDid string, context string, pctId uint64, claim string, expirationDate string) (string, *carrierapipb.TxInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) SubmitProposal(proposalType int, proposalUrl string, candidateAddress string, candidateServiceUrl string) (string, *carrierapipb.TxInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) WithdrawProposal(proposalId *big.Int) (bool, *carrierapipb.TxInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) VoteProposal(proposalId *big.Int) (bool, *carrierapipb.TxInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) EffectProposal(proposalId *big.Int) (bool, *carrierapipb.TxInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (receiver *MockApiBackend) DeleteSeedNode(addr string) error {
	return nil
}
