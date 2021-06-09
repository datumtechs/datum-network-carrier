package core

import "github.com/RosettaFlow/Carrier-Go/types"

type Scheduler interface {
	OnSchedule(task *types.TaskMsg) error
	OnError () error
	SchedulerName() string
}