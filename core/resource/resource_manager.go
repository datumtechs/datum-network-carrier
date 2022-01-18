package resource

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/fileutil"
	"github.com/RosettaFlow/Carrier-Go/core"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/service/discovery"
	"github.com/RosettaFlow/Carrier-Go/types"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

const (
	defaultRefreshOrgResourceInterval = 30 * time.Second
)

type Manager struct {
	dataCenter           core.CarrierDB // Low level persistent database to store final content.
	mockIdentityIdsFile  string
	mockIdentityIdsCache map[string]struct{}
	resourceClientSet    *grpclient.InternalResourceClientSet // internal resource node set (Fighter node grpc client set)
}

func NewResourceManager(dataCenter core.CarrierDB, resourceClientSet *grpclient.InternalResourceClientSet, mockIdentityIdsFile string) *Manager {
	m := &Manager{
		dataCenter:           dataCenter,
		resourceClientSet:    resourceClientSet,
		mockIdentityIdsFile:  mockIdentityIdsFile, //TODO for test
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
	return m.dataCenter.StoreLocalResourceTable(table)
}
func (m *Manager) FreeSlot(nodeId string, mem, bandwidth, disk uint64, processor uint32) error {
	table, err := m.QueryLocalResourceTable(nodeId)
	if nil != err {
		return fmt.Errorf("No found the resource table of node: %s, %s", nodeId, err)
	}
	if err := table.FreeSlot(mem, bandwidth, disk, processor); nil != err {
		return err
	}
	return m.dataCenter.StoreLocalResourceTable(table)
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

	// Update local resource resource information [increase resource usage]
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

	// Update the resource usage information of the local jobNode resource
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

	// Report resource usage to datacenter in real time [increase resource usage]
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

	if err := m.removePartyTaskPowerUsedOnJobNode(used); nil != err {
		log.WithError(err).Errorf("Failed to remove partyTaskPowerUsed on resourceManager.UnLockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, freeMemCount: {%d}, freeBandwidthCount: {%d}, freeDiskCount: {%d}, freeProcessorCount: {%d}",
			taskId, partyId, jobNodeId, freeMemCount, freeBandwidthCount, freeDiskCount, freeProcessorCount)
		return err
	}

	// Unlock local resource (jobNode)
	if err := m.FreeSlot(used.GetNodeId(), freeMemCount, freeBandwidthCount, freeDiskCount, freeProcessorCount); nil != err {
		log.WithError(err).Errorf("Failed to freeSlot withJobNodeId on resourceManager.UnLockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, freeMemCount: {%d}, freeBandwidthCount: {%d}, freeDiskCount: {%d}, freeProcessorCount: {%d}",
			taskId, partyId, jobNodeId, freeMemCount, freeBandwidthCount, freeDiskCount, freeProcessorCount)
		return err
	}

	// Update local resource resource information [release resource usage]
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

	// Report resource usage to datacenter in real time [release resource usage]
	if err := m.dataCenter.SyncPowerUsed(jobNodeResource); nil != err {
		log.WithError(err).Errorf("Failed to sync jobNodeResource to dataCenter on resourceManager.UnLockLocalResourceWithTask(), taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, freeMemCount: {%d}, freeBandwidthCount: {%d}, freeDiskCount: {%d}, freeProcessorCount: {%d}",
			taskId, partyId, jobNodeId, freeMemCount, freeBandwidthCount, freeDiskCount, freeProcessorCount)
		return err
	}

	log.Infof("Finished unlock local resource with, taskId: {%s}, partyId: {%s}, jobNodeId: {%s}, freeMemCount: {%d}, freeBandwidthCount: {%d}, freeDiskCount: {%d}, freeProcessorCount: {%d}",
		taskId, partyId, jobNodeId, freeMemCount, freeBandwidthCount, freeDiskCount, freeProcessorCount)
	return nil
}

func (m *Manager) ReleaseLocalResourceWithTask(logdesc, taskId, partyId string, option ReleaseResourceOption, isSender bool) {

	log.Debugf("Start ReleaseLocalResourceWithTask %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}, isSender: {%v}", logdesc, taskId, partyId, option, isSender)

	has, err := m.dataCenter.HasLocalTaskExecuteStatusByPartyId(taskId, partyId)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed to query local task exec status with task %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}, isSender: {%v}",
			logdesc, taskId, partyId, option, isSender)
		return
	}

	if has {
		log.Debugf("The local task have been executing, don't `ReleaseLocalResourceWithTask` %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}, isSender: {%v}",
			logdesc, taskId, partyId, option, isSender)
		return
	}

	if option.IsUnlockLocalResorce() {
		log.Debugf("start unlock local resource with task %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}, isSender: {%v}",
			logdesc, taskId, partyId, option, isSender)
		if err := m.UnLockLocalResourceWithTask(taskId, partyId); nil != err {
			log.WithError(err).Warnf("Warning unlock local resource with task %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}, isSender: {%v}",
				logdesc, taskId, partyId, option, isSender)
		}
	}

	if option.IsRemoveLocalTask() {

		var removeAll bool

		has, err := m.dataCenter.HasLocalTaskExecuteStatusParty(taskId)
		if nil != err {
			log.WithError(err).Errorf("Failed to call HasLocalTaskExecuteStatusParty(), when remove all things about this local task %s, taskId: {%s}, partyId: {%s}, releaseOption: {%d}, isSender: {%v}",
				logdesc, taskId, partyId, option, isSender)
		}
		// When tasks in current organization, including sender and other partners, do not have an 'executestatus' symbol.
		// It means that no one is handling the task
		if !isSender && (nil == err && !has) {
			removeAll = true
		}
		if isSender {
			removeAll = true
		}

		if removeAll {

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
		log.Debugf("start remove party event list of task  %s, taskId: {%s}, partyId: {%s}, isSender: {%v}", logdesc, taskId, partyId, isSender)
		if err := m.dataCenter.RemoveTaskEventListByPartyId(taskId, partyId); nil != err {
			log.WithError(err).Errorf("Failed to clean party event list of task  %s, taskId: {%s}, partyId: {%s}, isSender: {%v}", logdesc, taskId, partyId, isSender)
		}
	}
}

