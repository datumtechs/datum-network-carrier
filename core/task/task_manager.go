package task

import (
	"fmt"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
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
	sendTaskCh chan<- types.TaskMsgs
	// TODO 接收 被调度好的 task, 准备发给自己的  Fighter-Py
	recvSchedTaskCh      chan *types.ConsensusScheduleTaskWrap
	runningTaskCache     map[string]*types.ConsensusScheduleTask
	runningTaskCacheLock sync.RWMutex
}

func NewTaskManager(
	dataCenter core.CarrierDB,
	eventEngine *ev.EventEngine,
	resourceMng *resource.Manager,
	resourceClientSet *grpclient.InternalResourceClientSet,
	sendTaskCh chan types.TaskMsgs,
	recvSchedTaskCh chan *types.ConsensusScheduleTaskWrap,
) *Manager {

	m := &Manager{
		dataCenter:        dataCenter,
		eventEngine:       eventEngine,
		resourceMng:       resourceMng,
		resourceClientSet: resourceClientSet,
		parser:            newTaskParser(),
		validator:         newTaskValidator(),
		eventCh:           make(chan *types.TaskEventInfo, 10),
		sendTaskCh:        sendTaskCh,
		recvSchedTaskCh:   recvSchedTaskCh,
		runningTaskCache:  make(map[string]*types.ConsensusScheduleTask, 0),
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
			if task.TaskDir == ctypes.RecvTaskDir { //  需要 读出自己本地的 event 发给 task 的发起者
				eventList, err := m.dataCenter.GetTaskEventList(event.TaskId)
				if nil != err {
					log.Error("Failed to query all recv task event on myself", "taskId", event.TaskId, "err", err)
					return err
				}
				eventList = append(eventList, event)

			} else { //  如果是 自己的task, 认为任务终止 ... 发送到 dataCenter

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
		case task := <-m.recvSchedTaskCh:

			m.addRunningTaskCache(task.Task)

			switch task.SelfTaskRole {
			case types.TaskOnwer:
				switch task.Task.TaskState {
				case types.TaskStateFailed, types.TaskStateSuccess:

					// 判断是否 taskDir 决定是否直接 往 dataCenter 发送数据
					m.pulishFinishedTaskToDataCenter(task)

				case types.TaskStateRunning:

					if err := m.driveTaskForExecute(task.SelfTaskRole, task.Task); nil != err {
						log.Errorf("Failed to execute task on taskOnwer node, taskId: %s, %s", task.Task.SchedTask.TaskId, err)
						event := m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
							task.Task.SchedTask.TaskId, task.Task.SchedTask.Owner.IdentityId, fmt.Sprintf("failed to execute task"))
						// 因为是 自己的任务, 所以直接将 task  和 event list  发给 dataCenter
						m.dataCenter.StoreTaskEvent(event)
						m.pulishFinishedTaskToDataCenter(task)

					}
					// TODO 而执行最终[成功]的 根据 Fighter 上报的 event 在 handleEvent() 里面处理
				default:
					log.Error("Failed to handle unknown task", "taskId", task.Task.SchedTask.TaskId)
				}
			//case types.DataSupplier:
			//case types.PowerSupplier:
			//case types.ResultSupplier:
			default:
				switch task.Task.TaskState {
				case types.TaskStateFailed, types.TaskStateSuccess:
					// 因为是 task 参与者, 所以需要构造 taskResult 发送给 task 发起者..
					m.sendTaskResultMsgToConsensus(task)
				case types.TaskStateRunning:

					if err := m.driveTaskForExecute(task.SelfTaskRole, task.Task); nil != err {
						log.Errorf("Failed to execute task on taskOnwer node, taskId: %s, %s", task.Task.SchedTask.TaskId, err)
						identityId, _ := m.dataCenter.GetIdentityId()
						event := m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
							task.Task.SchedTask.TaskId, identityId, fmt.Sprintf("failed to execute task"))

						// 因为是 task 参与者, 所以需要构造 taskResult 发送给 task 发起者..
						m.dataCenter.StoreTaskEvent(event)
						m.sendTaskResultMsgToConsensus(task)
					}
				default:
					log.Error("Failed to handle unknown task", "taskId", task.Task.SchedTask.TaskId)
				}

			}

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
