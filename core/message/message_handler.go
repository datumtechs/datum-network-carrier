package message

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/event"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strings"
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
	GetIdentityId() (string, error)
	GetIdentity() (*types.NodeAlias, error)
}

type DataCenter interface {
	// identity
	InsertIdentity(identity *types.Identity) error
	RevokeIdentity(identity *types.Identity) error

	// power
	InsertResource(resource *types.Resource) error
	// metaData
	InsertMetadata(metadata *types.Metadata) error
}

type MessageHandler struct {
	pool        *Mempool
	dataHandler DataHandler
	center      DataCenter
	//// Consensuses
	//engines map[string]consensus.Engine

	taskCh chan<- types.TaskMsgs

	identityMsgCh       chan types.IdentityMsgEvent
	identityRevokeMsgCh chan types.IdentityRevokeMsgEvent
	powerMsgCh          chan types.PowerMsgEvent
	powerRevokeMsgCh    chan types.PowerRevokeMsgEvent
	metaDataMsgCh       chan types.MetaDataMsgEvent
	metaDataRevokeMsgCh chan types.MetaDataRevokeMsgEvent
	taskMsgCh           chan types.TaskMsgEvent

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

func NewHandler(pool *Mempool, dataHandler DataHandler, dataCenter DataCenter, taskCh chan<- types.TaskMsgs) *MessageHandler {
	m := &MessageHandler{
		pool:        pool,
		dataHandler: dataHandler,
		center:      dataCenter,
		taskCh:      taskCh,
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
				if err := m.BroadcastPowerMsgs(m.powerMsgCache); nil != err {
					log.Error(fmt.Sprintf("%s", err))
				}
				m.powerMsgCache = make(types.PowerMsgs, 0)
				powerTimer.Reset(defaultBroadcastPowerMsgInterval)
			}
		case event := <-m.powerRevokeMsgCh:
			tmp := make(map[string]int, len(event.Msgs))
			for i, msg := range event.Msgs {
				tmp[msg.PowerId] = i
			}

			// Remove local cache powerMsgs
			m.lockPower.Lock()
			for i := 0; i < len(m.powerMsgCache); i++ {
				msg := m.powerMsgCache[i]
				if _, ok := tmp[msg.PowerId]; ok {
					delete(tmp, msg.PowerId)
					m.powerMsgCache = append(m.powerMsgCache[:i], m.powerMsgCache[i+1:]...)
					i--
				}
			}
			m.lockPower.Unlock()

			// Revoke remote power
			if len(tmp) != 0 {
				msgs, index := make(types.PowerRevokeMsgs, len(tmp)), 0
				for _, i := range tmp {
					msgs[index] = event.Msgs[i]
					index++
				}
				if err := m.BroadcastPowerRevokeMsgs(msgs); nil != err {
					log.Error(fmt.Sprintf("%s", err))
				}
			}

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
			tmp := make(map[string]int, len(event.Msgs))
			for i, msg := range event.Msgs {
				tmp[msg.MetaDataId] = i
			}

			// Remove local cache metaDataMsgs
			m.lockMetaData.Lock()
			for i := 0; i < len(m.metaDataMsgCache); i++ {
				msg := m.metaDataMsgCache[i]
				if _, ok := tmp[msg.MetaDataId]; ok {
					delete(tmp, msg.MetaDataId)
					m.metaDataMsgCache = append(m.metaDataMsgCache[:i], m.metaDataMsgCache[i+1:]...)
					i--
				}
			}
			m.lockMetaData.Unlock()

			// Revoke remote metaData
			if len(tmp) != 0 {
				msgs, index := make(types.MetaDataRevokeMsgs, len(tmp)), 0
				for _, i := range tmp {
					msgs[index] = event.Msgs[i]
					index++
				}
				if err := m.BroadcastMetaDataRevokeMsgs(msgs); nil != err {
					log.Error(fmt.Sprintf("%s", err))
				}
			}

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
	errs := make([]string, 0)
	for _, power := range powerMsgs {
		err := m.center.InsertResource(types.NewResource(&libTypes.ResourceData{
			Identity: power.OwnerIdentityId(),
			NodeId:   power.OwnerNodeId(),
			NodeName: power.OwnerName(),
			DataId:   power.PowerId,
			// the status of data, N means normal, D means deleted.
			DataStatus: types.ResourceDataStatusN.String(),
			// resource status, eg: create/release/revoke
			State: types.PowerStateRelease.String(),
			// unit: byte
			TotalMem: power.Memory(),
			// unit: byte
			UsedMem: 0,
			// number of cpu cores.
			TotalProcessor: power.Processor(),

			UsedProcessor: 0,

			// unit: byte
			TotalBandWidth: power.Bandwidth(),
			UsedBandWidth:  0,
		}))
		errs = append(errs, fmt.Sprintf("powerId: %s, %s", power.PowerId, err))
	}
	if len(errs) != 0 {
		return fmt.Errorf("broadcast powerMsgs err: %s", strings.Join(errs, "\n"))
	}
	return nil
}

func (m *MessageHandler) BroadcastPowerRevokeMsgs(powerRevokeMsgs types.PowerRevokeMsgs) error {
	errs := make([]string, 0)
	for _, revoke := range powerRevokeMsgs {
		err := m.center.InsertResource(types.NewResource(&libTypes.ResourceData{
			Identity: revoke.IdentityId,
			NodeId:   revoke.NodeId,
			NodeName: revoke.Name,
			DataId:   revoke.PowerId,
			// the status of data, N means normal, D means deleted.
			DataStatus: types.ResourceDataStatusD.String(),
			// resource status, eg: create/release/revoke
			State: types.PowerStateRevoke.String(),
			// unit: byte
			TotalMem: 0,
			// unit: byte
			UsedMem: 0,
			// number of cpu cores.
			TotalProcessor: 0,
			// unit: byte
			TotalBandWidth: 0,
		}))
		errs = append(errs, fmt.Sprintf("powerId: %s, %s", revoke.PowerId, err))
	}
	if len(errs) != 0 {
		return fmt.Errorf("broadcast powerRevokeMsgs err: %s", strings.Join(errs, "\n"))
	}
	return nil
}

func (m *MessageHandler) BroadcastMetaDataMsgs(metaDataMsgs types.MetaDataMsgs) error {
	errs := make([]string, 0)
	for _, metaData := range metaDataMsgs {
		err := m.center.InsertMetadata(types.NewMetadata(&libTypes.MetaData{
			Identity:       metaData.OwnerIdentityId(),
			NodeId:         metaData.OwnerNodeId(),
			NodeName:       metaData.OwnerName(),
			DataId:         metaData.MetaDataId,
			OriginId:       metaData.OriginId(),
			TableName:      metaData.TableName(),
			FilePath:       metaData.FilePath(),
			FileType:       metaData.FileType(),
			Desc:           metaData.Desc(),
			Rows:           uint64(metaData.Rows()),
			Columns:        uint64(metaData.Columns()),
			Size_:          uint64(metaData.Size()),
			HasTitleRow:    metaData.HasTitle(),
			ColumnMetaList: metaData.ColumnMetas(),
			// the status of data, N means normal, D means deleted.
			DataStatus: types.ResourceDataStatusN.String(),
			// metaData status, eg: create/release/revoke
			State: types.MetaDataStateRelease.String(),
		}))
		errs = append(errs, fmt.Sprintf("metaDataId: %s, %s", metaData.MetaDataId, err))
	}
	if len(errs) != 0 {
		return fmt.Errorf("broadcast metaDataMsgs err: %s", strings.Join(errs, "\n"))
	}
	return nil
}

func (m *MessageHandler) BroadcastMetaDataRevokeMsgs(metaDataRevokeMsgs types.MetaDataRevokeMsgs) error {
	errs := make([]string, 0)
	for _, revoke := range metaDataRevokeMsgs {
		err := m.center.InsertMetadata(types.NewMetadata(&libTypes.MetaData{
			Identity: revoke.IdentityId,
			NodeId:   revoke.NodeId,
			NodeName: revoke.Name,
			DataId:   revoke.MetaDataId,
			// the status of data, N means normal, D means deleted.
			DataStatus: types.ResourceDataStatusD.String(),
			// metaData status, eg: create/release/revoke
			State: types.MetaDataStateRevoke.String(),
		}))
		errs = append(errs, fmt.Sprintf("metaDataId: %s, %s", revoke.MetaDataId, err))
	}
	if len(errs) != 0 {
		return fmt.Errorf("broadcast metaDataRevokeMsgs err: %s", strings.Join(errs, "\n"))
	}
	return nil
}

func (m *MessageHandler) BroadcastTaskMsgs(taskMsgs types.TaskMsgs) error {
	m.taskCh <- taskMsgs
	return nil
}
