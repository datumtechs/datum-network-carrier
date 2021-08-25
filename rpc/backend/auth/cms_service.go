package auth

import (
	"context"
	"errors"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (svr *AuthServiceServer) ApplyIdentityJoin(ctx context.Context, req *pb.ApplyIdentityJoinRequest) (*apipb.SimpleResponse, error) {

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

	identityMsg := new(types.IdentityMsg)
	if req.Member == nil {
		return nil, errors.New("Invalid Params, req.Member is nil")
	}

	if "" == strings.Trim(req.Member.IdentityId, "") ||
		"" == strings.Trim(req.Member.NodeName, "") {
		return nil, errors.New("Invalid Params, req.Member.IdentityId or req.Member.Name is empty")
	}

	identityMsg.NodeAlias = &types.NodeAlias{}
	identityMsg.Name = req.Member.NodeName
	identityMsg.IdentityId = req.Member.IdentityId
	//identityMsg.NodeId = req.Member.NodeId
	identityMsg.CreateAt = uint64(timeutils.UnixMsec())

	err = svr.B.SendMsg(identityMsg)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ApplyIdentityJoin failed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}",
			req.Member.IdentityId, req.Member.NodeId, req.Member.NodeName)
		return nil, ErrSendIdentityMsg
	}
	log.Debugf("RPC-API:ApplyIdentityJoin succeed SendMsg, identityId: {%s}, nodeId: {%s}, nodeName: {%s}",
		req.Member.IdentityId, req.Member.NodeId, req.Member.NodeName)
	return &apipb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *AuthServiceServer) RevokeIdentityJoin(ctx context.Context, req *emptypb.Empty) (*apipb.SimpleResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if rawdb.IsDBNotFoundErr(err) {
		log.WithError(err).Errorf("RPC-API:RevokeIdentityJoin failed, the identity was not exist, can not revoke identity")
		return nil, ErrSendIdentityRevokeMsg
	}

	identityRevokeMsg := new(types.IdentityRevokeMsg)
	identityRevokeMsg.CreateAt = uint64(timeutils.UnixMsec())
	err = svr.B.SendMsg(identityRevokeMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:RevokeIdentityJoin failed")
		return nil, ErrSendIdentityRevokeMsg
	}
	log.Debug("RPC-API:RevokeIdentityJoin succeed SendMsg")
	return &apipb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *AuthServiceServer) GetNodeIdentity(ctx context.Context, req *emptypb.Empty) (*pb.GetNodeIdentityResponse, error) {
	identity, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetNodeIdentity failed")
		return nil, ErrGetNodeIdentity
	}
	return &pb.GetNodeIdentityResponse{
		Status: 0,
		Msg:    backend.OK,
		Owner: &apipb.Organization{
			NodeName:       identity.Name(),
			NodeId:     identity.NodeId(),
			IdentityId: identity.IdentityId(),
		},
	}, nil
}

func (svr *AuthServiceServer) GetIdentityList(ctx context.Context, req *emptypb.Empty) (*pb.GetIdentityListResponse, error) {
	identityList, err := svr.B.GetIdentityList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetIdentityList failed")
		return nil, ErrGetIdentityList
	}
	arr := make([]*apipb.Organization, len(identityList))
	for i, identity := range identityList {
		iden := &apipb.Organization{
			NodeName:       identity.Name(),
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

// 数据授权申请
func (svr *AuthServiceServer) ApplyMetaDataAuthority(context.Context, *pb.ApplyMetaDataAuthorityRequest) (*pb.ApplyMetaDataAuthorityResponse, error) {
	return nil, nil
}

// 数据授权审核
func (svr *AuthServiceServer) AuditMetaDataAuthority(context.Context, *pb.AuditMetaDataAuthorityRequest) (*apipb.SimpleResponse, error){
	return nil, nil
}

// 获取数据授权申请列表
func (svr *AuthServiceServer)  GetMetaDataAuthorityList(context.Context, *emptypb.Empty) (*pb.GetMetaDataAuthorityListResponse, error){
	return nil, nil
}
