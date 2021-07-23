package auth

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

var (
	ErrSendIdentityMsg       = &backend.RpcBizErr{Msg: "Failed to send identityMsg"}
	ErrSendIdentityRevokeMsg = &backend.RpcBizErr{Msg: "Failed to send identityRevokeMsg"}
	ErrGetNodeIdentity       = &backend.RpcBizErr{Msg: "Failed to get node identityInfo"}
	ErrGetIdentityList       = &backend.RpcBizErr{Msg: "Failed to get all identityInfo list"}
)

type AuthServiceServer struct {
	B backend.Backend
}
