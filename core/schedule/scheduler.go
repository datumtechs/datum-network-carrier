package schedule

import (
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/types"
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
	ReplaySchedule(localPartyId string, localTaskRole carriertypespb.TaskRole, task *types.NeedReplayScheduleTask) *types.ReplayScheduleResult
}
