package types

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"sync"
	"time"
)

type TaskOperationCost struct {
	Processor uint32 `json:"processor"`
	Mem       uint64 `json:"mem"`
	Bandwidth uint64 `json:"bandwidth"`
	Duration  uint64 `json:"duration"`
}

func (cost *TaskOperationCost) String() string {
	return fmt.Sprintf(`{"mem": %d, "processor": %d, "bandwidth": %d, "duration": %d}`, cost.Mem, cost.Processor, cost.Bandwidth, cost.Duration)
}

type ProposalStatePeriod uint32

const (
	PeriodUnknown ProposalStatePeriod = 0
	PeriodPrepare ProposalStatePeriod = 1
	PeriodConfirm ProposalStatePeriod = 2
	PeriodCommit  ProposalStatePeriod = 3
	// The consensus message life cycle of `proposal` has reached its deadline,
	// but at this time the `GetState` of `proposal` itself has not reached the `Deadline`.
	PeriodFinished ProposalStatePeriod = 4

	PrepareMsgVotingTimeout = 3 * time.Second // 3s
	ConfirmMsgVotingTimeout = 1 * time.Second // 1s
	CommitMsgEndingTimeout  = 1 * time.Second // 1s

	//ConfirmEpochUnknown ConfirmEpoch = 0
	//ConfirmEpochFirst   ConfirmEpoch = 1
	//ConfirmEpochSecond  ConfirmEpoch = 2

)

var (
	// during 60s, if the proposal haven't been done, kill it
	ProposalDeadlineDuration = uint64(60 * (time.Second.Milliseconds()))
)

type ProposalState struct {
	proposalId common.Hash
	taskId     string
	taskSender *apicommonpb.TaskOrganization
	stateCache map[string]*OrgProposalState // partyId -> states
	lock       sync.RWMutex
}

var EmptyProposalState = new(ProposalState)

func NewProposalState(proposalId common.Hash, taskId string, sender *apicommonpb.TaskOrganization) *ProposalState {
	return &ProposalState{
		proposalId: proposalId,
		taskId:     taskId,
		taskSender: sender,
		stateCache: make(map[string]*OrgProposalState, 0),
	}
}
func RecoveryProposalState(proposalId common.Hash, taskId string, sender *apicommonpb.TaskOrganization, cache map[string]*OrgProposalState) *ProposalState {
	return &ProposalState{
		proposalId: proposalId,
		taskId:     taskId,
		taskSender: sender,
		stateCache: cache,
	}
}
func (pstate *ProposalState) MustLock() { pstate.lock.Lock() }
func (pstate *ProposalState) MustUnLock() { pstate.lock.Unlock() }
func (pstate *ProposalState) MustRLock() { pstate.lock.RLock() }
func (pstate *ProposalState) MustRUnLock() { pstate.lock.RUnlock() }
func (pstate *ProposalState) GetProposalId() common.Hash                   { return pstate.proposalId }
func (pstate *ProposalState) GetTaskId() string                            { return pstate.taskId }
func (pstate *ProposalState) GetTaskSender() *apicommonpb.TaskOrganization { return pstate.taskSender }
func (pstate *ProposalState) GetStateCache() map[string]*OrgProposalState  { return pstate.stateCache }
func (pstate *ProposalState) StoreOrgProposalState(orgState *OrgProposalState) {
	pstate.lock.Lock()
	_, ok := pstate.stateCache[orgState.TaskOrg.GetPartyId()]
	if !ok {
		pstate.stateCache[orgState.TaskOrg.GetPartyId()] = orgState
	}
	pstate.lock.Unlock()
}
func (pstate *ProposalState) StoreOrgProposalStateUnSafe(orgState *OrgProposalState) {
	pstate.stateCache[orgState.TaskOrg.GetPartyId()] = orgState
}
func (pstate *ProposalState) RemoveOrgProposalState(partyId string) {
	pstate.lock.Lock()
	delete(pstate.stateCache, partyId)
	pstate.lock.Unlock()
}
func (pstate *ProposalState) RemoveOrgProposalStateUnSafe(partyId string) {
	log.Debugf("Start Remove org proposalState whit unsafe, partyId: {%s}", partyId)
	delete(pstate.stateCache, partyId)
}

func (pstate *ProposalState) MustGetOrgProposalState(partyId string) *OrgProposalState {
	state, _ := pstate.GetOrgProposalState(partyId)
	return state
}
func (pstate *ProposalState) GetOrgProposalState(partyId string) (*OrgProposalState, bool) {
	pstate.lock.RLock()
	state, ok := pstate.stateCache[partyId]
	pstate.lock.RUnlock()
	return state, ok
}

func (pstate *ProposalState) IsEmpty() bool {
	if nil == pstate {
		return true
	}
	return len(pstate.stateCache) == 0
}
func (pstate *ProposalState) IsNotEmpty() bool { return !pstate.IsEmpty() }


type OrgProposalState struct {
	PrePeriodStartTime uint64
	PeriodStartTime    uint64 // the timestemp
	DeadlineDuration   uint64 // Clear `ProposalState` , when the current time is greater than the `DeadlineDuration` createAt of proposalState
	CreateAt           uint64
	TaskId             string
	TaskRole           apicommonpb.TaskRole
	TaskOrg            *apicommonpb.TaskOrganization
	PeriodNum          ProposalStatePeriod
}

