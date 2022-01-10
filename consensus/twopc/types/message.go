package types

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"math"
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
func (cost *TaskOperationCost) GetMem() uint64       { return cost.Mem }
func (cost *TaskOperationCost) GetBandwidth() uint64 { return cost.Bandwidth }
func (cost *TaskOperationCost) GetProcessor() uint32 { return cost.Processor }

type ProposalStatePeriod uint32

const (
	PeriodUnknown ProposalStatePeriod = 0
	PeriodPrepare ProposalStatePeriod = 1
	PeriodConfirm ProposalStatePeriod = 2
	PeriodCommit  ProposalStatePeriod = 3
	// The consensus message life cycle of `proposal` has reached its deadline,
	// but at this time the `GetState` of `proposal` itself has not reached the `Deadline`.
	//
	// (it will be changed on `refresh loop` only)
	PeriodFinished ProposalStatePeriod = 4

	// the consensus period:
	//
	//   prepare epoch [prepare start time] <-- prepare voting duration --> confirm epoch [confirm start time] <-- confirm voting duration --> commit epoch [commit start time] <-- commit ending duration --> finished
	//
	PrepareMsgVotingDuration = 3 * time.Second // 3s
	ConfirmMsgVotingDuration = 1 * time.Second // 1s
	CommitMsgEndingDuration  = 1 * time.Second // 1s

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

//func (pstate *ProposalState) MustLock()                                    { pstate.lock.Lock() }
//func (pstate *ProposalState) MustUnLock()                                  { pstate.lock.Unlock() }
//func (pstate *ProposalState) MustRLock()                                   { pstate.lock.RLock() }
//func (pstate *ProposalState) MustRUnLock()                                 { pstate.lock.RUnlock() }
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
	log.Debugf("Start Remove org proposalState whith unsafe, partyId: {%s}", partyId)
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
	DeadlineDuration uint64 // Clear `ProposalState` , when the current time is greater than the `DeadlineDuration` createAt of proposalState
	CreateAt         uint64 // the time is that the proposal state was created on current iden
	StartAt          uint64 // the time is that the proposal state was created on task sender iden
	TaskId           string
	TaskRole         apicommonpb.TaskRole
	TaskOrg          *apicommonpb.TaskOrganization
	PeriodNum        ProposalStatePeriod
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
		DeadlineDuration: ProposalDeadlineDuration,
		CreateAt:         timeutils.UnixMsecUint64(),
		StartAt:          startTime,
	}
}

func (pstate *OrgProposalState) GetTaskId() string                         { return pstate.TaskId }
func (pstate *OrgProposalState) GetTaskRole() apicommonpb.TaskRole         { return pstate.TaskRole }
func (pstate *OrgProposalState) GetTaskOrg() *apicommonpb.TaskOrganization { return pstate.TaskOrg }
func (pstate *OrgProposalState) GetPeriodNum() ProposalStatePeriod         { return pstate.PeriodNum }
func (pstate *OrgProposalState) GetDeadlineDuration() uint64               { return pstate.DeadlineDuration }
func (pstate *OrgProposalState) GetCreateAt() uint64                       { return pstate.CreateAt }
func (pstate *OrgProposalState) GetStartAt() uint64                        { return pstate.StartAt }
func (pstate *OrgProposalState) GetPrepareExpireTime() int64 {
	return int64(pstate.GetStartAt()) + PrepareMsgVotingDuration.Milliseconds()
}
func (pstate *OrgProposalState) GetConfirmExpireTime() int64 {
	return int64(pstate.GetStartAt()) + PrepareMsgVotingDuration.Milliseconds() + ConfirmMsgVotingDuration.Milliseconds()
}
func (pstate *OrgProposalState) GetCommitExpireTime() int64 {
	return int64(pstate.GetStartAt()) + PrepareMsgVotingDuration.Milliseconds() + ConfirmMsgVotingDuration.Milliseconds() + CommitMsgEndingDuration.Milliseconds()
}
func (pstate *OrgProposalState) GetDeadlineExpireTime() int64 {
	return int64(pstate.GetStartAt() + pstate.GetDeadlineDuration())
}

