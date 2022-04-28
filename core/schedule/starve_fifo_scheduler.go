package schedule

import (
	"container/heap"
	"fmt"
	auth2 "github.com/Metisnetwork/Metis-Carrier/ach/auth"
	"github.com/Metisnetwork/Metis-Carrier/common/bytesutil"
	"github.com/Metisnetwork/Metis-Carrier/common/timeutils"
	ctypes "github.com/Metisnetwork/Metis-Carrier/consensus/twopc/types"
	"github.com/Metisnetwork/Metis-Carrier/core/election"
	"github.com/Metisnetwork/Metis-Carrier/core/evengine"
	"github.com/Metisnetwork/Metis-Carrier/core/rawdb"
	"github.com/Metisnetwork/Metis-Carrier/core/resource"
	pb "github.com/Metisnetwork/Metis-Carrier/lib/api"
	libtypes "github.com/Metisnetwork/Metis-Carrier/lib/types"
	"github.com/Metisnetwork/Metis-Carrier/policy"
	"github.com/Metisnetwork/Metis-Carrier/types"
	"github.com/ethereum/go-ethereum/rlp"
	log "github.com/sirupsen/logrus"
	"sync"
)

const (
	ReschedMaxCount = 8
	StarveTerm      = 3
)

var (
	ErrRescheduleLargeThreshold             = fmt.Errorf("the reschedule count of task bullet is large than max threshold")
	ErrAbandonTaskWithNotFoundTask          = fmt.Errorf("the task must be abandoned when not found local task")
	ErrAbandonTaskWithNotFoundPowerPartyIds = fmt.Errorf("the task must be abandoned when not fetch partyIds from task powerPolicy")
)

