package node

import (
	"github.com/RosettaFlow/Carrier-Go/carrier"
	"github.com/RosettaFlow/Carrier-Go/common/flags"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/urfave/cli/v2"
)

const (
	clientIdentifier = "carrier" // Client identifier to advertise over the network
)

type carrierConfig struct {
	Carrier   carrier.Config
	Node      Config
	// more config modules
}

func makeConfig(cliCtx *cli.Context) carrierConfig {
	// Load defaults.
	cfg := carrierConfig{
		Carrier:   carrier.DefaultConfig,
		Node:      defaultNodeConfig(),
	}

	// todo: file conf load for config.

	// Apply flags.
	SetNodeConfig(cliCtx, &cfg.Node)
	SetCarrierConfig(cliCtx, &cfg.Carrier)

	return cfg
}

func defaultNodeConfig() Config {
	cfg := DefaultConfig
	return cfg
}

func configureNetwork(cliCtx *cli.Context) {
	if cliCtx.IsSet(flags.BootstrapNode.Name) {
		c := params.CarrierNetworkConfig()
		c.BootstrapNodes = cliCtx.StringSlice(flags.BootstrapNode.Name)
		params.OverrideCarrierNetworkConfig(c)
	}
}