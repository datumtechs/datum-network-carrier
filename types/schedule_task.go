package types

import (
	"bytes"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	msgcommonpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/common"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	"github.com/libp2p/go-libp2p-core/peer"
	"strings"
)

type ProposalTask struct {
	ProposalId common.Hash
	*Task
	CreateAt uint64
}

func NewProposalTask(proposalId common.Hash, task *Task, createAt uint64) *ProposalTask {
	return &ProposalTask{
		ProposalId: proposalId,
		Task:       task,
		CreateAt:   createAt,
	}
}

func (pt *ProposalTask) GetProposalId() common.Hash { return pt.ProposalId }
func (pt *ProposalTask) GetTask() *Task             { return pt.Task }
func (pt *ProposalTask) GetCreateAt() uint64        { return pt.CreateAt }

type TaskActionStatus uint16

func (t TaskActionStatus) String() string {
	switch t {
	case TaskConsensusFinished:
		return "TaskConsensusFinished"
	case TaskConsensusInterrupt:
		return "TaskConsensusInterrupt"
	case TaskTerminate:
		return "TaskTerminate"
	case TaskNeedExecute:
		return "TaskNeedExecute"
	case TaskExecutingInterrupt:
		return "TaskExecutingInterrupt"

	default:
		return "UnknownTaskResultStatus"
	}
}

