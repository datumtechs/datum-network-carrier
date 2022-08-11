package auth

import (
	"context"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/carrierdb/rawdb"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/p2p"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/rpc/backend"
	"github.com/datumtechs/datum-network-carrier/signsuite"
	"github.com/datumtechs/datum-network-carrier/types"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

// for organization identity

func (svr *Server) ApplyIdentityJoin(ctx context.Context, req *carrierapipb.ApplyIdentityJoinRequest) (*carriertypespb.SimpleResponse, error) {

	identity, err := svr.B.GetNodeIdentity()
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("RPC-API:ApplyIdentityJoin failed, query local identity failed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}, imgUrl: {%s}, details: {%s}",
			req.GetInformation().GetIdentityId(), req.GetInformation().GetNodeId(), req.GetInformation().GetNodeName(), req.GetInformation().GetImageUrl(), req.GetInformation().GetDetails())

		errMsg := fmt.Sprintf("query local identity failed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}, imgUrl: {%s}, details: {%s}",
			req.GetInformation().GetIdentityId(), req.GetInformation().GetNodeId(), req.GetInformation().GetNodeName(), req.GetInformation().GetImageUrl(), req.GetInformation().GetDetails())
		return &carriertypespb.SimpleResponse{Status: backend.ErrApplyIdentityMsg.ErrCode(), Msg: errMsg}, nil
	}

	if nil != identity {
		log.Errorf("RPC-API:ApplyIdentityJoin failed, identity was already exist, old identityId: {%s}, old nodeId: {%s}, old nodeName: {%s}, old imgUrl: {%s}, old details: {%s}",
			identity.GetIdentityId(), identity.GetNodeId(), identity.GetName(), identity.GetImageUrl(), identity.GetDetails())

		errMsg := fmt.Sprintf("identity was already exist, identityId: {%s}, nodeId: {%s}, nodeName: {%s}, imgUrl: {%s}, details: {%s}",
			identity.GetIdentityId(), identity.GetNodeId(), identity.GetName(), identity.GetImageUrl(), identity.GetDetails())
		return &carriertypespb.SimpleResponse{Status: backend.ErrApplyIdentityMsg.ErrCode(), Msg: errMsg}, nil
	}

	if nil == req.GetInformation() {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: backend.ErrRequireParams.Error()}, nil
	}

	if "" == strings.Trim(req.GetInformation().GetIdentityId(), "") ||
		"" == strings.Trim(req.GetInformation().GetNodeName(), "") {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: backend.ErrRequireParams.Error()}, nil
	}

	identityMsg := types.NewIdentityMessageFromRequest(req)
	if err := identityMsg.CheckLength(); nil != err {
		errMsg := fmt.Sprintf("check fields len failed, %s , identityId: {%s}, nodeId: {%s}, nodeName: {%s}, imgUrl: {%s}, details: {%s}",
			err, identity.GetIdentityId(), identity.GetNodeId(), identity.GetName(), identity.GetImageUrl(), identity.GetDetails())
		return &carriertypespb.SimpleResponse{Status: backend.ErrApplyIdentityMsg.ErrCode(), Msg: errMsg}, nil
	}

	if err = svr.B.SendMsg(identityMsg); nil != err {
		log.WithError(err).Errorf("RPC-API:ApplyIdentityJoin failed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}, imgUrl: {%s}, details: {%s}",
			req.GetInformation().GetIdentityId(), req.GetInformation().GetNodeId(), req.GetInformation().GetNodeName(), req.GetInformation().GetImageUrl(), req.GetInformation().GetDetails())

		errMsg := fmt.Sprintf("%s, identityId: {%s}, nodeId: {%s}, nodeName: {%s}", backend.ErrApplyIdentityMsg.Error(),
			req.GetInformation().GetIdentityId(), req.GetInformation().GetNodeId(), req.GetInformation().GetNodeName())
		return &carriertypespb.SimpleResponse{Status: backend.ErrApplyIdentityMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:ApplyIdentityJoin succeed SendMsg, identityId: {%s}, nodeId: {%s}, nodeName: {%s}, imgUrl: {%s}, details: {%s}",
		req.GetInformation().GetIdentityId(), req.GetInformation().GetNodeId(), req.GetInformation().GetNodeName(), req.GetInformation().GetImageUrl(), req.GetInformation().GetDetails())
	return &carriertypespb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) RevokeIdentityJoin(ctx context.Context, req *emptypb.Empty) (*carriertypespb.SimpleResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if rawdb.IsDBNotFoundErr(err) {
		log.WithError(err).Errorf("RPC-API:RevokeIdentityJoin failed, the identity was not exist, can not revoke identity")
		errMsg := fmt.Sprintf("%s, the identity was not exist, can not revoke identity", backend.ErrRevokeIdentityMsg.Error())
		return &carriertypespb.SimpleResponse{Status: backend.ErrRevokeIdentityMsg.ErrCode(), Msg: errMsg}, nil
	}

	// what if local task we can not revoke identity
	has, err := svr.B.HasLocalTask()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokeIdentityJoin failed, can not check has local task")

		errMsg := fmt.Sprintf("%s, can not check has local task", backend.ErrRevokeIdentityMsg.Error())
		return &carriertypespb.SimpleResponse{Status: backend.ErrRevokeIdentityMsg.ErrCode(), Msg: errMsg}, nil
	}

	if has {
		log.WithError(err).Errorf("RPC-API:RevokeIdentityJoin failed, don't revoke identity when has local task")

		errMsg := fmt.Sprintf("%s, don't revoke identity when has local task", backend.ErrRevokeIdentityMsg.Error())
		return &carriertypespb.SimpleResponse{Status: backend.ErrRevokeIdentityMsg.ErrCode(), Msg: errMsg}, nil
	}

	identityRevokeMsg := types.NewIdentityRevokeMessage()

	if err = svr.B.SendMsg(identityRevokeMsg); nil != err {
		log.WithError(err).Error("RPC-API:RevokeIdentityJoin failed")

		errMsg := fmt.Sprintf("%s, send identity revoke msg failed", backend.ErrRevokeIdentityMsg.Error())
		return &carriertypespb.SimpleResponse{Status: backend.ErrRevokeIdentityMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debug("RPC-API:RevokeIdentityJoin succeed SendMsg")
	return &carriertypespb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) GetNodeIdentity(ctx context.Context, req *emptypb.Empty) (*carrierapipb.GetNodeIdentityResponse, error) {
	identity, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetNodeIdentity failed")
		return &carrierapipb.GetNodeIdentityResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}
	return &carrierapipb.GetNodeIdentityResponse{
		Status: 0,
		Msg:    backend.OK,
		Information: &carriertypespb.Organization{
			NodeName:   identity.GetName(),
			NodeId:     identity.GetNodeId(),
			IdentityId: identity.GetIdentityId(),
			ImageUrl:   identity.GetImageUrl(),
			Details:    identity.GetDetails(),
			DataStatus: identity.GetDataStatus(),
			Status:     identity.GetStatus(),
			UpdateAt:   identity.GetUpdateAt(),
		},
	}, nil
}

