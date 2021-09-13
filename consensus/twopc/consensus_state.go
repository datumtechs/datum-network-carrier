package twopc

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/types"
	"sync"
)

type state struct {

	// Proposal being processed (proposalId -> proposalState)
	proposalSet map[common.Hash]*ctypes.ProposalState
	// About the voting state of prepareMsg for proposal
	prepareVotes map[common.Hash]*prepareVoteState
	// About the voting state of confirmMsg for proposal
	confirmVotes map[common.Hash]*confirmVoteState
	// cache
	proposalPeerInfoCache map[common.Hash]*twopcpb.ConfirmTaskPeerInfo

	proposalsLock       sync.RWMutex
	prepareVotesLock    sync.RWMutex
	confirmVotesLock    sync.RWMutex
	confirmPeerInfoLock sync.RWMutex

	// TODO 有些缓存需要持久化
}

func newState() *state {
	return &state{
		proposalSet:           make(map[common.Hash]*ctypes.ProposalState, 0),
		prepareVotes:          make(map[common.Hash]*prepareVoteState, 0),
		confirmVotes:          make(map[common.Hash]*confirmVoteState, 0),
		proposalPeerInfoCache: make(map[common.Hash]*twopcpb.ConfirmTaskPeerInfo, 0),
	}
}

func (s *state) IsEmpty() bool    { return nil == s }
func (s *state) IsNotEmpty() bool { return !s.IsEmpty() }

func (s *state) HasOrgProposal(proposalId common.Hash, partyId string) bool {
	s.proposalsLock.RLock()
	defer s.proposalsLock.RUnlock()

	if st, ok := s.proposalSet[proposalId]; ok {
		if _, has := st.GetOrgProposalState(partyId); has {
			return true
		}
	}
	return false
}
func (s *state) HasNotOrgProposal(proposalId common.Hash, partyId string) bool {
	return !s.HasOrgProposal(proposalId, partyId)
}

func (s *state) GetProposalState(proposalId common.Hash) *ctypes.ProposalState {
	s.proposalsLock.RLock()
	proposalState := s.proposalSet[proposalId]
	s.proposalsLock.RUnlock()
	return proposalState
}
func (s *state) StoreProposalState(proposalState *ctypes.ProposalState) {
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
func (s *state) RemoveProposalState(proposalId common.Hash) {
	s.proposalsLock.Lock()
	delete(s.proposalSet, proposalId)
	s.proposalsLock.Unlock()
}

func (s *state) ChangeToConfirm(proposalId common.Hash, partyId string, startTime uint64) {
	s.proposalsLock.Lock()
	defer s.proposalsLock.Unlock()

	proposalState, ok := s.proposalSet[proposalId]
	if !ok {
		return
	}
	orgProposalState, ok := proposalState.GetOrgProposalState(partyId)
	if !ok {
		return
	}
	orgProposalState.ChangeToConfirm(startTime)
	proposalState.StoreOrgProposalState(orgProposalState)

	log.Debugf("Succeed to call `ChangeToConfirm`, proposalId: {%s}, partyId: {%s}, startTime: {%d}", proposalId.String(), partyId, startTime)

	s.proposalSet[proposalId] = proposalState
}

func (s *state) ChangeToCommit(proposalId common.Hash, partyId string, startTime uint64) {
	s.proposalsLock.Lock()
	defer s.proposalsLock.Unlock()

	proposalState, ok := s.proposalSet[proposalId]
	if !ok {
		return
	}
	orgProposalState, ok := proposalState.GetOrgProposalState(partyId)
	if !ok {
		return
	}
	orgProposalState.ChangeToCommit(startTime)
	proposalState.StoreOrgProposalState(orgProposalState)

	log.Debugf("Succeed to call `ChangeToCommit`, proposalId: {%s}, partyId: {%s}, startTime: {%d}", proposalId.String(), partyId, startTime)

	s.proposalSet[proposalId] = proposalState
}

func (s *state) CleanProposalState(proposalId common.Hash) {
	s.RemoveProposalState(proposalId)
	s.RemovePrepareVoteState(proposalId)
	s.RemoveConfirmVoteState(proposalId)
}

// ---------------- PrepareVote ----------------
func (s *state) HasPrepareVoting(proposalId common.Hash, org *apicommonpb.TaskOrganization) bool {
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
	pvs, ok := s.prepareVotes[vote.MsgOption.ProposalId]
	if !ok {
		pvs = newPrepareVoteState()
	}
	pvs.addVote(vote)
	s.prepareVotes[vote.MsgOption.ProposalId] = pvs
	s.prepareVotesLock.Unlock()
}

func (s *state) RemovePrepareVote(proposalId common.Hash, partyId string, role apicommonpb.TaskRole) {
	s.prepareVotesLock.Lock()
	pvs, ok := s.prepareVotes[proposalId]
	if !ok {
		return
	}
	pvs.removeVote(partyId, role)
	if pvs.isNotEmptyVote() {
		s.prepareVotes[proposalId] = pvs
	} else {
		delete(s.prepareVotes, proposalId)
	}
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

func (s *state) GetPrepareVote(proposalId common.Hash, partyId string) *types.PrepareVote {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return nil
	}
	return pvs.getVote(partyId)
}

func (s *state) HasPrepareVoteState(proposalId common.Hash, partyId, identityId string) bool {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if ok {
		return pvs.hasPrepareVoting(partyId, identityId)
	}
	return false
}

func (s *state) RemovePrepareVoteState(proposalId common.Hash) {
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
	return pvs.voteYesCount(apicommonpb.TaskRole_TaskRole_DataSupplier)
}
func (s *state) GetTaskPowerSupplierPrepareYesVoteCount(proposalId common.Hash) uint32 {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return pvs.voteYesCount(apicommonpb.TaskRole_TaskRole_PowerSupplier)
}
func (s *state) GetTaskReceiverPrepareYesVoteCount(proposalId common.Hash) uint32 {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return pvs.voteYesCount(apicommonpb.TaskRole_TaskRole_Receiver)
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
	return pvs.voteTotalCount(apicommonpb.TaskRole_TaskRole_DataSupplier)
}
func (s *state) GetTaskPowerSupplierPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return pvs.voteTotalCount(apicommonpb.TaskRole_TaskRole_PowerSupplier)
}
func (s *state) GetTaskReceiverPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return pvs.voteTotalCount(apicommonpb.TaskRole_TaskRole_Receiver)
}

