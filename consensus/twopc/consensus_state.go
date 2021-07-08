package twopc

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"time"
)

type state struct {
	// About the voting state of prepareMsg for proposal
	prepareVotes map[common.Hash]*prepareVoteState
	// About the voting state of confirmMsg for proposal
	confirmVotes map[common.Hash]*confirmVoteState
	// Proposal being processed (proposalId -> proposalState)
	runningProposals map[common.Hash]*ctypes.ProposalState
	// the global empty proposalState
	empty *ctypes.ProposalState
}

func newState() *state {
	return &state{
		runningProposals: make(map[common.Hash]*ctypes.ProposalState, 0),
		empty:            ctypes.EmptyProposalState,
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

func (s *state) GetProposalState(proposalId common.Hash) *ctypes.ProposalState {
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

// ---------------- PrepareVote ----------------
func (s *state) AddPrepareVote(vote *types.PrepareVote) {
	state, ok := s.prepareVotes[vote.ProposalId]
	if !ok {
		state = newPrepareVoteState()
	}
	state.addVote(vote)
	s.prepareVotes[vote.ProposalId] = state
}
func (s *state) CleanPrepareVoteState(proposalId common.Hash) {
	delete(s.prepareVotes, proposalId)
}
func (s *state) GetTaskDataSupplierPrepareVoteCount(proposalId common.Hash) uint32 {
	state, ok := s.prepareVotes[proposalId]
	if !ok {
		return 0
	}
	return state.voteYesCount(types.DataSupplier)
}
func (s *state) GetTaskPowerSupplierPrepareVoteCount(proposalId common.Hash) uint32 {
	state, ok := s.prepareVotes[proposalId]
	if !ok {
		return 0
	}
	return state.voteYesCount(types.PowerSupplier)
}
func (s *state) GetTaskResulterPrepareVoteCount(proposalId common.Hash) uint32 {
	state, ok := s.prepareVotes[proposalId]
	if !ok {
		return 0
	}
	return state.voteYesCount(types.ResultSupplier)
}

// ---------------- ConfirmVote ----------------
func (s *state) AddConfirmVote(vote *types.ConfirmVote) {
	state, ok := s.confirmVotes[vote.ProposalId]
	if !ok {
		state = newConfirmVoteState()
	}
	state.addVote(vote)
	s.confirmVotes[vote.ProposalId] = state
}
func (s *state) CleanConfirmVoteState(proposalId common.Hash) {
	delete(s.confirmVotes, proposalId)
}
func (s *state) GetTaskDataSupplierConfirmVoteCount(proposalId common.Hash) uint32 {
	state, ok := s.confirmVotes[proposalId]
	if !ok {
		return 0
	}
	return state.voteYesCount(types.DataSupplier)
}
func (s *state) GetTaskPowerSupplierConfirmVoteCount(proposalId common.Hash) uint32 {
	state, ok := s.confirmVotes[proposalId]
	if !ok {
		return 0
	}
	return state.voteYesCount(types.PowerSupplier)
}
func (s *state) GetTaskResulterConfirmVoteCount(proposalId common.Hash) uint32 {
	state, ok := s.confirmVotes[proposalId]
	if !ok {
		return 0
	}
	return state.voteYesCount(types.ResultSupplier)
}

// about prepareVote
type prepareVoteState struct {
	votes    []*types.PrepareVote
	yesVotes map[types.TaskRole]uint32
}

func newPrepareVoteState() *prepareVoteState {
	return &prepareVoteState{
		votes:    make([]*types.PrepareVote, 0),
		yesVotes: make(map[types.TaskRole]uint32, 0),
	}
}

func (st *prepareVoteState) addVote(vote *types.PrepareVote) {
	st.votes = append(st.votes, vote)
	if count, ok := st.yesVotes[vote.TaskRole]; ok {
		if vote.VoteOption == types.Yes {
			st.yesVotes[vote.TaskRole] = count + 1
		}
	} else {
		if vote.VoteOption == types.Yes {
			st.yesVotes[vote.TaskRole] = 1
		}
	}
}
func (st *prepareVoteState) voteYesCount(role types.TaskRole) uint32 {
	if count, ok := st.yesVotes[role]; ok {
		return count
	} else {
		return 0
	}
}

// about confirmVote
type confirmVoteState struct {
	votes    []*types.ConfirmVote
	yesVotes map[types.TaskRole]uint32
}

func newConfirmVoteState() *confirmVoteState {
	return &confirmVoteState{
		votes:    make([]*types.ConfirmVote, 0),
		yesVotes: make(map[types.TaskRole]uint32, 0),
	}
}

func (st *confirmVoteState) addVote(vote *types.ConfirmVote) {
	st.votes = append(st.votes, vote)
	if count, ok := st.yesVotes[vote.TaskRole]; ok {
		if vote.VoteOption == types.Yes {
			st.yesVotes[vote.TaskRole] = count + 1
		}
	} else {
		if vote.VoteOption == types.Yes {
			st.yesVotes[vote.TaskRole] = 1
		}
	}
}
func (st *confirmVoteState) voteYesCount(role types.TaskRole) uint32 {
	if count, ok := st.yesVotes[role]; ok {
		return count
	} else {
		return 0
	}
}
