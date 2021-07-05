package carrier

import (
	"github.com/RosettaFlow/Carrier-Go/common/hexutil"
	"github.com/RosettaFlow/Carrier-Go/core"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"time"
)

// DefaultConfig contains default settings for use on the Carrier main.
var DefaultConfig = Config{
	DatabaseCache: 768,
	TrieCache:     32,
	TrieTimeout:   60 * time.Minute,
	DBDisabledGC:  false,
	DBGCInterval:  86400,
	DBGCTimeout:   time.Minute,
	DBGCMpt:       true,
	DBGCBlock:     10,

	BodyCacheLimit:    256,
	BlockCacheLimit:   256,
	MaxFutureBlocks:   256,
	BadBlockLimit:     10,
	TriesInMemory:     128,
	BlockChainVersion: 3,
}

//go:generate gencodec -type Config -field-override configMarshaling -formats toml -out gen_config.go

type Config struct {
	carrierDB core.CarrierDB
	p2p p2p.P2P
	NoPruning bool

	// Database options
	SkipBcVersionCheck bool `toml:"-"`
	DatabaseHandles    int  `toml:"-"`
	DatabaseCache      int
	TrieCache          int
	TrieTimeout        time.Duration
	TrieDBCache        int
	DBDisabledGC       bool
	DBGCInterval       uint64
	DBGCTimeout        time.Duration
	DBGCMpt            bool
	DBGCBlock          int

	// VM options
	VMWasmType        string
	VmTimeoutDuration uint64

	// block config
	BodyCacheLimit           int
	BlockCacheLimit          int
	MaxFutureBlocks          int
	BadBlockLimit            int
	TriesInMemory            int
	BlockChainVersion        int // BlockChainVersion ensures that an incompatible database forces a resync from scratch.
	DefaultTxsCacheSize      int
	DefaultBroadcastInterval time.Duration

	Debug bool
}

type configMarshaling struct {
	MinerExtraData hexutil.Bytes
}
