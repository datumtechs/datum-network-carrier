package twopc

import (
	"bytes"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/bytesutil"
	"github.com/RosettaFlow/Carrier-Go/common/rlputil"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	ev "github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"strings"
	"sync"
	"time"
)

type ConsensusMsgLocationSymbol bool

func (cls ConsensusMsgLocationSymbol) string() string {
	switch cls {
	case LocalConsensusMsg:
		return "local consensus msg"
	case RemoteConsensusMsg:
		return "remote consensus msg"
	default:
		return "Unknown consensus location symbol"

	}
}

const (
	//defaultCleanExpireProposalInterval  = 30 * time.Millisecond
	defaultRefreshProposalStateInternal = 300 * time.Millisecond

	LocalConsensusMsg  ConsensusMsgLocationSymbol = true
	RemoteConsensusMsg ConsensusMsgLocationSymbol = false
)

type Twopc struct {
	config      *Config
	p2p         p2p.P2P
	state       *state
	resourceMng *resource.Manager
	// send remote task to `Scheduler` to replay
	needReplayScheduleTaskCh chan *types.NeedReplayScheduleTask
	// send has was consensus remote tasks to taskManager
	needExecuteTaskCh chan *types.NeedExecuteTask
	asyncCallCh       chan func()
	quit              chan struct{}
	taskResultBusCh   chan *types.TaskConsResult
	taskResultChSet   map[string]chan<- *types.TaskConsResult
	taskResultLock    sync.Mutex
	wal               *walDB
	Errs              []error
}

func New(
	conf *Config,
	resourceMng *resource.Manager,
	p2p p2p.P2P,
	needReplayScheduleTaskCh chan *types.NeedReplayScheduleTask,
	needExecuteTaskCh chan *types.NeedExecuteTask,

) *Twopc {
	newWalDB := newWal(conf)
	return &Twopc{
		config:                   conf,
		p2p:                      p2p,
		state:                    newState(newWalDB),
		resourceMng:              resourceMng,
		needReplayScheduleTaskCh: needReplayScheduleTaskCh,
		needExecuteTaskCh:        needExecuteTaskCh,
		asyncCallCh:              make(chan func(), conf.PeerMsgQueueSize),
		quit:                     make(chan struct{}),
		taskResultBusCh:          make(chan *types.TaskConsResult, 100),
		taskResultChSet:          make(map[string]chan<- *types.TaskConsResult, 100),
		wal:                      newWalDB,
		Errs:                     make([]error, 0),
	}
}

func (t *Twopc) Start() error {
	go t.loop()
	log.Info("Started 2pc consensus engine ...")
	return nil
}
func (t *Twopc) Close() error {
	close(t.quit)
	return nil
}
func (t *Twopc) loop() {
	refreshProposalStateTicker := time.NewTicker(defaultRefreshProposalStateInternal)
	for {
		select {

		case fn := <-t.asyncCallCh:
			fn()

		case res := <-t.taskResultBusCh:
			if nil == res {
				return
			}
			t.handleTaskConsensusResult(res)

		case <-refreshProposalStateTicker.C:
			go t.refreshProposalState()

		case <-t.quit:
			log.Info("Stopped 2pc consensus engine ...")
			return
		}
	}
}

