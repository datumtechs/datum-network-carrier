package message

import (
	"encoding/json"
	auth2 "github.com/Metisnetwork/Metis-Carrier/ach/auth"
	"github.com/Metisnetwork/Metis-Carrier/common/feed"
	"github.com/Metisnetwork/Metis-Carrier/common/timeutils"
	"github.com/Metisnetwork/Metis-Carrier/core/rawdb"
	"github.com/Metisnetwork/Metis-Carrier/core/resource"
	"github.com/Metisnetwork/Metis-Carrier/core/task"
	"github.com/Metisnetwork/Metis-Carrier/event"
	pb "github.com/Metisnetwork/Metis-Carrier/lib/api"
	libtypes "github.com/Metisnetwork/Metis-Carrier/lib/types"
	"github.com/Metisnetwork/Metis-Carrier/rpc/backend"
	"github.com/Metisnetwork/Metis-Carrier/types"
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
	pool        *Mempool
	resourceMng *resource.Manager
	// Send taskMsg to taskManager
	taskManager *task.Manager
	authManager *auth2.AuthorityManager
	// internal resource node set (Fighter node grpc client set)
	msgChannel chan *feed.Event
	quit       chan struct{}

	msgSub event.Subscription

	powerMsgCache        types.PowerMsgArr
	metadataMsgCache     types.MetadataMsgArr
	metadataAuthMsgCache types.MetadataAuthorityMsgArr
	taskMsgCache         types.TaskMsgArr

	lockPower        sync.Mutex
	lockMetadata     sync.Mutex
	lockMetadataAuth sync.Mutex
	lockTask         sync.Mutex
}

func NewHandler(pool *Mempool, resourceMng *resource.Manager, taskManager *task.Manager, authManager *auth2.AuthorityManager) *MessageHandler {
	m := &MessageHandler{
		pool:        pool,
		resourceMng: resourceMng,
		taskManager: taskManager,
		authManager: authManager,
		msgChannel:  make(chan *feed.Event, 5),
		quit:        make(chan struct{}),
	}
	return m
}

func (m *MessageHandler) Start() error {
	m.msgSub = m.pool.SubscribeNewMessageEvent(m.msgChannel)
	m.recoveryCache()
	go m.loop()
	log.Info("Started messageHandler ...")
	return nil
}
func (m *MessageHandler) Stop() error {
	close(m.quit)
	return nil
}

