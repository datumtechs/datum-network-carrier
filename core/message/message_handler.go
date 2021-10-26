package message

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/auth"
	"github.com/RosettaFlow/Carrier-Go/common/feed"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/core/iface"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/core/task"
	"github.com/RosettaFlow/Carrier-Go/event"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"sync"
	"time"
)

const (
	defaultPowerMsgsCacheSize        = 3
	defaultMetadataMsgsCacheSize     = 3
	defaultMetadataAuthMsgsCacheSize = 3
	defaultTaskMsgsCacheSize         = 5

	defaultBroadcastPowerMsgInterval        = 30 * time.Second
	defaultBroadcastMetadataMsgInterval     = 30 * time.Second
	defaultBroadcastMetadataAuthMsgInterval = 30 * time.Second
	defaultBroadcastTaskMsgInterval         = 10 * time.Second
)

type MessageHandler struct {
	pool       *Mempool
	dataCenter iface.ForHandleDB
	// Send taskMsg to taskManager
	taskManager *task.Manager
	authManager *auth.AuthorityManager
	// internal resource node set (Fighter node grpc client set)
	resourceClientSet *grpclient.InternalResourceClientSet
	msgChannel        chan *feed.Event
	quit              chan struct{}

	msgSub event.Subscription

	powerMsgCache        types.PowerMsgArr
	metadataMsgCache     types.MetadataMsgArr
	metadataAuthMsgCache types.MetadataAuthorityMsgArr
	taskMsgCache         types.TaskMsgArr

	lockPower        sync.Mutex
	lockMetadata     sync.Mutex
	lockMetadataAuth sync.Mutex
	lockTask         sync.Mutex

	// TODO 有些缓存需要持久化
}

func NewHandler(pool *Mempool, dataCenter iface.ForHandleDB, taskManager *task.Manager, authManager *auth.AuthorityManager, resourceClientSet *grpclient.InternalResourceClientSet) *MessageHandler {
	m := &MessageHandler{
		pool:              pool,
		dataCenter:        dataCenter,
		taskManager:       taskManager,
		authManager:       authManager,
		resourceClientSet: resourceClientSet,
		msgChannel:        make(chan *feed.Event, 5),
		quit:              make(chan struct{}),
	}
	return m
}

func (m *MessageHandler) Start() error {
	m.msgSub = m.pool.SubscribeNewMessageEvent(m.msgChannel)
	m.recoveryCache()
	go m.loop()
	log.Info("Started message handler ...")
	return nil
}
func (m *MessageHandler) Stop() error {
	close(m.quit)
	return nil
}

