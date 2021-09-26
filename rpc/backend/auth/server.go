package auth

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

var (
	ErrSendIdentityMsg                           = &backend.RpcBizErr{Code: 10001, Msg: "%s, identityId: {%s}, nodeId: {%s}, nodeName: {%s}"}
	ErrSendIdentityRevokeMsg                     = &backend.RpcBizErr{Code: 10002, Msg: "RevokeIdentityJoin failed: %s"}
	ErrGetNodeIdentity                           = &backend.RpcBizErr{Code: 10003, Msg: "Failed to get node identityInfo"}
	ErrGetIdentityList                           = &backend.RpcBizErr{Code: 10004, Msg: "Failed to get all identityInfo list"}
	ErrGetAuthorityList                          = &backend.RpcBizErr{Code: 10005, Msg: "Failed to get all authorityList list"}
	ErrRevokeMetadataAuthority                   = &backend.RpcBizErr{Code: 10006, Msg: "RevokeMetadataAuthority failed, send metadata authority revoke msg failed, userType: {%s}, user: {%s}, MetadataAuthId: {%s}"}
	ErrGetAuthorityListByUserTypeAndUser         = &backend.RpcBizErr{Code: 10007, Msg: "GetMetadataAuthorityListByUser failed, userType: {%s}, user: {%s}"}
	ErrAuditMetadataAuth                         = &backend.RpcBizErr{Code: 10008, Msg: "Failed to audit metadataAuth, metadataAuthId: {%s}, audit option: {%s}, audit suggestion: {%s}"}
	ErrValidMetadataAuthMustCannotExist          = &backend.RpcBizErr{Code: 10009, Msg: "A valid metadata auth must cannot exist, userType: {%s}, user: {%s}, metadataId: {%s}"}
	ErrReqMemberParams                           = &backend.RpcBizErr{Code: 10010, Msg: "Invalid Params, req.Member is nil"}
	ErrReqMemberIdentityIdOrNameParams           = &backend.RpcBizErr{Code: 10011, Msg: "Invalid Params, req.Member.IdentityId or req.Member.Name is empty"}
	ErrReqGetUserForMetadataAuthApply            = &backend.RpcBizErr{Code: 10012, Msg: "Initiating an authorization application for using metadata:Failed to get use information"}
	ErrVerifyUserTypeForMetadataAuthApply        = &backend.RpcBizErr{Code: 10013, Msg: "Initiating an authorization application for using metadata:Verify that the user type is incorrect"}
	ErrReqAuthForMetadataAuthApply               = &backend.RpcBizErr{Code: 10014, Msg: "Initiating an authorization application for using metadata:Failed to get the metadata permission"}
	ErrReqUserSignForMetadataAuthApply           = &backend.RpcBizErr{Code: 10015, Msg: "Initiating an authorization application for using metadata:Failed to get the user signature"}
	ErrApplyMetadataAuthority                    = &backend.RpcBizErr{Code: 10016, Msg: "ApplyMetadataAuthority failed, %s, userType: {%s}, user: {%s}, metadataId: {%s}"}
	ErrValidAuditMetadataOptionMustCannotPending = &backend.RpcBizErr{Code: 10017, Msg: "A valid audit metadata option must cannot pending"}
	ErrReqGetUserForRevokeMetadataAuth           = &backend.RpcBizErr{Code: 10018, Msg: "Revoke Metadata Authority:Failed to get use information"}
	ErrVerifyUserTypeForRevokeMetadataAuth       = &backend.RpcBizErr{Code: 10019, Msg: "Revoke Metadata Authority:Verify that the user type is incorrect"}
	ErrReqAuthIDForRevokeMetadataAuth            = &backend.RpcBizErr{Code: 10020, Msg: "Revoke Metadata Authority:Failed to get the metadata permission ID"}
	ErrReqUserSignForRevokeMetadataAuth          = &backend.RpcBizErr{Code: 10021, Msg: "Revoke Metadata Authority:Failed to get the user signature"}
	ErrReqAuthIDForAuditMetadataAuth             = &backend.RpcBizErr{Code: 10022, Msg: "Organize the review of users' data authorization requests:Failed to get the metadata permission ID"}
	ErrReqGetUserForMetadataAuthListByUser       = &backend.RpcBizErr{Code: 10023, Msg: "Request for a list of authorization requests and audit results for all current metadata:Failed to get use information"}
	ErrVerifyUserTypeForMetadataAuthListByUser   = &backend.RpcBizErr{Code: 10024, Msg: "Request for a list of authorization requests and audit results for all current metadata:Verify that the user type is incorrect"}
)

type Server struct {
	B backend.Backend
}
