package types

import "errors"

var (

	// Task
	ErrProposalTaskNotFound = errors.New("The task of proposal not found")

	// Prepare
	ErrProposalIllegal                = errors.New("The proposal is illegal")
	ErrProposalAlreadyProcessed       = errors.New("The proposalId have been processing")
	ErrProposalNotFound               = errors.New("The proposal not found")
	ErrProposalInTheFuture            = errors.New("The proposal is in the future")
	ErrPrososalTaskRoleIsUnknown      = errors.New("The proposal's task role is unknown")
	//ErrPrososalTaskIdBelongSendTask   = errors.New("The proposal's taskId is belonging to send task")
	//ErrPrososalTaskIdBelongRecvTask   = errors.New("The proposal's taskId is belonging to recv task")
	ErrPrososalTaskIsProcessed  = errors.New("The task of proposal have been processing")

	// prepareVote
	ErrPrepareVoteInTheFuture     = errors.New("The prepareVote is in the future")


	ErrPrepareVoteIllegal                = errors.New("The prepareVote is illegal")
	ErrPrepareVoteOptionIsUnknown = errors.New("The prepareVote's option is unknown")
	ErrPrepareVoteParamsInvalid   = errors.New("The prepareVote's params is invalid")
	ErrPrepareVotehadVoted	 = errors.New("The prepareMsg we had voted")

	// confirm
	ErrConfirmMsgInTheFuture     = errors.New("The confirmMsg is in the future")
	ErrConfirmMsgEpochInvalid     = errors.New("The confirmMsg epoch is invalid")
	ErrConfirmMsgIllegal                = errors.New("The confirmMsg is illegal")

	// confirmVote
	ErrConfirmVoteInTheFuture     = errors.New("The confirmVote is in the future")
	ErrConfirmVoteEpochInvalid     = errors.New("The confirmVote epoch is invalid")
	ErrConfirmVoteIllegal                = errors.New("The confirmVote is illegal")
	ErrConfirmVoteOptionIsUnknown = errors.New("The prepareVote's option is unknown")


	// commit
	ErrCommitMsgInTheFuture     = errors.New("The commitMsg is in the future")
	ErrCommitMsgIllegal                = errors.New("The commitMsg is illegal")

	//ErrOrganizationIdentityNameEmpty          = errors.New("The organization identity name is empty")
	//ErrOrganizationIdentityNodeIdInvalid      = errors.New("The organization identity nodeId is invalid")
	ErrOrganizationIdentity                   = errors.New("The organization identity is invalid")
	ErrProposalTaskOperationCostInvalid       = errors.New("The proposal task OperationCost is invalid")
	ErrProposalTaskCalculateContractCodeEmpty = errors.New("The proposal task CalculateContractCode is empty")
	ErrProposalParamsInvalid                  = errors.New("The proposal params is invalid")



	ErrMsgSignInvalid          = errors.New("The msg signature is invalid")
	ErrMsgOwnerNodeIdInvalid     = errors.New("The msg's owner nodeId is invalid")

	ErrProposalWithoutPreparePeriod = errors.New("The proposal without prepare period")
	ErrProposalWithoutConfirmPeriod = errors.New("The proposal without confirm period")
	ErrProposalWithoutCommitPeriod = errors.New("The proposal without commit period")
	ErrProposalConfirmEpochInvalid = errors.New("The confirm epoch of proposal is invalid")

	ErrProposalPrepareVoteTimeout = errors.New("The prepareVote of proposal has timeout")
	ErrProposalConfirmVoteTimeout = errors.New("The confirmVote of proposal has timeout")

	ErrProposalPrepareVoteOwnerInvalid = errors.New("The owner of proposal's prepareVote is invalid")
	ErrProposalConfirmVoteVoteOwnerInvalid = errors.New("The owner of proposal's confirmVote is invalid")


	ErrVoteCountOverflow = errors.New("The vote count has overflow")
)
