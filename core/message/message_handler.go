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
	defaultPowerMsgsCacheSize    = 3
	defaultMetaDataMsgsCacheSize = 3
	defaultTaskMsgsCacheSize     = 5

	defaultBroadcastPowerMsgInterval    = 30 * time.Second
	defaultBroadcastMetaDataMsgInterval = 30 * time.Second
	defaultBroadcastTaskMsgInterval     = 10 * time.Second
)


type MessageHandler struct {
	pool       *Mempool
	dataCenter iface.ForHandleDB

	// Send taskMsg to taskManager
	taskManager *task.Manager

	msgChannel chan *feed.Event

	msgSub               event.Subscription

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
	powerTicker := time.NewTicker(defaultBroadcastPowerMsgInterval)
	metaDataTicker := time.NewTicker(defaultBroadcastMetaDataMsgInterval)
	taskTicker := time.NewTicker(defaultBroadcastTaskMsgInterval)

	for {
		select {
		case event := <-m.msgChannel:
			switch event.Type {
			case types.ApplyIdentity:
				eventMessage := event.Data.(*types.IdentityMsgEvent)
				if err := m.BroadcastIdentityMsg(eventMessage.Msg); nil != err {
					log.Errorf("Failed to call `BroadcastIdentityMsg` on MessageHandler, %s", err)
				}
			case types.RevokeIdentity:
				if err := m.BroadcastIdentityRevokeMsg(); nil != err {
					log.Errorf("Failed to call `BroadcastIdentityRevokeMsg` on MessageHandler, %s", err)
				}
			case types.ApplyPower:
				msg := event.Data.(*types.PowerMsgEvent)
				m.lockPower.Lock()
				m.powerMsgCache = append(m.powerMsgCache, msg.Msgs...)
				if len(m.powerMsgCache) >= defaultPowerMsgsCacheSize {
					if err := m.BroadcastPowerMsgs(m.powerMsgCache); nil != err {
						log.Error(fmt.Sprintf("Failed to call `BroadcastPowerMsgs` on MessageHandler, %s", err))
					}
					m.powerMsgCache = make(types.PowerMsgs, 0)
				}
				m.lockPower.Unlock()
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
						log.Errorf("Failed to call `BroadcastPowerRevokeMsgs` on MessageHandler, %s", err)
					}
				}
			case types.ApplyMetadata:
				eventMessage := event.Data.(*types.MetaDataMsgEvent)
				m.lockMetaData.Lock()
				m.metaDataMsgCache = append(m.metaDataMsgCache, eventMessage.Msgs...)
				if len(m.metaDataMsgCache) >= defaultMetaDataMsgsCacheSize {
					if err := m.BroadcastMetaDataMsgs(m.metaDataMsgCache); nil != err {
						log.Errorf("Failed to call `BroadcastMetaDataMsgs` on MessageHandler, %s", err)
					}
					m.metaDataMsgCache = make(types.MetaDataMsgs, 0)
				}
				m.lockMetaData.Unlock()
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
						log.Errorf("Failed to call `BroadcastMetaDataRevokeMsgs` on MessageHandler, %s", err)
					}
				}
			case types.ApplyTask:
				eventMessage := event.Data.(*types.TaskMsgEvent)
				m.taskMsgCache = append(m.taskMsgCache, eventMessage.Msgs...)
				if len(m.taskMsgCache) >= defaultTaskMsgsCacheSize {
					if err := m.BroadcastTaskMsgs(m.taskMsgCache); nil != err {
						log.Errorf("Failed to call `BroadcastTaskMsgs` on MessageHandler, %s", err)
					}
					m.taskMsgCache = make(types.TaskMsgs, 0)
				}
			}
		case <-powerTicker.C:

			if len(m.powerMsgCache) > 0 {
				if err := m.BroadcastPowerMsgs(m.powerMsgCache); nil != err {
					log.Errorf("Failed to call `BroadcastPowerMsgs` on MessageHandler with timer, %s", err)
				}
				m.powerMsgCache = make(types.PowerMsgs, 0)
			}

		case <-metaDataTicker.C:

			if len(m.metaDataMsgCache) > 0 {
				if err := 	m.BroadcastMetaDataMsgs(m.metaDataMsgCache); nil != err {
					log.Errorf("Failed to call `BroadcastMetaDataMsgs` on MessageHandler with timer, %s", err)
				}
				m.metaDataMsgCache = make(types.MetaDataMsgs, 0)
			}

		case <-taskTicker.C:

			if len(m.taskMsgCache) > 0 {
				if err := m.BroadcastTaskMsgs(m.taskMsgCache); nil != err {
					log.Errorf("Failed to call `BroadcastTaskMsgs` on MessageHandler with timer, %s", err)
				}
				m.taskMsgCache = make(types.TaskMsgs, 0)
			}

		// Err() channel will be closed when unsubscribing.
		case err := <-m.msgSub.Err():
			log.Errorf("Received err from msgSub, return loop, err: %s", err)
			return
		}
	}
}

