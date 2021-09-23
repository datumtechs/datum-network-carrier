package types

import (
	"encoding/json"
	pb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/taskmng"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)

const (
	ApplyIdentity = iota + 1
	RevokeIdentity
	ApplyMetadata
	RevokeMetadata
	ApplyMetadataAuth
	RevokeMetadataAuth
	ApplyPower
	RevokePower
	ApplyTask
	TerminateTask
)

type IdentityMsgEvent struct{ Msg *IdentityMsg }
type IdentityRevokeMsgEvent struct{ Msg *IdentityRevokeMsg }
type MetadataMsgEvent struct{ Msgs MetadataMsgArr }
type MetadataRevokeMsgEvent struct{ Msgs MetadataRevokeMsgArr }
type MetadataAuthMsgEvent struct{ Msgs MetadataAuthorityMsgArr }
type MetadataAuthRevokeMsgEvent struct{ Msgs MetadataAuthorityRevokeMsgArr }
type PowerMsgEvent struct{ Msgs PowerMsgArr }
type PowerRevokeMsgEvent struct{ Msgs PowerRevokeMsgArr }
type TaskMsgEvent struct{ Msgs TaskMsgArr }
type TaskTerminateMsgEvent struct{ Msgs TaskTerminateMsgArr }

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
func (msg *MetadataMsgEvent) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *MetadataRevokeMsgEvent) String() string {
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

func ConvertTaskEvent(event *libtypes.TaskEvent) *pb.TaskEvent {
	return &pb.TaskEvent{
		Type:       []byte(event.Type),
		TaskId:     []byte(event.TaskId),
		IdentityId: []byte(event.IdentityId),
		Content:    []byte(event.Content),
		CreateAt:   event.CreateAt,
	}
}

func FetchTaskEvent(event *pb.TaskEvent) *libtypes.TaskEvent {
	return &libtypes.TaskEvent{
		Type:       string(event.Type),
		TaskId:     string(event.TaskId),
		IdentityId: string(event.IdentityId),
		Content:    string(event.Content),
		CreateAt:   event.CreateAt,
	}
}

func ConvertTaskEventArr(events []*libtypes.TaskEvent) []*pb.TaskEvent {
	arr := make([]*pb.TaskEvent, len(events))
	for i, ev := range events {
		arr[i] = ConvertTaskEvent(ev)
	}
	return arr
}

func FetchTaskEventArr(events []*pb.TaskEvent) []*libtypes.TaskEvent {
	arr := make([]*libtypes.TaskEvent, len(events))
	for i, ev := range events {
		arr[i] = FetchTaskEvent(ev)
	}
	return arr
}

//func ConvertTaskEventToDataCenter(event *libtypes.TaskEvent) *libtypes.TaskEvent {
//	return &libtypes.TaskEvent{
//		GetTaskId:     event.GetTaskId,
//		Type:       event.Type,
//		IdentityId: event.IdentityId,
//		Content:    event.Content,
//		GetCreateAt:   event.GetCreateAt,
//	}
//}
//
//func FetchTaskEventFromDataCenter(event *libtypes.TaskEvent) *libtypes.TaskEvent {
//	return &libtypes.TaskEvent{
//		GetTaskId:     event.GetTaskId,
//		Type:       event.Type,
//		IdentityId: event.IdentityId,
//		Content:    event.Content,
//		GetCreateAt:   event.GetCreateAt,
//	}
//}
//
//func ConvertTaskEventArrToDataCenter(events []*libtypes.TaskEvent) []*libtypes.TaskEvent {
//	arr := make([]*libtypes.TaskEvent, len(events))
//	for i, ev := range events {
//		arr[i] = ConvertTaskEventToDataCenter(ev)
//	}
//	return arr
//}
//
//func FetchTaskEventArrFromDataCenter(events []*libtypes.TaskEvent) []*libtypes.TaskEvent {
//	arr := make([]*libtypes.TaskEvent, len(events))
//	for i, ev := range events {
//		arr[i] = FetchTaskEventFromDataCenter(ev)
//	}
//	return arr
//}
