package message

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/feed"
	"github.com/RosettaFlow/Carrier-Go/core/iface"
	"github.com/RosettaFlow/Carrier-Go/core/task"
	"github.com/RosettaFlow/Carrier-Go/event"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
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
	quit       chan struct{}

	msgSub event.Subscription

	powerMsgCache    types.PowerMsgArr
	metadataMsgCache types.MetadataMsgArr
	taskMsgCache     types.TaskMsgArr

	lockPower    sync.Mutex
	lockMetadata sync.Mutex

	// TODO 有些缓存需要持久化
}

func NewHandler(pool *Mempool, dataCenter iface.ForHandleDB, taskManager *task.Manager) *MessageHandler {
	m := &MessageHandler{
		pool:        pool,
		dataCenter:  dataCenter,
		taskManager: taskManager,
		msgChannel:  make(chan *feed.Event, 5),
		quit:        make(chan struct{}),
	}
	return m
}

func (m *MessageHandler) Start() error {
	m.msgSub = m.pool.SubscribeNewMessageEvent(m.msgChannel)
	go m.loop()
	log.Info("Started message handler ...")
	return nil
}
func (m *MessageHandler) Stop() error {
	close(m.quit)
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
					if err := m.BroadcastPowerMsgArr(m.powerMsgCache); nil != err {
						log.Error(fmt.Sprintf("Failed to call `BroadcastPowerMsgArr` on MessageHandler, %s", err))
					}
					m.powerMsgCache = make(types.PowerMsgArr, 0)
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
					msgs, index := make(types.PowerRevokeMsgArr, len(tmp)), 0
					for _, i := range tmp {
						msgs[index] = eventMessage.Msgs[i]
						index++
					}
					if err := m.BroadcastPowerRevokeMsgArr(msgs); nil != err {
						log.Errorf("Failed to call `BroadcastPowerRevokeMsgArr` on MessageHandler, %s", err)
					}
				}
			case types.ApplyMetadata:
				eventMessage := event.Data.(*types.MetadataMsgEvent)
				m.lockMetadata.Lock()
				m.metadataMsgCache = append(m.metadataMsgCache, eventMessage.Msgs...)
				if len(m.metadataMsgCache) >= defaultMetaDataMsgsCacheSize {
					if err := m.BroadcastMetadataMsgArr(m.metadataMsgCache); nil != err {
						log.Errorf("Failed to call `BroadcastMetadataMsgArr` on MessageHandler, %s", err)
					}
					m.metadataMsgCache = make(types.MetadataMsgArr, 0)
				}
				m.lockMetadata.Unlock()
			case types.RevokeMetadata:
				eventMessage := event.Data.(*types.MetadataRevokeMsgEvent)
				tmp := make(map[string]int, len(eventMessage.Msgs))
				for i, msg := range eventMessage.Msgs {
					tmp[msg.MetadataId] = i
				}

				// Remove local cache metaDataMsgs
				m.lockMetadata.Lock()
				for i := 0; i < len(m.metadataMsgCache); i++ {
					msg := m.metadataMsgCache[i]
					if _, ok := tmp[msg.MetadataId]; ok {
						delete(tmp, msg.MetadataId)
						m.metadataMsgCache = append(m.metadataMsgCache[:i], m.metadataMsgCache[i+1:]...)
						i--
					}
				}
				m.lockMetadata.Unlock()

				// Revoke remote metaData
				if len(tmp) != 0 {
					msgs, index := make(types.MetadataRevokeMsgArr, len(tmp)), 0
					for _, i := range tmp {
						msgs[index] = eventMessage.Msgs[i]
						index++
					}
					if err := m.BroadcastMetaDataRevokeMsgArr(msgs); nil != err {
						log.Errorf("Failed to call `BroadcastMetaDataRevokeMsgArr` on MessageHandler, %s", err)
					}
				}
			case types.ApplyTask:
				eventMessage := event.Data.(*types.TaskMsgEvent)
				m.taskMsgCache = append(m.taskMsgCache, eventMessage.Msgs...)
				if len(m.taskMsgCache) >= defaultTaskMsgsCacheSize {
					if err := m.BroadcastTaskMsgArr(m.taskMsgCache); nil != err {
						log.Errorf("Failed to call `BroadcastTaskMsgArr` on MessageHandler, %s", err)
					}
					m.taskMsgCache = make(types.TaskMsgArr, 0)
				}
			}
		case <-powerTicker.C:

			if len(m.powerMsgCache) > 0 {
				if err := m.BroadcastPowerMsgArr(m.powerMsgCache); nil != err {
					log.Errorf("Failed to call `BroadcastPowerMsgArr` on MessageHandler with timer, %s", err)
				}
				m.powerMsgCache = make(types.PowerMsgArr, 0)
			}

		case <-metaDataTicker.C:

			if len(m.metadataMsgCache) > 0 {
				if err := m.BroadcastMetadataMsgArr(m.metadataMsgCache); nil != err {
					log.Errorf("Failed to call `BroadcastMetadataMsgArr` on MessageHandler with timer, %s", err)
				}
				m.metadataMsgCache = make(types.MetadataMsgArr, 0)
			}

		case <-taskTicker.C:

			if len(m.taskMsgCache) > 0 {
				if err := m.BroadcastTaskMsgArr(m.taskMsgCache); nil != err {
					log.Errorf("Failed to call `BroadcastTaskMsgArr` on MessageHandler with timer, %s", err)
				}
				m.taskMsgCache = make(types.TaskMsgArr, 0)
			}

		// Err() channel will be closed when unsubscribing.
		case err := <-m.msgSub.Err():
			log.Errorf("Received err from msgSub, return loop, err: %s", err)
			return
		case <-m.quit:
			log.Infof("Stopped message handler ...")
			return
		}
	}
}

