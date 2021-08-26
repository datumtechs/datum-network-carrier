package types

import (
	"encoding/json"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/bytesutil"
	"github.com/RosettaFlow/Carrier-Go/common/rlputil"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"time"

	"github.com/RosettaFlow/Carrier-Go/consensus/twopc/utils"
	"github.com/RosettaFlow/Carrier-Go/types"
	"sync/atomic"
)

type taskOption struct {
	Role                  types.TaskRole           `json:"role"` // The role information of the current recipient of the task
	TaskId                string                   `json:"taskId"`
	TaskName              string                   `json:"taskName"`
	Owner                 *apipb.Organization         `json:"owner"`
	AlgoSupplier          *apipb.Organization         `json:"algoSupplier"`
	DataSupplier          []*dataSupplierOption    `json:"dataSupplier"`
	PowerSupplier         []*powerSupplierOption   `json:"powerSupplier"`
	Receivers             []*receiverOption        `json:"receivers"`
	OperationCost         *TaskOperationCost `json:"operationCost"`
	CalculateContractCode string                   `json:"calculateContractCode"`
	DataSplitContractCode string                   `json:"dataSplitContractCode"`
	CreateAt              uint64                   `json:"createat"`
}

func (t *taskOption) Hash() common.Hash {
	return rlputil.RlpHash(t)
}

type TaskOperationCost struct {
	Processor uint64 `json:"processor"`
	Mem       uint64 `json:"mem"`
	Bandwidth uint64 `json:"bandwidth"`
	Duration  uint64 `json:"duration"`
}

func (cost *TaskOperationCost) String() string {
	return fmt.Sprintf(`{"mem": %d, "processor": %d, "bandwidth": %d, "duration": %d}`, cost.Mem, cost.Processor, cost.Bandwidth, cost.Duration)
}

type dataSupplierOption struct {
	MemberInfo      *apipb.Organization `json:"memberInfo"`
	MetaDataId      string           `json:"metaDataId"`
	ColumnIndexList []uint64         `json:"columnIndexList"`
}
type powerSupplierOption struct {
	MemberInfo *apipb.Organization `json:"memberInfo"`
}
type receiverOption struct {
	MemberInfo *apipb.Organization   `json:"memberInfo"`
	Providers  []*apipb.Organization `json:"providers"`
}

