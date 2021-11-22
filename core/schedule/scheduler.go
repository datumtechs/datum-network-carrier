package schedule

import (
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
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
	TrySchedule() (*types.Task, string, error)
	ReplaySchedule(localPartyId string, localTaskRole apicommonpb.TaskRole, task *types.Task) *types.ReplayScheduleResult
}
