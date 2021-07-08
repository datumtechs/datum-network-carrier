package twopc

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"time"
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

func (s *state) CleanExpireProposal() ([]string, []string) {
	now := uint64(time.Now().UnixNano())
	sendTaskIds := make([]string, 0)
	recvTaskIds := make([]string, 0)
	for id, proposal := range s.runningProposals {
		if (now - proposal.CreateAt) > ctypes.ProposalDeadlineDuration {
			log.Info("Clean 2pc expire Proposal", "proposalId", id.String(), "taskId", proposal.TaskId)
			delete(s.runningProposals, id)
			if proposal.TaskDir == ctypes.SendTaskDir {
				sendTaskIds = append(sendTaskIds, proposal.TaskId)
			} else {
				recvTaskIds = append(recvTaskIds, proposal.TaskId)
			}

		}
	}
	return sendTaskIds, recvTaskIds
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
func (s *state) AddProposalState(proposalState *ctypes.ProposalState) {
	s.runningProposals[proposalState.ProposalId] = proposalState
}
func (s *state) UpdateProposalState(proposalState *ctypes.ProposalState) {
	if _, ok := s.runningProposals[proposalState.ProposalId]; ok {
		s.runningProposals[proposalState.ProposalId] = proposalState
	}
}
func (s *state) DelProposalState(proposalHash common.Hash) {
	delete(s.runningProposals, proposalHash)
}