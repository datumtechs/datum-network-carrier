package task

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/auth"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/consensus"
	ev "github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/core/schedule"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strings"
	"sync"
	"time"
)

const (
	defaultScheduleTaskInterval = 2 * time.Second
	taskMonitorInterval         = 30 * time.Second
	senderExecuteTaskExpire     = 10 * time.Second
)

//type Scheduler interface {
//	Start() error
//	Stop() error
//	Error() error
//	Name() string
//	AddTask(task *types.Task) error
//	RemoveTask(taskId string) error
//	TrySchedule() (*types.NeedConsensusTask, error)
//	ReplaySchedule(myPartyId string, myTaskRole apipb.TaskRole, task *types.Task) *types.ReplayScheduleResult
//}

type Manager struct {
	p2p             p2p.P2P
	scheduler       schedule.Scheduler
	consensusEngine consensus.Engine
	eventEngine     *ev.EventEngine
	resourceMng     *resource.Manager
	authMng         *auth.AuthorityManager
	parser          *TaskParser
	validator       *TaskValidator
	// internal resource node set (Fighter node grpc client set)
	resourceClientSet *grpclient.InternalResourceClientSet
	eventCh           chan *types.ReportTaskEvent
	quit              chan struct{}
	// send the validated taskMsgs to scheduler
	localTasksCh             chan types.TaskDataArray
	needReplayScheduleTaskCh chan *types.NeedReplayScheduleTask
	needExecuteTaskCh        chan *types.NeedExecuteTask
	runningTaskCache         map[string]map[string]*types.NeedExecuteTask //  taskId -> {partyId -> task}
	runningTaskCacheLock     sync.RWMutex

	// TODO 有些缓存需要持久化
}

func NewTaskManager(
	p2p p2p.P2P,
	scheduler schedule.Scheduler,
	consensusEngine consensus.Engine,
	eventEngine *ev.EventEngine,
	resourceMng *resource.Manager,
	authMng *auth.AuthorityManager,
	resourceClientSet *grpclient.InternalResourceClientSet,
	localTasksCh chan types.TaskDataArray,
	needReplayScheduleTaskCh chan *types.NeedReplayScheduleTask,
	needExecuteTaskCh chan *types.NeedExecuteTask,
) *Manager {

	m := &Manager{
		p2p:                      p2p,
		scheduler:                scheduler,
		consensusEngine:          consensusEngine,
		eventEngine:              eventEngine,
		resourceMng:              resourceMng,
		authMng:                  authMng,
		resourceClientSet:        resourceClientSet,
		parser:                   newTaskParser(resourceMng),
		validator:                newTaskValidator(resourceMng, authMng),
		eventCh:                  make(chan *types.ReportTaskEvent, 10),
		localTasksCh:             localTasksCh,
		needReplayScheduleTaskCh: needReplayScheduleTaskCh,
		needExecuteTaskCh:        needExecuteTaskCh,
		runningTaskCache:         make(map[string]map[string]*types.NeedExecuteTask, 0), // taskId -> partyId -> needExecuteTask
		quit:                     make(chan struct{}),
	}
	return m
}

func (m *Manager) Start() error {
	go m.loop()
	log.Info("Started taskManager ...")
	return nil
}
func (m *Manager) Stop() error {
	close(m.quit)
	return nil
}

