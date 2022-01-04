package types

import (
	"bytes"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	msgcommonpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/common"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	"github.com/libp2p/go-libp2p-core/peer"
	"math/big"
	"strings"
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

type ExecuteTaskTimeout struct {
	taskId  string
	partyId string
	when    time.Time
}

func NewExecuteTaskTimeout(taskId, partyId string, when time.Time) *ExecuteTaskTimeout {
	return &ExecuteTaskTimeout{
		taskId:  taskId,
		partyId: partyId,
		when:    when,
	}
}
func (ett *ExecuteTaskTimeout) GetTaskId() string  { return ett.taskId }
func (ett *ExecuteTaskTimeout) GetPartyId() string { return ett.partyId }
func (ett *ExecuteTaskTimeout) GetWhen() time.Time { return ett.when }

type ExecuteTaskTimeoutQueue []*ExecuteTaskTimeout

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
