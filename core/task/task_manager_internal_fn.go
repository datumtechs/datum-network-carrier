package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	ev "github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/core/schedule"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/lib/fighter/common"
	msgcommonpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/common"
	taskmngpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/taskmng"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	"strconv"
)

func (m *Manager) tryScheduleTask() error {
	nonConsTask, err := m.scheduler.TrySchedule()
	if nil != err && err != schedule.ErrRescheduleLargeThreshold {
		return err
	} else if nil != err && err == schedule.ErrRescheduleLargeThreshold {
		return m.storeFailedReScheduleTask(nonConsTask.GetTask().GetTaskId())
	} else if nil == err && nil == nonConsTask {
		return nil
	}

	go func(nonConsTask *types.NeedConsensusTask) {

		log.Debugf("Start consensus task on 2pc consensus engine, taskId: {%s}", nonConsTask.GetTask().GetTaskId())

		if err := m.consensusEngine.OnPrepare(nonConsTask.GetTask()); nil != err {
			log.Errorf("Failed to call `OnPrepare()` on 2pc consensus engine, taskId: {%s}, err: {%s}", nonConsTask.GetTask().GetTaskId(), err)
			nonConsTask.SendResult(types.NewTaskConsResult(nonConsTask.GetTask().GetTaskId(), types.TaskConsensusInterrupt, err))
			return
		}

		if err := m.consensusEngine.OnHandle(nonConsTask.GetTask(), nonConsTask.GetResultCh()); nil != err {
			log.Errorf("Failed to call `OnHandle()` on 2pc consensus engine, taskId: {%s}, err: {%s}", nonConsTask.GetTask().GetTaskId(), err)
		}

		result := nonConsTask.ReceiveResult()
		log.Debugf("Received task result from consensus, taskId: {%s}, result: {%s}", nonConsTask.GetTask().GetTaskId(), result.String())

		// Consensus failed, task needs to be suspended and rescheduled
		if result.Status != types.TaskConsensusSucceed {

			m.eventEngine.StoreEvent(m.eventEngine.GenerateEvent(ev.TaskFailedConsensus.Type,
				nonConsTask.GetTask().GetTaskId(), nonConsTask.GetTask().GetTaskSender().GetIdentityId(), result.GetErr().Error()))

			// re push task into queue
			if err := m.scheduler.AddTask(nonConsTask.GetTask()); err == schedule.ErrRescheduleLargeThreshold {
				m.storeFailedReScheduleTask(nonConsTask.GetTask().GetTaskId())
			}
		}

	}(nonConsTask)
	return nil
}

func (m *Manager) storeFailedReScheduleTask(taskId string) error {

	task, err := m.resourceMng.GetDB().QueryLocalTask(taskId)
	if nil != err {
		log.Errorf("Failed to query local task for sending datacenter on taskManager.storeFailedReScheduleTask(), taskId: {%s}, err: {%s}", task.GetTaskId(), err)
		return err
	}

	events, err := m.resourceMng.GetDB().QueryTaskEventList(task.GetTaskId())
	if nil != err {
		log.Errorf("Failed to query all task event list for sending datacenter on taskManager.storeFailedReScheduleTask(), taskId: {%s}, err: {%s}", task.GetTaskId(), err)
		return err
	}
	events = append(events, m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
		taskId, task.GetTaskSender().GetIdentityId(), "overflow reschedule threshold"))

	if err := m.storeBadTask(task, events, "overflow reschedule threshold"); nil != err {
		log.Errorf("Failed to sending task to datacenter on taskManager.storeFailedReScheduleTask(), taskId: {%s}, err: {%s}", taskId, err)
		return err
	}
	return nil
}

// To execute task
func (m *Manager) driveTaskForExecute(task *types.NeedExecuteTask) error {

	task.GetTask().GetTaskData().State = apicommonpb.TaskState_TaskState_Running
	task.GetTask().GetTaskData().StartAt = timeutils.UnixMsecUint64()
	if err := m.resourceMng.GetDB().StoreLocalTask(task.GetTask()); nil != err {
		log.Errorf("Failed to update local task state before executing task, taskId: {%s}, need update state: {%s}, err: {%s}",
			task.GetTask().GetTaskId(), apicommonpb.TaskState_TaskState_Running.String(), err)
	}
	// update local cache
	m.addNeedExecuteTaskCache(task)

	switch task.GetLocalTaskRole() {
	case apicommonpb.TaskRole_TaskRole_DataSupplier, apicommonpb.TaskRole_TaskRole_Receiver:
		return m.executeTaskOnDataNode(task)
	case apicommonpb.TaskRole_TaskRole_PowerSupplier:
		return m.executeTaskOnJobNode(task)
	default:
		log.Errorf("Faided to driveTaskForExecute(), Unknown task role, taskId: {%s}, taskRole: {%s}", task.GetTask().GetTaskId(), task.GetLocalTaskRole().String())
		return errors.New("Unknown resource node type")
	}
}

func (m *Manager) executeTaskOnDataNode(task *types.NeedExecuteTask) error {

	// 找到自己的投票
	dataNodeId := task.GetLocalResource().Id

	// clinet *grpclient.DataNodeClient,
	client, has := m.resourceClientSet.QueryDataNodeClient(dataNodeId)
	if !has {
		log.Errorf("Failed to query internal data node on `taskManager.executeTaskOnDataNode()`, taskId: {%s}, dataNodeId: {%s}",
			task.GetTask().GetTaskId(), dataNodeId)
		return errors.New("data node client not found")
	}
	if client.IsNotConnected() {
		if err := client.Reconnect(); nil != err {
			log.Errorf("Failed to connect internal data node on `taskManager.executeTaskOnDataNode()`, taskId: {%s}, dataNodeId: {%s}, err: {%s}",
				task.GetTask().GetTaskId(), dataNodeId, err)
			return err
		}
	}

	req, err := m.makeTaskReadyGoReq(task)
	if nil != err {
		log.Errorf("Falied to make TaskReadyGoReq on `taskManager.executeTaskOnDataNode()`, taskId: {%s}, dataNodeId: {%s}, err: {%s}",
			task.GetTask().GetTaskId(), dataNodeId, err)
		return err
	}

	resp, err := client.HandleTaskReadyGo(req)
	if nil != err {
		log.Errorf("Falied to call publish schedTask to `data-Fighter` node to executing, taskId: {%s}, dataNodeId: {%s}, err: {%s}",
			task.GetTask().GetTaskId(), dataNodeId, err)
		return err
	}
	if !resp.Ok {
		log.Errorf("Falied to executing task from `data-Fighter` node to executing the resp code is not `ok`, taskId: {%s}, dataNodeId: {%s}, resp: {%s}",
			task.GetTask().GetTaskId(), dataNodeId, resp.String())
		return nil
	}

	log.Infof("Success to publish schedTask to `data-Fighter` node to executing,  taskId: {%s}, dataNodeId: {%s}",
		task.GetTask().GetTaskId(), dataNodeId)
	return nil
}

