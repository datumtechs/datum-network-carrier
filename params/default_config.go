package params

import "github.com/mohae/deepcopy"

// use Default config
var carrierConfig = MainnetConfig()

// CarrierConfig retrieves carrier chain config.
func CarrierConfig() *CarrierChainConfig {
	return carrierConfig
}

// OverrideCarrierConfig by replacing the config. The preferred pattern is to
// call BeaconConfig(), change the specific parameters, and then call
// OverrideCarrierConfig(c). Any subsequent calls to params.CarrierConfig() will
// return this new configuration.
func OverrideCarrierConfig(c *CarrierChainConfig) {
	carrierConfig = c
}

// Copy returns a copy of the config object.
func (c *CarrierChainConfig) Copy() *CarrierChainConfig {
	config, ok := deepcopy.Copy(*c).(CarrierChainConfig)
	if !ok {
		config = *carrierConfig
	}
	return &config
}