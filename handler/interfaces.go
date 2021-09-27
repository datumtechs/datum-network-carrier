package handler

import (
	taskmngpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/taskmng"
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
	OnPrepare(task *types.Task) error
	OnHandle(task *types.Task, result chan<- *types.TaskConsResult) error
	ValidateConsensusMsg(pid peer.ID, msg types.ConsensusMsg) error
	OnConsensusMsg(pid peer.ID, msg types.ConsensusMsg) error
	OnError() error
}

type TaskManager interface {
	Start() error
	Stop() error
	ValidateTaskResultMsg(pid peer.ID, taskResultMsg *taskmngpb.TaskResultMsg) error
	OnTaskResultMsg(pid peer.ID, taskResultMsg *taskmngpb.TaskResultMsg) error
	SendTaskEvent(reportEvent *types.ReportTaskEvent) error
	SendTaskResourceUsage (usage *types.TaskResuorceUsage) error
}