func (m *Manager) UnLockLocalResourceWithJobNodeId(jobNodeId string) error {

	taskIdsAndPartyIdsPairs, err := m.dataCenter.QueryJobNodeRunningTaskIdsAndPartyIdsPairs(jobNodeId)
	if nil != err {
		log.WithError(err).Errorf("Failed to query jobNode running taskIds AND partyIds pairs on resourceManager.UnLockLocalResourceWithJobNodeId(), jobNodeId: {%s}", jobNodeId)
		return err
	}
	for taskId, partyIds := range taskIdsAndPartyIdsPairs {
		for _, partyId := range partyIds {
			if err := m.UnLockLocalResourceWithTask(taskId, partyId); nil != err {
				log.WithError(err).Errorf("Warning unlock local resource on old jobNode on resourceManager.UnLockLocalResourceWithJobNodeId(), jobNodeId: {%s}, taskId: {%s}, partyId: {%s}",
					jobNodeId, taskId, partyId)
				continue
			}
		}
	}
	return nil
}

func (m *Manager) UnLockLocalResourceWithPowerId(powerId string) error {
	jobNodeId, err := m.dataCenter.QueryJobNodeIdByPowerId(powerId)
	if nil != err {
		log.WithError(err).Errorf("Failed to query jobNodeId with powerId on resourceManager.UnLockLocalResourceWithPowerId(), powerId: {%s}", powerId)
		return err
	}
	return m.UnLockLocalResourceWithJobNodeId(jobNodeId)
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

func (m *Manager) StoreJobNodeExecuteTaskId(jobNodeId, taskId, partyId string) error {
	if err := m.dataCenter.StoreJobNodeTaskPartyId(jobNodeId, taskId, partyId); nil != err {
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
		if err := m.dataCenter.StoreJobNodeHistoryTaskId(jobNodeId, taskId); nil != err {
			log.WithError(err).Errorf("failed to inscrease JobNode history task count on StoreJobNodeExecuteTaskId(), jobNodeId: {%s}, taskId: {%s}, partyId: {%s}",
				jobNodeId, taskId, partyId)
			return err
		}
	}
	return nil
}

func (m *Manager) RemoveJobNodeExecuteTaskId(jobNodeId, taskId, partyId string) error {

	if err := m.dataCenter.RemoveJobNodeTaskPartyId(jobNodeId, taskId, partyId); nil != err {
		log.WithError(err).Errorf("failed to call RemoveJobNodeTaskPartyId on RemoveJobNodeExecuteTaskId(), jobNodeId: {%s}, taskId: {%s}, partyId: {%s}",
			jobNodeId, taskId, partyId)
		return err
	}
	log.Debugf("Succeed remove JobNodeId runningTask partyId on RemoveJobNodeExecuteTaskId(), jobNodeId: {%s}, taskId: {%s}, partyId: {%s}",
		jobNodeId, taskId, partyId)
	return nil
}

func (m *Manager) IsMockIdentityId(identityId string) bool {
	if _, ok := m.mockIdentityIdsCache[identityId]; ok {
		return true
	}
	return false
}

/// ======================  v 0.2.0
func (m *Manager) GetDB() core.CarrierDB { return m.dataCenter }

/// ======================  v 0.3.0
func (m *Manager) HasNotInternalJobNodeClientSet() bool { return !m.HasInternalJobNodeClientSet() }
func (m *Manager) HasInternalJobNodeClientSet() bool {
	if nil == m.resourceClientSet || 0 == m.resourceClientSet.JobNodeClientSize() {
		return false
	}
	return true
}

func (m *Manager) HasNotInternalDataNodeClientSet() bool { return !m.HasInternalDataNodeClientSet() }
func (m *Manager) HasInternalDataNodeClientSet() bool {
	if nil == m.resourceClientSet || 0 == m.resourceClientSet.DataNodeClientSize() {
		return false
	}
	return true
}

func (m *Manager) StoreJobNodeClient(nodeId string, client *grpclient.JobNodeClient) {
	m.resourceClientSet.StoreJobNodeClient(nodeId, client)
}

func (m *Manager) QueryJobNodeClient(nodeId string) (*grpclient.JobNodeClient, bool) {
	return m.resourceClientSet.QueryJobNodeClient(nodeId)
}

func (m *Manager) QueryJobNodeClients() []*grpclient.JobNodeClient {
	return m.resourceClientSet.QueryJobNodeClients()
}

func (m *Manager) RemoveJobNodeClient(nodeId string) {
	m.resourceClientSet.RemoveJobNodeClient(nodeId)
}

func (m *Manager) JobNodeClientSize() int {
	return m.resourceClientSet.JobNodeClientSize()
}

func (m *Manager) StoreDataNodeClient(nodeId string, client *grpclient.DataNodeClient) {
	m.resourceClientSet.StoreDataNodeClient(nodeId, client)
}

func (m *Manager) QueryDataNodeClient(nodeId string) (*grpclient.DataNodeClient, bool) {
	return m.resourceClientSet.QueryDataNodeClient(nodeId)
}

func (m *Manager) QueryDataNodeClients() []*grpclient.DataNodeClient {
	return m.resourceClientSet.QueryDataNodeClients()
}

func (m *Manager) RemoveDataNodeClient(nodeId string) {
	m.resourceClientSet.RemoveDataNodeClient(nodeId)
}

func (m *Manager) DataNodeClientSize() int {
	return m.resourceClientSet.DataNodeClientSize()
}

func (m *Manager) AddDiscoveryJobNodeResource(identity *apicommonpb.Organization, jobNodeId, jobNodeIP, jobNodePort, jobNodeExternalIP, jobNodeExternalPort string) error {

	log.Infof("Discovered a new jobNode from consul server, add jobNode resource on resourceManager.AddDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
		jobNodeId, jobNodeIP, jobNodePort)

	client, err := grpclient.NewJobNodeClient(context.Background(), fmt.Sprintf("%s:%s", jobNodeIP, jobNodePort), jobNodeId)
	if nil != err {
		log.WithError(err).Errorf("Failed to connect new jobNode on resourceManager.AddDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
			jobNodeId, jobNodeIP, jobNodePort)
		return err
	}
	jobNodeStatus, err := client.GetStatus()
	if nil != err {
		log.WithError(err).Errorf("Failed to connect jobNode to query status on resourceManager.AddDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
			jobNodeId, jobNodeIP, jobNodePort)
		return err
	}

	// 1. add local jobNode resource
	// add resource usage first, but not own power now (mem, proccessor, bandwidth)
	// store into local db
	if err := m.dataCenter.InsertLocalResource(types.NewLocalResource(&libtypes.LocalResourcePB{
		IdentityId: identity.GetIdentityId(),
		NodeId:     identity.GetNodeId(),
		NodeName:   identity.GetNodeName(),
		JobNodeId:  jobNodeId,
		DataId:     "", // can not own powerId now, because power have not publish
		// the status of data, N means normal, D means deleted.
		DataStatus: apicommonpb.DataStatus_DataStatus_Normal,
		// resource status, eg: create/release/revoke
		State: apicommonpb.PowerState_PowerState_Created,
		// unit: byte
		TotalMem: jobNodeStatus.GetTotalMemory(),
		UsedMem:  0,
		// number of cpu cores.
		TotalProcessor: jobNodeStatus.GetTotalCpu(),
		UsedProcessor:  0,
		// unit: byte
		TotalBandwidth: jobNodeStatus.GetTotalBandwidth(),
		UsedBandwidth:  0,
		TotalDisk:      jobNodeStatus.GetTotalDisk(),
		UsedDisk:       0,
	})); nil != err {
		log.WithError(err).Errorf("Failed to store jobNode local resource on resourceManager.AddDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
			jobNodeId, jobNodeIP, jobNodePort)
		return err
	}

	// 2. add rpc client
	m.StoreJobNodeClient(jobNodeId, client)

	// 3. add local jobNode info
	// build new jobNode info that was need to store local db
	if err = m.dataCenter.SetRegisterNode(pb.PrefixTypeJobNode,
		&pb.YarnRegisteredPeerDetail{
			Id:           strings.Join([]string{discovery.JobNodeConsulServiceIdPrefix, jobNodeIP, jobNodePort}, discovery.ConsulServiceIdSeparator),
			InternalIp:   jobNodeIP,
			InternalPort: jobNodePort,
			ExternalIp:   jobNodeExternalIP,
			ExternalPort: jobNodeExternalPort,
			ConnState:    pb.ConnState_ConnState_Connected,
		}); nil != err {
		log.WithError(err).Errorf("Failed to store registerNode into local db on resourceManager.AddDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
			jobNodeId, jobNodeIP, jobNodePort)
		return err
	}

	log.Infof("Succeed add a new jobNode resource from consul server, add jobNode resource on resourceManager.AddDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
		jobNodeId, jobNodeIP, jobNodePort)
	return nil
}

func (m *Manager) UpdateDiscoveryJobNodeResource(identity *apicommonpb.Organization, jobNodeId, jobNodeIP, jobNodePort, jobNodeExternalIP, jobNodeExternalPort string, old *pb.YarnRegisteredPeerDetail) error {

	// check the  via external ip and port comparing old infomation,
	// if it is, update the some things about jobNode.
	if old.GetExternalIp() != jobNodeExternalIP || old.GetExternalPort() != jobNodeExternalPort {

		oldIp := old.GetExternalIp()
		oldPort := old.GetExternalPort()

		// update jobNode info that was need to store local db
		old.ExternalIp = jobNodeExternalIP
		old.ExternalPort = jobNodeExternalPort
		// 1. update local jobNode info
		// update jobNode ip port into local db
		if err := m.dataCenter.SetRegisterNode(pb.PrefixTypeJobNode, old); nil != err {
			log.WithError(err).Errorf("Failed to update jobNode into local db on resourceManager.UpdateDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
				jobNodeId, jobNodeIP, jobNodePort)
			return err
		}

		log.Infof("Succeed update a old jobNode external ip and port from consul server on resourceManager.UpdateDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}, old externalIp: {%s}, old externalPort: {%s}, new externalIp: {%s}, new externalPort: {%s}",
			jobNodeId, jobNodeIP, jobNodePort, oldIp, oldPort, old.GetExternalIp(), old.GetExternalPort())
	}
	// check connection status,
	// if it be changed, update the connState value about jobNode
	if old.GetConnState() != pb.ConnState_ConnState_Connected {

		old.ConnState = pb.ConnState_ConnState_Connected
		if err := m.dataCenter.SetRegisterNode(pb.PrefixTypeJobNode, old); nil != err {
			log.WithError(err).Errorf("Failed to update jobNode into local db on resourceManager.UpdateDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
				jobNodeId, jobNodeIP, jobNodePort)
			return err
		}

		log.Infof("Succeed update jobNode ConnState to `connected` on resourceManager.UpdateDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
			jobNodeId, jobNodeIP, jobNodePort)
	}

	resourceTable, err := m.dataCenter.QueryLocalResourceTable(jobNodeId)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed to query local power resource on old jobNode on resourceManager.UpdateDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
			jobNodeId, jobNodeIP, jobNodePort)
		return err
	}
	if nil != resourceTable && !resourceTable.GetAlive() {
		resourceTable.SetAlive(true)
		if err := m.dataCenter.StoreLocalResourceTable(resourceTable); nil != err {
			log.WithError(err).Errorf("Failed to update alive flag of local jobNode resource on resourceManager.UpdateDiscoveryJobNodeResource(), powerId: {%s}, jobNodeId: {%s}",
				resourceTable.GetPowerId(), jobNodeId)
			return err
		}
	}

	// check jobNode resource wether have change?
	// query local resource
	resource, err := m.dataCenter.QueryLocalResource(jobNodeId)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed to query local resource on resourceManager.UpdateDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
			jobNodeId, jobNodeIP, jobNodePort)
		return err
	}
	if nil != resource {
		client, ok := m.QueryJobNodeClient(jobNodeId)
		if !ok {
			log.WithError(err).Errorf("can not find jobNode rpc client on resourceManager.UpdateDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
				jobNodeId, jobNodeIP, jobNodePort)
			return err
		}
		if client.IsNotConnected() {
			if err := client.Reconnect(); nil != err {
				log.WithError(err).Errorf("Failed to connect internal jobNode on resourceManager.UpdateDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
					jobNodeId, jobNodeIP, jobNodePort)
				return err
			}
		}
		jobNodeStatus, err := client.GetStatus()
		if nil != err {
			log.WithError(err).Errorf("Failed to connect jobNode to query status on resourceManager.UpdateDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
				jobNodeId, jobNodeIP, jobNodePort)
			return err
		}

		// update jobNode local resource total
		if jobNodeStatus.GetTotalBandwidth() != resource.GetData().GetTotalBandwidth() ||
			jobNodeStatus.GetTotalCpu() != resource.GetData().GetTotalProcessor() ||
			jobNodeStatus.GetTotalMemory() != resource.GetData().GetTotalMem() ||
			jobNodeStatus.GetUsedDisk() != resource.GetData().GetTotalDisk() {

			resource.GetData().TotalBandwidth = jobNodeStatus.GetTotalBandwidth()
			resource.GetData().TotalMem = jobNodeStatus.GetTotalMemory()
			resource.GetData().TotalProcessor = jobNodeStatus.GetTotalCpu()
			resource.GetData().TotalDisk = jobNodeStatus.GetTotalDisk()

			if err := m.dataCenter.InsertLocalResource(resource); nil != err {
				log.WithError(err).Errorf("Failed to update local resource when jobNode total resource had change from consul server on resourceManager.UpdateDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
					jobNodeId, jobNodeIP, jobNodePort)
				return err
			}
			log.Infof("Succeed update jobNode local total resource on resourceManager.UpdateDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
				jobNodeId, jobNodeIP, jobNodePort)
		}
	}
	return nil
}