func (svr *Server) GetIdentityList(ctx context.Context, req *carrierapipb.GetIdentityListRequest) (*carrierapipb.GetIdentityListResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = backend.DefaultPageSize
	}
	identityList, err := svr.B.GetIdentityList(req.GetLastUpdated(), pageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:QueryIdentityList failed")
		return &carrierapipb.GetIdentityListResponse{Status: backend.ErrQueryIdentityList.ErrCode(), Msg: backend.ErrQueryIdentityList.Error()}, nil
	}
	arr := make([]*carriertypespb.Organization, len(identityList))
	for i, identity := range identityList {
		iden := &carriertypespb.Organization{
			NodeName:   identity.GetName(),
			NodeId:     identity.GetNodeId(),
			IdentityId: identity.GetIdentityId(),
			ImageUrl:   identity.GetImageUrl(),
			Details:    identity.GetDetails(),
			UpdateAt:   identity.GetUpdateAt(),
			DataStatus: identity.GetDataStatus(),
			Status:     identity.GetStatus(),
		}
		arr[i] = iden
		if hexNodeId, err := p2p.HexPeerID(identity.GetNodeId()); err == nil {
			log.Debugf("the nodes that have entered the network have node id:{%s}", hexNodeId)
		}
	}
	log.Debugf("Query all org's identity list, len: {%d}", len(identityList))
	return &carrierapipb.GetIdentityListResponse{
		Status:    0,
		Msg:       backend.OK,
		Identitys: arr,
	}, nil
}

