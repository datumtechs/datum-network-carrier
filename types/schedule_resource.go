package types

type Slot struct {
	Mem       uint64
	Processor uint64
	Bandwidth uint64
}

func (s *Slot) CalculateSlotCount (mem, processor, bandwidth uint64) uint64 {
	memCount := mem / s.Mem
	processorCount := processor / s.Processor
	bandwidthCount := bandwidth / s.Bandwidth
	return min3number(memCount, processorCount, bandwidthCount)
}