func (m *Manager) RemoveDiscoveryJobNodeResource(identity *apicommonpb.Organization, jobNodeId, jobNodeIP, jobNodePort, jobNodeExternalIP, jobNodeExternalPort string, old *pb.YarnRegisteredPeerDetail) error {

	log.Infof("Disappeared a old jobNode from consul server, jobNodeId: {%s}", jobNodeId)


	resourceTable, err := m.dataCenter.QueryLocalResourceTable(jobNodeId)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed to query local power resource of old jobNode on resourceManager.RemoveDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
			jobNodeId, jobNodeIP, jobNodePort)
		return err
	}
	if nil != resourceTable && resourceTable.GetAlive() {
		log.Warnf("still have the published computing power information by the jobNode on resourceManager.RemoveDiscoveryJobNodeResource(), %s",
			resourceTable.String())
		// ##############################
		// A. update alive status of resource table about jobNode
		// ##############################

		// 1. update resource table
		resourceTable.SetAlive(false)
		if err := m.dataCenter.StoreLocalResourceTable(resourceTable); nil != err {
			log.WithError(err).Errorf("Failed to call StoreLocalResourceTable() to update local jobNode resource on resourceManager.RemoveDiscoveryJobNodeResource(), powerId: {%s}, jobNodeId: {%s}",
				resourceTable.GetPowerId(), jobNodeId)
			return err
		}
	}

	// ##############################
	// B. release resource about jobNode
	// ##############################

	// 1.  unlock local resource table used.
	if err = m.UnLockLocalResourceWithJobNodeId(jobNodeId); nil != err {
		log.WithError(err).Errorf("Failed to unlock local resource with jobNodeId on resourceManager.RemoveDiscoveryJobNodeResource(), jobNodeId: {%s}",
			jobNodeId)
		return err
	}

	//// 2. remove local jobNode reource
	//// remove jobNode local resource
	//if err = m.dataCenter.RemoveLocalResource(jobNodeId); nil != err {
	//	log.WithError(err).Errorf("Failed to remove jobNode local resource on resourceManager.RemoveDiscoveryJobNodeResource(), jobNodeId: {%s}",
	//		jobNodeId)
	//	return err
	//}

	// 3. remove rpc client
	if client, ok := m.QueryJobNodeClient(jobNodeId); ok {
		client.Close()
		m.RemoveJobNodeClient(jobNodeId)
	}
	// 4. update connState of local jobNode info
	if old.GetConnState() != pb.ConnState_ConnState_UnConnected {

		old.ConnState = pb.ConnState_ConnState_UnConnected
		if err := m.dataCenter.SetRegisterNode(pb.PrefixTypeJobNode, old); nil != err {
			log.WithError(err).Errorf("Failed to update jobNode into local db on resourceManager.RemoveDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
				jobNodeId, jobNodeIP, jobNodePort)
			return err
		}

		log.Infof("Succeed update jobNode ConnState to `unconnected` on resourceManager.RemoveDiscoveryJobNodeResource(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
			jobNodeId, jobNodeIP, jobNodePort)
	}

	log.Infof("Succeed remove a old jobNode, jobNodeId: {%s}", jobNodeId)

	return nil
}