func (m *MessageHandler) BroadcastIdentityMsg(msg *types.IdentityMsg) error {

	// add identity to local db
	if err := m.dataCenter.StoreIdentity(msg.NodeAlias); nil != err {
		log.Errorf("Failed to store local org identity on MessageHandler with broadcast, identityId: {%s}, err: {%s}", msg.IdentityId, err)
		return err
	}

	// send identity to datacenter
	if err := m.dataCenter.InsertIdentity(msg.ToDataCenter()); nil != err {
		log.Errorf("Failed to broadcast org org identity on MessageHandler with broadcast, identityId: {%s}, nodeId: {%s}, nodeName: {%s}, err: {%s}", msg.IdentityId, msg.NodeId, msg.Name, err)
		return err
	}
	log.Debugf("Registered identity succeed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}", msg.IdentityId, msg.NodeId, msg.Name)
	return nil
}

func (m *MessageHandler) BroadcastIdentityRevokeMsg() error {

	// remove identity from local db
	identity, err := m.dataCenter.GetIdentity()
	if nil != err {
		log.Errorf("Failed to get local org identity on MessageHandler with revoke, identityId: {%s}, err: {%s}", identity, err)
		return fmt.Errorf("query local identity failed, %s", err)
	}
	if err := m.dataCenter.RemoveIdentity(); nil != err {
		log.Errorf("Failed to delete org identity to local on MessageHandler with revoke, identityId: {%s}, err: {%s}", identity, err)
		return err
	}

	// remove identity from dataCenter
	if err := m.dataCenter.RevokeIdentity(
		types.NewIdentity(&libTypes.IdentityData{
			NodeName: identity.Name,
			NodeId:   identity.NodeId,
			Identity: identity.IdentityId,
		})); nil != err {
		log.Errorf("Failed to remove org identity to remote on MessageHandler with revoke, identityId: {%s}, err: {%s}", identity, err)
		return err
	}
	log.Debugf("Revoke identity succeed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}", identity.IdentityId, identity.NodeId, identity.Name)
	return nil
}

