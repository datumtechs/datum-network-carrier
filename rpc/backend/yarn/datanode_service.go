package yarn

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/rpc/backend"
	"github.com/datumtechs/datum-network-carrier/types"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (svr *Server) ReportUpFileSummary(ctx context.Context, req *pb.ReportUpFileSummaryRequest) (*carriertypespb.SimpleResponse, error) {

	if "" == req.GetOriginId() {
		return &carriertypespb.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require originId"}, nil
	}


	if "" == req.GetIp() || "" == req.GetPort() {
		return &carriertypespb.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require ip and port"}, nil
	}

	// add by v 0.4.0
	if "" == req.GetDataHash() {
		return &carriertypespb.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require dataHash"}, nil
	}
	// add by v 0.4.0
	if carriertypespb.OrigindataType_OrigindataType_Unknown == req.GetDataType() {
		return &carriertypespb.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require dataType"}, nil
	}
	// add by v 0.4.0
	if "" == req.GetMetadataOption() {
		return &carriertypespb.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataOption"}, nil
	}


	dataNodeList, err := svr.B.GetRegisterNodeList(pb.PrefixTypeDataNode)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportUpFileSummary failed, call QueryRegisterNodeList() failed, req.GetOriginId: {%s}, req.GetDataType: {%s}, req.Ip: {%s}, req.Port: {%s}, req.GetDataHash: {%s}, req.GetMetadataOption: %s",
			req.GetOriginId(), req.GetDataType().String(), req.GetIp(), req.GetPort(), req.GetDataHash(), req.GetMetadataOption())

		errMsg := fmt.Sprintf("%s,req.GetOriginId: {%s}, req.GetDataType: {%s}, req.Ip: {%s}, req.Port: {%s}, req.GetDataHash: {%s}, req.GetMetadataOption: %s", backend.ErrReportUpFileSummary.Error(),
			req.GetOriginId(), req.GetDataType().String(), req.GetIp(), req.GetPort(), req.GetDataHash(), req.GetMetadataOption())
		return &carriertypespb.SimpleResponse{ Status: backend.ErrReportUpFileSummary.ErrCode(), Msg: errMsg}, nil
	}
	var nodeId string
	for _, dataNode := range dataNodeList {
		if req.GetIp() == dataNode.GetInternalIp() && req.GetPort() == dataNode.GetInternalPort() {
			nodeId = dataNode.GetId()
			break
		}
	}
	if "" == strings.Trim(nodeId, "") {
		log.Errorf("RPC-API:ReportUpFileSummary failed, not found nodeId, req.GetOriginId: {%s}, req.GetDataType: {%s}, req.Ip: {%s}, req.Port: {%s}, req.GetDataHash: {%s}, req.GetMetadataOption: %s, found dataNodeId: {%s}",
			req.GetOriginId(), req.GetDataType().String(), req.GetIp(), req.GetPort(), req.GetDataHash(), req.GetMetadataOption(), nodeId)

		errMsg := fmt.Sprintf("%s, not found nodeId, req.GetOriginId: {%s}, req.GetDataType: {%s}, req.Ip: {%s}, req.Port: {%s}, req.GetDataHash: {%s}, req.GetMetadataOption: %s, dataNodeId: %s", backend.ErrReportUpFileSummary.Error(),
			req.GetOriginId(), req.GetDataType().String(), req.GetIp(), req.GetPort(), req.GetDataHash(), req.GetMetadataOption(), nodeId)
		return &carriertypespb.SimpleResponse{ Status: backend.ErrReportUpFileSummary.ErrCode(), Msg: errMsg}, nil
	}
	// store data upload summary when the file upload first.
	err = svr.B.StoreDataResourceFileUpload(types.NewDataResourceFileUpload(uint32(req.GetDataType()), nodeId, req.GetOriginId(), "", req.GetMetadataOption(), req.GetDataHash()))
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportUpFileSummary failed, call StoreDataResourceFileUpload() failed, req.GetOriginId: {%s}, req.GetDataType: {%s}, req.Ip: {%s}, req.Port: {%s}, req.GetDataHash: {%s}, req.GetMetadataOption: %s, found dataNodeId: {%s}",
			req.GetOriginId(), req.GetDataType().String(), req.GetIp(), req.GetPort(), req.GetDataHash(), req.GetMetadataOption(), nodeId)

		errMsg := fmt.Sprintf("%s, call StoreDataResourceFileUpload() failed, req.GetOriginId: {%s}, req.GetDataType: {%s}, req.Ip: {%s}, req.Port: {%s}, req.GetDataHash: {%s}, req.GetMetadataOption: %s, dataNodeId: %s", backend.ErrReportUpFileSummary.Error(),
			req.GetOriginId(), req.GetDataType().String(), req.GetIp(), req.GetPort(), req.GetDataHash(), req.GetMetadataOption(), nodeId)
		return &carriertypespb.SimpleResponse{ Status: backend.ErrReportUpFileSummary.ErrCode(), Msg: errMsg}, nil
	}

	log.Debugf("RPC-API:ReportUpFileSummary succeed, req.GetOriginId: {%s}, req.GetDataType: {%s}, req.Ip: {%s}, req.Port: {%s}, req.GetDataHash: {%s}, req.GetMetadataOption: %s, found dataNodeId: {%s}",
		req.GetOriginId(), req.GetDataType().String(), req.GetIp(), req.GetPort(), req.GetDataHash(), req.GetMetadataOption(), nodeId)

	return &carriertypespb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) ReportTaskResultFileSummary(ctx context.Context, req *pb.ReportTaskResultFileSummaryRequest) (*carriertypespb.SimpleResponse, error) {

	if "" == strings.Trim(req.GetTaskId(), "") {
		return &carriertypespb.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require taskId"}, nil
	}

	if "" == strings.Trim(req.GetOriginId(), "") {
		return &carriertypespb.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require originId"}, nil
	}

	if "" == strings.Trim(req.GetIp(), "") || "" == strings.Trim(req.GetPort(), "") {
		return &carriertypespb.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require ip and port"}, nil
	}
	// add by v 0.4.0
	//if "" == strings.Trim(req.GetExtra(), "") {
	//	return &carriertypespb.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require extra"}, nil
	//}
	// add by v 0.4.0
	//if "" == strings.Trim(req.GetDataHash(), "") {
	//	return &carriertypespb.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require dataHash"}, nil
	//}
	// add by v 0.4.0
	if carriertypespb.OrigindataType_OrigindataType_Unknown == req.GetDataType() {
		return &carriertypespb.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require dataType"}, nil
	}
	// add by v 0.4.0
	if "" == req.GetMetadataOption() {
		return &carriertypespb.SimpleResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataOption"}, nil
	}

	dataNodeList, err := svr.B.GetRegisterNodeList(pb.PrefixTypeDataNode)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportTaskResultFileSummary failed, call QueryRegisterNodeList() failed, req.TaskId: {%s}, req.GetOriginId: {%s}, req.GetDataType: {%s}, req.Ip: {%s}, req.Port: {%s}, req.GetMetadataOption: %s",
			req.GetTaskId(), req.GetOriginId(), req.GetDataType().String(), req.GetIp(), req.GetPort(), req.GetMetadataOption())

		errMsg := fmt.Sprintf("%s, call QueryRegisterNodeList() failed, originId: %s, dataType: %s, ip: %s, port: %s, metadataOption: %s", backend.ErrReportTaskResultFileSummary.Error(),
			req.GetOriginId(), req.GetDataType().String(), req.GetIp(), req.GetPort(), req.GetMetadataOption())
		return &carriertypespb.SimpleResponse{ Status: backend.ErrReportTaskResultFileSummary.ErrCode(), Msg: errMsg}, nil
	}
	var dataNodeId string
	for _, dataNode := range dataNodeList {
		if req.GetIp() == dataNode.GetInternalIp() && req.GetPort() == dataNode.GetInternalPort() {
			dataNodeId = dataNode.GetId()
			break
		}
	}
	if "" == strings.Trim(dataNodeId, "") {
		log.Errorf("RPC-API:ReportTaskResultFileSummary failed, not found dataNodeId, req.TaskId: {%s}, req.OriginId: {%s}, req.DataType: {%s}, req.Ip: {%s}, req.Port: {%s}, req.MetadataOption: %s, found dataNodeId: {%s}",
			req.GetTaskId(), req.GetOriginId(), req.GetDataType().String(), req.GetIp(), req.GetPort(), req.GetMetadataOption(), dataNodeId)

		errMsg := fmt.Sprintf("%s, not found dataNodeId, originId: %s, dataType: %s, ip: %s, port: %s, metadataOption: %s, dataNodeId: %s", backend.ErrReportTaskResultFileSummary.Error(),
			req.GetOriginId(), req.GetDataType().String(), req.GetIp(), req.GetPort(), req.GetMetadataOption(), dataNodeId)
		return &carriertypespb.SimpleResponse{ Status: backend.ErrReportTaskResultFileSummary.ErrCode(), Msg: errMsg}, nil
	}

	// the empty fileHash for task result file
	err = svr.B.StoreTaskResultFileSummary(req.GetTaskId(), req.GetOriginId(), req.GetDataHash(), req.GetMetadataOption(), dataNodeId, req.GetExtra(), req.GetDataType())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ReportTaskResultFileSummary failed, call StoreTaskResultFileSummary() failed, req.TaskId: {%s}, req.GetOriginId: {%s}, req.GetDataType: {%s}, req.Ip: {%s}, req.Port: {%s}, req.DataHash: {%s}, req.MetadataOption: %s, found dataNodeId: {%s}",
			req.GetTaskId(), req.GetOriginId(), req.GetDataType().String(), req.GetIp(), req.GetPort(), req.GetDataHash(), req.GetMetadataOption(), dataNodeId)

		errMsg := fmt.Sprintf("%s, call StoreTaskResultFileSummary() failed, originId: %s, dataType: %s, ip: %s, port: %s, dataHash: %s, metadataOption: %s, dataNodeId: %s", backend.ErrReportTaskResultFileSummary.Error(),
			req.GetOriginId(), req.GetDataType().String(), req.GetIp(), req.GetPort(), req.GetDataHash(), req.GetMetadataOption(), dataNodeId)
		return &carriertypespb.SimpleResponse{ Status: backend.ErrReportTaskResultFileSummary.ErrCode(), Msg: errMsg}, nil
	}

	log.Debugf("RPC-API:ReportTaskResultFileSummary succeed, req.TaskId: {%s}, req.GetOriginId: {%s}, req.GetDataType: {%s}, req.Ip: {%s}, req.Port: {%s}, req.DataHash: {%s}, req.MetadataOption: %s, found dataNodeId: {%s}",
		req.GetTaskId(), req.GetOriginId(), req.GetDataType().String(), req.GetIp(), req.GetPort(), req.GetDataHash(), req.GetMetadataOption(), dataNodeId)

	return &carriertypespb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) QueryAvailableDataNode(ctx context.Context, req *pb.QueryAvailableDataNodeRequest) (*pb.QueryAvailableDataNodeResponse, error) {

	if req.GetDataType() == carriertypespb.OrigindataType_OrigindataType_Unknown {
		return &pb.QueryAvailableDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown dataType"}, nil
	}

	if req.GetDataSize() == 0 {
		return &pb.QueryAvailableDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require dataSize"}, nil
	}

	dataResourceTables, err := svr.B.QueryDataResourceTables()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryAvailableDataNode-QueryDataResourceTables failed, dataType: {%s}, dataSize: {%d}",
			req.GetDataType(), req.GetDataSize())

		errMsg := fmt.Sprintf("%s, call QueryDataResourceTables() failed, %s, %d", backend.ErrQueryAvailableDataNode.Error(), req.GetDataType(), req.GetDataSize())
		return &pb.QueryAvailableDataNodeResponse{Status: backend.ErrQueryAvailableDataNode.ErrCode(), Msg: errMsg}, nil
	}

	var nodeId string
	for _, resource := range dataResourceTables {
		log.Debugf("QueryAvailableDataNode, dataResourceTable: %s, need disk: %d, remain disk: %d", resource.String(), req.GetDataSize(), resource.RemainDisk())
		if req.GetDataSize() < resource.RemainDisk() {
			nodeId = resource.GetNodeId()
			break
		}
	}

	if "" == strings.Trim(nodeId, "") {
		log.Errorf("RPC-API:QueryAvailableDataNode, not found available dataNodeId of dataNode, dataType: {%s}, dataSize: {%d}, dataNodeId: {%s}",
			req.GetDataType(), req.GetDataSize(), nodeId)

		errMsg := fmt.Sprintf("%s, not found available dataNodeId of dataNode, %s, %d, %s", backend.ErrQueryAvailableDataNode.Error(),
			req.GetDataType(), req.GetDataSize(), nodeId)
		return &pb.QueryAvailableDataNodeResponse{Status: backend.ErrQueryAvailableDataNode.ErrCode(), Msg: errMsg}, nil
	}

	dataNode, err := svr.B.GetRegisterNode(pb.PrefixTypeDataNode, nodeId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryAvailableDataNode-QueryRegisterNode failed, dataType: {%s}, dataSize: {%d}, dataNodeId: {%s}",
			req.GetDataType(), req.GetDataSize(), nodeId)

		errMsg := fmt.Sprintf("%s, call QueryRegisterNode() failed, %s, %d, %s", backend.ErrQueryAvailableDataNode.Error(),
			req.GetDataType(), req.GetDataSize(), nodeId)
		return &pb.QueryAvailableDataNodeResponse{Status: backend.ErrQueryAvailableDataNode.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:QueryAvailableDataNode succeed, dataType: {%s}, dataSize: {%d}, return dataNodeId: {%s}, dataNodeIp: {%s}, dataNodePort: {%s}",
		req.GetDataType(), req.GetDataSize(), dataNode.GetId(), dataNode.GetInternalIp(), dataNode.GetInternalPort())

	return &pb.QueryAvailableDataNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Information: &pb.QueryAvailableDataNode{
			Ip:   dataNode.GetInternalIp(),
			Port: dataNode.GetInternalPort(),
		},
	}, nil
}