func (m *Manager) loop() {

	taskMonitorTicker := time.NewTicker(taskMonitorInterval)
	//taskTicker := time.NewTicker(defaultScheduleTaskInterval)

	for {
		select {

		// handle reported event from fighter of self organization
		case event := <-m.eventCh:
			go func() {
				if err := m.handleTaskEvent(event.GetPartyId(), event.GetEvent()); nil != err {
					log.Errorf("Failed to call handleTaskEvent() on TaskManager, taskId: {%s}, event: %s", event.GetEvent().GetTaskId(), event.GetEvent().String())
				}
			}()

		// To schedule local task while received some local task msg
		case tasks := <-m.localTasksCh:

			for _, task := range tasks {
				if err := m.scheduler.AddTask(task); nil != err {
					log.Errorf("Failed to add local task into scheduler queue, err: {%s}", err)
					continue
				}

				if err := m.tryScheduleTask(); nil != err {
					log.Errorf("Failed to try schedule local task while received local tasks, err: {%s}", err)
					continue
				}

			}

		// To schedule local task interval
		//case <-taskTicker.C:
		//
		//	if err := m.tryScheduleTask(); nil != err {
		//		log.Errorf("Failed to try schedule local task when taskTicker, err: {%s}", err)
		//	}

		// handle the task of need replay scheduling while received from remote peer on consensus epoch
		case needReplayScheduleTask := <-m.needReplayScheduleTaskCh:

			go func() {

				// Do duplication check ...
				_, err := m.resourceMng.GetDB().QueryLocalTask(needReplayScheduleTask.GetTask().GetTaskId())
				if rawdb.IsNoDBNotFoundErr(err) {
					log.WithError(err).Errorf("Failed to query local task when received remote task, taskId: {%s}", needReplayScheduleTask.GetTask().GetTaskId())
					return
				}

				if rawdb.IsDBNotFoundErr(err) {
					// store metadata used taskId
					if err := m.storeMetaUsedTaskId(needReplayScheduleTask.GetTask()); nil != err {
						log.WithError(err).Errorf("Failed to store metadata used taskId when received remote task, taskId: {%s}", needReplayScheduleTask.GetTask().GetTaskId())
					}
					if err := m.resourceMng.GetDB().StoreLocalTask(needReplayScheduleTask.GetTask()); nil != err {
						log.WithError(err).Errorf("Failed to call StoreLocalTask when replay schedule remote task, taskId: {%s}", needReplayScheduleTask.GetTask().GetTaskId())
						needReplayScheduleTask.SendFailedResult(needReplayScheduleTask.GetTask().GetTaskId(), err)
						return
					}
				}

				// Start replay schedule remote task ...
				result := m.scheduler.ReplaySchedule(needReplayScheduleTask.GetLocalPartyId(), needReplayScheduleTask.GetLocalTaskRole(), needReplayScheduleTask.GetTask())
				needReplayScheduleTask.SendResult(result)
			}()

		// handle task of need to executing, and send it to fighter of myself organization or send the task result msg to remote peer
		case task := <-m.needExecuteTaskCh:

			switch task.GetConsStatus() {
			// sender and partner to handle needExecuteTask when consensus succeed.
			// sender need to store some cache, partner need to execute task.
			case types.TaskNeedExecute, types.TaskConsensusFinished:
				// to execute the task
				m.handleNeedExecuteTask(task)

			// sender and partner to clear local task things after received status: `scheduleFailed` and `interrupt` and `terminate`.
			// sender need to publish local task and event to datacenter,
			// partner need to send task's event to remote task's sender.
			default:
				m.storeTaskFinalEvent(task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(),
					task.GetLocalTaskOrganization().GetPartyId(),
					fmt.Sprintf("execute task: %s with %s", task.GetConsStatus().String(),
						task.GetLocalTaskOrganization().GetPartyId()), apicommonpb.TaskState_TaskState_Failed)
				switch task.GetLocalTaskRole() {
				case apicommonpb.TaskRole_TaskRole_Sender:
					m.publishFinishedTaskToDataCenter(task, true)
				default:
					m.sendTaskResultMsgToRemotePeer(task)
				}
			}

		// handle the executing expire tasks
		case <-taskMonitorTicker.C:

			m.expireTaskMonitor()

		case <-m.quit:
			log.Info("Stopped taskManager ...")
			return
		}
	}
}

