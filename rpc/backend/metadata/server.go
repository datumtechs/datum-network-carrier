package metadata

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

const (
	ErrSendMetaDataRevokeMsgStr = "Failed to send metaDataRevokeMsg"
	ErrGetMetaDataDetailStr     = "Failed to get metadata detail"
	ErrGetMetaDataDetailListStr = "Failed to get metadata detail list"
	ErrSendMetaDataMsgStr    = "Failed to send metaDataMsg"
)

type MetaDataServiceServer struct {
	B backend.Backend
}
