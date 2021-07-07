package backend

import (
	"context"
	"errors"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"github.com/RosettaFlow/Carrier-Go/types"
	"time"
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
		return nil, NewRpcBizErr(ErrGetTotalPowerListStr)
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
	return &pb.GetPowerTotalDetailListResponse{
		PowerList: respList,
	}, nil
}

func (svr *PowerServiceServer) GetPowerSingleDetailList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetPowerSingleDetailListResponse, error) {

	powerList, err := svr.B.GetPowerSingleDetailList()
	if nil != err {
		return nil, NewRpcBizErr(ErrGetSinglePowerListStr)
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
	return &pb.GetPowerSingleDetailListResponse{
		PowerList: respList,
	}, nil
}

func (svr *PowerServiceServer) PublishPower(ctx context.Context, req *pb.PublishPowerRequest) (*pb.PublishPowerResponse, error) {

	if req == nil || req.Owner == nil {
		return nil, errors.New("required owner")
	}
	if req.Information == nil {
		return nil, errors.New("required information")
	}
	
	powerMsg := types.NewPowerMessageFromRequest(req)
	/*powerMsg.Data.JobNodeId = req.JobNodeId
	powerMsg.Data.CreateAt = uint64(time.Now().UnixNano())
	powerMsg.Data.Name = req.Owner.Name
	powerMsg.Data.NodeId = req.Owner.NodeId
	powerMsg.Data.IdentityId = req.Owner.IdentityId
	powerMsg.Data.Information.Processor = req.Information.Processor
	powerMsg.Data.Information.Mem = req.Information.Mem
	powerMsg.Data.Information.Bandwidth = req.Information.Bandwidth*/
	powerId := powerMsg.GetPowerId()

	err := svr.B.SendMsg(powerMsg)
	if nil != err {
		return nil, NewRpcBizErr(ErrSendPowerMsgStr)
	}
	return &pb.PublishPowerResponse{
		Status:  0,
		Msg:     OK,
		PowerId: powerId,
	}, nil
}

func (svr *PowerServiceServer) RevokePower(ctx context.Context, req *pb.RevokePowerRequest) (*pb.SimpleResponseCode, error) {
	if req == nil || req.Owner == nil {
		return nil, errors.New("required owner")
	}
	if req.PowerId == "" {
		return nil, errors.New("required powerId")
	}
	powerRevokeMsg := types.NewPowerRevokeMessageFromRequest(req)
	powerRevokeMsg.CreateAt = uint64(time.Now().UnixNano())

	/*powerRevokeMsg.Name = req.Owner.Name
	powerRevokeMsg.NodeId = req.Owner.NodeId
	powerRevokeMsg.IdentityId = req.Owner.IdentityId
	powerRevokeMsg.PowerId = req.PowerId*/

	err := svr.B.SendMsg(powerRevokeMsg)
	if nil != err {
		return nil, NewRpcBizErr(ErrSendPowerRevokeMsgStr)
	}
	return &pb.SimpleResponseCode{
		Status: 0,
		Msg:    OK,
	}, nil
}
