package resource

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/db"
	"github.com/RosettaFlow/Carrier-Go/event"
	"github.com/RosettaFlow/Carrier-Go/types"
)

//type DataCenter interface {
//	GetResourceList() (types.ResourceArray, error)
//}

type Manager struct {
	// TODO 这里需要一个 config <SlotUnit 的>

	db              db.Database // Low level persistent database to store final content.
	eventCh         chan *event.TaskEvent
	slotUnit        *types.Slot
	localTables     map[string]*types.LocalResourceTable
	localTableQueue []*types.LocalResourceTable
	//remoteTables     map[string]*types.RemoteResourceTable
	remoteTableQueue []*types.RemoteResourceTable
}

func NewResourceManager(db db.Database) *Manager {
	m := &Manager{
		db:               db,
		eventCh:          make(chan *event.TaskEvent, 0),
		localTables:      make(map[string]*types.LocalResourceTable),
		localTableQueue:  make([]*types.LocalResourceTable, 0),
		remoteTableQueue: make([]*types.RemoteResourceTable, 0),
		slotUnit:         types.DefaultSlotUnit, // TODO for test
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

func (m *Manager) Start() error {
	m.SetSlotUnit(0, 0, 0)
	// load slotUnit
	slotUnit, err := queryNodeResourceSlotUnit(m.db)
	if nil != err {
		return err
	}
	m.slotUnit = slotUnit
	// load resource localTables
	resources, err := queryNodeResources(m.db)
	if nil != err {
		return err
	}
	tables := make(map[string]*types.LocalResourceTable, len(resources))
	for _, resource := range resources {
		tables[resource.GetNodeId()] = resource
	}
	m.localTables = tables
	m.localTableQueue = resources

	// load resource remoteTables

	return nil
}

func (m *Manager) Stop() error {
	// store slotUnit
	if err := storeNodeResourceSlotUnit(m.db, m.slotUnit); nil != err {
		return err
	}
	// store resource localTables
	if err := storeNodeResources(m.db, m.localTableQueue); nil != err {
		return err
	}
	return nil
}

func (m *Manager) SetSlotUnit(mem, p, b uint64) {
	//m.slotUnit = &types.Slot{
	//	Mem:       mem,
	//	Processor: p,
	//	Bandwidth: b,
	//}
	m.slotUnit = types.DefaultSlotUnit // TODO for test
	if len(m.localTables) != 0 {
		for _, re := range m.localTables {
			re.SetSlotUnit(m.slotUnit)
		}
	}
}
func (m *Manager) GetSlotUnit() *types.Slot { return m.slotUnit }
func (m *Manager) UseSlot(nodeId string, slotCount uint32) error {
	table, ok := m.localTables[nodeId]
	if !ok {
		return fmt.Errorf("No found the resource table of node: %s", nodeId)
	}
	if table.GetLockedSlot() < slotCount {
		return fmt.Errorf("Insufficient locked number of slots of node: %s", nodeId)
	}
	table.UseSlot(slotCount)
	return nil
}
func (m *Manager) LockSlot(nodeId string, slotCount uint32) error {
	table, ok := m.localTables[nodeId]
	if !ok {
		return fmt.Errorf("No found the resource table of node: %s", nodeId)
	}
	if table.RemianSlot() < slotCount {
		return fmt.Errorf("Insufficient remaining number of slots of node: %s", nodeId)
	}
	table.LockSlot(slotCount)
	return nil
}

func (m *Manager) SetResource(table *types.LocalResourceTable) {
	m.localTables[table.GetNodeId()] = table
	m.localTableQueue = append(m.localTableQueue, table)
}
func (m *Manager) GetResource(nodeId string) *types.LocalResourceTable { return m.localTables[nodeId] }
func (m *Manager) GetResources() []*types.LocalResourceTable           { return m.localTableQueue }
func (m *Manager) DelResource(nodeId string) {
	for i := 0; i < len(m.localTableQueue); i++ {
		table := m.localTableQueue[i]
		if table.GetNodeId() == nodeId {
			delete(m.localTables, nodeId)
			m.localTableQueue = append(m.localTableQueue[:i], m.localTableQueue[i+1:]...)
			i--
		}
	}
}
func (m *Manager) AddRemoteResourceTable(table *types.RemoteResourceTable)  {
	m.remoteTableQueue = append(m.remoteTableQueue, table)
}
func (m *Manager) UpdateRemoteResouceTable(table *types.RemoteResourceTable) {
	for i := 0; i < len(m.remoteTableQueue); i++ {
		if m.remoteTableQueue[i].GetIdentityId() == table.GetIdentityId() {
			m.remoteTableQueue[i] = table
		}
	}
}
func (m *Manager) AddOrUpdateRemoteResouceTable(table *types.RemoteResourceTable) {
	var has bool
	for i := 0; i < len(m.remoteTableQueue); i++ {
		if m.remoteTableQueue[i].GetIdentityId() == table.GetIdentityId() {
			m.remoteTableQueue[i] = table
			has = true
		}
	}
	if has {
		return
	}
	m.remoteTableQueue = append(m.remoteTableQueue, table)
}
func (m *Manager) DelRemoteResourceTable(identityId string)  {
	for i := 0; i < len(m.remoteTableQueue); i++ {
		if m.remoteTableQueue[i].GetIdentityId() == identityId {
			m.remoteTableQueue = append(m.remoteTableQueue[:i], m.remoteTableQueue[i+1:]...)
			i--
		}
	}
}

func (m *Manager) SendTaskEvent(event *event.TaskEvent) error {
	m.eventCh <- event
	return nil
}
