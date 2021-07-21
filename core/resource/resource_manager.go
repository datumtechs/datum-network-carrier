package resource

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/core/iface"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/types"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	defaultRefreshOrgResourceInterval = 60 * time.Millisecond
)

type Manager struct {
	// TODO 这里需要一个 config <SlotUnit 的>
	dataCenter iface.ForResourceDB // Low level persistent database to store final content.
	//eventCh                chan *types.TaskEventInfo
	slotUnit *types.Slot
	//remoteTables     map[string]*types.RemoteResourceTable
	remoteTableQueue []*types.RemoteResourceTable
}

func NewResourceManager(dataCenter iface.ForResourceDB) *Manager {
	m := &Manager{
		dataCenter: dataCenter,
		//eventCh:          make(chan *types.TaskEventInfo, 0),
		//localTables:      make(map[string]*types.LocalResourceTable),
		//localTableQueue:  make([]*types.LocalResourceTable, 0),
		remoteTableQueue: make([]*types.RemoteResourceTable, 0),
		slotUnit:         types.DefaultSlotUnit, // TODO for test
	}

	return m
}

func (m *Manager) loop() {
	refreshTimer := time.NewTimer(defaultRefreshOrgResourceInterval)
	for {
		select {
		case <-refreshTimer.C:
			if err := m.refreshOrgResourceTable(); nil != err {
				log.Errorf("Failed to refresh org resourceTables, err: %s", err)
			}
		default:
		}
	}
}

func (m *Manager) Start() error {

	slotUnit, err := m.dataCenter.QueryNodeResourceSlotUnit()
	if nil != err {
		log.Warn("Failed to load local slotUnit on resourceManager Start(), err: {%s}", err)
	} else {
		m.SetSlotUnit(slotUnit.Mem, slotUnit.Processor, slotUnit.Bandwidth)
	}

	// store slotUnit
	if err := m.dataCenter.StoreNodeResourceSlotUnit(m.slotUnit); nil != err {
		return err
	}
	// load remote org resource Tables
	remoteResources, err := m.dataCenter.QueryOrgResourceTables()
	if nil != err && err != rawdb.ErrNotFound {
		return err
	}
	m.remoteTableQueue = remoteResources
	go m.loop()
	log.Info("Started resourceManager ...")
	return nil
}

