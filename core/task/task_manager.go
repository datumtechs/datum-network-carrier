package task

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/auth"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
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
	eventCh           chan *libtypes.TaskEvent
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
		eventCh:                  make(chan *libtypes.TaskEvent, 10),
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
	taskTicker := time.NewTicker(defaultScheduleTaskInterval)

	for {
		select {

		// handle reported event from fighter of self organization
		case event := <-m.eventCh:
			go func() {
				if err := m.handleTaskEventWithCurrentIdentity(event); nil != err {
					log.WithError(err).Errorf("Failed to call handleTaskEventWithCurrentIdentity() on TaskManager, taskId: {%s}, event: %s", event.GetTaskId(), event.String())
				}
			}()

		// To schedule local task while received some local task msg
		case tasks := <-m.localTasksCh:

			for _, task := range tasks {
				if err := m.scheduler.AddTask(task); nil != err {
					log.WithError(err).Errorf("Failed to add local task into scheduler queue")
					continue
				}

				if err := m.tryScheduleTask(); nil != err {
					log.WithError(err).Errorf("Failed to try schedule local task while received local tasks")
					continue
				}

			}

		// To schedule local task interval
		case <-taskTicker.C:

			if err := m.tryScheduleTask(); nil != err {
				log.WithError(err).Errorf("Failed to try schedule local task when taskTicker")
			}

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
			//
			// NOTE:
			//	sender never received status `interrupt`, and partners never received status `scheduleFailed`.
			default:
				m.storeTaskFinalEvent(task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(),
					task.GetLocalTaskOrganization().GetPartyId(),
					fmt.Sprintf("execute task: %s with %s", task.GetConsStatus().String(),
						task.GetLocalTaskOrganization().GetPartyId()), apicommonpb.TaskState_TaskState_Failed)
				switch task.GetLocalTaskRole() {
				case apicommonpb.TaskRole_TaskRole_Sender:
					m.publishFinishedTaskToDataCenter(task, true)
				default:
					m.sendTaskResultMsgToTaskSender(task)
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

	// Why load 'local task' instead of 'needexecutetask'?
	//
	// Because the task may still be in the `consensus phase` rather than the `execution phase`,
	// the cached task cannot be found in the 'needexecutetaskcache' at this time.
	task, err := m.resourceMng.GetDB().QueryLocalTask(terminate.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task on `taskManager.TerminateTask()`, taskId: {%s}", terminate.GetTaskId())
		return
	}

	if nil == task {
		log.Errorf("Not found local task on `taskManager.TerminateTask()`, taskId: {%s}", terminate.GetTaskId())
		return
	}

	// The task sender only makes consensus, so interrupt consensus while need terminate task with task sender
	// While task is consensus or executing, can terminate.
	has, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusValConsByPartyId(task.GetTaskId(), task.GetTaskSender().GetPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task execute `cons` status on `taskManager.TerminateTask()`, taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetTaskSender().GetPartyId())
		return
	}
	if has {
		if err = m.consensusEngine.OnConsensusMsg(
			"", types.NewInterruptMsgWrap(task.GetTaskId(),
			types.MakeMsgOption(common.Hash{},
			apicommonpb.TaskRole_TaskRole_Sender,
			apicommonpb.TaskRole_TaskRole_Sender,
				task.GetTaskSender().GetPartyId(),
				task.GetTaskSender().GetPartyId(),
				task.GetTaskSender()))); nil != err {
			log.WithError(err).Errorf("Failed to call `OnConsensusMsg()` on `taskManager.TerminateTask()`, taskId: {%s}, partyId: {%s}",
				task.GetTaskId(), task.GetTaskSender().GetPartyId())
			return
		}
	}

	if err := m.onTerminateExecuteTask(task); nil != err {
		log.Errorf("Failed to call `onTerminateExecuteTask()` on `taskManager.TerminateTask()`, taskId: {%s}, err: \n%s", task.GetTaskId(), err)
	}

}

func (m *Manager) onTerminateExecuteTask(task *types.Task) error {
	// what if find the needExecuteTask(status: types.TaskConsensusFinished) with sender
	if m.hasNeedExecuteTaskCache(task.GetTaskId(), task.GetTaskSender().GetPartyId()) {

		// 1、 store task terminate (failed or succeed) event with current party
		m.resourceMng.GetDB().StoreTaskEvent(&libtypes.TaskEvent{
			Type:        ev.TaskTerminated.Type,
			TaskId:     task.GetTaskId(),
			IdentityId: task.GetTaskSender().GetIdentityId(),
			PartyId:    task.GetTaskSender().GetPartyId(),
			Content:    "task was terminated.",
			CreateAt:   timeutils.UnixMsecUint64(),
		})

		// 2、 remove needExecuteTask cache with sender
		m.removeNeedExecuteTaskCache(task.GetTaskId(), task.GetTaskSender().GetPartyId())
		// 3、 send a new needExecuteTask(status: types.TaskTerminate) for terminate with sender
		m.sendNeedExecuteTaskByAction(task,
			apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Sender,
			task.GetTaskSender(), task.GetTaskSender(),
			types.TaskTerminate)
	}

	return m.sendTaskTerminateMsg(task)
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

func (m *Manager) SendTaskEvent(event *libtypes.TaskEvent) error {
	m.sendTaskEvent(event)
	return nil
}

func (m *Manager) SendTaskResourceUsage (usage *types.TaskResuorceUsage) error {

	task, ok := m.queryNeedExecuteTaskCache(usage.GetTaskId(), usage.GetPartyId())
	if !ok {
		return fmt.Errorf("Can not find `need execute task` cache, taskId: {%s}, partyId: {%s}", usage.GetPartyId(), usage.GetPartyId())
	}
	running, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusValExecByPartyId(usage.GetTaskId(), usage.GetPartyId())
	if nil != err {
		return err
	}
	if !running {
		return fmt.Errorf("task is not executed, taskId: {%s}, partyId: {%s}", usage.GetPartyId(), usage.GetPartyId())
	}

	m.sendTaskResourceUsageMsgToTaskSender(task, usage)
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