func (m *MessageHandler) recoveryCache() {
	taskMsgCache, err := m.resourceMng.GetDB().QueryTaskMsgArr()
	if nil != err {
		if rawdb.IsNoDBNotFoundErr(err) {
			log.WithError(err).Errorf("Failed to get taskMsgCache from QueryTaskMsgArr")
		} else {
			log.WithError(err).Warnf("Not found taskMsgCache from QueryTaskMsgArr")
		}
	} else {
		m.taskMsgCache = taskMsgCache
	}

	metadataAuthMsgCache, err := m.resourceMng.GetDB().QueryMetadataAuthorityMsgArr()
	if nil != err {
		if rawdb.IsNoDBNotFoundErr(err) {
			log.WithError(err).Errorf("Failed to get metadataAuthMsgCache from QueryMetadataAuthorityMsgArr")
		} else {
			log.WithError(err).Warnf("Not found metadataAuthMsgCache from QueryMetadataAuthorityMsgArr")
		}
	} else {
		m.metadataAuthMsgCache = metadataAuthMsgCache
	}

	metadataMsgCache, err := m.resourceMng.GetDB().QueryMetadataMsgArr()
	if nil != err {
		if rawdb.IsNoDBNotFoundErr(err) {
			log.WithError(err).Errorf("Failed to get metadataMsgCache from QueryMetadataMsgArr")
		} else {
			log.WithError(err).Warnf("Not found metadataMsgCache from QueryMetadataMsgArr")
		}
	} else {
		m.metadataMsgCache = metadataMsgCache
	}

	powerMsgCache, err := m.resourceMng.GetDB().QueryPowerMsgArr()
	if nil != err {
		if rawdb.IsNoDBNotFoundErr(err) {
			log.WithError(err).Errorf("Failed to get powerMsgCache from QueryPowerMsgArr")
		} else {
			log.WithError(err).Warnf("Not found powerMsgCache from QueryPowerMsgArr")
		}
	} else {
		m.powerMsgCache = powerMsgCache
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
				m.powerMsgCache = append(m.powerMsgCache, msg.Msg)
				m.resourceMng.GetDB().StoreMessageCache(msg.Msg) // backup power msg into disk
				if len(m.powerMsgCache) >= defaultPowerMsgsCacheSize {
					m.BroadcastPowerMsgArr(m.powerMsgCache)
					m.powerMsgCache = make(types.PowerMsgArr, 0)
				}
				m.lockPower.Unlock()
			case types.RevokePower:
				revoke := event.Data.(*types.PowerRevokeMsgEvent)

				var flag bool
				// Remove local cache powerMsgs
				m.lockPower.Lock()
				for i := 0; i < len(m.powerMsgCache); i++ {
					msg := m.powerMsgCache[i]
					if revoke.Msg.GetPowerId() == msg.GetPowerId() {
						flag = true
						m.powerMsgCache = append(m.powerMsgCache[:i], m.powerMsgCache[i+1:]...)
						m.resourceMng.GetDB().RemovePowerMsg(msg.GetPowerId()) // remove from disk
						i--
					}
				}
				m.lockPower.Unlock()

				// Revoke remote power
				if !flag {
					m.BroadcastPowerRevokeMsgArr(types.PowerRevokeMsgArr{revoke.Msg})
				}
			case types.ApplyMetadata:
				msg := event.Data.(*types.MetadataMsgEvent)
				m.lockMetadata.Lock()
				m.metadataMsgCache = append(m.metadataMsgCache, msg.Msg)
				m.resourceMng.GetDB().StoreMessageCache(msg.Msg) // backup metadata msg into disk
				if len(m.metadataMsgCache) >= defaultMetadataMsgsCacheSize {
					m.BroadcastMetadataMsgArr(m.metadataMsgCache)
					m.metadataMsgCache = make(types.MetadataMsgArr, 0)
				}
				m.lockMetadata.Unlock()
			case types.RevokeMetadata:
				revoke := event.Data.(*types.MetadataRevokeMsgEvent)

				var flag bool

				// Remove local cache metadataMsgs
				m.lockMetadata.Lock()
				for i := 0; i < len(m.metadataMsgCache); i++ {
					msg := m.metadataMsgCache[i]
					if revoke.Msg.GetMetadataId() == msg.GetMetadataId() {
						flag = true
						m.metadataMsgCache = append(m.metadataMsgCache[:i], m.metadataMsgCache[i+1:]...)
						m.resourceMng.GetDB().RemoveMetadataMsg(msg.GetMetadataId()) // remove from disk
						i--
					}
				}
				m.lockMetadata.Unlock()

				// Revoke remote metadata
				if !flag {
					m.BroadcastMetadataRevokeMsgArr(types.MetadataRevokeMsgArr{revoke.Msg})
				}

			case types.ApplyMetadataAuth:
				msg := event.Data.(*types.MetadataAuthMsgEvent)
				m.lockMetadataAuth.Lock()
				m.metadataAuthMsgCache = append(m.metadataAuthMsgCache, msg.Msg)
				m.resourceMng.GetDB().StoreMessageCache(msg.Msg) // backup metadataAuth into disk
				if len(m.metadataAuthMsgCache) >= defaultMetadataAuthMsgsCacheSize {
					m.BroadcastMetadataAuthMsgArr(m.metadataAuthMsgCache)
					m.metadataAuthMsgCache = make(types.MetadataAuthorityMsgArr, 0)
				}
				m.lockMetadataAuth.Unlock()
			case types.RevokeMetadataAuth:
				revoke := event.Data.(*types.MetadataAuthRevokeMsgEvent)

				var flag bool

				// Remove local cache metadataAuthorityMsgs
				m.lockMetadataAuth.Lock()
				for i := 0; i < len(m.metadataAuthMsgCache); i++ {
					msg := m.metadataAuthMsgCache[i]
					if revoke.Msg.GetMetadataAuthId() == msg.GetMetadataAuthId() {
						flag = true
						m.metadataAuthMsgCache = append(m.metadataAuthMsgCache[:i], m.metadataAuthMsgCache[i+1:]...)
						m.resourceMng.GetDB().RemoveMetadataAuthMsg(msg.GetMetadataAuthId()) // remove from disk
						i--
					}
				}
				m.lockMetadataAuth.Unlock()

				// Revoke remote metadataAuthority
				if !flag {
					m.BroadcastMetadataAuthRevokeMsgArr(types.MetadataAuthorityRevokeMsgArr{revoke.Msg})
				}

			case types.ApplyTask:
				msg := event.Data.(*types.TaskMsgEvent)
				m.lockTask.Lock()
				m.taskMsgCache = append(m.taskMsgCache, msg.Msg)
				m.resourceMng.GetDB().StoreMessageCache(msg.Msg) // backup task msg into disk
				if len(m.taskMsgCache) >= defaultTaskMsgsCacheSize {
					m.BroadcastTaskMsgArr(m.taskMsgCache)
					m.taskMsgCache = make(types.TaskMsgArr, 0)
				}
				m.lockTask.Unlock()
			case types.TerminateTask:
				terminate := event.Data.(*types.TaskTerminateMsgEvent)

				var flag bool

				// Remove local cache taskMsgs
				m.lockTask.Lock()
				for i := 0; i < len(m.taskMsgCache); i++ {
					msg := m.taskMsgCache[i]
					if terminate.Msg.GetTaskId() == msg.GetTaskId() {
						flag = true
						m.taskMsgCache = append(m.taskMsgCache[:i], m.taskMsgCache[i+1:]...)
						m.resourceMng.GetDB().RemoveTaskMsg(msg.GetTaskId()) // remove from disk
						i--
					}
				}
				m.lockTask.Unlock()

				// Revoke remote task
				if !flag {
					m.BroadcastTaskTerminateMsgArr(types.TaskTerminateMsgArr{terminate.Msg})
				}
			}

		case <-powerTicker.C:

			if len(m.powerMsgCache) > 0 {
				m.BroadcastPowerMsgArr(m.powerMsgCache)
				m.powerMsgCache = make(types.PowerMsgArr, 0)
			}

		case <-metadataTicker.C:

			if len(m.metadataMsgCache) > 0 {
				m.BroadcastMetadataMsgArr(m.metadataMsgCache)
				m.metadataMsgCache = make(types.MetadataMsgArr, 0)
			}

		case <-metadataAuthTicker.C:

			if len(m.metadataAuthMsgCache) > 0 {
				m.BroadcastMetadataAuthMsgArr(m.metadataAuthMsgCache)
				m.metadataAuthMsgCache = make(types.MetadataAuthorityMsgArr, 0)
			}

		case <-taskTicker.C:

			if len(m.taskMsgCache) > 0 {
				m.BroadcastTaskMsgArr(m.taskMsgCache)
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

func (m *MessageHandler) BroadcastIdentityMsg(msg *types.IdentityMsg) {

	// add identity to local db
	identity := msg.GetOrganization()
	identity.DataStatus = libtypes.DataStatus_DataStatus_Valid
	identity.Status = libtypes.CommonStatus_CommonStatus_Valid
	if err := m.resourceMng.GetDB().StoreIdentity(identity); nil != err {
		log.WithError(err).Errorf("Failed to store local org identity on MessageHandler with broadcast identity, identityId: {%s}", msg.GetOwnerIdentityId())
		return
	}
	// TODO 填充 nonce
	// send identity to datacenter
	if err := m.resourceMng.GetDB().InsertIdentity(msg.ToDataCenter()); nil != err {
		log.WithError(err).Errorf("Failed to broadcast org org identity on MessageHandler with broadcast identity, identityId: {%s}, nodeId: {%s}, nodeName: {%s}",
			msg.GetOwnerIdentityId(), msg.GetOwnerNodeId(), msg.GetOwnerName())
		return
	}
	log.Debugf("broadcast identity msg succeed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}", msg.GetOwnerIdentityId(), msg.GetOwnerNodeId(), msg.GetOwnerName())
	return
}

func (m *MessageHandler) BroadcastIdentityRevokeMsg() {

	// query local identity
	identity, err := m.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to get local org identity on MessageHandler with revoke identity, identityId: {%s}", identity.GetIdentityId())
		return
	}

	// what if running task, can not revoke identity
	jobNodes, err := m.resourceMng.GetDB().QueryRegisterNodeList(pb.PrefixTypeJobNode)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("query all jobNode failed")
		return
	}
	for _, node := range jobNodes {
		runningTaskCount, err := m.resourceMng.GetDB().QueryJobNodeRunningTaskCount(node.Id)
		if rawdb.IsNoDBNotFoundErr(err) {
			log.Errorf("query local running taskCount on old jobNode failed on MessageHandler with revoke identity, %s", err)
			return
		}
		if runningTaskCount > 0 {
			log.Errorf("the old jobNode have been running {%d} task current, don't revoke identity it", runningTaskCount)
			return
		}
	}

	// TODO ###########################################################################
	// TODO ###########################################################################
	// TODO ###########################################################################
	// TODO ###########################################################################
	// TODO ###########################################################################
	// TODO ###########################################################################
	//
	// todo what if have any power publised, revoke the powers,
	// todo what metadata auth using on
	// todo what if metadata publish
	//
	// TODO ###########################################################################
	// TODO ###########################################################################
	// TODO ###########################################################################
	// TODO ###########################################################################
	// TODO ###########################################################################
	// TODO ###########################################################################

	// remove local identity
	if err := m.resourceMng.GetDB().RemoveIdentity(); nil != err {
		log.WithError(err).Errorf("Failed to delete org identity to local on MessageHandler with revoke identity, identityId: {%s}", identity.GetIdentityId())
		return
	}

	// remove identity from dataCenter
	if err := m.resourceMng.GetDB().RevokeIdentity(
		types.NewIdentity(&libtypes.IdentityPB{
			NodeName:   identity.GetNodeName(),
			NodeId:     identity.GetNodeId(),
			IdentityId: identity.GetIdentityId(),
			DataId:     "",
			DataStatus: libtypes.DataStatus_DataStatus_Invalid,
			Status:     libtypes.CommonStatus_CommonStatus_Invalid,
			Credential: "",
		})); nil != err {
		log.WithError(err).Errorf("Failed to remove org identity to remote on MessageHandler with revoke identity, identityId: {%s}", identity.GetIdentityId())
		return
	}
	log.Debugf("Revoke identity msg succeed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}", identity.GetIdentityId(), identity.GetNodeId(), identity.GetNodeName())
	return
}

