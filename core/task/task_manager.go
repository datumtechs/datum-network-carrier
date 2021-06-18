package task

import (
	"github.com/RosettaFlow/Carrier-Go/event"
	"github.com/RosettaFlow/carrier-go/core"
)

type Manager struct {
	eventCh   chan *event.TaskEvent
	dataChain *core.DataChain
}

func NewTaskManager(dataChain *core.DataChain) *Manager {

	m := &Manager{
		eventCh:   make(chan *event.TaskEvent, 0),
		dataChain: dataChain,
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
		case event := <-m.eventCh:
			m.HandleEvent(event)
		default:
		}
	}
}

func (m *Manager) SendTaskEvent(event *event.TaskEvent) error {
	m.eventCh <- event
	return nil
}
