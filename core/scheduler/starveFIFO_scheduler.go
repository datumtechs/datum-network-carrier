package scheduler

import (
	"container/heap"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/iface"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
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
func (sche *SchedulerStarveFIFO) OnStop() error  { return nil } // TODO 未实现 ...
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
		taskMsg := bullet.TaskMsg
		repushFn := func(bullet *types.TaskBullet) {

			bullet.IncreaseResched()
			if bullet.Resched > ReschedMaxCount {
				// 被丢弃掉的 task  也要清理掉  本地任务的资源, 并提交到数据中心 ...
				log.Error("The number of times the task has been rescheduled exceeds the expected threshold", "taskId", bullet.TaskId)
				sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskDiscarded.Type,
					taskMsg.TaskId, taskMsg.OwnerIdentityId(), fmt.Sprintf(
						"The number of times the task has been rescheduled exceeds the expected threshold")))

				sendTask := &types.DoneScheduleTaskChWrap{
					ProposalId:   common.Hash{},
					SelfTaskRole: types.TaskOnwer,
					// SelfPeerInfo:
					Task: &types.ConsensusScheduleTask{
						TaskDir:   types.SendTaskDir,
						TaskState: types.TaskStateFailed,
						SchedTask: types.ConvertTaskMsgToTask(taskMsg, nil),
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

		cost := &types.TaskOperationCost{Mem: taskMsg.OperationCost().CostMem, Processor: uint64(taskMsg.OperationCost().CostProcessor),
			Bandwidth: taskMsg.OperationCost().CostBandwidth}
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

		dataIdentityIdCache := map[string]struct{}{taskMsg.OwnerIdentityId(): struct{}{}}
		for _, dataSupplier := range taskMsg.TaskMetadataSupplierDatas() {
			dataIdentityIdCache[dataSupplier.Organization.Identity] = struct{}{}
		}
		for _, receiver := range taskMsg.TaskResultReceiverDatas() {
			dataIdentityIdCache[receiver.Receiver.Identity] = struct{}{}
		}
		// 【选出 其他组织的算力】
		powers, err := sche.electionConputeOrg(taskMsg.PowerPartyIds, dataIdentityIdCache, cost)
		if nil != err {
			log.Errorf("Failed to election power org, err: %s", err)
			sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskFailedConsensus.Type,
				taskMsg.TaskId, taskMsg.OwnerIdentityId(), err.Error()))
			repushFn(bullet)
			return
		}

		var dataResourceDiskUsed *types.DataResourceDiskUsed
		for _, dataSupplier := range taskMsg.TaskMetadataSupplierDatas() {
			// identity 和 partyId 都一致, 才是同一个人 ..
			if taskMsg.OwnerIdentityId() == dataSupplier.Organization.Identity && taskMsg.OwnerPartyId() == dataSupplier.Organization.PartyId {
				// 【选出 发起方 自己的 metaDataId 的 file 对应的  dataNode [ip:port]】
				dataResourceDiskUsed, err = sche.dataCenter.QueryDataResourceDiskUsed(dataSupplier.MetaId)
				if nil != err {
					log.Errorf("Failed to query localResourceId By MetaDataId of task owner: %s, err: %s", dataSupplier.MetaId, err)
					sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskFailedConsensus.Type,
						taskMsg.TaskId, taskMsg.OwnerIdentityId(), err.Error()))
					repushFn(bullet)
					return
				}
				break
			}
		}

		dataNodeResource, err := sche.dataCenter.GetRegisterNode(types.PREFIX_TYPE_DATANODE, dataResourceDiskUsed.GetNodeId())
		if nil != err {
			log.Errorf("Failed to query localResourceInfo By dataNodeId: %s, err: %s", dataResourceDiskUsed.GetNodeId(), err)
			sche.eventEngine.StoreEvent(sche.eventEngine.GenerateEvent(evengine.TaskFailedConsensus.Type,
				taskMsg.TaskId, taskMsg.OwnerIdentityId(), err.Error()))
			repushFn(bullet)
			return
		}

		// Send task to consensus Engine to consensus.
		scheduleTask := types.ConvertTaskMsgToTask(taskMsg, powers)
		toConsensusTask := &types.ConsensusTaskWrap{
			Task: scheduleTask,
			OwnerDataResource: &types.PrepareVoteResource{
				Id:   dataResourceDiskUsed.GetNodeId(),
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
				taskMsg.TaskId, taskMsg.Data.TaskData().Identity, consensusRes.Err.Error()))
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
func (sche *SchedulerStarveFIFO) replaySchedule(replayTask *types.ReplayScheduleTaskWrap) {

	cost := &types.TaskOperationCost{
		Mem: replayTask.Task.TaskData().TaskResource.CostMem,
		Processor: uint64(replayTask.Task.TaskData().TaskResource.CostProcessor),
		Bandwidth: replayTask.Task.TaskData().TaskResource.CostBandwidth,
	}

	self, err := sche.dataCenter.GetIdentity()
	if nil != err {
		log.Errorf("Failed to query self identityInfo, err: %s", err)
		replayTask.SendFailedResult(replayTask.Task.TaskId(), err)
		return
	}
	// 任务的 重演者 不应该是 任务的发起者
	if self.IdentityId == replayTask.Task.TaskData().Identity {
		log.Errorf("failed to validate task, self cannot be task owner")
		replayTask.SendFailedResult(replayTask.Task.TaskId(), fmt.Errorf("task ower can not replay schedule task"))
		return
	}

	switch replayTask.Role {

	// 如果 当前参与方为 DataSupplier   [重新 演算 选 powers]
	case types.DataSupplier:

		powerPartyIds := make([]string, len(replayTask.Task.TaskData().ResourceSupplier))
		for i, power := range replayTask.Task.TaskData().ResourceSupplier {
			powerPartyIds[i] = power.Organization.PartyId
		}

		dataIdentityIdCache := map[string]struct{}{replayTask.Task.TaskData().Identity: {}}
		for _, dataSupplier := range replayTask.Task.TaskData().MetadataSupplier {
			dataIdentityIdCache[dataSupplier.Organization.Identity] = struct{}{}
		}
		for _, receiver := range replayTask.Task.TaskData().Receivers {
			dataIdentityIdCache[receiver.Receiver.Identity] = struct{}{}
		}
		// mock election power orgs
		powers, err := sche.electionConputeOrg(powerPartyIds, dataIdentityIdCache, cost)
		if nil != err {
			log.Errorf("Failed to election power org, err: %s", err)
			replayTask.SendFailedResult(replayTask.Task.TaskId(), fmt.Errorf("failed to election power org, err: %s", err))
			return
		}

		// compare powerSuppliers of task And powerSuppliers of election
		if len(powers) != len(replayTask.Task.TaskData().ResourceSupplier) {
			log.Errorf("election powerSuppliers and task powerSuppliers is not match")
			replayTask.SendFailedResult(replayTask.Task.TaskId(),
				fmt.Errorf("election powerSuppliers and task powerSuppliers is not match"))
			return
		}

		tmp := make(map[string]struct{}, len(powers))

		for _, power := range powers {
			tmp[power.Organization.Identity] = struct{}{}
		}
		for _, power := range replayTask.Task.TaskData().ResourceSupplier {
			if _, ok := tmp[power.Organization.Identity]; !ok {
				log.Errorf("election powerSuppliers and task powerSuppliers is not match")
				replayTask.SendFailedResult(replayTask.Task.TaskId(),
					fmt.Errorf("election powerSuppliers and task powerSuppliers is not match"))
				return
			}
		}

		// 选出 关于自己 metaDataId 所在的 dataNode
		var metaDataId string
		for _, dataResource := range replayTask.Task.TaskData().MetadataSupplier{
			if self.IdentityId == dataResource.Organization.Identity {
				metaDataId = dataResource.MetaId
				break
			}
		}
		dataResourceDiskUsed, err := sche.dataCenter.QueryDataResourceDiskUsed(metaDataId)
		if nil != err {
			log.Errorf("failed query internal data node by metaDataId, metaDataId: %s", metaDataId)
			replayTask.SendFailedResult(replayTask.Task.TaskId(),
				fmt.Errorf("failed query internal data node by metaDataId on replay schedule task"))
			return
		}

		dataNode, err := sche.dataCenter.GetRegisterNode(types.PREFIX_TYPE_DATANODE, dataResourceDiskUsed.GetNodeId())
		if nil != err {
			log.Errorf("failed query internal data node by metaDataId, metaDataId: %s", metaDataId)
			replayTask.SendFailedResult(replayTask.Task.TaskId(),
				fmt.Errorf("failed query internal data node by metaDataId on replay schedule task"))
			return
		}

		replayTask.SendResult(&types.ScheduleResult{
			TaskId: replayTask.Task.TaskId(),
			Status: types.TaskSchedOk,
			Resource: &types.PrepareVoteResource{
				Id:   dataNode.Id,
				Ip:   dataNode.ExternalIp,
				Port: dataNode.ExternalPort,
			},
		})

	// 如果 当前参与方为 PowerSupplier  [选出自己的 内部 power 资源, 并锁定, todo 在最后 DoneXxxxWrap 中解锁]
	case types.PowerSupplier:
		needSlotCount := sche.resourceMng.GetSlotUnit().CalculateSlotCount(cost.Mem, cost.Processor, cost.Bandwidth)
		selfResourceInfo, err := sche.electionConputeNode(uint32(needSlotCount))
		if nil != err {
			log.Errorf("Failed to election internal power resource, err: %s", err)
			replayTask.SendFailedResult(replayTask.Task.TaskId(),
				fmt.Errorf("failed to replay sched myself local power on replay schedule task"))
			return
		}

		if err := sche.resourceMng.LockLocalResourceWithTask(selfResourceInfo.Id, needSlotCount,
			replayTask.Task); nil != err {
			log.Errorf("Failed to Lock LocalResource {%s} With Task {%s}, err: %s", selfResourceInfo.Id, replayTask.Task.TaskId, err)
			replayTask.SendFailedResult(replayTask.Task.TaskId(), err)
			return
		}

		replayTask.SendResult(&types.ScheduleResult{
			TaskId: replayTask.Task.TaskId(),
			Status: types.TaskSchedOk,
			Resource: &types.PrepareVoteResource{
				Id:   selfResourceInfo.Id,
				Ip:   selfResourceInfo.Ip,
				Port: selfResourceInfo.Port,
			},
		})

	// 如果 当前参与方为 ResultSupplier  [仅仅是选出自己可用的 dataNode]
	case types.ResultSupplier:

		localResourceTables, err := sche.resourceMng.GetLocalResourceTables()
		if nil != err {
			log.Errorf("Failed to election internal data resource with replay schedule task, err: %s", err)
			replayTask.SendFailedResult(replayTask.Task.TaskId(),
				fmt.Errorf("failed to election internal data resource with replay schedule task, err: %s", err))
			return
		}

		resource := localResourceTables[len(localResourceTables)-1]
		resourceInfo, err := sche.dataCenter.GetRegisterNode(types.PREFIX_TYPE_DATANODE, resource.GetNodeId())
		if nil != err {
			log.Errorf("Failed to query internal data node resource, err: %s", err)
			replayTask.SendFailedResult(replayTask.Task.TaskId(),
				fmt.Errorf("failed to query internal data node resource, err: %s", err))
			return
		}

		replayTask.SendResult(&types.ScheduleResult{
			TaskId: replayTask.Task.TaskId(),
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
func (sche *SchedulerStarveFIFO) electionConputeOrg(
	powerPartyIds []string,
	dataIdentityIdCache map[string]struct{},
	cost *types.TaskOperationCost,
) ([]*libTypes.TaskResourceSupplierData, error) {


	calculateCount := len(powerPartyIds)
	identityIds := make([]string, 0)

	for _, r := range sche.resourceMng.GetRemoteResouceTables() {
		// 计算方不可以是任务发起方 和 数据参与方 和 接收方
		if _, ok := dataIdentityIdCache[r.GetIdentityId()]; ok {
			continue
		}
		// 还需要有足够的 资源
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

	if len(identityIdTmp) != calculateCount {
		return nil, ErrEnoughResourceOrgCountLessCalculateCount
	}

	identityInfoArr, err := sche.dataCenter.GetIdentityList()
	if nil != err {
		return nil, err
	}
	identityInfoTmp := make(map[string]*types.Identity, calculateCount)
	for _, identityInfo := range identityInfoArr {
		if _, ok := identityIdTmp[identityInfo.IdentityId()]; ok {
			identityInfoTmp[identityInfo.IdentityId()] = identityInfo
		}
	}
	if len(identityInfoTmp) != calculateCount {
		return nil, ErrEnoughResourceOrgCountLessCalculateCount
	}

	resourceArr, err := sche.dataCenter.GetResourceList()
	if nil != err {
		return nil, err
	}

	orgs := make([]*libTypes.TaskResourceSupplierData, calculateCount)
	i := 0
	for _, iden := range resourceArr {

		if i == calculateCount {
			break
		}

		if info, ok := identityInfoTmp[iden.GetIdentityId()]; ok {
			orgs[i] = &libTypes.TaskResourceSupplierData{
				Organization: &libTypes.OrganizationData{
					PartyId:  powerPartyIds[i],
					NodeName: info.Name(),
					NodeId:   info.NodeId(),
					Identity: info.IdentityId(),
				},
				ResourceUsedOverview: &libTypes.ResourceUsedOverview{
					TotalMem:       iden.GetTotalMem(),
					UsedMem:        iden.GetUsedMem(),
					TotalProcessor: uint32(iden.GetTotalProcessor()),
					UsedProcessor:  uint32(iden.GetUsedProcessor()),
					TotalBandwidth: iden.GetTotalBandWidth(),
					UsedBandwidth:  iden.GetUsedBandWidth(),
				},
			}
			i++
			delete(identityInfoTmp, iden.GetIdentityId())
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
