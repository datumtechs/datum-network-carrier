package twopc

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/rlputil"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/core/iface"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/handler"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"strings"
	"sync"
	"time"
)

const (
	//defaultCleanExpireProposalInterval  = 30 * time.Millisecond
	defaultRefreshProposalStateInternal = 300 * time.Millisecond
)

type TwoPC struct {
	config      *Config
	p2p         p2p.P2P
	peerSet     *ctypes.PeerSet
	state       *state
	dataCenter  iface.ForResourceDB
	resourceMng *resource.Manager

	// fetch tasks scheduled from `Scheduler`
	schedTaskCh <-chan *types.ConsensusTaskWrap
	// send remote task to `Scheduler` to replay
	replayTaskCh chan<- *types.ReplayScheduleTaskWrap
	// send has consensused remote tasks to taskManager
	doneScheduleTaskCh chan<- *types.DoneScheduleTaskChWrap
	asyncCallCh        chan func()
	quit               chan struct{}
	// The task being processed by myself  (taskId -> task)
	sendTasks map[string]*types.Task
	// The task processing  that received someone else (taskId -> task)
	recvTasks map[string]*types.Task

	taskResultCh   chan *types.ConsensuResult
	taskResultChs  map[string]chan<- *types.ConsensuResult
	taskResultLock sync.Mutex

	Errs []error
}