func (m *MessageHandler) recoveryCache()  {
	taskMsgCache, err := m.dataCenter.QueryTaskMsgArr()
	if nil != err {
		log.Warning("Failed to get taskMsgCache from QueryTaskMsgArr.")
	} else {
		m.taskMsgCache = taskMsgCache
	}

	metadataAuthMsgCache,err:=m.dataCenter.QueryMetadataAuthorityMsgArr()
	if nil != err {
		log.Warning("Failed to get metadataAuthMsgCache from QueryMetadataAuthorityMsgArr.")
	} else {
		m.metadataAuthMsgCache = metadataAuthMsgCache
	}

	metadataMsgCache,err:=m.dataCenter.QueryMetadataMsgArr()
	if nil != err {
		log.Warning("Failed to get metadataMsgCache from QueryMetadataMsgArr.")
	} else {
		m.metadataMsgCache=metadataMsgCache
	}

	powerMsgCache,err:=m.dataCenter.QueryPowerMsgArr()
	if nil != err {
		log.Warning("Failed to get powerMsgCache from QueryPowerMsgArr.")
	} else {
		m.powerMsgCache=powerMsgCache
	}
}
func (m *MessageHandler) loop() {
	powerTicker := time.NewTicker(defaultBroadcastPowerMsgInterval)
	metadataTicker := time.NewTicker(defaultBroadcastMetadataMsgInterval)
	metadataAuthTicker := time.NewTicker(defaultBroadcastMetadataAuthMsgInterval)
	taskTicker := time.NewTicker(defaultBroadcastTaskMsgInterval)
	for {
		select {
		case event := <-m.msgChannel:
			switch event.Type {
			case types.ApplyIdentity:
				eventMessage := event.Data.(*types.IdentityMsgEvent)
				m.BroadcastIdentityMsg(eventMessage.Msg)
			case types.RevokeIdentity:
				m.BroadcastIdentityRevokeMsg()
			case types.ApplyPower:
				msg := event.Data.(*types.PowerMsgEvent)
				m.lockPower.Lock()
				m.powerMsgCache.Len()
				m.powerMsgCache = append(m.powerMsgCache, msg.Msgs...)
				m.dataCenter.StoreMessageCache(m.powerMsgCache)
				if len(m.powerMsgCache) >= defaultPowerMsgsCacheSize {
					m.BroadcastPowerMsgArr(m.powerMsgCache)
					m.powerMsgCache = make(types.PowerMsgArr, 0)
					m.dataCenter.StoreMessageCache(m.powerMsgCache)
				}
				m.lockPower.Unlock()
			case types.RevokePower:
				eventMessage := event.Data.(*types.PowerRevokeMsgEvent)
				tmp := make(map[string]int, len(eventMessage.Msgs))
				for i, msg := range eventMessage.Msgs {
					tmp[msg.GetPowerId()] = i
				}

				// Remove local cache powerMsgs
				m.lockPower.Lock()
				for i := 0; i < len(m.powerMsgCache); i++ {
					msg := m.powerMsgCache[i]
					if _, ok := tmp[msg.GetPowerId()]; ok {
						delete(tmp, msg.GetPowerId())
						m.powerMsgCache = append(m.powerMsgCache[:i], m.powerMsgCache[i+1:]...)
						m.dataCenter.StoreMessageCache(m.powerMsgCache)
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
					m.BroadcastPowerRevokeMsgArr(msgs)
				}
			case types.ApplyMetadata:
				eventMessage := event.Data.(*types.MetadataMsgEvent)
				m.lockMetadata.Lock()
				m.metadataMsgCache = append(m.metadataMsgCache, eventMessage.Msgs...)
				m.dataCenter.StoreMessageCache(m.metadataAuthMsgCache)
				if len(m.metadataMsgCache) >= defaultMetadataMsgsCacheSize {
					m.BroadcastMetadataMsgArr(m.metadataMsgCache)
					m.metadataMsgCache = make(types.MetadataMsgArr, 0)
					m.dataCenter.StoreMessageCache(m.metadataAuthMsgCache)
				}
				m.lockMetadata.Unlock()
			case types.RevokeMetadata:
				eventMessage := event.Data.(*types.MetadataRevokeMsgEvent)
				tmp := make(map[string]int, len(eventMessage.Msgs))
				for i, msg := range eventMessage.Msgs {
					tmp[msg.GetMetadataId()] = i
				}

				// Remove local cache metadataMsgs
				m.lockMetadata.Lock()
				for i := 0; i < len(m.metadataMsgCache); i++ {
					msg := m.metadataMsgCache[i]
					if _, ok := tmp[msg.GetMetadataId()]; ok {
						delete(tmp, msg.GetMetadataId())
						m.metadataMsgCache = append(m.metadataMsgCache[:i], m.metadataMsgCache[i+1:]...)
						m.dataCenter.StoreMessageCache(m.metadataAuthMsgCache)
						i--
					}
				}
				m.lockMetadata.Unlock()

				// Revoke remote metadata
				if len(tmp) != 0 {
					msgs, index := make(types.MetadataRevokeMsgArr, len(tmp)), 0
					for _, i := range tmp {
						msgs[index] = eventMessage.Msgs[i]
						index++
					}
					m.BroadcastMetadataRevokeMsgArr(msgs)
				}

			case types.ApplyMetadataAuth:
				eventMessage := event.Data.(*types.MetadataAuthMsgEvent)
				m.lockMetadataAuth.Lock()
				m.metadataAuthMsgCache = append(m.metadataAuthMsgCache, eventMessage.Msgs...)
				m.dataCenter.StoreMessageCache(m.metadataAuthMsgCache)
				if len(m.metadataAuthMsgCache) >= defaultMetadataAuthMsgsCacheSize {
					m.BroadcastMetadataAuthMsgArr(m.metadataAuthMsgCache)
					m.metadataAuthMsgCache = make(types.MetadataAuthorityMsgArr, 0)
					m.dataCenter.StoreMessageCache(m.metadataAuthMsgCache)
				}
				m.lockMetadataAuth.Unlock()
			case types.RevokeMetadataAuth:
				eventMessage := event.Data.(*types.MetadataAuthRevokeMsgEvent)
				tmp := make(map[string]int, len(eventMessage.Msgs))
				for i, msg := range eventMessage.Msgs {
					tmp[msg.GetMetadataAuthId()] = i
				}

				// Remove local cache metadataAuthorityMsgs
				m.lockMetadataAuth.Lock()
				for i := 0; i < len(m.metadataAuthMsgCache); i++ {
					msg := m.metadataAuthMsgCache[i]
					if _, ok := tmp[msg.GetMetadataAuthId()]; ok {
						delete(tmp, msg.GetMetadataAuthId())
						m.metadataAuthMsgCache = append(m.metadataAuthMsgCache[:i], m.metadataAuthMsgCache[i+1:]...)
						m.dataCenter.StoreMessageCache(m.metadataAuthMsgCache)
						i--
					}
				}
				m.lockMetadataAuth.Unlock()

				// Revoke remote metadataAuthority
				if len(tmp) != 0 {
					msgs, index := make(types.MetadataAuthorityRevokeMsgArr, len(tmp)), 0
					for _, i := range tmp {
						msgs[index] = eventMessage.Msgs[i]
						index++
					}
					m.BroadcastMetadataAuthRevokeMsgArr(msgs)
				}

			case types.ApplyTask:
				eventMessage := event.Data.(*types.TaskMsgEvent)
				m.lockTask.Lock()
				m.taskMsgCache = append(m.taskMsgCache, eventMessage.Msgs...)
				m.dataCenter.StoreMessageCache(m.taskMsgCache)
				if len(m.taskMsgCache) >= defaultTaskMsgsCacheSize {
					m.BroadcastTaskMsgArr(m.taskMsgCache)
					m.taskMsgCache = make(types.TaskMsgArr, 0)
					m.dataCenter.StoreMessageCache(m.taskMsgCache)
				}
				m.lockTask.Unlock()
			case types.TerminateTask:
				eventMessage := event.Data.(*types.TaskTerminateMsgEvent)
				tmp := make(map[string]int, len(eventMessage.Msgs))
				for i, msg := range eventMessage.Msgs {
					tmp[msg.GetTaskId()] = i
				}

				// Remove local cache taskMsgs
				m.lockTask.Lock()
				for i := 0; i < len(m.taskMsgCache); i++ {
					msg := m.taskMsgCache[i]
					if _, ok := tmp[msg.GetTaskId()]; ok {
						delete(tmp, msg.GetTaskId())
						m.taskMsgCache = append(m.taskMsgCache[:i], m.taskMsgCache[i+1:]...)
						m.dataCenter.StoreMessageCache(m.taskMsgCache)
						i--
					}
				}
				m.lockTask.Unlock()

				// Revoke remote task
				if len(tmp) != 0 {
					msgs, index := make(types.TaskTerminateMsgArr, len(tmp)), 0
					for _, i := range tmp {
						msgs[index] = eventMessage.Msgs[i]
						index++
					}
					m.BroadcastTaskTerminateMsgArr(msgs)
				}
			}

		case <-powerTicker.C:

			if len(m.powerMsgCache) > 0 {
				m.BroadcastPowerMsgArr(m.powerMsgCache)
				m.powerMsgCache = make(types.PowerMsgArr, 0)
				m.dataCenter.StoreMessageCache(m.powerMsgCache)
			}

		case <-metadataTicker.C:

			if len(m.metadataMsgCache) > 0 {
				m.BroadcastMetadataMsgArr(m.metadataMsgCache)
				m.metadataMsgCache = make(types.MetadataMsgArr, 0)
				m.dataCenter.StoreMessageCache(m.metadataAuthMsgCache)
			}

		case <-metadataAuthTicker.C:

			if len(m.metadataAuthMsgCache) > 0 {
				m.BroadcastMetadataAuthMsgArr(m.metadataAuthMsgCache)
				m.metadataAuthMsgCache = make(types.MetadataAuthorityMsgArr, 0)
				m.dataCenter.StoreMessageCache(m.metadataAuthMsgCache)
			}

		case <-taskTicker.C:

			if len(m.taskMsgCache) > 0 {
				m.BroadcastTaskMsgArr(m.taskMsgCache)
				m.taskMsgCache = make(types.TaskMsgArr, 0)
				m.dataCenter.StoreMessageCache(m.taskMsgCache)
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

func (m *MessageHandler) BroadcastIdentityMsg(msg *types.IdentityMsg) {

	// add identity to local db
	if err := m.dataCenter.StoreIdentity(msg.GetOrganization()); nil != err {
		log.Errorf("Failed to store local org identity on MessageHandler with broadcast identity, identityId: {%s}, err: {%s}", msg.GetOwnerIdentityId(), err)
		return
	}

	// send identity to datacenter
	if err := m.dataCenter.InsertIdentity(msg.ToDataCenter()); nil != err {
		log.Errorf("Failed to broadcast org org identity on MessageHandler with broadcast identity, identityId: {%s}, nodeId: {%s}, nodeName: {%s}, err: {%s}",
			msg.GetOwnerIdentityId(), msg.GetOwnerNodeId(), msg.GetOwnerName(), err)
		return
	}
	log.Debugf("broadcast identity msg succeed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}", msg.GetOwnerIdentityId(), msg.GetOwnerNodeId(), msg.GetOwnerName())
	return
}

func (m *MessageHandler) BroadcastIdentityRevokeMsg() {

	// query local identity
	identity, err := m.dataCenter.QueryIdentity()
	if nil != err {
		log.Errorf("Failed to get local org identity on MessageHandler with revoke identity, identityId: {%s}, err: {%s}", identity.GetIdentityId(), err)
		return
	}

	// what if running task, can not revoke identity
	jobNodes , err := m.dataCenter.QueryRegisterNodeList(pb.PrefixTypeJobNode)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.Errorf("query all jobNode failed, %s", err)
		return
	}
	for _, node := range jobNodes {
		runningTaskCount, err := m.dataCenter.QueryRunningTaskCountOnJobNode(node.Id)
		if rawdb.IsNoDBNotFoundErr(err) {
			log.Errorf("query local running taskCount on old jobNode failed on MessageHandler with revoke identity, %s", err)
			return
		}
		if runningTaskCount > 0 {
			log.Errorf("the old jobNode have been running {%d} task current, don't revoke identity it", runningTaskCount)
			return
		}
	}

	// todo what if have any power publised, revoke the powers,
	// todo what metadata auth using on
	// todo what if metadata publish


	// remove local identity
	if err := m.dataCenter.RemoveIdentity(); nil != err {
		log.Errorf("Failed to delete org identity to local on MessageHandler with revoke identity, identityId: {%s}, err: {%s}", identity.GetIdentityId(), err)
		return
	}

	// remove identity from dataCenter
	if err := m.dataCenter.RevokeIdentity(
		types.NewIdentity(&libtypes.IdentityPB{
			NodeName:   identity.GetNodeName(),
			NodeId:     identity.GetNodeId(),
			IdentityId: identity.GetIdentityId(),
			DataId: "",
			DataStatus: apicommonpb.DataStatus_DataStatus_Deleted,
			Status:     apicommonpb.CommonStatus_CommonStatus_NonNormal,
			Credential: "",
		})); nil != err {
		log.Errorf("Failed to remove org identity to remote on MessageHandler with revoke identity, identityId: {%s}, err: {%s}", identity.GetIdentityId(), err)
		return
	}
	log.Debugf("Revoke identity msg succeed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}", identity.GetIdentityId(), identity.GetNodeId(), identity.GetNodeName())
	return
}

func (m *MessageHandler) BroadcastPowerMsgArr(powerMsgArr types.PowerMsgArr) {

	identity, err := m.dataCenter.QueryIdentity()
	if nil != err {
		log.Errorf("query local identityInfo failed on MessageHandler with broadcast power, {%s}", err)
		return
	}

	slotUnit, err := m.dataCenter.QueryNodeResourceSlotUnit()
	if nil != err {
		log.Errorf("query local slotUnit failed, {%s}", err)
		return
	}

	for _, power := range powerMsgArr {

		// query local resource
		resource, err := m.dataCenter.QueryLocalResource(power.GetJobNodeId())
		if nil != err {
			log.Errorf("Failed to query local resource on MessageHandler with broadcast power, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.GetPowerId(), power.GetJobNodeId(), err)
			continue
		}

		// set powerId to resource and change state
		resource.GetData().DataId = power.GetPowerId()
		resource.GetData().State = apicommonpb.PowerState_PowerState_Released

		// store local resource with totally
		resourceTable := types.NewLocalResourceTable(
			power.GetJobNodeId(),
			power.GetPowerId(),
			resource.GetData().GetTotalMem(),
			resource.GetData().GetTotalBandwidth(),
			resource.GetData().GetTotalProcessor(),
		)
		resourceTable.SetSlotUnit(slotUnit)

		log.Debugf("Publish power, StoreLocalResourceTable, %s", resourceTable.String())
		if err := m.dataCenter.StoreLocalResourceTable(resourceTable); nil != err {
			log.Errorf("Failed to StoreLocalResourceTable on MessageHandler with broadcast power, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.GetPowerId(), power.GetJobNodeId(), err)
			continue
		}

		if err := m.dataCenter.StoreLocalResourceIdByPowerId(power.GetPowerId(), power.GetJobNodeId()); nil != err {
			log.Errorf("Failed to StoreLocalResourceIdByPowerId on MessageHandler with broadcast power,  powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.GetPowerId(), power.GetJobNodeId(), err)
			continue
		}

		// update local resource
		if err := m.dataCenter.InsertLocalResource(resource); nil != err {
			log.Errorf("Failed to update local resource with powerId to local on MessageHandler with broadcast power, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.GetPowerId(), power.GetJobNodeId(), err)
			continue
		}

		// publish to global
		if err := m.dataCenter.InsertResource(types.NewResource(&libtypes.ResourcePB{
			IdentityId: identity.GetIdentityId(),
			NodeId:     identity.GetNodeId(),
			NodeName:   identity.GetNodeName(),
			DataId:     power.GetPowerId(),
			// the status of data, N means normal, D means deleted.
			DataStatus: apicommonpb.DataStatus_DataStatus_Normal,
			// resource status, eg: create/release/revoke
			State: apicommonpb.PowerState_PowerState_Released,
			// unit: byte
			TotalMem: resource.GetData().GetTotalMem(),
			UsedMem: 0,
			// number of cpu cores.
			TotalProcessor: resource.GetData().GetTotalProcessor(),
			UsedProcessor:  0,
			// unit: byte
			TotalBandwidth: resource.GetData().GetTotalBandwidth(),
			UsedBandwidth:  0,
			PublishAt: timeutils.UnixMsecUint64(),
			UpdateAt:  timeutils.UnixMsecUint64(),
		})); nil != err {
			log.Errorf("Failed to store power to dataCenter on MessageHandler with broadcast power, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				power.GetPowerId(), power.GetJobNodeId(), err)
			continue
		}

		log.Debugf("broadcast power msg succeed, powerId: {%s}, jobNodeId: {%s}", power.GetPowerId(), power.GetJobNodeId())

	}
	return
}

func (m *MessageHandler) BroadcastPowerRevokeMsgArr(powerRevokeMsgArr types.PowerRevokeMsgArr) {

	identity, err := m.dataCenter.QueryIdentity()
	if nil != err {
		log.Errorf("failed to query local identityInfo failed on MessageHandler with revoke power, err: {%s}", err)
		return
	}

	for _, revoke := range powerRevokeMsgArr {

		jobNodeId, err := m.dataCenter.QueryLocalResourceIdByPowerId(revoke.GetPowerId())
		if nil != err {
			log.Errorf("Failed to QueryLocalResourceIdByPowerId on MessageHandler with revoke power, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.GetPowerId(), jobNodeId, err)
			continue
		}
		if err := m.dataCenter.RemoveLocalResourceIdByPowerId(revoke.GetPowerId()); nil != err {
			log.Errorf("Failed to RemoveLocalResourceIdByPowerId on MessageHandler with revoke power, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.GetPowerId(), jobNodeId, err)
			continue
		}
		if err := m.dataCenter.RemoveLocalResourceTable(jobNodeId); nil != err {
			log.Errorf("Failed to RemoveLocalResourceTable on MessageHandler with revoke power, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.GetPowerId(), jobNodeId, err)
			continue
		}


		// query local resource
		resource, err := m.dataCenter.QueryLocalResource(jobNodeId)
		if nil != err {
			log.Errorf("Failed to query local resource on MessageHandler with revoke power, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.GetPowerId(), jobNodeId, err)
			continue
		}

		// remove powerId from local resource and change state
		resource.GetData().DataId = ""
		resource.GetData().State = apicommonpb.PowerState_PowerState_Revoked
		// update local resource
		if err := m.dataCenter.InsertLocalResource(resource); nil != err {
			log.Errorf("Failed to update local resource with powerId to local on MessageHandler with revoke power, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.GetPowerId(), jobNodeId, err)
			continue
		}

		// remove from global
		if err := m.dataCenter.RevokeResource(types.NewResource(&libtypes.ResourcePB{
			IdentityId: identity.GetIdentityId(),
			NodeId:     identity.GetNodeId(),
			NodeName:   identity.GetNodeName(),
			DataId:     revoke.GetPowerId(),
			// the status of data, N means normal, D means deleted.
			DataStatus: apicommonpb.DataStatus_DataStatus_Deleted,
			// resource status, eg: create/release/revoke
			State: apicommonpb.PowerState_PowerState_Revoked,
			UpdateAt: timeutils.UnixMsecUint64(),
		})); nil != err {
			log.Errorf("Failed to remove dataCenter resource on MessageHandler with revoke power, powerId: {%s}, jobNodeId: {%s}, err: {%s}",
				revoke.GetPowerId(), jobNodeId, err)
			continue
		}

		log.Debugf("revoke power msg succeed, powerId: {%s}, jobNodeId: {%s}", revoke.GetPowerId(), jobNodeId)
	}
	return
}

func (m *MessageHandler) BroadcastMetadataMsgArr(metadataMsgArr types.MetadataMsgArr) {

	identity, err := m.dataCenter.QueryIdentity()
	if nil != err {
		log.Errorf("Failed to query local identity on MessageHandler with broadcast metadata, err: {%s}", err)
		return
	}

	for _, metadata := range metadataMsgArr {

		// maintain the orginId and metadataId relationship of the local data service
		dataResourceFileUpload, err := m.dataCenter.QueryDataResourceFileUpload(metadata.GetOriginId())
		if nil != err {
			log.Errorf("Failed to QueryDataResourceFileUpload on MessageHandler with broadcast metadata, originId: {%s}, metadataId: {%s}, err: {%s}",
				metadata.GetOriginId(), metadata.GetMetadataId(), err)
			continue
		}

		// Update metadataId in fileupload information
		dataResourceFileUpload.SetMetadataId(metadata.GetMetadataId())
		if err := m.dataCenter.StoreDataResourceFileUpload(dataResourceFileUpload); nil != err {
			log.Errorf("Failed to StoreDataResourceFileUpload on MessageHandler with broadcast metadata, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metadata.GetOriginId(), metadata.GetMetadataId(), dataResourceFileUpload.GetNodeId(), err)
			continue
		}
		// Record the size of the resources occupied by the original data
		dataResourceTable, err := m.dataCenter.QueryDataResourceTable(dataResourceFileUpload.GetNodeId())
		if nil != err {
			log.Errorf("Failed to QueryDataResourceTable on MessageHandler with broadcast metadata, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metadata.GetOriginId(), metadata.GetMetadataId(), dataResourceFileUpload.GetNodeId(), err)
			continue
		}
		// update disk used of data resource table
		dataResourceTable.UseDisk(metadata.GetSize())
		if err := m.dataCenter.StoreDataResourceTable(dataResourceTable); nil != err {
			log.Errorf("Failed to StoreDataResourceTable on MessageHandler with broadcast metadata, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metadata.GetOriginId(), metadata.GetMetadataId(), dataResourceFileUpload.GetNodeId(), err)
			continue
		}
		// Separately record the GetSize of the metaData and the dataNodeId where it is located
		if err := m.dataCenter.StoreDataResourceDiskUsed(types.NewDataResourceDiskUsed(
			metadata.GetMetadataId(), dataResourceFileUpload.GetNodeId(), metadata.GetSize())); nil != err {
			log.Errorf("Failed to StoreDataResourceDiskUsed on MessageHandler with broadcast metadata, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metadata.GetOriginId(), metadata.GetMetadataId(), dataResourceFileUpload.GetNodeId(), err)
			continue
		}

		// publish metadata information
		if err := m.dataCenter.InsertMetadata(metadata.ToDataCenter(identity)); nil != err {
			log.Errorf("Failed to store metadata to dataCenter on MessageHandler with broadcast metadata, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				metadata.GetOriginId(), metadata.GetMetadataId(), dataResourceFileUpload.GetNodeId(), err)

			m.dataCenter.RemoveDataResourceDiskUsed(metadata.GetMetadataId())
			dataResourceTable.FreeDisk(metadata.GetSize())
			m.dataCenter.StoreDataResourceTable(dataResourceTable)

			continue
		}

		log.Debugf("broadcast metadata msg succeed, originId: {%s}, metadataId: {%s}", metadata.GetOriginId(), metadata.GetMetadataId())
	}

	return
}

func (m *MessageHandler) BroadcastMetadataRevokeMsgArr(metadataRevokeMsgArr types.MetadataRevokeMsgArr) {

	identity, err := m.dataCenter.QueryIdentity()
	if nil != err {
		log.Errorf("Failed to query local identity on MessageHandler with revoke metadata, err: {%s}", err)
		return
	}

	for _, revoke := range metadataRevokeMsgArr {
		// (metaDataId -> {metaDataId, dataNodeId, diskUsed})
		dataResourceDiskUsed, err := m.dataCenter.QueryDataResourceDiskUsed(revoke.GetMetadataId())
		if nil != err {
			log.Errorf("Failed to QueryDataResourceDiskUsed on MessageHandler with revoke metadata, metadataId: {%s}, err: {%s}",
				revoke.GetMetadataId(), err)
			continue
		}
		// update dataNode table (dataNodeId -> {dataNodeId, totalDisk, usedDisk})
		dataResourceTable, err := m.dataCenter.QueryDataResourceTable(dataResourceDiskUsed.GetNodeId())
		if nil != err {
			log.Errorf("Failed to QueryDataResourceTable on MessageHandler with revoke metadata, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.GetMetadataId(), dataResourceDiskUsed.GetNodeId(), err)
			continue
		}
		dataResourceTable.FreeDisk(dataResourceDiskUsed.GetDiskUsed())
		if err := m.dataCenter.StoreDataResourceTable(dataResourceTable); nil != err {
			log.Errorf("Failed to StoreDataResourceTable on MessageHandler with revoke metadata, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.GetMetadataId(), dataResourceDiskUsed.GetNodeId(), err)
			continue
		}

		// remove dataNodeDiskUsed (metaDataId -> {metaDataId, dataNodeId, diskUsed})
		if err := m.dataCenter.RemoveDataResourceDiskUsed(revoke.GetMetadataId()); nil != err {
			log.Errorf("Failed to RemoveDataResourceDiskUsed on MessageHandler with revoke metadata, metadataId: {%s}, dataNodeId: {%s}, err: {%s}",
				revoke.GetMetadataId(), dataResourceDiskUsed.GetNodeId(), err)
			continue
		}

		// revoke from global
		if err := m.dataCenter.RevokeMetadata(revoke.ToDataCenter(identity)); nil != err {
			log.Errorf("Failed to store metadata to dataCenter on MessageHandler with revoke metadata, metadataId: {%s}, err: {%s}",
				revoke.GetMetadataId(), err)
			continue
		}
		log.Debugf("revoke metadata msg succeed, metadataId: {%s}", revoke.GetMetadataId())
	}
	return
}

func (m *MessageHandler) BroadcastMetadataAuthMsgArr(metadataAuthMsgArr types.MetadataAuthorityMsgArr) {
	for _, msg := range metadataAuthMsgArr {

		has, err := m.authManager.HasValidMetadataAuth(msg.GetUserType(), msg.GetUser(), msg.GetMetadataAuthorityOwnerIdentity(), msg.GetMetadataAuthorityMetadataId())
		if nil != err {
			log.Errorf("Failed to call HasValidLocalMetadataAuth on MessageHandler with broadcast metadataAuth, metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}, err: {%s}",
				msg.GetMetadataAuthId(), msg.GetMetadataAuthority().GetMetadataId(), msg.GetUserType(), msg.GetUser(), err)
			continue
		}

		if has {
			log.Errorf("Failed to broadcast metadataAuth, cause alreay has valid last metadataAuth on MessageHandler with broadcast metadataAuth, userType: {%s}, user: {%s}, metadataId: {%s}",
				msg.GetUserType(), msg.GetUser(), msg.GetMetadataAuthority().GetMetadataId())
			continue
		}

		// check usageType/endTime once again before store and pushlish
		var (
			expire bool
			state apicommonpb.MetadataAuthorityState
		)

		state = apicommonpb.MetadataAuthorityState_MAState_Released

		switch msg.GetMetadataAuthority().GetUsageRule().GetUsageType() {
		case apicommonpb.MetadataUsageType_Usage_Period:
			if timeutils.UnixMsecUint64() >= msg.GetMetadataAuthority().GetUsageRule().GetEndAt() {
				expire = true
				state = apicommonpb.MetadataAuthorityState_MAState_Invalid
			} else {
				expire = false
			}
		case apicommonpb.MetadataUsageType_Usage_Times:
			// do nothing
		default:
			log.Errorf("unknown usageType of the metadataAuth on MessageHandler with broadcast metadataAuth, userType: {%s}, user: {%s}, metadataId: {%s}, usageType: {%s}",
				msg.GetUserType().String(), msg.GetUser(), msg.GetMetadataAuthority().GetMetadataId(), msg.GetMetadataAuthority().GetUsageRule().GetUsageType().String())
			continue
		}

		// Store metadataAuthority
		if err := m.authManager.ApplyMetadataAuthority(types.NewMetadataAuthority(&libtypes.MetadataAuthorityPB{
			MetadataAuthId:  msg.GetMetadataAuthId(),
			User:            msg.GetUser(),
			UserType:        msg.GetUserType(),
			Auth:            msg.GetMetadataAuthority(),
			AuditOption:     apicommonpb.AuditMetadataOption_Audit_Pending,
			AuditSuggestion: "",
			UsedQuo: &libtypes.MetadataUsedQuo{
				UsageType: msg.GetMetadataAuthority().GetUsageRule().GetUsageType(),
				Expire: expire,
				UsedTimes: 0,
			},
			ApplyAt: msg.GetCreateAt(),
			AuditAt: 0,
			State:   state,
			Sign:    msg.GetSign(),
		})); nil != err {
			log.Errorf("Failed to store metadataAuth to dataCenter on MessageHandler with broadcast metadataAuth, metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}, err: {%s}",
				msg.GetMetadataAuthId(), msg.GetMetadataAuthority().GetMetadataId(), msg.GetUserType(), msg.GetUser(), err)
			continue
		}

		log.Debugf("broadcast metadataAuth msg succeed, metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}",
			msg.GetMetadataAuthId(), msg.GetMetadataAuthority().GetMetadataId(), msg.GetUserType(), msg.GetUser())
	}

	return
}

func (m *MessageHandler) BroadcastMetadataAuthRevokeMsgArr(metadataAuthRevokeMsgArr types.MetadataAuthorityRevokeMsgArr) {
	for _, revoke := range metadataAuthRevokeMsgArr {

		// verify
		metadataAuth, err := m.authManager.GetMetadataAuthority(revoke.GetMetadataAuthId())
		if nil != err {
			log.Errorf("Failed to query old metadataAuth on MessageHandler with revoke metadataAuth, metadataAuthId: {%s}, user:{%s}, userType: {%s}, err: {%s}",
				revoke.GetMetadataAuthId(), revoke.GetUser(), revoke.GetUserType().String(), err)
			continue
		}

		if metadataAuth.GetData().GetUser() != revoke.GetUser() || metadataAuth.GetData().GetUserType() != revoke.GetUserType() {
			log.Errorf("user of metadataAuth is wrong on MessageHandler with revoke metadataAuth, metadataAuthId: {%s}, user:{%s}, userType: {%s}",
				revoke.GetMetadataAuthId(), revoke.GetUser(), revoke.GetUserType().String())
			continue
		}

		if bytes.Compare(metadataAuth.GetData().GetSign(), revoke.GetSign()) != 0 {
			log.Errorf("user sign of metadataAuth is wrong on MessageHandler with revoke metadataAuth, metadataAuthId: {%s}, user:{%s}, userType: {%s}, metadataAuth's sign: {%v}, revoke msg's sign: {%v}",
				revoke.GetMetadataAuthId(), revoke.GetUser(), revoke.GetUserType().String(), metadataAuth.GetData().GetSign(), revoke.GetSign())
			continue
		}

		// The data authorization application information that has been `invalidated` or has been `revoked` is not allowed to be revoked
		if metadataAuth.GetData().GetState() == apicommonpb.MetadataAuthorityState_MAState_Revoked ||
			metadataAuth.GetData().GetState() == apicommonpb.MetadataAuthorityState_MAState_Invalid {
			log.Errorf("state of metadataAuth is wrong on MessageHandler with revoke metadataAuth, metadataAuthId: {%s}, user:{%s}, state: {%s}",
				revoke.GetMetadataAuthId(), revoke.GetUser(), metadataAuth.GetData().GetState().String())
			continue
		}

		// The data authorization application information that has been audited and cannot be revoked
		if metadataAuth.GetData().GetAuditOption() != apicommonpb.AuditMetadataOption_Audit_Pending {
			log.Errorf("the metadataAuth has audit on MessageHandler with revoke metadataAuth, metadataAuthId: {%s}, user:{%s}, state: {%s}",
				revoke.GetMetadataAuthId(), revoke.GetUser(), metadataAuth.GetData().GetAuditOption().String())
			continue
		}

		if err := m.dataCenter.UpdateMetadataAuthority(types.NewMetadataAuthority(&libtypes.MetadataAuthorityPB{
			MetadataAuthId:  revoke.GetMetadataAuthId(),
			User:            revoke.GetUser(),
			UserType:        revoke.GetUserType(),
			Auth:            &libtypes.MetadataAuthority{},
			AuditOption:     metadataAuth.GetData().GetAuditOption(),
			AuditSuggestion: metadataAuth.GetData().GetAuditSuggestion(),
			UsedQuo:         metadataAuth.GetData().GetUsedQuo(),
			ApplyAt:         metadataAuth.GetData().GetApplyAt(),
			AuditAt:         metadataAuth.GetData().GetAuditAt(),
			State:           apicommonpb.MetadataAuthorityState_MAState_Revoked,
		})); nil != err {
			log.Errorf("Failed to store metadataAuth to dataCenter on MessageHandler with revoke metadataAuth, metadataAuthId: {%s}, user:{%s}, err: {%s}",
				revoke.GetMetadataAuthId(), revoke.GetUser(), err)
			continue
		}
		log.Debugf("revoke metadataAuth msg succeed, metadataAuthId: {%s}", revoke.GetMetadataAuthId())
	}

	return
}

func (m *MessageHandler) BroadcastTaskMsgArr(taskMsgArr types.TaskMsgArr) {
	if err := m.taskManager.SendTaskMsgArr(taskMsgArr); nil != err{
		log.Errorf("Failed to call `BroadcastTaskMsgArr` on MessageHandler, %s", err)
	}
	return
}

func (m *MessageHandler) BroadcastTaskTerminateMsgArr(terminateMsgArr types.TaskTerminateMsgArr) {
	if err := m.taskManager.SendTaskTerminate(terminateMsgArr); nil != err{
		log.Errorf("Failed to call `BroadcastTaskTerminateMsgArr` on MessageHandler, %s", err)
	}
	return
}
