package types

import "errors"

var (
	ErrProposalIllegal                = errors.New("The proposal is illegal")
	ErrProposalAlreadyProcessed       = errors.New("The proposalId have been processing")
	ErrProposalInTheFuture            = errors.New("The proposal is in the future")
	ErrPrepareMsgNotHasSign           = errors.New("The prepareMsg has not signature")
	ErrPrepareMsgSignInvalid          = errors.New("The prepareMsg signature is invalid")
	ErrProposalOwnerNodeIdInvalid     = errors.New("The proposal's owner nodeId is invalid")
	ErrRecoverPubkeyFromProposalOwner = errors.New("Failed to recover pubkey from proposal owner")
	ErrPrososalTaskRoleIsUnknown      = errors.New("The proposal's taskRole is unknown")
	ErrPrososalTaskIdBelongSendTask   = errors.New("The proposal's taskId is belonging to send task")
	ErrPrososalTaskIdBelongRecvTask   = errors.New("The proposal's taskId is belonging to recv task")

	//ErrOrganizationIdentityNameEmpty          = errors.New("The organization identity name is empty")
	//ErrOrganizationIdentityNodeIdInvalid      = errors.New("The organization identity nodeId is invalid")
	ErrOrganizationIdentity = errors.New("The organization identity is invalid")
	ErrProposalTaskOperationCostInvalid       = errors.New("The proposal task OperationCost is invalid")
	ErrProposalTaskCalculateContractCodeEmpty = errors.New("The proposal task CalculateContractCode is empty")
	ErrProposalParamsInvalid                  = errors.New("The proposal params is invalid")
)
