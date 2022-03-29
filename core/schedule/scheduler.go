package schedule

import (
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type Scheduler interface {
	Start() error
	Stop() error
	Error() error
	Name() string
	AddTask(task *types.Task) error
	RepushTask(task *types.Task) error
	RemoveTask(taskId string) error
	TrySchedule() (*types.NeedConsensusTask, string, error)
	ReplaySchedule(localPartyId string, localTaskRole libtypes.TaskRole, task *types.NeedReplayScheduleTask) *types.ReplayScheduleResult
}