func (m *MessageHandler) BroadcastPowerMsgArr(powerMsgArr types.PowerMsgArr) {

	identity, err := m.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("query local identityInfo failed on MessageHandler with broadcast power msg")
		return
	}

	for _, msg := range powerMsgArr {

		m.resourceMng.GetDB().RemovePowerMsg(msg.GetPowerId()) // remove from disk if msg been handle

		// query local resource
		resource, err := m.resourceMng.GetDB().QueryLocalResource(msg.GetJobNodeId())
		if nil != err {
			log.WithError(err).Errorf("Failed to query local resource on MessageHandler with broadcast msg, powerId: {%s}, jobNodeId: {%s}",
				msg.GetPowerId(), msg.GetJobNodeId())
			continue
		}

		// set powerId to resource and change state
		resource.GetData().DataId = msg.GetPowerId()
		resource.GetData().State = libtypes.PowerState_PowerState_Released

		// check jobNode wether connected?
		client, ok := m.resourceMng.QueryJobNodeClient(msg.GetJobNodeId())

		// store local resource with totally
		resourceTable := types.NewLocalResourceTable(
			msg.GetJobNodeId(),
			msg.GetPowerId(),
			resource.GetData().GetTotalMem(),
			resource.GetData().GetTotalBandwidth(),
			resource.GetData().GetTotalDisk(),
			resource.GetData().GetTotalProcessor(),
			ok && client.IsConnected(), // true OR  false ?
		)

		log.Debugf("Publish msg, StoreLocalResourceTable, %s", resourceTable.String())
		if err := m.resourceMng.GetDB().StoreLocalResourceTable(resourceTable); nil != err {
			log.WithError(err).Errorf("Failed to StoreLocalResourceTable on MessageHandler with broadcast msg, powerId: {%s}, jobNodeId: {%s}",
				msg.GetPowerId(), msg.GetJobNodeId())
			continue
		}

		if err := m.resourceMng.GetDB().StoreJobNodeIdIdByPowerId(msg.GetPowerId(), msg.GetJobNodeId()); nil != err {
			log.WithError(err).Errorf("Failed to StoreJobNodeIdByPowerId on MessageHandler with broadcast msg,  powerId: {%s}, jobNodeId: {%s}",
				msg.GetPowerId(), msg.GetJobNodeId())
			continue
		}

		// TODO 填充 nonce
		// update local resource
		if err := m.resourceMng.GetDB().StoreLocalResource(resource); nil != err {
			log.WithError(err).Errorf("Failed to update local resource with powerId to local on MessageHandler with broadcast msg, powerId: {%s}, jobNodeId: {%s}",
				msg.GetPowerId(), msg.GetJobNodeId())
			continue
		}
		// TODO 填充 nonce
		// publish to global
		if err := m.resourceMng.GetDB().InsertResource(types.NewResource(&libtypes.ResourcePB{
			Owner:  identity,
			DataId: msg.GetPowerId(),
			// the status of data for local storage, 1 means valid, 2 means invalid
			DataStatus: libtypes.DataStatus_DataStatus_Valid,
			// resource status, eg: create/release/revoke
			State: libtypes.PowerState_PowerState_Released,
			// unit: byte
			TotalMem: resource.GetData().GetTotalMem(),
			UsedMem:  0,
			// unit: byte
			TotalBandwidth: resource.GetData().GetTotalBandwidth(),
			UsedBandwidth:  0,
			// disk sixe
			TotalDisk: resource.GetData().GetTotalDisk(),
			UsedDisk:  0,
			// number of cpu cores.
			TotalProcessor: resource.GetData().GetTotalProcessor(),
			UsedProcessor:  0,
			PublishAt:      timeutils.UnixMsecUint64(),
			UpdateAt:       timeutils.UnixMsecUint64(),
		})); nil != err {
			log.WithError(err).Errorf("Failed to store power msg to dataCenter on MessageHandler with broadcast msg, powerId: {%s}, jobNodeId: {%s}",
				msg.GetPowerId(), msg.GetJobNodeId())
			continue
		}

		log.Debugf("broadcast power msg succeed, powerId: {%s}, jobNodeId: {%s}", msg.GetPowerId(), msg.GetJobNodeId())

	}
	return
}