func (m *Manager) executeTaskOnJobNode(task *types.NeedExecuteTask) error {

	// 找到自己的投票
	jobNodeId := task.GetLocalResource().Id

	// clinet *grpclient.DataNodeClient,
	client, has := m.resourceClientSet.QueryJobNodeClient(jobNodeId)
	if !has {
		log.Errorf("Failed to query internal job node on `taskManager.executeTaskOnJobNode()`, taskId: {%s}, jobNodeId: {%s}",
			task.GetTask().GetTaskId(), jobNodeId)
		return errors.New("job node client not found")
	}
	if client.IsNotConnected() {
		if err := client.Reconnect(); nil != err {
			log.Errorf("Failed to connect internal job node on `taskManager.executeTaskOnJobNode()`, taskId: {%s}, jobNodeId: {%s}, err: {%s}",
				task.GetTask().GetTaskId(), jobNodeId, err)
			return err
		}
	}

	req, err := m.makeTaskReadyGoReq(task)
	if nil != err {
		log.Errorf("Falied to make TaskReadyGoReq on `taskManager.executeTaskOnJobNode()`, taskId: {%s}, jobNodeId: {%s}, err: {%s}",
			task.GetTask().GetTaskId(), jobNodeId, err)
		return err
	}

	resp, err := client.HandleTaskReadyGo(req)
	if nil != err {
		log.Errorf("Falied to publish schedTask to `job-Fighter` node to executing, taskId: {%s}, jobNodeId: {%s}, err: {%s}",
			task.GetTask().GetTaskId(), jobNodeId, err)
		return err
	}
	if !resp.Ok {
		log.Errorf("Falied to publish schedTask to `job-Fighter` node to executing the resp code is not `ok`, taskId: {%s}, jobNodeId: {%s}",
			task.GetTask().GetTaskId(), jobNodeId)
		return nil
	}

	log.Infof("Success to publish schedTask to `job-Fighter` node to executing, taskId: {%s}, jobNodeId: {%s}",
		task.GetTask().GetTaskId(), jobNodeId)
	return nil
}

// To terminate task
func (m *Manager) driveTaskForTerminate(task *types.NeedExecuteTask) error {

	m.storeTaskFinalEvent(task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(), fmt.Sprintf("task terminate"), apicommonpb.TaskState_TaskState_Failed)
	switch task.GetLocalTaskRole() {
	case apicommonpb.TaskRole_TaskRole_Sender:

		// send task terminate msg to remote peers
		m.sendTaskTerminateMsgToRemotePeer(task)
		//
		m.publishFinishedTaskToDataCenter(task)
	default:
		//
		m.sendTaskResultMsgToRemotePeer(task)
	}

	switch task.GetLocalTaskRole() {
	case apicommonpb.TaskRole_TaskRole_DataSupplier, apicommonpb.TaskRole_TaskRole_Receiver:
		return m.terminateTaskOnDataNode(task)
	case apicommonpb.TaskRole_TaskRole_PowerSupplier:
		return m.terminateTaskOnJobNode(task)
	default:
		log.Errorf("Faided to driveTaskForTerminate(), Unknown task role, taskId: {%s}, taskRole: {%s}", task.GetTask().GetTaskId(), task.GetLocalTaskRole().String())
		return errors.New("Unknown resource node type")
	}
}

func (m *Manager) terminateTaskOnDataNode(task *types.NeedExecuteTask) error {

	// 找到自己的投票
	dataNodeId := task.GetLocalResource().Id

	// clinet *grpclient.DataNodeClient,
	client, has := m.resourceClientSet.QueryDataNodeClient(dataNodeId)
	if !has {
		log.Errorf("Failed to query internal data node on `taskManager.terminateTaskOnDataNode()`, taskId: {%s}, dataNodeId: {%s}",
			task.GetTask().GetTaskId(), dataNodeId)
		return errors.New("data node client not found")
	}
	if client.IsNotConnected() {
		if err := client.Reconnect(); nil != err {
			log.Errorf("Failed to connect internal data node on `taskManager.terminateTaskOnDataNode()`, taskId: {%s}, dataNodeId: {%s}, err: {%s}",
				task.GetTask().GetTaskId(), dataNodeId, err)
			return err
		}
	}

	req, err := m.makeTerminateTaskReq(task)
	if nil != err {
		log.Errorf("Falied to make TaskCancelReq on `taskManager.terminateTaskOnDataNode()`, taskId: {%s}, dataNodeId: {%s}, err: {%s}",
			task.GetTask().GetTaskId(), dataNodeId, err)
		return err
	}

	resp, err := client.HandleCancelTask(req)
	if nil != err {
		log.Errorf("Falied to call publish schedTask to `data-Fighter` node to terminating, taskId: {%s}, dataNodeId: {%s}, err: {%s}",
			task.GetTask().GetTaskId(), dataNodeId, err)
		return err
	}
	if !resp.Ok {
		log.Errorf("Falied to executing task from `data-Fighter` node to terminating the resp code is not `ok`, taskId: {%s}, dataNodeId: {%s}, resp: {%s}",
			task.GetTask().GetTaskId(), dataNodeId, resp.String())
		return nil
	}

	log.Infof("Success to publish schedTask to `data-Fighter` node to terminating,  taskId: {%s}, dataNodeId: {%s}",
		task.GetTask().GetTaskId(), dataNodeId)
	return nil
}

