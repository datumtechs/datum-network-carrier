package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
)

type LocalResourceTable struct {
	nodeId  string    // resource node id
	powerId string    // powerId
	total   *resource // The total resource on the node
	used    *resource // the used resource on the node
	assign  bool      // whether to assign the slot tag
	alive   bool      // indicates whether the corresponding jobnode leaves (e.g. downtime)
}
type localResourceTableRlp struct {
	NodeId         string // node id
	PowerId        string
	TotalProcessor uint32
	UsedProcessor  uint32
	TotalMem       uint64
	TotalBandwidth uint64
	TotalDisk      uint64
	UsedMem        uint64
	UsedBandwidth  uint64
	UsedDisk       uint64
	Assign         bool // Whether to assign the slot tag
	Alive          bool // Indicates whether the corresponding jobnode leaves (e.g. downtime)
}

func NewLocalResourceTable(nodeId, powerId string, mem, bandwidth, disk uint64, processor uint32, alive bool) *LocalResourceTable {
	return &LocalResourceTable{
		nodeId:  nodeId,
		powerId: powerId,
		total:   newResource(mem, bandwidth, disk, processor),
		used:    &resource{},
		assign:  false,
		alive:   alive,
	}
}

func (r *LocalResourceTable) String() string {
	return fmt.Sprintf(`{"nodeId": "%s", "powerId": "%s", "total": %s, "used": %s, "assign": %v}`,
		r.nodeId, r.powerId, r.total.String(), r.used.String(), r.assign)
}
func (r *LocalResourceTable) GetNodeId() string         { return r.nodeId }
func (r *LocalResourceTable) GetPowerId() string        { return r.powerId }
func (r *LocalResourceTable) GetTotalMem() uint64       { return r.total.mem }
func (r *LocalResourceTable) GetTotalProcessor() uint32 { return r.total.processor }
func (r *LocalResourceTable) GetTotalBandwidth() uint64 { return r.total.bandwidth }
func (r *LocalResourceTable) GetTotalDisk() uint64      { return r.total.disk }
func (r *LocalResourceTable) GetUsedMem() uint64        { return r.used.mem }
func (r *LocalResourceTable) GetUsedProcessor() uint32  { return r.used.processor }
func (r *LocalResourceTable) GetUsedBandwidth() uint64  { return r.used.bandwidth }
func (r *LocalResourceTable) GetUsedDisk() uint64       { return r.used.disk }
func (r *LocalResourceTable) GetAssign() bool           { return r.assign }
func (r *LocalResourceTable) GetAlive() bool            { return r.alive }
func (r *LocalResourceTable) SetAlive(alive bool)       { r.alive = alive }
func (r *LocalResourceTable) RemainProcessor() uint32   { return r.total.processor - r.used.processor }
func (r *LocalResourceTable) RemainMem() uint64         { return r.total.mem - r.used.mem }
func (r *LocalResourceTable) RemainBandwidth() uint64   { return r.total.bandwidth - r.used.bandwidth }
func (r *LocalResourceTable) RemainDisk() uint64        { return r.total.disk - r.used.disk }
func (r *LocalResourceTable) UseSlot(mem, bandwidth, disk uint64, processor uint32) error {

	if mem == 0 && bandwidth == 0 && disk == 0 && processor == 0 {
		return nil
	}

	if r.RemainMem() < mem {
		return fmt.Errorf("Failed to lock local resource, mem remain {%d} less than need lock count {%d}", r.RemainMem(), mem)
	}
	if r.RemainBandwidth() < bandwidth {
		return fmt.Errorf("Failed to lock local resource, bandwidth remain {%d} less than need lock count {%d}", r.RemainBandwidth(), bandwidth)
	}
	if r.RemainDisk() < disk {
		return fmt.Errorf("Failed to lock local resource, disk remain {%d} less than need lock count {%d}", r.RemainDisk(), disk)
	}
	if r.RemainProcessor() < processor {
		return fmt.Errorf("Failed to lock local resource, processor remain {%d} less than need lock count {%d}", r.RemainProcessor(), processor)
	}

	log.Debugf("Call UseSlot JobNode LocalResourceTable before, table: %s, needMem: {%d}, needBandwidth: {%d}, needDisk: {%d}, needProcessor: {%d}",
		r.String(), mem, bandwidth, disk, processor)

	r.used.mem += mem
	r.used.bandwidth += bandwidth
	r.used.disk += disk
	r.used.processor += processor

	if r.used.mem > 0 || r.used.bandwidth > 0 || r.used.disk == 0 || r.used.processor > 0 {
		r.assign = true
	}
	log.Debugf("Call UseSlot JobNode LocalResourceTable after, table: %s, needMem: {%d}, needBandwidth: {%d}, needDisk: {%d}, needProcessor: {%d}",
		r.String(), mem, bandwidth, disk, processor)

	return nil
}
func (r *LocalResourceTable) FreeSlot(mem, bandwidth, disk uint64, processor uint32) error {
	if !r.assign {
		return nil
	}

	if mem == 0 && bandwidth == 0 && disk == 0 && processor == 0 {
		return nil
	}
	if r.used.mem == 0 && r.used.bandwidth == 0 && r.used.disk == 0 && r.used.processor == 0 {
		return nil
	}

	if r.used.mem < mem {
		return fmt.Errorf("Failed to unlock local resource, mem used {%d} less than need free count {%d}", r.used.mem, mem)
	}
	if r.used.bandwidth < bandwidth {
		return fmt.Errorf("Failed to unlock local resource, bandwidth used {%d} less than need free count {%d}", r.used.bandwidth, bandwidth)
	}
	if r.used.disk < disk {
		return fmt.Errorf("Failed to unlock local resource, disk used {%d} less than need free count {%d}", r.used.disk, disk)
	}
	if r.used.processor < processor {
		return fmt.Errorf("Failed to unlock local resource, processor used {%d} less than need free count {%d}", r.used.processor, processor)
	}

	log.Debugf("Call FreeSlot JobNode LocalResourceTable before, table: %s, freeMem: {%d}, freeBandwidth: {%d}, freeDisk: {%d}, freeProcessor: {%d}",
		r.String(), mem, bandwidth, disk, processor)

	r.used.mem -= mem
	r.used.bandwidth -= bandwidth
	r.used.disk -= disk
	r.used.processor -= processor

	if r.used.mem == 0 && r.used.bandwidth == 0 && r.used.disk == 0 && r.used.processor == 0 {
		r.assign = false
	}

	log.Debugf("Call FreeSlot JobNode LocalResourceTable after, table: %s, freeMem: {%d}, freeBandwidth: {%d}, freeDisk: {%d}, freeProcessor: {%d}",
		r.String(), mem, bandwidth, disk, processor)

	return nil
}