func (m *MessageHandler) BroadcastPowerMsgs(powerMsgs types.PowerMsgs) error {

	identity, err := m.dataCenter.GetIdentity()
	if nil != err {
		return fmt.Errorf("query local identityInfo failed, {%s}", err)
	}

	errs := make([]string, 0)

	slotUnit, err := m.dataCenter.QueryNodeResourceSlotUnit()
	if nil != err {
		return fmt.Errorf("query local slotUnit failed, {%s}", err)
	}

	for _, power := range powerMsgs {
		// 存储本地的 资源信息

		resourceTable := types.NewLocalResourceTable(
			power.JobNodeId,
			power.PowerId,
			types.GetDefaultResoueceMem(),
			types.GetDefaultResoueceProcessor(),
			types.GetDefaultResoueceBandwidth(),
			)
		resourceTable.SetSlotUnit(slotUnit)

		if err := m.dataCenter.StoreLocalResourceTable(resourceTable); nil != err {
			log.Errorf("Failed to StoreLocalResourceTable on MessageHandler with broadcast, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.PowerId, power.JobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to StoreLocalResourceTable on MessageHandler with broadcast, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.PowerId, power.JobNodeId, err))
			continue
		}


		if err := m.dataCenter.StoreLocalResourceIdByPowerId(power.PowerId, power.JobNodeId); nil != err {
			log.Errorf("Failed to StoreLocalResourceIdByPowerId on MessageHandler with broadcast,  powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.PowerId, power.JobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to StoreLocalResourceIdByPowerId on MessageHandler with broadcast,  powerId: {%s}, jobNodeId: {%s}, err: {%s}",
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
			DataStatus: types.DataStatusNormal.String(),
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
			log.Errorf("Failed to store power to local on MessageHandler with broadcast, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.PowerId, power.JobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to store power to local on MessageHandler with broadcast,  powerId: {%s}, jobNodeId: {%s}, err: {%s}",
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
			DataStatus: types.DataStatusNormal.String(),
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
			log.Errorf("Failed to store power to dataCenter on MessageHandler with broadcast, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.PowerId, power.JobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to store power to dataCenter on MessageHandler with broadcast,  powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.PowerId, power.JobNodeId, err))
			continue
		}

		log.Debugf("broadcast power msg succeed, powerId: {%s}, jobNodeId: {%s}", power.PowerId, power.JobNodeId)

	}
	if len(errs) != 0 {
		return fmt.Errorf("broadcast powerMsgs errs: \n%s", strings.Join(errs, "\n"))
	}
	return nil
}

func (m *MessageHandler) BroadcastPowerRevokeMsgs(powerRevokeMsgs types.PowerRevokeMsgs) error {

	identity, err := m.dataCenter.GetIdentity()
	if nil != err {
		return fmt.Errorf("failed to broadcast powerRevokeMsgs, query local identityInfo failed, {%s}", err)
	}

	errs := make([]string, 0)
	for _, revoke := range powerRevokeMsgs {

		jobNodeId, err := m.dataCenter.QueryLocalResourceIdByPowerId(revoke.PowerId)
		if nil != err {
			log.Errorf("Failed to QueryLocalResourceIdByPowerId on MessageHandler with revoke, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to QueryLocalResourceIdByPowerId on MessageHandler with revoke, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err))
			continue
		}
		if err := m.dataCenter.RemoveLocalResourceIdByPowerId(revoke.PowerId); nil != err {
			log.Errorf("Failed to RemoveLocalResourceIdByPowerId on MessageHandler with revoke, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to RemoveLocalResourceIdByPowerId on MessageHandler with revoke, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err))
			continue
		}
		if err := m.dataCenter.RemoveLocalResourceTable(jobNodeId); nil != err {
			log.Errorf("Failed to RemoveLocalResourceTable on MessageHandler with revoke, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to RemoveLocalResourceTable on MessageHandler with revoke, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err))
			continue
		}

		if err := m.dataCenter.RemoveLocalResource(jobNodeId); nil != err {
			log.Errorf("Failed to RemoveLocalResource on MessageHandler with revoke, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to RemoveLocalResource on MessageHandler with revoke,  powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err))
			continue
		}


		if err := m.dataCenter.RevokeResource(types.NewResource(&libTypes.ResourceData{
			Identity:  identity.IdentityId,
			NodeId:    identity.NodeId,
			NodeName:  identity.Name,
			DataId:    revoke.PowerId,
			// the status of data, N means normal, D means deleted.
			DataStatus: types.DataStatusDeleted.String(),
			// resource status, eg: create/release/revoke
			State: types.PowerStateRevoke.String(),
		})); nil != err {
			log.Errorf("Failed to remove dataCenter resource on MessageHandler with revoke, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to remove dataCenter resource on MessageHandler with revoke, jobNodeId: {%s}, err: {%s}",
				revoke.PowerId, jobNodeId, err))
			continue
		}

		log.Debugf("revoke power msg succeed, powerId: {%s}, jobNodeId: {%s}", revoke.PowerId, jobNodeId)
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
			log.Errorf("Failed to QueryDataResourceFileUpload on MessageHandler with broadcast, originId: {%s}, metaDataId: {%s}, err: {%s}",
				metaData.OriginId(), metaData.MetaDataId, err)
			errs = append(errs, fmt.Sprintf("failed to QueryDataResourceFileUpload on MessageHandler with broadcast, originId: {%s}, metaDataId: {%s}, err: {%s}",
				metaData.OriginId(), metaData.MetaDataId, err))
			continue
		}
		dataResourceFileUpload.SetMetaDataId(metaData.MetaDataId)
		if err := m.dataCenter.StoreDataResourceFileUpload(dataResourceFileUpload); nil != err {
			log.Errorf("Failed to StoreDataResourceFileUpload on MessageHandler with broadcast, originId: {%s}, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metaData.OriginId(),metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to StoreDataResourceFileUpload on MessageHandler with broadcast, originId: {%s}, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metaData.OriginId(),metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), err))
			continue
		}
		// 记录原始数据占用资源大小
		dataResourceTable, err := m.dataCenter.QueryDataResourceTable(dataResourceFileUpload.GetNodeId())
		if nil != err {
			log.Errorf("Failed to QueryDataResourceTable on MessageHandler with broadcast, originId: {%s}, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metaData.OriginId(),metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to QueryDataResourceTable on MessageHandler with broadcast, originId: {%s}, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metaData.OriginId(),metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), err))
			continue
		}
		dataResourceTable.UseDisk(uint64(metaData.Size()))
		if err := m.dataCenter.StoreDataResourceTable(dataResourceTable); nil != err {
			log.Errorf("Failed to StoreDataResourceTable on MessageHandler with broadcast, originId: {%s}, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metaData.OriginId(),metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to StoreDataResourceTable on MessageHandler with broadcast, originId: {%s}, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metaData.OriginId(),metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), err))
			continue
		}
		// 单独记录 metaData 的 Size 和所在 dataNodeId
		if err := m.dataCenter.StoreDataResourceDiskUsed(types.NewDataResourceDiskUsed(
			metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), uint64(metaData.Size()))); nil != err {
			log.Errorf("Failed to StoreDataResourceDiskUsed on MessageHandler with broadcast, originId: {%s}, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metaData.OriginId(),metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to StoreDataResourceDiskUsed on MessageHandler with broadcast, originId: {%s}, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metaData.OriginId(),metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), err))
			continue
		}

		if err := m.dataCenter.InsertMetadata(metaData.ToDataCenter()); nil != err {
			log.Errorf("Failed to store metaData to dataCenter on MessageHandler with broadcast, originId: {%s}, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metaData.OriginId(),metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to store metaData to dataCenter on MessageHandler with broadcast, originId: {%s}, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metaData.OriginId(),metaData.MetaDataId, dataResourceFileUpload.GetNodeId(), err))
			continue
		}

		log.Debugf("broadcast metaData msg succeed, originId: {%s}, metaDataId: {%s}", metaData.OriginId(), metaData.MetaDataId)
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
			log.Errorf("Failed to QueryDataResourceDiskUsed on MessageHandler with revoke, metaDataId: {%s}, err: {%s}",
				revoke.MetaDataId, err)
			errs = append(errs, fmt.Sprintf("failed to QueryDataResourceDiskUsed on MessageHandler with revoke, metaDataId: {%s}, err: {%s}",
				revoke.MetaDataId, err))
			continue
		}
		// 记录原始数据占用资源大小
		dataResourceTable, err := m.dataCenter.QueryDataResourceTable(dataResourceDiskUsed.GetNodeId())
		if nil != err {
			log.Errorf("Failed to QueryDataResourceTable on MessageHandler with revoke, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.MetaDataId, dataResourceDiskUsed.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to QueryDataResourceTable on MessageHandler with revoke, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.MetaDataId, dataResourceDiskUsed.GetNodeId(), err))
			continue
		}
		dataResourceTable.FreeDisk(dataResourceDiskUsed.GetDiskUsed())
		if err := m.dataCenter.StoreDataResourceTable(dataResourceTable); nil != err {
			log.Errorf("Failed to StoreDataResourceTable on MessageHandler with revoke, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.MetaDataId, dataResourceDiskUsed.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to StoreDataResourceTable on MessageHandler with revoke, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.MetaDataId, dataResourceDiskUsed.GetNodeId(), err))
			continue
		}

		// 移除 metaData 的 Size 和所在 dataNodeId 的单条记录
		if err := m.dataCenter.RemoveDataResourceDiskUsed(revoke.MetaDataId); nil != err {
			log.Errorf("Failed to RemoveDataResourceDiskUsed on MessageHandler with revoke, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.MetaDataId, dataResourceDiskUsed.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to RemoveDataResourceDiskUsed on MessageHandler with revoke, metaDataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.MetaDataId, dataResourceDiskUsed.GetNodeId(), err))
			continue
		}

		if err := m.dataCenter.RevokeMetadata(revoke.ToDataCenter()); nil != err {
			log.Errorf("Failed to store metaData to dataCenter on MessageHandler with revoke, metaDataId: {%s}, err: {%s}",
				revoke.MetaDataId, err)
			errs = append(errs, fmt.Sprintf("failed to store metaData to dataCenter on MessageHandler with revoke, metaDataId: {%s}, err: {%s}",
				revoke.MetaDataId, err))
			continue
		}
		log.Debugf("revoke metaData msg succeed, metaDataId: {%s}", revoke.MetaDataId)
	}
	if len(errs) != 0 {
		return fmt.Errorf("broadcast metaDataRevokeMsgs errs: \n%s", strings.Join(errs, "\n"))
	}
	return nil
}

func (m *MessageHandler) BroadcastTaskMsgs(taskMsgs types.TaskMsgs) error {
	return m.taskManager.SendTaskMsgs(taskMsgs)
}
