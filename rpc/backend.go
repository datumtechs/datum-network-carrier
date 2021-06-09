package rpc

import "github.com/RosettaFlow/Carrier-Go/p2p"

type Backend interface {
	Engine()

	Protocols() []p2p.Protocol

	Start() error

	Stop() error
}