type SchedulerStarveFIFO struct {
	elector     *election.VrfElector
	resourceMng *resource.Manager
	authMng     *auth2.AuthorityManager
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
	authMng *auth2.AuthorityManager,
) *SchedulerStarveFIFO {

	return &SchedulerStarveFIFO{
		elector:     elector,
		resourceMng: resourceMng,
		authMng:     authMng,
		queue:       new(types.TaskBullets),
		starveQueue: new(types.TaskBullets),
		schedulings: make(map[string]*types.TaskBullet),
		eventEngine: eventEngine,
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

	log.Debugf("Call SchedulerStarveFIFO.TrySchedule() start, taskId: {%s}, partyId: {%s}",
		task.GetTaskData().GetTaskId(), task.GetTaskData().GetSender().GetPartyId())

	// NOTE: fill the powerSupplier AND powerResourceOption into local task.
	partyIdAndIndexCache := make(map[string]int, 0)
	assignmentSymbolRandomElectionPowerPartyIds := make([]string, 0)
	dataNodeProviderPowerPartyIds := make([]string, 0)

	for i, policyType := range task.GetTaskData().GetPowerPolicyTypes() {

		switch policyType {
		case types.TASK_POWER_POLICY_ASSIGNMENT_SYMBOL_RANDOM_ELECTION_POWER:
			assignmentSymbolRandomElectionPowerPartyIds = append(assignmentSymbolRandomElectionPowerPartyIds, task.GetTaskData().GetPowerPolicyOptions()[i])
		case types.TASK_POWER_POLICY_DATANODE_PROVIDE_POWER:
			dataNodeProviderPowerPartyIds = append(dataNodeProviderPowerPartyIds, task.GetTaskData().GetPowerPolicyOptions()[i])
		// NOTE: unknown powerPolicyType
		default:
			log.WithError(err).Errorf("unknown powerPolicyType of task on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
			return types.NewNeedConsensusTask(task, ""), bullet.GetTaskId(), fmt.Errorf("unknown powerPolicyType of task, %d", policyType)
		}
		// collection partyId index into cache.
		partyIdAndIndexCache[task.GetTaskData().GetPowerPolicyOptions()[i]] = i
	}

	powerOrgs      := make([]*libtypes.TaskOrganization, len(task.GetTaskData().GetPowerPolicyTypes()))
	powerResources := make([]*libtypes.TaskPowerResourceOption, len(task.GetTaskData().GetPowerPolicyTypes()))

	var (
		evidence  string
	)

	// NOTE: many election policys
	if len(assignmentSymbolRandomElectionPowerPartyIds) != 0 {

		evidenceJson, orgs, resources, err := sche.scheduleVrfElectionPower(task.GetTaskId(), &ctypes.TaskOperationCost{
			Mem:       task.GetTaskData().GetOperationCost().GetMemory(),
			Processor: task.GetTaskData().GetOperationCost().GetProcessor(),
			Bandwidth: task.GetTaskData().GetOperationCost().GetBandwidth(),
			Duration:  task.GetTaskData().GetOperationCost().GetDuration(),
		}, assignmentSymbolRandomElectionPowerPartyIds)
		if nil != err {
			log.WithError(err).Errorf("not fetch partyIds from task powerPolicy on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
			return types.NewNeedConsensusTask(task, ""), bullet.GetTaskId(), err
		}
		for i, org := range orgs {
			if index, ok := partyIdAndIndexCache[org.GetPartyId()]; !ok {
				log.Errorf("not fetch partyIds from task powerPolicy on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
				return types.NewNeedConsensusTask(task, ""), bullet.GetTaskId(), fmt.Errorf("invalid partyId")
			} else {
				powerOrgs[index] = org
				powerResources[index] = resources[i]
			}
		}
		evidence = evidenceJson
	}

	if len(dataNodeProviderPowerPartyIds) != 0 {
		orgs, resources, err := sche.scheduleDataNodeProvidePower(task, dataNodeProviderPowerPartyIds)
		if nil != err {
			log.WithError(err).Errorf("not fetch partyIds from task powerPolicy on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
			return types.NewNeedConsensusTask(task, ""), bullet.GetTaskId(), err
		}

		for i, org := range orgs {
			if index, ok := partyIdAndIndexCache[org.GetPartyId()]; !ok {
				log.Errorf("not fetch partyIds from task powerPolicy on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
				return types.NewNeedConsensusTask(task, ""), bullet.GetTaskId(), fmt.Errorf("invalid partyId")
			} else {
				powerOrgs[index] = org
				powerResources[index] = resources[i]
			}
		}
	}

	log.Debugf("Succeed to election powers org on SchedulerStarveFIFO.TrySchedule(), taskId {%s}, powerOrgs: %s, powerResources: %s",
		task.GetTaskId(), types.UtilOrgPowerArrString(powerOrgs), types.UtilOrgPowerResourceArrString(powerResources))

	// Set elected powers into task info, and restore into local db.
	task.SetPowerSuppliers(powerOrgs)
	task.SetPowerResources(powerResources)
	// restore task by power
	if err := sche.resourceMng.GetDB().StoreLocalTask(task); nil != err {
		log.WithError(err).Errorf("Failed to update local task by election powers on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
		return types.NewNeedConsensusTask(task, ""), bullet.GetTaskId(), fmt.Errorf("update local task failed, %s", err)
	}
	return types.NewNeedConsensusTask(task, evidence), bullet.GetTaskId(), nil
}

func (sche *SchedulerStarveFIFO) ReplaySchedule(
	partyId string, taskRole libtypes.TaskRole,
	replayTask *types.NeedReplayScheduleTask) *types.ReplayScheduleResult {

	task := replayTask.GetTask()

	log.Debugf("Call SchedulerStarveFIFO.ReplaySchedule() start, taskId: {%s}, role: {%s}, partyId: {%s}",
		task.GetTaskId(), taskRole.String(), partyId)

	identityId, err := sche.resourceMng.GetDB().QueryIdentityId()
	if nil != err {
		log.WithError(err).Errorf("Failed to query self identityInfo on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}",
			task.GetTaskId(), taskRole.String(), partyId)
		return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("query local identity failed, %s", err), nil)
	}

	var result *types.ReplayScheduleResult

	switch taskRole {

	case libtypes.TaskRole_TaskRole_DataSupplier:

		// NOTE: verify the powerSupplier AND powerResourceOption into remote task.
		assignmentSymbolRandomElectionPowerSuppliers := make([]*libtypes.TaskOrganization, 0)
		assignmentSymbolRandomElectionPowerResources := make([]*libtypes.TaskPowerResourceOption, 0)
		dataNodeProviderPowerSuppliers := make([]*libtypes.TaskOrganization, 0)
		dataNodeProviderPowerResources := make([]*libtypes.TaskPowerResourceOption, 0)
		dataNodeProviderPowerPartyIds := make([]string, 0)
		for i, policyType := range task.GetTaskData().GetPowerPolicyTypes() {

			switch policyType {
			case types.TASK_POWER_POLICY_ASSIGNMENT_SYMBOL_RANDOM_ELECTION_POWER:
				assignmentSymbolRandomElectionPowerSuppliers = append(assignmentSymbolRandomElectionPowerSuppliers, task.GetTaskData().GetPowerSuppliers()[i])
				assignmentSymbolRandomElectionPowerResources = append(assignmentSymbolRandomElectionPowerResources, task.GetTaskData().GetPowerResourceOptions()[i])
			case types.TASK_POWER_POLICY_DATANODE_PROVIDE_POWER:
				dataNodeProviderPowerSuppliers = append(dataNodeProviderPowerSuppliers, task.GetTaskData().GetPowerSuppliers()[i])
				dataNodeProviderPowerResources = append(dataNodeProviderPowerResources, task.GetTaskData().GetPowerResourceOptions()[i])
				dataNodeProviderPowerPartyIds = append(dataNodeProviderPowerPartyIds, task.GetTaskData().GetPowerPolicyOptions()[i])
			// NOTE: unknown powerPolicyType
			default:
				log.WithError(err).Errorf("unknown powerPolicyType of task on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
				return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
			}
		}

		if len(assignmentSymbolRandomElectionPowerSuppliers) != 0 {
			if len(assignmentSymbolRandomElectionPowerSuppliers) != len(assignmentSymbolRandomElectionPowerResources) {
				log.WithError(err).Errorf("unknown powerPolicyType of task on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
				return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
			}
			if err := sche.reScheduleVrfElectionPower(task.GetTaskId(), task.GetTaskSender().GetNodeId(),
				assignmentSymbolRandomElectionPowerSuppliers, assignmentSymbolRandomElectionPowerResources,replayTask.GetEvidence()); nil != err {
				return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
			}
		}

		if len(dataNodeProviderPowerSuppliers) != 0 {
			if len(dataNodeProviderPowerSuppliers) != len(dataNodeProviderPowerResources) {
				log.WithError(err).Errorf("unknown powerPolicyType of task on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
				return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
			}
			if err := sche.reScheduleDataNodeProvidePower(task.GetTaskId(), task.GetTaskData().GetDataSuppliers(), dataNodeProviderPowerSuppliers, dataNodeProviderPowerPartyIds); nil != err {
				return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
			}
		}

		log.Debugf("Succeed to verify election powers org when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, agree: %v",
			task.GetTaskId(), taskRole.String(), partyId)

		// Find metadataId of current identyt of task with current partyId.
		var metadataId string
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			if identityId == dataSupplier.GetIdentityId() && partyId == dataSupplier.GetPartyId() {
				mId, err := policy.FetchMetedataIdByPartyId(partyId, task.GetTaskData().GetDataPolicyTypes(), task.GetTaskData().GetDataPolicyOptions())
				if nil != err {
					log.WithError(err).Errorf("not fetch metadataId from task dataPolicy on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, userType: {%s}, user: {%s}",
						task.GetTaskId(), taskRole.String(), partyId, task.GetTaskData().GetUserType(), task.GetTaskData().GetUser())
					return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("%s when fetch metadataId from task dataPolicy", err), nil)
				}
				metadataId = mId
				break
			}
		}

		// verify user metadataAuth about msg
		if err = sche.verifyUserMetadataAuthOnTask(task.GetTaskData().GetUserType(), task.GetTaskData().GetUser(), metadataId); nil != err {
			log.WithError(err).Errorf("failed verify user metadataAuth when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, userType: {%s}, user: {%s}, metadataId: {%s}",
				task.GetTaskId(), taskRole.String(), partyId, task.GetTaskData().GetUserType(), task.GetTaskData().GetUser(), metadataId)
			return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("%s when verify user metadataAuth", err), nil)
		}

		log.Debugf("Succeed verify user metadataAuth when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, userType: {%s}, user: {%s}, metadataId: {%s}",
			task.GetTaskId(), taskRole.String(), partyId, task.GetTaskData().GetUserType(), task.GetTaskData().GetUser(), metadataId)

		// Select the datanode where your metadata ID is located.
		var dataNodeId string
		// check the metadata whether internal metadata
		internalMetadataFlag, err := sche.resourceMng.GetDB().IsInternalMetadataById(metadataId)
		if nil != err {
			log.WithError(err).Errorf("failed check metadata whether internal metadata by metaDataId when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
				task.GetTaskId(), taskRole.String(), partyId, metadataId)
			return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("check metadata whether internal failed, %s", err), nil)
		}
		if internalMetadataFlag {
			internalMetadata, err := sche.resourceMng.GetDB().QueryInternalMetadataById(metadataId)
			if nil != err {
				log.WithError(err).Errorf("failed query internal metadataInfo by metaDataId when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
					task.GetTaskId(), taskRole.String(), partyId, metadataId)
				return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("query internal metadata by metadataId failed, %s", err), nil)
			}

			dataOriginId, err := policy.FetchOriginId(internalMetadata.GetData().GetDataType(), internalMetadata.GetData().GetMetadataOption())
			if nil != err {
				log.WithError(err).Errorf("not fetch internalMetadata originId from task metadataOption on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, userType: {%s}, user: {%s}, metadataId: {%s}",
					task.GetTaskId(), taskRole.String(), partyId, task.GetTaskData().GetUserType(), task.GetTaskData().GetUser(), internalMetadata.GetData().GetMetadataId())
				return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("%s when fetch metadataId from task dataPolicy", err), nil)
			}
			dataResourceFileUpload, err := sche.resourceMng.GetDB().QueryDataResourceFileUpload(dataOriginId)
			if nil != err {
				log.WithError(err).Errorf("Failed query dataResourceFileUpload on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
					task.GetTaskId(), taskRole.String(), partyId, metadataId)
				return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("query dataFileUpload by originId failed, %s", err), nil)
			}
			dataNodeId = dataResourceFileUpload.GetNodeId()
		} else {

			dataResourceDiskUsed, err := sche.resourceMng.GetDB().QueryDataResourceDiskUsed(metadataId)
			if nil != err {
				log.WithError(err).Errorf("failed query dataResourceDiskUsed by metaDataId when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
					task.GetTaskId(), taskRole.String(), partyId, metadataId)
				return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("query dataResourceDiskUsed by metadataId failed, %s", err), nil)
			}
			dataNodeId = dataResourceDiskUsed.GetNodeId()
		}

		dataNode, err := sche.resourceMng.GetDB().QueryRegisterNode(pb.PrefixTypeDataNode, dataNodeId)
		if nil != err {
			log.WithError(err).Errorf("failed query internal data resource by metaDataId when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
				task.GetTaskId(), taskRole.String(), partyId, metadataId)
			return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("query internal dataNode by dataNodeId failed, %s", err), nil)
		}

		log.Debugf("Succeed election internal data resource when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, dataNode: %s",
			task.GetTaskId(), taskRole.String(), partyId, dataNode.String())

		result = types.NewReplayScheduleResult(task.GetTaskId(), nil, types.NewPrepareVoteResource(dataNode.Id, dataNode.ExternalIp, dataNode.ExternalPort, partyId))

	// If the current participant is powersupplier.
	// select your own internal power resource and lock it,
	// and finally click 'publishfinishedtasktodatacenter'
	// or 'sendtaskresultmsgtotasksender' in taskmnager.
	case libtypes.TaskRole_TaskRole_PowerSupplier:

		log.Debugf("Succeed CalculateSlotCount when role is powerSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, cost.mem: {%d}, cost.Bandwidth: {%d}, cost.Processor: {%d}",
			task.GetTaskId(), taskRole.String(), partyId,
			task.GetTaskData().GetOperationCost().GetMemory(),
			task.GetTaskData().GetOperationCost().GetBandwidth(),
			task.GetTaskData().GetOperationCost().GetProcessor())

		jobNode, err := sche.elector.ElectionNode(task.GetTaskId(),
			task.GetTaskData().GetOperationCost().GetMemory(),
			task.GetTaskData().GetOperationCost().GetBandwidth(), 0,
			task.GetTaskData().GetOperationCost().GetProcessor(), "")
		if nil != err {
			log.WithError(err).Errorf("Failed to election internal power resource when role is powerSupplier on SchedulerStarveFIFO.ReplaySchedule(),taskId: {%s}, role: {%s}, partyId: {%s}",
				task.GetTaskId(), taskRole.String(), partyId)
			return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("%s when election internal jobNode", err), nil)
		}

		log.Debugf("Succeed election internal power resource when role is powerSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, jobNode: %s",
			task.GetTaskId(), taskRole.String(), partyId, jobNode.String())

		if err := sche.resourceMng.LockLocalResourceWithTask(partyId, jobNode.Id,
			task.GetTaskData().GetOperationCost().GetMemory(),
			task.GetTaskData().GetOperationCost().GetBandwidth(), 0,
			task.GetTaskData().GetOperationCost().GetProcessor(), task); nil != err {
			log.WithError(err).Errorf("Failed to Lock LocalResource when role is powerSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
				task.GetTaskId(), taskRole.String(), partyId, jobNode.Id)
			return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("%s when lock internal jobNode resource", err), nil)
		}

		result = types.NewReplayScheduleResult(task.GetTaskId(), nil, types.NewPrepareVoteResource(jobNode.Id, jobNode.ExternalIp, jobNode.ExternalPort, partyId))

	// If the current participant is resultsupplier.
	// just select their own available datanodes.
	case libtypes.TaskRole_TaskRole_Receiver:

		dataResourceTables, err := sche.resourceMng.GetDB().QueryDataResourceTables()
		if nil != err {
			log.WithError(err).Errorf("Failed to election internal data resource when role is receiver on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}",
				task.GetTaskId(), taskRole.String(), partyId)
			return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("query all internal dataResourceTables failed, %s", err), nil)
		}

		log.Debugf("QueryDataResourceTables when role is receiver on replaySchedule by taskRole is the resuler on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, dataResourceTables: %s",
			task.GetTaskId(), taskRole.String(), partyId, types.UtilDataResourceArrString(dataResourceTables))

		resource := dataResourceTables[len(dataResourceTables)-1]
		dataNode, err := sche.resourceMng.GetDB().QueryRegisterNode(pb.PrefixTypeDataNode, resource.GetNodeId())
		if nil != err {
			log.WithError(err).Errorf("Failed to query internal data resource when role is receiver on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
				task.GetTaskId(), taskRole.String(), partyId, resource.GetNodeId())
			return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("query internal dataNode by dataNodeId failed, %s", err), nil)
		}

		log.Debugf("Succeed election internal data resource when role is receiver on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, dataNode: %s",
			task.GetTaskId(), taskRole.String(), partyId, dataNode.String())

		result = types.NewReplayScheduleResult(task.GetTaskId(), nil, types.NewPrepareVoteResource(dataNode.Id, dataNode.ExternalIp, dataNode.ExternalPort, partyId))

	default:
		result = types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("task role: {%s} can not replay schedule task", taskRole.String()), nil)
	}

	return result
}

