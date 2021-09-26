package task

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

var (
	ErrGetNodeTaskList                        = &backend.RpcBizErr{Code: 13001, Msg: "Failed to get all task of current node"}
	ErrGetNodeTaskEventList                   = &backend.RpcBizErr{Code: 13002, Msg: "Failed to get all event of current node's task, taskId: {%s}"}
	ErrSendTaskMsg                            = &backend.RpcBizErr{Code: 13003, Msg: "TerminateTask failed, send task terminate msg failed, taskId: {%s}"}
	ErrReqOperationCostForPublishTask         = &backend.RpcBizErr{Code: 13004, Msg: "Publish task request:Failed to get OperationCost"}
	ErrReqReceiversForPublishTask             = &backend.RpcBizErr{Code: 13005, Msg: "Publish task request:Failed to get Receivers"}
	ErrReqDataSuppliersForPublishTask         = &backend.RpcBizErr{Code: 13006, Msg: "Publish task request:Failed to get DataSuppliers"}
	ErrReqCalculateContractCodeForPublishTask = &backend.RpcBizErr{Code: 13007, Msg: "Publish task request:Failed to get CalculateContractCode"}
	ErrReqMetadataDetailForPublishTask        = &backend.RpcBizErr{Code: 13008, Msg: "Publish task request:Failed to query metadata of partner, identityId: {%s}, metadataId: {%s}"}
	ErrReqMetadataByKeyColumn                 = &backend.RpcBizErr{Code: 13009, Msg: "Publish task request:Not found keyColumn on metadata, identityId: {%s}, metadataId: {%s}, columnIndex: {%d}"}
	ErrReqMetadataBySelectedColumn            = &backend.RpcBizErr{Code: 13010, Msg: "Publish task request:Not found selected column on metadata, identityId: {%s}, metadataId: {%s}, columnIndex: {%d}"}
	ErrSendTaskMsgByTaskId                    = &backend.RpcBizErr{Code: 13011, Msg: "Failed to send taskMsg, send task msg failed, taskId: {%s}"}
	ErrPublishTaskDeclare                     = &backend.RpcBizErr{Code: 13012, Msg: "PublishTaskDeclare failed, query local identity failed, can not publish task"}
	ErrTerminateTask                          = &backend.RpcBizErr{Code: 13013, Msg: "TerminateTask failed, query local identity failed, can not publish task"}
	ErrReqUserTypePublishTask                 = &backend.RpcBizErr{Code: 13014, Msg: "Publish task request:Failed to get UserType"}
	ErrReqUserPublishTask                     = &backend.RpcBizErr{Code: 13015, Msg: "Publish task request:Failed to get User"}
	ErrReqUserSignPublishTask                 = &backend.RpcBizErr{Code: 13016, Msg: "Publish task request:Failed to get UserSign"}
	ErrReqUserTypeTerminateTask               = &backend.RpcBizErr{Code: 13017, Msg: "Terminate task request:Failed to get UserType"}
	ErrReqUserTerminateTask                   = &backend.RpcBizErr{Code: 13018, Msg: "Terminate task request:Failed to get User"}
	ErrReqTaskIdTerminateTask                 = &backend.RpcBizErr{Code: 13019, Msg: "Terminate task request:Failed to get TaskId"}
	ErrReqUserSignTerminateTask               = &backend.RpcBizErr{Code: 13020, Msg: "Terminate task request:Failed to get UserSign"}
)

type Server struct {
	B backend.Backend
}
