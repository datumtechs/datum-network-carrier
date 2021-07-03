package resource

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/event"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type Manager struct {
	eventCh    chan *event.TaskEvent

	slotUnit   *types.Slot
	tables     map[string]*types.ResourceTable
	tableQueue []*types.ResourceTable
}

func NewResourceManager() *Manager {
	m := &Manager{
		eventCh: make(chan *event.TaskEvent, 0),
		tables:  make(map[string]*types.ResourceTable),
	}
	go m.loop()
	return m
}

func (m *Manager) loop() {

	for {
		select {
		case event := <-m.eventCh:
			_ = event // TODO add some logic about eventEngine
		default:
		}
	}
}

func (m *Manager) SetSlotUnit(mem, p, b uint64) {
	m.slotUnit = &types.Slot{
		Mem:       mem,
		Processor: p,
		Bandwidth: b,
	}
	if len(m.tables) != 0 {
		for _, re := range m.tables {
			re.SetSlotUnit(m.slotUnit)
		}
	}
}

func (m *Manager) UseSlot(nodeId string, slotCount uint32) error {
	table, ok := m.tables[nodeId]
	if !ok {
		return fmt.Errorf("No found the resource table of node: %s", nodeId)
	}
	if table.RemianSlot() < slotCount {
		return fmt.Errorf("Insufficient remaining number of slots of node: %s", nodeId)
	}
	table.UseSlot(slotCount)
	return nil
}

func (m *Manager) SetResource(table *types.ResourceTable) {
	m.tables[table.GetNodeId()] = table
	m.tableQueue = append(m.tableQueue, table)
}
func (m *Manager) GetResource(nodeId string) *types.ResourceTable { return m.tables[nodeId] }
func (m *Manager) GetResources() []*types.ResourceTable           { return m.tableQueue }
func (m *Manager) DelResource(nodeId string) {
	for i := 0; i < len(m.tableQueue); i++ {
		table := m.tableQueue[i]
		if table.GetNodeId() == nodeId {
			delete(m.tables, nodeId)
			m.tableQueue = append(m.tableQueue[:i], m.tableQueue[i+1:]...)
			i--
		}
	}
}

func (m *Manager) SendTaskEvent(event *event.TaskEvent) error {
	m.eventCh <- event
	return nil
}
