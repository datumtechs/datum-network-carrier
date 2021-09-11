package schedule

import (
	"fmt"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sync"
)

const (
	ReschedMaxCount             = 8
	StarveTerm                  = 3

	electionOrgCondition        = 10000
	electionLocalSeed           = 2
	//taskComputeOrgCount         = 3
)

var (
	ErrEnoughResourceOrgCountLessCalculateCount = fmt.Errorf("the enough resource org count is less calculate count")
	ErrEnoughInternalResourceCount              = fmt.Errorf("has not enough internal resource count")
	ErrRescheduleLargeThreshold = errors.New("The reschedule count of task bullet is large than max threshold")
)

type SchedulerStarveFIFO struct {
	internalNodeSet *grpclient.InternalResourceClientSet
	resourceMng     *resource.Manager
	// the local task into this queue, first
	queue *types.TaskBullets
	// the very very starve local task by priority
	starveQueue *types.TaskBullets
	// the scheduling task, it is ejected from the queue (taskId -> taskBullet)
	schedulings   map[string]*types.TaskBullet
	queueMutex  sync.Mutex

	//quit            chan struct{}
	eventEngine     *evengine.EventEngine
	//dataCenter      iface.ForResourceDB
	err             error

	// TODO 有些缓存需要持久化
}

func NewSchedulerStarveFIFO(
	internalNodeSet *grpclient.InternalResourceClientSet,
	eventEngine *evengine.EventEngine,
	mng *resource.Manager,
) *SchedulerStarveFIFO {

	return &SchedulerStarveFIFO{
		internalNodeSet:      internalNodeSet,
		resourceMng:          mng,
		queue:                new(types.TaskBullets),
		starveQueue:          new(types.TaskBullets),
		schedulings: 		  make(map[string]*types.TaskBullet),
		eventEngine:          eventEngine,
		//quit:                 make(chan struct{}),
	}
}

func (sche *SchedulerStarveFIFO) Start() error {
	//go sche.loop()
	log.Info("Started SchedulerStarveFIFO ...")
	return nil
}
func (sche *SchedulerStarveFIFO) Stop() error {
	//close(sche.quit)
	return nil
}
func (sche *SchedulerStarveFIFO) Error() error { return sche.err }
func (sche *SchedulerStarveFIFO) Name() string { return "SchedulerStarveFIFO" }

func (sche *SchedulerStarveFIFO) AddTask(task *types.Task) error {
	bullet := types.NewTaskBullet(task.GetTaskId())
	// todo 这里需要做一次 持久化
	return sche.pushTaskBullet(bullet)
}

func (sche *SchedulerStarveFIFO) RemoveTask(taskId string) error {
	return sche.removeTaskBullet(taskId)
}