func (svr *Server) QueryFilePosition(ctx context.Context, req *pb.QueryFilePositionRequest) (*pb.QueryFilePositionResponse, error) {

	if "" == req.GetOriginId() {
		return &pb.QueryFilePositionResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require originId"}, nil
	}

	dataResourceFileUpload, err := svr.B.QueryDataResourceFileUpload(req.GetOriginId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryFilePosition-QueryDataResourceFileUpload failed, originId: {%s}", req.GetOriginId())

		errMsg := fmt.Sprintf("%s, call QueryDataResourceFileUpload() failed, %s", backend.ErrQueryFilePosition.Error(), req.GetOriginId())
		return &pb.QueryFilePositionResponse{Status: backend.ErrQueryFilePosition.ErrCode(), Msg: errMsg}, nil
	}
	dataNode, err := svr.B.GetRegisterNode(pb.PrefixTypeDataNode, dataResourceFileUpload.GetNodeId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryFilePosition-QueryRegisterNode failed, originId: {%s}, dataNodeId: {%s}",
			req.GetOriginId(), dataResourceFileUpload.GetNodeId())

		errMsg := fmt.Sprintf("%s, call QueryRegisterNode() failed, %s, %s", backend.ErrQueryFilePosition.Error(),
			req.GetOriginId(), dataResourceFileUpload.GetNodeId())
		return &pb.QueryFilePositionResponse{Status: backend.ErrQueryFilePosition.ErrCode(), Msg: errMsg}, nil
	}

	log.Debugf("RPC-API:QueryFilePosition Succeed, originId: {%s}, return dataNodeIp: {%s}, dataNodePort: {%s}, dataType: {%s}, metadataOption: %s",
		req.GetOriginId(), dataNode.GetInternalIp(), dataNode.GetInternalPort(), carriertypespb.OrigindataType(dataResourceFileUpload.GetDataType()).String(), dataResourceFileUpload.GetMetadataOption())

	var dataPath string

	switch carriertypespb.OrigindataType(dataResourceFileUpload.GetDataType()) {
	case carriertypespb.OrigindataType_OrigindataType_CSV:
		var option *types.MetadataOptionCSV
		if err := json.Unmarshal([]byte(dataResourceFileUpload.GetMetadataOption()), &option); nil != err {
			return &pb.QueryFilePositionResponse{Status: backend.ErrQueryFilePosition.ErrCode(), Msg: fmt.Sprintf("unmashal metadataOption to csv failed, %s", err)}, nil
		}
		dataPath = option.GetDataPath()
	case carriertypespb.OrigindataType_OrigindataType_DIR:
		var option *types.MetadataOptionDIR
		if err := json.Unmarshal([]byte(dataResourceFileUpload.GetMetadataOption()), &option); nil != err {
			return &pb.QueryFilePositionResponse{Status: backend.ErrQueryFilePosition.ErrCode(), Msg: fmt.Sprintf("unmashal metadataOption to dir failed, %s", err)}, nil
		}
		dataPath = option.GetDirPath()
	case carriertypespb.OrigindataType_OrigindataType_BINARY:
		var option *types.MetadataOptionBINARY
		if err := json.Unmarshal([]byte(dataResourceFileUpload.GetMetadataOption()), &option); nil != err {
			return &pb.QueryFilePositionResponse{Status: backend.ErrQueryFilePosition.ErrCode(), Msg: fmt.Sprintf("unmashal metadataOption to binary failed, %s", err)}, nil
		}
		dataPath = option.GetDataPath()
	default:
		return &pb.QueryFilePositionResponse{Status: backend.ErrQueryFilePosition.ErrCode(), Msg: fmt.Sprint("cannot match metadata option")}, nil
	}

	return &pb.QueryFilePositionResponse{
		Status: 0,
		Msg:    backend.OK,
		Information: &pb.QueryFilePosition{
			Ip:       dataNode.GetInternalIp(),
			Port:     dataNode.GetInternalPort(),
			DataPath: dataPath,
		},
	}, nil
}

