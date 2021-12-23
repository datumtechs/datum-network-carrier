package auth

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

// for organization identity

func (svr *Server) ApplyIdentityJoin(ctx context.Context, req *pb.ApplyIdentityJoinRequest) (*apicommonpb.SimpleResponse, error) {

	identity, err := svr.B.GetNodeIdentity()
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("RPC-API:ApplyIdentityJoin failed, query local identity failed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}",
			req.GetMember().GetIdentityId(), req.GetMember().GetNodeId(), req.GetMember().GetNodeName())

		errMsg := fmt.Sprintf("query local identity failed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}",
			req.GetMember().GetIdentityId(), req.GetMember().GetNodeId(), req.GetMember().GetNodeName())
		return nil, backend.NewRpcBizErr(ErrSendIdentityMsg.Code, errMsg)
	}

	if nil != identity {
		log.Errorf("RPC-API:ApplyIdentityJoin failed, identity was already exist, old identityId: {%s}, old nodeId: {%s}, old nodeName: {%s}",
			identity.GetIdentityId(), identity.GetNodeId(), identity.GetName())

		errMsg := fmt.Sprintf("identity was already exist, identityId: {%s}, nodeId: {%s}, nodeName: {%s}",
			identity.GetIdentityId(), identity.GetNodeId(), identity.GetName())
		return nil, backend.NewRpcBizErr(ErrSendIdentityMsg.Code, errMsg)
	}

	if req.GetMember() == nil {
		return nil, ErrReqMemberParams
	}

	if "" == strings.Trim(req.GetMember().GetIdentityId(), "") ||
		"" == strings.Trim(req.GetMember().GetNodeName(), "") {
		return nil, ErrReqMemberIdentityIdOrNameParams
	}

	identityMsg := types.NewIdentityMessageFromRequest(req)
	err = svr.B.SendMsg(identityMsg)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ApplyIdentityJoin failed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}",
			req.GetMember().GetIdentityId(), req.GetMember().GetNodeId(), req.GetMember().GetNodeName())

		errMsg := fmt.Sprintf("%s, identityId: {%s}, nodeId: {%s}, nodeName: {%s}", ErrSendIdentityMsg.Msg,
			req.GetMember().GetIdentityId(), req.GetMember().GetNodeId(), req.GetMember().GetNodeName())
		return nil, backend.NewRpcBizErr(ErrSendIdentityMsg.Code, errMsg)
	}
	log.Debugf("RPC-API:ApplyIdentityJoin succeed SendMsg, identityId: {%s}, nodeId: {%s}, nodeName: {%s}",
		req.GetMember().GetIdentityId(), req.GetMember().GetNodeId(), req.GetMember().GetNodeName())
	return &apicommonpb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) RevokeIdentityJoin(ctx context.Context, req *emptypb.Empty) (*apicommonpb.SimpleResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if rawdb.IsDBNotFoundErr(err) {
		log.WithError(err).Errorf("RPC-API:RevokeIdentityJoin failed, the identity was not exist, can not revoke identity")

		errMsg := fmt.Sprintf("%s, the identity was not exist, can not revoke identity", ErrSendIdentityRevokeMsg.Msg)
		return nil, backend.NewRpcBizErr(ErrSendIdentityRevokeMsg.Code, errMsg)
	}

	// what if local task we can not revoke identity
	has, err := svr.B.HasLocalTask ()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokeIdentityJoin failed, can not check has local task")

		errMsg := fmt.Sprintf("%s, can not check has local task", ErrSendIdentityRevokeMsg.Msg)
		return nil, backend.NewRpcBizErr(ErrSendIdentityRevokeMsg.Code, errMsg)
	}

	if has {
		log.WithError(err).Errorf("RPC-API:RevokeIdentityJoin failed, don't revoke identity when has local task")

		errMsg := fmt.Sprintf("%s, don't revoke identity when has local task", ErrSendIdentityRevokeMsg.Msg)
		return nil, backend.NewRpcBizErr(ErrSendIdentityRevokeMsg.Code, errMsg)
	}

	identityRevokeMsg := types.NewIdentityRevokeMessage()
	err = svr.B.SendMsg(identityRevokeMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:RevokeIdentityJoin failed")

		errMsg := fmt.Sprintf("%s, send identity revoke msg failed", ErrSendIdentityRevokeMsg.Msg)
		return nil, backend.NewRpcBizErr(ErrSendIdentityRevokeMsg.Code, errMsg)
	}
	log.Debug("RPC-API:RevokeIdentityJoin succeed SendMsg")
	return &apicommonpb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) GetNodeIdentity(ctx context.Context, req *emptypb.Empty) (*pb.GetNodeIdentityResponse, error) {
	identity, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetNodeIdentity failed")
		return nil, ErrGetNodeIdentity
	}
	return &pb.GetNodeIdentityResponse{
		Status: 0,
		Msg:    backend.OK,
		Owner: &apicommonpb.Organization{
			NodeName:   identity.GetName(),
			NodeId:     identity.GetNodeId(),
			IdentityId: identity.GetIdentityId(),
		},
	}, nil
}