func (m *Manager) terminateTaskOnJobNode(task *types.NeedExecuteTask) error {

	// 找到自己的投票
	jobNodeId := task.GetLocalResource().Id

	// clinet *grpclient.DataNodeClient,
	client, has := m.resourceClientSet.QueryJobNodeClient(jobNodeId)
	if !has {
		log.Errorf("Failed to query internal job node on `taskManager.terminateTaskOnJobNode()`, taskId: {%s}, jobNodeId: {%s}",
			task.GetTask().GetTaskId(), jobNodeId)
		return errors.New("job node client not found")
	}
	if client.IsNotConnected() {
		if err := client.Reconnect(); nil != err {
			log.Errorf("Failed to connect internal job node on `taskManager.terminateTaskOnJobNode()`, taskId: {%s}, jobNodeId: {%s}, err: {%s}",
				task.GetTask().GetTaskId(), jobNodeId, err)
			return err
		}
	}

	req, err := m.makeTerminateTaskReq(task)
	if nil != err {
		log.Errorf("Falied to make TaskCancelReq on `taskManager.terminateTaskOnJobNode()`, taskId: {%s}, jobNodeId: {%s}, err: {%s}",
			task.GetTask().GetTaskId(), jobNodeId, err)
		return err
	}

	resp, err := client.HandleCancelTask(req)
	if nil != err {
		log.Errorf("Falied to publish schedTask to `job-Fighter` node to terminating, taskId: {%s}, jobNodeId: {%s}, err: {%s}",
			task.GetTask().GetTaskId(), jobNodeId, err)
		return err
	}
	if !resp.Ok {
		log.Errorf("Falied to publish schedTask to `job-Fighter` node to terminating the resp code is not `ok`, taskId: {%s}, jobNodeId: {%s}",
			task.GetTask().GetTaskId(), jobNodeId)
		return nil
	}

	log.Infof("Success to publish schedTask to `job-Fighter` node to terminating, taskId: {%s}, jobNodeId: {%s}",
		task.GetTask().GetTaskId(), jobNodeId)
	return nil
}

func (m *Manager) publishFinishedTaskToDataCenter(task *types.NeedExecuteTask) {

	eventList, err := m.resourceMng.GetDB().QueryTaskEventList(task.GetTask().GetTaskId())
	if nil != err {
		log.Errorf("Failed to Query all task event list for sending datacenter on publishFinishedTaskToDataCenter, taskId: {%s}, err: {%s}", task.GetTask().GetTaskId(), err)
		return
	}
	var isFailed bool
	for _, event := range eventList {
		if event.Type == ev.TaskFailed.Type {
			isFailed = true
			break
		}
	}
	var taskState apicommonpb.TaskState
	if isFailed {
		taskState = apicommonpb.TaskState_TaskState_Failed
	} else {
		taskState = apicommonpb.TaskState_TaskState_Succeed
	}

	log.Debugf("Start publishFinishedTaskToDataCenter, taskId: {%s}, taskState: {%s}", task.GetTask().GetTaskId(), taskState.String())

	finalTask := m.convertScheduleTaskToTask(task.GetTask(), eventList, taskState)
	if err := m.resourceMng.GetDB().InsertTask(finalTask); nil != err {
		log.Errorf("Failed to save task to datacenter on publishFinishedTaskToDataCenter, taskId: {%s}, err: {%s}", task.GetTask().GetTaskId(), err)
		return
	}

	if err := m.resourceMng.GetDB().RemoveLocalTaskExecuteStatus(task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetPartyId()); nil != err {
		log.Errorf("Failed to remove task executing status on publishFinishedTaskToDataCenter, taskId: {%s}, partyId: {%s}, err: {%s}",
			task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), err)
		return
	}
	// clean local task cache
	m.resourceMng.ReleaseLocalResourceWithTask("on taskManager.publishFinishedTaskToDataCenter()",
		task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), resource.SetAllReleaseResourceOption())

	log.Debugf("Finished pulishFinishedTaskToDataCenter, taskId: {%s}, partyId: {%s}, taskState: {%s}", task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), taskState)
}
func (m *Manager) sendTaskResultMsgToRemotePeer(task *types.NeedExecuteTask) {

	log.Debugf("Start sendTaskResultMsgToRemotePeer, taskId: {%s}, taskRole: {%s},  partyId: {%s}, remote pid: {%s}",
		task.GetTask().GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())

	if task.HasRemotePID() {

		//if err := handler.SendTaskResultMsg(context.TODO(), m.p2p, task.GetRemotePID(), m.makeTaskResultByEventList(task)); nil != err {
		if err := m.p2p.Broadcast(context.TODO(), m.makeTaskResultByEventList(task)); nil != err {
			log.Errorf("failed to call `SendTaskResultMsg`, taskId: {%s}, taskRole: {%s},  partyId: {%s}, remote pid: {%s}, err: {%s}",
				task.GetTask().GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID(), err)
			return
		}
	}

	if err := m.resourceMng.GetDB().RemoveLocalTaskExecuteStatus(task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetPartyId()); nil != err {
		log.Errorf("Failed to remove task executing status on sendTaskResultMsgToRemotePeer, taskId: {%s}, taskRole: {%s},  partyId: {%s}, remote pid: {%s}, err: {%s}",
			task.GetTask().GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID(), err)
		return
	}

	// clean local task cache
	m.resourceMng.ReleaseLocalResourceWithTask("on taskManager.sendTaskResultMsgToRemotePeer()",
		task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), resource.SetAllReleaseResourceOption())

	log.Debugf("Finished sendTaskResultMsgToRemotePeer,  taskId: {%s}, taskRole: {%s},  partyId: {%s}, remote pid: {%s}",
		task.GetTask().GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())
}

func (m *Manager) sendTaskResourceUsageMsgToRemotePeer(task *types.NeedExecuteTask, usage *types.TaskResuorceUsage) {

	msg := &taskmngpb.TaskResourceUsageMsg{
		MsgOption: &msgcommonpb.MsgOption{
			ProposalId:      task.GetProposalId().Bytes(),
			SenderRole:      uint64(task.GetLocalTaskRole()),
			SenderPartyId:   []byte(task.GetLocalTaskOrganization().GetPartyId()),
			ReceiverRole:    uint64(task.GetRemoteTaskRole()),
			ReceiverPartyId: []byte(task.GetRemoteTaskOrganization().GetPartyId()),
			MsgOwner: &msgcommonpb.TaskOrganizationIdentityInfo{
				Name:       []byte(task.GetLocalTaskOrganization().GetNodeName()),
				NodeId:     []byte(task.GetLocalTaskOrganization().GetNodeId()),
				IdentityId: []byte(task.GetLocalTaskOrganization().GetIdentityId()),
				PartyId:    []byte(task.GetLocalTaskOrganization().GetPartyId()),
			},
		},
		TaskId: []byte(task.GetTask().GetTaskId()),
		Usage: &msgcommonpb.ResourceUsage{
			TotalMem: usage.GetTotalMem(),
			UsedMem: usage.GetUsedMem(),
			TotalProcessor: uint64(usage.GetTotalProcessor()),
			UsedProcessor: uint64(usage.GetUsedProcessor()),
			TotalBandwidth: usage.GetTotalBandwidth(),
			UsedBandwidth: usage.GetUsedBandwidth(),
			TotalDisk: usage.GetTotalDisk(),
			UsedDisk: usage.GetUsedDisk(),
		},
		CreateAt: timeutils.UnixMsecUint64(),
		Sign: nil,
	}

	// send msg to remote target peer with broadcast
	if task.HasRemotePID() {
		//if err := handler.SendTaskResourceUsageMsg(context.TODO(), m.p2p, task.GetRemotePID(), msg); nil != err {
		if err := m.p2p.Broadcast(context.TODO(), msg); nil != err {
			log.Errorf("failed to call `SendTaskResourceUsageMsg`, taskId: {%s}, taskRole: {%s},  partyId: {%s}, remote pid: {%s}, err: {%s}",
				task.GetTask().GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID(), err)
			return
		}
	} else {

		// send msg to current peer
		if err := m.OnTaskResourceUsageMsg(task.GetRemotePID(), msg); nil != err {
			log.Errorf("failed to call `OnTaskResourceUsageMsg`, taskId: {%s}, taskRole: {%s},  partyId: {%s}, remote pid: {%s}, err: {%s}",
				task.GetTask().GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID(), err)
			return
		}
	}
}

