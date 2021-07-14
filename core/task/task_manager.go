package task

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/core"
	ev "github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	"github.com/RosettaFlow/Carrier-Go/types"
	"sync"
)

type Manager struct {
	dataCenter  core.CarrierDB
	eventEngine *ev.EventEngine
	resourceMng *resource.Manager
	parser      *TaskParser
	validator   *TaskValidator
	// internal resource node set (Fighter node grpc client set)
	resourceClientSet *grpclient.InternalResourceClientSet

	eventCh chan *types.TaskEventInfo
	// send the validated taskMsgs to scheduler
	localTaskMsgCh chan<- types.TaskMsgs
	// TODO 接收 被调度好的 task, 准备发给自己的  Fighter-Py
	doneScheduleTaskCh   chan *types.DoneScheduleTaskChWrap
	runningTaskCache     map[string]*types.DoneScheduleTaskChWrap
	runningTaskCacheLock sync.RWMutex
}

func NewTaskManager(
	dataCenter core.CarrierDB,
	eventEngine *ev.EventEngine,
	resourceMng *resource.Manager,
	resourceClientSet *grpclient.InternalResourceClientSet,
	localTaskMsgCh chan types.TaskMsgs,
	doneScheduleTaskCh chan *types.DoneScheduleTaskChWrap,
) *Manager {

	m := &Manager{
		dataCenter:         dataCenter,
		eventEngine:        eventEngine,
		resourceMng:        resourceMng,
		resourceClientSet:  resourceClientSet,
		parser:             newTaskParser(),
		validator:          newTaskValidator(),
		eventCh:            make(chan *types.TaskEventInfo, 10),
		localTaskMsgCh:     localTaskMsgCh,
		doneScheduleTaskCh: doneScheduleTaskCh,
		runningTaskCache:   make(map[string]*types.DoneScheduleTaskChWrap, 0),
	}
	go m.loop()
	return m
}

func (m *Manager) handleEvent(event *types.TaskEventInfo) error {
	eventType := event.Type
	if len(eventType) != ev.EventTypeCharLen {
		return ev.IncEventType
	}
	// TODO need to validate the task that have been processing ? Maybe~
	if event.Type == ev.TaskExecuteEOF.Type {
		if task, ok := m.queryRunningTaskCacheOk(event.TaskId); ok {
			defer func() {
				m.removeRunningTaskCache(event.TaskId)
			}()

			if task.Task.TaskDir == types.RecvTaskDir {
				// 因为是 task 参与者, 所以需要构造 taskResult 发送给 task 发起者..
				m.dataCenter.StoreTaskEvent(event)
				m.sendTaskResultMsgToConsensus(event.TaskId)

			} else {
				//  如果是 自己的task, 认为任务终止 ... 发送到 dataCenter
				m.pulishFinishedTaskToDataCenter(event.TaskId)
			}
		}
		return nil
	} else {
		return m.eventEngine.StoreEvent(event)
	}
}

func (m *Manager) loop() {

	for {
		select {
		// 自己组织的 Fighter 上报过来的 event
		case event := <-m.eventCh:
			if err := m.handleEvent(event); nil != err {
				log.Error("Failed to store task event on local", "taskId", event.TaskId, "event", event.String())
			}

		// 接收 被调度好的 task, 准备发给自己的  Fighter-Py 或者直接存到 dataCenter
		case task := <-m.doneScheduleTaskCh:

			// 添加本地缓存
			m.addRunningTaskCache(task)
			m.handleDoneScheduleTask(task.Task.SchedTask.TaskId)

		default:
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
				errtask.TaskId, errtask.Onwer().IdentityId, fmt.Sprintf("failed to parse taskMsg")))

			if e := m.storeErrTaskMsg(errtask, types.ConvertTaskEventArrToDataCenter(events), "failed to parse taskMsg"); nil != e {
				log.Error("Failed to store the err taskMsg", "taskId", errtask)
			}
		}
		return err
	}

	if errTasks, err := m.validator.validateTaskMsg(msgs); nil != err {
		for _, errtask := range errTasks {
			events, _ := m.dataCenter.GetTaskEventList(errtask.TaskId)
			events = append(events, m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
				errtask.TaskId, errtask.Onwer().IdentityId, fmt.Sprintf("failed to validate taskMsg")))

			if e := m.storeErrTaskMsg(errtask, types.ConvertTaskEventArrToDataCenter(events), "failed to validate taskMsg"); nil != e {
				log.Error("Failed to store the err taskMsg", "taskId", errtask)
			}
		}
		return err
	}
	// transfer `taskMsgs` to Scheduler
	go func() {
		m.sendTaskMsgsToScheduler(msgs)
	}()
	return nil
}

func (m *Manager) SendTaskEvent(event *types.TaskEventInfo) error {
	identityId, err := m.dataCenter.GetIdentityId()
	if nil != err {
		log.Errorf("Failed to query self identityId on SendTaskEvent, %s", err)
		return err
	}
	event.Identity = identityId
	m.sendTaskEvent(event)
	return nil
}
