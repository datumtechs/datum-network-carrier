package task

import (
	"fmt"
	ev "github.com/RosettaFlow/Carrier-Go/core/evengine"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/lib/fighter/common"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"time"
)



func (m *Manager) driveTaskForExecute (taskRole types.TaskRole , task *types.DoneScheduleTaskChWrap) error {

	switch taskRole {
	case types.TaskOnwer:

		dataNodeList, err := m.dataCenter.GetRegisterNodeList(types.PREFIX_TYPE_DATANODE)
		if nil != err {
			return err
		}
		ip := string(task.Task.Resources.OwnerPeerInfo.Ip)
		port := string(task.Task.Resources.OwnerPeerInfo.Port)

		var dataNodeId string
		for _, dataNode := range dataNodeList {
			if ip == dataNode.ExternalIp && port == dataNode.ExternalPort {
				dataNodeId = dataNode.Id
				break
			}
		}
		return m.executeTaskOnDataNode(dataNodeId, task)

	case types.DataSupplier:
		dataNodeList, err := m.dataCenter.GetRegisterNodeList(types.PREFIX_TYPE_DATANODE)
		if nil != err {
			return err
		}

		tmp := make(map[string]struct{}, len(task.Task.Resources.DataSupplierPeerInfoList))
		for _, dataNode := range task.Task.Resources.DataSupplierPeerInfoList {
			tmp[string(dataNode.Ip) + "_" + string(dataNode.Port)] = struct{}{}
		}

		// task 中的 dataNode 可能是 单个组织的 多个 dataNode,
		// 逐个下发 task
		for _, dataNode := range dataNodeList {
			if _, ok := tmp[dataNode.ExternalIp + "_" + dataNode.ExternalPort]; ok {
				if err := m.executeTaskOnDataNode(dataNode.Id, task); nil != err {
					log.Errorf("Failed to execute task on dataNode: %s, %s", dataNode.Id, err)
					return err
				}
			}
		}

	case types.PowerSupplier:
		jobNodeList, err := m.dataCenter.GetRegisterNodeList(types.PREFIX_TYPE_JOBNODE)
		if nil != err {
			return err
		}

		tmp := make(map[string]struct{}, len(task.Task.Resources.PowerSupplierPeerInfoList))
		for _, jobNode := range task.Task.Resources.PowerSupplierPeerInfoList {
			tmp[string(jobNode.Ip) + "_" + string(jobNode.Port)] = struct{}{}
		}

		// task 中的 jobNode 可能是 单个组织的 多个 jobNode,
		// 逐个下发 task
		for _, jobNode := range jobNodeList {
			if _, ok := tmp[jobNode.ExternalIp + "_" + jobNode.ExternalPort]; ok {
				if err := m.executeTaskOnJobNode(jobNode.Id, task); nil != err {
					log.Errorf("Failed to execute task on jobNode: %s, %s", jobNode.Id, err)
					return err
				}
			}
		}

	case types.ResultSupplier:
		dataNodeList, err := m.dataCenter.GetRegisterNodeList(types.PREFIX_TYPE_DATANODE)
		if nil != err {
			return err
		}

		tmp := make(map[string]struct{}, len(task.Task.Resources.ResultReceiverPeerInfoList))
		for _, receiveNode := range task.Task.Resources.ResultReceiverPeerInfoList {
			tmp[string(receiveNode.Ip) + "_" + string(receiveNode.Port)] = struct{}{}
		}

		// task 中的 dataNode 可能是 单个组织的 多个 dataNode,
		// 逐个下发 task
		for _, dataNode := range dataNodeList {
			if _, ok := tmp[dataNode.ExternalIp + "_" + dataNode.ExternalPort]; ok {
				if err := m.executeTaskOnDataNode(dataNode.Id, task); nil != err {
					log.Errorf("Failed to execute task on receiveNode: %s, %s", dataNode.Id, err)
					return err
				}
			}
		}
	}
	return nil
}

func (m *Manager) executeTaskOnDataNode(nodeId string, task *types.DoneScheduleTaskChWrap) error {

	// clinet *grpclient.DataNodeClient,
	client, isconn := m.resourceClientSet.QueryDataNodeClient(nodeId)
	if !isconn {
		if err := client.Reconnect(); nil != err {
			log.Error("Failed to connect internal data node", "nodeId", nodeId, "err", err)
			return err
		}
	}
	resp, err := client.HandleTaskReadyGo(m.convertScheduleTaskToTaskReadyGoReq(task.Task.SchedTask, task.Task.Resources))
	if nil != err {
		log.Errorf("Falied to publish schedTask to `data-Fighter` node to executing, taskId: %s, %s", task.Task.SchedTask.TaskId, err)
		return err
	}
	if !resp.Ok {
		log.Errorf("Falied to publish schedTask to `data-Fighter` node to executing, taskId: %s, %s", task.Task.SchedTask.TaskId, resp.Msg)
		return nil
	}

	task.Task.SchedTask.StartAt = uint64(time.Now().UnixNano())
	m.addRunningTaskCache(task)

	log.Infof("Success to publish schedTask to `data-Fighter` node to executing, taskId: %s", task.Task.SchedTask.TaskId)
	return nil
}

