package task

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/Metisnetwork/Metis-Carrier/ach/auth"
	"github.com/Metisnetwork/Metis-Carrier/ach/metispay"
	"github.com/Metisnetwork/Metis-Carrier/common"
	"github.com/Metisnetwork/Metis-Carrier/common/bytesutil"
	"github.com/Metisnetwork/Metis-Carrier/common/signutil"
	"github.com/Metisnetwork/Metis-Carrier/common/timeutils"
	"github.com/Metisnetwork/Metis-Carrier/common/traceutil"
	"github.com/Metisnetwork/Metis-Carrier/consensus"
	ev "github.com/Metisnetwork/Metis-Carrier/core/evengine"
	"github.com/Metisnetwork/Metis-Carrier/core/rawdb"
	"github.com/Metisnetwork/Metis-Carrier/core/resource"
	"github.com/Metisnetwork/Metis-Carrier/core/schedule"
	msgcommonpb "github.com/Metisnetwork/Metis-Carrier/lib/netmsg/common"
	twopcpb "github.com/Metisnetwork/Metis-Carrier/lib/netmsg/consensus/twopc"
	taskmngpb "github.com/Metisnetwork/Metis-Carrier/lib/netmsg/taskmng"
	libtypes "github.com/Metisnetwork/Metis-Carrier/lib/types"
	"github.com/Metisnetwork/Metis-Carrier/p2p"
	"github.com/Metisnetwork/Metis-Carrier/params"
	"github.com/Metisnetwork/Metis-Carrier/policy"
	"github.com/Metisnetwork/Metis-Carrier/types"
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
	metisPayMng     *metispay.MetisPayManager
	parser          *TaskParser
	validator       *TaskValidator
	eventCh         chan *libtypes.TaskEvent

	// send the validated taskMsgs to scheduler
	localTasksCh             chan types.TaskDataArray
	taskConsResultCh         chan *types.TaskConsResult
	needReplayScheduleTaskCh chan *types.NeedReplayScheduleTask
	needExecuteTaskCh        chan *types.NeedExecuteTask
	runningTaskCache         map[string]map[string]*types.NeedExecuteTask //  taskId -> {partyId -> task}
	syncExecuteTaskMonitors  *types.SyncExecuteTaskMonitorQueue
	quit                     chan struct{}
	runningTaskCacheLock     sync.RWMutex
	// add by v 0.4.0
	config   *params.TaskManagerConfig
	msgCache *TaskmngMsgCache
}

