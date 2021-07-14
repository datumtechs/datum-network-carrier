package yarn

import (
	"context"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strings"
)

func (svr *YarnServiceServer) ReportUpFileSummary(ctx context.Context, req *pb.ReportUpFileSummaryRequest) (*pb.SimpleResponseCode, error) {
	jobNodeList, err := svr.B.GetRegisterNodeList(types.PREFIX_TYPE_JOBNODE)
	if nil != err {
		return nil, backend.NewRpcBizErr(ErrGetJobNodeListStr)
	}
	var resourceId string
	for _, jobNode := range jobNodeList {
		if req.Ip == jobNode.InternalIp && req.Port == jobNode.InternalPort {
			resourceId = jobNode.Id
			break
		}
	}
	if "" == strings.Trim(resourceId, "") {
		return nil, backend.NewRpcBizErr(ErrGetJobNodeListStr)
	}
	err = svr.B.StoreDataResourceDataUsed(types.NewDataResourceDataUsed(resourceId, req.OriginId, "", req.FilePath))
	if nil != err {
		return nil, backend.NewRpcBizErr(ErrReportUpFileSummaryStr)
	}
	return &pb.SimpleResponseCode{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}
