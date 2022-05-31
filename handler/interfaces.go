package handler

import (
	carriernetmsgtaskmngpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/taskmng"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/types"
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
	ValidateTaskResultMsg(pid peer.ID, taskResultMsg *carriernetmsgtaskmngpb.TaskResultMsg) error
	OnTaskResultMsg(pid peer.ID, taskResultMsg *carriernetmsgtaskmngpb.TaskResultMsg) error
	ValidateTaskResourceUsageMsg(pid peer.ID, taskResourceUsageMsg *carriernetmsgtaskmngpb.TaskResourceUsageMsg) error
	OnTaskResourceUsageMsg(pid peer.ID, taskResourceUsageMsg *carriernetmsgtaskmngpb.TaskResourceUsageMsg) error
	ValidateTaskTerminateMsg(pid peer.ID, terminateMsg *carriernetmsgtaskmngpb.TaskTerminateMsg) error
	OnTaskTerminateMsg (pid peer.ID, terminateMsg *carriernetmsgtaskmngpb.TaskTerminateMsg) error
	SendTaskEvent(event *carriertypespb.TaskEvent) error
	HandleReportResourceUsage(usage *types.TaskResuorceUsage) error
}