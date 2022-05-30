package types

import (
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/mohae/deepcopy"
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
	// add by v 0.3.0
	DiscoveryServiceName string
	DiscoveryServiceTags []string
	// add by v 0.4.0
	BlockChainIdCache     map[commonconstantpb.UserType]*big.Int   // blockchain {userAddress -> chainId}
}

// Copy returns a copy of the config object.
func (c *CarrierChainConfig) Copy() *CarrierChainConfig {
	config, ok := deepcopy.Copy(*c).(CarrierChainConfig)
	if !ok {
		return nil
	}
	return &config
}
