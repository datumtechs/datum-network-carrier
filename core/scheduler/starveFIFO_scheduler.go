package scheduler

import (
	"container/heap"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/types"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	ReschedMaxCount             = 6
	StarveTerm                  = 6
	defaultScheduleTaskInterval = 20 * time.Millisecond
	electionOrgCondition        = 10000
	electionLocalSeed           = 2
	taskComputeOrgCount         = 3
)

var (
	ErrEnoughResourceOrgCountLessCalculateCount = fmt.Errorf("the enough resource org count is less calculate count")
	ErrEnoughInternalResourceCount              = fmt.Errorf("has not enough internal resource count")
)

type DataCenter interface {
	GetRegisterNode(typ types.RegisteredNodeType, id string) (*types.RegisteredNodeInfo, error)
	GetIdentityList() (types.IdentityArray, error)
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
	remoteTaskCh chan *types.ConsensusScheduleTaskWrap
	// todo  发送经过调度好的 task 交给 taskManager 去分发给自己的 Fighter-Py
	sendSchedTaskCh chan<- *types.ScheduleTask

	eventEngine *evengine.EventEngine
	dataCenter  DataCenter
	err         error
}

func NewSchedulerStarveFIFO(
	localTaskCh chan types.TaskMsgs, schedTaskCh chan *types.ConsensusTaskWrap,
	remoteTaskCh chan *types.ConsensusScheduleTaskWrap, dataCenter DataCenter,
	sendSchedTaskCh chan<- *types.ScheduleTask, mng *resource.Manager,
	eventEngine *evengine.EventEngine) *SchedulerStarveFIFO {

	return &SchedulerStarveFIFO{
		resourceMng:     mng,
		queue:           new(types.TaskBullets),
		starveQueue:     new(types.TaskBullets),
		scheduledQueue:  make([]*types.ScheduleTask, 0),
		localTaskCh:     localTaskCh,
		schedTaskCh:     schedTaskCh,
		remoteTaskCh:    remoteTaskCh,
		sendSchedTaskCh: sendSchedTaskCh,
		dataCenter:      dataCenter,
		eventEngine:     eventEngine,
	}
}
func (sche *SchedulerStarveFIFO) loop() {
	taskTimer := time.NewTimer(defaultScheduleTaskInterval)
	for {
		select {
		// From taskManager
		case tasks := <-sche.localTaskCh:

			for _, task := range tasks {
				bullet := types.NewTaskBullet(task)
				sche.addTaskBullet(bullet)
				sche.trySchedule()
			}

		// From Consensus Engine, from remote peer
		case task := <-sche.remoteTaskCh:
			// todo 让自己的Scheduler 重演选举
			_ = task
			//sche.replaySchedule()

			// todo 这里还需要写上 定时调度 队列中的任务信息
		case <-taskTimer.C:
			sche.trySchedule()
		}

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
	sche.inceaseTaskTerm()

	var bullet *types.TaskBullet

	if sche.starveQueue.Len() != 0 {
		x := heap.Pop(sche.starveQueue)
		bullet = x.(*types.TaskBullet)
	} else {
		x := heap.Pop(sche.queue)
		bullet = x.(*types.TaskBullet)
	}

	go func() {
		task := bullet.TaskMsg
		repushFn := func(bullet *types.TaskBullet) {

			bullet.IncreaseResched()
			if bullet.Resched > ReschedMaxCount {
				log.Error("The number of times the task has been rescheduled exceeds the expected threshold", "taskId", bullet.TaskId)
				sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskDiscarded.Type,
					task.TaskId, task.Onwer().IdentityId, fmt.Sprintf(
						"The number of times the task has been rescheduled exceeds the expected threshold")))
			} else {
				if bullet.Starve {
					heap.Push(sche.starveQueue, bullet)
				} else {
					heap.Push(sche.queue, bullet)
				}
			}
		}

		cost := &types.TaskOperationCost{Mem: task.OperationCost().Mem, Processor: task.OperationCost().Processor,
			Bandwidth: task.OperationCost().Bandwidth}

		selfResourceInfo, err := sche.electionConputeNode(cost)

		needSlotCount := sche.resourceMng.GetSlotUnit().CalculateSlotCount(task.OperationCost().Mem, task.OperationCost().Processor, task.OperationCost().Bandwidth)

		powers, err := sche.electionConputeOrg(taskComputeOrgCount, cost)
		if nil != err {
			log.Errorf("Failed to election power org, err: %s", err)
			sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskFailedConsensus.Type,
				task.TaskId, task.Onwer().IdentityId, err.Error()))
			sche.resourceMng.UnLockSlot(selfResourceInfo.Id, uint32(needSlotCount))
			repushFn(bullet)
			return
		}
		scheduleTask := buildScheduleTask(task, powers)
		resCh := make(chan *types.TaskConsResult, 0)
		sche.schedTaskCh <- &types.ConsensusTaskWrap{
			Task:         scheduleTask,
			SelfResource: selfResourceInfo,
			ResultCh:     resCh,
		}
		res := <-resCh
		// Consensus failed, task needs to be suspended and rescheduled
		if res.Status == types.TaskConsensusInterrupt {
			sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskFailedConsensus.Type,
				task.TaskId, task.Onwer().IdentityId, res.Err.Error()))
			sche.resourceMng.UnLockSlot(selfResourceInfo.Id, uint32(needSlotCount))
			repushFn(bullet)
			return
		}
	}()

	return nil
}