func (m *Manager) TerminateTask (terminate *types.TaskTerminateMsg) {

	localTask, err := m.resourceMng.GetDB().QueryLocalTask(terminate.GetTaskId())
	if nil != err {
		log.Errorf("Failed to query local task on `taskManager.TerminateTask()`, taskId: {%s}, err: {%s}", terminate.GetTaskId(), err)
		return
	}

	if nil == localTask {
		log.Errorf("Not found local task on `taskManager.TerminateTask()`, taskId: {%s}", terminate.GetTaskId())
		return
	}

	// check user
	if localTask.GetTaskData().GetUser() != terminate.GetUser() ||
		localTask.GetTaskData().GetUserType() != terminate.GetUserType() {

		log.Errorf("terminate task user and publish task user must be same on `taskManager.TerminateTask()`, taskId: {%s}, partyId: {%s}",
			localTask.GetTaskId(), localTask.GetTaskSender().GetPartyId())
		return
	}

	// todo verify user sign with terminate task

	// The task sender only makes consensus, so interrupt consensus while need terminate task with task sender
	// While task is consensus or executing, can terminate.
	has, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusValConsByPartyId(localTask.GetTaskId(), localTask.GetTaskSender().GetPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task execute `cons` status on `taskManager.TerminateTask()`, taskId: {%s}, partyId: {%s}",
			localTask.GetTaskId(), localTask.GetTaskSender().GetPartyId())
		return
	}
	if has {
		if err = m.consensusEngine.OnConsensusMsg(
			"", types.NewInterruptMsgWrap(localTask.GetTaskId(),
			types.MakeMsgOption(common.Hash{},
			apicommonpb.TaskRole_TaskRole_Sender,
			apicommonpb.TaskRole_TaskRole_Sender,
				localTask.GetTaskSender().GetPartyId(),
				localTask.GetTaskSender().GetPartyId(),
				localTask.GetTaskSender()))); nil != err {
			log.WithError(err).Errorf("Failed to call `OnConsensusMsg()` on `taskManager.TerminateTask()`, taskId: {%s}, partyId: {%s}",
				localTask.GetTaskId(), localTask.GetTaskSender().GetPartyId())
			return
		}
		// The task sender does not need to release local resources, but needs to remove the taskExecuteStatus for the consensus.
		m.resourceMng.GetDB().RemoveLocalTaskExecuteStatusByPartyId(localTask.GetTaskId(), localTask.GetTaskSender().GetPartyId())
	}

	// Anyway, need send terminateMsg to remote partners
	if err := m.sendTaskTerminateMsg(localTask); nil != err {
		log.Errorf("Failed to call `sendTaskTerminateMsg()` on `taskManager.TerminateTask()`, taskId: {%s}, err: \n%s", terminate.GetTaskId(), err)
	}
}

func (m *Manager) SendTaskMsgArr(msgArr types.TaskMsgArr) error {

	if len(msgArr) == 0 {
		return fmt.Errorf("Receive some empty task msgArr")
	}

	nonParsedMsgArr, parsedMsgArr, err := m.parser.ParseTask(msgArr)
	if nil != err {
		for _, badMsg := range nonParsedMsgArr {
			events := []*libtypes.TaskEvent{m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
				badMsg.GetTaskId(), badMsg.GetSenderIdentityId(), badMsg.GetSenderPartyId(), fmt.Sprintf("failed to parse local taskMsg"))}

			if e := m.storeBadTask(badMsg.Data, events, "failed to parse taskMsg"); nil != e {
				log.Errorf("Failed to store the err taskMsg on taskManager, taskId: {%s}", badMsg.GetTaskId())
			}
		}

		if len(nonParsedMsgArr) == len(msgArr) {
			return err
		}
	}

	nonValidatedMsgArr, validatedMsgArr, err := m.validator.validateTaskMsg(parsedMsgArr)
	if nil != err {

		for _, badMsg := range nonValidatedMsgArr {
			events := []*libtypes.TaskEvent{m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
				badMsg.GetTaskId(), badMsg.GetSenderIdentityId(), badMsg.GetSenderPartyId(), fmt.Sprintf("failed to validate local taskMsg"))}

			if e := m.storeBadTask(badMsg.Data, events, "failed to validate taskMsg"); nil != e {
				log.Errorf("Failed to store the err taskMsg on taskManager, taskId: {%s}", badMsg.GetTaskId())
			}
		}

		if len(nonValidatedMsgArr) == len(parsedMsgArr) {
			return err
		}
	}

	taskArr := make(types.TaskDataArray, 0)

	for _, msg := range validatedMsgArr {

		task := msg.Data

		// store metadata used taskId
		if err := m.storeMetaUsedTaskId(task); nil != err {
			log.WithError(err).Errorf("Failed to store metadata used taskId when received local task, taskId: {%s}", task.GetTaskId())
		}

		var storeErrs []string
		if err := m.resourceMng.GetDB().StoreLocalTask(task); nil != err {
			storeErrs = append(storeErrs, err.Error())
		}
		if err := m.resourceMng.GetDB().StoreTaskPowerPartyIds(task.GetTaskId(), msg.GetPowerPartyIds()); nil != err {
			storeErrs = append(storeErrs, err.Error())
		}

		if err := m.storeTaskHanlderPartyIds(task, msg.GetPowerPartyIds()); nil != err {
			storeErrs = append(storeErrs, err.Error())
		}

		if len(storeErrs) != 0 {

			log.Errorf("failed to call StoreLocalTask on taskManager with schedule task, err: %s",
				"\n["+strings.Join(storeErrs, ",")+"]")
			events := []*libtypes.TaskEvent{m.eventEngine.GenerateEvent(ev.TaskDiscarded.Type,
				task.GetTaskId(), task.GetTaskData().GetIdentityId(), task.GetTaskData().GetPartyId(),"store local task failed")}
			if err := m.storeBadTask(task, events, "store local task failed"); nil != err {
				log.WithError(err).Errorf("Failed to sending the task to datacenter on taskManager, taskId: {%s}", task.GetTaskId())
			}

		} else {

			taskArr = append(taskArr, task)
		}

	}

	// transfer `taskMsgs` to Scheduler
	go func() {
		m.sendLocalTaskToScheduler(taskArr)
	}()
	return nil
}