func New(
	conf *Config,
	dataCenter iface.ForResourceDB,
	resourceMng *resource.Manager,
	p2p p2p.P2P,
	schedTaskCh chan *types.ConsensusTaskWrap,
	replayTaskCh chan *types.ReplayScheduleTaskWrap,
	doneScheduleTaskCh chan *types.DoneScheduleTaskChWrap,
) *TwoPC {
	return &TwoPC{
		config:             conf,
		p2p:                p2p,
		peerSet:            ctypes.NewPeerSet(10), // TODO 暂时写死的
		state:              newState(),
		dataCenter:         dataCenter,
		resourceMng:        resourceMng,
		schedTaskCh:        schedTaskCh,
		replayTaskCh:       replayTaskCh,
		doneScheduleTaskCh: doneScheduleTaskCh,
		asyncCallCh:        make(chan func(), conf.PeerMsgQueueSize),
		quit:               make(chan struct{}),
		sendTasks:          make(map[string]*types.Task),
		recvTasks:          make(map[string]*types.Task),

		taskResultCh:  make(chan *types.ConsensuResult, 100),
		taskResultChs: make(map[string]chan<- *types.ConsensuResult, 100),

		Errs: make([]error, 0),
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
		case taskWrap := <-t.schedTaskCh:
			// Start a goroutine to process a new schedTask
			go func() {

				log.Debugf("Start consensus task on 2pc consensus engine, taskId: {%s}", taskWrap.Task.TaskId())

				if err := t.OnPrepare(taskWrap.Task); nil != err {
					log.Errorf("Failed to call `OnPrepare()` on 2pc consensus engine, taskId: {%s}, err: {%s}", taskWrap.Task.TaskId(), err)
					taskWrap.SendResult(&types.ConsensuResult{
						TaskConsResult: &types.TaskConsResult{
							TaskId: taskWrap.Task.TaskId(),
							Status: types.TaskConsensusInterrupt,
							Done:   false,
							Err:    fmt.Errorf("failed to OnPrepare 2pc, %s", err),
						},
					})
					return
				}
				if err := t.OnHandle(taskWrap.Task, taskWrap.OwnerDataResource, taskWrap.ResultCh); nil != err {
					log.Errorf("Failed to call `OnHandle()` on 2pc consensus engine, taskId: {%s}, err: {%s}", taskWrap.Task.TaskId(), err)
				}
			}()
		case fn := <-t.asyncCallCh:
			fn()

		case res := <-t.taskResultCh:
			if nil == res {
				return
			}
			t.sendConsensusTaskResultToSched(res)

		//case <-cleanExpireProposalTimer.C:
		//	t.cleanExpireProposal()

		case <-refreshProposalStateTicker.C:
			go t.refreshProposalState()

		case <-t.quit:
			log.Info("Stop 2pc consensus engine ...")
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
	//case *types.TaskResultMsgWrap:
	//	return t.validateTaskResultMsg(pid, msg)
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
	case *types.TaskResultMsgWrap:
		return t.onTaskResultMsg(pid, msg)
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
func (t *TwoPC) OnHandle(task *types.Task, selfPeerResource *types.PrepareVoteResource, result chan<- *types.ConsensuResult) error {

	if t.isProcessingTask(task.TaskId()) {
		return ctypes.ErrPrososalTaskIsProcessed
	}

	now := uint64(timeutils.UnixMsec())
	proposalHash := rlputil.RlpHash([]interface{}{
		t.config.Option.NodeID,
		now,
		task.TaskId(),
		task.TaskData().TaskName,
		task.TaskData().PartyId,
		task.TaskData().Identity,
		task.TaskData().NodeId,
		task.TaskData().NodeName,
		task.TaskData().DataId,
		task.TaskData().DataStatus,
		task.TaskData().State,
		task.TaskData().MetadataSupplier,
		task.TaskData().ResourceSupplier,
		task.TaskData().Receivers,
		task.TaskData().CalculateContractCode,
		task.TaskData().DataSplitContractCode,
		task.TaskData().ContractExtraParams,
		task.TaskData().TaskResource,
		task.TaskData().CreateAt,
	})

	proposalState := ctypes.NewProposalState(
		proposalHash,
		task.TaskId(),
		types.SendTaskDir,
		types.TaskOnwer,
		&types.TaskNodeAlias{
			PartyId: task.TaskData().PartyId,
			Name: task.TaskData().Identity,
			NodeId: task.TaskData().NodeId,
			IdentityId: task.TaskData().NodeName,
		},
		now)

	log.Debugf("Generate proposal, proposalId: {%s}, taskId: {%s}", proposalHash, task.TaskId())

	// add proposal
	t.addProposalState(proposalState)
	// add task
	// task 不论是 发起方 还是 参与方, 都应该是  一抵达, 就保存本地..
	t.addSendTask(task)
	// add ResultCh
	t.addTaskResultCh(task.TaskId(), result)
	// set myself peerInfo cache
	t.state.StoreSelfPeerInfo(proposalHash, selfPeerResource)

	// Start handle task ...
	go func() {

		if err := t.sendPrepareMsg(proposalHash, task, now); nil != err {
			log.Errorf("Failed to sendPrepareMsg, consensus epoch finished, proposalId: {%s}, taskId: {%s}, err: {%s}", proposalHash, task.TaskId(), err)
			// Send consensus result to Scheduler
			t.collectTaskResultWillSendToSched(&types.ConsensuResult{
				TaskConsResult: &types.TaskConsResult{
					TaskId: task.TaskId(),
					Status: types.TaskConsensusInterrupt,
					Done:   false,
					Err:    err,
				},
			})
			// clean some invalid data
			t.delProposalStateAndTask(proposalHash)
		}
	}()
	return nil
}

// TODO 问题: 自己接受的任务, 但是任务失败了, 由于任务信息存储在本地(自己正在参与的任务), 而任务发起方怎么通知 我这边删除调自己本地正在参与的任务信息 ??? 2pc 消息已经中断了 ...

// Handle the prepareMsg from the task pulisher peer (on Subscriber)
func (t *TwoPC) onPrepareMsg(pid peer.ID, prepareMsg *types.PrepareMsgWrap) error {


	msg, err := fetchPrepareMsg(prepareMsg)
	if nil != err {
		return err
	}
	log.Debugf("Received remote prepareMsg, remote pid: {%s}, prepareMsg: %s",
		pid, msg.String())

	proposal := fetchProposalFromPrepareMsg(msg)

	// 第一次接收到 发起方的 prepareMsg, 这时, 以作为接收方的身份处理msg并本地生成 proposalState
	if t.state.HasProposal(proposal.ProposalId) {
		return ctypes.ErrProposalAlreadyProcessed
	}

	// If you have already voted then we will not vote again
	if t.state.HasPrepareVoteState(proposal.ProposalId) {
		return ctypes.ErrPrepareVotehadVoted
	}


	self, err := t.dataCenter.GetIdentity()
	if nil != err {
		log.Errorf("Failed to call onPrepareMsg with `GetIdentity`, taskId: {%s}, err: {%s}", proposal.TaskId(), err)
		return fmt.Errorf("query local identity failed, %s", err)
	}

	// Create the proposal state from recv.Proposal
	proposalState := ctypes.NewProposalState(
		proposal.ProposalId,
		proposal.TaskId(),
		types.RecvTaskDir,
		types.TaskRoleFromBytes(prepareMsg.TaskRole),
		&types.TaskNodeAlias{
			PartyId: string(prepareMsg.TaskPartyId),
			Name: self.Name,
			NodeId: self.NodeId,
			IdentityId: self.IdentityId,
		},
		proposal.CreateAt)

	t.addProposalState(proposalState)

	task := proposal.Task
	// 将接受到的 task 保存本地
	// task 不论是 发起方 还是 参与方, 都应该是  一抵达, 就保存本地.. todo 让 超时检查 proposalState 机制去清除 task 和各类本地缓存
	t.addRecvTask(task)

	if err := t.validateRecvTask(task); nil != err {
		// clean some data
		t.delProposalStateAndTask(proposal.ProposalId)
		return err
	}

	// Send task to Scheduler to replay sched.
	replaySchedTask := types.NewReplayScheduleTaskWrap(
		types.TaskRoleFromBytes(prepareMsg.TaskRole),
		string(prepareMsg.TaskPartyId),
		task)

	// replay schedule task on myself ...
	t.sendReplaySchedTaskToScheduler(replaySchedTask)
	result := replaySchedTask.RecvResult()

	vote := &types.PrepareVote{
		ProposalId: proposal.ProposalId,
		TaskRole:   types.TaskRoleFromBytes(prepareMsg.TaskRole),
		Owner: &types.TaskNodeAlias{
			Name:       self.Name,
			NodeId:     self.NodeId,
			IdentityId: self.IdentityId,
			PartyId:    string(prepareMsg.TaskPartyId),
		},
		CreateAt: uint64(timeutils.UnixMsec()),
		//Sign:
	}

	if result.Status == types.TaskSchedFailed {
		vote.VoteOption = types.No
		log.Warnf("Failed to replay schedule task, will vote `NO`, taskId: {%s}, err: {%s}", result.TaskId, result.Err.Error())
	} else {
		vote.VoteOption = types.Yes
		vote.PeerInfo = &types.PrepareVoteResource{
			Ip:      result.Resource.Ip,
			Port:    result.Resource.Port,
			PartyId: result.Resource.PartyId,
		}
		log.Infof("Succeed to replay schedule task, will vote `YES`, taskId: {%s}", result.TaskId)
	}

	// store self vote state And Send vote to Other peer
	t.state.StorePrepareVoteState(vote)
	go func() {

		if err = handler.SendTwoPcPrepareVote(context.TODO(), t.p2p, pid, types.ConvertPrepareVote(vote)); nil != err {
			err := fmt.Errorf("failed to `SendTwoPcPrepareVote`, taskId: %s, taskRole: %s, nodeId: %s, err: %s",
				proposal.TaskId, prepareMsg.TaskRole, self.NodeId, err)
			log.Error(err)

			t.resourceMng.ReleaseLocalResourceWithTask("on onPrepareMsg", task.TaskId(), resource.SetAllReleaseResourceOption())
			// clean some data
			t.delProposalStateAndTask(proposal.ProposalId)
		}
	}()
	return nil
}

// (on Publisher)
func (t *TwoPC) onPrepareVote(pid peer.ID, prepareVote *types.PrepareVoteWrap) error {

	voteMsg, err := fetchPrepareVote(prepareVote)
	if nil != err {
		return err
	}

	log.Debugf("Received remote prepareVote, remote pid: {%s}, prepareVote: %s", pid, voteMsg.String())

	if t.state.HasNotProposal(voteMsg.ProposalId) {
		return ctypes.ErrProposalNotFound
	}
	proposalState := t.state.GetProposalState(voteMsg.ProposalId)
	// 只有 当前 state 是 prepare 状态才可以处理 prepare 阶段的 vote
	if !proposalState.IsPreparePeriod() {
		return ctypes.ErrProposalPrepareVoteTimeout
	}

	if t.state.IsRecvTaskOnProposalState(voteMsg.ProposalId) {
		return ctypes.ErrMsgTaskDirInvalid
	}
	// find the task of proposal on sendTasks
	task, ok := t.sendTasks[proposalState.TaskId]
	if !ok {
		return fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrProposalTaskNotFound, task.TaskData().TaskId, voteMsg.TaskRole.String(),
			voteMsg.Owner.IdentityId, voteMsg.Owner.PartyId)
	}

	// Voter voted repeatedly
	if t.state.HasPrepareVoting(voteMsg.Owner.IdentityId, voteMsg.ProposalId) {
		return fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrPrepareVoteRepeatedly, task.TaskData().TaskId, voteMsg.TaskRole.String(),
			voteMsg.Owner.IdentityId, voteMsg.Owner.PartyId)
	}

	dataSupplierCount := uint32(len(task.TaskData().MetadataSupplier) - 1)
	powerSupplierCount := uint32(len(task.TaskData().ResourceSupplier))
	resulterCount := uint32(len(task.TaskData().Receivers))
	taskMemCount := dataSupplierCount + powerSupplierCount + resulterCount

	var identityValid bool
	switch voteMsg.TaskRole {
	case types.DataSupplier:
		if dataSupplierCount == t.state.GetTaskDataSupplierPrepareTotalVoteCount(voteMsg.ProposalId) {
			return fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
				ctypes.ErrVoteCountOverflow, task.TaskData().TaskId, voteMsg.TaskRole.String(),
				voteMsg.Owner.IdentityId, voteMsg.Owner.PartyId)
		}
		for _, dataSupplier := range task.TaskData().MetadataSupplier {

			// identity + partyId
			if dataSupplier.Organization.Identity == voteMsg.Owner.IdentityId && dataSupplier.Organization.PartyId == voteMsg.Owner.PartyId {
				// TODO validate vote sign
				identityValid = true
				break
			}
		}
	case types.PowerSupplier:
		if powerSupplierCount == t.state.GetTaskPowerSupplierPrepareTotalVoteCount(voteMsg.ProposalId) {
			return fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
				ctypes.ErrVoteCountOverflow, task.TaskData().TaskId, voteMsg.TaskRole.String(),
				voteMsg.Owner.IdentityId, voteMsg.Owner.PartyId)
		}
		for _, powerSupplier := range task.TaskData().ResourceSupplier {

			// identity + partyId
			if powerSupplier.Organization.Identity == voteMsg.Owner.IdentityId && powerSupplier.Organization.PartyId == voteMsg.Owner.PartyId {
				// TODO validate vote sign
				identityValid = true
				break
			}
		}
	case types.ResultSupplier:
		if resulterCount == t.state.GetTaskResulterPrepareTotalVoteCount(voteMsg.ProposalId) {
			return fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
				ctypes.ErrVoteCountOverflow, task.TaskData().TaskId, voteMsg.TaskRole.String(),
				voteMsg.Owner.IdentityId, voteMsg.Owner.PartyId)
		}
		for _, resulter := range task.TaskData().Receivers {

			// identity + partyId
			if resulter.Receiver.Identity == voteMsg.Owner.IdentityId && resulter.Receiver.PartyId == voteMsg.Owner.PartyId {
				// TODO validate vote sign
				identityValid = true
				break
			}
		}
	default:
		return fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrMsgOwnerNodeIdInvalid, task.TaskData().TaskId, voteMsg.TaskRole.String(),
			voteMsg.Owner.IdentityId, voteMsg.Owner.PartyId)

	}
	if !identityValid {
		return fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrProposalPrepareVoteOwnerInvalid, task.TaskData().TaskId, voteMsg.TaskRole.String(),
			voteMsg.Owner.IdentityId, voteMsg.Owner.PartyId)
	}

	// Store vote
	t.state.StorePrepareVote(voteMsg)

	yesVoteCount := t.state.GetTaskPrepareYesVoteCount(voteMsg.ProposalId)
	totalVoteCount := t.state.GetTaskPrepareTotalVoteCount(voteMsg.ProposalId)

	if taskMemCount == totalVoteCount {
		// Change the propoState to `confirmPeriod`
		if taskMemCount == yesVoteCount {

			now := uint64(timeutils.UnixMsec())
			// 修改状态
			t.state.ChangeToConfirm(voteMsg.ProposalId, now)

			go func() {

				if err := t.sendConfirmMsg(voteMsg.ProposalId, task, now); nil != err {
					// Send consensus result
					t.collectTaskResultWillSendToSched(&types.ConsensuResult{
						TaskConsResult: &types.TaskConsResult{
							TaskId: task.TaskId(),
							Status: types.TaskConsensusInterrupt,
							Done:   false,
							Err:    err,
						},
						//Resources:
					})

					// clean some invalid data
					t.delProposalStateAndTask(voteMsg.ProposalId)
				}
			}()

		} else {

			// Send consensus result
			go func() {

				t.collectTaskResultWillSendToSched(&types.ConsensuResult{
					TaskConsResult: &types.TaskConsResult{
						TaskId: task.TaskId(),
						Status: types.TaskConsensusInterrupt,
						Done:   false,
						Err:    fmt.Errorf("The prepareMsg voting result was not passed"),
					},
					//Resources:
				})
				t.delProposalStateAndTask(voteMsg.ProposalId)
			}()
		}
	}
	return nil
}