func (m *Manager) sendTaskTerminateMsgToRemotePeer (task *types.NeedExecuteTask) {

	msg := &taskmngpb.TaskTerminateMsg{
		MsgOption: &msgcommonpb.MsgOption{
			ProposalId:      task.GetProposalId().Bytes(),
			SenderRole:      uint64(task.GetLocalTaskRole()),
			SenderPartyId:   []byte(task.GetLocalTaskOrganization().GetPartyId()),
			ReceiverRole:    uint64(task.GetRemoteTaskRole()),
			ReceiverPartyId: []byte(task.GetRemoteTaskOrganization().GetPartyId()),
			MsgOwner: &msgcommonpb.TaskOrganizationIdentityInfo{
				Name:       []byte(task.GetLocalTaskOrganization().GetNodeName()),
				NodeId:     []byte(task.GetLocalTaskOrganization().GetNodeId()),
				IdentityId: []byte(task.GetLocalTaskOrganization().GetIdentityId()),
				PartyId:    []byte(task.GetLocalTaskOrganization().GetPartyId()),
			},
		},
		TaskId: []byte(task.GetTask().GetTaskId()),
		CreateAt: timeutils.UnixMsecUint64(),
		Sign: nil,
	}

	// send msg to remote target peer with broadcast
	if task.HasRemotePID() {
		//if err := handler.SendTaskTerminateMsg(context.TODO(), m.p2p, task.GetRemotePID(), msg); nil != err {
		if err := m.p2p.Broadcast(context.TODO(), msg); nil != err {
			log.Errorf("failed to call `SendTaskTerminateMsg`, taskId: {%s}, taskRole: {%s},  partyId: {%s}, remote pid: {%s}, err: {%s}",
				task.GetTask().GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID(), err)
			return
		}
	} else {

		// send msg to current peer
		if err := m.OnTaskTerminateMsg(task.GetRemotePID(), msg); nil != err {
			log.Errorf("failed to call `OnTaskTerminateMsg`, taskId: {%s}, taskRole: {%s},  partyId: {%s}, remote pid: {%s}, err: {%s}",
				task.GetTask().GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID(), err)
			return
		}
	}
}

func (m *Manager) sendLocalTaskToScheduler(tasks types.TaskDataArray) {
	m.localTasksCh <- tasks
}
func (m *Manager) sendTaskEvent(reportEvent *types.ReportTaskEvent) {
	m.eventCh <- reportEvent
}

func (m *Manager) storeBadTask(task *types.Task, events []*libtypes.TaskEvent, reason string) error {
	task.GetTaskData().TaskEvents = events
	task.GetTaskData().EventCount = uint32(len(events))
	task.GetTaskData().State = apicommonpb.TaskState_TaskState_Failed
	task.GetTaskData().Reason = reason
	task.GetTaskData().EndAt = timeutils.UnixMsecUint64()

	m.resourceMng.GetDB().RemoveLocalTask(task.GetTaskId())
	m.resourceMng.GetDB().RemoveTaskPowerPartyIds(task.GetTaskId())
	m.resourceMng.GetDB().RemoveTaskEventList(task.GetTaskId())

	return m.resourceMng.GetDB().InsertTask(task)
}

// TODO MOCK TASK
func (m *Manager) storeMockTask(task *types.Task, events []*libtypes.TaskEvent, reason string) error {
	task.GetTaskData().TaskEvents = events
	task.GetTaskData().EventCount = uint32(len(events))
	task.GetTaskData().State = apicommonpb.TaskState_TaskState_Succeed
	task.GetTaskData().Reason = reason
	task.GetTaskData().StartAt = timeutils.UnixMsecUint64()
	task.GetTaskData().EndAt = timeutils.UnixMsecUint64()

	m.resourceMng.GetDB().RemoveLocalTask(task.GetTaskId())
	m.resourceMng.GetDB().RemoveTaskPowerPartyIds(task.GetTaskId())
	m.resourceMng.GetDB().RemoveTaskEventList(task.GetTaskId())

	return m.resourceMng.GetDB().InsertTask(task)
}

func (m *Manager) convertScheduleTaskToTask(task *types.Task, eventList []*libtypes.TaskEvent, state apicommonpb.TaskState) *types.Task {
	task.GetTaskData().TaskEvents = eventList
	task.GetTaskData().EventCount = uint32(len(eventList))
	task.GetTaskData().EndAt = timeutils.UnixMsecUint64()
	task.GetTaskData().State = state
	return task
}

