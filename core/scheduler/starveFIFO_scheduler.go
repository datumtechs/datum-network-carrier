package scheduler

import (
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/types"
)

const (
	StarveTerm = uint32(3)
)

type SchedulerStarveFIFO struct {
	resourceMng    *resource.Manager
	queue          types.TaskBullets
	starveQueue    types.TaskBullets
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
		queue:          make(types.TaskBullets, 0),
		scheduledQueue: make([]*types.ScheduleTask, 0),
		localTaskCh:    localTaskCh,
		schedTaskCh:    schedTaskCh,
		remoteTaskCh:   remoteTaskCh,
	}
}
func (sche *SchedulerStarveFIFO) loop() {
	for {
		select {
		case tasks := <-sche.localTaskCh:

			for _, task := range tasks {
				bullet := types.NewTaskBullet(task)
				sche.addTaskBullet(bullet)
				sche.trySchedule()
			}

		}

		// todo 这里还需要写上 定时调度 队列中的任务信息
	}
}

func (sche *SchedulerStarveFIFO) OnStart() error {
	err := sche.resourceMng.Start()
	if nil != err {
		return err
	}
	go sche.loop()
	return nil
}
func (sche *SchedulerStarveFIFO) OnError() error { return sche.err }
func (sche *SchedulerStarveFIFO) Name() string   { return "SchedulerStarveFIFO" }
func (sche *SchedulerStarveFIFO) addTaskBullet(bullet *types.TaskBullet) {
	sche.queue = append(sche.queue, bullet)
}
func (sche *SchedulerStarveFIFO) trySchedule() error {

	return nil
}
