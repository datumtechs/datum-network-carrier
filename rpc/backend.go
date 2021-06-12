package rpc

import (
	"github.com/RosettaFlow/Carrier-Go/event"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type Backend interface {
	//Start()
	//Stop() error
	//Status() error
	SendMsg (msg types.Msg) error
 	SetSeedNode (seed *types.SeedNodeInfo) (types.NodeConnStatus,error)
	DeleteSeedNode(id string) error
	GetSeedNode (id string) (*types.SeedNodeInfo, error)
	GetSeedNodeList () ([]*types.SeedNodeInfo, error)
	SetRegisterNode (typ types.RegisteredNodeType, node *types.RegisteredNodeInfo) (types.NodeConnStatus,error)
	DeleteRegisterNode (typ types.RegisteredNodeType, id string) error
	GetRegisterNode (typ types.RegisteredNodeType, id string) (*types.RegisteredNodeInfo, error)
	GetRegisterNodeList (typ types.RegisteredNodeType) ([]*types.RegisteredNodeInfo, error)

	SendTaskEvent(event *event.TaskEvent) error
}
