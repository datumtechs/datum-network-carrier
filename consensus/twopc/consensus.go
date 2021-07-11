package twopc

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/rlputil"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/handler"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"strings"
	"sync"
	"time"
)

const (
	defaultCleanExpireProposalInterval = 30 * time.Millisecond
	defaultRefreshProposalStateInternal = 300 * time.Millisecond
)

type DataCenter interface {
	// identity
	HasIdentity(identity *types.NodeAlias) (bool, error)
	GetIdentity() (*types.NodeAlias, error)
	StoreTaskEvent(event *types.TaskEventInfo) error
}

type TwoPC struct {
	config     *Config
	p2p        p2p.P2P
	peerSet    *ctypes.PeerSet
	state      *state
	dataCenter DataCenter
	// TODO 需要有一个地方 监听整个 共识结果 ...

	// fetch tasks scheduled from `Scheduler`
	schedTaskCh <-chan *types.ConsensusTaskWrap
	// send remote task to `Scheduler` to replay
	replayTaskCh chan<- *types.ScheduleTaskWrap
	// send has consensused remote tasks to taskManager
	recvSchedTaskCh chan<- *types.ConsensusScheduleTask
	asyncCallCh     chan func()
	quit            chan struct{}
	// The task being processed by myself  (taskId -> task)
	sendTasks map[string]*types.ScheduleTask
	// The task processing  that received someone else (taskId -> task)
	recvTasks map[string]*types.ScheduleTask

	taskResultCh   chan *types.ConsensuResult
	taskResultChs  map[string]chan<- *types.ConsensuResult
	taskResultLock sync.Mutex

	Errs []error
}

func New(conf *Config, dataCenter DataCenter, p2p p2p.P2P,
	schedTaskCh chan *types.ConsensusTaskWrap,
	replayTaskCh chan *types.ScheduleTaskWrap,
	recvSchedTaskCh chan*types.ConsensusScheduleTask,
	) *TwoPC {
	return &TwoPC{
		config:       conf,
		p2p:          p2p,
		peerSet:      ctypes.NewPeerSet(10), // TODO 暂时写死的
		state:        newState(),
		dataCenter:   dataCenter,
		schedTaskCh:  schedTaskCh,
		replayTaskCh: replayTaskCh,
		recvSchedTaskCh: recvSchedTaskCh,
		asyncCallCh:  make(chan func(), conf.PeerMsgQueueSize),
		quit:         make(chan struct{}),
		sendTasks:    make(map[string]*types.ScheduleTask),
		recvTasks:    make(map[string]*types.ScheduleTask),

		taskResultCh:  make(chan *types.ConsensuResult, 100),
		taskResultChs: make(map[string]chan<- *types.ConsensuResult, 100),

		Errs: make([]error, 0),
	}
}

func (t *TwoPC) Start() error {
	go t.loop()
	return nil
}
func (t *TwoPC) Close() error {
	close(t.quit)
	return nil
}
func (t *TwoPC) loop() {
	cleanExpireProposalTimer := time.NewTimer(defaultCleanExpireProposalInterval)
	refreshProposalStateTimer := time.NewTimer(defaultRefreshProposalStateInternal)
	for {
		select {
		case taskWrap := <-t.schedTaskCh:
			// Start a goroutine to process a new schedTask
			go func() {
				if err := t.OnPrepare(taskWrap.Task); nil != err {
					taskWrap.ResultCh <- &types.ConsensuResult{
						TaskConsResult: &types.TaskConsResult{
							TaskId: taskWrap.Task.TaskId,
							Status: types.TaskConsensusInterrupt,
							Done:   false,
							Err:    fmt.Errorf("failed to OnPrepare 2pc, %s", err),
						},
					}
					close(taskWrap.ResultCh)
					return
				}
				if err := t.OnHandle(taskWrap.Task, taskWrap.ResultCh); nil != err {
					log.Error("Failed to OnStart 2pc", "err", err)
				}
			}()
		case fn := <-t.asyncCallCh:
			fn()

		case res := <-t.taskResultCh:
			if nil == res {
				return
			}
			t.sendTaskResult(res)

		// TODO case : 需要做一次 confirmMsg 的超时 vote 重发机制, epoch 这时需要等于 2

		case <-cleanExpireProposalTimer.C:
			t.cleanExpireProposal()

		case <-refreshProposalStateTimer.C:
			go t.refreshProposalState()

		case <-t.quit:
			log.Info("Stop 2pc consensus engine ...")
			return
		}
	}
}

