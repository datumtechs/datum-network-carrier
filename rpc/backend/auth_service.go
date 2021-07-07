package backend

import (
	"context"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"github.com/RosettaFlow/Carrier-Go/types"
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
	identityMsg.NodeAlias = &types.NodeAlias{}
	identityMsg.Name = req.Member.Name
	identityMsg.IdentityId = req.Member.IdentityId
	identityMsg.NodeId = req.Member.NodeId
	identityMsg.CreateAt = uint64(time.Now().UnixNano())

	err := svr.B.SendMsg(identityMsg)
	if nil != err {
		return nil, NewRpcBizErr(ErrSendIdentityMsgStr)
	}
	return &pb.SimpleResponseCode{
		Status: 0,
		Msg:    OK,
	}, nil
}

func (svr *AuthServiceServer) RevokeIdentityJoin(ctx context.Context, req *pb.EmptyGetParams) (*pb.SimpleResponseCode, error) {

	identityRevokeMsg := new(types.IdentityRevokeMsg)
	identityRevokeMsg.CreateAt = uint64(time.Now().UnixNano())
	err := svr.B.SendMsg(identityRevokeMsg)
	if nil != err {
		return nil, NewRpcBizErr(ErrSendIdentityMsgStr)
	}
	return &pb.SimpleResponseCode{
		Status: 0,
		Msg:    OK,
	}, nil
}

func (svr *AuthServiceServer) GetNodeIdentity(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetNodeIdentityResponse, error) {

	identity, err := svr.B.GetNodeIdentity()
	if nil != err {
		return nil, NewRpcBizErr(ErrGetNodeIdentityStr)
	}
	return &pb.GetNodeIdentityResponse{
		Status: 0,
		Msg:    OK,
		Owner: &pb.OrganizationIdentityInfo{
			Name:       identity.Name(),
			NodeId:     identity.NodeId(),
			IdentityId: identity.IdentityId(),
		},
	}, nil
}

func (svr *AuthServiceServer) GetIdentityList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetIdentityListResponse, error) {
	identitys, err := svr.B.GetIdentityList()
	if nil != err {
		return nil, NewRpcBizErr(ErrGetIdentityListStr)
	}
	arr := make([]*pb.OrganizationIdentityInfo, len(identitys))
	for i, identity := range identitys {
		iden := &pb.OrganizationIdentityInfo{
			Name:       identity.Name(),
			NodeId:     identity.NodeId(),
			IdentityId: identity.IdentityId(),
		}
		arr[i] = iden
	}
	return &pb.GetIdentityListResponse{
		Status:     0,
		Msg:        OK,
		MemberList: arr,
	}, nil
}
