package twopc

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/types"
	"sync"
)

type state struct {

	// Proposal being processed (proposalId -> proposalState)
	runningProposals map[common.Hash]*ctypes.ProposalState
	selfPeerInfoCache map[common.Hash]*types.PrepareVoteResource
	//Proposal Vote State for self Org
	selfVoteState *ctypes.VoteState
	// About the voting state of prepareMsg for proposal
	prepareVotes map[common.Hash]*prepareVoteState
	// About the voting state of confirmMsg for proposal
	confirmVotes map[common.Hash]*confirmVoteState
	// cache
	proposalPeerInfoCache map[common.Hash]*pb.ConfirmTaskPeerInfo
	// the global empty proposalState
	empty *ctypes.ProposalState

	proposalsLock         sync.RWMutex
	selfPeerInfoCacheLock sync.RWMutex
	prepareVotesLock      sync.RWMutex
	confirmVotesLock      sync.RWMutex
	confirmPeerInfoLock   sync.RWMutex
}

func newState() *state {
	return &state{
		runningProposals:      make(map[common.Hash]*ctypes.ProposalState, 0),
		selfPeerInfoCache:     make(map[common.Hash]*types.PrepareVoteResource, 0),
		selfVoteState:         ctypes.NewVoteState(),
		prepareVotes:          make(map[common.Hash]*prepareVoteState, 0),
		confirmVotes:          make(map[common.Hash]*confirmVoteState, 0),
		proposalPeerInfoCache: make(map[common.Hash]*pb.ConfirmTaskPeerInfo, 0),
		empty:                 ctypes.EmptyProposalState,
	}
}

func (s *state) EmptyInfo() *ctypes.ProposalState { return s.empty }

func (s *state) HasProposal(proposalId common.Hash) bool {
	s.proposalsLock.RLock()
	defer s.proposalsLock.RUnlock()

	if _, ok := s.runningProposals[proposalId]; ok {
		return true
	}
	return false
}

func (s *state) HasNotProposal(proposalId common.Hash) bool {
	return !s.HasProposal(proposalId)
}

func (s *state) IsRecvTaskOnProposalState(proposalId common.Hash) bool {
	s.proposalsLock.RLock()
	defer s.proposalsLock.RUnlock()

	proposalState, ok := s.runningProposals[proposalId]
	if !ok {
		return false
	}
	if proposalState.TaskDir == types.RecvTaskDir {
		return true
	}
	return false
}

func (s *state) IsSendTaskOnProposalState(proposalId common.Hash) bool {
	s.proposalsLock.RLock()
	defer s.proposalsLock.RUnlock()

	proposalState, ok := s.runningProposals[proposalId]
	if !ok {
		return false
	}
	if proposalState.TaskDir == types.SendTaskDir {
		return true
	}
	return false
}

func (s *state) GetProposalState(proposalId common.Hash) *ctypes.ProposalState {
	s.proposalsLock.RLock()
	defer s.proposalsLock.RUnlock()

	proposalState, ok := s.runningProposals[proposalId]
	if !ok {
		return s.empty
	}
	return proposalState
}
func (s *state) AddProposalState(proposalState *ctypes.ProposalState) {
	s.proposalsLock.Lock()
	s.runningProposals[proposalState.ProposalId] = proposalState
	s.proposalsLock.Unlock()
}
func (s *state) UpdateProposalState(proposalState *ctypes.ProposalState) {
	s.proposalsLock.Lock()
	if _, ok := s.runningProposals[proposalState.ProposalId]; ok {
		s.runningProposals[proposalState.ProposalId] = proposalState
	}
	s.proposalsLock.Unlock()
}
func (s *state) DelProposalState(proposalId common.Hash) {
	s.proposalsLock.Lock()
	delete(s.runningProposals, proposalId)
	s.proposalsLock.Unlock()
}

func (s *state) GetProposalStates() map[common.Hash]*ctypes.ProposalState {
	return s.runningProposals
}