func (svr *Server) GetTaskResultFileSummary(ctx context.Context, req *pb.GetTaskResultFileSummaryRequest) (*pb.GetTaskResultFileSummaryResponse, error) {

	if "" == req.GetTaskId() {
		return &pb.GetTaskResultFileSummaryResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require taskId"}, nil
	}

	summary, err := svr.B.QueryTaskResultFileSummary(req.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetTaskResultFileSummary-QueryTaskResultFileSummary failed, taskId: {%s}", req.GetTaskId())

		errMsg := fmt.Sprintf("%s, call QueryTaskResultFileSummary() failed, %s", backend.ErrQueryTaskResultFileSummary.Error(), req.GetTaskId())
		return &pb.GetTaskResultFileSummaryResponse{Status: backend.ErrQueryTaskResultFileSummary.ErrCode(), Msg: errMsg}, nil
	}
	dataNode, err := svr.B.GetRegisterNode(pb.PrefixTypeDataNode, summary.GetNodeId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetTaskResultFileSummary-QueryRegisterNode failed, taskId: {%s}, dataNodeId: {%s}",
			req.GetTaskId(), summary.GetNodeId())

		errMsg := fmt.Sprintf("%s, call QueryRegisterNode() failed, %s, %s", backend.ErrQueryTaskResultFileSummary.Error(),
			req.GetTaskId(), summary.GetNodeId())
		return &pb.GetTaskResultFileSummaryResponse{Status: backend.ErrQueryTaskResultFileSummary.ErrCode(), Msg: errMsg}, nil
	}

	log.Debugf("RPC-API:GetTaskResultFileSummary Succeed, taskId: {%s}, return dataNodeIp: {%s}, dataNodePort: {%s}, metadataId: {%s}, originId: {%s}, metadataName: {%s}, dataHash: {%s}, dataType: {%s}, metadataOption: %s",
		req.GetTaskId(), dataNode.GetInternalIp(), dataNode.GetInternalPort(), summary.GetMetadataId(), summary.GetOriginId(), summary.GetMetadataName(), summary.GetDataHash(), carriertypespb.OrigindataType(summary.GetDataType()).String(), summary.GetMetadataOption())

	return &pb.GetTaskResultFileSummaryResponse{
		Status: 0,
		Msg:    backend.OK,
		Information: &pb.GetTaskResultFileSummary{
			/**
			TaskId               string
			MetadataName         string
			MetadataId           string
			OriginId             string
			Ip                   string
			Port                 string
			Extra                string
			DataHash             string
			DataType             types.OrigindataType
			MetadataOption       string
			*/
			TaskId:         summary.GetTaskId(),
			MetadataName:   summary.GetMetadataName(),
			MetadataId:     summary.GetMetadataId(),
			OriginId:       summary.GetOriginId(),
			Ip:             dataNode.GetInternalIp(),
			Port:           dataNode.GetInternalPort(),
			Extra:          summary.GetExtra(),
			DataHash:       summary.GetDataHash(),
			DataType:       carriertypespb.OrigindataType(summary.GetDataType()),
			MetadataOption: summary.GetMetadataOption(),
		},
	}, nil
}

