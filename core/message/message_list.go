package message

import (
	"container/heap"
	"github.com/RosettaFlow/Carrier-Go/types"
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
	all    *msgLookup
	items  map[uint64]*types.PowerMsg
	prioty *timeHeap
	cache  types.PowerMsgs // Cache of the powerMsgs already sorted
}

func newPowerMsgList(all *msgLookup) *PowerMsgList {
	return &PowerMsgList{
		all:    all,
		items:  make(map[uint64]*types.PowerMsg),
		prioty: new(timeHeap),
	}
}
func (lis *PowerMsgList) put(msg *types.PowerMsg) {
	heap.Push(lis.prioty, msg.CreateAt())
	lis.items[msg.CreateAt()], lis.cache = msg, nil
}

func (lis *PowerMsgList) get(createAt uint64) *types.PowerMsg {
	return lis.items[createAt]
}
func (lis *PowerMsgList) popBy(createAt uint64) bool {
	power, ok := lis.items[createAt]
	if !ok {
		return false
	}

	// Otherwise delete the msg and fix the heap prioty
	for i := 0; i < lis.prioty.Len(); i++ {
		if (*lis.prioty)[i] == createAt {
			heap.Remove(lis.prioty, i)
			break
		}
	}
	delete(lis.items, createAt)
	lis.all.removePowerMsg(power.PowerId)
	return true
}
func (lis *PowerMsgList) reheap() {
	*lis.prioty = make([]uint64, 0, len(lis.items))
	for time := range lis.items {
		*lis.prioty = append(*lis.prioty, time)
	}
	heap.Init(lis.prioty)
	lis.cache = nil
}
func (lis *PowerMsgList) flatten() types.PowerMsgs {
	// If the sorting was not cached yet, create and cache it
	if lis.cache == nil {
		lis.cache = make(types.PowerMsgs, 0, len(lis.items))
		for _, msg := range lis.items {
			lis.cache = append(lis.cache, msg)
		}
		sort.Sort(lis.cache)
	}
	return lis.cache
}

// Flatten returns the PowerMsgs sequence sorted by timestamp size order
func (lis *PowerMsgList) Flatten() types.PowerMsgs {
	// Copy the cache to prevent accidental modifications
	cache := lis.flatten()
	txs := make(types.PowerMsgs, len(cache))
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

type MetaDataMsgList struct {
	all    *msgLookup
	items  map[uint64]*types.MetaDataMsg
	prioty *timeHeap
	cache  types.MetaDataMsgs // Cache of the metaDataMsgs already sorted
}

func newMetaDataMsgList(all *msgLookup) *MetaDataMsgList {
	return &MetaDataMsgList{
		all:    all,
		items:  make(map[uint64]*types.MetaDataMsg),
		prioty: new(timeHeap),
	}
}

func (lis *MetaDataMsgList) put(msg *types.MetaDataMsg) {
	heap.Push(lis.prioty, msg.CreateAt())
	lis.items[msg.CreateAt()] = msg
}

func (lis *MetaDataMsgList) get(createAt uint64) *types.MetaDataMsg {
	return lis.items[createAt]
}
func (lis *MetaDataMsgList) popBy(createAt uint64) bool {
	metaData, ok := lis.items[createAt]
	if !ok {
		return false
	}

	// Otherwise delete the msg and fix the heap prioty
	for i := 0; i < lis.prioty.Len(); i++ {
		if (*lis.prioty)[i] == createAt {
			heap.Remove(lis.prioty, i)
			break
		}
	}
	delete(lis.items, createAt)
	lis.all.removeMetaDataMsg(metaData.MetaDataId)
	return true
}
func (lis *MetaDataMsgList) reheap() {
	*lis.prioty = make([]uint64, 0, len(lis.items))
	for time := range lis.items {
		*lis.prioty = append(*lis.prioty, time)
	}
	heap.Init(lis.prioty)
	lis.cache = nil
}
func (lis *MetaDataMsgList) flatten() types.MetaDataMsgs {
	// If the sorting was not cached yet, create and cache it
	if lis.cache == nil {
		lis.cache = make(types.MetaDataMsgs, 0, len(lis.items))
		for _, msg := range lis.items {
			lis.cache = append(lis.cache, msg)
		}
		sort.Sort(lis.cache)
	}
	return lis.cache
}

// Flatten returns the MetaDataMsgs sequence sorted by timestamp size order
func (lis *MetaDataMsgList) Flatten() types.MetaDataMsgs {
	// Copy the cache to prevent accidental modifications
	cache := lis.flatten()
	txs := make(types.MetaDataMsgs, len(cache))
	copy(txs, cache)
	return txs
}

// LastElement returns the last element of the flat list
func (lis *MetaDataMsgList) LastElement() *types.MetaDataMsg {
	cache := lis.flatten()
	return cache[len(cache)-1]
}

// FirstElement returns the first element of the flat list
func (lis *MetaDataMsgList) FirstElement() *types.MetaDataMsg {
	cache := lis.flatten()
	return cache[len(cache)-1]
}

type TaskMsgList struct {
	all    *msgLookup
	items  map[uint64]*types.TaskMsg
	prioty *timeHeap
	cache  types.TaskMsgs // Cache of the taskMsgs already sorted
}

func newTaskMsgList(all *msgLookup) *TaskMsgList {
	return &TaskMsgList{
		all:    all,
		items:  make(map[uint64]*types.TaskMsg),
		prioty: new(timeHeap),
	}
}

func (lis *TaskMsgList) put(msg *types.TaskMsg) {
	heap.Push(lis.prioty, msg.CreateAt())
	lis.items[msg.CreateAt()] = msg
}

func (lis *TaskMsgList) get(createAt uint64) *types.TaskMsg {
	return lis.items[createAt]
}
func (lis *TaskMsgList) popBy(createAt uint64) bool {
	task, ok := lis.items[createAt]
	if !ok {
		return false
	}

	// Otherwise delete the msg and fix the heap prioty
	for i := 0; i < lis.prioty.Len(); i++ {
		if (*lis.prioty)[i] == createAt {
			heap.Remove(lis.prioty, i)
			break
		}
	}
	delete(lis.items, createAt)
	lis.all.removeTaskMsg(task.TaskId)
	return true
}
func (lis *TaskMsgList) reheap() {
	*lis.prioty = make([]uint64, 0, len(lis.items))
	for time := range lis.items {
		*lis.prioty = append(*lis.prioty, time)
	}
	heap.Init(lis.prioty)
	lis.cache = nil
}
func (lis *TaskMsgList) flatten() types.TaskMsgs {
	// If the sorting was not cached yet, create and cache it
	if lis.cache == nil {
		lis.cache = make(types.TaskMsgs, 0, len(lis.items))
		for _, msg := range lis.items {
			lis.cache = append(lis.cache, msg)
		}
		sort.Sort(lis.cache)
	}
	return lis.cache
}

// Flatten returns the TaskMsgs sequence sorted by timestamp size order
func (lis *TaskMsgList) Flatten() types.TaskMsgs {
	// Copy the cache to prevent accidental modifications
	cache := lis.flatten()
	txs := make(types.TaskMsgs, len(cache))
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
