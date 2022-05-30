package evengine

import (
	"errors"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
)

type EventSysCode string

func (code EventSysCode) String() string { return string(code) }

const (
	EventTypeCharLen = 7
)

type EventType struct {
	Type string
	Msg  string
}

func NewEventType(Type string, text string) *EventType {
	return &EventType{Type: Type, Msg: text}
}

func (e *EventType) GetType() string { return e.Type }
func (e *EventType) GetMsg() string  { return e.Msg }

var IncEventType = errors.New("incorrect evengine type")

// 系统码
var (
	// 任务最终成功
	TaskExecuteSucceedEOF = NewEventType("0008000", "task was executed succeed EOF")
	TaskExecuteFailedEOF  = NewEventType("0008001", "task was executed failed EOF")
)

// 调度服务事件
var (
	TaskCreate                 = NewEventType("0100000", "created task")
	TaskNeedRescheduled        = NewEventType("0100001", "the task need rescheduled")
	TaskDiscarded              = NewEventType("0100002", "the task was discarded")
	TaskStarted                = NewEventType("0100003", "the task was started")
	TaskFailed                 = NewEventType("0100004", "the task was failed")
	TaskSucceed                = NewEventType("0100005", "the task was succeed")
	TaskResourceElectionFailed = NewEventType("0100006", "the resource of task was failed on election")
	TaskScheduleFailed         = NewEventType("0100007", "the task was scheduled failed")
	TaskTerminated             = NewEventType("0100008", "the task was terminated")
	TaskStartConsensus         = NewEventType("0101001", "the task was started consensus")
	TaskFailedConsensus        = NewEventType("0101002", "the task was failed consensus")
	TaskProposalStateDeadline  = NewEventType("0101003", "the task proposalState was deadline")
	TaskSucceedConsensus       = NewEventType("0101004", "the task was succeed consensus")
	TaskConsensusPrepareEpoch  = NewEventType("0101005", "the task was consensus prepare epoch")
	TaskConsensusConfirmEpoch  = NewEventType("0101006", "the task was consensus confirm epoch")
	TaskConsensusCommitEpoch   = NewEventType("0101007", "the task was consensus commit epoch")
)

var ScheduleEvent = map[string]string{
	TaskCreate.Type:          TaskCreate.Msg,
	TaskNeedRescheduled.Type: TaskNeedRescheduled.Msg,
	TaskDiscarded.Type:       TaskDiscarded.Msg,
	TaskFailed.Type:          TaskFailed.Msg,
	TaskSucceed.Type:         TaskSucceed.Msg,
	TaskStartConsensus.Type:  TaskStartConsensus.Msg,
	TaskFailedConsensus.Type: TaskFailedConsensus.Msg,
}

func MakeScheduleEventInfo(event *carriertypespb.TaskEvent) (*carriertypespb.TaskEvent, error) {
	if _, ok := ScheduleEvent[event.Type]; ok {
		event.Content = ScheduleEvent[event.Type]
		return event, nil
	}

	return event, IncEventType
}

// 数据服务事件
var (
	SourceUpLoadSucceed      = NewEventType("0207000", "source data uploaded successfully.")
	SourceUpLoadFailed       = NewEventType("0207001", "source data upload failed.")
	SourceDownloadSucceed    = NewEventType("0207002", "source data downloaded successfully.")
	SourceDownloadFailed     = NewEventType("0207003", "source data downloaded failed.")
	SourceDeleteSucceed      = NewEventType("0207004", "source data deleted successfully.")
	SourceDeleteFailed       = NewEventType("0207005", "source data deleted failed.")
	StartDataShard           = NewEventType("0207006", "start data sharding.")
	RetrievedDataFileSucceed = NewEventType("0207007", "data file/directory retrieved successfully.")
	RetrievedDataFileFailed  = NewEventType("0207008", "data file/directory retrieved failed.")
)

var DataServiceEvent = map[string]string{
	SourceUpLoadSucceed.Type:      SourceUpLoadSucceed.Msg,
	SourceUpLoadFailed.Type:       SourceUpLoadFailed.Msg,
	SourceDownloadSucceed.Type:    SourceDownloadSucceed.Msg,
	SourceDownloadFailed.Type:     SourceDownloadFailed.Msg,
	SourceDeleteSucceed.Type:      SourceDeleteSucceed.Msg,
	SourceDeleteFailed.Type:       SourceDeleteFailed.Msg,
	StartDataShard.Type:           StartDataShard.Msg,
	RetrievedDataFileSucceed.Type: RetrievedDataFileSucceed.Msg,
	RetrievedDataFileFailed.Type:  RetrievedDataFileFailed.Msg,
}

// 计算服务事件
var (
	ReportComputeRes      = NewEventType("0308000", "report computing resources.")
	StartNewTask          = NewEventType("0309000", "start a new task.")
	DownloadCodeSucceed   = NewEventType("0309001", "the contract code was downloaded successfully.")
	DownloadCodeFailed    = NewEventType("0309002", "the contract code was downloaded failed.")
	StartBuildTaskEnv     = NewEventType("0309003", "start building the computing task environment.")
	CreateIoSucceed       = NewEventType("0309004", "create network IO successfully.")
	CreateIoFailed        = NewEventType("0309005", "create network IO failed.")
	RegisterViaSucceed    = NewEventType("0309006", "successfully registered with VIA service.")
	RegisterViaFailed     = NewEventType("0309007", "failed to register with VIA service.")
	BuildTaskEnvSucceed   = NewEventType("0309008", "build computational task environment successfully.")
	StartComputeTask      = NewEventType("0309009", "starts the computation task.")
	CancelComputeTask     = NewEventType("0309010", "cancel the execution of a computation task.")
	ExecuteComputeSucceed = NewEventType("0309011", "the computation task executed successfully.")
	ExecuteComputeFailed  = NewEventType("0309012", "the computation task failed to execute.")
	ReportComputeResult   = NewEventType("0309013", "report of calculation results.")
	ReportTaskUsage       = NewEventType("0309014", "resource usage for reporting tasks.")
)

var ComputerServiceEvent = map[string]string{
	ReportComputeRes.Type:      ReportComputeRes.Msg,
	StartNewTask.Type:          StartNewTask.Msg,
	DownloadCodeSucceed.Type:   DownloadCodeSucceed.Msg,
	DownloadCodeFailed.Type:    DownloadCodeFailed.Msg,
	StartBuildTaskEnv.Type:     StartBuildTaskEnv.Msg,
	CreateIoSucceed.Type:       CreateIoSucceed.Msg,
	CreateIoFailed.Type:        CreateIoFailed.Msg,
	RegisterViaSucceed.Type:    RegisterViaSucceed.Msg,
	RegisterViaFailed.Type:     RegisterViaFailed.Msg,
	BuildTaskEnvSucceed.Type:   BuildTaskEnvSucceed.Msg,
	StartComputeTask.Type:      StartComputeTask.Msg,
	CancelComputeTask.Type:     CancelComputeTask.Msg,
	ExecuteComputeSucceed.Type: ExecuteComputeSucceed.Msg,
	ExecuteComputeFailed.Type:  ExecuteComputeFailed.Msg,
	ReportComputeResult.Type:   ReportComputeResult.Msg,
	ReportTaskUsage.Type:       ReportTaskUsage.Msg,
}
