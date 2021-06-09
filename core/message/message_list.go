package message

import "github.com/RosettaFlow/Carrier-Go/types"

type  timeHeap  []*uint64

type PowerMsgList struct {
	items 			map[uint64]*types.PowerMsg
	prioty			timeHeap
}

type MetaDataMsgList struct {
	items 			map[uint64]*types.MetaDataMsg
	prioty			timeHeap
}

type TaskMsgList struct {
	items 			map[uint64]*types.TaskMsg
	prioty			timeHeap
}
