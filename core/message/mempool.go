package message

import (
	"errors"
	"github.com/RosettaFlow/Carrier-Go/common/feed"
	"github.com/RosettaFlow/Carrier-Go/event"
	"github.com/RosettaFlow/Carrier-Go/types"
	"sync"
)

var (
	ErrIdentityMsgConvert       = errors.New("convert identity msg failed")
	ErrIdentityRevokeMsgConvert = errors.New("convert identity revokeMsg failed")

	ErrPowerMsgConvert       = errors.New("convert power msg failed")
	ErrPowerRevokeMsgConvert = errors.New("convert power revokeMsg failed")

	ErrMetaDataMsgConvert       = errors.New("convert metadata msg failed")
	ErrMetaDataRevokeMsgConvert = errors.New("convert metadata revokeMsg failed")

	ErrTaskMsgConvert = errors.New("convert task msg failed")

	ErrUnknownMsgType = errors.New("Unknown msg type")
)

type MempoolConfig struct {
	NodeId string
}

type Mempool struct {
	cfg *MempoolConfig

	msgFeed event.Feed
	scope   event.SubscriptionScope

	//metaDataMsgQueue *MetaDataMsgList
	//powerMsgQueue    *PowerMsgList
	//taskMsgQueue     *TaskMsgList
	//all *msgLookup
}

func NewMempool(cfg *MempoolConfig) *Mempool {
	return &Mempool{
		cfg: cfg,
		//all:               newMsgLookup(),
		//metaDataMsgQueue: newMetaDataMsgList(),
		//powerMsgQueue:    newPowerMsgList(),
		//taskMsgQueue:     newTaskMsgList(),
	}
}

// SubscribeNewTxsEvent registers a subscription of NewTxsEvent and
// starts sending evengine to the given channel.

func (pool *Mempool) SubscribeNewMessageEvent(ch chan<- *feed.Event) event.Subscription {
	return pool.scope.Track(pool.msgFeed.Subscribe(ch))
}

func (pool *Mempool) Add(msg types.Msg) error {

	switch msg.(type) {
	case *types.IdentityMsg:
		identity, ok := msg.(*types.IdentityMsg)
		if !ok {
			return ErrIdentityMsgConvert
		}
		// 先设置 本地节点的 nodeId
		identity.NodeId = pool.cfg.NodeId
		// We've directly injected a replacement identityMsg, notify subsystems
		pool.msgFeed.Send(&feed.Event{
			Type: types.ApplyIdentity,
			Data: &types.IdentityMsgEvent{Msg: identity},
		})

	case *types.IdentityRevokeMsg:
		identityRevoke, ok := msg.(*types.IdentityRevokeMsg)
		if !ok {
			return ErrIdentityRevokeMsgConvert
		}

		// We've directly injected a replacement identityMsg, notify subsystems
		pool.msgFeed.Send(&feed.Event{
			Type: types.RevokeIdentity,
			Data: &types.IdentityRevokeMsgEvent{Msg: identityRevoke},
		})

	case *types.PowerMsg:
		power, ok := msg.(*types.PowerMsg)
		if !ok {
			return ErrPowerMsgConvert
		}

		//pool.powerMsgQueue.put(power)
		// We've directly injected a replacement identityRevokeMsg, notify subsystems
		pool.msgFeed.Send(&feed.Event{
			Type: types.ApplyPower,
			Data: &types.PowerMsgEvent{Msgs: types.PowerMsgs{power}},
		})

	case *types.PowerRevokeMsg:
		powerRevoke, ok := msg.(*types.PowerRevokeMsg)
		if !ok {
			return ErrPowerRevokeMsgConvert
		}

		// We've directly injected a replacement powerRevokeMsg, notify subsystems
		pool.msgFeed.Send(&feed.Event{
			Type: types.RevokePower,
			Data: &types.PowerRevokeMsgEvent{Msgs: types.PowerRevokeMsgs{powerRevoke}},
		})

	case *types.MetaDataMsg:
		metaData, ok := msg.(*types.MetaDataMsg)
		if !ok {
			return ErrMetaDataMsgConvert
		}

		//pool.metaDataMsgQueue.put(metaData)
		// We've directly injected a replacement metaDataMsg, notify subsystems
		pool.msgFeed.Send(&feed.Event{
			Type: types.ApplyMetadata,
			Data: &types.MetaDataMsgEvent{Msgs: types.MetaDataMsgs{metaData}},
		})

	case *types.MetaDataRevokeMsg:
		metaDataRevoke, ok := msg.(*types.MetaDataRevokeMsg)
		if !ok {
			return ErrMetaDataRevokeMsgConvert
		}

		// We've directly injected a replacement metaDataRevokeMsg, notify subsystems
		pool.msgFeed.Send(&feed.Event{
			Type: types.RevokeMetadata,
			Data: &types.MetaDataRevokeMsgEvent{Msgs: types.MetaDataRevokeMsgs{metaDataRevoke}},
		})

	case *types.TaskMsg:
		task, ok := msg.(*types.TaskMsg)
		if !ok {
			return ErrTaskMsgConvert
		}

		// We've directly injected a replacement taskMsg, notify subsystems
		pool.msgFeed.Send(&feed.Event{
			Type: types.ApplyTask,
			Data: &types.TaskMsgEvent{Msgs: types.TaskMsgs{task}},
		})

	default:
		log.Fatalf("Failed to add msg, can not match the msg type")
		return ErrUnknownMsgType
	}
	return nil
}

