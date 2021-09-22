package auth

import (
	"context"
	"errors"
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
			req.Member.IdentityId, req.Member.NodeId, req.Member.NodeName)
		return nil, ErrSendIdentityMsg
	}

	if nil != identity {
		log.Errorf("RPC-API:ApplyIdentityJoin failed, identity was already exist, old identityId: {%s}, old nodeId: {%s}, old nodeName: {%s}",
			identity.IdentityId(), identity.NodeId(), identity.Name())
		return nil, ErrSendIdentityMsg
	}

	if req.Member == nil {
		return nil, errors.New("Invalid Params, req.Member is nil")
	}

	if "" == strings.Trim(req.Member.IdentityId, "") ||
		"" == strings.Trim(req.Member.NodeName, "") {
		return nil, errors.New("Invalid Params, req.Member.IdentityId or req.Member.Name is empty")
	}

	identityMsg := types.NewIdentityMessageFromRequest(req)
	err = svr.B.SendMsg(identityMsg)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ApplyIdentityJoin failed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}",
			req.Member.IdentityId, req.Member.NodeId, req.Member.NodeName)
		return nil, ErrSendIdentityMsg
	}
	log.Debugf("RPC-API:ApplyIdentityJoin succeed SendMsg, identityId: {%s}, nodeId: {%s}, nodeName: {%s}",
		req.Member.IdentityId, req.Member.NodeId, req.Member.NodeName)
	return &apicommonpb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) RevokeIdentityJoin(ctx context.Context, req *emptypb.Empty) (*apicommonpb.SimpleResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if rawdb.IsDBNotFoundErr(err) {
		log.WithError(err).Errorf("RPC-API:RevokeIdentityJoin failed, the identity was not exist, can not revoke identity")
		return nil, ErrSendIdentityRevokeMsg
	}

	identityRevokeMsg := types.NewIdentityRevokeMessage()
	err = svr.B.SendMsg(identityRevokeMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:RevokeIdentityJoin failed")
		return nil, ErrSendIdentityRevokeMsg
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
			NodeName:   identity.Name(),
			NodeId:     identity.NodeId(),
			IdentityId: identity.IdentityId(),
		},
	}, nil
}

func (svr *Server) GetIdentityList(ctx context.Context, req *emptypb.Empty) (*pb.GetIdentityListResponse, error) {
	identityList, err := svr.B.GetIdentityList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetIdentityList failed")
		return nil, ErrGetIdentityList
	}
	arr := make([]*apicommonpb.Organization, len(identityList))
	for i, identity := range identityList {
		iden := &apicommonpb.Organization{
			NodeName:   identity.Name(),
			NodeId:     identity.NodeId(),
			IdentityId: identity.IdentityId(),
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
		return nil, errors.New("required GetUser")
	}
	if !verifyUserType(req.GetUserType()) {
		return nil, errors.New("required right user type")
	}
	if req.GetAuth() == nil {
		return nil, errors.New("required metadata authority")
	}
	if len(req.GetSign()) == 0 {
		return nil, errors.New("required user sign")
	}

	has, err := svr.B.HasValidUserMetadataAuth(req.GetUserType(), req.GetUser(), req.GetAuth().GetMetadataId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ApplyMetadataAuthority failed, query valid user metadataAuth failed, userType: {%s}, user: {%s}, metadataId: {%s}",
			req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId())
		return nil, ErrSendMetadataAuthMsg
	}

	if has {
		log.Errorf("RPC-API:ApplyMetadataAuthority failed, has valid metadata, userType: {%s}, user: {%s}, metadataId: {%s}",
			req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId())
		return nil, ErrValidMetadataAuthMustCannotExist
	}

	metadataAuthorityMsg := types.NewMetadataAuthorityMessageFromRequest(req)
	metadataAuthId := metadataAuthorityMsg.GetMetadataAuthId()

	err = svr.B.SendMsg(metadataAuthorityMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:ApplyMetadataAuthority failed")
		return nil, ErrSendMetadataAuthMsg
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
		return nil, errors.New("required GetUser")
	}
	if !verifyUserType(req.GetUserType()) {
		return nil, errors.New("required right user type")
	}
	if req.GetMetadataAuthId() == "" {
		return nil, errors.New("required metadataAuthId")
	}
	if len(req.GetSign()) == 0 {
		return nil, errors.New("required user sign")
	}

	metadataAuthorityRevokeMsg := types.NewMetadataAuthorityRevokeMessageFromRequest(req)
	metadataAuthId := metadataAuthorityRevokeMsg.GetMetadataAuthId()

	err := svr.B.SendMsg(metadataAuthorityRevokeMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:RevokeMetadataAuthority failed")
		return nil, ErrSendMetadataAuthMsg
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
		return nil, errors.New("require metadataAuthId")
	}

	if req.GetAudit() == apicommonpb.AuditMetadataOption_Audit_Pending {
		return nil, errors.New("audit option can not is pending")
	}

	option, err := svr.B.AuditMetadataAuthority(types.NewMetadataAuthAudit(req.GetMetadataAuthId(), req.GetSuggestion(), req.GetAudit()))
	if nil != err {
		log.WithError(err).Error("RPC-API:AuditMetadataAuthority failed")
		return nil, ErrAuditMetadataAuth
	}
	log.Debugf("RPC-API:AuditMetadataAuthority succeed, metadataAuthId: {%s}, audit option: {%s}, audit suggestion: {%s}",
		req.GetMetadataAuthId(), req.GetAudit().String(), req.GetSuggestion())

	return &pb.AuditMetadataAuthorityResponse{
		Status: 0,
		Msg:    backend.OK,
		Audit: option,
	}, nil
}

