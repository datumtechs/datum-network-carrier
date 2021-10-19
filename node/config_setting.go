package node

import (
	"github.com/RosettaFlow/Carrier-Go/carrier"
	"github.com/RosettaFlow/Carrier-Go/common/flags"
	"github.com/RosettaFlow/Carrier-Go/common/tracing"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/urfave/cli/v2"
)

const (
	clientIdentifier = "carrier" // Client identifier to advertise over the network
)

type carrierConfig struct {
	Carrier carrier.Config
	Node    Config
	// more config modules
	MockIdentityIdsFile string
	ConsensusStateFile  string
}

func makeConfig(cliCtx *cli.Context) carrierConfig {
	// Load defaults.
	cfg := carrierConfig{
		Carrier:             carrier.DefaultConfig,
		Node:                defaultNodeConfig(),
		MockIdentityIdsFile: cliCtx.String(flags.MockIdentityIdFileFlag.Name),
		ConsensusStateFile:  cliCtx.String(flags.ConsensusStateStoreFlag.Name),
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

func configureTracing(cliCtx *cli.Context) error {
	return tracing.Setup(
		"carrier-network", // service name
		cliCtx.String(flags.TracingProcessNameFlag.Name),
		cliCtx.String(flags.TracingEndpointFlag.Name),
		cliCtx.Float64(flags.TraceSampleFractionFlag.Name),
		cliCtx.Bool(flags.EnableTracingFlag.Name),
	)
}
