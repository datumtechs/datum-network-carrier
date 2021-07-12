package message

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/core/iface"
	"github.com/RosettaFlow/Carrier-Go/core/task"
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


type MessageHandler struct {
	pool       *Mempool
	dataCenter iface.ForHandleDB

	// Send taskMsg to taskManager
	taskManager *task.Manager

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

func NewHandler(pool *Mempool, dataCenter iface.ForHandleDB, taskManager *task.Manager) *MessageHandler {
	m := &MessageHandler{
		pool:        pool,
		dataCenter:  dataCenter,
		taskManager: taskManager,
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

	// add identity to local db
	if err := m.dataCenter.StoreIdentity(msg.NodeAlias); nil != err {
		log.Error("Failed to store local org identity on MessageHandler, err:", err)
		return err
	}

	// send identity to datacenter
	if err := m.dataCenter.InsertIdentity(msg.ToDataCenter()); nil != err {
		log.Error("Failed to broadcast org org identity on MessageHandler, err:", err)
		return err
	}

	return nil
}

func (m *MessageHandler) BroadcastIdentityRevokeMsg() error {

	// remove identity from local db
	identity, err := m.dataCenter.GetIdentity()
	if nil != err {
		log.Error("Failed to get local org identity on MessageHandler, err:", err)
		return err
	}
	if err := m.dataCenter.RemoveIdentity(); nil != err {
		log.Error("Failed to delete org identity to local on MessageHandler, err:", err)
		return err
	}

	// remove identity from dataCenter
	if err := m.dataCenter.RevokeIdentity(
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
		// 存储本地的 资源信息
		if err := m.dataCenter.StoreLocalResourceTable(types.NewLocalResourceTable(power.JobNodeId(),
			power.Memory(), power.Processor(), power.Bandwidth())); nil != err {
			log.Error("Failed to StoreLocalResourceTable on MessageHandler, err:", err)
			errs = append(errs, fmt.Sprintf("powerId: %s, %s", power.PowerId, err))
			continue
		}

		if err := m.dataCenter.StoreLocalResourceIdByPowerId(power.PowerId, power.JobNodeId()); nil != err {
			log.Error("Failed to store powerId and jobNodeId mapping to local on MessageHandler, err:", err)
			errs = append(errs, fmt.Sprintf("powerId: %s, %s", power.PowerId, err))
			continue
		}
		if err := m.dataCenter.InsertLocalResource(power.ToLocal()); nil != err {
			log.Error("Failed to store power to local on MessageHandler, err:", err)
			errs = append(errs, fmt.Sprintf("powerId: %s, %s", power.PowerId, err))
			continue
		}

		// 发布到全网
		if err := m.dataCenter.InsertResource(power.ToDataCenter()); nil != err {
			log.Error("Failed to store power to dataCenter on MessageHandler, err:", err)
			errs = append(errs, fmt.Sprintf("powerId: %s, %s", power.PowerId, err))
			continue
		}

	}
	if len(errs) != 0 {
		return fmt.Errorf("broadcast powerMsgs err: %s", strings.Join(errs, "\n"))
	}
	return nil
}

func (m *MessageHandler) BroadcastPowerRevokeMsgs(powerRevokeMsgs types.PowerRevokeMsgs) error {
	errs := make([]string, 0)
	for _, revoke := range powerRevokeMsgs {

		jobNodeId, err := m.dataCenter.QueryLocalResourceIdByPowerId(revoke.PowerId)
		if nil != err {
			log.Error("Failed to query jobNodeId on local on MessageHandler, err:", err)
			errs = append(errs, fmt.Sprintf("powerId: %s, %s", revoke.PowerId, err))
			continue
		}

		if err := m.dataCenter.RemoveLocalResourceTable(jobNodeId); nil != err {
			log.Error("Failed to RemoveLocalResourceTable on MessageHandler, err:", err)
			errs = append(errs, fmt.Sprintf("powerId: %s, jobNodeId: %s, %s", revoke.PowerId, jobNodeId, err))
			continue
		}

		if err := m.dataCenter.RemoveLocalResource(jobNodeId); nil != err {
			log.Error("Failed to remove powerId and jobNodeId mapping to local on MessageHandler, err:", err)
			errs = append(errs, fmt.Sprintf("powerId: %s, jobNodeId: %s, %s", revoke.PowerId, jobNodeId, err))
			continue
		}
		if err := m.dataCenter.RemoveLocalResource(jobNodeId); nil != err {
			log.Error("Failed to remove local resource on MessageHandler, err:", err)
			errs = append(errs, fmt.Sprintf("powerId: %s, jobNodeId: %s, %s", revoke.PowerId, jobNodeId, err))
			continue
		}

		if err := m.dataCenter.InsertResource(revoke.ToDataCenter()); nil != err {
			log.Error("Failed to remove dataCenter resource on MessageHandler, err:", err)
			errs = append(errs, fmt.Sprintf("powerId: %s, jobNodeId: %s, %s", revoke.PowerId, jobNodeId, err))
			continue
		}
	}
	if len(errs) != 0 {
		return fmt.Errorf("broadcast powerRevokeMsgs err: %s", strings.Join(errs, "\n"))
	}
	return nil
}

func (m *MessageHandler) BroadcastMetaDataMsgs(metaDataMsgs types.MetaDataMsgs) error {
	errs := make([]string, 0)
	for _, metaData := range metaDataMsgs {

		// 维护本地 数据服务的 orginId  和 metaDataId 关系
		dataResourceDataUsed, err := m.dataCenter.QueryDataResourceDataUsed(metaData.OriginId())
		if nil != err {
			log.Error("Failed to QueryDataResourceDataUsed on MessageHandler, err:", err)
			errs = append(errs, fmt.Sprintf("originId: %s, %s", metaData.OriginId(), err))
			continue
		}
		dataResourceDataUsed.SetMetaDataId(metaData.MetaDataId)
		if err := m.dataCenter.StoreDataResourceDataUsed(dataResourceDataUsed); nil != err {
			log.Error("Failed to StoreDataResourceDataUsed on MessageHandler, err:", err)
			errs = append(errs, fmt.Sprintf("originId: %s, %s", metaData.OriginId(), err))
			continue
		}
		if err := m.dataCenter.InsertMetadata(metaData.ToDataCenter()); nil != err {
			log.Error("Failed to store metaData to dataCenter on MessageHandler, err:", err)
			errs = append(errs, fmt.Sprintf("originId: %s, %s", metaData.OriginId(), err))
			continue
		}
	}
	if len(errs) != 0 {
		return fmt.Errorf("broadcast metaDataMsgs err: %s", strings.Join(errs, "\n"))
	}
	return nil
}

func (m *MessageHandler) BroadcastMetaDataRevokeMsgs(metaDataRevokeMsgs types.MetaDataRevokeMsgs) error {
	errs := make([]string, 0)
	for _, revoke := range metaDataRevokeMsgs {
		// 不需要维护本地信息 ...
		if err := m.dataCenter.InsertMetadata(revoke.ToDataCenter()); nil != err {
			log.Error("Failed to store metaData to dataCenter on MessageHandler, err:", err)
			errs = append(errs, fmt.Sprintf("metaDataId: %s, %s", revoke.MetaDataId, err))
			continue
		}
	}
	if len(errs) != 0 {
		return fmt.Errorf("broadcast metaDataRevokeMsgs err: %s", strings.Join(errs, "\n"))
	}
	return nil
}

func (m *MessageHandler) BroadcastTaskMsgs(taskMsgs types.TaskMsgs) error {
	return m.taskManager.SendTaskMsgs(taskMsgs)
}