func NewTaskManager(
	nodePriKey *ecdsa.PrivateKey,
	p2p p2p.P2P,
	scheduler schedule.Scheduler,
	consensusEngine consensus.Engine,
	eventEngine *ev.EventEngine,
	resourceMng *resource.Manager,
	authMng *auth.AuthorityManager,
	metisPayMng *metispay.MetisPayManager,
	needReplayScheduleTaskCh chan *types.NeedReplayScheduleTask,
	needExecuteTaskCh chan *types.NeedExecuteTask,
	taskConsResultCh chan *types.TaskConsResult,
	config *params.TaskManagerConfig,
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
		metisPayMng:              metisPayMng,
		parser:                   newTaskParser(resourceMng),
		validator:                newTaskValidator(resourceMng, authMng),
		eventCh:                  make(chan *libtypes.TaskEvent, 800),
		needReplayScheduleTaskCh: needReplayScheduleTaskCh,
		needExecuteTaskCh:        needExecuteTaskCh,
		taskConsResultCh:         taskConsResultCh,
		config:                   config,
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

			var res libtypes.NeedExecuteTask

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
			go func(event *libtypes.TaskEvent) {
				if err := m.handleTaskEventWithCurrentOranization(event); nil != err {
					log.WithError(err).Errorf("Failed to call handleTaskEventWithCurrentOranization() on taskManager.loop(), taskId: {%s}, event: %s", event.GetTaskId(), event.String())
				}
			}(event)


		// To schedule local task interval
		case <-taskTicker.C:

			if err := m.tryScheduleTask(); nil != err {
				log.WithError(err).Errorf("Failed to try schedule local task when taskTicker")
			}

		case res := <-m.taskConsResultCh: // received the task consensus result by task sender

			if nil == res {
				continue
			}

			go func(result *types.TaskConsResult) {

				log.Debugf("Received `NEED-CONSENSUS` task result from 2pc consensus engine, taskId: {%s}, result: {%s}", result.GetTaskId(), result.String())

				task, err := m.resourceMng.GetDB().QueryLocalTask(result.GetTaskId())
				if nil != err {
					log.WithError(err).Errorf("Failed to query local task when received `NEED-CONSENSUS` task result, taskId: {%s}, result: {%s}", result.GetTaskId(), result.String())
					return
				}

				// received status must be `TaskConsensusFinished` & `TaskConsensusInterrupt` & `TaskTerminate` from consensus engine
				// never be `TaskNeedExecute|TaskScheduleFailed`
				//
				// Consensus failed, task needs to be suspended and rescheduled
				switch result.GetStatus() {
				case types.TaskTerminate, types.TaskConsensusFinished:

					if result.GetStatus() == types.TaskConsensusFinished {
						// store task consensus result (failed or succeed) event with sender party
						m.resourceMng.GetDB().StoreTaskEvent(&libtypes.TaskEvent{
							Type:       ev.TaskSucceedConsensus.GetType(),
							TaskId:     task.GetTaskId(),
							IdentityId: task.GetTaskSender().GetIdentityId(),
							PartyId:    task.GetTaskSender().GetPartyId(),
							Content:    "succeed consensus",
							CreateAt:   timeutils.UnixMsecUint64(),
						})
					}
					// remove task from scheduler.queue|starvequeue after task consensus succeed
					// Don't send needexecuteTask, because that will be handle in `2pc engine.driveTask()`
					if err := m.scheduler.RemoveTask(result.GetTaskId()); nil != err {
						log.WithError(err).Errorf("Failed to remove local task from queue/starve queue %s when received `NEED-CONSENSUS` task result, taskId: {%s}",
							result.GetStatus().String(), task.GetTaskId())
					}
					m.sendNeedExecuteTaskByAction(types.NewNeedExecuteTask(
						"",
						libtypes.TaskRole_TaskRole_Sender,
						libtypes.TaskRole_TaskRole_Sender,
						task.GetTaskSender(),
						task.GetTaskSender(),
						task.GetTaskId(),
						result.GetStatus(),
						&types.PrepareVoteResource{},   // zero value
						&twopcpb.ConfirmTaskPeerInfo{}, // zero value
						result.GetErr(),
					))
				case types.TaskConsensusInterrupt:

					// clean old powerSuppliers and update local task
					task.RemovePowerSuppliers()
					// restore task by power
					if err := m.resourceMng.GetDB().StoreLocalTask(task); nil != err {
						log.WithError(err).Errorf("Failed to update local task whit clean powers after consensus interrupted when received `NEED-CONSENSUS` task result, taskId: {%s}", task.GetTaskId())
					}

					// re push task into queue ,if anything else
					if err := m.scheduler.RepushTask(task); err == schedule.ErrRescheduleLargeThreshold {
						log.WithError(err).Errorf("Failed to repush local task into queue/starve queue %s when received `NEED-CONSENSUS` task result, taskId: {%s}",
							result.GetStatus().String(), task.GetTaskId())

						m.scheduler.RemoveTask(task.GetTaskId())
						m.sendNeedExecuteTaskByAction(types.NewNeedExecuteTask(
							"",
							libtypes.TaskRole_TaskRole_Sender,
							libtypes.TaskRole_TaskRole_Sender,
							task.GetTaskSender(),
							task.GetTaskSender(),
							task.GetTaskId(),
							types.TaskScheduleFailed,
							&types.PrepareVoteResource{},   // zero value
							&twopcpb.ConfirmTaskPeerInfo{}, // zero value
							fmt.Errorf("consensus interrupted: "+result.GetErr().Error()+" and "+schedule.ErrRescheduleLargeThreshold.Error()),
						))
					} else {
						log.Debugf("Succeed to repush local task into queue/starve queue %s when received `NEED-CONSENSUS` task result, taskId: {%s}",
							result.GetStatus().String(), task.GetTaskId())
					}
				}
			}(res)

		// handle the task of need replay scheduling while received from remote peer on consensus epoch
		case task := <-m.needReplayScheduleTaskCh:

			if nil == task {
				continue
			}

			go func(task *types.NeedReplayScheduleTask) {

				// Do duplication check ...
				has, err := m.resourceMng.GetDB().HasLocalTask(task.GetTask().GetTaskId())
				if nil != err {
					log.WithError(err).Errorf("Failed to query local task when received remote task, taskId: {%s}", task.GetTask().GetTaskId())
					return
				}

				// There is no need to store local tasks repeatedly
				if !has {

					log.Infof("Start to store local task on taskManager.loop() when received needReplayScheduleTask, taskId: {%s}", task.GetTask().GetTaskId())

					// store metadata used taskId
					if err := m.storeMetadataUsedTaskId(task.GetTask()); nil != err {
						log.WithError(err).Errorf("Failed to store metadata used taskId when received remote task, taskId: {%s}", task.GetTask().GetTaskId())
					}
					if err := m.resourceMng.GetDB().StoreLocalTask(task.GetTask()); nil != err {
						log.WithError(err).Errorf("Failed to call StoreLocalTask when replay schedule remote task, taskId: {%s}", task.GetTask().GetTaskId())
						task.SendFailedResult(task.GetTask().GetTaskId(), err)
						return
					}
					log.Infof("Finished to store local task on taskManager.loop() when received needReplayScheduleTask, taskId: {%s}", task.GetTask().GetTaskId())

				}

				// Start replay schedule remote task ...
				result := m.scheduler.ReplaySchedule(task.GetLocalPartyId(), task.GetLocalTaskRole(), task)
				task.SendResult(result)
			}(task)

		// handle task of need to executing, and send it to fighter of myself organization or send the task result msg to remote peer
		case task := <-m.needExecuteTaskCh:

			if nil == task {
				continue
			}

			go func(task *types.NeedExecuteTask) {

				localTask, err := m.resourceMng.GetDB().QueryLocalTask(task.GetTaskId())
				if nil != err {
					log.WithError(err).Errorf("Failed to query local task info on taskManager.loop() when received needExecuteTask, taskId: {%s}, partyId: {%s}, status: {%s}",
						task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), task.GetConsStatus().String())
					//continue
					return
				}

				// init consumeSpec of NeedExecuteTask first (by v0.4.0)
				m.initConsumeSpecByConsumeOption(task)

				switch task.GetConsStatus() {
				// sender and partner to handle needExecuteTask when consensus succeed.
				// sender need to store some cache, partner need to execute task.
				case types.TaskNeedExecute, types.TaskConsensusFinished:
					// to execute the task
					m.handleNeedExecuteTask(task, localTask)

				// sender and partner to clear local task things after received status: `scheduleFailed` and `interrupt` and `terminate`.
				// sender need to publish local task and event to datacenter,
				// partner need to send task's event to remote task's sender.
				//
				// NOTE:
				//	sender never received status `interrupt`, and partners never received status `scheduleFailed`.
				default:

					// store a bad event into local db before handle bad task.
					var eventTyp string
					if task.GetConsStatus() == types.TaskConsensusInterrupt {
						eventTyp = ev.TaskFailedConsensus.GetType()
					} else {
						eventTyp = ev.TaskFailed.GetType() // then the task status was `scheduleFailed` and `terminate`.
					}
					m.resourceMng.GetDB().StoreTaskEvent(m.eventEngine.GenerateEvent(eventTyp, task.GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(),
						task.GetLocalTaskOrganization().GetPartyId(), task.GetErr().Error()))

					switch task.GetLocalTaskRole() {
					case libtypes.TaskRole_TaskRole_Sender:
						m.publishFinishedTaskToDataCenter(task, localTask, true)
					default:
						m.sendTaskResultMsgToTaskSender(task, localTask)
					}
				}
			}(task)

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

	log.Infof("Start terminate task, userType: {%s}, user: {%s}, taskId: {%s}", terminate.GetUserType(), terminate.GetUser(), terminate.GetTaskId())

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

	// ## 1、 check wether task status is `terminate`
	terminating, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusTerminateByPartyId(task.GetTaskId(), task.GetTaskSender().GetPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task execute `termining` with task sender on `taskManager.TerminateTask()`, taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetTaskSender().GetPartyId())
		return
	}
	if terminating {
		log.Warnf("Warning query local task execute status has `termining` with task sender on `taskManager.TerminateTask()`, taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetTaskSender().GetPartyId())
		return
	}

	// ## 2、 check wether task status is `running`
	running, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusRunningByPartyId(task.GetTaskId(), task.GetTaskSender().GetPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task execute status has `running` with task sender on `taskManager.TerminateTask()`, taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetTaskSender().GetPartyId())
		return
	}
	// (While task is consensus or running, can terminate.)
	if running {
		if err := m.onTerminateExecuteTask(task.GetTaskId(), task.GetTaskSender().GetPartyId(), task); nil != err {
			log.Errorf("Failed to call `onTerminateExecuteTask()` on `taskManager.TerminateTask()`, taskId: {%s}, err: \n%s", task.GetTaskId(), err)
		}
		return
	}

	// ## 3、 check wether task status is `consensus`

	// The task sender only makes consensus, so interrupt consensus while need terminate task with task sender
	// (While task is consensus or executing, can terminate.)
	consensusing, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusConsensusByPartyId(task.GetTaskId(), task.GetTaskSender().GetPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task execute status has `cons` with task sender on `taskManager.TerminateTask()`, taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetTaskSender().GetPartyId())
		return
	}
	if consensusing {
		if err = m.consensusEngine.OnConsensusMsg(
			"", types.NewInterruptMsgWrap(task.GetTaskId(),
				types.MakeMsgOption(common.Hash{},
					libtypes.TaskRole_TaskRole_Sender,
					libtypes.TaskRole_TaskRole_Sender,
					task.GetTaskSender().GetPartyId(),
					task.GetTaskSender().GetPartyId(),
					task.GetTaskSender()))); nil != err {
			log.WithError(err).Errorf("Failed to call `OnConsensusMsg()` on `taskManager.TerminateTask()`, taskId: {%s}, partyId: {%s}",
				task.GetTaskId(), task.GetTaskSender().GetPartyId())
		}
	}
}

