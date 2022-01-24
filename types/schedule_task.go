package types

import (
	"bytes"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	msgcommonpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/common"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	"github.com/libp2p/go-libp2p-core/peer"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"
)

type TaskActionStatus uint16

func (t TaskActionStatus) String() string {
	switch t {
	case TaskConsensusFinished:
		return "task consensus finished"
	case TaskConsensusInterrupt:
		return "task consensus interrupted"
	case TaskTerminate:
		return "task terminated"
	case TaskNeedExecute:
		return "task need execute"
	case TaskScheduleFailed:
		return "task schedule failed"
	default:
		return "unknown task result status"
	}
}

const (
	TaskConsensusFinished  TaskActionStatus = 0x0000
	TaskConsensusInterrupt TaskActionStatus = 0x0001
	TaskTerminate          TaskActionStatus = 0x0010 // terminate task while consensus or executing
	TaskNeedExecute        TaskActionStatus = 0x0100
	TaskScheduleFailed     TaskActionStatus = 0x1000 // schedule failed final
)

type TaskConsResult struct {
	TaskId string
	Status TaskActionStatus
	Err    error
}

func NewTaskConsResult(taskId string, status TaskActionStatus, err error) *TaskConsResult {
	return &TaskConsResult{
		TaskId: taskId,
		Status: status,
		Err:    err,
	}
}
func (res *TaskConsResult) GetTaskId() string           { return res.TaskId }
func (res *TaskConsResult) GetStatus() TaskActionStatus { return res.Status }
func (res *TaskConsResult) GetErr() error               { return res.Err }
func (res *TaskConsResult) String() string {
	return fmt.Sprintf(`{"taskId": %s, "status": %s, "err": %s}`, res.TaskId, res.Status.String(), res.Err)
}

// ================================================= V2.0 =================================================

// Local tasks that need to be agreed (scheduled but not yet agreed)
type NeedConsensusTask struct {
	task       *Task
	nonce      []byte
	weights    [][]byte
	electionAt uint64
}

func NewNeedConsensusTask(task *Task, nonce []byte, weights [][]byte, electionAt uint64) *NeedConsensusTask {
	return &NeedConsensusTask{
		task:       task,
		nonce:      nonce,
		weights:    weights,
		electionAt: electionAt,
	}
}

func (nct *NeedConsensusTask) GetTask() *Task        { return nct.task }
func (nct *NeedConsensusTask) GetNonce() []byte      { return nct.nonce }
func (nct *NeedConsensusTask) GetWeights() [][]byte  { return nct.weights }
func (nct *NeedConsensusTask) GetElectionAt() uint64 { return nct.electionAt }
func (nct *NeedConsensusTask) String() string {
	taskStr := "{}"
	if nil != nct.task {
		taskStr = nct.task.GetTaskData().String()
	}

	nonceStr := "0x"
	if len(nct.nonce) != 0 {
		nonceStr = common.BytesToHash(nct.nonce).Hex()
	}

	weightsStr := "[]"
	if len(nct.weights) != 0 {
		arr := make([]string, len(nct.weights))
		for i, weight := range nct.weights {
			arr[i] = new(big.Int).SetBytes(weight).String()
		}
		weightsStr = "[" + strings.Join(arr, ",") + "]"
	}
	return fmt.Sprintf(`{"task": %s, "nonce": %s, "weights": %s, "electionAt": %d}`, taskStr, nonceStr, weightsStr, nct.electionAt)
}

// Remote tasks that need to be scheduled again
// (those that are in the process of consensus and need to be scheduled again after receiving the proposal from the opposite end)
type NeedReplayScheduleTask struct {
	localTaskRole apicommonpb.TaskRole
	localPartyId  string
	task          *Task
	nonce         []byte
	weights       [][]byte
	electionAt    uint64
	resultCh      chan *ReplayScheduleResult
}