func (m *MessageHandler) BroadcastIdentityMsg(msg *types.IdentityMsg) error {

	// add identity to local db
	if err := m.dataCenter.StoreIdentity(msg.GetOrganization()); nil != err {
		log.Errorf("Failed to store local org identity on MessageHandler with broadcast, identityId: {%s}, err: {%s}", msg.GetOwnerIdentityId(), err)
		return err
	}

	// send identity to datacenter
	if err := m.dataCenter.InsertIdentity(msg.ToDataCenter()); nil != err {
		log.Errorf("Failed to broadcast org org identity on MessageHandler with broadcast, identityId: {%s}, nodeId: {%s}, nodeName: {%s}, err: {%s}",
			msg.GetOwnerIdentityId(), msg.GetOwnerNodeId(), msg.GetOwnerName(), err)
		return err
	}
	log.Debugf("Registered identity succeed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}", msg.GetOwnerIdentityId(), msg.GetOwnerNodeId(), msg.GetOwnerName())
	return nil
}

func (m *MessageHandler) BroadcastIdentityRevokeMsg() error {

	// remove identity from local db
	identity, err := m.dataCenter.GetIdentity()
	if nil != err {
		log.Errorf("Failed to get local org identity on MessageHandler with revoke, identityId: {%s}, err: {%s}", identity.GetIdentityId(), err)
		return fmt.Errorf("query local identity failed, %s", err)
	}
	if err := m.dataCenter.RemoveIdentity(); nil != err {
		log.Errorf("Failed to delete org identity to local on MessageHandler with revoke, identityId: {%s}, err: {%s}", identity.GetIdentityId(), err)
		return err
	}

	// remove identity from dataCenter
	if err := m.dataCenter.RevokeIdentity(
		types.NewIdentity(&libtypes.IdentityPB{
			NodeName: identity.GetNodeName(),
			NodeId:   identity.GetNodeId(),
			IdentityId: identity.GetIdentityId(),
		})); nil != err {
		log.Errorf("Failed to remove org identity to remote on MessageHandler with revoke, identityId: {%s}, err: {%s}", identity.GetIdentityId(), err)
		return err
	}
	log.Debugf("Revoke identity succeed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}", identity.GetIdentityId(), identity.GetNodeId(), identity.GetNodeName())
	return nil
}

