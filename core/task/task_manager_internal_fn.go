package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/runutil"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	ev "github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/core/schedule"
	"github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	fightercommon "github.com/RosettaFlow/Carrier-Go/lib/fighter/common"
	msgcommonpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/common"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	taskmngpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/taskmng"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"strconv"
	"strings"
	"sync"
)

func (m *Manager) tryScheduleTask() error {

	// TODO 一: 需要提前判断下 本地资源 的实况 和 全网的资源实况, 再做决定是否开启调度 ??
	// TODO 二: 共识中是否区分 interrupt 中止, 和 end 终止 ?? ---- 中止的task可以 repush, 而 终止的task直接 remove ??

	task, taskId, err := m.scheduler.TrySchedule()
	if nil == err && nil == task {
		return nil
	} else if nil != err && err == schedule.ErrAbandonTaskWithNotFoundTask {
		m.scheduler.RemoveTask(taskId)
	}  else if nil != err && err == schedule.ErrAbandonTaskWithNotFoundPowerPartyIds {
		m.scheduler.RemoveTask(taskId)
		m.eventEngine.StoreEvent(m.eventEngine.GenerateEvent(ev.TaskScheduleFailed.GetType(),
			task.GetTaskId(), task.GetTaskSender().GetIdentityId(),
			task.GetTaskSender().GetPartyId(), schedule.ErrAbandonTaskWithNotFoundPowerPartyIds.Error()))
		return m.sendNeedExecuteTaskByAction(task.GetTaskId(),
			apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Sender,
			task.GetTaskSender(), task.GetTaskSender(),
			types.TaskScheduleFailed)

	} else if nil != err {
		if nil != task {

			m.eventEngine.StoreEvent(m.eventEngine.GenerateEvent(ev.TaskScheduleFailed.GetType(),
				task.GetTaskId(), task.GetTaskSender().GetIdentityId(),
				task.GetTaskSender().GetPartyId(), err.Error()))

			if er := m.scheduler.RepushTask(task); er == schedule.ErrRescheduleLargeThreshold {
				log.WithError(err).Errorf("Failed to repush local task into queue/starve queue after trySchedule failed %s on `taskManager.tryScheduleTask()`, taskId: {%s}",
					err, task.GetTaskId())

				m.scheduler.RemoveTask(task.GetTaskId())
				m.eventEngine.StoreEvent(m.eventEngine.GenerateEvent(ev.TaskScheduleFailed.GetType(),
					task.GetTaskId(), task.GetTaskSender().GetIdentityId(),
					task.GetTaskSender().GetPartyId(), schedule.ErrRescheduleLargeThreshold.Error()))
				m.sendNeedExecuteTaskByAction(task.GetTaskId(),
					apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Sender,
					task.GetTaskSender(), task.GetTaskSender(),
					types.TaskScheduleFailed)
			}
		}
		return err
	}


	go func(task *types.Task) {

		log.Debugf("Start `NEED-CONSENSUS` task to 2pc consensus engine on `taskManager.tryScheduleTask()`, taskId: {%s}", task.GetTaskId())

		if err := m.consensusEngine.OnPrepare(task); nil != err {
			log.WithError(err).Errorf("Failed to call `OnPrepare()` of 2pc consensus engine on `taskManager.tryScheduleTask()`, taskId: {%s}", task.GetTaskId())
			// re push task into queue ,if anything else
			if err := m.scheduler.RepushTask(task); err == schedule.ErrRescheduleLargeThreshold {
				log.WithError(err).Errorf("Failed to repush local task into queue/starve queue after call `consensus.onPrepare()` on `taskManager.tryScheduleTask()`, taskId: {%s}",
					task.GetTaskId())

				m.scheduler.RemoveTask(task.GetTaskId())
				m.eventEngine.StoreEvent(m.eventEngine.GenerateEvent(ev.TaskScheduleFailed.GetType(),
					task.GetTaskId(), task.GetTaskSender().GetIdentityId(),
					task.GetTaskSender().GetPartyId(), schedule.ErrRescheduleLargeThreshold.Error()))
				m.sendNeedExecuteTaskByAction(task.GetTaskId(),
					apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Sender,
					task.GetTaskSender(), task.GetTaskSender(),
					types.TaskScheduleFailed)
			} else {
				log.Debugf("Succeed to repush local task into queue/starve queue after call `consensus.onPrepare()` on `taskManager.tryScheduleTask()`, taskId: {%s}",
					task.GetTaskId())
			}
			return
		}


		nonConsTask := types.NewNeedConsensusTask(task)
		if err := m.consensusEngine.OnHandle(nonConsTask.GetTask(), nonConsTask.GetResultCh()); nil != err {
			log.WithError(err).Errorf("Failed to call `OnHandle()` of 2pc consensus engine on `taskManager.tryScheduleTask()`, taskId: {%s}", nonConsTask.GetTask().GetTaskId())
		}

		result := nonConsTask.ReceiveResult()
		log.Debugf("Received `NEED-CONSENSUS` task result from 2pc consensus engine on `taskManager.tryScheduleTask()`, taskId: {%s}, result: {%s}", nonConsTask.GetTask().GetTaskId(), result.String())

		// store task event
		var (
			content   string
			eventType string
		)

		switch result.GetStatus() {
		case types.TaskTerminate, types.TaskConsensusFinished:
			if result.GetStatus() == types.TaskTerminate {
				content = "task was terminated."
				eventType = ev.TaskTerminated.GetType()
			} else {
				content = "succeed consensus."
				eventType = ev.TaskSucceedConsensus.GetType()
			}
		case types.TaskConsensusInterrupt:
			if nil != result.GetErr() {
				content = result.GetErr().Error()
			} else {
				content = "consensus was interrupted."
			}
			eventType = ev.TaskFailedConsensus.GetType()
		}

		// store task consensus result (failed or succeed) event with sender party
		m.resourceMng.GetDB().StoreTaskEvent(&libtypes.TaskEvent{
			Type:       eventType,
			TaskId:     nonConsTask.GetTask().GetTaskId(),
			IdentityId: nonConsTask.GetTask().GetTaskSender().GetIdentityId(),
			PartyId:    nonConsTask.GetTask().GetTaskSender().GetPartyId(),
			Content:    content,
			CreateAt:   timeutils.UnixMsecUint64(),
		})

		// received status must be `TaskConsensusFinished|TaskConsensusInterrupt|TaskTerminate` from consensus engine
		// never be `TaskNeedExecute|TaskScheduleFailed`
		//
		// Consensus failed, task needs to be suspended and rescheduled
		switch result.GetStatus() {
		case types.TaskTerminate, types.TaskConsensusFinished:
			// remove task from scheduler.queue|starvequeue after task consensus succeed
			// Don't send needexecuteTask, because that will be handle in `2pc engine.driveTask()`
			if err := m.scheduler.RemoveTask(result.GetTaskId()); nil != err {
				log.WithError(err).Errorf("Failed to remove local task from queue/starve queue after %s on `taskManager.tryScheduleTask()`, taskId: {%s}",
					result.GetStatus().String(), nonConsTask.GetTask().GetTaskId())
			}
			m.sendNeedExecuteTaskByAction(nonConsTask.GetTask().GetTaskId(),
				apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Sender,
				nonConsTask.GetTask().GetTaskSender(), nonConsTask.GetTask().GetTaskSender(),
				result.GetStatus())
		case types.TaskConsensusInterrupt:
			// re push task into queue ,if anything else
			if err := m.scheduler.RepushTask(nonConsTask.GetTask()); err == schedule.ErrRescheduleLargeThreshold {
				log.WithError(err).Errorf("Failed to repush local task into queue/starve queue after task consensus %s on `taskManager.tryScheduleTask()`, taskId: {%s}",
					result.GetStatus().String(), nonConsTask.GetTask().GetTaskId())

				m.scheduler.RemoveTask(nonConsTask.GetTask().GetTaskId())
				m.eventEngine.StoreEvent(m.eventEngine.GenerateEvent(ev.TaskScheduleFailed.Type,
					nonConsTask.GetTask().GetTaskId(), nonConsTask.GetTask().GetTaskSender().GetIdentityId(),
					nonConsTask.GetTask().GetTaskSender().GetPartyId(), schedule.ErrRescheduleLargeThreshold.Error()))
				m.sendNeedExecuteTaskByAction(nonConsTask.GetTask().GetTaskId(),
					apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Sender,
					nonConsTask.GetTask().GetTaskSender(), nonConsTask.GetTask().GetTaskSender(),
					types.TaskScheduleFailed)
			} else {
				log.Debugf("Succeed to repush local task into queue/starve queue after task cnsensus %s on `taskManager.tryScheduleTask()`, taskId: {%s}",
					result.GetStatus().String(), nonConsTask.GetTask().GetTaskId())
			}
		}
	}(task)

	return nil
}

