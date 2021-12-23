package auth

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

var (
	ErrSendIdentityMsg                           = &backend.RpcBizErr{Code: 10001, Msg: "ApplyIdentityJoin failed"}
	ErrSendIdentityRevokeMsg                     = &backend.RpcBizErr{Code: 10002, Msg: "RevokeIdentityJoin failed"}
	ErrGetNodeIdentity                           = &backend.RpcBizErr{Code: 10003, Msg: "Failed to get node identityInfo"}
	ErrGetIdentityList                           = &backend.RpcBizErr{Code: 10004, Msg: "Failed to get all identityInfo list"}
	ErrGetAuthorityList                          = &backend.RpcBizErr{Code: 10005, Msg: "Failed to get all authority list"}
	ErrRevokeMetadataAuthority                   = &backend.RpcBizErr{Code: 10006, Msg: "RevokeMetadataAuthority failed, send metadata authority revoke msg failed"}
	ErrGetAuthorityListByUserTypeAndUser         = &backend.RpcBizErr{Code: 10007, Msg: "GetMetadataAuthorityListByUser failed"}
	ErrAuditMetadataAuth                         = &backend.RpcBizErr{Code: 10008, Msg: "Failed to audit metadataAuth"}
	ErrValidMetadataAuthMustCannotExist          = &backend.RpcBizErr{Code: 10009, Msg: "A valid metadata auth must cannot exist"}
	ErrReqMemberParams                           = &backend.RpcBizErr{Code: 10010, Msg: "Invalid Params, req.Member is nil"}
	ErrReqMemberIdentityIdOrNameParams           = &backend.RpcBizErr{Code: 10011, Msg: "Invalid Params, req.Member.IdentityId or req.Member.Name is empty"}
	ErrReqGetUserForMetadataAuthApply            = &backend.RpcBizErr{Code: 10012, Msg: "Initiating an authorization application for using metadata: the user is empty"}
	ErrVerifyUserTypeForMetadataAuthApply        = &backend.RpcBizErr{Code: 10013, Msg: "Initiating an authorization application for using metadata: verify that the user type is incorrect"}
	ErrReqAuthForMetadataAuthApply               = &backend.RpcBizErr{Code: 10014, Msg: "Initiating an authorization application for using metadata: the metadata auth is nil"}
	ErrReqUserSignForMetadataAuthApply           = &backend.RpcBizErr{Code: 10015, Msg: "Initiating an authorization application for using metadata: the user signature is nil"}
	ErrApplyMetadataAuthority                    = &backend.RpcBizErr{Code: 10016, Msg: "ApplyMetadataAuthority failed"}
	ErrValidAuditMetadataOptionMustCannotPending = &backend.RpcBizErr{Code: 10017, Msg: "A valid audit metadata option must cannot pending"}
	ErrReqGetUserForRevokeMetadataAuth           = &backend.RpcBizErr{Code: 10018, Msg: "Revoke Metadata Authority: the user is empty"}
	ErrVerifyUserTypeForRevokeMetadataAuth       = &backend.RpcBizErr{Code: 10019, Msg: "Revoke Metadata Authority: verify that the user type is incorrect"}
	ErrReqAuthIDForRevokeMetadataAuth            = &backend.RpcBizErr{Code: 10020, Msg: "Revoke Metadata Authority: the metadata auth is nil"}
	ErrReqUserSignForRevokeMetadataAuth          = &backend.RpcBizErr{Code: 10021, Msg: "Revoke Metadata Authority: the user signature is nil"}
	ErrReqAuthIDForAuditMetadataAuth             = &backend.RpcBizErr{Code: 10022, Msg: "Organize the review of users' data authorization requests: the metadataAuth Id is empty"}
	ErrReqGetUserForMetadataAuthListByUser       = &backend.RpcBizErr{Code: 10023, Msg: "Request for a list of authorization requests and audit results for all current metadata: the user is empty"}
	ErrVerifyUserTypeForMetadataAuthListByUser   = &backend.RpcBizErr{Code: 10024, Msg: "Request for a list of authorization requests and audit results for all current metadata: verify that the user type is incorrect"}
)

type Server struct {
	B backend.Backend
}
