package yarn

import (
	"context"
	"errors"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strings"
)

func (svr *Server) ReportUpFileSummary(ctx context.Context, req *pb.ReportUpFileSummaryRequest) (*apicommonpb.SimpleResponse, error) {

	if "" == req.GetOriginId() {
		return nil, errors.New("require originId")
	}

	if "" == req.GetFilePath() {
		return nil, errors.New("require filePath")
	}

	if "" == req.GetIp() || "" == req.GetPort() {
		return nil, errors.New("require ip or port")
	}

	dataNodeList, err := svr.B.GetRegisterNodeList(pb.PrefixTypeDataNode)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportUpFileSummary failed, call GetRegisterNodeList() failed, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
			req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort())
		return nil, ErrGetDataNodeList
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
		return nil, ErrGetDataNodeList
	}
	err = svr.B.StoreDataResourceFileUpload(types.NewDataResourceFileUpload(resourceId, req.GetOriginId(), "", req.GetFilePath()))
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportUpFileSummary failed, call StoreDataResourceFileUpload() failed, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
			req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)
		return nil, ErrReportUpFileSummary
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
		return nil, errors.New("require taskId")
	}

	if "" == req.GetOriginId() {
		return nil, errors.New("require originId")
	}

	if "" == req.GetFilePath() {
		return nil, errors.New("require filePath")
	}

	if "" == req.GetIp() || "" == req.GetPort() {
		return nil, errors.New("require ip or port")
	}

	dataNodeList, err := svr.B.GetRegisterNodeList(pb.PrefixTypeDataNode)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportTaskResultFileSummary failed, call GetRegisterNodeList() failed, req.TaskId: {%s}, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
			req.GetTaskId(), req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort())
		return nil, ErrGetDataNodeList
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
		return nil, ErrGetDataNodeList
	}

	err = svr.B.StoreTaskResultFileSummary(req.GetTaskId(), req.GetOriginId(), req.GetFilePath(), resourceId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportTaskResultFileSummary failed, call StoreDataResourceFileUpload() failed, req.TaskId: {%s}, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
			req.GetTaskId(), req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)
		return nil, ErrReportUpFileSummary
	}

	log.Debugf("RPC-API:ReportTaskResultFileSummary succeed, req.TaskId: {%s}, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
		req.GetTaskId(), req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)

	return &apicommonpb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}