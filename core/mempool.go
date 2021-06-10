package core

import (
	"github.com/RosettaFlow/Carrier-Go/core/message"
	"github.com/RosettaFlow/Carrier-Go/event"
	"github.com/RosettaFlow/Carrier-Go/types"
	"sync"
)

type Msg interface {
	Marshal() ([]byte, error)
	Unmarshal(b []byte) error
	String() string
	MsgType() string
}

type MempoolConfig struct {
}







type Mempool struct {
	cfg *MempoolConfig

	msgFeed event.Feed
	scope   event.SubscriptionScope

	mateDataMsgQueue *message.MetaDataMsgList
	powerMsgQueue    *message.PowerMsgList
	taskMsgQueue     *message.TaskMsgList

	all			*msgLookup
}

func New(cfg *MempoolConfig) *Mempool {
	return &Mempool{
		cfg:              cfg,
		mateDataMsgQueue: message.NewMetaDataMsgList(),
		powerMsgQueue:    message.NewPowerMsgList(),
		taskMsgQueue:     message.NewTaskMsgList(),
		all: newMsgLookup(),
	}
}

// SubscribeNewTxsEvent registers a subscription of NewTxsEvent and
// starts sending event to the given channel.
func (pool *Mempool) SubscribeNewMetaDataMsgsEvent(ch chan<- event.MetaDataMsgEvent) event.Subscription {
	return pool.scope.Track(pool.msgFeed.Subscribe(ch))
}
func (pool *Mempool) SubscribeNewPowerMsgsEvent(ch chan<- event.PowerMsgEvent) event.Subscription {
	return pool.scope.Track(pool.msgFeed.Subscribe(ch))
}
func (pool *Mempool) SubscribeNewTaskMsgsEvent(ch chan<- event.TaskMsgEvent) event.Subscription {
	return pool.scope.Track(pool.msgFeed.Subscribe(ch))
}

func (pool *Mempool) Add(msg Msg) error {

	switch msg.(type) {
	case *types.PowerMsg:
		power, _ := msg.(*types.PowerMsg)
		pool.powerMsgQueue.Put(power)
		// We've directly injected a replacement powerMsg, notify subsystems
		go pool.msgFeed.Send(event.PowerMsgEvent{types.PowerMsgs{power}})
	case *types.MetaDataMsg:
		metaData, _ := msg.(*types.MetaDataMsg)
		pool.mateDataMsgQueue.Put(metaData)
		// We've directly injected a replacement metaDataMsg, notify subsystems
		go pool.msgFeed.Send(event.MetaDataMsgEvent{types.MetaDataMsgs{metaData}})
	case *types.TaskMsg:
		task, _ := msg.(*types.TaskMsg)
		pool.taskMsgQueue.Put(task)
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
func (lookup *msgLookup) RangeMetaDataMsg(f func(metaDataId string, msg *types.MetaDataMsg) bool) {
	lookup.metaDataMsgLock.RLock()
	defer lookup.metaDataMsgLock.RUnlock()

	for key, value := range lookup.allMateDataMsg {
		if !f(key, value) {
			break
		}
	}
}
// RangePowerMsg calls f on each key and value present in the map.
func (lookup *msgLookup) RangePowerMsg(f func(powerId string, msg *types.PowerMsg) bool) {
	lookup.powerMsgLock.RLock()
	defer lookup.powerMsgLock.RUnlock()

	for key, value := range lookup.allPowerMsg {
		if !f(key, value) {
			break
		}
	}
}
// RangeTaskMsg calls f on each key and value present in the map.
func (lookup *msgLookup) RangeTaskMsg(f func(taskId string, msg *types.TaskMsg) bool) {
	lookup.taskMsgLock.RLock()
	defer lookup.taskMsgLock.RUnlock()

	for key, value := range lookup.allTaskMsg {
		if !f(key, value) {
			break
		}
	}
}

// Get returns a metaDataMsg if it exists in the lookup, or nil if not found.
func (lookup *msgLookup) GetMetaDataMsg(metaDataId string) *types.MetaDataMsg {
	lookup.metaDataMsgLock.RLock()
	defer lookup.metaDataMsgLock.RUnlock()

	return lookup.allMateDataMsg[metaDataId]
}
// Get returns a powerMsg if it exists in the lookup, or nil if not found.
func (lookup *msgLookup) GetPowerMsg(powerId string) *types.PowerMsg {
	lookup.powerMsgLock.RLock()
	defer lookup.powerMsgLock.RUnlock()

	return lookup.allPowerMsg[powerId]
}
// Get returns a taskMsg if it exists in the lookup, or nil if not found.
func (lookup *msgLookup) GetTaskMsg(taskId string) *types.TaskMsg {
	lookup.taskMsgLock.RLock()
	defer lookup.taskMsgLock.RUnlock()

	return lookup.allTaskMsg[taskId]
}

// Count returns the current number of items in the lookup.
func (lookup *msgLookup) MetaDataMsgCount() int {
	lookup.metaDataMsgLock.RLock()
	defer lookup.metaDataMsgLock.RUnlock()

	return len(lookup.allMateDataMsg)
}
// Count returns the current number of items in the lookup.
func (lookup *msgLookup) PowerMsgCount() int {
	lookup.powerMsgLock.RLock()
	defer lookup.powerMsgLock.RUnlock()

	return len(lookup.allPowerMsg)
}
// Count returns the current number of items in the lookup.
func (lookup *msgLookup) TaskMsgCount() int {
	lookup.taskMsgLock.RLock()
	defer lookup.taskMsgLock.RUnlock()

	return len(lookup.allTaskMsg)
}

// Add adds a metaDataMsg to the lookup.
func (lookup *msgLookup) AddMetaDataMsg(msg *types.MetaDataMsg) {
	lookup.metaDataMsgLock.RLock()
	defer lookup.metaDataMsgLock.RUnlock()

	lookup.allMateDataMsg[msg.MetaDataId] = msg
}
// Add adds a powerMsg to the lookup.
func (lookup *msgLookup) AddPowerMsg(msg *types.PowerMsg) {
	lookup.powerMsgLock.RLock()
	defer lookup.powerMsgLock.RUnlock()

	lookup.allPowerMsg[msg.PowerId] = msg
}
// Add adds a taskMsg to the lookup.
func (lookup *msgLookup) AddTaskMsg(msg *types.TaskMsg) {
	lookup.taskMsgLock.RLock()
	defer lookup.taskMsgLock.RUnlock()

	lookup.allTaskMsg[msg.TaskId] = msg
}

// Remove removes a metaDataMsg from the lookup.
func (lookup *msgLookup) RemoveMetaDataMsg(metaDataId string) {
	lookup.metaDataMsgLock.RLock()
	defer lookup.metaDataMsgLock.RUnlock()

	delete(lookup.allMateDataMsg, metaDataId)
}
// Remove removes a powerMsg from the lookup.
func (lookup *msgLookup) RemovePowerMsg(powerId string) {
	lookup.powerMsgLock.RLock()
	defer lookup.powerMsgLock.RUnlock()

	delete(lookup.allPowerMsg, powerId)
}
// Remove removes a taskMsg from the lookup.
func (lookup *msgLookup) RemoveTaskMsg(taskId string) {
	lookup.taskMsgLock.RLock()
	defer lookup.taskMsgLock.RUnlock()

	delete(lookup.allTaskMsg, taskId)
}