func (sche *SchedulerStarveFIFO) replaySchedule(schedTask *types.ScheduleTask) error {

	return nil
}

func (sche *SchedulerStarveFIFO) inceaseTaskTerm() {
	// handle starve queue
	sche.starveQueue.IncreaseTerm()

	// handle queue
	i := 0
	for {
		if i == sche.queue.Len() {
			return
		}
		bullet := (*(sche.queue))[i]
		bullet.IncreaseTerm()

		// When the task in the queue meets hunger, it will be transferred to starveQueue
		if bullet.Term >= StarveTerm {
			bullet.Starve = true
			heap.Push(sche.starveQueue, bullet)
			heap.Remove(sche.queue, i)
			i = 0
			continue
		}
		(*(sche.queue))[i] = bullet
		i++
	}
}
func (sche *SchedulerStarveFIFO) electionConputeNode(cost *types.TaskOperationCost) (*types.PrepareVoteResource, error) {

	needSlotCount := sche.resourceMng.GetSlotUnit().CalculateSlotCount(cost.Mem, cost.Processor, cost.Bandwidth)

	resourceNodeIdArr := make([]string, 0)

	for _, r := range sche.resourceMng.GetLocalResourceTables() {
		if r.IsEnough(uint32(needSlotCount)) {
			resourceNodeIdArr = append(resourceNodeIdArr, r.GetNodeId())
		}
	}
	if len(resourceNodeIdArr) == 0 {
		return nil, ErrEnoughInternalResourceCount
	}

	resourceId := resourceNodeIdArr[len(resourceNodeIdArr) % electionLocalSeed]
	internalNodeInfo, err := sche.dataCenter.GetRegisterNode(types.PREFIX_TYPE_JOBNODE, resourceId)
	if nil != err {
		return nil, err
	}

	// Lock local resource (jobNode)
	sche.resourceMng.LockSlot(resourceId, uint32(needSlotCount))

	return &types.PrepareVoteResource{
		Id:   resourceId,
		Ip:   internalNodeInfo.ExternalIp,
		Port: internalNodeInfo.ExternalPort,
	}, nil
}
func (sche *SchedulerStarveFIFO) electionConputeOrg(calculateCount int, cost *types.TaskOperationCost) ([]*types.NodeAlias, error) {

	orgs := make([]*types.NodeAlias, 0)
	identityIds := make([]string, 0)

	for _, r := range sche.resourceMng.GetRemoteResouceTables() {
		if r.IsEnough(cost.Mem, cost.Processor, cost.Bandwidth) {
			identityIds = append(identityIds, r.GetIdentityId())
			//identityIdTmp[r.GetIdentityId()] = struct{}{}
		}
	}
	if calculateCount > len(identityIds) {
		return nil, ErrEnoughResourceOrgCountLessCalculateCount
	}
	// Election
	index := electionOrgCondition % len(identityIds)
	identityIdTmp := make(map[string]struct{}, 0)
	for i := calculateCount; i > 0; i-- {

		identityIdTmp[identityIds[index]] = struct{}{}
		index++

	}

	identityArr, err := sche.dataCenter.GetIdentityList()
	if nil != err {
		return nil, err
	}
	for _, iden := range identityArr {
		if _, ok := identityIdTmp[iden.IdentityId()]; ok {
			orgs = append(orgs, &types.NodeAlias{
				Name:       iden.Name(),
				NodeId:     iden.NodeId(),
				IdentityId: iden.IdentityId(),
			})
		}
	}
	return orgs, nil
}

func buildScheduleTask(task *types.TaskMsg, powers []*types.NodeAlias) *types.ScheduleTask {

	partners := make([]*types.ScheduleTaskDataSupplier, len(task.PartnerTaskSuppliers()))
	for i, p := range task.PartnerTaskSuppliers() {
		partner := &types.ScheduleTaskDataSupplier{
			NodeAlias: &types.NodeAlias{
				Name:       p.Name,
				NodeId:     p.NodeId,
				IdentityId: p.IdentityId,
			},
			MetaData: p.MetaData,
		}
		partners[i] = partner
	}

	powerArr := make([]*types.ScheduleTaskPowerSupplier, len(powers))
	for i, p := range powers {
		power := &types.ScheduleTaskPowerSupplier{
			NodeAlias: p,
		}
		powerArr[i] = power
	}

	receivers := make([]*types.ScheduleTaskResultReceiver, len(task.ReceiverDetails()))
	for i, r := range task.ReceiverDetails() {
		receiver := &types.ScheduleTaskResultReceiver{
			NodeAlias: r.NodeAlias,
			Providers: r.Providers,
		}
		receivers[i] = receiver
	}
	return &types.ScheduleTask{
		TaskId:   task.TaskId,
		TaskName: task.TaskName(),
		Owner: &types.ScheduleTaskDataSupplier{
			NodeAlias: task.Onwer(),
			MetaData:  task.OwnerTaskSupplier().MetaData,
		},
		Partners:              partners,
		PowerSuppliers:        powerArr,
		Receivers:             receivers,
		CalculateContractCode: task.CalculateContractCode(),
		DataSplitContractCode: task.DataSplitContractCode(),
		OperationCost:         task.OperationCost(),
		CreateAt:              task.CreateAt(),
	}
}
