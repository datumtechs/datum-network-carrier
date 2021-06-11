package rpc

import (
	"github.com/RosettaFlow/Carrier-Go/types"
)

type Backend interface {
	Start()
	Stop() error
	Status() error
	SendMsg (msg types.Msg) error
 	SetSeedNode (seed *types.SeedNodeInfo) error
	GetSeedNode (id string) (*types.SeedNodeInfo, error)
	GetSeedNodeList () ([]*types.SeedNodeInfo, error)
	SetRegisterNode (node *types.RegisteredNodeInfo) error
	GetRegisterNode (id string) (*types.RegisteredNodeInfo, error)
	GetRegisterNodeList () ([]*types.RegisteredNodeInfo, error)
}
