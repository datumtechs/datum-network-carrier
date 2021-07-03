package scheduler

import (
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type SchedulerStarveFIFO struct {
	resourceMng    *resource.Manager
	queue          types.TaskMsgs
	scheduledQueue []*types.ScheduleTask
	localTaskCh    chan types.TaskMsgs
	schedTaskCh    chan *types.ConsensusTaskWrap
	remoteTaskCh   chan *types.ScheduleTaskWrap
	err            error
}

func (sche *SchedulerStarveFIFO) NewSchedulerStarveFIFO(
	localTaskCh chan types.TaskMsgs, schedTaskCh chan *types.ConsensusTaskWrap,
	remoteTaskCh chan *types.ScheduleTaskWrap) *SchedulerStarveFIFO {

	return &SchedulerStarveFIFO{
		resourceMng:    resource.NewResourceManager(),
		queue:          make(types.TaskMsgs, 0),
		scheduledQueue: make([]*types.ScheduleTask, 0),
		localTaskCh:    localTaskCh,
		schedTaskCh:    schedTaskCh,
		remoteTaskCh:   remoteTaskCh,
	}
}
func (sche *SchedulerStarveFIFO) OnStart() error {
	go sche.loop()
	return nil
}
func (sche *SchedulerStarveFIFO) OnError() error { return sche.err }
func (sche *SchedulerStarveFIFO) Name() string   { return "SchedulerStarveFIFO" }
func (sche *SchedulerStarveFIFO) loop() {
	for {
		select {
		case task := <-sche.localTaskCh:
			_ = task
			// TODO 还没写完 Ch
		}

		// todo 这里还需要写上 定时调度 队列中的任务信息
	}
}
