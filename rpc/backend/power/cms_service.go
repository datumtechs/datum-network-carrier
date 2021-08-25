package power

import (
	"context"
	"errors"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (svr *PowerServiceServer) GetPowerTotalDetailList(ctx context.Context, req *emptypb.Empty) (*pb.GetPowerTotalDetailListResponse, error) {
	powerList, err := svr.B.GetPowerTotalDetailList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetPowerTotalDetailList failed")
		return nil, ErrGetTotalPowerList
	}
	respList := make([]*pb.GetPowerTotalDetailResponse, len(powerList))

	for i, power := range powerList {
		resp := &pb.GetPowerTotalDetailResponse{
			Owner: types.ConvertNodeAliasToPB(power.Owner),
			Power: &libtypes.PowerTotalDetail{
				Information: types.ConvertResourceUsageToPB(power.PowerDetail.ResourceUsage),
				TotalTaskCount:power.PowerDetail.TotalTaskCount, // TODO 管理台查询 全网算力， 暂不显示 任务实况
				CurrentTaskCount: power.PowerDetail.CurrentTaskCount, // TODO 管理台查询 全网算力， 暂不显示 任务实况
				Tasks: types.ConvertPowerTaskArrToPB(power.PowerDetail.Tasks), // TODO 管理台查询 全网算力， 暂不显示 任务实况
				State: power.PowerDetail.State,
			},
		}
		respList[i] = resp
	}
	log.Debugf("RPC-API:GetPowerTotalDetailList succeed, powerList: {%d}, json: %s", len(respList), utilGetPowerTotalDetailResponseArrString(respList))
	return &pb.GetPowerTotalDetailListResponse{
		Status: 0,
		Msg: backend.OK,
		PowerList: respList,
	}, nil
}

func (svr *PowerServiceServer) GetPowerSingleDetailList(ctx context.Context, req *emptypb.Empty) (*pb.GetPowerSingleDetailListResponse, error) {
	powerList, err := svr.B.GetPowerSingleDetailList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetPowerSingleDetailList failed")
		return nil, ErrGetSinglePowerList
	}
	respList := make([]*pb.GetPowerSingleDetailResponse, len(powerList))

	for i, power := range powerList {

		resp := &pb.GetPowerSingleDetailResponse{
			Owner:  types.ConvertNodeAliasToPB(power.Owner),
			Power: &libtypes.PowerSingleDetail{
				JobNodeId: power.PowerDetail.JobNodeId,
				PowerId: power.PowerDetail.PowerId,
				Information:types.ConvertResourceUsageToPB(power.PowerDetail.ResourceUsage),
				TotalTaskCount:power.PowerDetail.TotalTaskCount,
				CurrentTaskCount: power.PowerDetail.CurrentTaskCount,
				Tasks: types.ConvertPowerTaskArrToPB(power.PowerDetail.Tasks),
				State: power.PowerDetail.State,
			},
		}
		respList[i] = resp
	}
	log.Debugf("RPC-API:GetPowerSingleDetailList succeed, powerList: {%d}, json: %s", len(respList), utilGetPowerSingleDetailResponseArrString(respList))
	return &pb.GetPowerSingleDetailListResponse{
		Status: 0,
		Msg: backend.OK,
		PowerList: respList,
	}, nil
}
func utilGetPowerTotalDetailResponseArrString(resp []*pb.GetPowerTotalDetailResponse) string {
	arr := make([]string, len(resp))
	for i, u := range resp {
		arr[i] = u.String()
	}
	if len(arr) != 0 {
		return "[" +  strings.Join(arr, ",") + "]"
	}
	return "[]"
}
func utilGetPowerSingleDetailResponseArrString(resp []*pb.GetPowerSingleDetailResponse) string {
	arr := make([]string, len(resp))
	for i, u := range resp {
		arr[i] = u.String()
	}
	if len(arr) != 0 {
		return "[" +  strings.Join(arr, ",") + "]"
	}
	return "[]"
}

func (svr *PowerServiceServer) PublishPower(ctx context.Context, req *pb.PublishPowerRequest) (*pb.PublishPowerResponse, error) {
	if req == nil {
		return nil, errors.New("required owner")
	}

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishPower failed, query local identity failed, can not publish power")
		return nil, ErrSendPowerMsg
	}

	//if identity.IdentityId() != req.Owner.IdentityId {
	//	return nil, errors.New("invalid identityId of req")
	//}
	//if identity.NodeId() != req.Owner.NodeId {
	//	return nil, errors.New("invalid nodeId of req")
	//}
	//if identity.Name() != req.Owner.Name {
	//	return nil, errors.New("invalid nodeName of req")
	//}


	powerMsg := types.NewPowerMessageFromRequest(req)
	powerId := powerMsg.SetPowerId()


	err = svr.B.SendMsg(powerMsg)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishPower failed, jobNodeId: {%s}, powerId: {%s}", req.JobNodeId, powerId)
		return nil, ErrSendPowerMsg
	}
	log.Debugf("RPC-API:PublishPower succeed, jobNodeId: {%s}, powerId: {%s}", req.JobNodeId, powerId)
	return &pb.PublishPowerResponse{
		Status:  0,
		Msg:     backend.OK,
		PowerId: powerId,
	}, nil
}

func (svr *PowerServiceServer) RevokePower(ctx context.Context, req *pb.RevokePowerRequest) (*apipb.SimpleResponse, error) {
	if req == nil {
		return nil, errors.New("required owner")
	}
	if req.PowerId == "" {
		return nil, errors.New("required powerId")
	}

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokePower failed, query local identity failed, can not revoke power")
		return nil, ErrSendPowerRevokeMsg
	}

	//if identity.IdentityId() != req.Owner.IdentityId {
	//	return nil, errors.New("invalid identityId of req")
	//}
	//if identity.NodeId() != req.Owner.NodeId {
	//	return nil, errors.New("invalid nodeId of req")
	//}
	//if identity.Name() != req.Owner.Name {
	//	return nil, errors.New("invalid nodeName of req")
	//}

	powerRevokeMsg := types.NewPowerRevokeMessageFromRequest(req)

	err = svr.B.SendMsg(powerRevokeMsg)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokePower failed, powerId: {%s}", req.PowerId)
		return nil, ErrSendPowerRevokeMsg
	}
	log.Debugf("RPC-API:RevokePower succeed, powerId: {%s}", req.PowerId)
	return &apipb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}