func (t *TwoPC) OnPrepare(task *types.ScheduleTask) error {

	return nil
}
func (t *TwoPC) OnHandle(task *types.ScheduleTask, result chan<- *types.ConsensuResult) error {

	if t.isProcessingTask(task.TaskId) {
		return ctypes.ErrPrososalTaskIsProcessed
	}

	now := uint64(time.Now().UnixNano())
	proposalHash := rlputil.RlpHash([]interface{}{
		t.config.Option.NodeID,
		now,
		task.TaskId,
		task.TaskName,
		task.Owner,
		task.Partners,
		task.PowerSuppliers,
		task.Receivers,
		task.OperationCost,
		task.CreateAt,
	})

	proposalState := ctypes.NewProposalState(proposalHash, task.TaskId, ctypes.SendTaskDir, now)
	// add proposal
	t.addProposalState(proposalState)
	// add task
	t.addSendTask(task)
	// add ResultCh
	t.addTaskResultCh(task.TaskId, result)

	// Start handle task ...
	if err := t.sendPrepareMsg(proposalHash, task, now); nil != err {

		// Send consensus result
		t.collectTaskResult(&types.ConsensuResult{
			TaskConsResult: &types.TaskConsResult{
				TaskId: task.TaskId,
				Status: types.TaskConsensusInterrupt,
				Done:   false,
				Err:    err,
			},
		})
		// clean some invalid data
		t.delProposalStateAndTask(proposalHash)
		return err
	}

	return nil
}