func (m *MessageHandler) BroadcastPowerRevokeMsgArr(powerRevokeMsgArr types.PowerRevokeMsgArr) {

	identity, err := m.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("failed to query local identityInfo failed on MessageHandler with revoke power msg")
		return
	}

	for _, revoke := range powerRevokeMsgArr {

		jobNodeId, err := m.resourceMng.GetDB().QueryJobNodeIdByPowerId(revoke.GetPowerId())
		if nil != err {
			log.WithError(err).Errorf("Failed to call QueryJobNodeIdByPowerId() on MessageHandler with revoke power, powerId: {%s}, jobNodeId: {%s}",
				revoke.GetPowerId(), jobNodeId)
			continue
		}
		if err := m.resourceMng.GetDB().RemoveJobNodeIdByPowerId(revoke.GetPowerId()); nil != err {
			log.WithError(err).Errorf("Failed to call RemoveJobNodeIdByPowerId() on MessageHandler with revoke power, powerId: {%s}, jobNodeId: {%s}",
				revoke.GetPowerId(), jobNodeId)
			continue
		}
		if err := m.resourceMng.GetDB().RemoveLocalResourceTable(jobNodeId); nil != err {
			log.WithError(err).Errorf("Failed to RemoveLocalResourceTable on MessageHandler with revoke power, powerId: {%s}, jobNodeId: {%s}",
				revoke.GetPowerId(), jobNodeId)
			continue
		}

		// query local resource
		resource, err := m.resourceMng.GetDB().QueryLocalResource(jobNodeId)
		if nil != err {
			log.WithError(err).Errorf("Failed to query local resource on MessageHandler with revoke power, powerId: {%s}, jobNodeId: {%s}",
				revoke.GetPowerId(), jobNodeId)
			continue
		}

		// remove powerId from local resource and change state
		resource.GetData().DataId = ""
		resource.GetData().State = libtypes.PowerState_PowerState_Revoked
		// clean used resource value
		resource.GetData().UsedBandwidth = 0
		resource.GetData().UsedDisk = 0
		resource.GetData().UsedProcessor = 0
		resource.GetData().UsedMem = 0

		// update local resource
		if err := m.resourceMng.GetDB().StoreLocalResource(resource); nil != err {
			log.WithError(err).Errorf("Failed to update local resource with powerId to local on MessageHandler with revoke power, powerId: {%s}, jobNodeId: {%s}",
				revoke.GetPowerId(), jobNodeId)
			continue
		}

		// remove from global
		if err := m.resourceMng.GetDB().RevokeResource(types.NewResource(&libtypes.ResourcePB{
			Owner:  identity,
			DataId: revoke.GetPowerId(),
			// the status of data for local storage, 1 means valid, 2 means invalid
			DataStatus: libtypes.DataStatus_DataStatus_Invalid,
			// resource status, eg: create/release/revoke
			State:    libtypes.PowerState_PowerState_Revoked,
			UpdateAt: timeutils.UnixMsecUint64(),
		})); nil != err {
			log.WithError(err).Errorf("Failed to remove dataCenter resource on MessageHandler with revoke power, powerId: {%s}, jobNodeId: {%s}",
				revoke.GetPowerId(), jobNodeId)
			continue
		}

		log.Debugf("revoke power msg succeed, powerId: {%s}, jobNodeId: {%s}", revoke.GetPowerId(), jobNodeId)
	}
	return
}

