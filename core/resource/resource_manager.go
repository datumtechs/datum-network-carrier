package resource

import "github.com/RosettaFlow/Carrier-Go/event"

type Manager struct {

	eventCh chan *event.TaskEvent

}

func NewResourceManager() *Manager {

	m := &Manager{
		eventCh: make(chan *event.TaskEvent, 0),
	}
	go m.loop()
	return m
}

func (m *Manager) loop () {

	for {
		select {
		case event := <- m.eventCh:
			_ = event  // TODO add some logic about eventEngine
		default:
		}
	}
}

func (m *Manager) SendTaskEvent(event *event.TaskEvent) error {
	m.eventCh <- event
	return nil
}