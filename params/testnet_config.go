package params

import (
	"github.com/datumtechs/datum-network-carrier/types"
	"math/big"
)

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

func TestnetConfig() *types.CarrierChainConfig {
	cfg := MainnetConfig().Copy()
	//TODO: could be set some extra config.
	cfg.ChainID = big.NewInt(1111)	// eg.
	return cfg
}