func (m *Manager) makeTaskReadyGoReq(task *types.NeedExecuteTask) (*common.TaskReadyGoReq, error) {

	ownerPort := string(task.GetResources().GetOwnerPeerInfo().GetPort())
	port, err := strconv.Atoi(ownerPort)
	if nil != err {
		return nil, err
	}

	var dataPartyArr []string
	var powerPartyArr []string
	var receiverPartyArr []string

	peerList := []*common.TaskReadyGoReq_Peer{
		&common.TaskReadyGoReq_Peer{
			Ip:      string(task.GetResources().GetOwnerPeerInfo().GetIp()),
			Port:    int32(port),
			PartyId: string(task.GetResources().GetOwnerPeerInfo().GetPartyId()),
		},
	}
	dataPartyArr = append(dataPartyArr, string(task.GetResources().GetOwnerPeerInfo().GetPartyId()))

	for _, dataSupplier := range task.GetResources().GetDataSupplierPeerInfoList() {
		portStr := string(dataSupplier.GetPort())
		port, err := strconv.Atoi(portStr)
		if nil != err {
			return nil, err
		}
		peerList = append(peerList, &common.TaskReadyGoReq_Peer{
			Ip:      string(dataSupplier.GetIp()),
			Port:    int32(port),
			PartyId: string(dataSupplier.GetPartyId()),
		})
		dataPartyArr = append(dataPartyArr, string(dataSupplier.GetPartyId()))
	}

	for _, powerSupplier := range task.GetResources().GetPowerSupplierPeerInfoList() {
		portStr := string(powerSupplier.GetPort())
		port, err := strconv.Atoi(portStr)
		if nil != err {
			return nil, err
		}
		peerList = append(peerList, &common.TaskReadyGoReq_Peer{
			Ip:      string(powerSupplier.GetIp()),
			Port:    int32(port),
			PartyId: string(powerSupplier.GetPartyId()),
		})

		powerPartyArr = append(powerPartyArr, string(powerSupplier.GetPartyId()))
	}

	for _, receiver := range task.GetResources().GetResultReceiverPeerInfoList() {
		portStr := string(receiver.GetPort())
		port, err := strconv.Atoi(portStr)
		if nil != err {
			return nil, err
		}
		peerList = append(peerList, &common.TaskReadyGoReq_Peer{
			Ip:      string(receiver.GetIp()),
			Port:    int32(port),
			PartyId: string(receiver.GetPartyId()),
		})

		receiverPartyArr = append(receiverPartyArr, string(receiver.GetPartyId()))
	}

	contractExtraParams, err := m.makeContractParams(task)
	if nil != err {
		return nil, err
	}
	log.Debugf("Succeed make contractCfg, taskId:{%s}, contractCfg: %s", task.GetTask().GetTaskId(), contractExtraParams)
	return &common.TaskReadyGoReq{
		TaskId:     task.GetTask().GetTaskId(),
		ContractId: task.GetTask().GetTaskData().GetCalculateContractCode(),
		//DataId: "",
		PartyId: task.GetLocalTaskOrganization().GetPartyId(),
		//EnvId: "",
		Peers:            peerList,
		ContractCfg:      contractExtraParams,
		DataParty:        dataPartyArr,
		ComputationParty: powerPartyArr,
		ResultParty:      receiverPartyArr,
	}, nil
}

func (m *Manager) makeContractParams(task *types.NeedExecuteTask) (string, error) {

	partyId := task.GetLocalTaskOrganization().GetPartyId()

	var filePath string
	var keyColumn string
	var selectedColumns []string

	if task.GetLocalTaskRole() == apicommonpb.TaskRole_TaskRole_DataSupplier {

		var find bool

		for _, dataSupplier := range task.GetTask().GetTaskData().GetDataSuppliers() {
			if partyId == dataSupplier.GetOrganization().GetPartyId() {

				userType := task.GetTask().GetTaskData().GetUserType()
				user := task.GetTask().GetTaskData().GetUser()
				metadataId := dataSupplier.GetMetadataId()

				// verify metadataAuth first
				if !m.authMng.VerifyMetadataAuth(userType, user, metadataId) {
					return "", fmt.Errorf("verify user metadataAuth failed, userType: {%s}, user: {%s}, metadataId: {%s}",
						userType, user, metadataId)
				}

				metadata, err := m.resourceMng.GetDB().QueryMetadataByDataId(metadataId)
				if nil != err {
					return "", err
				}
				filePath = metadata.GetData().GetFilePath()

				keyColumn = dataSupplier.GetKeyColumn().GetCName()
				selectedColumns = make([]string, len(dataSupplier.GetSelectedColumns()))
				for i, col := range dataSupplier.GetSelectedColumns() {
					selectedColumns[i] = col.GetCName()
				}

				// query metadataAuthId by metadataId
				metadataAuthId, err := m.authMng.QueryMetadataAuthIdByMetadataId(userType, user, metadataId)
				if nil != err {
					return "", err
				}
				// ConsumeMetadataAuthority
				if err = m.authMng.ConsumeMetadataAuthority(metadataAuthId); nil != err {
					return "", err
				}

				find = true
				break
			}
		}

		if !find {
			return "", fmt.Errorf("can not find the dataSupplier for find originFilePath, taskId: {%s}, self.IdentityId: {%s}, seld.PartyId: {%s}",
				task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(), task.GetLocalTaskOrganization().GetPartyId())
		}
	}

	req := &types.FighterTaskReadyGoReqContractCfg{
		PartyId: partyId,
		DataParty: struct {
			InputFile       string   `json:"input_file"`
			KeyColumn       string   `json:"key_column"`
			SelectedColumns []string `json:"selected_columns"`
		}{
			InputFile:       filePath,
			KeyColumn:       keyColumn, // only dataSupplier own, but power supplier never own
			SelectedColumns: selectedColumns,
		},
	}

	var dynamicParameter map[string]interface{}
	log.Debugf("Start json Unmarshal the `ContractExtraParams`, taskId: {%s}, ContractExtraParams: %s", task.GetTask().GetTaskId(), task.GetTask().GetTaskData().GetContractExtraParams())
	if "" != task.GetTask().GetTaskData().GetContractExtraParams() {
		if err := json.Unmarshal([]byte(task.GetTask().GetTaskData().GetContractExtraParams()), &dynamicParameter); nil != err {
			return "", fmt.Errorf("can not json Unmarshal the `ContractExtraParams` of task, taskId: {%s}, self.IdentityId: {%s}, seld.PartyId: {%s}, err: {%s}",
				task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(), task.GetLocalTaskOrganization().GetPartyId(), err)
		}
	}
	req.DynamicParameter = dynamicParameter

	b, err := json.Marshal(req)
	if nil != err {
		return "", fmt.Errorf("can not json Marshal the `FighterTaskReadyGoReqContractCfg`, taskId: {%s}, self.IdentityId: {%s}, seld.PartyId: {%s}, err: {%s}",
			task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(), task.GetLocalTaskOrganization().GetPartyId(), err)
	}
	return string(b), nil
}

// make terminate rpc req
func (m *Manager) makeTerminateTaskReq(task *types.NeedExecuteTask) (*common.TaskCancelReq, error) {
	return &common.TaskCancelReq{
		TaskId: task.GetTask().GetTaskId(),
	}, nil
}

