package task

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/auth"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/bytesutil"
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
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"
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
func (m *Manager) recoveryNeedExecuteTask() {
	prefix := rawdb.GetNeedExecuteTaskKeyPrefix()
	if err := m.resourceMng.GetDB().ForEachNeedExecuteTask(func(key, value []byte) error {
		if len(key) != 0 && len(value) != 0 {

			// task:${taskId hex} == 5 + 2 + 64 == "taskId:" + "0x" + "e33...fe4"
			taskId := string(key[len(prefix) : len(prefix)+71])
			partyId := string(key[len(prefix)+71:])

			//task, err := m.resourceMng.GetDB().QueryLocalTask(taskId)
			//if nil != err {
			//	return fmt.Errorf("query local task failed on recover needExecuteTask from db, %s, taskId: {%s}", err, taskId)
			//}

			var res libtypes.NeedExecuteTask

			if err := proto.Unmarshal(value, &res); nil != err {
				return fmt.Errorf("Unmarshal needExecuteTask failed, %s", err)
			}

			cache, ok := m.runningTaskCache[taskId]
			if !ok {
				cache = make(map[string]*types.NeedExecuteTask, 0)
			}
			cache[partyId] = types.NewNeedExecuteTask(
				peer.ID(res.GetRemotePid()),
				common.HexToHash(res.GetProposalId()),
				res.GetLocalTaskRole(),
				res.GetRemoteTaskRole(),
				res.GetLocalTaskOrganization(),
				res.GetRemoteTaskOrganization(),
				taskId,
				types.TaskActionStatus(bytesutil.BytesToUint16(res.GetConsStatus())),
				types.NewPrepareVoteResource(
					res.GetLocalResource().GetId(),
					res.GetLocalResource().GetIp(),
					res.GetLocalResource().GetPort(),
					res.GetLocalResource().GetPartyId(),
				),
				res.GetResources(),
			)
			m.runningTaskCache[taskId] = cache
		}
		return nil
	}); nil != err {
		log.WithError(err).Fatalf("recover needExecuteTask failed")
	}
}
func (m *Manager) Start() error {
	m.recoveryNeedExecuteTask()
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

			log.Infof("Received local task arr on taskManager.loop(), task arr len: %d", len(tasks))

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

				// There is no need to store local tasks repeatedly
				if rawdb.IsDBNotFoundErr(err) {

					log.Infof("Start to store local task on taskManager.loop() when received needReplayScheduleTask, taskId: {%s}", needReplayScheduleTask.GetTask().GetTaskId())

					// store metadata used taskId
					if err := m.storeMetaUsedTaskId(needReplayScheduleTask.GetTask()); nil != err {
						log.WithError(err).Errorf("Failed to store metadata used taskId when received remote task, taskId: {%s}", needReplayScheduleTask.GetTask().GetTaskId())
					}
					if err := m.resourceMng.GetDB().StoreLocalTask(needReplayScheduleTask.GetTask()); nil != err {
						log.WithError(err).Errorf("Failed to call StoreLocalTask when replay schedule remote task, taskId: {%s}", needReplayScheduleTask.GetTask().GetTaskId())
						needReplayScheduleTask.SendFailedResult(needReplayScheduleTask.GetTask().GetTaskId(), err)
						return
					}
					log.Infof("Finished to store local task on taskManager.loop() when received needReplayScheduleTask, taskId: {%s}", needReplayScheduleTask.GetTask().GetTaskId())

				}

				// Start replay schedule remote task ...
				result := m.scheduler.ReplaySchedule(needReplayScheduleTask.GetLocalPartyId(), needReplayScheduleTask.GetLocalTaskRole(), needReplayScheduleTask.GetTask())
				needReplayScheduleTask.SendResult(result)
			}()

		// handle task of need to executing, and send it to fighter of myself organization or send the task result msg to remote peer
		case task := <-m.needExecuteTaskCh:
			
			localTask, err := m.resourceMng.GetDB().QueryLocalTask(task.GetTaskId())
			if nil != err {
				log.WithError(err).Errorf("Failed to query local task info on taskManager.loop() when received needExecuteTask, taskId: {%s}, partyId: {%s}, status: {%s}",
					task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), task.GetConsStatus().String())
				return
			}

			switch task.GetConsStatus() {
			// sender and partner to handle needExecuteTask when consensus succeed.
			// sender need to store some cache, partner need to execute task.
			case types.TaskNeedExecute, types.TaskConsensusFinished:
				// to execute the task
				m.handleNeedExecuteTask(task, localTask)

			// sender and partner to clear local task things after received status: `scheduleFailed` and `interrupt` and `terminate`.
			// sender need to publish local task and event to datacenter,
			// partner need to send task's event to remote task's sender.
			//
			// NOTE:
			//	sender never received status `interrupt`, and partners never received status `scheduleFailed`.
			default:

				// store a bad event into local db before handle bad task.
				m.storeTaskFinalEvent(task.GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(),
					task.GetLocalTaskOrganization().GetPartyId(),
					fmt.Sprintf("execute task: %s with %s", task.GetConsStatus().String(),
						task.GetLocalTaskOrganization().GetPartyId()), apicommonpb.TaskState_TaskState_Failed)

				switch task.GetLocalTaskRole() {
				case apicommonpb.TaskRole_TaskRole_Sender:
					m.publishFinishedTaskToDataCenter(task, localTask, true)
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

func (m *Manager) TerminateTask(terminate *types.TaskTerminateMsg) {

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
			Type:       ev.TaskTerminated.Type,
			TaskId:     task.GetTaskId(),
			IdentityId: task.GetTaskSender().GetIdentityId(),
			PartyId:    task.GetTaskSender().GetPartyId(),
			Content:    "task was terminated.",
			CreateAt:   timeutils.UnixMsecUint64(),
		})

		// 2、 remove needExecuteTask cache with sender
		m.removeNeedExecuteTaskCache(task.GetTaskId(), task.GetTaskSender().GetPartyId())
		// 3、 send a new needExecuteTask(status: types.TaskTerminate) for terminate with sender
		m.sendNeedExecuteTaskByAction(task.GetTaskId(),
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

	nonParsedMsgArr, parsedMsgArr := m.parser.ParseTask(msgArr)
	for _, badMsg := range nonParsedMsgArr {

		go m.resourceMng.GetDB().RemoveTaskMsg(badMsg.GetTaskMsg().GetTaskId()) // remove from disk if task been non-parsed task
		events := []*libtypes.TaskEvent{m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
			badMsg.GetTaskMsg().GetTaskId(), badMsg.GetTaskMsg().GetSenderIdentityId(), badMsg.GetTaskMsg().GetSenderPartyId(), badMsg.GetErrStr())}

		if e := m.storeBadTask(badMsg.GetTaskMsg().GetTask(), events, "failed to parse taskMsg"); nil != e {
			log.WithError(e).Errorf("Failed to store the err taskMsg on taskManager.SendTaskMsgArr(), taskId: {%s}", badMsg.GetTaskMsg().GetTaskId())
		}
	}

	nonValidatedMsgArr, validatedMsgArr := m.validator.validateTaskMsg(parsedMsgArr)
	for _, badMsg := range nonValidatedMsgArr {
		go m.resourceMng.GetDB().RemoveTaskMsg(badMsg.GetTaskMsg().GetTaskId()) // remove from disk if task been non-validated task
		events := []*libtypes.TaskEvent{m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
			badMsg.GetTaskMsg().GetTaskId(), badMsg.GetTaskMsg().GetSenderIdentityId(), badMsg.GetTaskMsg().GetSenderPartyId(), badMsg.GetErrStr())}

		if e := m.storeBadTask(badMsg.GetTaskMsg().GetTask(), events, "failed to validate taskMsg"); nil != e {
			log.WithError(e).Errorf("Failed to store the err taskMsg on taskManager.SendTaskMsgArr(), taskId: {%s}", badMsg.GetTaskMsg().GetTaskId())
		}
	}

	taskArr := make(types.TaskDataArray, 0)

	for _, msg := range validatedMsgArr {

		log.Infof("Start to store local task on taskManager.SendTaskMsgArr(), taskId: {%s}", msg.GetTaskId())

		go m.resourceMng.GetDB().RemoveTaskMsg(msg.GetTaskId()) // remove from disk if task been handle task

		task := msg.Data

		// store metadata used taskId
		if err := m.storeMetaUsedTaskId(task); nil != err {
			log.WithError(err).Errorf("Failed to store metadata used taskId when received local task on taskManager.SendTaskMsgArr(), taskId: {%s}", task.GetTaskId())
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

			log.Errorf("failed to call StoreLocalTask on taskManager with schedule task on taskManager.SendTaskMsgArr(), err: %s",
				"\n["+strings.Join(storeErrs, ",")+"]")
			events := []*libtypes.TaskEvent{m.eventEngine.GenerateEvent(ev.TaskDiscarded.Type,
				task.GetTaskId(), task.GetTaskData().GetIdentityId(), task.GetTaskData().GetPartyId(), "store local task failed")}
			if err := m.storeBadTask(task, events, "store local task failed"); nil != err {
				log.WithError(err).Errorf("Failed to sending the task to datacenter on taskManager.SendTaskMsgArr(), taskId: {%s}", task.GetTaskId())
			}

		} else {

			taskArr = append(taskArr, task)
			log.Infof("Finished to store local task on taskManager.SendTaskMsgArr(), taskId: {%s}", msg.GetTaskId())
		}

	}

	// transfer `taskMsgs` to Scheduler
	log.Infof("Need send to scheduler tasks count on taskManager.SendTaskMsgArr(), task arr len: %d", len(taskArr))
	m.sendLocalTaskToScheduler(taskArr)
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

func (m *Manager) SendTaskResourceUsageToTaskSender(usage *types.TaskResuorceUsage) error {

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

func (m *Manager) storeTaskHanlderPartyIds(task *types.Task, powerPartyIds []string) error {
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
