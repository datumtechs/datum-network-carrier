package yarn

import (
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
)

type Server struct {
	B            backend.Backend
	RpcSvrIp     string
	RpcSvrPort   string
}
