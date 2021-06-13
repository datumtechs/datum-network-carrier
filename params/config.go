package params

import (
	"fmt"
	"math/big"
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