func (t *TwoPC) ValidateConsensusMsg(pid peer.ID, msg types.ConsensusMsg) error {
	if nil == msg {
		return fmt.Errorf("Failed to validate 2pc consensus msg, the msg is nil")
	}
	switch msg := msg.(type) {
	case *types.PrepareMsgWrap:
		return t.validatePrepareMsg(pid, msg)
	case *types.PrepareVoteWrap:
		return t.validatePrepareVote(pid, msg)
	case *types.ConfirmMsgWrap:
		return t.validateConfirmMsg(pid, msg)
	case *types.ConfirmVoteWrap:
		return t.validateConfirmVote(pid, msg)
	case *types.CommitMsgWrap:
		return t.validateCommitMsg(pid, msg)
	case *types.TaskResultMsgWrap:
		return t.validateTaskResultMsg(pid, msg)
	default:
		return fmt.Errorf("TaskRoleUnknown the 2pc msg type")
	}
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

// TODO 问题: 自己接受的任务, 但是任务失败了, 由于任务信息存储在本地(自己正在参与的任务), 而任务发起方怎么通知 我这边删除调自己本地正在参与的任务信息 ??? 2pc 消息已经中断了 ...

// Handle the prepareMsg from the task pulisher peer (on subscriber)
func (t *TwoPC) onPrepareMsg(pid peer.ID, prepareMsg *types.PrepareMsgWrap) error {

	proposal, err := fetchPrepareMsg(prepareMsg)
	if nil != err {
		return err
	}

	// 第一次接收到 发起方的 prepareMsg, 这时, 以作为接收方的身份处理msg并本地生成 proposalState
	if t.state.HasProposal(proposal.ProposalId) {
		return ctypes.ErrProposalAlreadyProcessed
	}

	// If you have already voted then we will not vote again
	if t.state.HasPrepareVoteState(proposal.ProposalId) {
		return ctypes.ErrPrepareVotehadVoted
	}

	// Create the proposal state from recv.Proposal
	proposalState := ctypes.NewProposalState(proposal.ProposalId,
		proposal.TaskId, ctypes.RecvTaskDir, proposal.CreateAt)

	task := proposal.ScheduleTask

	t.addProposalState(proposalState)
	t.addRecvTask(task)
	if err := t.validateRecvTask(task); nil != err {
		// clean some data
		t.delProposalStateAndTask(proposal.ProposalId)
		return err
	}

	// Send task to Scheduler to replay sched.
	replaySchedTask := &types.ScheduleTaskWrap{
		Role:     types.TaskRoleFromBytes(prepareMsg.TaskOption.TaskRole),
		Task:     task,
		ResultCh: make(chan *types.ScheduleResult),
	}
	t.sendReplaySchedTask(replaySchedTask)
	result := replaySchedTask.RecvResult()

	self, err := t.dataCenter.GetIdentity()
	if nil != err {
		return err
	}

	vote := &types.PrepareVote{
		ProposalId: proposal.ProposalId,
		TaskRole:   types.TaskRoleFromBytes(prepareMsg.TaskOption.TaskRole),
		Owner: &types.NodeAlias{
			Name:       self.Name,
			NodeId:     self.NodeId,
			IdentityId: self.IdentityId,
		},
		CreateAt: uint64(time.Now().UnixNano()),
	}

	if result.Status == types.TaskSchedFailed {
		vote.VoteOption = types.No
		log.Error("Failed to replay schedule task", "taskId", result.TaskId, "err", result.Err.Error())
	} else {
		vote.VoteOption = types.Yes
		vote.PeerInfo = &types.PrepareVoteResource{
			Ip:   result.Resource.Ip,
			Port: result.Resource.Port,
		}
		log.Info("Succeed to replay schedule task, will vote `YES`", "taskId", result.TaskId)
	}

	// store self vote state And Send vote to Other peer
	t.state.StorePrepareVoteState(vote)
	if err = handler.SendTwoPcPrepareVote(context.TODO(), t.p2p, pid, types.ConvertPrepareVote(vote)); nil != err {
		err := fmt.Errorf("failed to `SendTwoPcPrepareVote`, taskId: %s, taskRole: %s, nodeId: %s, err: %s",
			proposal.TaskId, prepareMsg.TaskOption.TaskRole, self.NodeId, err)
		log.Error(err)
		return err
	}

	return nil
}
func (t *TwoPC) onPrepareVote(pid peer.ID, prepareVote *types.PrepareVoteWrap) error {

	voteMsg, err := fetchPrepareVote(prepareVote)
	if nil != err {
		return err
	}
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
		return ctypes.ErrProposalTaskNotFound
	}

	// Voter voted repeatedly
	if t.state.HasPrepareVoting(voteMsg.Owner.IdentityId, voteMsg.ProposalId) {
		return ctypes.ErrPrepareVoteRepeatedly
	}

	dataSupplierCount := uint32(len(task.Partners))
	powerSupplierCount := uint32(len(task.PowerSuppliers))
	resulterCount := uint32(len(task.Receivers))

	taskMemCount := dataSupplierCount + powerSupplierCount + resulterCount

	var identityValid bool
	switch voteMsg.TaskRole {
	case types.DataSupplier:
		if dataSupplierCount == t.state.GetTaskDataSupplierPrepareTotalVoteCount(voteMsg.ProposalId) {
			return ctypes.ErrVoteCountOverflow
		}
		for _, dataSupplier := range task.Partners {
			if dataSupplier.IdentityId == voteMsg.Owner.IdentityId {

				// TODO validate vote sign

				identityValid = true
				break
			}
		}
	case types.PowerSupplier:
		if powerSupplierCount == t.state.GetTaskPowerSupplierPrepareTotalVoteCount(voteMsg.ProposalId) {
			return ctypes.ErrVoteCountOverflow
		}
		for _, powerSupplier := range task.PowerSuppliers {
			if powerSupplier.IdentityId == voteMsg.Owner.IdentityId {

				// TODO validate vote sign

				identityValid = true
				break
			}
		}
	case types.ResultSupplier:
		if resulterCount == t.state.GetTaskResulterPrepareTotalVoteCount(voteMsg.ProposalId) {
			return ctypes.ErrVoteCountOverflow
		}
		for _, resulter := range task.Receivers {
			if resulter.IdentityId == voteMsg.Owner.IdentityId {

				// TODO validate vote sign

				identityValid = true
				break
			}
		}
	default:
		return ctypes.ErrMsgOwnerNodeIdInvalid
	}
	if !identityValid {
		return ctypes.ErrProposalPrepareVoteOwnerInvalid
	}

	// Store vote
	t.state.StorePrepareVote(voteMsg)

	yesVoteCount := t.state.GetTaskDataSupplierPrepareYesVoteCount(voteMsg.ProposalId) +
		t.state.GetTaskPowerSupplierPrepareYesVoteCount(voteMsg.ProposalId) +
		t.state.GetTaskResulterPrepareYesVoteCount(voteMsg.ProposalId)

	// Change the propoState to `confirmPeriod`
	if taskMemCount == yesVoteCount {

		now := uint64(time.Now().UnixNano())

		t.state.ChangeToConfirm(voteMsg.ProposalId, now)




		// TODO ++++++++++++++++++++++++++++++++
		// TODO ++++++++++++++++++++++++++++++++
		// TODO ++++++++++++++++++++++++++++++++
		//
		// TODO 发送 confirmMsg 失败后需要重新在发一次
		//
		// TODO ++++++++++++++++++++++++++++++++
		// TODO ++++++++++++++++++++++++++++++++
		// TODO ++++++++++++++++++++++++++++++++

		if err := t.sendConfirmMsg(voteMsg.ProposalId, task, now); nil != err {
			// Send consensus result
			t.collectTaskResult(&types.ConsensuResult{
				TaskConsResult: &types.TaskConsResult{
					TaskId: task.TaskId,
					Status: types.TaskConsensusInterrupt,
					Done:   false,
					Err:    err,
				},
			})

			// todo 这里先考虑下是否直接清除掉 proposalState ????
			// clean some invalid data
			t.delProposalStateAndTask(voteMsg.ProposalId)
			return err
		}

	}

	return nil
}
func (t *TwoPC) onConfirmMsg(pid peer.ID, confirmMsg *types.ConfirmMsgWrap) error {

	msg, err := fetchConfirmMsg(confirmMsg)
	if nil != err {
		return err
	}

	if t.state.HasNotProposal(msg.ProposalId) {
		return ctypes.ErrProposalNotFound
	}

	// If you have already voted then we will not vote again
	if t.state.HasConfirmVoteState(msg.ProposalId, msg.Epoch) {
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
		return err
	}

	now := uint64(time.Now().UnixNano())

	vote := &types.ConfirmVote{
		ProposalId: msg.ProposalId,
		Epoch:      msg.Epoch,
		TaskRole:   msg.TaskRole,
		Owner:      msg.Owner,
		VoteOption: types.Yes,
		CreateAt:   now,
		//Sign
	}

	// 修改状态
	t.state.ChangeToConfirm(msg.ProposalId, msg.CreateAt)
	// store the proposal about all partner peerInfo of task to local cache
	t.state.StoreConfirmTaskPeerInfo(msg.ProposalId, confirmMsg.PeerDesc)
	// store self vote state And Send vote to Other peer
	t.state.StoreConfirmVoteState(vote)
	if err = handler.SendTwoPcConfirmVote(context.TODO(), t.p2p, pid, types.ConvertConfirmVote(vote)); nil != err {
		err := fmt.Errorf("failed to `SendTwoPcConfirmVote`, taskId: %s, taskRole: %s, nodeId: %s, err: %s",
			task.TaskId, msg.TaskRole, self.NodeId, err)
		log.Error(err)
		return err
	}
	return nil
}
func (t *TwoPC) onConfirmVote(pid peer.ID, confirmVote *types.ConfirmVoteWrap) error {

	voteMsg, err := fetchConfirmVote(confirmVote)
	if nil != err {
		return err
	}
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

	dataSupplierCount := uint32(len(task.Partners))
	powerSupplierCount := uint32(len(task.PowerSuppliers))
	resulterCount := uint32(len(task.Receivers))
	taskMemCount := dataSupplierCount + powerSupplierCount + resulterCount

	var identityValid bool
	switch voteMsg.TaskRole {
	case types.DataSupplier:
		if dataSupplierCount == t.state.GetTaskDataSupplierConfirmTotalVoteCount(voteMsg.ProposalId) {
			return ctypes.ErrVoteCountOverflow
		}
		for _, dataSupplier := range task.Partners {
			if dataSupplier.IdentityId == voteMsg.Owner.IdentityId {

				// TODO validate vote sign

				identityValid = true
				break
			}
		}
	case types.PowerSupplier:
		if powerSupplierCount == t.state.GetTaskPowerSupplierConfirmTotalVoteCount(voteMsg.ProposalId) {
			return ctypes.ErrVoteCountOverflow
		}
		for _, powerSupplier := range task.PowerSuppliers {
			if powerSupplier.IdentityId == voteMsg.Owner.IdentityId {

				// TODO validate vote sign

				identityValid = true
				break
			}
		}
	case types.ResultSupplier:
		if resulterCount == t.state.GetTaskResulterConfirmTotalVoteCount(voteMsg.ProposalId) {
			return ctypes.ErrVoteCountOverflow
		}
		for _, resulter := range task.Receivers {
			if resulter.IdentityId == voteMsg.Owner.IdentityId {

				// TODO validate vote sign

				identityValid = true
				break
			}
		}
	default:
		return ctypes.ErrMsgOwnerNodeIdInvalid
	}
	if !identityValid {
		return ctypes.ErrProposalConfirmVoteVoteOwnerInvalid
	}

	// Store vote
	t.state.StoreConfirmVote(voteMsg)

	yesVoteCount := t.state.GetTaskDataSupplierConfirmYesVoteCount(voteMsg.ProposalId) +
		t.state.GetTaskPowerSupplierConfirmYesVoteCount(voteMsg.ProposalId) +
		t.state.GetTaskResulterConfirmYesVoteCount(voteMsg.ProposalId)

	// Change the propoState to `commitPeriod`
	if taskMemCount == yesVoteCount {

		now := uint64(time.Now().UnixNano())

		t.state.ChangeToCommit(voteMsg.ProposalId, now)

		if err := t.sendCommitMsg(voteMsg.ProposalId, task, now); nil != err {
			// Send consensus result
			t.collectTaskResult(&types.ConsensuResult{
				TaskConsResult: &types.TaskConsResult{
					TaskId: task.TaskId,
					Status: types.TaskConsensusInterrupt,
					Done:   false,
					Err:    err,
				},
			})

			// todo 这里先考虑下是否直接清除掉 proposalState ????
			// clean some invalid data
			t.delProposalStateAndTask(voteMsg.ProposalId)
			return err
		}

		// If sending `CommitMsg` is successful, we will forward `schedTask` to `taskManager` to send it to `Fighter` to execute the task.
		if err := t.driveTask(ctypes.SendTaskDir, voteMsg.ProposalId, task.TaskId); nil != err {
			log.Error("Failed to drive the local schedTask to taskManager to execute", "proposald", voteMsg.ProposalId, "taskId", task.TaskId)
			// TODO 发送失败的的异常 先不处理了, 正常这里应该是 直接将任务结束 发给数据中心的 ...
		}
	}

	// TODO 不是 yes 票, 且满足任务结束, 需要直接结束任务

	return nil
}

