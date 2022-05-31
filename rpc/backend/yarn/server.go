package yarn

import (
	"github.com/datumtechs/datum-network-carrier/rpc/backend"
)

type Server struct {
	B            backend.Backend
	RpcSvrIp     string
	RpcSvrPort   string
}
