package backend

type RpcBizErr struct {
	Code int32
	Msg  string
}

func NewRpcBizErr(code int32, msg string) *RpcBizErr { return &RpcBizErr{Code: code, Msg: msg} }

func (e *RpcBizErr) Error() string {
	return e.Msg
}

func (e *RpcBizErr) ErrCode() int32 {
	return e.Code
}

// common
var (
	ErrSystemPanic   = &RpcBizErr{Code: 00001, Msg: "System error"}
	ErrRequireParams = &RpcBizErr{Code: 00002, Msg: "Require params"}
)

// auth
var (
	ErrApplyIdentityMsg        = &RpcBizErr{Code: 10001, Msg: "ApplyIdentityJoin failed"}
	ErrRevokeIdentityMsg       = &RpcBizErr{Code: 10002, Msg: "RevokeIdentityJoin failed"}
	ErrGetNodeIdentity         = &RpcBizErr{Code: 10003, Msg: "Query local identity failed"}
	ErrGetIdentityList         = &RpcBizErr{Code: 10004, Msg: "Query all identity list failed"}
	ErrGetAuthorityList        = &RpcBizErr{Code: 10005, Msg: "Query all authority list failed"}
	ErrApplyMetadataAuthority  = &RpcBizErr{Code: 10006, Msg: "ApplyMetadataAuthority failed"}
	ErrRevokeMetadataAuthority = &RpcBizErr{Code: 10007, Msg: "RevokeMetadataAuthority failed"}
	ErrAuditMetadataAuth       = &RpcBizErr{Code: 10008, Msg: "Audit metadataAuth failed"}
)

// metadata
var (
	ErrGetMetadataDetail           = &RpcBizErr{Code: 11001, Msg: "Query metadata detail failed"}
	ErrGetMetadataDetailList       = &RpcBizErr{Code: 11002, Msg: "Query metadata detail list failed"}
	ErrPublishMetadataMsg          = &RpcBizErr{Code: 11003, Msg: "PublishMetaData failed"}
	ErrRevokeMetadataMsg           = &RpcBizErr{Code: 11004, Msg: "RevokemetaData failed"}
	ErrQueryMetadataUsedTaskIdList = &RpcBizErr{Code: 11005, Msg: "Query metadata used taskId list failed"}
)

// power
var (
	ErrRevokePowerMsg       = &RpcBizErr{Code: 12001, Msg: "RevokePower failed"}
	ErrQueryGlobalPowerList = &RpcBizErr{Code: 12002, Msg: "Query global power list failed"}
	ErrQueryLocalPowerList  = &RpcBizErr{Code: 12003, Msg: "Query local power list failed"}
	ErrPublishPowerMsg      = &RpcBizErr{Code: 12004, Msg: "PublishPower failed"}
)

// task
var (
	ErrGetNodeTaskList      = &RpcBizErr{Code: 13001, Msg: "Failed to get all task of current node"}
	ErrGetNodeTaskEventList = &RpcBizErr{Code: 13002, Msg: "Failed to get all event of current node's task"}
	ErrPublishTaskMsg       = &RpcBizErr{Code: 13003, Msg: "PublishTask failed"}
	ErrTerminateTaskMsg     = &RpcBizErr{Code: 13004, Msg: "TerminateTask failed, send task terminate msg failed"}
)

