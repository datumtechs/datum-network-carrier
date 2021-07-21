package message

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/feed"
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
	defaultMetaDataMsgsCacheSize = 1
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

	msgChannel chan *feed.Event
	//identityMsgCh       chan types.IdentityMsgEvent
	//identityRevokeMsgCh chan types.IdentityRevokeMsgEvent
	//powerMsgCh          chan types.PowerMsgEvent
	//powerRevokeMsgCh    chan types.PowerRevokeMsgEvent
	//metaDataMsgCh       chan types.MetaDataMsgEvent
	//metaDataRevokeMsgCh chan types.MetaDataRevokeMsgEvent
	//taskMsgCh           chan types.TaskMsgEvent
	msgSub               event.Subscription
	//identityMsgSub       event.Subscription
	//identityRevokeMsgSub event.Subscription
	//powerMsgSub          event.Subscription
	//metaDataMsgSub       event.Subscription
	//taskMsgSub           event.Subscription

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
		msgChannel:  make(chan *feed.Event, 5),
	}
	return m
}

func (m *MessageHandler) Start() error {
	m.msgSub = m.pool.SubscribeNewMessageEvent(m.msgChannel)

	//m.identityMsgSub = m.pool.SubscribeNewIdentityMsgsEvent(m.msgChannel)
	//m.identityRevokeMsgSub = m.pool.SubscribeNewIdentityRevokeMsgsEvent(m.msgChannel)

	//m.powerMsgSub = m.pool.SubscribeNewPowerMsgsEvent(m.msgChannel)
	//m.pool.SubscribeNewPowerRevokeMsgsEvent(m.msgChannel)

	//m.metaDataMsgSub = m.pool.SubscribeNewMetaDataMsgsEvent(m.msgChannel)
	//m.pool.SubscribeNewMetaDataRevokeMsgsEvent(m.msgChannel)

	//m.taskMsgSub = m.pool.SubscribeNewTaskMsgsEvent(m.msgChannel)

	go m.loop()
	log.Info("Started messageManager ...")
	return nil
}

func (m *MessageHandler) loop() {
	powerTimer := time.NewTimer(defaultBroadcastPowerMsgInterval)
	metaDataTimer := time.NewTimer(defaultBroadcastMetaDataMsgInterval)
	taskTimer := time.NewTimer(defaultBroadcastTaskMsgInterval)

	for {
		select {
		case event := <-m.msgChannel:
			switch event.Type {
			case types.ApplyIdentity:
				eventMessage := event.Data.(*types.IdentityMsgEvent)
				if err := m.BroadcastIdentityMsg(eventMessage.Msg); nil != err {
					log.Error("Failed to broadcast org identityMsg  on MessageHandler, err:", err)
				}
			case types.RevokeIdentity:
				if err := m.BroadcastIdentityRevokeMsg(); nil != err {
					log.Error("Failed to remove org identity on MessageHandler, err:", err)
				}
			case types.ApplyPower:
				m.lockPower.Lock()
				msg := event.Data.(*types.PowerMsgEvent)
				m.powerMsgCache = append(m.powerMsgCache, msg.Msgs...)
				m.lockPower.Unlock()
				if len(m.powerMsgCache) >= defaultPowerMsgsCacheSize {
					if err := m.BroadcastPowerMsgs(m.powerMsgCache); nil != err {
						log.Error(fmt.Sprintf("%s", err))
					}
					m.powerMsgCache = make(types.PowerMsgs, 0)
					powerTimer.Reset(defaultBroadcastPowerMsgInterval)
				}
			case types.RevokePower:
				eventMessage := event.Data.(*types.PowerRevokeMsgEvent)
				tmp := make(map[string]int, len(eventMessage.Msgs))
				for i, msg := range eventMessage.Msgs {
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
						msgs[index] = eventMessage.Msgs[i]
						index++
					}
					if err := m.BroadcastPowerRevokeMsgs(msgs); nil != err {
						log.Error(fmt.Sprintf("%s", err))
					}
				}
			case types.ApplyMetadata:
				eventMessage := event.Data.(*types.MetaDataMsgEvent)
				m.lockMetaData.Lock()
				m.metaDataMsgCache = append(m.metaDataMsgCache, eventMessage.Msgs...)
				m.lockMetaData.Unlock()
				if len(m.metaDataMsgCache) >= defaultMetaDataMsgsCacheSize {
					m.BroadcastMetaDataMsgs(m.metaDataMsgCache)
					m.metaDataMsgCache = make(types.MetaDataMsgs, 0)
					metaDataTimer.Reset(defaultBroadcastMetaDataMsgInterval)
				}
			case types.RevokeMetadata:
				eventMessage := event.Data.(*types.MetaDataRevokeMsgEvent)
				tmp := make(map[string]int, len(eventMessage.Msgs))
				for i, msg := range eventMessage.Msgs {
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
						msgs[index] = eventMessage.Msgs[i]
						index++
					}
					if err := m.BroadcastMetaDataRevokeMsgs(msgs); nil != err {
						log.Error(fmt.Sprintf("%s", err))
					}
				}
			case types.ApplyTask:
				eventMessage := event.Data.(*types.TaskMsgEvent)
				m.taskMsgCache = append(m.taskMsgCache, eventMessage.Msgs...)
				if len(m.taskMsgCache) >= defaultTaskMsgsCacheSize {
					m.BroadcastTaskMsgs(m.taskMsgCache)
					m.taskMsgCache = make(types.TaskMsgs, 0)
					taskTimer.Reset(defaultBroadcastTaskMsgInterval)
				}
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
		case <-m.msgSub.Err():
			return
		}
	}
}