// (on Subscriber)
func (t *TwoPC) onConfirmMsg(pid peer.ID, confirmMsg *types.ConfirmMsgWrap) error {

	msg, err := fetchConfirmMsg(confirmMsg)
	if nil != err {
		return err
	}

	log.Debugf("Received remote confirmMsg, remote pid: {%s}, confirmMsg: %s", pid, msg.String())

	if t.state.HasNotProposal(msg.ProposalId) {
		return ctypes.ErrProposalNotFound
	}

	// If you have already voted then we will not vote again
	if t.state.HasConfirmVoteState(msg.ProposalId) {
		return ctypes.ErrConfirmVotehadVoted
	}

	proposalState := t.state.GetProposalState(msg.ProposalId)

	// 判断是第几轮 confirmMsg
	// 只有 当前 state 是 prepare <定时任务还未更新 proposalState>
	//和 confirm <定时任务还更新 proposalState> or <现在是第二epoch> 状态才可以处理 confirm 阶段的 Msg
	// 收到第一epoch confirmMsg 时, 我应该是 prepare 阶段或者confirm 阶段,
	// 收到第二epoch confirmMsg 时, 我应该是 confirm 阶段
	if proposalState.IsCommitPeriod() {
		return ctypes.ErrProposalConfirmMsgTimeout
	}

	// 判断任务方向
	if t.state.IsSendTaskOnProposalState(msg.ProposalId) {
		return ctypes.ErrMsgTaskDirInvalid
	}
	// find the task of proposal on recvTasks
	task, ok := t.recvTasks[proposalState.TaskId]
	if !ok {
		return ctypes.ErrProposalTaskNotFound
	}

	self, err := t.dataCenter.GetIdentity()
	if nil != err {
		t.resourceMng.ReleaseLocalResourceWithTask("on onConfirmMsg", task.TaskId(), resource.SetAllReleaseResourceOption())
		// clean some data
		t.delProposalStateAndTask(proposalState.ProposalId)
		return fmt.Errorf("query local identity failed, %s", err)
	}

	vote := &types.ConfirmVote{
		ProposalId: msg.ProposalId,
		TaskRole:   msg.TaskRole,
		Owner: &types.TaskNodeAlias{
			PartyId:    msg.TaskPartyId,
			Name:       self.Name,
			NodeId:     self.NodeId,
			IdentityId: self.IdentityId,
		},
		VoteOption: types.Yes,
		CreateAt: uint64(timeutils.UnixMsec()),
		//Sign
	}

	// 修改状态
	t.state.ChangeToConfirm(msg.ProposalId, msg.CreateAt)
	// store the proposal about all partner peerInfo of task to local cache
	t.state.StoreConfirmTaskPeerInfo(msg.ProposalId, confirmMsg.PeerDesc)
	// store self vote state And Send vote to Other peer
	t.state.StoreConfirmVoteState(vote)

	go func() {

		if err = handler.SendTwoPcConfirmVote(context.TODO(), t.p2p, pid, types.ConvertConfirmVote(vote)); nil != err {
			err := fmt.Errorf("failed to `SendTwoPcConfirmVote`, taskId: %s, taskRole: %s, nodeId: %s, err: %s",
				task.TaskId, msg.TaskRole, self.NodeId, err)
			log.Error(err)

			t.resourceMng.ReleaseLocalResourceWithTask("on onConfirmMsg", task.TaskId(), resource.SetAllReleaseResourceOption())
			// clean some data
			t.delProposalStateAndTask(proposalState.ProposalId)
		}
	}()

	return nil
}

