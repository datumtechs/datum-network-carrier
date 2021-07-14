package scheduler

import (
	"container/heap"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/iface"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/types"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	ReschedMaxCount             = 8
	StarveTerm                  = 3
	defaultScheduleTaskInterval = 20 * time.Millisecond
	electionOrgCondition        = 10000
	electionLocalSeed           = 2
	//taskComputeOrgCount         = 3
)

var (
	ErrEnoughResourceOrgCountLessCalculateCount = fmt.Errorf("the enough resource org count is less calculate count")
	ErrEnoughInternalResourceCount              = fmt.Errorf("has not enough internal resource count")
)

type SchedulerStarveFIFO struct {
	resourceMng *resource.Manager
	// the local task into this queue, first
	queue *types.TaskBullets
	// the very very starve local task by priority
	starveQueue *types.TaskBullets
	// the cache with scheduled local task, will be send to `Consensus`
	scheduledQueue []*types.ScheduleTask
	// fetch local task from taskManager`
	localTaskMsgCh chan types.TaskMsgs
	// send local task scheduled to `Consensus`
	needConsensusTaskCh chan<- *types.ConsensusTaskWrap
	// receive remote task to replay from `Consensus`
	replayScheduleTaskCh <-chan *types.ReplayScheduleTaskWrap
	// 发送经过调度好的 task 交给 taskManager 去分发给自己的 Fighter-Py
	doneSchedTaskCh chan<- *types.DoneScheduleTaskChWrap

	eventEngine *evengine.EventEngine
	dataCenter  iface.ForResourceDB
	err         error
}

