package handler

import (
	taskmngpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/taskmng"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
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
	Stop() error
	OnPrepare(task *types.NeedConsensusTask) error
	OnHandle(task *types.NeedConsensusTask) error
	OnConsensusMsg(pid peer.ID, msg types.ConsensusMsg) error
	OnError() error
}

type TaskManager interface {
	Start() error
	Stop() error
	ValidateTaskResultMsg(pid peer.ID, taskResultMsg *taskmngpb.TaskResultMsg) error
	OnTaskResultMsg(pid peer.ID, taskResultMsg *taskmngpb.TaskResultMsg) error
	ValidateTaskResourceUsageMsg(pid peer.ID, taskResourceUsageMsg *taskmngpb.TaskResourceUsageMsg) error
	OnTaskResourceUsageMsg(pid peer.ID, taskResourceUsageMsg *taskmngpb.TaskResourceUsageMsg) error
	ValidateTaskTerminateMsg(pid peer.ID, terminateMsg *taskmngpb.TaskTerminateMsg) error
	OnTaskTerminateMsg (pid peer.ID, terminateMsg *taskmngpb.TaskTerminateMsg) error
	SendTaskEvent(event *libtypes.TaskEvent) error
	HandleReportResourceUsage(usage *types.TaskResuorceUsage) error
}