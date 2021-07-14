package yarn

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

const (
	ErrGetRegisteredPeersStr         = "Failed to get all registeredNodes"
	ErrSetSeedNodeInfoStr            = "Failed to set seed node info"
	ErrDeleteSeedNodeInfoStr         = "Failed to delete seed node info"
	ErrGetSeedNodeListStr            = "Failed to get seed nodes"
	ErrSetDataNodeInfoStr            = "Failed to set data node info"
	ErrDeleteDataNodeInfoStr         = "Failed to delete data node info"
	ErrGetDataNodeListStr            = "Failed to get data nodes"
	ErrGetDataNodeInfoStr            = "Failed to get data node info"
	ErrSetJobNodeInfoStr             = "Failed to set job node info"
	ErrGetJobNodeListStr             = "Failed to get data nodes"
	ErrDeleteJobNodeInfoStr          = "Failed to delete job node info"

	ErrReportTaskEventStr            = "Failed to report taskEvent"

	ErrReportUpFileSummaryStr        = "Failed to ReportUpFileSummary"
	ErrQueryDataResourceTableListStr = "Failed to query dataResourceTableList"
	ErrQueryDataResourceDataUsedStr  = "Failed to query dataResourceDataUsed"
	ErrGetNodeInfoStr                = "Failed to get yarn node information"
)

type YarnServiceServer struct {
	B backend.Backend
}
