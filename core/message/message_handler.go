package message

import (
	"github.com/RosettaFlow/Carrier-Go/event"
	"github.com/RosettaFlow/Carrier-Go/types"
	"sync"
	"time"
)

const (

	defaultPowerMsgsCacheSize    = 5
	defaultMetaDataMsgsCacheSize = 5
	defaultTaskMsgsCacheSize     = 5

	defaultBroadcastPowerMsgInterval    = 100 * time.Millisecond
	defaultBroadcastMetaDataMsgInterval = 100 * time.Millisecond
	defaultBroadcastTaskMsgInterval     = 10 * time.Millisecond
)

type DataHandler interface {
	StoreYarnName(name string) error
	DelYarnName() error
	StoreIdentityId(identity string) error
	DelIdentityId() error
}

type MessageHandler struct {
	pool        *Mempool
	dataHandler DataHandler

	identityMsgCh       chan event.IdentityMsgEvent
	identityRevokeMsgCh chan event.IdentityRevokeMsgEvent
	powerMsgCh          chan event.PowerMsgEvent
	powerRevokeMsgCh    chan event.PowerRevokeMsgEvent
	metaDataMsgCh       chan event.MetaDataMsgEvent
	metaDataRevokeMsgCh chan event.MetaDataRevokeMsgEvent
	taskMsgCh           chan event.TaskMsgEvent

	identityMsgSub       event.Subscription
	identityRevokeMsgSub event.Subscription
	powerMsgSub          event.Subscription
	metaDataMsgSub       event.Subscription
	taskMsgSub           event.Subscription

	powerMsgCache    types.PowerMsgs
	metaDataMsgCache types.MetaDataMsgs
	taskMsgCache     types.TaskMsgs

	lockPower    sync.Mutex
	lockMetaData sync.Mutex
}

func NewHandler(pool *Mempool, dataHandler DataHandler) *MessageHandler {
	m := &MessageHandler{
		pool:        pool,
		dataHandler: dataHandler,
	}
	return m
}

func (m *MessageHandler) Start() error {

	m.identityMsgSub = m.pool.SubscribeNewIdentityMsgsEvent(m.identityMsgCh)
	m.identityRevokeMsgSub = m.pool.SubscribeNewIdentityRevokeMsgsEvent(m.identityRevokeMsgCh)

	m.powerMsgSub = m.pool.SubscribeNewPowerMsgsEvent(m.powerMsgCh)
	m.pool.SubscribeNewPowerRevokeMsgsEvent(m.powerRevokeMsgCh)

	m.metaDataMsgSub = m.pool.SubscribeNewMetaDataMsgsEvent(m.metaDataMsgCh)
	m.pool.SubscribeNewMetaDataRevokeMsgsEvent(m.metaDataRevokeMsgCh)

	m.taskMsgSub = m.pool.SubscribeNewTaskMsgsEvent(m.taskMsgCh)

	go m.loop()
	return nil
}

