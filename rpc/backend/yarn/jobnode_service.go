package yarn

import (
	"context"

	"errors"

	"fmt"

	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func (svr *Server) ReportTaskEvent(ctx context.Context, req *pb.ReportTaskEventRequest) (*apicommonpb.SimpleResponse, error) {

	log.Debugf("RPC-API:ReportTaskEvent, req: {%v}", req)
	err := svr.B.SendTaskEvent(types.NewReportTaskEvent(req.PartyId, req.GetTaskEvent()))
	if nil != err {
		log.WithError(err).Error("RPC-API:ReportTaskEvent failed")

		errMsg := fmt.Sprintf(ErrReportTaskEvent.Msg, req.PartyId)
		return nil, backend.NewRpcBizErr(ErrReportTaskEvent.Code, errMsg)
	}
	return &apicommonpb.SimpleResponse{Status: 0, Msg: backend.OK}, nil
}




func (svr *Server) ReportTaskResourceExpense (ctx context.Context, req *pb.ReportTaskResourceExpenseRequest) (*apicommonpb.SimpleResponse, error) {

	if req.GetTaskId() == "" {
		return nil, errors.New("require taskId")
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

	err := svr.B.ReportTaskResourceExpense(req.GetNodeType(), req.GetTaskId(), req.GetIp(), req.GetPort(), req.GetUsage())
	if nil != err {
		log.WithError(err).Error("RPC-API:ReportTaskResourceExpense failed")
		return nil, ErrReportTaskResourceExpense
	}

	return &apicommonpb.SimpleResponse{
		Status: 0,
		Msg: backend.OK,
	}, nil
}

