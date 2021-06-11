package message

import (
	"container/heap"
	"github.com/RosettaFlow/Carrier-Go/types"
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
}

func NewPowerMsgList() *PowerMsgList {
	return &PowerMsgList{
		items: 	make(map[uint64]*types.PowerMsg),
		prioty:	new(timeHeap),
	}
}

func (lis *PowerMsgList) Put (msg *types.PowerMsg) {
	heap.Push(lis.prioty, msg.CreateAt())
	lis.items[msg.CreateAt()] = msg
}

func (lis *PowerMsgList) Get (createAt uint64) *types.PowerMsg {
	return lis.items[createAt]
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