func (svr *Server) GetIdentityList(ctx context.Context, req *pb.GetIdentityListRequest) (*pb.GetIdentityListResponse, error) {
	identityList, err := svr.B.GetIdentityList(req.GetLastUpdated(), backend.DefaultPageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:QueryIdentityList failed")
		return nil, ErrGetIdentityList
	}
	arr := make([]*apicommonpb.Organization, len(identityList))
	for i, identity := range identityList {
		iden := &apicommonpb.Organization{
			NodeName:   identity.GetName(),
			NodeId:     identity.GetNodeId(),
			IdentityId: identity.GetIdentityId(),
		}
		arr[i] = iden
	}
	log.Debugf("Query all org's identity list, len: {%d}", len(identityList))
	return &pb.GetIdentityListResponse{
		Status:     0,
		Msg:        backend.OK,
		MemberList: arr,
	}, nil
}

// for metadata authority apply

func (svr *Server) ApplyMetadataAuthority(ctx context.Context, req *pb.ApplyMetadataAuthorityRequest) (*pb.ApplyMetadataAuthorityResponse, error) {
	if req.GetUser() == "" {
		return nil, ErrReqGetUserForMetadataAuthApply
	}
	if !verifyUserType(req.GetUserType()) {
		return nil, ErrVerifyUserTypeForMetadataAuthApply
	}
	if req.GetAuth() == nil {
		return nil, ErrReqAuthForMetadataAuthApply
	}
	if len(req.GetSign()) == 0 {
		return nil, ErrReqUserSignForMetadataAuthApply
	}

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ApplyMetadataAuthority failed, query local identity failed")
		return nil, fmt.Errorf("query local identity failed")
	}

	now := timeutils.UnixMsecUint64()
	switch req.GetAuth().GetUsageRule().GetUsageType() {
	case apicommonpb.MetadataUsageType_Usage_Period:
		if now >= req.GetAuth().GetUsageRule().GetEndAt() {
			log.Errorf("RPC-API:ApplyMetadataAuthority failed, usaageRule endTime of metadataAuth has expire, userType: {%s}, user: {%s}, metadataId: {%s}, usageType: {%s}, usageEndTime: {%d}, now: {%d}",
				req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId(), req.GetAuth().GetUsageRule().GetUsageType().String(), req.GetAuth().GetUsageRule().GetEndAt(), now)
			return nil, fmt.Errorf("usaageRule endTime of metadataAuth has expire")
		}
	case apicommonpb.MetadataUsageType_Usage_Times:
		if req.GetAuth().GetUsageRule().GetTimes() == 0 {
			log.Errorf("RPC-API:ApplyMetadataAuthority failed, usaageRule times of metadataAuth must be greater than zero, userType: {%s}, user: {%s}, metadataId: {%s}, usageType: {%s}, usageEndTime: {%d}, now: {%d}",
				req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId(), req.GetAuth().GetUsageRule().GetUsageType().String(), req.GetAuth().GetUsageRule().GetEndAt(), now)
			return nil, fmt.Errorf("usaageRule times of metadataAuth must be greater than zero")
		}
	default:
		log.Errorf("RPC-API:ApplyMetadataAuthority failed, unknown usageType of the metadataAuth, userType: {%s}, user: {%s}, metadataId: {%s}, usageType: {%s}",
			req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId(), req.GetAuth().GetUsageRule().GetUsageType().String())
		return nil, fmt.Errorf("unknown usageType of the metadataAuth")
	}

	// ############################################
	// ############################################
	// NOTE:
	// 		check metadataId and identity of authority
	// ############################################
	// ############################################
	ideneityList, err := svr.B.GetIdentityList(timeutils.BeforeYearUnixMsecUint64(), backend.DefaultMaxPageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:ApplyMetadataAuthority failed, query global identity list failed")
		return nil, backend.NewRpcBizErr(ErrApplyMetadataAuthority.Code, "query global identity list failed")
	}
	var valid bool // false
	for _, identity := range ideneityList {
		if identity.GetIdentityId() == req.GetAuth().GetOwner().GetIdentityId() {
			valid = true
			break
		}
	}
	if !valid {
		log.WithError(err).Errorf("RPC-API:ApplyMetadataAuthority failed, not found identity with identityId of auth, identityId: {%s}",
			req.GetAuth().GetOwner().GetIdentityId())
		return nil, backend.NewRpcBizErr(ErrApplyMetadataAuthority.Code, "not found identity with identityId of auth")
	}
	// reset val
	valid = false

	metadataList, err := svr.B.GetGlobalMetadataDetailList(timeutils.BeforeYearUnixMsecUint64(), backend.DefaultMaxPageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:ApplyMetadataAuthority failed, query global metadata list failed")
		return nil, backend.NewRpcBizErr(ErrApplyMetadataAuthority.Code, "query global metadata list failed")
	}
	for _, metadata := range metadataList {
		if metadata.GetInformation().GetMetadataSummary().GetMetadataId() == req.GetAuth().GetMetadataId() {
			valid = true
			break
		}
	}
	if !valid {
		log.WithError(err).Errorf("RPC-API:ApplyMetadataAuthority failed, not found metadata with metadataId of auth, metadataId: {%s}",
			req.GetAuth().GetMetadataId())
		return nil, backend.NewRpcBizErr(ErrApplyMetadataAuthority.Code, "not found metadata with metadataId of auth")
	}

	// check the metadataId whether has valid metadataAuth with current userType and user.
	has, err := svr.B.HasValidMetadataAuth(req.GetUserType(), req.GetUser(), req.GetAuth().GetOwner().GetIdentityId(), req.GetAuth().GetMetadataId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ApplyMetadataAuthority failed, query valid user metadataAuth failed, userType: {%s}, user: {%s}, metadataId: {%s}",
			req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId())

		errMsg := fmt.Sprintf(ErrApplyMetadataAuthority.Msg, "query valid user metadataAuth failed",
			req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId())
		return nil, backend.NewRpcBizErr(ErrApplyMetadataAuthority.Code, errMsg)
	}

	if has {
		log.Errorf("RPC-API:ApplyMetadataAuthority failed, has valid metadataAuth exists, userType: {%s}, user: {%s}, metadataId: {%s}",
			req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId())

		errMsg := fmt.Sprintf("%s, userType: {%s}, user: {%s}, metadataId: {%s}", ErrValidMetadataAuthMustCannotExist.Msg,
			req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId())
		return nil, backend.NewRpcBizErr(ErrValidMetadataAuthMustCannotExist.Code, errMsg)
	}

	metadataAuthorityMsg := types.NewMetadataAuthorityMessageFromRequest(req)
	metadataAuthId := metadataAuthorityMsg.GetMetadataAuthId()

	err = svr.B.SendMsg(metadataAuthorityMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:ApplyMetadataAuthority failed")

		errMsg := fmt.Sprintf(ErrApplyMetadataAuthority.Msg, "send metadata authority msg failed",
			req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId())
		return nil, backend.NewRpcBizErr(ErrApplyMetadataAuthority.Code, errMsg)
	}
	log.Debugf("RPC-API:ApplyMetadataAuthority succeed, userType: {%s}, user: {%s}, metadataOwner: {%s}, metadataId: {%s}, usageRule: {%s},  return metadataAuthId: {%s}",
		req.GetUserType().String(), req.GetUser(), req.GetAuth().GetOwner().String(), req.GetAuth().GetMetadataId(), req.GetAuth().GetUsageRule().String(), metadataAuthId)
	return &pb.ApplyMetadataAuthorityResponse{
		Status:         0,
		Msg:            backend.OK,
		MetadataAuthId: metadataAuthId,
	}, nil
}

