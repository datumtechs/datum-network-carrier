package types

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