func (r *LocalResourceTable) IsEnoughMem(mem uint64) bool {
	if r.RemainMem() < mem {
		return false
	}
	return true
}
func (r *LocalResourceTable) IsEnoughBandwidth(bandwidth uint64) bool {
	if r.RemainBandwidth() < bandwidth {
		return false
	}
	return true
}
func (r *LocalResourceTable) IsEnoughDisk(disk uint64) bool {
	if r.RemainDisk() < disk {
		return false
	}
	return true
}
func (r *LocalResourceTable) IsEnoughProcessor(processor uint32) bool {
	if r.RemainProcessor() < processor {
		return false
	}
	return true
}
func (r *LocalResourceTable) IsEnough(mem, bandwidth, disk uint64, processor uint32) bool {
	if !r.IsEnoughMem(mem) {
		return false
	}
	if !r.IsEnoughBandwidth(bandwidth) {
		return false
	}
	if !r.IsEnoughDisk(disk) {
		return false
	}
	if !r.IsEnoughProcessor(processor) {
		return false
	}
	return true
}

// EncodeRLP implements rlp.Encoder.
func (r *LocalResourceTable) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, localResourceTableRlp{
		NodeId:         r.nodeId,
		PowerId:        r.powerId,
		TotalProcessor: r.total.processor,
		TotalMem:       r.total.mem,
		TotalBandwidth: r.total.bandwidth,
		TotalDisk:      r.total.disk,
		UsedProcessor:  r.used.processor,
		UsedMem:        r.used.mem,
		UsedBandwidth:  r.used.bandwidth,
		UsedDisk:       r.used.disk,
		Assign:         r.assign,
		Alive:          r.alive,
	})
}

