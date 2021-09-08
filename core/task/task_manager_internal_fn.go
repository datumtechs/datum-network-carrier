package task

import (
	"encoding/json"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	ev "github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/handler"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/lib/fighter/common"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

func (m *Manager) driveTaskForExecute(task *types.NeedExecuteTask) error {

	task.Task.SchedTask.GetTaskData().State = types.TaskStateRunning.String()
	task.Task.SchedTask.GetTaskData().StartAt = uint64(timeutils.UnixMsec())
	if err := m.dataCenter.StoreLocalTask(task.Task.SchedTask); nil != err {
		log.Errorf("Failed to update local task state before executing task, taskId: {%s}, need update state: {%s}, err: {%s}",
			task.Task.SchedTask.GetTaskId(), types.TaskStateRunning.String(), err)
	}
	// update local cache
	m.addRunningTaskCache(task)

	//return fmt.Errorf("Mock task finished")

	switch task.SelfTaskRole {
	case types.TaskOwner, types.DataSupplier, types.ResultSupplier:
		return m.executeTaskOnDataNode(task)
	case types.PowerSupplier:
		return m.executeTaskOnJobNode(task)
	default:
		log.Errorf("Faided to driveTaskForExecute(), Unknown task role, taskId: {%s}, taskRole: {%s}", task.Task.SchedTask.GetTaskId(), task.SelfTaskRole.String())
		return errors.New("Unknown resource node type")
	}
}

func (m *Manager) executeTaskOnDataNode(task *types.DoneScheduleTaskChWrap) error {

	// 找到自己的投票
	dataNodeId := task.Task.SelfVotePeerInfo.Id

	// clinet *grpclient.DataNodeClient,
	client, has := m.resourceClientSet.QueryDataNodeClient(dataNodeId)
	if !has {
		log.Errorf("Failed to query internal data node, taskId: {%s}, dataNodeId: {%s}",
			task.Task.SchedTask.GetTaskId(), dataNodeId)
		return errors.New("data node client not found")
	}
	if client.IsNotConnected() {
		if err := client.Reconnect(); nil != err {
			log.Errorf("Failed to connect internal data node, taskId: {%s}, dataNodeId: {%s}, err: {%s}",
				task.Task.SchedTask.GetTaskId(), dataNodeId, err)
			return err
		}
	}

	req, err := m.makeTaskReadyGoReq(task)
	if nil != err {
		log.Errorf("Falied to make taskReadyGoReq, taskId: {%s}, dataNodeId: {%s}, err: {%s}",
			task.Task.SchedTask.GetTaskId(), dataNodeId, err)
		return err
	}

	resp, err := client.HandleTaskReadyGo(req)
	if nil != err {
		log.Errorf("Falied to call publish schedTask to `data-Fighter` node to executing, taskId: {%s}, dataNodeId: {%s}, err: {%s}",
			task.Task.SchedTask.GetTaskId(), dataNodeId, err)
		return err
	}
	if !resp.Ok {
		log.Errorf("Falied to executing task from `data-Fighter` node response, taskId: {%s}, dataNodeId: {%s}, resp: {%s}",
			task.Task.SchedTask.GetTaskId(), dataNodeId, resp.String())
		return nil
	}


	log.Infof("Success to publish schedTask to `data-Fighter` node to executing,  taskId: {%s}, dataNodeId: {%s}",
		task.Task.SchedTask.GetTaskId(), dataNodeId)
	return nil
}

func (m *Manager) executeTaskOnJobNode(task *types.DoneScheduleTaskChWrap) error {


	// 找到自己的投票
	jobNodeId := task.Task.SelfVotePeerInfo.Id


	// clinet *grpclient.DataNodeClient,
	client, has := m.resourceClientSet.QueryJobNodeClient(jobNodeId)
	if !has {
		log.Errorf("Failed to query internal job node, taskId: {%s}, jobNodeId: {%s}",
			task.Task.SchedTask.GetTaskId(), jobNodeId)
		return errors.New("job node client not found")
	}
	if client.IsNotConnected() {
		if err := client.Reconnect(); nil != err {
			log.Errorf("Failed to connect internal job node, taskId: {%s}, jobNodeId: {%s}, err: {%s}",
				task.Task.SchedTask.GetTaskId(), jobNodeId, err)
			return err
		}
	}

	req, err := m.makeTaskReadyGoReq(task)
	if nil != err {
		log.Errorf("Falied to make taskReadyGoReq, taskId: {%s}, jobNodeId: {%s}, err: {%s}",
			task.Task.SchedTask.GetTaskId(), jobNodeId, err)
		return err
	}

	resp, err := client.HandleTaskReadyGo(req)
	if nil != err {
		log.Errorf("Falied to publish schedTask to `job-Fighter` node to executing, taskId: {%s}, jobNodeId: {%s}, err: {%s}",
			task.Task.SchedTask.GetTaskId(), jobNodeId, err)
		return err
	}
	if !resp.Ok {
		log.Errorf("Falied to publish schedTask to `job-Fighter` node to executing, taskId: {%s}, jobNodeId: {%s}",
			task.Task.SchedTask.GetTaskId(), jobNodeId)
		return nil
	}


	log.Infof("Success to publish schedTask to `job-Fighter` node to executing, taskId: {%s}, jobNodeId: {%s}",
		task.Task.SchedTask.GetTaskId(), jobNodeId)
	return nil
}

func (m *Manager) publishFinishedTaskToDataCenter(taskId string) {
	task, ok := m.queryRunningTaskCacheOk(taskId)
	if !ok {
		return
	}

	time.Sleep(2*time.Second)  // todo 故意等待小段时间 防止 onTaskResultMsg 因为 网络延迟, 而没收集全 其他 peer 的 eventList

	eventList, err := m.resourceMng.GetDB().GetTaskEventList(taskId)
	if nil != err {
		log.Errorf("Failed to Query all task event list for sending datacenter on publishFinishedTaskToDataCenter, taskId: {%s}, err: {%s}", taskId, err)
		return
	}
	var isFailed bool
	for _, event := range eventList {
		if event.Type == ev.TaskFailed.Type {
			isFailed = true
			break
		}
	}
	var taskState apipb.TaskState
	if isFailed {
		taskState = apipb.TaskState_TaskState_Failed
	} else {
		taskState = apipb.TaskState_TaskState_Succeed
	}

	log.Debugf("Start publishFinishedTaskToDataCenter, taskId: {%s}, taskState: {%s}", taskId, taskState.String())

	finalTask := m.convertScheduleTaskToTask(task.GetTask(), eventList, taskState)
	if err := m.resourceMng.GetDB().InsertTask(finalTask); nil != err {
		log.Errorf("Failed to save task to datacenter on publishFinishedTaskToDataCenter, taskId: {%s}, err: {%s}", taskId, err)
		return
	}

	if err := m.resourceMng.GetDB().RemoveLocalTaskExecuteStatus(taskId); nil != err {
		log.Errorf("Failed to remove task executing status on publishFinishedTaskToDataCenter, taskId: {%s}, err: {%s}", taskId, err)
		return
	}
	// clean local task cache
	m.removeRunningTaskCache(taskId)
	m.resourceMng.ReleaseLocalResourceWithTask("on taskManager.publishFinishedTaskToDataCenter()", taskId, resource.SetAllReleaseResourceOption())

	log.Debugf("Finished pulishFinishedTaskToDataCenter, taskId: {%s}, taskState: {%s}", taskId, taskState)
}
func (m *Manager) sendTaskResultMsgToConsensus(taskId string) {

	taskWrap, ok := m.queryRunningTaskCacheOk(taskId)
	if !ok {
		log.Errorf("Not found taskwrap, taskId: %s", taskId)
		return
	}

	log.Debugf("Start sendTaskResultMsgToConsensus, taskId: {%s}", taskId)

	if err := m.dataCenter.RemoveLocalTaskExecuteStatus(taskId); nil != err {
		log.Errorf("Failed to remove task executing status on sendTaskResultMsgToConsensus, taskId: {%s}, err: {%s}", taskId, err)
		return
	}

	taskResultMsg := m.makeTaskResultByEventList(taskWrap)
	if nil != taskResultMsg {
		taskWrap.ResultCh <- taskResultMsg
	}
	close(taskWrap.ResultCh)

	// clean local task cache
	m.removeRunningTaskCache(taskWrap.Task.SchedTask.GetTaskId())

	log.Debugf("Finished sendTaskResultMsgToConsensus, taskId: {%s}", taskId)
}

func (m *Manager) sendTaskMsgsToScheduler(tasks types.TaskDataArray) {
	m.localTaskMsgCh <- msgs
}
func (m *Manager) sendTaskEvent(event *libTypes.TaskEvent) {
	m.eventCh <- event
}

func (m *Manager) storeBadTask(task *types.Task, events []*libTypes.TaskEvent, reason string) error {
	task.GetTaskData().TaskEvents = events
	task.GetTaskData().EventCount = uint32(len(events))
	task.GetTaskData().State = apipb.TaskState_TaskState_Failed
	task.GetTaskData().Reason = reason
	task.GetTaskData().EndAt = uint64(timeutils.UnixMsec())
	return m.resourceMng.GetDB().InsertTask(task)
}

func (m *Manager) convertScheduleTaskToTask(task *types.Task, eventList []*libTypes.TaskEvent, state apipb.TaskState) *types.Task {
	task.GetTaskData().TaskEvents = eventList
	task.GetTaskData().EventCount = uint32(len(eventList))
	task.GetTaskData().EndAt = uint64(timeutils.UnixMsec())
	task.GetTaskData().State = state
	return task
}

func (m *Manager) makeTaskReadyGoReq(task *types.NeedExecuteTask) (*common.TaskReadyGoReq, error) {

	ownerPort := string(task.Task.Resources.OwnerPeerInfo.Port)
	port, err := strconv.Atoi(ownerPort)
	if nil != err {
		return nil, err
	}

	var dataPartyArr []string
	var powerPartyArr []string
	var receiverPartyArr []string

	peerList := []*common.TaskReadyGoReq_Peer{
		&common.TaskReadyGoReq_Peer{
			Ip:      string(task.Task.Resources.OwnerPeerInfo.Ip),
			Port:    int32(port),
			PartyId: string(task.Task.Resources.OwnerPeerInfo.PartyId),
		},
	}
	dataPartyArr = append(dataPartyArr, string(task.Task.Resources.OwnerPeerInfo.PartyId))

	for _, dataSupplier := range task.Task.Resources.DataSupplierPeerInfoList {
		portStr := string(dataSupplier.Port)
		port, err := strconv.Atoi(portStr)
		if nil != err {
			return nil, err
		}
		peerList = append(peerList, &common.TaskReadyGoReq_Peer{
			Ip:      string(dataSupplier.Ip),
			Port:    int32(port),
			PartyId: string(dataSupplier.PartyId),
		})
		dataPartyArr = append(dataPartyArr, string(dataSupplier.PartyId))
	}

	for _, powerSupplier := range task.Task.Resources.PowerSupplierPeerInfoList {
		portStr := string(powerSupplier.Port)
		port, err := strconv.Atoi(portStr)
		if nil != err {
			return nil, err
		}
		peerList = append(peerList, &common.TaskReadyGoReq_Peer{
			Ip:      string(powerSupplier.Ip),
			Port:    int32(port),
			PartyId: string(powerSupplier.PartyId),
		})

		powerPartyArr = append(powerPartyArr, string(powerSupplier.PartyId))
	}

	for _, receiver := range task.Task.Resources.ResultReceiverPeerInfoList {
		portStr := string(receiver.Port)
		port, err := strconv.Atoi(portStr)
		if nil != err {
			return nil, err
		}
		peerList = append(peerList, &common.TaskReadyGoReq_Peer{
			Ip:      string(receiver.Ip),
			Port:    int32(port),
			PartyId: string(receiver.PartyId),
		})

		receiverPartyArr = append(receiverPartyArr, string(receiver.PartyId))
	}

	contractExtraParams, err := m.makeContractParams(task)
	if nil != err {
		return nil, err
	}
	log.Debugf("Succeed make contractCfg, taskId:{%s}, contractCfg: %s", task.Task.SchedTask.GetTaskId(), contractExtraParams)
	return &common.TaskReadyGoReq{
		TaskId:     task.Task.SchedTask.GetTaskId(),
		ContractId: task.Task.SchedTask.GetTaskData().CalculateContractCode,
		//DataId: "",
		PartyId: task.SelfIdentity.PartyId,
		//EnvId: "",
		Peers:            peerList,
		ContractCfg:      contractExtraParams,
		DataParty:        dataPartyArr,
		ComputationParty: powerPartyArr,
		ResultParty:      receiverPartyArr,
	}, nil
}

func (m *Manager) makeContractParams(task *types.NeedExecuteTask) (string, error) {

	partyId := task.SelfIdentity.PartyId

	var filePath string
	var idColumnName string

	if task.SelfTaskRole == types.TaskOwner || task.SelfTaskRole == types.DataSupplier {

		var find bool

		for _, dataSupplier := range task.Task.SchedTask.GetTaskData().DataSupplier {
			if partyId == dataSupplier.MemberInfo.PartyId {

				metaData, err := m.dataCenter.GetMetadataByDataId(dataSupplier.MetadataId)
				if nil != err {
					return "", err
				}
				filePath = metaData.MetadataData().FilePath

				// 目前只取 第一列 (对于 dataSupplier)
				if len(dataSupplier.ColumnList) != 0 {
					idColumnName = dataSupplier.ColumnList[0].CName
				}
				find = true
				break
			}
		}

		if !find {
			return "", fmt.Errorf("can not find the dataSupplier for find originFilePath, taskId: {%s}, self.IdentityId: {%s}, seld.PartyId: {%s}",
				task.Task.SchedTask.GetTaskId(), task.SelfIdentity.IdentityId, task.SelfIdentity.PartyId)
		}
	}

	// 目前 默认只会用一列, 后面再拓展 ..
	req := &types.FighterTaskReadyGoReqContractCfg{
		PartyId: partyId,
		DataParty: struct {
			InputFile    string `json:"input_file"`
			IdColumnName string `json:"id_column_name"`
		}{
			InputFile:    filePath,
			IdColumnName: idColumnName, // 目前 默认只会用一列, 后面再拓展 .. 只有 dataSupplier 才有, powerSupplier 不会有
		},
	}

	var dynamicParameter map[string]interface{}
	log.Debugf("Start json Unmarshal the `ContractExtraParams`, taskId: {%s}, ContractExtraParams: %s", task.Task.SchedTask.GetTaskId(), task.Task.SchedTask.GetTaskData().ContractExtraParams)
	if "" != task.Task.SchedTask.GetTaskData().ContractExtraParams {
		if err := json.Unmarshal([]byte(task.Task.SchedTask.GetTaskData().ContractExtraParams), &dynamicParameter); nil != err {
			return "", fmt.Errorf("can not json Unmarshal the `ContractExtraParams` of task, taskId: {%s}, self.IdentityId: {%s}, seld.PartyId: {%s}, err: {%s}",
				task.Task.SchedTask.GetTaskId(), task.SelfIdentity.IdentityId, task.SelfIdentity.PartyId, err)
		}
	}
	req.DynamicParameter = dynamicParameter

	b, err := json.Marshal(req)
	if nil != err {
		return "", fmt.Errorf("can not json Marshal the `FighterTaskReadyGoReqContractCfg`, taskId: {%s}, self.IdentityId: {%s}, seld.PartyId: {%s}, err: {%s}",
			task.Task.SchedTask.GetTaskId(), task.SelfIdentity.IdentityId, task.SelfIdentity.PartyId, err)
	}
	return string(b), nil
}

func (m *Manager) addRunningTaskCache(task *types.NeedExecuteTask) {
	m.runningTaskCacheLock.Lock()
	m.runningTaskCache[task.GetTask().GetTaskId()] = task
	m.runningTaskCacheLock.Unlock()
}

func (m *Manager) removeRunningTaskCache(taskId string) {
	m.runningTaskCacheLock.Lock()
	delete(m.runningTaskCache, taskId)
	m.runningTaskCacheLock.Unlock()
}

func (m *Manager) queryRunningTaskCacheOk(taskId string) (*types.NeedExecuteTask, bool) {
	m.runningTaskCacheLock.RLock()
	task, ok := m.runningTaskCache[taskId]
	m.runningTaskCacheLock.RUnlock()
	return task, ok
}

func (m *Manager) queryRunningTaskCache(taskId string) *types.NeedExecuteTask {
	task, _ := m.queryRunningTaskCacheOk(taskId)
	return task
}

func (m *Manager) ForEachRunningTaskCache (f  func(taskId string, task *types.NeedExecuteTask) bool ) {
	m.runningTaskCacheLock.Lock()
	for taskId, task := range m.runningTaskCache {
		if ok := f(taskId, task); ok {
		}
	}
	m.runningTaskCacheLock.Unlock()
}

func (m *Manager) makeTaskResultByEventList(taskWrap *types.NeedExecuteTask) *types.TaskResultMsgWrap {

	if taskWrap.Task.TaskDir == types.SendTaskDir || types.TaskOwner == taskWrap.SelfTaskRole {
		log.Errorf("send task OR task owner can not make TaskResult Msg")
		return nil
	}

	eventList, err := m.dataCenter.GetTaskEventList(taskWrap.Task.SchedTask.GetTaskId())
	if nil != err {
		log.Errorf("Failed to make TaskResultMsg with query task eventList, taskId {%s}, err {%s}", taskWrap.Task.SchedTask.GetTaskId(), err)
		return nil
	}
	return &types.TaskResultMsgWrap{
		TaskResultMsg: &pb.TaskResultMsg{
			ProposalId: taskWrap.ProposalId.Bytes(),
			TaskRole:   taskWrap.SelfTaskRole.Bytes(),
			TaskId:     []byte(taskWrap.Task.SchedTask.GetTaskId()),
			Owner: &pb.TaskOrganizationIdentityInfo{
				PartyId:    []byte(taskWrap.SelfIdentity.PartyId),
				Name:       []byte(taskWrap.SelfIdentity.NodeName),
				NodeId:     []byte(taskWrap.SelfIdentity.NodeId),
				IdentityId: []byte(taskWrap.SelfIdentity.IdentityId),
			},
			TaskEventList: types.ConvertTaskEventArr(eventList),
			CreateAt:      uint64(timeutils.UnixMsec()),
			Sign:          nil,
		},
	}
}

func (m *Manager) handleEvent(event *libTypes.TaskEvent) error {
	eventType := event.Type
	if len(eventType) != ev.EventTypeCharLen {
		return ev.IncEventType
	}
	// TODO need to validate the task that have been processing ? Maybe~
	if event.Type == ev.TaskExecuteSucceedEOF.Type || event.Type == ev.TaskExecuteFailedEOF.Type {
		if task, ok := m.queryRunningTaskCacheOk(event.TaskId); ok {

			log.Debugf("Start handleEvent, `event is the end`, event: %s, current taskDir: {%s}", event.String(), task.Task.TaskDir.String())

			// 先 缓存下 最终休止符 event
			m.resourceMng.GetDB().StoreTaskEvent(event)
			if event.Type == ev.TaskExecuteFailedEOF.Type {
				m.storeTaskFinalEvent(task.GetTask().GetTaskId(), task.GetTaskOrganization().GetIdentityId(), "", apipb.TaskState_TaskState_Failed)
			} else {
				m.storeTaskFinalEvent(task.GetTask().GetTaskId(), task.GetTaskOrganization().GetIdentityId(), "", apipb.TaskState_TaskState_Succeed)
			}

			if task.GetTaskRole() == apipb.TaskRole_TaskRole_Sender {
				m.sendTaskResultMsgToConsensus (event.TaskId)
			} else {
				m.publishFinishedTaskToDataCenter(event.TaskId)

			}
			return nil
		} else {
			return errors.New(fmt.Sprintf("Not found task cache, taskId: {%s}", event.TaskId))
		}

	} else {

		log.Debugf("Start handleEvent, `event is not the end`, event: %s", event.String())
		// 不是休止符 event, 任务还在继续, 保存 event
		return m.resourceMng.GetDB().StoreTaskEvent(event)
	}
}
func (m *Manager) handleNeedExecuteTask(task *types.NeedExecuteTask) {

	log.Debugf("Start handle needExecuteTask, remote pid: {%s} proposalId: {%s}, taskId: {%s}, self taskRole: {%s}, self taskOrganization: {%s}",
		task.GetRemotePID(), task.GetProposalId(), task.GetTask().GetTaskId(), task.GetTaskRole().String(), task.GetTaskOrganization().String())

	// Store task exec status
	if err := m.resourceMng.GetDB().StoreLocalTaskExecuteStatus(task.GetTask().GetTaskId()); nil != err {
		log.Errorf("Failed to store local task about exec status,  remote pid: {%s} proposalId: {%s}, taskId: {%s}, self taskRole: {%s}, self taskOrganization: {%s}, err: {%s}",
			task.GetRemotePID(), task.GetProposalId(), task.GetTask().GetTaskId(), task.GetTaskRole().String(), task.GetTaskOrganization().String(), err)
		return
	}

	switch task.GetTaskRole() {
	case apipb.TaskRole_TaskRole_Sender:
		if err := m.driveTaskForExecute(task); nil != err {
			log.Errorf("Failed to execute task on taskOnwer node, taskId:{%s}, %s", task.GetTask().GetTaskId(), err)

			m.storeTaskFinalEvent(task.GetTask().GetTaskId(), task.GetTaskOrganization().GetIdentityId(), fmt.Sprintf("failed to execute task"), apipb.TaskState_TaskState_Failed)
			m.publishFinishedTaskToDataCenter(task.GetTask().GetTaskId())
		}

	default:
		if err := m.driveTaskForExecute(task); nil != err {
			log.Errorf("Failed to execute task on %s node, taskId: {%s}, %s", task.GetTaskRole().String(), task.GetTask().GetTaskId(), err)

			// 因为是 task 参与者, 所以需要构造 taskResult 发送给 task 发起者.. (里面有解锁 本地资源 ...)
			m.storeTaskFinalEvent(task.GetTask().GetTaskId(), task.GetTaskOrganization().GetIdentityId(), fmt.Sprintf("failed to execute task"), apipb.TaskState_TaskState_Failed)
			m.sendTaskResultMsgToConsensus(task.GetTask().GetTaskId())
		}
	}
}

func (m *Manager) expireTaskMonitor () {

	for taskId, task := range m.runningTaskCache {
		if task.Task.SchedTask.GetTaskData().State == types.TaskStateRunning.String() && task.Task.SchedTask.GetTaskData().StartAt != 0 {

			// the task has running expire
			duration := uint64(timeutils.UnixMsec()) - task.Task.SchedTask.GetTaskData().StartAt
			if duration >= task.Task.SchedTask.GetTaskData().OperationCost.Duration {
				log.Infof("Has task running expire, taskId: {%s}, current running duration: {%d ms}, need running duration: {%d ms}",
					taskId, duration, task.Task.SchedTask.GetTaskData().OperationCost.Duration)

				identityId, _ := m.dataCenter.GetIdentityId()
				m.storeTaskFinalEvent(task.Task.SchedTask.GetTaskId(), identityId, fmt.Sprintf("task running expire"),types.TaskStateFailed)
				switch task.SelfTaskRole {
				case types.TaskOwner:
					m.publishFinishedTaskToDataCenter(taskId)
				default:
					m.sendTaskResultMsgToConsensus(taskId)
				}
			}
		}
	}
}


func (m *Manager) storeTaskFinalEvent(taskId, identityId, extra string, state apipb.TaskState) {
	var evTyp string
	var evMsg string
	if state == apipb.TaskState_TaskState_Failed {
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


// Subscriber 在完成任务时对 task 生成 taskResultMsg 反馈给 发起方
func (t *TwoPC) sendTaskResultMsg(pid peer.ID, msg *types.TaskResultMsgWrap) error {
	if err := handler.SendTwoPcTaskResultMsg(context.TODO(), t.p2p, pid, msg.TaskResultMsg); nil != err {
		err := fmt.Errorf("failed to call `SendTwoPcTaskResultMsg`, taskId: {%s}, taskRole: {%s}, task owner's identityId: {%s}, task owner's peerId: {%s}, err: {%s}",
			msg.TaskResultMsg.TaskId, types.TaskRoleFromBytes(msg.TaskRole).String(), string(msg.TaskResultMsg.Owner.IdentityId), pid, err)
		return err
	}
	return nil
}

// (on Publisher)
func (t *TwoPC) onTaskResultMsg(pid peer.ID, taskResultMsg *types.TaskResultMsgWrap) error {
	msg := fetchTaskResultMsg(taskResultMsg)

	log.Debugf("Received remote taskResultMsg, remote pid: {%s}, taskResultMsg: %s", pid, msg.String())

	has, err := t.resourceMng.GetDB().HasLocalTaskExecute(msg.TaskId)
	if nil != err {
		log.Errorf("Failed to query local task executing status on `onTaskResultMsg`, taskId: {%s}, err: {%s}",
			msg.TaskId, err)
		return fmt.Errorf("query local task failed")
	}

	if !has {
		log.Warnf("Warning not found local task executing status on `onTaskResultMsg`, taskId: {%s}", msg.TaskId)
		return fmt.Errorf("%s, the local task executing status is not found", ctypes.ErrTaskResultMsgInvalid)
	}
	t.storeTaskEvent(pid, msg.TaskId, msg.TaskEventList)
	return nil
}