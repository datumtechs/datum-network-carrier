package twopc

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/rlputil"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	"strings"
	"sync"
	"time"
)

const (
	//defaultCleanExpireProposalInterval  = 30 * time.Millisecond
	defaultRefreshProposalStateInternal = 300 * time.Millisecond
)

type TwoPC struct {
	config  *Config
	p2p     p2p.P2P
	peerSet *ctypes.PeerSet
	state   *state
	//resourceMng.GetDB()  iface.ForResourceDB
	resourceMng *resource.Manager

	// fetch tasks scheduled from `Scheduler`
	//needConsensusTaskCh chan *types.NeedConsensusTask
	// send remote task to `Scheduler` to replay
	needReplayScheduleTaskCh chan *types.NeedReplayScheduleTask
	// send has was consensus remote tasks to taskManager
	needExecuteTaskCh chan *types.NeedExecuteTask
	asyncCallCh       chan func()
	quit              chan struct{}
	proposalTaskCache map[string]*types.ProposalTask // (taskId -> task)

	taskResultBusCh  chan *types.TaskConsResult
	taskResultChSet  map[string]chan<- *types.TaskConsResult
	taskResultLock   sync.Mutex
	proposalTaskLock sync.RWMutex

	Errs []error
}

func New(
	conf *Config,
	resourceMng *resource.Manager,
	p2p p2p.P2P,
	//needConsensusTaskCh chan *types.NeedConsensusTask,
	needReplayScheduleTaskCh chan *types.NeedReplayScheduleTask,
	needExecuteTaskCh chan *types.NeedExecuteTask,
) *TwoPC {
	return &TwoPC{
		config:                   conf,
		p2p:                      p2p,
		peerSet:                  ctypes.NewPeerSet(10), // TODO 暂时写死的
		state:                    newState(),
		resourceMng:              resourceMng,
		needReplayScheduleTaskCh: needReplayScheduleTaskCh,
		needExecuteTaskCh:        needExecuteTaskCh,
		asyncCallCh:              make(chan func(), conf.PeerMsgQueueSize),
		quit:                     make(chan struct{}),
		proposalTaskCache:        make(map[string]*types.ProposalTask),
		taskResultBusCh: make(chan *types.TaskConsResult, 100),
		taskResultChSet: make(map[string]chan<- *types.TaskConsResult, 100),
		Errs:            make([]error, 0),
	}
}