func NewNeedReplayScheduleTask(role apicommonpb.TaskRole, partyId string, task *Task, nonce []byte, weights [][]byte, electionAt uint64) *NeedReplayScheduleTask {
	return &NeedReplayScheduleTask{
		localTaskRole: role,
		localPartyId:  partyId,
		task:          task,
		nonce:         nonce,
		weights:       weights,
		electionAt:    electionAt,
		resultCh:      make(chan *ReplayScheduleResult),
	}
}
func (nrst *NeedReplayScheduleTask) SendFailedResult(taskId string, err error) {
	nrst.SendResult(&ReplayScheduleResult{
		taskId: taskId,
		err:    err,
	})
}
func (nrst *NeedReplayScheduleTask) SendResult(result *ReplayScheduleResult) {
	nrst.resultCh <- result
	close(nrst.resultCh)
}
func (nrst *NeedReplayScheduleTask) ReceiveResult() *ReplayScheduleResult {
	return <-nrst.resultCh
}
func (nrst *NeedReplayScheduleTask) GetLocalTaskRole() apicommonpb.TaskRole  { return nrst.localTaskRole }
func (nrst *NeedReplayScheduleTask) GetLocalPartyId() string                 { return nrst.localPartyId }
func (nrst *NeedReplayScheduleTask) GetTask() *Task                          { return nrst.task }
func (nrst *NeedReplayScheduleTask) GetNonce() []byte                        { return nrst.nonce }
func (nrst *NeedReplayScheduleTask) GetWeights() [][]byte                    { return nrst.weights }
func (nrst *NeedReplayScheduleTask) GetElectionAt() uint64                   { return nrst.electionAt }
func (nrst *NeedReplayScheduleTask) GetResultCh() chan *ReplayScheduleResult { return nrst.resultCh }
func (nrst *NeedReplayScheduleTask) String() string {
	taskStr := "{}"
	if nil != nrst.task {
		taskStr = nrst.task.GetTaskData().String()
	}
	nonceStr := "0x"
	if len(nrst.nonce) != 0 {
		nonceStr = common.BytesToHash(nrst.nonce).Hex()
	}

	weightsStr := "[]"
	if len(nrst.weights) != 0 {
		arr := make([]string, len(nrst.weights))
		for i, weight := range nrst.weights {
			arr[i] = new(big.Int).SetBytes(weight).String()
		}
		weightsStr = "[" + strings.Join(arr, ",") + "]"
	}
	return fmt.Sprintf(`{"taskRole": %s, "localPartyId": %s, "task": %s, "nonce": %s, "weights": %s, "resultCh": %p}`,
		nrst.localTaskRole.String(), nrst.localPartyId, taskStr, nonceStr, weightsStr, nrst.resultCh)
}

type ReplayScheduleResult struct {
	taskId   string
	err      error
	resource *PrepareVoteResource
}

func NewReplayScheduleResult(taskId string, err error, resource *PrepareVoteResource) *ReplayScheduleResult {
	return &ReplayScheduleResult{
		taskId:   taskId,
		err:      err,
		resource: resource,
	}
}
func (rsr *ReplayScheduleResult) GetTaskId() string                 { return rsr.taskId }
func (rsr *ReplayScheduleResult) GetErr() error                     { return rsr.err }
func (rsr *ReplayScheduleResult) GetResource() *PrepareVoteResource { return rsr.resource }
func (rsr *ReplayScheduleResult) String() string {
	errStr := "nil"
	if nil != rsr.err {
		errStr = rsr.err.Error()
	}
	resourceStr := "{}"
	if nil != rsr.resource {
		resourceStr = rsr.resource.String()
	}
	return fmt.Sprintf(`{"taskId": %s, "err": %s, "resource": %s}`,
		rsr.taskId, errStr, resourceStr)
}

// Tasks to be executed (local and remote, which have been completed by consensus and can be executed by issuing fighter)
type NeedExecuteTask struct {
	remotepid              peer.ID
	localTaskRole          apicommonpb.TaskRole
	localTaskOrganization  *apicommonpb.TaskOrganization
	remoteTaskRole         apicommonpb.TaskRole
	remoteTaskOrganization *apicommonpb.TaskOrganization
	consStatus             TaskActionStatus
	localResource          *PrepareVoteResource
	resources              *twopcpb.ConfirmTaskPeerInfo
	taskId                 string
	err                    error
}

