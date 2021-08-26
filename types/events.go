package types

import (
	"encoding/json"
	twopb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)

const (
	ApplyIdentity = iota + 1
	RevokeIdentity
	ApplyMetadata
	RevokeMetadata
	ApplyPower
	RevokePower
	ApplyTask
)

type IdentityMsgEvent struct{ Msg *IdentityMsg }
type IdentityRevokeMsgEvent struct{ Msg *IdentityRevokeMsg }
type MetaDataMsgEvent struct{ Msgs MetaDataMsgs }
type MetaDataRevokeMsgEvent struct{ Msgs MetaDataRevokeMsgs }
type PowerMsgEvent struct{ Msgs PowerMsgs }
type PowerRevokeMsgEvent struct{ Msgs PowerRevokeMsgs }
type TaskMsgEvent struct{ Msgs TaskMsgs }

func (msg *IdentityMsgEvent) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *IdentityRevokeMsgEvent) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *MetaDataMsgEvent) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *MetaDataRevokeMsgEvent) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *PowerMsgEvent) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return ""
	}
	return string(result)
}
func (msg *PowerRevokeMsgEvent) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *TaskMsgEvent) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}

func ConvertTaskEvent(event *libTypes.TaskEvent) *twopb.TaskEvent {
	return &twopb.TaskEvent{
		Type:       []byte(event.Type),
		TaskId:     []byte(event.TaskId),
		IdentityId: []byte(event.IdentityId),
		Content:    []byte(event.Content),
		CreateAt:   event.CreateAt,
	}
}

func FetchTaskEvent(event *twopb.TaskEvent) *libTypes.TaskEvent {
	return &libTypes.TaskEvent{
		Type:       string(event.Type),
		TaskId:     string(event.TaskId),
		IdentityId: string(event.IdentityId),
		Content:    string(event.Content),
		CreateAt:   event.CreateAt,
	}
}

func ConvertTaskEventArr(events []*libTypes.TaskEvent) []*twopb.TaskEvent {
	arr := make([]*twopb.TaskEvent, len(events))
	for i, ev := range events {
		arr[i] = ConvertTaskEvent(ev)
	}
	return arr
}

func FetchTaskEventArr(events []*twopb.TaskEvent) []*libTypes.TaskEvent {
	arr := make([]*libTypes.TaskEvent, len(events))
	for i, ev := range events {
		arr[i] = FetchTaskEvent(ev)
	}
	return arr
}

//func ConvertTaskEventToDataCenter(event *libTypes.TaskEvent) *libTypes.TaskEvent {
//	return &libTypes.TaskEvent{
//		TaskId:     event.TaskId,
//		Type:       event.Type,
//		IdentityId: event.IdentityId,
//		Content:    event.Content,
//		CreateAt:   event.CreateAt,
//	}
//}

//func FetchTaskEventFromDataCenter(event *libTypes.TaskEvent) *libTypes.TaskEvent {
//	return &libTypes.TaskEvent{
//		TaskId:     event.TaskId,
//		Type:       event.Type,
//		IdentityId: event.IdentityId,
//		Content:    event.Content,
//		CreateAt:   event.CreateAt,
//	}
//}

//func ConvertTaskEventArrToDataCenter(events []*libTypes.TaskEvent) []*libTypes.TaskEvent {
//	arr := make([]*libTypes.TaskEvent, len(events))
//	for i, ev := range events {
//		arr[i] = ConvertTaskEventToDataCenter(ev)
//	}
//	return arr
//}

//func FetchTaskEventArrFromDataCenter(events []*libTypes.TaskEvent) []*libTypes.TaskEvent {
//	arr := make([]*libTypes.TaskEvent, len(events))
//	for i, ev := range events {
//		arr[i] = FetchTaskEventFromDataCenter(ev)
//	}
//	return arr
//}
