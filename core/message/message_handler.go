package message

import (
	"encoding/json"
	"github.com/datumtechs/datum-network-carrier/ach/auth"
	"github.com/datumtechs/datum-network-carrier/carrierdb/rawdb"
	"github.com/datumtechs/datum-network-carrier/common/feed"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	"github.com/datumtechs/datum-network-carrier/core/resource"
	"github.com/datumtechs/datum-network-carrier/core/task"
	"github.com/datumtechs/datum-network-carrier/core/workflow"
	"github.com/datumtechs/datum-network-carrier/event"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/types"
	"sync"
	"time"
)

const (
	defaultPowerMsgsCacheSize          = 3
	defaultMetadataMsgsCacheSize       = 3
	defaultMetadataUpdateMsgsCacheSize = 3
	defaultMetadataAuthMsgsCacheSize   = 3
	defaultTaskMsgsCacheSize           = 5
	defaultWorkflowCacheSize           = 3

	defaultBroadcastPowerMsgInterval        = 30 * time.Second
	defaultBroadcastMetadataMsgInterval     = 30 * time.Second
	defaultBroadcastMetadataAuthMsgInterval = 30 * time.Second
	defaultBroadcastTaskMsgInterval         = 10 * time.Second
	defaultBroadCastWorkflowMsgInterval     = 10 * time.Second
)

type MessageHandler struct {
	pool        *Mempool
	resourceMng *resource.Manager
	// Send taskMsg to taskManager
	taskManager     *task.Manager
	authManager     *auth.AuthorityManager
	workflowManager *workflow.Manager
	// internal resource node set (Fighter node grpc client set)
	msgChannel chan *feed.Event
	quit       chan struct{}

	msgSub event.Subscription

	powerMsgCache          types.PowerMsgArr
	metadataMsgCache       types.MetadataMsgArr
	metadataAuthMsgCache   types.MetadataAuthorityMsgArr
	taskMsgCache           types.TaskMsgArr
	metadataUpdateMsgCache types.MetadataUpdateMsgArr
	workflowMsgCache       types.WorkflowMsgArr
	lockPower              sync.Mutex
	lockMetadata           sync.Mutex
	lockMetadataAuth       sync.Mutex
	lockTask               sync.Mutex
	lockMetadataUpdate     sync.Mutex
	lockWorkflow           sync.Mutex
}

