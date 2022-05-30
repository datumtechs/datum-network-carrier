package types

import (
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"strings"
	"sync"
	"time"
)

type ProposalTask struct {
	ProposalId common.Hash
	TaskId     string
	CreateAt   uint64
}

func NewProposalTask(proposalId common.Hash, taskId string, createAt uint64) *ProposalTask {
	return &ProposalTask{
		ProposalId: proposalId,
		TaskId:     taskId,
		CreateAt:   createAt,
	}
}

func (pt *ProposalTask) GetProposalId() common.Hash { return pt.ProposalId }
func (pt *ProposalTask) GetTaskId() string          { return pt.TaskId }
func (pt *ProposalTask) GetCreateAt() uint64        { return pt.CreateAt }

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
	// during 10s, if the proposal haven't been done, kill it
	ProposalDeadlineDuration = uint64(10 * (time.Second.Milliseconds()))
)

type OrgProposalState struct {
	proposalId       common.Hash
	deadlineDuration uint64 // Clear `ProposalState` , when the current time is greater than the `DeadlineDuration` createAt of proposalState
	createAt         uint64 // the time is that the proposal state was created on current iden
	startAt          uint64 // the time is that the proposal state was created on task sender iden
	taskId           string
	taskRole         carriertypespb.TaskRole
	taskSender       *carriertypespb.TaskOrganization
	taskOrg          *carriertypespb.TaskOrganization
	periodNum        ProposalStatePeriod
}

func NewOrgProposalState(
	proposalId common.Hash,
	taskId string,
	taskRole carriertypespb.TaskRole,
	taskSender *carriertypespb.TaskOrganization,
	taskOrg *carriertypespb.TaskOrganization,
	startAt uint64,
) *OrgProposalState {

	return &OrgProposalState{
		proposalId:       proposalId,
		taskId:           taskId,
		taskRole:         taskRole,
		taskSender:       taskSender,
		taskOrg:          taskOrg,
		periodNum:        PeriodPrepare,
		deadlineDuration: ProposalDeadlineDuration,
		createAt:         timeutils.UnixMsecUint64(),
		startAt:          startAt,
	}
}


func NewOrgProposalStateWithFields(
	proposalId common.Hash,
	taskId string,
	taskRole carriertypespb.TaskRole,
	taskSender *carriertypespb.TaskOrganization,
	taskOrg *carriertypespb.TaskOrganization,
	periodNum ProposalStatePeriod,
	deadlineDuration uint64,
	createAt uint64,
	startAt uint64,
) *OrgProposalState {

	return &OrgProposalState{
		proposalId:       proposalId,
		taskId:           taskId,
		taskRole:         taskRole,
		taskSender:       taskSender,
		taskOrg:          taskOrg,
		periodNum:        periodNum,
		deadlineDuration: deadlineDuration,
		createAt:         createAt,
		startAt:          startAt,
	}
}


