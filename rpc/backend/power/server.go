package power

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

var (
	ErrSendPowerRevokeMsg = &backend.RpcBizErr{Msg: "Failed to send powerRevokeMsg"}
	ErrGetTotalPowerList  = &backend.RpcBizErr{Msg: "Failed to get total power list"}
	ErrGetSinglePowerList = &backend.RpcBizErr{Msg: "Failed to get current node power list"}
	ErrSendPowerMsg       = &backend.RpcBizErr{Msg: "Failed to send powerMsg"}
)

type PowerServiceServer struct {
	B backend.Backend
}
