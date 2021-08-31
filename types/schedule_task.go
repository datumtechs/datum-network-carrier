package types

import (
	"encoding/json"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"strings"
)

type ProposalTask struct {
	ProposalId common.Hash
	*Task
	CreateAt uint64
}

type ConsensusTaskWrap struct {
	Task              *Task
	OwnerDataResource *PrepareVoteResource
	ResultCh          chan *ConsensusResult
}

func (wrap *ConsensusTaskWrap) SendResult(result *ConsensusResult) {
	wrap.ResultCh <- result
	close(wrap.ResultCh)
}
func (wrap *ConsensusTaskWrap) RecvResult() *ConsensusResult {
	return <-wrap.ResultCh
}
func (wrap *ConsensusTaskWrap) String() string {
	result, err := json.Marshal(wrap)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}

type ReplayScheduleTaskWrap struct {
	Role     TaskRole
	PartyId  string
	Task     *Task
	ResultCh chan *ScheduleResult
}

func NewReplayScheduleTaskWrap(role TaskRole, partyId string, task *Task) *ReplayScheduleTaskWrap {
	return &ReplayScheduleTaskWrap{
		Role:     role,
		PartyId:  partyId,
		Task:     task,
		ResultCh: make(chan *ScheduleResult),
	}
}
func (wrap *ReplayScheduleTaskWrap) SendFailedResult(taskId string, err error) {
	wrap.SendResult(&ScheduleResult{
		TaskId: taskId,
		Status: TaskSchedFailed,
		Err:    err,
	})
}
func (wrap *ReplayScheduleTaskWrap) SendResult(result *ScheduleResult) {
	wrap.ResultCh <- result
	close(wrap.ResultCh)
}
func (wrap *ReplayScheduleTaskWrap) RecvResult() *ScheduleResult {
	return <-wrap.ResultCh
}
func (wrap *ReplayScheduleTaskWrap) String() string {
	result, err := json.Marshal(wrap)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}

type DoneScheduleTaskChWrap struct {
	ProposalId   common.Hash
	SelfTaskRole TaskRole
	SelfIdentity *apipb.TaskOrganization
	Task         *ConsensusScheduleTask
	ResultCh     chan *TaskResultMsgWrap
}
type ConsensusScheduleTask struct {
	TaskDir          ProposalTaskDir
	TaskState        TaskState
	SchedTask        *Task
	SelfVotePeerInfo *PrepareVoteResource
	Resources        *pb.ConfirmTaskPeerInfo
}

type TaskConsStatus uint16

func (t TaskConsStatus) String() string {
	switch t {
	case TaskSucceed:
		return "TaskSucceed"
	case TaskConsensusInterrupt:
		return "TaskConsensusInterrupt"
	case TaskRunningInterrupt:
		return "TaskRunningInterrupt"
	default:
		return "UnknownTaskResultStatus"
	}
}

const (
	TaskSucceed            TaskConsStatus = 0x0000
	TaskConsensusInterrupt TaskConsStatus = 0x0001
	TaskRunningInterrupt   TaskConsStatus = 0x0100
)

type TaskConsResult struct {
	TaskId string
	Status TaskConsStatus
	Done   bool
	Err    error
}

type TaskSchedStatus bool

func (status TaskSchedStatus) String() string {
	switch status {
	case TaskSchedOk:
		return "TaskSchedOk"
	case TaskSchedFailed:
		return "TaskSchedFailed"
	default:
		return "UnknownTaskSchedResult"
	}
}

const (
	TaskSchedOk     TaskSchedStatus = true
	TaskSchedFailed TaskSchedStatus = false
)

type ScheduleResult struct {
	TaskId   string
	Status   TaskSchedStatus
	Err      error
	Resource *PrepareVoteResource
}

func (res *ScheduleResult) String() string {
	return fmt.Sprintf(`{"taskId": %s, "status": %s, "err": %s, "resource": %s}`,
		res.TaskId, res.Status.String(), res.Err, res.Resource.String())
}

type ConsensusResult struct {
	*TaskConsResult
}

func (res *ConsensusResult) String() string {
	return fmt.Sprintf(`{"taskId": %s, "status": %s, "done": %v, "err": %s}`, res.TaskId, res.Status.String(), res.Done, res.Err)
}

// ================================================= V2.0 =================================================

// 需要被 进行共识的 local task (已经调度好的, 还未共识的)
type NeedConsensusTask struct {
	task     *Task
	need     bool
	resource *PrepareVoteResource
	resultCh chan *ConsensusResult
}

func NewNeedConsensusTask(task *Task, need bool, resource *PrepareVoteResource) *NeedConsensusTask {
	return &NeedConsensusTask{
		task:     task,
		need:     need,
		resource: resource,
		resultCh: make(chan *ConsensusResult),
	}
}
func (nct *NeedConsensusTask) Task() *Task                     { return nct.task }
func (nct *NeedConsensusTask) Need() bool                      { return nct.need }
func (nct *NeedConsensusTask) Resource() *PrepareVoteResource  { return nct.resource }
func (nct *NeedConsensusTask) ResultCh() chan *ConsensusResult { return nct.resultCh }
func (nct *NeedConsensusTask) String() string {
	taskStr := "{}"
	if nil != nct.task {
		taskStr = nct.task.TaskData().String()
	}
	resourceStr := "{}"
	if nil != nct.resource {
		resourceStr = nct.resource.String()
	}
	return fmt.Sprintf(`{"task": %s, "need": %v, "selfResource": %v, "resultCh": %p}`,
		taskStr, nct.need, resourceStr, nct.resultCh)
}

