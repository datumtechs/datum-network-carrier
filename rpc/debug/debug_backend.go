package debug

import rpcpb "github.com/datumtechs/datum-network-carrier/pb/carrier/rpc/debug/v1"

type TwopcBackend interface {
	// 2pc consensus debug
	Get2PcProposalStateByTaskId(taskId string) (*rpcpb.Get2PcProposalStateResponse, error)
	Get2PcProposalStateByProposalId(proposalId string) (*rpcpb.Get2PcProposalStateResponse, error)
	Get2PcProposalPrepare(proposalId string) (*rpcpb.Get2PcProposalPrepareResponse, error)
	Get2PcProposalConfirm(proposalId string) (*rpcpb.Get2PcProposalConfirmResponse, error)
}

type DebugBackend interface {
	TwopcBackend
}