func (m *MessageHandler) BroadcastMetadataMsgArr(metadataMsgArr types.MetadataMsgArr) {

	identity, err := m.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to query local identity on MessageHandler with broadcast metadata msg")
		return
	}

	for _, msg := range metadataMsgArr {

		m.resourceMng.GetDB().RemoveMetadataMsg(msg.GetMetadataId()) // remove from disk if msg been handle

		if types.IsCSVdata(msg.GetDataType()) {
			var option *types.MetadataOptionCSV
			if err := json.Unmarshal([]byte(msg.GetMetadataOption()), &option); nil != err {
				log.WithError(err).Errorf("Failed to unmashal metadataOption on MessageHandler with broadcast msg, metadataId: {%s}",
					msg.GetMetadataId())
				continue
			}

			// maintain the orginId and metadataId relationship of the local data service
			dataResourceFileUpload, err := m.resourceMng.GetDB().QueryDataResourceFileUpload(option.GetOriginId())
			if nil != err {
				log.WithError(err).Errorf("Failed to QueryDataResourceFileUpload on MessageHandler with broadcast msg, originId: {%s}, metadataId: {%s}",
					option.GetOriginId(), msg.GetMetadataId())
				continue
			}

			// Update metadataId in fileupload information
			dataResourceFileUpload.SetMetadataId(msg.GetMetadataId())
			if err := m.resourceMng.GetDB().StoreDataResourceFileUpload(dataResourceFileUpload); nil != err {
				log.WithError(err).Errorf("Failed to StoreDataResourceFileUpload on MessageHandler with broadcast msg, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}",
					option.GetOriginId(), msg.GetMetadataId(), dataResourceFileUpload.GetNodeId())
				continue
			}

			// Record the size of the resources occupied by the original data
			dataResourceTable, err := m.resourceMng.GetDB().QueryDataResourceTable(dataResourceFileUpload.GetNodeId())
			if nil != err {
				log.WithError(err).Errorf("Failed to QueryDataResourceTable on MessageHandler with broadcast msg, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}",
					option.GetOriginId(), msg.GetMetadataId(), dataResourceFileUpload.GetNodeId())
				continue
			}
			// update disk used of data resource table
			dataResourceTable.UseDisk(option.GetSize())
			if err := m.resourceMng.GetDB().StoreDataResourceTable(dataResourceTable); nil != err {
				log.WithError(err).Errorf("Failed to StoreDataResourceTable on MessageHandler with broadcast msg, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}",
					option.GetOriginId(), msg.GetMetadataId(), dataResourceFileUpload.GetNodeId())
				continue
			}

			// Separately record the GetSize of the metaData and the dataNodeId where it is located
			if err := m.resourceMng.GetDB().StoreDataResourceDiskUsed(types.NewDataResourceDiskUsed(
				msg.GetMetadataId(), dataResourceFileUpload.GetNodeId(), option.GetSize())); nil != err {
				log.WithError(err).Errorf("Failed to StoreDataResourceDiskUsed on MessageHandler with broadcast msg, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}",
					option.GetOriginId(), msg.GetMetadataId(), dataResourceFileUpload.GetNodeId())
				continue
			}
		}

		// TODO 填充 nonce
		// publish msg information
		if err := m.resourceMng.GetDB().InsertMetadata(msg.ToDataCenter(identity)); nil != err {
			log.WithError(err).Errorf("Failed to store msg to dataCenter on MessageHandler with broadcast msg, metadataId: {%s}",
				msg.GetMetadataId())

			//m.resourceMng.GetDB().RemoveDataResourceDiskUsed(msg.GetMetadataId())
			//dataResourceTable.FreeDisk(msg.GetSize())
			//m.resourceMng.GetDB().StoreDataResourceTable(dataResourceTable)

			continue
		}

		log.Debugf("broadcast metadata msg succeed, metadataId: {%s}", msg.GetMetadataId())
	}

	return
}

