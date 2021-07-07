package twopc

import (
	"bytes"
	"context"
	"crypto/elliptic"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/rlputil"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/handler"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
	"time"
)

type DataCenter interface {
	// identity
	HasIdentity(identity *types.NodeAlias) (bool, error)
}

type TwoPC struct {
	config  *Config
	p2p     p2p.P2P
	peerSet *ctypes.PeerSet
	state   *state
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

	dataCenter DataCenter

	Errs []error
}

func New(conf *Config, dataCenter DataCenter, p2p p2p.P2P, taskCh <-chan *types.ConsensusTaskWrap, replayTaskCh chan<- *types.ScheduleTaskWrap) *TwoPC {
	return &TwoPC{
		config:       conf,
		p2p:          p2p,
		peerSet:      ctypes.NewPeerSet(10), // TODO 暂时写死的
		state:        newState(),
		schedTaskCh:  taskCh,
		replayTaskCh: replayTaskCh,
		asyncCallCh:  make(chan func(), conf.PeerMsgQueueSize),
		sendTasks:    make(map[string]*types.ScheduleTask),
		recvTasks:    make(map[string]*types.ScheduleTask),
		dataCenter:   dataCenter,
		Errs:         make([]error, 0),
	}
}

func (t *TwoPC) Start() error {
	go t.loop()
	return nil
}
func (t *TwoPC) Close() error {return nil}
func (t *TwoPC) loop() {

	for {
		select {
		case taskWrap := <-t.schedTaskCh:
			if err := t.OnPrepare(taskWrap.Task); nil != err {
				taskWrap.ResultCh <- &types.ConsensuResult{
					TaskConsResult: &types.TaskConsResult{
						TaskId: taskWrap.Task.TaskId,
						Status: types.TaskConsensusInterrupt,
						Done:   false,
						Err:    fmt.Errorf("failed to OnPrepare 2pc, %s", err),
					},
				}
				continue
			}
			if err := t.OnHandle(taskWrap.Task, taskWrap.ResultCh); nil != err {
				log.Error("Failed to OnStart 2pc", "err", err)
			}
		case fn := <-t.asyncCallCh:
			fn()
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

	proposalState := ctypes.NewProposalState(proposalHash, now)
	t.state.AddProposalState(proposalState)
	t.addSendTask(task)
	prepareMsg := assemblingPrepareMsgForPB()

	sendTaskFn := func(taskRole types.TaskRole, nodeId string, prepareMsg *pb.PrepareMsg, errCh chan<- error) {
		nid, err := p2p.HexID(nodeId)
		if nil != err {
			errCh <- fmt.Errorf("taskRole: %s, nodeId: %s, err: %s", taskRole.String(), nodeId, err)
			return
		}
		peer := t.peerSet.GetPeer(nid)
		if err = handler.SendTwoPcPrepareMsg(context.TODO(), t.p2p, peer.GetPid(), prepareMsg); nil != err {
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
		err := fmt.Errorf(
			`failed to Send PrepareMsg for task:
%s`, strings.Join(errStrs, "\n"))
		result <- &types.ConsensuResult{
			TaskConsResult:  &types.TaskConsResult{
				TaskId: task.TaskId,
				Status: types.TaskConsensusInterrupt,
				Done:   false,
				Err:    err,
			},
		}
		// clean some invalid data
		t.state.DelProposalState(proposalHash)
		t.delSendTask(task.TaskId)

		return err
	}

	return nil
}

func (t *TwoPC) ValidateConsensusMsg(msg types.ConsensusMsg) error {
	if nil == msg {
		return fmt.Errorf("Failed to validate 2pc consensus msg, the msg is nil")
	}
	switch msg := msg.(type) {
	case *types.PrepareMsgWrap:
		return t.validatePrepareMsg(msg)
	case *types.PrepareVoteWrap:
		return t.validatePrepareVote(msg)
	case *types.ConfirmMsgWrap:
		return t.validateConfirmMsg(msg)
	case *types.ConfirmVoteWrap:
		return t.validateConfirmVote(msg)
	case *types.CommitMsgWrap:
		return t.validateCommitMsg(msg)
	default:
		return fmt.Errorf("TaskRoleUnknown the 2pc msg type")
	}
}

func (t *TwoPC) OnConsensusMsg(msg types.ConsensusMsg) error {

	switch msg := msg.(type) {
	case *types.PrepareMsgWrap:
		return t.onPrepareMsg(msg)
	case *types.PrepareVoteWrap:
		return t.onPrepareVote(msg)
	case *types.ConfirmMsgWrap:
		return t.onConfirmMsg(msg)
	case *types.ConfirmVoteWrap:
		return t.onConfirmVote(msg)
	case *types.CommitMsgWrap:
		return t.onCommitMsg(msg)

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

// With subscriber
func (t *TwoPC) validatePrepareMsg(prepareMsg *types.PrepareMsgWrap) error {
	proposalId := common.BytesToHash(prepareMsg.ProposalId)
	if t.state.HasProposal(proposalId) {
		return ctypes.ErrProposalAlreadyProcessed
	}
	// Now, we has not  proposalState (on subscriber)

	now := uint64(time.Now().UnixNano())
	if prepareMsg.CreateAt >= now {
		return ctypes.ErrProposalInTheFuture
	}

	if len(prepareMsg.TaskOption.TaskRole) != 1 ||
		types.TaskRoleFromBytes(prepareMsg.TaskOption.TaskRole) == types.TaskRoleUnknown {
		return ctypes.ErrPrososalTaskRoleIsUnknown
	}
	taskId := string(prepareMsg.TaskOption.TaskId)

	if t.isProcessingTask(taskId) {
		return ctypes.ErrPrososalTaskIsProcessed
	}

	// Verify the signature
	_, err := t.verifyMsgSigned(prepareMsg.TaskOption.Owner.NodeId, prepareMsg.SealHash().Bytes(), prepareMsg.Signature())
	if err != nil {
		return err
	}

	// validate the owner
	if err := t.validateOrganizationIdentity(prepareMsg.TaskOption.Owner); nil != err {
		log.Error("Failed to validate prepareMsg, the owner organization identity is invalid", "err", err)
		return err
	}

	// validate the algo supplier
	if err := t.validateOrganizationIdentity(prepareMsg.TaskOption.AlgoSupplier); nil != err {
		log.Error("Failed to validate prepareMsg, the algoSupplier organization identity is invalid", "err", err)
		return err
	}

	// validate data suppliers
	if len(prepareMsg.TaskOption.DataSupplier) == 0 {
		return fmt.Errorf("Failed to validate prepareMsg, the data suppliers is empty")
	}
	for _, supplier := range prepareMsg.TaskOption.DataSupplier {
		if err := t.validateOrganizationIdentity(supplier.MemberInfo); nil != err {
			log.Error("Failed to validate prepareMsg, the dataSupplier organization identity is invalid", "err", err)
			return err
		}
	}

	// validate power suppliers
	if len(prepareMsg.TaskOption.PowerSupplier) == 0 {
		return fmt.Errorf("Failed to validate prepareMsg, the data suppliers is empty")
	}
	powerSuppliers := make(map[string]struct{}, len(prepareMsg.TaskOption.PowerSupplier))
	for _, supplier := range prepareMsg.TaskOption.PowerSupplier {
		if err := t.validateOrganizationIdentity(supplier.MemberInfo); nil != err {
			log.Error("Failed to validate prepareMsg, the powerSupplier organization identity is invalid", "err", err)
			return err
		}
		powerSuppliers[string(supplier.MemberInfo.IdentityId)] = struct{}{}
	}

	// validate receicers
	if len(prepareMsg.TaskOption.Receivers) == 0 {
		return fmt.Errorf("Failed to validate prepareMsg, the data suppliers is empty")
	}

	for _, supplier := range prepareMsg.TaskOption.Receivers {
		if err := t.validateOrganizationIdentity(supplier.MemberInfo); nil != err {
			log.Error("Failed to validate prepareMsg, the powerSupplier organization identity is invalid", "err", err)
			return err
		}
		receiverIdentityId := string(supplier.MemberInfo.IdentityId)
		for _, provider := range supplier.Providers {
			providerIdentityId := string(provider.IdentityId)
			if _, ok := powerSuppliers[providerIdentityId]; !ok {
				return fmt.Errorf("Has invalid provider with receiver", "providerIdentityId", providerIdentityId, "receiverIndentityId", receiverIdentityId)
			}
		}
	}

	// validate OperationCost
	if 0 == prepareMsg.TaskOption.OperationCost.CostProcessor ||
		0 == prepareMsg.TaskOption.OperationCost.CostMem ||
		0 == prepareMsg.TaskOption.OperationCost.CostBandwidth ||
		0 == prepareMsg.TaskOption.OperationCost.Duration {
		return ctypes.ErrProposalTaskOperationCostInvalid
	}

	// validate contractCode
	if 0 == len(prepareMsg.TaskOption.CalculateContractCode) {
		return ctypes.ErrProposalTaskCalculateContractCodeEmpty
	}

	// validate task create time
	if prepareMsg.TaskOption.CreateAt >= prepareMsg.CreateAt {
		return ctypes.ErrProposalParamsInvalid
	}
	return nil
}

// With publisher
func (t *TwoPC) validatePrepareVote(prepareVote *types.PrepareVoteWrap) error {
	proposalId := common.BytesToHash(prepareVote.ProposalId)
	if t.state.HasNotProposal(proposalId) {
		return ctypes.ErrProposalNotFound
	}
	// (On publisher)
	// Now, we has proposalState, current period is `preparePeriod`
	// the PeriodStartTime is not zero, But the PeriodEndTime is zero
	proposalState := t.state.ProposalStates(proposalId)
	if proposalState.IsNotPreparePeriod() {
		return ctypes.ErrPrepareVoteIllegal
	}

	// earlier vote is invalid vote
	if proposalState.PeriodStartTime >= prepareVote.CreateAt {
		return ctypes.ErrPrepareVoteIllegal
	}

	now := uint64(time.Now().UnixNano())
	if prepareVote.CreateAt >= now {
		return ctypes.ErrPrepareVoteInTheFuture
	}

	if types.TaskRoleFromBytes(prepareVote.TaskRole) == types.TaskRoleUnknown {
		return ctypes.ErrPrososalTaskRoleIsUnknown
	}

	// Verify the signature
	_, err := t.verifyMsgSigned(prepareVote.Owner.NodeId, prepareVote.SealHash().Bytes(), prepareVote.Signature())
	if err != nil {
		return err
	}

	// validate the owner
	if err := t.validateOrganizationIdentity(prepareVote.Owner); nil != err {
		log.Error("Failed to validate prepareVote, the owner organization identity is invalid", "err", err)
		return err
	}

	// validate voteOption
	if len(prepareVote.VoteOption) != 1 ||
		types.VoteOptionFromBytes(prepareVote.VoteOption) == types.VoteUnknown {
		return ctypes.ErrPrepareVoteOptionIsUnknown
	}
	if 0 == len(prepareVote.PeerInfo.Ip) || 0 == len(prepareVote.PeerInfo.Port) {
		return ctypes.ErrPrepareVoteParamsInvalid
	}

	return nil
}

// With subscriber
func (t *TwoPC) validateConfirmMsg(confirmMsg *types.ConfirmMsgWrap) error {

	proposalId := common.BytesToHash(confirmMsg.ProposalId)
	if t.state.HasNotProposal(proposalId) {
		return ctypes.ErrProposalNotFound
	}

	// (On subscriber)
	// Now, current period should be `confirmPeriod`,
	// If confirmMsg is first epoch, so still stay preparePeriod
	// If confirmMsg is second epoch, so stay confirmPeriod
	proposalState := t.state.ProposalStates(proposalId)
	if proposalState.IsNotPreparePeriod() || proposalState.IsNotConfirmPeriod() {
		return ctypes.ErrConfirmMsgIllegal
	}
	// When comfirm first epoch
	if proposalState.IsPreparePeriod() && ctypes.ConfirmEpochFirst.Uint64() != confirmMsg.Epoch {
		return ctypes.ErrConfirmMsgIllegal
	}
	// When comfirm second epoch
	if proposalState.IsConfirmPeriod() && ctypes.ConfirmEpochSecond.Uint64() != confirmMsg.Epoch {
		return ctypes.ErrConfirmMsgIllegal
	}

	// earlier confirmMsg is invalid msg
	if proposalState.PeriodStartTime >= confirmMsg.CreateAt {
		return ctypes.ErrConfirmMsgIllegal
	}
	now := uint64(time.Now().UnixNano())
	if confirmMsg.CreateAt >= now {
		return ctypes.ErrConfirmMsgInTheFuture
	}

	// validate epoch number
	if confirmMsg.Epoch == 0 || confirmMsg.Epoch > types.MsgEpochMaxNumber {
		return ctypes.ErrConfirmMsgEpochInvalid
	}

	// Verify the signature
	_, err := t.verifyMsgSigned(confirmMsg.Owner.NodeId, confirmMsg.SealHash().Bytes(), confirmMsg.Signature())
	if err != nil {
		return err
	}

	// validate the owner
	if err := t.validateOrganizationIdentity(confirmMsg.Owner); nil != err {
		log.Error("Failed to validate confirmMsg, the owner organization identity is invalid", "err", err)
		return err
	}

	return nil
}

// With publisher
func (t *TwoPC) validateConfirmVote(confirmVote *types.ConfirmVoteWrap) error {
	proposalId := common.BytesToHash(confirmVote.ProposalId)
	if t.state.HasNotProposal(proposalId) {
		return ctypes.ErrProposalNotFound
	}

	// (On publisher)
	// Now, current period should be `confirmPeriod`,
	// No matter first epoch of confirmVote or second epoch of confirmVote
	// it is always confirmPeriod.
	proposalState := t.state.ProposalStates(proposalId)
	if proposalState.IsNotConfirmPeriod() {
		return ctypes.ErrConfirmVoteIllegal
	}

	// earlier confirmVote is invalid vote
	if proposalState.PeriodStartTime >= confirmVote.CreateAt {
		return ctypes.ErrConfirmVoteIllegal
	}
	now := uint64(time.Now().UnixNano())
	if confirmVote.CreateAt >= now {
		return ctypes.ErrConfirmVoteInTheFuture
	}

	// validate epoch number
	if confirmVote.Epoch == 0 || confirmVote.Epoch > types.MsgEpochMaxNumber {
		return ctypes.ErrConfirmVoteEpochInvalid
	}

	if types.TaskRoleFromBytes(confirmVote.TaskRole) == types.TaskRoleUnknown {
		return ctypes.ErrPrososalTaskRoleIsUnknown
	}

	// Verify the signature
	_, err := t.verifyMsgSigned(confirmVote.Owner.NodeId, confirmVote.SealHash().Bytes(), confirmVote.Signature())
	if err != nil {
		return err
	}

	// validate the owner
	if err := t.validateOrganizationIdentity(confirmVote.Owner); nil != err {
		log.Error("Failed to validate confirmVote, the owner organization identity is invalid", "err", err)
		return err
	}

	// validate voteOption
	if len(confirmVote.VoteOption) != 1 ||
		types.VoteOptionFromBytes(confirmVote.VoteOption) == types.VoteUnknown {
		return ctypes.ErrConfirmVoteOptionIsUnknown
	}

	return nil
}

// With subscriber
func (t *TwoPC) validateCommitMsg(commitMsg *types.CommitMsgWrap) error {

	proposalId := common.BytesToHash(commitMsg.ProposalId)
	if t.state.HasNotProposal(proposalId) {
		return ctypes.ErrProposalNotFound
	}

	// (On subscriber)
	// Now, current period should be `commitPeroid`,
	// But the subscriber still stay on `confirmPeriod`
	proposalState := t.state.ProposalStates(proposalId)
	if proposalState.IsNotConfirmPeriod() {
		return ctypes.ErrCommitMsgIllegal
	}

	// earlier commitMsg is invalid msg
	if proposalState.PeriodStartTime >= commitMsg.CreateAt {
		return ctypes.ErrCommitMsgIllegal
	}
	now := uint64(time.Now().UnixNano())
	if commitMsg.CreateAt >= now {
		return ctypes.ErrCommitMsgInTheFuture
	}

	// Verify the signature
	_, err := t.verifyMsgSigned(commitMsg.Owner.NodeId, commitMsg.SealHash().Bytes(), commitMsg.Signature())
	if err != nil {
		return err
	}

	// validate the owner
	if err := t.validateOrganizationIdentity(commitMsg.Owner); nil != err {
		log.Error("Failed to validate commitMsg, the owner organization identity is invalid", "err", err)
		return err
	}

	return nil
}

// Handle the prepareMsg from the task pulisher peer (on subscriber)
func (t *TwoPC) onPrepareMsg(prepareMsg *types.PrepareMsgWrap) error {
	proposal, err := t.fetchTaskFromPrepareMsg(prepareMsg)
	if nil != err {
		return err
	}
	if t.state.HasNotProposal(proposal.ProposalId) {
		return ctypes.ErrProposalAlreadyProcessed
	}
	proposalState := ctypes.NewProposalState(proposal.ProposalId, proposal.CreateAt)
	task := proposal.ScheduleTask
	t.state.AddProposalState(proposalState)
	t.addRecvTask(task)

	if err := t.validateRecvTask(task); nil != err {
		// clean some data
		t.state.DelProposalState(proposal.ProposalId)
		t.delRecvTask(task.TaskId)
		return err
	}

	resultCh := make(chan *types.ScheduleResult)
	t.replayTaskCh <- &types.ScheduleTaskWrap{
		Role:     types.TaskRoleFromBytes(prepareMsg.TaskOption.TaskRole),
		Task:     task,
		ResultCh: resultCh,
	}
	result := <-resultCh
	_ = result
	// todo  检查自己的资源是否符合 预期

	// todo 检查自己是否 意愿透出 yes 票

	// todo 调用 handler 投出自己的一票

	return nil
}
func (t *TwoPC) onPrepareVote(prepareVote *types.PrepareVoteWrap) error { return nil }
func (t *TwoPC) onConfirmMsg(confirmMsg *types.ConfirmMsgWrap) error    { return nil }
func (t *TwoPC) onConfirmVote(confirmVote *types.ConfirmVoteWrap) error { return nil }
func (t *TwoPC) onCommitMsg(cimmitMsg *types.CommitMsgWrap) error       { return nil }

func (t *TwoPC) validateOrganizationIdentity(identityInfo *pb.TaskOrganizationIdentityInfo) error {
	if "" == string(identityInfo.Name) {
		return ctypes.ErrOrganizationIdentity
	}
	_, err := p2p.BytesID(identityInfo.NodeId)
	if nil != err {
		return ctypes.ErrOrganizationIdentity
	}
	has, err := t.dataCenter.HasIdentity(&types.NodeAlias{
		Name:       string(identityInfo.Name),
		NodeId:     string(identityInfo.NodeId),
		IdentityId: string(identityInfo.IdentityId),
	})
	if nil != err {
		return fmt.Errorf("Failed to validate organization identity from all identity list")
	}
	if !has {
		return ctypes.ErrOrganizationIdentity
	}

	return nil
}

func (t *TwoPC) verifyMsgSigned(nodeId []byte, m []byte, sig []byte) (bool, error) {
	// Verify the signature
	if len(sig) < types.MsgSignLength {
		return false, ctypes.ErrMsgSignInvalid
	}

	recPubKey, err := crypto.Ecrecover(m, sig)
	if err != nil {
		return false, err
	}

	ownerNodeId, err := p2p.BytesID(nodeId)
	if nil != err {
		return false, ctypes.ErrMsgOwnerNodeIdInvalid
	}
	ownerPubKey, err := ownerNodeId.Pubkey()
	if nil != err {
		return false, ctypes.ErrMsgOwnerNodeIdInvalid
	}
	pbytes := elliptic.Marshal(ownerPubKey.Curve, ownerPubKey.X, ownerPubKey.Y)
	if !bytes.Equal(pbytes, recPubKey) {
		return false, ctypes.ErrMsgSignInvalid
	}
	return true, nil
}

func (t *TwoPC) verifySelfSigned(m []byte, sig []byte) bool {
	recPubKey, err := crypto.Ecrecover(m, sig)
	if err != nil {
		return false
	}
	pubKey := t.config.Option.NodePriKey.PublicKey
	pbytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
	return bytes.Equal(pbytes, recPubKey)
}

func (t *TwoPC) fetchTaskFromPrepareMsg(prepareMsg *types.PrepareMsgWrap) (*types.ProposalTask, error) {

	return nil, nil
}

func (t *TwoPC) isProcessingTask(taskId string) bool {
	if _, ok := t.sendTasks[taskId]; ok {
		return true
	}
	if _, ok := t.recvTasks[taskId]; ok {
		return true
	}
	return false
}

//// VerifyHeader verify block's header.
//func (vp *ValidatorPool) VerifyHeader(header *types.Header) error {
//	_, err := crypto.Ecrecover(header.SealHash().Bytes(), header.Signature())
//	if err != nil {
//		return err
//	}
//	// todo: need confirmed.
//	return vp.agency.VerifyHeader(header, nil)
//}

func (t *TwoPC) validateRecvTask(task *types.ScheduleTask) error {

	return nil
}
func (t *TwoPC) addSendTask(task *types.ScheduleTask) {
	t.sendTasks[task.TaskId] = task
}
func (t *TwoPC) addRecvTask(task *types.ScheduleTask) {
	t.recvTasks[task.TaskId] = task
}
func (t *TwoPC) delSendTask(taskId string) {
	delete(t.sendTasks, taskId)
}
func (t *TwoPC) delRecvTask(taskId string) {
	delete(t.recvTasks, taskId)
}
