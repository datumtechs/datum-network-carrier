package task

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/ach/auth"
	"github.com/datumtechs/datum-network-carrier/ach/token"
	"github.com/datumtechs/datum-network-carrier/carrierdb/rawdb"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/common/bytesutil"
	"github.com/datumtechs/datum-network-carrier/common/signutil"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	"github.com/datumtechs/datum-network-carrier/common/traceutil"
	"github.com/datumtechs/datum-network-carrier/consensus"
	ev "github.com/datumtechs/datum-network-carrier/core/evengine"
	"github.com/datumtechs/datum-network-carrier/core/policy"
	"github.com/datumtechs/datum-network-carrier/core/resource"
	"github.com/datumtechs/datum-network-carrier/core/schedule"
	"github.com/datumtechs/datum-network-carrier/p2p"
	"github.com/datumtechs/datum-network-carrier/params"
	carriernetmsgcommonpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/common"
	carriertwopcpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/consensus/twopc"
	carriernetmsgtaskmngpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/taskmng"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"
	"strings"
	"sync"
	"time"
)

const (
	defaultScheduleTaskInterval = 2 * time.Second
	senderExecuteTaskExpire     = 10 * time.Second
)

type Manager struct {
	nodePriKey      *ecdsa.PrivateKey
	p2p             p2p.P2P
	scheduler       schedule.Scheduler
	consensusEngine consensus.Engine
	eventEngine     *ev.EventEngine
	resourceMng     *resource.Manager
	authMng         *auth.AuthorityManager
	token20PayMng   *token.Token20PayManager
	parser          *TaskParser
	validator       *TaskValidator
	eventCh         chan *carriertypespb.TaskEvent

	// send the validated taskMsgs to scheduler
	localTasksCh             chan types.TaskDataArray
	needReplayScheduleTaskCh chan *types.NeedReplayScheduleTask
	needExecuteTaskCh        chan *types.NeedExecuteTask
	runningTaskCache         map[string]map[string]*types.NeedExecuteTask //  taskId -> {partyId -> task}
	syncExecuteTaskMonitors  *types.SyncExecuteTaskMonitorQueue
	quit                     chan struct{}
	runningTaskCacheLock     sync.RWMutex
	// add by v 0.4.0
	config   *params.TaskManagerConfig
	msgCache *TaskmngMsgCache
	policyEngine *policy.PolicyEngine
}