func (m *Manager) executeTaskOnJobNode(nodeId string, task *types.DoneScheduleTaskChWrap) error {

	//clinet *grpclient.JobNodeClient,
	client, isconn := m.resourceClientSet.QueryJobNodeClient(nodeId)
	if !isconn {
		if err := client.Reconnect(); nil != err {
			log.Error("Failed to connect internal job node", "nodeId", nodeId, "err", err)
			return err
		}
	}

	resp, err := client.HandleTaskReadyGo(m.convertScheduleTaskToTaskReadyGoReq(task.Task.SchedTask, task.Task.Resources))
	if nil != err {
		log.Errorf("Falied to publish schedTask to `job-Fighter` node to executing, taskId: %s, %s", task.Task.SchedTask.TaskId, err)
		return err
	}
	if !resp.Ok {
		log.Errorf("Falied to publish schedTask to `job-Fighter` node to executing, taskId: %s, %s", task.Task.SchedTask.TaskId, resp.Msg)
		return nil
	}

	task.Task.SchedTask.StartAt = uint64(time.Now().UnixNano())
	m.addRunningTaskCache(task)
	log.Infof("Success to publish schedTask to `job-Fighter` node to executing, taskId: %s", task.Task.SchedTask.TaskId)
	return nil
}


func (m *Manager) pulishFinishedTaskToDataCenter(taskId string) {
	taskWrap, ok := m.queryRunningTaskCacheOk(taskId)
	if !ok {
		return
	}

	eventList, err := m.dataCenter.GetTaskEventList(taskWrap.Task.SchedTask.TaskId)
	if nil != err {
		log.Error("Failed to Query all task event list for sending datacenter", "taskId", taskWrap.Task.SchedTask.TaskId)
		return
	}
	if err := m.dataCenter.InsertTask(m.convertScheduleTaskToTask(taskWrap.Task.SchedTask, eventList)); nil != err {
		log.Error("Failed to save task to datacenter", "taskId", taskWrap.Task.SchedTask.TaskId)
		return
	}
	close(taskWrap.ResultCh)
	// clean local task cache
	m.removeRunningTaskCache(taskId)
	// 解锁 本地 资源缓存
	m.resourceMng.UnLockLocalResourceWithTask(taskId)
}


func (m *Manager) sendTaskMsgsToScheduler(msgs types.TaskMsgs) {
	m.localTaskMsgCh <- msgs
}
func (m *Manager) sendTaskEvent(event *types.TaskEventInfo){
	m.eventCh <- event
}

func (m *Manager) sendTaskResultMsgToConsensus(taskId string) {
	taskWrap, ok := m.queryRunningTaskCacheOk(taskId)
	if !ok {
		log.Errorf( "Not found taskwrap, taskId: %s", taskId)
		return
	}
	taskResultMsg := m.makeTaskResult(taskWrap)
	if nil != taskResultMsg {
		taskWrap.ResultCh <- taskResultMsg
	}

	close(taskWrap.ResultCh)
	// clean local task cache
	m.removeRunningTaskCache(taskWrap.Task.SchedTask.TaskId)
	// 解锁 本地 资源缓存
	m.resourceMng.UnLockLocalResourceWithTask(taskId)
}