func (m *Manager) sendNeedExecuteTaskByAction(taskId string, senderRole, receiverRole apicommonpb.TaskRole,
	sender, receiver *apicommonpb.TaskOrganization, taskActionStatus types.TaskActionStatus) error {
	m.needExecuteTaskCh <- types.NewNeedExecuteTask(
		"",
		common.Hash{},
		senderRole,
		receiverRole,
		sender,
		receiver,
		taskId,
		taskActionStatus,
		&types.PrepareVoteResource{},   // zero value
		&twopcpb.ConfirmTaskPeerInfo{}, // zero value
	)
	return nil
}

// To execute task
func (m *Manager) driveTaskForExecute(task *types.NeedExecuteTask) error {

	switch task.GetLocalTaskRole() {
	case apicommonpb.TaskRole_TaskRole_DataSupplier, apicommonpb.TaskRole_TaskRole_Receiver:
		return m.executeTaskOnDataNode(task)
	case apicommonpb.TaskRole_TaskRole_PowerSupplier:
		return m.executeTaskOnJobNode(task)
	default:
		log.Errorf("Faided to driveTaskForExecute(), Unknown task role, taskId: {%s}, taskRole: {%s}", task.GetTaskId(), task.GetLocalTaskRole().String())
		return fmt.Errorf("Unknown fighter node type when ready to execute task")
	}
}