func (pstate *OrgProposalState) GetPeriodStr() string {
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
	now := timeutils.UnixMsecUint64()
	return (now - pstate.GetStartAt()) >= ProposalDeadlineDuration
}

func (pstate *OrgProposalState) IsPrepareTimeout() bool {

	now := timeutils.UnixMsecUint64()
	duration := uint64(PrepareMsgVotingDuration.Milliseconds())

	if pstate.IsPreparePeriod() && (now-pstate.GetStartAt()) >= duration {
		return true
	}
	return false
}

func (pstate *OrgProposalState) IsConfirmTimeout() bool {

	now := timeutils.UnixMsecUint64()
	duration := uint64(PrepareMsgVotingDuration.Milliseconds()) + uint64(ConfirmMsgVotingDuration.Milliseconds())

	if pstate.IsConfirmPeriod() && (now-pstate.GetStartAt()) >= duration {
		return true
	}
	return false
}

func (pstate *OrgProposalState) IsCommitTimeout() bool {

	now := timeutils.UnixMsecUint64()
	duration := uint64(PrepareMsgVotingDuration.Milliseconds()) + uint64(ConfirmMsgVotingDuration.Milliseconds()) + uint64(CommitMsgEndingDuration.Milliseconds())

	if pstate.IsCommitPeriod() && (now-pstate.GetStartAt()) >= duration {
		return true
	}
	return false
}

func (pstate *OrgProposalState) ChangeToConfirm() {
	if pstate.PeriodNum == PeriodPrepare {
		pstate.PeriodNum = PeriodConfirm
	}
}

func (pstate *OrgProposalState) ChangeToCommit() {
	if pstate.PeriodNum == PeriodConfirm {
		pstate.PeriodNum = PeriodCommit
	}
}
func (pstate *OrgProposalState) ChangeToFinished() {
	if pstate.PeriodNum == PeriodCommit {
		pstate.PeriodNum = PeriodFinished
	}
}

// v 0.3.0 proposal state expire monitor
type ProposalStateMonitor struct {
	proposalId common.Hash
	partyId    string
	sender     *apicommonpb.TaskOrganization
	orgState   *OrgProposalState
	when       int64 // target timestamp
	next       int64
	fn         func(proposalId common.Hash, sender *apicommonpb.TaskOrganization, orgState *OrgProposalState)
}

func NewProposalStateMonitor(proposalId common.Hash, partyId string, sender *apicommonpb.TaskOrganization,
	orgState *OrgProposalState, when, next int64,
	fn func(proposalId common.Hash, sender *apicommonpb.TaskOrganization, orgState *OrgProposalState)) *ProposalStateMonitor {

	return &ProposalStateMonitor{
		proposalId: proposalId,
		partyId:    partyId,
		sender:     sender,
		orgState:   orgState,
		when:       when,
		next:       next,
		fn:         fn,
	}
}
func (psm *ProposalStateMonitor) GetProposalId() common.Hash { return psm.proposalId }
func (psm *ProposalStateMonitor) GetPartyId() string         { return psm.partyId }
func (psm *ProposalStateMonitor) GetWhen() int64             { return psm.when }
func (psm *ProposalStateMonitor) SetCallBackFn(f func(proposalId common.Hash, sender *apicommonpb.TaskOrganization, orgState *OrgProposalState)) {
	psm.fn = f
}
func (psm *ProposalStateMonitor) SetWhen(when int64) { psm.when = when }
func (psm *ProposalStateMonitor) SetNext(next int64) { psm.next = next }

type proposalStateMonitorQueue []*ProposalStateMonitor

type SyncProposalStateMonitorQueue struct {
	lock  sync.Mutex
	timer *time.Timer
	queue *proposalStateMonitorQueue
}