func NewTaskManager(
	nodePriKey *ecdsa.PrivateKey,
	p2p p2p.P2P,
	scheduler schedule.Scheduler,
	consensusEngine consensus.Engine,
	eventEngine *ev.EventEngine,
	resourceMng *resource.Manager,
	authMng *auth.AuthorityManager,
	token20PayMng *token.Token20PayManager,
	needReplayScheduleTaskCh chan *types.NeedReplayScheduleTask,
	needExecuteTaskCh chan *types.NeedExecuteTask,
	config *params.TaskManagerConfig,
	policyEngine *policy.PolicyEngine,
) (*Manager, error) {

	cache, err := NewTaskmngMsgCache(defaultTaskmngMsgCacheSize)
	if nil != err {
		return nil, err
	}

	m := &Manager{
		nodePriKey:               nodePriKey,
		p2p:                      p2p,
		scheduler:                scheduler,
		consensusEngine:          consensusEngine,
		eventEngine:              eventEngine,
		resourceMng:              resourceMng,
		authMng:                  authMng,
		token20PayMng:            token20PayMng,
		parser:                   newTaskParser(resourceMng),
		validator:                newTaskValidator(resourceMng, authMng),
		eventCh:                  make(chan *carriertypespb.TaskEvent, 800),
		needReplayScheduleTaskCh: needReplayScheduleTaskCh,
		needExecuteTaskCh:        needExecuteTaskCh,
		config:                   config,
		policyEngine:             policyEngine,
		runningTaskCache:         make(map[string]map[string]*types.NeedExecuteTask, 0), // taskId -> partyId -> needExecuteTask
		syncExecuteTaskMonitors:  types.NewSyncExecuteTaskMonitorQueue(0),
		quit:                     make(chan struct{}),
		msgCache:                 cache,
	}
	return m, nil
}
func (m *Manager) recoveryNeedExecuteTask() {
	prefix := rawdb.GetNeedExecuteTaskKeyPrefix()
	if err := m.resourceMng.GetDB().ForEachNeedExecuteTask(func(key, value []byte) error {
		if len(key) != 0 && len(value) != 0 {

			// task:${taskId hex} == len(5 + 2 + 64) == "taskId:" + "0x" + "e33...fe4"
			taskId := string(key[len(prefix) : len(prefix)+71])
			partyId := string(key[len(prefix)+71:])

			var res carriertypespb.NeedExecuteTask

			if err := proto.Unmarshal(value, &res); nil != err {
				//return fmt.Errorf("Unmarshal needExecuteTask failed, %s", err)
				log.WithError(err).Errorf("Unmarshal needExecuteTask failed, taskId: {%s}, partyId: {%s}", taskId, partyId)
				return nil
			}

			localTask, err := m.resourceMng.GetDB().QueryLocalTask(taskId)
			if nil != err {
				//return fmt.Errorf("found local task about needExecuteTask failed, %s", err)
				log.WithError(err).Errorf("found local task about needExecuteTask failed, taskId: {%s}, partyId: {%s}", taskId, partyId)
				return nil
			}

			cache, ok := m.runningTaskCache[taskId]
			if !ok {
				cache = make(map[string]*types.NeedExecuteTask, 0)
			}
			var taskErr error
			if strings.Trim(res.GetErrStr(), "") != "" {
				taskErr = fmt.Errorf(strings.Trim(res.GetErrStr(), ""))
			}

			task := types.NewNeedExecuteTask(
				peer.ID(res.GetRemotePid()),
				res.GetLocalTaskRole(),
				res.GetRemoteTaskRole(),
				res.GetLocalTaskOrganization(),
				res.GetRemoteTaskOrganization(),
				taskId,
				types.TaskActionStatus(bytesutil.BytesToUint16(res.GetConsStatus())),
				types.NewPrepareVoteResource(
					res.GetLocalResource().GetId(),
					res.GetLocalResource().GetIp(),
					res.GetLocalResource().GetPort(),
					res.GetLocalResource().GetPartyId(),
				),
				res.GetResources(),
				taskErr,
			)
			// TODO 需要添加根据 consumeSpec 中的一些状态， 决定是否在重新启动时，是否继续做的某些事 ...
			task.SetConsumeQueryId(res.GetConsumeQueryId())
			task.SetConsumeSpec(res.GetConsumeSpec())

			cache[partyId] = task
			m.runningTaskCache[taskId] = cache

			// v 0.3.0 add NeedExecuteTask Expire Monitor
			m.addmonitor(task, int64(localTask.GetTaskData().GetStartAt()+localTask.GetTaskData().GetOperationCost().GetDuration()))
		}
		return nil
	}); nil != err {
		log.WithError(err).Fatalf("recover needExecuteTask failed")
	}
}
func (m *Manager) Start() error {
	m.recoveryNeedExecuteTask()
	go m.loop()
	log.Info("Started taskManager ...")
	return nil
}
func (m *Manager) Stop() error {
	close(m.quit)
	return nil
}