func NewHandler(pool *Mempool, resourceMng *resource.Manager, taskManager *task.Manager, authManager *auth.AuthorityManager, workflowManager *workflow.Manager) *MessageHandler {
	m := &MessageHandler{
		pool:            pool,
		resourceMng:     resourceMng,
		taskManager:     taskManager,
		authManager:     authManager,
		workflowManager: workflowManager,
		msgChannel:      make(chan *feed.Event, 5),
		quit:            make(chan struct{}),
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

	metadataUpdateMsgCache, err := m.resourceMng.GetDB().QueryMetadataUpdateMsgArr()
	if nil != err {
		if rawdb.IsNoDBNotFoundErr(err) {
			log.WithError(err).Errorf("Failed to get metadataUpdateMsgCache from QueryMetadataUpdateMsgArr")
		} else {
			log.WithError(err).Warnf("Not found metadataUpdateMsgCache from QueryMetadataUpdateMsgArr")
		}
	} else {
		m.metadataUpdateMsgCache = metadataUpdateMsgCache
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
	workflowMsgCache, err := m.resourceMng.GetDB().QueryWorkflowMsgArr()
	if nil != err {
		if rawdb.IsNoDBNotFoundErr(err) {
			log.WithError(err).Errorf("Failed to get workflowMsgCache from QueryWorkflowMsgArr")
		} else {
			log.WithError(err).Warnf("Not found workflowMsgCache from QueryWorkflowMsgArr")
		}
	} else {
		m.workflowMsgCache = workflowMsgCache
	}
}
func (m *MessageHandler) loop() {
	powerTicker := time.NewTicker(defaultBroadcastPowerMsgInterval)
	metadataTicker := time.NewTicker(defaultBroadcastMetadataMsgInterval)
	metadataAuthTicker := time.NewTicker(defaultBroadcastMetadataAuthMsgInterval)
	taskTicker := time.NewTicker(defaultBroadcastTaskMsgInterval)
	workflowTicker := time.NewTicker(defaultBroadCastWorkflowMsgInterval)
	for {
		select {
		case event := <-m.msgChannel:
			switch event.Type {
			case types.ApplyIdentity:
				eventMessage := event.Data.(*types.IdentityMsgEvent)
				m.BroadcastIdentityMsg(eventMessage.Msg)
			case types.UpdateIdentityCredential:
				eventMessage := event.Data.(*types.UpdateIdentityCredentialEvent)
				m.BroadcastUpdateIdentityCredentialMsg(eventMessage.Msg)
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
			case types.UpdateMetadata:
				msg := event.Data.(*types.MetadataUpdateMsgEvent)
				m.lockMetadataUpdate.Lock()
				m.metadataUpdateMsgCache = append(m.metadataUpdateMsgCache, msg.Msg)
				m.resourceMng.GetDB().StoreMessageCache(msg.Msg) // backup metadata msg into disk
				if len(m.metadataUpdateMsgCache) >= defaultMetadataUpdateMsgsCacheSize {
					m.BroadcastMetadataUpdateMsgArr(m.metadataUpdateMsgCache)
					m.metadataUpdateMsgCache = make(types.MetadataUpdateMsgArr, 0)
				}
				m.lockMetadataUpdate.Unlock()
			case types.ApplyWorkflow:
				msg := event.Data.(*types.WorkflowMsgEvent)
				m.lockWorkflow.Lock()
				m.workflowMsgCache = append(m.workflowMsgCache, msg.Msg)
				m.resourceMng.GetDB().StoreMessageCache(msg.Msg)
				if len(m.workflowMsgCache) >= defaultWorkflowCacheSize {
					m.BroadcastWorkflowMsgArr(m.workflowMsgCache)
					m.workflowMsgCache = make(types.WorkflowMsgArr, 0)
				}
				m.lockWorkflow.Unlock()
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
			if len(m.metadataUpdateMsgCache) > 0 {
				m.BroadcastMetadataUpdateMsgArr(m.metadataUpdateMsgCache)
				m.metadataUpdateMsgCache = make(types.MetadataUpdateMsgArr, 0)
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
		case <-workflowTicker.C:
			if len(m.workflowMsgCache) > 0 {
				m.BroadcastWorkflowMsgArr(m.workflowMsgCache)
				m.workflowMsgCache = make(types.WorkflowMsgArr, 0)
			}
		case taskMsg := <-m.workflowManager.TaskMsgToMessageManagerCh:
			log.Debugf("come from workflowManager task,taskId is {%s}", taskMsg.GetTaskData().GetTaskId())
			m.lockTask.Lock()
			m.taskMsgCache = append(m.taskMsgCache, taskMsg)
			m.resourceMng.GetDB().StoreMessageCache(taskMsg) // backup task msg into disk
			m.lockTask.Unlock()
		// Err() channel will be closed when unsubscribing.
		case err := <-m.msgSub.Err():
			log.Errorf("Received err from msgSub, return loop, err: %s", err)
			return
		case <-m.quit:
			log.Infof("Stopped message handler ...")
			powerTicker.Stop()
			metadataTicker.Stop()
			metadataAuthTicker.Stop()
			taskTicker.Stop()
			return
		}
	}
}

func (m *MessageHandler) BroadcastIdentityMsg(msg *types.IdentityMsg) {

	// add identity to local db
	identity := msg.GetOrganization()
	identity.DataStatus = commonconstantpb.DataStatus_DataStatus_Valid
	identity.Status = commonconstantpb.CommonStatus_CommonStatus_Valid

	// set new nonce into msg.
	nonce, err := m.resourceMng.GetDB().IncreaseIdentityMsgNonce()
	if nil != err {
		log.WithError(err).Errorf("Failed to increase identity msg nonce on MessageHandler with broadcast identity msg, identityId: {%s}, nodeId: {%s}, nodeName: {%s}",
			msg.GetOwnerIdentityId(), msg.GetOwnerNodeId(), msg.GetOwnerName())
		return
	}
	identity.Nonce = nonce              // add by v0.5.0
	msg.GetOrganization().Nonce = nonce // add by v0.5.0

	if err := m.resourceMng.GetDB().StoreIdentity(identity); nil != err {
		log.WithError(err).Errorf("Failed to store local org identity on MessageHandler with broadcast identity msg, identityId: {%s}", msg.GetOwnerIdentityId())
		return
	}
	// send identity to datacenter
	if err := m.resourceMng.GetDB().InsertIdentity(msg.ToDataCenter()); nil != err {
		log.WithError(err).Errorf("Failed to broadcast org org identity on MessageHandler with broadcast identity msg, identityId: {%s}, nodeId: {%s}, nodeName: {%s}",
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
		log.WithError(err).Errorf("Failed to get local org identity on MessageHandler with revoke identity msg, identityId: {%s}", identity.GetIdentityId())
		return
	}

	// what if running task, can not revoke identity
	jobNodes, err := m.resourceMng.GetDB().QueryRegisterNodeList(carrierapipb.PrefixTypeJobNode)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("query all jobNode failed")
		return
	}
	for _, node := range jobNodes {
		runningTaskCount, err := m.resourceMng.GetDB().QueryJobNodeRunningTaskCount(node.Id)
		if rawdb.IsNoDBNotFoundErr(err) {
			log.Errorf("query local running taskCount on old jobNode failed on MessageHandler with revoke identity msg, %s", err)
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
		log.WithError(err).Errorf("Failed to delete org identity to local on MessageHandler with revoke identity msg, identityId: {%s}", identity.GetIdentityId())
		return
	}

	// remove identity from dataCenter
	if err := m.resourceMng.GetDB().RevokeIdentity(
		types.NewIdentity(&carriertypespb.IdentityPB{
			NodeName:   identity.GetNodeName(),
			NodeId:     identity.GetNodeId(),
			IdentityId: identity.GetIdentityId(),
			DataId:     "",
			DataStatus: commonconstantpb.DataStatus_DataStatus_Invalid,
			Status:     commonconstantpb.CommonStatus_CommonStatus_Invalid,
			Credential: "",
		})); nil != err {
		log.WithError(err).Errorf("Failed to remove org identity to remote on MessageHandler with revoke identity msg, identityId: {%s}", identity.GetIdentityId())
		return
	}
	log.Debugf("Revoke identity msg succeed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}", identity.GetIdentityId(), identity.GetNodeId(), identity.GetNodeName())
	return
}
func (m *MessageHandler) BroadcastUpdateIdentityCredentialMsg(msg *types.UpdateIdentityCredentialMsg) {
	// query local identity
	identity, err := m.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("query local identity failed on MessageHandler with broadcast identity update credential msg")
		return
	}
	if identity.GetIdentityId() != msg.IdentityId {
		log.Warnf("local identityId %s not equal msg identityId %s on MessageHandler with broadcast identity update credential msg", msg.IdentityId, identity.GetIdentityId())
		return
	}
	err = m.resourceMng.GetDB().UpdateIdentityCredential(msg.IdentityId, msg.Credential)
	if err != nil {
		log.WithError(err).Errorf("call UpdateIdentityCredential() failed on MessageHandler with broadcast identity update credential msg")
		return
	}
	identity.Credential = msg.Credential
	err = m.resourceMng.GetDB().StoreIdentity(identity)
	if err != nil {
		log.WithError(err).Errorf("call StoreIdentity failed on MessageHandler with broadcast identity update credential msg")
		return
	}
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
			log.WithError(err).Errorf("Failed to query local resource on MessageHandler with broadcast power msg, powerId: {%s}, jobNodeId: {%s}",
				msg.GetPowerId(), msg.GetJobNodeId())
			continue
		}

		// set powerId to resource and change state
		resource.GetData().DataId = msg.GetPowerId()
		resource.GetData().State = commonconstantpb.PowerState_PowerState_Released

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
			log.WithError(err).Errorf("Failed to StoreLocalResourceTable on MessageHandler with broadcast power msg, powerId: {%s}, jobNodeId: {%s}",
				msg.GetPowerId(), msg.GetJobNodeId())
			continue
		}

		if err := m.resourceMng.GetDB().StoreJobNodeIdIdByPowerId(msg.GetPowerId(), msg.GetJobNodeId()); nil != err {
			log.WithError(err).Errorf("Failed to StoreJobNodeIdByPowerId on MessageHandler with broadcast power msg,  powerId: {%s}, jobNodeId: {%s}",
				msg.GetPowerId(), msg.GetJobNodeId())
			continue
		}

		// set new nonce into msg.
		nonce, err := m.resourceMng.GetDB().IncreasePowerMsgNonce()
		if nil != err {
			log.WithError(err).Errorf("Failed to increase power msg nonce on MessageHandler with broadcast power msg, powerId: {%s}, jobNodeId: {%s}",
				msg.GetPowerId(), msg.GetJobNodeId())
			continue
		}
		resource.GetData().Nonce = nonce // add by v0.5.0
		// update local resource
		if err := m.resourceMng.GetDB().StoreLocalResource(resource); nil != err {
			log.WithError(err).Errorf("Failed to update local resource with powerId to local on MessageHandler with broadcast power msg, powerId: {%s}, jobNodeId: {%s}",
				msg.GetPowerId(), msg.GetJobNodeId())
			continue
		}
		// publish to global
		if err := m.resourceMng.GetDB().InsertResource(types.NewResource(&carriertypespb.ResourcePB{
			Owner:  identity,
			DataId: msg.GetPowerId(),
			// the status of data for local storage, 1 means valid, 2 means invalid
			DataStatus: commonconstantpb.DataStatus_DataStatus_Valid,
			// resource status, eg: create/release/revoke
			State: commonconstantpb.PowerState_PowerState_Released,
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
			Nonce:          nonce, // add by v0.5.0
		})); nil != err {
			log.WithError(err).Errorf("Failed to store power msg to dataCenter on MessageHandler with broadcast power msg, powerId: {%s}, jobNodeId: {%s}",
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
			log.WithError(err).Errorf("Failed to call QueryJobNodeIdByPowerId() on MessageHandler with revoke power msg, powerId: {%s}, jobNodeId: {%s}",
				revoke.GetPowerId(), jobNodeId)
			continue
		}
		if err := m.resourceMng.GetDB().RemoveJobNodeIdByPowerId(revoke.GetPowerId()); nil != err {
			log.WithError(err).Errorf("Failed to call RemoveJobNodeIdByPowerId() on MessageHandler with revoke power msg, powerId: {%s}, jobNodeId: {%s}",
				revoke.GetPowerId(), jobNodeId)
			continue
		}
		if err := m.resourceMng.GetDB().RemoveLocalResourceTable(jobNodeId); nil != err {
			log.WithError(err).Errorf("Failed to RemoveLocalResourceTable on MessageHandler with revoke power msg, powerId: {%s}, jobNodeId: {%s}",
				revoke.GetPowerId(), jobNodeId)
			continue
		}

		// query local resource
		resource, err := m.resourceMng.GetDB().QueryLocalResource(jobNodeId)
		if nil != err {
			log.WithError(err).Errorf("Failed to query local resource on MessageHandler with revoke power msg, powerId: {%s}, jobNodeId: {%s}",
				revoke.GetPowerId(), jobNodeId)
			continue
		}

		// remove powerId from local resource and change state
		resource.GetData().DataId = ""
		resource.GetData().State = commonconstantpb.PowerState_PowerState_Revoked
		// clean used resource value
		resource.GetData().UsedBandwidth = 0
		resource.GetData().UsedDisk = 0
		resource.GetData().UsedProcessor = 0
		resource.GetData().UsedMem = 0

		// update local resource
		if err := m.resourceMng.GetDB().StoreLocalResource(resource); nil != err {
			log.WithError(err).Errorf("Failed to update local resource with powerId to local on MessageHandler with revoke power msg, powerId: {%s}, jobNodeId: {%s}",
				revoke.GetPowerId(), jobNodeId)
			continue
		}

		// remove from global
		if err := m.resourceMng.GetDB().RevokeResource(types.NewResource(&carriertypespb.ResourcePB{
			Owner:  identity,
			DataId: revoke.GetPowerId(),
			// the status of data for local storage, 1 means valid, 2 means invalid
			DataStatus: commonconstantpb.DataStatus_DataStatus_Invalid,
			// resource status, eg: create/release/revoke
			State:    commonconstantpb.PowerState_PowerState_Revoked,
			UpdateAt: timeutils.UnixMsecUint64(),
		})); nil != err {
			log.WithError(err).Errorf("Failed to remove dataCenter resource on MessageHandler with revoke power msg, powerId: {%s}, jobNodeId: {%s}",
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
				log.WithError(err).Errorf("Failed to unmashal metadataOption on MessageHandler with broadcast metadata msg, metadataId: {%s}",
					msg.GetMetadataId())
				continue
			}

			// maintain the orginId and metadataId relationship of the local data service
			dataResourceDataUpload, err := m.resourceMng.GetDB().QueryDataResourceDataUpload(option.GetOriginId())
			if nil != err {
				log.WithError(err).Errorf("Failed to QueryDataResourceDataUpload on MessageHandler with broadcast metadata msg, originId: {%s}, metadataId: {%s}",
					option.GetOriginId(), msg.GetMetadataId())
				continue
			}

			// Update metadataId in fileupload information
			dataResourceDataUpload.SetMetadataId(msg.GetMetadataId())
			if err := m.resourceMng.GetDB().StoreDataResourceDataUpload(dataResourceDataUpload); nil != err {
				log.WithError(err).Errorf("Failed to StoreDataResourceDataUpload on MessageHandler with broadcast metadata msg, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}",
					option.GetOriginId(), msg.GetMetadataId(), dataResourceDataUpload.GetNodeId())
				continue
			}

			// Record the size of the resources occupied by the original data
			dataResourceTable, err := m.resourceMng.GetDB().QueryDataResourceTable(dataResourceDataUpload.GetNodeId())
			if nil != err {
				log.WithError(err).Errorf("Failed to QueryDataResourceTable on MessageHandler with broadcast metadata msg, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}",
					option.GetOriginId(), msg.GetMetadataId(), dataResourceDataUpload.GetNodeId())
				continue
			}
			// update disk used of data resource table
			dataResourceTable.UseDisk(option.GetSize())
			if err := m.resourceMng.GetDB().StoreDataResourceTable(dataResourceTable); nil != err {
				log.WithError(err).Errorf("Failed to StoreDataResourceTable on MessageHandler with broadcast metadata msg, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}",
					option.GetOriginId(), msg.GetMetadataId(), dataResourceDataUpload.GetNodeId())
				continue
			}

			// Separately record the GetSize of the metaData and the dataNodeId where it is located
			if err := m.resourceMng.GetDB().StoreDataResourceDiskUsed(types.NewDataResourceDiskUsed(
				msg.GetMetadataId(), dataResourceDataUpload.GetNodeId(), option.GetSize())); nil != err {
				log.WithError(err).Errorf("Failed to StoreDataResourceDiskUsed on MessageHandler with broadcast metadata msg, originId: {%s}, metadataId: {%s}, dataNodeId: {%s}",
					option.GetOriginId(), msg.GetMetadataId(), dataResourceDataUpload.GetNodeId())
				continue
			}
		}

		// set new nonce into msg.
		nonce, err := m.resourceMng.GetDB().IncreaseMetadataMsgNonce()
		if nil != err {
			log.WithError(err).Errorf("Failed to increase metadata msg nonce on MessageHandler with broadcast metadata msg, metadataId: {%s}",
				msg.GetMetadataId())
			continue
		}
		msg.GetMetadataSummary().Nonce = nonce // add by v0.5.0
		// publish msg information
		if err := m.resourceMng.GetDB().InsertMetadata(msg.ToDataCenter(identity)); nil != err {
			log.WithError(err).Errorf("Failed to store msg to dataCenter on MessageHandler with broadcast metadata msg, metadataId: {%s}",
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

func (m *MessageHandler) BroadcastMetadataUpdateMsgArr(metadataUpdateMsgArr types.MetadataUpdateMsgArr) {
	identity, err := m.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to query local identity on MessageHandler with broadcast metadata update msg")
		return
	}
	for _, msg := range metadataUpdateMsgArr {
		switch msg.GetMetadataSummary().GetDataType() {
		case commonconstantpb.OrigindataType_OrigindataType_CSV:
			var option *types.MetadataOptionCSV
			if err := json.Unmarshal([]byte(msg.GetMetadataSummary().MetadataOption), &option); nil != err {
				log.WithError(err).Errorf("Failed to unmashal metadataOption on MessageHandler with broadcast metadata update msg, metadataId: {%s}",
					msg.GetMetadataId())
				continue
			}
			newMetadata := types.NewMetadata(&carriertypespb.MetadataPB{
				MetadataId:     msg.GetMetadataId(),
				Owner:          identity,
				DataId:         msg.GetMetadataId(),
				DataStatus:     commonconstantpb.DataStatus_DataStatus_Valid,
				MetadataName:   msg.GetMetadataName(),
				MetadataType:   msg.GetMetadataType(),
				DataHash:       msg.GetDataHash(),
				Desc:           msg.GetDesc(),
				LocationType:   msg.GetLocationType(),
				DataType:       msg.GetDataType(),
				Industry:       msg.GetIndustry(),
				State:          commonconstantpb.MetadataState_MetadataState_Released, // metaData status, eg: create/release/revoke
				PublishAt:      msg.GetPublishAt(),
				UpdateAt:       timeutils.UnixMsecUint64(),
				Nonce:          msg.GetNonce(),
				MetadataOption: msg.GetMetadataOption(),
				User:           msg.GetUser(),
				UserType:       msg.GetUserType(),
			})
			// set new nonce into msg.
			nonce, err := m.resourceMng.GetDB().IncreaseMetadataMsgNonce()
			if nil != err {
				log.WithError(err).Errorf("Failed to increase metadata msg nonce on MessageHandler with broadcast metadata update msg, metadataId: {%s}",
					msg.GetMetadataId())
				continue
			}
			newMetadata.GetData().Nonce = nonce // add by v0.5.0
			if err := m.resourceMng.GetDB().UpdateGlobalMetadata(newMetadata); nil != err {
				log.WithError(err).Errorf("Failed to store msg to dataCenter on MessageHandler with broadcast metadata update msg, metadataId: {%s}",
					msg.GetMetadataId())
				continue
			}

			log.Debugf("broadcast metadata update msg succeed, metadataId: {%s}", msg.GetMetadataId())
			m.resourceMng.GetDB().RemoveMetadataUpdateMsg(msg.GetMetadataId()) // remove from disk if msg been handle
		default:
			log.Warnf("BroadcastMetadataUpdateMsgArr Unknown type %s found", msg.GetMetadataSummary().GetDataType().String())
		}
	}
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
			log.WithError(err).Errorf("Failed to QueryDataResourceDiskUsed on MessageHandler with revoke metadata msg, metadataId: {%s}",
				revoke.GetMetadataId())
			continue
		}
		// update dataNode table (dataNodeId -> {dataNodeId, totalDisk, usedDisk})
		dataResourceTable, err := m.resourceMng.GetDB().QueryDataResourceTable(dataResourceDiskUsed.GetNodeId())
		if nil != err {
			log.WithError(err).Errorf("Failed to QueryDataResourceTable on MessageHandler with revoke metadata msg, metadataId: {%s}, dataNodeId: {%s}",
				revoke.GetMetadataId(), dataResourceDiskUsed.GetNodeId())
			continue
		}
		dataResourceTable.FreeDisk(dataResourceDiskUsed.GetDiskUsed())
		if err := m.resourceMng.GetDB().StoreDataResourceTable(dataResourceTable); nil != err {
			log.WithError(err).Errorf("Failed to StoreDataResourceTable on MessageHandler with revoke metadata msg, metadataId: {%s}, dataNodeId: {%s}",
				revoke.GetMetadataId(), dataResourceDiskUsed.GetNodeId())
			continue
		}

		// remove dataNodeDiskUsed (metaDataId -> {metaDataId, dataNodeId, diskUsed})
		if err := m.resourceMng.GetDB().RemoveDataResourceDiskUsed(revoke.GetMetadataId()); nil != err {
			log.WithError(err).Errorf("Failed to RemoveDataResourceDiskUsed on MessageHandler with revoke metadata msg, metadataId: {%s}, dataNodeId: {%s}",
				revoke.GetMetadataId(), dataResourceDiskUsed.GetNodeId())
			continue
		}

		// revoke from global
		if err := m.resourceMng.GetDB().RevokeMetadata(revoke.ToDataCenter(identity)); nil != err {
			log.WithError(err).Errorf("Failed to revoke metadata to dataCenter on MessageHandler with revoke metadata msg, metadataId: {%s}",
				revoke.GetMetadataId())
			continue
		}
		log.Debugf("revoke metadata msg succeed, metadataId: {%s}", revoke.GetMetadataId())
	}
	return
}

func (m *MessageHandler) BroadcastMetadataAuthMsgArr(metadataAuthMsgArr types.MetadataAuthorityMsgArr) {

	for _, msg := range metadataAuthMsgArr {

		// remove from disk if msg been handle
		m.resourceMng.GetDB().RemoveMetadataAuthMsg(msg.GetMetadataAuthId())

		// ############################################
		// ############################################
		// NOTE:
		// 		check whether is valid with metadata?
		//
		//		Because the message is asynchronous,
		//		the corresponding identity and metadata may have been deleted long ago,
		//		so the data legitimacy is verified again
		// ############################################
		// ############################################
		auth := types.NewMetadataAuthority(&carriertypespb.MetadataAuthorityPB{
			MetadataAuthId:  msg.GetMetadataAuthId(),
			User:            msg.GetUser(),
			UserType:        msg.GetUserType(),
			Auth:            msg.GetMetadataAuthority(),
			AuditOption:     commonconstantpb.AuditMetadataOption_Audit_Pending,
			AuditSuggestion: "",
			UsedQuo: &carriertypespb.MetadataUsedQuo{
				UsageType: msg.GetMetadataAuthority().GetUsageRule().GetUsageType(),
				Expire:    false, // Initialized zero value
				UsedTimes: 0,     // Initialized zero value
			},
			ApplyAt:   msg.GetCreateAt(),
			AuditAt:   0,
			State:     commonconstantpb.MetadataAuthorityState_MAState_Released,
			Sign:      msg.GetSign(),
			PublishAt: timeutils.UnixMsecUint64(),
			UpdateAt:  timeutils.UnixMsecUint64(),
		})
		pass, err := m.authManager.VerifyMetadataAuthWithMetadataOption(auth)
		if nil != err {
			log.WithError(err).Errorf("Failed to verify metadataAuth with metadataOption on MessageHandler with broadcast metadataAuth msg, metadataId: {%s}, auth: %s",
				msg.GetMetadataAuthority().GetMetadataId(), msg.GetMetadataAuthority().String())
			continue
		}
		if !pass {
			log.Errorf("invalid metadataAuth on MessageHandler with broadcast metadataAuth msg, metadataId: {%s}, auth: %s",
				msg.GetMetadataAuthority().GetMetadataId(), msg.GetMetadataAuthority().String())
			continue
		}
		// set new nonce into msg.
		nonce, err := m.resourceMng.GetDB().IncreaseMetadataAuthMsgNonce()
		if nil != err {
			log.WithError(err).Errorf("Failed to increase metadataAuth msg nonce on MessageHandler with broadcast metadataAuth msg, metadataAuthId: {%s}",
				msg.GetMetadataAuthId())
			continue
		}
		auth.GetData().Nonce = nonce // add by v0.5.0
		// Store metadataAuthority
		if err := m.authManager.ApplyMetadataAuthority(auth); nil != err {
			log.WithError(err).Errorf("Failed to store metadataAuth to dataCenter on MessageHandler with broadcast metadataAuth msg, metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}",
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
			log.WithError(err).Errorf("Failed to query old metadataAuth on MessageHandler with revoke metadataAuth msg, metadataAuthId: {%s}, user:{%s}, userType: {%s}",
				revoke.GetMetadataAuthId(), revoke.GetUser(), revoke.GetUserType().String())
			continue
		}

		if metadataAuth.GetData().GetUser() != revoke.GetUser() || metadataAuth.GetData().GetUserType() != revoke.GetUserType() {
			log.Errorf("user of metadataAuth is wrong on MessageHandler with revoke metadataAuth msg, metadataAuthId: {%s}, user:{%s}, userType: {%s}",
				revoke.GetMetadataAuthId(), revoke.GetUser(), revoke.GetUserType().String())
			continue
		}

		// The data authorization application information that has been audited and cannot be revoked
		pass, err := m.authManager.VerifyMetadataAuthInfo(metadataAuth, false, true)
		if nil != err {
			log.WithError(err).Errorf("Failed to verify metadataAuth on MessageHandler with revoke metadataAuth msg, metadataAuthId: {%s}, user:{%s}, state: {%s}",
				revoke.GetMetadataAuthId(), revoke.GetUser(), metadataAuth.GetData().GetAuditOption().String())
			continue
		}
		if !pass {
			log.Errorf("invalid metadataAuth on MessageHandler with revoke metadataAuth msg, metadataAuthId: {%s}, user:{%s}, state: {%s}",
				revoke.GetMetadataAuthId(), revoke.GetUser(), metadataAuth.GetData().GetAuditOption().String())
			continue
		}

		// change state of metadataAuth from `release` to `revoke`
		metadataAuth.GetData().State = commonconstantpb.MetadataAuthorityState_MAState_Revoked
		// update metadataAuth from datacenter
		if err := m.resourceMng.GetDB().UpdateMetadataAuthority(metadataAuth); nil != err {
			log.WithError(err).Errorf("Failed to update metadataAuth to dataCenter on MessageHandler with revoke metadataAuth msg, metadataAuthId: {%s}, user:{%s}",
				revoke.GetMetadataAuthId(), revoke.GetUser())
			continue
		}
		log.Debugf("revoke metadataAuth msg succeed, metadataAuthId: {%s}", revoke.GetMetadataAuthId())
	}

	return
}

func (m *MessageHandler) BroadcastTaskMsgArr(taskMsgArr types.TaskMsgArr) {
	// set new nonce into msg.
	for i, msg := range taskMsgArr {
		nonce, err := m.resourceMng.GetDB().IncreaseTaskMsgNonce()
		if nil != err {
			log.WithError(err).Errorf("Failed to increase task msg nonce on MessageHandler, taskId: {%s}",
				msg.GetTaskData().GetTaskId())
			continue
		}
		msg.GetTaskData().Nonce = nonce // add by v0.5.0
		taskMsgArr[i] = msg
	}
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

func (m *MessageHandler) BroadcastWorkflowMsgArr(workflowMsgArr types.WorkflowMsgArr) {
	for _, v := range workflowMsgArr {
		err := m.workflowManager.AddWorkflow(v.Data)
		if err != nil {
			log.WithError(err).Errorf("Failed to call `BroadcastWorkflowMsgArr` on MessageHandler")
		} else {
			if err := m.resourceMng.GetDB().RemoveWorkflowMsg(v.Data.GetWorkflowId()); err != nil {
				log.WithError(err).Errorf("Remove local workflowMsg fail,workflowId is %s", v.Data.GetWorkflowId())
			}
		}
	}
}
