package types

import (
	"encoding/json"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)

type ProposalTask struct {
	ProposalId common.Hash
	*Task
	CreateAt uint64
}

type ConsensusTaskWrap struct {
	Task              *Task
	OwnerDataResource *PrepareVoteResource
	ResultCh          chan *ConsensuResult
}



func (wrap *ConsensusTaskWrap) SendResult(result *ConsensuResult) {
	wrap.ResultCh <- result
	close(wrap.ResultCh)
}
func (wrap *ConsensusTaskWrap) RecvResult() *ConsensuResult {
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

func NewReplayScheduleTaskWrap (role TaskRole, partyId  string, task *Task) *ReplayScheduleTaskWrap{
	return &ReplayScheduleTaskWrap{
		Role: role,
		PartyId: partyId,
		Task: task,
		ResultCh: make(chan *ScheduleResult),
	}
}
func (wrap *ReplayScheduleTaskWrap) SendFailedResult(taskId string, err error) {
	wrap.SendResult(&ScheduleResult{
		TaskId:taskId,
		Status: TaskSchedFailed,
		Err: err,
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
	SelfIdentity *libTypes.OrganizationData
	Task         *ConsensusScheduleTask
	ResultCh     chan *TaskResultMsgWrap
}
type ConsensusScheduleTask struct {
	TaskDir   ProposalTaskDir
	TaskState TaskState
	SchedTask *Task
	Resources *pb.ConfirmTaskPeerInfo
}

//type ScheduleTask struct {
//	TaskId                string                        `json:"TaskId"`
//	TaskName              string                        `json:"taskName"`
//	Owner                 *ScheduleTaskDataSupplier     `json:"owner"`
//	Partners              []*ScheduleTaskDataSupplier   `json:"partners"`
//	PowerSuppliers        []*ScheduleTaskPowerSupplier  `json:"powerSuppliers"`
//	Receivers             []*ScheduleTaskResultReceiver `json:"receivers"`
//	CalculateContractCode string                        `json:"calculateContractCode"`
//	DataSplitContractCode string                        `json:"dataSplitContractCode"`
//	OperationCost         *TaskOperationCost            `json:"spend"`
//	CreateAt              uint64                        `json:"createAt"`
//	StartAt               uint64                        `json:"startAt"`
//}


type ScheduleTaskDataSupplier struct {
	*TaskNodeAlias
	MetaData *SupplierMetaData `json:"metaData"`
}
type ScheduleTaskPowerSupplier struct {
	*TaskNodeAlias
}

type ScheduleSupplierMetaData struct {
	MetaId          string   `json:"metaId"`
	ColumnIndexList []uint64 `json:"columnIndexList"`
}

type ScheduleTaskResultReceiver struct {
	*TaskNodeAlias
	Providers []*TaskNodeAlias `json:"providers"`
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

// Task consensus result
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
type ConsensuResult struct {
	*TaskConsResult
}

func (res *ConsensuResult) String() string {
	return fmt.Sprintf(`{"taskId": %s, "status": %s, "done": %v, "err": %s}`, res.TaskId, res.Status.String(), res.Done, res.Err)
}
