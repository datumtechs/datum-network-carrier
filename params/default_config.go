package params

import (
	"github.com/RosettaFlow/Carrier-Go/types"
)

// use Default config
var carrierConfig = MainnetConfig()

// CarrierConfig retrieves carrier chain config.
func CarrierConfig() *types.CarrierChainConfig {
	return carrierConfig
}

// OverrideCarrierConfig by replacing the config. The preferred pattern is to
// call BeaconConfig(), change the specific parameters, and then call
// OverrideCarrierConfig(c). Any subsequent calls to params.CarrierConfig() will
// return this new configuration.
func OverrideCarrierConfig(c *types.CarrierChainConfig) {
	carrierConfig = c
}