package yarn

import (
	"context"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
)

func (svr *YarnServiceServer) ReportTaskEvent(ctx context.Context, req *pb.ReportTaskEventRequest) (*apipb.SimpleResponse, error) {
	log.Debugf("RPC-API:ReportTaskEvent, req: {%v}", req)
	err := svr.B.SendTaskEvent(&libTypes.TaskEvent{
		Type:       req.TaskEvent.Type,
		IdentityId:   req.TaskEvent.IdentityId,
		TaskId:     req.TaskEvent.TaskId,
		Content:    req.TaskEvent.Content,
		CreateAt: req.TaskEvent.CreateAt,
	})
	if nil != err {
		log.WithError(err).Error("RPC-API:ReportTaskEvent failed")
		return nil, ErrReportTaskEvent
	}
	return &apipb.SimpleResponse{Status: 0, Msg: backend.OK}, nil
}

func (svr *YarnServiceServer) ReportTaskResourceExpense(ctx context.Context, req *pb.ReportTaskResourceExpenseRequest) (*apipb.SimpleResponse, error) {
	return nil, nil
}