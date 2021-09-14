package slotutil

import (
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/params"
	types "github.com/prysmaticlabs/eth2-types"
	"time"
)

// SlotStartTime returns the start time in terms of its unix epoch value.
func SlotStartTime(genesis uint64, slot types.Slot) time.Time {
	duration := time.Second * time.Duration(slot.Mul(params.CarrierChainConfig().SecondsPerSlot))
	startTime := time.Unix(int64(genesis), 0).Add(duration)
	return startTime
}

// SlotsSinceGenesis returns the number of slots since
// the provided genesis time.
func SlotsSinceGenesis(genesis time.Time) types.Slot {
	if genesis.After(timeutils.Now()) { // Genesis has not occurred yet.
		return 0
	}
	return types.Slot(uint64(timeutils.Since(genesis).Seconds()) / params.CarrierChainConfig().SecondsPerSlot)
}

// EpochsSinceGenesis returns the number of slots since
// the provided genesis time.
func EpochsSinceGenesis(genesis time.Time) types.Epoch {
	return types.Epoch(SlotsSinceGenesis(genesis) / params.CarrierChainConfig().SlotsPerEpoch)
}

// DivideSlotBy divides the SECONDS_PER_SLOT configuration
// parameter by a specified number. It returns a value of time.Duration
// in milliseconds, useful for dividing values such as 1 second into
// millisecond-based durations.
func DivideSlotBy(timesPerSlot int64) time.Duration {
	return time.Duration(int64(params.CarrierChainConfig().SecondsPerSlot*1000)/timesPerSlot) * time.Millisecond
}
