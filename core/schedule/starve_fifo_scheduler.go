package schedule

import (
	"container/heap"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/auth"
	"github.com/RosettaFlow/Carrier-Go/common/bytesutil"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/core/election"
	"github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/policy"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/ethereum/go-ethereum/rlp"
	log "github.com/sirupsen/logrus"
	"sync"
)

const (
	ReschedMaxCount = 8
	StarveTerm      = 3
)

var (
	ErrRescheduleLargeThreshold                 = fmt.Errorf("the reschedule count of task bullet is large than max threshold")
	ErrAbandonTaskWithNotFoundTask              = fmt.Errorf("the task must be abandoned with not found local task")
	ErrAbandonTaskWithNotFoundPowerPartyIds     = fmt.Errorf("the task must be abandoned with not found power partyIds of local task")
)

type SchedulerStarveFIFO struct {
	elector         *election.VrfElector
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
	elector *election.VrfElector,
	eventEngine *evengine.EventEngine,
	resourceMng *resource.Manager,
	authMng *auth.AuthorityManager,
) *SchedulerStarveFIFO {

	return &SchedulerStarveFIFO{
		elector: elector,
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
			var bullet types.TaskBullet
			taskId := string(key[len(prefix):])

			if err := rlp.DecodeBytes(value, &bullet); nil != err {
				return fmt.Errorf("decode taskBullet failed, %s, taskId: {%s}", err, taskId)
			}
			log.Debugf("Recovery taskBullet, taskId: {%s}, inQueueFlag: {%v}", taskId, (&bullet).GetInQueueFlag())
			sche.schedulings[taskId] = &bullet
			

			if !(&bullet).GetInQueueFlag() {
				if (&bullet).IsOverlowReschedThreshold(ReschedMaxCount) {
					delete(sche.schedulings, taskId)
					sche.resourceMng.GetDB().RemoveTaskBullet(taskId)
				} else {
					(&bullet).InQueueFlag = true
				}
			}
			// push into queue/starve queue
			if (&bullet).IsStarve() {
				heap.Push(sche.starveQueue, &bullet)
			} else {
				heap.Push(sche.queue, &bullet)
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

	sche.scheduleMutex.Lock()
	defer sche.scheduleMutex.Unlock()

	bullet, ok := sche.schedulings[task.GetTaskId()]
	if !ok {
		return nil
	}

	if bullet.IsOverlowReschedThreshold(ReschedMaxCount) {
		return ErrRescheduleLargeThreshold
	}

	return sche.repushTaskBullet(bullet)
}

func (sche *SchedulerStarveFIFO) RemoveTask(taskId string) error {
	return sche.removeTaskBullet(taskId)
}

func (sche *SchedulerStarveFIFO) TrySchedule() (resTask *types.NeedConsensusTask, taskId string, err error) {

	sche.increaseTotalTaskTerm()
	bullet := sche.popTaskBullet()
	if nil == bullet {
		return nil, "", nil
	}
	bullet.IncreaseResched()

	task, err := sche.resourceMng.GetDB().QueryLocalTask(bullet.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task, must abandon it on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", bullet.GetTaskId())
		return nil, bullet.GetTaskId(), ErrAbandonTaskWithNotFoundTask
	}

	// query the powerPartyIds of this task
	powerPartyIds, err := sche.resourceMng.GetDB().QueryTaskPowerPartyIds(task.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query power partyIds of local task, must abandon it on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
		return types.NewNeedConsensusTask(task, nil, nil, 0), bullet.GetTaskId(), ErrAbandonTaskWithNotFoundPowerPartyIds
	}

	cost := &ctypes.TaskOperationCost{
		Mem:       task.GetTaskData().GetOperationCost().GetMemory(),
		Processor: task.GetTaskData().GetOperationCost().GetProcessor(),
		Bandwidth: task.GetTaskData().GetOperationCost().GetBandwidth(),
		Duration:  task.GetTaskData().GetOperationCost().GetDuration(),
	}

	log.Debugf("Call SchedulerStarveFIFO.TrySchedule() start, taskId: {%s}, partyId: {%s}, taskCost: {%s}",
		task.GetTaskData().GetTaskId(), task.GetTaskData().GetSender().GetPartyId(), cost.String())

	now := timeutils.UnixMsecUint64()
	vrfInput := append([]byte(bullet.GetTaskId()), bytesutil.Uint64ToBytes(now)...)  // input == taskId + nowtime
	// election other org's power resources
	powerOrgs, powerResources, nonce, weights, err := sche.elector.ElectionOrganization (bullet.GetTaskId(), powerPartyIds,
		nil, cost.GetMem(), cost.GetBandwidth(), 0, cost.GetProcessor(), vrfInput)
	if nil != err {
		log.WithError(err).Errorf("Failed to election powers org on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
		return types.NewNeedConsensusTask(task, nonce, weights, now), bullet.GetTaskId(), fmt.Errorf("election powerOrg failed, %s", err)
	}

	log.Debugf("Succeed to election powers org on SchedulerStarveFIFO.TrySchedule(), taskId {%s}, powerOrgs: %s, powerResources: %s",
		task.GetTaskId(), types.UtilOrgPowerArrString(powerOrgs), types.UtilOrgPowerResourceArrString(powerResources))

	// Set elected powers into task info, and restore into local db.
	task.SetPowerSuppliers(powerOrgs)
	task.SetPowerResources(powerResources)
	// restore task by power
	if err := sche.resourceMng.GetDB().StoreLocalTask(task); nil != err {
		log.WithError(err).Errorf("Failed to update local task by election powers on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
		return types.NewNeedConsensusTask(task, nonce, weights, now), bullet.GetTaskId(), fmt.Errorf("update local task failed, %s", err)
	}
	return types.NewNeedConsensusTask(task, nonce, weights, now), bullet.GetTaskId(), nil
}
func (sche *SchedulerStarveFIFO) ReplaySchedule(localPartyId string, localTaskRole apicommonpb.TaskRole, replayTask *types.NeedReplayScheduleTask) *types.ReplayScheduleResult {

	task := replayTask.GetTask()

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
		return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("query local identity failed, %s", err), nil)
	}

	var result *types.ReplayScheduleResult

	switch localTaskRole {

	case apicommonpb.TaskRole_TaskRole_DataSupplier:

		powerPartyIds := make([]string, len(task.GetTaskData().GetPowerSuppliers()))
		for i, power := range task.GetTaskData().GetPowerSuppliers() {
			powerPartyIds[i] = power.GetPartyId()
		}
		vrfInput := append([]byte(task.GetTaskId()), bytesutil.Uint64ToBytes(replayTask.GetElectionAt())...)  // input == taskId + nowtime
		// verify power orgs of task
		agree, err := sche.elector.VerifyElectionOrganization(task.GetTaskId(), task.GetTaskData().GetPowerSuppliers(),
			task.GetTaskData().GetPowerResourceOptions(), task.GetTaskSender().GetNodeId(), vrfInput, replayTask.GetNonce(), replayTask.GetWeights())
		if nil != err {
			log.WithError(err).Errorf("Failed to verify election powers org when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}",
				task.GetTaskId(), localTaskRole.String(), localPartyId)
			return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("%s when election power organization", err), nil)
		}

		log.Debugf("Succeed to verify election powers org when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, agree: %v",
			task.GetTaskId(), localTaskRole.String(), localPartyId, agree)

		// Find metadataId of current identyt of task with current partyId.
		var metadataId string
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			if selfIdentityId == dataSupplier.GetIdentityId() && localPartyId == dataSupplier.GetPartyId() {
				mId, err := policy.FetchMetedataIdByPartyId(localPartyId, task.GetTaskData().GetDataPolicyType(), task.GetTaskData().GetDataPolicyOption())
				if nil != err {
					log.WithError(err).Errorf("not fetch metadataId from task dataPolicy on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, userType: {%s}, user: {%s}",
						task.GetTaskId(), localTaskRole.String(), localPartyId, task.GetTaskData().GetUserType(), task.GetTaskData().GetUser())
					return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("%s when fetch metadataId from task dataPolicy", err), nil)
				}
				metadataId = mId
				break
			}
		}

		// verify user metadataAuth about msg
		if err = sche.verifyUserMetadataAuthOnTask(task.GetTaskData().GetUserType(), task.GetTaskData().GetUser(), metadataId); nil != err {
			log.WithError(err).Errorf("failed verify user metadataAuth when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, userType: {%s}, user: {%s}, metadataId: {%s}",
				task.GetTaskId(), localTaskRole.String(), localPartyId, task.GetTaskData().GetUserType(), task.GetTaskData().GetUser(), metadataId)
			return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("%s when verify user metadataAuth", err), nil)
		}

		log.Debugf("Succeed verify user metadataAuth when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, userType: {%s}, user: {%s}, metadataId: {%s}",
			task.GetTaskId(), localTaskRole.String(), localPartyId, task.GetTaskData().GetUserType(), task.GetTaskData().GetUser(), metadataId)

		// Select the datanode where your metadata ID is located.
		var dataNodeId string
		// check the metadata whether internal metadata
		internalMetadataFlag, err := sche.resourceMng.GetDB().IsInternalMetadataByDataId(metadataId)
		if nil != err {
			log.WithError(err).Errorf("failed check metadata whether internal metadata by metaDataId when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
				task.GetTaskId(), localTaskRole.String(), localPartyId, metadataId)
			return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("check metadata whether internal failed, %s", err), nil)
		}
		if internalMetadataFlag {
			internalMetadata, err := sche.resourceMng.GetDB().QueryInternalMetadataByDataId(metadataId)
			if nil != err {
				log.WithError(err).Errorf("failed query internal metadataInfo by metaDataId when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
					task.GetTaskId(), localTaskRole.String(), localPartyId, metadataId)
				return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("query internal metadata by metadataId failed, %s", err), nil)
			}

			fileOriginId, err := policy.FetchOriginId(internalMetadata.GetData().GetFileType(), internalMetadata.GetData().GetMetadataOption())
			if nil != err {
				log.WithError(err).Errorf("not fetch internalMetadata originId from task metadataOption on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, userType: {%s}, user: {%s}, metadataId: {%s}",
					task.GetTaskId(), localTaskRole.String(), localPartyId, task.GetTaskData().GetUserType(), task.GetTaskData().GetUser(), internalMetadata.GetData().GetMetadataId())
				return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("%s when fetch metadataId from task dataPolicy", err), nil)
			}
			dataResourceFileUpload, err := sche.resourceMng.GetDB().QueryDataResourceFileUpload(fileOriginId)
			if nil != err {
				log.WithError(err).Errorf("Failed query dataResourceFileUpload on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
					task.GetTaskId(), localTaskRole.String(), localPartyId, metadataId)
				return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("query dataFileUpload by originId failed, %s", err), nil)
			}
			dataNodeId = dataResourceFileUpload.GetNodeId()
		} else {

			dataResourceDiskUsed, err := sche.resourceMng.GetDB().QueryDataResourceDiskUsed(metadataId)
			if nil != err {
				log.WithError(err).Errorf("failed query dataResourceDiskUsed by metaDataId when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
					task.GetTaskId(), localTaskRole.String(), localPartyId, metadataId)
				return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("query dataResourceDiskUsed by metadataId failed, %s", err), nil)
			}
			dataNodeId = dataResourceDiskUsed.GetNodeId()
		}

		dataNode, err := sche.resourceMng.GetDB().QueryRegisterNode(pb.PrefixTypeDataNode, dataNodeId)
		if nil != err {
			log.WithError(err).Errorf("failed query internal data resource by metaDataId when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
				task.GetTaskId(), localTaskRole.String(), localPartyId, metadataId)
			return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("query internal dataNode by dataNodeId failed, %s", err), nil)
		}

		log.Debugf("Succeed election internal data resource when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, dataNode: %s",
			task.GetTaskId(), localTaskRole.String(), localPartyId, dataNode.String())

		result = types.NewReplayScheduleResult(task.GetTaskId(), nil, types.NewPrepareVoteResource(dataNode.Id, dataNode.ExternalIp, dataNode.ExternalPort, localPartyId))

	// If the current participant is powersupplier.
	// select your own internal power resource and lock it,
	// and finally click 'publishfinishedtasktodatacenter'
	// or 'sendtaskresultmsgtotasksender' in taskmnager.
	case apicommonpb.TaskRole_TaskRole_PowerSupplier:

		log.Debugf("Succeed CalculateSlotCount when role is powerSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, cost.mem: {%d}, cost.Bandwidth: {%d}, cost.Processor: {%d}",
			task.GetTaskId(), localTaskRole.String(), localPartyId, cost.Mem, cost.Bandwidth, cost.Processor)

		jobNode, err := sche.elector.ElectionNode(task.GetTaskId(), cost.GetMem(), cost.GetBandwidth(), 0, cost.GetProcessor(), "")
		if nil != err {
			log.WithError(err).Errorf("Failed to election internal power resource when role is powerSupplier on SchedulerStarveFIFO.ReplaySchedule(),taskId: {%s}, role: {%s}, partyId: {%s}",
				task.GetTaskId(), localTaskRole.String(), localPartyId)
			return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("%s when election internal jobNode", err), nil)
		}

		log.Debugf("Succeed election internal power resource when role is powerSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, jobNode: %s",
			task.GetTaskId(), localTaskRole.String(), localPartyId, jobNode.String())

		if err := sche.resourceMng.LockLocalResourceWithTask(localPartyId, jobNode.Id, cost.Mem, cost.Bandwidth, 0, cost.Processor, task); nil != err {
			log.WithError(err).Errorf("Failed to Lock LocalResource when role is powerSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
				task.GetTaskId(), localTaskRole.String(), localPartyId, jobNode.Id)
			return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("%s when lock internal jobNode resource", err), nil)
		}

		result = types.NewReplayScheduleResult(task.GetTaskId(), nil, types.NewPrepareVoteResource(jobNode.Id, jobNode.ExternalIp, jobNode.ExternalPort, localPartyId))

	// If the current participant is resultsupplier.
	// just select their own available datanodes.
	case apicommonpb.TaskRole_TaskRole_Receiver:

		dataResourceTables, err := sche.resourceMng.GetDB().QueryDataResourceTables()
		if nil != err {
			log.WithError(err).Errorf("Failed to election internal data resource when role is receiver on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}",
				task.GetTaskId(), localTaskRole.String(), localPartyId)
			return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("query all internal dataResourceTables failed, %s", err), nil)
		}

		log.Debugf("QueryDataResourceTables when role is receiver on replaySchedule by taskRole is the resuler on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, dataResourceTables: %s",
			task.GetTaskId(), localTaskRole.String(), localPartyId, types.UtilDataResourceArrString(dataResourceTables))

		resource := dataResourceTables[len(dataResourceTables)-1]
		dataNode, err := sche.resourceMng.GetDB().QueryRegisterNode(pb.PrefixTypeDataNode, resource.GetNodeId())
		if nil != err {
			log.WithError(err).Errorf("Failed to query internal data resource when role is receiver on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
				task.GetTaskId(), localTaskRole.String(), localPartyId, resource.GetNodeId())
			return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("query internal dataNode by dataNodeId failed, %s", err), nil)
		}

		log.Debugf("Succeed election internal data resource when role is receiver on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, dataNode: %s",
			task.GetTaskId(), localTaskRole.String(), localPartyId, dataNode.String())

		result = types.NewReplayScheduleResult(task.GetTaskId(), nil, types.NewPrepareVoteResource(dataNode.Id, dataNode.ExternalIp, dataNode.ExternalPort, localPartyId))

	default:
		result = types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("task role: {%s} can not replay schedule task", localTaskRole.String()), nil)
	}

	return result
}
