package event

import "github.com/RosettaFlow/Carrier-Go/types"

type MetaDataMsgEvent struct{ Msgs types.MetaDataMsgs }
type PowerMsgEvent struct{ Msgs types.PowerMsgs }
type TaskMsgEvent struct{ Msgs types.TaskMsgs }

type TaskEvent struct {
	Type       string `json:"type"`
	Identity   string `json:"identity"`
	TaskId     string `json:"taskId"`
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
}