func (s *state) ChangeToConfirm(proposalId common.Hash, startTime uint64) {
	s.proposalsLock.Lock()
	defer s.proposalsLock.Unlock()

	log.Debugf("Start to call `ChangeToConfirm`, proposalId: {%s}, startTime: {%d}", proposalId.String(), startTime)

	proposalState, ok := s.runningProposals[proposalId]
	if !ok {
		return
	}
	proposalState.ChangeToConfirm(startTime)
	s.runningProposals[proposalId] = proposalState
}

func (s *state) ChangeToCommit(proposalId common.Hash, startTime uint64) {
	s.proposalsLock.Lock()
	defer s.proposalsLock.Unlock()

	log.Debugf("Start to call `ChangeToCommit`, proposalId: {%s}, startTime: {%d}", proposalId.String(), startTime)

	proposalState, ok := s.runningProposals[proposalId]
	if !ok {
		return
	}
	proposalState.ChangeToCommit(startTime)
	s.runningProposals[proposalId] = proposalState
}

//func (s *state) ChangeToFinished(proposalId common.Hash, startTime uint64) {
//	s.proposalsLock.Lock()
//	defer s.proposalsLock.Unlock()
//
//	log.Debugf("Start to call `ChangeToFinished`, proposalId: {%s}, startTime: {%d}", proposalId.String(), startTime)
//
//	proposalState, ok := s.runningProposals[proposalId]
//	if !ok {
//		return
//	}
//	proposalState.ChangeToFinished(startTime)
//	s.runningProposals[proposalId] = proposalState
//}

// 作为发起方时, 自己给当前 proposal 提供的资源信息 ... [根据 metaDataId 锁定的 dataNode资源]
func (s *state) StoreSelfPeerInfo(proposalId common.Hash, peerInfo *types.PrepareVoteResource) {
	log.Debugf("Start Store slefPeerInfo, proposalId: {%s}, peerInfo: {%s}", proposalId.String(), peerInfo.String())
	s.selfPeerInfoCacheLock.Lock()
	s.selfPeerInfoCache[proposalId] = peerInfo
	s.selfPeerInfoCacheLock.Unlock()
}
func (s *state) GetSelfPeerInfo(proposalId common.Hash) *types.PrepareVoteResource {
	s.selfPeerInfoCacheLock.RLock()
	r := s.selfPeerInfoCache[proposalId]
	s.selfPeerInfoCacheLock.RUnlock()
	return r
}

func (s *state) RemoveSelfPeerInfo(proposalId common.Hash) {
	s.selfPeerInfoCacheLock.Lock()
	delete(s.selfPeerInfoCache, proposalId)
	s.selfPeerInfoCacheLock.Unlock()
}

func (s *state) StoreConfirmTaskPeerInfo(proposalId common.Hash, peerDesc *pb.ConfirmTaskPeerInfo) {
	s.confirmPeerInfoLock.Lock()
	s.proposalPeerInfoCache[proposalId] = peerDesc
	s.confirmPeerInfoLock.Unlock()
}
func (s *state) GetConfirmTaskPeerInfo(proposalId common.Hash) *pb.ConfirmTaskPeerInfo {
	s.confirmPeerInfoLock.RLock()
	peer := s.proposalPeerInfoCache[proposalId]
	s.confirmPeerInfoLock.RUnlock()
	return peer
}

func (s *state) RemoveConfirmTaskPeerInfo(proposalId common.Hash) {
	s.confirmPeerInfoLock.Lock()
	delete(s.proposalPeerInfoCache, proposalId)
	s.confirmPeerInfoLock.Unlock()
}

func (s *state) StorePrepareVoteState(vote *types.PrepareVote) {
	s.selfVoteState.StorePrepareVote(vote)
}
func (s *state) StoreConfirmVoteState(vote *types.ConfirmVote) {
	s.selfVoteState.StoreConfirmVote(vote)
}

func (s *state) HasPrepareVoteState(proposalId common.Hash) bool {
	return s.selfVoteState.HasPrepareVote(proposalId)
}
func (s *state) HasConfirmVoteState(proposalId common.Hash) bool {
	return s.selfVoteState.HasConfirmVote(proposalId)
}