func (m *Manager) AddDiscoveryDataNodeResource(identity *apicommonpb.Organization, dataNodeId, dataNodeIP, dataNodePort, dataNodeExternalIP, dataNodeExternalPort string) error {

	log.Infof("Discovered a new dataNode from consul server, add dataNode resource on resourceManager.AddDiscoveryDataNodeResource(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
		dataNodeId, dataNodeIP, dataNodePort)

	client, err := grpclient.NewDataNodeClient(context.Background(), fmt.Sprintf("%s:%s", dataNodeIP, dataNodePort), dataNodeId)
	if nil != err {
		log.WithError(err).Errorf("Failed to connect new dataNode on resourceManager.AddDiscoveryDataNodeResource(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
			dataNodeId, dataNodeIP, dataNodePort)
		return err
	}
	dataNodeStatus, err := client.GetStatus()
	if nil != err {
		log.WithError(err).Errorf("Failed to connect jobNode to query status on resourceManager.AddDiscoveryDataNodeResource(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
			dataNodeId, dataNodeIP, dataNodePort)
		return err
	}
	// 1. add data resource  (disk)
	err = m.dataCenter.StoreDataResourceTable(types.NewDataResourceTable(dataNodeId, dataNodeStatus.GetTotalDisk(), dataNodeStatus.GetUsedDisk(), true))
	if nil != err {
		log.WithError(err).Errorf("Failed to store disk summary of new dataNode on resourceManager.AddDiscoveryDataNodeResource(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
			dataNodeId, dataNodeIP, dataNodePort)
		return err
	}

	// 2. add rpc client
	m.StoreDataNodeClient(dataNodeId, client)

	// 3. add local dataNode info
	// build new dataNode info that was need to store local db
	if err = m.dataCenter.SetRegisterNode(pb.PrefixTypeDataNode,
		&pb.YarnRegisteredPeerDetail{
			Id:           strings.Join([]string{discovery.DataNodeConsulServiceIdPrefix, dataNodeIP, dataNodePort}, discovery.ConsulServiceIdSeparator),
			InternalIp:   dataNodeIP,
			InternalPort: dataNodePort,
			ExternalIp:   dataNodeExternalIP,
			ExternalPort: dataNodeExternalPort,
			ConnState:    pb.ConnState_ConnState_Connected,
		}); nil != err {
		log.WithError(err).Errorf("Failed to store dataNode into local db on resourceManager.AddDiscoveryDataNodeResource(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
			dataNodeId, dataNodeIP, dataNodePort)
		return err
	}

	log.Infof("Succeed add a new dataNode from consul server, add dataNode resource  on resourceManager.AddDiscoveryDataNodeResource(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
		dataNodeId, dataNodeIP, dataNodePort)

	return nil
}

