package types

import (
	"github.com/ethereum/go-ethereum/rlp"
	"io"
)

type ResourceTable struct {
	nodeId       string    // node id
	nodeResource *resource // The total resource on the node
	assign       bool      // Whether to assign the slot tag
	slotTotal    uint32    // The total number of slots are allocated on the resource of this node
	slotUsed     uint32    // The number of slots that have been used on the resource of the node
}
type resourceTableRlp struct {
	NodeId    string // node id
	Mem       uint64
	Processor uint64
	Bandwidth uint64
	Assign    bool   // Whether to assign the slot tag
	SlotTotal uint32 // The total number of slots are allocated on the resource of this node
	SlotUsed  uint32 // The number of slots that have been used on the resource of the node
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

// EncodeRLP implements rlp.Encoder.
func (r *ResourceTable) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, resourceTableRlp{
		NodeId:    r.nodeId,
		Mem:       r.nodeResource.mem,
		Processor: r.nodeResource.processor,
		Bandwidth: r.nodeResource.bandwidth,
		Assign:    r.assign,
		SlotTotal: r.slotTotal,
		SlotUsed:  r.slotUsed,
	})
}

// DecodeRLP implements rlp.Decoder.
func (r *ResourceTable) DecodeRLP(s *rlp.Stream) error {
	var dec resourceTableRlp
	err := s.Decode(&dec)
	if err == nil {
		nodeResource := &resource{mem: dec.Mem, processor: dec.Processor, bandwidth: dec.Bandwidth}
		r.nodeId, r.assign, r.slotTotal, r.slotUsed, r.nodeResource =
			dec.NodeId, dec.Assign, dec.SlotTotal, dec.SlotUsed, nodeResource
	}
	return err
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
