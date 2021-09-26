package yarn

import (
	"context"
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

func (svr *Server) ReportTaskResourceExpense(ctx context.Context, req *pb.ReportTaskResourceExpenseRequest) (*apicommonpb.SimpleResponse, error) {
	return nil, nil
}