// DecodeRLP implements rlp.Decoder.
func (r *LocalResourceTable) DecodeRLP(s *rlp.Stream) error {
	var dec localResourceTableRlp
	err := s.Decode(&dec)
	if err == nil {
		total := newResource(dec.TotalMem, dec.TotalBandwidth, dec.TotalDisk, dec.TotalProcessor)
		used := newResource(dec.UsedMem, dec.UsedBandwidth, dec.UsedDisk, dec.UsedProcessor)
		r.nodeId, r.powerId, r.assign, r.total, r.used, r.alive =
			dec.NodeId, dec.PowerId, dec.Assign, total, used, dec.Alive
	}
	return err
}

type resource struct {
	mem       uint64
	bandwidth uint64
	disk      uint64
	processor uint32
}

func newResource(mem, bandwidth, disk uint64, processor uint32) *resource {
	return &resource{
		mem:       mem,
		bandwidth: bandwidth,
		disk:      disk,
		processor: processor,
	}
}

func (r *resource) String() string {
	return fmt.Sprintf(`{"mem": %d, "processor": %d, "bandwidth": %d, "disk": %d}`, r.mem, r.processor, r.bandwidth, r.disk)
}

type LocalTaskPowerUsed struct {
	taskId  string
	partyId string
	nodeId  string
	used    *resource
}
type localTaskPowerUsedRlp struct {
	TaskId        string
	PartyId       string
	NodeId        string
	UsedProcessor uint32
	UsedMem       uint64
	UsedBandwidth uint64
	UsedDisk      uint64
}

func NewLocalTaskPowerUsed(taskId, partyId, nodeId string, mem, bandwidth, disk uint64, processor uint32) *LocalTaskPowerUsed {
	return &LocalTaskPowerUsed{
		taskId:  taskId,
		partyId: partyId,
		nodeId:  nodeId,
		used:    newResource(mem, bandwidth, disk, processor),
	}
}

// EncodeRLP implements rlp.Encoder.
func (pcache *LocalTaskPowerUsed) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, localTaskPowerUsedRlp{
		TaskId:        pcache.taskId,
		PartyId:       pcache.partyId,
		NodeId:        pcache.nodeId,
		UsedProcessor: pcache.used.processor,
		UsedMem:       pcache.used.mem,
		UsedBandwidth: pcache.used.bandwidth,
		UsedDisk:      pcache.used.disk,
	})
}

// DecodeRLP implements rlp.Decoder.
func (pcache *LocalTaskPowerUsed) DecodeRLP(s *rlp.Stream) error {
	var dec localTaskPowerUsedRlp
	err := s.Decode(&dec)
	if err == nil {
		used := newResource(dec.UsedMem, dec.UsedBandwidth, dec.UsedDisk, dec.UsedProcessor)
		pcache.taskId, pcache.partyId, pcache.nodeId, pcache.used = dec.TaskId, dec.PartyId, dec.NodeId, used
	}
	return err
}
func (pcache *LocalTaskPowerUsed) GetTaskId() string        { return pcache.taskId }
func (pcache *LocalTaskPowerUsed) GetPartyId() string       { return pcache.partyId }
func (pcache *LocalTaskPowerUsed) GetNodeId() string        { return pcache.nodeId }
func (pcache *LocalTaskPowerUsed) GetUsedMem() uint64       { return pcache.used.mem }
func (pcache *LocalTaskPowerUsed) GetUsedBandwidth() uint64 { return pcache.used.bandwidth }
func (pcache *LocalTaskPowerUsed) GetUsedDisk() uint64      { return pcache.used.disk }
func (pcache *LocalTaskPowerUsed) GetUsedProcessor() uint32 { return pcache.used.processor }
func (pcache *LocalTaskPowerUsed) String() string {
	return fmt.Sprintf(`{"taskId": %s, "partyId": %s, "nodeId": %s, "used":, %s}`,
		pcache.taskId, pcache.partyId, pcache.nodeId, pcache.used.String())
}

type DataResourceTable struct {
	nodeId    string
	totalDisk uint64
	usedDisk  uint64
	isUsed    bool // when the node is used for the first time (when there is data / file upload), this value is given "true"
	alive     bool // indicates whether the corresponding datanode leaves (e.g. downtime)
}
type dataResourceTableRlp struct {
	NodeId    string
	TotalDisk uint64
	UsedDisk  uint64
	IsUsed    bool
	Alive     bool // Indicates whether the corresponding datanode leaves (e.g. downtime)
}

