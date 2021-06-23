package core

import (
	"github.com/RosettaFlow/Carrier-Go/consensus"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type Scheduler interface {
	SetTaskEngine(engine consensus.Consensus) error
	OnSchedule() error
	OnError () error
	SchedulerName() string
	PushTasks(tasks types.TaskMsgs) error
}