// for metadata authority apply

func (svr *Server) ApplyMetadataAuthority(ctx context.Context, req *carrierapipb.ApplyMetadataAuthorityRequest) (*carrierapipb.ApplyMetadataAuthorityResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ApplyMetadataAuthority failed, query local identity failed")
		return &carrierapipb.ApplyMetadataAuthorityResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if req.GetUser() == "" {
		return &carrierapipb.ApplyMetadataAuthorityResponse{Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "the user is empty"}, nil
	}
	if !verifyUserType(req.GetUserType()) {
		return &carrierapipb.ApplyMetadataAuthorityResponse{Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "the user type is wrong"}, nil
	}
	if req.GetAuth() == nil {
		return &carrierapipb.ApplyMetadataAuthorityResponse{Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "the metadata auth is nil"}, nil
	}
	if len(req.GetSign()) == 0 {
		return &carrierapipb.ApplyMetadataAuthorityResponse{Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "the user signature is nil"}, nil
	}

	metadataAuthorityMsg := types.NewMetadataAuthorityMessageFromRequest(req)
	if common.OpenMessageSignCheck {
		from, _, err := signsuite.Sender(req.GetUserType(), metadataAuthorityMsg.Hash(), req.GetSign())
		if nil != err {
			log.WithError(err).Errorf("RPC-API:ApplyMetadataAuthority failed, cannot fetch sender from sign, userType: {%s}, user: {%s}",
				req.GetUserType().String(), req.GetUser())
			return &carrierapipb.ApplyMetadataAuthorityResponse{Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "cannot fetch sender from sign"}, nil
		}
		if from != req.GetUser() {
			log.WithError(err).Errorf("RPC-API:ApplyMetadataAuthority failed, sender from sign and user is not sameone, userType: {%s}, user: {%s}, sender of sign: {%s}",
				req.GetUserType().String(), req.GetUser(), from)
			return &carrierapipb.ApplyMetadataAuthorityResponse{Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "the user sign is invalid"}, nil
		}
	}

	// verify special options of metadataAuth...

	// ############################################
	// ############################################
	// NOTE:
	// 		check metadataAuth with metadata
	// ############################################
	// ############################################
	pass, err := svr.B.VerifyMetadataAuthWithMetadataOption(types.NewMetadataAuthority(&carriertypespb.MetadataAuthorityPB{
		MetadataAuthId:  "",
		User:            req.GetUser(),
		UserType:        req.GetUserType(),
		Auth:            req.GetAuth(),
		AuditOption:     commonconstantpb.AuditMetadataOption_Audit_Pending,
		AuditSuggestion: "",
		UsedQuo: &carriertypespb.MetadataUsedQuo{
			UsageType: req.GetAuth().GetUsageRule().GetUsageType(),
			Expire:    false, // Initialized zero value
			UsedTimes: 0,     // Initialized zero value
		},
		ApplyAt: metadataAuthorityMsg.GetCreateAt(),
		AuditAt: 0,
		State:   commonconstantpb.MetadataAuthorityState_MAState_Released,
		Sign:    req.GetSign(),
	}))
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ApplyMetadataAuthority failed, can not verify metadataAuth with metadataOption, metadataId: {%s}, auth: %s",
			req.GetAuth().GetMetadataId(), req.GetAuth().String())
		return &carrierapipb.ApplyMetadataAuthorityResponse{Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "can not verify metadataAuth with metadataOption"}, nil
	}
	if !pass {
		log.WithError(err).Errorf("RPC-API:ApplyMetadataAuthority failed, invalid metadataAuth, metadataId: {%s}, auth: %s",
			req.GetAuth().GetMetadataId(), req.GetAuth().String())
		return &carrierapipb.ApplyMetadataAuthorityResponse{Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "invalid metadataAuth"}, nil
	}

	metadataAuthId := metadataAuthorityMsg.GenMetadataAuthId()

	if err = svr.B.SendMsg(metadataAuthorityMsg); nil != err {
		log.WithError(err).Error("RPC-API:ApplyMetadataAuthority failed")

		errMsg := fmt.Sprintf(backend.ErrApplyMetadataAuthority.Error(), "send metadata authority msg failed",
			req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId())
		return &carrierapipb.ApplyMetadataAuthorityResponse{Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:ApplyMetadataAuthority succeed, userType: {%s}, user: {%s}, metadataOwner: {%s}, metadataId: {%s}, usageRule: {%s},  return metadataAuthId: {%s}",
		req.GetUserType().String(), req.GetUser(), req.GetAuth().GetOwner().String(), req.GetAuth().GetMetadataId(), req.GetAuth().GetUsageRule().String(), metadataAuthId)
	return &carrierapipb.ApplyMetadataAuthorityResponse{
		Status:         0,
		Msg:            backend.OK,
		MetadataAuthId: metadataAuthId,
	}, nil
}

