package task

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/core"
	ev "github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
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
	recvSchedTaskCh <-chan *types.ConsensusScheduleTask

	// TODO 持有 己方的所有 Fighter-Py 的 grpc client

	// TODO 用于接收 己方已连接 或 断开连接的 Fighter-Py 的 grpc client

}

func NewTaskManager(dataCenter core.CarrierDB, eventEngine *ev.EventEngine,
	resourceMng *resource.Manager,
	taskCh chan types.TaskMsgs, sendTaskCh chan types.TaskMsgs,
	recvSchedTaskCh chan *types.ConsensusScheduleTask) *Manager {

	m := &Manager{
		eventCh:         make(chan *types.TaskEventInfo, 10),
		dataCenter:       dataCenter,
		eventEngine:     eventEngine,
		resourceMng:     resourceMng,
		parser: newTaskParser(),
		validator: newTaskValidator(),
		taskCh:          taskCh,
		sendTaskCh:      sendTaskCh,
		recvSchedTaskCh: recvSchedTaskCh,
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

	return m.eventEngine.StoreEvent(event)
}

func (m *Manager) loop() {

	for {
		select {
		case event := <-m.eventCh:
			m.handleEvent(event)
		case taskMsgs := <-m.taskCh:
			if len(taskMsgs) == 0 {
				continue
			}
			// 先对 task 解析 和 校验
			if errTasks, err := m.parser.ParseTask(taskMsgs); nil != err {
				for _, errtask := range errTasks {

					events, _ := m.dataCenter.GetTaskEventList(errtask.TaskId)
					events = append(events, m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
						errtask.TaskId, errtask.Onwer().IdentityId, fmt.Sprintf("failed to parse taskMsg")))

					if e := m.storeErrTaskMsg(errtask, types.ConvertTaskEventArrToDataCenter(events), "failed to parse taskMsg"); nil != e {
						log.Error("Failed to store the err taskMsg", "taskId", errtask)
					}
				}
				continue
			}

			if errTasks, err := m.validator.validateTaskMsg(taskMsgs); nil != err {
				for _, errtask := range errTasks {
					events, _ := m.dataCenter.GetTaskEventList(errtask.TaskId)
					events = append(events, m.eventEngine.GenerateEvent(ev.TaskFailed.Type,
						errtask.TaskId, errtask.Onwer().IdentityId, fmt.Sprintf("failed to validate taskMsg")))

					if e := m.storeErrTaskMsg(errtask, types.ConvertTaskEventArrToDataCenter(events), "failed to validate taskMsg"); nil != e {
						log.Error("Failed to store the err taskMsg", "taskId", errtask)
					}
				}
				continue
			}

			// 再 转发给 Scheduler 处理
			go func() {
				m.sendTaskCh <- taskMsgs
			}()
		case task := <-m.recvSchedTaskCh:
			// TODO 对接收到 经 Scheduler  调度好的 task  转发给自己的 Fighter-Py
			_ = task

		default:
		}
	}
}

func (m *Manager) SendTaskEvent(event *types.TaskEventInfo) error {
	m.eventCh <- event
	return nil
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