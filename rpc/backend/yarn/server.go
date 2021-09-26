package yarn

import (
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
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
	ErrReportTaskResourceExpense  = &backend.RpcBizErr{Msg: "Failed to ReportTaskResourceExpense"}
	ErrQueryDataResourceTableList = &backend.RpcBizErr{Msg: "Failed to query dataResourceTableList"}
	ErrQueryDataResourceDataUsed  = &backend.RpcBizErr{Msg: "Failed to query dataResourceDataUsed"}
	ErrQueryTaskResultFileSummary = &backend.RpcBizErr{Msg: "Failed to query taskResultFileSummary"}
	ErrGetNodeInfo                = &backend.RpcBizErr{Msg: "Failed to get yarn node information"}
)

type Server struct {
	B backend.Backend
}
