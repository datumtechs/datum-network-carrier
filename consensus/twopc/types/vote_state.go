package types

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type VoteState struct {
	prepareVoteSate       map[common.Hash]*types.PrepareVote
	confirmVoteState map[common.Hash]*types.ConfirmVote
}

func NewVoteState() *VoteState{
	return &VoteState{
		prepareVoteSate:       make(map[common.Hash]*types.PrepareVote, 0),
		confirmVoteState: make(map[common.Hash]*types.ConfirmVote, 0),
	}
}
func (state *VoteState) StorePrepareVote(vote *types.PrepareVote) {
	if _, ok := state.prepareVoteSate[vote.ProposalId]; !ok {
		state.prepareVoteSate[vote.ProposalId] = vote
	}
}
func (state *VoteState) StoreConfirmVote(vote *types.ConfirmVote) {
	if _, ok := state.confirmVoteState[vote.ProposalId]; !ok {
		state.confirmVoteState[vote.ProposalId] = vote
	}
}
func (state *VoteState) HasPrepareVote(proposalId common.Hash) bool {
	if _, ok := state.prepareVoteSate[proposalId]; ok {
		return true
	}
	return false
}
func (state *VoteState) HasConfirmVote(proposalId common.Hash) bool {

	if _, ok := state.confirmVoteState[proposalId]; ok {
		return true
	}
	return false
}
func (state *VoteState) RemovePrepareVote(proposalId common.Hash) {
	delete(state.prepareVoteSate, proposalId)
}
func (state *VoteState) RemoveConfirmVote(proposalId common.Hash) {
	delete(state.confirmVoteState, proposalId)
}