func NewSyncProposalStateMonitorQueue(size int) *SyncProposalStateMonitorQueue {
	queue := make(proposalStateMonitorQueue, size)
	timer := time.NewTimer(0)
	<- timer.C
	return &SyncProposalStateMonitorQueue{
		queue: &(queue),
		timer: timer,
	}
}

func (syncQueue *SyncProposalStateMonitorQueue) Len () int {
	syncQueue.lock.Lock()
	defer syncQueue.lock.Unlock()
	return len(*(syncQueue.queue))
}

func (syncQueue *SyncProposalStateMonitorQueue) Timer() *time.Timer {
	return syncQueue.timer
}

func (syncQueue *SyncProposalStateMonitorQueue) CheckMonitors(now int64) int64 {

	syncQueue.lock.Lock()
	defer syncQueue.lock.Unlock()
	// Note that runMonitor may temporarily unlock queue.Lock.
rerun:
	for len(*(syncQueue.queue)) > 0 {
		if future := syncQueue.runMonitor(now); future != 0 {
			if future > 0 {
				now = timeutils.UnixMsec()
				if future > now {
					return future
				} else {
					continue rerun
				}
			}
		}
	}
	return math.MaxInt32
}

func (syncQueue *SyncProposalStateMonitorQueue) Size() int { return len(*(syncQueue.queue)) }

func (syncQueue *SyncProposalStateMonitorQueue) AddMonitor(m *ProposalStateMonitor) {
	syncQueue.lock.Lock()
	defer syncQueue.lock.Unlock()
	// when must never be negative;
	if m.when-timeutils.UnixMsec() < 0 {
		panic("target time is negative number")
	}
	i := len(*(syncQueue.queue))
	*(syncQueue.queue) = append(*(syncQueue.queue), m)
	syncQueue.siftUpMonitor(i)

	// reset the timer
	var until int64
	if len(*(syncQueue.queue)) > 0 {
		until =  (*(syncQueue.queue))[0].when
	} else {
		until = -1
	}
	future := time.Duration(until - timeutils.UnixMsec())
	if future <= 0 {
		future = 0
	}
	syncQueue.timer.Reset(future * time.Millisecond)
}

func (syncQueue *SyncProposalStateMonitorQueue) UpdateMonitor(proposalId common.Hash, partyId string, when, next int64) {
	syncQueue.lock.Lock()
	defer syncQueue.lock.Unlock()
	// when must never be negative;
	if when-timeutils.UnixMsec() < 0 {
		panic("target time is negative number")
	}
	for i := 0; i < len(*(syncQueue.queue)); i++ {
		m := (*(syncQueue.queue))[i]
		if m.GetProposalId() == proposalId && m.GetPartyId() == partyId {
			// Leave in heap but adjust next time to fire.
			m.when = when
			m.next = next // maybe is zero value.
			syncQueue.siftDownMonitor(0)
			return
		}
	}
}

func (syncQueue *SyncProposalStateMonitorQueue) DelMonitor(proposalId common.Hash, partyId string) {
	syncQueue.lock.Lock()
	defer syncQueue.lock.Unlock()

	for i := 0; i < len(*(syncQueue.queue)); i++ {
		m := (*(syncQueue.queue))[i]
		if m.GetProposalId() == proposalId && m.GetPartyId() == partyId {
			syncQueue.delMonitorWithIndex(i)
			return
		}
	}
}

func (syncQueue *SyncProposalStateMonitorQueue) delMonitorWithIndex(i int) {

	last := len(*(syncQueue.queue)) - 1
	if i != last {
		(*(syncQueue.queue))[i] = (*(syncQueue.queue))[last]
	}
	(*(syncQueue.queue))[last] = nil
	*(syncQueue.queue) = (*(syncQueue.queue))[:last]
	if i != last {
		// Moving to i may have moved the last monitor to a new parent,
		// so sift up to preserve the heap guarantee.
		syncQueue.siftUpMonitor(i)
		syncQueue.siftDownMonitor(i)
	}
}

