package debug

import carrierrpcdebugpbv1 "github.com/datumtechs/datum-network-carrier/pb/carrier/rpc/debug/v1"

type TwopcBackend interface {
	// 2pc consensus debug
	Get2PcProposalStateByTaskId(taskId string) (*carrierrpcdebugpbv1.Get2PcProposalStateResponse, error)
	Get2PcProposalStateByProposalId(proposalId string) (*carrierrpcdebugpbv1.Get2PcProposalStateResponse, error)
	Get2PcProposalPrepare(proposalId string) (*carrierrpcdebugpbv1.Get2PcProposalPrepareResponse, error)
	Get2PcProposalConfirm(proposalId string) (*carrierrpcdebugpbv1.Get2PcProposalConfirmResponse, error)
	GetAllBlackOrg() (*carrierrpcdebugpbv1.GetConsensusBlackOrgResponse, error)
}

type DebugBackend interface {
	TwopcBackend
}


