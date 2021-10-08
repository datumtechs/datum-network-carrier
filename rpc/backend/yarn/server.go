package yarn

import (
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
)

var (

	ErrGetRegisteredPeers                                = &backend.RpcBizErr{Code: 14001, Msg: "Failed to get all registeredNodes"}
	ErrSetSeedNodeInfo                                   = &backend.RpcBizErr{Code: 14002, Msg: "%s failed, seedNodeId: {%s}, internalIp: {%s}, internalPort: {%s}"}
	ErrDeleteSeedNodeInfo                                = &backend.RpcBizErr{Code: 14003, Msg: "Failed to delete seed node info, seedNodeId: {%s}"}
	ErrGetSeedNodeList                                   = &backend.RpcBizErr{Code: 14004, Msg: "Failed to get seed nodes"}
	ErrSetDataNodeInfo                                   = &backend.RpcBizErr{Code: 14005, Msg: "%s failed, dataNodeId:{%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}"}
	ErrDeleteDataNodeInfo                                = &backend.RpcBizErr{Code: 14006, Msg: "Failed to delete data node info, dataNodeId: {%s}"}
	ErrGetDataNodeList                                   = &backend.RpcBizErr{Code: 14007, Msg: "Failed to get data nodes"}
	ErrGetDataNodeInfo                                   = &backend.RpcBizErr{Code: 14008, Msg: "Failed to get data node info"}
	ErrSetJobNodeInfo                                    = &backend.RpcBizErr{Code: 14009, Msg: "%s failed, jobNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}"}
	ErrGetJobNodeList                                    = &backend.RpcBizErr{Code: 14010, Msg: "Failed to get job nodes"}
	ErrDeleteJobNodeInfo                                 = &backend.RpcBizErr{Code: 14011, Msg: "Failed to delete job node info, jobNodeId: {%s}"}
	ErrReportTaskEvent                                   = &backend.RpcBizErr{Code: 14012, Msg: "Failed to report taskEvent, partyId: {%s}"}
	ErrReportUpFileSummary                               = &backend.RpcBizErr{Code: 14013, Msg: "Failed to ReportUpFileSummary"}
	ErrQueryDataResourceTableList                        = &backend.RpcBizErr{Code: 14014, Msg: "Failed to query dataResourceTableList, fileType: {%s}, fileSize: {%d}"}
	ErrQueryDataResourceDataUsed                         = &backend.RpcBizErr{Code: 14015, Msg: "Failed to query dataResourceDataUsed, originId: {%s}"}
	ErrQueryTaskResultFileSummary                        = &backend.RpcBizErr{Code: 14016, Msg: "Failed to query taskResultFileSummary, taskId: {%s}"}
	ErrGetNodeInfo                                       = &backend.RpcBizErr{Code: 14017, Msg: "Failed to get yarn node information"}
	ErrReqOriginIdForReportUpFileSummary                 = &backend.RpcBizErr{Code: 14018, Msg: "The original file upload request is reported:Failed to get the Id of the original file that was successfully uploaded"}
	ErrReqFilePathForReportUpFileSummary                 = &backend.RpcBizErr{Code: 14019, Msg: "The original file upload request is reported:Failed to get the path of the original file that was successfully uploaded"}
	ErrReqIpOrPortForReportUpFileSummary                 = &backend.RpcBizErr{Code: 14020, Msg: "The original file upload request is reported:Failed to obtain the IP address or port of the original file that was successfully uploaded"}
	ErrGetRegisterNodeListForReportUpFileSummary         = &backend.RpcBizErr{Code: 14021, Msg: "ReportUpFileSummary failed, call QueryRegisterNodeList() failed, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}"}
	ErrFoundResourceIdForReportUpFileSummary             = &backend.RpcBizErr{Code: 14022, Msg: "ReportUpFileSummary failed, not found resourceId, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}"}
	ErrStoreDataResourceFileUpload                       = &backend.RpcBizErr{Code: 14023, Msg: "ReportUpFileSummary failed, call StoreDataResourceFileUpload() failed, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}"}
	ErrReqTaskIdForReportTaskResultFileSummary           = &backend.RpcBizErr{Code: 14024, Msg: "Report task result file summary request:Failed to get task id"}
	ErrReqOriginIdForReportTaskResultFileSummary         = &backend.RpcBizErr{Code: 14025, Msg: "Report task result file summary request:Failed to get Origin id"}
	ErrReqFilePathForReportTaskResultFileSummary         = &backend.RpcBizErr{Code: 14026, Msg: "Report task result file summary request:Failed to get file path"}
	ErrReqIpOrPortForReportTaskResultFileSummary         = &backend.RpcBizErr{Code: 14027, Msg: "Report task result file summary request:Failed to get ip or port"}
	ErrGetRegisterNodeListForReportTaskResultFileSummary = &backend.RpcBizErr{Code: 14028, Msg: "ReportTaskResultFileSummary failed, call QueryRegisterNodeList() failed, req.TaskId: {%s}, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}"}
	ErrFoundResourceIdForReportTaskResultFileSummary     = &backend.RpcBizErr{Code: 14029, Msg: "ReportTaskResultFileSummary failed, not found resourceId, req.TaskId: {%s}, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}"}
	ErrStoreTaskResultFileSummary                        = &backend.RpcBizErr{Code: 14030, Msg: "ReportTaskResultFileSummary failed, call StoreTaskResultFileSummary() failed, req.TaskId: {%s}, req.GetOriginId: {%s}, req.GetFilePath: {%s}, req.Ip: {%s}, req.Port: {%s}, found dataNodeId: {%s}"}
	ErrGetDataNodeInfoForQueryAvailableDataNode          = &backend.RpcBizErr{Code: 14031, Msg: "QueryAvailableDataNode failed, call QueryRegisterNode() failed, fileType: {%s}, fileSize: {%d}, dataNodeId: {%s}"}
	ErrGetDataNodeInfoForQueryFilePosition               = &backend.RpcBizErr{Code: 14032, Msg: "QueryFilePosition failed, call QueryRegisterNode() failed, originId: {%s}, dataNodeId: {%s}"}
	ErrGetDataNodeInfoForTaskResultFileSummary           = &backend.RpcBizErr{Code: 14033, Msg: "GetTaskResultFileSummary failed, call QueryRegisterNode() failed, taskId: {%s}, dataNodeId: {%s}"}
	ErrReqOriginIdForQueryFilePosition                   = &backend.RpcBizErr{Code: 14034, Msg: "Query File Position Request: Failed to get OriginId"}
	ErrReqTaskIdForGetTaskResultFileSummary              = &backend.RpcBizErr{Code: 14035, Msg: "Get Task Result File Summary Request: require taskId"}
	ErrReportTaskResourceExpense  						 = &backend.RpcBizErr{Code: 14036, Msg: "Failed to ReportTaskResourceExpense"}

)

type Server struct {
	B backend.Backend
}