func (t *TwoPC) onCommitMsg(pid peer.ID, cimmitMsg *types.CommitMsgWrap) error {

	msg, err := fetchCommitMsg(cimmitMsg)
	if nil != err {
		return err
	}

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

	// If sending `CommitMsg` is successful, we will forward `schedTask` to `taskManager` to send it to `Fighter` to execute the task.
	if err := t.driveTask(ctypes.RecvTaskDir, msg.ProposalId, task.TaskId); nil != err {
		log.Error("Failed to drive the remote schedTask to taskManager to execute", "proposald", msg.ProposalId, "taskId", task.TaskId)
		// TODO 发送失败的的异常 先不处理了, 正常这里应该是 直接将任务结束 发给数据中心的 ...
	}

	return nil
}


func (t *TwoPC) onTaskResultMsg(pid peer.ID, taskResultMsg *types.TaskResultMsgWrap) error {
	msg, err := fetchTaskResultMsg(taskResultMsg)
	if nil != err {
		return err
	}

	if t.state.HasNotProposal(msg.ProposalId) {
		return ctypes.ErrProposalNotFound
	}
	proposalState := t.state.GetProposalState(msg.ProposalId)

	// 只有 当前 state 是 confirm <定时任务还未更新 proposalState>
	// 或 commit <定时任务更新了 proposalState> 状态才可以处理 commit 阶段的 Msg
	if proposalState.IsNotCommitPeriod() || proposalState.IsNotFinishedPeriod() {
		return ctypes.ErrProposalCommitMsgFuture
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

	t.storeTaskEvent(pid, proposalState.TaskId, msg.TaskEventList)
	return nil
}

