package task

import (
	"github.com/RosettaFlow/Carrier-Go/consensus"
	"github.com/RosettaFlow/Carrier-Go/core"
	"github.com/RosettaFlow/Carrier-Go/event"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type Scheduler interface {
	SetTaskEngine(engine consensus.Consensus) error
	OnSchedule() error
	OnError () error
	SchedulerName() string
	PushTasks(tasks types.TaskMsgs) error
}

type Manager struct {
	eventCh   chan *event.TaskEvent
	dataChain *core.DataChain
	taskCh    <- chan types.TaskMsgs
	// Consensuses
	engines 		  map[string]consensus.Consensus
	//sendScheTaskCh chan chan
	scheduler Scheduler
}

func NewTaskManager(dataChain *core.DataChain, taskCh <- chan types.TaskMsgs, engines map[string]consensus.Consensus) *Manager {

	m := &Manager{
		eventCh:   make(chan *event.TaskEvent, 0),
		dataChain: dataChain,
		taskCh: taskCh,
		engines: engines,
	}
	go m.loop()
	return m
}

func (m *Manager) HandleSystemEvent(event *event.TaskEvent) error {
	return nil
}

func (m *Manager) HandleDataServiceEvent(event *event.TaskEvent) error {
	event, err := MakeDataServiceEventInfo(event)
	if nil != err {
		return err
	}
	return m.dataChain.StoreTaskEvent(event)
}

func (m *Manager) HandleComputerServiceEvent(event *event.TaskEvent) error {
	event, err := MakeComputerServiceEventInfo(event)
	if nil != err {
		return err
	}
	return m.dataChain.StoreTaskEvent(event)
}

func (m *Manager) HandleSchedulerEvent(event *event.TaskEvent) error {
	event, err := MakeScheduleEventInfo(event)
	if nil != err {
		return err
	}
	return m.dataChain.StoreTaskEvent(event)
}

func (m *Manager) HandleEvent(event *event.TaskEvent) error {
	eventType := event.Type
	if len(eventType) != 7 {
		return IncEventType
	}

	sysCode := eventType[0:2]
	switch sysCode {
	case "00":
		return m.HandleSystemEvent(event)
	case "01":
		return m.HandleSchedulerEvent(event)
	case "02":
		return m.HandleDataServiceEvent(event)
	case "03":
		return m.HandleComputerServiceEvent(event)
	default:
		return IncEventType
	}

	return nil
}

func (m *Manager) loop() {

	for {
		select {
		case event := <- m.eventCh:
			m.HandleEvent(event)
		case taskMsgs := <- m.taskCh:
			if len(taskMsgs) == 0 {
				continue
			}
			if err := m.scheduler.PushTasks(taskMsgs); nil != err {
				log.Error("Failed to push task msgs into Scheduler queue", "err", err)
			}

		default:
		}
	}
}

func (m *Manager) SendTaskEvent(event *event.TaskEvent) error {
	m.eventCh <- event
	return nil
}