func NewSchedulerStarveFIFO(
	eventEngine *evengine.EventEngine,
	mng *resource.Manager,
	dataCenter iface.ForResourceDB,
	localTaskMsgCh chan types.TaskMsgs,
	needConsensusTaskCh chan *types.ConsensusTaskWrap,
	replayScheduleTaskCh chan *types.ReplayScheduleTaskWrap,
	doneSchedTaskCh chan *types.DoneScheduleTaskChWrap,
) *SchedulerStarveFIFO {

	return &SchedulerStarveFIFO{
		resourceMng:          mng,
		queue:                new(types.TaskBullets),
		starveQueue:          new(types.TaskBullets),
		scheduledQueue:       make([]*types.ScheduleTask, 0),
		localTaskMsgCh:       localTaskMsgCh,
		needConsensusTaskCh:  needConsensusTaskCh,
		replayScheduleTaskCh: replayScheduleTaskCh,
		doneSchedTaskCh:      doneSchedTaskCh,
		dataCenter:           dataCenter,
		eventEngine:          eventEngine,
	}
}
func (sche *SchedulerStarveFIFO) loop() {
	taskTimer := time.NewTimer(defaultScheduleTaskInterval)
	for {
		select {
		// From taskManager
		// 新task Msg 到来, 主动触发 调度
		case tasks := <-sche.localTaskMsgCh:

			for _, task := range tasks {
				bullet := types.NewTaskBullet(task)
				sche.addTaskBullet(bullet)
				sche.trySchedule()
			}

		// From Consensus Engine, from remote peer
		// 让自己的Scheduler 重演选举
		case task := <-sche.replayScheduleTaskCh:

			sche.replaySchedule(task)

		// 定时调度 队列中的任务信息
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
func (sche *SchedulerStarveFIFO) OnStop() error  { return nil }  // TODO 未实现 ...
func (sche *SchedulerStarveFIFO) OnError() error { return sche.err }
func (sche *SchedulerStarveFIFO) Name() string   { return "SchedulerStarveFIFO" }
func (sche *SchedulerStarveFIFO) addTaskBullet(bullet *types.TaskBullet) {
	heap.Push(sche.queue, bullet) //
}
func (sche *SchedulerStarveFIFO) trySchedule() error {

	sche.increaseTaskTerm()

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
				// 被丢弃掉的 task  也要清理掉  本地任务的资源, 并提交到数据中心 ...
				log.Error("The number of times the task has been rescheduled exceeds the expected threshold", "taskId", bullet.TaskId)
				sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskDiscarded.Type,
					task.TaskId, task.Onwer().IdentityId, fmt.Sprintf(
						"The number of times the task has been rescheduled exceeds the expected threshold")))

				sendTask := &types.DoneScheduleTaskChWrap{
					ProposalId: common.Hash{},
					SelfTaskRole: types.TaskOnwer,
					// SelfPeerInfo:
					Task: &types.ConsensusScheduleTask{
						TaskDir:   types.SendTaskDir,
						TaskState: types.TaskStateFailed,
						SchedTask: types.ConvertTaskMsgToScheduleTask(task, nil),
					},
					ResultCh: make(chan *types.TaskResultMsgWrap, 0),
				}
				sche.SendTaskToTaskManager(sendTask)
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
		//needSlotCount := sche.resourceMng.GetSlotUnit().CalculateSlotCount(cost.Mem, cost.Processor, cost.Bandwidth)
		//
		////  [选出其他组织的 算力] 如果自己不是 power 角色, 那么就不会与这一步
		//selfResourceInfo, err := sche.electionConputeNode(uint32(needSlotCount))
		//if nil != err {
		//	log.Errorf("Failed to election internal power resource, err: %s", err)
		//	sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskFailedConsensus.Type,
		//		task.TaskId, task.Onwer().IdentityId, err.Error()))
		//	repushFn(bullet)
		//	return
		//}
		//// Lock local resource (jobNode)
		//sche.resourceMng.LockSlot(selfResourceInfo.Id, uint32(needSlotCount))

		// 【选出 其他组织的算力】
		powers, err := sche.electionConputeOrg(task.PowerPartyIds(), task.OwnerIdentityId(), cost)
		if nil != err {
			log.Errorf("Failed to election power org, err: %s", err)
			sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskFailedConsensus.Type,
				task.TaskId, task.Onwer().IdentityId, err.Error()))
			repushFn(bullet)
			return
		}

		// 【选出 发起方 自己的 metaDataId 的 file 对应的  dataNode [ip:port]】
		dataNodeId, err := sche.dataCenter.QueryLocalResourceIdByMetaDataId(task.OwnerTaskSupplier().MetaData.MetaDataId)
		if nil != err {
			log.Errorf("Failed to query localResourceId By MetaDataId: %s, err: %s", task.OwnerTaskSupplier().MetaData.MetaDataId, err)
			sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskFailedConsensus.Type,
				task.TaskId, task.Onwer().IdentityId, err.Error()))
			repushFn(bullet)
			return
		}
		dataNodeResource, err := sche.dataCenter.GetRegisterNode(types.PREFIX_TYPE_DATANODE, dataNodeId)
		if nil != err {
			log.Errorf("Failed to query localResourceInfo By dataNodeId: %s, err: %s", dataNodeId, err)
			sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskFailedConsensus.Type,
				task.TaskId, task.Onwer().IdentityId, err.Error()))
			repushFn(bullet)
			return
		}
		// Send task to consensus Engine to consensus.
		scheduleTask := types.ConvertTaskMsgToScheduleTask(task, powers)
		toConsensusTask := &types.ConsensusTaskWrap{
			Task: scheduleTask,
			OwnerDataResource: &types.PrepareVoteResource{
				Id:   dataNodeId,
				Ip:   dataNodeResource.ExternalIp,
				Port: dataNodeResource.ExternalPort,
			},
			ResultCh: make(chan *types.ConsensuResult, 0),
		}
		sche.SendTaskToConsensus(toConsensusTask)
		consensusRes := toConsensusTask.RecvResult()

		// Consensus failed, task needs to be suspended and rescheduled
		if consensusRes.Status == types.TaskConsensusInterrupt {
			sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskFailedConsensus.Type,
				task.TaskId, task.Onwer().IdentityId, consensusRes.Err.Error()))
			repushFn(bullet)
			return
		}
		/*
			// the task has consensus succeed, need send `Fighter-Py` node
			// (On taskManager)
			sendTask := &types.DoneScheduleTaskChWrap{
				Task: &types.ConsensusScheduleTask{
					TaskDir:                ctypes.SendTaskDir,
					TaskState:              types.TaskStateRunning,
					SchedTask:              scheduleTask,
					//Resources:      		consensusRes.Resources,
				},
				ResultCh: nil,
			}
			sche.SendTaskToTaskManager(sendTask)*/
	}()

	return nil
}
func (sche *SchedulerStarveFIFO) replaySchedule(schedTask *types.ReplayScheduleTaskWrap)  {

	cost := &types.TaskOperationCost{Mem: schedTask.Task.OperationCost.Mem, Processor: schedTask.Task.OperationCost.Processor,
		Bandwidth: schedTask.Task.OperationCost.Bandwidth}

	self, err := sche.dataCenter.GetIdentity()
	if nil != err {
		log.Errorf("Failed to query self identityInfo, err: %s", err)
		schedTask.SendFailedResult(schedTask.Task.TaskId, err)
		return
	}
	// 任务的 重演者 不应该是 任务的发起者
	if self.IdentityId == schedTask.Task.Owner.IdentityId {
		log.Errorf("failed to validate task, self cannot be task owner")
		schedTask.SendFailedResult(schedTask.Task.TaskId, fmt.Errorf("task ower can not replay schedule task"))
		return
	}

	switch schedTask.Role {

	// 如果 当前参与方为 DataSupplier
	case types.DataSupplier:

		powerPartyIds := make([]string, len(schedTask.Task.PowerSuppliers))
		for i, power := range schedTask.Task.PowerSuppliers {
			powerPartyIds[i] = power.PartyId
		}
		// mock election power orgs
		powers, err := sche.electionConputeOrg(powerPartyIds, schedTask.Task.Owner.IdentityId, cost)
		if nil != err {
			log.Errorf("Failed to election power org, err: %s", err)
			sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskFailedConsensus.Type,
				schedTask.Task.TaskId, schedTask.Task.Owner.IdentityId, err.Error()))
			return
		}

		// compare powerSuppliers of task And powerSuppliers of election
		if len(powers) != len(schedTask.Task.PowerSuppliers) {
			err := fmt.Errorf("election powerSuppliers and task powerSuppliers is not match")
			log.Error(err)
			sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskFailedConsensus.Type,
				schedTask.Task.TaskId, schedTask.Task.Owner.IdentityId, err.Error()))

			schedTask.SendResult(&types.ScheduleResult{
				TaskId: schedTask.Task.TaskId,
				Status: types.TaskSchedFailed,
				Err:    fmt.Errorf("replay election powers is not matching recvTask.Powers on replay schedule task"),
				//Resource:
			})

			return
		}

		tmp := make(map[string]struct{}, len(powers))

		for _, power := range powers {
			tmp[power.IdentityId] = struct{}{}
		}
		for _, supplier := range schedTask.Task.PowerSuppliers {
			if _, ok := tmp[supplier.IdentityId]; !ok {
				err := fmt.Errorf("election powerSuppliers and task powerSuppliers is not match")
				log.Error(err)
				sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskFailedConsensus.Type,
					schedTask.Task.TaskId, schedTask.Task.Owner.IdentityId, err.Error()))

				schedTask.SendResult(&types.ScheduleResult{
					TaskId: schedTask.Task.TaskId,
					Status: types.TaskSchedFailed,
					Err:    fmt.Errorf("replay election powers is not matching recvTask.Powers on replay schedule task"),
				})
				return
			}
		}

		// TODO 选出 关于自己 metaDataId 所在的 dataNode
		var metaDataId string
		for _, dataResource := range schedTask.Task.Partners {
			if self.IdentityId == dataResource.IdentityId {
				metaDataId = dataResource.MetaData.MetaDataId
				break
			}
		}
		dataNodeId, err := sche.dataCenter.QueryLocalResourceIdByMetaDataId(metaDataId)
		if nil != err {
			err := fmt.Errorf("failed query internal data node by metaDataId, metaDataId: %s", metaDataId)
			log.Error(err)
			sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskFailedConsensus.Type,
				schedTask.Task.TaskId, schedTask.Task.Owner.IdentityId, err.Error()))

			schedTask.SendResult(&types.ScheduleResult{
				TaskId: schedTask.Task.TaskId,
				Status: types.TaskSchedFailed,
				Err:    fmt.Errorf("failed query internal data node by metaDataId on replay schedule task"),
			})
			return
		}
		dataNode, err := sche.dataCenter.GetRegisterNode(types.PREFIX_TYPE_DATANODE, dataNodeId)
		if nil != err {
			err := fmt.Errorf("failed query internal data node by metaDataId, metaDataId: %s", metaDataId)
			log.Error(err)
			sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskFailedConsensus.Type,
				schedTask.Task.TaskId, schedTask.Task.Owner.IdentityId, err.Error()))

			schedTask.SendResult(&types.ScheduleResult{
				TaskId: schedTask.Task.TaskId,
				Status: types.TaskSchedFailed,
				Err:    fmt.Errorf("failed query internal data node by metaDataId on replay schedule task"),
			})
			return
		}

		schedTask.SendResult(&types.ScheduleResult{
			TaskId: schedTask.Task.TaskId,
			Status: types.TaskSchedOk,
			Resource: &types.PrepareVoteResource{
				Id:   dataNode.Id,
				Ip:   dataNode.ExternalIp,
				Port: dataNode.ExternalPort,
			},
		})

	// 如果 当前参与方为 DataSupplier
	case types.PowerSupplier:
		needSlotCount := sche.resourceMng.GetSlotUnit().CalculateSlotCount(cost.Mem, cost.Processor, cost.Bandwidth)
		selfResourceInfo, err := sche.electionConputeNode(uint32(needSlotCount))
		if nil != err {
			log.Errorf("Failed to election internal power resource, err: %s", err)
			sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskFailedConsensus.Type,
				schedTask.Task.TaskId, schedTask.Task.Owner.IdentityId, err.Error()))

			schedTask.SendResult(&types.ScheduleResult{
				TaskId: schedTask.Task.TaskId,
				Status: types.TaskSchedFailed,
				Err:    fmt.Errorf("failed to replay sched myself local power on replay schedule task"),
			})
			return
		}
		// Lock local resource (jobNode)
		sche.resourceMng.LockSlot(selfResourceInfo.Id, uint32(needSlotCount))
		schedTask.SendResult(&types.ScheduleResult{
			TaskId: schedTask.Task.TaskId,
			Status: types.TaskSchedOk,
			Resource: &types.PrepareVoteResource{
				Id:   selfResourceInfo.Id,
				Ip:   selfResourceInfo.Ip,
				Port: selfResourceInfo.Port,
			},
		})

	// 如果 当前参与方为 ResultSupplier
	case types.ResultSupplier:

		localResourceTables, err := sche.resourceMng.GetLocalResourceTables()
		if nil != err {
			log.Errorf("Failed to election internal data resource, err: %s", err)
			sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskFailedConsensus.Type,
				schedTask.Task.TaskId, schedTask.Task.Owner.IdentityId, err.Error()))

			schedTask.SendResult(&types.ScheduleResult{
				TaskId: schedTask.Task.TaskId,
				Status: types.TaskSchedFailed,
				Err:    fmt.Errorf("failed to replay sched myself local data node on replay schedule task"),
			})
			return
		}

		resource := localResourceTables[len(localResourceTables)-1]
		resourceInfo, err := sche.dataCenter.GetRegisterNode(types.PREFIX_TYPE_DATANODE, resource.GetNodeId())
		if nil != err {
			log.Errorf("Failed to query internal data node resource, err: %s", err)
			sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskFailedConsensus.Type,
				schedTask.Task.TaskId, schedTask.Task.Owner.IdentityId, err.Error()))

			schedTask.SendResult(&types.ScheduleResult{
				TaskId: schedTask.Task.TaskId,
				Status: types.TaskSchedFailed,
				Err:    fmt.Errorf("failed to replay sched myself local data node on replay schedule task"),
			})
			return
		}

		schedTask.SendResult(&types.ScheduleResult{
			TaskId: schedTask.Task.TaskId,
			Status: types.TaskSchedOk,
			Resource: &types.PrepareVoteResource{
				Id:   resourceInfo.Id,
				Ip:   resourceInfo.ExternalIp,
				Port: resourceInfo.ExternalPort,
			},
		})
	}
	return
}

