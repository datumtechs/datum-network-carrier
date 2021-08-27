package task

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

var (
	ErrGetNodeTaskList      = &backend.RpcBizErr{Msg: "Failed to get all task of current node"}
	ErrGetNodeTaskEventList = &backend.RpcBizErr{Msg: "Failed to get all event of current node's task"}
	ErrSendTaskMsg          = &backend.RpcBizErr{Msg: "Failed to send taskMsg"}
)

type Server struct {
	B backend.Backend
}