type msgLookup struct {

	// metaDataId -> Msg
	allMateDataMsg map[string]*types.MetaDataMsg
	// powerId -> Msg
	allPowerMsg map[string]*types.PowerMsg
	//allTaskMsg     map[string]*types.TaskMsg

	metaDataMsgLock sync.RWMutex
	powerMsgLock    sync.RWMutex
	//taskMsgLock     sync.RWMutex
}

func newMsgLookup() *msgLookup {
	return &msgLookup{
		allMateDataMsg: make(map[string]*types.MetaDataMsg),
		allPowerMsg:    make(map[string]*types.PowerMsg),
		//allTaskMsg:     make(map[string]*types.TaskMsg),
	}
}

// RangeMetaDataMsg calls f on each key and value present in the map.
func (lookup *msgLookup) rangeMetaDataMsg(f func(metaDataId string, msg *types.MetaDataMsg) bool) {
	lookup.metaDataMsgLock.RLock()
	defer lookup.metaDataMsgLock.RUnlock()

	for key, value := range lookup.allMateDataMsg {
		if !f(key, value) {
			break
		}
	}
}

// RangePowerMsg calls f on each key and value present in the map.
func (lookup *msgLookup) rangePowerMsg(f func(powerId string, msg *types.PowerMsg) bool) {
	lookup.powerMsgLock.RLock()
	defer lookup.powerMsgLock.RUnlock()

	for key, value := range lookup.allPowerMsg {
		if !f(key, value) {
			break
		}
	}
}

//// RangeTaskMsg calls f on each key and value present in the map.
//func (lookup *msgLookup) rangeTaskMsg(f func(taskId string, msg *types.TaskMsg) bool) {
//	lookup.taskMsgLock.RLock()
//	defer lookup.taskMsgLock.RUnlock()
//
//	for key, value := range lookup.allTaskMsg {
//		if !f(key, value) {
//			break
//		}
//	}
//}

// Get returns a metaDataMsg if it exists in the lookup, or nil if not found.
func (lookup *msgLookup) getMetaDataMsg(metaDataId string) *types.MetaDataMsg {
	lookup.metaDataMsgLock.RLock()
	defer lookup.metaDataMsgLock.RUnlock()

	return lookup.allMateDataMsg[metaDataId]
}

// Get returns a powerMsg if it exists in the lookup, or nil if not found.
func (lookup *msgLookup) getPowerMsg(powerId string) *types.PowerMsg {
	lookup.powerMsgLock.RLock()
	defer lookup.powerMsgLock.RUnlock()

	return lookup.allPowerMsg[powerId]
}

//// Get returns a taskMsg if it exists in the lookup, or nil if not found.
//func (lookup *msgLookup) getTaskMsg(taskId string) *types.TaskMsg {
//	lookup.taskMsgLock.RLock()
//	defer lookup.taskMsgLock.RUnlock()
//
//	return lookup.allTaskMsg[taskId]
//}