func (svr *Server) GetTaskResultFileSummaryList(ctx context.Context, empty *emptypb.Empty) (*pb.GetTaskResultFileSummaryListResponse, error) {
	taskResultFileSummaryArr, err := svr.B.QueryTaskResultFileSummaryList()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetTaskResultFileSummaryList-QueryTaskResultFileSummaryList failed")
		return &pb.GetTaskResultFileSummaryListResponse{Status: backend.ErrQueryTaskResultFileSummaryList.ErrCode(), Msg: backend.ErrQueryTaskResultFileSummaryList.Error()}, nil
	}

	arr := make([]*pb.GetTaskResultFileSummary, 0)
	for _, summary := range taskResultFileSummaryArr {
		dataNode, err := svr.B.GetRegisterNode(pb.PrefixTypeDataNode, summary.GetNodeId())
		if nil != err {
			log.WithError(err).Errorf("RPC-API:GetTaskResultFileSummaryList-QueryRegisterNode failed, taskId: {%s}, dataNodeId: {%s}",
				summary.GetTaskId(), summary.GetNodeId())
			continue
		}
		arr = append(arr, &pb.GetTaskResultFileSummary{
			/**
			TaskId               string
			MetadataName         string
			MetadataId           string
			OriginId             string
			Ip                   string
			Port                 string
			Extra                string
			DataHash             string
			DataType             types.OrigindataType
			MetadataOption       string
			*/
			TaskId:         summary.GetTaskId(),
			MetadataName:   summary.GetMetadataName(),
			MetadataId:     summary.GetMetadataId(),
			OriginId:       summary.GetOriginId(),
			Ip:             dataNode.GetInternalIp(),
			Port:           dataNode.GetInternalPort(),
			Extra:          summary.GetExtra(),
			DataHash:       summary.GetDataHash(),
			DataType:       carriertypespb.OrigindataType(summary.GetDataType()),
			MetadataOption: summary.GetMetadataOption(),
		})
	}

	log.Debugf("RPC-API:GetTaskResultFileSummaryList Succeed, task result file summary list len: {%d}", len(arr))
	return &pb.GetTaskResultFileSummaryListResponse{
		Status:          0,
		Msg:             backend.OK,
		TaskResultFiles: arr,
	}, nil
}
