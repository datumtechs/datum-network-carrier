package types

import (
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)

type IdentityMsgEvent struct{ Msg *IdentityMsg }
type IdentityRevokeMsgEvent struct{ Msg *IdentityRevokeMsg }
type MetaDataMsgEvent struct{ Msgs MetaDataMsgs }
type MetaDataRevokeMsgEvent struct{ Msgs MetaDataRevokeMsgs }
type PowerMsgEvent struct{ Msgs PowerMsgs }
type PowerRevokeMsgEvent struct{ Msgs PowerRevokeMsgs }
type TaskMsgEvent struct{ Msgs TaskMsgs }

type TaskEventInfo struct {
	Type       string `json:"type"`
	Identity   string `json:"identity"`
	TaskId     string `json:"taskId"`
	Content    string `json:"content"`
	CreateTime uint64 `json:"createTime"`
}

func ConvertTaskEvent(event *TaskEventInfo) *pb.TaskEvent {
	return &pb.TaskEvent{
		Type: []byte(event.Type),
		TaskId: []byte(event.TaskId),
		IdentityId: []byte(event.Identity),
		Content: []byte(event.Content),
		CreateAt: event.CreateTime,
	}
}

func FetchTaskEvent(event *pb.TaskEvent) *TaskEventInfo {
	return &TaskEventInfo{
		Type: string(event.Type),
		TaskId: string(event.TaskId),
		Identity: string(event.IdentityId),
		Content: string(event.Content),
		CreateTime: event.CreateAt,
	}
}

func ConvertTaskEventArr(events []*TaskEventInfo) []*pb.TaskEvent {
	arr := make([]*pb.TaskEvent, len(events))
	for i, ev := range events {
		arr[i] = ConvertTaskEvent(ev)
	}
	return arr
}

func FetchTaskEventArr(events []*pb.TaskEvent) []*TaskEventInfo {
	arr := make([]*TaskEventInfo, len(events))
	for i, ev := range events {
		arr[i] = FetchTaskEvent(ev)
	}
	return arr
}


func ConvertTaskEventToDataCenter(event *TaskEventInfo) *libTypes.EventData {
	return &libTypes.EventData{
		TaskId       :event.TaskId,
		EventType    :event.Type,
		Identity     :event.Identity,
		EventContent :event.Content,
		EventAt      :event.CreateTime,
	}
}

func FetchTaskEventFromDataCenter(event *libTypes.EventData) *TaskEventInfo {
	return &TaskEventInfo{
		TaskId: event.TaskId,
		Type: event.EventType,
		Identity: event.Identity,
		Content: event.EventContent,
		CreateTime: event.EventAt,
	}
}

func ConvertTaskEventArrToDataCenter(events []*TaskEventInfo) []*libTypes.EventData {
	arr := make([]*libTypes.EventData, len(events))
	for i, ev := range events {
		arr[i] = ConvertTaskEventToDataCenter(ev)
	}
	return arr
}

func FetchTaskEventArrFromDataCenter(events []*libTypes.EventData) []*TaskEventInfo {
	arr := make([]*TaskEventInfo, len(events))
	for i, ev := range events {
		arr[i] = FetchTaskEventFromDataCenter(ev)
	}
	return arr
}