func (syncQueue *SyncProposalStateMonitorQueue) delMonitor0() {

	last := len(*(syncQueue.queue)) - 1
	if last > 0 {
		(*(syncQueue.queue))[0] = (*(syncQueue.queue))[last]
	}
	(*(syncQueue.queue))[last] = nil
	*(syncQueue.queue) = (*(syncQueue.queue))[:last]

	if last > 0 {
		syncQueue.siftDownMonitor(0)
	}
}

// NOTE: runMonitor() must be used in a logic between calling lock() and unlock().
func (syncQueue *SyncProposalStateMonitorQueue) runMonitor(now int64) int64 {

	if len(*(syncQueue.queue)) == 0 {
		return 0
	}

	m := (*(syncQueue.queue))[0]
	if m.when > now {
		// Not ready to run.
		return m.when
	}
	f := m.fn
	proposalId := m.proposalId
	sender := m.sender
	orgState := m.orgState

	if m.next > 0 {
		// when must never be negative;
		if m.next-timeutils.UnixMsec() < 0 {
			panic("target time is negative number")
		}
		// Leave in heap but adjust next time to fire.
		m.when = m.next
		m.next = 0 // NOTE: clean old next time, a new next time will set in the callback func (the monitor field `fn`)
		syncQueue.siftDownMonitor(0)

	} else {
		// Remove top member from heap.
		syncQueue.delMonitor0()
	}

	syncQueue.lock.Unlock()
	f(proposalId, sender, orgState)
	syncQueue.lock.Lock()
	return 0
}

func (syncQueue *SyncProposalStateMonitorQueue) TimeSleepUntil() int64 {
	syncQueue.lock.Lock()
	defer syncQueue.lock.Unlock()
	if len(*(syncQueue.queue)) > 0 {
		return (*(syncQueue.queue))[0].when
	} else {
		return -1
	}
}

func (syncQueue *SyncProposalStateMonitorQueue) siftUpMonitor(i int) {

	if i >= len(*(syncQueue.queue)) {
		panic("queue data corruption")
	}
	when := (*(syncQueue.queue))[i].when
	tmp := (*(syncQueue.queue))[i]
	for i > 0 {

		p := (i - 1) / 4 // parent
		if when >= (*(syncQueue.queue))[p].when {
			break
		}
		(*(syncQueue.queue))[i] = (*(syncQueue.queue))[p]
		i = p
	}
	if tmp != (*(syncQueue.queue))[i] {
		(*(syncQueue.queue))[i] = tmp
	}
}

func (syncQueue *SyncProposalStateMonitorQueue) siftDownMonitor(i int) {

	n := len(*(syncQueue.queue))
	if i >= n {
		panic("queue data corruption")
	}
	when := (*(syncQueue.queue))[i].when
	tmp := (*(syncQueue.queue))[i]
	for {
		c := i*4 + 1 // left child
		c3 := c + 2  // mid child
		if c >= n {
			break
		}
		w := (*(syncQueue.queue))[c].when
		if c+1 < n && (*(syncQueue.queue))[c+1].when < w {
			w = (*(syncQueue.queue))[c+1].when
			c++
		}
		if c3 < n {
			w3 := (*(syncQueue.queue))[c3].when
			if c3+1 < n && (*(syncQueue.queue))[c3+1].when < w3 {
				w3 = (*(syncQueue.queue))[c3+1].when
				c3++
			}
			if w3 < w {
				w = w3
				c = c3
			}
		}
		if w >= when {
			break
		}
		(*(syncQueue.queue))[i] = (*(syncQueue.queue))[c]
		i = c
	}
	if tmp != (*(syncQueue.queue))[i] {
		(*(syncQueue.queue))[i] = tmp
	}
}
