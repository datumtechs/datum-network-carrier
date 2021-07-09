package twopc

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/rlputil"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/handler"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"strings"
	"sync"
	"time"
)

const (
	defaultCleanExpireProposalInterval = 30 * time.Millisecond
)

type DataCenter interface {
	// identity
	HasIdentity(identity *types.NodeAlias) (bool, error)
	GetIdentity() (*types.NodeAlias, error)
}

type TwoPC struct {
	config  *Config
	p2p     p2p.P2P
	peerSet *ctypes.PeerSet
	state   *state
	dataCenter DataCenter
	// TODO 需要有一个地方 监听整个 共识结果 ...

	// fetch tasks scheduled from `Scheduler`
	schedTaskCh <-chan *types.ConsensusTaskWrap
	// send remote task to `Scheduler` to replay
	replayTaskCh chan<- *types.ScheduleTaskWrap
	asyncCallCh  chan func()
	quit         chan struct{}
	// The task being processed by myself  (taskId -> task)
	sendTasks map[string]*types.ScheduleTask
	// The task processing  that received someone else (taskId -> task)
	recvTasks map[string]*types.ScheduleTask

	taskResultCh   chan *types.ConsensuResult
	taskResultChs  map[string]chan<- *types.ConsensuResult
	taskResultLock sync.Mutex


	Errs []error
}