func (m *MessageHandler) BroadcastPowerMsgArr(powerMsgArr types.PowerMsgArr) error {

	identity, err := m.dataCenter.GetIdentity()
	if nil != err {
		return fmt.Errorf("query local identityInfo failed, {%s}", err)
	}

	errs := make([]string, 0)

	slotUnit, err := m.dataCenter.QueryNodeResourceSlotUnit()
	if nil != err {
		return fmt.Errorf("query local slotUnit failed, {%s}", err)
	}

	for _, power := range powerMsgArr {
		// 存储本地的 资源信息

		resourceTable := types.NewLocalResourceTable(
			power.GetJobNodeId(),
			power.GetPowerId(),
			types.GetDefaultResoueceMem(),
			types.GetDefaultResoueceBandwidth(),
			types.GetDefaultResoueceProcessor(),
		)
		resourceTable.SetSlotUnit(slotUnit)

		log.Debugf("Publish power, StoreLocalResourceTable, %s", resourceTable.String())
		if err := m.dataCenter.StoreLocalResourceTable(resourceTable); nil != err {
			log.Errorf("Failed to StoreLocalResourceTable on MessageHandler with broadcast, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.GetPowerId(), power.GetJobNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to StoreLocalResourceTable on MessageHandler with broadcast, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.GetPowerId(), power.GetJobNodeId(), err))
			continue
		}

		if err := m.dataCenter.StoreLocalResourceIdByPowerId(power.GetPowerId(), power.GetJobNodeId()); nil != err {
			log.Errorf("Failed to StoreLocalResourceIdByPowerId on MessageHandler with broadcast,  powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.GetPowerId(), power.GetJobNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to StoreLocalResourceIdByPowerId on MessageHandler with broadcast,  powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.GetPowerId(), power.GetJobNodeId(), err))
			continue
		}
		if err := m.dataCenter.InsertLocalResource(types.NewLocalResource(&libtypes.LocalResourcePB{
			IdentityId:  identity.GetIdentityId(),
			NodeId:    identity.GetNodeId(),
			NodeName:  identity.GetNodeName(),
			JobNodeId: power.GetJobNodeId(),
			DataId:    power.GetPowerId(),
			// the status of data, N means normal, D means deleted.
			DataStatus: apicommonpb.DataStatus_DataStatus_Normal,
			// resource status, eg: create/release/revoke
			State: apicommonpb.PowerState_PowerState_Released,
			// unit: byte
			TotalMem: types.GetDefaultResoueceMem(), // todo 使用 默认的资源大小
			// unit: byte
			UsedMem: 0,
			// number of cpu cores.
			TotalProcessor: uint32(types.GetDefaultResoueceProcessor()), // todo 使用 默认的资源大小
			UsedProcessor:  0,
			// unit: byte
			TotalBandwidth: types.GetDefaultResoueceBandwidth(), // todo 使用 默认的资源大小
			UsedBandwidth:  0,
		})); nil != err {
			log.Errorf("Failed to store power to local on MessageHandler with broadcast, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.GetPowerId(), power.GetJobNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to store power to local on MessageHandler with broadcast,  powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.GetPowerId(), power.GetJobNodeId(), err))
			continue
		}

		// 发布到全网
		if err := m.dataCenter.InsertResource(types.NewResource(&libtypes.ResourcePB{
			IdentityId: identity.GetIdentityId(),
			NodeId:   identity.GetNodeId(),
			NodeName: identity.GetNodeName(),
			DataId:   power.GetPowerId(),
			// the status of data, N means normal, D means deleted.
			DataStatus: apicommonpb.DataStatus_DataStatus_Normal,
			// resource status, eg: create/release/revoke
			State: apicommonpb.PowerState_PowerState_Released,
			// unit: byte
			TotalMem: types.GetDefaultResoueceMem(), // todo 使用 默认的资源大小
			// unit: byte
			UsedMem: 0,
			// number of cpu cores.
			TotalProcessor: uint32(types.GetDefaultResoueceProcessor()), // todo 使用 默认的资源大小
			UsedProcessor:  0,
			// unit: byte
			TotalBandwidth: types.GetDefaultResoueceBandwidth(), // todo 使用 默认的资源大小
			UsedBandwidth:  0,
		})); nil != err {
			log.Errorf("Failed to store power to dataCenter on MessageHandler with broadcast, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.GetPowerId(), power.GetJobNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to store power to dataCenter on MessageHandler with broadcast,  powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.GetPowerId(), power.GetJobNodeId(), err))
			continue
		}

		log.Debugf("broadcast power msg succeed, powerId: {%s}, jobNodeId: {%s}", power.GetPowerId(), power.GetJobNodeId())

	}
	if len(errs) != 0 {
		return fmt.Errorf("broadcast powerMsgArr errs: \n%s", strings.Join(errs, "\n"))
	}
	return nil
}

func (m *MessageHandler) BroadcastPowerRevokeMsgArr(powerRevokeMsgArr types.PowerRevokeMsgArr) error {

	identity, err := m.dataCenter.GetIdentity()
	if nil != err {
		return fmt.Errorf("failed to broadcast powerRevokeMsgArr, query local identityInfo failed, {%s}", err)
	}

	errs := make([]string, 0)
	for _, revoke := range powerRevokeMsgArr {

		jobNodeId, err := m.dataCenter.QueryLocalResourceIdByPowerId(revoke.GetPowerId())
		if nil != err {
			log.Errorf("Failed to QueryLocalResourceIdByPowerId on MessageHandler with revoke, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.GetPowerId(), jobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to QueryLocalResourceIdByPowerId on MessageHandler with revoke, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.GetPowerId(), jobNodeId, err))
			continue
		}
		if err := m.dataCenter.RemoveLocalResourceIdByPowerId(revoke.GetPowerId()); nil != err {
			log.Errorf("Failed to RemoveLocalResourceIdByPowerId on MessageHandler with revoke, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.GetPowerId(), jobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to RemoveLocalResourceIdByPowerId on MessageHandler with revoke, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.GetPowerId(), jobNodeId, err))
			continue
		}
		if err := m.dataCenter.RemoveLocalResourceTable(jobNodeId); nil != err {
			log.Errorf("Failed to RemoveLocalResourceTable on MessageHandler with revoke, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.GetPowerId(), jobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to RemoveLocalResourceTable on MessageHandler with revoke, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.GetPowerId(), jobNodeId, err))
			continue
		}

		if err := m.dataCenter.RemoveLocalResource(jobNodeId); nil != err {
			log.Errorf("Failed to RemoveLocalResource on MessageHandler with revoke, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.GetPowerId(), jobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to RemoveLocalResource on MessageHandler with revoke,  powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.GetPowerId(), jobNodeId, err))
			continue
		}

		if err := m.dataCenter.RevokeResource(types.NewResource(&libtypes.ResourcePB{
			IdentityId: identity.GetIdentityId(),
			NodeId:   identity.GetNodeId(),
			NodeName: identity.GetNodeName(),
			DataId:   revoke.GetPowerId(),
			// the status of data, N means normal, D means deleted.
			DataStatus: apicommonpb.DataStatus_DataStatus_Deleted,
			// resource status, eg: create/release/revoke
			State: apicommonpb.PowerState_PowerState_Revoked,
		})); nil != err {
			log.Errorf("Failed to remove dataCenter resource on MessageHandler with revoke, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.GetPowerId(), jobNodeId, err)
			errs = append(errs, fmt.Sprintf("failed to remove dataCenter resource on MessageHandler with revoke,  powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.GetPowerId(), jobNodeId, err))
			continue
		}

		log.Debugf("revoke power msg succeed, powerId: {%s}, jobNodeId: {%s}", revoke.GetPowerId(), jobNodeId)
	}
	if len(errs) != 0 {
		return fmt.Errorf("broadcast powerRevokeMsgArr errs: \n%s", strings.Join(errs, "\n"))
	}
	return nil
}

func (m *MessageHandler) BroadcastMetadataMsgArr(metadataMsgArr types.MetadataMsgArr) error {
	errs := make([]string, 0)
	for _, metadata := range metadataMsgArr {

		// 维护本地 数据服务的 orginId  和 metaDataId 关系
		dataResourceFileUpload, err := m.dataCenter.QueryDataResourceFileUpload(metadata.GetOriginId())
		if nil != err {
			log.Errorf("Failed to QueryDataResourceFileUpload on MessageHandler with broadcast, originId: {%s}, metadataId: {%s}, err: {%s}",
				metadata.GetOriginId(), metadata.GetMetadataId(), err)
			errs = append(errs, fmt.Sprintf("failed to QueryDataResourceFileUpload on MessageHandler with broadcast, originId: {%s}, metadataId: {%s}, err: {%s}",
				metadata.GetOriginId(), metadata.GetMetadataId(), err))
			continue
		}
		// 更新 fileupload 信息中的 metadataId
		dataResourceFileUpload.SetMetadataId(metadata.GetMetadataId())
		if err := m.dataCenter.StoreDataResourceFileUpload(dataResourceFileUpload); nil != err {
			log.Errorf("Failed to StoreDataResourceFileUpload on MessageHandler with broadcast, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metadata.GetOriginId(), metadata.GetMetadataId(), dataResourceFileUpload.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to StoreDataResourceFileUpload on MessageHandler with broadcast, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metadata.GetOriginId(), metadata.GetMetadataId(), dataResourceFileUpload.GetNodeId(), err))
			continue
		}
		// 记录原始数据占用资源大小
		dataResourceTable, err := m.dataCenter.QueryDataResourceTable(dataResourceFileUpload.GetNodeId())
		if nil != err {
			log.Errorf("Failed to QueryDataResourceTable on MessageHandler with broadcast, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metadata.GetOriginId(), metadata.GetMetadataId(), dataResourceFileUpload.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to QueryDataResourceTable on MessageHandler with broadcast, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metadata.GetOriginId(), metadata.GetMetadataId(), dataResourceFileUpload.GetNodeId(), err))
			continue
		}
		dataResourceTable.UseDisk(metadata.GetSize())
		if err := m.dataCenter.StoreDataResourceTable(dataResourceTable); nil != err {
			log.Errorf("Failed to StoreDataResourceTable on MessageHandler with broadcast, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metadata.GetOriginId(), metadata.GetMetadataId(), dataResourceFileUpload.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to StoreDataResourceTable on MessageHandler with broadcast, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metadata.GetOriginId(), metadata.GetMetadataId(), dataResourceFileUpload.GetNodeId(), err))
			continue
		}
		// 单独记录 metaData 的 GetSize 和所在 dataNodeId
		if err := m.dataCenter.StoreDataResourceDiskUsed(types.NewDataResourceDiskUsed(
			metadata.MetadataId, dataResourceFileUpload.GetNodeId(), metadata.GetSize())); nil != err {
			log.Errorf("Failed to StoreDataResourceDiskUsed on MessageHandler with broadcast, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metadata.GetOriginId(), metadata.GetMetadataId(), dataResourceFileUpload.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to StoreDataResourceDiskUsed on MessageHandler with broadcast, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metadata.GetOriginId(), metadata.GetMetadataId(), dataResourceFileUpload.GetNodeId(), err))
			continue
		}

		if err := m.dataCenter.InsertMetadata(metadata.ToDataCenter()); nil != err {
			log.Errorf("Failed to store metaData to dataCenter on MessageHandler with broadcast, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metadata.GetOriginId(), metadata.GetMetadataId(), dataResourceFileUpload.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to store metaData to dataCenter on MessageHandler with broadcast, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metadata.GetOriginId(), metadata.GetMetadataId(), dataResourceFileUpload.GetNodeId(), err))
			continue
		}

		log.Debugf("broadcast metaData msg succeed, originId: {%s}, metadataId: {%s}", metadata.GetOriginId(), metadata.GetMetadataId())
	}
	if len(errs) != 0 {
		return fmt.Errorf("broadcast metadataMsgs errs: \n%s", strings.Join(errs, "\n"))
	}

	return nil
}

func (m *MessageHandler) BroadcastMetaDataRevokeMsgArr(metadataRevokeMsgArr types.MetadataRevokeMsgArr) error {
	errs := make([]string, 0)
	for _, revoke := range metadataRevokeMsgArr {
		// 需要将 dataNode 的 disk 使用信息 加回来 ...
		dataResourceDiskUsed, err := m.dataCenter.QueryDataResourceDiskUsed(revoke.GetMetadataId())
		if nil != err {
			log.Errorf("Failed to QueryDataResourceDiskUsed on MessageHandler with revoke, metadataId: {%s}, err: {%s}",
				revoke.GetMetadataId(), err)
			errs = append(errs, fmt.Sprintf("failed to QueryDataResourceDiskUsed on MessageHandler with revoke, metadataId: {%s}, err: {%s}",
				revoke.GetMetadataId(), err))
			continue
		}
		// 记录原始数据占用资源大小
		dataResourceTable, err := m.dataCenter.QueryDataResourceTable(dataResourceDiskUsed.GetNodeId())
		if nil != err {
			log.Errorf("Failed to QueryDataResourceTable on MessageHandler with revoke, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.GetMetadataId(), dataResourceDiskUsed.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to QueryDataResourceTable on MessageHandler with revoke, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.GetMetadataId(), dataResourceDiskUsed.GetNodeId(), err))
			continue
		}
		dataResourceTable.FreeDisk(dataResourceDiskUsed.GetDiskUsed())
		if err := m.dataCenter.StoreDataResourceTable(dataResourceTable); nil != err {
			log.Errorf("Failed to StoreDataResourceTable on MessageHandler with revoke, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.GetMetadataId(), dataResourceDiskUsed.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to StoreDataResourceTable on MessageHandler with revoke, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.GetMetadataId(), dataResourceDiskUsed.GetNodeId(), err))
			continue
		}

		// 移除 metaData 的 GetSize 和所在 dataNodeId 的单条记录
		if err := m.dataCenter.RemoveDataResourceDiskUsed(revoke.GetMetadataId()); nil != err {
			log.Errorf("Failed to RemoveDataResourceDiskUsed on MessageHandler with revoke, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.GetMetadataId(), dataResourceDiskUsed.GetNodeId(), err)
			errs = append(errs, fmt.Sprintf("failed to RemoveDataResourceDiskUsed on MessageHandler with revoke, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.GetMetadataId(), dataResourceDiskUsed.GetNodeId(), err))
			continue
		}

		if err := m.dataCenter.RevokeMetadata(revoke.ToDataCenter()); nil != err {
			log.Errorf("Failed to store metaData to dataCenter on MessageHandler with revoke, metadataId: {%s}, err: {%s}",
				revoke.GetMetadataId(), err)
			errs = append(errs, fmt.Sprintf("failed to store metaData to dataCenter on MessageHandler with revoke, metadataId: {%s}, err: {%s}",
				revoke.GetMetadataId(), err))
			continue
		}
		log.Debugf("revoke metaData msg succeed, metadataId: {%s}", revoke.GetMetadataId())
	}
	if len(errs) != 0 {
		return fmt.Errorf("broadcast metadataRevokeMsgArr errs: \n%s", strings.Join(errs, "\n"))
	}
	return nil
}

func (m *MessageHandler) BroadcastTaskMsgArr(taskMsgArr types.TaskMsgArr) error {
	return m.taskManager.SendTaskMsgArr(taskMsgArr)
}
