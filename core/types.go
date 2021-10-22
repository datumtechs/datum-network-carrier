package core

import (
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/RosettaFlow/Carrier-Go/types"
)

// Processor is an interface for processing blocks.
type Processor interface {
	Process(block *types.Block, config *params.CarrierChainConfig) error
}