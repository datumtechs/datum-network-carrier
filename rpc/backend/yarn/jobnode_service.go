package yarn

import (
	"context"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func (svr *Server) ReportTaskEvent(ctx context.Context, req *pb.ReportTaskEventRequest) (*apipb.SimpleResponse, error) {
	log.Debugf("RPC-API:ReportTaskEvent, req: {%v}", req)
	err := svr.B.SendTaskEvent(types.NewReportTaskEvent(req.PartyId, req.GetTaskEvent()))
	if nil != err {
		log.WithError(err).Error("RPC-API:ReportTaskEvent failed")
		return nil, ErrReportTaskEvent
	}
	return &apipb.SimpleResponse{Status: 0, Msg: backend.OK}, nil
}

func (svr *Server) ReportTaskResourceExpense(ctx context.Context, req *pb.ReportTaskResourceExpenseRequest) (*apipb.SimpleResponse, error) {
	return nil, nil
}