// (on Publisher)
func (t *TwoPC) onConfirmVote(pid peer.ID, confirmVote *types.ConfirmVoteWrap) error {

	voteMsg, err := fetchConfirmVote(confirmVote)
	if nil != err {
		return err
	}

	log.Debugf("Received remote confirmVote, remote pid: {%s}, comfirmVote: %s", pid, voteMsg.String())

	if t.state.HasNotProposal(voteMsg.ProposalId) {
		return ctypes.ErrProposalNotFound
	}
	proposalState := t.state.GetProposalState(voteMsg.ProposalId)
	// 只有 当前 state 是 confirm 状态才可以处理 confirm 阶段的 vote
	if proposalState.IsPreparePeriod() {
		return ctypes.ErrProposalConfirmVoteFuture
	}
	if proposalState.IsCommitPeriod() {
		return ctypes.ErrProposalPrepareVoteTimeout
	}

	if t.state.IsRecvTaskOnProposalState(voteMsg.ProposalId) {
		return ctypes.ErrMsgTaskDirInvalid
	}
	// find the task of proposal on sendTasks
	task, ok := t.sendTasks[proposalState.TaskId]
	if !ok {
		return ctypes.ErrProposalTaskNotFound
	}

	// Voter voted repeatedly
	if t.state.HasConfirmVoting(voteMsg.Owner.IdentityId, voteMsg.ProposalId) {
		return ctypes.ErrPrepareVoteRepeatedly
	}

	dataSupplierCount := uint32(len(task.TaskData().MetadataSupplier) - 1)
	powerSupplierCount := uint32(len(task.TaskData().ResourceSupplier))
	resulterCount := uint32(len(task.TaskData().Receivers))
	taskMemCount := dataSupplierCount + powerSupplierCount + resulterCount

	var identityValid bool
	switch voteMsg.TaskRole {
	case types.DataSupplier:
		if dataSupplierCount == t.state.GetTaskDataSupplierConfirmTotalVoteCount(voteMsg.ProposalId) {
			return fmt.Errorf("%s, on the confirm vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
				ctypes.ErrVoteCountOverflow, task.TaskData().TaskId, voteMsg.TaskRole.String(),
				voteMsg.Owner.IdentityId, voteMsg.Owner.PartyId)
		}
		for _, dataSupplier := range task.TaskData().MetadataSupplier {

			// identity + partyId
			if dataSupplier.Organization.Identity == voteMsg.Owner.IdentityId && dataSupplier.Organization.PartyId == voteMsg.Owner.PartyId {
				// TODO validate vote sign
				identityValid = true
				break
			}
		}
	case types.PowerSupplier:
		if powerSupplierCount == t.state.GetTaskPowerSupplierConfirmTotalVoteCount(voteMsg.ProposalId) {
			return fmt.Errorf("%s, on the confirm vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
				ctypes.ErrVoteCountOverflow, task.TaskData().TaskId, voteMsg.TaskRole.String(),
				voteMsg.Owner.IdentityId, voteMsg.Owner.PartyId)
		}
		for _, powerSupplier := range task.TaskData().ResourceSupplier {

			// identity + partyId
			if powerSupplier.Organization.Identity == voteMsg.Owner.IdentityId && powerSupplier.Organization.PartyId == voteMsg.Owner.PartyId {
				// TODO validate vote sign
				identityValid = true
				break
			}
		}
	case types.ResultSupplier:
		if resulterCount == t.state.GetTaskResulterConfirmTotalVoteCount(voteMsg.ProposalId) {
			return fmt.Errorf("%s, on the confirm vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
				ctypes.ErrVoteCountOverflow, task.TaskData().TaskId, voteMsg.TaskRole.String(),
				voteMsg.Owner.IdentityId, voteMsg.Owner.PartyId)
		}
		for _, resulter := range task.TaskData().Receivers {

			// identity + partyId
			if resulter.Receiver.Identity == voteMsg.Owner.IdentityId && resulter.Receiver.PartyId == voteMsg.Owner.PartyId {
				// TODO validate vote sign
				identityValid = true
				break
			}
		}
	default:
		return fmt.Errorf("%s, on the confirm vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrMsgOwnerNodeIdInvalid, task.TaskData().TaskId, voteMsg.TaskRole.String(),
			voteMsg.Owner.IdentityId, voteMsg.Owner.PartyId)
	}
	if !identityValid {
		return fmt.Errorf("%s, on the confirm vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrProposalConfirmVoteVoteOwnerInvalid, task.TaskData().TaskId, voteMsg.TaskRole.String(),
			voteMsg.Owner.IdentityId, voteMsg.Owner.PartyId)
	}

	// Store vote
	t.state.StoreConfirmVote(voteMsg)

	yesVoteCount := t.state.GetTaskConfirmYesVoteCount(voteMsg.ProposalId)
	totalVoteCount := t.state.GetTaskConfirmTotalVoteCount(voteMsg.ProposalId)
	// Change the propoState to `commitPeriod`
	if taskMemCount == totalVoteCount {
		if taskMemCount == yesVoteCount {

			now := uint64(timeutils.UnixMsec())

			// 修改状态
			t.state.ChangeToCommit(voteMsg.ProposalId, now)

			go func() {

				if err := t.sendCommitMsg(voteMsg.ProposalId, task, now); nil != err {
					// Send consensus result
					t.collectTaskResultWillSendToSched(&types.ConsensuResult{
						TaskConsResult: &types.TaskConsResult{
							TaskId: task.TaskId(),
							Status: types.TaskConsensusInterrupt,
							Done:   false,
							Err:    err,
						},
						//Resources:
					})

					// clean some invalid data
					t.delProposalStateAndTask(voteMsg.ProposalId)
				}
			}()

			// If sending `CommitMsg` is successful,
			// we will forward `schedTask` to `taskManager` to send it to `Fighter` to execute the task.
			t.driveTask("", voteMsg.ProposalId, types.SendTaskDir, types.TaskStateRunning, types.TaskOnwer,
				&libTypes.OrganizationData{
					PartyId: task.TaskData().PartyId,
					Identity: task.TaskData().Identity,
					NodeId: task.TaskData().NodeId,
					NodeName: task.TaskData().NodeName,
				}, task)
		} else {

			// Send consensus result
			go func() {

				t.collectTaskResultWillSendToSched(&types.ConsensuResult{
					TaskConsResult: &types.TaskConsResult{
						TaskId: task.TaskId(),
						Status: types.TaskConsensusInterrupt,
						Done:   false,
						Err:    fmt.Errorf("The prepareMsg voting result was not passed"),
					},
					//Resources:
				})
				t.delProposalStateAndTask(voteMsg.ProposalId)
			}()


			// 共识 未达成. 删除本地 资源
			//// If the vote is not reached, we will clear the local `proposalState` related cache
			//// and end the task as a failure, and publish the task information to the datacenter.
			//t.driveTask("", voteMsg.ProposalId, types.SendTaskDir, types.TaskStateFailed, types.TaskOnwer, task)
			//// clean some invalid data
			//t.delProposalStateAndTask(voteMsg.ProposalId)
		}
	}
	return nil
}