func (m *MessageHandler) loop() {
	powerTimer := time.NewTimer(defaultBroadcastPowerMsgInterval)
	metaDataTimer := time.NewTimer(defaultBroadcastMetaDataMsgInterval)
	taskTimer := time.NewTimer(defaultBroadcastTaskMsgInterval)

	for {
		select {
		case event := <-m.identityMsgCh:
			if err := m.dataHandler.StoreYarnName(event.Msg.Name); nil != err {
				log.Error("Failed to store yarnNode name on MessageHandler, err:", err)
			}
			if err := m.dataHandler.StoreIdentityId(event.Msg.IdentityId); nil != err {
				log.Error("Failed to store identity on MessageHandler, err:", err)
			}
		case <-m.identityRevokeMsgCh:
			if err := m.dataHandler.DelYarnName(); nil != err {
				log.Error("Failed to delete yarnNode name on MessageHandler, err:", err)
			}
			if err := m.dataHandler.DelIdentityId(); nil != err {
				log.Error("Failed to delete identity on MessageHandler, err:", err)
			}
		case event := <-m.powerMsgCh:
			m.lockPower.Lock()
			m.powerMsgCache = append(m.powerMsgCache, event.Msgs...)
			m.lockPower.Unlock()
			if len(m.powerMsgCache) >= defaultPowerMsgsCacheSize {
				m.BroadcastPowerMsgs(m.powerMsgCache)
				m.powerMsgCache = make(types.PowerMsgs, 0)
				powerTimer.Reset(defaultBroadcastPowerMsgInterval)
			}
		case event := <-m.powerRevokeMsgCh:
			tmp := make(map[string]struct{}, len(event.Msgs))
			for _, msg := range event.Msgs {
				tmp[msg.PowerId] = struct{}{}
			}

			m.lockPower.Lock()
			for i := 0; i < len(m.powerMsgCache); i++ {
				msg := m.powerMsgCache[i]
				if _, ok := tmp[msg.PowerId]; ok {
					m.powerMsgCache = append(m.powerMsgCache[:i], m.powerMsgCache[i+1:]...)
					i--
				}
			}
			m.lockPower.Unlock()

		case event := <-m.metaDataMsgCh:
			m.lockMetaData.Lock()
			m.metaDataMsgCache = append(m.metaDataMsgCache, event.Msgs...)
			m.lockMetaData.Unlock()
			if len(m.metaDataMsgCache) >= defaultMetaDataMsgsCacheSize {
				m.BroadcastMetaDataMsgs(m.metaDataMsgCache)
				m.metaDataMsgCache = make(types.MetaDataMsgs, 0)
				metaDataTimer.Reset(defaultBroadcastMetaDataMsgInterval)
			}
		case event := <-m.metaDataRevokeMsgCh:
			tmp := make(map[string]struct{}, len(event.Msgs))
			for _, msg := range event.Msgs {
				tmp[msg.MetaDataId] = struct{}{}
			}

			m.lockMetaData.Lock()
			for i := 0; i < len(m.metaDataMsgCache); i++ {
				msg := m.metaDataMsgCache[i]
				if _, ok := tmp[msg.MetaDataId]; ok {
					m.metaDataMsgCache = append(m.metaDataMsgCache[:i], m.metaDataMsgCache[i+1:]...)
					i--
				}
			}
			m.lockMetaData.Unlock()

		case event := <-m.taskMsgCh:
			m.taskMsgCache = append(m.taskMsgCache, event.Msgs...)
			if len(m.taskMsgCache) >= defaultTaskMsgsCacheSize {
				m.BroadcastTaskMsgs(m.taskMsgCache)
				m.taskMsgCache = make(types.TaskMsgs, 0)
				taskTimer.Reset(defaultBroadcastTaskMsgInterval)
			}

		case <-powerTimer.C:
			if len(m.powerMsgCache) >= 0 {
				m.BroadcastPowerMsgs(m.powerMsgCache)
				m.powerMsgCache = make(types.PowerMsgs, 0)
				powerTimer.Reset(defaultBroadcastPowerMsgInterval)
			}

		case <-metaDataTimer.C:
			if len(m.metaDataMsgCache) >= 0 {
				m.BroadcastMetaDataMsgs(m.metaDataMsgCache)
				m.metaDataMsgCache = make(types.MetaDataMsgs, 0)
				powerTimer.Reset(defaultBroadcastMetaDataMsgInterval)
			}

		case <-taskTimer.C:
			if len(m.taskMsgCache) >= 0 {
				m.BroadcastTaskMsgs(m.taskMsgCache)
				m.taskMsgCache = make(types.TaskMsgs, 0)
				taskTimer.Reset(defaultBroadcastTaskMsgInterval)
			}
			// Err() channel will be closed when unsubscribing.
		case <-m.powerMsgSub.Err():
			return
		case <-m.metaDataMsgSub.Err():
			return
		case <-m.taskMsgSub.Err():
			return
		}
	}
}

func (m *MessageHandler) BroadcastPowerMsgs(powerMsgs types.PowerMsgs) error {

	return nil
}

func (m *MessageHandler) BroadcastMetaDataMsgs(metaDataMsgs types.MetaDataMsgs) error {

	return nil
}

func (m *MessageHandler) BroadcastTaskMsgs(taskMsgs types.TaskMsgs) error {

	return nil
}