func (m *Manager) executeTaskOnDataNode(task *types.NeedExecuteTask) error {

	// find dataNodeId with self vote
	var dataNodeId string
	dataNodes, err := m.resourceMng.GetDB().QueryRegisterNodeList(api.PrefixTypeDataNode)
	if nil != err {
		log.Errorf("Failed to query internal dataNode arr on `taskManager.executeTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId())
		return fmt.Errorf("query internal dataNodes failed when ready to execute task")
	}
	for _, dataNode := range dataNodes {
		if dataNode.GetExternalIp() == task.GetLocalResource().GetIp() && dataNode.GetExternalPort() == task.GetLocalResource().GetPort() {
			dataNodeId = dataNode.GetId()
			break
		}
	}

	if "" == strings.Trim(dataNodeId, "") {
		log.Errorf("Failed to find dataNodeId of self vote resource on `taskManager.executeTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
		return fmt.Errorf("not find dataNodeId of self vote resource when ready to execute task")
	}

	// clinet *grpclient.DataNodeClient,
	client, has := m.resourceClientSet.QueryDataNodeClient(dataNodeId)
	if !has {
		log.Errorf("Failed to query internal data node on `taskManager.executeTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
		return fmt.Errorf("dataNode client not found when ready to execute task")
	}
	if client.IsNotConnected() {
		if err := client.Reconnect(); nil != err {
			log.WithError(err).Errorf("Failed to connect internal data node on `taskManager.executeTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
				task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
			return fmt.Errorf("the dataNode is not connected when ready to execute task")
		}
	}

	req, err := m.makeTaskReadyGoReq(task)
	if nil != err {
		log.WithError(err).Errorf("Falied to make TaskReadyGoReq on `taskManager.executeTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
		return fmt.Errorf("make dataNode rpc req failed when ready to execute task")
	}

	resp, err := client.HandleTaskReadyGo(req)
	if nil != err {
		log.WithError(err).Errorf("Falied to call publish schedTask to `data-Fighter` node to executing, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
		return fmt.Errorf("call dataNode rpc api failed when ready to execute task")
	}
	if !resp.GetOk() {
		log.Errorf("Falied to executing task from `data-Fighter` node to executing the resp code is not `ok`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}, resp: %s",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId, resp.String())
		return fmt.Errorf("dataNode rpc api wrong resp %s when ready to execute task", resp.GetMsg())
	}

	log.Infof("Success to publish schedTask to `data-Fighter` node to executing, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
		task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
	return nil
}

func (m *Manager) executeTaskOnJobNode(task *types.NeedExecuteTask) error {

	// find jobNodeId with self vote
	var jobNodeId string
	jobNodes, err := m.resourceMng.GetDB().QueryRegisterNodeList(api.PrefixTypeJobNode)
	if nil != err {
		log.Errorf("Failed to query internal jobNode arr on `taskManager.executeTaskOnJobNode()`, taskId: {%s}, role: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId())
		return fmt.Errorf("query internal jobNodes failed when ready to execute task")
	}
	for _, jobNode := range jobNodes {
		if jobNode.GetExternalIp() == task.GetLocalResource().GetIp() && jobNode.GetExternalPort() == task.GetLocalResource().GetPort() {
			jobNodeId = jobNode.GetId()
			break
		}
	}

	if "" == strings.Trim(jobNodeId, "") {
		log.Errorf("Failed to find jobNodeId of self vote resource on `taskManager.executeTaskOnJobNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
		return fmt.Errorf("not find jobNodeId of self vote resource when ready to execute task")
	}

	// clinet *grpclient.JobNodeClient,
	client, has := m.resourceClientSet.QueryJobNodeClient(jobNodeId)
	if !has {
		log.Errorf("Failed to query internal job node on `taskManager.executeTaskOnJobNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
		return fmt.Errorf("jobNode client not found when ready to execute task")
	}
	if client.IsNotConnected() {
		if err := client.Reconnect(); nil != err {
			log.WithError(err).Errorf("Failed to connect internal job node on `taskManager.executeTaskOnJobNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
				task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
			return fmt.Errorf("the jobNode is not connected when ready to execute task")
		}
	}

	req, err := m.makeTaskReadyGoReq(task)
	if nil != err {
		log.WithError(err).Errorf("Falied to make TaskReadyGoReq on `taskManager.executeTaskOnJobNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
		return fmt.Errorf("make jobNode rpc req failed when ready to execute task")
	}

	resp, err := client.HandleTaskReadyGo(req)
	if nil != err {
		log.WithError(err).Errorf("Falied to publish schedTask to `job-Fighter` node to executing, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
		return fmt.Errorf("call jobNode rpc api failed when ready to execute task")
	}
	if !resp.GetOk() {
		log.Errorf("Falied to publish schedTask to `job-Fighter` node to executing the resp code is not `ok`, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}, resp: %s",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId, resp.String())
		return fmt.Errorf("jobNode rpc api wrong resp %s when ready to execute task", resp.GetMsg())
	}

	log.Infof("Success to publish schedTask to `job-Fighter` node to executing, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
		task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
	return nil
}

// To terminate task
func (m *Manager) driveTaskForTerminate(task *types.NeedExecuteTask) error {

	// annonce fighter processor to terminate this task
	switch task.GetLocalTaskRole() {
	case apicommonpb.TaskRole_TaskRole_DataSupplier, apicommonpb.TaskRole_TaskRole_Receiver:
		return m.terminateTaskOnDataNode(task)
	case apicommonpb.TaskRole_TaskRole_PowerSupplier:
		return m.terminateTaskOnJobNode(task)
	default:
		log.Errorf("Faided to driveTaskForTerminate(), Unknown task role, taskId: {%s}, taskRole: {%s}", task.GetTaskId(), task.GetLocalTaskRole().String())
		return fmt.Errorf("Unknown resource node type when ready to terminate task")
	}
}

func (m *Manager) terminateTaskOnDataNode(task *types.NeedExecuteTask) error {

	// find dataNodeId with self vote
	var dataNodeId string
	dataNodes, err := m.resourceMng.GetDB().QueryRegisterNodeList(api.PrefixTypeDataNode)
	if nil != err {
		log.Errorf("Failed to query internal dataNode arr on `taskManager.terminateTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId())
		return fmt.Errorf("query internal dataNodes failed when ready to terminate task")
	}
	for _, dataNode := range dataNodes {
		if dataNode.GetExternalIp() == task.GetLocalResource().GetIp() && dataNode.GetExternalPort() == task.GetLocalResource().GetPort() {
			dataNodeId = dataNode.GetId()
			break
		}
	}

	if "" == strings.Trim(dataNodeId, "") {
		log.Errorf("Failed to find dataNodeId of self vote resource on `taskManager.terminateTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
		return fmt.Errorf("not find dataNodeId of self vote resource when ready to terminate task")
	}

	// clinet *grpclient.DataNodeClient,
	client, has := m.resourceClientSet.QueryDataNodeClient(dataNodeId)
	if !has {
		log.Errorf("Failed to query internal data node on `taskManager.terminateTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
		return fmt.Errorf("dataNode client not found when ready to terminate task")
	}
	if client.IsNotConnected() {
		if err := client.Reconnect(); nil != err {
			log.WithError(err).Errorf("Failed to connect internal data node on `taskManager.terminateTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
				task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
			return fmt.Errorf("the dataNode is not connected when ready to terminate task")
		}
	}

	req, err := m.makeTerminateTaskReq(task)
	if nil != err {
		log.WithError(err).Errorf("Falied to make TaskCancelReq on `taskManager.terminateTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
		return fmt.Errorf("make dataNode rpc req failed when ready to terminate task")
	}

	resp, err := client.HandleCancelTask(req)
	if nil != err {
		log.WithError(err).Errorf("Falied to call publish schedTask to `data-Fighter` node to terminating, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
		return fmt.Errorf("call dataNode rpc api failed when ready to terminate task")
	}
	if !resp.GetOk() {
		log.Errorf("Falied to executing task from `data-Fighter` node to terminating the resp code is not `ok`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}, resp: %s",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId, resp.String())
		return fmt.Errorf("dataNode rpc api wrong resp %s when ready to terminate task", resp.GetMsg())
	}

	log.Infof("Success to publish schedTask to `data-Fighter` node to terminating, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
		task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
	return nil
}

func (m *Manager) terminateTaskOnJobNode(task *types.NeedExecuteTask) error {

	// find jobNodeId with self vote
	var jobNodeId string
	jobNodes, err := m.resourceMng.GetDB().QueryRegisterNodeList(api.PrefixTypeJobNode)
	if nil != err {
		log.Errorf("Failed to query internal jobNode arr on `taskManager.terminateTaskOnJobNode()`, taskId: {%s}, role: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId())
		return fmt.Errorf("query internal jobNodes failed when ready to terminate task")
	}
	for _, jobNode := range jobNodes {
		if jobNode.GetExternalIp() == task.GetLocalResource().GetIp() && jobNode.GetExternalPort() == task.GetLocalResource().GetPort() {
			jobNodeId = jobNode.GetId()
			break
		}
	}

	if "" == strings.Trim(jobNodeId, "") {
		log.Errorf("Failed to find jobNodeId of self vote resource on `taskManager.terminateTaskOnJobNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
		return fmt.Errorf("not find jobNodeId of self vote resource when ready to terminate task")
	}

	// clinet *grpclient.JobNodeClient,
	client, has := m.resourceClientSet.QueryJobNodeClient(jobNodeId)
	if !has {
		log.Errorf("Failed to query internal job node on `taskManager.terminateTaskOnJobNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
		return fmt.Errorf("jobNode client not found when ready to terminate task")
	}
	if client.IsNotConnected() {
		if err := client.Reconnect(); nil != err {
			log.WithError(err).Errorf("Failed to connect internal job node on `taskManager.terminateTaskOnJobNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
				task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
			return fmt.Errorf("the jobNode is not connected when ready to terminate task")
		}
	}

	req, err := m.makeTerminateTaskReq(task)
	if nil != err {
		log.WithError(err).Errorf("Falied to make TaskCancelReq on `taskManager.terminateTaskOnJobNode()`,taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
		return fmt.Errorf("make jobNode rpc req failed when ready to terminate task")
	}

	resp, err := client.HandleCancelTask(req)
	if nil != err {
		log.WithError(err).Errorf("Falied to publish schedTask to `job-Fighter` node to terminating, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
		return fmt.Errorf("call jobNode rpc api failed when ready to terminate task")
	}
	if !resp.GetOk() {
		log.Errorf("Falied to publish schedTask to `job-Fighter` node to terminating the resp code is not `ok`,taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}, resp: %s",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId, resp.String())
		return fmt.Errorf("jobNode rpc api wrong resp %s when ready to terminate task", resp.GetMsg())
	}

	log.Infof("Success to publish schedTask to `job-Fighter` node to terminating, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
		task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
	return nil
}

func (m *Manager) publishFinishedTaskToDataCenter(task *types.NeedExecuteTask, localTask *types.Task, delay bool) {

	/**
	++++++++++++++++++++++++++++
	NOTE: the needExecuteTask must be sender's here. (task.GetLocalTaskOrganization is sender's identity)
	++++++++++++++++++++++++++++
	*/

	handleFn := func() {

		eventList, err := m.resourceMng.GetDB().QueryTaskEventList(task.GetTaskId())
		if nil != err {
			log.WithError(err).Errorf("Failed to Query all task event list for sending datacenter on publishFinishedTaskToDataCenter, taskId: {%s}",
				task.GetTaskId())
			return
		}

		// check all events of this task, and change task state finally.
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

		log.Debugf("Start publishFinishedTaskToDataCenter, taskId: {%s}, taskState: {%s}", task.GetTaskId(), taskState.String())

		finalTask := m.convertScheduleTaskToTask(localTask, eventList, taskState)
		if err := m.resourceMng.GetDB().InsertTask(finalTask); nil != err {
			log.WithError(err).Errorf("Failed to save task to datacenter on publishFinishedTaskToDataCenter, taskId: {%s}, partyId: {%s}",
				task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())
			return
		}
		if err := m.RemoveExecuteTaskStateAfterExecuteTask("on taskManager.publishFinishedTaskToDataCenter()", task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), resource.SetAllReleaseResourceOption()); nil != err {
			log.WithError(err).Errorf("Failed to call RemoveExecuteTaskStateAfterExecuteTask() on publishFinishedTaskToDataCenter(), taskId: {%s}, taskRole: {%s},  partyId: {%s}, remote pid: {%s}",
				task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())
			return
		}

		log.Debugf("Finished pulishFinishedTaskToDataCenter, taskId: {%s}, partyId: {%s}, taskState: {%s}",
			task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), taskState)
	}

	if delay {
		// delays handling some logic. (default delay 10s)
		runutil.RunOnce(context.TODO(), senderExecuteTaskExpire, handleFn)
	} else {
		handleFn()
	}
}
func (m *Manager) sendTaskResultMsgToTaskSender(task *types.NeedExecuteTask) {

	log.Debugf("Start sendTaskResultMsgToTaskSender, taskId: {%s}, taskRole: {%s},  partyId: {%s}, remote pid: {%s}",
		task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())

	var option resource.ReleaseResourceOption

	// when other task partner and task sender is same identity,
	// we don't need to removed local task and local eventList
	if task.GetLocalTaskOrganization().GetIdentityId() == task.GetRemoteTaskOrganization().GetIdentityId() {
		option = resource.SetUnlockLocalResorce() // unlock local resource of partyId, but don't remove local task and events of partyId
	} else {
		option = resource.SetAllReleaseResourceOption() // unlock local resource and remove local task and events
		// broadcast `task result msg` to reply remote peer
		if err := m.p2p.Broadcast(context.TODO(), m.makeTaskResultMsgWithEventList(task)); nil != err {
			log.WithError(err).Errorf("failed to call `SendTaskResultMsg`, taskId: {%s}, taskRole: {%s},  partyId: {%s}, remote pid: {%s}",
				task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())
		}
	}
	if err := m.RemoveExecuteTaskStateAfterExecuteTask("on sendTaskResultMsgToTaskSender()", task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), option); nil != err {
		log.WithError(err).Errorf("Failed to call RemoveExecuteTaskStateAfterExecuteTask() on sendTaskResultMsgToTaskSender(), taskId: {%s}, taskRole: {%s},  partyId: {%s}, remote pid: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())
		return
	}
	log.Debugf("Finished sendTaskResultMsgToTaskSender, taskId: {%s}, taskRole: {%s},  partyId: {%s}, remote pid: {%s}",
		task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())
}


func (m *Manager) StoreExecuteTaskStateBeforeExecuteTask (logdesc, taskId, partyId string) error {
	// Store task exec status
	if err := m.resourceMng.GetDB().StoreLocalTaskExecuteStatusValExecByPartyId(taskId, partyId); nil != err {
		log.WithError(err).Errorf("Failed to store local task about `exec` status %s, taskId: {%s}, partyId: {%s}",
			logdesc, taskId, partyId)
		return err
	}
	used, err := m.resourceMng.GetDB().QueryLocalTaskPowerUsed(taskId, partyId)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed to call QueryLocalTaskPowerUsed %s, used: {%s}", logdesc, used.String())
		return err
	}
	if nil == err && nil != used {
		if err := m.resourceMng.StoreJobNodeExecuteTaskId(used.GetNodeId(), used.GetTaskId(), used.GetPartyId()); nil != err {
			log.WithError(err).Errorf("Failed to store execute taskId into jobNode %s, jobNodeId: {%s}, taskId: {%s}, partyId: {%s}",
				logdesc, used.GetNodeId(), taskId, partyId)
			return err
		}
	}
	return nil
}

func (m *Manager) RemoveExecuteTaskStateAfterExecuteTask(logdesc, taskId, partyId string, option resource.ReleaseResourceOption) error {
	if err := m.resourceMng.GetDB().RemoveLocalTaskExecuteStatusByPartyId(taskId, partyId); nil != err {
		log.WithError(err).Errorf("Failed to remove task executing status %s, taskId: {%s} partyId: {%s}",
			logdesc, taskId, partyId)
		return err
	}
	used, err := m.resourceMng.GetDB().QueryLocalTaskPowerUsed(taskId, partyId)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed to call QueryLocalTaskPowerUsed %s, used: {%s}", logdesc, used.String())
		return err
	}
	if nil == err && nil != used {
		if err := m.resourceMng.RemoveJobNodeExecuteTaskId(used.GetNodeId(), used.GetTaskId(), used.GetPartyId()); nil != err {
			log.WithError(err).Errorf("Failed to remove execute taskId into jobNode %s, jobNodeId: {%s}, taskId: {%s}, partyId: {%s}",
				logdesc, used.GetNodeId(), taskId, partyId)
			return err
		}
	}
	// clean local task cache
	m.resourceMng.ReleaseLocalResourceWithTask(logdesc, taskId, partyId, option)
	return nil
}

func (m *Manager) sendTaskTerminateMsg(task *types.Task) error {

	sender := task.GetTaskSender()

	sendTerminateMsgFn := func(wg *sync.WaitGroup, sender, receiver *apicommonpb.TaskOrganization, senderRole, receiverRole apicommonpb.TaskRole, errCh chan<- error) {

		defer wg.Done()

		pid, err := p2p.HexPeerID(receiver.NodeId)
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId, taskId: %s, other peer's taskRole: %s, other peer's partyId: %s, other identityId: %s, pid: %s, err: %s",
				task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid, err)
			return
		}

		terminateMsg := &taskmngpb.TaskTerminateMsg{
			MsgOption: &msgcommonpb.MsgOption{
				ProposalId:      common.Hash{}.Bytes(),
				SenderRole:      uint64(senderRole),
				SenderPartyId:   []byte(sender.GetPartyId()),
				ReceiverRole:    uint64(receiverRole),
				ReceiverPartyId: []byte(receiver.GetPartyId()),
				MsgOwner: &msgcommonpb.TaskOrganizationIdentityInfo{
					Name:       []byte(sender.GetNodeName()),
					NodeId:     []byte(sender.GetNodeId()),
					IdentityId: []byte(sender.GetIdentityId()),
					PartyId:    []byte(sender.GetPartyId()),
				},
			},
			TaskId:   []byte(task.GetTaskId()),
			CreateAt: timeutils.UnixMsecUint64(),
			Sign:     nil,
		}

		//var sendErr error
		var logdesc string
		if types.IsSameTaskOrg(sender, receiver) {
			m.onTaskTerminateMsg(pid, terminateMsg, types.LocalNetworkMsg)
			logdesc = "OnTaskTerminateMsg()"
		} else {
			m.p2p.Broadcast(context.TODO(), terminateMsg)
			logdesc = "Broadcast()"
		}

		//// Send the ConfirmMsg to other peer
		//if nil != sendErr {
		//	errCh <- fmt.Errorf("failed to call`sendTaskTerminateMsg.%s` taskId: %s, other peer's taskRole: %s, other peer's partyId: %s, other identityId: %s, pid: %s, err: %s",
		//		logdesc, task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid, err)
		//	errCh <- err
		//	return
		//}

		log.Debugf("Succceed to call`sendTaskTerminateMsg.%s` taskId: %s, other peer's taskRole: %s, other peer's partyId: %s, other identityId: %s, pid: %s",
			logdesc, task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid)

	}

	size := (len(task.GetTaskData().GetDataSuppliers())) + len(task.GetTaskData().GetPowerSuppliers()) + len(task.GetTaskData().GetReceivers())
	errCh := make(chan error, size)
	var wg sync.WaitGroup

	for i := 0; i < len(task.GetTaskData().GetDataSuppliers()); i++ {

		wg.Add(1)
		dataSupplier := task.GetTaskData().GetDataSuppliers()[i]
		receiver := dataSupplier.GetOrganization()
		go sendTerminateMsgFn(&wg, sender, receiver, apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_DataSupplier, errCh)

	}
	for i := 0; i < len(task.GetTaskData().GetPowerSuppliers()); i++ {

		wg.Add(1)
		powerSupplier := task.GetTaskData().GetPowerSuppliers()[i]
		receiver := powerSupplier.GetOrganization()
		go sendTerminateMsgFn(&wg, sender, receiver, apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_PowerSupplier, errCh)

	}

	for i := 0; i < len(task.GetTaskData().GetReceivers()); i++ {

		wg.Add(1)
		receiver := task.GetTaskData().GetReceivers()[i]
		go sendTerminateMsgFn(&wg, sender, receiver, apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Receiver, errCh)
	}

	wg.Wait()
	close(errCh)

	errStrs := make([]string, 0)

	for err := range errCh {
		if nil != err {
			errStrs = append(errStrs, err.Error())
		}
	}
	if len(errStrs) != 0 {
		return fmt.Errorf(
			"\n######################################################## \n%s\n########################################################\n",
			strings.Join(errStrs, "\n"))
	}
	return nil
}

func (m *Manager) sendLocalTaskToScheduler(tasks types.TaskDataArray) {
	m.localTasksCh <- tasks
	log.Infof("Succeed send to scheduler tasks count on taskManager.SendTaskMsgArr(), task arr len: %d", len(tasks))
}
func (m *Manager) sendTaskEvent(event *libtypes.TaskEvent) {
	m.eventCh <- event
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

func (m *Manager) makeTaskReadyGoReq(task *types.NeedExecuteTask) (*fightercommon.TaskReadyGoReq, error) {

	localTask, err := m.resourceMng.GetDB().QueryLocalTask(task.GetTaskId())
	if nil != err {
		return nil, fmt.Errorf("query local task info failed")
	}

	var dataPartyArr []string
	var powerPartyArr []string
	var receiverPartyArr []string

	peerList := make([]*fightercommon.TaskReadyGoReq_Peer, 0)

	for _, dataSupplier := range task.GetResources().GetDataSupplierPeerInfos() {
		portStr := string(dataSupplier.GetPort())
		port, err := strconv.Atoi(portStr)
		if nil != err {
			return nil, err
		}
		peerList = append(peerList, &fightercommon.TaskReadyGoReq_Peer{
			Ip:      string(dataSupplier.GetIp()),
			Port:    int32(port),
			PartyId: string(dataSupplier.GetPartyId()),
		})
		dataPartyArr = append(dataPartyArr, string(dataSupplier.GetPartyId()))
	}

	for _, powerSupplier := range task.GetResources().GetPowerSupplierPeerInfos() {
		portStr := string(powerSupplier.GetPort())
		port, err := strconv.Atoi(portStr)
		if nil != err {
			return nil, err
		}
		peerList = append(peerList, &fightercommon.TaskReadyGoReq_Peer{
			Ip:      string(powerSupplier.GetIp()),
			Port:    int32(port),
			PartyId: string(powerSupplier.GetPartyId()),
		})

		powerPartyArr = append(powerPartyArr, string(powerSupplier.GetPartyId()))
	}

	for _, receiver := range task.GetResources().GetResultReceiverPeerInfos() {
		portStr := string(receiver.GetPort())
		port, err := strconv.Atoi(portStr)
		if nil != err {
			return nil, err
		}
		peerList = append(peerList, &fightercommon.TaskReadyGoReq_Peer{
			Ip:      string(receiver.GetIp()),
			Port:    int32(port),
			PartyId: string(receiver.GetPartyId()),
		})

		receiverPartyArr = append(receiverPartyArr, string(receiver.GetPartyId()))
	}

	contractExtraParams, err := m.makeContractParams(task, localTask)
	if nil != err {
		return nil, fmt.Errorf("make contractParams failed, %s", err)
	}
	log.Debugf("Succeed make contractCfg, taskId:{%s}, contractCfg: %s", task.GetTaskId(), contractExtraParams)
	return &fightercommon.TaskReadyGoReq{
		TaskId:     task.GetTaskId(),
		ContractId: localTask.GetTaskData().GetCalculateContractCode(),
		//DataId: "",
		PartyId: task.GetLocalTaskOrganization().GetPartyId(),
		//EnvId: "",
		Peers:            peerList,
		ContractCfg:      contractExtraParams,
		DataParty:        dataPartyArr,
		ComputationParty: powerPartyArr,
		ResultParty:      receiverPartyArr,
		Duration:         localTask.GetTaskData().GetOperationCost().GetDuration(),
	}, nil
}

func (m *Manager) makeContractParams(task *types.NeedExecuteTask, localTask *types.Task) (string, error) {

	partyId := task.GetLocalTaskOrganization().GetPartyId()

	var filePath string
	var keyColumn string
	var selectedColumns []string

	if task.GetLocalTaskRole() == apicommonpb.TaskRole_TaskRole_DataSupplier {

		var find bool

		for _, dataSupplier := range localTask.GetTaskData().GetDataSuppliers() {
			if partyId == dataSupplier.GetOrganization().GetPartyId() {

				userType := localTask.GetTaskData().GetUserType()
				user := localTask.GetTaskData().GetUser()
				metadataId := dataSupplier.GetMetadataId()

				// verify metadataAuth first

				if err := m.authMng.VerifyMetadataAuth(userType, user, metadataId); nil != err {
					return "", fmt.Errorf("verify user metadataAuth failed, %s, userType: {%s}, user: {%s}, metadataId: {%s}",
						err, userType, user, metadataId)
				}

				internalMetadataFlag, err := m.resourceMng.GetDB().IsInternalMetadataByDataId(metadataId)
				if nil != err {
					return "", fmt.Errorf("check metadata whether internal metadata failed %s, metadataId: {%s}", err, metadataId)
				}

				var metadata *types.Metadata

				// whether the metadata is internal metadata ?
				if internalMetadataFlag {
					// query internal metadata
					metadata, err = m.resourceMng.GetDB().QueryInternalMetadataByDataId(metadataId)
					if nil != err {
						return "", fmt.Errorf("query internale metadata failed %s, metadataId: {%s}", err, metadataId)
					}
				} else {
					// query published metadata
					metadata, err = m.resourceMng.GetDB().QueryMetadataByDataId(metadataId)
					if nil != err {
						return "", fmt.Errorf("query publish metadata failed %s, metadataId: {%s}", err, metadataId)
					}
				}

				filePath = metadata.GetData().GetFilePath()

				keyColumn = dataSupplier.GetKeyColumn().GetCName()
				selectedColumns = make([]string, len(dataSupplier.GetSelectedColumns()))
				for i, col := range dataSupplier.GetSelectedColumns() {
					selectedColumns[i] = col.GetCName()
				}

				// only consume metadata auth when metadata is not internal metadata.
				if !internalMetadataFlag {
					// query metadataAuthId by metadataId
					metadataAuthId, err := m.authMng.QueryMetadataAuthIdByMetadataId(userType, user, metadataId)
					if nil != err {
						return "", fmt.Errorf("query metadataAuthId failed %s, metadataId: {%s}", err, metadataId)
					}
					// ConsumeMetadataAuthority
					if err = m.authMng.ConsumeMetadataAuthority(metadataAuthId); nil != err {
						return "", fmt.Errorf("consume metadataAuth failed %s, metadataAuthId: {%s}", err, metadataAuthId)
					}
				}

				find = true
				break
			}
		}

		if !find {
			return "", fmt.Errorf("can not find the dataSupplier for find originFilePath, taskId: {%s}, self.IdentityId: {%s}, seld.PartyId: {%s}",
				task.GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(), task.GetLocalTaskOrganization().GetPartyId())
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
	log.Debugf("Start json Unmarshal the `ContractExtraParams`, taskId: {%s}, ContractExtraParams: %s", task.GetTaskId(), localTask.GetTaskData().GetContractExtraParams())
	if "" != localTask.GetTaskData().GetContractExtraParams() {
		if err := json.Unmarshal([]byte(localTask.GetTaskData().GetContractExtraParams()), &dynamicParameter); nil != err {
			return "", fmt.Errorf("can not json Unmarshal the `ContractExtraParams` of task %s, taskId: {%s}, self.IdentityId: {%s}, seld.PartyId: {%s}",
				err, task.GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(), task.GetLocalTaskOrganization().GetPartyId())
		}
	}
	req.DynamicParameter = dynamicParameter

	b, err := json.Marshal(req)
	if nil != err {
		return "", fmt.Errorf("can not json Marshal the `FighterTaskReadyGoReqContractCfg` %s, taskId: {%s}, self.IdentityId: {%s}, seld.PartyId: {%s}",
			err, task.GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(), task.GetLocalTaskOrganization().GetPartyId())
	}
	return string(b), nil
}

// make terminate rpc req
func (m *Manager) makeTerminateTaskReq(task *types.NeedExecuteTask) (*fightercommon.TaskCancelReq, error) {
	return &fightercommon.TaskCancelReq{
		TaskId:  task.GetTaskId(),
		PartyId: task.GetLocalTaskOrganization().GetPartyId(),
	}, nil
}

func (m *Manager) addNeedExecuteTaskCache(task *types.NeedExecuteTask) {
	m.runningTaskCacheLock.Lock()
	cache, ok := m.runningTaskCache[task.GetTaskId()]
	if !ok {
		cache = make(map[string]*types.NeedExecuteTask, 0)
	}
	cache[task.GetLocalTaskOrganization().GetPartyId()] = task
	m.runningTaskCache[task.GetTaskId()] = cache
	if err := m.resourceMng.GetDB().StoreNeedExecuteTask(task); nil != err {
		log.WithError(err).Errorf("store needExecuteTask failed, taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())
	}
	m.runningTaskCacheLock.Unlock()
}

func (m *Manager) removeNeedExecuteTask(taskId string) {
	m.runningTaskCacheLock.Lock()
	delete(m.runningTaskCache, taskId)
	go m.resourceMng.GetDB().RemoveNeedExecuteTask(taskId)
	m.runningTaskCacheLock.Unlock()
}

func (m *Manager) removeNeedExecuteTaskCache(taskId, partyId string) {
	m.runningTaskCacheLock.Lock()
	defer m.runningTaskCacheLock.Unlock()
	cache, ok := m.runningTaskCache[taskId]
	if !ok {
		return
	}
	delete(cache, partyId)
	go m.resourceMng.GetDB().RemoveNeedExecuteTaskByPartyId(taskId, partyId)
	if len(cache) == 0 {
		delete(m.runningTaskCache, taskId) // delete empty map[partyId]task
	} else {
		m.runningTaskCache[taskId] = cache // restore map[partyId]task if it is not empty
	}
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
	defer m.runningTaskCacheLock.Unlock()
	for taskId, cache := range m.runningTaskCache {
		for _, task := range cache {
			if ok := f(taskId, task); ok {
			}
		}
	}

}

func (m *Manager) makeTaskResultMsgWithEventList(task *types.NeedExecuteTask) *taskmngpb.TaskResultMsg {

	if task.GetLocalTaskRole() == apicommonpb.TaskRole_TaskRole_Sender {
		log.Errorf("send task OR task owner can not make taskResultMsg")
		return nil
	}

	eventList, err := m.resourceMng.GetDB().QueryTaskEventListByPartyId(task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to make taskResultMsg with query task eventList, taskId {%s}, partyId {%s}",
			task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())
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

func (m *Manager) handleTaskEventWithCurrentIdentity(event *libtypes.TaskEvent) error {
	eventType := event.Type
	if len(eventType) != ev.EventTypeCharLen {
		return ev.IncEventType
	}

	identityId, err := m.resourceMng.GetDB().QueryIdentityId()
	if nil != err {
		log.WithError(err).Errorf("Failed to query self identityId on taskManager.SendTaskEvent()")
		return fmt.Errorf("query local identityId failed, %s", err)
	}
	event.IdentityId = identityId

	if task, ok := m.queryNeedExecuteTaskCache(event.GetTaskId(), event.GetPartyId()); ok {

		localTask, err := m.resourceMng.GetDB().QueryLocalTask(event.GetTaskId())
		if nil != err {
			log.WithError(err).Errorf("Failed to query local task info on `taskManager.handleTaskEventWithCurrentIdentity()`, taskId: {%s}, partyId: {%s}",
				event.GetTaskId(), event.GetPartyId())

			// remove wrong task cache
			m.removeNeedExecuteTaskCache(event.GetTaskId(), event.GetPartyId())
			return err
		}

		// need to validate the task that have been processing ? Maybe~
		// While task is consensus or executing, can terminate.
		has, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusByPartyId(event.GetTaskId(), event.GetPartyId())
		if nil != err {
			log.WithError(err).Errorf("Failed to check local task execute status whether exist on `taskManager.handleTaskEventWithCurrentIdentity()`, taskId: {%s}, partyId: {%s}",
				event.GetTaskId(), event.GetPartyId())

			// remove wrong task cache
			m.removeNeedExecuteTaskCache(event.GetTaskId(), event.GetPartyId())
			return err
		}

		if !has {
			log.Warnf("Warn ignore event, `event is the end` but not find party task executeStatus on `taskManager.handleTaskEventWithCurrentIdentity()`, event: %s",
				event.String())

			// remove wrong task cache
			m.removeNeedExecuteTaskCache(event.GetTaskId(), event.GetPartyId())
			return nil
		}

		switch event.Type {
		case ev.TaskExecuteSucceedEOF.Type:
			log.Infof("Started handle taskEvent with currentIdentity, `event is the task final succeed EOF finished`, event: %s", event.String())

			// store final event
			m.resourceMng.GetDB().StoreTaskEvent(event)
			m.resourceMng.GetDB().StoreTaskEvent(m.eventEngine.GenerateEvent(ev.TaskSucceed.Type, task.GetTaskId(), identityId, task.GetLocalTaskOrganization().GetPartyId(), "task execute succeed"))

			publish, err := m.checkTaskSenderPublishOpportunity(localTask, event)
			if nil != err {
				log.WithError(err).Errorf("Failed to check task sender publish opportunity on `taskManager.handleTaskEventWithCurrentIdentity()`, event: %s",
					event.GetPartyId())
				return err
			}

			if publish {

				//1、 handle last party
				// send this task result to remote target peer
				m.sendTaskResultMsgToTaskSender(task)
				m.removeNeedExecuteTaskCache(event.GetTaskId(), event.GetPartyId())

				//2、 handle sender party
				senderNeedTask, ok := m.queryNeedExecuteTaskCache(event.GetTaskId(), localTask.GetTaskSender().GetPartyId())
				if ok {
					log.Debugf("Need to call `publishFinishedTaskToDataCenter` on `taskManager.handleTaskEventWithCurrentIdentity()`, taskId: {%s}, sender partyId: {%s}",
						event.GetTaskId(), localTask.GetTaskSender().GetPartyId())
					// handle this task result with current peer
					m.publishFinishedTaskToDataCenter(senderNeedTask, localTask, true)
					m.removeNeedExecuteTaskCache(event.GetTaskId(), localTask.GetTaskSender().GetPartyId())
				}
			} else {
				// send this task result to remote target peer
				m.sendTaskResultMsgToTaskSender(task)
				m.removeNeedExecuteTaskCache(event.GetTaskId(), event.GetPartyId())
			}
		case ev.TaskExecuteFailedEOF.Type:
			log.Infof("Started handle taskEvent with currentIdentity, `event is the task final [FAILED] EOF finished` will terminate task, event: %s", event.String())

			//
			if err := m.resourceMng.GetDB().RemoveTaskPartnerPartyIds(event.GetTaskId()); nil != err {
				log.WithError(err).Errorf("Failed to remove all partyId of local task's partner arr on `taskManager.handleTaskEventWithCurrentIdentity()`, taskId: {%s}",
					event.GetTaskId())
			}

			// store final failed EOF event
			m.resourceMng.GetDB().StoreTaskEvent(event)
			if err := m.onTerminateExecuteTask(localTask); nil != err {
				log.Errorf("Failed to call `onTerminateExecuteTask()` on `taskManager.handleTaskEventWithCurrentIdentity()`, taskId: {%s}, err: \n%s", localTask.GetTaskId(), err)
			}
		default:
			log.Infof("Started handle taskEvent with currentIdentity, `event is not the end EOF`, event: %s", event.String())
			// It's not EOF event, then the task still executing, so store this event
			return m.resourceMng.GetDB().StoreTaskEvent(event)

		}
	}
	return nil // ignore event while task is not exist.
}

func (m *Manager) handleNeedExecuteTask(task *types.NeedExecuteTask, localTask *types.Task) {

	log.Debugf("Start handle needExecuteTask, taskId: {%s}, role: {%s}, partyId: {%s}",
		task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId())

	// Store task exec status
	if err := m.StoreExecuteTaskStateBeforeExecuteTask("on handleNeedExecuteTask()", task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId()); nil != err {
		log.WithError(err).Errorf("Failed to call StoreExecuteTaskStateBeforeExecuteTask() on handleNeedExecuteTask(), taskId: {%s}, role: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId())
		return
	}

	if localTask.GetTaskData().GetState() == apicommonpb.TaskState_TaskState_Pending &&
		localTask.GetTaskData().GetStartAt() == 0 {

		localTask.GetTaskData().State = apicommonpb.TaskState_TaskState_Running
		localTask.GetTaskData().StartAt = timeutils.UnixMsecUint64()
		if err := m.resourceMng.GetDB().StoreLocalTask(localTask); nil != err {
			log.WithError(err).Errorf("Failed to update local task state before executing task, taskId: {%s}, need update state: {%s}",
				task.GetTaskId(), apicommonpb.TaskState_TaskState_Running.String())
		}
	}

	// store local cache
	m.addNeedExecuteTaskCache(task)

	// The task sender will not execute the task
	if task.GetLocalTaskRole() != apicommonpb.TaskRole_TaskRole_Sender &&
		task.GetLocalTaskOrganization().GetPartyId() != localTask.GetTaskSender().GetPartyId() {
		// driving task to executing
		if err := m.driveTaskForExecute(task); nil != err {
			log.WithError(err).Errorf("Failed to execute task on internal node, taskId: {%s}, role: {%s}, partyId: {%s}",
				task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId())

			m.SendTaskEvent(m.eventEngine.GenerateEvent(ev.TaskFailed.GetType(), task.GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(), task.GetLocalTaskOrganization().GetPartyId(),
				fmt.Sprintf("execute task failed: %s with %s", err,
				task.GetLocalTaskOrganization().GetPartyId())))
		}
	}
}

func (m *Manager) expireTaskMonitor() {

	m.runningTaskCacheLock.Lock()
	defer m.runningTaskCacheLock.Unlock()

	for taskId, cache := range m.runningTaskCache {

		localTask, err := m.resourceMng.GetDB().QueryLocalTask(taskId)

		for partyId, task := range cache {

			if nil != err {

				log.WithError(err).Warnf("Warning query local task info failed, clean current party task cache short circuit AND skip it, on `taskManager.expireTaskMonitor()`, taskId: {%s}, partyId: {%s}",
					taskId, partyId)
				// clean current party task cache short circuit.
				delete(cache, partyId)
				go m.resourceMng.GetDB().RemoveNeedExecuteTaskByPartyId(taskId, partyId)
				if len(cache) == 0 {
					delete(m.runningTaskCache, taskId)
				} else {
					m.runningTaskCache[taskId] = cache
				}

				continue
			}

			if localTask.GetTaskData().GetState() == apicommonpb.TaskState_TaskState_Running && localTask.GetTaskData().GetStartAt() != 0 {

				// the task has running expire
				var duration uint64

				duration = timeutils.UnixMsecUint64() - localTask.GetTaskData().GetStartAt()

				if duration >= localTask.GetTaskData().GetOperationCost().GetDuration() {
					log.Infof("Has task running expire, taskId: {%s}, partyId: {%s}, current running duration: {%d ms}, need running duration: {%d ms}",
						taskId, partyId, duration, localTask.GetTaskData().GetOperationCost().GetDuration())

					// 1、 store task expired (failed) event with current party
					m.storeTaskFinalEvent(task.GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(),
						task.GetLocalTaskOrganization().GetPartyId(), fmt.Sprintf("task running expire"),
						apicommonpb.TaskState_TaskState_Failed)
					switch task.GetLocalTaskRole() {
					case apicommonpb.TaskRole_TaskRole_Sender:
						m.publishFinishedTaskToDataCenter(task, localTask, true)
					default:
						// 2、terminate fighter processor for this task with current party
						m.driveTaskForTerminate(task)
						m.sendTaskResultMsgToTaskSender(task)
					}

					// clean current party task cache
					delete(cache, partyId)
					go m.resourceMng.GetDB().RemoveNeedExecuteTaskByPartyId(taskId, partyId)
					if len(cache) == 0 {
						delete(m.runningTaskCache, taskId)
					} else {
						m.runningTaskCache[taskId] = cache
					}

				}
			}
		}
	}
}

func (m *Manager) storeTaskFinalEvent(taskId, identityId, partyId, extra string, state apicommonpb.TaskState) {
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
	m.resourceMng.GetDB().StoreTaskEvent(m.eventEngine.GenerateEvent(evTyp, taskId, identityId, partyId, evMsg))
}

func (m *Manager) storeMetaUsedTaskId(task *types.Task) error {
	identityId, err := m.resourceMng.GetDB().QueryIdentityId()
	if nil != err {
		return err
	}
	for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
		if dataSupplier.GetOrganization().GetIdentityId() == identityId {
			if err := m.resourceMng.GetDB().StoreMetadataHistoryTaskId(dataSupplier.GetMetadataId(), task.GetTaskId()); nil != err {
				return err
			}
		}
	}
	return nil
}

func (m *Manager) checkTaskSenderPublishOpportunity(task *types.Task, event *libtypes.TaskEvent) (bool, error) {

	if "" == strings.Trim(event.GetPartyId(), "") {
		log.Errorf("Failed to check partyId of event, partyId is empty on taskManager.checkTaskSenderPublishOpportunity(), event: %s", event.String())
		return false, fmt.Errorf("empty partyId of event")
	}

	identityId, err := m.resourceMng.GetDB().QueryIdentityId()
	if nil != err {
		log.WithError(err).Errorf("Failed to query self identityId on taskManager.checkTaskSenderPublishOpportunity()")
		return false, fmt.Errorf("query local identityId failed, %s", err)
	}

	// Collect the task result things of other organization, When the current organization is the task sender
	if task.GetTaskSender().GetIdentityId() != identityId {
		return false, nil
	}

	// Remove the currently processed partyId from the partyIds array of the task partner to be processed
	if err := m.resourceMng.GetDB().RemoveTaskPartnerPartyId(event.GetTaskId(), event.GetPartyId()); nil != err {
		log.WithError(err).Errorf("Failed to remove partyId of local task's partner arr on `taskManager.checkTaskSenderPublishOpportunity()`, taskId: {%s}, partyId: {%s}",
			event.GetTaskId(), event.GetPartyId())
	}
	// has, err := m.resourceMng.GetDB().HasTaskPartnerPartyIds(event.GetTaskId())
	partyIds, err := m.resourceMng.GetDB().QueryTaskPartnerPartyIds(event.GetTaskId())
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed to query task partner partyIds on `taskManager.checkTaskSenderPublishOpportunity()`, taskId: {%s}",
			event.GetTaskId())
		return false, err
	}

	log.Debugf("Query partyIds on `taskManager.checkTaskSenderPublishOpportunity()`, current partyId: {%s}, partyIds: %s",
		event.GetPartyId(), "["+strings.Join(partyIds, ",")+"]")

	if rawdb.IsDBNotFoundErr(err) {
		return true, nil
	}

	// if has {
	if len(partyIds) != 0 {
		return false, nil
	}
	return true, nil
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

	if len(msg.GetTaskEventList()) == 0 {
		return nil
	}

	taskId := msg.GetTaskEventList()[0].GetTaskId()

	// While task is consensus or executing, handle task resultMsg.
	has, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusByPartyId(taskId, msg.GetMsgOption().GetReceiverPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task executing status on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("query local task executing status failed")
	}

	if !has {
		log.Warnf("Warning not found local task executing status on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return nil
	}

	task, err := m.resourceMng.GetDB().QueryLocalTask(taskId)
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryLocalTask()` on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("query local task failed, %s", err)
	}

	/**
	+++++++++++++++++++++++++++++++++++++++++++
	+++++++++++++++++++++++++++++++++++++++++++
	NOTE: receiverPartyId must be task sender partyId.
	+++++++++++++++++++++++++++++++++++++++++++
	+++++++++++++++++++++++++++++++++++++++++++
	*/
	if msg.GetMsgOption().GetReceiverPartyId() != task.GetTaskSender().GetPartyId() {
		log.Errorf("Failed to check receiver partyId of msg must be task sender partyId, but it is not, on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}, taskSenderPartyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId(), task.GetTaskSender().GetPartyId())
		return fmt.Errorf("invalid taskResultMsg")
	}

	receiver := fetchOrgByPartyRole(msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetReceiverRole(), task)
	identity, err := m.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryIdentity()` on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("query local identity failed, %s", err)
	}
	// verify the receiver is myself ?
	if identity.GetIdentityId() != receiver.GetIdentityId() {
		log.Errorf("Failed to verify receiver identityId of taskResultMsg, receiver is not me on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("receiver is not me of taskResultMsg")
	}

	for _, event := range msg.GetTaskEventList() {

		if "" == strings.Trim(event.GetPartyId(), "") || msg.GetMsgOption().GetSenderPartyId() != strings.Trim(event.GetPartyId(), "") {
			continue
		}

		if err := m.resourceMng.GetDB().StoreTaskEvent(event); nil != err {
			log.WithError(err).Errorf("Failed to store local task event from remote peer on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}, event: %s",
				msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId(), event.String())
		}

		switch event.Type {
		case ev.TaskExecuteSucceedEOF.Type:
			log.Infof("Received task result msg `event is the task final succeed EOF finished`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}, event: %s",
				msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId(), event.String())

			publish, err := m.checkTaskSenderPublishOpportunity(task, event)
			if nil != err {
				log.WithError(err).Errorf("Failed to check task sender publish opportunity on `taskManager.OnTaskResultMsg()`, event: %s",
					event.GetPartyId())
				return err
			}

			if publish {
				senderNeedTask, ok := m.queryNeedExecuteTaskCache(event.GetTaskId(), msg.GetMsgOption().GetReceiverPartyId())
				if ok {
					log.Debugf("Need to call `publishFinishedTaskToDataCenter` on `taskManager.OnTaskResultMsg()`, taskId: {%s}, sender partyId: {%s}",
						event.GetTaskId(), msg.GetMsgOption().GetReceiverPartyId())
					// handle this task result with current peer
					m.publishFinishedTaskToDataCenter(senderNeedTask, task, true)
					m.removeNeedExecuteTaskCache(event.GetTaskId(), msg.GetMsgOption().GetReceiverPartyId())
				}
			}
		case ev.TaskExecuteFailedEOF.Type:

			log.Infof("Received task result msg `event is the task final [FAILED] EOF finished` will terminate task, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}, event: %s",
				msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId(), event.String())

			//
			if err := m.resourceMng.GetDB().RemoveTaskPartnerPartyIds(event.GetTaskId()); nil != err {
				log.WithError(err).Errorf("Failed to remove all partyId of local task's partner arr on `taskManager.OnTaskResultMsg()`, taskId: {%s}", event.GetTaskId())
			}

			// store final failed EOF event
			m.resourceMng.GetDB().StoreTaskEvent(event)
			if err := m.onTerminateExecuteTask(task); nil != err {
				log.Errorf("Failed to call `onTerminateExecuteTask()` on `taskManager.OnTaskResultMsg()`, taskId: {%s}, err: \n%s", task.GetTaskId(), err)
			}
		}
	}
	return nil
}

func (m *Manager) ValidateTaskResourceUsageMsg(pid peer.ID, taskResourceUsageMsg *taskmngpb.TaskResourceUsageMsg) error {
	return nil
}

func (m *Manager) OnTaskResourceUsageMsg(pid peer.ID, usageMsg *taskmngpb.TaskResourceUsageMsg) error {
	return m.onTaskResourceUsageMsg(pid, usageMsg, types.RemoteNetworkMsg)
}

func (m *Manager) onTaskResourceUsageMsg(pid peer.ID, usageMsg *taskmngpb.TaskResourceUsageMsg, nmls types.NetworkMsgLocationSymbol) error {

	msg := types.FetchTaskResourceUsageMsg(usageMsg)

	log.Debugf("Received taskResourceUsageMsg, consensusSymbol: {%s}, remote pid: {%s}, taskResourceUsageMsg: %s", nmls.String(), pid, msg.String())

	has, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusValExecByPartyId(msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task executing status on `taskManager.OnTaskResourceUsageMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("query local task executing status failed")
	}

	if !has {
		log.Warnf("Warning not found local task executing status on `taskManager.OnTaskResourceUsageMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return nil
	}

	// Update task resourceUsed of powerSuppliers of local task
	task, err := m.resourceMng.GetDB().QueryLocalTask(msg.GetUsage().GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryLocalTask()` on `taskManager.OnTaskResourceUsageMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("query local task failed, %s", err)
	}
	receiver := fetchOrgByPartyRole(msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetReceiverRole(), task)
	identity, err := m.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryIdentity()` on `taskManager.OnTaskResourceUsageMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("query local identity failed, %s", err)
	}
	// verify the receiver is myself ?
	if identity.GetIdentityId() != receiver.GetIdentityId() {

		log.Errorf("Failed to verify receiver identityId of taskResourceUsageMsg, receiver is not me on `taskManager.OnTaskResourceUsageMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("receiver is not me of taskResourceUsageMsg")
	}

	var needUpdate bool

	for i, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {

		// find power supplier info by identity and partyId with msg from reomte peer
		// (find the target power supplier, it maybe local power supplier or remote power supplier)
		// and update its' resource usage info.
		if msg.GetMsgOption().GetSenderPartyId() == powerSupplier.GetOrganization().GetPartyId() &&
			msg.GetMsgOption().GetOwner().GetIdentityId() == powerSupplier.GetOrganization().GetIdentityId() {

			resourceUsage := task.GetTaskData().GetPowerSuppliers()[i].GetResourceUsedOverview()
			// update ...
			if msg.GetUsage().GetUsedMem() > resourceUsage.GetUsedMem() {
				if msg.GetUsage().GetUsedMem() > task.GetTaskData().GetOperationCost().GetMemory() {
					resourceUsage.UsedMem = task.GetTaskData().GetOperationCost().GetMemory()
				} else {
					resourceUsage.UsedMem = msg.GetUsage().GetUsedMem()
				}
				needUpdate = true
			}
			if msg.GetUsage().GetUsedProcessor() > resourceUsage.GetUsedProcessor() {
				if msg.GetUsage().GetUsedProcessor() > task.GetTaskData().GetOperationCost().GetProcessor() {
					resourceUsage.UsedProcessor = task.GetTaskData().GetOperationCost().GetProcessor()
				} else {
					resourceUsage.UsedProcessor = msg.GetUsage().GetUsedProcessor()
				}
				needUpdate = true
			}
			if msg.GetUsage().GetUsedBandwidth() > resourceUsage.GetUsedBandwidth() {
				if msg.GetUsage().GetUsedBandwidth() > task.GetTaskData().GetOperationCost().GetBandwidth() {
					resourceUsage.UsedBandwidth = task.GetTaskData().GetOperationCost().GetBandwidth()
				} else {
					resourceUsage.UsedBandwidth = msg.GetUsage().GetUsedBandwidth()
				}
				needUpdate = true
			}
			if msg.GetUsage().GetUsedDisk() > resourceUsage.GetUsedDisk() {
				resourceUsage.UsedDisk = msg.GetUsage().GetUsedDisk()
				needUpdate = true
			}
			// update ...
			task.GetTaskData().GetPowerSuppliers()[i].ResourceUsedOverview = resourceUsage
		}
	}

	if needUpdate {

		log.Debugf("Need to update local task on `taskManager.OnTaskResourceUsageMsg()`, usage: %s", msg.GetUsage().String())

		// Updata task when resourceUsed change.
		if err = m.resourceMng.GetDB().StoreLocalTask(task); nil != err {
			log.WithError(err).Errorf("Failed to store local task info on `taskManager.OnTaskResourceUsageMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
				msg.GetMsgOption().GetProposalId().String(), msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
			return fmt.Errorf("update local task by usage change failed")
		}
	}

	return nil
}

func (m *Manager) ValidateTaskTerminateMsg(pid peer.ID, terminateMsg *taskmngpb.TaskTerminateMsg) error {
	return nil
}

func (m *Manager) OnTaskTerminateMsg(pid peer.ID, terminateMsg *taskmngpb.TaskTerminateMsg) error {
	return m.onTaskTerminateMsg(pid, terminateMsg, types.RemoteNetworkMsg)
}

func (m *Manager) onTaskTerminateMsg(pid peer.ID, terminateMsg *taskmngpb.TaskTerminateMsg, nmls types.NetworkMsgLocationSymbol) error {
	msg := types.FetchTaskTerminateTaskMngMsg(terminateMsg)
	log.Debugf("Received taskTerminateMsg, consensusSymbol: {%s}, remote pid: {%s}, taskTerminateMsg: %s", nmls.String(), pid, msg.String())

	task, err := m.resourceMng.GetDB().QueryLocalTask(msg.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryLocalTask()` on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("query local task failed, %s", err)
	}
	if nil == task {
		log.Errorf("Not found local task on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return err
	}

	receiver := fetchOrgByPartyRole(msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetReceiverRole(), task)
	identity, err := m.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryIdentity()` on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("query local identity failed, %s", err)
	}
	// verify the receiver is myself ?
	if identity.GetIdentityId() != receiver.GetIdentityId() {
		log.Errorf("Failed to verify receiver identityId of taskResultMsg, receiver is not me on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("receiver is not me of taskResultMsg")
	}

	needExecuteTask, ok := m.queryNeedExecuteTaskCache(task.GetTaskId(), msg.GetMsgOption().GetReceiverPartyId())
	if !ok {
		log.Warnf("Warning query needExecuteTask failed, task not find on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("needExecuteTask not find")
	}

	// While task is consensus or executing, can terminate.
	has, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusValConsByPartyId(task.GetTaskId(), msg.GetMsgOption().GetReceiverPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task execute `cons` status on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return err
	}

	// interrupt consensus with sender AND send terminateMsg to remote partners
	// OR terminate executing task AND send terminateMsg to remote partners
	if has {
		if err = m.consensusEngine.OnConsensusMsg(pid, types.NewInterruptMsgWrap(task.GetTaskId(), terminateMsg.GetMsgOption())); nil != err {
			log.WithError(err).Errorf("Failed to call `OnConsensusMsg()` on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
				msg.GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
			return err
		}
	} else {
		has, err = m.resourceMng.GetDB().HasLocalTaskExecuteStatusValExecByPartyId(task.GetTaskId(), msg.GetMsgOption().GetReceiverPartyId())
		if rawdb.IsNoDBNotFoundErr(err) {
			log.WithError(err).Errorf("Failed to query local task execute `exec` status on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
				msg.GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
			return err
		}

		if has {
			// 1、terminate fighter processor for this task with current party
			if err := m.driveTaskForTerminate(needExecuteTask); nil != err {
				log.WithError(err).Errorf("Failed to call driveTaskForTerminate() on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
					msg.GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
				return err
			}

			// 2、 store task terminate (failed or succeed) event with current party
			m.resourceMng.GetDB().StoreTaskEvent(&libtypes.TaskEvent{
				Type:       ev.TaskTerminated.GetType(),
				TaskId:     task.GetTaskId(),
				IdentityId: receiver.GetIdentityId(),
				PartyId:    receiver.GetPartyId(),
				Content:    "task was terminated.",
				CreateAt:   timeutils.UnixMsecUint64(),
			})

			// 3、 remove needExecuteTask cache with current party
			m.removeNeedExecuteTaskCache(task.GetTaskId(), msg.GetMsgOption().GetReceiverPartyId())
			// 4、 send a new needExecuteTask(status: types.TaskTerminate) for terminate with current party
			m.sendNeedExecuteTaskByAction(task.GetTaskId(),
				msg.GetMsgOption().GetReceiverRole(), msg.GetMsgOption().GetSenderRole(),
				receiver, task.GetTaskSender(),
				types.TaskTerminate)
		}
	}
	return nil
}

func fetchOrgByPartyRole(partyId string, role apicommonpb.TaskRole, task *types.Task) *apicommonpb.TaskOrganization {

	switch role {
	case apicommonpb.TaskRole_TaskRole_Sender:
		if partyId == task.GetTaskSender().GetPartyId() {
			return task.GetTaskSender()
		}
	case apicommonpb.TaskRole_TaskRole_DataSupplier:
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			if partyId == dataSupplier.GetOrganization().GetPartyId() {
				return dataSupplier.GetOrganization()
			}
		}
	case apicommonpb.TaskRole_TaskRole_PowerSupplier:
		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			if partyId == powerSupplier.GetOrganization().GetPartyId() {
				return powerSupplier.GetOrganization()
			}
		}
	case apicommonpb.TaskRole_TaskRole_Receiver:
		for _, receiver := range task.GetTaskData().GetReceivers() {
			if partyId == receiver.GetPartyId() {
				return receiver
			}
		}
	}
	return nil
}
