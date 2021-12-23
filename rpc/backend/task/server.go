package task

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

var (
	ErrGetNodeTaskList                        = &backend.RpcBizErr{Code: 13001, Msg: "Failed to get all task of current node"}
	ErrGetNodeTaskEventList                   = &backend.RpcBizErr{Code: 13002, Msg: "Failed to get all event of current node's task"}
	ErrTerminateTaskMsg                       = &backend.RpcBizErr{Code: 13003, Msg: "TerminateTask failed, send task terminate msg failed"}
	ErrReqOperationCostForPublishTask         = &backend.RpcBizErr{Code: 13004, Msg: "Publish task request: the OperationCost is nil"}
	ErrReqReceiversForPublishTask             = &backend.RpcBizErr{Code: 13005, Msg: "Publish task request: the Receivers is nil"}
	ErrReqDataSuppliersForPublishTask         = &backend.RpcBizErr{Code: 13006, Msg: "Publish task request: the DataSuppliers is nil"}
	ErrReqCalculateContractCodeForPublishTask = &backend.RpcBizErr{Code: 13007, Msg: "Publish task request: the CalculateContractCode is empty"}
	ErrReqMetadataDetailForPublishTask        = &backend.RpcBizErr{Code: 13008, Msg: "Publish task request: query metadata of partner failed"}
	ErrReqMetadataByKeyColumn                 = &backend.RpcBizErr{Code: 13009, Msg: "Publish task request: not found keyColumn on metadata"}
	ErrReqMetadataBySelectedColumn            = &backend.RpcBizErr{Code: 13010, Msg: "Publish task request: not found selected column on metadata"}
	ErrSendTaskMsgByTaskId                    = &backend.RpcBizErr{Code: 13011, Msg: "Failed to send taskMsg, send task msg failed"}
	ErrPublishTaskDeclare                     = &backend.RpcBizErr{Code: 13012, Msg: "PublishTaskDeclare failed"}
	ErrTerminateTask                          = &backend.RpcBizErr{Code: 13013, Msg: "TerminateTask failed"}
	ErrReqUserTypePublishTask                 = &backend.RpcBizErr{Code: 13014, Msg: "Publish task request: the userType is empty"}
	ErrReqUserPublishTask                     = &backend.RpcBizErr{Code: 13015, Msg: "Publish task request: the user is empty"}
	ErrReqUserSignPublishTask                 = &backend.RpcBizErr{Code: 13016, Msg: "Publish task request: the userSign is empty"}
	ErrReqUserTypeTerminateTask               = &backend.RpcBizErr{Code: 13017, Msg: "Terminate task request: the userType is empty"}
	ErrReqUserTerminateTask                   = &backend.RpcBizErr{Code: 13018, Msg: "Terminate task request: the user is empty"}
	ErrReqTaskIdTerminateTask                 = &backend.RpcBizErr{Code: 13019, Msg: "Terminate task request: the taskId is empty"}
	ErrReqUserSignTerminateTask               = &backend.RpcBizErr{Code: 13020, Msg: "Terminate task request: the userSign is empty"}
)

type Server struct {
	B backend.Backend
}
