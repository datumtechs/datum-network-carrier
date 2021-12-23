package yarn

import (
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
)

var (
	ErrGetRegisteredPeers                                = &backend.RpcBizErr{Code: 14001, Msg: "Failed to get all registeredNodes"}
	ErrSetSeedNodeInfo                                   = &backend.RpcBizErr{Code: 14002, Msg: "Failed to set seed node"}
	ErrDeleteSeedNodeInfo                                = &backend.RpcBizErr{Code: 14003, Msg: "Failed to delete seed node"}
	ErrGetSeedNodeList                                   = &backend.RpcBizErr{Code: 14004, Msg: "Failed to get seed nodes"}
	ErrSetDataNodeInfo                                   = &backend.RpcBizErr{Code: 14005, Msg: "Failed to set data node"}
	ErrDeleteDataNodeInfo                                = &backend.RpcBizErr{Code: 14006, Msg: "Failed to delete data node"}
	ErrGetDataNodeList                                   = &backend.RpcBizErr{Code: 14007, Msg: "Failed to get data nodes"}
	ErrGetDataNodeInfo                                   = &backend.RpcBizErr{Code: 14008, Msg: "Failed to get data node"}
	ErrSetJobNodeInfo                                    = &backend.RpcBizErr{Code: 14009, Msg: "Failed to set job node"}
	ErrGetJobNodeList                                    = &backend.RpcBizErr{Code: 14010, Msg: "Failed to get job nodes"}
	ErrDeleteJobNodeInfo                                 = &backend.RpcBizErr{Code: 14011, Msg: "Failed to delete job node"}
	ErrReportTaskEvent                                   = &backend.RpcBizErr{Code: 14012, Msg: "Failed to report taskEvent"}
	ErrReportUpFileSummary                               = &backend.RpcBizErr{Code: 14013, Msg: "Failed to ReportUpFileSummary"}
	ErrQueryDataResourceTableList                        = &backend.RpcBizErr{Code: 14014, Msg: "Failed to query dataResourceTableList"}
	ErrQueryDataResourceDataUsed                         = &backend.RpcBizErr{Code: 14015, Msg: "Failed to query dataResourceDataUsed"}
	ErrQueryTaskResultFileSummary                        = &backend.RpcBizErr{Code: 14016, Msg: "Failed to query taskResultFileSummary"}
	ErrGetNodeInfo                                       = &backend.RpcBizErr{Code: 14017, Msg: "Failed to get yarn node information"}
	ErrReqOriginIdForReportUpFileSummary                 = &backend.RpcBizErr{Code: 14018, Msg: "The original file upload request is reported: the originId is empty"}
	ErrReqFilePathForReportUpFileSummary                 = &backend.RpcBizErr{Code: 14019, Msg: "The original file upload request is reported: the filepath is empty"}
	ErrReqIpOrPortForReportUpFileSummary                 = &backend.RpcBizErr{Code: 14020, Msg: "The original file upload request is reported: the ip address or port is empty"}
	ErrGetRegisterNodeListForReportUpFileSummary         = &backend.RpcBizErr{Code: 14021, Msg: "ReportUpFileSummary failed, call QueryRegisterNodeList() failed"}
	ErrFoundResourceIdForReportUpFileSummary             = &backend.RpcBizErr{Code: 14022, Msg: "ReportUpFileSummary failed, not found resourceId"}
	ErrStoreDataResourceFileUpload                       = &backend.RpcBizErr{Code: 14023, Msg: "ReportUpFileSummary failed, call StoreDataResourceFileUpload() failed"}
	ErrReqTaskIdForReportTaskResultFileSummary           = &backend.RpcBizErr{Code: 14024, Msg: "Report task result file summary request: the taskId is empty"}
	ErrReqOriginIdForReportTaskResultFileSummary         = &backend.RpcBizErr{Code: 14025, Msg: "Report task result file summary request: the originId is empty"}
	ErrReqFilePathForReportTaskResultFileSummary         = &backend.RpcBizErr{Code: 14026, Msg: "Report task result file summary request: the filepath is empty"}
	ErrReqIpOrPortForReportTaskResultFileSummary         = &backend.RpcBizErr{Code: 14027, Msg: "Report task result file summary request: the ip address or port is empty"}
	ErrGetRegisterNodeListForReportTaskResultFileSummary = &backend.RpcBizErr{Code: 14028, Msg: "ReportTaskResultFileSummary failed, call QueryRegisterNodeList() failed"}
	ErrFoundResourceIdForReportTaskResultFileSummary     = &backend.RpcBizErr{Code: 14029, Msg: "ReportTaskResultFileSummary failed, not found resourceId"}
	ErrStoreTaskResultFileSummary                        = &backend.RpcBizErr{Code: 14030, Msg: "ReportTaskResultFileSummary failed, call StoreTaskResultFileSummary() failed"}
	ErrGetDataNodeInfoForQueryAvailableDataNode          = &backend.RpcBizErr{Code: 14031, Msg: "QueryAvailableDataNode failed, call QueryRegisterNode() failed"}
	ErrGetDataNodeInfoForQueryFilePosition               = &backend.RpcBizErr{Code: 14032, Msg: "QueryFilePosition failed, call QueryRegisterNode() failed"}
	ErrGetDataNodeInfoForTaskResultFileSummary           = &backend.RpcBizErr{Code: 14033, Msg: "GetTaskResultFileSummary failed, call QueryRegisterNode() failed"}
	ErrReqOriginIdForQueryFilePosition                   = &backend.RpcBizErr{Code: 14034, Msg: "Query File Position Request: the originId is empty"}
	ErrReqTaskIdForGetTaskResultFileSummary              = &backend.RpcBizErr{Code: 14035, Msg: "Get Task Result File Summary Request: require taskId"}
	ErrReportTaskResourceExpense                         = &backend.RpcBizErr{Code: 14036, Msg: "Failed to ReportTaskResourceExpense"}
)

type Server struct {
	B            backend.Backend
	RpcSvrIp     string
	RpcSvrPort   string
}
