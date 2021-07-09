package twopc

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"time"
)

type state struct {

	//Proposal Vote State for self Org
	selfVoteState *ctypes.VoteState
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
		selfVoteState:    ctypes.NewVoteState(),
		prepareVotes:     make(map[common.Hash]*prepareVoteState, 0),
		confirmVotes:     make(map[common.Hash]*confirmVoteState, 0),
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
func (s *state) DelProposalState(proposalId common.Hash) { delete(s.runningProposals, proposalId) }

func (s *state) ChangeToConfirm(proposalId common.Hash, startTime uint64) {
	proposalState, ok := s.runningProposals[proposalId]
	if !ok {
		return
	}
	proposalState.ChangeToConfirm(startTime)
}
func (s *state) ChangeToConfirmSecondEpoch(proposalId common.Hash, startTime uint64) {
	proposalState, ok := s.runningProposals[proposalId]
	if !ok {
		return
	}
	proposalState.ChangeToConfirmSecondEpoch(startTime)
}
func (s *state) ChangeToCommit(proposalId common.Hash, startTime uint64) {
	proposalState, ok := s.runningProposals[proposalId]
	if !ok {
		return
	}
	proposalState.ChangeToCommit(startTime)
}

func  (s *state) StorePrepareVoteState(vote *types.PrepareVote) {
	s.selfVoteState.StorePrepareVote(vote)
}
func  (s *state) StoreConfirmVoteState(vote *types.ConfirmVote) {
	s.selfVoteState.StoreConfirmVote(vote)
}

func  (s *state) HasPrepareVoteState(proposalId common.Hash) bool {
	return s.selfVoteState.HasPrepareVote(proposalId)
}
func  (s *state) HasConfirmVoteState(proposalId common.Hash) bool {
	return s.selfVoteState.HasConfirmVote(proposalId)
}

func  (s *state) RemovePrepareVoteState(proposalId common.Hash) {
	s.selfVoteState.RemovePrepareVote(proposalId)
}
func  (s *state) RemoveConfirmVoteState(proposalId common.Hash) {
	s.selfVoteState.RemoveConfirmVote(proposalId)
}


// ---------------- PrepareVote ----------------
func (s *state) StorePrepareVote(vote *types.PrepareVote) {
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
func (s *state) GetTaskDataSupplierPrepareYesVoteCount(proposalId common.Hash) uint32 {
	state, ok := s.prepareVotes[proposalId]
	if !ok {
		return 0
	}
	return state.voteYesCount(types.DataSupplier)
}
func (s *state) GetTaskPowerSupplierPrepareYesVoteCount(proposalId common.Hash) uint32 {
	state, ok := s.prepareVotes[proposalId]
	if !ok {
		return 0
	}
	return state.voteYesCount(types.PowerSupplier)
}
func (s *state) GetTaskResulterPrepareYesVoteCount(proposalId common.Hash) uint32 {
	state, ok := s.prepareVotes[proposalId]
	if !ok {
		return 0
	}
	return state.voteYesCount(types.ResultSupplier)
}
func (s *state) GetTaskDataSupplierPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	state, ok := s.prepareVotes[proposalId]
	if !ok {
		return 0
	}
	return state.voteTotalCount(types.DataSupplier)
}
func (s *state) GetTaskPowerSupplierPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	state, ok := s.prepareVotes[proposalId]
	if !ok {
		return 0
	}
	return state.voteTotalCount(types.PowerSupplier)
}
func (s *state) GetTaskResulterPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	state, ok := s.prepareVotes[proposalId]
	if !ok {
		return 0
	}
	return state.voteTotalCount(types.ResultSupplier)
}

// ---------------- ConfirmVote ----------------
func (s *state) StoreConfirmVote(vote *types.ConfirmVote) {
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
func (s *state) GetTaskDataSupplierConfirmYesVoteCount(proposalId common.Hash) uint32 {
	state, ok := s.confirmVotes[proposalId]
	if !ok {
		return 0
	}
	return state.voteYesCount(types.DataSupplier)
}
func (s *state) GetTaskPowerSupplierConfirmYesVoteCount(proposalId common.Hash) uint32 {
	state, ok := s.confirmVotes[proposalId]
	if !ok {
		return 0
	}
	return state.voteYesCount(types.PowerSupplier)
}
func (s *state) GetTaskResulterConfirmYesVoteCount(proposalId common.Hash) uint32 {
	state, ok := s.confirmVotes[proposalId]
	if !ok {
		return 0
	}
	return state.voteYesCount(types.ResultSupplier)
}
func (s *state) GetTaskDataSupplierConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	state, ok := s.confirmVotes[proposalId]
	if !ok {
		return 0
	}
	return state.voteTotalCount(types.DataSupplier)
}
func (s *state) GetTaskPowerSupplierConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	state, ok := s.confirmVotes[proposalId]
	if !ok {
		return 0
	}
	return state.voteTotalCount(types.PowerSupplier)
}
func (s *state) GetTaskResulterConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	state, ok := s.confirmVotes[proposalId]
	if !ok {
		return 0
	}
	return state.voteTotalCount(types.ResultSupplier)
}

// about prepareVote
type prepareVoteState struct {
	votes      []*types.PrepareVote
	yesVotes   map[types.TaskRole]uint32
	voteStatus map[types.TaskRole]uint32
}

func newPrepareVoteState() *prepareVoteState {
	return &prepareVoteState{
		votes:      make([]*types.PrepareVote, 0),
		yesVotes:   make(map[types.TaskRole]uint32, 0),
		voteStatus: make(map[types.TaskRole]uint32, 0),
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

	if count, ok := st.voteStatus[vote.TaskRole]; ok {
		st.voteStatus[vote.TaskRole] = count + 1
	} else {
		st.voteStatus[vote.TaskRole] = 1
	}
}
func (st *prepareVoteState) voteTotalCount(role types.TaskRole) uint32 {
	if count, ok := st.voteStatus[role]; ok {
		return count
	} else {
		return 0
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
	votes      []*types.ConfirmVote
	yesVotes   map[types.TaskRole]uint32
	voteStatus map[types.TaskRole]uint32
}

func newConfirmVoteState() *confirmVoteState {
	return &confirmVoteState{
		votes:      make([]*types.ConfirmVote, 0),
		yesVotes:   make(map[types.TaskRole]uint32, 0),
		voteStatus: make(map[types.TaskRole]uint32, 0),
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

	if count, ok := st.voteStatus[vote.TaskRole]; ok {
		st.voteStatus[vote.TaskRole] = count + 1
	} else {
		st.voteStatus[vote.TaskRole] = 1
	}
}
func (st *confirmVoteState) voteYesCount(role types.TaskRole) uint32 {
	if count, ok := st.yesVotes[role]; ok {
		return count
	} else {
		return 0
	}
}
func (st *confirmVoteState) voteTotalCount(role types.TaskRole) uint32 {
	if count, ok := st.voteStatus[role]; ok {
		return count
	} else {
		return 0
	}
}
