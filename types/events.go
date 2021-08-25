package types

import (
	"encoding/json"
	"fmt"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
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

type TaskEventInfo struct {
	Type       string `json:"type"`
	Identity   string `json:"identity"`
	TaskId     string `json:"taskId"`
	Content    string `json:"content"`
	CreateTime uint64 `json:"createTime"`
}

func (ev *TaskEventInfo) String() string {
	return fmt.Sprintf(`{"taskId": %s, "identity", %s, type": %s, "content": %s, "createTime": %d}`,
		ev.TaskId, ev.Identity, ev.Type, ev.Content, ev.CreateTime)
}

func (msg *IdentityMsgEvent) String() string {
	result, err := json.Marshal(msg)
	if err != nil{
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *IdentityRevokeMsgEvent) String() string {
	result, err := json.Marshal(msg)
	if err != nil{
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *MetaDataMsgEvent) String() string {
	result, err := json.Marshal(msg)
	if err != nil{
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *MetaDataRevokeMsgEvent) String() string {
	result, err := json.Marshal(msg)
	if err != nil{
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *PowerMsgEvent) String() string {
	result, err := json.Marshal(msg)
	if err != nil{
		return ""
	}
	return string(result)
}
func (msg *PowerRevokeMsgEvent) String() string {
	result, err := json.Marshal(msg)
	if err != nil{
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *TaskMsgEvent) String() string {
	result, err := json.Marshal(msg)
	if err != nil{
		return "Failed to generate string"
	}
	return string(result)
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


func ConvertTaskEventToDataCenter(event *TaskEventInfo) *libTypes.TaskEvent {
	return &libTypes.TaskEvent{
		TaskId       :event.TaskId,
		Type    :event.Type,
		IdentityId     :event.Identity,
		Content :event.Content,
		CreateAt      :event.CreateTime,
	}
}

func FetchTaskEventFromDataCenter(event *libTypes.TaskEvent) *TaskEventInfo {
	return &TaskEventInfo{
		TaskId: event.TaskId,
		Type: event.Type,
		Identity: event.IdentityId,
		Content: event.Content,
		CreateTime: event.CreateAt,
	}
}

func ConvertTaskEventArrToDataCenter(events []*TaskEventInfo) []*libTypes.TaskEvent {
	arr := make([]*libTypes.TaskEvent, len(events))
	for i, ev := range events {
		arr[i] = ConvertTaskEventToDataCenter(ev)
	}
	return arr
}

func FetchTaskEventArrFromDataCenter(events []*libTypes.TaskEvent) []*TaskEventInfo {
	arr := make([]*TaskEventInfo, len(events))
	for i, ev := range events {
		arr[i] = FetchTaskEventFromDataCenter(ev)
	}
	return arr
}