func (m *MessageHandler) BroadcastIdentityMsg(msg *types.IdentityMsg) error {

	// add identity to local db
	if err := m.dataCenter.StoreIdentity(msg.NodeAlias); nil != err {
		log.Errorf("Failed to store local org identity on MessageHandler, identityId: {%s}, err: {%s}", msg.IdentityId, err)
		return err
	}

	// send identity to datacenter
	if err := m.dataCenter.InsertIdentity(msg.ToDataCenter()); nil != err {
		log.Errorf("Failed to broadcast org org identity on MessageHandler, identityId: {%s}, nodeId: {%s}, nodeName: {%s}, err: {%s}", msg.IdentityId, msg.NodeId, msg.Name, err)
		return err
	}
	log.Debugf("Registered identity succeed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}", msg.IdentityId, msg.NodeId, msg.Name)
	return nil
}

func (m *MessageHandler) BroadcastIdentityRevokeMsg() error {

	// remove identity from local db
	identity, err := m.dataCenter.GetIdentity()
	if nil != err {
		log.Errorf("Failed to get local org identity on MessageHandler, identityId: {%s}, err: {%s}", identity, err)
		return err
	}
	if err := m.dataCenter.RemoveIdentity(); nil != err {
		log.Errorf("Failed to delete org identity to local on MessageHandler, identityId: {%s}, err: {%s}", identity, err)
		return err
	}

	// remove identity from dataCenter
	if err := m.dataCenter.RevokeIdentity(
		types.NewIdentity(&libTypes.IdentityData{
			NodeName: identity.Name,
			NodeId:   identity.NodeId,
			Identity: identity.IdentityId,
		})); nil != err {
		log.Errorf("Failed to remove org identity to remote on MessageHandler, identityId: {%s}, err: {%s}", identity, err)
		return err
	}
	log.Debugf("Revoke identity succeed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}", identity.IdentityId, identity.NodeId, identity.Name)
	return nil
}

