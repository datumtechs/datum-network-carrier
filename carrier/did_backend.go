package carrier

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/ach/tk"
	"github.com/datumtechs/datum-network-carrier/common/hashutil"
	"github.com/datumtechs/datum-network-carrier/common/hexutil"
	"github.com/datumtechs/datum-network-carrier/grpclient"
	"github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	"github.com/datumtechs/datum-network-carrier/rpc/backend"
	didsdkgocrypto "github.com/datumtechs/did-sdk-go/crypto"
	"github.com/datumtechs/did-sdk-go/did"
	"github.com/datumtechs/did-sdk-go/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/net/context"
	"math/big"
	"strconv"
)

func toApiTxInfo(didTxInfo *did.TransactionInfo) *api.TxInfo {
	if didTxInfo == nil {
		return nil
	}
	return &api.TxInfo{BlockNumber: didTxInfo.BlockNumber, TxIndex: uint32(didTxInfo.TransactionIndex), TxHash: didTxInfo.TxHash}
}

func (s *CarrierAPIBackend) CreateDID() (string, *api.TxInfo, error) {
	req := did.CreateDidReq{}
	req.PrivateKey = tk.WalletManagerInstance().GetPrivateKey()
	req.PublicKey = hex.EncodeToString(crypto.FromECDSAPub(tk.WalletManagerInstance().GetPublicKey()))
	req.PublicKeyType = types.PublicKey_SECP256K1

	response := s.carrier.didService.DocumentService.CreateDID(req)

	if response.Status != did.Response_SUCCESS {
		return "", nil, errors.New(response.Msg)
	}
	return response.Data, toApiTxInfo(&response.TxInfo), nil
}

//接收本地组织admin的VC申请，用本地私钥签名，并调用远端carrier的ApplyVCRemote
func (s *CarrierAPIBackend) ApplyVCLocal(issuerDid, issuerUrl, applicantDid string, pctId uint64, claim, expirationDate, vccontext, extInfo string) error {
	rawData := applicantDid + claim
	reqHash := hashutil.HashSHA256([]byte(rawData))
	sig := didsdkgocrypto.SignSecp256k1(reqHash, tk.WalletManagerInstance().GetPrivateKey())

	reqRemote := new(api.ApplyVCReq)
	reqRemote.Claim = claim
	reqRemote.IssuerDid = issuerDid
	reqRemote.IssuerUrl = issuerUrl
	reqRemote.Context = vccontext
	reqRemote.ApplicantDid = applicantDid
	reqRemote.PctId = pctId
	reqRemote.ExpirationDate = expirationDate
	reqRemote.ExtInfo = extInfo
	reqRemote.ReqDigest = hex.EncodeToString(reqHash)
	reqRemote.ReqSignature = hex.EncodeToString(sig)

	/*issuerAddress, err := types.ParseToAddress(issuerDid)
	if err != nil {
		return err
	}
	response := s.carrier.didService.ProposalService.GetAuthority(issuerAddress)
	if response.Status != did.Response_SUCCESS {
		return errors.New(response.Msg)
	}*/
	/*conn, err := grpc.Dial(issuerUrl, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		log.WithError(err).Error("failed to dial issuer service")
		return err
	}*/

	ctx, cancelFn := context.WithTimeout(context.Background(), grpclient.DefaultGrpcDialTimeout)
	defer cancelFn()

	conn, err := grpclient.DialContext(ctx, issuerUrl, true)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := api.NewVcServiceClient(conn)

	// 请求issuer的carrier
	simpleResp, err := client.ApplyVCRemote(context.Background(), reqRemote)
	if err != nil {
		return err
	}
	if simpleResp.Status != 0 {
		return errors.New(fmt.Sprintf("apply vc error: %d", simpleResp.Status))
	}
	return nil
}