const (
	TaskConsensusFinished  TaskActionStatus = 0x0000
	TaskConsensusInterrupt TaskActionStatus = 0x0001
	TaskTerminate          TaskActionStatus = 0x0010 // terminate task while consensus or executing
	TaskNeedExecute        TaskActionStatus = 0x0100
	TaskExecutingInterrupt TaskActionStatus = 0x1000
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

// 需要被 进行共识的 local task (已经调度好的, 还未共识的)
type NeedConsensusTask struct {
	task     *Task
	resultCh chan *TaskConsResult
}

func NewNeedConsensusTask(task *Task) *NeedConsensusTask {
	return &NeedConsensusTask{
		task:     task,
		resultCh: make(chan *TaskConsResult),
	}
}
func (nct *NeedConsensusTask) GetTask() *Task { return nct.task }

func (nct *NeedConsensusTask) GetResultCh() chan *TaskConsResult { return nct.resultCh }
func (nct *NeedConsensusTask) String() string {
	taskStr := "{}"
	if nil != nct.task {
		taskStr = nct.task.GetTaskData().String()
	}
	return fmt.Sprintf(`{"task": %s, "resultCh": %p}`, taskStr, nct.resultCh)
}
func (nct *NeedConsensusTask) SendResult(result *TaskConsResult) {
	nct.resultCh <- result
	close(nct.resultCh)
}
func (nct *NeedConsensusTask) ReceiveResult() *TaskConsResult {
	return <-nct.resultCh
}

// 需要 重演调度的 remote task (接收到对端发来的 proposal 中的, 处于共识过程中的, 需要重演调度的)
type NeedReplayScheduleTask struct {
	localTaskRole apicommonpb.TaskRole
	localPartyId  string
	task          *Task
	resultCh      chan *ReplayScheduleResult
}

func NewNeedReplayScheduleTask(role apicommonpb.TaskRole, partyId string, task *Task) *NeedReplayScheduleTask {
	return &NeedReplayScheduleTask{
		localTaskRole: role,
		localPartyId:  partyId,
		task:          task,
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
func (nrst *NeedReplayScheduleTask) GetResultCh() chan *ReplayScheduleResult { return nrst.resultCh }
func (nrst *NeedReplayScheduleTask) String() string {
	taskStr := "{}"
	if nil != nrst.task {
		taskStr = nrst.task.GetTaskData().String()
	}
	return fmt.Sprintf(`{"taskRole": %s, "localPartyId": %s, "task": %s, "resultCh": %p}`,
		nrst.localTaskRole.String(), nrst.localPartyId, taskStr, nrst.resultCh)
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

// 需要被执行的 task (local 和 remote, 已经被共识完成的, 可以下发 fighter 去执行的)
type NeedExecuteTask struct {
	remotepid              peer.ID
	proposalId             common.Hash
	localTaskRole          apicommonpb.TaskRole
	localTaskOrganization  *apicommonpb.TaskOrganization
	remoteTaskRole         apicommonpb.TaskRole
	remoteTaskOrganization *apicommonpb.TaskOrganization
	task                   *Task
	consStatus             TaskActionStatus
	localResource          *PrepareVoteResource
	resources              *twopcpb.ConfirmTaskPeerInfo
}

func NewNeedExecuteTask(
	remotepid peer.ID,
	proposalId common.Hash,
	localTaskRole apicommonpb.TaskRole,
	localTaskOrganization *apicommonpb.TaskOrganization,
	remoteTaskRole apicommonpb.TaskRole,
	remoteTaskOrganization *apicommonpb.TaskOrganization,
	task *Task,
	consStatus TaskActionStatus,
	localResource *PrepareVoteResource,
	resources *twopcpb.ConfirmTaskPeerInfo,
) *NeedExecuteTask {
	return &NeedExecuteTask{
		remotepid:              remotepid,
		proposalId:             proposalId,
		localTaskRole:          localTaskRole,
		localTaskOrganization:  localTaskOrganization,
		remoteTaskRole:         remoteTaskRole,
		remoteTaskOrganization: remoteTaskOrganization,
		task:                   task,
		consStatus:             consStatus,
		localResource:          localResource,
		resources:              resources,
	}
}
func (net *NeedExecuteTask) HasRemotePID() bool                     { return strings.Trim(string(net.remotepid), "") != "" }
func (net *NeedExecuteTask) HasEmptyRemotePID() bool                { return !net.HasRemotePID() }
func (net *NeedExecuteTask) GetRemotePID() peer.ID                  { return net.remotepid }
func (net *NeedExecuteTask) GetProposalId() common.Hash             { return net.proposalId }
func (net *NeedExecuteTask) GetLocalTaskRole() apicommonpb.TaskRole { return net.localTaskRole }
func (net *NeedExecuteTask) GetLocalTaskOrganization() *apicommonpb.TaskOrganization {
	return net.localTaskOrganization
}
func (net *NeedExecuteTask) GetRemoteTaskRole() apicommonpb.TaskRole { return net.remoteTaskRole }
func (net *NeedExecuteTask) GetRemoteTaskOrganization() *apicommonpb.TaskOrganization {
	return net.remoteTaskOrganization
}
func (net *NeedExecuteTask) GetTask() *Task                             { return net.task }
func (net *NeedExecuteTask) GetConsStatus() TaskActionStatus            { return net.consStatus }
func (net *NeedExecuteTask) GetLocalResource() *PrepareVoteResource     { return net.localResource }
func (net *NeedExecuteTask) GetResources() *twopcpb.ConfirmTaskPeerInfo { return net.resources }
func (net *NeedExecuteTask) String() string {
	taskStr := "{}"
	if nil != net.task {
		taskStr = net.task.GetTaskData().String()
	}
	localIdentityStr := "{}"
	if nil != net.localTaskOrganization {
		localIdentityStr = net.localTaskOrganization.String()
	}
	remoteIdentityStr := "{}"
	if nil != net.remoteTaskOrganization {
		remoteIdentityStr = net.remoteTaskOrganization.String()
	}
	localResourceStr := "{}"
	if nil != net.localResource {
		localResourceStr = net.localResource.String()
	}
	return fmt.Sprintf(`{"remotepid": %s, "proposalId": %s, "localTaskRole": %s, "localTaskOrganization": %s, "remoteTaskRole": %s, "remoteTaskOrganization": %s, "task": %s, "localResource": %s, "resources": %s}`,
		net.remotepid, net.proposalId.String(), net.localTaskRole.String(), localIdentityStr, net.remoteTaskRole.String(), remoteIdentityStr, taskStr, localResourceStr, ConfirmTaskPeerInfoString(net.resources))
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
