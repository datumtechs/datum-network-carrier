package consensus

import (
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/libp2p/go-libp2p-core/peer"
)



type Engine interface {
	Start() error
	Stop() error
	OnPrepare(task *types.NeedConsensusTask) error
	OnHandle(task *types.NeedConsensusTask) error
	OnConsensusMsg(pid peer.ID, msg types.ConsensusMsg) error
	OnError() error
}



