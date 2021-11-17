package yarn

import (
	"context"
	"strings"

	"errors"

	"fmt"

	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func (svr *Server) ReportTaskEvent(ctx context.Context, req *pb.ReportTaskEventRequest) (*apicommonpb.SimpleResponse, error) {

	if "" == strings.Trim(req.GetTaskEvent().GetTaskId(), "") {
		return nil, backend.NewRpcBizErr(ErrReportTaskEvent.Code, "require taskId")
	}

	if "" == strings.Trim(req.GetTaskEvent().GetPartyId(), "") {
		return nil, backend.NewRpcBizErr(ErrReportTaskEvent.Code, "require partyId")
	}

	log.Debugf("RPC-API:ReportTaskEvent, req: {%v}", req)
	err := svr.B.SendTaskEvent(req.GetTaskEvent())
	if nil != err {
		log.WithError(err).Error("RPC-API:ReportTaskEvent failed")

		errMsg := fmt.Sprintf("%s, %s", ErrReportTaskEvent.Msg, req.GetTaskEvent().GetPartyId())
		return nil, backend.NewRpcBizErr(ErrReportTaskEvent.Code, errMsg)
	}
	return &apicommonpb.SimpleResponse{Status: 0, Msg: backend.OK}, nil
}


func (svr *Server) ReportTaskResourceUsage (ctx context.Context, req *pb.ReportTaskResourceUsageRequest) (*apicommonpb.SimpleResponse, error) {

	if req.GetTaskId() == "" {
		return nil, errors.New("require taskId")
	}

	if req.GetPartyId() == "" {
		return nil, errors.New("require partyId")
	}

	if req.GetIp() == "" || req.GetPort() == "" {
		return nil, errors.New("require ip and port")
	}

	if req.GetNodeType() != pb.NodeType_NodeType_JobNode && req.GetNodeType() != pb.NodeType_NodeType_DataNode {
		return nil, errors.New("invalid node type")
	}

	if nil == req.GetUsage() {
		return nil, errors.New("require resource usage")
	}

	err := svr.B.ReportTaskResourceUsage(req.GetNodeType(), req.GetIp(), req.GetPort(),
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
			req.GetUsage().GetUsedProcessor()))
	if nil != err {
		log.WithError(err).Error("RPC-API:ReportTaskResourceUsage failed")
		return nil, ErrReportTaskResourceExpense
	}

	return &apicommonpb.SimpleResponse{
		Status: 0,
		Msg: backend.OK,
	}, nil
}

