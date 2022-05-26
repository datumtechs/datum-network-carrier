package message

import (
	"container/heap"
	"github.com/datumtechs/datum-network-carrier/types"
	"sort"
)

type timeHeap []uint64

func (h timeHeap) Len() int           { return len(h) }
func (h timeHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h timeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *timeHeap) Push(x interface{}) {
	*h = append(*h, x.(uint64))
}

func (h *timeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type PowerMsgList struct {
	items  map[uint64]*types.PowerMsg
	prioty *timeHeap
	cache  types.PowerMsgArr // Cache of the powerMsgs already sorted
}

func newPowerMsgList() *PowerMsgList {
	return &PowerMsgList{
		items:  make(map[uint64]*types.PowerMsg),
		prioty: new(timeHeap),
	}
}
func (lis *PowerMsgList) put(msg *types.PowerMsg) {
	heap.Push(lis.prioty, msg.CreateAt)
	lis.items[msg.CreateAt], lis.cache = msg, nil
}

func (lis *PowerMsgList) get(createAt uint64) *types.PowerMsg {
	return lis.items[createAt]
}

func (lis *PowerMsgList) reheap() {
	*lis.prioty = make([]uint64, 0, len(lis.items))
	for time := range lis.items {
		*lis.prioty = append(*lis.prioty, time)
	}
	heap.Init(lis.prioty)
	lis.cache = nil
}
func (lis *PowerMsgList) flatten() types.PowerMsgArr {
	// If the sorting was not cached yet, create and cache it
	if lis.cache == nil {
		lis.cache = make(types.PowerMsgArr, 0, len(lis.items))
		for _, msg := range lis.items {
			lis.cache = append(lis.cache, msg)
		}
		sort.Sort(lis.cache)
	}
	return lis.cache
}

// Flatten returns the PowerMsgArr sequence sorted by timestamp size order
func (lis *PowerMsgList) Flatten() types.PowerMsgArr {
	// Copy the cache to prevent accidental modifications
	cache := lis.flatten()
	txs := make(types.PowerMsgArr, len(cache))
	copy(txs, cache)
	return txs
}

// LastElement returns the last element of the flat list
func (lis *PowerMsgList) LastElement() *types.PowerMsg {
	cache := lis.flatten()
	return cache[len(cache)-1]
}

// FirstElement returns the first element of the flat list
func (lis *PowerMsgList) FirstElement() *types.PowerMsg {
	cache := lis.flatten()
	return cache[len(cache)-1]
}

type MetadataMsgList struct {
	items  map[uint64]*types.MetadataMsg
	prioty *timeHeap
	cache  types.MetadataMsgArr // Cache of the metaDataMsgs already sorted
}

func newMetadataMsgList() *MetadataMsgList {
	return &MetadataMsgList{
		items:  make(map[uint64]*types.MetadataMsg),
		prioty: new(timeHeap),
	}
}

func (lis *MetadataMsgList) put(msg *types.MetadataMsg) {
	heap.Push(lis.prioty, msg.GetCreateAt())
	lis.items[msg.GetCreateAt()] = msg
}

func (lis *MetadataMsgList) get(createAt uint64) *types.MetadataMsg {
	return lis.items[createAt]
}

func (lis *MetadataMsgList) reheap() {
	*lis.prioty = make([]uint64, 0, len(lis.items))
	for time := range lis.items {
		*lis.prioty = append(*lis.prioty, time)
	}
	heap.Init(lis.prioty)
	lis.cache = nil
}
func (lis *MetadataMsgList) flatten() types.MetadataMsgArr {
	// If the sorting was not cached yet, create and cache it
	if lis.cache == nil {
		lis.cache = make(types.MetadataMsgArr, 0, len(lis.items))
		for _, msg := range lis.items {
			lis.cache = append(lis.cache, msg)
		}
		sort.Sort(lis.cache)
	}
	return lis.cache
}

// Flatten returns the MetadataMsgArr sequence sorted by timestamp size order
func (lis *MetadataMsgList) Flatten() types.MetadataMsgArr {
	// Copy the cache to prevent accidental modifications
	cache := lis.flatten()
	txs := make(types.MetadataMsgArr, len(cache))
	copy(txs, cache)
	return txs
}

// LastElement returns the last element of the flat list
func (lis *MetadataMsgList) LastElement() *types.MetadataMsg {
	cache := lis.flatten()
	return cache[len(cache)-1]
}

// FirstElement returns the first element of the flat list
func (lis *MetadataMsgList) FirstElement() *types.MetadataMsg {
	cache := lis.flatten()
	return cache[len(cache)-1]
}

type TaskMsgList struct {
	items  map[uint64]*types.TaskMsg
	prioty *timeHeap
	cache  types.TaskMsgArr // Cache of the taskMsgs already sorted
}

func newTaskMsgList() *TaskMsgList {
	return &TaskMsgList{
		items:  make(map[uint64]*types.TaskMsg),
		prioty: new(timeHeap),
	}
}

func (lis *TaskMsgList) put(msg *types.TaskMsg) {
	heap.Push(lis.prioty, msg.GetCreateAt())
	lis.items[msg.GetCreateAt()] = msg
}

func (lis *TaskMsgList) get(createAt uint64) *types.TaskMsg {
	return lis.items[createAt]
}

func (lis *TaskMsgList) reheap() {
	*lis.prioty = make([]uint64, 0, len(lis.items))
	for time := range lis.items {
		*lis.prioty = append(*lis.prioty, time)
	}
	heap.Init(lis.prioty)
	lis.cache = nil
}
func (lis *TaskMsgList) flatten() types.TaskMsgArr {
	// If the sorting was not cached yet, create and cache it
	if lis.cache == nil {
		lis.cache = make(types.TaskMsgArr, 0, len(lis.items))
		for _, msg := range lis.items {
			lis.cache = append(lis.cache, msg)
		}
		sort.Sort(lis.cache)
	}
	return lis.cache
}

// Flatten returns the TaskMsgArr sequence sorted by timestamp size order
func (lis *TaskMsgList) Flatten() types.TaskMsgArr {
	// Copy the cache to prevent accidental modifications
	cache := lis.flatten()
	txs := make(types.TaskMsgArr, len(cache))
	copy(txs, cache)
	return txs
}

// LastElement returns the last element of the flat list
func (lis *TaskMsgList) LastElement() *types.TaskMsg {
	cache := lis.flatten()
	return cache[len(cache)-1]
}

// FirstElement returns the first element of the flat list
func (lis *TaskMsgList) FirstElement() *types.TaskMsg {
	cache := lis.flatten()
	return cache[len(cache)-1]
}
