package rpc

import (
	"github.com/RosettaFlow/Carrier-Go/types"
)

type Backend interface {
	Start()
	Stop() error
	Status() error
	SendMsg (msg types.Msg) error
}
