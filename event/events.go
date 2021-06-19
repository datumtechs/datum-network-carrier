package event

import "github.com/RosettaFlow/Carrier-Go/types"

type IdentityMsgEvent struct{ Msg *types.IdentityMsg }
type IdentityRevokeMsgEvent struct{ Msg *types.IdentityRevokeMsg }
type MetaDataMsgEvent struct{ Msgs types.MetaDataMsgs }
type MetaDataRevokeMsgEvent struct{ Msgs types.MetaDataRevokeMsgs }
type PowerMsgEvent struct{ Msgs types.PowerMsgs }
type PowerRevokeMsgEvent struct{ Msgs types.PowerRevokeMsgs }
type TaskMsgEvent struct{ Msgs types.TaskMsgs }

type TaskEvent struct {
	Type       string `json:"type"`
	Identity   string `json:"identity"`
	TaskId     string `json:"taskId"`
	Content    string `json:"content"`
	CreateTime uint64 `json:"createTime"`
}