func (svr *Server) GetMetadataAuthorityList(context.Context, *emptypb.Empty) (*pb.GetMetadataAuthorityListResponse, error) {
	authorityList, err := svr.B.GetMetadataAuthorityList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetMetadataAuthorityList failed")
		return nil, ErrGetAuthorityList
	}
	arr := make([]*pb.GetMetadataAuthority, len(authorityList))
	for i, auth := range authorityList {
		arr[i] = &pb.GetMetadataAuthority{
			MetadataAuthId:  auth.GetData().GetMetadataAuthId(),
			User:            auth.GetData().GetUser(),
			UserType:        auth.GetData().GetUserType(),
			Auth:            auth.GetData().GetAuth(),
			AuditSuggestion: auth.GetData().GetAuditSuggestion(),
			ApplyAt:         auth.GetData().GetApplyAt(),
			AuditAt:         auth.GetData().GetAuditAt(),
		}
	}
	log.Debugf("Query all authority list, len: {%d}", len(authorityList))
	return &pb.GetMetadataAuthorityListResponse{
		Status: 0,
		Msg:    backend.OK,
		List:   arr,
	}, nil
}

func (svr *Server) GetMetadataAuthorityListByUser(ctx context.Context, req *pb.GetMetadataAuthorityListByUserRequest) (*pb.GetMetadataAuthorityListResponse, error) {

	if "" == req.GetUser() {
		return nil, errors.New("require user")
	}

	if !verifyUserType(req.GetUserType()) {
		return nil, errors.New("userType is invalid")
	}

	authorityList, err := svr.B.GetMetadataAuthorityListByUser(req.GetUserType(), req.GetUser())
	if nil != err {
		log.WithError(err).Error("RPC-API:GetMetadataAuthorityListByUser failed")
		return nil, ErrGetAuthorityList
	}
	arr := make([]*pb.GetMetadataAuthority, len(authorityList))
	for i, auth := range authorityList {
		arr[i] = &pb.GetMetadataAuthority{
			MetadataAuthId:  auth.GetData().GetMetadataAuthId(),
			User:            auth.GetData().GetUser(),
			UserType:        auth.GetData().GetUserType(),
			Auth:            auth.GetData().GetAuth(),
			AuditSuggestion: auth.GetData().GetAuditSuggestion(),
			ApplyAt:         auth.GetData().GetApplyAt(),
			AuditAt:         auth.GetData().GetAuditAt(),
		}
	}
	log.Debugf("Query all authority list by user, userType: {%s}, user: {%s}, len: {%d}", req.GetUserType().String(), req.GetUser(), len(authorityList))
	return &pb.GetMetadataAuthorityListResponse{
		Status: 0,
		Msg:    backend.OK,
		List:   arr,
	}, nil
}

func verifyUserType(userType apicommonpb.UserType) bool {
	switch userType {
	case apicommonpb.UserType_User_ETH:
		return true
	case apicommonpb.UserType_User_ATP:
		return true
	case apicommonpb.UserType_User_LAT:
		return true
	default:
		return false
	}
}
