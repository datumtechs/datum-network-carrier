package schedule

import (
	"container/heap"
	"encoding/json"
	"fmt"
	auth2 "github.com/Metisnetwork/Metis-Carrier/ach/auth"
	"github.com/Metisnetwork/Metis-Carrier/common/bytesutil"
	"github.com/Metisnetwork/Metis-Carrier/common/timeutils"
	ctypes "github.com/Metisnetwork/Metis-Carrier/consensus/twopc/types"
	"github.com/Metisnetwork/Metis-Carrier/core/election"
	"github.com/Metisnetwork/Metis-Carrier/core/evengine"
	"github.com/Metisnetwork/Metis-Carrier/core/rawdb"
	"github.com/Metisnetwork/Metis-Carrier/core/resource"
	libapipb "github.com/Metisnetwork/Metis-Carrier/lib/api"
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

// NOTE: on sender only
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
	caches := make(map[uint32]ScheduleCollecter, 0)

	for i, policyType := range task.GetTaskData().GetPowerPolicyTypes() {

		switch policyType {
		case types.TASK_POWER_POLICY_ASSIGNMENT_SYMBOL_RANDOM_ELECTION:

			var collecter *ScheduleWithSymbolRandomElectionPower
			cache, ok := caches[policyType]
			if !ok {
				collecter = &ScheduleWithSymbolRandomElectionPower{
					partyIds: make([]string, 0),
				}
			} else {
				collecter = cache.(*ScheduleWithSymbolRandomElectionPower)
			}
			collecter.AppendPartyId(task.GetTaskData().GetPowerPolicyOptions()[i])
			caches[policyType] = collecter

			// collection partyId index into cache.
			partyIdAndIndexCache[task.GetTaskData().GetPowerPolicyOptions()[i]] = i
		case types.TASK_POWER_POLICY_DATANODE_PROVIDE:

			var provide *types.TaskPowerPolicyDataNodeProvide
			if err := json.Unmarshal([]byte(task.GetTaskData().GetPowerPolicyOptions()[i]), &provide); nil != err {
				log.WithError(err).Errorf("can not unmarshal powerPolicyType, on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
				return types.NewNeedConsensusTask(task, ""), bullet.GetTaskId(), fmt.Errorf("can not unmarshal powerPolicyType of task, %d", policyType)
			}

			var collecter *ScheduleWithDataNodeProvidePower
			cache, ok := caches[policyType]
			if !ok {
				collecter = &ScheduleWithDataNodeProvidePower{
					provides: make([]*types.TaskPowerPolicyDataNodeProvide, 0),
				}
			} else {
				collecter = cache.(*ScheduleWithDataNodeProvidePower)
			}
			collecter.AppendProvide(provide)
			caches[policyType] = collecter

			// collection partyId index into cache.
			partyIdAndIndexCache[provide.GetPowerPartyId()] = i
		// NOTE: unknown powerPolicyType
		default:
			log.WithError(err).Errorf("unknown powerPolicyType of task on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
			return types.NewNeedConsensusTask(task, ""), bullet.GetTaskId(), fmt.Errorf("unknown powerPolicyType of task, %d", policyType)
		}

	}

	powerOrgs := make([]*libtypes.TaskOrganization, len(task.GetTaskData().GetPowerPolicyTypes()))
	powerResources := make([]*libtypes.TaskPowerResourceOption, len(task.GetTaskData().GetPowerPolicyTypes()))

	var (
		evidence string
	)

	// NOTE: many election policys
	for _, coller := range caches {
		switch coller.(type) {
		// 1、for vrf election power policy
		case *ScheduleWithSymbolRandomElectionPower:

			evidenceJson, orgs, resources, err := sche.scheduleVrfElectionPower(task.GetTaskId(), &ctypes.TaskOperationCost{
				Mem:       task.GetTaskData().GetOperationCost().GetMemory(),
				Processor: task.GetTaskData().GetOperationCost().GetProcessor(),
				Bandwidth: task.GetTaskData().GetOperationCost().GetBandwidth(),
				Duration:  task.GetTaskData().GetOperationCost().GetDuration(),
			}, coller.(*ScheduleWithSymbolRandomElectionPower).GetPartyIds())
			if nil != err {
				log.WithError(err).Errorf("vrf election powerSupplier failed on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
				return types.NewNeedConsensusTask(task, ""), bullet.GetTaskId(), err
			}
			for i, org := range orgs {
				if index, ok := partyIdAndIndexCache[org.GetPartyId()]; !ok {
					log.Errorf("not found partyId of powerSupplier in partyIdAndIndexCache with vrf election powerSupplier on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}, partyId: {%s}", task.GetTaskId(), org.GetPartyId())
					return types.NewNeedConsensusTask(task, ""), bullet.GetTaskId(), fmt.Errorf("invalid partyId")
				} else {
					powerOrgs[index] = org
					powerResources[index] = resources[i]
				}
			}
			evidence = evidenceJson

		// 2、for dataProvider election power policy
		case *ScheduleWithDataNodeProvidePower:

			orgs, resources, err := sche.scheduleDataNodeProvidePower(task, coller.(*ScheduleWithDataNodeProvidePower).Getprovides())
			if nil != err {
				log.WithError(err).Errorf("dataNodeProvider election powerSupplier failed on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}", task.GetTaskId())
				return types.NewNeedConsensusTask(task, ""), bullet.GetTaskId(), err
			}

			for i, org := range orgs {
				if index, ok := partyIdAndIndexCache[org.GetPartyId()]; !ok {
					log.Errorf("not found partyId of powerSupplier in partyIdAndIndexCache with dataNodeProvider election powerSupplier on SchedulerStarveFIFO.TrySchedule(), taskId: {%s}, partyId: {%s}", task.GetTaskId(), org.GetPartyId())
					return types.NewNeedConsensusTask(task, ""), bullet.GetTaskId(), fmt.Errorf("invalid partyId")
				} else {
					powerOrgs[index] = org
					powerResources[index] = resources[i]
				}
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

	findDataNodeByMetadataIdFn := func(metadataId string) (*libapipb.YarnRegisteredPeerDetail, error) {
		// Select the datanode where your metadata ID is located.
		var dataNodeId string
		// check the metadata whether internal metadata
		internalMetadataFlag, err := sche.resourceMng.GetDB().IsInternalMetadataById(metadataId)
		if nil != err {
			log.WithError(err).Errorf("failed check metadata whether internal metadata by metaDataId when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
				task.GetTaskId(), taskRole.String(), partyId, metadataId)
			return nil, fmt.Errorf("check metadata whether internal failed, %s", err)
		}
		if internalMetadataFlag {
			internalMetadata, err := sche.resourceMng.GetDB().QueryInternalMetadataById(metadataId)
			if nil != err {
				log.WithError(err).Errorf("failed query internal metadataInfo by metaDataId when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
					task.GetTaskId(), taskRole.String(), partyId, metadataId)
				return nil, fmt.Errorf("query internal metadata by metadataId failed, %s", err)
			}

			dataOriginId, err := policy.FetchOriginId(internalMetadata.GetData().GetDataType(), internalMetadata.GetData().GetMetadataOption())
			if nil != err {
				log.WithError(err).Errorf("not fetch internalMetadata originId from task metadataOption on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, userType: {%s}, user: {%s}, metadataId: {%s}",
					task.GetTaskId(), taskRole.String(), partyId, task.GetTaskData().GetUserType(), task.GetTaskData().GetUser(), internalMetadata.GetData().GetMetadataId())
				return nil, fmt.Errorf("%s when fetch metadataId from task dataPolicy", err)
			}
			dataResourceFileUpload, err := sche.resourceMng.GetDB().QueryDataResourceFileUpload(dataOriginId)
			if nil != err {
				log.WithError(err).Errorf("Failed query dataResourceFileUpload on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
					task.GetTaskId(), taskRole.String(), partyId, metadataId)
				return nil, fmt.Errorf("query dataFileUpload by originId failed, %s", err)
			}
			dataNodeId = dataResourceFileUpload.GetNodeId()
		} else {

			dataResourceDiskUsed, err := sche.resourceMng.GetDB().QueryDataResourceDiskUsed(metadataId)
			if nil != err {
				log.WithError(err).Errorf("failed query dataResourceDiskUsed by metaDataId when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
					task.GetTaskId(), taskRole.String(), partyId, metadataId)
				return nil, fmt.Errorf("query dataResourceDiskUsed by metadataId failed, %s", err)
			}
			dataNodeId = dataResourceDiskUsed.GetNodeId()
		}

		dataNode, err := sche.resourceMng.GetDB().QueryRegisterNode(libapipb.PrefixTypeDataNode, dataNodeId)
		if nil != err {
			log.WithError(err).Errorf("failed query internal data resource by metaDataId when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, metadataId: {%s}",
				task.GetTaskId(), taskRole.String(), partyId, metadataId)
			return nil, fmt.Errorf("query internal dataNode by dataNodeId failed, %s", err)
		}
		return dataNode, nil
	}

	var result *types.ReplayScheduleResult

	switch taskRole {

	case libtypes.TaskRole_TaskRole_DataSupplier:

		// ## 1、verify the powsuppliers election result
		caches := make(map[uint32]ReScheduleCollecter, 0)

		for i, policyType := range task.GetTaskData().GetPowerPolicyTypes() {

			switch policyType {
			case types.TASK_POWER_POLICY_ASSIGNMENT_SYMBOL_RANDOM_ELECTION:

				var collecter *ReScheduleWithSymbolRandomElectionPower
				cache, ok := caches[policyType]
				if !ok {
					collecter = &ReScheduleWithSymbolRandomElectionPower{
						suppliers: make([]*libtypes.TaskOrganization, 0),
						resources: make([]*libtypes.TaskPowerResourceOption, 0),
					}
				} else {
					collecter = cache.(*ReScheduleWithSymbolRandomElectionPower)
				}
				collecter.AppendSupplier(task.GetTaskData().GetPowerSuppliers()[i])
				collecter.AppendResource(task.GetTaskData().GetPowerResourceOptions()[i])
				caches[policyType] = collecter

			case types.TASK_POWER_POLICY_DATANODE_PROVIDE:

				var provide *types.TaskPowerPolicyDataNodeProvide
				if err := json.Unmarshal([]byte(task.GetTaskData().GetPowerPolicyOptions()[i]), &provide); nil != err {
					log.WithError(err).Errorf("can not unmarshal powerPolicyType, on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}", task.GetTaskId())
					return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("can not unmarshal powerPolicyType of task, %d", policyType), nil)
				}
				var collecter *ReScheduleWithDataNodeProvidePower
				cache, ok := caches[policyType]
				if !ok {
					collecter = &ReScheduleWithDataNodeProvidePower{
						suppliers: make([]*libtypes.TaskOrganization, 0),
						provides: make([]*types.TaskPowerPolicyDataNodeProvide, 0),
					}
				} else {
					collecter = cache.(*ReScheduleWithDataNodeProvidePower)
				}
				collecter.AppendSupplier(task.GetTaskData().GetPowerSuppliers()[i])
				collecter.AppendProvide(provide)
				caches[policyType] = collecter

			// NOTE: unknown powerPolicyType
			default:
				log.WithError(err).Errorf("unknown powerPolicyType of task on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}", task.GetTaskId())
				return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
			}
		}

		// NOTE:vefiry many election policys
		for _, coller := range caches {
			switch coller.(type) {
			// a、for vrf election power policy
			case *ReScheduleWithSymbolRandomElectionPower:

				c := coller.(*ReScheduleWithSymbolRandomElectionPower)
				if len(c.GetSuppliers()) != len(c.GetResources()) {
					log.Errorf("assignmentSymbolRandom: election powerSuppliers count and election powerResources count is not same on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, powerSuppliers len: %d, powerResources len: %d",
						task.GetTaskId(), len(c.GetSuppliers()), len(c.GetResources()))
					return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
				}
				if err := sche.reScheduleVrfElectionPower(task.GetTaskId(), task.GetTaskSender().GetNodeId(),
					c.GetSuppliers(), c.GetResources(), replayTask.GetEvidence()); nil != err {
					return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
				}

			// b、for dataProvide election power policy
			case *ReScheduleWithDataNodeProvidePower:

				c := coller.(*ReScheduleWithDataNodeProvidePower)

				if len(c.GetSuppliers()) != len(c.GetProvides()) {
					log.Errorf("dataNodeProvider: election powerSuppliers count and election providePolicys count is not same on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, powerSuppliers len: %d, providePolicys len: %d",
						task.GetTaskId(), len(c.GetSuppliers()), len(c.GetProvides()))
					return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
				}
				if err := sche.reScheduleDataNodeProvidePower(task.GetTaskId(), task.GetTaskData().GetDataSuppliers(), c.GetSuppliers(), c.GetProvides()); nil != err {
					return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
				}
			}
		}

		log.Debugf("Succeed to verify election powers org when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, agree: %v",
			task.GetTaskId(), taskRole.String(), partyId)

		// ## 2、check metadata on current party AND choosing the dataNode

		// Find metadataId of current identyt of task with current partyId.
		var metadataId string
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			if identityId == dataSupplier.GetIdentityId() && partyId == dataSupplier.GetPartyId() {
				mId, err := policy.FetchMetedataIdByPartyIdFromDataPolicy(partyId, task.GetTaskData().GetDataPolicyTypes(), task.GetTaskData().GetDataPolicyOptions())
				if nil != err {
					log.WithError(err).Errorf("not fetch metadataId from task dataPolicy on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, userType: {%s}, user: {%s}",
						task.GetTaskId(), taskRole.String(), partyId, task.GetTaskData().GetUserType(), task.GetTaskData().GetUser())
					return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("%s when fetch metadataId from task dataPolicy", err), nil)
				}
				metadataId = mId
				break
			}
		}

		dataNode, err := findDataNodeByMetadataIdFn(metadataId)
		if nil != err {
			return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
		}

		log.Debugf("Succeed election internal data resource when role is dataSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, dataNode: %s",
			task.GetTaskId(), taskRole.String(), partyId, dataNode.String())

		result = types.NewReplayScheduleResult(task.GetTaskId(), nil, types.NewPrepareVoteResource(dataNode.Id, dataNode.ExternalIp, dataNode.ExternalPort, partyId))

	// If the current participant is powersupplier.
	// select your own internal power resource and lock it,
	// and finally click 'publishfinishedtasktodatacenter'
	// or 'sendtaskresultmsgtotasksender' in taskmnager.
	case libtypes.TaskRole_TaskRole_PowerSupplier:

		// ## 1、choosing powerSupplier jobNode

		var (
			node *libapipb.YarnRegisteredPeerDetail
			err  error
		)

		for i, policyType := range task.GetTaskData().GetPowerPolicyTypes() {

			switch policyType {
			case types.TASK_POWER_POLICY_ASSIGNMENT_SYMBOL_RANDOM_ELECTION:

				if task.GetTaskData().GetPowerPolicyOptions()[i] == partyId && task.GetTaskData().GetPowerSuppliers()[i].GetPartyId() == partyId {

					log.Debugf("Succeed CalculateSlotCount when role is powerSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, cost.mem: {%d}, cost.Bandwidth: {%d}, cost.Processor: {%d}",
						task.GetTaskId(), taskRole.String(), partyId,
						task.GetTaskData().GetOperationCost().GetMemory(),
						task.GetTaskData().GetOperationCost().GetBandwidth(),
						task.GetTaskData().GetOperationCost().GetProcessor())

					// election jobNode machine
					node, err = sche.elector.ElectionNode(
						task.GetTaskId(), partyId,
						task.GetTaskData().GetOperationCost().GetMemory(),
						task.GetTaskData().GetOperationCost().GetBandwidth(), 0,
						task.GetTaskData().GetOperationCost().GetProcessor(), "")
					if nil != err {
						log.WithError(err).Errorf("Failed to election internal power resource when role is powerSupplier on SchedulerStarveFIFO.ReplaySchedule(),taskId: {%s}, role: {%s}, partyId: {%s}",
							task.GetTaskId(), taskRole.String(), partyId)
						return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("%s when election internal jobNode", err), nil)
					}

					// Lock the jobNode resource
					if err := sche.resourceMng.LockLocalResourceWithTask(
						partyId,
						node.GetId(),
						task.GetTaskData().GetOperationCost().GetMemory(),
						task.GetTaskData().GetOperationCost().GetBandwidth(),
						0,
						task.GetTaskData().GetOperationCost().GetProcessor(),
						task); nil != err {
						log.WithError(err).Errorf("Failed to Lock LocalResource when role is powerSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
							task.GetTaskId(), taskRole.String(), partyId, node.GetId())
						return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("%s when lock internal jobNode resource", err), nil)
					}

					break // break the external loop
				}

			case types.TASK_POWER_POLICY_DATANODE_PROVIDE:

				var powerPolicy *types.TaskPowerPolicyDataNodeProvide
				if err := json.Unmarshal([]byte(task.GetTaskData().GetPowerPolicyOptions()[i]), &powerPolicy); nil != err {
					log.WithError(err).Errorf("can not unmarshal powerPolicyType, on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}", task.GetTaskId())
					return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("can not unmarshal powerPolicyType of task, %d", policyType), nil)
				}

				if powerPolicy.GetPowerPartyId() == partyId && task.GetTaskData().GetPowerSuppliers()[i].GetPartyId() == partyId {

					for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {

						// ### NOTE:
						//		Only when the receiver and the dataSupplier belong to the same organization
						//		that can the datasupplier's datanode be specified as the receiver's datanode
						if powerPolicy.GetProviderPartyId() == dataSupplier.GetPartyId() && identityId == dataSupplier.GetIdentityId() {

							metadataId, err := policy.FetchMetedataIdByPartyIdFromDataPolicy(dataSupplier.GetPartyId(), task.GetTaskData().GetDataPolicyTypes(), task.GetTaskData().GetDataPolicyOptions())
							if nil != err {
								log.WithError(err).Errorf("not fetch metadataId from task dataPolicy when dataNode provider election jobNode on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, powerSupplier partyId: {%s}, power provider dataSupplier partyId: {%s}, userType: {%s}, user: {%s}",
									task.GetTaskId(), taskRole.String(), partyId, dataSupplier.GetPartyId(), task.GetTaskData().GetUserType(), task.GetTaskData().GetUser())
								return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("%s when fetch metadataId from task dataPolicy", err), nil)
							}

							// election dataNode machine to be jobNode
							node, err = findDataNodeByMetadataIdFn(metadataId)
							if nil != err {
								return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
							}
							break // break the internal loop
						}
					}

					break // break the external loop
				}

			// NOTE: unknown powerPolicyType
			default:
				log.WithError(err).Errorf("unknown powerPolicyType of task on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}", task.GetTaskId())
				return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
			}
		}

		if nil == node {
			return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("cannot found jobNode"), nil)
		}

		log.Debugf("Succeed election internal power resource when role is powerSupplier on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, node(jobNode|dataNode): %s",
			task.GetTaskId(), taskRole.String(), partyId, node.String())

		result = types.NewReplayScheduleResult(task.GetTaskId(), nil, types.NewPrepareVoteResource(node.Id, node.GetExternalIp(), node.GetExternalPort(), partyId))

	// If the current participant is resultsupplier.
	// just select their own available datanodes.
	case libtypes.TaskRole_TaskRole_Receiver:

		// ## 1、choosing receiver dataNode

		var (
			dataNode *libapipb.YarnRegisteredPeerDetail
			err      error
		)

		for i, policyType := range task.GetTaskData().GetReceiverPolicyTypes() {

			switch policyType {
			case types.TASK_RECEIVER_POLICY_RANDOM_ELECTION:

				if task.GetTaskData().GetReceiverPolicyOptions()[i] == partyId && task.GetTaskData().GetReceivers()[i].GetPartyId() == partyId {

					dataResourceTables, err := sche.resourceMng.GetDB().QueryDataResourceTables()
					if nil != err {
						log.WithError(err).Errorf("Failed to election internal data resource when role is receiver on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}",
							task.GetTaskId(), taskRole.String(), partyId)
						return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("query all internal dataResourceTables failed, %s", err), nil)
					}

					log.Debugf("QueryDataResourceTables when role is receiver on replaySchedule by taskRole is the resuler on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, dataResourceTables: %s",
						task.GetTaskId(), taskRole.String(), partyId, types.UtilDataResourceArrString(dataResourceTables))

					resource := dataResourceTables[len(dataResourceTables)-1]
					dataNode, err = sche.resourceMng.GetDB().QueryRegisterNode(libapipb.PrefixTypeDataNode, resource.GetNodeId())
					if nil != err {
						log.WithError(err).Errorf("Failed to query internal data resource when role is receiver on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
							task.GetTaskId(), taskRole.String(), partyId, resource.GetNodeId())
						return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("query internal dataNode by dataNodeId failed, %s", err), nil)
					}

					break // break the external loop
				}

			case types.TASK_RECEIVER_POLICY_DATANODE_PROVIDE:

				var receiverPolicy *types.TaskReceiverPolicyDataNodeProvide
				if err = json.Unmarshal([]byte(task.GetTaskData().GetReceiverPolicyOptions()[i]), &receiverPolicy); nil != err {
					log.WithError(err).Errorf("can not unmarshal receiverPolicyType, on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}", task.GetTaskId())
					return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("can not unmarshal receiverPolicyType of task, %d", policyType), nil)
				}

				if receiverPolicy.GetReceiverPartyId() == partyId && task.GetTaskData().GetReceivers()[i].GetPartyId() == partyId {

					for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {

						// ### NOTE:
						//		Only when the receiver and the dataSupplier belong to the same organization
						//		that can the datasupplier's datanode be specified as the receiver's datanode
						if receiverPolicy.GetProviderPartyId() == dataSupplier.GetPartyId() && identityId == dataSupplier.GetIdentityId() {

							metadataId, err := policy.FetchMetedataIdByPartyIdFromDataPolicy(dataSupplier.GetPartyId(), task.GetTaskData().GetDataPolicyTypes(), task.GetTaskData().GetDataPolicyOptions())
							if nil != err {
								log.WithError(err).Errorf("not fetch metadataId from task dataPolicy when dataNode provider receiver on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, receiver partyId: {%s}, receiver provider dataSupplier partyId: {%s}, userType: {%s}, user: {%s}",
									task.GetTaskId(), taskRole.String(), partyId, dataSupplier.GetPartyId(), task.GetTaskData().GetUserType(), task.GetTaskData().GetUser())
								return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("%s when fetch metadataId from task dataPolicy", err), nil)
							}
							dataNode, err = findDataNodeByMetadataIdFn(metadataId)
							if nil != err {
								return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
							}
							break // break the internal loop
						}
					}

					break // break the external loop
				}

			// NOTE: unknown powerPolicyType
			default:
				log.WithError(err).Errorf("unknown receiverPolicyType of task on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}", task.GetTaskId())
				return types.NewReplayScheduleResult(task.GetTaskId(), err, nil)
			}
		}

		if nil == dataNode {
			return types.NewReplayScheduleResult(task.GetTaskId(), fmt.Errorf("cannot found dataNode"), nil)
		}

		//NEXT:

		log.Debugf("Succeed election internal data resource when role is receiver on SchedulerStarveFIFO.ReplaySchedule(), taskId: {%s}, role: {%s}, partyId: {%s}, dataNode: %s",
			task.GetTaskId(), taskRole.String(), partyId, dataNode.String())

		result = types.NewReplayScheduleResult(task.GetTaskId(), nil, types.NewPrepareVoteResource(dataNode.Id, dataNode.GetExternalIp(), dataNode.GetExternalPort(), partyId))

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
		log.WithError(err).Errorf("Failed to election powers org on scheduleVrfElectionPower(), taskId: {%s}", taskId)
		return "", nil, nil, fmt.Errorf("election powerOrg failed, %s", err)
	}

	evidenceJson, err := policy.NewVRFElectionEvidence(nonce, weights, now).EncodeJson()
	if nil != err {
		log.WithError(err).Errorf("Failed to encode evidence on scheduleVrfElectionPower(), taskId: {%s}", taskId)
		return "", nil, nil, fmt.Errorf("encode evidence failed, %s", err)
	}
	return evidenceJson, powers, resources, nil
}

