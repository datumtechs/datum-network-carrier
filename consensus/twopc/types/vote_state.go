package types

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/types"
	"sync"
)

type VoteState struct {
	prepareVoteSate       map[common.Hash]*types.PrepareVote
	confirmVoteState map[common.Hash]*types.ConfirmVote

	preparelock sync.RWMutex
	confirmlock sync.RWMutex
}

func NewVoteState() *VoteState{
	return &VoteState{
		prepareVoteSate:       make(map[common.Hash]*types.PrepareVote, 0),
		confirmVoteState: make(map[common.Hash]*types.ConfirmVote, 0),
	}
}
func (state *VoteState) StorePrepareVote(vote *types.PrepareVote) {
	state.preparelock.Lock()
	if _, ok := state.prepareVoteSate[vote.MsgOption.ProposalId]; !ok {
		state.prepareVoteSate[vote.MsgOption.ProposalId] = vote
	}
	state.preparelock.Unlock()
}
func (state *VoteState) StoreConfirmVote(vote *types.ConfirmVote) {
	state.confirmlock.Lock()
	if _, ok := state.confirmVoteState[vote.MsgOption.ProposalId]; !ok {
		state.confirmVoteState[vote.MsgOption.ProposalId] = vote
	}
	state.confirmlock.Unlock()
}
func (state *VoteState) HasPrepareVote(proposalId common.Hash) bool {
	state.preparelock.RLock()
	_, ok := state.prepareVoteSate[proposalId]
	state.preparelock.RUnlock()
	if ok {
		return true
	}
	return false
}
func (state *VoteState) HasConfirmVote(proposalId common.Hash) bool {
	state.confirmlock.RLock()
	_, ok := state.confirmVoteState[proposalId]
	state.confirmlock.RUnlock()
	if ok {
		return true
	}
	return false
}
func (state *VoteState) RemovePrepareVote(proposalId common.Hash) {
	state.preparelock.Lock()
	delete(state.prepareVoteSate, proposalId)
	state.preparelock.Unlock()
}
func (state *VoteState) RemoveConfirmVote(proposalId common.Hash) {
	state.confirmlock.Lock()
	delete(state.confirmVoteState, proposalId)
	state.confirmlock.Unlock()
}