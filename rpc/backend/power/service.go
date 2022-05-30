package power

import (
	"context"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/core/rawdb"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/rpc/backend"
	"github.com/datumtechs/datum-network-carrier/types"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (svr *Server) GetGlobalPowerSummaryList(ctx context.Context, req *emptypb.Empty) (*carrierapipb.GetGlobalPowerSummaryListResponse, error) {
	powerList, err := svr.B.GetGlobalPowerSummaryList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetGlobalPowerSummaryList failed")
		return &carrierapipb.GetGlobalPowerSummaryListResponse { Status: backend.ErrQueryGlobalPowerList.ErrCode(), Msg: backend.ErrQueryGlobalPowerList.Error()}, nil
	}
	log.Debugf("RPC-API:GetGlobalPowerSummaryList succeed, powerList: {%d}", len(powerList))
	return &carrierapipb.GetGlobalPowerSummaryListResponse{
		Status:    0,
		Msg:       backend.OK,
		Powers: powerList,
	}, nil
}

func (svr *Server) GetGlobalPowerDetailList(ctx context.Context, req *carrierapipb.GetGlobalPowerDetailListRequest) (*carrierapipb.GetGlobalPowerDetailListResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = backend.DefaultPageSize
	}
	powerList, err := svr.B.GetGlobalPowerDetailList(req.GetLastUpdated(), pageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetGlobalPowerDetailList failed")
		return &carrierapipb.GetGlobalPowerDetailListResponse { Status: backend.ErrQueryGlobalPowerList.ErrCode(), Msg: backend.ErrQueryGlobalPowerList.Error()}, nil
	}
	log.Debugf("RPC-API:GetGlobalPowerDetailList succeed, powerList: {%d}", len(powerList))
	return &carrierapipb.GetGlobalPowerDetailListResponse{
		Status:    0,
		Msg:       backend.OK,
		Powers: powerList,
	}, nil
}