// (on Subscriber)
func (t *TwoPC) onCommitMsg(pid peer.ID, cimmitMsg *types.CommitMsgWrap) error {

	msg, err := fetchCommitMsg(cimmitMsg)
	if nil != err {
		return err
	}

	log.Debugf("Received remote commitMsg, remote pid: {%s}, commitMsg: %s", pid, msg.String())

	if t.state.HasNotProposal(msg.ProposalId) {
		return ctypes.ErrProposalNotFound
	}

	proposalState := t.state.GetProposalState(msg.ProposalId)

	// 只有 当前 state 是 confirm <定时任务还未更新 proposalState>
	// 或 commit <定时任务更新了 proposalState> 状态才可以处理 commit 阶段的 Msg
	if proposalState.IsPreparePeriod() {
		return ctypes.ErrProposalCommitMsgFuture
	}
	if proposalState.IsFinishedPeriod() {
		return ctypes.ErrProposalCommitMsgTimeout
	}

	// 判断任务方向
	if t.state.IsSendTaskOnProposalState(msg.ProposalId) {
		return ctypes.ErrMsgTaskDirInvalid
	}
	// find the task of proposal on recvTasks
	task, ok := t.recvTasks[proposalState.TaskId]
	if !ok {
		return ctypes.ErrProposalTaskNotFound
	}

	self, err := t.dataCenter.GetIdentity()
	if nil != err {
		t.resourceMng.ReleaseLocalResourceWithTask("on onCommitMsg", task.TaskId(), resource.SetAllReleaseResourceOption())
		// clean some data
		t.delProposalStateAndTask(proposalState.ProposalId)
		return fmt.Errorf("query local identity failed, %s", err)
	}

	// 修改状态
	t.state.ChangeToCommit(msg.ProposalId, msg.CreateAt)
	// If sending `CommitMsg` is successful,
	// we will forward `schedTask` to `taskManager` to send it to `Fighter` to execute the task.
	go func() {

		t.driveTask(pid, msg.ProposalId, types.RecvTaskDir, types.TaskStateRunning, msg.TaskRole,
			&libTypes.OrganizationData{
				PartyId:  msg.TaskPartyId,
				Identity: self.IdentityId,
				NodeId: self.NodeId,
				NodeName: self.Name,
			},task)
	}()
	// 最后留给 定时器 清除本地 proposalState 香瓜内心戏
	return nil
}

