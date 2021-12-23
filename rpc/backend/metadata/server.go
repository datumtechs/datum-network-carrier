package metadata

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

var (
	ErrSendMetadataRevokeMsg              = &backend.RpcBizErr{Code: 11001, Msg: "Failed to send metaDataRevokeMsg"}
	ErrGetMetadataDetail                  = &backend.RpcBizErr{Code: 11002, Msg: "Failed to get metadata detail"}
	ErrGetMetadataDetailList              = &backend.RpcBizErr{Code: 11003, Msg: "Failed to get metadata detail list"}
	ErrSendMetadataMsg                    = &backend.RpcBizErr{Code: 11004, Msg: "Failed to send metaDataMsg"}
	ErrReqInfoForPublishMetadata          = &backend.RpcBizErr{Code: 11005, Msg: "Publish Metadata request: failed to get metadata information"}
	ErrReqMetaSummaryForPublishMetadata   = &backend.RpcBizErr{Code: 11006, Msg: "Publish Metadata request: failed to get metadata summary"}
	ErrReqMetaColumnsForPublishMetadata   = &backend.RpcBizErr{Code: 11007, Msg: "Publish Metadata request: failed to get metadata columns information"}
	ErrReqMetaIdForMetadataUsedTaskIdList = &backend.RpcBizErr{Code: 11008, Msg: "the metadataId is empty"}
	ErrReqListForMetadataUsedTaskIdList   = &backend.RpcBizErr{Code: 11009, Msg: "failed to get taskId list"}
)

type Server struct {
	B backend.Backend
}