func (m *Manager) onTerminateExecuteTask(taskId, partyId string, task *types.Task) error {

	if taskId != task.GetTaskId() {
		log.Warnf("taskId has not match, %s and %s is not same one", taskId, task.GetTaskId())
		return fmt.Errorf("taskId has not match, %s and %s is not same one", taskId, task.GetTaskId())
	}

	// what if find the needExecuteTask(status: types.TaskConsensusFinished) with sender
	if m.hasNeedExecuteTaskCache(task.GetTaskId(), task.GetTaskSender().GetPartyId()) {

		// 1、 store task terminate (failed or succeed) event with current party
		m.resourceMng.GetDB().StoreTaskEvent(&libtypes.TaskEvent{
			Type:       ev.TaskTerminated.Type,
			TaskId:     task.GetTaskId(),
			IdentityId: task.GetTaskSender().GetIdentityId(),
			PartyId:    task.GetTaskSender().GetPartyId(),
			Content:    "task was terminated.",
			CreateAt:   timeutils.UnixMsecUint64(),
		})

		// 2、 remove needExecuteTask cache with sender
		m.removeNeedExecuteTaskCache(task.GetTaskId(), task.GetTaskSender().GetPartyId())
		// 3、Set the execution status of the task to being terminated`
		if err := m.resourceMng.GetDB().StoreLocalTaskExecuteStatusValTerminateByPartyId(task.GetTaskId(), task.GetTaskSender().GetPartyId()); nil != err {
			log.WithError(err).Errorf("Failed to store needExecute task status to `terminate` with task sender, taskId: {%s}, partyId: {%s}",
				task.GetTaskId(), task.GetTaskSender().GetPartyId())
		}
		// 4、 send a new needExecuteTask(status: types.TaskTerminate) for terminate with sender
		m.sendNeedExecuteTaskByAction(types.NewNeedExecuteTask(
			"",
			libtypes.TaskRole_TaskRole_Sender,
			libtypes.TaskRole_TaskRole_Sender,
			task.GetTaskSender(),
			task.GetTaskSender(),
			task.GetTaskId(),
			types.TaskTerminate,
			&types.PrepareVoteResource{},   // zero value
			&twopcpb.ConfirmTaskPeerInfo{}, // zero value
			fmt.Errorf("task was terminated"),
		))

		if err := m.sendTaskTerminateMsg(task); nil != err {
			log.WithError(err).Errorf("Failed to store needExecute task status to `terminate` with task sender, taskId: {%s}, partyId: {%s}",
				task.GetTaskId(), task.GetTaskSender().GetPartyId())
		}
	} else {

		// When the partner and sender of the current task do not belong to the same organization,
		// we terminate the task of our own party.
		if needExecuteTask, has := m.queryNeedExecuteTaskCache(taskId, partyId); has {
			return m.startTerminateWithNeedExecuteTask(needExecuteTask)
		}
	}

	return nil
}

