package message

import (
	"container/heap"
	"github.com/RosettaFlow/Carrier-Go/types"
	"sort"
)

type  timeHeap  []uint64

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
	items 			map[uint64]*types.PowerMsg
	prioty			*timeHeap
	cache 			types.PowerMsgs            // Cache of the powerMsgs already sorted
}
func NewPowerMsgList() *PowerMsgList {
	return &PowerMsgList{
		items: 	make(map[uint64]*types.PowerMsg),
		prioty:	new(timeHeap),
	}
}
func (lis *PowerMsgList) Put (msg *types.PowerMsg) {
	heap.Push(lis.prioty, msg.CreateAt())
	lis.items[msg.CreateAt()], lis.cache = msg, nil
}

func (lis *PowerMsgList) Get (createAt uint64) *types.PowerMsg {
	return lis.items[createAt]
}
func (lis *PowerMsgList) Pop (createAt uint64) bool {
	_, ok := lis.items[createAt]
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
func  (lis *PowerMsgList) FirstElement() *types.PowerMsg {
	cache := lis.flatten()
	return cache[len(cache)-1]
}

func (lis *PowerMsgList) arrange () {

}




type MetaDataMsgList struct {
	items 			map[uint64]*types.MetaDataMsg
	prioty			*timeHeap
}

func NewMetaDataMsgList() *MetaDataMsgList {
	return &MetaDataMsgList{
		items: 	make(map[uint64]*types.MetaDataMsg),
		prioty:	new(timeHeap),
	}
}

func (lis *MetaDataMsgList) Put (msg *types.MetaDataMsg) {
	heap.Push(lis.prioty, msg.CreateAt())
	lis.items[msg.CreateAt()] = msg
}

func (lis *MetaDataMsgList) Get (createAt uint64) *types.MetaDataMsg {
	return lis.items[createAt]
}

type TaskMsgList struct {
	items 			map[uint64]*types.TaskMsg
	prioty			*timeHeap
}

func NewTaskMsgList() *TaskMsgList {
	return &TaskMsgList{
		items: 	make(map[uint64]*types.TaskMsg),
		prioty:	new(timeHeap),
	}
}

func (lis *TaskMsgList) Put (msg *types.TaskMsg) {
	heap.Push(lis.prioty, msg.CreateAt())
	lis.items[msg.CreateAt()] = msg
}

func (lis *TaskMsgList) Get (createAt uint64) *types.TaskMsg {
	return lis.items[createAt]
}