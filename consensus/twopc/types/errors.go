package types

import "errors"

var (

	// Prepare
	ErrProposalIllegal                = errors.New("The proposal is illegal")
	ErrProposalAlreadyProcessed       = errors.New("The proposalId have been processing")
	ErrProposalNotFound               = errors.New("The proposal is not found")
	ErrProposalInTheFuture            = errors.New("The proposal is in the future")
	ErrPrepareMsgSignInvalid          = errors.New("The prepareMsg signature is invalid")
	ErrProposalOwnerNodeIdInvalid     = errors.New("The proposal's owner nodeId is invalid")
	ErrRecoverPubkeyFromProposalOwner = errors.New("Failed to recover pubkey from proposal owner")
	ErrPrososalTaskRoleIsUnknown      = errors.New("The proposal's task role is unknown")
	ErrPrososalTaskIdBelongSendTask   = errors.New("The proposal's taskId is belonging to send task")
	ErrPrososalTaskIdBelongRecvTask   = errors.New("The proposal's taskId is belonging to recv task")

	// prepareVote
	ErrPrepareVoteInTheFuture     = errors.New("The prepareVote is in the future")
	ErrPrepareVoteOwnerNodeIdInvalid     = errors.New("The prepareVote's owner nodeId is invalid")
	ErrRecoverPubkeyFromPrepareVoteOwner = errors.New("Failed to recover pubkey from prepareVote owner")
	ErrPrepareVoteIllegal                = errors.New("The prepareVote is illegal")
	ErrPrepareVoteOptionIsUnknown = errors.New("The prepareVote's option is unknown")
	ErrPrepareVoteParamsInvalid   = errors.New("The prepareVote's params is invalid")
	ErrPrepareVoteSignInvalid          = errors.New("The prepareVote signature is invalid")

	// confirm
	ErrConfirmMsgInTheFuture     = errors.New("The confirmMsg is in the future")
	ErrConfirmMsgEpochInvalid     = errors.New("The confirmMsg epoch is invalid")
	ErrConfirmMsgOwnerNodeIdInvalid     = errors.New("The confirmMsg's owner nodeId is invalid")
	ErrRecoverPubkeyFromConfirmMsgOwner = errors.New("Failed to recover pubkey from confirmMsg owner")
	ErrConfirmMsgIllegal                = errors.New("The confirmMsg is illegal")
	ErrConfirmMsgSignInvalid          = errors.New("The confirmMsg signature is invalid")

	// confirmVote
	ErrConfirmVoteInTheFuture     = errors.New("The confirmVote is in the future")
	ErrConfirmVoteEpochInvalid     = errors.New("The confirmVote epoch is invalid")
	ErrConfirmVoteOwnerNodeIdInvalid     = errors.New("The confirmVote's owner nodeId is invalid")
	ErrRecoverPubkeyFromConfirmVoteOwner = errors.New("Failed to recover pubkey from confirmVote owner")
	ErrConfirmVoteIllegal                = errors.New("The confirmVote is illegal")
	ErrConfirmVoteSignInvalid          = errors.New("The confirmVote signature is invalid")
	ErrConfirmVoteOptionIsUnknown = errors.New("The prepareVote's option is unknown")


	// commit
	ErrCommitMsgInTheFuture     = errors.New("The commitMsg is in the future")
	ErrCommitMsgOwnerNodeIdInvalid     = errors.New("The commitMsg's owner nodeId is invalid")
	ErrRecoverPubkeyFromCommitMsgOwner = errors.New("Failed to recover pubkey from commitMsg owner")
	ErrCommitMsgIllegal                = errors.New("The commitMsg is illegal")
	ErrCommitMsgSignInvalid          = errors.New("The commitMsg signature is invalid")

	//ErrOrganizationIdentityNameEmpty          = errors.New("The organization identity name is empty")
	//ErrOrganizationIdentityNodeIdInvalid      = errors.New("The organization identity nodeId is invalid")
	ErrOrganizationIdentity                   = errors.New("The organization identity is invalid")
	ErrProposalTaskOperationCostInvalid       = errors.New("The proposal task OperationCost is invalid")
	ErrProposalTaskCalculateContractCodeEmpty = errors.New("The proposal task CalculateContractCode is empty")
	ErrProposalParamsInvalid                  = errors.New("The proposal params is invalid")
)
