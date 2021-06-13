package params

import (
	"fmt"
	"math/big"
)

const (
	DefaultHTTPHost = "localhost" // Default host interface for the HTTP RPC server
	DefaultHTTPPort = 8545        // Default TCP port for the HTTP RPC server
)


// DataChainConfig is the core config which determines the datachain settings.
type DataChainConfig struct {
	ChainID *big.Int `json:"chainId"`
}

// String implements the fmt.Stringer interface.
func (c *DataChainConfig) String() string {
	return fmt.Sprintf("{%v}", "config")
}

// DataCenterConfig is the datacenter service config.
type DataCenterConfig struct {
	GrpcUrl string
	Port uint64
}
