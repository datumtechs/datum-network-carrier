package carrier

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/datumtechs/datum-network-carrier/ach/tk"
	"github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	"github.com/datumtechs/did-sdk-go/did"
	"github.com/datumtechs/did-sdk-go/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func toApiTxInfo(didTxInfo *did.TransactionInfo) *api.TxInfo {
	if didTxInfo == nil {
		return nil
	}
	return &api.TxInfo{BlockNumber: didTxInfo.BlockNumber, TxIndex: uint32(didTxInfo.TransactionIndex), TxHash: didTxInfo.TxHash.String()}
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
	req.Context = ""

	response := s.carrier.didService.VcService.CreateCredential(req)
	if response.Status != did.Response_SUCCESS {
		return "", nil, errors.New(response.Msg)
	}

	vcBytes, err := json.Marshal(response.Data)
	if err != nil {
		return "", nil, err
	} else {
		return string(vcBytes), toApiTxInfo(&response.TxInfo), nil
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
