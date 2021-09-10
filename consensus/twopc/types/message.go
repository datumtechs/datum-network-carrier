package types

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
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
	// but at this time the `State` of `proposal` itself has not reached the `Deadline`.
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

type OrgProposalState struct {
	PrePeriodStartTime uint64
	PeriodStartTime    uint64 // the timestemp
	DeadlineDuration   uint64 // Clear `ProposalState` , when the current time is greater than the `DeadlineDuration` createAt of proposalState
	CreateAt           uint64
	TaskId             string
	TaskRole           apipb.TaskRole
	TaskOrg            *apipb.TaskOrganization
	PeriodNum          ProposalStatePeriod
}

func NewOrgProposalState(
	taskId string,
	taskRole apipb.TaskRole,
	taskOrg   *apipb.TaskOrganization,
	startTime uint64,
) *OrgProposalState {

	return &OrgProposalState{
		TaskId:           taskId,
		TaskRole:         taskRole,
		TaskOrg:          taskOrg,
		PeriodNum:        PeriodPrepare,
		PeriodStartTime:  startTime,
		DeadlineDuration: ProposalDeadlineDuration,
		CreateAt:         uint64(timeutils.UnixMsec()),
	}
}

type ProposalState struct {
	proposalId common.Hash
	taskId     string
	taskSender *apipb.TaskOrganization
	stateCache map[string]*OrgProposalState // partyId -> states
	lock       sync.RWMutex
}

var EmptyProposalState = new(ProposalState)

func NewProposalState(proposalId common.Hash, taskId string, sender *apipb.TaskOrganization) *ProposalState {
	return &ProposalState{
		proposalId: proposalId,
		taskId: taskId,
		taskSender: sender,
		stateCache: make(map[string]*OrgProposalState, 0),
	}
}
func (pstate *ProposalState) GetProposalId() common.Hash { return pstate.proposalId }
func (pstate *ProposalState) GetTaskId() string          { return pstate.taskId }
func (pstate *ProposalState) GetTaskSender() *apipb.TaskOrganization { return pstate.taskSender }

func (pstate *ProposalState) StoreOrgProposalState(orgState *OrgProposalState) {

	pstate.lock.Lock()
	_, ok := pstate.stateCache[orgState.TaskOrg.GetPartyId()]
	if !ok {
		pstate.stateCache[orgState.TaskOrg.GetPartyId()] = orgState
	}
	pstate.lock.Unlock()
}
func (pstate *ProposalState) RemoveOrgProposalState(partyId string) {
	pstate.lock.Lock()
	delete(pstate.stateCache, partyId)
	pstate.lock.Unlock()
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

func (pstate *ProposalState) RefreshProposalState() {

	pstate.lock.Lock()
	for partyId, orgState := range pstate.stateCache {

		if orgState.IsDeadline() {
			log.Debugf("Started refresh proposalState loop, the proposalState direct be deadline, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
				pstate.GetProposalId().String(), pstate.GetTaskId(), partyId)

			delete(pstate.stateCache, partyId)
			continue
		}

		switch orgState.CurrPeriodNum() {

		case PeriodPrepare:

			if orgState.IsPrepareTimeout() {
				log.Debugf("Started refresh org proposalState, the org proposalState was prepareTimeout, change to confirm epoch, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					pstate.GetProposalId().String(), pstate.GetTaskId(), partyId)

				orgState.ChangeToConfirm(orgState.PeriodStartTime + uint64(PrepareMsgVotingTimeout.Milliseconds()))
				pstate.stateCache[partyId] = orgState
			}

		case PeriodConfirm:

			if orgState.IsConfirmTimeout() {
				log.Debugf("Started refresh org proposalState, the org proposalState was confirmTimeout, change to commit epoch, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					pstate.GetProposalId().String(), pstate.GetTaskId(), partyId)

				orgState.ChangeToCommit(orgState.PeriodStartTime + uint64(ConfirmMsgVotingTimeout.Milliseconds()))
				pstate.stateCache[partyId] = orgState
			}

		case PeriodCommit:

			if orgState.IsCommitTimeout() {
				log.Debugf("Started refresh org proposalState, the org proposalState was commitTimeout, change to finished epoch, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					pstate.GetProposalId().String(), pstate.GetTaskId(), partyId)

				orgState.ChangeToFinished(orgState.PeriodStartTime + uint64(CommitMsgEndingTimeout.Milliseconds()))
				pstate.stateCache[partyId] = orgState
			}

		case PeriodFinished:

			if orgState.IsDeadline() {
				log.Debugf("Started refresh org proposalState, the org proposalState was finished, but coming deadline now, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					pstate.GetProposalId().String(), pstate.GetTaskId(), partyId)

				delete(pstate.stateCache, partyId)


			}

		default:
			log.Errorf("Unknown the proposalState period,  proposalId: {%s}, taskId: {%s}, partyId: {%s}", pstate.GetProposalId().String(), pstate.GetTaskId(), partyId)
		}
	}

	pstate.lock.Unlock()

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
func (pstate *OrgProposalState) IsPreparePeriod() bool     { return pstate.PeriodNum == PeriodPrepare }
func (pstate *OrgProposalState) IsConfirmPeriod() bool     { return pstate.PeriodNum == PeriodConfirm }
func (pstate *OrgProposalState) IsCommitPeriod() bool      { return pstate.PeriodNum == PeriodCommit }
func (pstate *OrgProposalState) IsFinishedPeriod() bool    { return pstate.PeriodNum == PeriodFinished }
func (pstate *OrgProposalState) IsNotPreparePeriod() bool  { return !pstate.IsPreparePeriod() }
func (pstate *OrgProposalState) IsNotConfirmPeriod() bool  { return !pstate.IsConfirmPeriod() }
func (pstate *OrgProposalState) IsNotCommitPeriod() bool   { return !pstate.IsCommitPeriod() }
func (pstate *OrgProposalState) IsNotFinishedPeriod() bool { return !pstate.IsFinishedPeriod() }
func (pstate *OrgProposalState) IsDeadline() bool {
	now := uint64(timeutils.UnixMsec())
	return (now - pstate.CreateAt) >= ProposalDeadlineDuration
}

func (pstate *OrgProposalState) IsPrepareTimeout() bool {

	if !pstate.IsPreparePeriod() {
		return true
	}

	now := uint64(timeutils.UnixMsec())
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

	now := uint64(timeutils.UnixMsec())
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

	now := uint64(timeutils.UnixMsec())
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