// yarn
var (
	ErrGetRegisteredPeers                                = &RpcBizErr{Code: 14001, Msg: "Failed to get all registeredNodes"}
	ErrSetSeedNodeInfo                                   = &RpcBizErr{Code: 14002, Msg: "Failed to set seed node"}
	ErrDeleteSeedNodeInfo                                = &RpcBizErr{Code: 14003, Msg: "Failed to delete seed node"}
	ErrGetSeedNodeList                                   = &RpcBizErr{Code: 14004, Msg: "Failed to get seed nodes"}
	ErrSetDataNodeInfo                                   = &RpcBizErr{Code: 14005, Msg: "Failed to set data node"}
	ErrDeleteDataNodeInfo                                = &RpcBizErr{Code: 14006, Msg: "Failed to delete data node"}
	ErrGetDataNodeList                                   = &RpcBizErr{Code: 14007, Msg: "Failed to get data nodes"}
	ErrGetDataNodeInfo                                   = &RpcBizErr{Code: 14008, Msg: "Failed to get data node"}
	ErrSetJobNodeInfo                                    = &RpcBizErr{Code: 14009, Msg: "Failed to set job node"}
	ErrGetJobNodeList                                    = &RpcBizErr{Code: 14010, Msg: "Failed to get job nodes"}
	ErrDeleteJobNodeInfo                                 = &RpcBizErr{Code: 14011, Msg: "Failed to delete job node"}
	ErrReportTaskEvent                                   = &RpcBizErr{Code: 14012, Msg: "Failed to report taskEvent"}
	ErrReportUpFileSummary                               = &RpcBizErr{Code: 14013, Msg: "Failed to ReportUpFileSummary"}
	ErrQueryDataResourceTableList                        = &RpcBizErr{Code: 14014, Msg: "Failed to query dataResourceTableList"}
	ErrQueryDataResourceDataUsed                         = &RpcBizErr{Code: 14015, Msg: "Failed to query dataResourceDataUsed"}
	ErrQueryTaskResultFileSummary                        = &RpcBizErr{Code: 14016, Msg: "Failed to query taskResultFileSummary"}
	ErrGetNodeInfo                                       = &RpcBizErr{Code: 14017, Msg: "Failed to get yarn node information"}
	ErrReqOriginIdForReportUpFileSummary                 = &RpcBizErr{Code: 14018, Msg: "The original file upload request is reported: the originId is empty"}
	ErrReqFilePathForReportUpFileSummary                 = &RpcBizErr{Code: 14019, Msg: "The original file upload request is reported: the filepath is empty"}
	ErrReqIpOrPortForReportUpFileSummary                 = &RpcBizErr{Code: 14020, Msg: "The original file upload request is reported: the ip address or port is empty"}
	ErrGetRegisterNodeListForReportUpFileSummary         = &RpcBizErr{Code: 14021, Msg: "ReportUpFileSummary failed, call QueryRegisterNodeList() failed"}
	ErrFoundResourceIdForReportUpFileSummary             = &RpcBizErr{Code: 14022, Msg: "ReportUpFileSummary failed, not found resourceId"}
	ErrStoreDataResourceFileUpload                       = &RpcBizErr{Code: 14023, Msg: "ReportUpFileSummary failed, call StoreDataResourceFileUpload() failed"}
	ErrReportTaskResultFileSummary                       = &RpcBizErr{Code: 14024, Msg: "Report task result file summary failed"}
	ErrReqOriginIdForReportTaskResultFileSummary         = &RpcBizErr{Code: 14025, Msg: "Report task result file summary request: the originId is empty"}
	ErrReqFilePathForReportTaskResultFileSummary         = &RpcBizErr{Code: 14026, Msg: "Report task result file summary request: the filepath is empty"}
	ErrReqIpOrPortForReportTaskResultFileSummary         = &RpcBizErr{Code: 14027, Msg: "Report task result file summary request: the ip address or port is empty"}
	ErrGetRegisterNodeListForReportTaskResultFileSummary = &RpcBizErr{Code: 14028, Msg: "ReportTaskResultFileSummary failed, call QueryRegisterNodeList() failed"}
	ErrFoundResourceIdForReportTaskResultFileSummary     = &RpcBizErr{Code: 14029, Msg: "ReportTaskResultFileSummary failed, not found resourceId"}
	ErrStoreTaskResultFileSummary                        = &RpcBizErr{Code: 14030, Msg: "ReportTaskResultFileSummary failed, call StoreTaskResultFileSummary() failed"}
	ErrGetDataNodeInfoForQueryAvailableDataNode          = &RpcBizErr{Code: 14031, Msg: "QueryAvailableDataNode failed, call QueryRegisterNode() failed"}
	ErrGetDataNodeInfoForQueryFilePosition               = &RpcBizErr{Code: 14032, Msg: "QueryFilePosition failed, call QueryRegisterNode() failed"}
	ErrGetDataNodeInfoForTaskResultFileSummary           = &RpcBizErr{Code: 14033, Msg: "GetTaskResultFileSummary failed, call QueryRegisterNode() failed"}
	ErrReqOriginIdForQueryFilePosition                   = &RpcBizErr{Code: 14034, Msg: "Query File Position Request: the originId is empty"}
	ErrReqTaskIdForGetTaskResultFileSummary              = &RpcBizErr{Code: 14035, Msg: "Get Task Result File Summary Request: require taskId"}
	ErrReportTaskResourceExpense                         = &RpcBizErr{Code: 14036, Msg: "Failed to ReportTaskResourceExpense"}
)
