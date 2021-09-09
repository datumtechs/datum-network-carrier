package types

import (
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)
//type TaskEvent struct {
//	GetTaskId   string     `json:"taskId"`
//	Type     string     `json:"type"`
//	CreateAt uint64     `json:"createAt"`
//	Content  string     `json:"content"`
//	Owner    *NodeAlias `json:"owner"`
//}
//func ConvertTaskEventToPB (event *TaskEvent) *pb.TaskEventShow {
//	return &pb.TaskEventShow{
//		Type:     event.Type,
//		GetTaskId:   event.GetTaskId,
//		Owner:    ConvertNodeAliasToPB(event.Owner),
//		Content:  event.Content,
//		CreateAt: event.CreateAt,
//	}
//}

//func ConvertTaskEventFromPB (event *pb.TaskEventShow) *TaskEvent {
//	return &TaskEvent{
//		Type:     event.Type,
//		GetTaskId:   event.GetTaskId,
//		Owner:    ConvertNodeAliasFromPB(event.Owner),
//		Content:  event.Content,
//		CreateAt: event.CreateAt,
//	}
//}


//func ConvertTaskEventArrToPB (events []*libTypes.TaskEvent) []*pb.TaskEventShow {
//
//	arr := make([]*pb.TaskEventShow, len(events))
//	for i, event := range events {
//		e := &pb.TaskEventShow{
//			Type:     event.Type,
//			GetTaskId:   event.GetTaskId,
//			Owner:    ConvertNodeAliasToPB(event.Owner),
//			Content:  event.Content,
//			CreateAt: event.CreateAt,
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
//			Owner:    ConvertNodeAliasFromPB(event.Owner),
//			Content:  event.Content,
//			CreateAt: event.CreateAt,
//		}
//		arr[i] = e
//	}
//
//	return arr
//}


type ReportTaskEvent struct {
	 PartyId   string
	 Event     *libTypes.TaskEvent
}

func NewReportTaskEvent(partyId string, event *libTypes.TaskEvent) *ReportTaskEvent {
	return &ReportTaskEvent{
		PartyId: partyId,
		Event: event,
	}
}
func (rte *ReportTaskEvent) GetPartyId() string { return rte.PartyId }
func (rte *ReportTaskEvent) GetEvent() *libTypes.TaskEvent { return rte.Event }