func (pstate *OrgProposalState) String() string {
	return fmt.Sprintf(`{"proposalId": %s, "taskId": %s, "taskRole": %s, "taskSender": %s, "taskOrg": %s, "periodNum": %s, "deadlineDuration": %d, "startAt": %d, "createAt": %d}`,
		pstate.GetProposalId().String(), pstate.GetTaskId(), pstate.GetTaskRole().String(), pstate.GetTaskSender().String(), pstate.GetTaskOrg().String(), pstate.GetPeriodStr(), pstate.GetDeadlineDuration(), pstate.GetStartAt(), pstate.GetCreateAt())
}
func (pstate *OrgProposalState) GetProposalId() common.Hash        { return pstate.proposalId }
func (pstate *OrgProposalState) GetTaskId() string                 { return pstate.taskId }
func (pstate *OrgProposalState) GetTaskRole() carriertypespb.TaskRole { return pstate.taskRole }
func (pstate *OrgProposalState) GetTaskSender() *carriertypespb.TaskOrganization {
	return pstate.taskSender
}
func (pstate *OrgProposalState) GetTaskOrg() *carriertypespb.TaskOrganization { return pstate.taskOrg }
func (pstate *OrgProposalState) GetPeriodNum() ProposalStatePeriod         { return pstate.periodNum }
func (pstate *OrgProposalState) GetDeadlineDuration() uint64               { return pstate.deadlineDuration }
func (pstate *OrgProposalState) GetCreateAt() uint64                       { return pstate.createAt }
func (pstate *OrgProposalState) GetStartAt() uint64                        { return pstate.startAt }
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
	switch pstate.periodNum {
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

func (pstate *OrgProposalState) IsPreparePeriod() bool     { return pstate.periodNum == PeriodPrepare }
func (pstate *OrgProposalState) IsConfirmPeriod() bool     { return pstate.periodNum == PeriodConfirm }
func (pstate *OrgProposalState) IsCommitPeriod() bool      { return pstate.periodNum == PeriodCommit }
func (pstate *OrgProposalState) IsFinishedPeriod() bool    { return pstate.periodNum == PeriodFinished }
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

func printTime(loghead string, pstate *OrgProposalState) {
	now := time.Now()
	log.Debugf("%s: now {%d<==>%s}, startAt: {%d<==>%s}, duration: %d ms, taskId: {%s}, partyId: {%s}, proposalId: {%s}",
		loghead, now.UnixNano()/1e6, now.Format("2006-01-02 15:04:05"),
		pstate.GetStartAt(), time.Unix(int64(pstate.GetStartAt())/1000, 0).Format("2006-01-02 15:04:05"),
		now.UnixNano()/1e6-int64(pstate.GetStartAt()), pstate.GetTaskId(), pstate.GetTaskOrg().GetPartyId(), pstate.GetProposalId().String(),
	)
}

func (pstate *OrgProposalState) ChangeToConfirm() {
	if pstate.periodNum == PeriodPrepare {
		pstate.periodNum = PeriodConfirm
		printTime("Change Prepare To Confirm", pstate)
	}
}

func (pstate *OrgProposalState) ChangeToCommit() {
	if pstate.periodNum == PeriodConfirm {
		pstate.periodNum = PeriodCommit
		printTime("Change Confirm To Commit", pstate)
	}
}
func (pstate *OrgProposalState) ChangeToFinished() {
	if pstate.periodNum == PeriodCommit {
		pstate.periodNum = PeriodFinished
		printTime("Change Commit To Finished", pstate)
	}
}

// v 0.3.0 proposal state expire monitor
type ProposalStateMonitor struct {
	proposalId common.Hash
	partyId    string
	sender     *carriertypespb.TaskOrganization
	orgState   *OrgProposalState
	when       int64 // target timestamp
	next       int64
	index      int
	fn         func(orgState *OrgProposalState)
}

func NewProposalStateMonitor(orgState *OrgProposalState, when, next int64,
	fn func(orgState *OrgProposalState)) *ProposalStateMonitor {

	return &ProposalStateMonitor{
		proposalId: orgState.GetProposalId(),
		partyId:    orgState.GetTaskOrg().GetPartyId(),
		sender:     orgState.GetTaskSender(),
		orgState:   orgState,
		when:       when,
		next:       next,
		fn:         fn,
	}
}

func (psm *ProposalStateMonitor) String() string {
	return fmt.Sprintf(`{"index": %d, "proposalId":, %s, "partyId": %s, "sender": %s, "orgState": %s, "when": %d, "next": %d}`,
		psm.GetIndex(), psm.GetProposalId().String(), psm.GetPartyId(), psm.GetTaskSender().String(), psm.GetOrgState().String(), psm.GetWhen(), psm.GetNext())
}
func (psm *ProposalStateMonitor) GetIndex() int                                { return psm.index }
func (psm *ProposalStateMonitor) GetProposalId() common.Hash                   { return psm.proposalId }
func (psm *ProposalStateMonitor) GetPartyId() string                           { return psm.partyId }
func (psm *ProposalStateMonitor) GetTaskSender() *carriertypespb.TaskOrganization { return psm.sender }
func (psm *ProposalStateMonitor) GetOrgState() *OrgProposalState               { return psm.orgState }
func (psm *ProposalStateMonitor) GetWhen() int64                               { return psm.when }
func (psm *ProposalStateMonitor) GetNext() int64                               { return psm.next }
func (psm *ProposalStateMonitor) SetCallBackFn(f func(orgState *OrgProposalState)) {
	psm.fn = f
}
func (psm *ProposalStateMonitor) SetWhen(when int64) { psm.when = when }
func (psm *ProposalStateMonitor) SetNext(next int64) { psm.next = next }

type proposalStateMonitorQueue []*ProposalStateMonitor

func (queue *proposalStateMonitorQueue) String() string {
	arr := make([]string, len(*queue))
	for i, ett := range *queue {
		arr[i] = ett.String()
	}
	return "[" + strings.Join(arr, ",") + "]"
}

type SyncProposalStateMonitorQueue struct {
	lock  sync.Mutex
	timer *time.Timer
	queue *proposalStateMonitorQueue
}

func NewSyncProposalStateMonitorQueue(size int) *SyncProposalStateMonitorQueue {
	queue := make(proposalStateMonitorQueue, size)
	timer := time.NewTimer(0)
	<-timer.C
	return &SyncProposalStateMonitorQueue{
		queue: &(queue),
		timer: timer,
	}
}

func (syncQueue *SyncProposalStateMonitorQueue) QueueString() string {
	syncQueue.lock.Lock()
	defer syncQueue.lock.Unlock()
	return syncQueue.queue.String()
}

func (syncQueue *SyncProposalStateMonitorQueue) Len() int {
	syncQueue.lock.Lock()
	defer syncQueue.lock.Unlock()
	return len(*(syncQueue.queue))
}

func (syncQueue *SyncProposalStateMonitorQueue) Timer() *time.Timer {
	return syncQueue.timer
}

func (syncQueue *SyncProposalStateMonitorQueue) CheckMonitors(now int64, syncCall bool) int64 {

	syncQueue.lock.Lock()
	defer syncQueue.lock.Unlock()
	// Note that runMonitor may temporarily unlock queue.Lock.
rerun:
	for len(*(syncQueue.queue)) > 0 {
		if future := syncQueue.runMonitor(now, syncCall); future > 0  {
			now = timeutils.UnixMsec()
			if future > now {
				return future
			} else {
				continue rerun
			}
		}
	}
	// when no one monitor, return 0 duration value
	return 0
}

func (syncQueue *SyncProposalStateMonitorQueue) Size() int { return len(*(syncQueue.queue)) }

func (syncQueue *SyncProposalStateMonitorQueue) AddMonitor(m *ProposalStateMonitor) {

	// when must never be negative;
	if m.GetWhen()-timeutils.UnixMsec() < 0 {
		log.Warnf("Warning add proposalState monitor, target when time is negative number, proposalId: %s, taskId: %s, partyId: %s, when: %d, now: %d",
			m.GetProposalId().String(), m.GetOrgState().GetTaskId(), m.GetPartyId(), m.GetWhen(), timeutils.UnixMsec())
	}

	syncQueue.lock.Lock()
	defer syncQueue.lock.Unlock()

	i := len(*(syncQueue.queue))
	m.index = i
	*(syncQueue.queue) = append(*(syncQueue.queue), m)
	syncQueue.siftUpMonitor(i)

	// reset the timer
	var until int64
	if len(*(syncQueue.queue)) > 0 {
		until = (*(syncQueue.queue))[0].when
	} else {
		until = -1
	}
	future := time.Duration(until - timeutils.UnixMsec())
	if future <= 0 {
		future = 0
	}
	syncQueue.timer.Reset(future * time.Millisecond)
	log.Debugf("Add proposalState monitor, proposalId: {%s}, taskId: {%s}, partyId: {%s}, when: {%d}, next: {%d}, now: {%d}",
		m.GetProposalId().String(), m.GetOrgState().GetTaskId(), m.GetPartyId(), m.GetWhen(), m.GetNext(), timeutils.UnixMsec())
}

func (syncQueue *SyncProposalStateMonitorQueue) UpdateMonitor(proposalId common.Hash, partyId string, when, next int64) {

	// when must never be negative;
	if when-timeutils.UnixMsec() < 0 {
		log.Errorf("Failed to update proposalState monitor, target time is negative number, proposalId: %s, partyId: %s, when: %d, next: %d, now: %d",
			proposalId.String(), partyId, when, next, timeutils.UnixMsec())
		return
	}

	syncQueue.lock.Lock()
	defer syncQueue.lock.Unlock()

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
			log.Debugf("Delete proposalState monitor, proposalId: {%s}, taskId: {%s}, partyId: {%s}, when: {%d}, next: {%d}, now: {%d}, index: {%d}",
				m.GetProposalId().String(), m.GetOrgState().GetTaskId(), m.GetPartyId(), m.GetWhen(), m.GetNext(), timeutils.UnixMsec(), i)
			return
		}
	}
}