func (m *Manager) loop() {

	taskTicker := time.NewTicker(defaultScheduleTaskInterval) // 2 s
	taskMonitorTimer := m.needExecuteTaskMonitorTimer()

	for {
		select {

		// handle reported event from fighter of self organization
		case event := <-m.eventCh:
			go func(event *carriertypespb.TaskEvent) {
				if err := m.handleTaskEventWithCurrentOranization(event); nil != err {
					log.WithError(err).Errorf("Failed to call handleTaskEventWithCurrentOranization() on taskManager.loop(), taskId: {%s}, event: %s", event.GetTaskId(), event.String())
				}
			}(event)


		// To schedule local task interval
		case <-taskTicker.C:

			if err := m.tryScheduleTask(); nil != err {
				log.WithError(err).Errorf("Failed to try schedule local task when taskTicker")
			}

		// handle the task of need replay scheduling while received from remote peer on consensus epoch
		case task := <-m.needReplayScheduleTaskCh:

			if nil == task {
				continue
			}
			go m.handleNeedReplayScheduleTask(task)

		// handle task of need to executing, and send it to fighter of myself organization or send the task result msg to remote peer
		case task := <-m.needExecuteTaskCh:

			if nil == task {
				continue
			}
			go m.handleNeedExecuteTask(task)

		// handle the executing expire tasks
		case <-taskMonitorTimer.C:

			future := m.checkNeedExecuteTaskMonitors(timeutils.UnixMsec(), true)
			now := timeutils.UnixMsec()
			if future > now {
				taskMonitorTimer.Reset(time.Duration(future-now) * time.Millisecond)
			} else if future < now {
				taskMonitorTimer.Reset(time.Duration(now) * time.Millisecond)
			}
			// when future value is 0, we do nothing

		case <-m.quit:
			log.Info("Stopped taskManager ...")
			taskTicker.Stop()
			taskMonitorTimer.Stop()
			return
		}
	}
}

func (m *Manager) TerminateTask(terminate *types.TaskTerminateMsg) {

	log.Infof("Start [terminate task], userType: {%s}, user: {%s}, taskId: {%s}", terminate.GetUserType(), terminate.GetUser(), terminate.GetTaskId())

	// NOTE:
	//
	// The sender of the terminatemsg must be the sender of the task
	//
	// Why load 'local task' instead of 'needexecutetask'?
	//
	// Because the task may still be in the `consensus phase` rather than the `execution phase`,
	// the cached task cannot be found in the 'needexecutetaskcache' at this time.
	task, err := m.resourceMng.GetDB().QueryLocalTask(terminate.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task on `taskManager.TerminateTask()`, taskId: {%s}", terminate.GetTaskId())
		return
	}
	if nil == task {
		log.Warnf("Not found local task on `taskManager.TerminateTask()`, taskId: {%s}", terminate.GetTaskId())
		return
	}
	identity, err := m.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task on `taskManager.TerminateTask()`, taskId: {%s}", terminate.GetTaskId())
		return
	}
	if task.GetTaskSender().GetIdentityId() != identity.GetIdentityId() {
		log.WithError(err).Errorf("The sender of the terminatemsg must be the sender of the task on `taskManager.TerminateTask()`, taskId: {%s}, taskSender: {%s}, localIdentity: {%s}",
			terminate.GetTaskId(), task.GetTaskSender().GetIdentityId(), identity.GetIdentityId())
		return
	}

	if err := m.terminateExecuteTaskBySender(task); nil != err {
		log.Errorf("Failed to call `terminateExecuteTaskBySender()` on `taskManager.TerminateTask()`, taskId: {%s}, err: \n%s", task.GetTaskId(), err)
	}
}