func (m *Manager) UpdateDiscoveryDataNodeResource(identity *apicommonpb.Organization, dataNodeId, dataNodeIP, dataNodePort, dataNodeExternalIP, dataNodeExternalPort string, old *pb.YarnRegisteredPeerDetail) error {
	// check the  via external ip and port comparing old infomation,
	// if it is, update the some things about dataNode.
	if old.GetExternalIp() != dataNodeExternalIP || old.GetExternalPort() != dataNodeExternalPort {

		oldIp := old.GetExternalIp()
		oldPort := old.GetExternalPort()

		// update dataNode info that was need to store local db
		old.ExternalIp = dataNodeExternalIP
		old.ExternalPort = dataNodeExternalPort
		// 1. update local dataNode info
		// update dataNode ip port into local db
		if err := m.dataCenter.SetRegisterNode(pb.PrefixTypeDataNode, old); nil != err {
			log.WithError(err).Errorf("Failed to update dataNode into local db on resourceManager.UpdateDiscoveryDataNodeResource(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
				dataNodeId, dataNodeIP, dataNodePort)
			return err
		}

		log.Infof("Succeed update a old dataNode external ip and port from consul server on resourceManager.UpdateDiscoveryDataNodeResource(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}, old externalIp: {%s}, old externalPort: {%s}, new externalIp: {%s}, new externalPort: {%s}",
			dataNodeId, dataNodeIP, dataNodePort, oldIp, oldPort, old.GetExternalIp(), old.GetExternalPort())
	}

	resourceTable, err := m.dataCenter.QueryDataResourceTable (dataNodeId)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed to query disk summary of old dataNode on resourceManager.UpdateDiscoveryDataNodeResource(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
			dataNodeId, dataNodeIP, dataNodePort)
		return err
	}
	if nil != resourceTable && !resourceTable.GetAlive() {
		resourceTable.SetAlive(true)
		if err := m.dataCenter.StoreDataResourceTable(resourceTable); nil != err {
			log.WithError(err).Errorf("Failed to update alive flag of local dataNode resource on resourceManager.UpdateDiscoveryDataNodeResource(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
				dataNodeId, dataNodeIP, dataNodePort)
			return err
		}
	}

	// check connection status,
	// if it be changed, update the connState value about jobNode
	if old.GetConnState() != pb.ConnState_ConnState_Connected {

		old.ConnState = pb.ConnState_ConnState_Connected
		if err := m.dataCenter.SetRegisterNode(pb.PrefixTypeDataNode, old); nil != err {
			log.WithError(err).Errorf("Failed to update dataNode into local db on resourceManager.UpdateDiscoveryDataNodeResource(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
				dataNodeId, dataNodeIP, dataNodePort)
			return err

		}

		log.Infof("Succeed update dataNode ConnState to `connected` on resourceManager.UpdateDiscoveryDataNodeResource(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
			dataNodeId, dataNodeIP, dataNodePort)
	}

	return nil
}

