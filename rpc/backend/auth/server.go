package auth

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

const (
	ErrSendIdentityMsgStr       = "Failed to send identityMsg"
	ErrSendIdentityRevokeMsgStr = "Failed to send identityRevokeMsg"
	ErrGetNodeIdentityStr       = "Failed to get node identityInfo"
	ErrGetIdentityListStr       = "Failed to get all identityInfo list"
)

var (
	ErrSendIdentityMessage       = &backend.RpcBizErr{Msg: "Failed to send identityMsg"}
	ErrSendIdentityRevokeMessage = &backend.RpcBizErr{Msg: "Failed to send identityRevokeMsg"}
	ErrGetNodeIdentity           = &backend.RpcBizErr{Msg: "Failed to get node identityInfo"}
	ErrGetIdentityList           = &backend.RpcBizErr{Msg: "Failed to get all identityInfo list"}
)

type AuthServiceServer struct {
	B backend.Backend
}
