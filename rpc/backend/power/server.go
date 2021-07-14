package power

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

const (
	ErrSendPowerRevokeMsgStr = "Failed to send powerRevokeMsg"
	ErrGetTotalPowerListStr  = "Failed to get total power list"
	ErrGetSinglePowerListStr = "Failed to get current node power list"
	ErrSendPowerMsgStr       = "Failed to send powerMsg"
)

type PowerServiceServer struct {
	B backend.Backend
}
