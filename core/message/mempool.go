package message

import (
	"github.com/RosettaFlow/Carrier-Go/event"
	"github.com/RosettaFlow/Carrier-Go/types"
	"sync"
)



type MempoolConfig struct {
}







type Mempool struct {
	cfg *MempoolConfig

	msgFeed event.Feed
	scope   event.SubscriptionScope

	metaDataMsgQueue *MetaDataMsgList
	powerMsgQueue    *PowerMsgList
	taskMsgQueue     *TaskMsgList

	all			*msgLookup

	identityErrCh   chan error
}

func NewMempool(cfg *MempoolConfig) *Mempool {
	lookup := newMsgLookup()
	return &Mempool{
		cfg:              cfg,
		all:              lookup,
		metaDataMsgQueue: newMetaDataMsgList(lookup),
		powerMsgQueue:    newPowerMsgList(lookup),
		taskMsgQueue:     newTaskMsgList(lookup),
	}
}

// SubscribeNewTxsEvent registers a subscription of NewTxsEvent and
// starts sending event to the given channel.
func (pool *Mempool) SubscribeNewIdentityMsgsEvent(ch chan<- event.IdentityMsgEvent) event.Subscription {
	return pool.scope.Track(pool.msgFeed.Subscribe(ch))
}
func (pool *Mempool) SubscribeNewIdentityRevokeMsgsEvent(ch chan<- event.IdentityRevokeMsgEvent) event.Subscription {
	return pool.scope.Track(pool.msgFeed.Subscribe(ch))
}
func (pool *Mempool) SubscribeNewMetaDataMsgsEvent(ch chan<- event.MetaDataMsgEvent) event.Subscription {
	return pool.scope.Track(pool.msgFeed.Subscribe(ch))
}
func (pool *Mempool) SubscribeNewMetaDataRevokeMsgsEvent(ch chan<- event.MetaDataRevokeMsgEvent) event.Subscription {
	return pool.scope.Track(pool.msgFeed.Subscribe(ch))
}
func (pool *Mempool) SubscribeNewPowerMsgsEvent(ch chan<- event.PowerMsgEvent) event.Subscription {
	return pool.scope.Track(pool.msgFeed.Subscribe(ch))
}
func (pool *Mempool) SubscribeNewPowerRevokeMsgsEvent(ch chan<- event.PowerRevokeMsgEvent) event.Subscription {
	return pool.scope.Track(pool.msgFeed.Subscribe(ch))
}
func (pool *Mempool) SubscribeNewTaskMsgsEvent(ch chan<- event.TaskMsgEvent) event.Subscription {
	return pool.scope.Track(pool.msgFeed.Subscribe(ch))
}