// ---------------- ConfirmVote ----------------
func (s *state) HasConfirmVoting(proposalId common.Hash, org *apicommonpb.TaskOrganization) bool {
	s.prepareVotesLock.RLock()
	pvs, ok := s.prepareVotes[proposalId]
	s.prepareVotesLock.RUnlock()
	if !ok {
		return false
	}
	return pvs.hasPrepareVoting(org.PartyId, org.IdentityId)
}
func (s *state) StoreConfirmVote(vote *types.ConfirmVote) {
	s.confirmVotesLock.Lock()
	cvs, ok := s.confirmVotes[vote.MsgOption.ProposalId]
	if !ok {
		cvs = newConfirmVoteState()
	}
	cvs.addVote(vote)
	s.confirmVotes[vote.MsgOption.ProposalId] = cvs
	s.confirmVotesLock.Unlock()
}

func (s *state) RemoveConfirmVote(proposalId common.Hash, partyId string, role apicommonpb.TaskRole) {
	s.confirmVotesLock.Lock()
	cvs, ok := s.confirmVotes[proposalId]
	if !ok {
		return
	}
	cvs.removeVote(partyId, role)
	if cvs.isNotEmptyVote() {
		s.confirmVotes[proposalId] = cvs
	} else {
		delete(s.confirmVotes, proposalId)
	}
	s.confirmVotesLock.Unlock()
}

func (s *state) GetConfirmVoteArr(proposalId common.Hash) []*types.ConfirmVote {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return nil
	}
	return cvs.getVotes()
}

func (s *state) GetConfirmVote(proposalId common.Hash, partyId string) *types.ConfirmVote {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return nil
	}
	return cvs.getVote(partyId)
}

func (s *state) HasConfirmVoteState(proposalId common.Hash, partyId, identityId string) bool {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if ok {
		return cvs.hasConfirmVoting(partyId, identityId)
	}
	return false
}