func (svr *Server) RevokeMetadataAuthority(ctx context.Context, req *pb.RevokeMetadataAuthorityRequest) (*apicommonpb.SimpleResponse, error) {
	if req.GetUser() == "" {
		return nil, ErrReqGetUserForRevokeMetadataAuth
	}
	if !verifyUserType(req.GetUserType()) {
		return nil, ErrVerifyUserTypeForRevokeMetadataAuth
	}
	if req.GetMetadataAuthId() == "" {
		return nil, ErrReqAuthIDForRevokeMetadataAuth
	}
	if len(req.GetSign()) == 0 {
		return nil, ErrReqUserSignForRevokeMetadataAuth
	}

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, query local identity failed")
		return nil, fmt.Errorf("query local identity failed")
	}

	authorityList, err := svr.B.GetLocalMetadataAuthorityList()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, query local metadataAuth list failed, %s", err)
		return nil, backend.NewRpcBizErr(ErrRevokeMetadataAuthority.Code, "query local metadataAuth list failed")
	}

	for _, auth := range authorityList {

		if auth.GetUserType() == req.GetUserType() &&
			auth.GetUser() == req.GetUser() &&
			auth.GetData().GetMetadataAuthId() == req.GetMetadataAuthId() {


			// The data authorization application information that has been audited and cannot be revoked
			if auth.GetData().GetAuditOption() != apicommonpb.AuditMetadataOption_Audit_Pending {
				log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, the metadataAuth was audited")
				return nil, backend.NewRpcBizErr(ErrRevokeMetadataAuthority.Code, "the metadataAuth state was audited")
			}

			if auth.GetData().GetState() != apicommonpb.MetadataAuthorityState_MAState_Released {
				log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, the metadataAuth state was not released")
				return nil, backend.NewRpcBizErr(ErrRevokeMetadataAuthority.Code, "the metadataAuth state was not released")
			}

			switch auth.GetData().GetAuth().GetUsageRule().GetUsageType() {
			case apicommonpb.MetadataUsageType_Usage_Period:
				if timeutils.UnixMsecUint64() >= auth.GetData().GetAuth().GetUsageRule().GetEndAt() {
					log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, the metadataAuth had been expire")
					return nil, backend.NewRpcBizErr(ErrRevokeMetadataAuthority.Code, "the metadataAuth had been expire")
				}
			case apicommonpb.MetadataUsageType_Usage_Times:
				if auth.GetData().GetUsedQuo().GetUsedTimes() >= auth.GetData().GetAuth().GetUsageRule().GetTimes() {
					log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, the metadataAuth had been not enough times")
					return nil, backend.NewRpcBizErr(ErrRevokeMetadataAuthority.Code, "the metadataAuth had been not enough times")
				}
			default:
				log.Errorf("unknown usageType of the old metadataAuth on AuthorityManager.filterMetadataAuth(), metadataAuthId: {%s}", auth.GetData().GetMetadataAuthId())
				return nil, fmt.Errorf("unknown usageType of the old metadataAuth")
			}
		}
	}

	metadataAuthorityRevokeMsg := types.NewMetadataAuthorityRevokeMessageFromRequest(req)
	metadataAuthId := metadataAuthorityRevokeMsg.GetMetadataAuthId()

	if err := svr.B.SendMsg(metadataAuthorityRevokeMsg); nil != err {
		log.WithError(err).Error("RPC-API:RevokeMetadataAuthority failed")

		errMsg := fmt.Sprintf("%s, userType: {%s}, user: {%s}, MetadataAuthId: {%s}", ErrRevokeMetadataAuthority.Msg,
			req.GetUserType().String(), req.GetUser(), req.GetMetadataAuthId())
		return nil, backend.NewRpcBizErr(ErrRevokeMetadataAuthority.Code, errMsg)
	}
	log.Debugf("RPC-API:RevokeMetadataAuthority succeed, userType: {%s}, user: {%s}, metadataAuthId: {%s}",
		req.GetUserType().String(), req.GetUser(), metadataAuthId)
	return &apicommonpb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) AuditMetadataAuthority(ctx context.Context, req *pb.AuditMetadataAuthorityRequest) (*pb.AuditMetadataAuthorityResponse, error) {

	if "" == req.GetMetadataAuthId() {
		return nil, ErrReqAuthIDForAuditMetadataAuth
	}

	if req.GetAudit() == apicommonpb.AuditMetadataOption_Audit_Pending {
		return nil, ErrValidAuditMetadataOptionMustCannotPending
	}

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:AuditMetadataAuthority failed, query local identity failed")
		return nil, fmt.Errorf("query local identity failed")
	}

	option, err := svr.B.AuditMetadataAuthority(types.NewMetadataAuthAudit(req.GetMetadataAuthId(), req.GetSuggestion(), req.GetAudit()))
	if nil != err {
		log.WithError(err).Error("RPC-API:AuditMetadataAuthority failed")

		errMsg := fmt.Sprintf("%s, metadataAuthId: {%s}, audit option: {%s}, audit suggestion: {%s}", ErrAuditMetadataAuth.Msg,
			req.GetMetadataAuthId(), req.GetAudit().String(), req.GetSuggestion())
		return nil, backend.NewRpcBizErr(ErrAuditMetadataAuth.Code, errMsg)
	}
	log.Debugf("RPC-API:AuditMetadataAuthority succeed, metadataAuthId: {%s}, audit option: {%s}, audit suggestion: {%s}",
		req.GetMetadataAuthId(), req.GetAudit().String(), req.GetSuggestion())

	return &pb.AuditMetadataAuthorityResponse{
		Status: 0,
		Msg:    backend.OK,
		Audit:  option,
	}, nil
}

