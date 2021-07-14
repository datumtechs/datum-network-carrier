package auth

import "github.com/RosettaFlow/Carrier-Go/rpc/backend"

const (
	ErrSendIdentityMsgStr       = "Failed to send identityMsg"
	ErrSendIdentityRevokeMsgStr = "Failed to send identityRevokeMsg"
	ErrGetNodeIdentityStr       = "Failed to get node identityInfo"
	ErrGetIdentityListStr       = "Failed to get all identityInfo list"
)

type AuthServiceServer struct {
	B backend.Backend
}