//接收远端组织carier的VC申请，校验申请签名，并调用本地admin的applyVC接口
func (s *CarrierAPIBackend) ApplyVCRemote(issuerDid, applicantDid string, pctId uint64, claim, expirationDate, vccontext, extInfo, reqDigest, reqSignature string) error {
	//从签名恢复的publicKey，必须和document中的一致
	publicKey, err := crypto.SigToPub(hexutil.MustDecode(reqDigest), hexutil.MustDecode(reqSignature))
	if err != nil {
		return errors.New("cannot recover public key from signature")
	}
	// 申请人的document
	docReponse := s.carrier.didService.DocumentService.QueryDidDocument(applicantDid)
	if docReponse.Status != did.Response_SUCCESS {
		return errors.New(fmt.Sprintf("cannot find did document:%s", applicantDid))
	}

	// publicKey是否存在
	didPublicKey := docReponse.Data.FindDidPublicKeyByPublicKey(hex.EncodeToString(crypto.FromECDSAPub(publicKey)))
	if didPublicKey == nil {
		return errors.New("cannot find public key in did document")
	}

	//验证签名
	rawData := applicantDid + claim
	reqHash := hashutil.HashSHA256([]byte(rawData))

	ok := didsdkgocrypto.VerifySecp256k1Signature(reqHash, reqSignature, publicKey)
	if !ok {
		return errors.New("cannot verify the apply VC signature")
	}

	sig := didsdkgocrypto.SignSecp256k1(reqHash, tk.WalletManagerInstance().GetPrivateKey())

	// vc请求校验OK，转发到本地admin;
	reqRemote := new(api.ApplyVCReq)
	reqRemote.Claim = claim
	reqRemote.IssuerDid = issuerDid
	reqRemote.Context = vccontext
	reqRemote.ApplicantDid = applicantDid
	reqRemote.PctId = pctId
	reqRemote.ExpirationDate = expirationDate
	reqRemote.ExtInfo = extInfo
	reqRemote.ReqDigest = hex.EncodeToString(reqHash)
	reqRemote.ReqSignature = hex.EncodeToString(sig)

	//查找本地admin服务端口
	adminGrpcService, err := s.carrier.consulManager.QueryAdminService()
	if err != nil {
		log.WithError(err).Error("cannot find local admin gRPC service")
		return errors.New("cannot find local admin gRPC service")
	}
	if adminGrpcService == nil {
		return errors.New("cannot find local admin gRPC service")
	}
	log.Debugf("adminGrpcService info:%+v", adminGrpcService)

	ctx, cancelFn := context.WithTimeout(context.Background(), grpclient.DefaultGrpcDialTimeout)
	defer cancelFn()

	adminGrpcServiceEndpoint := adminGrpcService.Address + ":" + strconv.Itoa(adminGrpcService.Port)
	conn, err := grpclient.DialContext(ctx, adminGrpcServiceEndpoint, true)
	defer conn.Close()

	if err != nil {
		log.WithError(err).Error("failed to dial admin gRPC service")
		return err
	}
	client := api.NewVcServiceClient(conn)
	simpleResp, err := client.ApplyVCRemote(context.Background(), reqRemote)
	if err != nil {
		log.WithError(err).Error("failed to forward VC apply to admin gRPC service")
		return errors.New("failed to forward VC apply to admin gRPC service")
	}

	if simpleResp.Status != 0 {
		return errors.New(fmt.Sprintf("forward VC apply error: %d", simpleResp.Status))
	}
	return nil
}

//接收本地组织admin的VC申请，用本地私钥签名，并调用远端carrier的ApplyVCRemote
func (s *CarrierAPIBackend) DownloadVCLocal(issuerDid, issuerUrl, applicantDid string) *api.DownloadVCResponse {
	rawData := applicantDid
	reqHash := hashutil.HashSHA256([]byte(rawData))
	sig := didsdkgocrypto.SignSecp256k1(reqHash, tk.WalletManagerInstance().GetPrivateKey())

	reqRemote := new(api.DownloadVCReq)
	reqRemote.IssuerDid = issuerDid
	reqRemote.ApplicantDid = applicantDid
	reqRemote.ReqDigest = hex.EncodeToString(reqHash)
	reqRemote.ReqSignature = hex.EncodeToString(sig)

	/*issuerAddress, err := types.ParseToAddress(issuerDid)
	if err != nil {
		return err
	}
	response := s.carrier.didService.ProposalService.GetAuthority(issuerAddress)
	if response.Status != did.Response_SUCCESS {
		return errors.New(response.Msg)
	}
	conn, err := grpc.Dial(response.Data.Url, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		log.WithError(err).Error("failed to dial issuer service")
		return err
	}*/
	ctx, cancelFn := context.WithTimeout(context.Background(), grpclient.DefaultGrpcDialTimeout)
	defer cancelFn()

	conn, err := grpclient.DialContext(ctx, issuerUrl, true)
	defer conn.Close()
	client := api.NewVcServiceClient(conn)

	// 请求issuer的carrier
	downloadVcResp, err := client.DownloadVCRemote(context.Background(), reqRemote)
	if err != nil {
		log.WithError(err).Errorf("download VC: failed to call issuer, issuerDid:%s, issuerUrl:%s, applicantDid:%s", issuerDid, issuerUrl, applicantDid)
		return &api.DownloadVCResponse{
			Status: 0,
			Msg:    "failed to download VC",
		}
	}
	return downloadVcResp
}

