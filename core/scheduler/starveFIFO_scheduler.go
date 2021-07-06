package scheduler

import (
	"container/heap"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/db"
	"github.com/RosettaFlow/Carrier-Go/types"
	"time"
)

const (
	StarveTerm                  = uint32(3)
	defaultScheduleTaskInterval = 20 * time.Millisecond
)

type DataCenter interface {
	GetRegisterNodeList(typ types.RegisteredNodeType) ([]*types.RegisteredNodeInfo, error)
	GetIdentity() (*types.NodeAlias, error)
	GetResourceList() (types.ResourceArray, error)
}
type SchedulerStarveFIFO struct {
	resourceMng *resource.Manager
	// the local task into this queue, first
	queue *types.TaskBullets
	// the very very starve local task by priority
	starveQueue *types.TaskBullets
	// the cache with scheduled local task, will be send to `Consensus`
	scheduledQueue []*types.ScheduleTask
	// fetch local task from taskManager`
	localTaskCh chan types.TaskMsgs
	// send local task scheduled to `Consensus`
	schedTaskCh chan *types.ConsensusTaskWrap
	// receive remote task to replay from `Consensus`
	remoteTaskCh chan *types.ScheduleTaskWrap
	dataCenter   DataCenter
	err          error
}

func NewSchedulerStarveFIFO(
	localTaskCh chan types.TaskMsgs, schedTaskCh chan *types.ConsensusTaskWrap,
	remoteTaskCh chan *types.ScheduleTaskWrap, db db.Database, dataCenter DataCenter) *SchedulerStarveFIFO {

	return &SchedulerStarveFIFO{
		resourceMng:    resource.NewResourceManager(db),
		queue:          new(types.TaskBullets),
		starveQueue:    new(types.TaskBullets),
		scheduledQueue: make([]*types.ScheduleTask, 0),
		localTaskCh:    localTaskCh,
		schedTaskCh:    schedTaskCh,
		remoteTaskCh:   remoteTaskCh,
		dataCenter:     dataCenter,
	}
}
func (sche *SchedulerStarveFIFO) loop() {
	taskTimer := time.NewTimer(defaultScheduleTaskInterval)
	for {
		select {
		case tasks := <-sche.localTaskCh:

			for _, task := range tasks {
				bullet := types.NewTaskBullet(task)
				sche.addTaskBullet(bullet)
				sche.trySchedule()
			}

		case task := <-sche.remoteTaskCh:
			// todo 让自己的Scheduler 重演选举


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
func (sche *SchedulerStarveFIFO) OnStop() error  { return nil }
func (sche *SchedulerStarveFIFO) OnError() error { return sche.err }
func (sche *SchedulerStarveFIFO) Name() string   { return "SchedulerStarveFIFO" }
func (sche *SchedulerStarveFIFO) addTaskBullet(bullet *types.TaskBullet) {
	heap.Push(sche.queue, bullet) //
}
func (sche *SchedulerStarveFIFO) trySchedule() error {

	if sche.starveQueue.Len() != 0 {
		x := heap.Pop(sche.starveQueue)
		task := x.(*types.TaskBullet).TaskMsg

	}

	return nil
}

func (sche *SchedulerStarveFIFO) replaySchedule() error {

}

func (sche *SchedulerStarveFIFO) election(calculateNum int, cost *types.TaskOperationCost) []*types.NodeAlias {
	slot := sche.resourceMng.GetSlotUnit()

	slotCount
}
