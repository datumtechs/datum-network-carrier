package auth

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

var (
	ErrSendIdentityMsg                  = &backend.RpcBizErr{Msg: "Failed to send identityMsg"}
	ErrSendIdentityRevokeMsg            = &backend.RpcBizErr{Msg: "Failed to send identityRevokeMsg"}
	ErrGetNodeIdentity                  = &backend.RpcBizErr{Msg: "Failed to get node identityInfo"}
	ErrGetIdentityList                  = &backend.RpcBizErr{Msg: "Failed to get all identityInfo list"}
	ErrGetAuthorityList                 = &backend.RpcBizErr{Msg: "Failed to get all authorityList list"}
	ErrSendMetadataAuthMsg              = &backend.RpcBizErr{Msg: "Failed to send metadataAuthMsg"}
	ErrSendMetadataAuthRevokeMsg        = &backend.RpcBizErr{Msg: "Failed to send metadataAuthRevokeMsg"}
	ErrAuditMetadataAuth                = &backend.RpcBizErr{Msg: "Failed to audit metadataAuth"}
	ErrValidMetadataAuthMustCannotExist = &backend.RpcBizErr{Msg: "A valid metadata auth must cannot exist"}
)

type Server struct {
	B backend.Backend
}