func (s *CarrierAPIBackend) DownloadVCRemote(issuerDid, applicantDid string, reqDigest, reqSignature string) *api.DownloadVCResponse {
	//从签名恢复的publicKey，必须和document中的一致
	failedResponse := &api.DownloadVCResponse{
		Status: backend.ErrDownloadVC.ErrCode(),
		Msg:    backend.ErrDownloadVC.Error(),
	}

	publicKey, err := crypto.SigToPub(hexutil.MustDecode(reqDigest), hexutil.MustDecode(reqSignature))
	if err != nil {
		log.WithError(err).Errorf("download VC: failed to recover public key from signature, issuerDid:%s, applicantDid:%s", issuerDid, applicantDid)
		return failedResponse
	}
	// 申请人的document
	docResponse := s.carrier.didService.DocumentService.QueryDidDocument(applicantDid)
	if docResponse.Status != did.Response_SUCCESS {
		log.Errorf("download VC: failed to find doc document, issuerDid:%s, applicantDid:%s", issuerDid, applicantDid)
		return failedResponse
	}

	// publicKey是否存在
	didPublicKey := docResponse.Data.FindDidPublicKeyByPublicKey(hex.EncodeToString(crypto.FromECDSAPub(publicKey)))
	if didPublicKey == nil {
		log.Errorf("download VC: failed to find public key in document , issuerDid:%s, applicantDid:%s", issuerDid, applicantDid)
		return failedResponse
	}

	//验证签名
	rawData := applicantDid
	reqHash := hashutil.HashSHA256([]byte(rawData))
	ok := didsdkgocrypto.VerifySecp256k1Signature(reqHash, reqSignature, publicKey)
	if !ok {
		log.Errorf("download VC: failed to verify signature, issuerDid:%s, applicantDid:%s", issuerDid, applicantDid)
		return failedResponse
	}

	sig := didsdkgocrypto.SignSecp256k1(reqHash, tk.WalletManagerInstance().GetPrivateKey())

	// vc请求校验OK，转发到本地admin;
	reqRemote := new(api.DownloadVCReq)
	reqRemote.IssuerDid = issuerDid
	reqRemote.ApplicantDid = applicantDid
	reqRemote.ReqDigest = hex.EncodeToString(reqHash)
	reqRemote.ReqSignature = hex.EncodeToString(sig)

	//查找本地admin服务端口
	adminService, err := s.carrier.consulManager.QueryAdminService()
	if err != nil {
		log.WithError(err).Errorf("download VC: failed to find local admin RPC service, issuerDid:%s, applicantDid:%s", issuerDid, applicantDid)
		return failedResponse
	}
	if adminService == nil {
		log.Errorf("download VC: local admin RPC service is none, issuerDid:%s, applicantDid:%s", issuerDid, applicantDid)
		return failedResponse
	}
	log.Debugf("download VC: adminService info:%+v", adminService)

	ctx, cancelFn := context.WithTimeout(context.Background(), grpclient.DefaultGrpcDialTimeout)
	defer cancelFn()

	adminServiceUrl := adminService.Address + ":" + strconv.Itoa(adminService.Port)
	conn, err := grpclient.DialContext(ctx, adminServiceUrl, true)
	defer conn.Close()

	if err != nil {
		log.WithError(err).Errorf("download VC: failed to dial admin RPC service, issuerDid:%s, applicantDid:%s", issuerDid, applicantDid)
		return failedResponse
	}
	client := api.NewVcServiceClient(conn)
	downloadVcResp, err := client.DownloadVCRemote(context.Background(), reqRemote)
	if err != nil {
		log.WithError(err).Errorf("download VC: failed to call admin service, issuerDid:%s, applicantDid:%s", issuerDid, applicantDid)
		return failedResponse
	}
	return downloadVcResp
}