func (s *state) RemovePrepareVoteState(proposalId common.Hash) {
	s.selfVoteState.RemovePrepareVote(proposalId)
}
func (s *state) RemoveConfirmVoteState(proposalId common.Hash) {
	s.selfVoteState.RemoveConfirmVote(proposalId)
}
func (s *state) CleanProposalState(proposalId common.Hash) {
	s.DelProposalState(proposalId)
	s.RemovePrepareVoteState(proposalId)
	s.RemoveConfirmVoteState(proposalId)
	s.CleanPrepareVoteState(proposalId)
	s.CleanConfirmVoteState(proposalId)
	s.RemoveSelfPeerInfo(proposalId)
	s.RemoveConfirmTaskPeerInfo(proposalId)
}

// ---------------- PrepareVote ----------------
func (s *state) HasPrepareVoting(identityId string, proposalId common.Hash) bool {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return false
	}
	return pvs.hasPrepareVoting(identityId)
}
func (s *state) StorePrepareVote(vote *types.PrepareVote) {
	s.prepareVotesLock.Lock()
	pvs, ok := s.prepareVotes[vote.ProposalId]
	if !ok {
		pvs = newPrepareVoteState()
	}
	pvs.addVote(vote)
	s.prepareVotes[vote.ProposalId] = pvs
	s.prepareVotesLock.Unlock()
}

func (s *state) GetPrepareVoteArr(proposalId common.Hash) []*types.PrepareVote {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return nil
	}
	return pvs.getVotes()
}

func (s *state) CleanPrepareVoteState(proposalId common.Hash) {
	s.prepareVotesLock.Lock()
	delete(s.prepareVotes, proposalId)
	s.prepareVotesLock.Unlock()
}
func (s *state) GetTaskPrepareYesVoteCount(proposalId common.Hash) uint32 {
	return s.GetTaskDataSupplierPrepareYesVoteCount(proposalId) +
		s.GetTaskPowerSupplierPrepareYesVoteCount(proposalId) +
		s.GetTaskResulterPrepareYesVoteCount(proposalId)
}
func (s *state) GetTaskDataSupplierPrepareYesVoteCount(proposalId common.Hash) uint32 {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return pvs.voteYesCount(types.DataSupplier)
}
func (s *state) GetTaskPowerSupplierPrepareYesVoteCount(proposalId common.Hash) uint32 {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return pvs.voteYesCount(types.PowerSupplier)
}
func (s *state) GetTaskResulterPrepareYesVoteCount(proposalId common.Hash) uint32 {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return pvs.voteYesCount(types.ResultSupplier)
}

func (s *state) GetTaskPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	return s.GetTaskDataSupplierPrepareTotalVoteCount(proposalId) +
		s.GetTaskPowerSupplierPrepareTotalVoteCount(proposalId) +
		s.GetTaskResulterPrepareTotalVoteCount(proposalId)
}
func (s *state) GetTaskDataSupplierPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return pvs.voteTotalCount(types.DataSupplier)
}
func (s *state) GetTaskPowerSupplierPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return pvs.voteTotalCount(types.PowerSupplier)
}
func (s *state) GetTaskResulterPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return pvs.voteTotalCount(types.ResultSupplier)
}

