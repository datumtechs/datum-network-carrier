package schedule

import (
	"container/heap"
	"errors"
	"fmt"
	towTypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/types"
	log "github.com/sirupsen/logrus"
)
var (
	ErrRescheduleLargeThreshold = errors.New("The reschedule count of task bullet is large than max threshold")
)


func (sche *SchedulerStarveFIFO) pushTaskBullet(bullet *types.TaskBullet) error {
	sche.queueMutex.Lock()
	defer sche.queueMutex.Unlock()
	// The bullet is first into queue
	old, ok := sche.schedulings[bullet.TaskId]
	if !ok {
		heap.Push(sche.queue, bullet)
		return nil
	}

	bullet = old
	delete(sche.schedulings, bullet.TaskId)

	if bullet.Resched >= ReschedMaxCount {
		return ErrRescheduleLargeThreshold
	} else {
		if bullet.Starve {
			log.Debugf("Task repush  into starve queue, taskId: {%s}, reschedCount: {%d}, max threshold: {%d}",
				bullet.TaskId, bullet.Resched, ReschedMaxCount)
			heap.Push(sche.starveQueue, bullet)
		} else {
			log.Debugf("Task repush  into queue, taskId: {%s}, reschedCount: {%d}, max threshold: {%d}",
				bullet.TaskId, bullet.Resched, ReschedMaxCount)
			heap.Push(sche.queue, bullet)
		}
	}
	return nil
}

func (sche *SchedulerStarveFIFO) removeTaskBullet(taskId string) error {
	sche.queueMutex.Lock()
	defer sche.queueMutex.Unlock()

	// traversal the queue to remove task bullet, first.
	i := 0
	for {
		if i == sche.queue.Len() {
			break
		}
		bullet := (*(sche.queue))[i]

		// When found the bullet with taskId, removed it from queue.
		if bullet.TaskId == taskId {
			heap.Remove(sche.queue, i)
			delete(sche.schedulings, taskId)
			return nil  // todo 这里需要做一次 持久化
		}
		(*(sche.queue))[i] = bullet
		i++
	}

	// otherwise, traversal the starveQueue to remove task bullet, second.
	i = 0
	for {
		if i == sche.starveQueue.Len() {
			break
		}
		bullet := (*(sche.starveQueue))[i]

		// When found the bullet with taskId, removed it from starveQueue.
		if bullet.TaskId == taskId {
			heap.Remove(sche.starveQueue, i)
			delete(sche.schedulings, taskId)
			return nil // todo 这里需要做一次 持久化
		}
		(*(sche.starveQueue))[i] = bullet
		i++
	}
	return nil
}

func (sche *SchedulerStarveFIFO) popTaskBullet() (*types.TaskBullet, error) {
	sche.queueMutex.Lock()
	defer sche.queueMutex.Unlock()

	var bullet *types.TaskBullet

	if sche.starveQueue.Len() != 0 {
		x := heap.Pop(sche.starveQueue)
		bullet = x.(*types.TaskBullet)
	} else {
		if sche.queue.Len() != 0 {
			x := heap.Pop(sche.queue)
			bullet = x.(*types.TaskBullet)
		} else {
			return nil, nil
		}
	}
	bullet.IncreaseResched()
	sche.schedulings[bullet.TaskId] = bullet

	return bullet, nil
}

func (sche *SchedulerStarveFIFO) AddTask(task *types.Task) {
	bullet := types.NewTaskBullet(task.TaskId())
	// todo 这里需要做一次 持久化
	sche.pushTaskBullet(bullet)
}

func (sche *SchedulerStarveFIFO) RemoveTask(taskId string) error {
	return sche.removeTaskBullet(taskId)
}

