package params

import (
	"fmt"
	types "github.com/prysmaticlabs/eth2-types"
	"math/big"
)

const (
	DefaultHTTPHost = "localhost" // Default host interface for the HTTP RPC server
	DefaultHTTPPort = 8545        // Default TCP port for the HTTP RPC server
)

// DataChainConfig is the core config which determines the datachain settings.
type DataChainConfig struct {
	ChainID *big.Int `json:"chainId"`

	SlotsPerEpoch                    types.Slot  `yaml:"SLOTS_PER_EPOCH" spec:"true"`                     // SlotsPerEpoch is the number of slots in an epoch.
}

// String implements the fmt.Stringer interface.
func (c *DataChainConfig) String() string {
	return fmt.Sprintf("{%v}", "config")
}

func CarrierChainConfig() *DataChainConfig {
	//TODO: need to coding....
	return &DataChainConfig{}
}

// DataCenterConfig is the datacenter service config.
type DataCenterConfig struct {
	GrpcUrl string
	Port uint64
}