func (svr *Server) GetLocalMetadataAuthorityList(ctx context.Context, req *pb.GetMetadataAuthorityListRequest) (*pb.GetMetadataAuthorityListResponse, error) {
	authorityList, err := svr.B.GetLocalMetadataAuthorityList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetLocalMetadataAuthorityList failed")
		return nil, ErrGetAuthorityList
	}
	arr := make([]*pb.GetMetadataAuthority, len(authorityList))
	for i, auth := range authorityList {
		arr[i] = &pb.GetMetadataAuthority{
			MetadataAuthId:  auth.GetData().GetMetadataAuthId(),
			User:            auth.GetData().GetUser(),
			UserType:        auth.GetData().GetUserType(),
			Auth:            auth.GetData().GetAuth(),
			AuditOption: 	 auth.GetData().GetAuditOption(),
			AuditSuggestion: auth.GetData().GetAuditSuggestion(),
			UsedQuo:         auth.GetData().GetUsedQuo(),
			ApplyAt:         auth.GetData().GetApplyAt(),
			AuditAt:         auth.GetData().GetAuditAt(),
			State: 			 auth.GetData().GetState(),
		}
	}
	log.Debugf("RPC-API:GetLocalMetadataAuthorityList succeed, metadata authority list, len: {%d}", len(authorityList))
	return &pb.GetMetadataAuthorityListResponse{
		Status: 0,
		Msg:    backend.OK,
		List:   arr,
	}, nil
}