func (m *Manager) SendTaskTerminate(msgArr types.TaskTerminateMsgArr) error {
	for _, terminate := range msgArr {
		go m.TerminateTask(terminate)
	}
	return nil
}

func (m *Manager) SendTaskEvent(reportEvent *types.ReportTaskEvent) error {
	identityId, err := m.resourceMng.GetDB().QueryIdentityId()
	if nil != err {
		log.Errorf("Failed to query self identityId on taskManager.SendTaskEvent(), %s", err)
		return fmt.Errorf("query local identityId failed, %s", err)
	}
	reportEvent.GetEvent().IdentityId = identityId
	m.sendTaskEvent(reportEvent)
	return nil
}

func (m *Manager) SendTaskResourceUsage (usage *types.TaskResuorceUsage) error {

	task, ok :=m.queryNeedExecuteTaskCache(usage.GetTaskId(), usage.GetPartyId())
	if !ok {
		return fmt.Errorf("Can not find task cache, taskId: {%s}, partyId: {%s}", usage.GetPartyId(), usage.GetPartyId())
	}
	running, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusValExecByPartyId(usage.GetTaskId(), usage.GetPartyId())
	if nil != err {
		return err
	}
	if !running {
		return fmt.Errorf("task is not executed, taskId: {%s}, partyId: {%s}", usage.GetPartyId(), usage.GetPartyId())
	}

	m.sendTaskResourceUsageMsgToRemotePeer(task, usage)
	return nil
}

func (m *Manager) storeTaskHanlderPartyIds (task *types.Task, powerPartyIds []string) error {
	partyIds := make([]string, 0)

	for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
		partyIds = append(partyIds, dataSupplier.GetOrganization().GetPartyId())
	}
	partyIds = append(partyIds, powerPartyIds...)
	for _, receiver := range task.GetTaskData().GetReceivers() {
		partyIds = append(partyIds, receiver.GetPartyId())
	}
	if len(partyIds) == 0 {
		log.Warnf("Warn to store task handler partyIds, partyIds is empty, taskId: {%s}", task.GetTaskId())
		return nil
	}
	if err := m.resourceMng.GetDB().StoreTaskPartnerPartyIds(task.GetTaskId(), partyIds); nil != err {
		log.WithError(err).Errorf("Failed to store task handler partyIds, taskId: {%s}", task.GetTaskId())
		return err
	}
	return nil
}