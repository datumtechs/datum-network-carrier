package twopc

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/types"
	"sync"
	"time"
)

type state struct {

	selfPeerInfoCache  map[common.Hash]*types.PrepareVoteResource
	//Proposal Vote State for self Org
	selfVoteState *ctypes.VoteState
	// About the voting state of prepareMsg for proposal
	prepareVotes map[common.Hash]*prepareVoteState
	// About the voting state of confirmMsg for proposal
	confirmVotes map[common.Hash]*confirmVoteState
	// Proposal being processed (proposalId -> proposalState)
	runningProposals map[common.Hash]*ctypes.ProposalState
	// cache
	proposalPeerInfoCache map[common.Hash]*pb.ConfirmTaskPeerInfo
	// the global empty proposalState
	empty *ctypes.ProposalState

	proposalsLock   sync.RWMutex
}

func newState() *state {
	return &state{
		selfPeerInfoCache: make(map[common.Hash]*types.PrepareVoteResource, 0),
		selfVoteState:    ctypes.NewVoteState(),
		prepareVotes:     make(map[common.Hash]*prepareVoteState, 0),
		confirmVotes:     make(map[common.Hash]*confirmVoteState, 0),
		runningProposals: make(map[common.Hash]*ctypes.ProposalState, 0),
		proposalPeerInfoCache: make(map[common.Hash]*pb.ConfirmTaskPeerInfo, 0),
		empty:            ctypes.EmptyProposalState,
	}
}

func (s *state) EmptyInfo () *ctypes.ProposalState { return s.empty }

