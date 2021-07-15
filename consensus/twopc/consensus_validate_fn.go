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
	"github.com/libp2p/go-libp2p-core/peer"
	"time"
)



// With subscriber
func (t *TwoPC) validatePrepareMsg(pid peer.ID, prepareMsg *types.PrepareMsgWrap) error {
	//proposalId := common.BytesToHash(prepareMsg.ProposalId)
	//if t.state.HasProposal(proposalId) {
	//	return ctypes.ErrProposalAlreadyProcessed
	//}
	//// Now, we has not  proposalState (on subscriber)
	//
	//now := uint64(time.Now().UnixNano())
	//if prepareMsg.CreateAt >= now {
	//	return ctypes.ErrProposalInTheFuture
	//}
	//
	//if len(prepareMsg.TaskRole) != 1 ||
	//	types.TaskRoleFromBytes(prepareMsg.TaskRole) == types.TaskRoleUnknown {
	//	return ctypes.ErrPrososalTaskRoleIsUnknown
	//}
	//taskId := string(prepareMsg.TaskOption.TaskId)
	//
	//if t.isProcessingTask(taskId) {
	//	return ctypes.ErrPrososalTaskIsProcessed
	//}
	//
	//// Verify the signature
	//_, err := t.verifyMsgSigned(prepareMsg.TaskOption.Owner.NodeId, prepareMsg.SealHash().Bytes(), prepareMsg.Signature())
	//if err != nil {
	//	return err
	//}
	//
	//// validate the owner
	//if err := t.validateOrganizationIdentity(prepareMsg.TaskOption.Owner); nil != err {
	//	log.Error("Failed to validate prepareMsg, the owner organization identity is invalid", "err", err)
	//	return err
	//}
	//
	//// validate the algo supplier
	//if err := t.validateOrganizationIdentity(prepareMsg.TaskOption.AlgoSupplier); nil != err {
	//	log.Error("Failed to validate prepareMsg, the algoSupplier organization identity is invalid", "err", err)
	//	return err
	//}
	//
	//// validate data suppliers
	//if len(prepareMsg.TaskOption.DataSupplier) == 0 {
	//	return fmt.Errorf("Failed to validate prepareMsg, the data suppliers is empty")
	//}
	//for _, supplier := range prepareMsg.TaskOption.DataSupplier {
	//	if err := t.validateOrganizationIdentity(supplier.MemberInfo); nil != err {
	//		log.Error("Failed to validate prepareMsg, the dataSupplier organization identity is invalid", "err", err)
	//		return err
	//	}
	//}
	//
	//// validate power suppliers
	//if len(prepareMsg.TaskOption.PowerSupplier) == 0 {
	//	return fmt.Errorf("Failed to validate prepareMsg, the data suppliers is empty")
	//}
	//powerSuppliers := make(map[string]struct{}, len(prepareMsg.TaskOption.PowerSupplier))
	//for _, supplier := range prepareMsg.TaskOption.PowerSupplier {
	//	if err := t.validateOrganizationIdentity(supplier.MemberInfo); nil != err {
	//		log.Error("Failed to validate prepareMsg, the powerSupplier organization identity is invalid", "err", err)
	//		return err
	//	}
	//	powerSuppliers[string(supplier.MemberInfo.IdentityId)] = struct{}{}
	//}
	//
	//// validate receicers
	//if len(prepareMsg.TaskOption.Receivers) == 0 {
	//	return fmt.Errorf("Failed to validate prepareMsg, the data suppliers is empty")
	//}
	//
	//for _, supplier := range prepareMsg.TaskOption.Receivers {
	//	if err := t.validateOrganizationIdentity(supplier.MemberInfo); nil != err {
	//		log.Error("Failed to validate prepareMsg, the powerSupplier organization identity is invalid", "err", err)
	//		return err
	//	}
	//	receiverIdentityId := string(supplier.MemberInfo.IdentityId)
	//	for _, provider := range supplier.Providers {
	//		providerIdentityId := string(provider.IdentityId)
	//		if _, ok := powerSuppliers[providerIdentityId]; !ok {
	//			return fmt.Errorf("Has invalid provider with receiver", "providerIdentityId", providerIdentityId, "receiverIndentityId", receiverIdentityId)
	//		}
	//	}
	//}
	//
	//// validate OperationCost
	//if 0 == prepareMsg.TaskOption.OperationCost.CostProcessor ||
	//	0 == prepareMsg.TaskOption.OperationCost.CostMem ||
	//	0 == prepareMsg.TaskOption.OperationCost.CostBandwidth ||
	//	0 == prepareMsg.TaskOption.OperationCost.Duration {
	//	return ctypes.ErrProposalTaskOperationCostInvalid
	//}
	//
	//// validate contractCode
	//if 0 == len(prepareMsg.TaskOption.CalculateContractCode) {
	//	return ctypes.ErrProposalTaskCalculateContractCodeEmpty
	//}
	//
	//// validate task create time
	//if prepareMsg.TaskOption.CreateAt >= prepareMsg.CreateAt {
	//	return ctypes.ErrProposalParamsInvalid
	//}
	return nil
}

