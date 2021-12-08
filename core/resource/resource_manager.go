package resource

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/fileutil"
	"github.com/RosettaFlow/Carrier-Go/core"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/types"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	defaultRefreshOrgResourceInterval = 30 * time.Second
)

type Manager struct {
	dataCenter  core.CarrierDB // Low level persistent database to store final content.
	mockIdentityIdsFile  string
	mockIdentityIdsCache map[string]struct{}
}

func NewResourceManager(dataCenter core.CarrierDB, mockIdentityIdsFile string) *Manager {
	m := &Manager{
		dataCenter: dataCenter,
		mockIdentityIdsFile: mockIdentityIdsFile,   //TODO for test
		mockIdentityIdsCache: make(map[string]struct{}, 0),
	}

	return m
}

func (m *Manager) loop() {
}

func (m *Manager) Start() error {

	// build mock identityIds cache
	if "" != m.mockIdentityIdsFile {
		var identityIdList []string
		if err := fileutil.LoadJSON(m.mockIdentityIdsFile, &identityIdList); err != nil {
			log.WithError(err).Errorf("Failed to load `--mock-identity-file` on Start resourceManager, file: {%s}", m.mockIdentityIdsFile)
			return err
		}

		for _, iden := range identityIdList {
			m.mockIdentityIdsCache[iden] = struct{}{}
		}
	}


	//go m.loop()
	log.Info("Started resourceManager ...")
	return nil
}

func (m *Manager) Stop() error {
	log.Infof("Stopped resource manager ...")
	return nil
}

func (m *Manager) UseSlot(nodeId string, mem, bandwidth, disk uint64, processor uint32) error {
	table, err := m.QueryLocalResourceTable(nodeId)
	if nil != err {
		return fmt.Errorf("No found the resource table of node: %s, %s", nodeId, err)
	}
	if err := table.UseSlot(mem, bandwidth, disk, processor); nil != err {
		return err
	}
	return m.StoreLocalResourceTable(table)
}
func (m *Manager) FreeSlot(nodeId string, mem, bandwidth, disk uint64, processor uint32) error {
	table, err := m.QueryLocalResourceTable(nodeId)
	if nil != err {
		return fmt.Errorf("No found the resource table of node: %s, %s", nodeId, err)
	}
	if err := table.FreeSlot(mem, bandwidth, disk, processor); nil != err {
		return err
	}
	return m.StoreLocalResourceTable(table)
}