func NewDataResourceTable(nodeId string, totalDisk, usedDisk uint64, alive bool) *DataResourceTable {
	return &DataResourceTable{
		nodeId:    nodeId,
		totalDisk: totalDisk,
		usedDisk:  usedDisk,
		isUsed:    false,
		alive:     alive,
	}
}

// EncodeRLP implements rlp.Encoder.
func (drt *DataResourceTable) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, dataResourceTableRlp{
		NodeId:    drt.nodeId,
		TotalDisk: drt.totalDisk,
		UsedDisk:  drt.usedDisk,
		IsUsed:    drt.isUsed,
		Alive:     drt.alive,
	})
}

// DecodeRLP implements rlp.Decoder.
func (drt *DataResourceTable) DecodeRLP(s *rlp.Stream) error {
	var dec dataResourceTableRlp
	err := s.Decode(&dec)
	if err == nil {
		drt.nodeId, drt.totalDisk, drt.usedDisk, drt.isUsed, drt.alive = dec.NodeId, dec.TotalDisk, dec.UsedDisk, dec.IsUsed, dec.Alive
	}
	return err
}

func (drt *DataResourceTable) String() string {
	return fmt.Sprintf(`{"nodeId": %s, "totalDisk": %d, "usedDisk": %d, "isUsed": %v}`,
		drt.GetNodeId(), drt.GetTotalDisk(), drt.GetUsedDisk(), drt.GetIsUsed())
}
func (drt *DataResourceTable) GetNodeId() string             { return drt.nodeId }
func (drt *DataResourceTable) GetTotalDisk() uint64          { return drt.totalDisk }
func (drt *DataResourceTable) SetTotalDisk(totalDisk uint64) { drt.totalDisk = totalDisk }
func (drt *DataResourceTable) GetUsedDisk() uint64           { return drt.usedDisk }
func (drt *DataResourceTable) GetIsUsed() bool               { return drt.isUsed }
func (drt *DataResourceTable) RemainDisk() uint64            { return drt.totalDisk - drt.usedDisk }
func (drt *DataResourceTable) IsUsed() bool                  { return drt.isUsed && drt.usedDisk != 0 }
func (drt *DataResourceTable) IsNotUsed() bool               { return !drt.IsUsed() }
func (drt *DataResourceTable) GetAlive() bool                { return drt.alive }
func (drt *DataResourceTable) SetAlive(alive bool)           { drt.alive = alive }
func (drt *DataResourceTable) IsEmpty() bool                 { return nil == drt }
func (drt *DataResourceTable) IsNotEmpty() bool              { return !drt.IsEmpty() }
func (drt *DataResourceTable) UseDisk(use uint64) {
	if drt.RemainDisk() > use {
		drt.usedDisk += use
	} else {
		drt.usedDisk = drt.totalDisk
	}
	if !drt.isUsed {
		drt.isUsed = true
	}
}
func (drt *DataResourceTable) FreeDisk(use uint64) {
	if drt.usedDisk > use {
		drt.usedDisk -= use
	} else {
		drt.usedDisk = 0
	}
	if drt.usedDisk == 0 {
		drt.isUsed = false
	}
}

type DataResourceFileUpload struct {
	originId   string // db key
	nodeId     string
	metadataId string
	filePath   string
}

type dataResourceFileUploadRlp struct {
	NodeId     string
	OriginId   string
	MetadataId string
	FilePath   string
}

func NewDataResourceFileUpload(nodeId, originId, metaDataId, filePath string) *DataResourceFileUpload {
	return &DataResourceFileUpload{
		nodeId:     nodeId,
		originId:   originId,
		metadataId: metaDataId,
		filePath:   filePath,
	}
}

// EncodeRLP implements rlp.Encoder.
func (drt *DataResourceFileUpload) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, dataResourceFileUploadRlp{
		NodeId:     drt.nodeId,
		OriginId:   drt.originId,
		MetadataId: drt.metadataId,
		FilePath:   drt.filePath,
	})
}

