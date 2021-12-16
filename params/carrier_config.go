package params

import (
	types "github.com/prysmaticlabs/eth2-types"
	"math/big"
)

// CarrierChainConfig is the core config which determines the datachain settings.
type CarrierChainConfig struct {
	ChainID        *big.Int
	GrpcUrl        string     // GrpcUrl is the host of the data service center.
	Port           uint64     // Port is the port of the data service center.
	SlotsPerEpoch  types.Slot // SlotsPerEpoch is the number of slots in an epoch.
	SecondsPerSlot uint64     // SecondsPerSlot is how many seconds are in a single slot.

	DiscoveryServiceName string
	DiscoveryServiceTags []string
}