func (m *Manager) storeErrTaskMsg(msg *types.TaskMsg, events []*libTypes.EventData, reason string) error {

	// make dataSupplierArr
	metadataSupplierArr := make([]*libTypes.TaskMetadataSupplierData, len(msg.PartnerTaskSuppliers()))
	for i, dataSupplier := range msg.PartnerTaskSuppliers() {

		data, err := m.dataCenter.GetMetadataByDataId(dataSupplier.MetaData.MetaDataId)
		if nil != err {
			return err
		}
		metaData := types.NewOrgMetaDataInfoFromMetadata(data)
		mclist := metaData.MetaData.ColumnMetas

		columnList := make([]*libTypes.ColumnMeta, len(dataSupplier.MetaData.ColumnIndexList))
		for j, index := range dataSupplier.MetaData.ColumnIndexList {
			columnList[j] = &libTypes.ColumnMeta{
				Cindex: uint32(index),
				Cname: mclist[index].Cname,
				Ctype: mclist[index].Ctype,
				// unit:
				Csize: mclist[index].Csize,
				Ccomment: mclist[index].Ccomment,
			}
		}

		metadataSupplierArr[i] = &libTypes.TaskMetadataSupplierData{
			Organization: &libTypes.OrganizationData{
				PartyId: dataSupplier.PartyId,
				Identity: dataSupplier.IdentityId,
				NodeId: dataSupplier.NodeId,
				NodeName: dataSupplier.Name,
			},
			MetaId: metaData.MetaData.MetaDataSummary.MetaDataId,
			MetaName:  metaData.MetaData.MetaDataSummary.TableName,
			ColumnList: columnList,
		}
	}

	// make powerSupplierArr (Empty powerSupplierArr)

	// make receiverArr
	receiverArr := make([]*libTypes.TaskResultReceiverData, len(msg.ReceiverDetails()))
	for i, recv := range msg.ReceiverDetails() {
		receiverArr[i] = &libTypes.TaskResultReceiverData{
			Receiver: &libTypes.OrganizationData{
				PartyId: recv.PartyId,
				Identity: recv.IdentityId,
				NodeId: recv.NodeId,
				NodeName: recv.Name,
			},
			Provider: make([]*libTypes.OrganizationData, 0),
		}
	}


	task :=  types.NewTask(&libTypes.TaskData{
		Identity: msg.OwnerNodeId(),
		NodeId: msg.OwnerNodeId(),
		NodeName:msg.OwnerName(),
		DataId: "",
		// the status of data, N means normal, D means deleted.
		DataStatus: types.ResourceDataStatusN.String(),
		TaskId: msg.TaskId,
		TaskName: msg.TaskName(),
		State: types.TaskStateFailed.String(),
		Reason: reason,
		EventCount: uint32(len(events)),
		// Desc
		CreateAt: msg.CreateAt(),
		EndAt: uint64(time.Now().UnixNano()),
		// 少了 StartAt
		AlgoSupplier: &libTypes.OrganizationData{
			PartyId:  msg.OwnerPartyId(),
			Identity: msg.OwnerIdentityId(),
			NodeId:   msg.OwnerNodeId(),
			NodeName: msg.OwnerName(),
		},
		TaskResource: &libTypes.TaskResourceData{
			CostMem: msg.OperationCost().Mem,
			CostProcessor: uint32(msg.OperationCost().Processor),
			CostBandwidth: msg.OperationCost().Bandwidth,
			Duration: msg.OperationCost().Duration,
		},
		MetadataSupplier: metadataSupplierArr,
		ResourceSupplier: make([]*libTypes.TaskResourceSupplierData, 0),
		Receivers: receiverArr,
		//PartnerList:
		EventDataList: events,
	})
	return m.dataCenter.InsertTask(task)
}


// TODO 需要实现
func (m *Manager) convertScheduleTaskToTask(task *types.ScheduleTask, eventList []*types.TaskEventInfo)  *types.Task {

	//
	//types.NewTask(&libTypes.TaskData{
	//
	//})
	//
	//partners := make([]*libTypes.TaskMetadataSupplierData, len(task.Partners))
	//for i, p := range task.Partners {
	//	partner := &libTypes.TaskMetadataSupplierData {
	//
	//	}
	//	partners[i] = partner
	//}
	//
	//powerArr := make([]*types.ScheduleTaskPowerSupplier, len(powers))
	//for i, p := range powers {
	//	power := &types.ScheduleTaskPowerSupplier{
	//		NodeAlias: p,
	//	}
	//	powerArr[i] = power
	//}
	//
	//receivers := make([]*types.ScheduleTaskResultReceiver, len(task.ReceiverDetails()))
	//for i, r := range task.ReceiverDetails() {
	//	receiver := &types.ScheduleTaskResultReceiver{
	//		NodeAlias: r.NodeAlias,
	//		Providers: r.Providers,
	//	}
	//	receivers[i] = receiver
	//}
	//return &types.ScheduleTask{
	//	TaskId:   task.TaskId,
	//	TaskName: task.TaskName(),
	//	Owner: &types.ScheduleTaskDataSupplier{
	//		NodeAlias: task.Onwer(),
	//		MetaData:  task.OwnerTaskSupplier().MetaData,
	//	},
	//	Partners:              partners,
	//	PowerSuppliers:        powerArr,
	//	Receivers:             receivers,
	//	CalculateContractCode: task.CalculateContractCode(),
	//	DataSplitContractCode: task.DataSplitContractCode(),
	//	OperationCost:         task.OperationCost(),
	//	CreateAt:              task.CreateAt(),
	//}

	//EndAt:  time.Now().UnixNano()
	return nil
}
// TODO 需要实现
func (m *Manager) convertScheduleTaskToTaskReadyGoReq(task *types.ScheduleTask, resources  *pb.ConfirmTaskPeerInfo) *common.TaskReadyGoReq {

	return &common.TaskReadyGoReq{
		//TaskId
		//ContractId
		//DataId
		//PartyId
		//EnvId
		//Peers
		//ContractCfg
		//DataParty
		//ComputationParty
		//ResultParty
	}
}

