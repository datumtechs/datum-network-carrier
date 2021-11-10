package consensus

import (
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
)



type Engine interface {
	Start() error
	Stop() error
	OnPrepare(task *types.Task) error
	OnHandle(task *types.Task, result chan<- *types.TaskConsResult) error
	OnConsensusMsg(pid peer.ID, msg types.ConsensusMsg) error
	OnError() error
}