func (m *MessageHandler) BroadcastMetadataRevokeMsgArr(metadataRevokeMsgArr types.MetadataRevokeMsgArr) {

	identity, err := m.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to query local identity on MessageHandler with revoke metadata msg")
		return
	}

	for _, revoke := range metadataRevokeMsgArr {
		// (metaDataId -> {metaDataId, dataNodeId, diskUsed})
		dataResourceDiskUsed, err := m.resourceMng.GetDB().QueryDataResourceDiskUsed(revoke.GetMetadataId())
		if nil != err {
			log.WithError(err).Errorf("Failed to QueryDataResourceDiskUsed on MessageHandler with revoke metadata, metadataId: {%s}",
				revoke.GetMetadataId())
			continue
		}
		// update dataNode table (dataNodeId -> {dataNodeId, totalDisk, usedDisk})
		dataResourceTable, err := m.resourceMng.GetDB().QueryDataResourceTable(dataResourceDiskUsed.GetNodeId())
		if nil != err {
			log.WithError(err).Errorf("Failed to QueryDataResourceTable on MessageHandler with revoke metadata, metadataId: {%s}, dataNodeId: {%s}",
				revoke.GetMetadataId(), dataResourceDiskUsed.GetNodeId())
			continue
		}
		dataResourceTable.FreeDisk(dataResourceDiskUsed.GetDiskUsed())
		if err := m.resourceMng.GetDB().StoreDataResourceTable(dataResourceTable); nil != err {
			log.WithError(err).Errorf("Failed to StoreDataResourceTable on MessageHandler with revoke metadata, metadataId: {%s}, dataNodeId: {%s}",
				revoke.GetMetadataId(), dataResourceDiskUsed.GetNodeId())
			continue
		}

		// remove dataNodeDiskUsed (metaDataId -> {metaDataId, dataNodeId, diskUsed})
		if err := m.resourceMng.GetDB().RemoveDataResourceDiskUsed(revoke.GetMetadataId()); nil != err {
			log.WithError(err).Errorf("Failed to RemoveDataResourceDiskUsed on MessageHandler with revoke metadata, metadataId: {%s}, dataNodeId: {%s}",
				revoke.GetMetadataId(), dataResourceDiskUsed.GetNodeId())
			continue
		}

		// revoke from global
		if err := m.resourceMng.GetDB().RevokeMetadata(revoke.ToDataCenter(identity)); nil != err {
			log.WithError(err).Errorf("Failed to revoke metadata to dataCenter on MessageHandler with revoke metadata, metadataId: {%s}",
				revoke.GetMetadataId())
			continue
		}
		log.Debugf("revoke metadata msg succeed, metadataId: {%s}", revoke.GetMetadataId())
	}
	return
}