func (m *MessageHandler) BroadcastPowerMsgs(powerMsgs types.PowerMsgs) error {

	identity, err := m.dataCenter.GetIdentity()
	if nil != err {
		return fmt.Errorf("Failed to broadcast powerMsgs, query local identityInfo failed, {%s}", err)
	}

	errs := make([]string, 0)
	for _, power := range powerMsgs {
		// 存储本地的 资源信息
		if err := m.dataCenter.StoreLocalResourceTable(types.NewLocalResourceTable(power.JobNodeId, power.PowerId,
			types.GetDefaultResoueceMem(), types.GetDefaultResoueceProcessor(), types.GetDefaultResoueceBandwidth())); nil != err {
			log.Errorf("Failed to StoreLocalResourceTable on MessageHandler, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.PowerId, power.JobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to StoreLocalResourceTable on MessageHandler, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.PowerId, power.JobNodeId, err))
			continue
		}

		if err := m.dataCenter.StoreLocalResourceIdByPowerId(power.PowerId, power.JobNodeId); nil != err {
			log.Errorf("Failed to StoreLocalResourceIdByPowerId on MessageHandler,  powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.PowerId, power.JobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to StoreLocalResourceIdByPowerId on MessageHandler,  powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.PowerId, power.JobNodeId, err))
			continue
		}
		if err := m.dataCenter.InsertLocalResource(types.NewLocalResource(&libTypes.LocalResourceData{
			Identity:  identity.IdentityId,
			NodeId:    identity.NodeId,
			NodeName:  identity.Name,
			JobNodeId: power.JobNodeId,
			DataId:    power.PowerId,
			// the status of data, N means normal, D means deleted.
			DataStatus: types.ResourceDataStatusN.String(),
			// resource status, eg: create/release/revoke
			State: types.PowerStateRelease.String(),
			// unit: byte
			TotalMem: types.GetDefaultResoueceMem(),  // todo 使用 默认的资源大小
			// unit: byte
			UsedMem: 0,
			// number of cpu cores.
			TotalProcessor: types.GetDefaultResoueceProcessor(),  // todo 使用 默认的资源大小
			UsedProcessor:  0,
			// unit: byte
			TotalBandWidth: types.GetDefaultResoueceBandwidth(),    // todo 使用 默认的资源大小
			UsedBandWidth:  0,
		})); nil != err {
			log.Errorf("Failed to store power to local on MessageHandler, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.PowerId, power.JobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to store power to local on MessageHandler,  powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.PowerId, power.JobNodeId, err))
			continue
		}

		// 发布到全网
		if err := m.dataCenter.InsertResource(types.NewResource(&libTypes.ResourceData{
			Identity:  identity.IdentityId,
			NodeId:    identity.NodeId,
			NodeName:  identity.Name,
			DataId:    power.PowerId,
			// the status of data, N means normal, D means deleted.
			DataStatus: types.ResourceDataStatusN.String(),
			// resource status, eg: create/release/revoke
			State: types.PowerStateRelease.String(),
			// unit: byte
			TotalMem: types.GetDefaultResoueceMem(),  // todo 使用 默认的资源大小
			// unit: byte
			UsedMem: 0,
			// number of cpu cores.
			TotalProcessor: types.GetDefaultResoueceProcessor(),  // todo 使用 默认的资源大小
			UsedProcessor:  0,
			// unit: byte
			TotalBandWidth: types.GetDefaultResoueceBandwidth(),    // todo 使用 默认的资源大小
			UsedBandWidth:  0,
		})); nil != err {
			log.Errorf("Failed to store power to dataCenter on MessageHandler,  powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.PowerId, power.JobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to store power to dataCenter on MessageHandler,  powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.PowerId, power.JobNodeId, err))
			continue
		}

	}
	if len(errs) != 0 {
		return fmt.Errorf("broadcast powerMsgs errs: \n%s", strings.Join(errs, "\n"))
	}
	return nil
}

func (m *MessageHandler) BroadcastPowerRevokeMsgs(powerRevokeMsgs types.PowerRevokeMsgs) error {

	identity, err := m.dataCenter.GetIdentity()
	if nil != err {
		return fmt.Errorf("Failed to broadcast powerRevokeMsgs, query local identityInfo failed, {%s}", err)
	}

	errs := make([]string, 0)
	for _, revoke := range powerRevokeMsgs {

		jobNodeId, err := m.dataCenter.QueryLocalResourceIdByPowerId(revoke.PowerId)
		if nil != err {
			log.Errorf("Failed to QueryLocalResourceIdByPowerId on MessageHandler, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to QueryLocalResourceIdByPowerId on MessageHandler, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err))
			continue
		}
		if err := m.dataCenter.RemoveLocalResourceIdByPowerId(revoke.PowerId); nil != err {
			log.Errorf("Failed to RemoveLocalResourceIdByPowerId on MessageHandler, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to RemoveLocalResourceIdByPowerId on MessageHandler, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err))
			continue
		}
		if err := m.dataCenter.RemoveLocalResourceTable(jobNodeId); nil != err {
			log.Errorf("Failed to RemoveLocalResourceTable on MessageHandler, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to RemoveLocalResourceTable on MessageHandler, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err))
			continue
		}

		if err := m.dataCenter.RemoveLocalResource(jobNodeId); nil != err {
			log.Errorf("Failed to RemoveLocalResource on MessageHandler, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to RemoveLocalResource on MessageHandler,  powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err))
			continue
		}


		if err := m.dataCenter.InsertResource(types.NewResource(&libTypes.ResourceData{
			Identity:  identity.IdentityId,
			NodeId:    identity.NodeId,
			NodeName:  identity.Name,
			DataId:    revoke.PowerId,
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
			UsedProcessor:  0,
			// unit: byte
			TotalBandWidth: 0,
			UsedBandWidth:  0,
		})); nil != err {
			log.Error("Failed to remove dataCenter resource on MessageHandler, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to remove dataCenter resource on MessageHandler, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err))
			continue
		}
	}
	if len(errs) != 0 {
		return fmt.Errorf("broadcast powerRevokeMsgs errs: \n%s", strings.Join(errs, "\n"))
	}
	return nil
}

func (m *MessageHandler) BroadcastMetaDataMsgs(metaDataMsgs types.MetaDataMsgs) error {
	errs := make([]string, 0)
	for _, metaData := range metaDataMsgs {

		// 维护本地 数据服务的 orginId  和 metaDataId 关系
		dataResourceFileUpload, err := m.dataCenter.QueryDataResourceFileUpload(metaData.OriginId())
		if nil != err {
			log.Errorf("Failed to QueryDataResourceFileUpload on MessageHandler, originId: {%s}, metaDataId: {%s}, err: {%s}",
				metaData.OriginId(), metaData.MetaDataId, err)
			errs = append(errs, fmt.Sprintf("failed to QueryDataResourceFileUpload on MessageHandler, originId: {%s}, metaDataId: {%s}, err: {%s}",
				metaData.OriginId(), metaData.MetaDataId, err))
			continue
		}
		dataResourceFileUpload.SetMetaDataId(metaData.MetaDataId)
		if err := m.dataCenter.StoreDataResourceFileUpload(dataResourceFileUpload); nil != err {
			log.Errorf("Failed to StoreDataResourceFileUpload on MessageHandler, originId: {%s}, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metaData.OriginId(),metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to StoreDataResourceFileUpload on MessageHandler, originId: {%s}, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metaData.OriginId(),metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), err))
			continue
		}
		// 记录原始数据占用资源大小
		dataResourceTable, err := m.dataCenter.QueryDataResourceTable(dataResourceFileUpload.GetNodeId())
		if nil != err {
			log.Errorf("Failed to QueryDataResourceTable on MessageHandler, originId: {%s}, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metaData.OriginId(),metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to QueryDataResourceTable on MessageHandler, originId: {%s}, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metaData.OriginId(),metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), err))
			continue
		}
		dataResourceTable.UseDisk(uint64(metaData.Size()))
		if err := m.dataCenter.StoreDataResourceTable(dataResourceTable); nil != err {
			log.Errorf("Failed to StoreDataResourceTable on MessageHandler, originId: {%s}, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metaData.OriginId(),metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to StoreDataResourceTable on MessageHandler, originId: {%s}, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metaData.OriginId(),metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), err))
			continue
		}
		// 单独记录 metaData 的 Size 和所在 dataNodeId
		if err := m.dataCenter.StoreDataResourceDiskUsed(types.NewDataResourceDiskUsed(
			metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), uint64(metaData.Size()))); nil != err {
			log.Errorf("Failed to StoreDataResourceDiskUsed on MessageHandler, originId: {%s}, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metaData.OriginId(),metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to StoreDataResourceDiskUsed on MessageHandler, originId: {%s}, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metaData.OriginId(),metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), err))
			continue
		}

		if err := m.dataCenter.InsertMetadata(metaData.ToDataCenter()); nil != err {
			log.Errorf("Failed to store metaData to dataCenter on MessageHandler, err:", err)
			errs = append(errs, fmt.Sprintf("originId: %s, %s", metaData.OriginId(), err))
			continue
		}
	}
	if len(errs) != 0 {
		return fmt.Errorf("broadcast metaDataMsgs errs: \n%s", strings.Join(errs, "\n"))
	}
	return nil
}

func (m *MessageHandler) BroadcastMetaDataRevokeMsgs(metaDataRevokeMsgs types.MetaDataRevokeMsgs) error {
	errs := make([]string, 0)
	for _, revoke := range metaDataRevokeMsgs {
		// 需要将 dataNode 的 disk 使用信息 加回来 ...
		dataResourceDiskUsed, err := m.dataCenter.QueryDataResourceDiskUsed(revoke.MetaDataId)
		if nil != err {
			log.Errorf("Failed to QueryDataResourceDiskUsed on MessageHandler, metaDataId: {%s}, err: {%s}",
				revoke.MetaDataId, err)
			errs = append(errs, fmt.Sprintf("failed to QueryDataResourceDiskUsed on MessageHandler, metaDataId: {%s}, err: {%s}",
				revoke.MetaDataId, err))
			continue
		}
		// 记录原始数据占用资源大小
		dataResourceTable, err := m.dataCenter.QueryDataResourceTable(dataResourceDiskUsed.GetNodeId())
		if nil != err {
			log.Errorf("Failed to QueryDataResourceTable on MessageHandler, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.MetaDataId, dataResourceDiskUsed.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to QueryDataResourceTable on MessageHandler, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.MetaDataId, dataResourceDiskUsed.GetNodeId(), err))
			continue
		}
		dataResourceTable.FreeDisk(dataResourceDiskUsed.GetDiskUsed())
		if err := m.dataCenter.StoreDataResourceTable(dataResourceTable); nil != err {
			log.Errorf("Failed to StoreDataResourceTable on MessageHandler, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.MetaDataId, dataResourceDiskUsed.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to StoreDataResourceTable on MessageHandler, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.MetaDataId, dataResourceDiskUsed.GetNodeId(), err))
			continue
		}

		// 移除 metaData 的 Size 和所在 dataNodeId 的单条记录
		if err := m.dataCenter.RemoveDataResourceDiskUsed(revoke.MetaDataId); nil != err {
			log.Errorf("Failed to RemoveDataResourceDiskUsed on MessageHandler, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.MetaDataId, dataResourceDiskUsed.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to RemoveDataResourceDiskUsed on MessageHandler, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.MetaDataId, dataResourceDiskUsed.GetNodeId(), err))
			continue
		}

		if err := m.dataCenter.InsertMetadata(revoke.ToDataCenter()); nil != err {
			log.Errorf("Failed to store metaData to dataCenter on MessageHandler,metaDataId: {%s}, err: {%s}",
				revoke.MetaDataId, err)
			errs = append(errs, fmt.Sprintf("failed to store metaData to dataCenter on MessageHandler,metaDataId: {%s}, err: {%s}",
				revoke.MetaDataId, err))
			continue
		}
	}
	if len(errs) != 0 {
		return fmt.Errorf("broadcast metaDataRevokeMsgs errs: \n%s", strings.Join(errs, "\n"))
	}
	return nil
}

func (m *MessageHandler) BroadcastTaskMsgs(taskMsgs types.TaskMsgs) error {
	return m.taskManager.SendTaskMsgs(taskMsgs)
}