func (s *CarrierAPIBackend) CreateVC(didString string, context string, pctId uint64, claimJson string, expirationDate string) (string, *api.TxInfo, error) {
	var claimMap types.Claim
	err := json.Unmarshal([]byte(claimJson), &claimMap)
	if err != nil {
		return "", nil, err
	}
	req := did.CreateCredentialReq{}
	req.PrivateKey = tk.WalletManagerInstance().GetPrivateKey()
	req.Did = didString
	req.Issuer = types.BuildDid(tk.WalletManagerInstance().GetAddress())
	req.PctId = new(big.Int).SetUint64(pctId)
	req.Context = context
	req.ExpirationDate = expirationDate
	req.Claim = claimMap
	req.Type = types.CREDENTIAL_TYPE_VC
	req.Context = types.DEFAULT_CREDENTIAL_CONTEXT

	response := s.carrier.didService.CredentialService.CreateCredential(req)
	if response.Status != did.Response_SUCCESS {
		return "", nil, errors.New(response.Msg)
	}

	createEvidenceReq := new(did.CreateEvidenceReq)
	createEvidenceReq.PrivateKey = tk.WalletManagerInstance().GetPrivateKey()
	createEvidenceReq.Credential = response.Data
	saveProofResp := s.carrier.didService.CredentialService.CreateEvidence(*createEvidenceReq)

	if saveProofResp.Status != did.Response_SUCCESS {
		return "", nil, errors.New(response.Msg)
	}

	vcBytes, err := json.Marshal(response.Data)
	if err != nil {
		return "", nil, err
	} else {
		return string(vcBytes), toApiTxInfo(&saveProofResp.TxInfo), nil
	}
}

func (s *CarrierAPIBackend) SubmitProposal(proposalType int, proposalUrl string, candidateAddress string, candidateServiceUrl string) (string, *api.TxInfo, error) {
	req := did.SubmitProposalReq{}
	req.PrivateKey = tk.WalletManagerInstance().GetPrivateKey()
	req.ProposalType = uint8(proposalType)
	req.ProposalUrl = proposalUrl
	req.Candidate = ethcommon.HexToAddress(candidateAddress)
	req.CandidateServiceUrl = candidateServiceUrl

	response := s.carrier.didService.ProposalService.SubmitProposal(req)
	if response.Status != did.Response_SUCCESS {
		return "", nil, errors.New(response.Msg)
	}

	return response.Data, toApiTxInfo(&response.TxInfo), nil
}
func (s *CarrierAPIBackend) WithdrawProposal(proposalId *big.Int) (bool, *api.TxInfo, error) {
	req := did.WithdrawProposalReq{}
	req.PrivateKey = tk.WalletManagerInstance().GetPrivateKey()
	req.ProposalId = proposalId
	response := s.carrier.didService.ProposalService.WithdrawProposal(req)
	if response.Status != did.Response_SUCCESS {
		return false, nil, errors.New(response.Msg)
	}

	return response.Data, toApiTxInfo(&response.TxInfo), nil
}
func (s *CarrierAPIBackend) VoteProposal(proposalId *big.Int) (bool, *api.TxInfo, error) {
	req := did.VoteProposalReq{}
	req.PrivateKey = tk.WalletManagerInstance().GetPrivateKey()
	req.ProposalId = proposalId

	response := s.carrier.didService.ProposalService.VoteProposal(req)
	if response.Status != did.Response_SUCCESS {
		return false, nil, errors.New(response.Msg)
	}

	return response.Data, toApiTxInfo(&response.TxInfo), nil
}
func (s *CarrierAPIBackend) EffectProposal(proposalId *big.Int) (bool, *api.TxInfo, error) {
	req := did.EffectProposalReq{}
	req.PrivateKey = tk.WalletManagerInstance().GetPrivateKey()
	req.ProposalId = proposalId

	response := s.carrier.didService.ProposalService.EffectProposal(req)
	if response.Status != did.Response_SUCCESS {
		return false, nil, errors.New(response.Msg)
	}

	return response.Data, toApiTxInfo(&response.TxInfo), nil
}