func (t *Twopc) OnConsensusMsg(pid peer.ID, msg types.ConsensusMsg) error {

	switch msg := msg.(type) {
	case *types.PrepareMsgWrap:
		return t.onPrepareMsg(pid, msg, RemoteConsensusMsg)
	case *types.PrepareVoteWrap:
		return t.onPrepareVote(pid, msg, RemoteConsensusMsg)
	case *types.ConfirmMsgWrap:
		return t.onConfirmMsg(pid, msg, RemoteConsensusMsg)
	case *types.ConfirmVoteWrap:
		return t.onConfirmVote(pid, msg, RemoteConsensusMsg)
	case *types.CommitMsgWrap:
		return t.onCommitMsg(pid, msg, RemoteConsensusMsg)
	case *types.InterruptMsgWrap: // Must be  local msg
		return t.onTerminateTaskConsensus(pid, msg)
	default:
		return fmt.Errorf("Unknown the 2pc msg type")

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

func (t *Twopc) OnPrepare(task *types.Task) error {

	return nil
}
func (t *Twopc) OnHandle(task *types.Task, result chan<- *types.TaskConsResult) error {

	t.addTaskResultCh(task.GetTaskId(), result)

	if t.state.HasProposalTaskWithPartyId(task.GetTaskId(), task.GetTaskSender().GetPartyId()) {
		log.Errorf("Failed to check org proposalTask whether have been not exist on OnHandle, but it's alreay exist, taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetTaskSender().GetPartyId())
		t.stopTaskConsensus(ctypes.ErrPrososalTaskIsProcessed.Error(), common.Hash{}, task.GetTaskId(),
			apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Sender, task.GetTaskSender(), task.GetTaskSender(), task, types.TaskConsensusInterrupt)
		return ctypes.ErrPrososalTaskIsProcessed
	}

	// Store task execute status `cons` before consensus when send task prepareMsg to remote peers
	if err := t.resourceMng.GetDB().StoreLocalTaskExecuteStatusValConsByPartyId(task.GetTaskId(), task.GetTaskSender().GetPartyId()); nil != err {
		log.WithError(err).Errorf("Failed to store local task about `cons` status on OnHandle,  taskId: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetTaskSender().GetPartyId())
		t.stopTaskConsensus("store task executeStatus about `cons` failed", common.Hash{}, task.GetTaskId(),
			apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Sender, task.GetTaskSender(), task.GetTaskSender(), task, types.TaskConsensusInterrupt)
		return err
	}

	now := timeutils.UnixMsecUint64()

	var buf bytes.Buffer
	buf.Write(t.config.Option.NodeID.Bytes())
	buf.Write([]byte(task.GetTaskId()))
	buf.Write([]byte(task.GetTaskData().GetTaskName()))
	buf.Write(bytesutil.Uint64ToBytes(task.GetTaskData().GetCreateAt()))
	buf.Write(bytesutil.Uint64ToBytes(now))
	proposalId := rlputil.RlpHash(buf.Bytes())

	log.Infof("Generate proposal, proposalId: {%s}, taskId: {%s}, partyId: {%s}", proposalId.String(), task.GetTaskId(), task.GetTaskSender().GetPartyId())

	// Store some local cache
	t.storeOrgProposalState(
		proposalId,
		task.GetTaskId(),
		task.GetTaskSender(),
		ctypes.NewOrgProposalState(task.GetTaskId(), apicommonpb.TaskRole_TaskRole_Sender, task.GetTaskSender(), now),
	)
	t.state.StoreProposalTaskWithPartyId(task.GetTaskSender().GetPartyId(), types.NewProposalTask(proposalId, task, now))

	// Start handle task ...
	go func() {

		if err := t.sendPrepareMsg(proposalId, task, now); nil != err {
			log.Errorf("Failed to call `sendPrepareMsg`, consensus epoch finished, proposalId: {%s}, taskId: {%s}, partyId: {%s}, err: \n%s",
				proposalId.String(), task.GetTaskId(), task.GetTaskSender().GetPartyId(), err)
			// Send consensus result to Scheduler
			t.stopTaskConsensus("failed to call `sendPrepareMsg`", proposalId, task.GetTaskId(),
				apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Sender, task.GetTaskSender(), task.GetTaskSender(), task, types.TaskConsensusInterrupt)
			// clean some invalid data
			t.removeOrgProposalStateAndTask(proposalId, task.GetTaskSender().GetPartyId())
		}
	}()
	return nil
}

// Handle the prepareMsg from the task pulisher peer (on Subscriber)
func (t *Twopc) onPrepareMsg(pid peer.ID, prepareMsg *types.PrepareMsgWrap, consensusSymbol ConsensusMsgLocationSymbol) error {

	msg, err := fetchPrepareMsg(prepareMsg)
	if nil != err {
		return err
	}
	log.Debugf("Received remote prepareMsg, remote pid: {%s}, consensusSymbol: {%s}, prepareMsg: %s", pid, consensusSymbol.string(), msg.String())

	if t.state.HasOrgProposalWithPartyId(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId) {
		log.Errorf("Failed to check org proposalState whether have been not exist on onPrepareMsg, but it's alreay exist, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			msg.MsgOption.ProposalId.String(), msg.TaskInfo.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
		return ctypes.ErrProposalAlreadyProcessed
	}

	identity, err := t.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryIdentity()` on onPrepareMsg, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			msg.MsgOption.ProposalId.String(), msg.TaskInfo.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
		return fmt.Errorf("query local identity failed, %s", err)
	}

	sender := fetchOrgByPartyRole(msg.MsgOption.SenderPartyId, msg.MsgOption.SenderRole, msg.TaskInfo)
	receiver := fetchOrgByPartyRole(msg.MsgOption.ReceiverPartyId, msg.MsgOption.ReceiverRole, msg.TaskInfo)
	if nil == sender || nil == receiver {
		log.Errorf("\"Failed to check msg.MsgOption sender and receiver on onPrepareMsg, some one is empty, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			msg.MsgOption.ProposalId.String(), msg.TaskInfo.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
		return ctypes.ErrConsensusMsgInvalid
	}

	// verify the receiver is myself ?
	if identity.GetIdentityId() != receiver.GetIdentityId() {
		log.Errorf("Failed to verify receiver identityId of prepareMsg, receiver is not me, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, my identityId: {%s}, receiver identityId: {%s}",
			msg.MsgOption.ProposalId.String(), msg.TaskInfo.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId, identity.GetIdentityId(), receiver.GetIdentityId())
		return ctypes.ErrConsensusMsgInvalid
	}

	org := &apicommonpb.TaskOrganization{
		PartyId:    msg.MsgOption.ReceiverPartyId,
		NodeName:   identity.GetNodeName(),
		NodeId:     identity.GetNodeId(),
		IdentityId: identity.GetIdentityId(),
	}

	// If you have already voted then we will not vote again.
	// Cause the local message will only call the local function once,
	// and the remote message needs to prevent receiving the repeated forwarded consensus message.
	if consensusSymbol == RemoteConsensusMsg && t.state.HasPrepareVoting(msg.MsgOption.ProposalId, org) {
		log.Errorf("Failed to check remote peer prepare vote wether exist on onPrepareMsg, it's exist alreay, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			msg.MsgOption.ProposalId.String(), msg.TaskInfo.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
		return ctypes.ErrPrepareVotehadVoted
	}

	// Store task execute status `cons` before consensus when received a remote task prepareMsg
	if err := t.resourceMng.GetDB().StoreLocalTaskExecuteStatusValConsByPartyId(msg.TaskInfo.GetTaskId(), msg.MsgOption.ReceiverPartyId); nil != err {
		log.WithError(err).Errorf("Failed to store local task about `cons` status on onPrepareMsg, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			msg.MsgOption.ProposalId.String(), msg.TaskInfo.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
		return err
	}

	// Store some local cache
	t.storeOrgProposalState(
		msg.MsgOption.ProposalId,
		msg.TaskInfo.GetTaskId(),
		sender,
		ctypes.NewOrgProposalState(msg.TaskInfo.GetTaskId(), msg.MsgOption.ReceiverRole, receiver, msg.CreateAt),
	)
	t.state.StoreProposalTaskWithPartyId(msg.MsgOption.ReceiverPartyId, types.NewProposalTask(msg.MsgOption.ProposalId, msg.TaskInfo, msg.CreateAt))

	// Send task to Scheduler to replay sched.
	needReplayScheduleTask := types.NewNeedReplayScheduleTask(msg.MsgOption.ReceiverRole, msg.MsgOption.ReceiverPartyId, msg.TaskInfo)
	t.sendNeedReplayScheduleTask(needReplayScheduleTask)
	replayTaskResult := needReplayScheduleTask.ReceiveResult()

	log.Debugf("Received the reschedule task result from `schedule.ReplaySchedule()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, the result: %s",
		msg.MsgOption.ProposalId.String(), msg.TaskInfo.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId, replayTaskResult.String())

	var vote *twopcpb.PrepareVote

	if nil != replayTaskResult.GetErr() {
		vote = makePrepareVote(
			msg.MsgOption.ProposalId,
			msg.MsgOption.ReceiverRole,
			msg.MsgOption.SenderRole,
			msg.MsgOption.ReceiverPartyId,
			msg.MsgOption.SenderPartyId,
			receiver,
			types.No,
			&types.PrepareVoteResource{},
			timeutils.UnixMsecUint64(),
		)
		log.WithError(replayTaskResult.GetErr()).Warnf("Failed to replay schedule task on onPrepareMsg, replay result has err, will vote `No`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			msg.MsgOption.ProposalId.String(), replayTaskResult.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
	} else {
		vote = makePrepareVote(
			msg.MsgOption.ProposalId,
			msg.MsgOption.ReceiverRole,
			msg.MsgOption.SenderRole,
			msg.MsgOption.ReceiverPartyId,
			msg.MsgOption.SenderPartyId,
			receiver,
			types.Yes,
			types.NewPrepareVoteResource(
				replayTaskResult.GetResource().Id,
				replayTaskResult.GetResource().Ip,
				replayTaskResult.GetResource().Port,
				replayTaskResult.GetResource().PartyId,
			),
			timeutils.UnixMsecUint64(),
		)
		log.Infof("Succeed to replay schedule task on onPrepareMsg, will vote `YES`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			msg.MsgOption.ProposalId.String(), replayTaskResult.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
	}

	// Store current peer own vote for checking whether to vote already
	if consensusSymbol == RemoteConsensusMsg {
		t.state.StorePrepareVote(types.FetchPrepareVote(vote))
	}
	go func() {
		if err := t.sendPrepareVote(pid, receiver, sender, vote); nil != err {
			log.Errorf("failed to call `sendPrepareVote`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, receiver role:{%s}, receiver partyId:{%s}, receiver peerId: {%s}, err: \n%s",
				msg.MsgOption.ProposalId.String(), msg.TaskInfo.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId, msg.MsgOption.SenderRole.String(), msg.MsgOption.SenderPartyId, pid, err)

			// release local resource and clean some data  (on task partner)
			t.stopTaskConsensus("on onPrepareMsg", msg.MsgOption.ProposalId, msg.TaskInfo.GetTaskId(),
				msg.MsgOption.ReceiverRole, msg.MsgOption.SenderRole, receiver, sender, msg.TaskInfo, types.TaskConsensusInterrupt)
			t.removeOrgProposalStateAndTask(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId)
		} else {
			log.Debugf("Succceed to call `sendPrepareVote`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, receiver role:{%s}, receiver partyId:{%s}, receiver peerId: {%s}",
				msg.MsgOption.ProposalId.String(), msg.TaskInfo.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId, msg.MsgOption.SenderRole.String(), msg.MsgOption.SenderPartyId, pid)
		}
	}()
	return nil
}

// (on Publisher)
func (t *Twopc) onPrepareVote(pid peer.ID, prepareVote *types.PrepareVoteWrap, consensusSymbol ConsensusMsgLocationSymbol) error {

	vote := fetchPrepareVote(prepareVote)

	log.Debugf("Received remote prepareVote, remote pid: {%s}, consensusSymbol: {%s}, prepareVote: %s", pid, consensusSymbol.string(), vote.String())

	if t.state.HasNotOrgProposalWithPartyId(vote.MsgOption.ProposalId, vote.MsgOption.ReceiverPartyId) {
		log.Errorf("Failed to check org proposalState whether have been exist on onPrepareVote, but it's not exist, proposalId: {%s}, role: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId)
		return fmt.Errorf("%s onPrepareVote", ctypes.ErrProposalNotFound)
	}
	orgProposalState := t.mustGetOrgProposalState(vote.MsgOption.ProposalId, vote.MsgOption.ReceiverPartyId)

	// The vote in the consensus prepare epoch can be processed only if the current state is the prepare state
	if orgProposalState.IsNotPreparePeriod() {
		log.Errorf("Failed to check org proposalState priod on onPrepareVote, it's not prepare epoch now, proposalId: {%s}, role: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId)
		return ctypes.ErrProposalPrepareVoteTimeout
	}

	// find the task of proposal on proposalTask
	proposalTask, ok := t.state.GetProposalTaskWithPartyId(orgProposalState.GetTaskId(), vote.MsgOption.ReceiverPartyId)
	if !ok {
		log.Errorf("%s on onPrepareVote, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			ctypes.ErrProposalTaskNotFound, vote.MsgOption.ProposalId.String(), orgProposalState.GetTaskId(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId)
		return fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrProposalTaskNotFound, orgProposalState.GetTaskId(), vote.MsgOption.ReceiverRole.String(),
			vote.MsgOption.Owner.GetIdentityId(), vote.MsgOption.ReceiverPartyId)
	}

	sender := fetchOrgByPartyRole(vote.MsgOption.SenderPartyId, vote.MsgOption.SenderRole, proposalTask.Task)
	receiver := fetchOrgByPartyRole(vote.MsgOption.ReceiverPartyId, vote.MsgOption.ReceiverRole, proposalTask.Task)
	if nil == sender || nil == receiver {
		log.Errorf("Failed to check vote.MsgOption sender and receiver of prepareVote on onPrepareVote, some one is empty, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId)
		return ctypes.ErrConsensusMsgInvalid
	}

	identity, err := t.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryIdentity()` on onPrepareVote, some one is empty, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId)

		// Send consensus result to interrupt consensus epoch and clean some data (on task sender)
		t.stopTaskConsensus("query local identity failed on onPrepareVote", vote.MsgOption.ProposalId, proposalTask.GetTaskId(),
			apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Sender, receiver, receiver, proposalTask.GetTask(), types.TaskConsensusInterrupt)
		t.removeOrgProposalStateAndTask(vote.MsgOption.ProposalId, vote.MsgOption.ReceiverPartyId)
		return fmt.Errorf("query local identity failed, %s", err)
	}
	// verify the receiver is myself ?
	if identity.GetIdentityId() != receiver.GetIdentityId() {
		log.Errorf("Failed to verify receiver identityId of prepareVote, receiver is not me, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId)
		return ctypes.ErrConsensusMsgInvalid
	}

	// Voter <the vote sender> voted repeatedly
	if t.state.HasPrepareVoting(vote.MsgOption.ProposalId, sender) {
		log.Errorf("%s on onPrepareVote, they are not same, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, vote sender role: {%s}, vote sender partyId: {%s}",
			ctypes.ErrPrepareVoteRepeatedly, vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId,
			vote.MsgOption.SenderRole.String(), vote.MsgOption.SenderPartyId)
		return fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrPrepareVoteRepeatedly, proposalTask.GetTaskData().GetTaskId(), vote.MsgOption.ReceiverRole.String(),
			vote.MsgOption.Owner.GetIdentityId(), vote.MsgOption.ReceiverPartyId)
	}

	identityValid, err := t.verifyPrepareVoteRoleIsTaskPartner(sender.GetIdentityId(), sender.GetPartyId(), vote.MsgOption.SenderRole, proposalTask.Task)
	if nil != err {
		log.WithError(err).Errorf("Failed to call `verifyPrepareVoteRoleIsTaskPartner()` verify prepare vote role on onPrepareVote, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId)
		return err
	}
	if !identityValid {
		log.Errorf("The prepare vote role is not include task partners on onPrepareVote, they are not same, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId)
		return fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrProposalPrepareVoteOwnerInvalid, proposalTask.GetTaskData().GetTaskId(), vote.MsgOption.SenderRole.String(),
			sender.GetIdentityId(), sender.GetPartyId())
	}

	// verify resource of `YES` vote
	if vote.VoteOption == types.Yes && vote.PeerInfoEmpty() {
		log.Errorf("%s on onPrepareVote, they are not same, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			ctypes.ErrProposalPrepareVoteResourceInvalid, vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(),
			vote.MsgOption.ReceiverPartyId)
		return fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrProposalPrepareVoteResourceInvalid, proposalTask.GetTaskData().GetTaskId(), vote.MsgOption.SenderRole.String(),
			sender.GetIdentityId(), sender.GetPartyId())
	}

	// Store vote
	t.state.StorePrepareVote(vote)

	totalNeedVoteCount := uint32(len(proposalTask.Task.GetTaskData().GetDataSuppliers()) +
		len(proposalTask.Task.GetTaskData().GetPowerSuppliers()) +
		len(proposalTask.Task.GetTaskData().GetReceivers()))

	yesVoteCount := t.state.GetTaskPrepareYesVoteCount(vote.MsgOption.ProposalId)
	totalVotedCount := t.state.GetTaskPrepareTotalVoteCount(vote.MsgOption.ProposalId)

	if totalNeedVoteCount == totalVotedCount {

		now := timeutils.UnixMsecUint64()

		// send confirm msg by option `start` to other remote peers,
		// (announce other peer to continue consensus epoch to confirm epoch)
		// and change proposal state from prepare epoch to confirm epoch
		if totalNeedVoteCount == yesVoteCount {

			// change state from prepare epoch to confirm epoch
			t.state.ChangeToConfirm(vote.MsgOption.ProposalId, vote.MsgOption.ReceiverPartyId, now)

			// store confirm peers resource info
			peers := t.makeConfirmTaskPeerDesc(vote.MsgOption.ProposalId)
			t.storeConfirmTaskPeerInfo(vote.MsgOption.ProposalId, peers)

			go func() {

				log.Infof("PrepareVoting succeed on consensus prepare epoch, the `Yes` vote count has enough, will send `Start` confirm msg, the `Yes` vote count: {%d}, need total count: {%d}, with proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
					yesVoteCount, totalNeedVoteCount, vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(),
					vote.MsgOption.ReceiverPartyId)
				if err := t.sendConfirmMsg(vote.MsgOption.ProposalId, proposalTask.Task, peers, types.TwopcMsgStart, now); nil != err {
					log.Errorf("Failed to call `sendConfirmMsg` with `start` consensus prepare epoch on `onPrepareVote`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, err: \n%s",
						vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId, err)
					// Send consensus result to interrupt consensus epoch and clean some data (on task sender)
					t.stopTaskConsensus("failed to call `sendConfirmMsg`", vote.MsgOption.ProposalId, proposalTask.GetTaskId(),
						apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Sender, receiver, receiver, proposalTask.GetTask(), types.TaskConsensusInterrupt)
					t.removeOrgProposalStateAndTask(vote.MsgOption.ProposalId, vote.MsgOption.ReceiverPartyId)
				}
			}()

		} else {

			// send confirm msg by option `stop` to other remote peers,
			// (announce other peer to interrupt consensus epoch)
			// and remove local cache (task/proposal state/prepare vote) about proposal and task
			go func() {

				log.Infof("PrepareVoting failed on consensus prepare epoch, the `Yes` vote count is no enough, will send `Stop` confirm msg, the `Yes` vote count: {%d}, need total count: {%d}, with proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
					yesVoteCount, totalNeedVoteCount, vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(),
					vote.MsgOption.ReceiverPartyId)

				if err := t.sendConfirmMsg(vote.MsgOption.ProposalId, proposalTask.Task, t.makeEmptyConfirmTaskPeerDesc(), types.TwopcMsgStop, now); nil != err {
					log.Errorf("Failed to call `sendConfirmMsg` with `stop` consensus prepare epoch on `onPrepareVote`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, err: \n%s",
						vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId, err)
				}
				// Send consensus result to interrupt consensus epoch and clean some data (on task sender)
				t.stopTaskConsensus("the prepareMsg voting result was not passed", vote.MsgOption.ProposalId, proposalTask.GetTaskId(),
					apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Sender, receiver, receiver, proposalTask.GetTask(), types.TaskConsensusInterrupt)
				t.removeOrgProposalStateAndTask(vote.MsgOption.ProposalId, vote.MsgOption.ReceiverPartyId)
			}()
		}
	}
	return nil
}

// (on Subscriber)
func (t *Twopc) onConfirmMsg(pid peer.ID, confirmMsg *types.ConfirmMsgWrap, consensusSymbol ConsensusMsgLocationSymbol) error {

	msg := fetchConfirmMsg(confirmMsg)

	log.Debugf("Received remote confirmMsg, remote pid: {%s}, consensusSymbol: {%s}, confirmMsg: %s", pid, consensusSymbol.string(), msg.String())

	if t.state.HasNotOrgProposalWithPartyId(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId) {
		log.Errorf("Failed to check org proposalState whether have been exist on onConfirmMsg, but it's not exist, proposalId: {%s}, role: {%s}, partyId: {%s}",
			msg.MsgOption.ProposalId.String(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
		return fmt.Errorf("%s onConfirmMsg", ctypes.ErrProposalNotFound)
	}

	orgProposalState := t.mustGetOrgProposalState(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId)

	// The vote in the consensus prepare epoch or confirm epoch can be processed just if the current state is the prepare state or confirm state.
	if orgProposalState.IsCommitPeriod() {
		log.Errorf("Failed to check org proposalState priod on onConfirmMsg, it's commit epoch now, proposalId: {%s}, role: {%s}, partyId: {%s}",
			msg.MsgOption.ProposalId.String(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
		return ctypes.ErrProposalConfirmMsgTimeout
	}

	// find the task of proposal on proposalTask
	proposalTask, ok := t.state.GetProposalTaskWithPartyId(orgProposalState.GetTaskId(), msg.MsgOption.ReceiverPartyId)
	if !ok {
		log.Errorf("%s on onConfirmMsg, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			ctypes.ErrProposalTaskNotFound, msg.MsgOption.ProposalId.String(), orgProposalState.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
		return fmt.Errorf("%s, on the confirm msg [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrProposalTaskNotFound, proposalTask.GetTaskData().GetTaskId(), msg.MsgOption.ReceiverRole.String(),
			msg.MsgOption.Owner.GetIdentityId(), msg.MsgOption.ReceiverPartyId)
	}

	sender := fetchOrgByPartyRole(msg.MsgOption.SenderPartyId, msg.MsgOption.SenderRole, proposalTask.Task)
	receiver := fetchOrgByPartyRole(msg.MsgOption.ReceiverPartyId, msg.MsgOption.ReceiverRole, proposalTask.Task)
	if nil == sender || nil == receiver {
		log.Errorf("Failed to check msg.MsgOption sender and receiver of confirmMsg on onConfirmMsg, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			msg.MsgOption.ProposalId.String(), orgProposalState.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
		return ctypes.ErrConsensusMsgInvalid
	}

	identity, err := t.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryIdentity()` on onConfirmMsg, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			msg.MsgOption.ProposalId.String(), orgProposalState.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
		// release local resource and clean some data  (on task partner)
		t.stopTaskConsensus("on onConfirmMsg", msg.MsgOption.ProposalId, proposalTask.GetTask().GetTaskId(),
			msg.MsgOption.ReceiverRole, msg.MsgOption.SenderRole, receiver, sender, proposalTask.GetTask(), types.TaskConsensusInterrupt)
		t.removeOrgProposalStateAndTask(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId)
		return fmt.Errorf("query local identity failed, %s", err)
	}

	// verify the receiver is myself ?
	if identity.GetIdentityId() != receiver.GetIdentityId() {

		log.Errorf("Failed to verify receiver identityId of confirmMsg, receiver is not me, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			msg.MsgOption.ProposalId.String(), orgProposalState.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
		return ctypes.ErrConsensusMsgInvalid
	}

	org := &apicommonpb.TaskOrganization{
		PartyId:    msg.MsgOption.ReceiverPartyId,
		NodeName:   identity.GetNodeName(),
		NodeId:     identity.GetNodeId(),
		IdentityId: identity.GetIdentityId(),
	}

	// If you have already voted then we will not vote again.
	// Cause the local message will only call the local function once,
	// and the remote message needs to prevent receiving the repeated forwarded consensus message.
	if consensusSymbol == RemoteConsensusMsg && t.state.HasConfirmVoting(msg.MsgOption.ProposalId, org) {
		log.Errorf("Failed to check remote peer confirm vote wether voting on onConfirmMsg, it's voting alreay, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, confirmMsgOption: {%s}",
			msg.MsgOption.ProposalId.String(), orgProposalState.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId, msg.ConfirmOption.String())
		return ctypes.ErrConfirmVotehadVoted
	}

	// check msg confirm option value is `start` or `stop` ?
	if msg.ConfirmOption == types.TwopcMsgStop || msg.ConfirmOption == types.TwopcMsgUnknown {
		log.Errorf("Failed to verify confirmMsgOption of confirmMsg on onConfirmMsg, confirmMsgOption is not `Start`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, confirmMsgOption: {%s}",
			msg.MsgOption.ProposalId.String(), orgProposalState.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId, msg.ConfirmOption.String())
		// release local resource and clean some data  (on task partner)
		t.stopTaskConsensus("on onConfirmMsg", msg.MsgOption.ProposalId, proposalTask.GetTask().GetTaskId(),
			msg.MsgOption.ReceiverRole, msg.MsgOption.SenderRole, receiver, sender, proposalTask.GetTask(), types.TaskConsensusInterrupt)
		t.removeOrgProposalStateAndTask(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId)
		return ctypes.ErrConsensusMsgInvalid
	}

	var vote *twopcpb.ConfirmVote

	// verify peers resources
	if msg.PeersEmpty() {

		vote = makeConfirmVote(
			proposalTask.ProposalId,
			msg.MsgOption.ReceiverRole,
			msg.MsgOption.SenderRole,
			msg.MsgOption.ReceiverPartyId,
			msg.MsgOption.SenderPartyId,
			receiver,
			types.No,
			timeutils.UnixMsecUint64(),
		)

		log.Warnf("Failed to verify peers resources of confirmMsg on onConfirmMsg, the peerDesc reources is empty, will vote `No`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, confirmMsgOption: {%s}",
			msg.MsgOption.ProposalId.String(), orgProposalState.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId, msg.ConfirmOption.String())

	} else {
		// store confirm peers resource info
		t.storeConfirmTaskPeerInfo(msg.MsgOption.ProposalId, msg.Peers)
		//if consensusSymbol == RemoteConsensusMsg {
		//	t.storeConfirmTaskPeerInfo(msg.MsgOption.ProposalId, msg.Peers)
		//}
		vote = makeConfirmVote(
			proposalTask.ProposalId,
			msg.MsgOption.ReceiverRole,
			msg.MsgOption.SenderRole,
			msg.MsgOption.ReceiverPartyId,
			msg.MsgOption.SenderPartyId,
			receiver,
			types.Yes,
			timeutils.UnixMsecUint64(),
		)

		log.Infof("Succeed to verify peers resources of confirmMsg on onConfirmMsg, will vote `Yes`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, confirmMsgOption: {%s}",
			msg.MsgOption.ProposalId.String(), orgProposalState.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId, msg.ConfirmOption.String())

	}

	// Store current peer own vote for checking whether to vote already
	if consensusSymbol == RemoteConsensusMsg {
		t.state.StoreConfirmVote(types.FetchConfirmVote(vote))
	}

	// change state from prepare epoch to confirm epoch
	t.state.ChangeToConfirm(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId, msg.CreateAt)
	go func() {

		if err := t.sendConfirmVote(pid, receiver, sender, vote); nil != err {
			log.Errorf("failed to call `sendConfirmVote`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, receiver role:{%s}, receiver partyId:{%s}, receiver peerId: {%s}, \n%s",
				msg.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId, msg.MsgOption.SenderRole.String(), msg.MsgOption.SenderPartyId, pid, err)
			// release local resource and clean some data  (on task partner)
			t.stopTaskConsensus("on onConfirmMsg", msg.MsgOption.ProposalId, proposalTask.GetTask().GetTaskId(),
				msg.MsgOption.ReceiverRole, msg.MsgOption.SenderRole, receiver, sender, proposalTask.GetTask(), types.TaskConsensusInterrupt)
			t.removeOrgProposalStateAndTask(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId)
		} else {

			// In any case, as long as voting 'No', Need to clean the local cache
			if msg.PeersEmpty() {
				// release local resource and clean some data  (on task partner)
				t.stopTaskConsensus("on onConfirmMsg", msg.MsgOption.ProposalId, proposalTask.GetTask().GetTaskId(),
					msg.MsgOption.ReceiverRole, msg.MsgOption.SenderRole, receiver, sender, proposalTask.GetTask(), types.TaskConsensusInterrupt)
				t.removeOrgProposalStateAndTask(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId)
			}
			log.Debugf("Succceed to call `sendConfirmVote`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, receiver role:{%s}, receiver partyId:{%s}, receiver peerId: {%s}",
				msg.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId, msg.MsgOption.SenderRole.String(), msg.MsgOption.SenderPartyId, pid)

		}
	}()

	return nil
}

// (on Publisher)
func (t *Twopc) onConfirmVote(pid peer.ID, confirmVote *types.ConfirmVoteWrap, consensusSymbol ConsensusMsgLocationSymbol) error {

	vote := fetchConfirmVote(confirmVote)

	log.Debugf("Received remote confirmVote, remote pid: {%s}, consensusSymbol: {%s}, comfirmVote: %s", pid, consensusSymbol.string(), vote.String())

	if t.state.HasNotOrgProposalWithPartyId(vote.MsgOption.ProposalId, vote.MsgOption.ReceiverPartyId) {
		log.Errorf("Failed to check org proposalState whether have been exist on onConfirmVote, but it's not exist, proposalId: {%s}, role: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId)
		return fmt.Errorf("%s onConfirmVote", ctypes.ErrProposalNotFound)
	}
	orgProposalState := t.mustGetOrgProposalState(vote.MsgOption.ProposalId, vote.MsgOption.ReceiverPartyId)

	// The vote in the consensus confirm epoch can be processed only if the current state is the confirm state
	if orgProposalState.IsPreparePeriod() {
		log.Errorf("Failed to check org proposalState priod on onConfirmVote, it's not confirm epoch and is prepare epoch now, proposalId: {%s}, role: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId)
		return ctypes.ErrProposalConfirmVoteFuture
	}
	if orgProposalState.IsCommitPeriod() {
		log.Errorf("Failed to check org proposalState priod on onConfirmVote, it's not confirm epoch and is commit epoch now, proposalId: {%s}, role: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId)
		return ctypes.ErrProposalPrepareVoteTimeout
	}

	// find the task of proposal on proposalTask
	proposalTask, ok := t.state.GetProposalTaskWithPartyId(orgProposalState.GetTaskId(), vote.MsgOption.ReceiverPartyId)
	if !ok {
		log.Errorf("%s on onConfirmVote, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			ctypes.ErrProposalTaskNotFound, vote.MsgOption.ProposalId.String(), orgProposalState.GetTaskId(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId)
		return fmt.Errorf("%s, on the confirm vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrProposalTaskNotFound, proposalTask.GetTaskData().GetTaskId(), vote.MsgOption.ReceiverRole.String(),
			vote.MsgOption.Owner.GetIdentityId(), vote.MsgOption.ReceiverPartyId)
	}

	sender := fetchOrgByPartyRole(vote.MsgOption.SenderPartyId, vote.MsgOption.SenderRole, proposalTask.Task)
	receiver := fetchOrgByPartyRole(vote.MsgOption.ReceiverPartyId, vote.MsgOption.ReceiverRole, proposalTask.Task)
	if nil == sender || nil == receiver {
		log.Errorf("Failed to check vote.MsgOption sender and receiver of confirmVote on onConfirmVote, some one is empty, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId)
		return ctypes.ErrConsensusMsgInvalid
	}

	identity, err := t.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryIdentity()` on onConfirmVote, some one is empty, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId)

		// Send consensus result to interrupt consensus epoch and clean some data (on task sender)
		t.stopTaskConsensus("query local identity failed on onConfirmVote", vote.MsgOption.ProposalId, proposalTask.GetTaskId(),
			apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Sender, receiver, receiver, proposalTask.GetTask(), types.TaskConsensusInterrupt)
		t.removeOrgProposalStateAndTask(vote.MsgOption.ProposalId, vote.MsgOption.ReceiverPartyId)
		return fmt.Errorf("query local identity failed, %s", err)
	}
	// verify the receiver is myself ?
	if identity.GetIdentityId() != receiver.GetIdentityId() {
		log.Errorf("Failed to verify receiver identityId of confirmVote, receiver is not me, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId)
		return ctypes.ErrConsensusMsgInvalid
	}

	// Voter <the vote sender> voted repeatedly
	if t.state.HasConfirmVoting(vote.MsgOption.ProposalId, sender) {
		log.Errorf("%s on onConfirmVote, they are not same, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, vote sender role: {%s}, vote sender partyId: {%s}",
			ctypes.ErrConfirmVoteRepeatedly, vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId,
			vote.MsgOption.SenderRole.String(), vote.MsgOption.SenderPartyId)
		return ctypes.ErrConfirmVoteRepeatedly
	}

	identityValid, err := t.verifyConfirmVoteRoleIsTaskPartner(sender.GetIdentityId(), sender.GetPartyId(), vote.MsgOption.SenderRole, proposalTask.Task)
	if nil != err {
		log.WithError(err).Errorf("Failed to call `verifyConfirmVoteRoleIsTaskPartner()` verify confirm vote role on onConfirmVote, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId)
		return err
	}
	if !identityValid {
		log.Errorf("The confirm vote role is not include task partners on onConfirmVote, they are not same, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId)
		return fmt.Errorf("%s, on the confirm vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrProposalConfirmVoteVoteOwnerInvalid, proposalTask.GetTaskData().GetTaskId(), vote.MsgOption.SenderRole.String(), sender.GetIdentityId(), sender.GetPartyId())
	}

	// Store vote
	t.state.StoreConfirmVote(vote)

	totalNeedVoteCount := uint32(len(proposalTask.Task.GetTaskData().GetDataSuppliers()) +
		len(proposalTask.Task.GetTaskData().GetPowerSuppliers()) +
		len(proposalTask.Task.GetTaskData().GetReceivers()))

	yesVoteCount := t.state.GetTaskConfirmYesVoteCount(vote.MsgOption.ProposalId)
	totalVotedCount := t.state.GetTaskConfirmTotalVoteCount(vote.MsgOption.ProposalId)

	if totalNeedVoteCount == totalVotedCount {

		now := timeutils.UnixMsecUint64()

		// send commit msg by option `start` to other remote peers,
		// (announce other peer to continue consensus epoch to commit epoch)
		// and change proposal state from confirm epoch to commit epoch
		if totalNeedVoteCount == yesVoteCount {

			// change state from confirm epoch to commit epoch
			t.state.ChangeToCommit(vote.MsgOption.ProposalId, vote.MsgOption.ReceiverPartyId, now)

			go func() {

				log.Debugf("ConfirmVoting succeed on consensus confirm epoch, the `Yes` vote count has enough, will send `Start` commit msg, the `Yes` vote count: {%d}, need total count: {%d}, with proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
					yesVoteCount, totalNeedVoteCount, vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(),
					vote.MsgOption.ReceiverPartyId)

				if err := t.sendCommitMsg(vote.MsgOption.ProposalId, proposalTask.Task, types.TwopcMsgStart, now); nil != err {
					log.Errorf("Failed to call `sendCommitMsg` with `start` consensus confirm epoch on `onConfirmVote`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, err: \n%s",
						vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId, err)
					// Send consensus result (on task sender)
					t.stopTaskConsensus("failed to call `sendCommitMsg`", vote.MsgOption.ProposalId, proposalTask.GetTaskId(),
						apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Sender, receiver, receiver, proposalTask.GetTask(), types.TaskConsensusInterrupt)
				} else {
					// Send consensus result (on task sender)
					t.replyTaskConsensusResult(types.NewTaskConsResult(proposalTask.GetTaskId(), types.TaskConsensusFinished, nil))
				}
				// Finally, whether the commitmsg is sent successfully or not, the local cache needs to be cleared
				t.removeOrgProposalStateAndTask(vote.MsgOption.ProposalId, vote.MsgOption.ReceiverPartyId)

			}()

		} else {

			// send commit msg by option `stop` to other remote peers,
			// (announce other peer to interrupt consensus epoch)
			// and remove local cache (task/proposal state/prepare vote/confirm vote/peerDesc) about proposal and task
			go func() {

				log.Debugf("ConfirmVoting failed on consensus confirm epoch, the `Yes` vote count is no enough, will send `Stop` commit msg, the `Yes` vote count: {%d}, need total count: {%d}, with proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
					yesVoteCount, totalNeedVoteCount, vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(),
					vote.MsgOption.ReceiverPartyId)

				if err := t.sendCommitMsg(vote.MsgOption.ProposalId, proposalTask.Task, types.TwopcMsgStop, now); nil != err {
					log.Errorf("Failed to call `sendCommitMsg` with `stop` consensus confirm epoch on `onConfirmVote`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, err: \n%s",
						vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), vote.MsgOption.ReceiverRole.String(), vote.MsgOption.ReceiverPartyId, err)
				}
				// Send consensus result to interrupt consensus epoch and clean some data (on task sender)
				t.stopTaskConsensus("the cofirmMsg voting result was not passed", vote.MsgOption.ProposalId, proposalTask.GetTaskId(),
					apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Sender, receiver, receiver, proposalTask.GetTask(), types.TaskConsensusInterrupt)
				t.removeOrgProposalStateAndTask(vote.MsgOption.ProposalId, vote.MsgOption.ReceiverPartyId)
			}()
		}
	}
	return nil
}

// (on Subscriber)
func (t *Twopc) onCommitMsg(pid peer.ID, cimmitMsg *types.CommitMsgWrap, consensusSymbol ConsensusMsgLocationSymbol) error {

	msg := fetchCommitMsg(cimmitMsg)

	log.Debugf("Received remote commitMsg, remote pid: {%s}, consensusSymbol: {%s}, commitMsg: %s", pid, consensusSymbol.string(), msg.String())

	if t.state.HasNotOrgProposalWithPartyId(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId) {
		log.Errorf("Failed to check org proposalState whether have been exist on onCommitMsg, but it's not exist, proposalId: {%s}, role: {%s}, partyId: {%s}",
			msg.MsgOption.ProposalId.String(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
		return fmt.Errorf("%s onCommitMsg", ctypes.ErrProposalNotFound)
	}

	orgProposalState := t.mustGetOrgProposalState(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId)

	// The vote in the consensus confirm epoch or commit epoch can be processed just if the current state is the confirm state or commit state
	if orgProposalState.IsPreparePeriod() {
		log.Errorf("Failed to check org proposalState priod on onCommitMsg, it's not commit epoch and is prepare epoch now, proposalId: {%s}, role: {%s}, partyId: {%s}",
			msg.MsgOption.ProposalId.String(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
		return ctypes.ErrProposalCommitMsgFuture
	}
	if orgProposalState.IsFinishedPeriod() {
		log.Errorf("Failed to check org proposalState priod on onCommitMsg, it's not commit epoch and is finished epoch now, proposalId: {%s}, role: {%s}, partyId: {%s}",
			msg.MsgOption.ProposalId.String(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
		return ctypes.ErrProposalCommitMsgTimeout
	}

	// find the task of proposal on proposalTask
	proposalTask, ok := t.state.GetProposalTaskWithPartyId(orgProposalState.GetTaskId(), msg.MsgOption.ReceiverPartyId)
	if !ok {
		log.Errorf("%s on onCommitMsg, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			ctypes.ErrProposalTaskNotFound, msg.MsgOption.ProposalId.String(), orgProposalState.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
		return fmt.Errorf("%s, on the commit msg [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrProposalTaskNotFound, proposalTask.GetTaskData().GetTaskId(), msg.MsgOption.ReceiverRole.String(),
			msg.MsgOption.Owner.GetIdentityId(), msg.MsgOption.ReceiverPartyId)
	}

	sender := fetchOrgByPartyRole(msg.MsgOption.SenderPartyId, msg.MsgOption.SenderRole, proposalTask.Task)
	receiver := fetchOrgByPartyRole(msg.MsgOption.ReceiverPartyId, msg.MsgOption.ReceiverRole, proposalTask.Task)
	if nil == sender || nil == receiver {
		log.Errorf("Failed to check msg.MsgOption sender and receiver of commitMsg on onCommitMsg, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			msg.MsgOption.ProposalId.String(), orgProposalState.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
		return ctypes.ErrConsensusMsgInvalid
	}

	identity, err := t.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryIdentity()` on onCommitMsg, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			msg.MsgOption.ProposalId.String(), orgProposalState.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
		// release local resource and clean some data  (on task partner)
		t.stopTaskConsensus("on onCommitMsg", msg.MsgOption.ProposalId, proposalTask.GetTask().GetTaskId(),
			msg.MsgOption.ReceiverRole, msg.MsgOption.SenderRole, receiver, sender, proposalTask.GetTask(), types.TaskConsensusInterrupt)
		t.removeOrgProposalStateAndTask(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId)
		return fmt.Errorf("query local identity failed, %s", err)
	}

	// verify the receiver is myself ?
	if identity.GetIdentityId() != receiver.GetIdentityId() {
		log.Errorf("Failed to verify receiver identityId of commitMsg, receiver is not me, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			msg.MsgOption.ProposalId.String(), orgProposalState.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId)
		return ctypes.ErrConsensusMsgInvalid
	}

	// check msg commit option value is `start` or `stop` ?
	if msg.CommitOption == types.TwopcMsgStop || msg.CommitOption == types.TwopcMsgUnknown {
		log.Errorf("Failed to verify commitMsgOption of commitMsg on onCommitMsg, commitMsgOption is not `Start`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, confirmMsgOption: {%s}",
			msg.MsgOption.ProposalId.String(), orgProposalState.GetTaskId(), msg.MsgOption.ReceiverRole.String(), msg.MsgOption.ReceiverPartyId, msg.CommitOption.String())
		// release local resource and clean some data  (on task partner)
		t.stopTaskConsensus("on onCommitMsg", msg.MsgOption.ProposalId, proposalTask.GetTask().GetTaskId(),
			msg.MsgOption.ReceiverRole, msg.MsgOption.SenderRole, receiver, sender, proposalTask.GetTask(), types.TaskConsensusInterrupt)
		t.removeOrgProposalStateAndTask(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId)
		return ctypes.ErrConsensusMsgInvalid
	}

	// change state from confirm epoch to commit epoch
	t.state.ChangeToCommit(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId, msg.CreateAt)

	go func() {

		// store succeed consensus event for partyId
		t.resourceMng.GetDB().StoreTaskEvent(&libtypes.TaskEvent{
			Type:       ev.TaskFinishedConsensus.Type,
			TaskId:     orgProposalState.GetTaskId(),
			IdentityId: receiver.GetIdentityId(),
			PartyId:    receiver.GetPartyId(),
			Content:    fmt.Sprintf("finished consensus succeed."),
			CreateAt:   timeutils.UnixMsecUint64(),
		})
		// If receiving `CommitMsg` is successful,
		// we will forward `schedTask` to `taskManager` to send it to `Fighter` to execute the task.
		t.driveTask(pid, msg.MsgOption.ProposalId, msg.MsgOption.ReceiverRole, receiver, msg.MsgOption.SenderRole, sender, proposalTask.Task)
		t.removeOrgProposalStateAndTask(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId)
	}()

	// Finally, it is left 'taskmanager' to call 'releaselocalresourcewithtask()' to release local resources after handle `driveTask()`.
	// No more processing here.
	return nil
}

func (t *Twopc) onTerminateTaskConsensus(pid peer.ID, msg *types.InterruptMsgWrap) error {

	msgOption := types.FetchMsgOption(msg.MsgOption)
	log.Infof("Start interrupt task consensus, taskId: {%s}, partyId: {%s}", msg.GetTaskId(), msgOption.ReceiverPartyId)

	// find the task of proposal on proposalTask
	proposalTask, ok := t.state.GetProposalTaskWithPartyId(msg.GetTaskId(), msgOption.ReceiverPartyId)
	if !ok {
		log.Errorf("%s on onTerminateTaskConsensus, taskId: {%s}, partyId: {%s}", ctypes.ErrProposalTaskNotFound, msg.GetTaskId(), msgOption.ReceiverPartyId)
		return fmt.Errorf("%s, on the interrupt consensus [taskId: %s, partyId: %s]",
			ctypes.ErrProposalTaskNotFound, msg.GetTaskId(), msgOption.ReceiverPartyId)
	}

	if t.state.HasNotOrgProposalWithPartyId(proposalTask.GetProposalId(), msgOption.ReceiverPartyId) {
		log.Errorf("Failed to check org proposalState whether have been exist on onTerminateTaskConsensus, but it's not exist, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
			proposalTask.GetProposalId().String(), msg.GetTaskId(), msgOption.ReceiverPartyId)
		return fmt.Errorf("%s, on the interrupt consensus", ctypes.ErrProposalNotFound)
	}

	sender := fetchOrgByPartyRole(msgOption.SenderPartyId, msgOption.SenderRole, proposalTask.Task)
	receiver := fetchOrgByPartyRole(msgOption.ReceiverPartyId, msgOption.ReceiverRole, proposalTask.Task)
	if nil == sender || nil == receiver {
		log.Errorf("Failed to check msg.MsgOption sender and receiver of interruptMsg on onTerminateTaskConsensus, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
			proposalTask.GetProposalId().String(), msg.GetTaskId(), msgOption.ReceiverRole, msgOption.ReceiverPartyId)
		return ctypes.ErrConsensusMsgInvalid
	}

	orgProposalState := t.mustGetOrgProposalState(proposalTask.GetProposalId(), msgOption.ReceiverPartyId)
	switch orgProposalState.CurrPeriodNum() {
	case ctypes.PeriodPrepare:
		// remove `proposal state` and `task cache` AND inerrupt consensus with sender OR release local locked resource with partner
		t.stopTaskConsensus("interrupt consensus with terminate task while prepare epoch", proposalTask.GetProposalId(), msg.GetTaskId(),
			msgOption.ReceiverRole, msgOption.SenderRole, receiver, sender, proposalTask.GetTask(), types.TaskTerminate)
		t.removeOrgProposalStateAndTask(proposalTask.GetProposalId(), proposalTask.GetTaskId())
	case ctypes.PeriodConfirm:
		// remove `proposal state` and `task cache` AND inerrupt consensus with sender OR release local locked resource with partner
		t.stopTaskConsensus("interrupt consensus with terminate task while confirm epoch", proposalTask.GetProposalId(), msg.GetTaskId(),
			msgOption.ReceiverRole, msgOption.SenderRole, receiver, sender, proposalTask.GetTask(), types.TaskTerminate)
		t.removeOrgProposalStateAndTask(proposalTask.GetProposalId(), proposalTask.GetTaskId())
	case ctypes.PeriodCommit, ctypes.PeriodFinished:
		// need send terminate msg with task manager
		// so do nothing here
	default:
		log.Errorf("unknown org proposalState priod on onTerminateTaskConsensus, proposalId: {%s}, taskId: {%s}, partyId: {%s}, peroid: {%s}",
			proposalTask.GetProposalId().String(), msg.GetTaskId(), msgOption.ReceiverPartyId, orgProposalState.GetPeriodStr())
		return fmt.Errorf("unknown org proposalState priod, on the interrupt consensus")
	}
	return nil
}