func (sche *SchedulerStarveFIFO) scheduleDataNodeProvidePower(task *types.Task, provides []*types.TaskPowerPolicyDataNodeProvide) (
	[]*libtypes.TaskOrganization, []*libtypes.TaskPowerResourceOption, error) {

	dataSupplierCache := make(map[string]*libtypes.TaskOrganization)
	for _, supplier := range task.GetTaskData().GetDataSuppliers() {
		newPowerSupplier := &libtypes.TaskOrganization{
			IdentityId: supplier.GetIdentityId(),
			NodeId:     supplier.GetNodeId(),
			NodeName:   supplier.GetNodeName(),
			PartyId:    "",
		}
		dataSupplierCache[supplier.GetPartyId()] = newPowerSupplier
	}

	powers := make([]*libtypes.TaskOrganization, len(provides))
	resources := make([]*libtypes.TaskPowerResourceOption, len(provides))
	//
	for i, provide := range provides {
		supplier, ok := dataSupplierCache[provide.GetProviderPartyId()]
		if !ok {
			log.Errorf("not found oranization of dataSupplier with partyId on scheduleDataNodeProvidePower(), taskId: {%s}, provide partyId: {%s},  power partyId: {%s}",
				task.GetTaskId(), provide.GetProviderPartyId(), provide.GetPowerPartyId())
			return nil, nil, fmt.Errorf("not found oranization of dataSupplier with partyId, %s to %s", provide.GetProviderPartyId(), provide.GetPowerPartyId())
		}
		supplier.PartyId = provide.GetPowerPartyId()
		powers[i] = supplier
		resources[i] = &libtypes.TaskPowerResourceOption{
			PartyId: provide.GetPowerPartyId(),
			ResourceUsedOverview: &libtypes.ResourceUsageOverview{
				TotalMem:       0, // total resource value of org.
				UsedMem:        0, // used resource of this task (real time max used)
				TotalBandwidth: 0,
				UsedBandwidth:  0, // used resource of this task (real time max used)
				TotalDisk:      0,
				UsedDisk:       0,
				TotalProcessor: 0,
				UsedProcessor:  0, // used resource of this task (real time max used)
			},
		}
	}

	return powers, resources, nil
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

func (sche *SchedulerStarveFIFO) reScheduleDataNodeProvidePower(taskId string, dataSuppliers, powerSuppliers []*libtypes.TaskOrganization, provides []*types.TaskPowerPolicyDataNodeProvide) error {

	if len(powerSuppliers) != len(provides) {
		return fmt.Errorf("powerSuppliers count AND partyIds count is not same, powerSuppliers: {%d}, provides: {%d}",
			len(powerSuppliers), len(provides))
	}

	dataSupplierCache := make(map[string]*libtypes.TaskOrganization)
	for _, supplier := range dataSuppliers {
		dataSupplierCache[supplier.GetPartyId()] = supplier
	}
	powerSupplierCache := make(map[string]*libtypes.TaskOrganization)
	for _, supplier := range powerSuppliers {
		powerSupplierCache[supplier.GetPartyId()] = supplier
	}

	for _, provide := range provides {
		// check from dataSuppliers
		dataSupplier, ok := dataSupplierCache[provide.GetProviderPartyId()]
		if !ok {
			log.Errorf("not found oranization of dataSupplier with partyId on reScheduleDataNodeProvidePower(), taskId: {%s}, provide partyId: {%s},  power partyId: {%s}",
				taskId, provide.GetProviderPartyId(), provide.GetPowerPartyId())
			return fmt.Errorf("not found oranization of dataSupplier with partyId, %s to %s", provide.GetProviderPartyId(), provide.GetPowerPartyId())
		}
		// check from powerSuppliers
		powerSupplier, ok := powerSupplierCache[provide.GetPowerPartyId()]
		if !ok {
			log.Errorf("not found oranization of powerSupplier with partyId on reScheduleDataNodeProvidePower(), taskId: {%s}, provide partyId: {%s},  power partyId: {%s}",
				taskId, provide.GetProviderPartyId(), provide.GetPowerPartyId())
			return fmt.Errorf("not found oranization of dataSupplier with partyId, %s to %s", provide.GetProviderPartyId(), provide.GetPowerPartyId())
		}
		// dataSupplier must equal powerSupplier
		if dataSupplier.GetIdentityId() != powerSupplier.GetIdentityId() {
			log.Errorf("the corresponding power provider of dataSupplier organization is inconsistent with the selected powerSupplier organization on reScheduleDataNodeProvidePower(), taskId: {%s}, provide partyId: {%s},  power partyId: {%s}",
				taskId, provide.GetProviderPartyId(), provide.GetPowerPartyId())
			return fmt.Errorf("the corresponding power provider of dataSupplier organization is inconsistent with the selected powerSupplier organization, dataSupplier: [partyId: %s, identityId: %s], powerSupplier:[partyId: %s, identityId: %s]",
				dataSupplier.GetPartyId(), dataSupplier.GetIdentityId(), powerSupplier.GetPartyId(), powerSupplier.GetIdentityId())
		}
	}

	return nil
}