func (m *Manager) addNeedExecuteTaskCache(task *types.NeedExecuteTask) {
	m.runningTaskCacheLock.Lock()
	cache, ok := m.runningTaskCache[task.GetTask().GetTaskId()]
	if !ok {
		cache = make(map[string]*types.NeedExecuteTask, 0)
	}
	cache[task.GetLocalTaskOrganization().GetPartyId()] = task
	m.runningTaskCache[task.GetTask().GetTaskId()] = cache
	m.runningTaskCacheLock.Unlock()
}

func (m *Manager) removeNeedExecuteTask(taskId string) {
	m.runningTaskCacheLock.Lock()
	delete(m.runningTaskCache, taskId)
	m.runningTaskCacheLock.Unlock()
}

func (m *Manager) removeNeedExecuteTaskCache(taskId, partyId string) {
	m.runningTaskCacheLock.Lock()
	cache, ok := m.runningTaskCache[taskId]
	if !ok {
		return
	}
	delete(cache, partyId)
	if len(cache) == 0 {
		delete(m.runningTaskCache, taskId)
	} else {
		m.runningTaskCache[taskId] = cache
	}
	m.runningTaskCacheLock.Unlock()
}

func (m *Manager) hasNeedExecuteTaskCache(taskId, partyId string) bool {
	m.runningTaskCacheLock.RLock()
	defer m.runningTaskCacheLock.RUnlock()
	cache, ok := m.runningTaskCache[taskId]
	if !ok {
		return false
	}
	_, ok = cache[partyId]
	return ok
}

func (m *Manager) queryNeedExecuteTaskCache(taskId, partyId string) (*types.NeedExecuteTask, bool) {
	m.runningTaskCacheLock.RLock()
	defer m.runningTaskCacheLock.RUnlock()
	cache, ok := m.runningTaskCache[taskId]
	if !ok {
		return nil, false
	}
	task, ok := cache[partyId]
	return task, ok
}

func (m *Manager) mustQueryNeedExecuteTaskCache(taskId, partyId string) *types.NeedExecuteTask {
	task, _ := m.queryNeedExecuteTaskCache(taskId, partyId)
	return task
}

func (m *Manager) ForEachRunningTaskCache(f func(taskId string, task *types.NeedExecuteTask) bool) {
	m.runningTaskCacheLock.Lock()
	for taskId, cache := range m.runningTaskCache {
		for _, task := range cache {
			if ok := f(taskId, task); ok {
			}
		}
	}
	m.runningTaskCacheLock.Unlock()
}

func (m *Manager) makeTaskResultByEventList(task *types.NeedExecuteTask) *taskmngpb.TaskResultMsg {

	if task.GetLocalTaskRole() == apicommonpb.TaskRole_TaskRole_Sender {
		log.Errorf("send task OR task owner can not make TaskResult Msg")
		return nil
	}

	eventList, err := m.resourceMng.GetDB().QueryTaskEventList(task.GetTask().GetTaskId())
	if nil != err {
		log.Errorf("Failed to make TaskResultMsg with query task eventList, taskId {%s}, err {%s}", task.GetTask().GetTaskId(), err)
		return nil
	}
	return &taskmngpb.TaskResultMsg{
		MsgOption: &msgcommonpb.MsgOption{
			ProposalId:      task.GetProposalId().Bytes(),
			SenderRole:      uint64(task.GetLocalTaskRole()),
			SenderPartyId:   []byte(task.GetLocalTaskOrganization().GetPartyId()),
			ReceiverRole:    uint64(task.GetRemoteTaskRole()),
			ReceiverPartyId: []byte(task.GetRemoteTaskOrganization().GetPartyId()),
			MsgOwner: &msgcommonpb.TaskOrganizationIdentityInfo{
				Name:       []byte(task.GetLocalTaskOrganization().GetNodeName()),
				NodeId:     []byte(task.GetLocalTaskOrganization().GetNodeId()),
				IdentityId: []byte(task.GetLocalTaskOrganization().GetIdentityId()),
				PartyId:    []byte(task.GetLocalTaskOrganization().GetPartyId()),
			},
		},
		TaskEventList: types.ConvertTaskEventArr(eventList),
		CreateAt:      timeutils.UnixMsecUint64(),
		Sign:          nil,
	}
}

func (m *Manager) handleTaskEvent(partyId string, event *libtypes.TaskEvent) error {
	eventType := event.Type
	if len(eventType) != ev.EventTypeCharLen {
		return ev.IncEventType
	}
	// TODO need to validate the task that have been processing ? Maybe~
	if event.Type == ev.TaskExecuteSucceedEOF.Type || event.Type == ev.TaskExecuteFailedEOF.Type {
		if task, ok := m.queryNeedExecuteTaskCache(event.GetTaskId(), partyId); ok {

			log.Debugf("Start handleTaskEvent, `event is the end`, current partyId: {%s}, event: %s", partyId, event.String())

			// store EOF event first
			m.resourceMng.GetDB().StoreTaskEvent(event)
			if event.Type == ev.TaskExecuteFailedEOF.Type {
				m.storeTaskFinalEvent(task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(), "task execute failed eof", apicommonpb.TaskState_TaskState_Failed)
			} else {
				m.storeTaskFinalEvent(task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(), "task execute succeed eof", apicommonpb.TaskState_TaskState_Succeed)
			}

			if task.GetLocalTaskRole() == apicommonpb.TaskRole_TaskRole_Sender {
				// announce remote peer to terminate this task
				m.sendTaskTerminateMsgToRemotePeer(task)
				// handle this task result with current peer
				m.publishFinishedTaskToDataCenter(task)
				m.removeNeedExecuteTaskCache(event.GetTaskId(), partyId)
			} else {
				// send this task result to remote target peer
				m.sendTaskResultMsgToRemotePeer(task)
				m.removeNeedExecuteTaskCache(event.GetTaskId(), partyId)
			}
			return nil
		} else {
			return errors.New(fmt.Sprintf("Not found task cache, taskId: {%s}", event.GetTaskId()))
		}

	} else {

		log.Debugf("Start handleTaskEvent, `event is not the end`, event: %s", event.String())
		// It's not EOF event, then the task still executing, so store this event
		return m.resourceMng.GetDB().StoreTaskEvent(event)
	}
}

