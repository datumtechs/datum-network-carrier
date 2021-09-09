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

	task.GetTask().GetTaskData().State = apipb.TaskState_TaskState_Running
	task.GetTask().GetTaskData().StartAt = uint64(timeutils.UnixMsec())
	if err := m.resourceMng.GetDB().StoreLocalTask(task.GetTask()); nil != err {
		log.Errorf("Failed to update local task state before executing task, taskId: {%s}, need update state: {%s}, err: {%s}",
			task.GetTask().GetTaskId(), apipb.TaskState_TaskState_Running.String(), err)
	}
	// update local cache
	m.addNeedExecuteTaskCache(task)

	//return fmt.Errorf("Mock task finished")

	switch task.GetTaskRole() {
	case apipb.TaskRole_TaskRole_DataSupplier, apipb.TaskRole_TaskRole_Receiver:
		return m.executeTaskOnDataNode(task)
	case apipb.TaskRole_TaskRole_PowerSupplier:
		return m.executeTaskOnJobNode(task)
	default:
		log.Errorf("Faided to driveTaskForExecute(), Unknown task role, taskId: {%s}, taskRole: {%s}", task.GetTask().GetTaskId(), task.GetTaskRole().String())
		return errors.New("Unknown resource node type")
	}
}

func (m *Manager) executeTaskOnDataNode(task *types.NeedExecuteTask) error {

	// 找到自己的投票
	dataNodeId := task.GetSelfResource().Id

	// clinet *grpclient.DataNodeClient,
	client, has := m.resourceClientSet.QueryDataNodeClient(dataNodeId)
	if !has {
		log.Errorf("Failed to query internal data node, taskId: {%s}, dataNodeId: {%s}",
			task.GetTask().GetTaskId(), dataNodeId)
		return errors.New("data node client not found")
	}
	if client.IsNotConnected() {
		if err := client.Reconnect(); nil != err {
			log.Errorf("Failed to connect internal data node, taskId: {%s}, dataNodeId: {%s}, err: {%s}",
				task.GetTask().GetTaskId(), dataNodeId, err)
			return err
		}
	}

	req, err := m.makeTaskReadyGoReq(task)
	if nil != err {
		log.Errorf("Falied to make taskReadyGoReq, taskId: {%s}, dataNodeId: {%s}, err: {%s}",
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
		log.Errorf("Falied to executing task from `data-Fighter` node response, taskId: {%s}, dataNodeId: {%s}, resp: {%s}",
			task.GetTask().GetTaskId(), dataNodeId, resp.String())
		return nil
	}

	log.Infof("Success to publish schedTask to `data-Fighter` node to executing,  taskId: {%s}, dataNodeId: {%s}",
		task.GetTask().GetTaskId(), dataNodeId)
	return nil
}

func (m *Manager) executeTaskOnJobNode(task *types.NeedExecuteTask) error {

	// 找到自己的投票
	jobNodeId := task.GetSelfResource().Id

	// clinet *grpclient.DataNodeClient,
	client, has := m.resourceClientSet.QueryJobNodeClient(jobNodeId)
	if !has {
		log.Errorf("Failed to query internal job node, taskId: {%s}, jobNodeId: {%s}",
			task.GetTask().GetTaskId(), jobNodeId)
		return errors.New("job node client not found")
	}
	if client.IsNotConnected() {
		if err := client.Reconnect(); nil != err {
			log.Errorf("Failed to connect internal job node, taskId: {%s}, jobNodeId: {%s}, err: {%s}",
				task.GetTask().GetTaskId(), jobNodeId, err)
			return err
		}
	}

	req, err := m.makeTaskReadyGoReq(task)
	if nil != err {
		log.Errorf("Falied to make taskReadyGoReq, taskId: {%s}, jobNodeId: {%s}, err: {%s}",
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
		log.Errorf("Falied to publish schedTask to `job-Fighter` node to executing, taskId: {%s}, jobNodeId: {%s}",
			task.GetTask().GetTaskId(), jobNodeId)
		return nil
	}

	log.Infof("Success to publish schedTask to `job-Fighter` node to executing, taskId: {%s}, jobNodeId: {%s}",
		task.GetTask().GetTaskId(), jobNodeId)
	return nil
}

func (m *Manager) publishFinishedTaskToDataCenter(taskId, partyId string) {
	task, ok := m.queryNeedExecuteTaskCache(taskId, partyId)
	if !ok {
		return
	}

	time.Sleep(2 * time.Second) // todo 故意等待小段时间 防止 onTaskResultMsg 因为 网络延迟, 而没收集全 其他 peer 的 eventList

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
	m.removeNeedExecuteTaskCache(taskId, partyId)
	m.resourceMng.ReleaseLocalResourceWithTask("on taskManager.publishFinishedTaskToDataCenter()", taskId, resource.SetAllReleaseResourceOption())

	log.Debugf("Finished pulishFinishedTaskToDataCenter, taskId: {%s}, taskState: {%s}", taskId, taskState)
}
func (m *Manager) sendTaskResultMsgToConsensus(taskId, partyId string) {

	if ok := m.hasNeedExecuteTaskCache(taskId, partyId); !ok {
		log.Errorf("Not found taskwrap, taskId: %s", taskId)
		return
	}

	log.Debugf("Start sendTaskResultMsgToConsensus, taskId: {%s}", taskId)

	if err := m.resourceMng.GetDB().RemoveLocalTaskExecuteStatus(taskId); nil != err {
		log.Errorf("Failed to remove task executing status on sendTaskResultMsgToConsensus, taskId: {%s}, err: {%s}", taskId, err)
		return
	}

	//taskResultMsg := m.makeTaskResultByEventList(taskWrap)
	//if nil != taskResultMsg {
	//	taskWrap.ResultCh <- taskResultMsg
	//}
	//close(taskWrap.ResultCh)

	// clean local task cache
	m.removeNeedExecuteTaskCache(taskId, partyId)

	log.Debugf("Finished sendTaskResultMsgToConsensus, taskId: {%s}, partyId: {%s}", taskId, partyId)
}

func (m *Manager) sendLocalTaskToScheduler(tasks types.TaskDataArray) {
	m.localTasksCh <- tasks
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
		PartyId: task.GetTaskOrganization().GetPartyId(),
		//EnvId: "",
		Peers:            peerList,
		ContractCfg:      contractExtraParams,
		DataParty:        dataPartyArr,
		ComputationParty: powerPartyArr,
		ResultParty:      receiverPartyArr,
	}, nil
}

func (m *Manager) makeContractParams(task *types.NeedExecuteTask) (string, error) {

	partyId := task.GetTaskOrganization().GetPartyId()

	var filePath string
	var idColumnName string

	if task.GetTaskRole() == apipb.TaskRole_TaskRole_DataSupplier {

		var find bool

		for _, dataSupplier := range task.GetTask().GetTaskData().GetDataSuppliers() {
			if partyId == dataSupplier.GetOrganization().PartyId {

				metaData, err := m.resourceMng.GetDB().GetMetadataByDataId(dataSupplier.MetadataId)
				if nil != err {
					return "", err
				}
				filePath = metaData.MetadataData().FilePath

				// 目前只取 第一列 (对于 dataSupplier)
				if len(dataSupplier.GetColumns()) != 0 {
					idColumnName = dataSupplier.GetColumns()[0].CName
				}
				find = true
				break
			}
		}

		if !find {
			return "", fmt.Errorf("can not find the dataSupplier for find originFilePath, taskId: {%s}, self.IdentityId: {%s}, seld.PartyId: {%s}",
				task.GetTask().GetTaskId(), task.GetTaskOrganization().GetIdentityId(), task.GetTaskOrganization().GetPartyId())
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
	log.Debugf("Start json Unmarshal the `ContractExtraParams`, taskId: {%s}, ContractExtraParams: %s", task.GetTask().GetTaskId(), task.GetTask().GetTaskData().GetContractExtraParams())
	if "" != task.GetTask().GetTaskData().GetContractExtraParams() {
		if err := json.Unmarshal([]byte(task.GetTask().GetTaskData().GetContractExtraParams()), &dynamicParameter); nil != err {
			return "", fmt.Errorf("can not json Unmarshal the `ContractExtraParams` of task, taskId: {%s}, self.IdentityId: {%s}, seld.PartyId: {%s}, err: {%s}",
				task.GetTask().GetTaskId(), task.GetTaskOrganization().GetIdentityId(), task.GetTaskOrganization().GetPartyId(), err)
		}
	}
	req.DynamicParameter = dynamicParameter

	b, err := json.Marshal(req)
	if nil != err {
		return "", fmt.Errorf("can not json Marshal the `FighterTaskReadyGoReqContractCfg`, taskId: {%s}, self.IdentityId: {%s}, seld.PartyId: {%s}, err: {%s}",
			task.GetTask().GetTaskId(), task.GetTaskOrganization().GetIdentityId(), task.GetTaskOrganization().GetPartyId(), err)
	}
	return string(b), nil
}

func (m *Manager) addNeedExecuteTaskCache(task *types.NeedExecuteTask) {
	m.runningTaskCacheLock.Lock()
	cache, ok := m.runningTaskCache[task.GetTask().GetTaskId()]
	if !ok {
		cache = make(map[string]*types.NeedExecuteTask, 0)
	}
	cache[task.GetTaskOrganization().GetPartyId()] = task
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

func (m *Manager) makeTaskResultByEventList(task *types.NeedExecuteTask) *types.TaskResultMsgWrap {

	if task.GetTaskRole() == apipb.TaskRole_TaskRole_Sender {
		log.Errorf("send task OR task owner can not make TaskResult Msg")
		return nil
	}

	eventList, err := m.resourceMng.GetDB().GetTaskEventList(task.GetTask().GetTaskId())
	if nil != err {
		log.Errorf("Failed to make TaskResultMsg with query task eventList, taskId {%s}, err {%s}", task.GetTask().GetTaskId(), err)
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
		if task, ok := m.queryNeedExecuteTaskCache(event.TaskId); ok {

			log.Debugf("Start handleEvent, `event is the end`, event: %s, current taskDir: {%s}", event.String(), task.Task.TaskDir.String())

			// 先 缓存下 最终休止符 event
			m.resourceMng.GetDB().StoreTaskEvent(event)
			if event.Type == ev.TaskExecuteFailedEOF.Type {
				m.storeTaskFinalEvent(task.GetTask().GetTaskId(), task.GetTaskOrganization().GetIdentityId(), "", apipb.TaskState_TaskState_Failed)
			} else {
				m.storeTaskFinalEvent(task.GetTask().GetTaskId(), task.GetTaskOrganization().GetIdentityId(), "", apipb.TaskState_TaskState_Succeed)
			}

			if task.GetTaskRole() == apipb.TaskRole_TaskRole_Sender {
				m.sendTaskResultMsgToConsensus(event.TaskId)
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

func (m *Manager) expireTaskMonitor() {

	for taskId, task := range m.runningTaskCache {
		if task.Task.SchedTask.GetTaskData().State == types.TaskStateRunning.String() && task.Task.SchedTask.GetTaskData().StartAt != 0 {

			// the task has running expire
			duration := uint64(timeutils.UnixMsec()) - task.Task.SchedTask.GetTaskData().StartAt
			if duration >= task.Task.SchedTask.GetTaskData().OperationCost.Duration {
				log.Infof("Has task running expire, taskId: {%s}, current running duration: {%d ms}, need running duration: {%d ms}",
					taskId, duration, task.Task.SchedTask.GetTaskData().OperationCost.Duration)

				identityId, _ := m.dataCenter.GetIdentityId()
				m.storeTaskFinalEvent(task.Task.SchedTask.GetTaskId(), identityId, fmt.Sprintf("task running expire"), types.TaskStateFailed)
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
