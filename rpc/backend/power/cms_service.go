package power

import (
	"context"
	"errors"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
)

//
//func (svr *PowerServiceServer) GetPowerTotalSummaryList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetPowerTotalSummaryListResponse, error) {
//	return nil, nil
//}
//func (svr *PowerServiceServer) GetPowerSingleSummaryList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetPowerSingleSummaryListResponse, error) {
//	return nil, nil
//}
//func (svr *PowerServiceServer) GetPowerTotalSummaryByState(ctx context.Context, req *pb.GetPowerTotalSummaryByStateRequest) (*pb.GetPowerTotalSummaryListResponse, error) {
//	return nil, nil
//}
//func (svr *PowerServiceServer) GetPowerSingleSummaryByState(ctx context.Context, req *pb.GetPowerSingleSummaryByStateRequest) (*pb.GetPowerSingleSummaryListResponse, error) {
//	return nil, nil
//}
//func (svr *PowerServiceServer) GetPowerTotalSummaryByOwner(ctx context.Context, req *pb.GetPowerTotalSummaryByOwnerRequest) (*pb.GetPowerTotalSummaryResponse, error) {
//	return nil, nil
//}
//func (svr *PowerServiceServer) GetPowerSingleSummaryByOwner(ctx context.Context, req *pb.GetPowerSingleSummaryByOwnerRequest) (*pb.GetPowerSingleSummaryListResponse, error) {
//	return nil, nil
//}
//func (svr *PowerServiceServer) GetPowerSingleDetail(ctx context.Context, req *pb.GetPowerSingleDetailRequest) (*pb.GetPowerSingleDetailResponse, error) {
//	return nil, nil
//}
func (svr *PowerServiceServer) GetPowerTotalDetailList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetPowerTotalDetailListResponse, error) {
	powerList, err := svr.B.GetPowerTotalDetailList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetPowerTotalDetailList failed")
		return nil, backend.NewRpcBizErr(ErrGetTotalPowerListStr)
	}
	respList := make([]*pb.GetPowerTotalDetailResponse, len(powerList))

	for i, power := range powerList {
		resp := &pb.GetPowerTotalDetailResponse{
			Owner: types.ConvertNodeAliasToPB(power.Owner),
			Power: &pb.PowerTotalDetail{
				Information: types.ConvertResourceUsageToPB(power.PowerDetail.ResourceUsage),
				TotalTaskCount:power.PowerDetail.TotalTaskCount, // TODO 管理台查询 全网算力， 暂不显示 任务实况
				CurrentTaskCount: power.PowerDetail.CurrentTaskCount, // TODO 管理台查询 全网算力， 暂不显示 任务实况
				Tasks: types.ConvertPowerTaskArrToPB(power.PowerDetail.Tasks), // TODO 管理台查询 全网算力， 暂不显示 任务实况
				State: power.PowerDetail.State,
			},
		}
		respList[i] = resp
	}
	log.Debugf("RPC-API:GetPowerTotalDetailList succeed, powerList: {%d}", len(respList))
	return &pb.GetPowerTotalDetailListResponse{
		Status: 0,
		Msg: backend.OK,
		PowerList: respList,
	}, nil
}

func (svr *PowerServiceServer) GetPowerSingleDetailList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetPowerSingleDetailListResponse, error) {
	powerList, err := svr.B.GetPowerSingleDetailList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetPowerSingleDetailList failed")
		return nil, backend.NewRpcBizErr(ErrGetSinglePowerListStr)
	}
	respList := make([]*pb.GetPowerSingleDetailResponse, len(powerList))

	for i, power := range powerList {

		resp := &pb.GetPowerSingleDetailResponse{
			Owner:  types.ConvertNodeAliasToPB(power.Owner),
			Power: &pb.PowerSingleDetail{
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
	log.Debugf("RPC-API:GetPowerSingleDetailList succeed, powerList: {%d}", len(respList))
	return &pb.GetPowerSingleDetailListResponse{
		Status: 0,
		Msg: backend.OK,
		PowerList: respList,
	}, nil
}

func (svr *PowerServiceServer) PublishPower(ctx context.Context, req *pb.PublishPowerRequest) (*pb.PublishPowerResponse, error) {
	if req == nil {
		return nil, errors.New("required owner")
	}

	_, err := svr.B.GetNodeIdentity()
	if rawdb.IsDBNotFoundErr(err) {
		log.Errorf("RPC-API:PublishPower failed, the identity was not exist, can not revoke identity")
		return nil, backend.NewRpcBizErr(ErrSendPowerMsgStr)
	}

	powerMsg := types.NewPowerMessageFromRequest(req)
	powerId := powerMsg.GetPowerId()

	err = svr.B.SendMsg(powerMsg)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishPower failed, jobNodeId: {%s}, powerId: {%s}", req.JobNodeId, powerId)
		return nil, backend.NewRpcBizErr(ErrSendPowerMsgStr)
	}
	log.Debugf("RPC-API:PublishPower succeed, jobNodeId: {%s}, powerId: {%s}", req.JobNodeId, powerId)
	return &pb.PublishPowerResponse{
		Status:  0,
		Msg:     backend.OK,
		PowerId: powerId,
	}, nil
}

func (svr *PowerServiceServer) RevokePower(ctx context.Context, req *pb.RevokePowerRequest) (*pb.SimpleResponseCode, error) {
	if req == nil {
		return nil, errors.New("required owner")
	}
	if req.PowerId == "" {
		return nil, errors.New("required powerId")
	}

	_, err := svr.B.GetNodeIdentity()
	if rawdb.IsDBNotFoundErr(err) {
		log.Errorf("RPC-API:RevokePower failed, the identity was not exist, can not revoke identity")
		return nil, backend.NewRpcBizErr(ErrSendPowerRevokeMsgStr)
	}

	powerRevokeMsg := types.NewPowerRevokeMessageFromRequest(req)

	err = svr.B.SendMsg(powerRevokeMsg)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokePower failed, powerId: {%s}", req.PowerId)
		return nil, backend.NewRpcBizErr(ErrSendPowerRevokeMsgStr)
	}
	log.Debugf("RPC-API:RevokePower succeed, powerId: {%s}", req.PowerId)
	return &pb.SimpleResponseCode{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}
