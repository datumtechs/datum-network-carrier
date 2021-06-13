package carrier

import (
	"github.com/RosettaFlow/Carrier-Go/event"
	"github.com/RosettaFlow/Carrier-Go/types"
)

// CarrierAPIBackend implements rpc.Backend for Carrier
type CarrierAPIBackend struct {
	carrier *Service
}

func (s *CarrierAPIBackend) SendMsg(msg types.Msg) error {
	return s.carrier.mempool.Add(msg)
}

func (s *CarrierAPIBackend) SetSeedNode(seed *types.SeedNodeInfo) (types.NodeConnStatus, error) {
	return s.carrier.datachain.SetSeedNode(seed)
}

func (s *CarrierAPIBackend) DeleteSeedNode(id string) error {
	return s.carrier.datachain.DeleteSeedNode(id)
}

func (s *CarrierAPIBackend) GetSeedNode(id string) (*types.SeedNodeInfo, error) {
	return s.carrier.datachain.GetSeedNode(id)
}

func (s *CarrierAPIBackend) GetSeedNodeList() ([]*types.SeedNodeInfo, error) {
	return s.carrier.datachain.GetSeedNodeList()
}

func (s *CarrierAPIBackend) SetRegisterNode(typ types.RegisteredNodeType, node *types.RegisteredNodeInfo) (types.NodeConnStatus, error) {
	return s.carrier.datachain.SetRegisterNode(typ, node)
}

func (s *CarrierAPIBackend) DeleteRegisterNode(typ types.RegisteredNodeType, id string) error {
	return s.carrier.datachain.DeleteRegisterNode(typ, id)
}

func (s *CarrierAPIBackend) GetRegisterNode(typ types.RegisteredNodeType, id string) (*types.RegisteredNodeInfo, error) {
	return s.carrier.datachain.GetRegisterNode(typ, id)
}

func (s *CarrierAPIBackend) GetRegisterNodeList(typ types.RegisteredNodeType) ([]*types.RegisteredNodeInfo, error) {
	return s.carrier.datachain.GetRegisterNodeList(typ)
}

func (s *CarrierAPIBackend) SendTaskEvent(event *event.TaskEvent) error  {
	return s.carrier.resourceManager.SendTaskEvent(event)
}

