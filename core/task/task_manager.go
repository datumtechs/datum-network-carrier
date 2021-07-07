package task

import (
	"github.com/RosettaFlow/Carrier-Go/core"
	ev "github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/types"
)

//type Scheduler interface {
//	OnStart() error
//	OnStop() error
//	OnError () error
//	Name() string
//}

type Manager struct {
	eventCh   chan *types.TaskEventInfo
	dataChain *core.DataChain
	eventEngine *ev.EventEngine
	resourceMng *resource.Manager
	// recv the taskMsgs from messageHandler
	taskCh <-chan types.TaskMsgs
	// send the validated taskMsgs to scheduler
	sendTaskCh chan<- types.TaskMsgs
	// TODO 接收 被调度好的 task, 准备发给自己的  Fighter-Py
	recvSchedTaskCh <-chan *types.ConsensusScheduleTask

	// TODO 持有 己方的所有 Fighter-Py 的 grpc client

	// TODO 用于接收 己方已连接 或 断开连接的 Fighter-Py 的 grpc client


}

func NewTaskManager(dataChain *core.DataChain, eventEngine *ev.EventEngine,
	resourceMng *resource.Manager,
	taskCh <-chan types.TaskMsgs, sendTaskCh chan<- types.TaskMsgs,
	recvSchedTaskCh <-chan *types.ConsensusScheduleTask) *Manager {

	m := &Manager{
		eventCh:   make(chan *types.TaskEventInfo, 10),
		dataChain: dataChain,
		eventEngine: eventEngine,
		resourceMng: resourceMng,
		taskCh:    taskCh,
		sendTaskCh: sendTaskCh,
		recvSchedTaskCh: recvSchedTaskCh,

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

	return m.eventEngine.StoreEvent(event)
}

func (m *Manager) loop() {

	for {
		select {
		case event := <-m.eventCh:
			m.handleEvent(event)
		case taskMsgs := <-m.taskCh:
			if len(taskMsgs) == 0 {
				continue
			}
			// TODO 先对 task 做校验 和解析

			// TODO 再 转发给 Scheduler 处理

			case task := <- m.recvSchedTaskCh:
				// TODO 对接收到 经 Scheduler  调度好的 task  转发给自己的 Fighter-Py
			_=task

		default:
		}
	}
}

func (m *Manager) SendTaskEvent(event *types.TaskEventInfo) error {
	m.eventCh <- event
	return nil
}
