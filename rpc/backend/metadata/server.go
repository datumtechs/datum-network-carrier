package metadata

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

const (
	ErrSendMetaDataRevokeMsgStr = "Failed to send metaDataRevokeMsg"
	ErrGetMetaDataDetailStr     = "Failed to get metadata detail"
	ErrGetMetaDataDetailListStr = "Failed to get metadata detail list"
	ErrSendMetaDataMsgStr       = "Failed to send metaDataMsg"
)

var (
	ErrSendMetaDataRevokeMessage = &backend.RpcBizErr{Msg: "Failed to send metaDataRevokeMsg"}
	ErrGetMetaDataDetail         = &backend.RpcBizErr{Msg: "Failed to get metadata detail"}
	ErrSendIdentityMessage       = &backend.RpcBizErr{Msg: "Failed to get metadata detail list"}
	ErrGetMetaDataDetailList     = &backend.RpcBizErr{Msg: "Failed to send metaDataMsg"}
)

type MetaDataServiceServer struct {
	B backend.Backend
}