func (m *Manager) handleNeedExecuteTask(task *types.NeedExecuteTask) {

	log.Debugf("Start handle needExecuteTask, remote pid: {%s} proposalId: {%s}, taskId: {%s}, self taskRole: {%s}, self taskOrganization: {%s}",
		task.GetRemotePID(), task.GetProposalId(), task.GetTask().GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().String())

	// Store task exec status
	if err := m.resourceMng.GetDB().StoreLocalTaskExecuteStatus(task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetPartyId()); nil != err {
		log.Errorf("Failed to store local task about exec status,  remote pid: {%s} proposalId: {%s}, taskId: {%s}, self taskRole: {%s}, self taskOrganization: {%s}, err: {%s}",
			task.GetRemotePID(), task.GetProposalId(), task.GetTask().GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().String(), err)
		return
	}

	switch task.GetLocalTaskRole() {
	case apicommonpb.TaskRole_TaskRole_Sender:
		if err := m.driveTaskForExecute(task); nil != err {
			log.Errorf("Failed to execute task on taskOnwer node, taskId:{%s}, %s", task.GetTask().GetTaskId(), err)

			m.storeTaskFinalEvent(task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(), fmt.Sprintf("failed to execute task"), apicommonpb.TaskState_TaskState_Failed)
			m.publishFinishedTaskToDataCenter(task)
			m.removeNeedExecuteTaskCache(task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())

		}

	default:
		if err := m.driveTaskForExecute(task); nil != err {
			log.Errorf("Failed to execute task on %s node, taskId: {%s}, %s", task.GetLocalTaskRole().String(), task.GetTask().GetTaskId(), err)

			// 因为是 task 参与者, 所以需要构造 taskResult 发送给 task 发起者.. (里面有解锁 本地资源 ...)
			m.storeTaskFinalEvent(task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(), fmt.Sprintf("failed to execute task"), apicommonpb.TaskState_TaskState_Failed)
			m.sendTaskResultMsgToRemotePeer(task)
			m.removeNeedExecuteTaskCache(task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())
		}
	}
}

func (m *Manager) expireTaskMonitor() {

	m.runningTaskCacheLock.Lock()

	for taskId, cache := range m.runningTaskCache {

		for partyId, task := range cache {

			if task.GetTask().GetTaskData().State == apicommonpb.TaskState_TaskState_Running && task.GetTask().GetTaskData().GetStartAt() != 0 {

				// the task has running expire
				var duration uint64

				switch task.GetLocalTaskRole() {
				case apicommonpb.TaskRole_TaskRole_Sender:
					duration = timeutils.UnixMsecUint64() - task.GetTask().GetTaskData().GetStartAt() + uint64(senderExecuteTaskExpire.Milliseconds())
				default:
					duration = timeutils.UnixMsecUint64() - task.GetTask().GetTaskData().GetStartAt()
				}

				if duration >= task.GetTask().GetTaskData().GetOperationCost().GetDuration() {
					log.Infof("Has task running expire, taskId: {%s}, current running duration: {%d ms}, need running duration: {%d ms}",
						taskId, duration, task.GetTask().GetTaskData().GetOperationCost().GetDuration())

					m.storeTaskFinalEvent(task.GetTask().GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(), fmt.Sprintf("task running expire"), apicommonpb.TaskState_TaskState_Failed)
					switch task.GetLocalTaskRole() {
					case apicommonpb.TaskRole_TaskRole_Sender:
						m.publishFinishedTaskToDataCenter(task)

					default:
						m.sendTaskResultMsgToRemotePeer(task)
					}

					// clean current party task cache
					delete(cache, partyId)
					if len(cache) == 0 {
						delete(m.runningTaskCache, taskId)
					} else {
						m.runningTaskCache[taskId] = cache
					}

				}
			}
		}
	}

	m.runningTaskCacheLock.Unlock()
}

func (m *Manager) storeTaskFinalEvent(taskId, identityId, extra string, state apicommonpb.TaskState) {
	var evTyp string
	var evMsg string
	if state == apicommonpb.TaskState_TaskState_Failed {
		evTyp = ev.TaskFailed.Type
		evMsg = ev.TaskFailed.Msg
	} else {
		evTyp = ev.TaskSucceed.Type
		evMsg = ev.TaskSucceed.Msg
	}
	if "" != extra {
		evMsg = extra
	}
	m.resourceMng.GetDB().StoreTaskEvent(m.eventEngine.GenerateEvent(evTyp, taskId, identityId, evMsg))
}

func (m *Manager) storeMetaUsedTaskId (task *types.Task) error {
	identityId, err := m.resourceMng.GetDB().QueryIdentityId()
	if nil != err {
		return err
	}
	for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
		if dataSupplier.GetOrganization().GetIdentityId() == identityId {
			if err := m.resourceMng.GetDB().StoreMetadataUsedTaskId(dataSupplier.GetMetadataId(), task.GetTaskId()); nil != err {
				return err
			}
		}
	}
	return nil
}

func (m *Manager) ValidateTaskResultMsg(pid peer.ID, taskResultMsg *taskmngpb.TaskResultMsg) error {
	msg := types.FetchTaskResultMsg(taskResultMsg) // fetchTaskResultMsg(taskResultMsg)

	log.Debugf("Received remote taskResultMsg on ValidateTaskResultMsg(), remote pid: {%s}, taskResultMsg: %s", pid, msg.String())

	if len(msg.TaskEventList) == 0 {
		return nil
	}

	taskId := msg.TaskEventList[0].GetTaskId()

	for _, event := range msg.TaskEventList {
		if taskId != event.GetTaskId() {
			return fmt.Errorf("Received event failed, has invalid taskId: {%s}, right taskId: {%s}", event.GetTaskId(), taskId)
		}
	}

	return nil
}

func (m *Manager) OnTaskResultMsg(pid peer.ID, taskResultMsg *taskmngpb.TaskResultMsg) error {

	msg := types.FetchTaskResultMsg(taskResultMsg)

	log.Debugf("Received remote taskResultMsg, remote pid: {%s}, taskResultMsg: %s", pid, msg.String())

	if len(msg.TaskEventList) == 0 {
		return nil
	}

	taskId := msg.TaskEventList[0].GetTaskId()

	has, err := m.resourceMng.GetDB().HasLocalTaskExecute(taskId, msg.MsgOption.ReceiverPartyId)
	if nil != err {
		log.Errorf("Failed to query local task executing status on `onTaskResultMsg`, taskId: {%s}, err: {%s}", taskId, err)
		return fmt.Errorf("query local task executing status failed")
	}

	if !has {
		log.Warnf("Warning not found local task executing status on `onTaskResultMsg`, taskId: {%s}", taskId)
		//return fmt.Errorf("%s, the local task executing status is not found", ctypes.ErrTaskResultMsgInvalid)
		return nil
	}
	for _, event := range msg.TaskEventList {
		if err := m.resourceMng.GetDB().StoreTaskEvent(event); nil != err {
			log.Errorf("Failed to store local task event from remote peer, remote peerId: {%s}, event: {%s}", pid, event.String())
		}
	}

	return nil
}

