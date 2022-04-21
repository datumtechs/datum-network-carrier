package debug

import (
	"context"
	rpcpb "github.com/Metisnetwork/Metis-Carrier/lib/rpc/debug/v1"
)

func (ds *Server) Get2PcProposalStateByTaskId(_ context.Context, req *rpcpb.Get2PcProposalStateByTaskIdRequest) (*rpcpb.Get2PcProposalStateResponse, error) {
	taskId := req.GetTaskId()
	return ds.DebugAPI.Get2PcProposalStateByTaskId(taskId)
}
func (ds *Server) Get2PcProposalStateByProposalId(_ context.Context, req *rpcpb.Get2PcProposalStateByProposalIdRequest) (*rpcpb.Get2PcProposalStateResponse, error) {
	proposalId:=req.GetProposalId()
	return ds.DebugAPI.Get2PcProposalStateByProposalId(proposalId)
}
func (ds *Server) Get2PcProposalPrepare(_ context.Context, req *rpcpb.Get2PcProposalPrepareRequest) (*rpcpb.Get2PcProposalPrepareResponse, error) {
	proposalId:=req.GetProposalId()
	return ds.DebugAPI.Get2PcProposalPrepare(proposalId)
}
func (ds *Server) Get2PcProposalConfirm(_ context.Context, req *rpcpb.Get2PcProposalConfirmRequest) (*rpcpb.Get2PcProposalConfirmResponse, error) {
	proposalId:=req.GetProposalId()
	return ds.DebugAPI.Get2PcProposalConfirm(proposalId)
}
