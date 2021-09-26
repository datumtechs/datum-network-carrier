package power

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

var (
	ErrSendPowerRevokeMsg            = &backend.RpcBizErr{Code: 12001, Msg: "Failed to send powerRevokeMsg: query local identity failed, can not revoke power"}
	ErrGetTotalPowerList             = &backend.RpcBizErr{Code: 12002, Msg: "Failed to get total power list"}
	ErrGetSinglePowerList            = &backend.RpcBizErr{Code: 12003, Msg: "Failed to get current node power list"}
	ErrSendPowerMsg                  = &backend.RpcBizErr{Code: 12004, Msg: "Failed to send powerMsg: query local identity failed, can not publish power"}
	ErrSendPowerMsgByNidAndPowerId   = &backend.RpcBizErr{Code: 12005, Msg: "Failed to send powerMsg, jobNodeId:{%s}, powerId:{%s}"}
	ErrReqEmptyForPublishPower       = &backend.RpcBizErr{Code: 12006, Msg: "Publish Power Request: request is empty, enabling failed"}
	ErrReqEmptyForRevokePower        = &backend.RpcBizErr{Code: 12007, Msg: "Revoke Power Request: request is empty, revoke failed"}
	ErrReqEmptyPowerIdForRevokePower = &backend.RpcBizErr{Code: 12008, Msg: "Revoke Power Request: Power ID is empty, revoke failed"}
	ErrSendPowerRevokeMsgByPowerId   = &backend.RpcBizErr{Code: 12009, Msg: "Failed to send powerRevokeMsg, powerId:{%s}"}
)

type Server struct {
	B backend.Backend
}
