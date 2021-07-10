package params

import (
	"github.com/mohae/deepcopy"
	types "github.com/prysmaticlabs/eth2-types"
	"math/big"
)

// DataChainConfig is the core config which determines the datachain settings.
type DataChainConfig struct {
	ChainID *big.Int `json:"chainId"`

	SlotsPerEpoch  types.Slot `yaml:"SLOTS_PER_EPOCH" spec:"true"`  // SlotsPerEpoch is the number of slots in an epoch.
	SecondsPerSlot uint64     `yaml:"SECONDS_PER_SLOT" spec:"true"` // SecondsPerSlot is how many seconds are in a single slot.
}

// DataCenterConfig is the datacenter service config.
type DataCenterConfig struct {
	GrpcUrl string
	Port    uint64
}

var carrierConfig = MainnetConfig()

// MainnetConfig returns the configuration to be used in the main network.
func MainnetConfig() *DataChainConfig {
	return mainnetCarrierConfig
}

func TestnetConfig() *DataChainConfig {
	return testnetCarrierConfig
}

func UseMainnetConfig() {
	carrierConfig = MainnetConfig()
}

// CarrierConfig retrieves carrier chain config.
func CarrierChainConfig() *DataChainConfig {
	return carrierConfig
}

// OverrideCarrierConfig by replacing the config. The preferred pattern is to
// call BeaconConfig(), change the specific parameters, and then call
// OverrideCarrierConfig(c). Any subsequent calls to params.CarrierConfig() will
// return this new configuration.
func OverrideCarrierConfig(c *DataChainConfig) {
	carrierConfig = c
}

// Copy returns a copy of the config object.
func (c *DataChainConfig) Copy() *DataChainConfig {
	config, ok := deepcopy.Copy(*c).(DataChainConfig)
	if !ok {
		config = *carrierConfig
	}
	return &config
}
