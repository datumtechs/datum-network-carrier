package metadata

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

var (
	ErrSendMetaDataRevokeMsg = &backend.RpcBizErr{Msg: "Failed to send metaDataRevokeMsg"}
	ErrGetMetaDataDetail     = &backend.RpcBizErr{Msg: "Failed to get metadata detail"}
	ErrGetMetaDataDetailList = &backend.RpcBizErr{Msg: "Failed to get metadata detail list"}
	ErrSendMetaDataMsg       = &backend.RpcBizErr{Msg: "Failed to send metaDataMsg"}
)

type Server struct {
	B backend.Backend
}
