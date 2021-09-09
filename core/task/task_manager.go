package task

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	ev "github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"sync"
	"time"
)

const (
	defaultScheduleTaskInterval = 2 * time.Second
	taskMonitorInterval         = 30 * time.Second
)

type Scheduler interface {
	Start() error
	Stop() error
	Error() error
	Name() string
	AddTask(task *types.Task)
	RemoveTask(taskId string) error
	TrySchedule() (*types.NeedConsensusTask, error)
	ReplaySchedule(myPartyId string, myTaskRole apipb.TaskRole, task *types.Task) *types.ReplayScheduleResult
}

type Manager struct {
	scheduler   Scheduler
	eventEngine *ev.EventEngine
	resourceMng *resource.Manager
	parser      *TaskParser
	validator   *TaskValidator
	// internal resource node set (Fighter node grpc client set)
	resourceClientSet *grpclient.InternalResourceClientSet

	eventCh chan *libTypes.TaskEvent
	quit    chan struct{}
	// send the validated taskMsgs to scheduler
	localTasksCh             chan types.TaskDataArray
	needConsensusTaskCh      chan *types.NeedConsensusTask
	needReplayScheduleTaskCh chan *types.NeedReplayScheduleTask
	needExecuteTaskCh        chan *types.NeedExecuteTask
	runningTaskCache         map[string]map[string]*types.NeedExecuteTask //  taskId -> {partyId -> task}
	runningTaskCacheLock     sync.RWMutex
}

func NewTaskManager(
	scheduler Scheduler,
	eventEngine *ev.EventEngine,
	resourceMng *resource.Manager,
	resourceClientSet *grpclient.InternalResourceClientSet,
	localTasksCh chan types.TaskDataArray,
	needConsensusTaskCh chan *types.NeedConsensusTask,
	needReplayScheduleTaskCh chan *types.NeedReplayScheduleTask,
	needExecuteTaskCh chan *types.NeedExecuteTask,
) *Manager {

	m := &Manager{
		scheduler:                scheduler,
		eventEngine:              eventEngine,
		resourceMng:              resourceMng,
		resourceClientSet:        resourceClientSet,
		parser:                   newTaskParser(),
		validator:                newTaskValidator(),
		eventCh:                  make(chan *libTypes.TaskEvent, 10),
		localTasksCh:             localTasksCh,
		needConsensusTaskCh:      needConsensusTaskCh,
		needReplayScheduleTaskCh: needReplayScheduleTaskCh,
		needExecuteTaskCh:        needExecuteTaskCh,
		runningTaskCache:         make(map[string]map[string]*types.NeedExecuteTask, 0),
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
		// 自己组织的 Fighter 上报过来的 event
		case event := <-m.eventCh:
			go func() {
				if err := m.handleEvent(event); nil != err {
					log.Error("Failed to call handleEvent() on TaskManager", "taskId", event.TaskId, "event", event.String())
				}
			}()

		// 接收 被调度好的 task, 准备发给自己的  Fighter-Py 或者直接存到 dataCenter
		case task := <-m.needExecuteTaskCh:

			// 添加本地缓存
			m.addNeedExecuteTaskCache(task)
			m.handleNeedExecuteTask(task)

		case <-taskMonitorTicker.C:
			m.expireTaskMonitor()

		case tasks := <-m.localTasksCh:

			for _, task := range tasks {
				m.scheduler.AddTask(task)
				m.scheduler.TrySchedule()
			}
			// 定时调度 队列中的任务信息
		case <-taskTicker.C:
			m.scheduler.TrySchedule()

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
			events := []*libTypes.TaskEvent{m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
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
			events := []*libTypes.TaskEvent{m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
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
			log.Errorf("failed to call StoreLocalTask on SchedulerStarveFIFO with schedule task, err: {%s}", e.Error())

			events := []*libTypes.TaskEvent{m.eventEngine.GenerateEvent(ev.TaskDiscarded.Type, task.GetTaskId(), task.GetTaskData().GetIdentityId(), e.Error())}

			task.GetTaskData().EndAt = uint64(timeutils.UnixMsec())
			task.GetTaskData().Reason = e.Error()
			task.GetTaskData().State = apipb.TaskState_TaskState_Failed
			task.GetTaskData().EventCount = uint32(len(events))
			task.GetTaskData().TaskEvents = events

			if err = m.resourceMng.GetDB().InsertTask(task); nil != err {
				log.Errorf("Failed to save task to datacenter, taskId: {%s}", task.GetTaskId())
				continue
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

func (m *Manager) SendTaskEvent(event *libTypes.TaskEvent) error {
	identityId, err := m.resourceMng.GetDB().GetIdentityId()
	if nil != err {
		log.Errorf("Failed to query self identityId on taskManager.SendTaskEvent(), %s", err)
		return fmt.Errorf("query local identityId failed, %s", err)
	}
	event.IdentityId = identityId
	m.sendTaskEvent(event)
	return nil
}
