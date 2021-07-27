package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
)

var (
	// TODO 写死的 资源固定消耗 ...
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

	DefaultDisk = uint64(10*1024*1024*1024*1024*1024)  //10pb
)
func GetDefaultResoueceMem () uint64 { return DefaultResouece.mem }
func GetDefaultResoueceProcessor () uint64 { return DefaultResouece.processor }
func GetDefaultResoueceBandwidth () uint64 { return DefaultResouece.bandwidth }

type LocalResourceTable struct {
	nodeId       string    // Resource node id
	powerId      string
	nodeResource *resource // The total resource on the node
	assign       bool      // Whether to assign the slot tag
	slotTotal    uint32    // The total number of slots are allocated on the resource of this node
	//slotLocked   uint32    // Maybe we will to use, so lock first.
	slotUsed     uint32    // The number of slots that have been used on the resource of the node
}
type localResourceTableRlp struct {
	NodeId    string // node id
	PowerId   string
	Mem       uint64
	Processor uint64
	Bandwidth uint64
	Assign    bool   // Whether to assign the slot tag
	SlotTotal uint32 // The total number of slots are allocated on the resource of this node
	SlotUsed  uint32 // The number of slots that have been used on the resource of the node
}

func NewLocalResourceTable(nodeId, powerId string, mem, processor, bandwidth uint64) *LocalResourceTable {
	return &LocalResourceTable{
		nodeId: nodeId,
		powerId: powerId,
		//nodeResource: &resource{
		//	mem:       mem,
		//	processor: processor,
		//	bandwidth: bandwidth,
		//},
		nodeResource: DefaultResouece, // TODO for test
		assign:       false,
	}
}

func (r *LocalResourceTable) String() string {
	return fmt.Sprintf(`{"nodeId": "%s", "powerId": "%s", "nodeResource": {"mem": %d, "processor": %d, "bandwidth": %d}, "assign": %v}`,
		r.nodeId, r.powerId, r.nodeResource.mem, r.nodeResource.processor, r.nodeResource.bandwidth, r.assign)
}
func (r *LocalResourceTable) GetNodeId() string    { return r.nodeId }
func (r *LocalResourceTable) GetPowerId() string    { return r.powerId }
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

func (r *LocalResourceTable) RemianSlot() uint32 { return r.slotTotal - r.slotUsed /*- r.slotLocked*/ }
func (r *LocalResourceTable) UseSlot(count uint32) error {

	if r.RemianSlot() < count {
		return fmt.Errorf("Failed to lock local resource, slotRemain {%s} less than need lock count {%s}", r.RemianSlot(), count)
	}
	r.slotUsed += count
	r.assign = true
	return nil
}
func (r *LocalResourceTable) FreeSlot(count uint32) error {
	if !r.assign {
		return nil
	}
	if r.slotUsed == 0 || r.slotUsed < count {
		return fmt.Errorf("Failed to unlock local resource, slotUsed {%s} less than need free count {%s}", r.slotUsed, count)
	} else {
		r.slotUsed -= count
	}
	if r.slotUsed == 0 {
		r.assign = false
	}
	return nil
}
//func (r *LocalResourceTable) LockSlot(count uint32) {
//	if count > r.RemianSlot() {
//		r.slotLocked = r.slotTotal - r.slotUsed
//	} else {
//		r.slotLocked += count
//	}
//}
//func (r *LocalResourceTable) UnLockSlot(count uint32) {
//	if r.slotLocked <= count {
//		r.slotLocked = 0
//	} else {
//		r.slotLocked = r.slotLocked - count
//	}
//}
func (r *LocalResourceTable) GetTotalSlot() uint32  { return r.slotTotal }
func (r *LocalResourceTable) GetUsedSlot() uint32   { return r.slotUsed }
//func (r *LocalResourceTable) GetLockedSlot() uint32 { return r.slotLocked }
func (r *LocalResourceTable) IsEnough(slotCount uint32) bool {
	if r.RemianSlot() < slotCount {
		return false
	}
	return true
}

