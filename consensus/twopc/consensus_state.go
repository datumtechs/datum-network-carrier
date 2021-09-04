package twopc

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/types"
	"sync"
)

type state struct {

	// Proposal being processed (proposalId -> proposalState)
	proposalSet map[common.Hash]*ctypes.ProposalState
	//selfPeerInfoCache map[common.Hash]*types.PrepareVoteResource
	//Proposal Vote State for self Org
	//selfVoteState *ctypes.VoteState
	// About the voting state of prepareMsg for proposal
	prepareVotes map[common.Hash]*prepareVoteState
	// About the voting state of confirmMsg for proposal
	confirmVotes map[common.Hash]*confirmVoteState
	//// cache
	//proposalPeerInfoCache map[common.Hash]*pb.ConfirmTaskPeerInfo
	// the global empty proposalState
	empty *ctypes.ProposalState

	proposalsLock         sync.RWMutex
	//selfPeerInfoCacheLock sync.RWMutex
	prepareVotesLock      sync.RWMutex
	confirmVotesLock      sync.RWMutex
	//confirmPeerInfoLock   sync.RWMutex
}

func newState() *state {
	return &state{
		proposalSet:           make(map[common.Hash]*ctypes.ProposalState, 0),
		//selfPeerInfoCache:     make(map[common.Hash]*types.PrepareVoteResource, 0),
		//selfVoteState:         ctypes.NewVoteState(),
		prepareVotes:          make(map[common.Hash]*prepareVoteState, 0),
		confirmVotes:          make(map[common.Hash]*confirmVoteState, 0),
		//proposalPeerInfoCache: make(map[common.Hash]*pb.ConfirmTaskPeerInfo, 0),
		empty:                 ctypes.EmptyProposalState,
	}
}

func (s *state) EmptyInfo() *ctypes.ProposalState { return s.empty }

func (s *state) HasProposal(proposalId common.Hash) bool {
	s.proposalsLock.RLock()
	defer s.proposalsLock.RUnlock()

	if _, ok := s.proposalSet[proposalId]; ok {
		return true
	}
	return false
}

func (s *state) HasNotProposal(proposalId common.Hash) bool {
	return !s.HasProposal(proposalId)
}

//func (s *state) IsRecvTaskOnProposalState(proposalId common.Hash) bool {
//	s.proposalsLock.RLock()
//	defer s.proposalsLock.RUnlock()
//
//	proposalState, ok := s.proposalSet[proposalId]
//	if !ok {
//		return false
//	}
//	if proposalState.TaskDir == types.RecvTaskDir {
//		return true
//	}
//	return false
//}
//
//func (s *state) IsSendTaskOnProposalState(proposalId common.Hash) bool {
//	s.proposalsLock.RLock()
//	defer s.proposalsLock.RUnlock()
//
//	proposalState, ok := s.proposalSet[proposalId]
//	if !ok {
//		return false
//	}
//	if proposalState.TaskDir == types.SendTaskDir {
//		return true
//	}
//	return false
//}

func (s *state) GetProposalState(proposalId common.Hash) *ctypes.ProposalState {
	s.proposalsLock.RLock()
	defer s.proposalsLock.RUnlock()

	proposalState, ok := s.proposalSet[proposalId]
	if !ok {
		return s.empty
	}
	return proposalState
}
func (s *state) AddProposalState(proposalState *ctypes.ProposalState) {
	s.proposalsLock.Lock()
	s.proposalSet[proposalState.GetProposalId()] = proposalState
	s.proposalsLock.Unlock()
}
func (s *state) UpdateProposalState(proposalState *ctypes.ProposalState) {
	s.proposalsLock.Lock()
	if _, ok := s.proposalSet[proposalState.GetProposalId()]; ok {
		s.proposalSet[proposalState.GetProposalId()] = proposalState
	}
	s.proposalsLock.Unlock()
}
func (s *state) DelProposalState(proposalId common.Hash) {
	s.proposalsLock.Lock()
	delete(s.proposalSet, proposalId)
	s.proposalsLock.Unlock()
}

