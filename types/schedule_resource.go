package types

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/mathutil"
)

type Slot struct {
	Mem       uint64
	Processor uint32
	Bandwidth uint64
}

func (s *Slot) CalculateSlotCount (mem, bandwidth uint64, processor uint32) uint64 {

	memCount := mathutil.DivCeilUint64(mem, s.GetMem())
	processorCount := mathutil.DivCeilUint32(processor, s.GetProcessor())
	bandwidthCount := mathutil.DivCeilUint64(bandwidth, s.GetBandwidth())

	return  mathutil.Max3number(memCount, uint64(processorCount), bandwidthCount)
}

func (s *Slot)String() string  {
	return fmt.Sprintf(`{"mem": %d, "processor": %d, "bandwidth": %d}`, s.Mem, s.Processor, s.Bandwidth)
}

func (s *Slot) GetMem() uint64  { return s.Mem }
func (s *Slot) GetProcessor() uint32  { return s.Processor }
func (s *Slot) GetBandwidth() uint64  { return s.Bandwidth }
