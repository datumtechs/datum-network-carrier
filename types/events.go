package types

import (
	"encoding/json"
	carriernetmsgtaskmngpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/taskmng"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
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
type MetadataMsgEvent struct{ Msg *MetadataMsg }
type MetadataRevokeMsgEvent struct{ Msg *MetadataRevokeMsg }
type MetadataAuthMsgEvent struct{ Msg *MetadataAuthorityMsg }
type MetadataAuthRevokeMsgEvent struct{ Msg *MetadataAuthorityRevokeMsg }
type PowerMsgEvent struct{ Msg *PowerMsg }
type PowerRevokeMsgEvent struct{ Msg *PowerRevokeMsg }
type TaskMsgEvent struct{ Msg *TaskMsg }
type TaskTerminateMsgEvent struct{ Msg *TaskTerminateMsg }

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

func ConvertTaskEvent(event *carriertypespb.TaskEvent) *carriernetmsgtaskmngpb.TaskEvent {
	return &carriernetmsgtaskmngpb.TaskEvent{
		Type:       []byte(event.GetType()),
		TaskId:     []byte(event.GetTaskId()),
		IdentityId: []byte(event.GetIdentityId()),
		PartyId:    []byte(event.GetPartyId()),
		Content:    []byte(event.GetContent()),
		CreateAt:   event.GetCreateAt(),
	}
}

func FetchTaskEvent(event *carriernetmsgtaskmngpb.TaskEvent) *carriertypespb.TaskEvent {
	return &carriertypespb.TaskEvent{
		Type:       string(event.GetType()),
		TaskId:     string(event.GetTaskId()),
		IdentityId: string(event.GetIdentityId()),
		PartyId:    string(event.GetPartyId()),
		Content:    string(event.GetContent()),
		CreateAt:   event.GetCreateAt(),
	}
}

func ConvertTaskEventArr(events []*carriertypespb.TaskEvent) []*carriernetmsgtaskmngpb.TaskEvent {
	arr := make([]*carriernetmsgtaskmngpb.TaskEvent, len(events))
	for i, ev := range events {
		arr[i] = ConvertTaskEvent(ev)
	}
	return arr
}

func FetchTaskEventArr(events []*carriernetmsgtaskmngpb.TaskEvent) []*carriertypespb.TaskEvent {
	arr := make([]*carriertypespb.TaskEvent, len(events))
	for i, ev := range events {
		arr[i] = FetchTaskEvent(ev)
	}
	return arr
}