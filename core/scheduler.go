package core

import "github.com/RosettaFlow/Carrier-Go/types"

type Scheduler interface {
	OnSchedule() error
	OnError () error
	SchedulerName() string
	PushTasks(tasks types.TaskMsgs) error
}