func NewOrgProposalState(
	taskId string,
	taskRole apicommonpb.TaskRole,
	taskOrg *apicommonpb.TaskOrganization,
	startTime uint64,
) *OrgProposalState {

	return &OrgProposalState{
		TaskId:           taskId,
		TaskRole:         taskRole,
		TaskOrg:          taskOrg,
		PeriodNum:        PeriodPrepare,
		PeriodStartTime:  startTime,
		DeadlineDuration: ProposalDeadlineDuration,
		CreateAt:         timeutils.UnixMsecUint64(),
	}
}

func (pstate *OrgProposalState) CurrPeriodNum() ProposalStatePeriod {
	return pstate.PeriodNum
}

func (pstate *OrgProposalState) GetTaskId() string { return pstate.TaskId }
func (pstate *OrgProposalState) GetPeriod() string {

	switch pstate.PeriodNum {
	case PeriodPrepare:
		return "PeriodPrepare"
	case PeriodConfirm:
		return "PeriodConfirm"
	case PeriodCommit:
		return "PeriodCommit"
	case PeriodFinished:
		return "PeriodFinished"
	default:
		return "PeriodUnknown"
	}
}

func (pstate *OrgProposalState) SetDeadlineDuration(deadline uint64) {
	pstate.DeadlineDuration = deadline
}
func (pstate *OrgProposalState) AddDeadlineDuration(duration uint64) {
	pstate.DeadlineDuration += duration
}
func (pstate *OrgProposalState) IsPreparePeriod() bool     { return pstate.PeriodNum == PeriodPrepare }
func (pstate *OrgProposalState) IsConfirmPeriod() bool     { return pstate.PeriodNum == PeriodConfirm }
func (pstate *OrgProposalState) IsCommitPeriod() bool      { return pstate.PeriodNum == PeriodCommit }
func (pstate *OrgProposalState) IsFinishedPeriod() bool    { return pstate.PeriodNum == PeriodFinished }
func (pstate *OrgProposalState) IsNotPreparePeriod() bool  { return !pstate.IsPreparePeriod() }
func (pstate *OrgProposalState) IsNotConfirmPeriod() bool  { return !pstate.IsConfirmPeriod() }
func (pstate *OrgProposalState) IsNotCommitPeriod() bool   { return !pstate.IsCommitPeriod() }
func (pstate *OrgProposalState) IsNotFinishedPeriod() bool { return !pstate.IsFinishedPeriod() }
func (pstate *OrgProposalState) IsDeadline() bool {
	now := timeutils.UnixMsecUint64()
	return (now - pstate.CreateAt) >= ProposalDeadlineDuration
}

func (pstate *OrgProposalState) IsPrepareTimeout() bool {

	if !pstate.IsPreparePeriod() {
		return true
	}

	now := timeutils.UnixMsecUint64()
	duration := uint64(PrepareMsgVotingTimeout.Milliseconds())

	// Due to the time boundary problem, the value `==`
	if pstate.IsPreparePeriod() && (now-pstate.PeriodStartTime) >= duration {
		return true
	}
	return false
}
func (pstate *OrgProposalState) IsConfirmTimeout() bool {

	if pstate.IsPreparePeriod() {
		return false
	}
	if pstate.IsConfirmPeriod() {
		return false
	}
	if pstate.IsCommitPeriod() {
		return true
	}
	if pstate.IsFinishedPeriod() {
		return true
	}

	now := timeutils.UnixMsecUint64()
	duration := uint64(ConfirmMsgVotingTimeout.Milliseconds())

	if pstate.IsConfirmPeriod() && (now-pstate.PeriodStartTime) >= duration {
		return true
	}
	return false
}

func (pstate *OrgProposalState) IsCommitTimeout() bool {
	if pstate.IsPreparePeriod() {
		return false
	}
	if pstate.IsConfirmPeriod() {
		return false
	}
	if pstate.IsFinishedPeriod() {
		return true
	}

	now := timeutils.UnixMsecUint64()
	duration := uint64(CommitMsgEndingTimeout.Milliseconds())

	// Due to the time boundary problem, the value `==`
	if pstate.IsCommitPeriod() && (now-pstate.PeriodStartTime) >= duration {
		return true
	}
	return false
}

func (pstate *OrgProposalState) ChangeToConfirm(startTime uint64) {
	if pstate.PeriodNum == PeriodPrepare {
		pstate.PrePeriodStartTime = pstate.PeriodStartTime
		pstate.PeriodStartTime = startTime
		pstate.PeriodNum = PeriodConfirm
	}
	//pstate.ConfirmEpoch = ConfirmEpochFirst
}

func (pstate *OrgProposalState) ChangeToCommit(startTime uint64) {
	if pstate.PeriodNum == PeriodConfirm {
		pstate.PrePeriodStartTime = pstate.PeriodStartTime
		pstate.PeriodStartTime = startTime
		pstate.PeriodNum = PeriodCommit
	}
	//pstate.ConfirmEpoch = ConfirmEpochUnknown
}
func (pstate *OrgProposalState) ChangeToFinished(startTime uint64) {
	if pstate.PeriodNum == PeriodCommit {
		pstate.PrePeriodStartTime = pstate.PeriodStartTime
		pstate.PeriodStartTime = startTime
		pstate.PeriodNum = PeriodFinished
	}
	//pstate.ConfirmEpoch = ConfirmEpochUnknown
}