// EncodeRLP implements rlp.Encoder.
func (r *LocalResourceTable) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, localResourceTableRlp{
		NodeId:    r.nodeId,
		PowerId:   r.powerId,
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
		r.nodeId, r.powerId, r.assign, r.slotTotal, r.slotUsed, r.nodeResource =
			dec.NodeId, dec.PowerId, dec.Assign, dec.SlotTotal, dec.SlotUsed, nodeResource
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
func (r *RemoteResourceTable) Remain() (uint64, uint64, uint64) {
	return r.total.mem - r.used.mem, r.total.processor - r.used.processor, r.total.bandwidth - r.used.bandwidth
}
func (r *RemoteResourceTable) IsEnough(mem, processor, bandwidth uint64) bool {
	rMem, rProcessor, rBandwidth := r.Remain()
	if rMem < mem {
		return false
	}
	if rProcessor < processor {
		return false
	}
	if rBandwidth < bandwidth {
		return false
	}
	return true
}
func (r *RemoteResourceTable) GetIdentityId() string     { return r.identityId }
func (r *RemoteResourceTable) GetTotalMem() uint64       { return r.total.mem }
func (r *RemoteResourceTable) GetTotalProcessor() uint64 { return r.total.processor }
func (r *RemoteResourceTable) GetTotalBandwidth() uint64 { return r.total.bandwidth }
func (r *RemoteResourceTable) GetUsedMem() uint64        { return r.used.mem }
func (r *RemoteResourceTable) GetUsedProcessor() uint64  { return r.used.processor }
func (r *RemoteResourceTable) GetUsedBandwidth() uint64  { return r.used.bandwidth }

// 给本地 缓存用的

// 本地任务所占用的 资源缓存
type LocalTaskPowerUsed struct {
	taskId    string
	nodeId    string
	slotCount uint64
}
type localTaskPowerUsedRlp struct {
	TaskId    string
	NodeId    string
	SlotCount uint64
}

func NewLocalTaskPowerUsed(taskId, nodeId string, slotCount uint64) *LocalTaskPowerUsed {
	return &LocalTaskPowerUsed{
		taskId:    taskId,
		nodeId:    nodeId,
		slotCount: slotCount,
	}
}

// EncodeRLP implements rlp.Encoder.
func (pcache *LocalTaskPowerUsed) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, localTaskPowerUsedRlp{
		TaskId:    pcache.taskId,
		NodeId:    pcache.nodeId,
		SlotCount: pcache.slotCount,
	})
}

// DecodeRLP implements rlp.Decoder.
func (pcache *LocalTaskPowerUsed) DecodeRLP(s *rlp.Stream) error {
	var dec localTaskPowerUsedRlp
	err := s.Decode(&dec)
	if err == nil {
		pcache.taskId, pcache.nodeId, pcache.slotCount = dec.TaskId, dec.NodeId, dec.SlotCount
	}
	return err
}
func (pcache *LocalTaskPowerUsed) GetTaskId() string    { return pcache.taskId }
func (pcache *LocalTaskPowerUsed) GetNodeId() string    { return pcache.nodeId }
func (pcache *LocalTaskPowerUsed) GetSlotCount() uint64 { return pcache.slotCount }

type DataResourceTable struct {
	nodeId    string
	totalDisk uint64
	usedDisk  uint64
}
type dataResourceTableRlp struct {
	NodeId    string
	TotalDisk uint64
	UsedDisk  uint64
}

func NewDataResourceTable(nodeId string, totalDisk, usedDisk uint64) *DataResourceTable {
	return &DataResourceTable{
		nodeId:    nodeId,
		//totalDisk: totalDisk,
		//usedDisk:  usedDisk,
		totalDisk: DefaultDisk,
		usedDisk:  0,
	}
}

// EncodeRLP implements rlp.Encoder.
func (drt *DataResourceTable) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, dataResourceTableRlp{
		NodeId:    drt.nodeId,
		TotalDisk: drt.totalDisk,
		UsedDisk:  drt.usedDisk,
	})
}