func NewNeedExecuteTask(
	remotepid peer.ID,
	localTaskRole, remoteTaskRole apicommonpb.TaskRole,
	localTaskOrganization, remoteTaskOrganization *apicommonpb.TaskOrganization,
	taskId string,
	consStatus TaskActionStatus,
	localResource *PrepareVoteResource,
	resources *twopcpb.ConfirmTaskPeerInfo,
	err error,
) *NeedExecuteTask {
	return &NeedExecuteTask{
		remotepid:              remotepid,
		localTaskRole:          localTaskRole,
		localTaskOrganization:  localTaskOrganization,
		remoteTaskRole:         remoteTaskRole,
		remoteTaskOrganization: remoteTaskOrganization,
		taskId:                 taskId,
		consStatus:             consStatus,
		localResource:          localResource,
		resources:              resources,
		err:                    err,
	}
}
func (net *NeedExecuteTask) HasRemotePID() bool                      { return strings.Trim(string(net.remotepid), "") != "" }
func (net *NeedExecuteTask) HasEmptyRemotePID() bool                 { return !net.HasRemotePID() }
func (net *NeedExecuteTask) GetRemotePID() peer.ID                   { return net.remotepid }
func (net *NeedExecuteTask) GetLocalTaskRole() apicommonpb.TaskRole  { return net.localTaskRole }
func (net *NeedExecuteTask) GetRemoteTaskRole() apicommonpb.TaskRole { return net.remoteTaskRole }
func (net *NeedExecuteTask) GetLocalTaskOrganization() *apicommonpb.TaskOrganization {
	return net.localTaskOrganization
}
func (net *NeedExecuteTask) GetRemoteTaskOrganization() *apicommonpb.TaskOrganization {
	return net.remoteTaskOrganization
}
func (net *NeedExecuteTask) GetTaskId() string                          { return net.taskId }
func (net *NeedExecuteTask) GetConsStatus() TaskActionStatus            { return net.consStatus }
func (net *NeedExecuteTask) GetLocalResource() *PrepareVoteResource     { return net.localResource }
func (net *NeedExecuteTask) GetResources() *twopcpb.ConfirmTaskPeerInfo { return net.resources }
func (net *NeedExecuteTask) GetErr() error                              { return net.err }
func (net *NeedExecuteTask) String() string {
	localIdentityStr := "{}"
	if nil != net.GetLocalTaskOrganization() {
		localIdentityStr = net.GetLocalTaskOrganization().String()
	}
	remoteIdentityStr := "{}"
	if nil != net.GetRemoteTaskOrganization() {
		remoteIdentityStr = net.GetRemoteTaskOrganization().String()
	}
	localResourceStr := "{}"
	if nil != net.GetLocalResource() {
		localResourceStr = net.GetLocalResource().String()
	}
	return fmt.Sprintf(`{"remotepid": %s, "localTaskRole": %s, "localTaskOrganization": %s, "remoteTaskRole": %s, "remoteTaskOrganization": %s, "taskId": %s, "localResource": %s, "resources": %s, "err": %s}`,
		net.GetRemotePID(), net.GetLocalTaskRole().String(), localIdentityStr, net.GetRemoteTaskRole().String(), remoteIdentityStr, net.GetTaskId(), localResourceStr, ConfirmTaskPeerInfoString(net.GetResources()), net.GetErr())
}

type ExecuteTaskMonitor struct {
	taskId  string
	partyId string
	when    int64 // target timestamp
	index   int
	fn      func()
}

func NewExecuteTaskMonitor(taskId, partyId string, when int64, fn func()) *ExecuteTaskMonitor {
	return &ExecuteTaskMonitor{
		taskId:  taskId,
		partyId: partyId,
		when:    when,
		fn:      fn,
	}
}


func (etm *ExecuteTaskMonitor) String() string      {
	return fmt.Sprintf(`{"index": %d, "taskId": %s, "partyId": %s, "when": %d}`,
		etm.GetIndex(), etm.GetTaskId(), etm.GetPartyId(), etm.GetWhen())
}
func (etm *ExecuteTaskMonitor) GetIndex() int      { return etm.index }
func (etm *ExecuteTaskMonitor) GetTaskId() string  { return etm.taskId }
func (etm *ExecuteTaskMonitor) GetPartyId() string { return etm.partyId }
func (etm *ExecuteTaskMonitor) GetWhen() int64     { return etm.when }

