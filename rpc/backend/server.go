package backend

const (
	OK = "ok"
)

var (
	ErrGetNodeInfoStr        = "Failed to get yarn node information"
	ErrGetRegisteredPeersStr = "Failed to get all registeredNodes"

	ErrSetSeedNodeInfoStr    = "Failed to set seed node info"
	ErrDeleteSeedNodeInfoStr = "Failed to delete seed node info"
	ErrGetSeedNodeListStr    = "Failed to get seed nodes"
	ErrSetDataNodeInfoStr    = "Failed to set data node info"
	ErrDeleteDataNodeInfoStr = "Failed to delete data node info"
	ErrGetDataNodeListStr    = "Failed to get data nodes"
	ErrGetDataNodeInfoStr    = "Failed to get data node info"
	ErrSetJobNodeInfoStr     = "Failed to set job node info"
	ErrGetJobNodeListStr    = "Failed to get data nodes"
	ErrDeleteJobNodeInfoStr  = "Failed to delete job node info"
	ErrSendPowerMsgStr       = "Failed to send powerMsg"
	ErrSendMetaDataMsgStr    = "Failed to send metaDataMsg"
	ErrSendTaskMsgStr        = "Failed to send taskMsg"

	ErrReportTaskEventStr = "Failed to report taskEvent"

	ErrSendPowerRevokeMsgStr    = "Failed to send powerRevokeMsg"
	ErrSendMetaDataRevokeMsgStr = "Failed to send metaDataRevokeMsg"

	ErrSendIdentityMsgStr    = "Failed to send identityMsg"
	ErrSendIdentityRevokeMsgStr    = "Failed to send identityRevokeMsg"


	ErrGetMetaDataDetailStr = "Failed to get metadata detail"
	ErrGetMetaDataDetailListStr = "Failed to get metadata detail list"

	ErrGetTotalPowerListStr = "Failed to get total power list"
	ErrGetSinglePowerListStr = "Failed to get current node power list"

	ErrGetNodeIdentityStr = "Failed to get node identityInfo"
	ErrGetIdentityListStr = "Failed to get all identityInfo list"

	ErrGetNodeTaskListStr = "Failed to get all task of current node"
	ErrGetNodeTaskEventListStr = "Failed to get all event of current node's task"

	ErrFetchRpcMetadataCtxStr = "Failed to fetch rpc metadata ctx"
	ErrReportUpFileSummaryStr = "Failed to ReportUpFileSummary"
	ErrQueryDataResourceTableListStr = "Failed to query dataResourceTableList"
	ErrQueryDataResourceDataUsedStr = "Failed to query dataResourceDataUsed"
)

type YarnServiceServer struct {
	B Backend
}

type MetaDataServiceServer struct {
	B Backend
}

type PowerServiceServer struct {
	B Backend
}

type AuthServiceServer struct {
	B Backend
}

type TaskServiceServer struct {
	B Backend
}