func (svr *Server) GetLocalPowerDetailList(ctx context.Context, req *emptypb.Empty) (*carrierapipb.GetLocalPowerDetailListResponse, error) {
	powerList, err := svr.B.GetLocalPowerDetailList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetLocalPowerDetailList failed")
		return &carrierapipb.GetLocalPowerDetailListResponse { Status: backend.ErrQueryLocalPowerList.ErrCode(), Msg: backend.ErrQueryLocalPowerList.Error()}, nil
	}
	log.Debugf("RPC-API:GetLocalPowerDetailList succeed, powerList: {%d}", len(powerList))
	return &carrierapipb.GetLocalPowerDetailListResponse{
		Status:    0,
		Msg:       backend.OK,
		Powers: powerList,
	}, nil
}
func utilGetGlobalPowerSummaryResponseArrString(resp []*carrierapipb.GetGlobalPowerSummary) string {
	arr := make([]string, len(resp))
	for i, u := range resp {
		arr[i] = u.String()
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}
func utilGetGlobalPowerDetailResponseArrString(resp []*carrierapipb.GetGlobalPowerDetail) string {
	arr := make([]string, len(resp))
	for i, u := range resp {
		arr[i] = u.String()
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}
func utilGetLocalPowerDetailResponseArrString(resp []*carrierapipb.GetLocalPowerDetail) string {
	arr := make([]string, len(resp))
	for i, u := range resp {
		arr[i] = u.String()
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}

func (svr *Server) PublishPower(ctx context.Context, req *carrierapipb.PublishPowerRequest) (*carrierapipb.PublishPowerResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishPower failed, query local identity failed, can not publish power, jonNodeId: {%s}", req.GetJobNodeId())
		return &carrierapipb.PublishPowerResponse { Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if nil == req {
		return &carrierapipb.PublishPowerResponse { Status: backend.ErrRequireParams.ErrCode(), Msg: backend.ErrRequireParams.Error()}, nil
	}

	if "" == strings.Trim(req.GetJobNodeId(), "") {
		log.Error("RPC-API:PublishPower failed, jobNodeId must be not empty")
		return &carrierapipb.PublishPowerResponse { Status: backend.ErrRequireParams.ErrCode(), Msg: "require jobNodeId"}, nil
	}

	jobNode, err := svr.B.GetRegisterNode(carrierapipb.PrefixTypeJobNode, req.GetJobNodeId())
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("RPC-API:PublishPower failed, query jobNode failed, can not publish power, jonNodeId: {%s}", req.GetJobNodeId())
		return &carrierapipb.PublishPowerResponse { Status: backend.ErrPublishPowerMsg.ErrCode(), Msg: "query jobNode failed"}, nil
	}
	if jobNode.GetConnState() != carrierapipb.ConnState_ConnState_Connected {
		log.WithError(err).Errorf("RPC-API:PublishPower failed, jobNode was not connected, can not publish power, jonNodeId: {%s}， connState: {%s}",
			req.GetJobNodeId(), jobNode.GetConnState().String())
		return &carrierapipb.PublishPowerResponse { Status: backend.ErrPublishPowerMsg.ErrCode(), Msg: "jobNode was not connected"}, nil
	}

	powerMsg := types.NewPowerMessageFromRequest(req)
	powerId := powerMsg.GenPowerId()

	if err = svr.B.SendMsg(powerMsg); nil != err {
		log.WithError(err).Errorf("RPC-API:PublishPower failed, jobNodeId: {%s}, powerId: {%s}", req.GetJobNodeId(), powerId)
		errMsg := fmt.Sprintf("%s, jobNodeId:{%s}, powerId:{%s}", backend.ErrPublishPowerMsg.Error(), req.GetJobNodeId(), powerId)
		return &carrierapipb.PublishPowerResponse { Status: backend.ErrPublishPowerMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:PublishPower succeed, jobNodeId: {%s}, powerId: {%s}", req.GetJobNodeId(), powerId)
	return &carrierapipb.PublishPowerResponse{
		Status:  0,
		Msg:     backend.OK,
		PowerId: powerId,
	}, nil
}

func (svr *Server) RevokePower(ctx context.Context, req *carrierapipb.RevokePowerRequest) (*carriertypespb.SimpleResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokePower failed, query local identity failed, can not revoke power")
		return &carriertypespb.SimpleResponse{ Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if nil == req {
		return &carriertypespb.SimpleResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: backend.ErrRequireParams.Error()}, nil
	}
	if "" == strings.Trim(req.GetPowerId(), "") {
		return &carriertypespb.SimpleResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "require powerId"}, nil
	}

	// First check whether there is a task being executed on jobNode
	taskIdList, err := svr.B.QueryPowerRunningTaskList(req.GetPowerId())
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("RPC-API:RevokePower failed, query local running taskIdList failed, powerId: {%s}", req.GetPowerId())
		errMsg := fmt.Sprintf("query local running taskIdList failed, powerId:{%s}", req.GetPowerId())
		return &carriertypespb.SimpleResponse{ Status: backend.ErrRevokePowerMsg.ErrCode(), Msg: errMsg}, nil
	}
	if len(taskIdList) > 0 {
		log.WithError(err).Errorf("RPC-API:RevokePower failed, the old jobNode have been running {%d} task current, don't revoke it, powerId: {%s}", len(taskIdList), req.GetPowerId())
		errMsg := fmt.Sprintf("the old jobNode have been running {%d} task current, don't revoke it, powerId: {%s}", len(taskIdList), req.GetPowerId())
		return &carriertypespb.SimpleResponse{ Status: backend.ErrRevokePowerMsg.ErrCode(), Msg: errMsg}, nil
	}

	powerRevokeMsg := types.NewPowerRevokeMessageFromRequest(req)

	if err = svr.B.SendMsg(powerRevokeMsg); nil != err {
		log.WithError(err).Errorf("RPC-API:RevokePower failed, powerId: {%s}", req.GetPowerId())
		errMsg := fmt.Sprintf("%s, powerId:{%s}", backend.ErrRevokePowerMsg.Error(), req.GetPowerId())
		return &carriertypespb.SimpleResponse{ Status: backend.ErrRevokePowerMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:RevokePower succeed, powerId: {%s}", req.GetPowerId())
	return &carriertypespb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}