func (m *Manager) terminateExecuteTaskBySender(task *types.Task) error {

	// ## 1、 check whether task status is terminate (with task sender) ?
	terminating, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusTerminateByPartyId(task.GetTaskId(), task.GetTaskSender().GetPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to call HasLocalTaskExecuteStatusTerminateByPartyId() on taskManager.terminateExecuteTaskBySender(), taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetTaskSender().GetPartyId())
		return err
	}
	if terminating {
		log.Warnf("The localTask execute status has `terminate` on taskManager.terminateExecuteTaskBySender(), taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetTaskSender().GetPartyId())
		return fmt.Errorf("task was terminated")
	}

	// ## 2、 check whether task status is running (with task sender) ?
	running, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusRunningByPartyId(task.GetTaskId(), task.GetTaskSender().GetPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to call HasLocalTaskExecuteStatusRunningByPartyId() on taskManager.terminateExecuteTaskBySender(), taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetTaskSender().GetPartyId())
		return err
	}
	if !running {
		log.Warnf("Not found localTask execute status `running` on taskManager.terminateExecuteTaskBySender(), taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetTaskSender().GetPartyId())
		return fmt.Errorf("task is not executed")
	}

	// what if find the needExecuteTask(status: types.TaskConsensusFinished) with sender
	if m.hasNeedExecuteTaskCache(task.GetTaskId(), task.GetTaskSender().GetPartyId()) {

		// 1、 remove needExecuteTask cache with sender #### Need remove after call publishFinishedTaskToDataCenter
		//m.removeNeedExecuteTaskCache(task.GetTaskId(), task.GetTaskSender().GetPartyId())
		// 2、Set the execution status of the task to being terminated`
		if err := m.resourceMng.GetDB().StoreLocalTaskExecuteStatusValTerminateByPartyId(task.GetTaskId(), task.GetTaskSender().GetPartyId()); nil != err {
			log.WithError(err).Errorf("Failed to store needExecute task status to `terminate` with task sender, taskId: {%s}, partyId: {%s}",
				task.GetTaskId(), task.GetTaskSender().GetPartyId())
		}
		// 3、 send a new needExecuteTask(status: types.TaskTerminate) for terminate with sender
		m.sendNeedExecuteTaskByAction(types.NewNeedExecuteTask(
			"",
			commonconstantpb.TaskRole_TaskRole_Sender,
			commonconstantpb.TaskRole_TaskRole_Sender,
			task.GetTaskSender(),
			task.GetTaskSender(),
			task.GetTaskId(),
			types.TaskTerminate,
			&types.PrepareVoteResource{},   // zero value
			&carriertwopcpb.ConfirmTaskPeerInfo{}, // zero value
			fmt.Errorf("task was terminated"),
		))
		// 4、 broadcast terminateMsg to other party of task
		if err := m.broadcastTaskTerminateMsg(task); nil != err {
			log.WithError(err).Errorf("Failed to broadcast terminateMsg to other party of task with task sender, taskId: {%s}, partyId: {%s}",
				task.GetTaskId(), task.GetTaskSender().GetPartyId())
		}
	}

	return nil
}
func (m *Manager) terminateExecuteTaskByPartner(task *types.Task, needExecuteTask *types.NeedExecuteTask) error {
	// ## 1、 check whether task status is terminate (with current party self)?
	terminating, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusTerminateByPartyId(task.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to call HasLocalTaskExecuteStatusTerminateByPartyId() on taskManager.terminateExecuteTaskByPartner(), taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId())
		return err
	}
	if terminating {
		log.Warnf("The localTask execute status has `terminate` on taskManager.terminateExecuteTaskByPartner(), taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId())
		return fmt.Errorf("task was terminated")
	}

	// ## 2、 check whether task status is running (with current party self)?
	running, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusRunningByPartyId(task.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to call HasLocalTaskExecuteStatusRunningByPartyId() on taskManager.terminateExecuteTaskByPartner(), taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId())
		return err
	}
	if !running {
		log.Warnf("Not found localTask execute status `running` on taskManager.terminateExecuteTaskByPartner(), taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId())
		return fmt.Errorf("task is not executed")
	}
	return m.startTerminateWithNeedExecuteTask(needExecuteTask)
}

