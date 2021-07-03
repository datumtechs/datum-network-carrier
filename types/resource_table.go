package types

type ResourceTable struct {
	nodeId       string    // node id
	nodeResource *resource // The total resource on the node
	assign       bool      // Whether to assign the slot tag
	slotTotal    uint32    // The total number of slots are allocated on the resource of this node
	slotUsed     uint32    // The number of slots that have been used on the resource of the node
}

func NewResourceTable(nodeId string, mem, processor, bandwidth uint64) *ResourceTable {
	return &ResourceTable{
		nodeId: nodeId,
		nodeResource: &resource{
			mem:       mem,
			processor: processor,
			bandwidth: bandwidth,
		},
		assign: false,
	}
}
func (r *ResourceTable) GetNodeId() string    { return r.nodeId }
func (r *ResourceTable) GetMem() uint64       { return r.nodeResource.mem }
func (r *ResourceTable) GetProcessor() uint64 { return r.nodeResource.processor }
func (r *ResourceTable) GetBandwidth() uint64 { return r.nodeResource.bandwidth }
func (r *ResourceTable) GetAssign() bool      { return r.assign }
func (r *ResourceTable) GetSlotTotal() uint32 { return r.slotTotal }
func (r *ResourceTable) GetSlotUsed() uint32  { return r.slotUsed }
func (r *ResourceTable) SetSlotUnit(slot *Slot) {
	memCount := r.nodeResource.mem / slot.Mem
	processorCount := r.nodeResource.processor / slot.Processor
	bandwidthCount := r.nodeResource.bandwidth / slot.Bandwidth

	min := min3number(memCount, processorCount, bandwidthCount)

	r.slotTotal = uint32(min)
}

func (r *ResourceTable) RemianSlot() uint32 { return r.slotTotal - r.slotUsed }
func (r *ResourceTable) UseSlot(count uint32) {
	if count > r.slotTotal-r.slotUsed {
		r.slotUsed = r.slotTotal
	} else {
		r.slotUsed += count
	}
}

func min3number(a, b, c uint64) uint64 {
	var min uint64
	if a > b {
		min = b
	} else {
		min = a
	}

	if min > c {
		min = c
	}
	return min
}

type resource struct {
	mem       uint64
	processor uint64
	bandwidth uint64
}