type executeTaskMonitorQueue []*ExecuteTaskMonitor

func (queue *executeTaskMonitorQueue) String() string {
	arr := make([]string, len(*queue))
	for i, ett := range *queue {
		arr[i] = ett.String()
	}
	return "[" + strings.Join(arr, ",") +  "]"
}

type SyncExecuteTaskMonitorQueue struct {
	lock  sync.Mutex
	timer *time.Timer
	queue *executeTaskMonitorQueue
}

func NewSyncExecuteTaskMonitorQueue(size int) *SyncExecuteTaskMonitorQueue {
	queue := make(executeTaskMonitorQueue, size)
	timer := time.NewTimer(0)
	<-timer.C
	return &SyncExecuteTaskMonitorQueue{
		queue: &(queue),
		timer: timer,
	}
}

func (syncQueue *SyncExecuteTaskMonitorQueue) QueueString() string {
	syncQueue.lock.Lock()
	defer syncQueue.lock.Unlock()
	return syncQueue.queue.String()
}

func (syncQueue *SyncExecuteTaskMonitorQueue) Len() int {
	syncQueue.lock.Lock()
	defer syncQueue.lock.Unlock()
	return len(*(syncQueue.queue))
}

func (syncQueue *SyncExecuteTaskMonitorQueue) Timer() *time.Timer {
	return syncQueue.timer
}

