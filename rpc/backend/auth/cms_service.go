package auth

import (
	"context"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strings"
	"time"
)

func (svr *AuthServiceServer) ApplyIdentityJoin(ctx context.Context, req *pb.ApplyIdentityJoinRequest) (*pb.SimpleResponseCode, error) {
	identityMsg := new(types.IdentityMsg)
	if req.Member == nil {
		return &pb.SimpleResponseCode{
			Status: 0,
			Msg:    "Invalid Params",
		}, nil
	}

	if "" == strings.Trim(req.Member.IdentityId, "") ||
		"" == strings.Trim(req.Member.NodeId, "") ||
		"" == strings.Trim(req.Member.Name, "") {
		return &pb.SimpleResponseCode{
			Status: 0,
			Msg:    "Invalid Params",
		}, nil
	}

	identityMsg.NodeAlias = &types.NodeAlias{}
	identityMsg.Name = req.Member.Name
	identityMsg.IdentityId = req.Member.IdentityId
	identityMsg.NodeId = req.Member.NodeId
	identityMsg.CreateAt = uint64(time.Now().UnixNano())

	err := svr.B.SendMsg(identityMsg)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ApplyIdentityJoin failed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}",
			req.Member.IdentityId, req.Member.NodeId, req.Member.Name)
		return nil, backend.NewRpcBizErr(ErrSendIdentityMsgStr)
	}
	log.Debugf("RPC-API:ApplyIdentityJoin succeed SendMsg, identityId: {%s}, nodeId: {%s}, nodeName: {%s}",
		req.Member.IdentityId, req.Member.NodeId, req.Member.Name)
	return &pb.SimpleResponseCode{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *AuthServiceServer) RevokeIdentityJoin(ctx context.Context, req *pb.EmptyGetParams) (*pb.SimpleResponseCode, error) {
	identityRevokeMsg := new(types.IdentityRevokeMsg)
	identityRevokeMsg.CreateAt = uint64(time.Now().UnixNano())
	err := svr.B.SendMsg(identityRevokeMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:RevokeIdentityJoin failed")
		return nil, backend.NewRpcBizErr(ErrSendIdentityMsgStr)
	}
	log.Debug("RPC-API:RevokeIdentityJoin succeed SendMsg")
	return &pb.SimpleResponseCode{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *AuthServiceServer) GetNodeIdentity(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetNodeIdentityResponse, error) {
	identity, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetNodeIdentity failed")
		return nil, backend.NewRpcBizErr(ErrGetNodeIdentityStr)
	}
	return &pb.GetNodeIdentityResponse{
		Status: 0,
		Msg:    backend.OK,
		Owner: &pb.OrganizationIdentityInfo{
			Name:       identity.Name(),
			NodeId:     identity.NodeId(),
			IdentityId: identity.IdentityId(),
		},
	}, nil
}

func (svr *AuthServiceServer) GetIdentityList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetIdentityListResponse, error) {
	identityList, err := svr.B.GetIdentityList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetIdentityList failed")
		return nil, backend.NewRpcBizErr(ErrGetIdentityListStr)
	}
	arr := make([]*pb.OrganizationIdentityInfo, len(identityList))
	for i, identity := range identityList {
		iden := &pb.OrganizationIdentityInfo{
			Name:       identity.Name(),
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
