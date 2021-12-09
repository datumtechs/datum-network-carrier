package types

import "errors"

var (
	ErrConsensusMsgInvalid = errors.New("the consensus msg is invalid")
	ErrMsgTaskInvalid      = errors.New("the task of msg is invalid")

	// GetTask
	ErrProposalTaskNotFound = errors.New("the task of proposal not found")

	// Prepare
	ErrProposalIllegal           = errors.New("the proposal is illegal")
	ErrProposalAlreadyProcessed  = errors.New("the proposalId have been processing")
	ErrProposalNotFound          = errors.New("the proposal not found")
	ErrProposalInTheFuture       = errors.New("the proposal is in the future")
	ErrPrososalTaskRoleIsUnknown = errors.New("the proposal's task role is unknown")
	ErrPrososalTaskIsProcessed   = errors.New("the task of proposal have been processing")

	// prepareVote
	ErrPrepareVoteInTheFuture = errors.New("the prepareVote is in the future")

	ErrPrepareVoteIllegal         = errors.New("the prepareVote is illegal")
	ErrPrepareVoteOptionIsUnknown = errors.New("the prepareVote's option is unknown")
	ErrPrepareVoteParamsInvalid   = errors.New("the prepareVote's params is invalid")
	ErrPrepareVotehadVoted        = errors.New("the prepareMsg we had voted")
	ErrPrepareVoteRepeatedly      = errors.New("the voter had voted prepareVote repeatedly")

	// confirm
	ErrConfirmMsgInTheFuture  = errors.New("the confirmMsg is in the future")
	ErrConfirmMsgEpochInvalid = errors.New("the confirmMsg epoch is invalid")
	ErrConfirmMsgIllegal      = errors.New("the confirmMsg is illegal")

	// confirmVote
	ErrConfirmVoteInTheFuture     = errors.New("the confirmVote is in the future")
	ErrConfirmVoteEpochInvalid    = errors.New("the confirmVote epoch is invalid")
	ErrConfirmVoteIllegal         = errors.New("the confirmVote is illegal")
	ErrConfirmVoteOptionIsUnknown = errors.New("the prepareVote's option is unknown")
	ErrConfirmVotehadVoted        = errors.New("the confirmMsg we had voted")
	ErrConfirmVoteRepeatedly      = errors.New("the voter had voted confirmVote repeatedly")

	// commit
	ErrCommitMsgInTheFuture = errors.New("the commitMsg is in the future")
	ErrCommitMsgIllegal     = errors.New("the commitMsg is illegal")

	ErrOrganizationIdentity                   = errors.New("the organization identity is invalid")
	ErrProposalTaskOperationCostInvalid       = errors.New("the proposal task OperationCost is invalid")
	ErrProposalTaskCalculateContractCodeEmpty = errors.New("the proposal task CalculateContractCode is empty")
	ErrProposalParamsInvalid                  = errors.New("the proposal params is invalid")

	ErrMsgSignInvalid        = errors.New("the msg signature is invalid")
	ErrMsgOwnerNodeIdInvalid = errors.New("the msg's owner nodeId is invalid")

	ErrProposalWithoutPreparePeriod = errors.New("the proposal without prepare period")
	ErrProposalWithoutConfirmPeriod = errors.New("the proposal without confirm period")
	ErrProposalWithoutCommitPeriod  = errors.New("the proposal without commit period")
	ErrProposalConfirmEpochInvalid  = errors.New("the confirm epoch of proposal is invalid")

	ErrProposalPrepareVoteTimeout  = errors.New("received prepareVote of proposal timeout")
	ErrProposalConfirmMsgTimeout   = errors.New("received confirmMsg of proposal timeout")
	ErrProposalConfirmVoteTimeout  = errors.New("received confirmVote of proposal timeout")
	ErrProposalCommitMsgTimeout    = errors.New("received commitMsg of proposal timeout")
	ErrTaskResultMsgInvalid        = errors.New("received taskResultMsg is invalid")
	ErrTaskResourceUsageMsgInvalid = errors.New("received taskResourceUsageMsg is invalid")

	ErrProposalPrepareVoteFuture = errors.New("received prepareVote of proposal is future msg")
	ErrProposalConfirmMsgFuture  = errors.New("received confirmMsg of proposal is future msg")
	ErrProposalConfirmVoteFuture = errors.New("received confirmVote of proposal is future msg")
	ErrProposalCommitMsgFuture   = errors.New("received commitMsg of proposal is future msg")

	ErrProposalPrepareVoteOwnerInvalid     = errors.New("the owner of proposal's prepareVote is invalid")
	ErrProposalConfirmVoteVoteOwnerInvalid = errors.New("the owner of proposal's confirmVote is invalid")

	ErrProposalPrepareVoteResourceInvalid = errors.New("the resource of proposal's prepareVote is invalid")

	ErrVoteCountOverflow = errors.New("the vote count has overflow")
)