func (syncQueue *SyncExecuteTaskMonitorQueue) CheckMonitors(now int64) int64 {

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

func (syncQueue *SyncExecuteTaskMonitorQueue) Size() int { return len(*(syncQueue.queue)) }

func (syncQueue *SyncExecuteTaskMonitorQueue) AddMonitor(m *ExecuteTaskMonitor) {
	syncQueue.lock.Lock()
	defer syncQueue.lock.Unlock()
	// when must never be negative;
	if m.when-timeutils.UnixMsec() < 0 {
		log.Warnf("target time is negative number, taskId: %s, partyId: %s, when: %d, now: %d",
			 m.GetTaskId(), m.GetPartyId(), m.when, timeutils.UnixMsec())
	}
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

	log.Debugf("Add needExecuteTask monitor, taskId: {%s}, partyId: {%s}, when: {%d}, now: {%d}",
		m.GetTaskId(), m.GetPartyId(), m.GetWhen(), timeutils.UnixMsec())
}

func (syncQueue *SyncExecuteTaskMonitorQueue) DelMonitor(taskId, partyId string) {
	syncQueue.lock.Lock()
	defer syncQueue.lock.Unlock()

	for i := 0; i < len(*(syncQueue.queue)); i++ {
		m := (*(syncQueue.queue))[i]
		if m.GetTaskId() == taskId && m.GetPartyId() == partyId {
			syncQueue.delMonitorWithIndex(i)
			log.Debugf("Delete needExecuteTask monitor, taskId: {%s}, partyId: {%s}, when: {%d}, now: {%d}, index: {%d}",
				m.GetTaskId(), m.GetPartyId(), m.GetWhen(), timeutils.UnixMsec(), i)
			return
		}
	}
}

func (syncQueue *SyncExecuteTaskMonitorQueue) delMonitorWithIndex(i int) {

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

func (syncQueue *SyncExecuteTaskMonitorQueue) delMonitor0() {

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
func (syncQueue *SyncExecuteTaskMonitorQueue) runMonitor(now int64) int64 {

	if len(*(syncQueue.queue)) == 0 {
		return 0
	}

	m := (*(syncQueue.queue))[0]
	if m.when > now {
		// Not ready to run.
		return m.when
	}
	f := m.fn
	// Remove top member from heap.
	syncQueue.delMonitor0()
	log.Debugf("Delete heap top0 needExecuteTask monitor, taskId: {%s}, partyId: {%s}, when: {%d}, now: {%d}",
		m.GetTaskId(), m.GetPartyId(), m.GetWhen(), timeutils.UnixMsec())

	syncQueue.lock.Unlock()
	f()
	syncQueue.lock.Lock()
	return 0
}

func (syncQueue *SyncExecuteTaskMonitorQueue) TimeSleepUntil() int64 {
	syncQueue.lock.Lock()
	defer syncQueue.lock.Unlock()
	if len(*(syncQueue.queue)) > 0 {
		return (*(syncQueue.queue))[0].when
	} else {
		return -1
	}
}

func (syncQueue *SyncExecuteTaskMonitorQueue) siftUpMonitor(i int) {

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

func (syncQueue *SyncExecuteTaskMonitorQueue) siftDownMonitor(i int) {

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

func ConfirmTaskPeerInfoString(resources *twopcpb.ConfirmTaskPeerInfo) string {
	if nil == resources {
		return "{}"
	}
	dataSupplierList := make([]string, len(resources.GetDataSupplierPeerInfos()))
	for i, peerInfo := range resources.GetDataSupplierPeerInfos() {
		var resource *PrepareVoteResource
		if nil == peerInfo {
			resource = &PrepareVoteResource{}
		} else {
			resource = &PrepareVoteResource{
				Ip:      string(peerInfo.Ip),
				Port:    string(peerInfo.Port),
				PartyId: string(peerInfo.PartyId),
			}
		}
		dataSupplierList[i] = resource.String()
	}
	dataSupplierListStr := "[" + strings.Join(dataSupplierList, ",") + "]"

	powerSupplierList := make([]string, len(resources.GetPowerSupplierPeerInfos()))
	for i, peerInfo := range resources.GetPowerSupplierPeerInfos() {
		var resource *PrepareVoteResource
		if nil == peerInfo {
			resource = &PrepareVoteResource{}
		} else {
			resource = &PrepareVoteResource{
				Ip:      string(peerInfo.Ip),
				Port:    string(peerInfo.Port),
				PartyId: string(peerInfo.PartyId),
			}
		}
		powerSupplierList[i] = resource.String()
	}
	powerSupplierListStr := "[" + strings.Join(powerSupplierList, ",") + "]"

	receiverList := make([]string, len(resources.GetResultReceiverPeerInfos()))
	for i, peerInfo := range resources.GetResultReceiverPeerInfos() {
		var resource *PrepareVoteResource
		if nil == peerInfo {
			resource = &PrepareVoteResource{}
		} else {
			resource = &PrepareVoteResource{
				Ip:      string(peerInfo.Ip),
				Port:    string(peerInfo.Port),
				PartyId: string(peerInfo.PartyId),
			}
		}
		receiverList[i] = resource.String()
	}
	receiverListStr := "[" + strings.Join(receiverList, ",") + "]"

	return fmt.Sprintf(`{"dataSupplierPeerInfoList": %s, "powerSupplierPeerInfoList": %s, "resultReceiverPeerInfoList": %s}`,
		dataSupplierListStr, powerSupplierListStr, receiverListStr)
}

func IsSameTaskOrgByte(org1, org2 *msgcommonpb.TaskOrganizationIdentityInfo) bool {
	if bytes.Compare(org1.GetPartyId(), org2.GetPartyId()) == 0 && bytes.Compare(org1.GetIdentityId(), org2.GetIdentityId()) == 0 {
		return true
	}
	return false
}
func IsNotSameTaskOrgByte(org1, org2 *msgcommonpb.TaskOrganizationIdentityInfo) bool {
	return !IsSameTaskOrgByte(org1, org2)
}

func IsSameTaskOrgParty(org1, org2 *apicommonpb.TaskOrganization) bool {
	if org1.GetPartyId() == org2.GetPartyId() && org1.GetIdentityId() == org2.GetIdentityId() {
		return true
	}
	return false
}
func IsNotSameTaskOrgParty(org1, org2 *apicommonpb.TaskOrganization) bool {
	return !IsSameTaskOrgParty(org1, org2)
}

func IsSameTaskOrg(org1, org2 *apicommonpb.TaskOrganization) bool {
	if org1.GetIdentityId() == org2.GetIdentityId() {
		return true
	}
	return false
}
func IsNotSameTaskOrg(org1, org2 *apicommonpb.TaskOrganization) bool {
	return !IsSameTaskOrg(org1, org2)
}
