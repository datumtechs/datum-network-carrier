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
	localTaskMsgCh chan<- types.TaskMsgs
	// 接收 被调度好的 task, 准备发给自己的  Fighter-Py 或者 发给 dataCenter
	doneScheduleTaskCh   chan *types.DoneScheduleTaskChWrap
	runningTaskCache     map[string]*types.DoneScheduleTaskChWrap
	runningTaskCacheLock sync.RWMutex
}

func NewTaskManager(
	scheduler Scheduler,
	eventEngine *ev.EventEngine,
	resourceMng *resource.Manager,
	resourceClientSet *grpclient.InternalResourceClientSet,
	localTaskMsgCh chan types.TaskMsgs,
	doneScheduleTaskCh chan *types.DoneScheduleTaskChWrap,
) *Manager {

	m := &Manager{
		scheduler:          scheduler,
		eventEngine:        eventEngine,
		resourceMng:        resourceMng,
		resourceClientSet:  resourceClientSet,
		parser:             newTaskParser(),
		validator:          newTaskValidator(),
		eventCh:            make(chan *libTypes.TaskEvent, 10),
		localTaskMsgCh:     localTaskMsgCh,
		doneScheduleTaskCh: doneScheduleTaskCh,
		runningTaskCache:   make(map[string]*types.DoneScheduleTaskChWrap, 0),
		quit:               make(chan struct{}),
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
		case task := <-m.doneScheduleTaskCh:

			// 添加本地缓存
			m.addRunningTaskCache(task)
			m.handleDoneScheduleTask(task.Task.SchedTask.TaskId())

		case <-taskMonitorTicker.C:
			m.expireTaskMonitor()

			// 定时调度 队列中的任务信息
		case <-taskTicker.C:
			m.scheduler.TrySchedule()

		case <-m.quit:
			log.Info("Stopped taskManager ...")
			return
		}
	}
}

func (m *Manager) SendTaskMsgs(msgs types.TaskMsgs) error {
	if len(msgs) == 0 {
		return fmt.Errorf("Receive some empty task msgs")
	}



	if errTasks, err := m.parser.ParseTask(msgs); nil != err {
		for _, errtask := range errTasks {

			events, _ := m.dataCenter.GetTaskEventList(errtask.TaskId)
			events = append(events, m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
				errtask.TaskId, errtask.OwnerIdentityId(), fmt.Sprintf("failed to parse taskMsg")))

			if e := m.storeErrTaskMsg(errtask, events, "failed to parse taskMsg"); nil != e {
				log.Error("Failed to store the err taskMsg on taskManager", "taskId", errtask.TaskId)
			}
		}
		return err
	}

	if errTasks, err := m.validator.validateTaskMsg(msgs); nil != err {
		for _, errtask := range errTasks {
			events, _ := m.dataCenter.GetTaskEventList(errtask.TaskId)
			events = append(events, m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
				errtask.TaskId, errtask.OwnerIdentityId(), fmt.Sprintf("failed to validate taskMsg")))

			if e := m.storeErrTaskMsg(errtask, events, "failed to validate taskMsg"); nil != e {
				log.Error("Failed to store the err taskMsg on taskManager", "taskId", errtask.TaskId)
			}
		}
		return err
	}

	for _, msg := range msgs {
		task := msg.Data
		if err := m.resourceMng.GetDB().StoreLocalTask(task); nil != err {
			e := fmt.Errorf("store local task failed, taskId {%s}, %s", task.TaskData().TaskId, err)

			log.Errorf("failed to call StoreLocalTask on SchedulerStarveFIFO with schedule task, err: {%s}", e.Error())

			task.TaskData().EndAt = uint64(timeutils.UnixMsec())
			task.TaskData().Reason = e.Error()
			task.TaskData().State = apipb.TaskState_TaskState_Failed

			identityId, _ := sche.dataCenter.GetIdentityId()
			event := sche.eventEngine.GenerateEvent(evengine.TaskDiscarded.Type, task.Data.TaskData().TaskId, identityId, e.Error())
			task.Data.TaskData().EventCount = 1
			task.Data.TaskData().TaskEventList = []*libTypes.TaskEvent{event}

			if err = sche.dataCenter.InsertTask(task.Data); nil != err {
				log.Errorf("Failed to save task to datacenter, taskId: {%s}", task.Data.TaskData().TaskId)
				continue
			}
		}
	}


	// transfer `taskMsgs` to Scheduler
	go func() {
		m.sendTaskMsgsToScheduler(msgs)
	}()
	return nil
}

func (m *Manager) SendTaskEvent(event *libTypes.TaskEvent) error {
	identityId, err := m.dataCenter.GetIdentityId()
	if nil != err {
		log.Errorf("Failed to query self identityId on taskManager.SendTaskEvent(), %s", err)
		return fmt.Errorf("query local identityId failed, %s", err)
	}
	event.IdentityId = identityId
	m.sendTaskEvent(event)
	return nil
}
