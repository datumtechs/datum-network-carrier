package carrier

import (
	"github.com/RosettaFlow/Carrier-Go/consensus/twopc"
	rpcpb "github.com/RosettaFlow/Carrier-Go/lib/rpc/debug/v1"
)

// CarrierDebugAPIBackend implements rpc.Backend for Carrier
type CarrierDebugAPIBackend struct {
	engine  *twopc.Twopc
}

func NewCarrierDebugAPIBackend(engine  *twopc.Twopc) *CarrierDebugAPIBackend {
	return &CarrierDebugAPIBackend{engine: engine}
}



func (c *CarrierDebugAPIBackend)Get2PcProposalStateByTaskId (taskId string) (*rpcpb.Get2PcProposalStateResponse, error) {
	return c.engine.Get2PcProposalStateByProposalId(taskId)
}
func (c *CarrierDebugAPIBackend)Get2PcProposalStateByProposalId (proposalId string) (*rpcpb.Get2PcProposalStateResponse, error) {
	return c.engine.Get2PcProposalStateByProposalId(proposalId)
}
func (c *CarrierDebugAPIBackend)Get2PcProposalPrepare (proposalId string) (*rpcpb.Get2PcProposalPrepareResponse, error) {
	return c.engine.Get2PcProposalPrepare(proposalId)
}
func (c *CarrierDebugAPIBackend)Get2PcProposalConfirm (proposalId string) (*rpcpb.Get2PcProposalConfirmResponse, error) {
	return c.engine.Get2PcProposalConfirm(proposalId)
}