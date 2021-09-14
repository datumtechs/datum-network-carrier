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
		log.WithError(err).Errorf("RPC-API:ReportUpFileSummary failed, call GetRegisterNodeList() failed, req.OriginId: {%s}, req.FilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
			req.OriginId, req.FilePath, req.Ip, req.Port)
		return nil, ErrGetDataNodeList
	}
	var resourceId string
	for _, dataNode := range dataNodeList {
		if req.Ip == dataNode.InternalIp && req.Port == dataNode.InternalPort {
			resourceId = dataNode.Id
			break
		}
	}
	if "" == strings.Trim(resourceId, "") {
		log.Errorf("RPC-API:ReportUpFileSummary failed, not found resourceId, req.OriginId: {%s}, req.FilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
			req.OriginId, req.FilePath, req.Ip, req.Port, resourceId)
		return nil, ErrGetDataNodeList
	}
	err = svr.B.StoreDataResourceFileUpload(types.NewDataResourceFileUpload(resourceId, req.OriginId, "", req.FilePath))
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportUpFileSummary failed, call StoreDataResourceFileUpload() failed, req.OriginId: {%s}, req.FilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
			req.OriginId, req.FilePath, req.Ip, req.Port, resourceId)
		return nil, ErrReportUpFileSummary
	}

	log.Debugf("RPC-API:ReportUpFileSummary succeed, req.OriginId: {%s}, req.FilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
		req.OriginId, req.FilePath, req.Ip, req.Port, resourceId)

	return &pb.SimpleResponseCode{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}
