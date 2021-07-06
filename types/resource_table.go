package types

import (
	"github.com/ethereum/go-ethereum/rlp"
	"io"
)

var (
	multipleNumber  = uint64(3)
	DefaultSlotUnit = &Slot{
		Mem:       1024,
		Processor: 1,
		Bandwidth: 1024,
	}
	DefaultResouece = &resource{
		mem:       1024 * multipleNumber,
		processor: 1 * multipleNumber,
		bandwidth: 1024 * multipleNumber,
	}
)

type LocalResourceTable struct {
	nodeId       string    // node id
	nodeResource *resource // The total resource on the node
	assign       bool      // Whether to assign the slot tag
	slotTotal    uint32    // The total number of slots are allocated on the resource of this node
	slotLocked   uint32    // Maybe we will to use, so lock first.
	slotUsed     uint32    // The number of slots that have been used on the resource of the node
}
type localResourceTableRlp struct {
	NodeId    string // node id
	Mem       uint64
	Processor uint64
	Bandwidth uint64
	Assign    bool   // Whether to assign the slot tag
	SlotTotal uint32 // The total number of slots are allocated on the resource of this node
	SlotUsed  uint32 // The number of slots that have been used on the resource of the node
}

func NewLocalResourceTable(nodeId string, mem, processor, bandwidth uint64) *LocalResourceTable {
	return &LocalResourceTable{
		nodeId: nodeId,
		//nodeResource: &resource{
		//	mem:       mem,
		//	processor: processor,
		//	bandwidth: bandwidth,
		//},
		nodeResource: DefaultResouece, // TODO for test
		assign:       false,
	}
}
func (r *LocalResourceTable) GetNodeId() string    { return r.nodeId }
func (r *LocalResourceTable) GetMem() uint64       { return r.nodeResource.mem }
func (r *LocalResourceTable) GetProcessor() uint64 { return r.nodeResource.processor }
func (r *LocalResourceTable) GetBandwidth() uint64 { return r.nodeResource.bandwidth }
func (r *LocalResourceTable) GetAssign() bool      { return r.assign }
func (r *LocalResourceTable) GetSlotTotal() uint32 { return r.slotTotal }
func (r *LocalResourceTable) GetSlotUsed() uint32  { return r.slotUsed }
func (r *LocalResourceTable) SetSlotUnit(slot *Slot) {
	memCount := r.nodeResource.mem / slot.Mem
	processorCount := r.nodeResource.processor / slot.Processor
	bandwidthCount := r.nodeResource.bandwidth / slot.Bandwidth

	min := min3number(memCount, processorCount, bandwidthCount)

	r.slotTotal = uint32(min)
}

func (r *LocalResourceTable) RemianSlot() uint32 { return r.slotTotal - r.slotUsed - r.slotLocked }
func (r *LocalResourceTable) UseSlot(count uint32) {
	if r.slotLocked < count {
		return
	}
	r.slotUsed += count
	r.slotLocked -= count

}
func (r *LocalResourceTable) LockSlot(count uint32) {
	if count > r.slotTotal-r.slotUsed-r.slotLocked {
		r.slotLocked = r.slotTotal - r.slotUsed
	} else {
		r.slotLocked += count
	}
}
func (r *LocalResourceTable) GetTotalSlot() uint32  { return r.slotTotal }
func (r *LocalResourceTable) GetUsedSlot() uint32   { return r.slotUsed }
func (r *LocalResourceTable) GetLockedSlot() uint32 { return r.slotLocked }

// EncodeRLP implements rlp.Encoder.
func (r *LocalResourceTable) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, localResourceTableRlp{
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
func (r *LocalResourceTable) DecodeRLP(s *rlp.Stream) error {
	var dec localResourceTableRlp
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

// Other org total resource item
type RemoteResourceTable struct {
	// other org identityId
	identityId string
	// other org total resource
	total *resource
	// other org be used resource
	used *resource
}
type remoteResourceTableRlp struct {
	// other org identityId
	IdentityId     string
	TotalMem       uint64
	TotalProcessor uint64
	TotalBandwidth uint64
	UsedMem        uint64
	UsedProcessor  uint64
	UsedBandwidth  uint64
}

func NewRemoteResourceTable(identityId string, mem, usedMem, processor, usedProcessor, bandwidth, usedBandwidth uint64) *RemoteResourceTable {
	return &RemoteResourceTable{
		identityId: identityId,
		total: &resource{
			mem:       mem,
			processor: processor,
			bandwidth: bandwidth,
		},
		used: &resource{
			mem:       usedMem,
			processor: usedProcessor,
			bandwidth: usedBandwidth,
		},
	}
}

// EncodeRLP implements rlp.Encoder.
func (r *RemoteResourceTable) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, remoteResourceTableRlp{
		IdentityId:     r.identityId,
		TotalMem:       r.total.mem,
		TotalProcessor: r.total.processor,
		TotalBandwidth: r.total.bandwidth,
		UsedMem:        r.used.mem,
		UsedProcessor:  r.used.processor,
		UsedBandwidth:  r.used.bandwidth,
	})
}

// DecodeRLP implements rlp.Decoder.
func (r *RemoteResourceTable) DecodeRLP(s *rlp.Stream) error {
	var dec remoteResourceTableRlp
	err := s.Decode(&dec)
	if err == nil {
		totalResource := &resource{mem: dec.TotalMem, processor: dec.TotalProcessor, bandwidth: dec.TotalBandwidth}
		usedResource := &resource{mem: dec.UsedMem, processor: dec.UsedProcessor, bandwidth: dec.UsedBandwidth}
		r.identityId, r.total, r.used = dec.IdentityId, totalResource, usedResource
	}
	return err
}
func (r *RemoteResourceTable) remain() (uint64, uint64, uint64) {
	return r.total.mem - r.used.mem, r.total.processor - r.used.processor, r.total.bandwidth - r.used.bandwidth
}
func (r *RemoteResourceTable) GetIdentityId() string     { return r.identityId }
func (r *RemoteResourceTable) GetTotalMem() uint64       { return r.total.mem }
func (r *RemoteResourceTable) GetTotalProcessor() uint64 { return r.total.processor }
func (r *RemoteResourceTable) GetTotalBandwidth() uint64 { return r.total.bandwidth }
func (r *RemoteResourceTable) GetUsedMem() uint64        { return r.used.mem }
func (r *RemoteResourceTable) GetUsedProcessor() uint64  { return r.used.processor }
func (r *RemoteResourceTable) GetUsedBandwidth() uint64  { return r.used.bandwidth }