func (m *Manager) HandleTaskMsgs(msgArr types.TaskMsgArr) error {

	if len(msgArr) == 0 {
		return fmt.Errorf("receive some empty task msgArr")
	}

	nonParsedMsgArr, parsedMsgArr := m.parser.ParseTask(msgArr)
	for _, badMsg := range nonParsedMsgArr {

		events := []*libtypes.TaskEvent{m.eventEngine.GenerateEvent(ev.TaskFailed.GetType(),
			badMsg.GetTaskMsg().GetTaskId(), badMsg.GetTaskMsg().GetSenderIdentityId(), badMsg.GetTaskMsg().GetSenderPartyId(), badMsg.GetErrStr())}

		if e := m.publishBadTaskToDataCenter(badMsg.GetTaskMsg().GetTask(), events, "failed to parse taskMsg"); nil != e {
			log.WithError(e).Errorf("Failed to store the err taskMsg on taskManager.HandleTaskMsgs(), taskId: {%s}", badMsg.GetTaskMsg().GetTaskId())
		}
		m.resourceMng.GetDB().RemoveTaskMsg(badMsg.GetTaskMsg().GetTaskId()) // remove from disk if task been non-parsed task
	}

	nonValidatedMsgArr, validatedMsgArr := m.validator.validateTaskMsg(parsedMsgArr)
	for _, badMsg := range nonValidatedMsgArr {

		events := []*libtypes.TaskEvent{m.eventEngine.GenerateEvent(ev.TaskFailed.GetType(),
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

			events := []*libtypes.TaskEvent{m.eventEngine.GenerateEvent(ev.TaskDiscarded.Type,
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

func (m *Manager) SendTaskEvent(event *libtypes.TaskEvent) error {
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
		return fmt.Errorf("Can not find `need execute task` cache")
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

		usageMsg := &taskmngpb.TaskResourceUsageMsg{
			MsgOption: &msgcommonpb.MsgOption{
				ProposalId:      common.Hash{}.Bytes(),
				SenderRole:      uint64(needExecuteTask.GetLocalTaskRole()),
				SenderPartyId:   []byte(needExecuteTask.GetLocalTaskOrganization().GetPartyId()),
				ReceiverRole:    uint64(needExecuteTask.GetRemoteTaskRole()),
				ReceiverPartyId: []byte(needExecuteTask.GetRemoteTaskOrganization().GetPartyId()),
				MsgOwner: &msgcommonpb.TaskOrganizationIdentityInfo{
					Name:       []byte(needExecuteTask.GetLocalTaskOrganization().GetNodeName()),
					NodeId:     []byte(needExecuteTask.GetLocalTaskOrganization().GetNodeId()),
					IdentityId: []byte(needExecuteTask.GetLocalTaskOrganization().GetIdentityId()),
					PartyId:    []byte(needExecuteTask.GetLocalTaskOrganization().GetPartyId()),
				},
			},
			TaskId: []byte(task.GetTaskId()),
			Usage: &msgcommonpb.ResourceUsage{
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
		msg.Sign = sign

		// broadcast `task resource usage msg` to reply remote peer
		// send resource usage quo to remote peer that it will update power supplier resource usage info of task.
		//
		if err := m.p2p.Broadcast(context.TODO(), usageMsg); nil != err {
			log.WithError(err).Errorf("failed to call `SendTaskResourceUsageMsg` on taskManager.HandleReportResourceUsage(), taskId: {%s},  partyId: {%s}, remote pid: {%s}",
				task.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId(), needExecuteTask.GetRemotePID())
		} else {
			log.WithField("traceId", traceutil.GenerateTraceID(usageMsg)).Debugf("Succeed to call `SendTaskResourceUsageMsg` on taskManager.HandleReportResourceUsage(), taskId: {%s},  partyId: {%s}, remote pid: {%s}",
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
	partyIds, err := policy.FetchPowerPartyIdsFromPowerPolicy(task.GetTaskData().GetPowerPolicyTypes(), task.GetTaskData().GetPowerPolicyOptions())
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