func (m *Manager) HandleTaskMsgs(msgArr types.TaskMsgArr) error {

	if len(msgArr) == 0 {
		return fmt.Errorf("receive some empty task msgArr")
	}

	nonParsedMsgArr, parsedMsgArr := m.parser.ParseTask(msgArr)
	for _, badMsg := range nonParsedMsgArr {

		events := []*carriertypespb.TaskEvent{m.eventEngine.GenerateEvent(ev.TaskFailed.GetType(),
			badMsg.GetTaskMsg().GetTaskId(), badMsg.GetTaskMsg().GetSenderIdentityId(), badMsg.GetTaskMsg().GetSenderPartyId(), badMsg.GetErrStr())}

		if e := m.publishBadTaskToDataCenter(badMsg.GetTaskMsg().GetTask(), events, "failed to parse taskMsg"); nil != e {
			log.WithError(e).Errorf("Failed to store the err taskMsg on taskManager.HandleTaskMsgs(), taskId: {%s}", badMsg.GetTaskMsg().GetTaskId())
		}
		m.resourceMng.GetDB().RemoveTaskMsg(badMsg.GetTaskMsg().GetTaskId()) // remove from disk if task been non-parsed task
	}

	nonValidatedMsgArr, validatedMsgArr := m.validator.validateTaskMsg(parsedMsgArr)
	for _, badMsg := range nonValidatedMsgArr {

		events := []*carriertypespb.TaskEvent{m.eventEngine.GenerateEvent(ev.TaskFailed.GetType(),
			badMsg.GetTaskMsg().GetTaskId(), badMsg.GetTaskMsg().GetSenderIdentityId(), badMsg.GetTaskMsg().GetSenderPartyId(), badMsg.GetErrStr())}

		if e := m.publishBadTaskToDataCenter(badMsg.GetTaskMsg().GetTask(), events, "failed to validate taskMsg"); nil != e {
			log.WithError(e).Errorf("Failed to store the err taskMsg on taskManager.HandleTaskMsgs(), taskId: {%s}", badMsg.GetTaskMsg().GetTaskId())
		}
		m.resourceMng.GetDB().RemoveTaskMsg(badMsg.GetTaskMsg().GetTaskId()) // remove from disk if task been non-validated task
	}

	log.Infof("Start handle local task msgs on taskManager.HandleTaskMsgs(), task msg arr len: %d", len(validatedMsgArr))

	for _, msg := range validatedMsgArr {

		log.Infof("Start to store local task msg on taskManager.HandleTaskMsgs(), taskId: {%s}", msg.GetTaskId())

		task := msg.GetTask()

		// store metadata used taskId
		if err := m.storeMetadataUsedTaskId(task); nil != err {
			log.WithError(err).Errorf("Failed to store metadata used taskId when received local task on taskManager.HandleTaskMsgs(), taskId: {%s}", task.GetTaskId())
		}

		var storeErrs []string
		if err := m.resourceMng.GetDB().StoreLocalTask(task); nil != err {
			storeErrs = append(storeErrs, err.Error())
		}
		if err := m.storeTaskAllPartnerPartyIds(task); nil != err {
			storeErrs = append(storeErrs, err.Error())
		}

		if len(storeErrs) != 0 {

			log.Errorf("failed to call StoreLocalTask on taskManager with schedule task on taskManager.HandleTaskMsgs(), err: %s",
				"\n["+strings.Join(storeErrs, ",")+"]")

			events := []*carriertypespb.TaskEvent{m.eventEngine.GenerateEvent(ev.TaskDiscarded.Type,
				task.GetTaskId(), task.GetTaskSender().GetIdentityId(), task.GetTaskSender().GetPartyId(), "store local task failed")}
			if err := m.publishBadTaskToDataCenter(task, events, "store local task failed"); nil != err {
				log.WithError(err).Errorf("Failed to sending the task to datacenter on taskManager.HandleTaskMsgs(), taskId: {%s}", task.GetTaskId())
			}

		} else {

			if err := m.scheduler.AddTask(task); nil != err {
				log.WithError(err).Errorf("Failed to add local task into scheduler queue")
			}
			if err := m.tryScheduleTask(); nil != err {
				log.WithError(err).Errorf("Failed to try schedule local task while received local tasks")
			}
		}

		m.resourceMng.GetDB().RemoveTaskMsg(msg.GetTaskId()) // remove from disk if task been handle task
	}
	return nil
}

func (m *Manager) HandleTaskTerminateMsgs(msgArr types.TaskTerminateMsgArr) error {
	for _, terminate := range msgArr {
		go m.TerminateTask(terminate)
	}
	return nil
}

func (m *Manager) SendTaskEvent(event *carriertypespb.TaskEvent) error {
	m.sendTaskEvent(event)
	return nil
}