// DecodeRLP implements rlp.Decoder.
func (drt *DataResourceFileUpload) DecodeRLP(s *rlp.Stream) error {
	var dec dataResourceFileUploadRlp
	err := s.Decode(&dec)
	if err == nil {
		drt.nodeId, drt.originId, drt.metadataId, drt.filePath = dec.NodeId, dec.OriginId, dec.MetadataId, dec.FilePath
	}
	return err
}
func (drt *DataResourceFileUpload) GetNodeId() string               { return drt.nodeId }
func (drt *DataResourceFileUpload) GetOriginId() string             { return drt.originId }
func (drt *DataResourceFileUpload) SetMetadataId(metaDataId string) { drt.metadataId = metaDataId }
func (drt *DataResourceFileUpload) GetMetadataId() string           { return drt.metadataId }
func (drt *DataResourceFileUpload) GetFilePath() string             { return drt.filePath }

type DataResourceDiskUsed struct {
	metadataId string // db key
	nodeId     string
	diskUsed   uint64
}

type dataResourceDiskUsedRlp struct {
	MetadataId string // db key
	NodeId     string
	DiskUsed   uint64
}

func NewDataResourceDiskUsed(metaDataId, nodeId string, diskUsed uint64) *DataResourceDiskUsed {
	return &DataResourceDiskUsed{
		metadataId: metaDataId,
		nodeId:     nodeId,
		diskUsed:   diskUsed,
	}
}

// EncodeRLP implements rlp.Encoder.
func (drt *DataResourceDiskUsed) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, dataResourceDiskUsedRlp{
		MetadataId: drt.metadataId,
		NodeId:     drt.nodeId,
		DiskUsed:   drt.diskUsed,
	})
}

// DecodeRLP implements rlp.Decoder.
func (drt *DataResourceDiskUsed) DecodeRLP(s *rlp.Stream) error {
	var dec dataResourceDiskUsedRlp
	err := s.Decode(&dec)
	if err == nil {
		drt.metadataId, drt.nodeId, drt.diskUsed = dec.MetadataId, dec.NodeId, dec.DiskUsed
	}
	return err
}
func (drt *DataResourceDiskUsed) GetMetadataId() string { return drt.metadataId }
func (drt *DataResourceDiskUsed) GetNodeId() string     { return drt.nodeId }
func (drt *DataResourceDiskUsed) GetDiskUsed() uint64   { return drt.diskUsed }

// v 2.0

type TaskUpResultFile struct {
	taskId     string
	originId   string
	metadataId string
}

type taskUpResultFileRlp struct {
	TaskId     string
	OriginId   string
	MetadataId string
}

func NewTaskUpResultFile(taskId, originId, metadataId string) *TaskUpResultFile {
	return &TaskUpResultFile{
		taskId:     taskId,
		originId:   originId,
		metadataId: metadataId,
	}
}

// EncodeRLP implements rlp.Encoder.
func (turf *TaskUpResultFile) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, taskUpResultFileRlp{
		TaskId:     turf.taskId,
		OriginId:   turf.originId,
		MetadataId: turf.metadataId,
	})
}

// DecodeRLP implements rlp.Decoder.
func (turf *TaskUpResultFile) DecodeRLP(s *rlp.Stream) error {
	var dec taskUpResultFileRlp
	err := s.Decode(&dec)
	if err == nil {
		turf.taskId, turf.originId, turf.metadataId = dec.TaskId, dec.OriginId, dec.MetadataId
	}
	return err
}

func (turf *TaskUpResultFile) GetTaskId() string     { return turf.taskId }
func (turf *TaskUpResultFile) GetOriginId() string   { return turf.originId }
func (turf *TaskUpResultFile) GetMetadataId() string { return turf.metadataId }

type TaskResuorceUsage struct {
	taskId         string
	partyId        string
	totalMem       uint64
	totalProcessor uint32
	totalBandwidth uint64
	totalDisk      uint64
	usedMem        uint64
	usedProcessor  uint32
	usedBandwidth  uint64
	usedDisk       uint64
}

type taskResuorceUsageRlp struct {
	TaskId         string
	PartyId        string
	TotalMem       uint64
	TotalProcessor uint32
	TotalBandwidth uint64
	TotalDisk      uint64
	UsedMem        uint64
	UsedProcessor  uint32
	UsedBandwidth  uint64
	UsedDisk       uint64
}