// NOTE: schedule powerSuppliers by powerPolicy of task
func (sche *SchedulerStarveFIFO) scheduleVrfElectionPower(taskId string, cost *ctypes.TaskOperationCost, partyIds []string) (
	string, []*libtypes.TaskOrganization, []*libtypes.TaskPowerResourceOption, error) {

	now := timeutils.UnixMsecUint64()
	vrfInput := append([]byte(taskId), bytesutil.Uint64ToBytes(now)...) // input == taskId + nowtime
	// election other org's power resources
	powers, resources, nonce, weights, err := sche.elector.ElectionOrganization(taskId, partyIds,
		nil, cost.GetMem(), cost.GetBandwidth(), 0, cost.GetProcessor(), vrfInput)
	if nil != err {
		log.WithError(err).Errorf("Failed to election powers org on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", taskId)
		return "", nil, nil, fmt.Errorf("election powerOrg failed, %s", err)
	}

	evidenceJson, err := policy.NewVRFElectionEvidence(nonce, weights, now).EncodeJson()
	if nil != err {
		log.WithError(err).Errorf("Failed to encode evidence on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", taskId)
		return "", nil, nil, fmt.Errorf("encode evidence failed, %s", err)
	}
	return evidenceJson, powers, resources, nil
}

func (sche *SchedulerStarveFIFO) scheduleDataNodeProvidePower(task *types.Task, partyIds []string) (
	[]*libtypes.TaskOrganization, []*libtypes.TaskPowerResourceOption, error) {

	dataSupplierCache := make(map[string]*libtypes.TaskOrganization)
	for _, supplier := range task.GetTaskData().GetDataSuppliers() {
		dataSupplierCache[supplier.GetPartyId()] = supplier
	}

	powers := make([]*libtypes.TaskOrganization, len(partyIds))
	//
	for i, partyId := range partyIds {
		supplier, ok := dataSupplierCache[partyId]
		if !ok {
			log.Errorf("not found oranization of dataSupplier with partyId <when the dataNode just jobNode> on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}, partyId: {%s}",
				task.GetTaskId(), partyId)
			return nil, nil, fmt.Errorf("not found oranization of dataSupplier with partyId, %s", partyId)
		}
		powers[i] = supplier
	}

	return powers, nil, nil
}

