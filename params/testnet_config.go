package params

import "math/big"

func UseTestnetNetworkConfig() {
	cfg := CarrierNetworkConfig().Copy()
	//TODO: Here, set specific parameters for the testnet, eg:
	cfg.MinimumPeersInSubnet = 0
	OverrideCarrierNetworkConfig(cfg)
}

func UseTestnetConfig() {
	cfg := TestnetConfig()
	OverrideCarrierConfig(cfg)
}

func TestnetConfig() *CarrierChainConfig {
	cfg := MainnetConfig().Copy()
	//TODO: could be set some extra config.
	cfg.ChainID = big.NewInt(1111)	// eg.
	return cfg
}

