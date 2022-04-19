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

	ErrMetadataMsgConvert       = errors.New("convert metadata msg failed")
	ErrMetadataRevokeMsgConvert = errors.New("convert metadata revokeMsg failed")

	ErrMetadataAuthMsgConvert       = errors.New("convert metadata authority msg failed")
	ErrMetadataAuthRevokeMsgConvert = errors.New("convert metadata authority revokeMsg failed")

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
}

func NewMempool(cfg *MempoolConfig) *Mempool {
	return &Mempool{
		cfg: cfg,
	}
}

// SubscribeNewTxsEvent registers a subscription of NewTxsEvent and
// starts sending evengine to the given channel.

func (pool *Mempool) SubscribeNewMessageEvent(ch chan<- *feed.Event) event.Subscription {
	return pool.scope.Track(pool.msgFeed.Subscribe(ch))
}

func (pool *Mempool) Add(msg types.Msg) error {

	switch m := msg.(type) {
	case *types.IdentityMsg:

		// set local nodeId first
		m.SetOwnerNodeId(pool.cfg.NodeId)
		// We've directly injected a replacement identityMsg, notify subsystems
		pool.msgFeed.Send(&feed.Event{
			Type: types.ApplyIdentity,
			Data: &types.IdentityMsgEvent{Msg: m},
		})

	case *types.IdentityRevokeMsg:

		// We've directly injected a replacement identityMsg, notify subsystems
		pool.msgFeed.Send(&feed.Event{
			Type: types.RevokeIdentity,
			Data: &types.IdentityRevokeMsgEvent{Msg: m},
		})

	case *types.PowerMsg:

		// We've directly injected a replacement identityRevokeMsg, notify subsystems
		pool.msgFeed.Send(&feed.Event{
			Type: types.ApplyPower,
			Data: &types.PowerMsgEvent{Msg: m},
		})

	case *types.PowerRevokeMsg:

		// We've directly injected a replacement powerRevokeMsg, notify subsystems
		pool.msgFeed.Send(&feed.Event{
			Type: types.RevokePower,
			Data: &types.PowerRevokeMsgEvent{Msg: m},
		})

	case *types.MetadataMsg:

		// We've directly injected a replacement metaDataMsg, notify subsystems
		pool.msgFeed.Send(&feed.Event{
			Type: types.ApplyMetadata,
			Data: &types.MetadataMsgEvent{Msg: m},
		})

	case *types.MetadataRevokeMsg:

		// We've directly injected a replacement metaDataRevokeMsg, notify subsystems
		pool.msgFeed.Send(&feed.Event{
			Type: types.RevokeMetadata,
			Data: &types.MetadataRevokeMsgEvent{Msg: m},
		})

	case *types.MetadataAuthorityMsg:

		// We've directly injected a replacement metadata authority msg, notify subsystems
		pool.msgFeed.Send(&feed.Event{
			Type: types.ApplyMetadataAuth,
			Data: &types.MetadataAuthMsgEvent{Msg: m},
		})

	case *types.MetadataAuthorityRevokeMsg:

		// We've directly injected a replacement metadata authority rovkeMsg, notify subsystems
		pool.msgFeed.Send(&feed.Event{
			Type: types.RevokeMetadataAuth,
			Data: &types.MetadataAuthRevokeMsgEvent{Msg: m},
		})

	case *types.TaskMsg:

		// We've directly injected a replacement taskMsg, notify subsystems
		pool.msgFeed.Send(&feed.Event{
			Type: types.ApplyTask,
			Data: &types.TaskMsgEvent{Msg: m},
		})

	case *types.TaskTerminateMsg:

		// We've directly injected a replacement taskTerminate msg, notify subsystems
		pool.msgFeed.Send(&feed.Event{
			Type: types.TerminateTask,
			Data: &types.TaskTerminateMsgEvent{Msg: m},
		})

	default:
		log.Errorf("Failed to add msg, can not match the msg type")
		return ErrUnknownMsgType
	}
	return nil
}

type msgLookup struct {

	// metaDataId -> Msg
	allMateDataMsg map[string]*types.MetadataMsg
	// powerId -> Msg
	allPowerMsg map[string]*types.PowerMsg
	//allTaskMsg     map[string]*types.TaskMsg

	metaDataMsgLock sync.RWMutex
	powerMsgLock    sync.RWMutex

}

func newMsgLookup() *msgLookup {
	return &msgLookup{
		allMateDataMsg: make(map[string]*types.MetadataMsg),
		allPowerMsg:    make(map[string]*types.PowerMsg),
		//allTaskMsg:     make(map[string]*types.TaskMsg),
	}
}

// RangeMetadataMsg calls f on each key and value present in the map.
func (lookup *msgLookup) rangeMetadataMsg(f func(metaDataId string, msg *types.MetadataMsg) bool) {
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
func (lookup *msgLookup) getMetadataMsg(metaDataId string) *types.MetadataMsg {
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
func (lookup *msgLookup) addMetadataMsg(msg *types.MetadataMsg) {
	lookup.metaDataMsgLock.RLock()
	defer lookup.metaDataMsgLock.RUnlock()

	lookup.allMateDataMsg[msg.GetMetadataId()] = msg
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
//	lookup.allTaskMsg[msg.GetTaskId] = msg
//}

// Remove removes a metaDataMsg from the lookup.
func (lookup *msgLookup) removeMetadataMsg(metaDataId string) {
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
func (lookup *msgLookup) extractMetadataMsg(metaDataId string) (*types.MetadataMsg, bool) {
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