func (syncQueue *SyncProposalStateMonitorQueue) delMonitorWithIndex(i int) {

	last := len(*(syncQueue.queue)) - 1
	if i != last {
		(*(syncQueue.queue))[last].index = i
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
		(*(syncQueue.queue))[last].index = 0
		(*(syncQueue.queue))[0] = (*(syncQueue.queue))[last]
	}
	(*(syncQueue.queue))[last] = nil
	*(syncQueue.queue) = (*(syncQueue.queue))[:last]

	if last > 0 {
		syncQueue.siftDownMonitor(0)
	}
}

// NOTE: runMonitor() must be used in a logic between calling lock() and unlock().
func (syncQueue *SyncProposalStateMonitorQueue) runMonitor(now int64, syncCall bool) int64 {

	if len(*(syncQueue.queue)) == 0 {
		return 0
	}

	m := (*(syncQueue.queue))[0]
	if m.when > now {
		// Not ready to run.
		return m.when
	}
	f := m.fn
	orgState := m.orgState

	if m.next > 0 {
		// when must never be negative;
		if m.next-timeutils.UnixMsec() < 0 {
			log.Warnf("target next time is negative number, proposalId: %s, taskId: %s, partyId: %s, next: %d, now: %d",
				m.GetProposalId().String(), m.GetOrgState().GetTaskId(), m.GetPartyId(), m.next, timeutils.UnixMsec())
		}
		// Leave in heap but adjust next time to fire.
		m.when = m.next
		m.next = 0 // NOTE: clean old next time, a new next time will set in the callback func (the monitor field `fn`)
		syncQueue.siftDownMonitor(0)

	} else {
		// Remove top member from heap.
		syncQueue.delMonitor0()
		log.Debugf("Delete heap top0 proposalState monitor, proposalId: {%s}, taskId: {%s}, partyId: {%s}, when: {%d}, next: {%d}, now: {%d}",
			m.GetProposalId().String(), m.GetOrgState().GetTaskId(), m.GetPartyId(), m.GetWhen(), m.GetNext(), timeutils.UnixMsec())
	}

	syncQueue.lock.Unlock()
	if syncCall {
		go f(orgState)
	} else {
		f(orgState)
	}
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
		(*(syncQueue.queue))[p].index = i
		(*(syncQueue.queue))[i] = (*(syncQueue.queue))[p]
		i = p
	}
	if tmp != (*(syncQueue.queue))[i] {
		tmp.index = i
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
		(*(syncQueue.queue))[c].index = i
		(*(syncQueue.queue))[i] = (*(syncQueue.queue))[c]
		i = c
	}
	if tmp != (*(syncQueue.queue))[i] {
		tmp.index = i
		(*(syncQueue.queue))[i] = tmp
	}
}