func (t *TwoPC) Start() error {
	go t.loop()
	log.Info("Started 2pc consensus engine ...")
	return nil
}
func (t *TwoPC) Close() error {
	close(t.quit)
	return nil
}
func (t *TwoPC) loop() {
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

func (t *TwoPC) ValidateConsensusMsg(pid peer.ID, msg types.ConsensusMsg) error {
	if nil == msg {
		return fmt.Errorf("Failed to validate 2pc consensus msg, the msg is nil")
	}
	//switch msg := msg.(type) {
	//case *types.PrepareMsgWrap:
	//	return t.validatePrepareMsg(pid, msg)
	//case *types.PrepareVoteWrap:
	//	return t.validatePrepareVote(pid, msg)
	//case *types.ConfirmMsgWrap:
	//	return t.validateConfirmMsg(pid, msg)
	//case *types.ConfirmVoteWrap:
	//	return t.validateConfirmVote(pid, msg)
	//case *types.CommitMsgWrap:
	//	return t.validateCommitMsg(pid, msg)
	////case *types.TaskResultMsgWrap:
	////	return t.validateTaskResultMsg(pid, msg)
	//default:
	//	return fmt.Errorf("TaskRoleUnknown the 2pc msg type")
	//}
	return nil
}

func (t *TwoPC) OnConsensusMsg(pid peer.ID, msg types.ConsensusMsg) error {

	switch msg := msg.(type) {
	case *types.PrepareMsgWrap:
		return t.onPrepareMsg(pid, msg)
	case *types.PrepareVoteWrap:
		return t.onPrepareVote(pid, msg)
	case *types.ConfirmMsgWrap:
		return t.onConfirmMsg(pid, msg)
	case *types.ConfirmVoteWrap:
		return t.onConfirmVote(pid, msg)
	case *types.CommitMsgWrap:
		return t.onCommitMsg(pid, msg)
	//case *types.TaskResultMsgWrap:
	//	return t.onTaskResultMsg(pid, msg)
	default:
		return fmt.Errorf("TaskRoleUnknown the 2pc msg type")

	}
}

func (t *TwoPC) OnError() error {
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

func (t *TwoPC) OnPrepare(task *types.Task) error {

	return nil
}
func (t *TwoPC) OnHandle(task *types.Task, result chan<- *types.TaskConsResult) error {

	t.addTaskResultCh(task.GetTaskId(), result)

	if t.isNotProposalTask(task.GetTaskId()) {
		t.replyTaskConsensusResult(types.NewTaskConsResult(task.GetTaskId(), types.TaskConsensusInterrupt, ctypes.ErrPrososalTaskIsProcessed))
		return ctypes.ErrPrososalTaskIsProcessed
	}
	now := uint64(timeutils.UnixMsec())
	proposalId := rlputil.RlpHash([]interface{}{
		t.config.Option.NodeID,
		now,
		task.GetTaskId(),
		task.GetTaskData().TaskName,
		task.GetTaskData().CreateAt,
		uint64(time.Now().Nanosecond()),
	})
	proposalState := ctypes.NewProposalState(proposalId)
	proposalState.StoreOrgProposalState(ctypes.NewOrgProposalState(task.GetTaskId(), apipb.TaskRole_TaskRole_Sender, task.GetTaskSender().GetPartyId(), now))

	log.Debugf("Generate proposal, proposalId: {%s}, taskId: {%s}", proposalId, task.GetTaskId())

	// add some local cache
	t.storeProposalState(proposalState)
	t.addProposalTask(types.NewProposalTask(proposalId, task, now))


	// Start handle task ...
	go func() {

		if err := t.sendPrepareMsg(proposalId, task, now); nil != err {
			log.Errorf("Failed to call `SendTwoPcPrepareMsg`, consensus epoch finished, proposalId: {%s}, taskId: {%s}, err: \n%s", proposalId, task.GetTaskId(), err)
			// Send consensus result to Scheduler
			t.replyTaskConsensusResult(types.NewTaskConsResult(task.GetTaskId(), types.TaskConsensusInterrupt, errors.New("failed to call `SendTwoPcPrepareMsg`")))
			// clean some invalid data
			t.removeProposalStateAndTask(proposalId)
		}
	}()
	return nil
}

// TODO 问题: 自己接受的任务, 但是发起方那边任务失败了, 由于任务信息存储在本地(自己正在参与的任务), 而任务发起方怎么通知 我这边删除调自己本地正在参与的任务信息 ??? 2pc 消息已经中断了 ...

// Handle the prepareMsg from the task pulisher peer (on Subscriber)
func (t *TwoPC) onPrepareMsg(pid peer.ID, prepareMsg *types.PrepareMsgWrap) error {

	msg, err := fetchPrepareMsg(prepareMsg)
	if nil != err {
		return err
	}
	log.Debugf("Received remote prepareMsg, remote pid: {%s}, prepareMsg: %s", pid, msg.String())

	proposal := fetchProposalFromPrepareMsg(msg)
	if t.hasOrgProposal(proposal.ProposalId, msg.MsgOption.ReceiverPartyId) {
		return ctypes.ErrProposalAlreadyProcessed
	}

	identity, err := t.resourceMng.GetDB().GetIdentity()
	if nil != err {
		log.Errorf("Failed to call onPrepareMsg with `GetIdentity`, taskId: {%s}, err: {%s}", proposal.GetTaskId(), err)
		return fmt.Errorf("query local identity failed, %s", err)
	}

	sender := fetchOrgByPartyRole(msg.MsgOption.SenderPartyId, msg.MsgOption.SenderRole, msg.TaskInfo)
	receiver := fetchOrgByPartyRole(msg.MsgOption.ReceiverPartyId, msg.MsgOption.ReceiverRole, msg.TaskInfo)
	if nil == sender || nil == receiver {
		log.Errorf("Failed to verify partyId and taskRole on task, proposalId: {%s}, taskId: {%s}", msg.MsgOption.ProposalId.String(), msg.TaskInfo.GetTaskId())
		return ctypes.ErrConsensusMsgInvalid
	}

	// verify the receiver is myself ?
	if identity.GetIdentityId() != receiver.GetIdentityId() {
		log.Errorf("Failed to verify receiver identityId of prepareMsg, receiver is not me, my identityId: {%s}, receiver identityId: {%s}", identity.GetIdentityId(), receiver.GetIdentityId())
		return ctypes.ErrConsensusMsgInvalid
	}

	org := &apipb.TaskOrganization{
		PartyId:    msg.MsgOption.ReceiverPartyId,
		NodeName:   identity.GetNodeName(),
		NodeId:     identity.GetNodeId(),
		IdentityId: identity.GetIdentityId(),
	}

	// If you have already voted then we will not vote again
	if t.hasPrepareVoting(proposal.ProposalId, org) {
		return ctypes.ErrPrepareVotehadVoted
	}

	t.storeOrgProposalState(msg.MsgOption.ProposalId, ctypes.NewOrgProposalState(msg.TaskInfo.GetTaskId(), msg.MsgOption.ReceiverRole, msg.MsgOption.ReceiverPartyId, msg.CreateAt))
	t.addProposalTask(types.NewProposalTask(proposal.ProposalId, proposal.Task, proposal.CreateAt))

	// Send task to Scheduler to replay sched.
	needReplayScheduleTask := types.NewNeedReplayScheduleTask(msg.MsgOption.ReceiverRole, msg.MsgOption.ReceiverPartyId, proposal.Task)
	t.sendNeedReplayScheduleTask(needReplayScheduleTask)
	replayTaskResult := needReplayScheduleTask.ReceiveResult()

	log.Debugf("Received the reschedule task result from `schedule.ReplaySchedule()`, the result: %s", replayTaskResult.String())

	var vote *pb.PrepareVote

	if nil != replayTaskResult.GetErr() {
		vote = makePrepareVote(
			proposal.ProposalId,
			msg.MsgOption.ReceiverRole,
			msg.MsgOption.SenderRole,
			msg.MsgOption.ReceiverPartyId,
			msg.MsgOption.SenderPartyId,
			proposal.Task,
			types.No,
			&types.PrepareVoteResource{},
			uint64(timeutils.UnixMsec()),
		)
		log.Warnf("Failed to replay schedule task, will vote `NO`, taskId: {%s}, err: {%s}", replayTaskResult.GetTaskId(), replayTaskResult.GetErr().Error())
	} else {
		vote = makePrepareVote(
			proposal.ProposalId,
			msg.MsgOption.ReceiverRole,
			msg.MsgOption.SenderRole,
			msg.MsgOption.ReceiverPartyId,
			msg.MsgOption.SenderPartyId,
			proposal.Task,
			types.Yes,
			types.NewPrepareVoteResource(
				replayTaskResult.GetResource().Id,
				replayTaskResult.GetResource().Ip,
				replayTaskResult.GetResource().Port,
				replayTaskResult.GetResource().PartyId,
			),
			uint64(timeutils.UnixMsec()),
		)
		log.Infof("Succeed to replay schedule task, will vote `YES`, taskId: {%s}", replayTaskResult.GetTaskId())
	}

	// store self vote state And Send vote to Other peer
	t.storePrepareVote(types.FetchPrepareVote(vote))
	go func() {
		if err := t.sendPrepareVote(pid, receiver, sender, vote); nil != err {
			log.Errorf("failed to call `SendTwoPcPrepareVote`, proposalId: {%s}, taskId: {%s}, receiver taskPartyId:{%s}, receiver taskRole:{%s}, receiver peerId: {%s}, err: \n%s",
				proposal.ProposalId.String(), msg.TaskInfo.GetTaskId(), msg.MsgOption.SenderPartyId, msg.MsgOption.SenderRole.String(), pid, err)

			t.resourceMng.ReleaseLocalResourceWithTask("on onPrepareMsg", msg.TaskInfo.GetTaskId(), resource.SetAllReleaseResourceOption())
			// clean some data
			t.removeProposalStateAndTask(proposal.ProposalId)
		} else {
			log.Debugf("Succceed to call `SendTwoPcPrepareVote`, proposalId: {%s}, taskId: {%s}, receiver taskPartyId:{%s}, receiver taskRole:{%s}, receiver peerId: {%s}",
				proposal.ProposalId.String(), msg.TaskInfo.GetTaskId(), msg.MsgOption.SenderPartyId, msg.MsgOption.SenderRole.String(), pid)
		}
	}()
	return nil
}

// (on Publisher)
func (t *TwoPC) onPrepareVote(pid peer.ID, prepareVote *types.PrepareVoteWrap) error {

	vote := fetchPrepareVote(prepareVote)

	log.Debugf("Received remote prepareVote, remote pid: {%s}, prepareVote: %s", pid, vote.String())

	if t.hasNotOrgProposal(vote.MsgOption.ProposalId, vote.MsgOption.ReceiverPartyId) {
		return fmt.Errorf("%s onPrepareVote", ctypes.ErrProposalNotFound)
	}
	orgProposalState := t.mustGetOrgProposalState(vote.MsgOption.ProposalId, vote.MsgOption.ReceiverPartyId)

	// 只有 当前 state 是 prepare 状态才可以处理 prepare 阶段的 vote
	if orgProposalState.IsNotPreparePeriod() {
		return ctypes.ErrProposalPrepareVoteTimeout
	}

	// find the task of proposal on proposalTask
	proposalTask, ok := t.getProposalTask(orgProposalState.GetTaskId())
	if !ok {
		return fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrProposalTaskNotFound, proposalTask.GetTaskData().GetTaskId(), vote.MsgOption.ReceiverRole.String(),
			vote.MsgOption.Owner.GetIdentityId(), vote.MsgOption.ReceiverPartyId)
	}

	identity, err := t.resourceMng.GetDB().GetIdentity()
	if nil != err {
		log.Errorf("Failed to call onPrepareVote with `GetIdentity`, taskId: {%s}, err: {%s}", proposalTask.GetTaskId(), err)
		return fmt.Errorf("query local identity failed, %s", err)
	}

	sender := fetchOrgByPartyRole(vote.MsgOption.SenderPartyId, vote.MsgOption.SenderRole, proposalTask.Task)
	receiver := fetchOrgByPartyRole(vote.MsgOption.ReceiverPartyId, vote.MsgOption.ReceiverRole, proposalTask.Task)
	if nil == sender || nil == receiver {
		log.Errorf("Failed to verify partyId and taskRole of task on onPrepareVote, proposalId: {%s}, taskId: {%s}", vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId())
		return ctypes.ErrConsensusMsgInvalid
	}

	// verify the receiver is myself ?
	if identity.GetIdentityId() != receiver.GetIdentityId() {
		log.Errorf("Failed to verify receiver identityId of prepareVote, receiver is not me, my identityId: {%s}, receiver identityId: {%s}", identity.GetIdentityId(), receiver.GetIdentityId())
		return ctypes.ErrConsensusMsgInvalid
	}

	// Voter <the vote sender> voted repeatedly
	if t.hasPrepareVoting(vote.MsgOption.ProposalId, sender) {
		return fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrPrepareVoteRepeatedly, proposalTask.GetTaskData().GetTaskId(), vote.MsgOption.ReceiverRole.String(),
			vote.MsgOption.Owner.GetIdentityId(), vote.MsgOption.ReceiverPartyId)
	}

	identityValid, err := t.verifyPrepareVoteRole(vote.MsgOption.ProposalId, sender.GetPartyId(), sender.GetIdentityId(), vote.MsgOption.SenderRole, proposalTask.Task)
	if nil != err {
		return err
	}
	if !identityValid {
		return fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrProposalPrepareVoteOwnerInvalid, proposalTask.GetTaskData().GetTaskId(), vote.MsgOption.SenderRole.String(),
			sender.GetIdentityId(), sender.GetPartyId())
	}

	// verify resource of `YES` vote
	if vote.VoteOption == types.Yes && vote.PeerInfoEmpty() {
		return fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrProposalPrepareVoteResourceInvalid, proposalTask.GetTaskData().GetTaskId(), vote.MsgOption.SenderRole.String(),
			sender.GetIdentityId(), sender.GetPartyId())
	}

	// Store vote
	t.storePrepareVote(vote)

	totalNeedVoteCount := t.getNeedVotingCount(apipb.TaskRole_TaskRole_DataSupplier, proposalTask.Task) +
		t.getNeedVotingCount(apipb.TaskRole_TaskRole_PowerSupplier, proposalTask.Task) +
		t.getNeedVotingCount(apipb.TaskRole_TaskRole_Receiver, proposalTask.Task)
	yesVoteCount := t.getTaskPrepareYesVoteCount(vote.MsgOption.ProposalId)
	totalVotedCount := t.getTaskPrepareTotalVoteCount(vote.MsgOption.ProposalId)

	if totalNeedVoteCount == totalVotedCount {
		// Change the proposalState to `confirmPeriod`
		if totalNeedVoteCount == yesVoteCount {

			now := uint64(timeutils.UnixMsec())
			// 修改状态
			t.changeToConfirm(vote.MsgOption.ProposalId, vote.MsgOption.ReceiverPartyId, now)

			peers := t.makeConfirmTaskPeerDesc(vote.MsgOption.ProposalId)
			t.storeConfirmTaskPeerInfo(vote.MsgOption.ProposalId, peers)

			go func() {

				if err := t.sendConfirmMsg(vote.MsgOption.ProposalId, proposalTask.Task, peers, now); nil != err {
					log.Errorf("Failed to call `SendTwoPcConfirmMsg` proposalId: {%s}, taskId: {%s}, err: \n%s",
						vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId(), err)
					// Send consensus result
					t.replyTaskConsensusResult(types.NewTaskConsResult(proposalTask.GetTaskId(), types.TaskConsensusInterrupt, errors.New("failed to call `SendTwoPcConfirmMsg`")))
					t.removeProposalStateAndTask(vote.MsgOption.ProposalId)
				}
			}()

		} else {

			// Send consensus result
			go func() {

				log.Debugf("PrepareVoting failed on consensus's prepare epoch, the `YES` vote count is no enough, `YES` vote count: {%d}, need total count: {%d}",
					yesVoteCount, totalNeedVoteCount)

				t.replyTaskConsensusResult(types.NewTaskConsResult(proposalTask.GetTaskId(), types.TaskConsensusInterrupt, errors.New("The prepareMsg voting result was not passed")))
				t.removeProposalStateAndTask(vote.MsgOption.ProposalId)
			}()
		}
	}
	return nil
}