func (m *Manager) StoreLocalResourceTable(table *types.LocalResourceTable) error {
	return m.dataCenter.StoreLocalResourceTable(table)
}
func (m *Manager) QueryLocalResourceTable(nodeId string) (*types.LocalResourceTable, error) {
	return m.dataCenter.QueryLocalResourceTable(nodeId)
}
func (m *Manager) QueryLocalResourceTables() ([]*types.LocalResourceTable, error) {
	return m.dataCenter.QueryLocalResourceTables()
}
func (m *Manager) RemoveLocalResourceTable(nodeId string) error {
	return m.dataCenter.RemoveLocalResourceTable(nodeId)
}
func (m *Manager) RemoveLocalResourceTables() error {
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

func (m *Manager) LockLocalResourceWithTask(partyId, jobNodeId string, mem, bandwidth, disk uint64, processor uint32, task *types.Task) error {

	log.Infof("Start lock local resource on resourceManager.LockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, needMem: {%d}, needBandwidth: {%d}, needDisk: {%d}, needProcessor: {%d}",
		task.GetTaskId(), partyId, jobNodeId, mem, bandwidth, disk, processor)
	// Lock local resource (jobNode)
	if err := m.UseSlot(jobNodeId, mem, bandwidth, disk, processor); nil != err {
		log.WithError(err).Errorf("Failed to lock internal power resource on resourceManager.LockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, needMem: {%d}, needBandwidth: {%d}, needDisk: {%d}, needProcessor: {%d}",
			task.GetTaskId(), partyId, jobNodeId, mem, bandwidth, disk, processor)
		return err
	}

	used := types.NewLocalTaskPowerUsed(task.GetTaskId(), partyId, jobNodeId, mem, bandwidth, disk, processor)
	if err := m.addPartyTaskPowerUsedOnJobNode(used); nil != err {
		// rollback useSlot => freeSlot
		m.FreeSlot(jobNodeId, mem, bandwidth, disk, processor)
		return err
	}

	// 更新本地 resource 资源信息 [添加资源使用情况]
	jobNodeResource, err := m.dataCenter.QueryLocalResource(jobNodeId)
	if nil != err {
		// rollback useSlot => freeSlot
		// rollback addPartyTaskPowerUsedOnJobNode
		m.FreeSlot(jobNodeId, mem, bandwidth, disk, processor)
		m.removePartyTaskPowerUsedOnJobNode(used)

		log.WithError(err).Errorf("Failed to query local jobNodeResource on resourceManager.LockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, needMem: {%d}, needBandwidth: {%d}, needDisk: {%d}, needProcessor: {%d}",
			task.GetTaskId(), partyId, jobNodeId, mem, bandwidth, disk, processor)
		return err
	}

	jobNodeRunningTaskCount, err := m.dataCenter.QueryJobNodeRunningTaskCount(jobNodeId)
	if nil != err {
		// rollback useSlot => freeSlot
		// rollback addPartyTaskPowerUsedOnJobNode
		m.FreeSlot(jobNodeId, mem, bandwidth, disk, processor)
		m.removePartyTaskPowerUsedOnJobNode(used)

		log.WithError(err).Errorf("Failed to query task runningCount in jobNode on resourceManager.LockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, needMem: {%d}, needBandwidth: {%d}, needDisk: {%d}, needProcessor: {%d}",
			task.GetTaskId(), partyId, jobNodeId, mem, bandwidth, disk, processor)
		return err
	}

	// 更新 本地 jobNodeResource 的资源使用信息
	jobNodeResource.GetData().UsedMem += mem
	jobNodeResource.GetData().UsedProcessor += processor
	jobNodeResource.GetData().UsedBandwidth += bandwidth
	jobNodeResource.GetData().UsedDisk += disk
	if jobNodeRunningTaskCount > 0 {
		jobNodeResource.GetData().State = apicommonpb.PowerState_PowerState_Occupation
	}
	if err := m.dataCenter.InsertLocalResource(jobNodeResource); nil != err {
		// rollback useSlot => freeSlot
		// rollback addPartyTaskPowerUsedOnJobNode
		m.FreeSlot(jobNodeId, mem, bandwidth, disk, processor)
		m.removePartyTaskPowerUsedOnJobNode(used)

		log.WithError(err).Errorf("Failed to update local jobNodeResource on resourceManager.LockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, needMem: {%d}, needBandwidth: {%d}, needDisk: {%d}, needProcessor: {%d}",
			task.GetTaskId(), partyId, jobNodeId, mem, bandwidth, disk, processor)
		return err
	}

	// 还需要 将资源使用实况 实时上报给  dataCenter  [添加资源使用情况]
	if err := m.dataCenter.SyncPowerUsed(jobNodeResource); nil != err {
		log.WithError(err).Errorf("Failed to sync jobNodeResource to dataCenter on resourceManager.LockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, needMem: {%d}, needBandwidth: {%d}, needDisk: {%d}, needProcessor: {%d}",
			task.GetTaskId(), partyId, jobNodeId, mem, bandwidth, disk, processor)
		return err
	}

	log.Infof("Finished lock local resource with, taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, needMem: {%d}, needBandwidth: {%d}, needDisk: {%d}, needProcessor: {%d}",
		task.GetTaskId(), partyId, jobNodeId, mem, bandwidth, disk, processor)
	return nil
}

func (m *Manager) UnLockLocalResourceWithTask(taskId, partyId string) error {
	used, err := m.dataCenter.QueryLocalTaskPowerUsed(taskId, partyId)
	if nil != err {
		log.WithError(err).Warnf("Warning query local task powerUsed on resourceManager.UnLockLocalResourceWithTask(), taskId {%s}, partyId: {%s}", taskId, partyId)
		return err
	}

	jobNodeId := used.GetNodeId()
	freeMemCount := used.GetUsedMem()
	freeBandwidthCount := used.GetUsedBandwidth()
	freeDiskCount := used.GetUsedDisk()
	freeProcessorCount := used.GetUsedProcessor()

	log.Infof("Start unlock local resource on resourceManager.UnLockLocalResourceWithTask(), taskId {%s}, partyId: {%s}, jobNodeId {%s}, freeMemCount: {%d}, freeBandwidthCount: {%d}, freeDiskCount: {%d}, freeProcessorCount: {%d}, used: %s",
		taskId, partyId, jobNodeId, freeMemCount, freeBandwidthCount, freeDiskCount, freeProcessorCount, used.String())

	// Unlock local resource (jobNode)
	if err := m.FreeSlot(used.GetNodeId(), freeMemCount, freeBandwidthCount, freeDiskCount, freeProcessorCount); nil != err {
		log.WithError(err).Errorf("Failed to freeSlot withJobNodeId on resourceManager.UnLockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, freeMemCount: {%d}, freeBandwidthCount: {%d}, freeDiskCount: {%d}, freeProcessorCount: {%d}",
			taskId, partyId, jobNodeId, freeMemCount, freeBandwidthCount, freeDiskCount, freeProcessorCount)
		return err
	}

	if err := m.removePartyTaskPowerUsedOnJobNode(used); nil != err {
		log.WithError(err).Errorf("Failed to remove partyTaskPowerUsed on resourceManager.UnLockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, freeMemCount: {%d}, freeBandwidthCount: {%d}, freeDiskCount: {%d}, freeProcessorCount: {%d}",
			taskId, partyId, jobNodeId, freeMemCount, freeBandwidthCount, freeDiskCount, freeProcessorCount)
		return err
	}

	// 更新本地 resource 资源信息 [释放资源使用情况]
	jobNodeResource, err := m.dataCenter.QueryLocalResource(jobNodeId)
	if nil != err {
		log.WithError(err).Errorf("Failed to query local jobNodeResource on resourceManager.UnLockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, freeMemCount: {%d}, freeBandwidthCount: {%d}, freeDiskCount: {%d}, freeProcessorCount: {%d}",
			taskId, partyId, jobNodeId, freeMemCount, freeBandwidthCount, freeDiskCount, freeProcessorCount)
		return err
	}

	jobNodeRunningTaskCount, err := m.dataCenter.QueryJobNodeRunningTaskCount(jobNodeId)
	if nil != err {
		log.WithError(err).Errorf("Failed to query task runningCount in jobNode on resourceManager.UnLockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, freeMemCount: {%d}, freeBandwidthCount: {%d}, freeDiskCount: {%d}, freeProcessorCount: {%d}",
			taskId, partyId, jobNodeId, freeMemCount, freeBandwidthCount, freeDiskCount, freeProcessorCount)
		return err
	}

	jobNodeResource.GetData().UsedMem -= freeMemCount
	jobNodeResource.GetData().UsedProcessor -= freeProcessorCount
	jobNodeResource.GetData().UsedBandwidth -= freeBandwidthCount
	jobNodeResource.GetData().UsedDisk -= freeDiskCount
	if jobNodeRunningTaskCount == 0 {
		jobNodeResource.GetData().State = apicommonpb.PowerState_PowerState_Released
	}
	if err := m.dataCenter.InsertLocalResource(jobNodeResource); nil != err {
		log.WithError(err).Errorf("Failed to update local jobNodeResource on resourceManager.UnLockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, freeMemCount: {%d}, freeBandwidthCount: {%d}, freeDiskCount: {%d}, freeProcessorCount: {%d}",
			taskId, partyId, jobNodeId, freeMemCount, freeBandwidthCount, freeDiskCount, freeProcessorCount)
		return err
	}

	// 还需要 将资源使用实况 实时上报给  dataCenter  [释放资源使用情况]
	if err := m.dataCenter.SyncPowerUsed(jobNodeResource); nil != err {
		log.WithError(err).Errorf("Failed to sync jobNodeResource to dataCenter on resourceManager.UnLockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, freeMemCount: {%d}, freeBandwidthCount: {%d}, freeDiskCount: {%d}, freeProcessorCount: {%d}",
			taskId, partyId, jobNodeId, freeMemCount, freeBandwidthCount, freeDiskCount, freeProcessorCount)
		return err
	}

	log.Infof("Finished unlock local resource with, taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, freeMemCount: {%d}, freeBandwidthCount: {%d}, freeDiskCount: {%d}, freeProcessorCount: {%d}",
		taskId, partyId, jobNodeId, freeMemCount, freeBandwidthCount, freeDiskCount, freeProcessorCount)
	return nil
}

func (m *Manager) ReleaseLocalResourceWithTask(logdesc, taskId, partyId string, option ReleaseResourceOption) {

	log.Debugf("Start ReleaseLocalResourceWithTask %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}", logdesc, taskId, partyId, option)

	has, err := m.dataCenter.HasLocalTaskExecuteStatusByPartyId(taskId, partyId)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed to query local task exec status with task %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}",
			logdesc, taskId, partyId, option)
		return
	}

	if has {
		log.Debugf("The local task have been executing, don't `ReleaseLocalResourceWithTask` %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}",
			logdesc, taskId, partyId, option)
		return
	}

	if option.IsUnlockLocalResorce() {
		log.Debugf("start unlock local resource with task %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}",
			logdesc, taskId, partyId, option)
		if err := m.UnLockLocalResourceWithTask(taskId, partyId); nil != err {
			log.WithError(err).Warnf("Warning unlock local resource with task %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}",
				logdesc, taskId, partyId, option)
		}
	}

	if option.IsRemoveLocalTask() {

		has, err := m.dataCenter.HasLocalTaskExecuteStatusParty(taskId)
		// When tasks in current organization, including sender and other partners, do not have an 'executestatus' symbol.
		// It means that no one is handling the task
		if nil == err && !has {

			log.Debugf("start remove all things about this local task %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}",
				logdesc, taskId, partyId, option)

			// Remove the only task that everyone refers to together
			if err := m.dataCenter.RemoveLocalTask(taskId); nil != err {
				log.WithError(err).Errorf("Failed to remove local task  %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}",
					logdesc, taskId, partyId, option)
			}

			// Remove the only things in task that everyone refers to together
			if err := m.dataCenter.RemoveTaskPowerPartyIds(taskId); nil != err {
				log.WithError(err).Errorf("Failed to remove power's partyIds of local task  %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}",
					logdesc, taskId, partyId, option)
			}

			// Remove the partyId list of current task participants saved by the task sender
			if err := m.dataCenter.RemoveTaskPartnerPartyIds(taskId); nil != err {
				log.WithError(err).Errorf("Failed to remove handler partner's partyIds of local task  %s,taskId: {%s}, partyId: {%s}, releaseOption: {%d}",
					logdesc, taskId, partyId, option)
			}

			// Remove the task event of all partys
			if err := m.dataCenter.RemoveTaskEventList(taskId); nil != err {
				log.WithError(err).Errorf("Failed to clean all event list of task  %s, taskId: {%s}", logdesc, taskId)
			}
		}
	}

	if option.IsRemoveLocalTaskEvents() {
		log.Debugf("start remove party event list of task  %s, taskId: {%s}, partyId: {%s}", logdesc, taskId, partyId)
		if err := m.dataCenter.RemoveTaskEventListByPartyId(taskId, partyId); nil != err {
			log.WithError(err).Errorf("Failed to clean party event list of task  %s, taskId: {%s}, partyId: {%s}", logdesc, taskId, partyId)
		}
	}
}

