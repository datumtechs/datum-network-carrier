package handler

import (
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

// This defines the interface for interacting with block chain service
type blockchainService interface {

}

// Checker defines a struct which can verify whether a node is currently
// synchronizing a chain with the rest of peers in the network.
type Checker interface {
	Initialized() bool
	Syncing() bool
	Status() error
	Resync() error
}


type Engine interface {
	Start() error
	Close() error
	OnPrepare(task *types.ScheduleTask) error
	OnHandle(task *types.ScheduleTask, result chan<- *types.ConsensuResult) error
	ValidateConsensusMsg(pid peer.ID, msg types.ConsensusMsg) error
	OnConsensusMsg(pid peer.ID, msg types.ConsensusMsg) error
	OnError() error
}