func NewTaskResuorceUsage(
	taskId,
	partyId string,
	totalMem,
	totalBandwidth,
	totalDisk,
	usedMem,
	usedBandwidth,
	usedDisk uint64,
	totalProcessor,
	usedProcessor uint32,
) *TaskResuorceUsage {
	return &TaskResuorceUsage{
		taskId:         taskId,
		partyId:        partyId,
		totalMem:       totalMem,
		totalProcessor: totalProcessor,
		totalBandwidth: totalBandwidth,
		totalDisk:      totalDisk,
		usedMem:        usedMem,
		usedProcessor:  usedProcessor,
		usedBandwidth:  usedBandwidth,
		usedDisk:       usedDisk,
	}
}

// EncodeRLP implements rlp.Encoder.
func (tru *TaskResuorceUsage) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, taskResuorceUsageRlp{
		TaskId:         tru.taskId,
		PartyId:        tru.partyId,
		TotalMem:       tru.totalMem,
		TotalProcessor: tru.totalProcessor,
		TotalBandwidth: tru.totalBandwidth,
		TotalDisk:      tru.totalDisk,
		UsedMem:        tru.usedMem,
		UsedProcessor:  tru.usedProcessor,
		UsedBandwidth:  tru.usedBandwidth,
		UsedDisk:       tru.usedDisk,
	})
}

// DecodeRLP implements rlp.Decoder.
func (tru *TaskResuorceUsage) DecodeRLP(s *rlp.Stream) error {
	var dec taskResuorceUsageRlp
	err := s.Decode(&dec)
	if err == nil {
		tru.taskId, tru.partyId, tru.totalMem, tru.totalBandwidth, tru.totalDisk, tru.totalProcessor, tru.usedMem, tru.usedBandwidth, tru.usedDisk, tru.usedProcessor =
			dec.TaskId, dec.PartyId, dec.TotalMem, dec.TotalBandwidth, dec.TotalDisk, dec.TotalProcessor, dec.UsedMem, dec.UsedBandwidth, dec.UsedDisk, dec.UsedProcessor
	}
	return err
}

func (tru *TaskResuorceUsage) String() string {
	return fmt.Sprintf(`{"taskId": %s, "partyId": %s, "totalMem": %d, "totalBandwidth": %d, "totalDisk": %d, "totalProcessor": %d, "usedMem": %d, "usedBandwidth": %d, "usedDisk": %d, "usedProcessor": %d}`,
		tru.taskId, tru.partyId, tru.totalMem, tru.totalBandwidth, tru.totalDisk, tru.totalProcessor, tru.usedMem, tru.usedBandwidth, tru.usedDisk, tru.usedProcessor)
}

func (tru *TaskResuorceUsage) GetTaskId() string         { return tru.taskId }
func (tru *TaskResuorceUsage) GetPartyId() string        { return tru.partyId }
func (tru *TaskResuorceUsage) GetTotalMem() uint64       { return tru.totalMem }
func (tru *TaskResuorceUsage) GetTotalProcessor() uint32 { return tru.totalProcessor }
func (tru *TaskResuorceUsage) GetTotalBandwidth() uint64 { return tru.totalBandwidth }
func (tru *TaskResuorceUsage) GetTotalDisk() uint64      { return tru.totalDisk }
func (tru *TaskResuorceUsage) GetUsedMem() uint64        { return tru.usedMem }
func (tru *TaskResuorceUsage) GetUsedProcessor() uint32  { return tru.usedProcessor }
func (tru *TaskResuorceUsage) GetUsedBandwidth() uint64  { return tru.usedBandwidth }
func (tru *TaskResuorceUsage) GetUsedDisk() uint64       { return tru.usedDisk }

func (tru *TaskResuorceUsage) SetUsedMem(usedMem uint64) { tru.usedMem = usedMem }
func (tru *TaskResuorceUsage) SetUsedProcessor(usedProcessor uint32) {
	tru.usedProcessor = usedProcessor
}
func (tru *TaskResuorceUsage) SetUsedBandwidth(usedBandwidth uint64) {
	tru.usedBandwidth = usedBandwidth
}
func (tru *TaskResuorceUsage) SetUsedDisk(usedDisk uint64) { tru.usedDisk = usedDisk }