func (m *Manager) addPartyTaskPowerUsedOnJobNode(used *types.LocalTaskPowerUsed) error {

	hasPowerUsed, err := m.dataCenter.HasLocalTaskPowerUsed(used.GetTaskId(), used.GetPartyId())
	if nil != err {
		log.WithError(err).Errorf("failed to call HasLocalTaskPowerUsed on addPartyTaskPowerUsedOnJobNode(), used: {%s}", used.String())
		return err
	}

	if !hasPowerUsed {
		if err := m.dataCenter.StoreLocalTaskPowerUsed(used); nil != err {
			log.WithError(err).Errorf("failed to call StoreLocalTaskPowerUsed on addPartyTaskPowerUsedOnJobNode(), used: {%s}", used.String())
			return err
		}
		log.Debugf("Succeed store powerUsed on addPartyTaskPowerUsedOnJobNode(), used: {%s}", used.String())
	}
	return nil
}

func (m *Manager) removePartyTaskPowerUsedOnJobNode(used *types.LocalTaskPowerUsed) error {
	hasPowerUsed, err := m.dataCenter.HasLocalTaskPowerUsed(used.GetTaskId(), used.GetPartyId())
	if nil != err {
		log.WithError(err).Errorf("failed to call HasLocalTaskPowerUsed on removePartyTaskPowerUsedOnJobNode(), used: {%s}", used.String())
		return err
	}

	if hasPowerUsed {
		if err := m.dataCenter.RemoveLocalTaskPowerUsed(used.GetTaskId(), used.GetPartyId()); nil != err {
			log.WithError(err).Errorf("failed to call RemoveLocalTaskPowerUsed on removePartyTaskPowerUsedOnJobNode(), used: {%s}", used.String())
			return err
		}
		log.Debugf("Succeed remove powerUsed on removePartyTaskPowerUsedOnJobNode(), used: {%s}", used.String())
	}
	return nil
}