func (svr *Server) RevokeMetadataAuthority(ctx context.Context, req *carrierapipb.RevokeMetadataAuthorityRequest) (*carriertypespb.SimpleResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, query local identity failed")
		return &carriertypespb.SimpleResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if req.GetUser() == "" {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "the user is empty"}, nil
	}
	if !verifyUserType(req.GetUserType()) {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "the user type is wrong"}, nil
	}
	if req.GetMetadataAuthId() == "" {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "the metadata auth is nil"}, nil
	}
	if len(req.GetSign()) == 0 {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "the user signature is nil"}, nil
	}

	metadataAuthorityRevokeMsg := types.NewMetadataAuthorityRevokeMessageFromRequest(req)
	if common.OpenMessageSignCheck {
		from, _, err := signsuite.Sender(req.GetUserType(), metadataAuthorityRevokeMsg.Hash(), req.GetSign())
		if nil != err {
			log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, cannot fetch sender from sign, userType: {%s}, user: {%s}",
				req.GetUserType().String(), req.GetUser())
			return &carriertypespb.SimpleResponse{Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "cannot fetch sender from sign"}, nil
		}
		if from != req.GetUser() {
			log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, sender from sign and user is not sameone, userType: {%s}, user: {%s}, sender of sign: {%s}",
				req.GetUserType().String(), req.GetUser(), from)
			return &carriertypespb.SimpleResponse{Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "the user sign is invalid"}, nil
		}
	}

	metadataAuthId := metadataAuthorityRevokeMsg.GetMetadataAuthId()

	auth, err := svr.B.GetMetadataAuthority(metadataAuthId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, query local metadataAuth failed, metadataAuthId: {%s}", metadataAuthId)
		return &carriertypespb.SimpleResponse{Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "query local metadataAuth failed"}, nil
	}

	if auth.GetUserType() != req.GetUserType() || auth.GetUser() != req.GetUser() {
		log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, user or userType of local metadataAuth and req's is not same, metadataAuthId: {%s}", metadataAuthId)
		return &carriertypespb.SimpleResponse{Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "user or userType of local metadataAuth and req's is not same"}, nil
	}

	// The data authorization application information that has been audited and cannot be revoked
	pass, err := svr.B.VerifyMetadataAuthInfo(auth)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, can not verify metadataAuth, metadataAuthId: {%s}", metadataAuthId)
		return &carriertypespb.SimpleResponse{Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "can not verify metadataAuth"}, nil
	}
	if !pass {
		log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, invalid metadataAuth, metadataAuthId: {%s}", metadataAuthId)
		return &carriertypespb.SimpleResponse{Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "invalid metadataAuth"}, nil
	}

	if err := svr.B.SendMsg(metadataAuthorityRevokeMsg); nil != err {
		log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, metadataAuthId: {%s}", metadataAuthId)

		errMsg := fmt.Sprintf("%s, userType: {%s}, user: {%s}, metadataAuthId: {%s}", backend.ErrRevokeMetadataAuthority.Error(),
			req.GetUserType().String(), req.GetUser(), req.GetMetadataAuthId())
		return &carriertypespb.SimpleResponse{Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:RevokeMetadataAuthority succeed, userType: {%s}, user: {%s}, metadataAuthId: {%s}",
		req.GetUserType().String(), req.GetUser(), metadataAuthId)
	return &carriertypespb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) AuditMetadataAuthority(ctx context.Context, req *carrierapipb.AuditMetadataAuthorityRequest) (*carrierapipb.AuditMetadataAuthorityResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:AuditMetadataAuthority failed, query local identity failed")
		return &carrierapipb.AuditMetadataAuthorityResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if "" == req.GetMetadataAuthId() {
		return &carrierapipb.AuditMetadataAuthorityResponse{Status: backend.ErrAuditMetadataAuth.ErrCode(), Msg: "the metadataAuth Id is empty"}, nil
	}

	if req.GetAudit() == commonconstantpb.AuditMetadataOption_Audit_Pending {
		return &carrierapipb.AuditMetadataAuthorityResponse{Status: backend.ErrAuditMetadataAuth.ErrCode(), Msg: "the valid audit metadata option must cannot pending"}, nil
	}

	option, err := svr.B.AuditMetadataAuthority(types.NewMetadataAuthAudit(req.GetMetadataAuthId(), req.GetSuggestion(), req.GetAudit()))
	if nil != err {
		log.WithError(err).Error("RPC-API:AuditMetadataAuthority failed")

		errMsg := fmt.Sprintf("%s, metadataAuthId: {%s}, audit option: {%s}, audit suggestion: {%s}", backend.ErrAuditMetadataAuth.Error(),
			req.GetMetadataAuthId(), req.GetAudit().String(), req.GetSuggestion())
		return &carrierapipb.AuditMetadataAuthorityResponse{Status: backend.ErrAuditMetadataAuth.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:AuditMetadataAuthority succeed, metadataAuthId: {%s}, audit option: {%s}, audit suggestion: {%s}",
		req.GetMetadataAuthId(), req.GetAudit().String(), req.GetSuggestion())

	return &carrierapipb.AuditMetadataAuthorityResponse{
		Status: 0,
		Msg:    backend.OK,
		Audit:  option,
	}, nil
}

func (svr *Server) GetLocalMetadataAuthorityList(ctx context.Context, req *carrierapipb.GetMetadataAuthorityListRequest) (*carrierapipb.GetMetadataAuthorityListResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = backend.DefaultPageSize
	}
	authorityList, err := svr.B.GetLocalMetadataAuthorityList(req.GetLastUpdated(), pageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetLocalMetadataAuthorityList failed")
		return &carrierapipb.GetMetadataAuthorityListResponse{Status: backend.ErrQueryAuthorityList.ErrCode(), Msg: backend.ErrQueryAuthorityList.Error()}, nil
	}
	arr := make([]*carriertypespb.MetadataAuthorityDetail, len(authorityList))
	for i, auth := range authorityList {
		arr[i] = &carriertypespb.MetadataAuthorityDetail{
			MetadataAuthId:  auth.GetData().GetMetadataAuthId(),
			User:            auth.GetData().GetUser(),
			UserType:        auth.GetData().GetUserType(),
			Auth:            auth.GetData().GetAuth(),
			AuditOption:     auth.GetData().GetAuditOption(),
			AuditSuggestion: auth.GetData().GetAuditSuggestion(),
			UsedQuo:         auth.GetData().GetUsedQuo(),
			ApplyAt:         auth.GetData().GetApplyAt(),
			AuditAt:         auth.GetData().GetAuditAt(),
			State:           auth.GetData().GetState(),
			UpdateAt:        auth.GetData().GetUpdateAt(),
			Nonce:           auth.GetData().GetNonce(),
		}
	}
	log.Debugf("RPC-API:GetLocalMetadataAuthorityList succeed, metadata authority list, len: {%d}", len(authorityList))
	return &carrierapipb.GetMetadataAuthorityListResponse{
		Status:        0,
		Msg:           backend.OK,
		MetadataAuths: arr,
	}, nil
}

