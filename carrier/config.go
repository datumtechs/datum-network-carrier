package carrier

import (
	"github.com/RosettaFlow/Carrier-Go/core"
	"github.com/RosettaFlow/Carrier-Go/p2p"
)

// DefaultConfig contains default settings for use on the Carrier main.
var DefaultConfig = Config{
	DatabaseCache: 768,
}

//go:generate gencodec -type Config -formats toml -out gen_config.go

type Config struct {
	CarrierDB core.CarrierDB
	P2P       p2p.P2P

	// Database options
	DatabaseHandles    int  `toml:"-"`
	DatabaseCache      int
}