func (s *state) CleanExpireProposal() ([]string, []string) {
	now := uint64(time.Now().UnixNano())
	sendTaskIds := make([]string, 0)
	recvTaskIds := make([]string, 0)

	// This function does not need to be locked,
	// and the small function inside handles the lock by itself
	for id, proposal := range s.runningProposals {
		if (now - proposal.CreateAt) > ctypes.ProposalDeadlineDuration {
			log.Info("Clean 2pc expire Proposal", "proposalId", id.String(), "taskId",
				proposal.TaskId, "taskDir", proposal.TaskDir.String())
			s.CleanProposalState(id)
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
	if proposalState.TaskDir == ctypes.RecvTaskDir {
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
	if proposalState.TaskDir == ctypes.SendTaskDir {
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

func (s *state) GetProposalStates () map[common.Hash]*ctypes.ProposalState {
	return s.runningProposals
}

func (s *state) ChangeToConfirm(proposalId common.Hash, startTime uint64) {
	s.proposalsLock.Lock()
	defer s.proposalsLock.Unlock()

	proposalState, ok := s.runningProposals[proposalId]
	if !ok {
		return
	}
	proposalState.ChangeToConfirm(startTime)
	s.runningProposals[proposalId] = proposalState
}
//func (s *state) ChangeToConfirmSecondEpoch(proposalId common.Hash, startTime uint64) {
//	proposalState, ok := s.runningProposals[proposalId]
//	if !ok {
//		return
//	}
//	proposalState.ChangeToConfirmSecondEpoch(startTime)
//}
func (s *state) ChangeToCommit(proposalId common.Hash, startTime uint64) {
	s.proposalsLock.Lock()
	defer s.proposalsLock.Unlock()

	proposalState, ok := s.runningProposals[proposalId]
	if !ok {
		return
	}
	proposalState.ChangeToCommit(startTime)
	s.runningProposals[proposalId] = proposalState
}

func (s *state) ChangeToFinised(proposalId common.Hash, startTime uint64) {
	s.proposalsLock.Lock()
	defer s.proposalsLock.Unlock()

	proposalState, ok := s.runningProposals[proposalId]
	if !ok {
		return
	}
	proposalState.ChangeToFinished(startTime)
	s.runningProposals[proposalId] = proposalState
}

// 作为发起方时, 自己给当前 proposal 提供的资源信息 ...
func (s *state) StoreSelfPeerInfo(proposalId common.Hash, peerInfo *types.PrepareVoteResource) {
	s.selfPeerInfoCache[proposalId] = peerInfo
}
func (s *state) GetSelfPeerInfo(proposalId common.Hash)  *types.PrepareVoteResource{
	return s.selfPeerInfoCache[proposalId]
}

func (s *state) RemoveSelfPeerInfo(proposalId common.Hash) {
	delete(s.selfPeerInfoCache, proposalId)
}

func (s *state) StoreConfirmTaskPeerInfo(proposalId common.Hash, peerDesc *pb.ConfirmTaskPeerInfo) {
	s.proposalPeerInfoCache[proposalId] = peerDesc
}
func (s *state) GetConfirmTaskPeerInfo(proposalId common.Hash) *pb.ConfirmTaskPeerInfo {
	return s.proposalPeerInfoCache[proposalId]
}

func (s *state) RemoveConfirmTaskPeerInfo(proposalId common.Hash) {
	delete(s.proposalPeerInfoCache, proposalId)
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
func  (s *state) HasConfirmVoteState(proposalId common.Hash, epoch uint64) bool {
	return s.selfVoteState.HasConfirmVote(proposalId, epoch)
}

func  (s *state) RemovePrepareVoteState(proposalId common.Hash) {
	s.selfVoteState.RemovePrepareVote(proposalId)
}
func  (s *state) RemoveConfirmVoteState(proposalId common.Hash) {
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
	pvs, ok := s.prepareVotes[proposalId]
	if !ok {
		return false
	}
	return pvs.hasPrepareVoting(identityId)
}
func (s *state) StorePrepareVote(vote *types.PrepareVote) {
	pvs, ok := s.prepareVotes[vote.ProposalId]
	if !ok {
		pvs = newPrepareVoteState()
	}
	pvs.addVote(vote)
	s.prepareVotes[vote.ProposalId] = pvs
}

func (s *state) GetPrepareVoteArr(proposalId common.Hash) []*types.PrepareVote {
	pvs, ok := s.prepareVotes[proposalId]
	if !ok {
		return nil
	}
	return  pvs.getVotes()
}

func (s *state) CleanPrepareVoteState(proposalId common.Hash) {
	delete(s.prepareVotes, proposalId)
}
func (s *state) GetTaskDataSupplierPrepareYesVoteCount(proposalId common.Hash) uint32 {
	pvs, ok := s.prepareVotes[proposalId]
	if !ok {
		return 0
	}
	return pvs.voteYesCount(types.DataSupplier)
}
func (s *state) GetTaskPowerSupplierPrepareYesVoteCount(proposalId common.Hash) uint32 {
	pvs, ok := s.prepareVotes[proposalId]
	if !ok {
		return 0
	}
	return pvs.voteYesCount(types.PowerSupplier)
}
func (s *state) GetTaskResulterPrepareYesVoteCount(proposalId common.Hash) uint32 {
	pvs, ok := s.prepareVotes[proposalId]
	if !ok {
		return 0
	}
	return pvs.voteYesCount(types.ResultSupplier)
}
func (s *state) GetTaskDataSupplierPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	pvs, ok := s.prepareVotes[proposalId]
	if !ok {
		return 0
	}
	return pvs.voteTotalCount(types.DataSupplier)
}
func (s *state) GetTaskPowerSupplierPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	pvs, ok := s.prepareVotes[proposalId]
	if !ok {
		return 0
	}
	return pvs.voteTotalCount(types.PowerSupplier)
}
func (s *state) GetTaskResulterPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	pvs, ok := s.prepareVotes[proposalId]
	if !ok {
		return 0
	}
	return pvs.voteTotalCount(types.ResultSupplier)
}

// ---------------- ConfirmVote ----------------
func (s *state) HasConfirmVoting(identityId string, proposalId common.Hash) bool {
	cvs, ok := s.confirmVotes[proposalId]
	if !ok {
		return false
	}
	return cvs.hasConfirmVoting(identityId)
}
func (s *state) StoreConfirmVote(vote *types.ConfirmVote) {
	cvs, ok := s.confirmVotes[vote.ProposalId]
	if !ok {
		cvs = newConfirmVoteState()
	}
	cvs.addVote(vote)
	s.confirmVotes[vote.ProposalId] = cvs
}
func (s *state) CleanConfirmVoteState(proposalId common.Hash) {
	delete(s.confirmVotes, proposalId)
}
func (s *state) GetTaskDataSupplierConfirmYesVoteCount(proposalId common.Hash) uint32 {
	cvs, ok := s.confirmVotes[proposalId]
	if !ok {
		return 0
	}
	return cvs.voteYesCount(types.DataSupplier)
}
func (s *state) GetTaskPowerSupplierConfirmYesVoteCount(proposalId common.Hash) uint32 {
	cvs, ok := s.confirmVotes[proposalId]
	if !ok {
		return 0
	}
	return cvs.voteYesCount(types.PowerSupplier)
}
func (s *state) GetTaskResulterConfirmYesVoteCount(proposalId common.Hash) uint32 {
	cvs, ok := s.confirmVotes[proposalId]
	if !ok {
		return 0
	}
	return cvs.voteYesCount(types.ResultSupplier)
}
func (s *state) GetTaskDataSupplierConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	cvs, ok := s.confirmVotes[proposalId]
	if !ok {
		return 0
	}
	return cvs.voteTotalCount(types.DataSupplier)
}
func (s *state) GetTaskPowerSupplierConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	cvs, ok := s.confirmVotes[proposalId]
	if !ok {
		return 0
	}
	return cvs.voteTotalCount(types.PowerSupplier)
}
func (s *state) GetTaskResulterConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	cvs, ok := s.confirmVotes[proposalId]
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
func (st *prepareVoteState) getVotes () []*types.PrepareVote { return st.votes }
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
func (st *confirmVoteState) hasConfirmVoting(identityId string) bool {
	for _, vote := range st.votes {
		if vote.Owner.IdentityId == identityId {
			return true
		}
	}
	return false
}
