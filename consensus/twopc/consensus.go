package twopc

import (
	"bytes"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/blacklist"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/common/bytesutil"
	"github.com/datumtechs/datum-network-carrier/common/rlputil"
	"github.com/datumtechs/datum-network-carrier/common/signutil"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	"github.com/datumtechs/datum-network-carrier/common/traceutil"
	ctypes "github.com/datumtechs/datum-network-carrier/consensus/twopc/types"
	ev "github.com/datumtechs/datum-network-carrier/core/evengine"
	"github.com/datumtechs/datum-network-carrier/core/resource"
	"github.com/datumtechs/datum-network-carrier/p2p"
	carriernetmsgcommonpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/common"
	carriertwopcpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/consensus/twopc"
	carrierrpcdebugpbv1 "github.com/datumtechs/datum-network-carrier/pb/carrier/rpc/debug/v1"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"strings"
	"sync"
	"time"
)

const (
	//defaultCleanExpireProposalInterval  = 30 * time.Millisecond
	defaultRefreshProposalStateInternal = 300 * time.Millisecond
)
const thresholdCount = 10

type OrganizationTaskInfo struct {
	taskId     string
	nodeId     string
	proposalId string
	partyId    string
}

type Twopc struct {
	config                   *Config
	p2p                      p2p.P2P
	state                    *state
	resourceMng              *resource.Manager
	needReplayScheduleTaskCh chan *types.NeedReplayScheduleTask // send remote task to `Scheduler` to replay
	needExecuteTaskCh        chan *types.NeedExecuteTask        // send has was consensus remote tasks to taskManager
	asyncCallCh              chan func()
	quit                     chan struct{}
	taskConsResultCh         chan *types.TaskConsResult
	wal                      *walDB
	Errs                     []error
	orgBlacklistLock         sync.RWMutex
	orgBlacklistCache        map[string][]*OrganizationTaskInfo
	identityBlackListCache   *blacklist.IdentityBackListCache
}

func New(
	conf *Config,
	resourceMng *resource.Manager,
	p2p p2p.P2P,
	needReplayScheduleTaskCh chan *types.NeedReplayScheduleTask,
	needExecuteTaskCh chan *types.NeedExecuteTask,
	taskConsResultCh chan *types.TaskConsResult,
	identityBlackListCache *blacklist.IdentityBackListCache,
) (*Twopc, error) {
	newWalDB := newWal(conf)
	state, err := newState(newWalDB)
	if nil != err {
		return nil, err
	}
	engine := &Twopc{
		config:                   conf,
		p2p:                      p2p,
		state:                    state,
		resourceMng:              resourceMng,
		needReplayScheduleTaskCh: needReplayScheduleTaskCh,
		needExecuteTaskCh:        needExecuteTaskCh,
		asyncCallCh:              make(chan func(), conf.PeerMsgQueueSize),
		quit:                     make(chan struct{}),
		taskConsResultCh:         taskConsResultCh,
		wal:                      newWalDB,
		Errs:                     make([]error, 0),
		identityBlackListCache:   identityBlackListCache,
	}
	identityBlackListCache.SetEngineAndWal(engine, newWalDB)
	return engine, nil
}

func (t *Twopc) Start() error {
	t.recoverCache()
	go t.loop()
	log.Info("Started 2pc consensus engine ...")
	return nil
}
func (t *Twopc) Stop() error {
	close(t.quit)
	return nil
}
func (t *Twopc) loop() {

	proposalStateMonitorTimer := t.proposalStateMonitorTimer()
	for {
		select {

		// Force serial execution of calls initiated by each goroutine,
		// simplifying the error-proneness of concurrent logic
		case fn := <-t.asyncCallCh:

			fn()

		case <-proposalStateMonitorTimer.C:

			future := t.checkProposalStateMonitors(timeutils.UnixMsec(), true)
			now := timeutils.UnixMsec()
			if future > now {
				proposalStateMonitorTimer.Reset(time.Duration(future-now) * time.Millisecond)
			} else if future < now {
				proposalStateMonitorTimer.Reset(time.Duration(now) * time.Millisecond)
			}
			// when future value is 0, we do nothing

		case <-t.quit:
			log.Info("Stopped 2pc consensus engine ...")
			proposalStateMonitorTimer.Stop()
			return
		}
	}
}

func (t *Twopc) OnConsensusMsg(pid peer.ID, msg types.ConsensusMsg) error {

	switch msg := msg.(type) {
	case *types.PrepareMsgWrap:
		return t.onPrepareMsg(pid, msg, types.RemoteNetworkMsg)
	case *types.PrepareVoteWrap:
		return t.onPrepareVote(pid, msg, types.RemoteNetworkMsg)
	case *types.ConfirmMsgWrap:
		return t.onConfirmMsg(pid, msg, types.RemoteNetworkMsg)
	case *types.ConfirmVoteWrap:
		return t.onConfirmVote(pid, msg, types.RemoteNetworkMsg)
	case *types.CommitMsgWrap:
		return t.onCommitMsg(pid, msg, types.RemoteNetworkMsg)
	case *types.TerminateConsensusMsgWrap: // Must be  local msg
		return t.onTerminateTaskConsensus(pid, msg)
	default:
		return fmt.Errorf("unknown the 2pc msg type")

	}
}

func (t *Twopc) OnError() error {
	if len(t.Errs) == 0 {
		return nil
	}
	errStrs := make([]string, len(t.Errs))
	for _, err := range t.Errs {
		errStrs = append(errStrs, err.Error())
	}
	// reset Errs
	t.Errs = make([]error, 0)
	return fmt.Errorf("%s", strings.Join(errStrs, "\n"))
}

func (t *Twopc) OnPrepare(task *types.NeedConsensusTask) error { return nil }

