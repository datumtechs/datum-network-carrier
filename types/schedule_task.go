package types

import (
	"encoding/json"
	"github.com/RosettaFlow/Carrier-Go/common"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
)

type ProposalTask struct {
	ProposalId common.Hash
	*ScheduleTask
	CreateAt uint64
}

type ConsensusTaskWrap struct {
	Task              *ScheduleTask
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
	Task     *ScheduleTask
	ResultCh chan *ScheduleResult
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
	SelfPeerInfo *pb.TaskPeerInfo
	Task         *ConsensusScheduleTask
	ResultCh     chan *TaskResultMsgWrap
}
type ConsensusScheduleTask struct {
	TaskDir   ProposalTaskDir
	TaskState TaskState
	SchedTask *ScheduleTask
	Resources *pb.ConfirmTaskPeerInfo
}

type ScheduleTask struct {
	TaskId                string                        `json:"TaskId"`
	TaskName              string                        `json:"taskName"`
	Owner                 *ScheduleTaskDataSupplier     `json:"owner"`
	Partners              []*ScheduleTaskDataSupplier   `json:"partners"`
	PowerSuppliers        []*ScheduleTaskPowerSupplier  `json:"powerSuppliers"`
	Receivers             []*ScheduleTaskResultReceiver `json:"receivers"`
	CalculateContractCode string                        `json:"calculateContractCode"`
	DataSplitContractCode string                        `json:"dataSplitContractCode"`
	OperationCost         *TaskOperationCost            `json:"spend"`
	CreateAt              uint64                        `json:"createAt"`
	StartAt               uint64                        `json:"startAt"`
}

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

type ConsensuResult struct {
	*TaskConsResult
	//Resources *pb.ConfirmTaskPeerInfo
	//OwnerResource          *PrepareVoteResource
	//PartnersResource       []*PrepareVoteResource
	//PowerSuppliersResource []*PrepareVoteResource
	//ReceiversResource      []*PrepareVoteResource
}

//type ScheduleNodeAlias struct {
//	Name         string `json:"name"`
//	NodeId       string `json:"nodeId"`
//	IdentityId   string `json:"identityId"`
//	PeerIp   string `json:"peerIp"`
//	PeerPort string `json:"peerPort"`
//}