// ---------------- ConfirmVote ----------------
func (s *state) HasConfirmVoting(identityId string, proposalId common.Hash) bool {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return false
	}
	return cvs.hasConfirmVoting(identityId)
}
func (s *state) StoreConfirmVote(vote *types.ConfirmVote) {
	s.confirmVotesLock.Lock()
	cvs, ok := s.confirmVotes[vote.ProposalId]
	if !ok {
		cvs = newConfirmVoteState()
	}
	cvs.addVote(vote)
	s.confirmVotes[vote.ProposalId] = cvs
	s.confirmVotesLock.Unlock()
}
func (s *state) CleanConfirmVoteState(proposalId common.Hash) {
	s.confirmVotesLock.Lock()
	delete(s.confirmVotes, proposalId)
	s.confirmVotesLock.Unlock()
}
func (s *state) GetTaskConfirmYesVoteCount(proposalId common.Hash) uint32 {
	return s.GetTaskDataSupplierConfirmYesVoteCount(proposalId) +
		s.GetTaskPowerSupplierConfirmYesVoteCount(proposalId) +
		s.GetTaskResulterConfirmYesVoteCount(proposalId)
}
func (s *state) GetTaskDataSupplierConfirmYesVoteCount(proposalId common.Hash) uint32 {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return cvs.voteYesCount(types.DataSupplier)
}
func (s *state) GetTaskPowerSupplierConfirmYesVoteCount(proposalId common.Hash) uint32 {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return cvs.voteYesCount(types.PowerSupplier)
}
func (s *state) GetTaskResulterConfirmYesVoteCount(proposalId common.Hash) uint32 {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return cvs.voteYesCount(types.ResultSupplier)
}
func (s *state) GetTaskConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	return s.GetTaskDataSupplierConfirmTotalVoteCount(proposalId) +
		s.GetTaskPowerSupplierConfirmTotalVoteCount(proposalId) +
		s.GetTaskResulterConfirmTotalVoteCount(proposalId)
}
func (s *state) GetTaskDataSupplierConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return cvs.voteTotalCount(types.DataSupplier)
}
func (s *state) GetTaskPowerSupplierConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return cvs.voteTotalCount(types.PowerSupplier)
}
func (s *state) GetTaskResulterConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return cvs.voteTotalCount(types.ResultSupplier)
}

// about prepareVote
type prepareVoteState struct {
	votes      []*types.PrepareVote
	yesVotes   map[types.TaskRole]uint32
	voteStatus map[types.TaskRole]uint32
	lock       sync.Mutex
}

func newPrepareVoteState() *prepareVoteState {
	return &prepareVoteState{
		votes:      make([]*types.PrepareVote, 0),
		yesVotes:   make(map[types.TaskRole]uint32, 0),
		voteStatus: make(map[types.TaskRole]uint32, 0),
	}
}

func (st *prepareVoteState) addVote(vote *types.PrepareVote) {
	st.lock.Lock()
	defer st.lock.Unlock()

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
func (st *prepareVoteState) getVotes() []*types.PrepareVote { return st.votes }
func (st *prepareVoteState) voteTotalCount(role types.TaskRole) uint32 {
	st.lock.Lock()
	defer st.lock.Unlock()

	if count, ok := st.voteStatus[role]; ok {
		return count
	} else {
		return 0
	}
}
func (st *prepareVoteState) voteYesCount(role types.TaskRole) uint32 {
	st.lock.Lock()
	defer st.lock.Unlock()

	if count, ok := st.yesVotes[role]; ok {
		return count
	} else {
		return 0
	}
}
func (st *prepareVoteState) hasPrepareVoting(identityId string) bool {
	for _, vote := range st.votes {
		if vote.Owner.IdentityId == identityId {
			return true
		}
	}
	return false
}

// about confirmVote
type confirmVoteState struct {
	votes      []*types.ConfirmVote
	yesVotes   map[types.TaskRole]uint32
	voteStatus map[types.TaskRole]uint32
	lock       sync.Mutex
}

func newConfirmVoteState() *confirmVoteState {
	return &confirmVoteState{
		votes:      make([]*types.ConfirmVote, 0),
		yesVotes:   make(map[types.TaskRole]uint32, 0),
		voteStatus: make(map[types.TaskRole]uint32, 0),
	}
}

func (st *confirmVoteState) addVote(vote *types.ConfirmVote) {
	st.lock.Lock()
	defer st.lock.Unlock()

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
	st.lock.Lock()
	defer st.lock.Unlock()

	if count, ok := st.yesVotes[role]; ok {
		return count
	} else {
		return 0
	}
}
func (st *confirmVoteState) voteTotalCount(role types.TaskRole) uint32 {
	st.lock.Lock()
	defer st.lock.Unlock()

	if count, ok := st.voteStatus[role]; ok {
		return count
	} else {
		return 0
	}
}
func (st *confirmVoteState) hasConfirmVoting(identityId string) bool {
	for _, vote := range st.votes {
		if vote.Owner.IdentityId == identityId {
			return true
		}
	}
	return false
}
