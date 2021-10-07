package twopc
//
//import (
//	"bytes"
//	"crypto/elliptic"
//	"fmt"
//	"github.com/RosettaFlow/Carrier-Go/common"
//	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
//	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
//	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
//	"github.com/RosettaFlow/Carrier-Go/p2p"
//	"github.com/RosettaFlow/Carrier-Go/types"
//	"github.com/ethereum/go-ethereum/crypto"
//	"github.com/libp2p/go-libp2p-core/peer"
//)
//
//// With subscriber
//func (t *Twopc) validatePrepareMsg(pid peer.ID, prepareMsg *types.PrepareMsgWrap) error {
//	//msg, err := fetchPrepareMsg(prepareMsg)
//	//if nil != err {
//	//	return err
//	//}
//	//if t.state.HasProposal(proposalId) {
//	//	return ctypes.ErrProposalAlreadyProcessed
//	//}
//	//// Now, we has not  proposalState (on subscriber)
//	//
//	//now := uint64(timeutils.UnixMsec())
//	//if prepareMsg.GetCreateAt >= now {
//	//	return ctypes.ErrProposalInTheFuture
//	//}
//	//
//	//if len(prepareMsg.TaskRole) != 1 ||
//	//	types.TaskRoleFromBytes(prepareMsg.TaskRole) == types.TaskRoleUnknown {
//	//	return ctypes.ErrPrososalTaskRoleIsUnknown
//	//}
//	//taskId := string(prepareMsg.TaskOption.GetTaskId)
//	//
//	//if t.isProcessingTask(taskId) {
//	//	return ctypes.ErrPrososalTaskIsProcessed
//	//}
//	//
//	//// Verify the signature
//	//_, err := t.verifyMsgSigned(prepareMsg.TaskOption.GetSender.NodeId, prepareMsg.SealHash().Bytes(), prepareMsg.Signature())
//	//if err != nil {
//	//	return err
//	//}
//	//
//	//// validate the owner
//	//if err := t.validateOrganizationIdentity(prepareMsg.TaskOption.GetSender); nil != err {
//	//	log.Error("Failed to validate prepareMsg, the owner organization identity is invalid", "err", err)
//	//	return err
//	//}
//	//
//	//// validate the algo supplier
//	//if err := t.validateOrganizationIdentity(prepareMsg.TaskOption.AlgoSupplier); nil != err {
//	//	log.Error("Failed to validate prepareMsg, the algoSupplier organization identity is invalid", "err", err)
//	//	return err
//	//}
//	//
//	//// validate data suppliers
//	//if len(prepareMsg.TaskOption.DataSupplier) == 0 {
//	//	return fmt.Errorf("Failed to validate prepareMsg, the data suppliers is empty")
//	//}
//	//for _, supplier := range prepareMsg.TaskOption.DataSupplier {
//	//	if err := t.validateOrganizationIdentity(supplier.MemberInfo); nil != err {
//	//		log.Error("Failed to validate prepareMsg, the dataSupplier organization identity is invalid", "err", err)
//	//		return err
//	//	}
//	//}
//	//
//	//// validate power suppliers
//	//if len(prepareMsg.TaskOption.PowerSupplier) == 0 {
//	//	return fmt.Errorf("Failed to validate prepareMsg, the data suppliers is empty")
//	//}
//	//powerSuppliers := make(map[string]struct{}, len(prepareMsg.TaskOption.PowerSupplier))
//	//for _, supplier := range prepareMsg.TaskOption.PowerSupplier {
//	//	if err := t.validateOrganizationIdentity(supplier.MemberInfo); nil != err {
//	//		log.Error("Failed to validate prepareMsg, the powerSupplier organization identity is invalid", "err", err)
//	//		return err
//	//	}
//	//	powerSuppliers[string(supplier.MemberInfo.IdentityId)] = struct{}{}
//	//}
//	//
//	//// validate receicers
//	//if len(prepareMsg.TaskOption.Receivers) == 0 {
//	//	return fmt.Errorf("Failed to validate prepareMsg, the data suppliers is empty")
//	//}
//	//
//	//for _, supplier := range prepareMsg.TaskOption.Receivers {
//	//	if err := t.validateOrganizationIdentity(supplier.MemberInfo); nil != err {
//	//		log.Error("Failed to validate prepareMsg, the powerSupplier organization identity is invalid", "err", err)
//	//		return err
//	//	}
//	//	receiverIdentityId := string(supplier.MemberInfo.IdentityId)
//	//	for _, provider := range supplier.Providers {
//	//		providerIdentityId := string(provider.IdentityId)
//	//		if _, ok := powerSuppliers[providerIdentityId]; !ok {
//	//			return fmt.Errorf("Has invalid provider with receiver", "providerIdentityId", providerIdentityId, "receiverIndentityId", receiverIdentityId)
//	//		}
//	//	}
//	//}
//	//
//	//// validate OperationCost
//	//if 0 == prepareMsg.TaskOption.OperationCost.CostProcessor ||
//	//	0 == prepareMsg.TaskOption.OperationCost.CostMem ||
//	//	0 == prepareMsg.TaskOption.OperationCost.CostBandwidth ||
//	//	0 == prepareMsg.TaskOption.OperationCost.Duration {
//	//	return ctypes.ErrProposalTaskOperationCostInvalid
//	//}
//	//
//	//// validate contractCode
//	//if 0 == len(prepareMsg.TaskOption.CalculateContractCode) {
//	//	return ctypes.ErrProposalTaskCalculateContractCodeEmpty
//	//}
//	//
//	//// validate task create time
//	//if prepareMsg.TaskOption.GetCreateAt >= prepareMsg.GetCreateAt {
//	//	return ctypes.ErrProposalParamsInvalid
//	//}
//	return nil
//}
//
//// With publisher
//func (t *Twopc) validatePrepareVote(pid peer.ID, prepareVote *types.PrepareVoteWrap) error {
//
//	vote := fetchPrepareVote(prepareVote)
//
//	proposalId := vote.MsgOption.ProposalId
//	partyId := vote.MsgOption.ReceiverPartyId
//	if t.state.HasNotOrgProposal(proposalId, partyId) {
//		return fmt.Errorf("%s validatePrepareVote", ctypes.ErrProposalNotFound)
//	}
//	// (On publisher)
//	// Now, we has proposalState, current period is `preparePeriod`
//	// the PeriodStartTime is not zero, But the PeriodEndTime is zero
//	orgProposalState := t.mustGetOrgProposalState(proposalId, partyId)
//	if orgProposalState.IsNotPreparePeriod() {
//		return ctypes.ErrPrepareVoteIllegal
//	}
//
//	// earlier vote is invalid vote
//	if orgProposalState.PeriodStartTime >= vote.GetCreateAt {
//		return ctypes.ErrPrepareVoteIllegal
//	}
//
//	now := uint64(timeutils.UnixMsec())
//	if vote.GetCreateAt >= now {
//		return ctypes.ErrPrepareVoteInTheFuture
//	}
//
//	if vote.VoteOption == types.VoteUnknown && vote.PeerInfoEmpty() {
//		return ctypes.ErrPrososalTaskRoleIsUnknown
//	}
//
//	// Verify the signature
//	_, err := t.verifyMsgSigned(vote.MsgOption.GetSender.GetNodeId(), prepareVote.SealHash().Bytes(), prepareVote.Signature())
//	if err != nil {
//		return err
//	}
//
//	// validate the owner
//	if err := t.validateOrganizationIdentity(vote.MsgOption.GetSender); nil != err {
//		log.Error("Failed to validate prepareVote, the owner organization identity is invalid", "err", err)
//		return err
//	}
//
//	// validate voteOption
//	if len(prepareVote.VoteOption) != 1 ||
//		types.VoteOptionFromBytes(prepareVote.VoteOption) == types.VoteUnknown {
//		return ctypes.ErrPrepareVoteOptionIsUnknown
//	}
//	if 0 == len(prepareVote.PeerInfo.Ip) || 0 == len(prepareVote.PeerInfo.Port) {
//		return ctypes.ErrPrepareVoteParamsInvalid
//	}
//
//	return nil
//}
//
//// With subscriber
//func (t *Twopc) validateConfirmMsg(pid peer.ID, confirmMsg *types.ConfirmMsgWrap) error {
//
//	proposalId := common.BytesToHash(confirmMsg.ProposalId)
//	if t.state.HasNotOrgProposal(proposalId) {
//		return fmt.Errorf("%s validateConfirmMsg", ctypes.ErrProposalNotFound)
//	}
//
//	// (On subscriber)
//	// Now, current period should be `confirmPeriod`,
//	// If confirmMsg is first epoch, so still stay preparePeriod
//	// If confirmMsg is second epoch, so stay confirmPeriod
//	proposalState := t.state.GetProposalState(proposalId)
//	if proposalState.IsNotPreparePeriod() || proposalState.IsNotConfirmPeriod() {
//		return ctypes.ErrConfirmMsgIllegal
//	}
//	// When comfirm first epoch
//	//TODO: ConfirmEpochFirst and ConfirmEpochSecond is missing.
//	/*if proposalState.IsPreparePeriod() && ctypes.ConfirmEpochFirst.Uint64() != confirmMsg.Epoch {
//		return ctypes.ErrConfirmMsgIllegal
//	}
//	// When comfirm second epoch
//	if proposalState.IsConfirmPeriod() && ctypes.ConfirmEpochSecond.Uint64() != confirmMsg.Epoch {
//		return ctypes.ErrConfirmMsgIllegal
//	}*/
//
//	// earlier confirmMsg is invalid msg
//	if proposalState.PeriodStartTime >= confirmMsg.GetCreateAt {
//		return ctypes.ErrConfirmMsgIllegal
//	}
//	now := uint64(timeutils.UnixMsec())
//	if confirmMsg.GetCreateAt >= now {
//		return ctypes.ErrConfirmMsgInTheFuture
//	}
//
//	// Verify the signature
//	_, err := t.verifyMsgSigned(confirmMsg.GetSender.NodeId, confirmMsg.SealHash().Bytes(), confirmMsg.Signature())
//	if err != nil {
//		return err
//	}
//
//	// validate the owner
//	if err := t.validateOrganizationIdentity(confirmMsg.GetSender); nil != err {
//		log.Error("Failed to validate confirmMsg, the owner organization identity is invalid", "err", err)
//		return err
//	}
//
//	return nil
//}
//
//// With publisher
//func (t *Twopc) validateConfirmVote(pid peer.ID, confirmVote *types.ConfirmVoteWrap) error {
//	proposalId := common.BytesToHash(confirmVote.ProposalId)
//	if t.state.HasNotOrgProposal(proposalId) {
//		return fmt.Errorf("%s validateConfirmVote", ctypes.ErrProposalNotFound)
//	}
//
//	// (On publisher)
//	// Now, current period should be `confirmPeriod`,
//	// No matter first epoch of confirmVote or second epoch of confirmVote
//	// it is always confirmPeriod.
//	proposalState := t.state.GetProposalState(proposalId)
//	if proposalState.IsNotConfirmPeriod() {
//		return ctypes.ErrConfirmVoteIllegal
//	}
//
//	// earlier confirmVote is invalid vote
//	if proposalState.PeriodStartTime >= confirmVote.GetCreateAt {
//		return ctypes.ErrConfirmVoteIllegal
//	}
//	now := uint64(timeutils.UnixMsec())
//	if confirmVote.GetCreateAt >= now {
//		return ctypes.ErrConfirmVoteInTheFuture
//	}
//
//	if types.TaskRoleFromBytes(confirmVote.TaskRole) == types.TaskRoleUnknown {
//		return ctypes.ErrPrososalTaskRoleIsUnknown
//	}
//
//	// Verify the signature
//	_, err := t.verifyMsgSigned(confirmVote.GetSender.NodeId, confirmVote.SealHash().Bytes(), confirmVote.Signature())
//	if err != nil {
//		return err
//	}
//
//	// validate the owner
//	if err := t.validateOrganizationIdentity(confirmVote.GetSender); nil != err {
//		log.Error("Failed to validate confirmVote, the owner organization identity is invalid", "err", err)
//		return err
//	}
//
//	// validate voteOption
//	if len(confirmVote.VoteOption) != 1 ||
//		types.VoteOptionFromBytes(confirmVote.VoteOption) == types.VoteUnknown {
//		return ctypes.ErrConfirmVoteOptionIsUnknown
//	}
//
//	return nil
//}
//
//// With subscriber
//func (t *Twopc) validateCommitMsg(pid peer.ID, commitMsg *types.CommitMsgWrap) error {
//
//	proposalId := common.BytesToHash(commitMsg.ProposalId)
//	if t.state.HasNotOrgProposal(proposalId) {
//		return fmt.Errorf("%s validateCommitMsg", ctypes.ErrProposalNotFound)
//	}
//
//	// (On subscriber)
//	// Now, current period should be `commitPeroid`,
//	// But the subscriber still stay on `confirmPeriod`
//	proposalState := t.state.GetProposalState(proposalId)
//	if proposalState.IsNotConfirmPeriod() {
//		return ctypes.ErrCommitMsgIllegal
//	}
//
//	// earlier commitMsg is invalid msg
//	if proposalState.PeriodStartTime >= commitMsg.GetCreateAt {
//		return ctypes.ErrCommitMsgIllegal
//	}
//	now := uint64(timeutils.UnixMsec())
//	if commitMsg.GetCreateAt >= now {
//		return ctypes.ErrCommitMsgInTheFuture
//	}
//
//	// Verify the signature
//	_, err := t.verifyMsgSigned(commitMsg.GetSender.NodeId, commitMsg.SealHash().Bytes(), commitMsg.Signature())
//	if err != nil {
//		return err
//	}
//
//	// validate the owner
//	if err := t.validateOrganizationIdentity(commitMsg.GetSender); nil != err {
//		log.Error("Failed to validate commitMsg, the owner organization identity is invalid", "err", err)
//		return err
//	}
//
//	return nil
//}
//
//func (t *Twopc) validateOrganizationIdentity(identityInfo *apicommonpb.TaskOrganization) error {
//	if "" == identityInfo.GetNodeName() {
//		return ctypes.ErrOrganizationIdentity
//	}
//	_, err := p2p.HexID(identityInfo.GetNodeId())
//	if nil != err {
//		return ctypes.ErrOrganizationIdentity
//	}
//	has, err := t.resourceMng.GetDB().HasIdentity(&apicommonpb.Organization{
//		NodeName:   identityInfo.GetNodeName(),
//		NodeId:     identityInfo.GetNodeId(),
//		IdentityId: identityInfo.QueryIdentityId(),
//	})
//	if nil != err {
//		return fmt.Errorf("Failed to validate organization identity from all identity list")
//	}
//	if !has {
//		return ctypes.ErrOrganizationIdentity
//	}
//
//	return nil
//}
//
//func (t *Twopc) verifyMsgSigned(nodeId string, m []byte, sig []byte) (bool, error) {
//	// Verify the signature
//	if len(sig) < types.MsgSignLength {
//		return false, ctypes.ErrMsgSignInvalid
//	}
//
//	recPubKey, err := crypto.Ecrecover(m, sig)
//	if err != nil {
//		return false, err
//	}
//
//	ownerNodeId, err := p2p.HexID(nodeId)
//	if nil != err {
//		return false, ctypes.ErrMsgOwnerNodeIdInvalid
//	}
//	ownerPubKey, err := ownerNodeId.Pubkey()
//	if nil != err {
//		return false, ctypes.ErrMsgOwnerNodeIdInvalid
//	}
//	pbytes := elliptic.Marshal(ownerPubKey.Curve, ownerPubKey.X, ownerPubKey.Y)
//	if !bytes.Equal(pbytes, recPubKey) {
//		return false, ctypes.ErrMsgSignInvalid
//	}
//	return true, nil
//}
//
//func (t *Twopc) verifySelfSigned(m []byte, sig []byte) bool {
//	recPubKey, err := crypto.Ecrecover(m, sig)
//	if err != nil {
//		return false
//	}
//	pubKey := t.config.Option.NodePriKey.PublicKey
//	pbytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
//	return bytes.Equal(pbytes, recPubKey)
//}
//
////// VerifyHeader verify block's header.
////func (vp *ValidatorPool) VerifyHeader(header *types.Header) error {
////	_, err := crypto.Ecrecover(header.SealHash().Bytes(), header.Signature())
////	if err != nil {
////		return err
////	}
////	// todo: need confirmed.
////	return vp.agency.VerifyHeader(header, nil)
////}
//
//func (t *Twopc) validateRecvTask(task *types.Task) error {
//
//	return nil
//}