// (on Subscriber)
func (t *TwoPC) onConfirmMsg(pid peer.ID, confirmMsg *types.ConfirmMsgWrap) error {

	msg := fetchConfirmMsg(confirmMsg)

	log.Debugf("Received remote confirmMsg, remote pid: {%s}, confirmMsg: %s", pid, msg.String())

	if t.hasNotOrgProposal(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId) {
		return fmt.Errorf("%s onConfirmMsg", ctypes.ErrProposalNotFound)
	}

	orgProposalState := t.mustGetOrgProposalState(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId)

	// 判断是第几轮 confirmMsg
	// 只有 当前 state 是 prepare <定时任务还未更新 proposalState>
	//和 confirm <定时任务还更新 proposalState> or <现在是第二epoch> 状态才可以处理 confirm 阶段的 Msg
	// 收到第一epoch confirmMsg 时, 我应该是 prepare 阶段或者confirm 阶段,
	// 收到第二epoch confirmMsg 时, 我应该是 confirm 阶段
	if orgProposalState.IsCommitPeriod() {
		return ctypes.ErrProposalConfirmMsgTimeout
	}

	// find the task of proposal on proposalTask
	proposalTask, ok := t.getProposalTask(orgProposalState.GetTaskId())
	if !ok {
		return fmt.Errorf("%s, on the confirm msg [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrProposalTaskNotFound, proposalTask.GetTaskData().GetTaskId(), msg.MsgOption.ReceiverRole.String(),
			msg.MsgOption.Owner.GetIdentityId(), msg.MsgOption.ReceiverPartyId)
	}

	identity, err := t.resourceMng.GetDB().GetIdentity()
	if nil != err {
		t.resourceMng.ReleaseLocalResourceWithTask("on onConfirmMsg", proposalTask.GetTaskId(), resource.SetAllReleaseResourceOption())
		t.removeProposalStateAndTask(msg.MsgOption.ProposalId)

		log.Errorf("Failed to call onConfirmMsg with `GetIdentity`, taskId: {%s}, err: {%s}", proposalTask.GetTaskId(), err)
		return fmt.Errorf("query local identity failed, %s", err)
	}

	sender := fetchOrgByPartyRole(msg.MsgOption.SenderPartyId, msg.MsgOption.SenderRole, proposalTask.Task)
	receiver := fetchOrgByPartyRole(msg.MsgOption.ReceiverPartyId, msg.MsgOption.ReceiverRole, proposalTask.Task)
	if nil == sender || nil == receiver {
		t.resourceMng.ReleaseLocalResourceWithTask("on onConfirmMsg", proposalTask.GetTaskId(), resource.SetAllReleaseResourceOption())
		t.removeProposalStateAndTask(msg.MsgOption.ProposalId)

		log.Errorf("Failed to verify partyId and taskRole of task on onConfirmMsg, proposalId: {%s}, taskId: {%s}", msg.MsgOption.ProposalId.String(), proposalTask.GetTaskId())
		return ctypes.ErrConsensusMsgInvalid
	}

	// verify the receiver is myself ?
	if identity.GetIdentityId() != receiver.GetIdentityId() {
		t.resourceMng.ReleaseLocalResourceWithTask("on onConfirmMsg", proposalTask.GetTaskId(), resource.SetAllReleaseResourceOption())
		t.removeProposalStateAndTask(msg.MsgOption.ProposalId)

		log.Errorf("Failed to verify receiver identityId of confirmMsg, receiver is not me, my identityId: {%s}, receiver identityId: {%s}", identity.GetIdentityId(), receiver.GetIdentityId())
		return ctypes.ErrConsensusMsgInvalid
	}

	// verify peers resources
	if msg.PeersEmpty() {
		t.resourceMng.ReleaseLocalResourceWithTask("on onConfirmMsg", proposalTask.GetTaskId(), resource.SetAllReleaseResourceOption())
		t.removeProposalStateAndTask(msg.MsgOption.ProposalId)

		log.Errorf("Failed to verify peers resources of confirmMsg, receiver is not me, my identityId: {%s}, receiver identityId: {%s}", identity.GetIdentityId(), receiver.GetIdentityId())
		return ctypes.ErrConfirmMsgIllegal
	}

	t.storeConfirmTaskPeerInfo(msg.MsgOption.ProposalId, msg.Peers)

	vote := makeConfirmVote(
		proposalTask.ProposalId,
		msg.MsgOption.ReceiverRole,
		msg.MsgOption.SenderRole,
		msg.MsgOption.ReceiverPartyId,
		msg.MsgOption.SenderPartyId,
		proposalTask.Task,
		types.Yes,
		uint64(timeutils.UnixMsec()),
		)

	t.changeToConfirm(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId, msg.CreateAt)

	go func() {

		if err := t.sendConfirmVote(pid, receiver, sender, vote); nil != err {
			log.Errorf("failed to call `SendTwoPcConfirmVote`, proposalId: {%s}, taskId: {%s}, taskRole:{%s}, other identityId: {%s}, other peerId: {%s}, err: \n%s",
				proposalTask.ProposalId.String(), proposalTask.GetTaskId(), msg.MsgOption.SenderPartyId, msg.MsgOption.SenderRole.String(), pid, err)

			t.resourceMng.ReleaseLocalResourceWithTask("on onConfirmMsg", proposalTask.GetTaskId(), resource.SetAllReleaseResourceOption())
			// clean some data
			t.removeProposalStateAndTask(proposalTask.ProposalId)
		} else {
			log.Debugf("Succceed to call `SendTwoPcConfirmVote`, proposalId: {%s}, taskId: {%s}, taskRole:{%s}, other identityId: {%s}, other peerId: {%s}",
				proposalTask.ProposalId.String(), proposalTask.GetTaskId(), msg.MsgOption.SenderPartyId, msg.MsgOption.SenderRole.String(), pid)
		}
	}()

	return nil
}

