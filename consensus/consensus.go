package consensus

import (
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
)



type Engine interface {
	Start() error
	Close() error
	OnPrepare(task *types.Task) error
	OnHandle(task *types.Task, result chan<- *types.TaskConsResult) error
	ValidateConsensusMsg(pid peer.ID, msg types.ConsensusMsg) error
	OnConsensusMsg(pid peer.ID, msg types.ConsensusMsg) error
	OnError() error
}



