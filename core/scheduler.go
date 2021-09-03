package core

import (
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type Scheduler interface {
	Start() error
	Stop() error
	Error () error
	Name() string
	AddTask(task *types.Task)
	RemoveTask(taskId string) error
	TrySchedule() (*types.NeedConsensusTask, error)
	ReplaySchedule(myPartyId string, myTaskRole apipb.TaskRole, task *types.Task) *types.ReplayScheduleResult
}