func (m *Manager) StoreJobNodeExecuteTaskId (jobNodeId, taskId, partyId string) error {
	if err :=  m.dataCenter.StoreJobNodeTaskPartyId(jobNodeId, taskId, partyId); nil != err {
		log.WithError(err).Errorf("failed to call StoreJobNodeTaskPartyId on StoreJobNodeExecuteTaskId(), jobNodeId: {%s}, taskId: {%s}, partyId: {%s}",
			jobNodeId, taskId, partyId)
		return err
	}
	log.Debugf("Succeed store JobNodeId runningTask partyId on StoreJobNodeExecuteTaskId(), jobNodeId: {%s}, taskId: {%s}, partyId: {%s}",
		jobNodeId, taskId, partyId)

	hasHistoryTaskId, err := m.dataCenter.HasJobNodeHistoryTaskId(jobNodeId, taskId)
	if nil != err {
		log.WithError(err).Errorf("failed to check JobNode taskId whether exists on StoreJobNodeExecuteTaskId(), jobNodeId: {%s}, taskId: {%s}, partyId: {%s}",
			jobNodeId, taskId, partyId)
		return err
	}
	if !hasHistoryTaskId {
		if err := m.dataCenter.StoreJobNodeHistoryTaskId (jobNodeId, taskId); nil != err {
			log.WithError(err).Errorf("failed to inscrease JobNode history task count on StoreJobNodeExecuteTaskId(), jobNodeId: {%s}, taskId: {%s}, partyId: {%s}",
				jobNodeId, taskId, partyId)
			return err
		}
	}
	return nil
}

func (m *Manager) RemoveJobNodeExecuteTaskId (jobNodeId, taskId, partyId string) error {

	if err :=  m.dataCenter.RemoveJobNodeTaskPartyId(jobNodeId, taskId, partyId); nil != err {
		log.WithError(err).Errorf("failed to call RemoveJobNodeTaskPartyId on RemoveJobNodeExecuteTaskId(), jobNodeId: {%s}, taskId: {%s}, partyId: {%s}",
			jobNodeId, taskId, partyId)
		return err
	}
	log.Debugf("Succeed remove JobNodeId runningTask partyId on RemoveJobNodeExecuteTaskId(), jobNodeId: {%s}, taskId: {%s}, partyId: {%s}",
		jobNodeId, taskId, partyId)
	return nil
}

func (m *Manager) IsMockIdentityId (identityId string) bool {
	if _, ok := m.mockIdentityIdsCache[identityId]; ok {
		return true
	}
	return false
}


/// ======================  v 2.0
func (m *Manager) GetDB() core.CarrierDB { return m.dataCenter }
