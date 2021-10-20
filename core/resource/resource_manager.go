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
	// TODO 这里需要一个 config <SlotUnit 的>
	dataCenter  core.CarrierDB // Low level persistent database to store final content.
	slotUnit *types.Slot
	mockIdentityIdsFile  string
	mockIdentityIdsCache map[string]struct{}
}

func NewResourceManager(dataCenter core.CarrierDB, mockIdentityIdsFile string) *Manager {
	m := &Manager{
		dataCenter: dataCenter,
		//remoteTableQueue:    make([]*types.RemoteResourceTable, 0),
		slotUnit:            types.DefaultSlotUnit,
		mockIdentityIdsFile: mockIdentityIdsFile,   //TODO for test
		mockIdentityIdsCache: make(map[string]struct{}, 0),
	}

	return m
}

func (m *Manager) loop() {
}

func (m *Manager) Start() error {

	slotUnit, err := m.dataCenter.QueryNodeResourceSlotUnit()
	if nil != err {
		log.Warnf("Failed to load local slotUnit on resourceManager Start(), err: {%s}", err)
	} else {
		m.SetSlotUnit(slotUnit.Mem, slotUnit.Bandwidth, slotUnit.Processor)
	}

	// store slotUnit
	if err := m.dataCenter.StoreNodeResourceSlotUnit(m.slotUnit); nil != err {
		return err
	}

	// build mock identityIds cache
	if "" != m.mockIdentityIdsFile {
		var identityIdList []string
		if err := fileutil.LoadJSON(m.mockIdentityIdsFile, &identityIdList); err != nil {
			log.Errorf("Failed to load `--mock-identity-file` on Start resourceManager, file: {%s}, err: {%s}", m.mockIdentityIdsFile, err)
			return err
		}

		for _, iden := range identityIdList {
			m.mockIdentityIdsCache[iden] = struct{}{}
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
	log.Infof("Stopped resource manager ...")
	return nil
}

func (m *Manager) SetSlotUnit(mem, b uint64, p uint32) {
	m.slotUnit = &types.Slot{
		Mem:       mem,
		Processor: p,
		Bandwidth: b,
	}
}
func (m *Manager) GetSlotUnit() *types.Slot { return m.slotUnit }

func (m *Manager) UseSlot(nodeId string, slotCount uint32) error {
	table, err := m.QueryLocalResourceTable(nodeId)
	if nil != err {
		return fmt.Errorf("No found the resource table of node: %s, %s", nodeId, err)
	}
	if err := table.UseSlot(slotCount); nil != err {
		return err
	}
	return m.StoreLocalResourceTable(table)
}
func (m *Manager) FreeSlot(nodeId string, slotCount uint32) error {
	table, err := m.QueryLocalResourceTable(nodeId)
	if nil != err {
		return fmt.Errorf("No found the resource table of node: %s, %s", nodeId, err)
	}
	if err := table.FreeSlot(slotCount); nil != err {
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

func (m *Manager) LockLocalResourceWithTask(partyId, jobNodeId string, needSlotCount uint64, task *types.Task) error {

	log.Infof("Start lock local resource with taskId {%s}, partyId: {%s}, jobNodeId {%s}, slotCount {%d}", task.GetTaskId(), partyId, jobNodeId, needSlotCount)

	// Lock local resource (jobNode)
	if err := m.UseSlot(jobNodeId, uint32(needSlotCount)); nil != err {
		log.Errorf("Failed to lock internal power resource on resourceManager.LockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, usedSlotCount: {%d}, err: {%s}",
			task.GetTaskId(), partyId, jobNodeId, needSlotCount, err)
		return err
	}

	used := types.NewLocalTaskPowerUsed(task.GetTaskId(), partyId, jobNodeId, needSlotCount)
	if err := m.addPartyTaskPowerUsedOnJobNode(used); nil != err {

		m.FreeSlot(jobNodeId, uint32(needSlotCount))

		return err
	}

	// 更新本地 resource 资源信息 [添加资源使用情况]
	jobNodeResource, err := m.dataCenter.QueryLocalResource(jobNodeId)
	if nil != err {

		m.FreeSlot(jobNodeId, uint32(needSlotCount))
		m.removePartyTaskPowerUsedOnJobNode(used)

		log.Errorf("Failed to query local jobNodeResource on resourceManager.LockLocalResourceWithTask(), taskId: {%s}, jobNodeId: {%s}, usedSlotCount: {%d}, err: {%s}",
			task.GetTaskId(), jobNodeId, needSlotCount, err)
		return err
	}

	jobNodeRunningTaskCount, err := m.dataCenter.QueryRunningTaskCountOnJobNode(jobNodeId)
	if nil != err {

		m.FreeSlot(jobNodeId, uint32(needSlotCount))
		m.removePartyTaskPowerUsedOnJobNode(used)

		log.Errorf("Failed to query task runningCount in jobNode on resourceManager.LockLocalResourceWithTask(), taskId: {%s}, jobNodeId: {%s}, usedSlotCount: {%d}, err: {%s}",
			task.GetTaskId(), jobNodeId, needSlotCount, err)
		return err
	}

	// 更新 本地 jobNodeResource 的资源使用信息
	usedMem := m.slotUnit.Mem * needSlotCount
	usedProcessor := m.slotUnit.Processor * uint32(needSlotCount)
	usedBandwidth := m.slotUnit.Bandwidth * needSlotCount

	jobNodeResource.GetData().UsedMem += usedMem
	jobNodeResource.GetData().UsedProcessor +=usedProcessor
	jobNodeResource.GetData().UsedBandwidth += usedBandwidth
	if jobNodeRunningTaskCount > 0 {
		jobNodeResource.GetData().State = apicommonpb.PowerState_PowerState_Occupation
	}
	if err := m.dataCenter.InsertLocalResource(jobNodeResource); nil != err {

		m.FreeSlot(jobNodeId, uint32(needSlotCount))
		m.removePartyTaskPowerUsedOnJobNode(used)

		log.Errorf("Failed to update local jobNodeResource on resourceManager.LockLocalResourceWithTask(), taskId: {%s}, jobNodeId: {%s}, usedSlotCount: {%d}, err: {%s}",
			task.GetTaskId(), jobNodeId, needSlotCount, err)
		return err
	}

	// 还需要 将资源使用实况 实时上报给  dataCenter  [添加资源使用情况]
	if err := m.dataCenter.SyncPowerUsed(jobNodeResource); nil != err {
		log.Errorf("Failed to sync jobNodeResource to dataCenter on resourceManager.LockLocalResourceWithTask(), taskId: {%s}, jobNodeId: {%s}, usedSlotCount: {%d}, err: {%s}",
			task.GetTaskId(), jobNodeId, needSlotCount, err)
		return err
	}

	log.Infof("Finished lock local resource with taskId {%s}, jobNodeId {%s}, slotCount {%d}", task.GetTaskId(), jobNodeId, needSlotCount)
	return nil
}

func (m *Manager) UnLockLocalResourceWithTask(taskId, partyId string) error {

	used, err := m.dataCenter.QueryLocalTaskPowerUsed(taskId, partyId)
	if nil != err {
		log.Errorf("Failed to query local task powerUsed on resourceManager.UnLockLocalResourceWithTask(), taskId {%s}, partyId: {%s}, err: {%s}", taskId, partyId, err)
		return err
	}

	jobNodeId := used.GetNodeId()
	freeSlotUnitCount := used.GetSlotCount()

	log.Infof("Start unlock local resource on resourceManager.UnLockLocalResourceWithTask(), taskId {%s}, partyId: {%s}, jobNodeId {%s}, slotCount {%d}", taskId, partyId, jobNodeId, freeSlotUnitCount)

	// Unlock local resource (jobNode)
	if err := m.FreeSlot(used.GetNodeId(), uint32(freeSlotUnitCount)); nil != err {
		log.Errorf("Failed to unlock internal power resource on resourceManager.UnLockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, freeSlotUnitCount: {%d}, err: {%s}",
			taskId, partyId, jobNodeId, freeSlotUnitCount, err)
		return err
	}

	if err := m.removePartyTaskPowerUsedOnJobNode(used); nil != err {
		return err
	}

	// 移除partyId 对应的本地 该task 的 ResourceUsage
	if err := m.dataCenter.RemoveTaskResuorceUsage(taskId, partyId); nil != err {
		log.Errorf("Failed to remove local task resourceUsage on resourceManager.UnLockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, freeSlotUnitCount: {%d}, err: {%s}",
			taskId, partyId, jobNodeId, freeSlotUnitCount, err)
		return err
	}

	// 更新本地 resource 资源信息 [释放资源使用情况]
	jobNodeResource, err := m.dataCenter.QueryLocalResource(jobNodeId)
	if nil != err {
		log.Errorf("Failed to query local jobNodeResource on resourceManager.UnLockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, freeSlotUnitCount: {%d}, err: {%s}",
			taskId, partyId, jobNodeId, freeSlotUnitCount, err)
		return err
	}

	jobNodeRunningTaskCount, err := m.dataCenter.QueryRunningTaskCountOnJobNode(jobNodeId)
	if nil != err {
		log.Errorf("Failed to query task runningCount in jobNode on resourceManager.UnLockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, freeSlotUnitCount: {%d}, err: {%s}",
			taskId, partyId, jobNodeId, freeSlotUnitCount, err)
		return err
	}

	// 更新 本地 jobNodeResource 的资源使用信息
	usedMem := m.slotUnit.Mem * freeSlotUnitCount
	usedProcessor := m.slotUnit.Processor * uint32(freeSlotUnitCount)
	usedBandwidth := m.slotUnit.Bandwidth * freeSlotUnitCount

	jobNodeResource.GetData().UsedMem -= usedMem
	jobNodeResource.GetData().UsedProcessor -= usedProcessor
	jobNodeResource.GetData().UsedBandwidth -= usedBandwidth
	if jobNodeRunningTaskCount == 0 {
		jobNodeResource.GetData().State = apicommonpb.PowerState_PowerState_Released
	}
	if err := m.dataCenter.InsertLocalResource(jobNodeResource); nil != err {
		log.Errorf("Failed to update local jobNodeResource on resourceManager.UnLockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, freeSlotUnitCount: {%d}, err: {%s}",
			taskId, partyId, jobNodeId, freeSlotUnitCount, err)
		return err
	}

	// 还需要 将资源使用实况 实时上报给  dataCenter  [释放资源使用情况]
	if err := m.dataCenter.SyncPowerUsed(jobNodeResource); nil != err {
		log.Errorf("Failed to sync jobNodeResource to dataCenter on resourceManager.UnLockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, freeSlotUnitCount: {%d}, err: {%s}",
			taskId, partyId, jobNodeId, freeSlotUnitCount, err)
		return err
	}

	log.Infof("Finished unlock local resource with taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, freeSlotUnitCount: {%d}",
		taskId, partyId, jobNodeId, freeSlotUnitCount)
	return nil
}

func (m *Manager) ReleaseLocalResourceWithTask(logdesc, taskId, partyId string, option ReleaseResourceOption) {

	log.Debugf("Start ReleaseLocalResourceWithTask %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}", logdesc, taskId, partyId, option)

	has, err := m.dataCenter.HasLocalTaskExecuteStatusByPartyId(taskId, partyId)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.Errorf("Failed to query local task exec status with task %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}, err: {%s}",
			logdesc, taskId, partyId, option, err)
		return
	}

	if has {
		log.Debugf("The local task have been executing, don't `ReleaseLocalResourceWithTask` %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}",
			logdesc, taskId, partyId, option)
		return
	}


	//used, err := m.dataCenter.QueryLocalTaskPowerUsed(taskId, partyId)
	//if nil != err {
	//	log.Errorf("Failed to query local task powerUsed,taskId {%s}, partyId: {%s}, err: {%s}", taskId, partyId, err)
	//	return
	//}
	//// query partyId count on jobNode with jobNodeId and taskId.
	//count, err := m.dataCenter.QueryResourceTaskPartyIdCount(used.GetNodeId(), used.GetTaskId())
	//if nil != err {
	//	log.Errorf("failed to query resuorce task party count, used: {%s}, err: {%s}", used.String(), err)
	//	return
	//}

	if option.IsUnlockLocalResorce() {
		log.Debugf("start unlock local resource with task %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}",
			logdesc, taskId, partyId, option)
		if err := m.UnLockLocalResourceWithTask(taskId, partyId); nil != err {
			log.Errorf("Failed to unlock local resource with task %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}, err: {%s}",
				logdesc, taskId, partyId, option, err)
		}
	}

	if option.IsRemoveLocalTask() {
		log.Debugf("start remove local task %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}, err: {%s}",
			logdesc, taskId, partyId, option, err)
		// When tasks in current organization, including sender and other partners, do not have an 'executestatus' symbol.
		has, err := m.dataCenter.HasLocalTaskExecuteStatusParty(taskId)
		if nil == err && !has {
			if err := m.dataCenter.RemoveLocalTask(taskId); nil != err {
				log.Errorf("Failed to remove local task  %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}, err: {%s}",
					logdesc, taskId, partyId, option, err)
			}
			if err := m.dataCenter.RemoveTaskPowerPartyIds(taskId); nil != err {
				log.Errorf("Failed to remove power's partyIds of local task  %s, taskId: {%s}, last partyId: {%s}, releaseOption: {%d}, err: {%s}",
					logdesc, taskId, partyId, option, err)
			}
			if err := m.dataCenter.RemoveTaskPartnerPartyIds(taskId); nil != err {
				log.Errorf("Failed to remove handler partner's partyIds of local task  %s, taskId: {%s}, last partyId: {%s}, releaseOption: {%d}, err: {%s}",
					logdesc, taskId, partyId, option, err)
			}
			if err := m.dataCenter.RemoveTaskEventListByPartyId(taskId, partyId); nil != err {
				log.WithError(err).Errorf("Failed to clean all event list of task  %s, taskId: {%s}, partyId: {%s}", logdesc, taskId, partyId)
			}
		}
	}

	if option.IsRemoveLocalTaskEvents() {
		log.Debugf("start clean party event list of task  %s, taskId: {%s}, partyId: {%s}", logdesc, taskId, partyId)
		if err := m.dataCenter.RemoveTaskEventListByPartyId(taskId, partyId); nil != err {
			log.WithError(err).Errorf("Failed to clean party event list of task  %s, taskId: {%s}, partyId: {%s}", logdesc, taskId, partyId)
		}
	}
}

func (m *Manager) addPartyTaskPowerUsedOnJobNode(used *types.LocalTaskPowerUsed) error {
	has, err := m.dataCenter.HasLocalTaskPowerUsed(used.GetTaskId(), used.GetPartyId())
	if nil != err {
		log.Errorf("failed to call has local task powerUsed, used: {%s}, err: {%s}",
			used.String(), err)
		return err
	}
	if !has {
		count, err := m.dataCenter.QueryResourceTaskPartyIdCount(used.GetNodeId(), used.GetTaskId())
		if nil != err {
			log.Errorf("failed to query resuorce task party count, used: {%s}, err: {%s}",
				used.String(), err)
			return err
		}

		if err := m.dataCenter.IncreaseResourceTaskPartyIdCount(used.GetNodeId(), used.GetTaskId()); nil != err {

			log.Errorf("failed to increase resource task party count, used: {%s}, err: {%s}",
				used.String(), err)
			return err
		}

		if count == 0 {
			if err := m.dataCenter.StoreJobNodeRunningTaskId(used.GetNodeId(), used.GetTaskId()); nil != err {

				m.dataCenter.DecreaseResourceTaskPartyIdCount(used.GetNodeId(), used.GetTaskId())

				log.Errorf("failed to store local taskId and jobNodeId index, used: {%s}, err: {%s}",
					used.String(), err)
				return err
			}

			if err := m.dataCenter.IncreaseResourceTaskTotalCount(used.GetNodeId()); nil != err {

				m.dataCenter.DecreaseResourceTaskPartyIdCount(used.GetNodeId(), used.GetTaskId())
				m.dataCenter.RemoveJobNodeRunningTaskId(used.GetNodeId(), used.GetTaskId())

				log.Errorf("failed to increase taskTotalCount on jobNode, used: {%s}, err: {%s}",
					used.String(), err)
				return err
			}

		}

		if err := m.dataCenter.StoreLocalTaskPowerUsed(used); nil != err {

			m.dataCenter.DecreaseResourceTaskPartyIdCount(used.GetNodeId(), used.GetTaskId())
			if count == 0 {
				m.dataCenter.RemoveJobNodeRunningTaskId(used.GetNodeId(), used.GetTaskId())
			}

			log.Errorf("failed to store local taskId use jobNode slot, used: {%s}, err: {%s}",
				used.String(), err)
			return err
		}
	}

	return nil

}

func (m *Manager) removePartyTaskPowerUsedOnJobNode(used *types.LocalTaskPowerUsed) error {
	has, err := m.dataCenter.HasLocalTaskPowerUsed(used.GetTaskId(), used.GetPartyId())
	if nil != err {
		log.Errorf("failed to call has local task powerUsed, used: {%s}, err: {%s}",
			used.String(), err)
		return err
	}

	if has {

		oldCount, err := m.dataCenter.QueryResourceTaskPartyIdCount(used.GetNodeId(), used.GetTaskId())
		if nil != err {
			log.Errorf("failed to query resuorce task party count, used: {%s}, err: {%s}",
				used.String(), err)
			return err
		}
		if oldCount != 0 {
			if err := m.dataCenter.DecreaseResourceTaskPartyIdCount(used.GetNodeId(), used.GetTaskId()); nil != err {

				log.Errorf("failed to decrease resource task party count, used: {%s}, err: {%s}",
					used.String(), err)
				return err
			}
		}

		newCount, err := m.dataCenter.QueryResourceTaskPartyIdCount(used.GetNodeId(), used.GetTaskId())
		if nil != err {
			log.Errorf("failed to query resuorce task party count, used: {%s}, err: {%s}",
				used.String(), err)
			return err
		}

		if oldCount != 0 && newCount == 0 {

			if err := m.dataCenter.RemoveTaskEventList(used.GetTaskId()); nil != err {
				log.Errorf("failed to remove local task event list, used: {%s}, err: {%s}",
					used.String(), err)
				return err
			}

			if err := m.dataCenter.RemoveJobNodeRunningTaskId(used.GetNodeId(), used.GetTaskId()); nil != err {

				log.Errorf("failed to remove local taskId and jobNodeId index, used: {%s}, err: {%s}",
					used.String(), err)
				return err
			}
		}

		if err := m.dataCenter.RemoveLocalTaskPowerUsed(used.GetTaskId(), used.GetPartyId()); nil != err {

			log.Errorf("failed to remove local taskId use jobNode slot, used: {%s}, err: {%s}",
				used.String(), err)
			return err
		}
	}

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