func (sche *SchedulerStarveFIFO) increaseTaskTerm() {
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
func (sche *SchedulerStarveFIFO) electionConputeNode(needSlotCount uint32) (*types.PrepareVoteResource, error) {

	resourceNodeIdArr := make([]string, 0)

	tables, err := sche.resourceMng.GetLocalResourceTables()
	if nil != err {
		return nil, err
	}
	for _, r := range tables {
		if r.IsEnough(needSlotCount) {
			resourceNodeIdArr = append(resourceNodeIdArr, r.GetNodeId())
		}
	}
	if len(resourceNodeIdArr) == 0 {
		return nil, ErrEnoughInternalResourceCount
	}

	resourceId := resourceNodeIdArr[len(resourceNodeIdArr)%electionLocalSeed]
	internalNodeInfo, err := sche.dataCenter.GetRegisterNode(types.PREFIX_TYPE_JOBNODE, resourceId)
	if nil != err {
		return nil, err
	}
	return &types.PrepareVoteResource{
		Id:   resourceId,
		Ip:   internalNodeInfo.ExternalIp,
		Port: internalNodeInfo.ExternalPort,
	}, nil
}
func (sche *SchedulerStarveFIFO) electionConputeOrg(powerPartyIds []string, ownerIdentity string, cost *types.TaskOperationCost) ([]*types.TaskNodeAlias, error) {
	calculateCount := len(powerPartyIds)
	orgs := make([]*types.TaskNodeAlias, calculateCount)
	identityIds := make([]string, 0)

	for _, r := range sche.resourceMng.GetRemoteResouceTables() {
		if r.GetIdentityId() == ownerIdentity {
			continue
		}
		if r.IsEnough(cost.Mem, cost.Processor, cost.Bandwidth) {
			identityIds = append(identityIds, r.GetIdentityId())
		}
	}
	if calculateCount > len(identityIds) {
		return nil, ErrEnoughResourceOrgCountLessCalculateCount
	}
	// Election
	index := electionOrgCondition % len(identityIds)
	identityIdTmp := make(map[string]struct{}, calculateCount)
	for i := calculateCount; i > 0; i-- {
		identityIdTmp[identityIds[index]] = struct{}{}
		index++

	}

	identityArr, err := sche.dataCenter.GetIdentityList()
	if nil != err {
		return nil, err
	}

	i := 0
	for _, iden := range identityArr {
		if _, ok := identityIdTmp[iden.IdentityId()]; ok {
			orgs[i] =  &types.TaskNodeAlias{
				PartyId:    powerPartyIds[i],
				Name:       iden.Name(),
				NodeId:     iden.NodeId(),
				IdentityId: iden.IdentityId(),
			}
		}
	}
	return orgs, nil
}

func (sche *SchedulerStarveFIFO) SendTaskToConsensus(task *types.ConsensusTaskWrap) {
	sche.needConsensusTaskCh <- task
}

func (sche *SchedulerStarveFIFO) SendTaskToTaskManager(task *types.DoneScheduleTaskChWrap) {
	sche.doneSchedTaskCh <- task
}
