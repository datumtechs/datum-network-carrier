package scheduler

import "github.com/RosettaFlow/Carrier-Go/types"

type SchedulerStarveFIFO struct {

}

func (sche *SchedulerStarveFIFO) OnSchedule(task *types.TaskMsg) error {
	return nil
}
func (sche *SchedulerStarveFIFO) OnError () error {
	return nil
}
func (sche *SchedulerStarveFIFO) SchedulerName() string {
	return ""
}