//func (s *state) GetProposalStates() map[common.Hash]*ctypes.ProposalState {
//	s.proposalsLock.RLock()
//	proposals := s.proposalSet
//	proposalMap := make(map[common.Hash]*ctypes.ProposalState, len(proposals))
//	for hash, proposalState := range proposals {
//		proposalMap[hash] = &ctypes.ProposalState{
//			ProposalId: proposalState.ProposalId,
//			TaskDir:    proposalState.TaskDir,
//			TaskRole:   proposalState.TaskRole,
//			TaskOrg: &apipb.TaskOrganization{
//				PartyId:    proposalState.TaskOrg.PartyId,
//				NodeName:   proposalState.TaskOrg.NodeName,
//				NodeId:     proposalState.TaskOrg.NodeId,
//				IdentityId: proposalState.TaskOrg.IdentityId,
//			},
//			TaskId:             proposalState.TaskId,
//			PeriodNum:          proposalState.PeriodNum,
//			PrePeriodStartTime: proposalState.PrePeriodStartTime,
//			PeriodStartTime:    proposalState.PeriodStartTime,
//			DeadlineDuration:   proposalState.DeadlineDuration,
//			CreateAt:           proposalState.CreateAt,
//		}
//	}
//	s.proposalsLock.RUnlock()
//	return proposalMap
//}

func (s *state) ChangeToConfirm(proposalId common.Hash, startTime uint64) {
	s.proposalsLock.Lock()
	defer s.proposalsLock.Unlock()

	log.Debugf("Start to call `ChangeToConfirm`, proposalId: {%s}, startTime: {%d}", proposalId.String(), startTime)

	proposalState, ok := s.proposalSet[proposalId]
	if !ok {
		return
	}
	proposalState.ChangeToConfirm(startTime)
	s.proposalSet[proposalId] = proposalState
}

func (s *state) ChangeToCommit(proposalId common.Hash, startTime uint64) {
	s.proposalsLock.Lock()
	defer s.proposalsLock.Unlock()

	log.Debugf("Start to call `ChangeToCommit`, proposalId: {%s}, startTime: {%d}", proposalId.String(), startTime)

	proposalState, ok := s.proposalSet[proposalId]
	if !ok {
		return
	}
	proposalState.ChangeToCommit(startTime)
	s.proposalSet[proposalId] = proposalState
}

//func (s *state) StoreConfirmTaskPeerInfo(proposalId common.Hash, peerDesc *pb.ConfirmTaskPeerInfo) {
//	s.confirmPeerInfoLock.Lock()
//	s.proposalPeerInfoCache[proposalId] = peerDesc
//	s.confirmPeerInfoLock.Unlock()
//}
//func (s *state) GetConfirmTaskPeerInfo(proposalId common.Hash) *pb.ConfirmTaskPeerInfo {
//	s.confirmPeerInfoLock.RLock()
//	peer := s.proposalPeerInfoCache[proposalId]
//	s.confirmPeerInfoLock.RUnlock()
//	return peer
//}
//
//func (s *state) RemoveConfirmTaskPeerInfo(proposalId common.Hash) {
//	s.confirmPeerInfoLock.Lock()
//	delete(s.proposalPeerInfoCache, proposalId)
//	s.confirmPeerInfoLock.Unlock()
//}
func (s *state) CleanProposalState(proposalId common.Hash) {
	s.DelProposalState(proposalId)
	//s.RemovePrepareVoteState(proposalId)
	//s.RemoveConfirmVoteState(proposalId)
	s.CleanPrepareVoteState(proposalId)
	s.CleanConfirmVoteState(proposalId)
	//s.RemoveConfirmTaskPeerInfo(proposalId)
}

// ---------------- PrepareVote ----------------
func (s *state) HasPrepareVoting(proposalId common.Hash, org *apipb.TaskOrganization) bool {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return false
	}
	return pvs.hasPrepareVoting(org.PartyId, org.IdentityId)
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
		s.GetTaskReceiverPrepareYesVoteCount(proposalId)
}
func (s *state) GetTaskDataSupplierPrepareYesVoteCount(proposalId common.Hash) uint32 {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return pvs.voteYesCount(apipb.TaskRole_TaskRole_DataSupplier)
}
func (s *state) GetTaskPowerSupplierPrepareYesVoteCount(proposalId common.Hash) uint32 {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return pvs.voteYesCount(apipb.TaskRole_TaskRole_PowerSupplier)
}
func (s *state) GetTaskReceiverPrepareYesVoteCount(proposalId common.Hash) uint32 {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return pvs.voteYesCount(apipb.TaskRole_TaskRole_Receiver)
}

