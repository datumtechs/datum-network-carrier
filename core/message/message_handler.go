package message

import (
	"github.com/RosettaFlow/Carrier-Go/event"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
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
	// TODO 本地存储当前调度服务自身的  identity
	StoreIdentity(identity *types.NodeAlias) error
	DelIdentity() error
	GetYarnName() (string, error)
	GetIdentityID() (string, error)
	GetIdentity() (*types.NodeAlias, error)
}

type DataCenter interface {
	// identity
	InsertIdentity(identity *types.Identity) error
	RevokeIdentity(identity *types.Identity) error

	// power
	InsertResource(resource *types.Resource) error
}

type MessageHandler struct {
	pool        *Mempool
	dataHandler DataHandler
	center      DataCenter

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

func NewHandler(pool *Mempool, dataHandler DataHandler, dataCenter DataCenter) *MessageHandler {
	m := &MessageHandler{
		pool:        pool,
		dataHandler: dataHandler,
		center:      dataCenter,
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
			if err := m.BroadcastIdentityMsg(event.Msg); nil != err {
				log.Error("Failed to broadcast org identityMsg  on MessageHandler, err:", err)
			}
		case <-m.identityRevokeMsgCh:
			if err := m.BroadcastIdentityRevokeMsg(); nil != err {
				log.Error("Failed to remove org identity on MessageHandler, err:", err)
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

func (m *MessageHandler) BroadcastIdentityMsg(msg *types.IdentityMsg) error {
	if err := m.dataHandler.StoreIdentity(msg.NodeAlias); nil != err {
		log.Error("Failed to store local org identity on MessageHandler, err:", err)
		return err
	}

	if err := m.center.InsertIdentity(
		types.NewIdentity(&libTypes.IdentityData{
			NodeName: msg.Name,
			NodeId:   msg.NodeId,
			Identity: msg.IdentityId,
		})); nil != err {
		log.Error("Failed to broadcast org org identity on MessageHandler, err:", err)
		return err
	}

	return nil
}

func (m *MessageHandler) BroadcastIdentityRevokeMsg() error {
	identity, err := m.dataHandler.GetIdentity()
	if nil != err {
		log.Error("Failed to get local org identity on MessageHandler, err:", err)
		return err
	}
	if err := m.dataHandler.DelIdentity(); nil != err {
		log.Error("Failed to delete org identity to local on MessageHandler, err:", err)
		return err
	}

	if err := m.center.RevokeIdentity(
		types.NewIdentity(&libTypes.IdentityData{
			NodeName: identity.Name,
			NodeId:   identity.NodeId,
			Identity: identity.IdentityId,
		})); nil != err {
		log.Error("Failed to remove org identity to remote on MessageHandler, err:", err)
		return err
	}
	return nil
}

func (m *MessageHandler) BroadcastPowerMsgs(powerMsgs types.PowerMsgs) error {

	for _, power := range powerMsgs {
		m.center.InsertResource(types.NewResource(&libTypes.ResourceData{
			//Identity: power.Data.IdentityId,
			//NodeId: power.Data.NodeId,
			//NodeName: power.Data.NodeId,
			//DataId: power.PowerId,
			//// the status of data, N means normal, D means deleted.
			//DataStatus: "N",
			//// resource status, eg: create/release/revoke
			//State:
			//// unit: byte
			//TotalMem uint64 `protobuf:"varint,7,opt,name=totalMem,proto3" json:"totalMem,omitempty"`
			//// unit: byte
			//UsedMem uint64 `protobuf:"varint,8,opt,name=usedMem,proto3" json:"usedMem,omitempty"`
			//// number of cpu cores.
			//TotalProcessor uint32 `protobuf:"varint,9,opt,name=totalProcessor,proto3" json:"totalProcessor,omitempty"`
			//// unit: byte
			//TotalBandWidth       uint64   `protobuf:"varint,10,opt,name=totalBandWidth,proto3" json:"totalBandWidth,omitempty"`
		}))
	}

	return nil
}

func (m *MessageHandler) BroadcastMetaDataMsgs(metaDataMsgs types.MetaDataMsgs) error {

	return nil
}

func (m *MessageHandler) BroadcastTaskMsgs(taskMsgs types.TaskMsgs) error {

	return nil
}
