package types

import pb "github.com/RosettaFlow/Carrier-Go/lib/api"
type TaskEvent struct {
	TaskId   string     `json:"taskId"`
	Type     string     `json:"type"`
	CreateAt uint64     `json:"createAt"`
	Content  string     `json:"content"`
	Owner    *NodeAlias `json:"owner"`
}
func ConvertTaskEventToPB (event *TaskEvent) *pb.TaskEventShow {
	return &pb.TaskEventShow{
		Type: event.Type,
		TaskId: event.TaskId,
		Owner: ConvertNodeAliasToPB(event.Owner),
		Content: event.Content,
		CreateAt: event.CreateAt,
	}
}

func ConvertTaskEventFromPB (event *pb.TaskEventShow) *TaskEvent {
	return &TaskEvent{
		Type: event.Type,
		TaskId: event.TaskId,
		Owner: ConvertNodeAliasFromPB(event.Owner),
		Content: event.Content,
		CreateAt: event.CreateAt,
	}
}


func ConvertTaskEventArrToPB (events []*TaskEvent) []*pb.TaskEventShow {

	arr := make([]*pb.TaskEventShow, len(events))
	for i, event := range events {
		e := &pb.TaskEventShow{
			Type: event.Type,
			TaskId: event.TaskId,
			Owner: ConvertNodeAliasToPB(event.Owner),
			Content: event.Content,
			CreateAt: event.CreateAt,
		}
		arr[i] = e
	}

	return arr
}

func ConvertTaskEventArrFromPB (events []*pb.TaskEventShow) []*TaskEvent {
	arr := make([]*TaskEvent, len(events))
	for i, event := range events {
		e := &TaskEvent{
			Type: event.Type,
			TaskId: event.TaskId,
			Owner: ConvertNodeAliasFromPB(event.Owner),
			Content: event.Content,
			CreateAt: event.CreateAt,
		}
		arr[i] = e
	}

	return arr
}