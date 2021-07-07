package resource

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/types"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	defaultRefreshOrgResourceInterval = 60 * time.Millisecond
)


type CarrierDB interface {

	InsertResource(resource *types.Resource) error
	//GetResourceByDataId(powerId string) (*types.Resource, error)
	GetResourceListByNodeId(nodeId string) (types.ResourceArray, error)
	GetResourceList() (types.ResourceArray, error)

	SetRegisterNode(typ types.RegisteredNodeType, node *types.RegisteredNodeInfo) (types.NodeConnStatus, error)
	DeleteRegisterNode(typ types.RegisteredNodeType, id string) error
	GetRegisterNode(typ types.RegisteredNodeType, id string) (*types.RegisteredNodeInfo, error)
	GetRegisterNodeList(typ types.RegisteredNodeType) ([]*types.RegisteredNodeInfo, error)
	//StoreRunningTask(task *types.Task) error
	//StoreJobNodeRunningTaskId(jobNodeId, taskId string) error
	//IncreaseRunningTaskCountOnOrg() uint32
	//IncreaseRunningTaskCountOnJobNode(jobNodeId string) uint32
	//GetRunningTaskCountOnOrg() uint32
	//GetRunningTaskCountOnJobNode(jobNodeId string) uint32
	//GetJobNodeRunningTaskIdList(jobNodeId string) []string

	// For ResourceManager
	StoreLocalResourceTables(resources []*types.LocalResourceTable) error
	QueryLocalResourceTables() ([]*types.LocalResourceTable, error)
	StoreOrgResourceTables(resources []*types.RemoteResourceTable) error
	QueryOrgResourceTables() ([]*types.RemoteResourceTable, error)
	StoreNodeResourceSlotUnit(slot *types.Slot) error
	QueryNodeResourceSlotUnit() (*types.Slot, error)
}

type Manager struct {
	// TODO 这里需要一个 config <SlotUnit 的>
	db              CarrierDB // Low level persistent database to store final content.
	eventCh         chan *types.TaskEventInfo
	slotUnit        *types.Slot
	localTables     map[string]*types.LocalResourceTable
	localTableQueue []*types.LocalResourceTable
	//remoteTables     map[string]*types.RemoteResourceTable
	remoteTableQueue []*types.RemoteResourceTable

	localLock     sync.RWMutex
	remoteLock    sync.RWMutex

}

func NewResourceManager(db CarrierDB) *Manager {
	m := &Manager{
		db:               db,
		eventCh:          make(chan *types.TaskEventInfo, 0),
		localTables:      make(map[string]*types.LocalResourceTable),
		localTableQueue:  make([]*types.LocalResourceTable, 0),
		remoteTableQueue: make([]*types.RemoteResourceTable, 0),
		slotUnit:         types.DefaultSlotUnit, // TODO for test
	}

	return m
}

func (m *Manager) loop() {
	refreshTimer := time.NewTimer(defaultRefreshOrgResourceInterval)
	for {
		select {
		case event := <-m.eventCh:
			_ = event // TODO add some logic about eventEngine
		case <- refreshTimer.C:
			if err := m.refreshOrgResourceTable(); nil != err {
				log.Errorf("Failed to refresh org resourceTables, err: %s", err)
			}
		default:
		}
	}
}

func (m *Manager) Start() error {

	m.SetSlotUnit(0, 0, 0)
	// load slotUnit
	slotUnit, err := m.db.QueryNodeResourceSlotUnit()
	if nil != err {
		return err
	}
	m.slotUnit = slotUnit
	// load local resource Tables
	localResources, err :=  m.db.QueryLocalResourceTables()
	if nil != err {
		return err
	}
	tables := make(map[string]*types.LocalResourceTable, len(localResources))
	for _, resource := range localResources {
		tables[resource.GetNodeId()] = resource
	}
	m.localTables = tables
	m.localTableQueue = localResources

	// load remote org resource Tables
	remoteResources, err :=  m.db.QueryOrgResourceTables()
	if nil != err {
		return err
	}
	m.remoteTableQueue = remoteResources
	go m.loop()
	return nil
}

func (m *Manager) Stop() error {
	// store slotUnit
	if err := m.db.StoreNodeResourceSlotUnit(m.slotUnit); nil != err {
		return err
	}
	// store local resource Tables
	if err := m.db.StoreLocalResourceTables(m.localTableQueue); nil != err {
		return err
	}
	// store remote org resource Tables
	if err := m.db.StoreOrgResourceTables(m.remoteTableQueue); nil != err {
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
func (m *Manager) refreshLocalResourceTable() error {

	return nil
}

func (m *Manager) addRemoteResourceTable(table *types.RemoteResourceTable)  {
	m.remoteTableQueue = append(m.remoteTableQueue, table)
}
func (m *Manager) updateRemoteResouceTable(table *types.RemoteResourceTable) {
	for i := 0; i < len(m.remoteTableQueue); i++ {
		if m.remoteTableQueue[i].GetIdentityId() == table.GetIdentityId() {
			m.remoteTableQueue[i] = table
		}
	}
}
func (m *Manager) addOrUpdateRemoteResouceTable(table *types.RemoteResourceTable) {
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
func (m *Manager) delRemoteResourceTable(identityId string)  {
	for i := 0; i < len(m.remoteTableQueue); i++ {
		if m.remoteTableQueue[i].GetIdentityId() == identityId {
			m.remoteTableQueue = append(m.remoteTableQueue[:i], m.remoteTableQueue[i+1:]...)
			i--
		}
	}
}
func (m *Manager) CleanRemoteResourceTables() {
	m.remoteTableQueue = make([]*types.RemoteResourceTable, 0)
}
func (m *Manager) GetRemoteResouceTables() []*types.RemoteResourceTable { return m.remoteTableQueue }
func (m *Manager) refreshOrgResourceTable() error {
	resources, err := m.db.GetResourceList()
	if nil != err {
		return err
	}
	tmp := make(map[string]*types.Resource, len(resources))
	for _, r := range resources {
		tmp[r.GetIdentityId()] = r
	}
	for i, resource := range m.remoteTableQueue {
		// If has, update
		if r, ok := tmp[resource.GetIdentityId()]; ok {
			m.remoteTableQueue[i] = types.NewOrgResourceFromResource(r)
			delete(tmp, resource.GetIdentityId())
		} else {
		// If no has, delete
			m.remoteTableQueue = append(m.remoteTableQueue[:i], m.remoteTableQueue[i+1:]...)
			i--
		}
	}
	for _, r := range tmp {
		m.remoteTableQueue = append(m.remoteTableQueue, types.NewOrgResourceFromResource(r))
	}
	return nil
}
func (m *Manager) SendTaskEvent(event *types.TaskEventInfo) error {
	m.eventCh <- event
	return nil
}
