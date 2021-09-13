package task

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/consensus"
	ev "github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/core/schedule"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/types"

	"sync"
	"time"
)

const (
	defaultScheduleTaskInterval = 2 * time.Second
	taskMonitorInterval         = 30 * time.Second
	senderExecuteTaskExpire     = 6 * time.Second
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
	parser          *TaskParser
	validator       *TaskValidator
	// internal resource node set (Fighter node grpc client set)
	resourceClientSet *grpclient.InternalResourceClientSet

	eventCh chan *types.ReportTaskEvent
	quit    chan struct{}
	// send the validated taskMsgs to scheduler
	localTasksCh chan types.TaskDataArray
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
	resourceClientSet *grpclient.InternalResourceClientSet,
	localTasksCh chan types.TaskDataArray,
	needReplayScheduleTaskCh chan *types.NeedReplayScheduleTask,
	needExecuteTaskCh chan *types.NeedExecuteTask,
) *Manager {

	m := &Manager{
		p2p:               p2p,
		scheduler:         scheduler,
		consensusEngine:   consensusEngine,
		eventEngine:       eventEngine,
		resourceMng:       resourceMng,
		resourceClientSet: resourceClientSet,
		parser:            newTaskParser(),
		validator:         newTaskValidator(),
		eventCh:           make(chan *types.ReportTaskEvent, 10),
		localTasksCh:      localTasksCh,
		//needConsensusTaskCh:      needConsensusTaskCh,
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
		case <-taskTicker.C:

			if err := m.tryScheduleTask(); nil != err {
				log.Errorf("Failed to try schedule local task when taskTicker, err: {%s}", err)
			}

		// handle the task of need replay scheduling while received from remote peer on consensus epoch
		case needReplayScheduleTask := <- m.needReplayScheduleTaskCh:

			go func() {

				if err := m.resourceMng.GetDB().StoreLocalTask(needReplayScheduleTask.GetTask()); nil != err {

					log.Errorf("failed to call StoreLocalTask when replay schedule remote task, taskId: {%s}, err: {%s}", needReplayScheduleTask.GetTask().GetTaskId(), err)

					needReplayScheduleTask.SendFailedResult(needReplayScheduleTask.GetTask().GetTaskId(), err)

				} else {
					result := m.scheduler.ReplaySchedule(needReplayScheduleTask.GetLocalPartyId(), needReplayScheduleTask.GetLocalTaskRole(), needReplayScheduleTask.GetTask())
					needReplayScheduleTask.SendResult(result)
				}
			}()

		// handle task of need to executing, and send it to fighter of myself organization or send the task result msg to remote peer
		case task := <-m.needExecuteTaskCh:

			if task.GetConsStatus() == types.TaskNeedExecute {
				// add local cache first
				m.addNeedExecuteTaskCache(task)
				// to execute the task
				m.handleNeedExecuteTask(task)
			} else {
				// send task result msg to remote peer, short circuit.
				m.sendTaskResultMsgToRemotePeer(task)
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

func (m *Manager) SendTaskMsgArr(msgArr types.TaskMsgArr) error {
	if len(msgArr) == 0 {
		return fmt.Errorf("Receive some empty task msgArr")
	}

	nonParsedMsgArr, parsedMsgArr, err := m.parser.ParseTask(msgArr)
	if nil != err {
		for _, badMsg := range nonParsedMsgArr {
			events := []*libtypes.TaskEvent{m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
				badMsg.TaskId(), badMsg.OwnerIdentityId(), fmt.Sprintf("failed to parse local taskMsg"))}

			if e := m.storeBadTask(badMsg.Data, events, "failed to parse taskMsg"); nil != e {
				log.Errorf("Failed to store the err taskMsg on taskManager, taskId: {%s}", badMsg.TaskId())
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
				badMsg.TaskId(), badMsg.OwnerIdentityId(), fmt.Sprintf("failed to validate local taskMsg"))}

			if e := m.storeBadTask(badMsg.Data, events, "failed to validate taskMsg"); nil != e {
				log.Errorf("Failed to store the err taskMsg on taskManager, taskId: {%s}", badMsg.TaskId())
			}
		}

		if len(nonValidatedMsgArr) == len(parsedMsgArr) {
			return err
		}
	}

	taskArr := make(types.TaskDataArray, 0)

	for _, msg := range validatedMsgArr {
		task := msg.Data
		if err := m.resourceMng.GetDB().StoreLocalTask(task); nil != err {

			e := fmt.Errorf("store local task failed, taskId {%s}, %s", task.GetTaskData().TaskId, err)
			log.Errorf("failed to call StoreLocalTask on taskManager with schedule task, err: {%s}", e.Error())
			events := []*libtypes.TaskEvent{m.eventEngine.GenerateEvent(ev.TaskDiscarded.Type, task.GetTaskId(), task.GetTaskData().GetIdentityId(), e.Error())}
			if e := m.storeBadTask(task, events, e.Error()); nil != e {
				log.Errorf("Failed to sending the task to datacenter on taskManager, taskId: {%s}", task.GetTaskId())
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

func (m *Manager) SendTaskEvent(reportEvent *types.ReportTaskEvent) error {
	identityId, err := m.resourceMng.GetDB().GetIdentityId()
	if nil != err {
		log.Errorf("Failed to query self identityId on taskManager.SendTaskEvent(), %s", err)
		return fmt.Errorf("query local identityId failed, %s", err)
	}
	reportEvent.GetEvent().IdentityId = identityId
	m.sendTaskEvent(reportEvent)
	return nil
}
