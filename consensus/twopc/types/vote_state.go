package types

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type VoteState struct {
	prepareVoteSate       map[common.Hash]*types.PrepareVote
	confirmFirstVoteState map[common.Hash]*types.ConfirmVote
	confirmSecondVoteState map[common.Hash]*types.ConfirmVote
}

func NewVoteState() *VoteState{
	return &VoteState{
		prepareVoteSate:       make(map[common.Hash]*types.PrepareVote, 0),
		confirmFirstVoteState: make(map[common.Hash]*types.ConfirmVote, 0),
	}
}
func (state *VoteState) StorePrepareVote(vote *types.PrepareVote) {
	if _, ok := state.prepareVoteSate[vote.ProposalId]; !ok {
		state.prepareVoteSate[vote.ProposalId] = vote
	}
}
func (state *VoteState) StoreConfirmVote(vote *types.ConfirmVote) {
	if vote.Epoch ==  ConfirmEpochFirst.Uint64() {
		if _, ok := state.confirmFirstVoteState[vote.ProposalId]; !ok {
			state.confirmFirstVoteState[vote.ProposalId] = vote
		}
	} else {
		if _, ok := state.confirmSecondVoteState[vote.ProposalId]; !ok {
			state.confirmSecondVoteState[vote.ProposalId] = vote
		}
	}
}
func (state *VoteState) HasPrepareVote(proposalId common.Hash) bool {
	if _, ok := state.prepareVoteSate[proposalId]; ok {
		return true
	}
	return false
}
func (state *VoteState) HasConfirmVote(proposalId common.Hash, epoch uint64) bool {

	if epoch ==  ConfirmEpochFirst.Uint64() {
		if _, ok := state.confirmFirstVoteState[proposalId]; ok {
			return true
		}
	} else {
		if _, ok := state.confirmSecondVoteState[proposalId]; ok {
			return true
		}
	}
	return false
}
func (state *VoteState) RemovePrepareVote(proposalId common.Hash) {
	delete(state.prepareVoteSate, proposalId)
}
func (state *VoteState) RemoveConfirmVote(proposalId common.Hash) {
	delete(state.confirmFirstVoteState, proposalId)
	delete(state.confirmSecondVoteState, proposalId)
}