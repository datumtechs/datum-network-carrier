package task

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/core"
	ev "github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
)



type Manager struct {
	eventCh     chan *types.TaskEventInfo
	dataCenter   core.CarrierDB
	eventEngine *ev.EventEngine
	resourceMng *resource.Manager
	parser      *TaskParser
	validator   *TaskValidator
	// recv the taskMsgs from messageHandler
	taskCh <-chan types.TaskMsgs
	// send the validated taskMsgs to scheduler
	sendTaskCh chan<- types.TaskMsgs
	// TODO 接收 被调度好的 task, 准备发给自己的  Fighter-Py
	recvSchedTaskCh  chan *types.ConsensusScheduleTask
	runningTaskCache map[string]*types.ConsensusScheduleTask

	// internal resource node set (Fighter node grpc client set)
	resourceClientSet *grpclient.InternalResourceClientSet

	// TODO 用于接收 己方已连接 或 断开连接的 Fighter-Py 的 grpc client

}

func NewTaskManager(dataCenter core.CarrierDB, eventEngine *ev.EventEngine,
	resourceMng *resource.Manager, resourceClientSet *grpclient.InternalResourceClientSet,
	taskCh chan types.TaskMsgs, sendTaskCh chan types.TaskMsgs,
	recvSchedTaskCh chan *types.ConsensusScheduleTask) *Manager {

	m := &Manager{
		eventCh:           make(chan *types.TaskEventInfo, 10),
		dataCenter:        dataCenter,
		eventEngine:       eventEngine,
		resourceMng:       resourceMng,
		resourceClientSet: resourceClientSet,
		parser:            newTaskParser(),
		validator:         newTaskValidator(),
		taskCh:            taskCh,
		sendTaskCh:        sendTaskCh,
		recvSchedTaskCh:   recvSchedTaskCh,
		runningTaskCache:  make(map[string]*types.ConsensusScheduleTask, 0),
	}
	go m.loop()
	return m
}

func (m *Manager) handleEvent(event *types.TaskEventInfo) error {
	eventType := event.Type
	if len(eventType) != ev.EventTypeCharLen {
		return ev.IncEventType
	}
	// TODO need to validate the task that have been processing ? Maybe~
	if event.Type == ev.ExecuteComputeSucceed.Type || event.Type == ev.ExecuteComputeFailed.Type {
		m.eventEngine.StoreEvent(event)
////// 
		return nil
	} else {
		return m.eventEngine.StoreEvent(event)
	}
}

func (m *Manager) loop() {

	for {
		select {
		case event := <-m.eventCh:
			if err := m.handleEvent(event); nil != err {
				log.Error("Failed to store task event on local", "taskId", event.TaskId, "event", event.String())
			}
		case task := <-m.recvSchedTaskCh:
			// 对接收到 经 Scheduler  调度好的 task  转发给自己的 Fighter-Py
			switch task.TaskState {
			case types.TaskStateFailed, types.TaskStateSuccess:
				eventList, err := m.dataCenter.GetTaskEventList(task.SchedTask.TaskId)
				if nil != err {
					log.Error("Failed to Query all task event list for sending datacenter", "taskId", task.SchedTask.TaskId)
					continue
				}
				if err := m.dataCenter.InsertTask(m.convertScheduleTaskToTask(task.SchedTask, eventList)); nil != err {
					log.Error("Failed to save task to datacenter", "taskId", task.SchedTask.TaskId)
					continue
				}
				// clean local task cache
				delete(m.runningTaskCache, task.SchedTask.TaskId)
			case types.TaskStateRunning:
				m.runningTaskCache[task.SchedTask.TaskId] = task
				// TODO 发给 Fighter-py

			default:
				log.Error("Failed to handle unknown task", "taskId", task.SchedTask.TaskId)
			}
		default:
		}
	}
}


func (m *Manager) SendTaskMsgs(msgs types.TaskMsgs) error {
	if len(msgs) == 0 {
		return fmt.Errorf("Receive some empty task msgs")
	}

	if errTasks, err := m.parser.ParseTask(msgs); nil != err {
		for _, errtask := range errTasks {

			events, _ := m.dataCenter.GetTaskEventList(errtask.TaskId)
			events = append(events, m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
				errtask.TaskId, errtask.Onwer().IdentityId, fmt.Sprintf("failed to parse taskMsg")))

			if e := m.storeErrTaskMsg(errtask, types.ConvertTaskEventArrToDataCenter(events), "failed to parse taskMsg"); nil != e {
				log.Error("Failed to store the err taskMsg", "taskId", errtask)
			}
		}
		return err
	}

	if errTasks, err := m.validator.validateTaskMsg(msgs); nil != err {
		for _, errtask := range errTasks {
			events, _ := m.dataCenter.GetTaskEventList(errtask.TaskId)
			events = append(events, m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
				errtask.TaskId, errtask.Onwer().IdentityId, fmt.Sprintf("failed to validate taskMsg")))

			if e := m.storeErrTaskMsg(errtask, types.ConvertTaskEventArrToDataCenter(events), "failed to validate taskMsg"); nil != e {
				log.Error("Failed to store the err taskMsg", "taskId", errtask)
			}
		}
		return err
	}
	// transfer `taskMsgs` to Scheduler
	go func() {
		m.sendTaskCh <- msgs
	}()
	return nil
}

func (m *Manager) SendTaskEvent(event *types.TaskEventInfo) error {
	m.eventCh <- event
	return nil
}

//func (m *Manager) driveTask(task *types.ConsensusScheduleTask) error {
//
//}

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
				Alias: "",
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
				Alias: "",
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
		// EndAt
		// 少了 StartAt
		AlgoSupplier: &libTypes.OrganizationData{
			Alias: "",
			Identity: msg.OwnerIdentityId(),
			NodeId: msg.OwnerNodeId(),
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

// TODO 转换
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
	return nil
}