func (m *Manager) RemoveDiscoveryDataNodeResource(identity *apicommonpb.Organization, dataNodeId, dataNodeIP, dataNodePort, dataNodeExternalIP, dataNodeExternalPort string, old *pb.YarnRegisteredPeerDetail) error {

	log.Infof("Disappeared a old dataNode from consul server, dataNodeId: {%s}", dataNodeId)

	resourceTable, err := m.dataCenter.QueryDataResourceTable(dataNodeId)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed to query disk summary of old dataNode on resourceManager.RemoveDiscoveryDataNodeResource(), dataNodeServiceId: {%s}, dataNodeService: {%s:%s}",
			dataNodeId, dataNodeIP, dataNodePort)
		return err
	}
	if nil != resourceTable && resourceTable.GetAlive() {
		log.Warnf("still have used dataNode information on resourceManager.RemoveDiscoveryDataNodeResource(), %s",
			resourceTable.String())
		// ##############################
		// A. update alive status of resource table about dataNode
		// ##############################

		// 1. update resource table
		resourceTable.SetAlive(false)
		if err := m.dataCenter.StoreDataResourceTable(resourceTable); nil != err {
			log.WithError(err).Errorf("Failed to update alive flag of local dataNode resource on resourceManager.RemoveDiscoveryDataNodeResource(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
				dataNodeId, dataNodeIP, dataNodePort)
			return err
		}
	}

	// ##############################
	// B. release resource about dataNode
	// ##############################

	// 2. remove rpc client
	if client, ok := m.QueryDataNodeClient(dataNodeId); ok {
		client.Close()
		m.RemoveDataNodeClient(dataNodeId)
	}
	// 3. update connState of local dataNode info
	if old.GetConnState() != pb.ConnState_ConnState_UnConnected {

		old.ConnState = pb.ConnState_ConnState_UnConnected
		if err := m.dataCenter.SetRegisterNode(pb.PrefixTypeDataNode, old); nil != err {
			log.WithError(err).Errorf("Failed to update dataNode into local db on resourceManager.RemoveDiscoveryDataNodeResource(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
				dataNodeId, dataNodeIP, dataNodePort)
			return err
		}

		log.Infof("Succeed update dataNode ConnState to `unconnected` on resourceManager.RemoveDiscoveryDataNodeResource(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
			dataNodeId, dataNodeIP, dataNodePort)
	}

	log.Infof("Succeed remove a old dataNode, dataNodeId: {%s}", dataNodeId)

	return nil
}
