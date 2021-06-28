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

type TwoPC struct {
	config *Config
	Errs   []error
	p2p    p2p.P2P
	// The task being processed by myself  (taskId -> task)
	sendTasks map[string]*types.ScheduleTask
	// The task processing  that received someone else (taskId -> task)
	recvTasks map[string]*types.ScheduleTask
	// Proposal being processed
	runningProposals map[common.Hash]string
}

func New(conf *Config) *TwoPC {

	t := &TwoPC{
		config: conf,
		Errs:   make([]error, 0),
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
	switch msg.(type) {
	case *types.PrepareMsgWrap:
		return t.validatePrepareMsg(msg.(*types.PrepareMsgWrap))
	case *types.PrepareVoteWrap:

	case *types.ConfirmMsgWrap:

	case *types.ConfirmVoteWrap:

	case *types.CommitMsgWrap:

	default:
		return fmt.Errorf("Unknown the 2pc msg type")

	}

	return nil
}

func (t *TwoPC) OnConsensusMsg(msg types.ConsensusMsg) error {

	switch msg.(type) {
	case *types.PrepareMsgWrap:

	case *types.PrepareVoteWrap:

	case *types.ConfirmMsgWrap:

	case *types.ConfirmVoteWrap:

	case *types.CommitMsgWrap:

	default:
		return fmt.Errorf("Unknown the 2pc msg type")

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

func (t *TwoPC) hasProposal(proposalId common.Hash) bool {
	if _, ok := t.runningProposals[proposalId]; ok {
		return true
	}
	return false
}

func (t *TwoPC) hasNotProposal(proposalId common.Hash) bool {
	return !t.hasProposal(proposalId)
}

func (t *TwoPC) validatePrepareMsg(prepareMsg *types.PrepareMsgWrap) error {
	proposalId := common.BytesToHash(prepareMsg.ProposalId)
	if t.hasProposal(proposalId) {
		return ctypes.ErrProposalAlreadyProcessed
	}
	now := uint64(time.Now().UnixNano())
	if prepareMsg.CreateAt >= now {
		return ctypes.ErrProposalInTheFuture
	}

	if TaskRoleFromBytes(prepareMsg.TaskOption.TaskRole) == Unknown {
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
	if len(prepareMsg.Signature()) < MsgSignLength {
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
	if err := t.validateOrganizationIdentity(prepareMsg.TaskOption.Owner); nil != err {
		log.Error("Failed to validate prepareMsg, the owner organization identity is invalid", "err", err)
		return err
	}
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

func (t *TwoPC) validateOrganizationIdentity(identityInfo *pb.TaskOrganizationIdentityInfo) error {
	if "" == string(identityInfo.Name) {
		return ctypes.ErrOrganizationIdentityNameEmpty
	}
	_, err := p2p.BytesID(identityInfo.NodeId)
	if nil != err {
		return ctypes.ErrOrganizationIdentityNodeIdInvalid
	}
	// TODO ...

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

//// VerifyHeader verify block's header.
//func (vp *ValidatorPool) VerifyHeader(header *types.Header) error {
//	_, err := crypto.Ecrecover(header.SealHash().Bytes(), header.Signature())
//	if err != nil {
//		return err
//	}
//	// todo: need confirmed.
//	return vp.agency.VerifyHeader(header, nil)
//}
