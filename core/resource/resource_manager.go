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
	defaultRefreshOrgResourceInterval = 30 * time.Second
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
	refreshTicker := time.NewTicker(defaultRefreshOrgResourceInterval)
	for {
		select {
		case <-refreshTicker.C:
			if err := m.refreshOrgResourceTable(); nil != err {
				log.Errorf("Failed to refresh org resourceTables on loop, err: %s", err)
			}
		}
	}
}

func (m *Manager) Start() error {

	slotUnit, err := m.dataCenter.QueryNodeResourceSlotUnit()
	if nil != err {
		log.Warnf("Failed to load local slotUnit on resourceManager Start(), err: {%s}", err)
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
	if len(remoteResources) != 0 {
		m.remoteTableQueue = remoteResources
	} else {
		if err := m.refreshOrgResourceTable(); nil != err {
			log.Errorf("Failed to refresh org resourceTables on Start resourceManager, err: %s", err)
		}
	}

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

//func (m *Manager) AddRemoteResourceTable(table *types.RemoteResourceTable) {
//	m.remoteTableQueue = append(m.remoteTableQueue, table)
//}
//func (m *Manager) UpdateRemoteResouceTable(table *types.RemoteResourceTable) {
//	for i := 0; i < len(m.remoteTableQueue); i++ {
//		if m.remoteTableQueue[i].GetIdentityId() == table.GetIdentityId() {
//			m.remoteTableQueue[i] = table
//		}
//	}
//}
//func (m *Manager) AddOrUpdateRemoteResouceTable(table *types.RemoteResourceTable) {
//	var has bool
//	for i := 0; i < len(m.remoteTableQueue); i++ {
//		if m.remoteTableQueue[i].GetIdentityId() == table.GetIdentityId() {
//			m.remoteTableQueue[i] = table
//			has = true
//		}
//	}
//	if has {
//		return
//	}
//	m.remoteTableQueue = append(m.remoteTableQueue, table)
//}
//func (m *Manager) DelRemoteResourceTable(identityId string) {
//	for i := 0; i < len(m.remoteTableQueue); i++ {
//		if m.remoteTableQueue[i].GetIdentityId() == identityId {
//			m.remoteTableQueue = append(m.remoteTableQueue[:i], m.remoteTableQueue[i+1:]...)
//			i--
//		}
//	}
//}
//func (m *Manager) CleanRemoteResourceTables() {
//	m.remoteTableQueue = make([]*types.RemoteResourceTable, 0)
//}
func (m *Manager) GetRemoteResouceTables() []*types.RemoteResourceTable { return m.remoteTableQueue }
func (m *Manager) refreshOrgResourceTable() error {
	resources, err := m.dataCenter.GetResourceList()
	if nil != err {
		return err
	}

	remoteResourceArr := make([]*types.RemoteResourceTable, len(resources))

	for i, r := range resources {
		remoteResourceArr[i] = types.NewOrgResourceFromResource(r)
	}

	m.remoteTableQueue = remoteResourceArr

	//tmpMap := make(map[string]int, len(resources))
	//tmpArr := make([]*types.Resource, len(resources))
	//for i, r := range resources {
	//	tmpMap[r.GetIdentityId()] = i
	//	tmpArr[i] = r
	//}
	//
	//for i := 0; i < len(m.remoteTableQueue); i++ {
	//
	//	resource := m.remoteTableQueue[i]
	//
	//	// If has, update
	//	if index, ok := tmpMap[resource.GetIdentityId()]; ok {
	//		r := tmpArr[index]
	//		m.remoteTableQueue[i] = types.NewOrgResourceFromResource(r)
	//		delete(tmpMap, resource.GetIdentityId())
	//	} else {
	//		// If no has, delete
	//		m.remoteTableQueue = append(m.remoteTableQueue[:i], m.remoteTableQueue[i+1:]...)
	//		i--
	//	}
	//}
	//// If is new one, add
	//if len(tmpMap) != 0 {
	//	for _, r := range tmpArr {
	//		if _, ok := tmpMap[r.GetIdentityId()]; ok {
	//			m.remoteTableQueue = append(m.remoteTableQueue, types.NewOrgResourceFromResource(r))
	//		}
	//	}
	//}
	return nil
}

// TODO 有变更 RegisterNode mem  processor bandwidth 的 接口咩 ？？？
func (m *Manager) LockLocalResourceWithTask(jobNodeId string, needSlotCount uint64, task *types.Task) error {

	log.Infof("Start lock local resource with taskId {%s}, jobNodeId {%s}, slotCount {%d}", task.TaskId(), jobNodeId, needSlotCount)

	// Lock local resource (jobNode)
	if err := m.UseSlot(jobNodeId, uint32(needSlotCount)); nil != err {
		log.Errorf("Failed to lock internal power resource, taskId: {%s}, jobNodeId: {%s}, usedSlotCount: {%s}, err: {%s}",
			task.TaskId(), jobNodeId, needSlotCount, err)
		return fmt.Errorf("failed to lock internal power resource, {%s}", err)
	}


	if err := m.dataCenter.StoreJobNodeRunningTaskId(jobNodeId, task.TaskId()); nil != err {

		m.FreeSlot(jobNodeId, uint32(needSlotCount))

		log.Errorf("Failed to store local taskId and jobNodeId index, taskId: {%s}, jobNodeId: {%s}, usedSlotCount: {%s}, err: {%s}",
			task.TaskId(), jobNodeId, needSlotCount, err)
		return fmt.Errorf("failed to store local taskId and jobNodeId index, {%s}", err)
	}
	if err := m.dataCenter.StoreLocalTaskPowerUsed(types.NewLocalTaskPowerUsed(task.TaskId(), jobNodeId, needSlotCount)); nil != err {

		m.FreeSlot(jobNodeId, uint32(needSlotCount))
		m.dataCenter.RemoveJobNodeRunningTaskId(jobNodeId, task.TaskId())

		log.Errorf("Failed to store local taskId use jobNode slot, taskId: {%s}, jobNodeId: {%s}, usedSlotCount: {%s}, err: {%s}",
			task.TaskId(), jobNodeId, needSlotCount, err)
		return fmt.Errorf("failed to store local taskId use jobNode slot, {%s}", err)
	}


	// 更新本地 resource 资源信息 [添加资源使用情况]
	jobNodeResource, err := m.dataCenter.GetLocalResource(jobNodeId)
	if nil != err {

		m.FreeSlot(jobNodeId, uint32(needSlotCount))
		m.dataCenter.RemoveJobNodeRunningTaskId(jobNodeId, task.TaskId())
		m.dataCenter.RemoveLocalTaskPowerUsed(task.TaskId())

		log.Errorf("Failed to query local jobNodeResource, taskId: {%s}, jobNodeId: {%s}, usedSlotCount: {%s}, err: {%s}",
			task.TaskId(), jobNodeId, needSlotCount, err)
		return fmt.Errorf("failed to query local jobNodeResource, {%s}", err)
	}

	// 更新 本地 jobNodeResource 的资源使用信息
	usedMem := m.slotUnit.Mem * needSlotCount
	usedProcessor := m.slotUnit.Processor * needSlotCount
	usedBandwidth := m.slotUnit.Bandwidth * needSlotCount

	jobNodeResource.GetData().UsedMem += usedMem
	jobNodeResource.GetData().UsedProcessor += usedProcessor
	jobNodeResource.GetData().UsedBandWidth += usedBandwidth
	if err := m.dataCenter.InsertLocalResource(jobNodeResource); nil != err {

		m.FreeSlot(jobNodeId, uint32(needSlotCount))
		m.dataCenter.RemoveJobNodeRunningTaskId(jobNodeId, task.TaskId())
		m.dataCenter.RemoveLocalTaskPowerUsed(task.TaskId())

		log.Errorf("Failed to update local jobNodeResource, taskId: {%s}, jobNodeId: {%s}, usedSlotCount: {%s}, err: {%s}",
			task.TaskId(), jobNodeId, needSlotCount, err)
		return fmt.Errorf("failed to update local jobNodeResource, {%s}", err)
	}

	// 还需要 将资源使用实况 实时上报给  dataCenter  [添加资源使用情况]
	if err := m.dataCenter.SyncPowerUsed(jobNodeResource); nil != err {
		log.Errorf("Failed to sync jobNodeResource to dataCenter, taskId: {%s}, jobNodeId: {%s}, usedSlotCount: {%s}, err: {%s}",
			task.TaskId(), jobNodeId, needSlotCount, err)
		return fmt.Errorf("failed tosync jobNodeResource to dataCenter, {%s}", err)
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

	jobNodeId := localTaskPowerUsed.GetNodeId()
	freeSlotUnitCount := localTaskPowerUsed.GetSlotCount()

	log.Infof("Start unlock local resource with taskId {%s}, jobNodeId {%s}, slotCount {%d}", taskId, jobNodeId, localTaskPowerUsed.GetSlotCount())

	// Lock local resource (jobNode)
	if err := m.FreeSlot(localTaskPowerUsed.GetNodeId(), uint32(freeSlotUnitCount)); nil != err {
		log.Errorf("Failed to unlock internal power resource, taskId: {%s}, jobNodeId: {%s}, freeSlotUnitCount: {%s}, err: {%s}",
			taskId, jobNodeId, freeSlotUnitCount, err)
		return fmt.Errorf("failed to unlock internal power resource, {%s}", err)
	}

	if err := m.dataCenter.RemoveTaskEventList(taskId); nil != err {
		log.Errorf("Failed to remove local task event list, taskId: {%s}, jobNodeId: {%s}, freeSlotUnitCount: {%s}, err: {%s}",
			taskId, jobNodeId, freeSlotUnitCount, err)
		return fmt.Errorf("failed to remove local task event list, {%s}", err)
	}

	if err := m.dataCenter.RemoveJobNodeRunningTaskId(jobNodeId, taskId); nil != err {
		log.Errorf("Failed to remove local taskId and jobNodeId index, taskId: {%s}, jobNodeId: {%s}, freeSlotUnitCount: {%s}, err: {%s}",
			taskId, jobNodeId, freeSlotUnitCount, err)
		return fmt.Errorf("failed to remove local taskId and jobNodeId index, {%s}", err)
	}
	if err := m.dataCenter.RemoveLocalTaskPowerUsed(taskId); nil != err {
		log.Errorf("Failed to remove local taskId use jobNode slot, taskId: {%s}, jobNodeId: {%s}, freeSlotUnitCount: {%s}, err: {%s}",
			taskId, jobNodeId, freeSlotUnitCount, err)
		return fmt.Errorf("failed to remove local taskId use jobNode slot, {%s}", err)
	}

	// 更新本地 resource 资源信息 [释放资源使用情况]
	jobNodeResource, err := m.dataCenter.GetLocalResource(jobNodeId)
	if nil != err {
		log.Errorf("Failed to query local jobNodeResource, taskId: {%s}, jobNodeId: {%s}, freeSlotUnitCount: {%s}, err: {%s}",
			taskId, jobNodeId, freeSlotUnitCount, err)
		return fmt.Errorf("failed to query local jobNodeResource, {%s}", err)
	}

	// 更新 本地 jobNodeResource 的资源使用信息
	usedMem := m.slotUnit.Mem * freeSlotUnitCount
	usedProcessor := m.slotUnit.Processor * freeSlotUnitCount
	usedBandwidth := m.slotUnit.Bandwidth * freeSlotUnitCount

	jobNodeResource.GetData().UsedMem -= usedMem
	jobNodeResource.GetData().UsedProcessor -= usedProcessor
	jobNodeResource.GetData().UsedBandWidth -= usedBandwidth

	if err := m.dataCenter.InsertLocalResource(jobNodeResource); nil != err {
		log.Errorf("Failed to update local jobNodeResource, taskId: {%s}, jobNodeId: {%s}, freeSlotUnitCount: {%s}, err: {%s}",
			taskId, jobNodeId, freeSlotUnitCount, err)
		return fmt.Errorf("failed to update local jobNodeResource, {%s}", err)
	}

	// 还需要 将资源使用实况 实时上报给  dataCenter  [释放资源使用情况]
	if err := m.dataCenter.SyncPowerUsed(jobNodeResource); nil != err {
		log.Errorf("Failed to sync jobNodeResource to dataCenter, taskId: {%s}, jobNodeId: {%s}, freeSlotUnitCount: {%s}, err: {%s}",
			taskId, jobNodeId, freeSlotUnitCount, err)
		return fmt.Errorf("failed tosync jobNodeResource to dataCenter, {%s}", err)
	}

	log.Infof("Finished unlock local resource with taskId {%s}, jobNodeId {%s}, slotCount {%d}", taskId, localTaskPowerUsed.GetNodeId(), localTaskPowerUsed.GetSlotCount())
	return nil
}



func (m *Manager) ReleaseLocalResourceWithTask (logdesc, taskId string, option ReleaseResourceOption) {

	log.Debugf("Start ReleaseLocalResourceWithTask %s, taskId: {%s}, releaseOption: {%d}", logdesc, taskId, option)

	has, err := m.dataCenter.HasLocalTaskExecute(taskId)
	if nil != err {
		log.Errorf("Failed to query local task exec status with task %s, taskId: {%s}, err: {%s}", logdesc, taskId, err)
		return
	}

	if has {
		log.Debugf("The local task have been executing, don't `ReleaseLocalResourceWithTask` %s, taskId: {%s}, releaseOption: {%d}", logdesc, taskId, option)
		return
	}

	if option.IsUnlockLocalResorce() {
		log.Debugf("start unlock local resource with task %s, taskId: {%s}", logdesc, taskId)
		if err := m.UnLockLocalResourceWithTask(taskId); nil != err {
			log.Errorf("Failed to unlock local resource with task %s, taskId: {%s}, err: {%s}", logdesc, taskId, err)
		}
	}
	if option.IsRemoveLocalTask() {
		log.Debugf("start remove local task  %s, taskId: {%s}", logdesc, taskId)
		// 因为在 scheduler 那边已经对 task 做了 StoreLocalTask
		if err := m.dataCenter.RemoveLocalTask(taskId); nil != err {
			log.Errorf("Failed to remove local task  %s, taskId: {%s}, err: {%s}", logdesc, taskId, err)
		}
	}
	if option.IsCleanTaskEvents() {
		log.Debugf("start clean event list of task  %s, taskId: {%s}", logdesc, taskId)
		if err := m.dataCenter.RemoveTaskEventList(taskId); nil != err {
			log.Errorf("Failed to clean event list of task  %s, taskId: {%s}, err: {%s}", logdesc, taskId, err)
		}
	}
}

// todo 构造一些假的 本地任务信息
//func (m *Manager) mockLocalTaskList(){
//	identity, err := m.dataCenter.GetIdentity()
//	if nil != err {
//		log.Warnf("failed to query identityInfo, err: {%s}", err)
//		return
//	}
//
//
//
//}