// Count returns the current number of items in the lookup.
func (lookup *msgLookup) metaDataMsgCount() int {
	lookup.metaDataMsgLock.RLock()
	defer lookup.metaDataMsgLock.RUnlock()

	return len(lookup.allMateDataMsg)
}

// Count returns the current number of items in the lookup.
func (lookup *msgLookup) powerMsgCount() int {
	lookup.powerMsgLock.RLock()
	defer lookup.powerMsgLock.RUnlock()

	return len(lookup.allPowerMsg)
}

//// Count returns the current number of items in the lookup.
//func (lookup *msgLookup) taskMsgCount() int {
//	lookup.taskMsgLock.RLock()
//	defer lookup.taskMsgLock.RUnlock()
//
//	return len(lookup.allTaskMsg)
//}

// Add adds a metaDataMsg to the lookup.
func (lookup *msgLookup) addMetaDataMsg(msg *types.MetaDataMsg) {
	lookup.metaDataMsgLock.RLock()
	defer lookup.metaDataMsgLock.RUnlock()

	lookup.allMateDataMsg[msg.MetaDataId] = msg
}

// Add adds a powerMsg to the lookup.
func (lookup *msgLookup) addPowerMsg(msg *types.PowerMsg) {
	lookup.powerMsgLock.RLock()
	defer lookup.powerMsgLock.RUnlock()

	lookup.allPowerMsg[msg.PowerId] = msg
}

//// Add adds a taskMsg to the lookup.
//func (lookup *msgLookup) addTaskMsg(msg *types.TaskMsg) {
//	lookup.taskMsgLock.RLock()
//	defer lookup.taskMsgLock.RUnlock()
//
//	lookup.allTaskMsg[msg.TaskId] = msg
//}

// Remove removes a metaDataMsg from the lookup.
func (lookup *msgLookup) removeMetaDataMsg(metaDataId string) {
	lookup.metaDataMsgLock.RLock()
	delete(lookup.allMateDataMsg, metaDataId)
	lookup.metaDataMsgLock.RUnlock()
}

// Remove removes a powerMsg from the lookup.
func (lookup *msgLookup) removePowerMsg(powerId string) {
	lookup.powerMsgLock.RLock()
	delete(lookup.allPowerMsg, powerId)
	lookup.powerMsgLock.RUnlock()
}

//// Remove removes a taskMsg from the lookup.
//func (lookup *msgLookup) removeTaskMsg(taskId string) {
//	lookup.taskMsgLock.RLock()
//	delete(lookup.allTaskMsg, taskId)
//	lookup.taskMsgLock.RUnlock()
//}

// Extract removes a metaDataMsg from the lookup, and return.
func (lookup *msgLookup) extractMetaDataMsg(metaDataId string) (*types.MetaDataMsg, bool) {
	lookup.metaDataMsgLock.RLock()
	defer lookup.metaDataMsgLock.RUnlock()
	metaDataMsg, ok := lookup.allMateDataMsg[metaDataId]
	if !ok {
		return nil, false
	}
	delete(lookup.allMateDataMsg, metaDataId)
	return metaDataMsg, true
}

// Extract removes a powerMsg from the lookup, and return.
func (lookup *msgLookup) extractPowerMsg(powerId string) (*types.PowerMsg, bool) {
	lookup.powerMsgLock.RLock()
	defer lookup.powerMsgLock.RUnlock()
	powerMsg, ok := lookup.allPowerMsg[powerId]
	if !ok {
		return nil, false
	}
	delete(lookup.allPowerMsg, powerId)
	return powerMsg, true
}

//// Extract removes a taskMsg from the lookup, and return.
//func (lookup *msgLookup) extractTaskMsg(taskId string) (*types.TaskMsg, bool) {
//	lookup.taskMsgLock.RLock()
//	defer lookup.taskMsgLock.RUnlock()
//	taskMsg, ok := lookup.allTaskMsg[taskId]
//	if !ok {
//		return nil, false
//	}
//	delete(lookup.allTaskMsg, taskId)
//	return taskMsg, true
//}
