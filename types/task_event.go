package types

import (
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)
//type TaskEvent struct {
//	GetTaskId   string     `json:"taskId"`
//	Type     string     `json:"type"`
//	GetCreateAt uint64     `json:"createAt"`
//	Content  string     `json:"content"`
//	GetSender    *NodeAlias `json:"owner"`
//}
//func ConvertTaskEventToPB (event *TaskEvent) *pb.TaskEventShow {
//	return &pb.TaskEventShow{
//		Type:     event.Type,
//		GetTaskId:   event.GetTaskId,
//		GetSender:    ConvertNodeAliasToPB(event.GetSender),
//		Content:  event.Content,
//		GetCreateAt: event.GetCreateAt,
//	}
//}

//func ConvertTaskEventFromPB (event *pb.TaskEventShow) *TaskEvent {
//	return &TaskEvent{
//		Type:     event.Type,
//		GetTaskId:   event.GetTaskId,
//		GetSender:    ConvertNodeAliasFromPB(event.GetSender),
//		Content:  event.Content,
//		GetCreateAt: event.GetCreateAt,
//	}
//}


//func ConvertTaskEventArrToPB (events []*libtypes.TaskEvent) []*pb.TaskEventShow {
//
//	arr := make([]*pb.TaskEventShow, len(events))
//	for i, event := range events {
//		e := &pb.TaskEventShow{
//			Type:     event.Type,
//			GetTaskId:   event.GetTaskId,
//			GetSender:    ConvertNodeAliasToPB(event.GetSender),
//			Content:  event.Content,
//			GetCreateAt: event.GetCreateAt,
//		}
//		arr[i] = e
//	}
//
//	return arr
//}

//func ConvertTaskEventArrFromPB (events []*pb.TaskEventShow) []*TaskEvent {
//	arr := make([]*TaskEvent, len(events))
//	for i, event := range events {
//		e := &TaskEvent{
//			Type:     event.Type,
//			GetTaskId:   event.GetTaskId,
//			GetSender:    ConvertNodeAliasFromPB(event.GetSender),
//			Content:  event.Content,
//			GetCreateAt: event.GetCreateAt,
//		}
//		arr[i] = e
//	}
//
//	return arr
//}


type ReportTaskEvent struct {
	 PartyId   string
	 Event     *libtypes.TaskEvent
}

func NewReportTaskEvent(partyId string, event *libtypes.TaskEvent) *ReportTaskEvent {
	return &ReportTaskEvent{
		PartyId: partyId,
		Event: event,
	}
}
func (rte *ReportTaskEvent) GetPartyId() string            { return rte.PartyId }
func (rte *ReportTaskEvent) GetEvent() *libtypes.TaskEvent { return rte.Event }