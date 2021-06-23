package consensus

import "github.com/RosettaFlow/Carrier-Go/types"

type Consensus interface {
	OnPrepare(task *types.ScheduleTask) error
	OnStart(task *types.ScheduleTask, result chan<- *types.ScheduleResult) error
	OnError() error
}