// NOTE: reSchedule powerSuppliers by powerPolicy of task

func (sche *SchedulerStarveFIFO) reScheduleVrfElectionPower(taskId, nodeId string, powerSuppliers []*libtypes.TaskOrganization, powerResources []*libtypes.TaskPowerResourceOption, evidenceJson string) error {
	var evidence *policy.VRFElectionEvidence
	if err := evidence.DecodeJson(evidenceJson); nil != err {
		return fmt.Errorf("decode evidence failed, %s", err)
	}
	vrfInput := append([]byte(taskId), bytesutil.Uint64ToBytes(evidence.GetElectionAt())...) // input == taskId + nowtime
	// verify orgs of power of task
	if err := sche.elector.VerifyElectionOrganization(
		taskId, powerSuppliers,
		powerResources, nodeId,
		vrfInput, evidence.GetNonce(), evidence.GetWeights()); nil != err {
		return fmt.Errorf("%s when election power organization", err)
	}
	return nil
}

func (sche *SchedulerStarveFIFO) reScheduleDataNodeProvidePower(taskId string, dataSuppliers, powerSuppliers []*libtypes.TaskOrganization, partyIds []string) error {

	if len(powerSuppliers) != len(partyIds) {
		return fmt.Errorf("powerSuppliers count AND partyIds count is not same, powerSuppliers: {%d}, partyIds: {%d}",
			len(powerSuppliers), len(partyIds))
	}

	dataSupplierCache := make(map[string]*libtypes.TaskOrganization)
	for _, supplier := range dataSuppliers {
		dataSupplierCache[supplier.GetPartyId()] = supplier
	}
	powerSupplierCache := make(map[string]*libtypes.TaskOrganization)
	for _, supplier := range powerSuppliers {
		powerSupplierCache[supplier.GetPartyId()] = supplier
	}

	for _, partyId := range partyIds {
		// check from dataSuppliers
		if _, ok := dataSupplierCache[partyId]; !ok {
			log.Errorf("not found oranization of dataSupplier with partyId <when the dataNode just jobNode> on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}, partyId: {%s}",
				taskId, partyId)
			return fmt.Errorf("not found oranization of dataSupplier with partyId, %s", partyId)
		}
		// check from powerSuppliers
		if _, ok := powerSupplierCache[partyId]; !ok {
			log.Errorf("not found oranization of powerSupplier with partyId <when the dataNode just jobNode> on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}, partyId: {%s}",
				taskId, partyId)
			return fmt.Errorf("not found oranization of dataSupplier with partyId, %s", partyId)
		}
	}
	return nil
}