// (on Publisher)
func (t *TwoPC) onConfirmVote(pid peer.ID, confirmVote *types.ConfirmVoteWrap) error {

	vote := fetchConfirmVote(confirmVote)

	log.Debugf("Received remote confirmVote, remote pid: {%s}, comfirmVote: %s", pid, vote.String())

	if t.hasNotOrgProposal(vote.MsgOption.ProposalId, vote.MsgOption.ReceiverPartyId) {
		return fmt.Errorf("%s onConfirmVote", ctypes.ErrProposalNotFound)
	}
	orgProposalState := t.mustGetOrgProposalState(vote.MsgOption.ProposalId, vote.MsgOption.ReceiverPartyId)
	// 只有 当前 state 是 confirm 状态才可以处理 confirm 阶段的 vote
	if orgProposalState.IsPreparePeriod() {
		return ctypes.ErrProposalConfirmVoteFuture
	}
	if orgProposalState.IsCommitPeriod() {
		return ctypes.ErrProposalPrepareVoteTimeout
	}

	// find the task of proposal on proposalTask
	proposalTask, ok := t.getProposalTask(orgProposalState.GetTaskId())
	if !ok {
		return fmt.Errorf("%s, on the confirm vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrProposalTaskNotFound, proposalTask.GetTaskData().GetTaskId(), vote.MsgOption.ReceiverRole.String(),
			vote.MsgOption.Owner.GetIdentityId(), vote.MsgOption.ReceiverPartyId)
	}

	sender := fetchOrgByPartyRole(vote.MsgOption.SenderPartyId, vote.MsgOption.SenderRole, proposalTask.Task)
	receiver := fetchOrgByPartyRole(vote.MsgOption.ReceiverPartyId, vote.MsgOption.ReceiverRole, proposalTask.Task)
	if nil == sender || nil == receiver {
		log.Errorf("Failed to verify partyId and taskRole of task on onConfirmVote, proposalId: {%s}, taskId: {%s}", vote.MsgOption.ProposalId.String(), proposalTask.GetTaskId())
		return ctypes.ErrConsensusMsgInvalid
	}

	// Voter <the vote sender> voted repeatedly
	if t.hasConfirmVoting(vote.MsgOption.ProposalId, sender) {
		return ctypes.ErrConfirmVoteRepeatedly
	}

	identityValid, err := t.verifyConfirmVoteRole(vote.MsgOption.ProposalId, sender.GetPartyId(), sender.GetIdentityId(), vote.MsgOption.SenderRole, proposalTask.Task)
	if nil != err {
		return err
	}
	if !identityValid {
		return fmt.Errorf("%s, on the confirm vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrProposalConfirmVoteVoteOwnerInvalid, proposalTask.GetTaskData().GetTaskId(), vote.MsgOption.SenderRole.String(), sender.GetIdentityId(), sender.GetPartyId())
	}

	// Store vote
	t.storeConfirmVote(vote)

	totalNeedVoteCount := t.getNeedVotingCount(apipb.TaskRole_TaskRole_DataSupplier, proposalTask.Task) +
		t.getNeedVotingCount(apipb.TaskRole_TaskRole_PowerSupplier, proposalTask.Task) +
		t.getNeedVotingCount(apipb.TaskRole_TaskRole_Receiver, proposalTask.Task)
	yesVoteCount := t.getTaskConfirmYesVoteCount(vote.MsgOption.ProposalId)
	totalVotedCount := t.getTaskConfirmTotalVoteCount(vote.MsgOption.ProposalId)

	if totalNeedVoteCount == totalVotedCount {
		// Change the proposalState to `confirmPeriod`
		if totalNeedVoteCount == yesVoteCount {


			now := uint64(timeutils.UnixMsec())

			// 修改状态
			t.changeToCommit(vote.MsgOption.ProposalId, vote.MsgOption.ReceiverPartyId, now)

			go func() {

				if err := t.sendCommitMsg(vote.MsgOption.ProposalId, proposalTask.Task, now); nil != err {

					log.Errorf("Failed to call`SendTwoPcCommitMsg` proposalId: {%s}, taskId: {%s}, err: \n%s",
						vote.MsgOption.ProposalId, proposalTask.GetTaskId(), err)

					// Send consensus result
					t.replyTaskConsensusResult(types.NewTaskConsResult(proposalTask.GetTaskId(), types.TaskConsensusInterrupt, errors.New("failed to call `SendTwoPcCommitMsg`")))

				} else {
					// Send consensus result
					t.replyTaskConsensusResult(types.NewTaskConsResult(proposalTask.GetTaskId(), types.TaskConsensusSucceed, nil))

					// If sending `CommitMsg` is successful,
					// we will forward `schedTask` to `taskManager` to send it to `Fighter` to execute the task.
					t.driveTask("", vote.MsgOption.ProposalId, vote.MsgOption.ReceiverRole, receiver, vote.MsgOption.SenderRole, sender, proposalTask.Task)
				}

				// 不成功, commitMsg 有失败的
				// 成功,  commitMsg 全发出去之后，预示着 共识完成
				t.removeProposalStateAndTask(vote.MsgOption.ProposalId)

			}()

		} else {

			// Send consensus result
			go func() {

				log.Debugf("ConfirmVoting failed on consensus's confirm epoch, the `YES` vote count is no enough, `YES` vote count: {%d}, need total count: {%d}",
					yesVoteCount, totalNeedVoteCount)

				t.replyTaskConsensusResult(types.NewTaskConsResult(proposalTask.GetTaskId(), types.TaskConsensusInterrupt, errors.New("The cofirmMsg voting result was not passed")))
				t.removeProposalStateAndTask(vote.MsgOption.ProposalId)
			}()

			// 共识 未达成. 删除本地 资源
			//// If the vote is not reached, we will clear the local `proposalState` related cache
			//// and end the task as a failure, and publish the task information to the resourceMng.GetDB().
			//t.driveTask("", vote.ProposalId, types.SendTaskDir, types.TaskStateFailed, types.TaskOwner, task)
			//// clean some invalid data
			//t.removeProposalStateAndTask(vote.ProposalId)
		}
	}
	return nil
}

// (on Subscriber)
func (t *TwoPC) onCommitMsg(pid peer.ID, cimmitMsg *types.CommitMsgWrap) error {

	msg := fetchCommitMsg(cimmitMsg)

	log.Debugf("Received remote commitMsg, remote pid: {%s}, commitMsg: %s", pid, msg.String())

	if t.hasNotOrgProposal(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId) {
		return fmt.Errorf("%s onCommitMsg", ctypes.ErrProposalNotFound)
	}

	orgProposalState := t.mustGetOrgProposalState(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId)

	// 只有 当前 state 是 confirm <定时任务还未更新 proposalState>
	// 或 commit <定时任务更新了 proposalState> 状态才可以处理 commit 阶段的 Msg
	if orgProposalState.IsPreparePeriod() {
		return ctypes.ErrProposalCommitMsgFuture
	}
	if orgProposalState.IsFinishedPeriod() {
		return ctypes.ErrProposalCommitMsgTimeout
	}

	// find the task of proposal on proposalTask
	proposalTask, ok := t.getProposalTask(orgProposalState.GetTaskId())
	if !ok {
		return fmt.Errorf("%s, on the commit msg [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrProposalTaskNotFound, proposalTask.GetTaskData().GetTaskId(), msg.MsgOption.ReceiverRole.String(),
			msg.MsgOption.Owner.GetIdentityId(), msg.MsgOption.ReceiverPartyId)
	}

	identity, err := t.resourceMng.GetDB().GetIdentity()
	if nil != err {
		t.resourceMng.ReleaseLocalResourceWithTask("on onCommitMsg", proposalTask.GetTaskId(), resource.SetAllReleaseResourceOption())
		t.removeProposalStateAndTask(msg.MsgOption.ProposalId)

		log.Errorf("Failed to call onCommitMsg with `GetIdentity`, taskId: {%s}, err: {%s}", proposalTask.GetTaskId(), err)
		return fmt.Errorf("query local identity failed, %s", err)
	}

	sender := fetchOrgByPartyRole(msg.MsgOption.SenderPartyId, msg.MsgOption.SenderRole, proposalTask.Task)
	receiver := fetchOrgByPartyRole(msg.MsgOption.ReceiverPartyId, msg.MsgOption.ReceiverRole, proposalTask.Task)
	if nil == sender || nil == receiver {
		t.resourceMng.ReleaseLocalResourceWithTask("on onCommitMsg", proposalTask.GetTaskId(), resource.SetAllReleaseResourceOption())
		t.removeProposalStateAndTask(msg.MsgOption.ProposalId)

		log.Errorf("Failed to verify partyId and taskRole of task on onCommitMsg, proposalId: {%s}, taskId: {%s}", msg.MsgOption.ProposalId.String(), proposalTask.GetTaskId())
		return ctypes.ErrConsensusMsgInvalid
	}

	// verify the receiver is myself ?
	if identity.GetIdentityId() != receiver.GetIdentityId() {
		t.resourceMng.ReleaseLocalResourceWithTask("on onCommitMsg", proposalTask.GetTaskId(), resource.SetAllReleaseResourceOption())
		t.removeProposalStateAndTask(msg.MsgOption.ProposalId)

		log.Errorf("Failed to verify receiver identityId of commitMsg, receiver is not me, my identityId: {%s}, receiver identityId: {%s}", identity.GetIdentityId(), receiver.GetIdentityId())
		return ctypes.ErrConsensusMsgInvalid
	}

	// 修改状态
	t.changeToCommit(msg.MsgOption.ProposalId, msg.MsgOption.ReceiverPartyId, msg.CreateAt)
	// If sending `CommitMsg` is successful,
	// we will forward `schedTask` to `taskManager` to send it to `Fighter` to execute the task.
	go func() {

		t.driveTask(pid, msg.MsgOption.ProposalId, msg.MsgOption.ReceiverRole, receiver, msg.MsgOption.SenderRole, sender, proposalTask.Task)
		t.removeProposalStateAndTask(msg.MsgOption.ProposalId)
	}()
	// 最后留给 定时器 清除本地 proposalState 香瓜内心戏
	return nil
}