func New(conf *Config, dataCenter DataCenter, p2p p2p.P2P, schedTaskCh <-chan *types.ConsensusTaskWrap, replayTaskCh chan<- *types.ScheduleTaskWrap) *TwoPC {
	return &TwoPC{
		config:       conf,
		p2p:          p2p,
		peerSet:      ctypes.NewPeerSet(10), // TODO 暂时写死的
		state:        newState(),
		dataCenter:   dataCenter,
		schedTaskCh:  schedTaskCh,
		replayTaskCh: replayTaskCh,
		asyncCallCh:  make(chan func(), conf.PeerMsgQueueSize),
		quit:         make(chan struct{}),
		sendTasks:    make(map[string]*types.ScheduleTask),
		recvTasks:    make(map[string]*types.ScheduleTask),

		taskResultCh:  make(chan *types.ConsensuResult, 100),
		taskResultChs: make(map[string]chan<- *types.ConsensuResult, 100),

		Errs:         make([]error, 0),

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

		case res := <- t.taskResultCh:
			if nil == res {
				return
			}
			t.sendTaskResult(res)

		case <-cleanExpireProposalTimer.C:
			t.cleanExpireProposal()
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
	if err := t.sendPrepareMsg(task); nil != err {

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
		t.delProposalState(proposalHash)
		t.delSendTask(task.TaskId)
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

func (t *TwoPC) sendPrepareMsg(task *types.ScheduleTask) error {

	prepareMsg := makePrepareMsg()

	sendTaskFn := func(taskRole types.TaskRole, nodeId string, prepareMsg *pb.PrepareMsg, errCh chan<- error) {
		pid, err := p2p.HexPeerID(nodeId)
		if nil != err {
			errCh <- fmt.Errorf("taskRole: %s, nodeId: %s, err: %s", taskRole.String(), nodeId, err)
			return
		}

		if err = handler.SendTwoPcPrepareMsg(context.TODO(), t.p2p, pid, prepareMsg); nil != err {
			errCh <- fmt.Errorf("taskRole: %s, nodeId: %s, err: %s", taskRole.String(), nodeId, err)
			return
		}
		errCh <- nil
	}

	errCh := make(chan error, len(task.Partners)+len(task.PowerSuppliers)+len(task.Receivers))

	go func() {
		for _, partner := range task.Partners {
			go sendTaskFn(types.DataSupplier, partner.NodeId, prepareMsg, errCh)
		}
	}()
	go func() {
		for _, powerSupplier := range task.PowerSuppliers {
			go sendTaskFn(types.PowerSupplier, powerSupplier.NodeId, prepareMsg, errCh)
		}
	}()

	go func() {
		for _, receiver := range task.Receivers {
			go sendTaskFn(types.ResultSupplier, receiver.NodeId, prepareMsg, errCh)
		}
	}()
	errStrs := make([]string, 0)
	for err := range errCh {
		if nil != err {
			errStrs = append(errStrs, err.Error())
		}
	}
	if len(errStrs) != 0 {
		return fmt.Errorf(
			`failed to Send PrepareMsg for task:
%s`, strings.Join(errStrs, "\n"))
	}

	return nil
}

// Handle the prepareMsg from the task pulisher peer (on subscriber)
func (t *TwoPC) onPrepareMsg(pid peer.ID, prepareMsg *types.PrepareMsgWrap) error {


	proposal, err := t.fetchPrepareMsg(prepareMsg)
	if nil != err {
		return err
	}
	if t.state.HasProposal(proposal.ProposalId) {
		return ctypes.ErrProposalAlreadyProcessed
	}

	// If you have already voted then you will not vote again
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
		t.delProposalState(proposal.ProposalId)
		t.delRecvTask(task.TaskId)
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

	vote := &pb.PrepareVote{
		ProposalId: proposal.ProposalId.Bytes(),
		TaskRole:   prepareMsg.TaskOption.TaskRole,
		Owner: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(self.Name),
			NodeId:     []byte(self.NodeId),
			IdentityId: []byte(self.IdentityId),
		},
		CreateAt: uint64(time.Now().UnixNano()),
	}

	if result.Status == types.TaskSchedFailed {
		vote.VoteOption = types.No.Bytes()
		log.Error("Failed to replay schedule task", "taskId", result.TaskId, "err", result.Err.Error())
	} else {
		vote.VoteOption = types.Yes.Bytes()
		vote.PeerInfo = &pb.TaskPeerInfo{
			Ip:   []byte(result.Resource.Ip),
			Port: []byte(result.Resource.Port),
		}
		log.Info("Succeed to replay schedule task, will vote `YES`", "taskId", result.TaskId)
	}

	errCh := make(chan error, 0)
	go func() {
		if err = handler.SendTwoPcPrepareVote(context.TODO(), t.p2p, pid, vote); nil != err {
			err := fmt.Errorf("failed to `SendTwoPcPrepareVote`, taskId: %s, taskRole: %s, nodeId: %s, err: %s",
				proposal.TaskId, prepareMsg.TaskOption.TaskRole, self.NodeId, err)
			log.Error(err)
			errCh <- err
			return
		}
		close(errCh)
	}()

	return <-errCh
}
func (t *TwoPC) onPrepareVote(pid peer.ID, prepareVote *types.PrepareVoteWrap) error {

	voteMsg, err := t.fetchPrepareVote(prepareVote)
	if nil != err {
		return err
	}
	if t.state.HasNotProposal(voteMsg.ProposalId) {
		return ctypes.ErrProposalNotFound
	}
	proposalState := t.state.GetProposalState(voteMsg.ProposalId)
	if proposalState.IsConfirmPeriod() || proposalState.IsCommitPeriod() {
		return ctypes.ErrProposalPrepareVoteTimeout
	}

	// find the task of proposal on sendTasks
	task, ok := t.sendTasks[proposalState.TaskId]
	if !ok {
		return ctypes.ErrProposalTaskNotFound
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

	voteCount := t.state.GetTaskDataSupplierPrepareTotalVoteCount(voteMsg.ProposalId) +
		t.state.GetTaskPowerSupplierPrepareTotalVoteCount(voteMsg.ProposalId) +
		t.state.GetTaskResulterPrepareTotalVoteCount(voteMsg.ProposalId)

	now := uint64(time.Now().UnixNano())

	// Change the propoState to `confirmPeriod`
	if taskMemCount == voteCount {
		t.state.ChangeToConfirm(voteMsg.ProposalId, now)

		// TODO Send the ConfirmMsg

	}

	return nil
}
func (t *TwoPC) onConfirmMsg(pid peer.ID, confirmMsg *types.ConfirmMsgWrap) error    { return nil }
func (t *TwoPC) onConfirmVote(pid peer.ID, confirmVote *types.ConfirmVoteWrap) error { return nil }
func (t *TwoPC) onCommitMsg(pid peer.ID, cimmitMsg *types.CommitMsgWrap) error       { return nil }