func (m *Manager) ValidateTaskResourceUsageMsg(pid peer.ID, taskResourceUsageMsg *taskmngpb.TaskResourceUsageMsg) error {
	return nil
}

func (m *Manager) OnTaskResourceUsageMsg(pid peer.ID, usageMsg *taskmngpb.TaskResourceUsageMsg) error {
	msg := types.FetchTaskResourceUsageMsg(usageMsg)

	log.Debugf("Received remote taskResourceUsageMsg, remote pid: {%s}, taskResultMsg: %s", pid, msg.String())

	has, err := m.resourceMng.GetDB().HasLocalTaskExecute(msg.GetUsage().GetTaskId(), msg.MsgOption.ReceiverPartyId)
	if nil != err {
		log.Errorf("Failed to query local task executing status on `OnTaskResourceUsageMsg`, taskId: {%s}, remote partyId: {%s}, err: {%s}",
			msg.GetUsage().GetTaskId(), msg.GetUsage().GetPartyId(), err)
		return fmt.Errorf("query local task executing status failed")
	}

	if !has {
		log.Warnf("Warning not found local task executing status on `OnTaskResourceUsageMsg`, taskId: {%s}, remote partyId: {%s}",
			msg.GetUsage().GetTaskId(), msg.GetUsage().GetPartyId())
		//return fmt.Errorf("%s, the local task executing status is not found", ctypes.ErrTaskResultMsgInvalid)
		return nil
	}

	// todo 还不清楚需不需要更新本地的  ResourceUsage
	//if err := m.resourceMng.GetDB().StoreTaskResuorceUsage(msg.GetUsage()); nil != err {
	//	log.Errorf("Failed to store task resource usage on `OnTaskResourceUsageMsg`, taskId: {%s}, remote partyId: {%s}, err: {%s}",
	//		msg.GetUsage().GetTaskId(), msg.GetUsage().GetPartyId(), err)
	//	return fmt.Errorf("%s, the local task executing status is not found", ctypes.ErrTaskResourceUsageMsgInvalid)
	//}

	// Update task resourceUUsed of powerSuppliers of local task
	task, err := m.resourceMng.GetDB().QueryLocalTask(msg.GetUsage().GetTaskId())
	if nil != err {
		log.Errorf("Failed to query local task info on `OnTaskResourceUsageMsg`, taskId: {%s}, remote partyId: {%s}, err: {%s}",
			msg.GetUsage().GetTaskId(), msg.GetUsage().GetPartyId(), err)
		return fmt.Errorf("query local task failed")
	}
	for i, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {

		if msg.GetMsgOption().SenderPartyId == powerSupplier.GetOrganization().GetPartyId() &&
			msg.GetMsgOption().Owner.GetIdentityId() == powerSupplier.GetOrganization().GetIdentityId() {
			task.GetTaskData().PowerSuppliers[i].ResourceUsedOverview.UsedMem = msg.GetUsage().GetUsedMem()
			task.GetTaskData().PowerSuppliers[i].ResourceUsedOverview.UsedProcessor = msg.GetUsage().GetUsedProcessor()
			task.GetTaskData().PowerSuppliers[i].ResourceUsedOverview.UsedBandwidth = msg.GetUsage().GetUsedBandwidth()
			task.GetTaskData().PowerSuppliers[i].ResourceUsedOverview.UsedDisk = msg.GetUsage().GetUsedDisk()
		}
	}
	err = m.resourceMng.GetDB().StoreLocalTask(task)
	if nil != err {
		log.Errorf("Failed to store local task info on `OnTaskResourceUsageMsg`, taskId: {%s}, remote partyId: {%s}, err: {%s}",
			msg.GetUsage().GetTaskId(), msg.GetUsage().GetPartyId(), err)
		return fmt.Errorf("store local task failed")
	}

	return nil
}

func (m *Manager) ValidateTaskTerminateMsg(pid peer.ID, terminateMsg *taskmngpb.TaskTerminateMsg) error { return nil }

func (m *Manager) OnTaskTerminateMsg (pid peer.ID, terminateMsg *taskmngpb.TaskTerminateMsg) error {
	msg := types.FetchTaskTerminateTaskMngMsg(terminateMsg)
	log.Debugf("Received remote taskTerminateMsg, remote pid: {%s}, taskTerminateMsg: %s", pid, msg.String())


	localTask, err := m.resourceMng.GetDB().QueryLocalTask(msg.GetTaskId())
	if nil != err {
		log.Errorf("Failed to query local task on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, partyId: {%s}, err: {%s}",
			msg.GetTaskId(), msg.GetMsgOption().ReceiverPartyId, err)
		return err
	}

	if nil == localTask {
		log.Errorf("Not found local task on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, partyId: {%s}",
			msg.GetTaskId(), msg.GetMsgOption().ReceiverPartyId)
		return err
	}

	needExecuteTask, ok := m.queryNeedExecuteTaskCache(msg.GetTaskId(), msg.GetMsgOption().ReceiverPartyId)
	if ok {
		has, err := m.resourceMng.GetDB().HasLocalTaskExecute(msg.GetTaskId(), msg.GetMsgOption().ReceiverPartyId)
		if nil != err {
			log.Errorf("Failed to query local task execute status on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, partyId: {%s}, err: {%s}",
				msg.GetTaskId(), msg.GetMsgOption().ReceiverPartyId, err)
			return err
		}
		if has {
			// call terminate task
			if err := m.driveTaskForTerminate(needExecuteTask); nil != err {
				log.Errorf("Failed to call driveTaskForTerminate() on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, partyId: {%s}, err: {%s}",
					msg.GetTaskId(), msg.GetMsgOption().ReceiverPartyId, err)
				return err
			}

		} else {
			// remove the task on scheduler (maybe task on consensus now)
			m.sendTaskResultMsgToRemotePeer(needExecuteTask)
		}
	} else {
		// remove local task short circuit
		m.resourceMng.ReleaseLocalResourceWithTaskShortCircuit("on `taskManager.OnTaskTerminateMsg()`",
			msg.GetTaskId(), msg.GetMsgOption().ReceiverPartyId, resource.SetAllReleaseResourceOption())
	}
	return nil
}