package twopc

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
)

type state struct {
	// Proposal being processed (proposalId -> proposalState)
	runningProposals map[common.Hash]*ctypes.ProposalState
	// the global empty proposalState
	empty *ctypes.ProposalState
}

func newState() *state {
	return &state{
		runningProposals: make(map[common.Hash]*ctypes.ProposalState, 0),
		empty: ctypes.EmptyProposalState,
	}
}

func (s *state) HasProposal(proposalId common.Hash) bool {
	if _, ok := s.runningProposals[proposalId]; ok {
		return true
	}
	return false
}

func (s *state) HasNotProposal(proposalId common.Hash) bool {
	return !s.HasProposal(proposalId)
}

func (s *state) ProposalStates(proposalId common.Hash)  *ctypes.ProposalState {
	proposalState, ok := s.runningProposals[proposalId]
	if !ok {
		return s.empty
	}
	return proposalState
}