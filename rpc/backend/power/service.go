package power

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (svr *Server) GetGlobalPowerSummaryList(ctx context.Context, req *emptypb.Empty) (*pb.GetGlobalPowerSummaryListResponse, error) {
	powerList, err := svr.B.GetGlobalPowerSummaryList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetGlobalPowerSummaryList failed")
		return nil, ErrGetTotalPowerList
	}
	//log.Debugf("RPC-API:GetGlobalPowerSummaryList succeed, powerList: {%d}, json: %s", len(powerList), utilGetGlobalPowerSummaryResponseArrString(powerList))
	log.Debugf("RPC-API:GetGlobalPowerSummaryList succeed, powerList: {%d}", len(powerList))
	return &pb.GetGlobalPowerSummaryListResponse{
		Status:    0,
		Msg:       backend.OK,
		PowerList: powerList,
	}, nil
}

func (svr *Server) GetGlobalPowerDetailList(ctx context.Context, req *pb.GetGlobalPowerDetailListRequest) (*pb.GetGlobalPowerDetailListResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = backend.DefaultPageSize
	}
	powerList, err := svr.B.GetGlobalPowerDetailList(req.GetLastUpdated(), pageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetGlobalPowerDetailList failed")
		return nil, ErrGetTotalPowerList
	}
	//log.Debugf("RPC-API:GetGlobalPowerDetailList succeed, powerList: {%d}, json: %s", len(powerList), utilGetGlobalPowerDetailResponseArrString(powerList))
	log.Debugf("RPC-API:GetGlobalPowerDetailList succeed, powerList: {%d}", len(powerList))
	return &pb.GetGlobalPowerDetailListResponse{
		Status:    0,
		Msg:       backend.OK,
		PowerList: powerList,
	}, nil
}

func (svr *Server) GetLocalPowerDetailList(ctx context.Context, req *emptypb.Empty) (*pb.GetLocalPowerDetailListResponse, error) {
	powerList, err := svr.B.GetLocalPowerDetailList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetLocalPowerDetailList failed")
		return nil, ErrGetSinglePowerList
	}
	//log.Debugf("RPC-API:GetLocalPowerDetailList succeed, powerList: {%d}, json: %s", len(powerList), utilGetLocalPowerDetailResponseArrString(powerList))
	log.Debugf("RPC-API:GetLocalPowerDetailList succeed, powerList: {%d}", len(powerList))
	return &pb.GetLocalPowerDetailListResponse{
		Status:    0,
		Msg:       backend.OK,
		PowerList: powerList,
	}, nil
}
func utilGetGlobalPowerSummaryResponseArrString(resp []*pb.GetGlobalPowerSummaryResponse) string {
	arr := make([]string, len(resp))
	for i, u := range resp {
		arr[i] = u.String()
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}
func utilGetGlobalPowerDetailResponseArrString(resp []*pb.GetGlobalPowerDetailResponse) string {
	arr := make([]string, len(resp))
	for i, u := range resp {
		arr[i] = u.String()
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}
func utilGetLocalPowerDetailResponseArrString(resp []*pb.GetLocalPowerDetailResponse) string {
	arr := make([]string, len(resp))
	for i, u := range resp {
		arr[i] = u.String()
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}

func (svr *Server) PublishPower(ctx context.Context, req *pb.PublishPowerRequest) (*pb.PublishPowerResponse, error) {
	if req == nil {
		return nil, ErrReqEmptyForPublishPower
	}

	if "" == strings.Trim(req.GetJobNodeId(), "") {
		log.Error("RPC-API:PublishPower failed, jobNodeId must be not empty")
		return nil, ErrSendPowerMsg
	}

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishPower failed, query local identity failed, can not publish power")
		return nil, ErrSendPowerMsg
	}

	powerMsg := types.NewPowerMessageFromRequest(req)
	powerId := powerMsg.GenPowerId()

	err = svr.B.SendMsg(powerMsg)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishPower failed, jobNodeId: {%s}, powerId: {%s}", req.GetJobNodeId(), powerId)
		errMsg := fmt.Sprintf("%s, jobNodeId:{%s}, powerId:{%s}", ErrSendPowerMsgByNidAndPowerId.Msg, req.GetJobNodeId(), powerId)
		return nil, backend.NewRpcBizErr(ErrSendPowerMsgByNidAndPowerId.Code, errMsg)
	}
	log.Debugf("RPC-API:PublishPower succeed, jobNodeId: {%s}, powerId: {%s}", req.GetJobNodeId(), powerId)
	return &pb.PublishPowerResponse{
		Status:  0,
		Msg:     backend.OK,
		PowerId: powerId,
	}, nil
}

func (svr *Server) RevokePower(ctx context.Context, req *pb.RevokePowerRequest) (*apicommonpb.SimpleResponse, error) {
	if req == nil {
		return nil, ErrReqEmptyForRevokePower
	}
	if "" == strings.Trim(req.GetPowerId(), "") {
		return nil, ErrReqEmptyPowerIdForRevokePower
	}

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokePower failed, query local identity failed, can not revoke power")
		return nil, ErrSendPowerRevokeMsg
	}

	// First check whether there is a task being executed on jobNode
	taskIdList, err := svr.B.QueryPowerRunningTaskList(req.GetPowerId())
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("RPC-API:RevokePower failed, query local running taskIdList failed, powerId: {%s}", req.GetPowerId())
		errMsg := fmt.Sprintf("query local running taskIdList failed, powerId:{%s}", req.GetPowerId())
		return nil, backend.NewRpcBizErr(ErrSendPowerRevokeMsgByPowerId.Code, errMsg)
	}
	if len(taskIdList) > 0 {
		log.WithError(err).Errorf("RPC-API:RevokePower failed, the old jobNode have been running {%d} task current, don't revoke it, powerId: {%s}", req.GetPowerId())
		errMsg := fmt.Sprintf("the old jobNode have been running {%d} task current, don't revoke it, powerId: {%s}", len(taskIdList), req.GetPowerId())
		return nil, backend.NewRpcBizErr(ErrSendPowerRevokeMsgByPowerId.Code, errMsg)
	}

	powerRevokeMsg := types.NewPowerRevokeMessageFromRequest(req)

	err = svr.B.SendMsg(powerRevokeMsg)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokePower failed, powerId: {%s}", req.GetPowerId())
		errMsg := fmt.Sprintf("%s, powerId:{%s}", ErrSendPowerRevokeMsgByPowerId.Msg, req.GetPowerId())
		return nil, backend.NewRpcBizErr(ErrSendPowerRevokeMsgByPowerId.Code, errMsg)
	}
	log.Debugf("RPC-API:RevokePower succeed, powerId: {%s}", req.GetPowerId())
	return &apicommonpb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}