// With publisher
func (t *TwoPC) validatePrepareVote(pid peer.ID, prepareVote *types.PrepareVoteWrap) error {
	proposalId := common.BytesToHash(prepareVote.ProposalId)
	if t.state.HasNotProposal(proposalId) {
		return ctypes.ErrProposalNotFound
	}
	// (On publisher)
	// Now, we has proposalState, current period is `preparePeriod`
	// the PeriodStartTime is not zero, But the PeriodEndTime is zero
	proposalState := t.state.GetProposalState(proposalId)
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
func (t *TwoPC) validateConfirmMsg(pid peer.ID, confirmMsg *types.ConfirmMsgWrap) error {

	proposalId := common.BytesToHash(confirmMsg.ProposalId)
	if t.state.HasNotProposal(proposalId) {
		return ctypes.ErrProposalNotFound
	}

	// (On subscriber)
	// Now, current period should be `confirmPeriod`,
	// If confirmMsg is first epoch, so still stay preparePeriod
	// If confirmMsg is second epoch, so stay confirmPeriod
	proposalState := t.state.GetProposalState(proposalId)
	if proposalState.IsNotPreparePeriod() || proposalState.IsNotConfirmPeriod() {
		return ctypes.ErrConfirmMsgIllegal
	}
	// When comfirm first epoch
	//TODO: ConfirmEpochFirst and ConfirmEpochSecond is missing.
	/*if proposalState.IsPreparePeriod() && ctypes.ConfirmEpochFirst.Uint64() != confirmMsg.Epoch {
		return ctypes.ErrConfirmMsgIllegal
	}
	// When comfirm second epoch
	if proposalState.IsConfirmPeriod() && ctypes.ConfirmEpochSecond.Uint64() != confirmMsg.Epoch {
		return ctypes.ErrConfirmMsgIllegal
	}*/

	// earlier confirmMsg is invalid msg
	if proposalState.PeriodStartTime >= confirmMsg.CreateAt {
		return ctypes.ErrConfirmMsgIllegal
	}
	now := uint64(time.Now().UnixNano())
	if confirmMsg.CreateAt >= now {
		return ctypes.ErrConfirmMsgInTheFuture
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
func (t *TwoPC) validateConfirmVote(pid peer.ID, confirmVote *types.ConfirmVoteWrap) error {
	proposalId := common.BytesToHash(confirmVote.ProposalId)
	if t.state.HasNotProposal(proposalId) {
		return ctypes.ErrProposalNotFound
	}

	// (On publisher)
	// Now, current period should be `confirmPeriod`,
	// No matter first epoch of confirmVote or second epoch of confirmVote
	// it is always confirmPeriod.
	proposalState := t.state.GetProposalState(proposalId)
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
func (t *TwoPC) validateCommitMsg(pid peer.ID, commitMsg *types.CommitMsgWrap) error {

	proposalId := common.BytesToHash(commitMsg.ProposalId)
	if t.state.HasNotProposal(proposalId) {
		return ctypes.ErrProposalNotFound
	}

	// (On subscriber)
	// Now, current period should be `commitPeroid`,
	// But the subscriber still stay on `confirmPeriod`
	proposalState := t.state.GetProposalState(proposalId)
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

// With publisher
func (t *TwoPC) validateTaskResultMsg(pid peer.ID, taskResultMsg *types.TaskResultMsgWrap) error {
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






//// VerifyHeader verify block's header.
//func (vp *ValidatorPool) VerifyHeader(header *types.Header) error {
//	_, err := crypto.Ecrecover(header.SealHash().Bytes(), header.Signature())
//	if err != nil {
//		return err
//	}
//	// todo: need confirmed.
//	return vp.agency.VerifyHeader(header, nil)
//}

// TODO 需要实现
func (t *TwoPC) validateRecvTask(task *types.Task) error {

	return nil
}