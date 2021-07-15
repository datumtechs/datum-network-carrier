package yarn

import (
	"context"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strings"
)

func (svr *YarnServiceServer) ReportUpFileSummary(ctx context.Context, req *pb.ReportUpFileSummaryRequest) (*pb.SimpleResponseCode, error) {
	dataNodeList, err := svr.B.GetRegisterNodeList(types.PREFIX_TYPE_DATANODE)
	if nil != err {
		return nil, backend.NewRpcBizErr(ErrGetDataNodeListStr)
	}
	var resourceId string
	for _, dataNode := range dataNodeList {
		if req.Ip == dataNode.InternalIp && req.Port == dataNode.InternalPort {
			resourceId = dataNode.Id
			break
		}
	}
	if "" == strings.Trim(resourceId, "") {
		return nil, backend.NewRpcBizErr(ErrGetDataNodeListStr)
	}
	err = svr.B.StoreDataResourceFileUpload(types.NewDataResourceFileUpload(resourceId, req.OriginId, "", req.FilePath))
	if nil != err {
		return nil, backend.NewRpcBizErr(ErrReportUpFileSummaryStr)
	}
	return &pb.SimpleResponseCode{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}
