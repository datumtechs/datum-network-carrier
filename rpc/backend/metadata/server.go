package metadata

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

var (
	ErrSendMetadataRevokeMsg = &backend.RpcBizErr{Msg: "Failed to send metaDataRevokeMsg"}
	ErrGetMetadataDetail     = &backend.RpcBizErr{Msg: "Failed to get metadata detail"}
	ErrGetMetadataDetailList = &backend.RpcBizErr{Msg: "Failed to get metadata detail list"}
	ErrSendMetadataMsg       = &backend.RpcBizErr{Msg: "Failed to send metaDataMsg"}
)

type Server struct {
	B backend.Backend
}