func (t *Twopc) OnHandle(nonConsTask *types.NeedConsensusTask) error {

	task := nonConsTask.GetTask()
	if t.state.HasProposalTaskWithTaskIdAndPartyId(task.GetTaskId(), task.GetTaskSender().GetPartyId()) {
		log.Errorf("Failed to check org proposalTask whether have been not exist on OnHandle, but it's alreay exist, taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetTaskSender().GetPartyId())
		t.stopTaskConsensus(ctypes.ErrPrososalTaskIsProcessed.Error(), common.Hash{}, task.GetTaskId(),
			commonconstantpb.TaskRole_TaskRole_Sender, commonconstantpb.TaskRole_TaskRole_Sender, task.GetTaskSender(), task.GetTaskSender(),
			types.TaskConsensusInterrupt)
		return ctypes.ErrPrososalTaskIsProcessed
	}

	// Store task execute status `cons` before consensus when send task prepareMsg to remote peers
	if err := t.resourceMng.GetDB().StoreLocalTaskExecuteStatusValConsByPartyId(task.GetTaskId(), task.GetTaskSender().GetPartyId()); nil != err {
		log.WithError(err).Errorf("Failed to store local task about `cons` status on OnHandle,  taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetTaskSender().GetPartyId())
		t.stopTaskConsensus("store task executeStatus about `cons` failed", common.Hash{}, task.GetTaskId(),
			commonconstantpb.TaskRole_TaskRole_Sender, commonconstantpb.TaskRole_TaskRole_Sender, task.GetTaskSender(), task.GetTaskSender(),
			types.TaskConsensusInterrupt)
		return err
	}

	createAt := timeutils.UnixMsecUint64()

	var buf bytes.Buffer
	buf.Write(t.config.Option.NodeID.Bytes())
	buf.Write([]byte(task.GetTaskId()))
	buf.Write([]byte(task.GetTaskData().GetTaskName()))
	buf.Write(bytesutil.Uint64ToBytes(task.GetTaskData().GetCreateAt()))
	buf.Write(bytesutil.Uint64ToBytes(createAt))
	proposalId := rlputil.RlpHash(buf.Bytes())

	log.Infof("Generate proposal, proposalId: {%s}, taskId: {%s}, partyId: {%s}", proposalId.String(), task.GetTaskId(), task.GetTaskSender().GetPartyId())

	// Store some local cache
	t.storeOrgProposalState(
		ctypes.NewOrgProposalState(proposalId, task.GetTaskId(),
			commonconstantpb.TaskRole_TaskRole_Sender,
			task.GetTaskSender(), task.GetTaskSender(),
			createAt),
	)

	proposalTask := ctypes.NewProposalTask(proposalId, task.GetTaskId(), createAt)
	t.state.StoreProposalTaskWithPartyId(task.GetTaskSender().GetPartyId(), proposalTask)
	t.wal.StoreProposalTask(task.GetTaskSender().GetPartyId(), proposalTask)
	// Start handle task ...
	go func() {

		if err := t.sendPrepareMsg(proposalId, nonConsTask, createAt); nil != err {
			log.Errorf("Failed to call `sendPrepareMsg`, consensus epoch finished, proposalId: {%s}, taskId: {%s}, partyId: {%s}, err: \n%s",
				proposalId.String(), task.GetTaskId(), task.GetTaskSender().GetPartyId(), err)
			// Send consensus result to Scheduler
			t.stopTaskConsensus("send prepareMsg failed", proposalId, task.GetTaskId(),
				commonconstantpb.TaskRole_TaskRole_Sender, commonconstantpb.TaskRole_TaskRole_Sender, task.GetTaskSender(), task.GetTaskSender(), types.TaskConsensusInterrupt)
			// clean some invalid data
			t.removeOrgProposalStateAndTask(proposalId, task.GetTaskSender().GetPartyId())
		}
	}()
	return nil
}

// Handle the prepareMsg from the task pulisher peer (on Subscriber)
func (t *Twopc) onPrepareMsg(pid peer.ID, prepareMsg *types.PrepareMsgWrap, nmls types.NetworkMsgLocationSymbol) error {

	if err := t.state.ContainsOrAddMsg(prepareMsg.GetData()); nil != err {
		return err
	}

	msg, err := fetchPrepareMsg(prepareMsg)
	if nil != err {
		return err
	}

	// the prepareMsg is future msg.
	now := timeutils.UnixMsecUint64()
	jitterValue := uint64(100)
	if now+jitterValue < msg.GetCreateAt() { // maybe it be allowed to overflow 100ms for timewindows
		log.Errorf("received the prepareMsg is future msg when received prepareMsg, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, now: {%d}, msgCreateAt: {%d}",
			msg.GetMsgOption().GetProposalId().String(), msg.GetTask().GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), now, msg.GetCreateAt())
		return fmt.Errorf("%s when received prepareMsg", ctypes.ErrProposalIllegal)
	}
	// the prepareMsg is too late.
	if (now + jitterValue - msg.GetCreateAt()) >= uint64(ctypes.PrepareMsgVotingDuration.Milliseconds()) {
		log.Errorf("received the prepareMsg is too late when received prepareMsg, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, now: {%d}, msgCreateAt: {%d}, duration: {%d}, valid duration: {%d}",
			msg.GetMsgOption().GetProposalId().String(), msg.GetTask().GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(),
			now, msg.GetCreateAt(), now-msg.GetCreateAt(), ctypes.PrepareMsgVotingDuration.Milliseconds())
		return fmt.Errorf("%s when received prepareMsg", ctypes.ErrProposalIllegal)
	}

	// Verify the signature
	_, err = signutil.VerifyMsgSign(msg.GetMsgOption().GetOwner().GetNodeId(), msg.Hash().Bytes(), msg.GetSign())
	if err != nil {
		return fmt.Errorf("verify prepareMsg sign %s", err)
	}

	if err := t.validateTaskOfPrepareMsg(pid, msg); nil != err {
		return err
	}

	errCh := make(chan error, 0)

	t.asyncCallCh <- func() {

		if t.state.HasOrgProposalWithProposalId(msg.GetMsgOption().GetProposalId()) {
			log.Errorf("Failed to check org proposalState whether have been not exist when received prepareMsg, but it's alreay exist, proposalId: {%s}, taskId: {%s}",
				msg.GetMsgOption().GetProposalId().String(), msg.GetTask().GetTaskId())
			errCh <- fmt.Errorf("%s when received prepareMsg", ctypes.ErrProposalAlreadyProcessed)
			return
		}

		identity, err := t.resourceMng.GetDB().QueryIdentity()
		if nil != err {
			log.WithError(err).Errorf("Failed to call `QueryIdentity()` when received prepareMsg, proposalId: {%s}, taskId: {%s}",
				msg.GetMsgOption().GetProposalId().String(), msg.GetTask().GetTaskId())
			errCh <- fmt.Errorf("query local identity failed when received prepareMsg, %s", err)
			return
		}

		sender := fetchOrgByPartyRole(msg.GetMsgOption().GetSenderPartyId(), msg.GetMsgOption().GetSenderRole(), msg.GetTask())
		if nil == sender {
			log.Errorf("Failed to check msg.MsgOption sender and receiver when received prepareMsg, some one is empty, proposalId: {%s}, taskId: {%s}",
				msg.GetMsgOption().GetProposalId().String(), msg.GetTask().GetTaskId())
			errCh <- fmt.Errorf("%s when received prepareMsg", ctypes.ErrConsensusMsgInvalid)
			return
		}

		// Check whether the sender of the message is the same organization as the sender of the task.
		// If not, this message is illegal.
		if msg.GetTask().GetTaskSender().GetIdentityId() != sender.GetIdentityId() ||
			msg.GetTask().GetTaskSender().GetPartyId() != sender.GetPartyId() {
			log.Warnf("Warning the sender of the message is not the same organization as the sender of the task when received prepareMsg, proposalId: {%s}, taskId: {%s}, msg sender: %s, task sender: %s",
				msg.GetMsgOption().GetProposalId().String(), msg.GetTask().GetTaskId(), sender.String(), msg.GetTask().GetTaskSender().String())
			errCh <- fmt.Errorf("%s when received prepareMsg", ctypes.ErrConsensusMsgInvalid)
			return
		}

		log.WithField("traceId", traceutil.GenerateTraceID(prepareMsg.GetData())).Debugf("Received prepareMsg, consensusSymbol: {%s}, remote pid: {%s}, prepareMsg: %s", nmls.String(), pid, msg.String())

		votingFn := func(party *carriertypespb.TaskOrganization, role commonconstantpb.TaskRole) error {

			org := &carriertypespb.TaskOrganization{
				PartyId:    party.GetPartyId(),
				NodeName:   identity.GetNodeName(),
				NodeId:     identity.GetNodeId(),
				IdentityId: identity.GetIdentityId(),
			}

			// If you have already voted, we will no longer vote.
			// Need to receive a forwarded objection consensus message.
			if t.state.HasPrepareVoting(msg.GetMsgOption().GetProposalId(), org) {
				log.Errorf("Failed to check remote peer prepare vote wether exist when received prepareMsg, it's exist alreay, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					msg.GetMsgOption().GetProposalId().String(), msg.GetTask().GetTaskId(), party.GetPartyId())
				return fmt.Errorf("%s when received prepareMsg", ctypes.ErrPrepareVotehadVoted)
			}

			// Store task execute status `cons` before consensus when received a remote task prepareMsg
			if err := t.resourceMng.GetDB().StoreLocalTaskExecuteStatusValConsByPartyId(msg.GetTask().GetTaskId(), party.GetPartyId()); nil != err {
				log.WithError(err).Errorf("Failed to store local task about `cons` status when received prepareMsg, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					msg.GetMsgOption().GetProposalId().String(), msg.GetTask().GetTaskId(), party.GetPartyId())
				return fmt.Errorf("store task execute status failed when received prepareMsg, %s", err)
			}

			log.Infof("Store proposal from task sender, proposalId: {%s}, taskId: {%s}, partyId: {%s}", msg.GetMsgOption().String(), msg.GetTask().GetTaskId(), party.GetPartyId())

			// Store some local cache
			t.storeOrgProposalState(
				ctypes.NewOrgProposalState(msg.GetMsgOption().GetProposalId(),
					msg.GetTask().GetTaskId(),
					msg.GetMsgOption().GetReceiverRole(), msg.GetTask().GetTaskSender(), party,
					msg.GetCreateAt()),
			)

			proposalTask := ctypes.NewProposalTask(msg.GetMsgOption().GetProposalId(), msg.GetTask().GetTaskId(), msg.GetCreateAt())
			t.state.StoreProposalTaskWithPartyId(msg.GetMsgOption().GetReceiverPartyId(), proposalTask)
			t.wal.StoreProposalTask(msg.GetMsgOption().GetReceiverPartyId(), proposalTask)

			// Send task to Scheduler to replay sched.
			needReplayScheduleTask := types.NewNeedReplayScheduleTask(msg.GetMsgOption().GetReceiverRole(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetTask(), msg.GetEvidence(), msg.GetBlackOrg())
			t.sendNeedReplayScheduleTask(needReplayScheduleTask)
			replayTaskResult := needReplayScheduleTask.ReceiveResult()

			log.Debugf("Received the reschedule task result from `schedule.ReplaySchedule()`, proposalId: {%s}, taskId: {%s}, partyId: {%s}, the result: %s",
				msg.GetMsgOption().GetProposalId().String(), msg.GetTask().GetTaskId(), party.GetPartyId(), replayTaskResult.String())

			var (
				vote       *carriertwopcpb.PrepareVote
				content    string
				voteOption types.VoteOption
				resource   *types.PrepareVoteResource
			)

			if nil != replayTaskResult.GetErr() {
				voteOption = types.NO
				resource = &types.PrepareVoteResource{}
				content = fmt.Sprintf("will prepare voting `NO` for proposal '%s', as %s", msg.GetMsgOption().GetProposalId().TerminalString(), replayTaskResult.GetErr())

				log.WithError(replayTaskResult.GetErr()).Warnf("Failed to replay schedule task when received prepareMsg, replay result has err, will vote `NO`, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					msg.GetMsgOption().GetProposalId().String(), msg.GetTask().GetTaskId(), party.GetPartyId())
			} else {
				voteOption = types.YES
				resource = types.NewPrepareVoteResource(
					replayTaskResult.GetResource().GetId(),
					replayTaskResult.GetResource().GetIp(),
					replayTaskResult.GetResource().GetPort(),
					replayTaskResult.GetResource().GetPartyId(),
				)
				content = fmt.Sprintf("will prepare voting `YES` for proposal '%s'", msg.GetMsgOption().GetProposalId().TerminalString())

				log.Infof("Succeed to replay schedule task when received prepareMsg, will vote `YES`, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					msg.GetMsgOption().GetProposalId().String(), msg.GetTask().GetTaskId(), party.GetPartyId())
			}
			vote = makePrepareVote(
				msg.GetMsgOption().GetProposalId(),
				role,
				msg.GetMsgOption().GetSenderRole(),
				party.GetPartyId(),
				msg.GetMsgOption().GetSenderPartyId(),
				party,
				voteOption,
				resource,
				timeutils.UnixMsecUint64(),
			)

			// store event about prepare vote
			t.resourceMng.GetDB().StoreTaskEvent(&carriertypespb.TaskEvent{
				Type:       ev.TaskConsensusPrepareEpoch.GetType(),
				TaskId:     proposalTask.GetTaskId(),
				IdentityId: party.GetIdentityId(),
				PartyId:    party.GetPartyId(),
				Content:    content,
				CreateAt:   timeutils.UnixMsecUint64(),
			})

			// Store current peer own vote for checking whether to vote already
			if nmls == types.RemoteNetworkMsg {
				t.state.StorePrepareVote(types.FetchPrepareVote(vote))
			}
			go func() {
				err := t.sendPrepareVote(pid, party, sender, vote)

				var errStr string

				if voteOption == types.NO { // In any case, as long as voting 'NO', Need to clean the local cache
					errStr = "send `NO` prepareVote when replay schedule task failed"
				}
				if nil != err {
					errStr = fmt.Sprintf("send prepareVote `%s` failed", voteOption.String())
					log.WithField("traceId", traceutil.GenerateTraceID(vote)).Errorf("failed to call `sendPrepareVote`, proposalId: {%s}, taskId: {%s}, partyId: {%s}, receiver partyId:{%s}, receiver peerId: {%s}, err: \n%s",
						msg.GetMsgOption().GetProposalId().String(), msg.GetTask().GetTaskId(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderPartyId(), pid, err)
				} else {
					log.WithField("traceId", traceutil.GenerateTraceID(vote)).Debugf("Succeed to call `sendPrepareVote`, proposalId: {%s}, taskId: {%s}, partyId: {%s}, receiver partyId:{%s}, receiver peerId: {%s}",
						msg.GetMsgOption().GetProposalId().String(), msg.GetTask().GetTaskId(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderPartyId(), pid)
				}

				if "" != errStr {
					// release local resource and clean some data  (on task partner)
					t.stopTaskConsensus(errStr, msg.GetMsgOption().GetProposalId(), msg.GetTask().GetTaskId(),
						msg.GetMsgOption().GetReceiverRole(), msg.GetMsgOption().GetSenderRole(), party, sender, types.TaskConsensusInterrupt)
					t.removeOrgProposalStateAndTask(msg.GetMsgOption().GetProposalId(), msg.GetMsgOption().GetReceiverPartyId())
				}
			}()

			return nil
		}

		failedPartyIds := make([]string, 0)

		for _, data := range msg.GetTask().GetTaskData().GetDataSuppliers() {
			if identity.GetIdentityId() == data.GetIdentityId() {
				if err := votingFn(data, commonconstantpb.TaskRole_TaskRole_DataSupplier); nil != err {
					failedPartyIds = append(failedPartyIds, data.GetPartyId())
				}
			}
		}
		for _, data := range msg.GetTask().GetTaskData().GetPowerSuppliers() {
			if identity.GetIdentityId() == data.GetIdentityId() {
				if err := votingFn(data, commonconstantpb.TaskRole_TaskRole_PowerSupplier); nil != err {
					failedPartyIds = append(failedPartyIds, data.GetPartyId())
				}
			}
		}

		for _, data := range msg.GetTask().GetTaskData().GetReceivers() {
			if identity.GetIdentityId() == data.GetIdentityId() {
				if err := votingFn(data, commonconstantpb.TaskRole_TaskRole_Receiver); nil != err {
					failedPartyIds = append(failedPartyIds, data.GetPartyId())
				}
			}
		}
		if len(failedPartyIds) != 0 {
			errCh <- fmt.Errorf("prepare voting failed by [%s], proposaId: {%s}, taskId: {%s}",
				strings.Join(failedPartyIds, ","), msg.GetMsgOption().GetProposalId(), msg.GetTask().GetTaskId())
		}
	}
	close(errCh)
	return <-errCh
}

// (on Publisher)
func (t *Twopc) onPrepareVote(pid peer.ID, prepareVote *types.PrepareVoteWrap, nmls types.NetworkMsgLocationSymbol) error {

	if err := t.state.ContainsOrAddMsg(prepareVote.GetData()); nil != err {
		return err
	}

	vote := fetchPrepareVote(prepareVote)

	// Verify the signature
	_, err := signutil.VerifyMsgSign(vote.GetMsgOption().GetOwner().GetNodeId(), vote.Hash().Bytes(), vote.GetSign())
	if err != nil {
		return fmt.Errorf("verify prepareVote sign %s", err)
	}

	errCh := make(chan error, 0)

	t.asyncCallCh <- func() {

		if t.state.HasNotOrgProposalWithProposalId(vote.GetMsgOption().GetProposalId()) {
			log.Errorf("Failed to check org proposalState whether have been exist when received prepareVote, but it's not exist, proposalId: {%s}",
				vote.GetMsgOption().GetProposalId().String())
			errCh <- fmt.Errorf("%s when received prepareVote", ctypes.ErrProposalNotFound)
			return
		}

		identity, err := t.resourceMng.GetDB().QueryIdentity()
		if nil != err {
			log.WithError(err).Errorf("Failed to call `QueryIdentity()` when received prepareVote, some one is empty, proposalId: {%s}",
				vote.GetMsgOption().GetProposalId().String())

			errCh <- fmt.Errorf("query local identity failed when received prepareVote, %s", err)
			return
		}

		log.WithField("traceId", traceutil.GenerateTraceID(prepareVote.GetData())).Debugf("Received prepareVote, consensusSymbol: {%s}, remote pid: {%s}, prepareVote: %s", nmls.String(), pid, vote.String())

		randomSt, ok := t.state.RandomOrgProposalStateWithProposalId(vote.GetMsgOption().GetProposalId())
		if !ok {
			log.Errorf("Failed to check org proposalState whether have been exist when received prepareVote, but it's not exist, proposalId: {%s}",
				vote.GetMsgOption().GetProposalId().String())
			errCh <- fmt.Errorf("%s when received prepareVote", ctypes.ErrProposalNotFound)
			return
		}

		task, err := t.resourceMng.GetDB().QueryLocalTask(randomSt.GetTaskId())
		if nil != err {
			log.WithError(err).Errorf("Failed to query local task when received prepareVote, proposalId: {%s}, taskId: {%s}",
				vote.GetMsgOption().GetProposalId().String(), task.GetTaskId())
			errCh <- fmt.Errorf("not found local task, on the prepareVote [proposalId: %s, taskId: %s]",
				vote.GetMsgOption().GetProposalId().String(), task.GetTaskId())
			return
		}
		sender := fetchOrgByPartyRole(vote.GetMsgOption().GetSenderPartyId(), vote.GetMsgOption().GetSenderRole(), task)
		receiver := fetchOrgByPartyRole(vote.GetMsgOption().GetReceiverPartyId(), vote.GetMsgOption().GetReceiverRole(), task)
		if nil == sender || nil == receiver {
			log.Errorf("Failed to check vote.MsgOption sender and receiver of prepareVote when received prepareVote, some one is empty, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
				vote.GetMsgOption().GetProposalId().String(), task.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())
			errCh <- fmt.Errorf("%s when received prepareVote", ctypes.ErrConsensusMsgInvalid)
			return
		}
		// verify the receiver is myself ?
		if identity.GetIdentityId() != receiver.GetIdentityId() {
			log.Warnf("Warning verify receiver identityId of prepareVote, receiver is not me, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
				vote.GetMsgOption().GetProposalId().String(), task.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())
			errCh <- fmt.Errorf("%s when received prepareVote", ctypes.ErrConsensusMsgInvalid)
			return
		}

		// find the task sender party proposal state
		orgProposalState, ok := t.state.QueryOrgProposalStateWithProposalIdAndPartyId(vote.GetMsgOption().GetProposalId(), vote.GetMsgOption().GetReceiverPartyId())
		if !ok {
			log.Errorf("Failed to check org proposalState whether have been exist when received prepareVote, but it's not exist, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
				vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())
			errCh <- fmt.Errorf("%s when received prepareVote", ctypes.ErrProposalNotFound)
			return
		}

		// The vote in the consensus prepare epoch can be processed only if the current state is the prepare state
		if orgProposalState.IsNotPreparePeriod() {
			log.Errorf("Failed to check org proposalState priod when received prepareVote, it's not prepare epoch now, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
				vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())
			errCh <- fmt.Errorf("%s when received prepareVote", ctypes.ErrProposalPrepareVoteTimeout)
			return
		}

		// check the proposalTask of proposal exists
		if t.state.HasNotProposalTaskWithTaskIdAndPartyId(orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId()) {
			log.Errorf("%s when received confirmMsg, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
				ctypes.ErrProposalTaskNotFound, vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())
			errCh <- fmt.Errorf("%s, on the confirm msg [proposalId: %s, taskId: %s, partyId: %s]",
				ctypes.ErrProposalTaskNotFound, vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())
			return
		}

		// Voter <the vote sender> voted repeatedly
		if t.state.HasPrepareVoting(vote.GetMsgOption().GetProposalId(), sender) {
			log.Errorf("%s when received prepareVote, they are not same, proposalId: {%s}, taskId: {%s}, partyId: {%s}, vote sender partyId: {%s}",
				ctypes.ErrPrepareVoteRepeatedly, vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId(),
				vote.GetMsgOption().GetSenderPartyId())
			errCh <- fmt.Errorf("%s, on the prepare vote [proposalId: {%s}, taskId: %s, partyId: %s]",
				ctypes.ErrPrepareVoteRepeatedly, vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())
			return
		}

		identityValid, err := t.verifyPartyAndTaskPartner(vote.GetMsgOption().GetSenderRole(), sender, task)
		if nil != err {
			log.WithError(err).Errorf("Failed to call `verifyPrepareVoteRoleIsTaskPartner()` verify prepare vote role when received prepareVote, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
				vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())
			errCh <- fmt.Errorf("check task role of prepareVote failed when received prepareVote, %s", err)
			return
		}
		if !identityValid {
			log.Errorf("The prepare vote role is not include task partners when received prepareVote, they are not same, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
				vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())
			errCh <- fmt.Errorf("%s, on the prepare vote [proposalId: %s, taskId: %s, sender partyId: %s]",
				ctypes.ErrProposalPrepareVoteOwnerInvalid, vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), sender.GetPartyId())
			return
		}

		// verify resource of `YES` vote
		if vote.GetVoteOption() == types.YES && vote.PeerInfoEmpty() {
			log.Errorf("%s when received prepareVote, they are not same, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
				ctypes.ErrProposalPrepareVoteResourceInvalid, vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())
			errCh <- fmt.Errorf("%s, on the prepare vote [proposalId: %s, taskId: %s, sender partyId: %s]",
				ctypes.ErrProposalPrepareVoteOwnerInvalid, vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), sender.GetPartyId())
			return
		}

		// Store vote
		t.state.StorePrepareVote(vote)

		totalNeedVoteCount := uint32(len(task.GetTaskData().GetDataSuppliers()) +
			len(task.GetTaskData().GetPowerSuppliers()) +
			len(task.GetTaskData().GetReceivers()))

		yesVoteCount := t.state.GetTaskPrepareYesVoteCount(vote.GetMsgOption().GetProposalId())
		totalVotedCount := t.state.GetTaskPrepareTotalVoteCount(vote.GetMsgOption().GetProposalId())

		if totalNeedVoteCount == totalVotedCount {

			now := timeutils.UnixMsecUint64()

			// send confirm msg by option `start` to other remote peers,
			// (announce other peer to continue consensus epoch to confirm epoch)
			// and change proposal state from prepare epoch to confirm epoch
			if totalNeedVoteCount == yesVoteCount {

				// change state from prepare epoch to confirm epoch
				t.state.ChangeToConfirm(vote.GetMsgOption().GetProposalId(), vote.GetMsgOption().GetReceiverPartyId(), now)
				t.wal.StoreOrgProposalState(orgProposalState)

				// store confirm peers resource info
				peers := t.makeConfirmTaskPeerDesc(vote.GetMsgOption().GetProposalId())
				t.storeConfirmTaskPeerInfo(vote.GetMsgOption().GetProposalId(), peers)

				go func() {

					log.Infof("PrepareVoting succeed on consensus prepare epoch, the `YES` vote count has enough, will send `START` confirm msg, the `YES` vote count: {%d}, need total count: {%d}, with proposalId: {%s}, taskId: {%s}, partyId: {%s}",
						yesVoteCount, totalNeedVoteCount, vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())
					if err := t.sendConfirmMsg(vote.GetMsgOption().GetProposalId(), task, peers, types.TwopcMsgStart, now); nil != err {
						log.Errorf("Failed to call `sendConfirmMsg` with `start` consensus prepare epoch on `onPrepareVote`, proposalId: {%s}, taskId: {%s}, partyId: {%s}, err: \n%s",
							vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId(), err)
						// Send consensus result to interrupt consensus epoch and clean some data (on task sender)
						t.stopTaskConsensus("send confirmMsg failed", vote.GetMsgOption().GetProposalId(), orgProposalState.GetTaskId(),
							commonconstantpb.TaskRole_TaskRole_Sender, commonconstantpb.TaskRole_TaskRole_Sender, receiver, receiver, types.TaskConsensusInterrupt)
						t.removeOrgProposalStateAndTask(vote.GetMsgOption().GetProposalId(), vote.GetMsgOption().GetReceiverPartyId())
					}
				}()

			} else {

				// send confirm msg by option `stop` to other remote peers,
				// (announce other peer to interrupt consensus epoch)
				// and remove local cache (task/proposal state/prepare vote) about proposal and task
				go func() {

					log.Infof("PrepareVoting failed on consensus prepare epoch, the `YES` vote count is no enough, will send `STOP` confirm msg, the `YES` vote count: {%d}, need total count: {%d}, with proposalId: {%s}, taskId: {%s}, partyId: {%s}",
						yesVoteCount, totalNeedVoteCount, vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())

					if err := t.sendConfirmMsg(vote.GetMsgOption().GetProposalId(), task, t.makeEmptyConfirmTaskPeerDesc(), types.TwopcMsgStop, now); nil != err {
						log.Errorf("Failed to call `sendConfirmMsg` with `stop` consensus prepare epoch on `onPrepareVote`, proposalId: {%s}, taskId: {%s}, partyId: {%s}, err: \n%s",
							vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId(), err)
					}
					// Send consensus result to interrupt consensus epoch and clean some data (on task sender)
					t.stopTaskConsensus("the prepareMsg voting result was not passed", vote.GetMsgOption().GetProposalId(), orgProposalState.GetTaskId(),
						commonconstantpb.TaskRole_TaskRole_Sender, commonconstantpb.TaskRole_TaskRole_Sender, receiver, receiver, types.TaskConsensusInterrupt)
					t.removeOrgProposalStateAndTask(vote.GetMsgOption().GetProposalId(), vote.GetMsgOption().GetReceiverPartyId())
				}()
			}
		}
	}

	close(errCh)
	return <-errCh
}

// (on Subscriber)
func (t *Twopc) onConfirmMsg(pid peer.ID, confirmMsg *types.ConfirmMsgWrap, nmls types.NetworkMsgLocationSymbol) error {

	if err := t.state.ContainsOrAddMsg(confirmMsg.GetData()); nil != err {
		return err
	}

	msg := fetchConfirmMsg(confirmMsg)

	// Verify the signature
	_, err := signutil.VerifyMsgSign(msg.GetMsgOption().GetOwner().GetNodeId(), msg.Hash().Bytes(), msg.GetSign())
	if err != nil {
		return fmt.Errorf("verify confirmMsg sign %s", err)
	}

	errCh := make(chan error, 0)

	t.asyncCallCh <- func() {
		if t.state.HasNotOrgProposalWithProposalId(msg.GetMsgOption().GetProposalId()) {
			log.Errorf("Failed to check org proposalState whether have been exist when received confirmMsg, but it's not exist, proposalId: {%s}",
				msg.GetMsgOption().GetProposalId().String())
			errCh <- fmt.Errorf("%s when received confirmMsg", ctypes.ErrProposalNotFound)
			return
		}

		identity, err := t.resourceMng.GetDB().QueryIdentity()
		if nil != err {
			log.WithError(err).Errorf("Failed to call `QueryIdentity()` when received confirmMsg, proposalId: {%s}",
				msg.GetMsgOption().GetProposalId().String())
			errCh <- fmt.Errorf("query local identity failed when received confirmMsg, %s", err)
			return
		}

		log.WithField("traceId", traceutil.GenerateTraceID(confirmMsg.GetData())).Debugf("Received remote confirmMsg, consensusSymbol: {%s}, remote pid: {%s}, confirmMsg: %s", nmls.String(), pid, msg.String())

		randomSt, ok := t.state.RandomOrgProposalStateWithProposalId(msg.GetMsgOption().GetProposalId())
		if !ok {
			log.Errorf("Failed to check org proposalState whether have been exist when received confirmMsg, but it's not exist, proposalId: {%s}",
				msg.GetMsgOption().GetProposalId().String())
			errCh <- fmt.Errorf("%s when received confirmMsg", ctypes.ErrProposalNotFound)
			return
		}

		task, err := t.resourceMng.GetDB().QueryLocalTask(randomSt.GetTaskId())
		if nil != err {
			log.WithError(err).Errorf("Failed to query local task when received confirmMsg, proposalId: {%s}, taskId: {%s}",
				msg.GetMsgOption().GetProposalId().String(), task.GetTaskId())
			errCh <- fmt.Errorf("not found local task, on the confirmMsg [proposalId: %s, taskId: %s]",
				msg.GetMsgOption().GetProposalId().String(), task.GetTaskId())
			return
		}

		votingFn := func(party *carriertypespb.TaskOrganization, role commonconstantpb.TaskRole) error {

			orgProposalState, ok := t.state.QueryOrgProposalStateWithProposalIdAndPartyId(msg.GetMsgOption().GetProposalId(), party.GetPartyId())
			if !ok {
				log.Errorf("Failed to check org proposalState whether have been exist when received confirmMsg, but it's not exist, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					msg.GetMsgOption().GetProposalId().String(), task.GetTaskId(), party.GetPartyId())
				return fmt.Errorf("%s when received confirmMsg", ctypes.ErrProposalNotFound)
			}

			// The vote in the consensus prepare epoch or confirm epoch can be processed just if the current state is the prepare state or confirm state.
			if orgProposalState.IsCommitPeriod() {
				log.Errorf("Failed to check org proposalState priod when received confirmMsg, it's commit epoch now, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					msg.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), party.GetPartyId())
				return fmt.Errorf("%s when received confirmMsg", ctypes.ErrProposalConfirmMsgTimeout)
			}

			// check the proposalTask of proposal exists
			if t.state.HasNotProposalTaskWithTaskIdAndPartyId(orgProposalState.GetTaskId(), party.GetPartyId()) {
				log.Errorf("%s when received confirmMsg, proposalId: {%s}, taskId: {%s}, , partyId: {%s}",
					ctypes.ErrProposalTaskNotFound, msg.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), party.GetPartyId())
				return fmt.Errorf("%s, on the confirm msg [proposalId: %s, taskId: %s, partyId: %s]",
					ctypes.ErrProposalTaskNotFound, msg.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), party.GetPartyId())
			}

			sender := fetchOrgByPartyRole(msg.GetMsgOption().GetSenderPartyId(), msg.GetMsgOption().GetSenderRole(), task)
			if nil == sender {
				log.Errorf("Failed to check msg.MsgOption sender and receiver of confirmMsg when received confirmMsg, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					msg.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), party.GetPartyId())
				return fmt.Errorf("%s when received confirmMsg", ctypes.ErrConsensusMsgInvalid)
			}

			// Check whether the sender of the message is the same organization as the sender of the task.
			// If not, this message is illegal.
			if task.GetTaskSender().GetIdentityId() != sender.GetIdentityId() ||
				task.GetTaskSender().GetPartyId() != sender.GetPartyId() {
				log.Warnf("Warning the sender of the message is not the same organization as the sender of the task when received confirmMsg, proposalId: {%s}, taskId: {%s}, partyId: {%s}, msg sender: %s, task sender: %s",
					msg.GetMsgOption().GetProposalId().String(), task.GetTaskId(), party.GetPartyId(), sender.String(), task.GetTaskSender().String())
				return fmt.Errorf("%s when received confirmMsg", ctypes.ErrConsensusMsgInvalid)
			}

			// verify the receiver is myself ?
			if identity.GetIdentityId() != party.GetIdentityId() {
				log.Warnf("Warning verify receiver identityId of confirmMsg, receiver is not me, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					msg.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), party.GetPartyId())
				return fmt.Errorf("%s when received confirmMsg", ctypes.ErrConsensusMsgInvalid)
			}

			org := &carriertypespb.TaskOrganization{
				PartyId:    party.GetPartyId(),
				NodeName:   identity.GetNodeName(),
				NodeId:     identity.GetNodeId(),
				IdentityId: identity.GetIdentityId(),
			}

			// If you have already voted then we will not vote again.
			// Cause the local message will only call the local function once,
			// and the remote message needs to prevent receiving the repeated forwarded consensus message.
			if nmls == types.RemoteNetworkMsg && t.state.HasConfirmVoting(msg.GetMsgOption().GetProposalId(), org) {
				log.Errorf("Failed to check remote peer confirm vote wether voting when received confirmMsg, it's voting alreay, proposalId: {%s}, taskId: {%s}, partyId: {%s}, confirmMsgOption: {%s}",
					msg.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), party.GetPartyId(), msg.GetConfirmOption().String())
				return fmt.Errorf("%s when received confirmMsg", ctypes.ErrConfirmVotehadVoted)
			}

			// check msg confirm option value is `start` or `stop` ?
			if msg.GetConfirmOption() == types.TwopcMsgStop || msg.GetConfirmOption() == types.TwopcMsgUnknown {
				log.Warnf("verify confirmMsgOption is not `Start` of confirmMsg when received confirmMsg, proposalId: {%s}, taskId: {%s}, partyId: {%s}, confirmMsgOption: {%s}",
					msg.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), party.GetPartyId(), msg.GetConfirmOption().String())
				// release local resource and clean some data  (on task partner)
				t.stopTaskConsensus(fmt.Sprintf("check confirm option is %s when received confirmMsg",
					msg.GetConfirmOption().String()), msg.GetMsgOption().GetProposalId(), orgProposalState.GetTaskId(),
					role, msg.GetMsgOption().GetSenderRole(), party, sender, types.TaskConsensusInterrupt)
				t.removeOrgProposalStateAndTask(msg.GetMsgOption().GetProposalId(), party.GetPartyId())
				return fmt.Errorf("%s when received confirmMsg", ctypes.ErrConsensusMsgInvalid)
			}

			var (
				vote       *carriertwopcpb.ConfirmVote
				content    string
				voteOption types.VoteOption
			)

			// verify peers resources
			if msg.PeersEmpty() {
				voteOption = types.NO
				content = fmt.Sprintf("will confirm voting `NO` for proposal '%s', as received empty peers on confirm msg", msg.GetMsgOption().GetProposalId().TerminalString())

				log.Warnf("Failed to verify peers resources of confirmMsg when received confirmMsg, the peerDesc reources is empty, will vote `NO`, proposalId: {%s}, taskId: {%s}, partyId: {%s}, confirmMsgOption: {%s}",
					msg.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), party.GetPartyId(), msg.GetConfirmOption().String())

			} else {
				// store confirm peers resource info
				t.storeConfirmTaskPeerInfo(msg.GetMsgOption().GetProposalId(), msg.GetPeers())
				voteOption = types.YES
				content = fmt.Sprintf("will confirm voting `YES` for proposal '%s'", msg.GetMsgOption().GetProposalId().TerminalString())

				log.Infof("Succeed to verify peers resources of confirmMsg when received confirmMsg, will vote `YES`, proposalId: {%s}, taskId: {%s}, partyId: {%s}, confirmMsgOption: {%s}",
					msg.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), party.GetPartyId(), msg.GetConfirmOption().String())

			}
			vote = makeConfirmVote(
				orgProposalState.GetProposalId(),
				role,
				msg.GetMsgOption().GetSenderRole(),
				party.GetPartyId(),
				msg.GetMsgOption().GetSenderPartyId(),
				party,
				voteOption,
				timeutils.UnixMsecUint64(),
			)

			// store event about confirm vote
			t.resourceMng.GetDB().StoreTaskEvent(&carriertypespb.TaskEvent{
				Type:       ev.TaskConsensusConfirmEpoch.GetType(),
				TaskId:     orgProposalState.GetTaskId(),
				IdentityId: party.GetIdentityId(),
				PartyId:    party.GetPartyId(),
				Content:    content,
				CreateAt:   timeutils.UnixMsecUint64(),
			})

			// Store current peer own vote for checking whether to vote already
			if nmls == types.RemoteNetworkMsg {
				t.state.StoreConfirmVote(types.FetchConfirmVote(vote))
			}

			// change state from prepare epoch to confirm epoch
			t.state.ChangeToConfirm(msg.GetMsgOption().GetProposalId(), party.GetPartyId(), msg.GetCreateAt())
			t.wal.StoreOrgProposalState(orgProposalState)

			go func() {
				err := t.sendConfirmVote(pid, party, sender, vote)

				var errStr string

				if voteOption == types.NO { // In any case, as long as voting 'NO', Need to clean the local cache
					errStr = "send `NO` confirmVote when received empty peers confirmMsg"
				}
				if nil != err {
					errStr = fmt.Sprintf("send confirmVote `%s` failed", voteOption.String())
					log.WithField("traceId", traceutil.GenerateTraceID(vote)).Errorf("failed to call `sendConfirmVote`, proposalId: {%s}, taskId: {%s}, partyId: {%s}, \n%s",
						msg.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), party.GetPartyId(), err)
				} else {
					log.WithField("traceId", traceutil.GenerateTraceID(vote)).Debugf("Succeed to call `sendConfirmVote`, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
						msg.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), party.GetPartyId())
				}

				if "" != errStr {
					// release local resource and clean some data  (on task partner)
					t.stopTaskConsensus(errStr, msg.GetMsgOption().GetProposalId(), orgProposalState.GetTaskId(),
						role, msg.GetMsgOption().GetSenderRole(), party, sender, types.TaskConsensusInterrupt)
					t.removeOrgProposalStateAndTask(msg.GetMsgOption().GetProposalId(), party.GetPartyId())
				}
			}()

			return nil
		}

		failedPartyIds := make([]string, 0)

		for _, data := range task.GetTaskData().GetDataSuppliers() {
			if identity.GetIdentityId() == data.GetIdentityId() {
				if err := votingFn(data, commonconstantpb.TaskRole_TaskRole_DataSupplier); nil != err {
					failedPartyIds = append(failedPartyIds, data.GetPartyId())
				}
			}
		}
		for _, data := range task.GetTaskData().GetPowerSuppliers() {
			if identity.GetIdentityId() == data.GetIdentityId() {
				if err := votingFn(data, commonconstantpb.TaskRole_TaskRole_PowerSupplier); nil != err {
					failedPartyIds = append(failedPartyIds, data.GetPartyId())
				}
			}
		}

		for _, data := range task.GetTaskData().GetReceivers() {
			if identity.GetIdentityId() == data.GetIdentityId() {
				if err := votingFn(data, commonconstantpb.TaskRole_TaskRole_Receiver); nil != err {
					failedPartyIds = append(failedPartyIds, data.GetPartyId())
				}
			}
		}
		if len(failedPartyIds) != 0 {
			errCh <- fmt.Errorf("confirm voting failed by [%s], proposaId: {%s}, taskId: {%s}",
				strings.Join(failedPartyIds, ","), msg.GetMsgOption().GetProposalId(), task.GetTaskId())
		}
	}

	close(errCh)
	return <-errCh
}

// (on Publisher)
func (t *Twopc) onConfirmVote(pid peer.ID, confirmVote *types.ConfirmVoteWrap, nmls types.NetworkMsgLocationSymbol) error {

	if err := t.state.ContainsOrAddMsg(confirmVote.GetData()); nil != err {
		return err
	}

	vote := fetchConfirmVote(confirmVote)

	// Verify the signature
	_, err := signutil.VerifyMsgSign(vote.GetMsgOption().GetOwner().GetNodeId(), vote.Hash().Bytes(), vote.GetSign())
	if err != nil {
		return fmt.Errorf("verify confirmVote sign %s", err)
	}

	errCh := make(chan error, 0)

	t.asyncCallCh <- func() {

		if t.state.HasNotOrgProposalWithProposalId(vote.GetMsgOption().GetProposalId()) {
			log.Errorf("Failed to check org proposalState whether have been exist when received confirmVote, but it's not exist, proposalId: {%s}",
				vote.GetMsgOption().GetProposalId().String())
			errCh <- fmt.Errorf("%s when received confirmVote", ctypes.ErrProposalNotFound)
			return
		}

		identity, err := t.resourceMng.GetDB().QueryIdentity()
		if nil != err {
			log.WithError(err).Errorf("Failed to call `QueryIdentity()` when received confirmVote, some one is empty, proposalId: {%s}",
				vote.GetMsgOption().GetProposalId().String())

			errCh <- fmt.Errorf("query local identity failed when received confirmVote, %s", err)
			return
		}

		log.WithField("traceId", traceutil.GenerateTraceID(confirmVote.GetData())).Debugf("Received confirmVote, consensusSymbol: {%s}, remote pid: {%s}, confirmVote: %s", nmls.String(), pid, vote.String())

		randomSt, ok := t.state.RandomOrgProposalStateWithProposalId(vote.GetMsgOption().GetProposalId())
		if !ok {
			log.Errorf("Failed to check org proposalState whether have been exist when received confirmVote, but it's not exist, proposalId: {%s}",
				vote.GetMsgOption().GetProposalId().String())
			errCh <- fmt.Errorf("%s when received confirmVote", ctypes.ErrProposalNotFound)
			return
		}

		task, err := t.resourceMng.GetDB().QueryLocalTask(randomSt.GetTaskId())
		if nil != err {
			log.WithError(err).Errorf("Failed to query local task when received confirmVote, proposalId: {%s}, taskId: {%s}",
				vote.GetMsgOption().GetProposalId().String(), task.GetTaskId())
			errCh <- fmt.Errorf("not found local task, on the confirmVote [proposalId: %s, taskId: %s]",
				vote.GetMsgOption().GetProposalId().String(), task.GetTaskId())
			return
		}
		sender := fetchOrgByPartyRole(vote.GetMsgOption().GetSenderPartyId(), vote.GetMsgOption().GetSenderRole(), task)
		receiver := fetchOrgByPartyRole(vote.GetMsgOption().GetReceiverPartyId(), vote.GetMsgOption().GetReceiverRole(), task)
		if nil == sender || nil == receiver {
			log.Errorf("Failed to check vote.MsgOption sender and receiver of confirmVote when received confirmVote, some one is empty, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
				vote.GetMsgOption().GetProposalId().String(), task.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())
			errCh <- fmt.Errorf("%s when received confirmVote", ctypes.ErrConsensusMsgInvalid)
			return
		}
		// verify the receiver is myself ?
		if identity.GetIdentityId() != receiver.GetIdentityId() {
			log.Warnf("Warning verify receiver identityId of confirmVote, receiver is not me, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
				vote.GetMsgOption().GetProposalId().String(), task.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())
			errCh <- fmt.Errorf("%s when received confirmVote", ctypes.ErrConsensusMsgInvalid)
			return
		}

		// find the task sender party proposal state
		orgProposalState, ok := t.state.QueryOrgProposalStateWithProposalIdAndPartyId(vote.GetMsgOption().GetProposalId(), vote.GetMsgOption().GetReceiverPartyId())
		if !ok {
			log.Errorf("Failed to check org proposalState whether have been exist when received confirmVote, but it's not exist, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
				vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())
			errCh <- fmt.Errorf("%s when received confirmVote", ctypes.ErrProposalNotFound)
			return
		}

		// The vote in the consensus confirm epoch can be processed only if the current state is the confirm state
		if orgProposalState.IsPreparePeriod() {
			log.Errorf("Failed to check org proposalState priod when received confirmVote, it's not confirm epoch and is prepare epoch now, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
				vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())
			errCh <- fmt.Errorf("%s when received confirmVote", ctypes.ErrProposalConfirmVoteFuture)
			return
		}
		if orgProposalState.IsCommitPeriod() {
			log.Errorf("Failed to check org proposalState priod when received confirmVote, it's not confirm epoch and is commit epoch now, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
				vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())
			errCh <- fmt.Errorf("%s when received confirmVote", ctypes.ErrProposalPrepareVoteTimeout)
			return
		}

		// Voter <the vote sender> voted repeatedly
		if t.state.HasConfirmVoting(vote.GetMsgOption().GetProposalId(), sender) {
			log.Errorf("%s when received confirmVote, they are not same, proposalId: {%s}, taskId: {%s}, partyId: {%s}, vote sender partyId: {%s}",
				ctypes.ErrConfirmVoteRepeatedly, vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId(),
				vote.GetMsgOption().GetSenderPartyId())
			errCh <- fmt.Errorf("%s when received confirmVote", ctypes.ErrConfirmVoteRepeatedly)
			return
		}

		identityValid, err := t.verifyPartyAndTaskPartner(vote.GetMsgOption().GetSenderRole(), sender, task)
		if nil != err {
			log.WithError(err).Errorf("Failed to call `verifyPartyAndTaskPartner()` verify confirm vote role when received confirmVote, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
				vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())
			errCh <- fmt.Errorf("check task role of confirmVote failed when received confirmVote, %s", err)
			return
		}
		if !identityValid {
			log.Errorf("The confirm vote role is not include task partners when received confirmVote, they are not same, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
				vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())
			errCh <- fmt.Errorf("%s, on the confirm vote [proposalId: %s, taskId: %s, sender partyId: %s]",
				ctypes.ErrProposalConfirmVoteVoteOwnerInvalid, orgProposalState.GetProposalId(), orgProposalState.GetTaskId(), sender.GetPartyId())
			return
		}

		// Store vote
		t.state.StoreConfirmVote(vote)

		totalNeedVoteCount := uint32(len(task.GetTaskData().GetDataSuppliers()) +
			len(task.GetTaskData().GetPowerSuppliers()) +
			len(task.GetTaskData().GetReceivers()))

		yesVoteCount := t.state.GetTaskConfirmYesVoteCount(vote.GetMsgOption().GetProposalId())
		totalVotedCount := t.state.GetTaskConfirmTotalVoteCount(vote.GetMsgOption().GetProposalId())

		if totalNeedVoteCount == totalVotedCount {

			now := timeutils.UnixMsecUint64()

			// send commit msg by option `start` to other remote peers,
			// (announce other peer to continue consensus epoch to commit epoch)
			// and change proposal state from confirm epoch to commit epoch
			if totalNeedVoteCount == yesVoteCount {

				// change state from confirm epoch to commit epoch
				t.state.ChangeToCommit(vote.GetMsgOption().GetProposalId(), vote.GetMsgOption().GetReceiverPartyId(), now)
				t.wal.StoreOrgProposalState(orgProposalState)

				go func() {

					log.Debugf("ConfirmVoting succeed on consensus confirm epoch, the `YES` vote count has enough, will send `START` commit msg, the `YES` vote count: {%d}, need total count: {%d}, with proposalId: {%s}, taskId: {%s}, partyId: {%s}",
						yesVoteCount, totalNeedVoteCount, vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())

					if err := t.sendCommitMsg(vote.GetMsgOption().GetProposalId(), task, types.TwopcMsgStart, now); nil != err {
						log.Errorf("Failed to call `sendCommitMsg` with `start` consensus confirm epoch on `onConfirmVote`, proposalId: {%s}, taskId: {%s}, partyId: {%s}, err: \n%s",
							vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId(), err)
						// Send consensus result (on task sender)
						t.stopTaskConsensus("send commitMsg failed", vote.GetMsgOption().GetProposalId(), orgProposalState.GetTaskId(),
							commonconstantpb.TaskRole_TaskRole_Sender, commonconstantpb.TaskRole_TaskRole_Sender, receiver, receiver, types.TaskConsensusInterrupt)
					} else {
						// Send consensus result (on task sender)
						t.replyTaskConsensusResult(types.NewTaskConsResult(orgProposalState.GetTaskId(), types.TaskConsensusFinished, nil))
						task, err := t.resourceMng.GetDB().QueryLocalTask(orgProposalState.GetTaskId())
						if err != nil {
							t.identityBlackListCache.CheckConsensusResultOfNoVote(orgProposalState.GetProposalId(), task)
						} else {
							log.Warn("not found task ,taskId is:", orgProposalState.GetTaskId())
						}
					}
					// Finally, whether the commitmsg is sent successfully or not, the local cache needs to be cleared
					t.removeOrgProposalStateAndTask(vote.GetMsgOption().GetProposalId(), vote.GetMsgOption().GetReceiverPartyId())

				}()

			} else {

				// send commit msg by option `stop` to other remote peers,
				// (announce other peer to interrupt consensus epoch)
				// and remove local cache (task/proposal state/prepare vote/confirm vote/peerDesc) about proposal and task
				go func() {

					log.Debugf("ConfirmVoting failed on consensus confirm epoch, the `YES` vote count is no enough, will send `STOP` commit msg, the `YES` vote count: {%d}, need total count: {%d}, with proposalId: {%s}, taskId: {%s}, partyId: {%s}",
						yesVoteCount, totalNeedVoteCount, vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId())

					if err := t.sendCommitMsg(vote.GetMsgOption().GetProposalId(), task, types.TwopcMsgStop, now); nil != err {
						log.Errorf("Failed to call `sendCommitMsg` with `stop` consensus confirm epoch on `onConfirmVote`, proposalId: {%s}, taskId: {%s}, partyId: {%s}, err: \n%s",
							vote.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), vote.GetMsgOption().GetReceiverPartyId(), err)
					}
					// Send consensus result to interrupt consensus epoch and clean some data (on task sender)
					t.stopTaskConsensus("the cofirmMsg voting result was not passed", vote.GetMsgOption().GetProposalId(), orgProposalState.GetTaskId(),
						commonconstantpb.TaskRole_TaskRole_Sender, commonconstantpb.TaskRole_TaskRole_Sender, receiver, receiver, types.TaskConsensusInterrupt)
					t.removeOrgProposalStateAndTask(vote.GetMsgOption().GetProposalId(), vote.GetMsgOption().GetReceiverPartyId())
				}()
			}
		}
	}

	close(errCh)
	return <-errCh
}

// (on Subscriber)
func (t *Twopc) onCommitMsg(pid peer.ID, cimmitMsg *types.CommitMsgWrap, nmls types.NetworkMsgLocationSymbol) error {

	if err := t.state.ContainsOrAddMsg(cimmitMsg.GetData()); nil != err {
		return err
	}

	msg := fetchCommitMsg(cimmitMsg)

	// Verify the signature
	_, err := signutil.VerifyMsgSign(msg.GetMsgOption().GetOwner().GetNodeId(), msg.Hash().Bytes(), msg.GetSign())
	if err != nil {
		return fmt.Errorf("verify confirmVote sign %s", err)
	}

	errCh := make(chan error, 0)

	t.asyncCallCh <- func() {

		if t.state.HasNotOrgProposalWithProposalId(msg.GetMsgOption().GetProposalId()) {
			log.Errorf("Failed to check org proposalState whether have been exist when received commitMsg, but it's not exist, proposalId: {%s}",
				msg.GetMsgOption().GetProposalId().String())
			errCh <- fmt.Errorf("%s when received commitMsg", ctypes.ErrProposalNotFound)
			return
		}

		identity, err := t.resourceMng.GetDB().QueryIdentity()
		if nil != err {
			log.WithError(err).Errorf("Failed to call `QueryIdentity()` when received confirmMsg, proposalId: {%s}",
				msg.GetMsgOption().GetProposalId().String())
			errCh <- fmt.Errorf("query local identity failed when received confirmMsg, %s", err)
			return
		}

		log.WithField("traceId", traceutil.GenerateTraceID(cimmitMsg.GetData())).Debugf("Received commitMsg, consensusSymbol: {%s}, remote pid: {%s}, commitMsg: %s", nmls.String(), pid, msg.String())

		randomSt, ok := t.state.RandomOrgProposalStateWithProposalId(msg.GetMsgOption().GetProposalId())
		if !ok {
			log.Errorf("Failed to check org proposalState whether have been exist when received commitMsg, but it's not exist, proposalId: {%s}",
				msg.GetMsgOption().GetProposalId().String())
			errCh <- fmt.Errorf("%s when received commitMsg", ctypes.ErrProposalNotFound)
			return
		}

		task, err := t.resourceMng.GetDB().QueryLocalTask(randomSt.GetTaskId())
		if nil != err {
			log.WithError(err).Errorf("Failed to query local task when received commitMsg, proposalId: {%s}, taskId: {%s}",
				msg.GetMsgOption().GetProposalId().String(), task.GetTaskId())
			errCh <- fmt.Errorf("not found local task, on the commitMsg [proposalId: %s, taskId: %s]",
				msg.GetMsgOption().GetProposalId().String(), task.GetTaskId())
			return
		}

		driveTaskFn := func(party *carriertypespb.TaskOrganization, role commonconstantpb.TaskRole) error {

			orgProposalState, ok := t.state.QueryOrgProposalStateWithProposalIdAndPartyId(msg.GetMsgOption().GetProposalId(), party.GetPartyId())
			if !ok {
				log.Errorf("Failed to check org proposalState whether have been exist when received commitMsg, but it's not exist, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					msg.GetMsgOption().GetProposalId().String(), task.GetTaskId(), party.GetPartyId())
				return fmt.Errorf("%s when received commitMsg", ctypes.ErrProposalNotFound)
			}

			// The vote in the consensus confirm epoch or commit epoch can be processed just if the current state is the confirm state or commit state
			if orgProposalState.IsPreparePeriod() {
				log.Errorf("Failed to check org proposalState priod when received commitMsg, it's not commit epoch and is prepare epoch now, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					msg.GetMsgOption().GetProposalId().String(), task.GetTaskId(), party.GetPartyId())
				return fmt.Errorf("%s when received commitMsg", ctypes.ErrProposalCommitMsgFuture)
			}
			if orgProposalState.IsFinishedPeriod() {
				log.Errorf("Failed to check org proposalState priod when received commitMsg, it's not commit epoch and is finished epoch now, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					msg.GetMsgOption().GetProposalId().String(), task.GetTaskId(), party.GetPartyId())
				return fmt.Errorf("%s when received commitMsg", ctypes.ErrProposalCommitMsgTimeout)
			}

			// check the proposalTask of proposal exists
			if t.state.HasNotProposalTaskWithTaskIdAndPartyId(orgProposalState.GetTaskId(), party.GetPartyId()) {
				log.Errorf("%s when received commitMsg, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					ctypes.ErrProposalTaskNotFound, msg.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), party.GetPartyId())
				return fmt.Errorf("%s, on the commit msg [proposalId: %s, taskId: %s, partyId: %s]",
					ctypes.ErrProposalTaskNotFound, msg.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), party.GetPartyId())
			}

			sender := fetchOrgByPartyRole(msg.GetMsgOption().GetSenderPartyId(), msg.GetMsgOption().GetSenderRole(), task)

			if nil == sender {
				log.Errorf("Failed to check msg.MsgOption sender and receiver of commitMsg when received commitMsg, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					msg.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), party.GetPartyId())
				return fmt.Errorf("%s when received commitMsg", ctypes.ErrConsensusMsgInvalid)
			}

			// Check whether the sender of the message is the same organization as the sender of the task.
			// If not, this message is illegal.
			if task.GetTaskSender().GetIdentityId() != sender.GetIdentityId() ||
				task.GetTaskSender().GetPartyId() != sender.GetPartyId() {
				log.Warnf("Warning the sender of the message is not the same organization as the sender of the task when received commitMsg, proposalId: {%s}, taskId: {%s}, partyId: {%s}, msg sender: %s, task sender: %s",
					msg.GetMsgOption().GetProposalId().String(), task.GetTaskId(), party.GetPartyId(), sender.String(), task.GetTaskSender().String())
				return fmt.Errorf("%s when received commitMsg", ctypes.ErrConsensusMsgInvalid)
			}

			// verify the receiver is myself ?
			if identity.GetIdentityId() != party.GetIdentityId() {
				log.Warnf("Warning verify receiver identityId of commitMsg, receiver is not me, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					msg.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), party.GetPartyId())
				return fmt.Errorf("%s when received commitMsg", ctypes.ErrConsensusMsgInvalid)
			}

			// check msg commit option value is `start` or `stop` ?
			if msg.GetCommitOption() == types.TwopcMsgStop || msg.GetCommitOption() == types.TwopcMsgUnknown {
				log.Warnf("verify commitMsgOption is not `Start` of commitMsg when received commitMsg, proposalId: {%s}, taskId: {%s}, partyId: {%s}, confirmMsgOption: {%s}",
					msg.GetMsgOption().GetProposalId().String(), orgProposalState.GetTaskId(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetCommitOption().String())
				// release local resource and clean some data  (on task partner)
				t.stopTaskConsensus(fmt.Sprintf("check commit option is %s when received commitMsg", msg.GetCommitOption().String()), msg.GetMsgOption().GetProposalId(), orgProposalState.GetTaskId(),
					role, msg.GetMsgOption().GetSenderRole(), party, sender, types.TaskConsensusInterrupt)
				t.removeOrgProposalStateAndTask(msg.GetMsgOption().GetProposalId(), msg.GetMsgOption().GetReceiverPartyId())
				return fmt.Errorf("%s when received commitMsg", ctypes.ErrConsensusMsgInvalid)
			}

			// change state from confirm epoch to commit epoch
			t.state.ChangeToCommit(msg.GetMsgOption().GetProposalId(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetCreateAt())
			t.wal.StoreOrgProposalState(orgProposalState)

			go func() {

				// store succeed consensus event for partyId
				t.resourceMng.GetDB().StoreTaskEvent(&carriertypespb.TaskEvent{
					Type:       ev.TaskSucceedConsensus.GetType(),
					TaskId:     orgProposalState.GetTaskId(),
					IdentityId: party.GetIdentityId(),
					PartyId:    party.GetPartyId(),
					Content:    fmt.Sprintf("succeed consensus."),
					CreateAt:   timeutils.UnixMsecUint64(),
				})
				// If receiving `CommitMsg` is successful,
				// we will forward `schedTask` to `taskManager` to send it to `Fighter` to execute the task.
				t.driveTask(pid, msg.GetMsgOption().GetProposalId(), role, party, msg.GetMsgOption().GetSenderRole(), sender, orgProposalState.GetTaskId())
				t.removeOrgProposalStateAndTask(msg.GetMsgOption().GetProposalId(), msg.GetMsgOption().GetReceiverPartyId())
			}()

			// Finally, it is left 'taskmanager' to call 'releaselocalresourcewithtask()' to release local resources after handle `driveTask()`.
			// No more processing here.
			return nil
		}

		failedPartyIds := make([]string, 0)

		for _, data := range task.GetTaskData().GetDataSuppliers() {
			if identity.GetIdentityId() == data.GetIdentityId() {
				if err := driveTaskFn(data, commonconstantpb.TaskRole_TaskRole_DataSupplier); nil != err {
					failedPartyIds = append(failedPartyIds, data.GetPartyId())
				}
			}
		}
		for _, data := range task.GetTaskData().GetPowerSuppliers() {
			if identity.GetIdentityId() == data.GetIdentityId() {
				if err := driveTaskFn(data, commonconstantpb.TaskRole_TaskRole_PowerSupplier); nil != err {
					failedPartyIds = append(failedPartyIds, data.GetPartyId())
				}
			}
		}

		for _, data := range task.GetTaskData().GetReceivers() {
			if identity.GetIdentityId() == data.GetIdentityId() {
				if err := driveTaskFn(data, commonconstantpb.TaskRole_TaskRole_Receiver); nil != err {
					failedPartyIds = append(failedPartyIds, data.GetPartyId())
				}
			}
		}
		if len(failedPartyIds) != 0 {
			errCh <- fmt.Errorf("driving task failed by [%s], proposaId: {%s}, taskId: {%s}",
				strings.Join(failedPartyIds, ","), msg.GetMsgOption().GetProposalId(), task.GetTaskId())
		}
	}

	close(errCh)
	return <-errCh
}

func (t *Twopc) onTerminateTaskConsensus(pid peer.ID, msg *types.TerminateConsensusMsgWrap) error {

	msgOption := types.FetchMsgOption(msg.GetMsgOption())
	log.Infof("Start terminate task consensus when received terminateConsensusMsg , taskId: {%s}, partyId: {%s}", msg.GetTaskId(), msgOption.GetReceiverPartyId())

	// find the task of proposal on proposalTask
	proposalTask, ok := t.state.QueryProposalTaskWithTaskIdAndPartyId(msg.GetTaskId(), msgOption.GetReceiverPartyId())
	if !ok {
		log.Errorf("%s when received terminateConsensusMsg, taskId: {%s}, partyId: {%s}", ctypes.ErrProposalTaskNotFound, msg.GetTaskId(), msgOption.GetReceiverPartyId())
		return fmt.Errorf("%s, when received terminateConsensusMsg [taskId: %s, partyId: %s]",
			ctypes.ErrProposalTaskNotFound, msg.GetTaskId(), msgOption.GetReceiverPartyId())
	}

	task, err := t.resourceMng.GetDB().QueryLocalTask(proposalTask.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task when received terminateConsensusMsg, taskId: {%s}, partyId: {%s}", msg.GetTaskId(), msgOption.GetReceiverPartyId())
		return fmt.Errorf("not found local task when received terminateConsensusMsg [proposalId: %s, taskId: %s, partyId: %s]",
			proposalTask.GetProposalId().String(), msg.GetTaskId(), msgOption.GetReceiverPartyId())
	}

	if t.state.HasNotOrgProposalWithProposalIdAndPartyId(proposalTask.GetProposalId(), msgOption.GetReceiverPartyId()) {
		log.Errorf("Failed to check org proposalState whether have been exist when received terminateConsensusMsg, but it's not exist, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
			proposalTask.GetProposalId().String(), msg.GetTaskId(), msgOption.GetReceiverPartyId())
		return fmt.Errorf("%s, when received terminateConsensusMsg", ctypes.ErrProposalNotFound)
	}

	sender := fetchOrgByPartyRole(msgOption.GetSenderPartyId(), msgOption.GetSenderRole(), task)
	receiver := fetchOrgByPartyRole(msgOption.GetReceiverPartyId(), msgOption.GetReceiverRole(), task)
	if nil == sender || nil == receiver {
		log.Errorf("Failed to check msg.MsgOption sender and receiver of interruptMsg when received terminateConsensusMsg, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
			proposalTask.GetProposalId().String(), msg.GetTaskId(), msgOption.GetReceiverPartyId())
		return ctypes.ErrConsensusMsgInvalid
	}

	identity, err := t.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryIdentity()` when received terminateConsensusMsg, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
			msgOption.GetProposalId().String(), proposalTask.GetTaskId(), msgOption.GetReceiverPartyId())
		// release local resource and clean some data  (on task partner)
		t.stopTaskConsensus(fmt.Sprintf("query local identity failed %s, when received terminateConsensusMsg", err), msgOption.GetProposalId(), proposalTask.GetTaskId(),
			msgOption.GetReceiverRole(), msgOption.GetSenderRole(), receiver, sender, types.TaskConsensusInterrupt)
		t.removeOrgProposalStateAndTask(msgOption.GetProposalId(), msgOption.GetReceiverPartyId())
		return fmt.Errorf("query local identity failed when received terminateConsensusMsg, %s", err)
	}

	// verify the receiver is myself ?
	if identity.GetIdentityId() != receiver.GetIdentityId() {
		log.Warnf("Warning verify receiver identityId of terminateConsensusMsg, receiver is not me, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
			msgOption.GetProposalId().String(), proposalTask.GetTaskId(), msgOption.GetReceiverPartyId())
		return fmt.Errorf("%s when received terminateConsensusMsg", ctypes.ErrConsensusMsgInvalid)
	}

	// Check whether the sender of the message is the same organization as the sender of the task.
	// If not, this message is illegal.
	if task.GetTaskSender().GetIdentityId() != sender.GetIdentityId() ||
		task.GetTaskSender().GetPartyId() != sender.GetPartyId() {
		log.Warnf("Warning the sender of the message is not the same organization as the sender of the task when received terminateConsensusMsg, proposalId: {%s}, taskId: {%s}, partyId: {%s}, msg sender: %s, task sender: %s",
			msgOption.GetProposalId().String(), task.GetTaskId(), msgOption.GetReceiverPartyId(), sender.String(), task.GetTaskSender().String())
		return fmt.Errorf("%s when received terminateConsensusMsg", ctypes.ErrConsensusMsgInvalid)
	}

	orgProposalState := t.state.MustQueryOrgProposalStateWithProposalIdAndPartyId(proposalTask.GetProposalId(), msgOption.GetReceiverPartyId())
	switch orgProposalState.GetPeriodNum() {
	case ctypes.PeriodPrepare:
		// remove `proposal state` and `task cache` AND inerrupt consensus with sender OR release local locked resource with partner
		t.stopTaskConsensus("interrupt consensus with terminate task while prepare epoch", proposalTask.GetProposalId(), msg.GetTaskId(),
			msgOption.GetReceiverRole(), msgOption.GetSenderRole(), receiver, sender, types.TaskTerminate)
		t.removeOrgProposalStateAndTask(proposalTask.GetProposalId(), proposalTask.GetTaskId())
	case ctypes.PeriodConfirm:
		// remove `proposal state` and `task cache` AND inerrupt consensus with sender OR release local locked resource with partner
		t.stopTaskConsensus("interrupt consensus with terminate task while confirm epoch", proposalTask.GetProposalId(), msg.GetTaskId(),
			msgOption.GetReceiverRole(), msgOption.GetSenderRole(), receiver, sender, types.TaskTerminate)
		t.removeOrgProposalStateAndTask(proposalTask.GetProposalId(), proposalTask.GetTaskId())
	case ctypes.PeriodCommit, ctypes.PeriodFinished:
		// need send terminate msg with task manager
		// so do nothing here
	default:
		log.Errorf("unknown org proposalState priod when received terminateConsensusMsg, proposalId: {%s}, taskId: {%s}, partyId: {%s}, peroid: {%s}",
			proposalTask.GetProposalId().String(), msg.GetTaskId(), msgOption.GetReceiverPartyId(), orgProposalState.GetPeriodStr())
		return fmt.Errorf("unknown org proposalState priod, on the interrupt consensus")
	}
	return nil
}

func (t *Twopc) Get2PcProposalStateByTaskId(taskId string) (*carrierrpcdebugpbv1.Get2PcProposalStateResponse, error) {
	t.state.proposalTaskLock.RLock()
	defer t.state.proposalTaskLock.RUnlock()
	taskObj := t.state.proposalTaskCache[taskId]
	var proposalId common.Hash
	for partyId := range taskObj {
		if obj, err := taskObj[partyId]; err {
			proposalId = obj.ProposalId
			break
		}
	}

	currentTime := time.Now().UnixNano()
	proposalStateInfo := make(map[string]*carrierrpcdebugpbv1.ProposalState, 0)
	if proposalState, ok := t.state.proposalSet[proposalId]; ok {
		for partyId, obj := range proposalState {
			proposalStateInfo[partyId] = &carrierrpcdebugpbv1.ProposalState{
				PeriodNum:            uint32(obj.GetPeriodNum()),
				TaskId:               obj.GetTaskId(),
				ConsumeTime:          uint64(currentTime) - obj.GetStartAt(),
				TaskSenderIdentityId: obj.GetTaskSender().IdentityId,
			}
		}
	} else {
		return &carrierrpcdebugpbv1.Get2PcProposalStateResponse{}, nil
	}

	return &carrierrpcdebugpbv1.Get2PcProposalStateResponse{
		ProposalId: proposalId.String(),
		State:      proposalStateInfo,
	}, nil
}
func (t *Twopc) Get2PcProposalStateByProposalId(proposalId string) (*carrierrpcdebugpbv1.Get2PcProposalStateResponse, error) {
	currentTime := time.Now().UnixNano()
	proposalStateInfo := make(map[string]*carrierrpcdebugpbv1.ProposalState, 0)
	t.state.proposalsLock.RLock()
	defer t.state.proposalsLock.RUnlock()
	proposalState, ok := t.state.proposalSet[common.HexToHash(proposalId)]
	if ok {
		for partyId, obj := range proposalState {
			proposalStateInfo[partyId] = &carrierrpcdebugpbv1.ProposalState{
				PeriodNum:            uint32(obj.GetPeriodNum()),
				TaskId:               obj.GetTaskId(),
				ConsumeTime:          uint64(currentTime) - obj.GetStartAt(),
				TaskSenderIdentityId: obj.GetTaskSender().IdentityId,
			}
		}
	} else {
		return &carrierrpcdebugpbv1.Get2PcProposalStateResponse{}, nil
	}
	return &carrierrpcdebugpbv1.Get2PcProposalStateResponse{
		ProposalId: proposalId,
		State:      proposalStateInfo,
	}, nil
}
func (t *Twopc) Get2PcProposalPrepare(proposalId string) (*carrierrpcdebugpbv1.Get2PcProposalPrepareResponse, error) {
	t.state.prepareVotesLock.RLock()
	defer t.state.prepareVotesLock.RUnlock()
	prepareVoteInfo, ok := t.state.prepareVotes[common.HexToHash(proposalId)]
	if !ok {
		return &carrierrpcdebugpbv1.Get2PcProposalPrepareResponse{}, nil
	}
	votes := make(map[string]*carriertwopcpb.PrepareVote, 0)
	for partyId, obj := range prepareVoteInfo.votes {
		votes[partyId] = &carriertwopcpb.PrepareVote{
			MsgOption: &carriernetmsgcommonpb.MsgOption{
				ProposalId:      obj.MsgOption.ProposalId.Bytes(),
				SenderRole:      uint64(obj.MsgOption.SenderRole),
				SenderPartyId:   []byte(obj.MsgOption.SenderPartyId),
				ReceiverRole:    uint64(obj.MsgOption.ReceiverRole),
				ReceiverPartyId: []byte(obj.MsgOption.ReceiverPartyId),
				MsgOwner: &carriernetmsgcommonpb.TaskOrganizationIdentityInfo{
					Name:       []byte(obj.MsgOption.Owner.GetNodeName()),
					NodeId:     []byte(obj.MsgOption.Owner.GetNodeId()),
					IdentityId: []byte(obj.MsgOption.Owner.GetIdentityId()),
					PartyId:    []byte(obj.MsgOption.Owner.GetPartyId()),
				},
			},
			VoteOption: obj.VoteOption.Bytes(),
			CreateAt:   obj.CreateAt,
			Sign:       obj.Sign,
		}
	}

	yesVotes := make(map[string]uint32, 0)
	for role, voteCount := range prepareVoteInfo.yesVotes {
		yesVotes[role.String()] = voteCount
	}

	voteStatus := make(map[string]uint32, 0)
	for role, voteCount := range prepareVoteInfo.voteStatus {
		voteStatus[role.String()] = voteCount
	}
	return &carrierrpcdebugpbv1.Get2PcProposalPrepareResponse{
		Votes:      votes,
		YesVotes:   yesVotes,
		VoteStatus: voteStatus,
	}, nil
}
func (t *Twopc) Get2PcProposalConfirm(proposalId string) (*carrierrpcdebugpbv1.Get2PcProposalConfirmResponse, error) {
	t.state.confirmVotesLock.RLock()
	defer t.state.confirmVotesLock.RUnlock()
	confirmVoteInfo, ok := t.state.confirmVotes[common.HexToHash(proposalId)]
	if !ok {
		return &carrierrpcdebugpbv1.Get2PcProposalConfirmResponse{}, nil
	}
	votes := make(map[string]*carriertwopcpb.ConfirmVote, 0)
	for partyId, obj := range confirmVoteInfo.votes {
		votes[partyId] = &carriertwopcpb.ConfirmVote{
			MsgOption: &carriernetmsgcommonpb.MsgOption{
				ProposalId:      obj.MsgOption.ProposalId.Bytes(),
				SenderRole:      uint64(obj.MsgOption.SenderRole),
				SenderPartyId:   []byte(obj.MsgOption.SenderPartyId),
				ReceiverRole:    uint64(obj.MsgOption.ReceiverRole),
				ReceiverPartyId: []byte(obj.MsgOption.ReceiverPartyId),
				MsgOwner: &carriernetmsgcommonpb.TaskOrganizationIdentityInfo{
					Name:       []byte(obj.MsgOption.Owner.GetNodeName()),
					NodeId:     []byte(obj.MsgOption.Owner.GetNodeId()),
					IdentityId: []byte(obj.MsgOption.Owner.GetIdentityId()),
					PartyId:    []byte(obj.MsgOption.Owner.GetPartyId()),
				},
			},
			VoteOption: obj.VoteOption.Bytes(),
			CreateAt:   obj.CreateAt,
			Sign:       obj.Sign,
		}
	}

	yesVotes := make(map[string]uint32, 0)
	for role, voteCount := range confirmVoteInfo.yesVotes {
		yesVotes[role.String()] = voteCount
	}

	voteStatus := make(map[string]uint32, 0)
	for role, voteCount := range confirmVoteInfo.voteStatus {
		voteStatus[role.String()] = voteCount
	}
	return &carrierrpcdebugpbv1.Get2PcProposalConfirmResponse{
		Votes:      votes,
		YesVotes:   yesVotes,
		VoteStatus: voteStatus,
	}, nil
}
func (t *Twopc) HasPrepareVoting(proposalId common.Hash, org *carriertypespb.TaskOrganization) bool {
	return t.state.HasPrepareVoting(proposalId, org)
}
func (t *Twopc) HasConfirmVoting(proposalId common.Hash, org *carriertypespb.TaskOrganization) bool {
	return t.state.HasConfirmVoting(proposalId, org)
}
