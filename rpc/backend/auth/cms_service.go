package auth

import (
	"context"
	"errors"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strings"
	"time"
)

func (svr *AuthServiceServer) ApplyIdentityJoin(ctx context.Context, req *pb.ApplyIdentityJoinRequest) (*pb.SimpleResponseCode, error) {

	identity, err := svr.B.GetNodeIdentity()
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("RPC-API:ApplyIdentityJoin failed, query local identity failed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}",
			req.Member.IdentityId, req.Member.NodeId, req.Member.Name)
		return nil, backend.NewRpcBizErr(ErrSendIdentityMsgStr)
	}

	if nil != identity {
		log.Errorf("RPC-API:ApplyIdentityJoin failed, identity was already exist, old identityId: {%s}, old nodeId: {%s}, old nodeName: {%s}",
			identity.IdentityId(), identity.NodeId(), identity.Name())
		return nil, backend.NewRpcBizErr(ErrSendIdentityMsgStr)
	}

	identityMsg := new(types.IdentityMsg)
	if req.Member == nil {
		return nil, errors.New("Invalid Params, req.Member is nil")
	}

	if "" == strings.Trim(req.Member.IdentityId, "") ||
		"" == strings.Trim(req.Member.Name, "") {
		return nil, errors.New("Invalid Params, req.Member.IdentityId or req.Member.Name is empty")
	}

	identityMsg.NodeAlias = &types.NodeAlias{}
	identityMsg.Name = req.Member.Name
	identityMsg.IdentityId = req.Member.IdentityId
	//identityMsg.NodeId = req.Member.NodeId
	identityMsg.CreateAt = uint64(time.Now().UnixNano())

	err = svr.B.SendMsg(identityMsg)
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

	_, err := svr.B.GetNodeIdentity()
	if rawdb.IsDBNotFoundErr(err) {
		log.Errorf("RPC-API:RevokeIdentityJoin failed, the identity was not exist, can not revoke identity")
		return nil, backend.NewRpcBizErr(ErrSendIdentityRevokeMsgStr)
	}

	identityRevokeMsg := new(types.IdentityRevokeMsg)
	identityRevokeMsg.CreateAt = uint64(time.Now().UnixNano())
	err = svr.B.SendMsg(identityRevokeMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:RevokeIdentityJoin failed")
		return nil, backend.NewRpcBizErr(ErrSendIdentityRevokeMsgStr)
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