func (svr *Server) GetGlobalMetadataAuthorityList(ctx context.Context, req *pb.GetMetadataAuthorityListRequest) (*pb.GetMetadataAuthorityListResponse, error) {

	authorityList, err := svr.B.GetGlobalMetadataAuthorityList(req.GetLastUpdated(), backend.DefaultPageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetGlobalMetadataAuthorityList failed")
		return nil, ErrGetAuthorityList
	}
	arr := make([]*pb.GetMetadataAuthority, len(authorityList))
	for i, auth := range authorityList {
		arr[i] = &pb.GetMetadataAuthority{
			MetadataAuthId:  auth.GetData().GetMetadataAuthId(),
			User:            auth.GetData().GetUser(),
			UserType:        auth.GetData().GetUserType(),
			Auth:            auth.GetData().GetAuth(),
			AuditOption: 	 auth.GetData().GetAuditOption(),
			AuditSuggestion: auth.GetData().GetAuditSuggestion(),
			UsedQuo:         auth.GetData().GetUsedQuo(),
			ApplyAt:         auth.GetData().GetApplyAt(),
			AuditAt:         auth.GetData().GetAuditAt(),
			State: 			 auth.GetData().GetState(),
		}
	}
	log.Debugf("RPC-API:GetGlobalMetadataAuthorityList succeed, metadata authority list, len: {%d}", len(authorityList))
	return &pb.GetMetadataAuthorityListResponse{
		Status: 0,
		Msg:    backend.OK,
		List:   arr,
	}, nil
}

func verifyUserType(userType apicommonpb.UserType) bool {
	switch userType {
	case apicommonpb.UserType_User_1:
		return true
	case apicommonpb.UserType_User_2:
		return true
	case apicommonpb.UserType_User_3:
		return true
	default:
		return false
	}
}