func (s *state) GetTaskPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	return s.GetTaskDataSupplierPrepareTotalVoteCount(proposalId) +
		s.GetTaskPowerSupplierPrepareTotalVoteCount(proposalId) +
		s.GetTaskReceiverPrepareTotalVoteCount(proposalId)
}
func (s *state) GetTaskDataSupplierPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return pvs.voteTotalCount(apipb.TaskRole_TaskRole_DataSupplier)
}
func (s *state) GetTaskPowerSupplierPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return pvs.voteTotalCount(apipb.TaskRole_TaskRole_PowerSupplier)
}
func (s *state) GetTaskReceiverPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return pvs.voteTotalCount(apipb.TaskRole_TaskRole_Receiver)
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
		s.GetTaskReceiverConfirmYesVoteCount(proposalId)
}
func (s *state) GetTaskDataSupplierConfirmYesVoteCount(proposalId common.Hash) uint32 {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return cvs.voteYesCount(apipb.TaskRole_TaskRole_DataSupplier)
}
func (s *state) GetTaskPowerSupplierConfirmYesVoteCount(proposalId common.Hash) uint32 {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return cvs.voteYesCount(apipb.TaskRole_TaskRole_PowerSupplier)
}
func (s *state) GetTaskReceiverConfirmYesVoteCount(proposalId common.Hash) uint32 {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return cvs.voteYesCount(apipb.TaskRole_TaskRole_Receiver)
}
func (s *state) GetTaskConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	return s.GetTaskDataSupplierConfirmTotalVoteCount(proposalId) +
		s.GetTaskPowerSupplierConfirmTotalVoteCount(proposalId) +
		s.GetTaskReceiverConfirmTotalVoteCount(proposalId)
}
func (s *state) GetTaskDataSupplierConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return cvs.voteTotalCount(apipb.TaskRole_TaskRole_DataSupplier)
}
func (s *state) GetTaskPowerSupplierConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return cvs.voteTotalCount(apipb.TaskRole_TaskRole_PowerSupplier)
}
func (s *state) GetTaskReceiverConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return cvs.voteTotalCount(apipb.TaskRole_TaskRole_Receiver)
}

// about prepareVote
type prepareVoteState struct {
	votes      map[string]*types.PrepareVote // partyId -> vote
	yesVotes   map[apipb.TaskRole]uint32
	voteStatus map[apipb.TaskRole]uint32 // total vote count
	lock       sync.Mutex
}

func newPrepareVoteState() *prepareVoteState {
	return &prepareVoteState{
		votes:      make(map[string]*types.PrepareVote, 0),
		yesVotes:   make(map[apipb.TaskRole]uint32, 0),
		voteStatus: make(map[apipb.TaskRole]uint32, 0),
	}
}

func (st *prepareVoteState) addVote(vote *types.PrepareVote) {
	st.lock.Lock()
	defer st.lock.Unlock()

	if _, ok := st.votes[vote.PartyId]; ok {
		return
	}
	st.votes[vote.PartyId] = vote
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

func (st *prepareVoteState) getVote(partId string) *types.PrepareVote { return st.votes[partId] }
func (st *prepareVoteState) getVotes() []*types.PrepareVote {
	arr := make([]*types.PrepareVote, 0, len(st.votes))
	for _, vote := range st.votes {
		arr = append(arr, vote)
	}
	return arr
}
func (st *prepareVoteState) voteTotalCount(role apipb.TaskRole) uint32 {
	st.lock.Lock()
	defer st.lock.Unlock()

	if count, ok := st.voteStatus[role]; ok {
		return count
	} else {
		return 0
	}
}
func (st *prepareVoteState) voteYesCount(role apipb.TaskRole) uint32 {
	st.lock.Lock()
	defer st.lock.Unlock()

	if count, ok := st.yesVotes[role]; ok {
		return count
	} else {
		return 0
	}
}
func (st *prepareVoteState) hasPrepareVoting(partyId, identityId string) bool {
	if vote, ok := st.votes[partyId]; ok {
		if vote.PartyId == partyId && vote.Owner.IdentityId == identityId {
			return true
		}
	}
	return false
}

// about confirmVote
type confirmVoteState struct {
	votes      map[string]*types.ConfirmVote // partyId -> vote
	yesVotes   map[apipb.TaskRole]uint32
	voteStatus map[apipb.TaskRole]uint32
	lock       sync.Mutex
}

func newConfirmVoteState() *confirmVoteState {
	return &confirmVoteState{
		votes:      make(map[string]*types.ConfirmVote, 0),
		yesVotes:   make(map[apipb.TaskRole]uint32, 0),
		voteStatus: make(map[apipb.TaskRole]uint32, 0),
	}
}

func (st *confirmVoteState) addVote(vote *types.ConfirmVote) {
	st.lock.Lock()
	defer st.lock.Unlock()

	if _, ok := st.votes[vote.PartyId]; ok {
		return
	}

	st.votes[vote.PartyId] = vote
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
func (st *confirmVoteState) voteYesCount(role apipb.TaskRole) uint32 {
	st.lock.Lock()
	defer st.lock.Unlock()

	if count, ok := st.yesVotes[role]; ok {
		return count
	} else {
		return 0
	}
}
func (st *confirmVoteState) voteTotalCount(role apipb.TaskRole) uint32 {
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
