package task

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

const (
	ErrGetNodeTaskListStr      = "Failed to get all task of current node"
	ErrGetNodeTaskEventListStr = "Failed to get all event of current node's task"
	ErrSendTaskMsgStr          = "Failed to send taskMsg"
)

var (
	ErrGetNodeTaskList      = &backend.RpcBizErr{Msg: "Failed to get all task of current node"}
	ErrGetNodeTaskEventList = &backend.RpcBizErr{Msg: "Failed to get all event of current node's task"}
	ErrSendTaskMsg          = &backend.RpcBizErr{Msg: "Failed to send taskMsg"}
)

type TaskServiceServer struct {
	B backend.Backend
}