func (svr *Server) GetGlobalMetadataAuthorityList(ctx context.Context, req *carrierapipb.GetMetadataAuthorityListRequest) (*carrierapipb.GetMetadataAuthorityListResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = backend.DefaultPageSize
	}
	authorityList, err := svr.B.GetGlobalMetadataAuthorityList(req.GetLastUpdated(), pageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetGlobalMetadataAuthorityList failed")
		return &carrierapipb.GetMetadataAuthorityListResponse{Status: backend.ErrQueryAuthorityList.ErrCode(), Msg: backend.ErrQueryAuthorityList.Error()}, nil
	}
	arr := make([]*carriertypespb.MetadataAuthorityDetail, len(authorityList))
	for i, auth := range authorityList {
		arr[i] = &carriertypespb.MetadataAuthorityDetail{
			MetadataAuthId:  auth.GetData().GetMetadataAuthId(),
			User:            auth.GetData().GetUser(),
			UserType:        auth.GetData().GetUserType(),
			Auth:            auth.GetData().GetAuth(),
			AuditOption:     auth.GetData().GetAuditOption(),
			AuditSuggestion: auth.GetData().GetAuditSuggestion(),
			UsedQuo:         auth.GetData().GetUsedQuo(),
			ApplyAt:         auth.GetData().GetApplyAt(),
			AuditAt:         auth.GetData().GetAuditAt(),
			State:           auth.GetData().GetState(),
			UpdateAt:        auth.GetData().GetUpdateAt(),
			Nonce:           auth.GetData().GetNonce(),
		}
	}
	log.Debugf("RPC-API:GetGlobalMetadataAuthorityList succeed, metadata authority list, len: {%d}", len(authorityList))
	return &carrierapipb.GetMetadataAuthorityListResponse{
		Status:        0,
		Msg:           backend.OK,
		MetadataAuths: arr,
	}, nil
}

