package yarn

import (
	"context"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func (svr *YarnServiceServer) ReportTaskEvent(ctx context.Context, req *pb.ReportTaskEventRequest) (*pb.SimpleResponseCode, error) {
	var err error
	go func() {
		err = svr.B.SendTaskEvent(&types.TaskEventInfo{
			Type:       req.TaskEvent.Type,
			Identity:   req.TaskEvent.IdentityId,
			TaskId:     req.TaskEvent.TaskId,
			Content:    req.TaskEvent.Content,
			CreateTime: req.TaskEvent.CreateAt,
		})
	}()
	if nil != err {
		return nil, backend.NewRpcBizErr(ErrReportTaskEventStr)
	}
	return &pb.SimpleResponseCode{Status: 0, Msg: backend.OK}, nil
}

func (svr *YarnServiceServer) ReportTaskResourceExpense(ctx context.Context, req *pb.ReportTaskResourceExpenseRequest) (*pb.SimpleResponseCode, error) {
	return nil, nil
}