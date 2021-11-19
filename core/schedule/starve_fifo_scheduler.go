package schedule

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/auth"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sync"
)

const (
	ReschedMaxCount = 8
	StarveTerm      = 3
)

var (
	ErrEnoughResourceOrgCountLessCalculateCount = fmt.Errorf("the enough resource org count is less calculate count")
	ErrEnoughInternalResourceCount              = fmt.Errorf("has not enough internal resource count")
	ErrRescheduleLargeThreshold                 = errors.New("The reschedule count of task bullet is large than max threshold")
)

type SchedulerStarveFIFO struct {
	internalNodeSet *grpclient.InternalResourceClientSet
	resourceMng     *resource.Manager
	authMng         *auth.AuthorityManager
	// the local task into this queue, first
	queue *types.TaskBullets
	// the very very starve local task by priority
	starveQueue *types.TaskBullets
	// the scheduling task, it is ejected from the queue (taskId -> taskBullet)
	schedulings   map[string]*types.TaskBullet
	scheduleMutex sync.Mutex

	//quit            chan struct{}
	eventEngine *evengine.EventEngine
	//dataCenter      iface.ForResourceDB
	err error

}

func NewSchedulerStarveFIFO(
	internalNodeSet *grpclient.InternalResourceClientSet,
	eventEngine *evengine.EventEngine,
	resourceMng *resource.Manager,
	authMng *auth.AuthorityManager,
) *SchedulerStarveFIFO {

	return &SchedulerStarveFIFO{
		internalNodeSet: internalNodeSet,
		resourceMng:     resourceMng,
		authMng:         authMng,
		queue:           new(types.TaskBullets),
		starveQueue:     new(types.TaskBullets),
		schedulings:     make(map[string]*types.TaskBullet),
		eventEngine:     eventEngine,
		//quit:                 make(chan struct{}),
	}
}

func (sche *SchedulerStarveFIFO) recoveryQueueSchedulings() {

	prefix := rawdb.GetTaskBulletKeyPrefix()
	if err := sche.resourceMng.GetDB().ForEachTaskBullets(func(key, value []byte) error {
		if len(key) != 0 && len(value) != 0 {
			var result types.TaskBullet
			taskId := string(key[len(prefix):])

			if err := rlp.DecodeBytes(value, &result); nil != err {
				return fmt.Errorf("decode taskBullet failed, %s, taskId: {%s}", err, taskId)
			}
			sche.schedulings[taskId] = &result
			if result.Starve == true {
				sche.starveQueue.Push(result)
			} else {
				sche.queue.Push(result)
			}
		}
		return nil
	}); nil != err {
		log.WithError(err).Fatalf("recover taskBullet queue/starveQueue/schedulings failed")
		return
	}

	// todo need fill `pending` task into queue ?

}
func (sche *SchedulerStarveFIFO) Start() error {
	//go sche.loop()
	sche.recoveryQueueSchedulings()
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
	return sche.pushTaskBullet(bullet)
}

func (sche *SchedulerStarveFIFO) RepushTask(task *types.Task) error {
	bullet, ok := sche.schedulings[task.GetTaskId()]
	if !ok {
		return nil
	}
	if bullet.IsOverlowReschedThreshold(ReschedMaxCount) {
		sche.removeTaskBullet(bullet.GetTaskId())
		return ErrRescheduleLargeThreshold
	}
	return sche.repushTaskBullet(bullet)
}

func (sche *SchedulerStarveFIFO) RemoveTask(taskId string) error {
	return sche.removeTaskBullet(taskId)
}

