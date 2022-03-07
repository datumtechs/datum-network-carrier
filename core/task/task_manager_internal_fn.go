package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/runutil"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/common/traceutil"
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
	"time"
)

func (m *Manager) tryScheduleTask() error {

	nonConsTask, taskId, err := m.scheduler.TrySchedule()
	if nil == err && nil == nonConsTask {
		return nil
	} else if nil != err && err == schedule.ErrAbandonTaskWithNotFoundTask {
		m.scheduler.RemoveTask(taskId)
		log.WithError(err).Errorf("Failed to call scheduler.TrySchedule(), then remove bullet task, taskId: {%s}", taskId)
	} else if nil != err && err == schedule.ErrAbandonTaskWithNotFoundPowerPartyIds {
		m.scheduler.RemoveTask(taskId)
		log.WithError(err).Errorf("Failed to call scheduler.TrySchedule(), then remove bullet task, taskId: {%s}", taskId)
		m.sendNeedExecuteTaskByAction(types.NewNeedExecuteTask(
			"",
			apicommonpb.TaskRole_TaskRole_Sender,
			apicommonpb.TaskRole_TaskRole_Sender,
			nonConsTask.GetTask().GetTaskSender(),
			nonConsTask.GetTask().GetTaskSender(),
			nonConsTask.GetTask().GetTaskId(),
			types.TaskScheduleFailed,
			&types.PrepareVoteResource{},   // zero value
			&twopcpb.ConfirmTaskPeerInfo{}, // zero value
			fmt.Errorf("schedule failed: "+schedule.ErrAbandonTaskWithNotFoundPowerPartyIds.Error()),
		))
		return err
	} else if nil != err {
		if nil != nonConsTask {

			m.resourceMng.GetDB().StoreTaskEvent(m.eventEngine.GenerateEvent(ev.TaskScheduleFailed.GetType(),
				nonConsTask.GetTask().GetTaskId(), nonConsTask.GetTask().GetTaskSender().GetIdentityId(),
				nonConsTask.GetTask().GetTaskSender().GetPartyId(), "schedule failed: "+err.Error()))

			if er := m.scheduler.RepushTask(nonConsTask.GetTask()); er == schedule.ErrRescheduleLargeThreshold {
				log.WithError(err).Errorf("Failed to repush local task into queue/starve queue after trySchedule failed %s on `taskManager.tryScheduleTask()`, taskId: {%s}",
					err, nonConsTask.GetTask().GetTaskId())

				m.scheduler.RemoveTask(nonConsTask.GetTask().GetTaskId())
				m.sendNeedExecuteTaskByAction(types.NewNeedExecuteTask(
					"",
					apicommonpb.TaskRole_TaskRole_Sender,
					apicommonpb.TaskRole_TaskRole_Sender,
					nonConsTask.GetTask().GetTaskSender(),
					nonConsTask.GetTask().GetTaskSender(),
					nonConsTask.GetTask().GetTaskId(),
					types.TaskScheduleFailed,
					&types.PrepareVoteResource{},   // zero value
					&twopcpb.ConfirmTaskPeerInfo{}, // zero value
					fmt.Errorf("schedule failed: "+err.Error()+" and "+schedule.ErrRescheduleLargeThreshold.Error()),
				))
			}
		}
		return err
	}

	go func(nonConsTask *types.NeedConsensusTask) {

		log.Debugf("Start `NEED-CONSENSUS` task to 2pc consensus engine on `taskManager.tryScheduleTask()`, taskId: {%s}", nonConsTask.GetTask().GetTaskId())

		if err := m.consensusEngine.OnPrepare(nonConsTask); nil != err {
			log.WithError(err).Errorf("Failed to call `OnPrepare()` of 2pc consensus engine on `taskManager.tryScheduleTask()`, taskId: {%s}", nonConsTask.GetTask().GetTaskId())
			// re push task into queue ,if anything else
			if err := m.scheduler.RepushTask(nonConsTask.GetTask()); err == schedule.ErrRescheduleLargeThreshold {
				log.WithError(err).Errorf("Failed to repush local task into queue/starve queue after call `consensus.onPrepare()` on `taskManager.tryScheduleTask()`, taskId: {%s}",
					nonConsTask.GetTask().GetTaskId())

				m.scheduler.RemoveTask(nonConsTask.GetTask().GetTaskId())
				m.sendNeedExecuteTaskByAction(types.NewNeedExecuteTask(
					"",
					apicommonpb.TaskRole_TaskRole_Sender,
					apicommonpb.TaskRole_TaskRole_Sender,
					nonConsTask.GetTask().GetTaskSender(),
					nonConsTask.GetTask().GetTaskSender(),
					nonConsTask.GetTask().GetTaskId(),
					types.TaskScheduleFailed,
					&types.PrepareVoteResource{},   // zero value
					&twopcpb.ConfirmTaskPeerInfo{}, // zero value
					fmt.Errorf("consensus onPrepare failed: "+err.Error()+" and "+schedule.ErrRescheduleLargeThreshold.Error()),
				))
			} else {
				log.Debugf("Succeed to repush local task into queue/starve queue after call `consensus.onPrepare()` on `taskManager.tryScheduleTask()`, taskId: {%s}",
					nonConsTask.GetTask().GetTaskId())
			}
			return
		}
		if err := m.consensusEngine.OnHandle(nonConsTask); nil != err {
			log.WithError(err).Errorf("Failed to call `OnHandle()` of 2pc consensus engine on `taskManager.tryScheduleTask()`, taskId: {%s}", nonConsTask.GetTask().GetTaskId())
		}
	}(nonConsTask)
	return nil
}

func (m *Manager) sendNeedExecuteTaskByAction(task *types.NeedExecuteTask) {
	go func(task *types.NeedExecuteTask) { // asynchronous transmission to reduce Chan blocking
		m.needExecuteTaskCh <- task
	}(task)
}

