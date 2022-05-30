package debug

import (
	"context"
	carrierrpcdebugpbv1 "github.com/datumtechs/datum-network-carrier/pb/carrier/rpc/debug/v1"
)

func (ds *Server) Get2PcProposalStateByTaskId(_ context.Context, req *carrierrpcdebugpbv1.Get2PcProposalStateByTaskIdRequest) (*carrierrpcdebugpbv1.Get2PcProposalStateResponse, error) {
	taskId := req.GetTaskId()
	return ds.DebugAPI.Get2PcProposalStateByTaskId(taskId)
}
func (ds *Server) Get2PcProposalStateByProposalId(_ context.Context, req *carrierrpcdebugpbv1.Get2PcProposalStateByProposalIdRequest) (*carrierrpcdebugpbv1.Get2PcProposalStateResponse, error) {
	proposalId:=req.GetProposalId()
	return ds.DebugAPI.Get2PcProposalStateByProposalId(proposalId)
}
func (ds *Server) Get2PcProposalPrepare(_ context.Context, req *carrierrpcdebugpbv1.Get2PcProposalPrepareRequest) (*carrierrpcdebugpbv1.Get2PcProposalPrepareResponse, error) {
	proposalId:=req.GetProposalId()
	return ds.DebugAPI.Get2PcProposalPrepare(proposalId)
}
func (ds *Server) Get2PcProposalConfirm(_ context.Context, req *carrierrpcdebugpbv1.Get2PcProposalConfirmRequest) (*carrierrpcdebugpbv1.Get2PcProposalConfirmResponse, error) {
	proposalId:=req.GetProposalId()
	return ds.DebugAPI.Get2PcProposalConfirm(proposalId)
}
