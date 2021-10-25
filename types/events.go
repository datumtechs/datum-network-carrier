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
		Type:       []byte(event.GetType()),
		TaskId:     []byte(event.GetTaskId()),
		IdentityId: []byte(event.GetIdentityId()),
		PartyId:    []byte(event.GetPartyId()),
		Content:    []byte(event.GetContent()),
		CreateAt:   event.GetCreateAt(),
	}
}

func FetchTaskEvent(event *pb.TaskEvent) *libtypes.TaskEvent {
	return &libtypes.TaskEvent{
		Type:       string(event.GetType()),
		TaskId:     string(event.GetTaskId()),
		IdentityId: string(event.GetIdentityId()),
		PartyId:    string(event.GetPartyId()),
		Content:    string(event.GetContent()),
		CreateAt:   event.GetCreateAt(),
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