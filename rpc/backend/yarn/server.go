package yarn

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

const (
	ErrGetRegisteredPeersStr = "Failed to get all registeredNodes"
	ErrSetSeedNodeInfoStr    = "Failed to set seed node info"
	ErrDeleteSeedNodeInfoStr = "Failed to delete seed node info"
	ErrGetSeedNodeListStr    = "Failed to get seed nodes"
	ErrSetDataNodeInfoStr    = "Failed to set data node info"
	ErrDeleteDataNodeInfoStr = "Failed to delete data node info"
	ErrGetDataNodeListStr    = "Failed to get data nodes"
	ErrGetDataNodeInfoStr    = "Failed to get data node info"
	ErrSetJobNodeInfoStr     = "Failed to set job node info"
	ErrGetJobNodeListStr     = "Failed to get data nodes"
	ErrDeleteJobNodeInfoStr  = "Failed to delete job node info"

	ErrReportTaskEventStr = "Failed to report taskEvent"

	ErrReportUpFileSummaryStr        = "Failed to ReportUpFileSummary"
	ErrQueryDataResourceTableListStr = "Failed to query dataResourceTableList"
	ErrQueryDataResourceDataUsedStr  = "Failed to query dataResourceDataUsed"
	ErrGetNodeInfoStr                = "Failed to get yarn node information"
)

var (
	ErrGetRegisteredPeers         = &backend.RpcBizErr{Msg: "Failed to get all registeredNodes"}
	ErrSetSeedNodeInfo            = &backend.RpcBizErr{Msg: "Failed to set seed node info"}
	ErrDeleteSeedNodeInfo         = &backend.RpcBizErr{Msg: "Failed to delete seed node info"}
	ErrGetSeedNodeList            = &backend.RpcBizErr{Msg: "Failed to get seed nodes"}
	ErrSetDataNodeInfo            = &backend.RpcBizErr{Msg: "Failed to set data node info"}
	ErrDeleteDataNodeInfo         = &backend.RpcBizErr{Msg: "Failed to delete data node info"}
	ErrGetDataNodeList            = &backend.RpcBizErr{Msg: "Failed to get data nodes"}
	ErrGetDataNodeInfo            = &backend.RpcBizErr{Msg: "Failed to get data node info"}
	ErrSetJobNodeInfo             = &backend.RpcBizErr{Msg: "Failed to set job node info"}
	ErrGetJobNodeList             = &backend.RpcBizErr{Msg: "Failed to get data nodes"}
	ErrDeleteJobNodeInfo          = &backend.RpcBizErr{Msg: "Failed to delete job node info"}
	ErrReportTaskEvent            = &backend.RpcBizErr{Msg: "Failed to report taskEvent"}
	ErrReportUpFileSummary        = &backend.RpcBizErr{Msg: "Failed to ReportUpFileSummary"}
	ErrQueryDataResourceTableList = &backend.RpcBizErr{Msg: "Failed to query dataResourceTableList"}
	ErrQueryDataResourceDataUsed  = &backend.RpcBizErr{Msg: "Failed to query dataResourceDataUsed"}
	ErrGetNodeInfo                = &backend.RpcBizErr{Msg: "Failed to get yarn node information"}
)

type YarnServiceServer struct {
	B backend.Backend
}