func (s *state) RemoveConfirmVoteState(proposalId common.Hash) {
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
	return cvs.voteYesCount(apicommonpb.TaskRole_TaskRole_DataSupplier)
}
func (s *state) GetTaskPowerSupplierConfirmYesVoteCount(proposalId common.Hash) uint32 {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return cvs.voteYesCount(apicommonpb.TaskRole_TaskRole_PowerSupplier)
}
func (s *state) GetTaskReceiverConfirmYesVoteCount(proposalId common.Hash) uint32 {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return cvs.voteYesCount(apicommonpb.TaskRole_TaskRole_Receiver)
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
	return cvs.voteTotalCount(apicommonpb.TaskRole_TaskRole_DataSupplier)
}
func (s *state) GetTaskPowerSupplierConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return cvs.voteTotalCount(apicommonpb.TaskRole_TaskRole_PowerSupplier)
}
func (s *state) GetTaskReceiverConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	s.confirmVotesLock.RLock()
	cvs, ok := s.confirmVotes[proposalId]
	s.confirmVotesLock.RUnlock()
	if !ok {
		return 0
	}
	return cvs.voteTotalCount(apicommonpb.TaskRole_TaskRole_Receiver)
}

// about prepareVote
type prepareVoteState struct {
	votes      map[string]*types.PrepareVote // partyId -> vote
	yesVotes   map[apicommonpb.TaskRole]uint32
	voteStatus map[apicommonpb.TaskRole]uint32 // total vote count
	lock       sync.Mutex
}

func newPrepareVoteState() *prepareVoteState {
	return &prepareVoteState{
		votes:      make(map[string]*types.PrepareVote, 0),
		yesVotes:   make(map[apicommonpb.TaskRole]uint32, 0),
		voteStatus: make(map[apicommonpb.TaskRole]uint32, 0),
	}
}

func (st *prepareVoteState) addVote(vote *types.PrepareVote) {
	st.lock.Lock()
	defer st.lock.Unlock()

	if _, ok := st.votes[vote.MsgOption.SenderPartyId]; ok {
		return
	}
	st.votes[vote.MsgOption.SenderPartyId] = vote
	if count, ok := st.yesVotes[vote.MsgOption.SenderRole]; ok {
		if vote.VoteOption == types.Yes {
			st.yesVotes[vote.MsgOption.SenderRole] = count + 1
		}
	} else {
		if vote.VoteOption == types.Yes {
			st.yesVotes[vote.MsgOption.SenderRole] = 1
		}
	}

	if count, ok := st.voteStatus[vote.MsgOption.SenderRole]; ok {
		st.voteStatus[vote.MsgOption.SenderRole] = count + 1
	} else {
		st.voteStatus[vote.MsgOption.SenderRole] = 1
	}
}
func (st *prepareVoteState) removeVote(partyId string, role apicommonpb.TaskRole) {
	st.lock.Lock()
	defer st.lock.Unlock()

	vote, ok := st.votes[partyId]
	if !ok {
		return
	}

	delete(st.votes, partyId)

	if count, ok := st.yesVotes[role]; ok {
		if vote.VoteOption == types.Yes && count != 0 {
			st.yesVotes[role] = count - 1
		}
	}

	if count, ok := st.voteStatus[role]; ok {
		if count != 0 {
			st.voteStatus[role] = count - 1
		}
	}
}
func (st *prepareVoteState) isEmptyVote() bool                        { return len(st.votes) == 0 }
func (st *prepareVoteState) isNotEmptyVote() bool                     { return !st.isEmptyVote() }
func (st *prepareVoteState) getVote(partId string) *types.PrepareVote { return st.votes[partId] }
func (st *prepareVoteState) getVotes() []*types.PrepareVote {
	arr := make([]*types.PrepareVote, 0, len(st.votes))
	for _, vote := range st.votes {
		arr = append(arr, vote)
	}
	return arr
}
func (st *prepareVoteState) voteTotalCount(role apicommonpb.TaskRole) uint32 {
	st.lock.Lock()
	defer st.lock.Unlock()

	if count, ok := st.voteStatus[role]; ok {
		return count
	} else {
		return 0
	}
}
func (st *prepareVoteState) voteYesCount(role apicommonpb.TaskRole) uint32 {
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
		if vote.MsgOption.SenderPartyId == partyId && vote.MsgOption.Owner.GetIdentityId() == identityId {
			return true
		}
	}
	return false
}

// about confirmVote
type confirmVoteState struct {
	votes      map[string]*types.ConfirmVote // partyId -> vote
	yesVotes   map[apicommonpb.TaskRole]uint32
	voteStatus map[apicommonpb.TaskRole]uint32
	lock       sync.Mutex
}

func newConfirmVoteState() *confirmVoteState {
	return &confirmVoteState{
		votes:      make(map[string]*types.ConfirmVote, 0),
		yesVotes:   make(map[apicommonpb.TaskRole]uint32, 0),
		voteStatus: make(map[apicommonpb.TaskRole]uint32, 0),
	}
}

