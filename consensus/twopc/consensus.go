package twopc

import (
	"bytes"
	"crypto/elliptic"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
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
	Errs    []error
	p2p     p2p.P2P
	peerSet *ctypes.PeerSet

	// The task being processed by myself  (taskId -> task)
	sendTasks map[string]*types.ScheduleTask
	// The task processing  that received someone else (taskId -> task)
	recvTasks map[string]*types.ScheduleTask
	state    *state

	dataCenter DataCenter
}

func New(conf *Config) *TwoPC {

	t := &TwoPC{
		config: conf,
		Errs:   make([]error, 0),
		state: newState(),
	}

	return t
}

func (t *TwoPC) OnPrepare(task *types.ScheduleTask) error {
	return nil
}
func (t *TwoPC) OnStart(task *types.ScheduleTask, result chan<- *types.ScheduleResult) error {

	return nil
}

func (t *TwoPC) ValidateConsensusMsg(msg types.ConsensusMsg) error {
	if nil == msg {
		return fmt.Errorf("Failed to validate 2pc consensus msg, the msg is nil")
	}
	switch msg.(type) {
	case *types.PrepareMsgWrap:
		return t.validatePrepareMsg(msg.(*types.PrepareMsgWrap))
	case *types.PrepareVoteWrap:
		return t.validatePrepareVote(msg.(*types.PrepareVoteWrap))
	case *types.ConfirmMsgWrap:
		return t.validateConfirmMsg(msg.(*types.ConfirmMsgWrap))
	case *types.ConfirmVoteWrap:
		return t.validateConfirmVote(msg.(*types.ConfirmVoteWrap))
	case *types.CommitMsgWrap:
		return t.validateCommitMsg(msg.(*types.CommitMsgWrap))
	default:
		return fmt.Errorf("TaskRoleUnknown the 2pc msg type")
	}
}