func (sche *SchedulerStarveFIFO) TrySchedule() (*types.NeedConsensusTask, error) {


	repushFn := func(bullet *types.TaskBullet) {

		if err := sche.pushTaskBullet(bullet); err == ErrRescheduleLargeThreshold {

		}
	}

	sche.increaseTaskTerm()
	bullet, err := sche.popTaskBullet()
	if nil != err {
		log.Errorf("Failed to popTaskBullet on SchedulerStarveFIFO.TrySchedule(), err: {%s}", err)
		return nil, err
	}

	task, err := sche.resourceMng.GetDB().GetLocalTask(bullet.TaskId)
	if nil != err {
		log.Errorf("Failed to QueryLocalTask on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}, err: {%s}", bullet.TaskId, err)

		repushFn(bullet)
		return nil, err
	}


	cost := &towTypes.TaskOperationCost{
		Mem:       task.TaskData().GetOperationCost().GetCostMem(),
		Processor: task.TaskData().GetOperationCost().GetCostProcessor(),
		Bandwidth: task.TaskData().GetOperationCost().GetCostBandwidth(),
		Duration:  task.TaskData().GetOperationCost().GetDuration(),
	}

	log.Debugf("Call SchedulerStarveFIFO.TrySchedule() start, taskId: {%s}, partyId: {%s}, taskCost: {%s}",
		task.TaskData().TaskId, task.TaskData().PartyId, cost.String())


	// 获取 powerPartyIds 标签 TODO

	// 【选出 其他组织的算力】
	powers, err := sche.electionConputeOrg(nil, nil, cost)
	if nil != err {
		log.Errorf("Failed to election powers org on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}, err: {%s}", task.TaskId(), err)

		repushFn(bullet)
		return nil, err
	}

	log.Debugf("Succeed to election powers org on SchedulerStarveFIFO.TrySchedule(), taskId {%s}, powers: %s", task.TaskId(), utilOrgPowerArrString(powers))

	// Set elected powers into task info, and restore into local db.
	task = types.ConvertTaskMsgToTaskWithPowers(task, powers)
	// restore task by power
	if err := sche.dataCenter.StoreLocalTask(task); nil != err {
		log.Errorf("Failed tp update local task by election powers on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}, err: {%s}", task.TaskId(), err)

		repushFn(bullet)
		return nil, err
	}
	return types.NewNeedConsensusTask(task), nil
}
func (sche *SchedulerStarveFIFO) ReplaySchedule(myPartyId string, myTaskRole apipb.TaskRole, task *types.Task) *types.ReplayScheduleResult {

	cost := &towTypes.TaskOperationCost{
		Mem:       task.TaskData().GetOperationCost().GetCostMem(),
		Processor: task.TaskData().GetOperationCost().GetCostProcessor(),
		Bandwidth: task.TaskData().GetOperationCost().GetCostBandwidth(),
		Duration:  task.TaskData().GetOperationCost().GetDuration(),
	}

	log.Debugf("Call SchedulerStarveFIFO.ReplaySchedule() start, taskId: {%s}, myTaskRole: {%s}, myPartyId: {%s}, taskCost: {%s}",
		task.TaskId(), myTaskRole.String(), myPartyId, cost.String())

	selfIdentityId, err := sche.dataCenter.GetIdentityId()
	if nil != err {
		log.Errorf("Failed to query self identityInfo on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, err: {%s}", task.TaskId(), err)
		return types.NewReplayScheduleResult(task.TaskId(), err, nil)
	}


	var result *types.ReplayScheduleResult

	switch myTaskRole {

	// 如果 当前参与方为 DataSupplier   [重新 演算 选 powers]
	case apipb.TaskRole_TaskRole_DataSupplier:

		var isSender bool
		if myPartyId == task.TaskSender().GetPartyId() && selfIdentityId == task.TaskSender().GetIdentityId() {
			isSender = true
		}

		if !isSender {

			powerPartyIds := make([]string, len(task.TaskData().GetPowerSuppliers()))
			for i, power := range task.TaskData().GetPowerSuppliers() {
				powerPartyIds[i] = power.GetOrganization().GetPartyId()
			}
			// mock election power orgs
			powers, err := sche.electionConputeOrg(powerPartyIds, nil, cost)
			if nil != err {
				log.Errorf("Failed to election powers org on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, err: {%s}",
					task.TaskId(), err)

				return types.NewReplayScheduleResult(task.TaskId(), err, nil)
			}

			log.Debugf("Succeed to election powers org on SchedulerStarveFIFO.ReplaySchedule(), taskId {%s}, powers: %s",
				task.TaskId(), utilOrgPowerArrString(powers))

			// compare powerSuppliers of task And powerSuppliers of election
			if len(powers) != len(task.TaskData().GetPowerSuppliers()) {
				log.Errorf("reschedule powers len and task powers len is not match on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, reschedule power len: {%d}, task powers len: {%d}",
					task.TaskId(), len(powers), len(task.TaskData().GetPowerSuppliers()))

				return types.NewReplayScheduleResult(task.TaskId(), err, nil)
			}

			tmp := make(map[string]struct{}, len(powers))

			for _, power := range powers {
				tmp[power.GetOrganization().GetIdentityId()] = struct{}{}
			}
			for _, power := range task.TaskData().GetPowerSuppliers() {
				if _, ok := tmp[power.GetOrganization().GetIdentityId()]; !ok {
					log.Errorf("task power identityId not found on reschedule powers on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, task power identityId: {%s}",
						task.TaskId(), power.GetOrganization().GetIdentityId())
					return types.NewReplayScheduleResult(task.TaskId(), err, nil)
				}
			}

		}


		// 选出 关于自己 metaDataId 所在的 dataNode
		var metaDataId string
		for _, dataSupplier := range task.TaskData().GetDataSuppliers() {
			if selfIdentityId == dataSupplier.GetMemberInfo().GetIdentityId() && myPartyId == dataSupplier.GetMemberInfo().GetPartyId() {
				metaDataId = dataSupplier.MetadataId
			}
		}


		// 获取 metaData 所在的dataNode 资源
		dataResourceDiskUsed, err := sche.dataCenter.QueryDataResourceDiskUsed(metaDataId)
		if nil != err {
			log.Errorf("failed query internal data node by metaDataId on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, metaDataId: {%s}",
				task.TaskId(), metaDataId)

			return types.NewReplayScheduleResult(task.TaskId(), err, nil)
		}

		dataNode, err := sche.dataCenter.GetRegisterNode(pb.PrefixTypeDataNode, dataResourceDiskUsed.GetNodeId())
		if nil != err {
			log.Errorf("failed query internal data node by metaDataId on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, metaDataId: {%s}",
				task.TaskId(), metaDataId)

			return types.NewReplayScheduleResult(task.TaskId(), err, nil)
		}

		log.Debugf("Succeed dataSupplier dataNode on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, dataNode: %s",
			task.TaskId(), dataNode.String())

		result = types.NewReplayScheduleResult(task.TaskId(), nil, types.NewPrepareVoteResource(dataNode.Id, dataNode.ExternalIp, dataNode.ExternalPort, myPartyId))

	// 如果 当前参与方为 PowerSupplier  [选出自己的 内部 power 资源, 并锁定, todo 在最后 DoneXxxxWrap 中解锁]
	case apipb.TaskRole_TaskRole_PowerSupplier:

		needSlotCount := sche.resourceMng.GetSlotUnit().CalculateSlotCount(cost.Mem, cost.Bandwidth, cost.Processor)
		jobNode, err := sche.electionComputeNode(needSlotCount)
		if nil != err {
			log.Errorf("Failed to election internal power resource on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, err: {%s}",
				task.TaskId(), err)

			return types.NewReplayScheduleResult(task.TaskId(), err, nil)
		}

		log.Debugf("Succeed powerSupplier jobNode on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, jobNode: %s",
			task.TaskId(), jobNode.String())

		if err := sche.resourceMng.LockLocalResourceWithTask(jobNode.Id, needSlotCount, task); nil != err {
			log.Errorf("Failed to Lock LocalResource {%s} With Task {%s} on SchedulerStarveFIFO.ReplaySchedule(), err: {%s}",
				jobNode.Id, task.TaskId(), err)

			return types.NewReplayScheduleResult(task.TaskId(), err, nil)
		}

		result = types.NewReplayScheduleResult(task.TaskId(), nil, types.NewPrepareVoteResource(jobNode.Id, jobNode.ExternalIp, jobNode.ExternalPort, myPartyId))

	// 如果 当前参与方为 ResultSupplier  [仅仅是选出自己可用的 dataNode]
	case apipb.TaskRole_TaskRole_Receiver:

		dataResourceTables, err := sche.dataCenter.QueryDataResourceTables()
		if nil != err {
			log.Errorf("Failed to election internal data resource on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, err: {%s}",
				task.TaskId(), err)

			return types.NewReplayScheduleResult(task.TaskId(), err, nil)
		}

		log.Debugf("QueryDataResourceTables on replaySchedule by taskRole is the resuler on SchedulerStarveFIFO.ReplaySchedule(), dataResourceTables: %s", utilDataResourceArrString(dataResourceTables))

		resource := dataResourceTables[len(dataResourceTables)-1]
		dataNode, err := sche.dataCenter.GetRegisterNode(pb.PrefixTypeDataNode, resource.GetNodeId())
		if nil != err {
			log.Errorf("Failed to query internal data node resource on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, dataNodeId: {%s}, err: {%s}",
				task.TaskId(), resource.GetNodeId(), err)

			return types.NewReplayScheduleResult(task.TaskId(), err, nil)
		}

		log.Debugf("Succeed resultReceiver dataNode on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, dataNode: %s",
			task.TaskId(), dataNode.String())

		result = types.NewReplayScheduleResult(task.TaskId(), nil, types.NewPrepareVoteResource(dataNode.Id, dataNode.ExternalIp, dataNode.ExternalPort, myPartyId))

	default:
		result = types.NewReplayScheduleResult(task.TaskId(), fmt.Errorf("task role: {%s} can not replay schedule task", myTaskRole.String()), nil)
	}

	return result
}