func (m *Manager) addRunningTaskCache(task *types.DoneScheduleTaskChWrap) {
	m.runningTaskCacheLock.Lock()
	m.runningTaskCache[task.Task.SchedTask.TaskId] = task
	m.runningTaskCacheLock.Unlock()
}

func (m *Manager) removeRunningTaskCache(taskId string) {
	m.runningTaskCacheLock.Lock()
	delete(m.runningTaskCache, taskId)
	m.runningTaskCacheLock.Unlock()
}

func (m *Manager) queryRunningTaskCacheOk(taskId string) (*types.DoneScheduleTaskChWrap, bool) {
	task, ok := m.runningTaskCache[taskId]
	return task, ok
}

func (m *Manager) queryRunningTaskCache(taskId string) *types.DoneScheduleTaskChWrap {
	task, _ := m.queryRunningTaskCacheOk(taskId)
	return task
}
// TODO 需要实现
func (m *Manager) makeTaskResult (taskWrap *types.DoneScheduleTaskChWrap) *types.TaskResultMsgWrap {


	// TODO 需要查出自己存在本地的所有 task 信息 和event 信息, 并删除本地 task  和 event 内容

	return &types.TaskResultMsgWrap{

	}
}


func (m *Manager) handleDoneScheduleTask(taskId string) {

	task, ok := m.queryRunningTaskCacheOk(taskId)
	if !ok {
		return
	}

	switch task.SelfTaskRole {
	case types.TaskOnwer:
		switch task.Task.TaskState {
		case types.TaskStateFailed, types.TaskStateSuccess:

			// 发起方直接 往 dataCenter 发送数据 (里面有解锁 本地资源 ...)
			m.pulishFinishedTaskToDataCenter(taskId)

		case types.TaskStateRunning:

			if err := m.driveTaskForExecute(task.SelfTaskRole, task); nil != err {
				log.Errorf("Failed to execute task on taskOnwer node, taskId: %s, %s", task.Task.SchedTask.TaskId, err)
				event := m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
					task.Task.SchedTask.TaskId, task.Task.SchedTask.Owner.IdentityId, fmt.Sprintf("failed to execute task"))
				// 因为是 自己的任务, 所以直接将 task  和 event list  发给 dataCenter  (里面有解锁 本地资源 ...)
				m.dataCenter.StoreTaskEvent(event)
				m.pulishFinishedTaskToDataCenter(taskId)  //
			}
			// TODO 而执行最终[成功]的 根据 Fighter 上报的 event 在 handleEvent() 里面处理
		default:
			log.Error("Failed to handle unknown task", "taskId", task.Task.SchedTask.TaskId)
		}
	//case types.DataSupplier:
	//case types.PowerSupplier:
	//case types.ResultSupplier:
	default:
		switch task.Task.TaskState {
		case types.TaskStateFailed, types.TaskStateSuccess:
			// 因为是 task 参与者, 所以需要构造 taskResult 发送给 task 发起者..  (里面有解锁 本地资源 ...)
			m.sendTaskResultMsgToConsensus(taskId)
		case types.TaskStateRunning:

			if err := m.driveTaskForExecute(task.SelfTaskRole, task); nil != err {
				log.Errorf("Failed to execute task on taskOnwer node, taskId: %s, %s", task.Task.SchedTask.TaskId, err)
				identityId, _ := m.dataCenter.GetIdentityId()
				event := m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
					task.Task.SchedTask.TaskId, identityId, fmt.Sprintf("failed to execute task"))

				// 因为是 task 参与者, 所以需要构造 taskResult 发送给 task 发起者.. (里面有解锁 本地资源 ...)
				m.dataCenter.StoreTaskEvent(event)
				m.sendTaskResultMsgToConsensus(taskId)
			}
		default:
			log.Error("Failed to handle unknown task", "taskId", task.Task.SchedTask.TaskId)
		}
	}
}