// To execute task
func (m *Manager) driveTaskForExecute(task *types.NeedExecuteTask) error {

	//// todo mock
	//// ######
	//log.Debugf("Satrt mock flow, taskId: {%s}, partyId: {%s}", task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())
	//localTask, _ := m.resourceMng.GetDB().QueryLocalTask(task.GetTaskId())
	//m.handleResourceUsage(task, task.GetLocalTaskOrganization().GetIdentityId(),
	//	types.NewTaskResuorceUsage(
	//		task.GetTaskId(),
	//		task.GetLocalTaskOrganization().GetPartyId(),
	//		0, 0, 0,
	//		1147483648, 30000000, 100,
	//		0, 2),
	//		localTask)
	////m.RemoveExecuteTaskStateAfterExecuteTask()
	//m.removeNeedExecuteTaskCache(task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())
	//return nil
	////  ######

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
	client, has := m.resourceMng.QueryDataNodeClient(dataNodeId)
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
	client, has := m.resourceMng.QueryJobNodeClient(jobNodeId)
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
	client, has := m.resourceMng.QueryDataNodeClient(dataNodeId)
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
	client, has := m.resourceMng.QueryJobNodeClient(jobNodeId)
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
			if event.GetType() == ev.TaskFailed.GetType() {
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
		}
		if err := m.RemoveExecuteTaskStateAfterExecuteTask("on taskManager.publishFinishedTaskToDataCenter()", task.GetTaskId(),
			task.GetLocalTaskOrganization().GetPartyId(), resource.SetAllReleaseResourceOption(), true); nil != err {
			log.WithError(err).Errorf("Failed to call RemoveExecuteTaskStateAfterExecuteTask() on publishFinishedTaskToDataCenter(), taskId: {%s},  partyId: {%s}, remote pid: {%s}",
				task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())
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

	log.Debugf("Start sendTaskResultMsgToTaskSender, taskId: {%s}, partyId: {%s}, remote pid: {%s}",
		task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())

	var option resource.ReleaseResourceOption

	// when other task partner and task sender is same identity,
	// we don't need to removed local task and local eventList
	if task.GetLocalTaskOrganization().GetIdentityId() == task.GetRemoteTaskOrganization().GetIdentityId() {
		option = resource.SetUnlockLocalResorce() // unlock local resource of partyId, but don't remove local task and events of partyId
	} else {
		option = resource.SetAllReleaseResourceOption() // unlock local resource and remove local task and events
		// broadcast `task result msg` to reply remote peer
		taskResultMsg := m.makeTaskResultMsgWithEventList(task)
		if nil != taskResultMsg {
			if err := m.p2p.Broadcast(context.TODO(), taskResultMsg); nil != err {
				log.WithError(err).Errorf("failed to call `SendTaskResultMsg` on sendTaskResultMsgToTaskSender(), taskId: {%s}, partyId: {%s}, remote pid: {%s}",
					task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())
			} else {
				log.WithField("traceId", traceutil.GenerateTraceID(taskResultMsg)).Debugf("Succeed broadcast taskResultMsg to taskSender on sendTaskResultMsgToTaskSender(), taskId: {%s}, partyId: {%s}, remote pid: {%s}",
					task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())
			}
		}
	}
	if err := m.RemoveExecuteTaskStateAfterExecuteTask("on sendTaskResultMsgToTaskSender()", task.GetTaskId(),
		task.GetLocalTaskOrganization().GetPartyId(), option, false); nil != err {
		log.WithError(err).Errorf("Failed to call RemoveExecuteTaskStateAfterExecuteTask() on sendTaskResultMsgToTaskSender(), taskId: {%s}, partyId: {%s}, remote pid: {%s}",
			task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())
		return
	}
	log.Debugf("Finished sendTaskResultMsgToTaskSender, taskId: {%s}, partyId: {%s}, remote pid: {%s}",
		task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())
}

func (m *Manager) StoreExecuteTaskStateBeforeExecuteTask(logdesc, taskId, partyId string) error {
	// Store task exec status
	if err := m.resourceMng.GetDB().StoreLocalTaskExecuteStatusValExecByPartyId(taskId, partyId); nil != err {
		log.WithError(err).Errorf("Failed to store local task about `running` status %s, taskId: {%s}, partyId: {%s}",
			logdesc, taskId, partyId)
		return err
	}
	log.Debugf("Succeed store local task about `running` status %s, taskId: {%s}, partyId: {%s}",
		logdesc, taskId, partyId)
	// do anythings else?
	return nil
}

func (m *Manager) RemoveExecuteTaskStateAfterExecuteTask(logdesc, taskId, partyId string, option resource.ReleaseResourceOption, isSender bool) error {
	if err := m.resourceMng.GetDB().RemoveLocalTaskExecuteStatusByPartyId(taskId, partyId); nil != err {
		log.WithError(err).Errorf("Failed to remove task executing status %s, taskId: {%s} partyId: {%s}, isSender: {%v}",
			logdesc, taskId, partyId, isSender)
		return err
	}

	log.Debugf("Succeed remove task executing status %s, taskId: {%s} partyId: {%s}, isSender: {%v}",
		logdesc, taskId, partyId, isSender)

	// clean local task cache
	m.resourceMng.ReleaseLocalResourceWithTask(logdesc, taskId, partyId, option, isSender)
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

		log.WithField("traceId", traceutil.GenerateTraceID(terminateMsg)).Debugf("Succeed to call`sendTaskTerminateMsg.%s` taskId: %s, other peer's taskRole: %s, other peer's partyId: %s, other identityId: %s, pid: %s",
			logdesc, task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid)

	}

	size := (len(task.GetTaskData().GetDataSuppliers())) + len(task.GetTaskData().GetPowerSuppliers()) + len(task.GetTaskData().GetReceivers())
	errCh := make(chan error, size)
	var wg sync.WaitGroup

	for i := 0; i < len(task.GetTaskData().GetDataSuppliers()); i++ {

		wg.Add(1)
		dataSupplier := task.GetTaskData().GetDataSuppliers()[i]
		receiver := dataSupplier
		go sendTerminateMsgFn(&wg, sender, receiver, apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_DataSupplier, errCh)

	}
	for i := 0; i < len(task.GetTaskData().GetPowerSuppliers()); i++ {

		wg.Add(1)
		powerSupplier := task.GetTaskData().GetPowerSuppliers()[i]
		receiver := powerSupplier
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
	log.Debugf("Succeed send to manager.loop() tasks count on taskManager.SendTaskMsgArr(), task arr len: %d", len(tasks))
}
func (m *Manager) sendTaskEvent(event *libtypes.TaskEvent) {
	go func(event *libtypes.TaskEvent) {
		m.eventCh <- event
		log.Debugf("Succeed send to manager.loop() a task event on taskManager.sendTaskEvent(), event: %s", event.String())
	}(event)
}

func (m *Manager) storeBadTask(task *types.Task, events []*libtypes.TaskEvent, reason string) error {
	task.GetTaskData().TaskEvents = events
	task.GetTaskData().State = apicommonpb.TaskState_TaskState_Failed
	task.GetTaskData().Reason = reason
	task.GetTaskData().EndAt = timeutils.UnixMsecUint64()

	m.resourceMng.GetDB().RemoveLocalTask(task.GetTaskId())
	m.resourceMng.GetDB().RemoveTaskPowerPartyIds(task.GetTaskId())
	m.resourceMng.GetDB().RemoveTaskEventList(task.GetTaskId())

	return m.resourceMng.GetDB().InsertTask(task)
}

func (m *Manager) convertScheduleTaskToTask(task *types.Task, eventList []*libtypes.TaskEvent, state apicommonpb.TaskState) *types.Task {
	task.GetTaskData().TaskEvents = eventList
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
		ContractId: localTask.GetTaskData().GetAlgorithmCode(),
		//DataId: "",
		PartyId: task.GetLocalTaskOrganization().GetPartyId(),
		//EnvId: "",
		Peers:            peerList,
		ContractCfg:      contractExtraParams,
		DataParty:        dataPartyArr,
		ComputationParty: powerPartyArr,
		ResultParty:      receiverPartyArr,
		Duration:         localTask.GetTaskData().GetOperationCost().GetDuration(),
		Memory:           localTask.GetTaskData().GetOperationCost().GetMemory(),
		Processor:        localTask.GetTaskData().GetOperationCost().GetProcessor(),
		Bandwidth:        localTask.GetTaskData().GetOperationCost().GetBandwidth(),
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
			if partyId == dataSupplier.GetPartyId() {

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
					} else {
						log.Debugf("Succeed consume metadataAuth, taskId: {%s}, metadataAuthId: {%s}", task.GetTaskId(), metadataAuthId)
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

func (m *Manager) addNeedExecuteTaskCache(task *types.NeedExecuteTask, when int64) {
	m.runningTaskCacheLock.Lock()

	taskId, partyId := task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId()

	cache, ok := m.runningTaskCache[taskId]
	if !ok {
		cache = make(map[string]*types.NeedExecuteTask, 0)
	}
	cache[partyId] = task
	m.runningTaskCache[taskId] = cache
	if err := m.resourceMng.GetDB().StoreNeedExecuteTask(task); nil != err {
		log.WithError(err).Errorf("store needExecuteTask failed, taskId: {%s}, partyId: {%s}", taskId, partyId)
	}
	// v0.3.0 add NeedExecuteTask Expire Monitor
	m.addmonitor(task, when)

	log.Debugf("Succeed call addNeedExecuteTaskCache, taskId: {%s}, partyId: {%s}", taskId, partyId)
	m.runningTaskCacheLock.Unlock()
}

func (m *Manager) addmonitor(task *types.NeedExecuteTask, when int64) {

	taskId, partyId := task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId()

	m.syncExecuteTaskMonitors.AddMonitor(types.NewExecuteTaskMonitor(taskId, partyId, when, func() {
		m.runningTaskCacheLock.Lock()
		defer m.runningTaskCacheLock.Unlock()

		cache, ok := m.runningTaskCache[taskId]
		if !ok {
			return
		}

		// 1、check local task from taskId
		localTask, err := m.resourceMng.GetDB().QueryLocalTask(taskId)
		if nil != err {
			for pid, _ := range cache {
				log.WithError(err).Warnf("Can not query local task info, clean current party task cache short circuit AND skip it, on `taskManager.expireTaskMonitor()`, taskId: {%s}, partyId: {%s}",
					taskId, pid)
				// clean current party task cache short circuit.
				delete(cache, pid)
				go m.resourceMng.GetDB().RemoveNeedExecuteTaskByPartyId(taskId, pid)
				if len(cache) == 0 {
					delete(m.runningTaskCache, taskId)
				} else {
					m.runningTaskCache[taskId] = cache
				}
				log.Debugf("Call expireTaskMonitor remove NeedExecuteTask as query local task info failed when task was expired, taskId: {%s}, partyId: {%s}", taskId, partyId)
				continue
			}
			return
		}

		// 2、 check partyId from cache
		if _, ok := cache[partyId]; !ok {
			return
		}

		// 3、 handle ExpireTask
		if localTask.GetTaskData().GetState() == apicommonpb.TaskState_TaskState_Running && localTask.GetTaskData().GetStartAt() != 0 {
			var duration uint64

			duration = timeutils.UnixMsecUint64() - localTask.GetTaskData().GetStartAt()

			log.Infof("Has task running expire, taskId: {%s}, partyId: {%s}, current running duration: {%d ms}, need running duration: {%d ms}",
				taskId, partyId, duration, localTask.GetTaskData().GetOperationCost().GetDuration())

			// 1、 store task expired (failed) event with current party
			m.resourceMng.GetDB().StoreTaskEvent(m.eventEngine.GenerateEvent(ev.TaskFailed.GetType(), taskId,
				task.GetLocalTaskOrganization().GetIdentityId(), partyId,
				fmt.Sprintf("task running expire")))

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
			log.Debugf("Call expireTaskMonitor remove NeedExecuteTask when task was expired, taskId: {%s}, partyId: {%s}", taskId, partyId)
		}
	}))
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
	// v3.0  remove executeTask monitor
	m.syncExecuteTaskMonitors.DelMonitor(taskId, partyId)

	log.Debugf("Call removeNeedExecuteTaskCache, taskId: {%s}, partyId: {%s}",
		taskId, partyId)
}

func (m *Manager) hasNeedExecuteTaskCache(taskId, partyId string) bool {
	m.runningTaskCacheLock.RLock()
	defer m.runningTaskCacheLock.RUnlock()
	cache, ok := m.runningTaskCache[taskId]
	if !ok {
		log.Debugf("Call queryNeedExecuteTaskCache, taskId: {%s}, partyId: {%s}, has: {%v}",
			taskId, partyId, ok)
		return false
	}
	_, ok = cache[partyId]
	log.Debugf("Call queryNeedExecuteTaskCache, taskId: {%s}, partyId: {%s}, has: {%v}",
		taskId, partyId, ok)
	return ok
}

func (m *Manager) queryNeedExecuteTaskCache(taskId, partyId string) (*types.NeedExecuteTask, bool) {
	m.runningTaskCacheLock.RLock()
	defer m.runningTaskCacheLock.RUnlock()
	cache, ok := m.runningTaskCache[taskId]
	if !ok {
		log.Debugf("Call queryNeedExecuteTaskCache, taskId: {%s}, partyId: {%s}, has: {%v}",
			taskId, partyId, ok)
		return nil, false
	}
	task, ok := cache[partyId]
	log.Debugf("Call queryNeedExecuteTaskCache, taskId: {%s}, partyId: {%s}, has: {%v}",
		taskId, partyId, ok)
	return task, ok
}

func (m *Manager) mustQueryNeedExecuteTaskCache(taskId, partyId string) *types.NeedExecuteTask {
	task, _ := m.queryNeedExecuteTaskCache(taskId, partyId)
	return task
}

func (m *Manager) makeTaskResultMsgWithEventList(task *types.NeedExecuteTask) *taskmngpb.TaskResultMsg {

	if task.GetLocalTaskRole() == apicommonpb.TaskRole_TaskRole_Sender {
		log.Errorf("the task sender can not make taskResultMsg")
		return nil
	}

	eventList, err := m.resourceMng.GetDB().QueryTaskEventListByPartyId(task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to make taskResultMsg with query task eventList, taskId {%s}, partyId {%s}",
			task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())
		return nil
	}

	if len(eventList) == 0 {
		log.Errorf("Failed to make taskResultMsg with query task eventList is empty, taskId {%s}, partyId {%s}",
			task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())
		return nil
	}

	return &taskmngpb.TaskResultMsg{
		MsgOption: &msgcommonpb.MsgOption{
			ProposalId:      common.Hash{}.Bytes(),
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

func (m *Manager) handleTaskEventWithCurrentOranization(event *libtypes.TaskEvent) error {
	if len(event.GetType()) != ev.EventTypeCharLen {
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
			log.WithError(err).Errorf("Failed to query local task info on `taskManager.handleTaskEventWithCurrentOranization()`, taskId: {%s}, partyId: {%s}",
				event.GetTaskId(), event.GetPartyId())
			// remove wrong task cache
			m.removeNeedExecuteTaskCache(event.GetTaskId(), event.GetPartyId())
			return err
		}

		// need to validate the task that have been processing ? Maybe~
		// While task is consensus or executing, can terminate.
		has, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusByPartyId(event.GetTaskId(), event.GetPartyId())
		if nil != err {
			log.WithError(err).Errorf("Failed to check local task execute status whether exist on `taskManager.handleTaskEventWithCurrentOranization()`, taskId: {%s}, partyId: {%s}",
				event.GetTaskId(), event.GetPartyId())
			// remove wrong task cache
			m.removeNeedExecuteTaskCache(event.GetTaskId(), event.GetPartyId())
			return err
		}

		if !has {
			log.Warnf("Warn ignore event, `event is the end` but not find party task executeStatus on `taskManager.handleTaskEventWithCurrentOranization()`, event: %s",
				event.String())
			// remove wrong task cache
			m.removeNeedExecuteTaskCache(event.GetTaskId(), event.GetPartyId())
			return nil
		}

		return m.executeTaskEvent("on `taskManager.handleTaskEventWithCurrentOranization()`", types.LocalNetworkMsg, event, task, localTask)
	}
	return nil // ignore event while task is not exist.
}

func (m *Manager) handleNeedExecuteTask(task *types.NeedExecuteTask, localTask *types.Task) {

	log.Debugf("Start handle needExecuteTask on handleNeedExecuteTask(), taskId: {%s}, role: {%s}, partyId: {%s}",
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
			log.WithError(err).Errorf("Failed to update local task state before executing task on handleNeedExecuteTask(), taskId: {%s}, need update state: {%s}",
				task.GetTaskId(), apicommonpb.TaskState_TaskState_Running.String())
		}
	}

	// store local cache
	m.addNeedExecuteTaskCache(task, int64(localTask.GetTaskData().GetStartAt()+localTask.GetTaskData().GetOperationCost().GetDuration()))

	// The task sender will not execute the task
	if task.GetLocalTaskRole() != apicommonpb.TaskRole_TaskRole_Sender &&
		task.GetLocalTaskOrganization().GetPartyId() != localTask.GetTaskSender().GetPartyId() {
		// driving task to executing
		if err := m.driveTaskForExecute(task); nil != err {
			log.WithError(err).Errorf("Failed to execute task on internal node on handleNeedExecuteTask(), taskId: {%s}, role: {%s}, partyId: {%s}",
				task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId())

			m.SendTaskEvent(m.eventEngine.GenerateEvent(ev.TaskExecuteFailedEOF.GetType(), task.GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(),
				task.GetLocalTaskOrganization().GetPartyId(), fmt.Sprintf("%s, %s with %s", ev.TaskExecuteFailedEOF.GetMsg(), err,
					task.GetLocalTaskOrganization().GetPartyId())))
		}
	}
}

func (m *Manager) executeTaskEvent(logkeyword string, symbol types.NetworkMsgLocationSymbol, event *libtypes.TaskEvent, localNeedtask *types.NeedExecuteTask, localTask *types.Task) error {


	if err := m.resourceMng.GetDB().StoreTaskEvent(event); nil != err {
		log.WithError(err).Errorf("Failed to store %s taskEvent %s, event: %s", symbol.String(), logkeyword, event.String())
	} else {
		log.Infof("Started store %s taskEvent %s, event: %s", symbol.String(), logkeyword, event.String())
	}

	switch event.GetType() {
	case ev.TaskExecuteSucceedEOF.GetType():

		log.Infof("Started handle taskEvent with currentIdentity, `event is the task final succeed EOF finished` %s, event: %s", logkeyword, event.String())

		if symbol == types.LocalNetworkMsg {
			m.resourceMng.GetDB().StoreTaskEvent(m.eventEngine.GenerateEvent(ev.TaskSucceed.GetType(), event.GetTaskId(),
				event.GetIdentityId(), event.GetPartyId(), "task execute succeed"))
		}

		publish, err := m.checkTaskSenderPublishOpportunity(localTask, event)
		if nil != err {
			log.WithError(err).Errorf("Failed to check task sender publish opportunity %s, event: %s",
				logkeyword, event.GetPartyId())
			return err
		}

		if publish {

			// 1、 handle last party
			//   send this task result to remote target peer,
			//   but they belong to same organization, call local msg.
			if symbol == types.LocalNetworkMsg {
				m.sendTaskResultMsgToTaskSender(localNeedtask)
				m.removeNeedExecuteTaskCache(event.GetTaskId(), event.GetPartyId())
			}

			// 2、 handle sender party
			senderNeedTask, ok := m.queryNeedExecuteTaskCache(event.GetTaskId(), localTask.GetTaskSender().GetPartyId())
			if ok {
				log.Debugf("Need to call `publishFinishedTaskToDataCenter` %s, taskId: {%s}, sender partyId: {%s}",
					logkeyword, event.GetTaskId(), localTask.GetTaskSender().GetPartyId())
				// handle this task result with current peer
				m.publishFinishedTaskToDataCenter(senderNeedTask, localTask, true)
				m.removeNeedExecuteTaskCache(event.GetTaskId(), localTask.GetTaskSender().GetPartyId())
			}
		} else {
			if symbol == types.LocalNetworkMsg {
				// send this task result to remote target peer
				m.sendTaskResultMsgToTaskSender(localNeedtask)
				m.removeNeedExecuteTaskCache(event.GetTaskId(), event.GetPartyId())
			}
		}
	case ev.TaskExecuteFailedEOF.GetType():
		log.Infof("Started handle taskEvent with currentIdentity, `event is the task final [FAILED] EOF finished` will terminate task %s, event: %s", logkeyword, event.String())

		//
		if err := m.resourceMng.GetDB().RemoveTaskPartnerPartyIds(event.GetTaskId()); nil != err {
			log.WithError(err).Errorf("Failed to remove all partyId of local task's partner arr %s, taskId: {%s}",
				logkeyword, event.GetTaskId())
		}

		// ## 1、 check wether task status is `terminate`
		terminating, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusTerminateByPartyId(localTask.GetTaskId(), localTask.GetTaskSender().GetPartyId())
		if nil != err {
			log.WithError(err).Errorf("Failed to query local task execute `termining` status with task sender %s, taskId: {%s}, partyId: {%s}",
				logkeyword, localTask.GetTaskId(), localTask.GetTaskSender().GetPartyId())
			return err
		}
		if terminating {
			log.Warnf("Warning query local task execute status has `termining` with task sender %s, taskId: {%s}, partyId: {%s}",
				logkeyword, localTask.GetTaskId(), localTask.GetTaskSender().GetPartyId())
			return nil
		}

		// ## 2、 check wether task status is `running`
		running, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusRunningByPartyId(localTask.GetTaskId(), localTask.GetTaskSender().GetPartyId())
		if nil != err {
			log.WithError(err).Errorf("Failed to query local task execute status has `running` with task sender %s, taskId: {%s}, partyId: {%s}",
				logkeyword, localTask.GetTaskId(), localTask.GetTaskSender().GetPartyId())
			return err
			return err
		}
		// (While task is consensus or running, can terminate.)
		if running {
			if err := m.onTerminateExecuteTask(event.GetTaskId(), event.GetPartyId(), localTask); nil != err {
				log.Errorf("Failed to call `onTerminateExecuteTask()` %s, taskId: {%s}, err: \n%s", logkeyword, localTask.GetTaskId(), err)
			}
		}
	}
	return nil
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

func (m *Manager) handleResourceUsage(keyword, usageIdentityId string, usage *types.TaskResuorceUsage, localTask *types.Task, nmls types.NetworkMsgLocationSymbol) (bool, error) {

	var (
		terminating bool
		running bool
	)

	if nmls == types.LocalNetworkMsg {
		// ## 1、 check whether task status is terminate (with party self) ?
		tflag, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusTerminateByPartyId(usage.GetTaskId(), usage.GetPartyId())
		if nil != err {
			log.WithError(err).Errorf("Failed to call HasLocalTaskExecuteStatusTerminateByPartyId() on taskManager.handleResourceUsage() %s, taskId: {%s}, partyId: {%s}",
				keyword, usage.GetTaskId(), usage.GetPartyId())
			return false, fmt.Errorf("check current party has `terminate` status needExecuteTask failed, %s", err)
		}
		terminating = tflag
	} else {
		// ## 1、 check whether task status is terminate (with task sender)?
		tflag, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusTerminateByPartyId(usage.GetTaskId(), localTask.GetTaskSender().GetPartyId())
		if nil != err {
			log.WithError(err).Errorf("Failed to call HasLocalTaskExecuteStatusTerminateByPartyId() on taskManager.handleResourceUsage() %s, taskId: {%s}, partyId: {%s}, task sender partyId: {%s}",
				keyword, usage.GetTaskId(), usage.GetPartyId(), localTask.GetTaskSender().GetPartyId())
			return false, fmt.Errorf("check task sender has `terminate` status needExecuteTask failed, %s", err)
		}
		terminating = tflag
	}

	if terminating {
		log.Warnf("the localTask execute status has `terminate` on taskManager.handleResourceUsage() %s, taskId: {%s}, partyId: {%s}, task sender partyId: {%s}",
			keyword, usage.GetTaskId(), usage.GetPartyId(), localTask.GetTaskSender().GetPartyId())
		return false, fmt.Errorf("task was terminated")
	}

	if nmls == types.LocalNetworkMsg {
		// ## 2、 check whether task status is running (with current party self)?
		rflag, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusRunningByPartyId(usage.GetTaskId(), usage.GetPartyId())
		if nil != err {
			log.WithError(err).Errorf("Failed to call HasLocalTaskExecuteStatusRunningByPartyId() on taskManager.handleResourceUsage() %s, taskId: {%s}, partyId: {%s}",
				keyword, usage.GetTaskId(), usage.GetPartyId())
			return false, fmt.Errorf("check current party has `running` status needExecuteTask failed, %s", err)
		}
		running = rflag
	} else {
		// ## 2、 check whether task status is running (with task sender)?
		rflag, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusRunningByPartyId(usage.GetTaskId(), localTask.GetTaskSender().GetPartyId())
		if nil != err {
			log.WithError(err).Errorf("Failed to call HasLocalTaskExecuteStatusRunningByPartyId() on taskManager.handleResourceUsage() %s, taskId: {%s}, partyId: {%s}, task sender partyId: {%s}",
				keyword, usage.GetTaskId(), usage.GetPartyId(), localTask.GetTaskSender().GetPartyId())
			return false, fmt.Errorf("check task sender has `running` status needExecuteTask failed, %s", err)
		}
		running = rflag
	}

	if !running {
		log.Warnf("Not found localTask execute status `running` on taskManager.handleResourceUsage() %s, taskId: {%s}, partyId: {%s}, task sender partyId: {%s}",
			keyword, usage.GetTaskId(), usage.GetPartyId(), localTask.GetTaskSender().GetPartyId())
		return false, fmt.Errorf("task is not executed")
	}

	var needUpdate bool

	for i, powerSupplier := range localTask.GetTaskData().GetPowerSuppliers() {

		// find power supplier info by identity and partyId with msg from reomte peer
		// (find the target power supplier, it maybe local power supplier or remote power supplier)
		// and update its' resource usage info.
		if usage.GetPartyId() == powerSupplier.GetOrganization().GetPartyId() &&
			usageIdentityId == powerSupplier.GetOrganization().GetIdentityId() {

			resourceUsage := localTask.GetTaskData().GetPowerSuppliers()[i].GetResourceUsedOverview()
			// update ...
			if usage.GetUsedMem() > resourceUsage.GetUsedMem() {
				if usage.GetUsedMem() > localTask.GetTaskData().GetOperationCost().GetMemory() {
					resourceUsage.UsedMem = localTask.GetTaskData().GetOperationCost().GetMemory()
				} else {
					resourceUsage.UsedMem = usage.GetUsedMem()
				}
				needUpdate = true
			}
			if usage.GetUsedProcessor() > resourceUsage.GetUsedProcessor() {
				if usage.GetUsedProcessor() > localTask.GetTaskData().GetOperationCost().GetProcessor() {
					resourceUsage.UsedProcessor = localTask.GetTaskData().GetOperationCost().GetProcessor()
				} else {
					resourceUsage.UsedProcessor = usage.GetUsedProcessor()
				}
				needUpdate = true
			}
			if usage.GetUsedBandwidth() > resourceUsage.GetUsedBandwidth() {
				if usage.GetUsedBandwidth() > localTask.GetTaskData().GetOperationCost().GetBandwidth() {
					resourceUsage.UsedBandwidth = localTask.GetTaskData().GetOperationCost().GetBandwidth()
				} else {
					resourceUsage.UsedBandwidth = usage.GetUsedBandwidth()
				}
				needUpdate = true
			}
			if usage.GetUsedDisk() > resourceUsage.GetUsedDisk() {
				resourceUsage.UsedDisk = usage.GetUsedDisk()
				needUpdate = true
			}
			// update ...
			localTask.GetTaskData().GetPowerSuppliers()[i].ResourceUsedOverview = resourceUsage
		}
	}

	if needUpdate {
		log.Debugf("Need to update local task on taskManager.handleResourceUsage() %s, usage: %s", keyword, usage.String())

		// Updata task when resourceUsed change.
		if err := m.resourceMng.GetDB().StoreLocalTask(localTask); nil != err {
			log.WithError(err).Errorf("Failed to call StoreLocalTask() on taskManager.handleResourceUsage() %s, taskId: {%s}, partyId: {%s}",
				keyword, usage.GetTaskId(), usage.GetPartyId())
			return false, fmt.Errorf("update local task by usage change failed, %s", err)
		}
	}

	return needUpdate, nil
}


func (m *Manager) ValidateTaskResultMsg(pid peer.ID, taskResultMsg *taskmngpb.TaskResultMsg) error {
	msg := types.FetchTaskResultMsg(taskResultMsg) // fetchTaskResultMsg(taskResultMsg)

	if len(msg.GetTaskEventList()) == 0 {
		return nil
	}

	taskId := msg.GetTaskEventList()[0].GetTaskId()

	for _, event := range msg.GetTaskEventList() {
		if taskId != event.GetTaskId() {
			return fmt.Errorf("Received event failed, has invalid taskId: {%s}, right taskId: {%s}", event.GetTaskId(), taskId)
		}
	}

	return nil
}

func (m *Manager) OnTaskResultMsg(pid peer.ID, taskResultMsg *taskmngpb.TaskResultMsg) error {

	msg := types.FetchTaskResultMsg(taskResultMsg)

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
		log.Warnf("Not found local task executing status on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return nil
	}

	localTask, err := m.resourceMng.GetDB().QueryLocalTask(taskId)
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
	if msg.GetMsgOption().GetReceiverPartyId() != localTask.GetTaskSender().GetPartyId() {
		log.Errorf("Failed to check receiver partyId of msg must be task sender partyId, but it is not, on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}, taskSenderPartyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId(),
			localTask.GetTaskSender().GetPartyId())
		return fmt.Errorf("invalid taskResultMsg")
	}

	receiver := fetchOrgByPartyRole(msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetReceiverRole(), localTask)
	identity, err := m.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryIdentity()` on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("query local identity failed, %s", err)
	}
	// verify the receiver is myself ?
	if identity.GetIdentityId() != receiver.GetIdentityId() {
		log.Warnf("Warning verify receiver identityId of taskResultMsg, receiver is not me on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("receiver is not me of taskResultMsg")
	}

	log.WithField("traceId", traceutil.GenerateTraceID(taskResultMsg)).Debugf("Received remote taskResultMsg, remote pid: {%s}, taskId: {%s}, taskResultMsg: %s", pid, taskId, msg.String())

	for _, event := range msg.GetTaskEventList() {

		if "" == strings.Trim(event.GetPartyId(), "") || msg.GetMsgOption().GetSenderPartyId() != strings.Trim(event.GetPartyId(), "") {
			continue
		}

		if err := m.executeTaskEvent("on `taskManager.OnTaskResultMsg()`", types.RemoteNetworkMsg, event, nil, localTask); nil != err {
			return err
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

	if msg.GetMsgOption().GetSenderPartyId() != msg.GetUsage().GetPartyId() {
		log.Errorf("sender partyId of usageMsg AND partyId of usageMsg is not same when received taskResourceUsageMsg, taskId: {%s}, sender partyId: {%s}, usagemsgPartyId: {%s}",
			msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetSenderPartyId(), msg.GetUsage().GetPartyId())
		return fmt.Errorf("invalid partyId of usageMsg")
	}

	// Note: the needexecutetask obtained here is generally the needexecutetask of the task sender, so the remoteorganization is also the task sender.
	_, ok := m.queryNeedExecuteTaskCache(msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverPartyId())
	if !ok {
		log.Warnf("Not found needExecuteTask when received taskResourceUsageMsg, taskId: {%s}, partyId: {%s}",
			msg.GetUsage().GetTaskId(),  msg.GetUsage().GetPartyId())
		return fmt.Errorf("Can not find `need execute task` cache")
	}

	// Update task resourceUsed of powerSuppliers of local task
	task, err := m.resourceMng.GetDB().QueryLocalTask(msg.GetUsage().GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryLocalTask()` when received taskResourceUsageMsg, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("query local task failed, %s", err)
	}
	receiver := fetchOrgByPartyRole(msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetReceiverRole(), task)
	identity, err := m.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryIdentity()` when received taskResourceUsageMsg, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("query local identity failed, %s", err)
	}
	// verify the receiver is myself ?
	if identity.GetIdentityId() != receiver.GetIdentityId() {
		log.Warnf("Warning verify receiver identityId of taskResourceUsageMsg, receiver is not me when received taskResourceUsageMsg, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("receiver is not me of taskResourceUsageMsg")
	}

	// Check whether the receiver of the message is the same organization as the sender of the task.
	// If not, this message is illegal.
	if task.GetTaskSender().GetIdentityId() != receiver.GetIdentityId() ||
		task.GetTaskSender().GetPartyId() != receiver.GetPartyId() {
		log.Warnf("Warning the receiver of the message is not the same organization as the sender of the task when received taskResourceUsageMsg, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, msg receiver: %s, task sender: %s",
			msg.GetMsgOption().GetProposalId().String(), task.GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), receiver.String(), task.GetTaskSender().String())
		return fmt.Errorf("receiver is not task sender of taskResourceUsageMsg")
	}

	log.WithField("traceId", traceutil.GenerateTraceID(usageMsg)).Debugf("Received taskResourceUsageMsg when received taskResourceUsageMsg, consensusSymbol: {%s}, remote pid: {%s}, taskResourceUsageMsg: %s",
		nmls.String(), pid, msg.String())

	needUpdate, err := m.handleResourceUsage("when received remote resourceUsage", msg.GetMsgOption().GetOwner().GetIdentityId(), msg.GetUsage(), task, types.RemoteNetworkMsg)
	if nil != err {
		return err
	}

	if needUpdate {
		log.Debugf("Succeed handle remote resourceUsage when received taskResourceUsageMsg, consensusSymbol: {%s}, remote pid: {%s}, taskResourceUsageMsg: %s", nmls.String(), pid, msg.String())
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
		log.Warnf("Warning verify receiver identityId of taskTerminateMsg, receiver is not me on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("receiver is not me of taskTerminateMsg")
	}

	log.WithField("traceId", traceutil.GenerateTraceID(terminateMsg)).Debugf("Received taskTerminateMsg, consensusSymbol: {%s}, remote pid: {%s}, taskTerminateMsg: %s", nmls.String(), pid, msg.String())

	// ## 1、 check whether the task has been terminated

	terminating, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusTerminateByPartyId(task.GetTaskId(), msg.GetMsgOption().GetReceiverPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task execute `termining` status on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetTaskSender().GetPartyId())
		return err
	}
	// If so, we will directly short circuit
	if terminating {
		log.Warnf("Warning query local task execute status has `termining` on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetTaskSender().GetPartyId())
		return nil
	}

	// ## 2、 check whether the task is running.
	if needExecuteTask, ok := m.queryNeedExecuteTaskCache(task.GetTaskId(), msg.GetMsgOption().GetReceiverPartyId()); ok {
		return m.startTerminateWithNeedExecuteTask(needExecuteTask)
	}

	// ## 3、 check whether the task is in consensus

	// While task is consensus or executing, can terminate.
	consensusing, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusConsensusByPartyId(task.GetTaskId(), msg.GetMsgOption().GetReceiverPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task execute `cons` status on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return err
	}
	// interrupt consensus with sender AND send terminateMsg to remote partners
	// OR terminate executing task AND send terminateMsg to remote partners
	if consensusing {
		if err = m.consensusEngine.OnConsensusMsg(pid, types.NewInterruptMsgWrap(task.GetTaskId(), terminateMsg.GetMsgOption())); nil != err {
			log.WithError(err).Errorf("Failed to call `OnConsensusMsg()` on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
				msg.GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
			return err
		}
		return nil
	}

	return nil
}


func (m *Manager) startTerminateWithNeedExecuteTask (needExecuteTask *types.NeedExecuteTask) error {

	// ## 1、 check whether the task has been terminated

	terminating, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusTerminateByPartyId(needExecuteTask.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task execute `termining` status on `taskManager.startTerminateWithNeedExecuteTask()`, taskId: {%s}, partyId: {%s}",
			needExecuteTask.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId())
		return err
	}
	// If so, we will directly short circuit
	if terminating {
		log.Warnf("Warning query local task execute status has `termining` on `taskManager.startTerminateWithNeedExecuteTask()`, taskId: {%s}, partyId: {%s}",
			needExecuteTask.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId())
		return nil
	}

	// ## 2、 check whether the task is running.
	running, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusRunningByPartyId(needExecuteTask.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId())
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed to query local task execute `exec` status on `taskManager.startTerminateWithNeedExecuteTask()`, taskId: {%s}, partyId: {%s}",
			needExecuteTask.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId())
		return err
	}
	// If it is, we will terminate the task
	if !running {
		log.Warnf("the local task execute status is not `running` on `taskManager.startTerminateWithNeedExecuteTask()`, taskId: {%s}, partyId: {%s}",
			needExecuteTask.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId())
		return fmt.Errorf("the task is not running")
	}

	// 1、terminate fighter processor for this task with current party
	if err := m.driveTaskForTerminate(needExecuteTask); nil != err {
		log.WithError(err).Errorf("Failed to call driveTaskForTerminate() on `taskManager.startTerminateWithNeedExecuteTask()`, taskId: {%s}, role: {%s}, partyId: {%s}",
			needExecuteTask.GetTaskId(), needExecuteTask.GetLocalTaskRole().String(), needExecuteTask.GetLocalTaskOrganization().GetPartyId())
		return err
	}
	// 2、 store task terminate (failed or succeed) event with current party
	m.resourceMng.GetDB().StoreTaskEvent(&libtypes.TaskEvent{
		Type:       ev.TaskTerminated.GetType(),
		TaskId:     needExecuteTask.GetTaskId(),
		IdentityId: needExecuteTask.GetLocalTaskOrganization().GetIdentityId(),
		PartyId:    needExecuteTask.GetLocalTaskOrganization().GetPartyId(),
		Content:    "task was terminated.",
		CreateAt:   timeutils.UnixMsecUint64(),
	})

	// 3、 remove needExecuteTask cache with current party
	m.removeNeedExecuteTaskCache(needExecuteTask.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId())
	// 4、Set the execution status of the task to being terminated`
	if err := m.resourceMng.GetDB().StoreLocalTaskExecuteStatusValTerminateByPartyId(needExecuteTask.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId()); nil != err {
		log.WithError(err).Errorf("Failed to store needExecute task status to `terminate` on `taskManager.startTerminateWithNeedExecuteTask()`, taskId: {%s}, role: {%s}, partyId: {%s}",
			needExecuteTask.GetTaskId(), needExecuteTask.GetLocalTaskRole().String(), needExecuteTask.GetLocalTaskOrganization().GetPartyId())
	}
	// 5、 send a new needExecuteTask(status: types.TaskTerminate) for terminate with current party
	m.sendNeedExecuteTaskByAction(types.NewNeedExecuteTask(
		"",
		needExecuteTask.GetLocalTaskRole(),
		needExecuteTask.GetRemoteTaskRole(),
		needExecuteTask.GetLocalTaskOrganization(),
		needExecuteTask.GetRemoteTaskOrganization(),
		needExecuteTask.GetTaskId(),
		types.TaskTerminate,
		&types.PrepareVoteResource{},   // zero value
		&twopcpb.ConfirmTaskPeerInfo{}, // zero value
		fmt.Errorf("task was terminated."),
	))
	return nil
}


func (m *Manager) checkNeedExecuteTaskMonitors(now int64) int64 {
	return m.syncExecuteTaskMonitors.CheckMonitors(now)
}
func (m *Manager) needExecuteTaskMonitorTimer() *time.Timer {
	return m.syncExecuteTaskMonitors.Timer()
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
