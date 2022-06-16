package carrier

import (
	"github.com/datumtechs/datum-network-carrier/blacklist"
	"github.com/datumtechs/datum-network-carrier/consensus/twopc"
	carrierrpcdebugpbv1 "github.com/datumtechs/datum-network-carrier/pb/carrier/rpc/debug/v1"
)

// CarrierDebugAPIBackend implements rpc.Backend for Carrier
type CarrierDebugAPIBackend struct {
	engine  *twopc.Twopc
	blackOrg *blacklist.IdentityBackListCache
}

func NewCarrierDebugAPIBackend(engine *twopc.Twopc, blackOrg *blacklist.IdentityBackListCache) *CarrierDebugAPIBackend {
	return &CarrierDebugAPIBackend{
		engine:   engine,
		blackOrg: blackOrg,
	}
}



func (c *CarrierDebugAPIBackend)Get2PcProposalStateByTaskId (taskId string) (*carrierrpcdebugpbv1.Get2PcProposalStateResponse, error) {
	return c.engine.Get2PcProposalStateByProposalId(taskId)
}
func (c *CarrierDebugAPIBackend)Get2PcProposalStateByProposalId (proposalId string) (*carrierrpcdebugpbv1.Get2PcProposalStateResponse, error) {
	return c.engine.Get2PcProposalStateByProposalId(proposalId)
}
func (c *CarrierDebugAPIBackend)Get2PcProposalPrepare (proposalId string) (*carrierrpcdebugpbv1.Get2PcProposalPrepareResponse, error) {
	return c.engine.Get2PcProposalPrepare(proposalId)
}
func (c *CarrierDebugAPIBackend)Get2PcProposalConfirm (proposalId string) (*carrierrpcdebugpbv1.Get2PcProposalConfirmResponse, error) {
	return c.engine.Get2PcProposalConfirm(proposalId)
}
func (c *CarrierDebugAPIBackend) GetAllBlackOrg() (*carrierrpcdebugpbv1.GetConsensusBlackOrgResponse, error) {
	return c.blackOrg.GetAllBlackOrg()
}