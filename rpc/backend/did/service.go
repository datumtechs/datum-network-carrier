package did

import (
	"context"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	"github.com/datumtechs/datum-network-carrier/rpc/backend"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/big"
)

func (svr *Server) CreateDID(ctx context.Context, req *emptypb.Empty) (*carrierapipb.CreateDIDResponse, error) {
	didString, err := svr.B.CreateDID()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:CreateDID failed")
		return &carrierapipb.CreateDIDResponse{Status: backend.ErrCreateDID.ErrCode(), Msg: backend.ErrCreateDID.Error(), Did: ""}, nil
	}
	log.Debugf("RPC-API:CreateDID Succeed: didString {%S}", didString)
	return &carrierapipb.CreateDIDResponse{
		Status: 0,
		Msg:    backend.OK,
		Did:    didString,
	}, nil
}

func (svr *Server) CreateVC(ctx context.Context, req *carrierapipb.CreateVCRequest) (*carrierapipb.CreateVCResponse, error) {
	vcJsonString, err := svr.B.CreateVC(req.Did, req.Context, req.PctId, req.Claim, req.ExpirationDate)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:CreateVC failed")
		return &carrierapipb.CreateVCResponse{Status: backend.ErrCreateVC.ErrCode(), Msg: backend.ErrCreateVC.Error(), Vc: ""}, nil
	}
	log.Debugf("RPC-API:CreateVC Succeed: didString {%s}", vcJsonString)
	return &carrierapipb.CreateVCResponse{
		Status: 0,
		Msg:    backend.OK,
		Vc:     vcJsonString,
	}, nil
}

func (svr *Server) SubmitProposal(ctx context.Context, req *carrierapipb.SubmitProposalRequest) (*carrierapipb.SubmitProposalResponse, error) {
	proposalId, err := svr.B.SubmitProposal(int(req.ProposalType), req.ProposalUrl, req.CandidateAddress, req.CandidateServiceUrl)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:SubmitProposal failed")
		return &carrierapipb.SubmitProposalResponse{Status: backend.ErrSubmitProposal.ErrCode(), Msg: backend.ErrSubmitProposal.Error(), ProposalId: ""}, nil
	}
	log.Debugf("RPC-API:SubmitProposal Succeed: proposalId {%s}", proposalId)
	return &carrierapipb.SubmitProposalResponse{
		Status:     0,
		Msg:        backend.OK,
		ProposalId: proposalId,
	}, nil
}

func (svr *Server) WithdrawProposal(ctx context.Context, req *carrierapipb.WithdrawProposalRequest) (*carrierapipb.WithdrawProposalResponse, error) {
	id, ok := new(big.Int).SetString(req.ProposalId, 10)
	if !ok {
		log.Error("RPC-API:WithdrawProposal failed, proposalId is not a valid number")
	}
	result, err := svr.B.WithdrawProposal(id)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:SubmitProposal failed")
		return &carrierapipb.WithdrawProposalResponse{Status: backend.ErrWithdrawProposal.ErrCode(), Msg: backend.ErrWithdrawProposal.Error(), Result: result}, nil
	}
	log.Debugf("RPC-API:WithdrawProposal Succeed: proposalId {%s}", req.ProposalId)
	return &carrierapipb.WithdrawProposalResponse{
		Status: 0,
		Msg:    backend.OK,
		Result: result,
	}, nil
}

func (svr *Server) VoteProposal(ctx context.Context, req *carrierapipb.VoteProposalRequest) (*carrierapipb.VoteProposalResponse, error) {
	id, ok := new(big.Int).SetString(req.ProposalId, 10)
	if !ok {
		log.Error("RPC-API:VoteProposal failed, proposalId is not a valid number")
	}
	result, err := svr.B.VoteProposal(id)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:VoteProposal failed")
		return &carrierapipb.VoteProposalResponse{Status: backend.ErrVoteProposal.ErrCode(), Msg: backend.ErrVoteProposal.Error(), Result: result}, nil
	}
	log.Debugf("RPC-API:VoteProposal Succeed: proposalId {%s}", req.ProposalId)
	return &carrierapipb.VoteProposalResponse{
		Status: 0,
		Msg:    backend.OK,
		Result: result,
	}, nil
}

func (svr *Server) EffectProposal(ctx context.Context, req *carrierapipb.EffectProposalRequest) (*carrierapipb.EffectProposalResponse, error) {
	id, ok := new(big.Int).SetString(req.ProposalId, 10)
	if !ok {
		log.Error("RPC-API:EffectProposal failed, proposalId is not a valid number")
	}
	result, err := svr.B.EffectProposal(id)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:EffectProposal failed")
		return &carrierapipb.EffectProposalResponse{Status: backend.ErrEffectProposal.ErrCode(), Msg: backend.ErrEffectProposal.Error(), Result: result}, nil
	}
	log.Debugf("RPC-API:EffectProposal Succeed: proposalId {%s}", req.ProposalId)
	return &carrierapipb.EffectProposalResponse{
		Status: 0,
		Msg:    backend.OK,
		Result: result,
	}, nil
}