type taskPeerInfo struct {
	// Used to connect when running task, internal network resuorce of org.
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

type PrepareMsg struct {
	ProposalID  common.Hash   `json:"proposalId"`
	TaskOption  *taskOption   `json:"taskOption"`
	CreateAt    uint64        `json:"createAt"`
	Sign        types.MsgSign `json:"sign"`
	messageHash atomic.Value  `rlp:"-"`
}

func (msg *PrepareMsg) String() string {
	b, _ := json.Marshal(msg)
	return string(b)
}
func (msg *PrepareMsg) MsgHash() common.Hash {
	if mhash := msg.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(PrepareProposalMsg, utils.MergeBytes(msg.ProposalID.Bytes(),
		msg.TaskOption.Hash().Bytes(), msg.Sign.Bytes(), bytesutil.Uint64ToBytes(msg.CreateAt)))
	msg.messageHash.Store(v)
	return v
}

type PrepareVote struct {
	ProposalID  common.Hash      `json:"proposalId"`
	Role        types.TaskRole   `json:"role"` // The role information of the current recipient of the task
	Owner       *apipb.Organization `json:"owner"`
	VoteOption  types.VoteOption `json:"voteOption"`
	PeerInfo    *taskPeerInfo    `json:"peerInfo"`
	CreateAt    uint64           `json:"createAt"`
	Sign        types.MsgSign    `json:"sign"`
	messageHash atomic.Value     `rlp:"-"`
}

func (msg *PrepareVote) String() string {
	b, _ := json.Marshal(msg)
	return string(b)
}

func (msg *PrepareVote) MsgHash() common.Hash {
	if mhash := msg.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(PrepareVoteMsg, utils.MergeBytes(msg.ProposalID.Bytes(), /*msg.VoteNodeID.Bytes(), */ // TODO 编码 NodeAlias
		msg.VoteOption.Bytes(), msg.Sign.Bytes(), bytesutil.Uint64ToBytes(msg.CreateAt)))
	msg.messageHash.Store(v)
	return v
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

//type ConfirmEpoch uint64
//
//func (c ConfirmEpoch) Uint64() uint64 { return uint64(c) }

type ProposalState struct {
	ProposalId         common.Hash
	TaskDir            types.ProposalTaskDir
	TaskRole           types.TaskRole
	SelfIdentity       *apipb.TaskOrganization
	TaskId             string
	PeriodNum          ProposalStatePeriod
	PrePeriodStartTime uint64
	PeriodStartTime    uint64 // the timestemp
	// Clear `ProposalState` ,
	// when the current time is greater than the `DeadlineDuration` createAt of proposalState
	DeadlineDuration uint64
	CreateAt         uint64
}

var EmptyProposalState = new(ProposalState)

func NewProposalState(proposalId common.Hash, taskId string,
	TaskDir types.ProposalTaskDir, taskRole types.TaskRole, selfIdentity *apipb.TaskOrganization, startTime uint64) *ProposalState {

	return &ProposalState{
		ProposalId:       proposalId,
		TaskId:           taskId,
		TaskDir:          TaskDir,
		TaskRole:         taskRole,
		SelfIdentity:     selfIdentity,
		PeriodNum:        PeriodPrepare,
		PeriodStartTime:  startTime,
		DeadlineDuration: ProposalDeadlineDuration,
		CreateAt:         uint64(timeutils.UnixMsec()),
	}
}

func (pstate *ProposalState) GetProposalId() common.Hash         { return pstate.ProposalId }
func (pstate *ProposalState) CurrPeriodNum() ProposalStatePeriod { return pstate.PeriodNum }

//func (pstate *ProposalState) CurrPeriodDuration() uint64 {
//	return pstate.PeriodStartTime - pstate.PeriodEndTime
//}
func (pstate *ProposalState) GetPeriod() string {

	switch pstate.PeriodNum  {
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
func (pstate *ProposalState) IsPreparePeriod() bool          { return pstate.PeriodNum == PeriodPrepare }
func (pstate *ProposalState) IsConfirmPeriod() bool          { return pstate.PeriodNum == PeriodConfirm }
func (pstate *ProposalState) IsCommitPeriod() bool           { return pstate.PeriodNum == PeriodCommit }
func (pstate *ProposalState) IsFinishedPeriod() bool         { return pstate.PeriodNum == PeriodFinished }
func (pstate *ProposalState) IsNotPreparePeriod() bool       { return !pstate.IsPreparePeriod() }
func (pstate *ProposalState) IsNotConfirmPeriod() bool       { return !pstate.IsConfirmPeriod() }
func (pstate *ProposalState) IsNotCommitPeriod() bool        { return !pstate.IsCommitPeriod() }
func (pstate *ProposalState) IsNotFinishedPeriod() bool      { return !pstate.IsFinishedPeriod() }
func (pstate *ProposalState) IsDeadline() bool {
	now := uint64(timeutils.UnixMsec())
	return (now - pstate.CreateAt) >= ProposalDeadlineDuration
}

//func (pstate *ProposalState) IsFirstConfirmEpoch() bool {
//	if pstate.ConfirmEpoch == ConfirmEpochFirst {
//		return true
//	}
//	return false
//}
//func (pstate *ProposalState) IsSecondConfirmEpoch() bool {
//	if pstate.ConfirmEpoch == ConfirmEpochSecond {
//		return true
//	}
//	return false
//}
func (pstate *ProposalState) IsEmpty() bool {
	if pstate == EmptyProposalState {
		return true
	}
	return false
}

//func (pstate *ProposalState) GetConfirmEpoch() ConfirmEpoch { return pstate.ConfirmEpoch }
func (pstate *ProposalState) IsPrepareTimeout() bool {

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
func (pstate *ProposalState) IsConfirmTimeout() bool {

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

func (pstate *ProposalState) IsCommitTimeout() bool {
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

func (pstate *ProposalState) ChangeToConfirm(startTime uint64) {
	if pstate.PeriodNum == PeriodPrepare {
		pstate.PrePeriodStartTime = pstate.PeriodStartTime
		pstate.PeriodStartTime = startTime
		pstate.PeriodNum = PeriodConfirm
	}
	//pstate.ConfirmEpoch = ConfirmEpochFirst
}

func (pstate *ProposalState) ChangeToCommit(startTime uint64) {
	if pstate.PeriodNum == PeriodConfirm {
		pstate.PrePeriodStartTime = pstate.PeriodStartTime
		pstate.PeriodStartTime = startTime
		pstate.PeriodNum = PeriodCommit
	}
	//pstate.ConfirmEpoch = ConfirmEpochUnknown
}
func (pstate *ProposalState) ChangeToFinished(startTime uint64) {
	if pstate.PeriodNum == PeriodCommit {
		pstate.PrePeriodStartTime = pstate.PeriodStartTime
		pstate.PeriodStartTime = startTime
		pstate.PeriodNum = PeriodFinished
	}
	//pstate.ConfirmEpoch = ConfirmEpochUnknown
}