// 需要 重演调度的 remote task (接收到对端发来的 proposal 中的, 处于共识过程中的, 需要重演调度的)
type NeedReplayScheduleTask struct {
	selfTaskRole TaskRole
	selfPartyId  string
	task         *Task
	resultCh     chan *ReplayScheduleResult
}

func NewNeedReplayScheduleTask(role TaskRole, partyId string, task *Task) *NeedReplayScheduleTask {
	return &NeedReplayScheduleTask{
		selfTaskRole: role,
		selfPartyId:  partyId,
		task:         task,
		resultCh:     make(chan *ReplayScheduleResult),
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
func (nrst *NeedReplayScheduleTask) RecvResult() *ReplayScheduleResult {
	return <-nrst.resultCh
}
func (nrst *NeedReplayScheduleTask) SelfTaskRole() TaskRole               { return nrst.selfTaskRole }
func (nrst *NeedReplayScheduleTask) SelfPartyId() string                  { return nrst.selfPartyId }
func (nrst *NeedReplayScheduleTask) Task() *Task                          { return nrst.task }
func (nrst *NeedReplayScheduleTask) ResultCh() chan *ReplayScheduleResult { return nrst.resultCh }
func (nrst *NeedReplayScheduleTask) String() string {
	taskStr := "{}"
	if nil != nrst.task {
		taskStr = nrst.task.TaskData().String()
	}
	return fmt.Sprintf(`{"selfTaskRole": %s, "selfPartyId": %s, "task": %s, "resultCh": %p}`,
		nrst.selfTaskRole.String(), nrst.selfPartyId, taskStr, nrst.resultCh)
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
func (rsr *ReplayScheduleResult) TaskId() string                 { return rsr.taskId }
func (rsr *ReplayScheduleResult) Err() error                     { return rsr.err }
func (rsr *ReplayScheduleResult) Resource() *PrepareVoteResource { return rsr.resource }
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
	proposalId   common.Hash
	selfTaskRole TaskRole
	selfIdentity *apipb.TaskOrganization
	task         *Task
	selfResource *PrepareVoteResource
	resources    *pb.ConfirmTaskPeerInfo
}

func NewNeedExecuteTask(
	proposalId common.Hash,
	selfTaskRole TaskRole,
	selfIdentity *apipb.TaskOrganization,
	task *Task,
	selfResource *PrepareVoteResource,
	resources *pb.ConfirmTaskPeerInfo,
) *NeedExecuteTask {
	return &NeedExecuteTask{
		proposalId:   proposalId,
		selfTaskRole: selfTaskRole,
		selfIdentity: selfIdentity,
		task:         task,
		selfResource: selfResource,
		resources:    resources,
	}
}
func (net *NeedExecuteTask) ProposalId() common.Hash                   { return net.proposalId }
func (net *NeedExecuteTask) TaskRole() TaskRole                        { return net.selfTaskRole }
func (net *NeedExecuteTask) TaskOrganization() *apipb.TaskOrganization { return net.selfIdentity }
func (net *NeedExecuteTask) Task() *Task                               { return net.task }
func (net *NeedExecuteTask) SelfResource() *PrepareVoteResource        { return net.selfResource }
func (net *NeedExecuteTask) Resources() *pb.ConfirmTaskPeerInfo        { return net.resources }
func (net *NeedExecuteTask) String() string {
	taskStr := "{}"
	if nil != net.task {
		taskStr = net.task.TaskData().String()
	}
	identityStr := "{}"
	if nil != net.selfIdentity {
		identityStr = net.selfIdentity.String()
	}
	selfResourceStr := "{}"
	if nil != net.selfResource {
		selfResourceStr = net.selfResource.String()
	}
	return fmt.Sprintf(`{"proposalId": %s, "selfTaskRole": %s, "selfIdentity": %s, "task": %s, "selfResource": %s, "resources": %s}`,
		net.proposalId.String(), net.selfTaskRole.String(), identityStr, taskStr, selfResourceStr, ConfirmTaskPeerInfoString(net.resources))
}

func ConfirmTaskPeerInfoString(resources *pb.ConfirmTaskPeerInfo) string {
	if nil == resources {
		return "{}"
	}
	ownerPeerInfoStr := FetchTaskPeerInfo(resources.GetOwnerPeerInfo()).String()
	dataSupplierList := make([]string, len(resources.GetDataSupplierPeerInfoList()))
	for i, peerInfo := range resources.GetDataSupplierPeerInfoList() {
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

	powerSupplierList := make([]string, len(resources.GetPowerSupplierPeerInfoList()))
	for i, peerInfo := range resources.GetPowerSupplierPeerInfoList() {
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

	receiverList := make([]string, len(resources.GetResultReceiverPeerInfoList()))
	for i, peerInfo := range resources.GetResultReceiverPeerInfoList() {
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

	return fmt.Sprintf(`{"ownerPeerInfo": %s, "dataSupplierPeerInfoList": %s, "powerSupplierPeerInfoList": %s, "resultReceiverPeerInfoList": %s}`,
		ownerPeerInfoStr, dataSupplierListStr, powerSupplierListStr, receiverListStr)
}