func (sche *SchedulerStarveFIFO) TrySchedule() (needConsensusTask *types.NeedConsensusTask, err error) {

	sche.increaseTotalTaskTerm()
	bullet := sche.popTaskBullet()

	if nil == bullet {
		return nil, nil
	}
	bullet.IncreaseResched()

	defer func() {
		if er := recover(); nil != er {
			if bullet.IsOverlowReschedThreshold(ReschedMaxCount) {
				needConsensusTask, err =  nil, ErrRescheduleLargeThreshold
			} else {
				err =  fmt.Errorf("%s", er)
			}
		}
	}()

	task, err := sche.resourceMng.GetDB().QueryLocalTask(bullet.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("Failed to QueryLocalTask on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", bullet.GetTaskId())
		sche.removeTaskBullet(bullet.GetTaskId())
		panic(err.Error())
	}

	needConsensusTask = types.NewNeedConsensusTask(task)

	// query the powerPartyIds of this task
	powerPartyIds, err := sche.resourceMng.GetDB().QueryTaskPowerPartyIds(task.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query powerPartyIds of task on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
		panic(err.Error())
	}

	cost := &ctypes.TaskOperationCost{
		Mem:       task.GetTaskData().GetOperationCost().GetMemory(),
		Processor: task.GetTaskData().GetOperationCost().GetProcessor(),
		Bandwidth: task.GetTaskData().GetOperationCost().GetBandwidth(),
		Duration:  task.GetTaskData().GetOperationCost().GetDuration(),
	}

	log.Debugf("Call SchedulerStarveFIFO.TrySchedule() start, taskId: {%s}, partyId: {%s}, taskCost: {%s}",
		task.GetTaskData().TaskId, task.GetTaskData().PartyId, cost.String())

	// election other org's power resources
	powers, err := sche.electionComputeOrg(powerPartyIds, nil, cost)
	if nil != err {
		log.WithError(err).Errorf("Failed to election powers org on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
		panic(err.Error())
	}

	log.Debugf("Succeed to election powers org on SchedulerStarveFIFO.TrySchedule(), taskId {%s}, powers: %s", task.GetTaskId(), utilOrgPowerArrString(powers))

	// Set elected powers into task info, and restore into local db.
	task = types.ConvertTaskMsgToTaskWithPowers(task, powers)
	// restore task by power
	if err := sche.resourceMng.GetDB().StoreLocalTask(task); nil != err {
		log.WithError(err).Errorf("Failed tp update local task by election powers on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
		panic(err.Error())
	}
	return needConsensusTask, nil
}
func (sche *SchedulerStarveFIFO) ReplaySchedule(localPartyId string, localTaskRole apicommonpb.TaskRole, task *types.Task) *types.ReplayScheduleResult {

	cost := &ctypes.TaskOperationCost{
		Mem:       task.GetTaskData().GetOperationCost().GetMemory(),
		Processor: task.GetTaskData().GetOperationCost().GetProcessor(),
		Bandwidth: task.GetTaskData().GetOperationCost().GetBandwidth(),
		Duration:  task.GetTaskData().GetOperationCost().GetDuration(),
	}

	log.Debugf("Call SchedulerStarveFIFO.ReplaySchedule() start, taskId: {%s}, role: {%s}, partyId: {%s}, taskCost: {%s}",
		task.GetTaskId(), localTaskRole.String(), localPartyId, cost.String())

	selfIdentityId, err := sche.resourceMng.GetDB().QueryIdentityId()
	if nil != err {
		log.WithError(err).Errorf("Failed to query self identityInfo on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}",
			task.GetTaskId(), localTaskRole.String(), localPartyId)
		return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
	}

	var result *types.ReplayScheduleResult

	switch localTaskRole {

	// 如果 当前参与方为 DataSupplier   [重新 演算 选 powers]
	case apicommonpb.TaskRole_TaskRole_DataSupplier:

		powerPartyIds := make([]string, len(task.GetTaskData().GetPowerSuppliers()))
		for i, power := range task.GetTaskData().GetPowerSuppliers() {
			powerPartyIds[i] = power.GetOrganization().GetPartyId()
		}
		// mock election power orgs
		powers, err := sche.electionComputeOrg(powerPartyIds, nil, cost)
		if nil != err {
			log.WithError(err).Errorf("Failed to election powers org when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}",
				task.GetTaskId(), localTaskRole.String(), localPartyId)

			return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
		}

		log.Debugf("Succeed to election powers org when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, powers: %s",
			task.GetTaskId(), localTaskRole.String(), localPartyId, utilOrgPowerArrString(powers))

		// compare powerSuppliers of task And powerSuppliers of election
		if len(powers) != len(task.GetTaskData().GetPowerSuppliers()) {
			log.Errorf("reschedule powers len and task powers len is not match when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, reschedule power len: {%d}, task powers len: {%d}",
				task.GetTaskId(), localTaskRole.String(), localPartyId, len(powers), len(task.GetTaskData().GetPowerSuppliers()))

			return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
		}

		tmp := make(map[string]struct{}, len(powers))

		for _, power := range powers {
			tmp[power.GetOrganization().GetIdentityId()] = struct{}{}
		}
		for _, power := range task.GetTaskData().GetPowerSuppliers() {
			if _, ok := tmp[power.GetOrganization().GetIdentityId()]; !ok {
				log.Errorf("task power identityId not found with reschedule powers when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, task power identityId: {%s}",
					task.GetTaskId(), localTaskRole.String(), localPartyId, power.GetOrganization().GetIdentityId())
				return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
			}
		}

		// Find metadataId of current identyt of task with current partyId.
		var metadataId string
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			if selfIdentityId == dataSupplier.GetOrganization().GetIdentityId() && localPartyId == dataSupplier.GetOrganization().GetPartyId() {
				metadataId = dataSupplier.MetadataId
			}
		}

		// verify user metadataAuth about msg
		if err = sche.verifyUserMetadataAuthOnTask(task.GetTaskData().GetUserType(), task.GetTaskData().GetUser(), metadataId); nil != err {
			log.WithError(err).Errorf("failed verify user metadataAuth when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, userType: {%s}, user: {%s}, metadataId: {%s}",
				task.GetTaskId(), localTaskRole.String(), localPartyId, task.GetTaskData().GetUserType(), task.GetTaskData().GetUser(), metadataId)
			return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("verify user metadataAuth failed"), nil)
		}

		log.Debugf("Succeed verify user metadataAuth when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, userType: {%s}, user: {%s}, metadataId: {%s}",
			task.GetTaskId(), localTaskRole.String(), localPartyId, task.GetTaskData().GetUserType(), task.GetTaskData().GetUser(), metadataId)

		// Select the datanode where your metadata ID is located.
		var dataNodeId string
		// check the metadata whether internal metadata
		internalMetadataFlag, err := sche.resourceMng.GetDB().IsInternalMetadataByDataId(metadataId)
		if nil != err {
			log.Errorf("failed check metadata whether internal metadata by metaDataId when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
				task.GetTaskId(), localTaskRole.String(), localPartyId, metadataId)

			return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
		}
		if internalMetadataFlag {

			internalMetadata, err := sche.resourceMng.GetDB().QueryInternalMetadataByDataId(metadataId)
			if nil != err {
				log.WithError(err).Errorf("failed query internal metadataInfo by metaDataId when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
					task.GetTaskId(), localTaskRole.String(), localPartyId, metadataId)
				return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
			}
			dataResourceFileUpload, err := sche.resourceMng.GetDB().QueryDataResourceFileUpload(internalMetadata.GetData().GetOriginId())
			if nil != err {
				log.WithError(err).Errorf("Failed query dataResourceFileUpload on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
					task.GetTaskId(), localTaskRole.String(), localPartyId, metadataId)
				return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
			}
			dataNodeId = dataResourceFileUpload.GetNodeId()
		} else {

			dataResourceDiskUsed, err := sche.resourceMng.GetDB().QueryDataResourceDiskUsed(metadataId)
			if nil != err {
				log.WithError(err).Errorf("failed query dataResourceDiskUsed by metaDataId when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, selfTaskRole: {%s}, selfPartyId: {%s}, metadataId: {%s}",
					task.GetTaskId(), localTaskRole.String(), localPartyId, metadataId)

				return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
			}
			dataNodeId = dataResourceDiskUsed.GetNodeId()
		}

		dataNode, err := sche.resourceMng.GetDB().QueryRegisterNode(pb.PrefixTypeDataNode, dataNodeId)
		if nil != err {
			log.WithError(err).Errorf("failed query internal data resource by metaDataId when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
				task.GetTaskId(), localTaskRole.String(), localPartyId, metadataId)

			return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
		}

		log.Debugf("Succeed election internal data resource when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, dataNode: %s",
			task.GetTaskId(), localTaskRole.String(), localPartyId, dataNode.String())

		result = types.NewReplayScheduleResult(task.GetTaskId(), nil, types.NewPrepareVoteResource(dataNode.Id, dataNode.ExternalIp, dataNode.ExternalPort, localPartyId))

	// 如果 当前参与方为 PowerSupplier  [选出自己的 内部 power 资源, 并锁定, todo 在最后 DoneXxxxWrap 中解锁]
	case apicommonpb.TaskRole_TaskRole_PowerSupplier:

		needSlotCount := sche.resourceMng.GetSlotUnit().CalculateSlotCount(cost.Mem, cost.Bandwidth, cost.Processor)
		jobNode, err := sche.electionComputeNode(needSlotCount)
		if nil != err {
			log.WithError(err).Errorf("Failed to election internal power resource when role is powerSupplier on SchedulerStarveFIFO.ReplaySchedule(),taskId: {%s}, role: {%s}, partyId: {%s}",
				task.GetTaskId(), localTaskRole.String(), localPartyId)

			return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
		}

		log.Debugf("Succeed election internal power resource when role is powerSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, jobNode: %s",
			task.GetTaskId(), localTaskRole.String(), localPartyId, jobNode.String())

		if err := sche.resourceMng.LockLocalResourceWithTask(localPartyId, jobNode.Id, needSlotCount, task); nil != err {
			log.WithError(err).Errorf("Failed to Lock LocalResource when role is powerSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
				task.GetTaskId(), localTaskRole.String(), localPartyId, jobNode.Id)

			return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
		}

		result = types.NewReplayScheduleResult(task.GetTaskId(), nil, types.NewPrepareVoteResource(jobNode.Id, jobNode.ExternalIp, jobNode.ExternalPort, localPartyId))

	// 如果 当前参与方为 ResultSupplier  [仅仅是选出自己可用的 dataNode]
	case apicommonpb.TaskRole_TaskRole_Receiver:

		dataResourceTables, err := sche.resourceMng.GetDB().QueryDataResourceTables()
		if nil != err {
			log.WithError(err).Errorf("Failed to election internal data resource when role is receiver on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}",
				task.GetTaskId(), localTaskRole.String(), localPartyId)

			return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
		}

		log.Debugf("QueryDataResourceTables when role is receiver on replaySchedule by taskRole is the resuler on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, dataResourceTables: %s",
			task.GetTaskId(), localTaskRole.String(), localPartyId, utilDataResourceArrString(dataResourceTables))

		resource := dataResourceTables[len(dataResourceTables)-1]
		dataNode, err := sche.resourceMng.GetDB().QueryRegisterNode(pb.PrefixTypeDataNode, resource.GetNodeId())
		if nil != err {
			log.WithError(err).Errorf("Failed to query internal data resource when role is receiver on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
				task.GetTaskId(), localTaskRole.String(), localPartyId, resource.GetNodeId())

			return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
		}

		log.Debugf("Succeed election internal data resource when role is receiver on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, dataNode: %s",
			task.GetTaskId(), localTaskRole.String(), localPartyId, dataNode.String())

		result = types.NewReplayScheduleResult(task.GetTaskId(), nil, types.NewPrepareVoteResource(dataNode.Id, dataNode.ExternalIp, dataNode.ExternalPort, localPartyId))

	default:
		result = types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("task role: {%s} can not replay schedule task", localTaskRole.String()), nil)
	}

	return result
}
