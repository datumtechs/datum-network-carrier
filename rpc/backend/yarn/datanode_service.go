package yarn

import (
	"context"
	"fmt"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strings"
)

func (svr *Server) ReportUpFileSummary(ctx context.Context, req *pb.ReportUpFileSummaryRequest) (*libtypes.SimpleResponse, error) {

	if "" == req.GetOriginId() {
		return &libtypes.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require originId"}, nil
	}

	if "" == req.GetFilePath() {
		return &libtypes.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require filePath"}, nil
	}

	if "" == req.GetIp() || "" == req.GetPort() {
		return &libtypes.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require ip and port"}, nil
	}

	dataNodeList, err := svr.B.GetRegisterNodeList(pb.PrefixTypeDataNode)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportUpFileSummary failed, call QueryRegisterNodeList() failed, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}",
			req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort())

		errMsg := fmt.Sprintf("%s, %s, %s, %s, %s", backend.ErrReportUpFileSummary.Error(),
			req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort())
		return &libtypes.SimpleResponse{ Status: backend.ErrReportUpFileSummary.ErrCode(), Msg: errMsg}, nil
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

		errMsg := fmt.Sprintf("%s, not found resourceId, originId: %s, filePath: %s, ip: %s, port: %s, resourceId: %s", backend.ErrReportUpFileSummary.Error(),
			req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)
		return &libtypes.SimpleResponse{ Status: backend.ErrReportUpFileSummary.ErrCode(), Msg: errMsg}, nil
	}
	err = svr.B.StoreDataResourceFileUpload(types.NewDataResourceFileUpload(resourceId, req.GetOriginId(), "", req.GetFilePath()))
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportUpFileSummary failed, call StoreDataResourceFileUpload() failed, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
			req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)

		errMsg := fmt.Sprintf("%s, call StoreDataResourceFileUpload() failed, originId: %s, filePath: %s, ip: %s, port: %s, resourceId: %s", backend.ErrReportUpFileSummary.Error(),
			req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)
		return &libtypes.SimpleResponse{ Status: backend.ErrReportUpFileSummary.ErrCode(), Msg: errMsg}, nil
	}

	log.Debugf("RPC-API:ReportUpFileSummary succeed, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
		req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)

	return &libtypes.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) ReportTaskResultFileSummary(ctx context.Context, req *pb.ReportTaskResultFileSummaryRequest) (*libtypes.SimpleResponse, error) {

	if "" == strings.Trim(req.GetTaskId(), "") {
		return &libtypes.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require taskId"}, nil
	}

	if "" == strings.Trim(req.GetOriginId(), "") {
		return &libtypes.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require originId"}, nil
	}

	if "" == strings.Trim(req.GetFilePath(), "") {
		return &libtypes.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require filePath"}, nil
	}

	if "" == strings.Trim(req.GetIp(), "") || "" == strings.Trim(req.GetPort(), "") {
		return &libtypes.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require ip and port"}, nil
	}

	dataNodeList, err := svr.B.GetRegisterNodeList(pb.PrefixTypeDataNode)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportTaskResultFileSummary failed, call QueryRegisterNodeList() failed, req.TaskId: {%s}, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}",
			req.GetTaskId(), req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort())

		errMsg := fmt.Sprintf("%s, call QueryRegisterNodeList() failed, originId: %s, filePath: %s, ip: %s, port: %s", backend.ErrReportTaskResultFileSummary.Error(),
			req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort())
		return &libtypes.SimpleResponse{ Status: backend.ErrReportTaskResultFileSummary.ErrCode(), Msg: errMsg}, nil
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

		errMsg := fmt.Sprintf("%s, not found resourceId, originId: %s, filePath: %s, ip: %s, port: %s, resourceId: %s", backend.ErrReportTaskResultFileSummary.Error(),
			req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)
		return &libtypes.SimpleResponse{ Status: backend.ErrReportTaskResultFileSummary.ErrCode(), Msg: errMsg}, nil
	}

	// the empty fileHash for task result file
	err = svr.B.StoreTaskResultFileSummary(req.GetTaskId(), req.GetOriginId(), "", req.GetFilePath(), resourceId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportTaskResultFileSummary failed, call StoreTaskResultFileSummary() failed, req.TaskId: {%s}, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
			req.GetTaskId(), req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)

		errMsg := fmt.Sprintf("%s, call StoreTaskResultFileSummary() failed, originId: %s, filePath: %s, ip: %s, port: %s, resourceId: %s", backend.ErrReportTaskResultFileSummary.Error(),
			req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)
		return &libtypes.SimpleResponse{ Status: backend.ErrReportTaskResultFileSummary.ErrCode(), Msg: errMsg}, nil
	}

	log.Debugf("RPC-API:ReportTaskResultFileSummary succeed, req.TaskId: {%s}, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}",
		req.GetTaskId(), req.GetOriginId(), req.GetFilePath(), req.GetIp(), req.GetPort(), resourceId)

	return &libtypes.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}