// DecodeRLP implements rlp.Decoder.
func (drt *DataResourceTable) DecodeRLP(s *rlp.Stream) error {
	var dec dataResourceTableRlp
	err := s.Decode(&dec)
	if err == nil {
		drt.nodeId, drt.totalDisk, drt.usedDisk = dec.NodeId, dec.TotalDisk, dec.UsedDisk
	}
	return err
}
func (drt *DataResourceTable) GetNodeId() string    { return drt.nodeId }
func (drt *DataResourceTable) GetTotalDisk() uint64 { return drt.totalDisk }
func (drt *DataResourceTable) GetUsedDisk() uint64  { return drt.usedDisk }
func (drt *DataResourceTable) RemainDisk() uint64 { return drt.totalDisk - drt.usedDisk }
func (drt *DataResourceTable) IsUsed() bool { return drt.usedDisk != 0 }
func (drt *DataResourceTable) IsNotUsed() bool { return !drt.IsUsed() }
func (drt *DataResourceTable) IsEmpty() bool { return nil == drt }
func (drt *DataResourceTable) IsNotEmpty() bool { return !drt.IsEmpty() }
func (drt *DataResourceTable) UseDisk(use uint64)  {
	if drt.RemainDisk() > use {
		drt.usedDisk += use
	} else {
		drt.usedDisk = drt.totalDisk
	}
}
func (drt *DataResourceTable) FreeDisk(use uint64) {
	if drt.usedDisk > use {
		drt.usedDisk -= use
	} else {
		drt.usedDisk = 0
	}
}

type DataResourceFileUpload struct {
	originId   string   // db key
	nodeId     string
	metaDataId string
	filePath   string
}

type dataResourceFileUploadRlp struct {
	NodeId     string
	OriginId   string
	MetaDataId string
	FilePath   string
}

func NewDataResourceFileUpload(nodeId, originId, metaDataId, filePath string) *DataResourceFileUpload {
	return &DataResourceFileUpload{
		nodeId:     nodeId,
		originId:   originId,
		metaDataId: metaDataId,
		filePath:   filePath,
	}
}

// EncodeRLP implements rlp.Encoder.
func (drt *DataResourceFileUpload) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, dataResourceFileUploadRlp{
		NodeId:     drt.nodeId,
		OriginId:   drt.originId,
		MetaDataId: drt.metaDataId,
		FilePath:   drt.filePath,
	})
}

// DecodeRLP implements rlp.Decoder.
func (drt *DataResourceFileUpload) DecodeRLP(s *rlp.Stream) error {
	var dec dataResourceFileUploadRlp
	err := s.Decode(&dec)
	if err == nil {
		drt.nodeId, drt.originId, drt.metaDataId, drt.filePath = dec.NodeId, dec.OriginId, dec.MetaDataId, dec.FilePath
	}
	return err
}
func (drt *DataResourceFileUpload) GetNodeId() string               { return drt.nodeId }
func (drt *DataResourceFileUpload) GetOriginId() string             { return drt.originId }
func (drt *DataResourceFileUpload) SetMetaDataId(metaDataId string) { drt.metaDataId = metaDataId }
func (drt *DataResourceFileUpload) GetMetaDataId() string           { return drt.metaDataId }
func (drt *DataResourceFileUpload) GetFilePath() string             { return drt.filePath }



type DataResourceDiskUsed struct {
	metaDataId string   // db key
	nodeId     string
	diskUsed   uint64
}

type dataResourceDiskUsedRlp struct {
	MetaDataId string   // db key
	NodeId     string
	DiskUsed   uint64
}

func NewDataResourceDiskUsed(metaDataId, nodeId string, diskUsed uint64) *DataResourceDiskUsed {
	return &DataResourceDiskUsed{
		metaDataId: metaDataId,
		nodeId:     nodeId,
		diskUsed:  diskUsed,
	}
}

// EncodeRLP implements rlp.Encoder.
func (drt *DataResourceDiskUsed) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, dataResourceDiskUsedRlp{
		MetaDataId: drt.metaDataId,
		NodeId:     drt.nodeId,
		DiskUsed:   drt.diskUsed,
	})
}

// DecodeRLP implements rlp.Decoder.
func (drt *DataResourceDiskUsed) DecodeRLP(s *rlp.Stream) error {
	var dec dataResourceDiskUsedRlp
	err := s.Decode(&dec)
	if err == nil {
		drt.metaDataId, drt.nodeId,  drt.diskUsed = dec.MetaDataId,  dec.NodeId, dec.DiskUsed
	}
	return err
}
func (drt *DataResourceDiskUsed) GetMetaDataId() string { return drt.metaDataId }
func (drt *DataResourceDiskUsed) GetNodeId() string     { return drt.nodeId }
func (drt *DataResourceDiskUsed) GetDiskUsed() uint64   { return drt.diskUsed }

