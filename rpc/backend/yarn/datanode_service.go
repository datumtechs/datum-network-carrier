package yarn

import (
	"context"
	"fmt"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strings"
)

func (svr *Server) ReportUpFileSummary(ctx context.Context, req *pb.ReportUpFileSummaryRequest) (*apicommonpb.SimpleResponse, error) {

	if "" == req.GetOriginId() {
		return nil, ErrReqOriginIdForReportUpFileSummary
	}

	if "" == req.GetFilePath() {
		return nil, ErrReqFilePathForReportUpFileSummary
	}

	if "" == req.GetIp() || "" == req.GetPort() {
		return nil, ErrReqIpOrPortForReportUpFileSummary
	}

	dataNodeList, err := svr.B.GetRegisterNodeList(pb.PrefixTypeDataNode)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportUpFileSummary failed, call QueryRegisterNodeList() failed, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}",
			req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort())

		errMsg := fmt.Sprintf("%s, %s, %s, %s, %s", ErrGetRegisterNodeListForReportUpFileSummary.Msg,
			req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort())
		return nil, backend.NewRpcBizErr(ErrGetRegisterNodeListForReportUpFileSummary.Code, errMsg)
	}
	var resourceId string
	for _, dataNode := range dataNodeList {
		if req.GetIp() == dataNode.GetInternalIp() && req.GetPort() == dataNode.GetInternalPort() {
			resourceId = dataNode.GetId()
			break
		}
	}
	if "" == strings.Trim(resourceId, "") {
		log.Errorf("RPC-API:ReportUpFileSummary failed, not found resourceId, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
			req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)

		errMsg := fmt.Sprintf("%s, %s, %s, %s, %s, %s, ", ErrFoundResourceIdForReportUpFileSummary.Msg,
			req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)
		return nil, backend.NewRpcBizErr(ErrFoundResourceIdForReportUpFileSummary.Code, errMsg)
	}
	err = svr.B.StoreDataResourceFileUpload(types.NewDataResourceFileUpload(resourceId, req.GetOriginId(), "", req.GetFilePath()))
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportUpFileSummary failed, call StoreDataResourceFileUpload() failed, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
			req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)

		errMsg := fmt.Sprintf("%s, %s, %s, %s, %s, %s", ErrStoreDataResourceFileUpload.Msg,
			req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)
		return nil, backend.NewRpcBizErr(ErrStoreDataResourceFileUpload.Code, errMsg)
	}

	log.Debugf("RPC-API:ReportUpFileSummary succeed, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
		req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)

	return &apicommonpb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) ReportTaskResultFileSummary(ctx context.Context, req *pb.ReportTaskResultFileSummaryRequest) (*apicommonpb.SimpleResponse, error) {

	if "" == req.GetTaskId() {
		return nil, ErrReqTaskIdForReportTaskResultFileSummary
	}

	if "" == req.GetOriginId() {
		return nil, ErrReqOriginIdForReportTaskResultFileSummary
	}

	if "" == req.GetFilePath() {
		return nil, ErrReqFilePathForReportTaskResultFileSummary
	}

	if "" == req.GetIp() || "" == req.GetPort() {
		return nil, ErrReqIpOrPortForReportTaskResultFileSummary
	}

	dataNodeList, err := svr.B.GetRegisterNodeList(pb.PrefixTypeDataNode)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportTaskResultFileSummary failed, call QueryRegisterNodeList() failed, req.TaskId: {%s}, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}",
			req.GetTaskId(), req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort())

		errMsg := fmt.Sprintf("%s, %s, %s, %s, %s, %s", ErrGetRegisterNodeListForReportTaskResultFileSummary.Msg,
			req.GetTaskId(), req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort())
		return nil, backend.NewRpcBizErr(ErrGetRegisterNodeListForReportTaskResultFileSummary.Code, errMsg)
	}
	var resourceId string
	for _, dataNode := range dataNodeList {
		if req.GetIp() == dataNode.GetInternalIp() && req.GetPort() == dataNode.GetInternalPort() {
			resourceId = dataNode.GetId()
			break
		}
	}
	if "" == strings.Trim(resourceId, "") {
		log.Errorf("RPC-API:ReportTaskResultFileSummary failed, not found resourceId, req.TaskId: {%s}, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
			req.GetTaskId(), req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)

		errMsg := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s", ErrFoundResourceIdForReportTaskResultFileSummary.Msg,
			req.GetTaskId(), req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)
		return nil, backend.NewRpcBizErr(ErrFoundResourceIdForReportTaskResultFileSummary.Code, errMsg)
	}

	err = svr.B.StoreTaskResultFileSummary(req.GetTaskId(), req.GetOriginId(), req.GetFilePath(), resourceId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportTaskResultFileSummary failed, call StoreTaskResultFileSummary() failed, req.TaskId: {%s}, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
			req.GetTaskId(), req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)

		errMsg := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s", ErrStoreTaskResultFileSummary.Msg,
			req.GetTaskId(), req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)
		return nil, backend.NewRpcBizErr(ErrStoreTaskResultFileSummary.Code, errMsg)
	}

	log.Debugf("RPC-API:ReportTaskResultFileSummary succeed, req.TaskId: {%s}, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
		req.GetTaskId(), req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)

	return &apicommonpb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}