func (t *TwoPC) OnConsensusMsg(msg types.ConsensusMsg) error {

	switch msg.(type) {
	case *types.PrepareMsgWrap:
		//task, err := t.fetchTaskFromPrepareMsg(msg.(*types.PrepareMsgWrap))
		//if nil != err {
		//
		//}

	case *types.PrepareVoteWrap:

	case *types.ConfirmMsgWrap:

	case *types.ConfirmVoteWrap:

	case *types.CommitMsgWrap:

	default:
		return fmt.Errorf("TaskRoleUnknown the 2pc msg type")

	}

	return nil
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

func (t *TwoPC) OnPrepareMsg(proposal *ctypes.PrepareMsg) error {

	return nil
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
		ctypes.TaskRoleFromBytes(prepareMsg.TaskOption.TaskRole) == ctypes.TaskRoleUnknown {
		return ctypes.ErrPrososalTaskRoleIsUnknown
	}
	taskId := string(prepareMsg.TaskOption.TaskId)
	if _, ok := t.sendTasks[taskId]; ok {
		return ctypes.ErrPrososalTaskIdBelongSendTask
	}
	if _, ok := t.recvTasks[taskId]; ok {
		return ctypes.ErrPrososalTaskIdBelongRecvTask
	}

	// Verify the signature
	if len(prepareMsg.Signature()) < ctypes.MsgSignLength {
		return ctypes.ErrPrepareMsgSignInvalid
	}

	recPubKey, err := crypto.Ecrecover(prepareMsg.SealHash().Bytes(), prepareMsg.Signature())
	if err != nil {
		return err
	}
	proposalOwnerNodeId, err := p2p.BytesID(prepareMsg.TaskOption.Owner.NodeId)
	if nil != err {
		return ctypes.ErrProposalOwnerNodeIdInvalid
	}
	ownerPubKey, err := proposalOwnerNodeId.Pubkey()
	if nil != err {
		return ctypes.ErrRecoverPubkeyFromProposalOwner
	}
	pbytes := elliptic.Marshal(ownerPubKey.Curve, ownerPubKey.X, ownerPubKey.Y)
	if !bytes.Equal(pbytes, recPubKey) {
		return ctypes.ErrProposalIllegal
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

	if ctypes.TaskRoleFromBytes(prepareVote.TaskRole) == ctypes.TaskRoleUnknown {
		return ctypes.ErrPrososalTaskRoleIsUnknown
	}

	// Verify the signature
	if len(prepareVote.Signature()) < ctypes.MsgSignLength {
		return ctypes.ErrPrepareVoteSignInvalid
	}

	recPubKey, err := crypto.Ecrecover(prepareVote.SealHash().Bytes(), prepareVote.Signature())
	if err != nil {
		return err
	}
	prepareVoteOwnerNodeId, err := p2p.BytesID(prepareVote.Owner.NodeId)
	if nil != err {
		return ctypes.ErrPrepareVoteOwnerNodeIdInvalid
	}
	ownerPubKey, err := prepareVoteOwnerNodeId.Pubkey()
	if nil != err {
		return ctypes.ErrRecoverPubkeyFromPrepareVoteOwner
	}
	pbytes := elliptic.Marshal(ownerPubKey.Curve, ownerPubKey.X, ownerPubKey.Y)
	if !bytes.Equal(pbytes, recPubKey) {
		return ctypes.ErrPrepareVoteIllegal
	}

	// validate the owner
	if err := t.validateOrganizationIdentity(prepareVote.Owner); nil != err {
		log.Error("Failed to validate prepareVote, the owner organization identity is invalid", "err", err)
		return err
	}

	// validate voteOption
	if len(prepareVote.VoteOption) != 1 ||
		ctypes.VoteOptionFromBytes(prepareVote.VoteOption) == ctypes.VoteUnknown {
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
		return ctypes.ErrConfirmMsgSignInvalid
	}
	// When comfirm first epoch
	if proposalState.IsPreparePeriod() && ctypes.ConfirmEpochFirst.Uint64() != confirmMsg.Epoch {
		return ctypes.ErrConfirmMsgSignInvalid
	}
	// When comfirm second epoch
	if proposalState.IsConfirmPeriod() && ctypes.ConfirmEpochSecond.Uint64() != confirmMsg.Epoch {
		return ctypes.ErrConfirmMsgSignInvalid
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
	if confirmMsg.Epoch == 0 || confirmMsg.Epoch > ctypes.MsgEpochMaxNumber {
		return ctypes.ErrConfirmMsgEpochInvalid
	}

	// Verify the signature
	if len(confirmMsg.Signature()) < ctypes.MsgSignLength {
		return ctypes.ErrConfirmMsgSignInvalid
	}

	recPubKey, err := crypto.Ecrecover(confirmMsg.SealHash().Bytes(), confirmMsg.Signature())
	if err != nil {
		return err
	}
	confirmMsgOwnerNodeId, err := p2p.BytesID(confirmMsg.Owner.NodeId)
	if nil != err {
		return ctypes.ErrConfirmMsgOwnerNodeIdInvalid
	}
	ownerPubKey, err := confirmMsgOwnerNodeId.Pubkey()
	if nil != err {
		return ctypes.ErrRecoverPubkeyFromConfirmMsgOwner
	}
	pbytes := elliptic.Marshal(ownerPubKey.Curve, ownerPubKey.X, ownerPubKey.Y)
	if !bytes.Equal(pbytes, recPubKey) {
		return ctypes.ErrConfirmMsgIllegal
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
		return ctypes.ErrConfirmVoteSignInvalid
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
	if confirmVote.Epoch == 0 || confirmVote.Epoch > ctypes.MsgEpochMaxNumber {
		return ctypes.ErrConfirmVoteEpochInvalid
	}


	if ctypes.TaskRoleFromBytes(confirmVote.TaskRole) == ctypes.TaskRoleUnknown {
		return ctypes.ErrPrososalTaskRoleIsUnknown
	}

	// Verify the signature
	if len(confirmVote.Signature()) < ctypes.MsgSignLength {
		return ctypes.ErrConfirmVoteSignInvalid
	}

	recPubKey, err := crypto.Ecrecover(confirmVote.SealHash().Bytes(), confirmVote.Signature())
	if err != nil {
		return err
	}
	confirmVoteOwnerNodeId, err := p2p.BytesID(confirmVote.Owner.NodeId)
	if nil != err {
		return ctypes.ErrConfirmVoteOwnerNodeIdInvalid
	}
	ownerPubKey, err := confirmVoteOwnerNodeId.Pubkey()
	if nil != err {
		return ctypes.ErrRecoverPubkeyFromConfirmVoteOwner
	}
	pbytes := elliptic.Marshal(ownerPubKey.Curve, ownerPubKey.X, ownerPubKey.Y)
	if !bytes.Equal(pbytes, recPubKey) {
		return ctypes.ErrConfirmVoteIllegal
	}

	// validate the owner
	if err := t.validateOrganizationIdentity(confirmVote.Owner); nil != err {
		log.Error("Failed to validate confirmVote, the owner organization identity is invalid", "err", err)
		return err
	}

	// validate voteOption
	if len(confirmVote.VoteOption) != 1 ||
		ctypes.VoteOptionFromBytes(confirmVote.VoteOption) == ctypes.VoteUnknown {
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
		return ctypes.ErrCommitMsgSignInvalid
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
	if len(commitMsg.Signature()) < ctypes.MsgSignLength {
		return ctypes.ErrCommitMsgSignInvalid
	}

	recPubKey, err := crypto.Ecrecover(commitMsg.SealHash().Bytes(), commitMsg.Signature())
	if err != nil {
		return err
	}
	commitMsgOwnerNodeId, err := p2p.BytesID(commitMsg.Owner.NodeId)
	if nil != err {
		return ctypes.ErrCommitMsgOwnerNodeIdInvalid
	}
	ownerPubKey, err := commitMsgOwnerNodeId.Pubkey()
	if nil != err {
		return ctypes.ErrRecoverPubkeyFromCommitMsgOwner
	}
	pbytes := elliptic.Marshal(ownerPubKey.Curve, ownerPubKey.X, ownerPubKey.Y)
	if !bytes.Equal(pbytes, recPubKey) {
		return ctypes.ErrCommitMsgIllegal
	}

	// validate the owner
	if err := t.validateOrganizationIdentity(commitMsg.Owner); nil != err {
		log.Error("Failed to validate commitMsg, the owner organization identity is invalid", "err", err)
		return err
	}

	return nil
}

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

func (t *TwoPC) verifySelfSigned(m []byte, sig []byte) bool {
	recPubKey, err := crypto.Ecrecover(m, sig)
	if err != nil {
		return false
	}
	pubKey := t.config.Option.NodePriKey.PublicKey
	pbytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
	return bytes.Equal(pbytes, recPubKey)
}

func (t *TwoPC) fetchTaskFromPrepareMsg(prepareMsg *types.PrepareMsgWrap) (*types.ScheduleTask, error) {

	return nil, nil
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