func (pool *Mempool) Add(msg types.Msg) error {

	switch msg.(type) {
	case *types.IdentityMsg:
		identity, _ := msg.(*types.IdentityMsg)
		// We've directly injected a replacement identityMsg, notify subsystems
		go pool.msgFeed.Send(event.IdentityMsgEvent{identity})

	case *types.IdentityRevokeMsg:
		identityRevoke, _ := msg.(*types.IdentityRevokeMsg)

		// We've directly injected a replacement identityMsg, notify subsystems
		go pool.msgFeed.Send(event.IdentityRevokeMsgEvent{identityRevoke})

	case *types.PowerMsg:
		power, _ := msg.(*types.PowerMsg)
		pool.powerMsgQueue.put(power)
		// We've directly injected a replacement identityRevokeMsg, notify subsystems
		go pool.msgFeed.Send(event.PowerMsgEvent{types.PowerMsgs{power}})

	case *types.PowerRevokeMsg:
		powerRevoke, _ := msg.(*types.PowerRevokeMsg)
		power := pool.all.getPowerMsg(powerRevoke.PowerId)
		if nil == power {
			return nil
		}
		pool.powerMsgQueue.popBy(power.CreateAt())
		// We've directly injected a replacement powerRevokeMsg, notify subsystems
		go pool.msgFeed.Send(event.PowerRevokeMsgEvent{types.PowerRevokeMsgs{powerRevoke}})

	case *types.MetaDataMsg:
		metaData, _ := msg.(*types.MetaDataMsg)
		pool.metaDataMsgQueue.put(metaData)
		// We've directly injected a replacement metaDataMsg, notify subsystems
		go pool.msgFeed.Send(event.MetaDataMsgEvent{types.MetaDataMsgs{metaData}})

	case *types.MetaDataRevokeMsg:
		metaDataRevoke, _ := msg.(*types.MetaDataRevokeMsg)
		metaData := pool.all.getMetaDataMsg(metaDataRevoke.MetaDataId)
		if nil == metaData {
			return nil
		}
		pool.metaDataMsgQueue.popBy(metaData.CreateAt())
		// We've directly injected a replacement metaDataRevokeMsg, notify subsystems
		go pool.msgFeed.Send(event.MetaDataRevokeMsgEvent{types.MetaDataRevokeMsgs{metaDataRevoke}})

	case *types.TaskMsg:
		task, _ := msg.(*types.TaskMsg)
		pool.taskMsgQueue.put(task)
		// We've directly injected a replacement taskMsg, notify subsystems
		go pool.msgFeed.Send(event.TaskMsgEvent{types.TaskMsgs{task}})

	default:
		log.Fatalf("Failed to add msg, can not match the msg type")
	}

	return nil
}




type msgLookup struct {
	allMateDataMsg map[string]*types.MetaDataMsg
	allPowerMsg    map[string]*types.PowerMsg
	allTaskMsg     map[string]*types.TaskMsg

	metaDataMsgLock sync.RWMutex
	powerMsgLock    sync.RWMutex
	taskMsgLock     sync.RWMutex
}

func newMsgLookup() *msgLookup {
	return &msgLookup{
		allMateDataMsg: make(map[string]*types.MetaDataMsg),
		allPowerMsg:    make(map[string]*types.PowerMsg),
		allTaskMsg:     make(map[string]*types.TaskMsg),
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
// RangeTaskMsg calls f on each key and value present in the map.
func (lookup *msgLookup) rangeTaskMsg(f func(taskId string, msg *types.TaskMsg) bool) {
	lookup.taskMsgLock.RLock()
	defer lookup.taskMsgLock.RUnlock()

	for key, value := range lookup.allTaskMsg {
		if !f(key, value) {
			break
		}
	}
}

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
// Get returns a taskMsg if it exists in the lookup, or nil if not found.
func (lookup *msgLookup) getTaskMsg(taskId string) *types.TaskMsg {
	lookup.taskMsgLock.RLock()
	defer lookup.taskMsgLock.RUnlock()

	return lookup.allTaskMsg[taskId]
}

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
// Count returns the current number of items in the lookup.
func (lookup *msgLookup) taskMsgCount() int {
	lookup.taskMsgLock.RLock()
	defer lookup.taskMsgLock.RUnlock()

	return len(lookup.allTaskMsg)
}

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
// Add adds a taskMsg to the lookup.
func (lookup *msgLookup) addTaskMsg(msg *types.TaskMsg) {
	lookup.taskMsgLock.RLock()
	defer lookup.taskMsgLock.RUnlock()

	lookup.allTaskMsg[msg.TaskId] = msg
}

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
// Remove removes a taskMsg from the lookup.
func (lookup *msgLookup) removeTaskMsg(taskId string) {
	lookup.taskMsgLock.RLock()
	delete(lookup.allTaskMsg, taskId)
	lookup.taskMsgLock.RUnlock()
}

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
// Extract removes a taskMsg from the lookup, and return.
func (lookup *msgLookup) extractTaskMsg(taskId string) (*types.TaskMsg, bool) {
	lookup.taskMsgLock.RLock()
	defer lookup.taskMsgLock.RUnlock()
	taskMsg, ok := lookup.allTaskMsg[taskId]
	if !ok {
		return nil, false
	}
	delete(lookup.allTaskMsg, taskId)
	return taskMsg, true
}