func (m *MessageHandler) BroadcastMetadataAuthMsgArr(metadataAuthMsgArr types.MetadataAuthorityMsgArr) {

	// ############################################
	// ############################################
	// NOTE:
	// 		check metadataId and identity of authority whether is valid ?
	//
	//		Because the message is asynchronous,
	//		the corresponding identity and metadata may have been deleted long ago,
	//		so the data legitimacy is verified again
	// ############################################
	// ############################################
	ideneityList, err := m.resourceMng.GetDB().QueryIdentityList(timeutils.BeforeYearUnixMsecUint64(), backend.DefaultMaxPageSize)
	if nil != err {
		log.WithError(err).Errorf("Failed to query global identity list on MessageHandler with broadcast metadataAuth")
		return
	}

	for _, msg := range metadataAuthMsgArr {

		m.resourceMng.GetDB().RemoveMetadataAuthMsg(msg.GetMetadataAuthId()) // remove from disk if msg been handle

		// ############################################
		// ############################################
		// NOTE:
		// 		check metadataId and identity of authority whether is valid ?
		//
		//		Because the message is asynchronous,
		//		the corresponding identity and metadata may have been deleted long ago,
		//		so the data legitimacy is verified again
		// ############################################
		// ############################################
		var valid bool // false
		for _, identity := range ideneityList {
			if identity.GetIdentityId() == msg.GetMetadataAuthority().GetOwner().GetIdentityId() {
				valid = true
				break
			}
		}
		if !valid {
			log.WithError(err).Errorf("not found identity with identityId of auth on MessageHandler with broadcast metadataAuth, userType: {%s}, user: {%s}, metadataId: {%s}",
				msg.GetUserType(), msg.GetUser(), msg.GetMetadataAuthority().GetMetadataId())
			continue
		}

		//todo: need checking...
		metadataList, err := m.resourceMng.GetDB().QueryMetadataListByIdentity(msg.GetMetadataAuthority().GetOwner().GetIdentityId(), timeutils.BeforeYearUnixMsecUint64(), backend.DefaultMaxPageSize)
		if nil != err {
			log.WithError(err).Errorf("Failed to query global metadata list by identityId on MessageHandler with broadcast metadataAuth, identityId: {%s}",
				msg.GetMetadataAuthority().GetOwner().GetIdentityId())
			return
		}

		// reset val
		valid = false
		for _, metadata := range metadataList {
			if metadata.GetData().GetMetadataId() == msg.GetMetadataAuthority().GetMetadataId() {
				valid = true
				break
			}
		}
		if !valid {
			log.WithError(err).Errorf("not found metadata with metadataId of auth on MessageHandler with broadcast metadataAuth, userType: {%s}, user: {%s}, metadataId: {%s}",
				msg.GetUserType(), msg.GetUser(), msg.GetMetadataAuthority().GetMetadataId())
			continue
		}

		// check the metadataId whether has valid metadataAuth with current userType and user.
		has, err := m.authManager.HasValidMetadataAuth(msg.GetUserType(), msg.GetUser(), msg.GetMetadataAuthorityOwnerIdentityId(), msg.GetMetadataAuthorityMetadataId())
		if nil != err {
			log.WithError(err).Errorf("Failed to call HasValidLocalMetadataAuth on MessageHandler with broadcast metadataAuth, metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}",
				msg.GetMetadataAuthId(), msg.GetMetadataAuthority().GetMetadataId(), msg.GetUserType(), msg.GetUser())
			continue
		}

		if has {
			log.Errorf("Failed to broadcast metadataAuth, cause alreay has valid last metadataAuth on MessageHandler with broadcast metadataAuth, userType: {%s}, user: {%s}, metadataId: {%s}",
				msg.GetUserType(), msg.GetUser(), msg.GetMetadataAuthority().GetMetadataId())
			continue
		}
		// TODO 填充 nonce
		// Store metadataAuthority
		if err := m.authManager.ApplyMetadataAuthority(types.NewMetadataAuthority(&libtypes.MetadataAuthorityPB{
			MetadataAuthId:  msg.GetMetadataAuthId(),
			User:            msg.GetUser(),
			UserType:        msg.GetUserType(),
			Auth:            msg.GetMetadataAuthority(),
			AuditOption:     libtypes.AuditMetadataOption_Audit_Pending,
			AuditSuggestion: "",
			UsedQuo: &libtypes.MetadataUsedQuo{
				UsageType: msg.GetMetadataAuthority().GetUsageRule().GetUsageType(),
				Expire:    false, // Initialized zero value
				UsedTimes: 0,     // Initialized zero value
			},
			ApplyAt: msg.GetCreateAt(),
			AuditAt: 0,
			State:   libtypes.MetadataAuthorityState_MAState_Released,
			Sign:    msg.GetSign(),
		})); nil != err {
			log.WithError(err).Errorf("Failed to store metadataAuth to dataCenter on MessageHandler with broadcast metadataAuth, metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}",
				msg.GetMetadataAuthId(), msg.GetMetadataAuthority().GetMetadataId(), msg.GetUserType(), msg.GetUser())
			continue
		}

		log.Debugf("broadcast metadataAuth msg succeed, metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}",
			msg.GetMetadataAuthId(), msg.GetMetadataAuthority().GetMetadataId(), msg.GetUserType(), msg.GetUser())
	}

	return
}

