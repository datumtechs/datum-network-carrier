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
	// 接收 被调度好的 task, 准备发给自己的  Fighter-Py 或者 发给 dataCenter
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
	return m
}

func (m *Manager) Start() error {
	go m.loop()
	log.Info("Started taskManager ...")
	return nil
}
func (m *Manager)Stop() error { return nil }

func (m *Manager) loop() {

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

			if e := m.storeErrTaskMsg(errtask, types.ConvertTaskEventArrToDataCenter(events), "failed to parse taskMsg"); nil != e {
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

			if e := m.storeErrTaskMsg(errtask, types.ConvertTaskEventArrToDataCenter(events), "failed to validate taskMsg"); nil != e {
				log.Error("Failed to store the err taskMsg on taskManager", "taskId", errtask.TaskId)
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
		log.Errorf("Failed to query self identityId on taskManager.SendTaskEvent(), %s", err)
		return fmt.Errorf("query local identityId failed, %s", err)
	}
	event.Identity = identityId
	m.sendTaskEvent(event)
	return nil
}