func (svr *Server) UpdateIdentityCredential(ctx context.Context, req *carrierapipb.UpdateIdentityCredentialRequest) (*carriertypespb.SimpleResponse, error) {
	identityId := req.GetIdentityId()
	credential := req.GetCredential()
	if identityId == "" {
		return &carriertypespb.SimpleResponse{Status: backend.ErrUpdateIdentityCredential.Code, Msg: "identityId can not be empty "}, nil
	}
	if len(identityId) != common.IdentityIdLength {
		return &carriertypespb.SimpleResponse{Status: backend.ErrUpdateIdentityCredential.Code, Msg: fmt.Sprintf("identityId len not equal %d ", common.IdentityIdLength)}, nil
	}
	updateIdentityCredentialMsg := &types.UpdateIdentityCredentialMsg{
		IdentityId: identityId,
		Credential: credential,
	}
	if err := svr.B.SendMsg(updateIdentityCredentialMsg); nil != err {
		log.WithError(err).Error("RPC-API:UpdateIdentityCredential failed")
		errMsg := fmt.Sprintf("UpdateIdentityCredential SendMsg fail, identity is %s,Credential is %s", identityId, credential)
		return &carriertypespb.SimpleResponse{Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: errMsg}, nil
	}
	return &carriertypespb.SimpleResponse{Status: 0, Msg: backend.OK}, nil
}

func verifyUserType(userType commonconstantpb.UserType) bool {
	switch userType {
	case commonconstantpb.UserType_User_1: // PlatON
		return true
	case commonconstantpb.UserType_User_2: // Alaya
		return true
	case commonconstantpb.UserType_User_3: // Ethereum
		return true
	default:
		return false
	}
}