func (m *Manager) Stop() error {
	// store slotUnit
	if err := m.dataCenter.StoreNodeResourceSlotUnit(m.slotUnit); nil != err {
		return err
	}
	// store local resource Tables
	//if err := m.db.StoreLocalResourceTables(m.localTableQueue); nil != err {
	//	return err
	//}
	// store remote org resource Tables
	if err := m.dataCenter.StoreOrgResourceTables(m.remoteTableQueue); nil != err {
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
	//if len(m.localTables) != 0 {
	//	for _, re := range m.localTables {
	//		re.SetSlotUnit(m.slotUnit)
	//	}
	//}
}
func (m *Manager) GetSlotUnit() *types.Slot { return m.slotUnit }

//func (m *Manager) UseSlot(nodeId string, slotCount uint32) error {
//
//	table, err := m.GetLocalResourceTable(nodeId)
//	if nil != err {
//		return fmt.Errorf("No found the resource table of node: %s, %s", nodeId, err)
//	}
//	if table.GetLockedSlot() < slotCount {
//		return fmt.Errorf("Insufficient locked number of slots of node: %s", nodeId)
//	}
//	table.UseSlot(slotCount)
//	return m.SetLocalResourceTable(table)
//}
//func (m *Manager) FreeSlot(nodeId string, slotCount uint32) error {
//	table, err := m.GetLocalResourceTable(nodeId)
//	if nil != err {
//		return fmt.Errorf("No found the resource table of node: %s, %s", nodeId, err)
//	}
//	table.FreeSlot(slotCount)
//	return m.SetLocalResourceTable(table)
//}
func (m *Manager) UseSlot(nodeId string, slotCount uint32) error {
	table, err := m.GetLocalResourceTable(nodeId)
	if nil != err {
		return fmt.Errorf("No found the resource table of node: %s, %s", nodeId, err)
	}
	//if table.RemianSlot() < slotCount {
	//	return fmt.Errorf("Insufficient slotRemain {%s} less than need lock count {%s} slots of node: %s", table.RemianSlot(),slotCount , nodeId)
	//}
	if err := table.UseSlot(slotCount); nil != err {
		return err
	}
	return m.SetLocalResourceTable(table)
}
func (m *Manager) FreeSlot(nodeId string, slotCount uint32) error {
	table, err := m.GetLocalResourceTable(nodeId)
	if nil != err {
		return fmt.Errorf("No found the resource table of node: %s, %s", nodeId, err)
	}
	if err := table.FreeSlot(slotCount); nil != err {
		return err
	}
	return m.SetLocalResourceTable(table)
}

func (m *Manager) SetLocalResourceTable(table *types.LocalResourceTable) error {
	return m.dataCenter.StoreLocalResourceTable(table)
}
func (m *Manager) GetLocalResourceTable(nodeId string) (*types.LocalResourceTable, error) {
	return m.dataCenter.QueryLocalResourceTable(nodeId)
}
func (m *Manager) GetLocalResourceTables() ([]*types.LocalResourceTable, error) {
	return m.dataCenter.QueryLocalResourceTables()
}
func (m *Manager) DelLocalResourceTable(nodeId string) error {
	return m.dataCenter.RemoveLocalResourceTable(nodeId)
}
func (m *Manager) CleanLocalResourceTables() error {
	localResourceTableArr, err := m.dataCenter.QueryLocalResourceTables()
	if nil != err {
		return err
	}
	for _, table := range localResourceTableArr {
		if err := m.dataCenter.RemoveLocalResourceTable(table.GetNodeId()); nil != err {
			return err
		}
	}
	return nil
}

func (m *Manager) AddRemoteResourceTable(table *types.RemoteResourceTable) {
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
func (m *Manager) DelRemoteResourceTable(identityId string) {
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
	resources, err := m.dataCenter.GetResourceList()
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
	// If is new one, add
	for _, r := range tmp {
		m.remoteTableQueue = append(m.remoteTableQueue, types.NewOrgResourceFromResource(r))
	}
	return nil
}

// TODO 有变更 RegisterNode mem  processor bandwidth 的 接口咩 ？？？
func (m *Manager) LockLocalResourceWithTask(jobNodeId string, needSlotCount uint64, task *types.Task) error {

	log.Infof("Start lock local resource with taskId {%s}, jobNodeId {%s}, slotCount {%d}", task.TaskId(), jobNodeId, needSlotCount)

	// Lock local resource (jobNode)
	if err := m.UseSlot(jobNodeId, uint32(needSlotCount)); nil != err {
		log.Errorf("Failed to lock internal power resource, err: %s", err)
		return fmt.Errorf("failed to lock internal power resource, err: %s", err)
	}

	//// task 不论是 发起方 还是 参与方, 都应该是  一抵达, 就保存本地..
	//if err := m.dataCenter.StoreLocalTask(task); nil != err {
	//
	//	m.FreeSlot(jobNodeId, uint32(needSlotCount))
	//
	//	log.Errorf("Failed to store local task, err: %s", err)
	//	return fmt.Errorf("failed to store local task, err: %s", err)
	//}

	if err := m.dataCenter.StoreJobNodeRunningTaskId(jobNodeId, task.TaskId()); nil != err {

		m.FreeSlot(jobNodeId, uint32(needSlotCount))
		m.dataCenter.RemoveLocalTask(task.TaskId())

		log.Errorf("Failed to store local taskId and jobNodeId index, err: %s", err)
		return fmt.Errorf("failed to store local taskId and jobNodeId index, err: %s", err)
	}
	if err := m.dataCenter.StoreLocalTaskPowerUsed(types.NewLocalTaskPowerUsed(task.TaskId(), jobNodeId, needSlotCount)); nil != err {

		m.FreeSlot(jobNodeId, uint32(needSlotCount))
		m.dataCenter.RemoveLocalTask(task.TaskId())
		m.dataCenter.RemoveJobNodeRunningTaskId(jobNodeId, task.TaskId())

		log.Errorf("Failed to store local taskId use jobNode slot, err: %s", err)
		return fmt.Errorf("failed to store local taskId use jobNode slot, err: %s", err)
	}

	log.Infof("Finished lock local resource with taskId {%s}, jobNodeId {%s}, slotCount {%d}", task.TaskId(), jobNodeId, needSlotCount)
	return nil
}

// TODO 有变更 RegisterNode mem  processor bandwidth 的 接口咩 ？？？
func (m *Manager) UnLockLocalResourceWithTask(taskId string) error {

	localTaskPowerUsed, err := m.dataCenter.QueryLocalTaskPowerUsed(taskId)
	if nil != err {
		log.Errorf("Failed to query local taskId and jobNodeId index, err: %s", err)
		return fmt.Errorf("failed to query local taskId and jobNodeId index, err: %s", err)
	}

	log.Infof("Start unlock local resource with taskId {%s}, jobNodeId {%s}, slotCount {%d}", taskId, localTaskPowerUsed.GetNodeId(), localTaskPowerUsed.GetSlotCount())

	// Lock local resource (jobNode)
	if err := m.FreeSlot(localTaskPowerUsed.GetNodeId(), uint32(localTaskPowerUsed.GetSlotCount())); nil != err {
		log.Errorf("Failed to unlock internal power resource, err: %s", err)
		return fmt.Errorf("failed to unlock internal power resource, err: %s", err)
	}
	//// 移除 本地任务
	//if err := m.dataCenter.RemoveLocalTask(taskId); nil != err {
	//	log.Errorf("Failed to remove local task, err: %s", err)
	//	return fmt.Errorf("failed to remove local task, err: %s", err)
	//}

	if err := m.dataCenter.CleanTaskEventList(taskId); nil != err {
		log.Errorf("Failed to remove local task event list, err: %s", err)
		return fmt.Errorf("failed to remove local task event list, err: %s", err)
	}

	if err := m.dataCenter.RemoveJobNodeRunningTaskId(localTaskPowerUsed.GetNodeId(), taskId); nil != err {
		log.Errorf("Failed to remove local taskId and jobNodeId index, err: %s", err)
		return fmt.Errorf("failed to remove local taskId and jobNodeId index, err: %s", err)
	}
	if err := m.dataCenter.RemoveLocalTaskPowerUsed(taskId); nil != err {

		log.Errorf("Failed to remove local taskId use jobNode slot, err: %s", err)
		return fmt.Errorf("failed to remove local taskId use jobNode slot, err: %s", err)
	}
	log.Infof("Finished unlock local resource with taskId {%s}, jobNodeId {%s}, slotCount {%d}", taskId, localTaskPowerUsed.GetNodeId(), localTaskPowerUsed.GetSlotCount())
	return nil
}