func (m *MessageHandler) BroadcastMetadataAuthRevokeMsgArr(metadataAuthRevokeMsgArr types.MetadataAuthorityRevokeMsgArr) {
	for _, revoke := range metadataAuthRevokeMsgArr {

		// checking ...
		metadataAuth, err := m.authManager.GetMetadataAuthority(revoke.GetMetadataAuthId())
		if nil != err {
			log.WithError(err).Errorf("Failed to query old metadataAuth on MessageHandler with revoke metadataAuth, metadataAuthId: {%s}, user:{%s}, userType: {%s}",
				revoke.GetMetadataAuthId(), revoke.GetUser(), revoke.GetUserType().String())
			continue
		}

		if metadataAuth.GetData().GetUser() != revoke.GetUser() || metadataAuth.GetData().GetUserType() != revoke.GetUserType() {
			log.Errorf("user of metadataAuth is wrong on MessageHandler with revoke metadataAuth, metadataAuthId: {%s}, user:{%s}, userType: {%s}",
				revoke.GetMetadataAuthId(), revoke.GetUser(), revoke.GetUserType().String())
			continue
		}

		// The data authorization application information that has been audited and cannot be revoked
		if metadataAuth.GetData().GetAuditOption() != libtypes.AuditMetadataOption_Audit_Pending {
			log.Errorf("the metadataAuth has audit on MessageHandler with revoke metadataAuth, metadataAuthId: {%s}, user:{%s}, state: {%s}",
				revoke.GetMetadataAuthId(), revoke.GetUser(), metadataAuth.GetData().GetAuditOption().String())
			continue
		}

		// The data authorization application information that has been `invalidated` or has been `revoked` is not allowed to be revoked
		if metadataAuth.GetData().GetState() != libtypes.MetadataAuthorityState_MAState_Released {
			log.Errorf("state of metadataAuth is wrong on MessageHandler with revoke metadataAuth, metadataAuthId: {%s}, user:{%s}, state: {%s}",
				revoke.GetMetadataAuthId(), revoke.GetUser(), metadataAuth.GetData().GetState().String())
			continue
		}

		// change state of metadataAuth from `release` to `revoke`
		metadataAuth.GetData().State = libtypes.MetadataAuthorityState_MAState_Revoked
		// update metadataAuth from datacenter
		if err := m.resourceMng.GetDB().UpdateMetadataAuthority(metadataAuth); nil != err {
			log.WithError(err).Errorf("Failed to update metadataAuth to dataCenter on MessageHandler with revoke metadataAuth, metadataAuthId: {%s}, user:{%s}",
				revoke.GetMetadataAuthId(), revoke.GetUser())
			continue
		}
		log.Debugf("revoke metadataAuth msg succeed, metadataAuthId: {%s}", revoke.GetMetadataAuthId())
	}

	return
}

func (m *MessageHandler) BroadcastTaskMsgArr(taskMsgArr types.TaskMsgArr) {
	// TODO 填充 nonce
	if err := m.taskManager.HandleTaskMsgs(taskMsgArr); nil != err {
		log.WithError(err).Errorf("Failed to call `BroadcastTaskMsgArr` on MessageHandler")
	}
	return
}

func (m *MessageHandler) BroadcastTaskTerminateMsgArr(terminateMsgArr types.TaskTerminateMsgArr) {
	if err := m.taskManager.HandleTaskTerminateMsgs(terminateMsgArr); nil != err {
		log.WithError(err).Errorf("Failed to call `BroadcastTaskTerminateMsgArr` on MessageHandler")
	}
	return
}