func (st *confirmVoteState) addVote(vote *types.ConfirmVote) {
	st.lock.Lock()
	defer st.lock.Unlock()

	if _, ok := st.votes[vote.MsgOption.SenderPartyId]; ok {
		return
	}

	st.votes[vote.MsgOption.SenderPartyId] = vote
	if count, ok := st.yesVotes[vote.MsgOption.SenderRole]; ok {
		if vote.VoteOption == types.Yes {
			st.yesVotes[vote.MsgOption.SenderRole] = count + 1
		}
	} else {
		if vote.VoteOption == types.Yes {
			st.yesVotes[vote.MsgOption.SenderRole] = 1
		}
	}

	if count, ok := st.voteStatus[vote.MsgOption.SenderRole]; ok {
		st.voteStatus[vote.MsgOption.SenderRole] = count + 1
	} else {
		st.voteStatus[vote.MsgOption.SenderRole] = 1
	}
}

func (st *confirmVoteState) removeVote(partyId string, role apicommonpb.TaskRole) {
	st.lock.Lock()
	defer st.lock.Unlock()

	vote, ok := st.votes[partyId]
	if !ok {
		return
	}

	delete(st.votes, partyId)

	if count, ok := st.yesVotes[role]; ok {
		if vote.VoteOption == types.Yes && count != 0 {
			st.yesVotes[role] = count - 1
		}
	}

	if count, ok := st.voteStatus[role]; ok {
		if count != 0 {
			st.voteStatus[role] = count - 1
		}
	}
}
func (st *confirmVoteState) isEmptyVote() bool                        { return len(st.votes) == 0 }
func (st *confirmVoteState) isNotEmptyVote() bool                     { return !st.isEmptyVote() }
func (st *confirmVoteState) getVote(partId string) *types.ConfirmVote { return st.votes[partId] }
func (st *confirmVoteState) getVotes() []*types.ConfirmVote {
	arr := make([]*types.ConfirmVote, 0, len(st.votes))
	for _, vote := range st.votes {
		arr = append(arr, vote)
	}
	return arr
}

func (st *confirmVoteState) voteYesCount(role apicommonpb.TaskRole) uint32 {
	st.lock.Lock()
	defer st.lock.Unlock()

	if count, ok := st.yesVotes[role]; ok {
		return count
	} else {
		return 0
	}
}
func (st *confirmVoteState) voteTotalCount(role apicommonpb.TaskRole) uint32 {
	st.lock.Lock()
	defer st.lock.Unlock()

	if count, ok := st.voteStatus[role]; ok {
		return count
	} else {
		return 0
	}
}
func (st *confirmVoteState) hasConfirmVoting(partyId, identityId string) bool {
	if vote, ok := st.votes[partyId]; ok {
		if vote.MsgOption.SenderPartyId == partyId && vote.MsgOption.Owner.GetIdentityId() == identityId {
			return true
		}
	}
	return false
}

func (s *state) StoreConfirmTaskPeerInfo(proposalId common.Hash, peerDesc *twopcpb.ConfirmTaskPeerInfo) {
	s.confirmPeerInfoLock.Lock()
	_, ok := s.proposalPeerInfoCache[proposalId]
	if !ok {
		s.proposalPeerInfoCache[proposalId] = peerDesc
	}
	s.confirmPeerInfoLock.Unlock()
}

func (s *state) HasConfirmTaskPeerInfo(proposalId common.Hash) bool {
	s.confirmPeerInfoLock.RLock()
	_, ok := s.proposalPeerInfoCache[proposalId]
	s.confirmPeerInfoLock.RUnlock()
	return ok
}

func (s *state) GetConfirmTaskPeerInfo(proposalId common.Hash) (*twopcpb.ConfirmTaskPeerInfo, bool) {
	s.confirmPeerInfoLock.RLock()
	peers, ok := s.proposalPeerInfoCache[proposalId]
	s.confirmPeerInfoLock.RUnlock()
	return peers, ok
}

func (s *state) MustGetConfirmTaskPeerInfo(proposalId common.Hash) *twopcpb.ConfirmTaskPeerInfo {
	peers, _ := s.GetConfirmTaskPeerInfo(proposalId)
	return peers
}

func (s *state) RemoveConfirmTaskPeerInfo(proposalId common.Hash) {
	s.confirmPeerInfoLock.Lock()
	delete(s.proposalPeerInfoCache, proposalId)
	s.confirmPeerInfoLock.Unlock()
}
