package carrier

import (
	"github.com/RosettaFlow/Carrier-Go/core"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"time"
)

// DefaultConfig contains default settings for use on the Carrier main.
var DefaultConfig = Config{
	DatabaseCache: 768,
	DBDisabledGC:  false,
	DBGCInterval:  86400,
	DBGCTimeout:   time.Minute,
	DBGCMpt:       true,
	DBGCBlock:     10,
}

//go:generate gencodec -type Config -formats toml -out gen_config.go

type Config struct {
	CarrierDB core.CarrierDB
	P2P       p2p.P2P

	// Database options
	SkipBcVersionCheck bool `toml:"-"`
	DatabaseHandles    int  `toml:"-"`
	DatabaseCache      int
	DBDisabledGC       bool
	DBGCInterval       uint64
	DBGCTimeout        time.Duration
	DBGCMpt            bool
	DBGCBlock          int
}