func (m *Manager) HandleReportResourceUsage(usage *types.TaskResuorceUsage) error {

	identityId, err := m.resourceMng.GetDB().QueryIdentityId()
	if nil != err {
		return fmt.Errorf("query local identityId failed, %s", err)
	}

	needExecuteTask, ok := m.queryNeedExecuteTaskCache(usage.GetTaskId(), usage.GetPartyId())
	if !ok {
		log.Warnf("Not found needExecuteTask on taskManager.HandleReportResourceUsage(), taskId: {%s}, partyId: {%s}",
			usage.GetTaskId(), usage.GetPartyId())
		return fmt.Errorf("can not find `need execute task` cache")
	}

	if identityId != needExecuteTask.GetLocalTaskOrganization().GetIdentityId() {
		log.Errorf("iedntityId od local identity AND identityId of needExecuteTask is not same on taskManager.HandleReportResourceUsage(), taskId: {%s}, partyId: {%s}",
			usage.GetTaskId(), usage.GetPartyId())
		return fmt.Errorf("invalid needExecuteTask local identityId")
	}

	task, err := m.resourceMng.GetDB().QueryLocalTask(usage.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("Failed to call QueryLocalTask() on taskManager.HandleReportResourceUsage(), taskId: {%s}, partyId: {%s}",
			usage.GetTaskId(), usage.GetPartyId())
		return fmt.Errorf("query local task failed, %s", err)
	}

	needUpdate, err := m.handleResourceUsage("when handle report local resourceUsage", needExecuteTask.GetLocalTaskOrganization().GetIdentityId(), usage, task, types.LocalNetworkMsg)
	if nil != err {
		return err
	}

	log.Debugf("Finished handle local resourceUsage on taskManager.HandleReportResourceUsage(), taskId: {%s}, partyId: {%s}, needUpdate: %v, currentOrganization: %s, taskSender: %s, remoteOrganization: %s",
		usage.GetTaskId(), usage.GetPartyId(), needUpdate, identityId, task.GetTaskSender().GetIdentityId(), needExecuteTask.GetRemoteTaskOrganization().GetIdentityId())

	// Update local task AND announce task sender to update task.
	if needUpdate &&
		// If remote organization AND current organization is not same.
		// send usageMsg to remote organization. (NOTE: the remote organization must be task sender)
		needExecuteTask.GetRemoteTaskOrganization().GetIdentityId() != identityId &&
		needExecuteTask.GetRemoteTaskOrganization().GetIdentityId() == task.GetTaskSender().GetIdentityId() {

		usageMsg := &carriernetmsgtaskmngpb.TaskResourceUsageMsg{
			MsgOption: &carriernetmsgcommonpb.MsgOption{
				ProposalId:      common.Hash{}.Bytes(),
				SenderRole:      uint64(needExecuteTask.GetLocalTaskRole()),
				SenderPartyId:   []byte(needExecuteTask.GetLocalTaskOrganization().GetPartyId()),
				ReceiverRole:    uint64(needExecuteTask.GetRemoteTaskRole()),
				ReceiverPartyId: []byte(needExecuteTask.GetRemoteTaskOrganization().GetPartyId()),
				MsgOwner: &carriernetmsgcommonpb.TaskOrganizationIdentityInfo{
					Name:       []byte(needExecuteTask.GetLocalTaskOrganization().GetNodeName()),
					NodeId:     []byte(needExecuteTask.GetLocalTaskOrganization().GetNodeId()),
					IdentityId: []byte(needExecuteTask.GetLocalTaskOrganization().GetIdentityId()),
					PartyId:    []byte(needExecuteTask.GetLocalTaskOrganization().GetPartyId()),
				},
			},
			TaskId: []byte(task.GetTaskId()),
			Usage: &carriernetmsgcommonpb.ResourceUsage{
				TotalMem:       usage.GetTotalMem(),
				UsedMem:        usage.GetUsedMem(),
				TotalProcessor: uint64(usage.GetTotalProcessor()),
				UsedProcessor:  uint64(usage.GetUsedProcessor()),
				TotalBandwidth: usage.GetTotalBandwidth(),
				UsedBandwidth:  usage.GetUsedBandwidth(),
				TotalDisk:      usage.GetTotalDisk(),
				UsedDisk:       usage.GetUsedDisk(),
			},
			CreateAt: timeutils.UnixMsecUint64(),
			Sign:     nil,
		}
		msg := &types.TaskResourceUsageMsg{
			MsgOption: types.FetchMsgOption(usageMsg.GetMsgOption()),
			Usage: types.NewTaskResuorceUsage(
				string(usageMsg.GetTaskId()),
				string(usageMsg.GetMsgOption().GetSenderPartyId()),
				usageMsg.GetUsage().GetTotalMem(),
				usageMsg.GetUsage().GetTotalBandwidth(),
				usageMsg.GetUsage().GetTotalDisk(),
				usageMsg.GetUsage().GetUsedMem(),
				usageMsg.GetUsage().GetUsedBandwidth(),
				usageMsg.GetUsage().GetUsedDisk(),
				uint32(usageMsg.GetUsage().GetTotalProcessor()),
				uint32(usageMsg.GetUsage().GetUsedProcessor()),
			),
			CreateAt: usageMsg.GetCreateAt(),
		}

		// signature the msg and fill sign field of taskResourceUsageMsg
		sign, err := signutil.SignMsg(msg.Hash().Bytes(), m.nodePriKey)
		if nil != err {
			return fmt.Errorf("sign taskResourceUsageMsg, %s", err)
		}
		usageMsg.Sign = sign

		// broadcast `task resource usage msg` to reply remote peer
		// send resource usage quo to remote peer that it will update power supplier resource usage info of task.
		//
		if err := m.p2p.Broadcast(context.TODO(), usageMsg); nil != err {
			log.WithError(err).Errorf("failed to call `SendTaskResourceUsageMsg` on taskManager.HandleReportResourceUsage(), taskId: {%s},  partyId: {%s}, remote pid: {%s}",
				task.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId(), needExecuteTask.GetRemotePID())
		} else {
			log.WithField("traceId", traceutil.GenerateTraceID(usageMsg)).Debugf("Succeed to call broadcast `ResourceUsageMsg` on taskManager.HandleReportResourceUsage(), taskId: {%s},  partyId: {%s}, remote pid: {%s}",
				task.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId(), needExecuteTask.GetRemotePID())
		}
	}

	// The local `resourceUasgeMsg` will not be handled,
	// because the local MSG has already performed the 'Update' operation on the 'localtask'
	// when it has received the reported 'usage', so there is no need to repeat the operation.

	return nil
}

func (m *Manager) storeTaskAllPartnerPartyIds(task *types.Task) error {

	// partyId of powerSuppliers
	partyIds, err := m.policyEngine.FetchPowerPartyIdsFromPowerPolicy(task.GetTaskData().GetPowerPolicyTypes(), task.GetTaskData().GetPowerPolicyOptions())
	if nil != err {
		log.WithError(err).Errorf("not fetch partyIds from task powerPolicy, taskId: {%s}", task.GetTaskId())
		return err
	}

	// partyId of dataSuppliers
	for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
		partyIds = append(partyIds, dataSupplier.GetPartyId())
	}

	// partyId of receivers
	for _, receiver := range task.GetTaskData().GetReceivers() {
		partyIds = append(partyIds, receiver.GetPartyId())
	}
	if len(partyIds) == 0 {
		log.Warnf("have no one partner partyId of task when call storeTaskAllPartnerPartyIds(), taskId: {%s}", task.GetTaskId())
		return nil
	}
	if err := m.resourceMng.GetDB().StoreTaskPartnerPartyIds(task.GetTaskId(), partyIds); nil != err {
		log.WithError(err).Errorf("Failed to store all partner partyIds of task, taskId: {%s}", task.GetTaskId())
		return err
	}
	return nil
}