func (sche *SchedulerStarveFIFO) TrySchedule() (*types.NeedConsensusTask, error) {

	sche.increaseTaskTerm()
	bullet := sche.popTaskBullet()

	if nil == bullet {
		return nil, nil
	}

	task, err := sche.resourceMng.GetDB().GetLocalTask(bullet.TaskId)
	if nil != err {
		log.Errorf("Failed to QueryLocalTask on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}, err: {%s}", bullet.TaskId, err)

		if e := sche.repushTaskBullet(bullet.TaskId); nil != e {
			err = e
		}
		return types.NewNeedConsensusTask(types.NewTask(&libTypes.TaskPB{TaskId: bullet.TaskId})), err
	}

	cost := &ctypes.TaskOperationCost{
		Mem:       task.GetTaskData().GetOperationCost().GetMemory(),
		Processor: task.GetTaskData().GetOperationCost().GetProcessor(),
		Bandwidth: task.GetTaskData().GetOperationCost().GetBandwidth(),
		Duration:  task.GetTaskData().GetOperationCost().GetDuration(),
	}

	log.Debugf("Call SchedulerStarveFIFO.TrySchedule() start, taskId: {%s}, partyId: {%s}, taskCost: {%s}",
		task.GetTaskData().TaskId, task.GetTaskData().PartyId, cost.String())

	// 获取 powerPartyIds 标签 TODO

	// 【选出 其他组织的算力】
	powers, err := sche.electionConputeOrg(nil, nil, cost)
	if nil != err {
		log.Errorf("Failed to election powers org on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}, err: {%s}", task.GetTaskId(), err)
		//repushFn(task)
		return types.NewNeedConsensusTask(task), err
	}

	log.Debugf("Succeed to election powers org on SchedulerStarveFIFO.TrySchedule(), taskId {%s}, powers: %s", task.GetTaskId(), utilOrgPowerArrString(powers))

	// Set elected powers into task info, and restore into local db.
	task = types.ConvertTaskMsgToTaskWithPowers(task, powers)
	// restore task by power
	if err := sche.resourceMng.GetDB().StoreLocalTask(task); nil != err {
		log.Errorf("Failed tp update local task by election powers on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}, err: {%s}", task.GetTaskId(), err)

		//repushFn(task)
		return types.NewNeedConsensusTask(task), err
	}
	return types.NewNeedConsensusTask(task), nil
}
func (sche *SchedulerStarveFIFO) ReplaySchedule(localPartyId string, localTaskRole apipb.TaskRole, task *types.Task) *types.ReplayScheduleResult {

	cost := &ctypes.TaskOperationCost{
		Mem:       task.GetTaskData().GetOperationCost().GetMemory(),
		Processor: task.GetTaskData().GetOperationCost().GetProcessor(),
		Bandwidth: task.GetTaskData().GetOperationCost().GetBandwidth(),
		Duration:  task.GetTaskData().GetOperationCost().GetDuration(),
	}

	log.Debugf("Call SchedulerStarveFIFO.ReplaySchedule() start, taskId: {%s}, localTaskRole: {%s}, myPartyId: {%s}, taskCost: {%s}",
		task.GetTaskId(), localTaskRole.String(), localPartyId, cost.String())

	selfIdentityId, err := sche.resourceMng.GetDB().GetIdentityId()
	if nil != err {
		log.Errorf("Failed to query self identityInfo on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, err: {%s}", task.GetTaskId(), err)
		return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
	}

	var result *types.ReplayScheduleResult

	switch localTaskRole {

	// 如果 当前参与方为 DataSupplier   [重新 演算 选 powers]
	case apipb.TaskRole_TaskRole_DataSupplier:

		var isSender bool
		if localPartyId == task.GetTaskSender().GetPartyId() && selfIdentityId == task.GetTaskSender().GetIdentityId() {
			isSender = true
		}

		if !isSender {

			powerPartyIds := make([]string, len(task.GetTaskData().GetPowerSuppliers()))
			for i, power := range task.GetTaskData().GetPowerSuppliers() {
				powerPartyIds[i] = power.GetOrganization().GetPartyId()
			}
			// mock election power orgs
			powers, err := sche.electionConputeOrg(powerPartyIds, nil, cost)
			if nil != err {
				log.Errorf("Failed to election powers org on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, err: {%s}",
					task.GetTaskId(), err)

				return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
			}

			log.Debugf("Succeed to election powers org on SchedulerStarveFIFO.ReplaySchedule(), taskId {%s}, powers: %s",
				task.GetTaskId(), utilOrgPowerArrString(powers))

			// compare powerSuppliers of task And powerSuppliers of election
			if len(powers) != len(task.GetTaskData().GetPowerSuppliers()) {
				log.Errorf("reschedule powers len and task powers len is not match on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, reschedule power len: {%d}, task powers len: {%d}",
					task.GetTaskId(), len(powers), len(task.GetTaskData().GetPowerSuppliers()))

				return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
			}

			tmp := make(map[string]struct{}, len(powers))

			for _, power := range powers {
				tmp[power.GetOrganization().GetIdentityId()] = struct{}{}
			}
			for _, power := range task.GetTaskData().GetPowerSuppliers() {
				if _, ok := tmp[power.GetOrganization().GetIdentityId()]; !ok {
					log.Errorf("task power identityId not found on reschedule powers on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, task power identityId: {%s}",
						task.GetTaskId(), power.GetOrganization().GetIdentityId())
					return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
				}
			}

		}

		// 选出 关于自己 metaDataId 所在的 dataNode
		var metaDataId string
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			if selfIdentityId == dataSupplier.GetOrganization().GetIdentityId() && localPartyId == dataSupplier.GetOrganization().GetPartyId() {
				metaDataId = dataSupplier.MetadataId
			}
		}

		// 获取 metaData 所在的dataNode 资源
		dataResourceDiskUsed, err := sche.resourceMng.GetDB().QueryDataResourceDiskUsed(metaDataId)
		if nil != err {
			log.Errorf("failed query internal data node by metaDataId on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, metaDataId: {%s}",
				task.GetTaskId(), metaDataId)

			return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
		}

		dataNode, err := sche.resourceMng.GetDB().GetRegisterNode(pb.PrefixTypeDataNode, dataResourceDiskUsed.GetNodeId())
		if nil != err {
			log.Errorf("failed query internal data node by metaDataId on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, metaDataId: {%s}",
				task.GetTaskId(), metaDataId)

			return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
		}

		log.Debugf("Succeed dataSupplier dataNode on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, dataNode: %s",
			task.GetTaskId(), dataNode.String())

		result = types.NewReplayScheduleResult(task.GetTaskId(), nil, types.NewPrepareVoteResource(dataNode.Id, dataNode.ExternalIp, dataNode.ExternalPort, localPartyId))

	// 如果 当前参与方为 PowerSupplier  [选出自己的 内部 power 资源, 并锁定, todo 在最后 DoneXxxxWrap 中解锁]
	case apipb.TaskRole_TaskRole_PowerSupplier:

		needSlotCount := sche.resourceMng.GetSlotUnit().CalculateSlotCount(cost.Mem, cost.Bandwidth, cost.Processor)
		jobNode, err := sche.electionComputeNode(needSlotCount)
		if nil != err {
			log.Errorf("Failed to election internal power resource on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, err: {%s}",
				task.GetTaskId(), err)

			return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
		}

		log.Debugf("Succeed powerSupplier jobNode on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, jobNode: %s",
			task.GetTaskId(), jobNode.String())

		if err := sche.resourceMng.LockLocalResourceWithTask(jobNode.Id, needSlotCount, task); nil != err {
			log.Errorf("Failed to Lock LocalResource {%s} With GetTask {%s} on SchedulerStarveFIFO.ReplaySchedule(), err: {%s}",
				jobNode.Id, task.GetTaskId(), err)

			return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
		}

		result = types.NewReplayScheduleResult(task.GetTaskId(), nil, types.NewPrepareVoteResource(jobNode.Id, jobNode.ExternalIp, jobNode.ExternalPort, localPartyId))

	// 如果 当前参与方为 ResultSupplier  [仅仅是选出自己可用的 dataNode]
	case apipb.TaskRole_TaskRole_Receiver:

		dataResourceTables, err := sche.resourceMng.GetDB().QueryDataResourceTables()
		if nil != err {
			log.Errorf("Failed to election internal data resource on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, err: {%s}",
				task.GetTaskId(), err)

			return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
		}

		log.Debugf("QueryDataResourceTables on replaySchedule by taskRole is the resuler on SchedulerStarveFIFO.ReplaySchedule(), dataResourceTables: %s", utilDataResourceArrString(dataResourceTables))

		resource := dataResourceTables[len(dataResourceTables)-1]
		dataNode, err := sche.resourceMng.GetDB().GetRegisterNode(pb.PrefixTypeDataNode, resource.GetNodeId())
		if nil != err {
			log.Errorf("Failed to query internal data node resource on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, dataNodeId: {%s}, err: {%s}",
				task.GetTaskId(), resource.GetNodeId(), err)

			return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
		}

		log.Debugf("Succeed resultReceiver dataNode on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, dataNode: %s",
			task.GetTaskId(), dataNode.String())

		result = types.NewReplayScheduleResult(task.GetTaskId(), nil, types.NewPrepareVoteResource(dataNode.Id, dataNode.ExternalIp, dataNode.ExternalPort, localPartyId))

	default:
		result = types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("task role: {%s} can not replay schedule task", localTaskRole.String()), nil)
	}

	return result
}
