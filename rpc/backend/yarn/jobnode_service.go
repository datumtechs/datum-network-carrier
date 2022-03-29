package yarn

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strings"

	"fmt"

	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
)

func (svr *Server) ReportTaskEvent(ctx context.Context, req *pb.ReportTaskEventRequest) (*libtypes.SimpleResponse, error) {

	if "" == strings.Trim(req.GetTaskEvent().GetTaskId(), "") {
		return &libtypes.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require taskId"}, nil
	}

	if "" == strings.Trim(req.GetTaskEvent().GetPartyId(), "") {
		return &libtypes.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require partyId"}, nil
	}

	log.Debugf("RPC-API:ReportTaskEvent, req: {%v}", req)

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportTaskEvent failed, query local identity failed, can not handle report event")
		return nil, fmt.Errorf("query local identity failed")
	}

	if err := svr.B.SendTaskEvent(req.GetTaskEvent()); nil != err {
		log.WithError(err).Error("RPC-API:ReportTaskEvent failed")

		errMsg := fmt.Sprintf("%s, %s", backend.ErrReportTaskEvent.Error(), req.GetTaskEvent().GetPartyId())
		return &libtypes.SimpleResponse{ Status: backend.ErrReportTaskEvent.ErrCode(), Msg: errMsg}, nil
	}
	return &libtypes.SimpleResponse{Status: 0, Msg: backend.OK}, nil
}


func (svr *Server) ReportTaskResourceUsage (ctx context.Context, req *pb.ReportTaskResourceUsageRequest) (*libtypes.SimpleResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportTaskResourceUsage failed, query local identity failed, can not handle report usage, taskId: {%s}, partyId: {%s}",
			req.GetTaskId(), req.GetPartyId())
		return &libtypes.SimpleResponse{ Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if "" == strings.Trim(req.GetTaskId(), "") {
		return &libtypes.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require taskId"}, nil
	}

	if "" == strings.Trim(req.GetPartyId(), "") {
		return &libtypes.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require partyId"}, nil
	}

	if "" == strings.Trim(req.GetIp(), "") || "" == strings.Trim(req.GetPort(), "") {
		return &libtypes.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require ip and port"}, nil
	}

	if req.GetNodeType() != pb.NodeType_NodeType_JobNode && req.GetNodeType() != pb.NodeType_NodeType_DataNode {
		return &libtypes.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown nodeType"}, nil
	}

	if nil == req.GetUsage() {
		return &libtypes.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown resourceUsage"}, nil
	}

	if err := svr.B.ReportTaskResourceUsage(req.GetNodeType(), req.GetIp(), req.GetPort(),
		types.NewTaskResuorceUsage(
			req.GetTaskId(),
			req.GetPartyId(),
			req.GetUsage().GetTotalMem(),
			req.GetUsage().GetTotalBandwidth(),
			req.GetUsage().GetTotalDisk(),
			req.GetUsage().GetUsedMem(),
			req.GetUsage().GetUsedBandwidth(),
			req.GetUsage().GetUsedDisk(),
			req.GetUsage().GetTotalProcessor(),
			req.GetUsage().GetUsedProcessor())); nil != err {
		log.WithError(err).Error("RPC-API:ReportTaskResourceUsage failed")
		return &libtypes.SimpleResponse{ Status: backend.ErrReportTaskResourceExpense.ErrCode(), Msg: backend.ErrReportTaskResourceExpense.Error()}, nil
	}

	return &libtypes.SimpleResponse{
		Status: 0,
		Msg: backend.OK,
	}, nil
}

