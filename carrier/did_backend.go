package carrier

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/datumtechs/datum-network-carrier/ach/tk"
	"github.com/datumtechs/did-sdk-go/did"
	"github.com/datumtechs/did-sdk-go/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func (s *CarrierAPIBackend) CreateDID() (string, error) {
	req := did.CreateDidReq{}
	req.PrivateKey = tk.WalletManagerInstance().GetPrivateKey()
	req.PublicKey = hex.EncodeToString(crypto.FromECDSAPub(tk.WalletManagerInstance().GetPublicKey()))
	req.PublicKeyType = types.PublicKey_SECP256K1

	response := s.carrier.didService.DocumentService.CreateDID(req)

	if response.Status != did.Response_SUCCESS {
		return "", errors.New(response.Msg)
	}
	return response.Data, nil
}

func (s *CarrierAPIBackend) CreateVC(didString string, context string, pctId uint64, claimJson string, expirationDate string) (string, error) {
	var claimMap types.Claim
	err := json.Unmarshal([]byte(claimJson), &claimMap)
	if err != nil {
		return "", err
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
		return "", errors.New(response.Msg)
	}

	vcBytes, err := json.Marshal(response.Data)
	if err != nil {
		return "", err
	} else {
		return string(vcBytes), nil
	}
}

func (s *CarrierAPIBackend) SubmitProposal(proposalType int, proposalUrl string, candidateAddress string, candidateServiceUrl string) (string, error) {
	req := did.SubmitProposalReq{}
	req.PrivateKey = tk.WalletManagerInstance().GetPrivateKey()
	req.ProposalType = uint8(proposalType)
	req.ProposalUrl = proposalUrl
	req.Candidate = ethcommon.HexToAddress(candidateAddress)
	req.CandidateServiceUrl = candidateServiceUrl

	response := s.carrier.didService.ProposalService.SubmitProposal(req)
	if response.Status != did.Response_SUCCESS {
		return "", errors.New(response.Msg)
	}

	return response.Data, nil
}
func (s *CarrierAPIBackend) WithdrawProposal(proposalId *big.Int) (bool, error) {
	req := did.WithdrawProposalReq{}
	req.PrivateKey = tk.WalletManagerInstance().GetPrivateKey()
	req.ProposalId = proposalId
	response := s.carrier.didService.ProposalService.WithdrawProposal(req)
	if response.Status != did.Response_SUCCESS {
		return false, errors.New(response.Msg)
	}

	return response.Data, nil
}
func (s *CarrierAPIBackend) VoteProposal(proposalId *big.Int) (bool, error) {

	req := did.VoteProposalReq{}
	req.PrivateKey = tk.WalletManagerInstance().GetPrivateKey()
	req.ProposalId = proposalId

	response := s.carrier.didService.ProposalService.VoteProposal(req)
	if response.Status != did.Response_SUCCESS {
		return false, errors.New(response.Msg)
	}

	return response.Data, nil
}
func (s *CarrierAPIBackend) EffectProposal(proposalId *big.Int) (bool, error) {
	req := did.EffectProposalReq{}
	req.PrivateKey = tk.WalletManagerInstance().GetPrivateKey()
	req.ProposalId = proposalId

	response := s.carrier.didService.ProposalService.EffectProposal(req)
	if response.Status != did.Response_SUCCESS {
		return false, errors.New(response.Msg)
	}

	return response.Data, nil
}