// Subscriber 在完成任务时对 task 生成 taskResultMsg 反馈给 发起方
func (t *TwoPC) sendTaskResultMsg(pid peer.ID, msg *types.TaskResultMsgWrap) error {
	if err := handler.SendTwoPcTaskResultMsg(context.TODO(), t.p2p, pid, msg.TaskResultMsg); nil != err {
		err := fmt.Errorf("failed to `SendTwoPcTaskResultMsg` to task owner, taskId: {%s}, taskRole: {%s}, other identityId: {%s}, other peerId: {%s}, err: {%s}",
			msg.TaskResultMsg.TaskId, msg.TaskRole, msg.TaskResultMsg.Owner.IdentityId, pid, err)
		return err
	}
	return nil
}

// (on Publisher)
func (t *TwoPC) onTaskResultMsg(pid peer.ID, taskResultMsg *types.TaskResultMsgWrap) error {
	msg, err := fetchTaskResultMsg(taskResultMsg)
	if nil != err {
		return err
	}

	log.Debugf("Received remote taskResultMsg, remote pid: {%s}, taskResultMsg: %s", pid, msg.String())

	if t.state.HasNotProposal(msg.ProposalId) {
		return ctypes.ErrProposalNotFound
	}
	proposalState := t.state.GetProposalState(msg.ProposalId)

	// 只有 当前 state 是 commit <定时任务还未更新 proposalState>
	// 或 finished <定时任务更新了 proposalState> 状态才可以处理 commit 阶段的 Msg
	if proposalState.IsNotCommitPeriod() || proposalState.IsNotFinishedPeriod() {
		return ctypes.ErrProposalTaskResultMsgTimeout
	}

	// 判断任务方向
	if t.state.IsRecvTaskOnProposalState(msg.ProposalId) {
		return ctypes.ErrMsgTaskDirInvalid
	}
	// find the task of proposal on recvTasks
	_, ok := t.sendTasks[proposalState.TaskId]
	if !ok {
		return ctypes.ErrProposalTaskNotFound
	}
	go t.storeTaskEvent(pid, proposalState.TaskId, msg.TaskEventList)
	return nil
}
