package core

import (
	"github.com/datumtechs/datum-network-carrier/types"
)

// Processor is an interface for processing blocks.
type Processor interface {
	Process(block *types.Block, config *types.CarrierChainConfig) error
}