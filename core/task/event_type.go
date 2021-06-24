package task

import (
	"errors"
	"github.com/RosettaFlow/Carrier-Go/event"
)

type EventSysCode string
func (code EventSysCode) String() string {return string(code)}
const (
	SysCode_Common EventSysCode = "00"
	SysCode_YarnNode EventSysCode = "01"
	SysCode_DataNode EventSysCode = "02"
	SysCode_PowerNode EventSysCode = "03"

)

type EventType struct {
	Type string
	Msg  string
}

func (e *EventType) EventInfo() string {
	return e.Msg
}

func NewEventType(Type string, text string) *EventType {
	return &EventType{Type: Type, Msg: text}
}

var IncEventType = errors.New("incorrect event type")

// 调度服务事件
var (
	OriginatingTask            = NewEventType("0100000", "Originating Task.")
	SuspendedTaskAgain         = NewEventType("0100001", "The task is suspended again.")
	DiscardedTask              = NewEventType("0100002", "The task was discarded.")
	FailTask                   = NewEventType("0100003", "The task was failed.")
	SucceedTask                = NewEventType("0100004", "The task was succeed.")
	UpdateComputeRes           = NewEventType("0100005", "Update computing resources.")
	StartTaskConsensus         = NewEventType("0101001", "Start of task consensus.")
	ConsensusTimeOutPhaseOne   = NewEventType("0101002", "Task Consensus 2PC Phase 1 timed out.")
	ConsensusNotMetPhaseOne    = NewEventType("0101003", "Mission Consensus 2PC Phase 1 votes are not met.")
	ConsensusCompletedPhaseOne = NewEventType("0101004", "Mission Consensus 2PC Phase 1 completed.")
	ConsensusPreLockedFail     = NewEventType("0101005", "Task consensus 2pc pre-locked resource failed.")
	ConsensusPreLockedSucceed  = NewEventType("0101006", "Task consensus 2pc pre-locked resources successfully.")
	ConsensusTimeOutPhaseTwo   = NewEventType("0101007", "Task Consensus 2PC Phase 2 timed out.")
	ConsensusNotMetPhaseTwo    = NewEventType("0101008", "Mission Consensus 2PC Phase 2 votes are not met.")
	ConsensusCompletedPhaseTwo = NewEventType("0101009", "Mission Consensus 2PC Phase 2 completed.")
	ConsensusPreLockedRelease  = NewEventType("0101010", "Task Consensus 2PC pre-locked resources are released.")
)

var ScheduleEvent = map[string]string{
	OriginatingTask.Type:            OriginatingTask.Msg,
	SuspendedTaskAgain.Type:         SuspendedTaskAgain.Msg,
	DiscardedTask.Type:              DiscardedTask.Msg,
	FailTask.Type:                   FailTask.Msg,
	SucceedTask.Type:                SucceedTask.Msg,
	UpdateComputeRes.Type:           UpdateComputeRes.Msg,
	StartTaskConsensus.Type:         StartTaskConsensus.Msg,
	ConsensusTimeOutPhaseOne.Type:   ConsensusTimeOutPhaseOne.Msg,
	ConsensusNotMetPhaseOne.Type:    ConsensusNotMetPhaseOne.Msg,
	ConsensusCompletedPhaseOne.Type: ConsensusCompletedPhaseOne.Msg,
	ConsensusPreLockedFail.Type:     ConsensusPreLockedFail.Msg,
	ConsensusPreLockedSucceed.Type:  ConsensusPreLockedSucceed.Msg,
	ConsensusTimeOutPhaseTwo.Type:   ConsensusTimeOutPhaseTwo.Msg,
	ConsensusNotMetPhaseTwo.Type:    ConsensusNotMetPhaseTwo.Msg,
	ConsensusCompletedPhaseTwo.Type: ConsensusCompletedPhaseTwo.Msg,
	ConsensusPreLockedRelease.Type:  ConsensusPreLockedRelease.Msg,
}

func MakeScheduleEventInfo(event *event.TaskEvent) (*event.TaskEvent, error) {
	if _, ok := ScheduleEvent[event.Type]; ok {
		event.Content = ScheduleEvent[event.Type]
		return event, nil
	}

	return event, IncEventType
}

// 数据服务事件
var (
	SourceUpLoadSucceed   = NewEventType("0207000", "Source data uploaded successfully.")
	SourceUpLoadFailed    = NewEventType("0207001", "Source data upload failed.")
	SourceDownloadSucceed = NewEventType("0207002", "Source data downloaded successfully.")
	SourceDownloadFailed  = NewEventType("0207003", "Source data downloaded failed.")
	SourceDeleteSucceed   = NewEventType("0207004", "Source data deleted successfully.")
	SourceDeleteFailed    = NewEventType("0207005", "Source data deleted failed.")
	StartDataShard        = NewEventType("0207006", "Start data sharding.")
	GetDataFileSucceed    = NewEventType("0207007", "Data file/directory retrieved successfully.")
	GetDataFileFailed     = NewEventType("0207008", "Data file/directory retrieved failed.")
)

var DataServiceEvent = map[string]string{
	SourceUpLoadSucceed.Type:   SourceUpLoadSucceed.Msg,
	SourceUpLoadFailed.Type:    SourceUpLoadFailed.Msg,
	SourceDownloadSucceed.Type: SourceDownloadSucceed.Msg,
	SourceDownloadFailed.Type:  SourceDownloadFailed.Msg,
	SourceDeleteSucceed.Type:   SourceDeleteSucceed.Msg,
	SourceDeleteFailed.Type:    SourceDeleteFailed.Msg,
	StartDataShard.Type:        StartDataShard.Msg,
	GetDataFileSucceed.Type:    GetDataFileSucceed.Msg,
	GetDataFileFailed.Type:     GetDataFileFailed.Msg,
}

func MakeDataServiceEventInfo(event *event.TaskEvent) (*event.TaskEvent, error) {
	if _, ok := DataServiceEvent[event.Type]; ok {
		event.Content = DataServiceEvent[event.Type]
		return event, nil
	}

	return event, IncEventType
}

// 计算服务事件
var (
	ReportComputeRes      = NewEventType("0308000", "Report computing resources.")
	StartNewTask          = NewEventType("0309000", "Start a new task.")
	DownloadCodeSucceed   = NewEventType("0309001", "The contract code was downloaded successfully.")
	DownloadCodeFailed    = NewEventType("0309002", "The contract code was downloaded failed.")
	StartBuildTaskEnv     = NewEventType("0309003", "Start building the computing task environment.")
	CreateIoSucceed       = NewEventType("0309004", "Create network IO successfully.")
	CreateIoFailed        = NewEventType("0309005", "Create network IO failed.")
	RegisterViaSucceed    = NewEventType("0309006", "Successfully registered with VIA service.")
	RegisterViaFailed     = NewEventType("0309007", "Failed to register with VIA service.")
	BuildTaskEnvSucceed   = NewEventType("0309008", "Build computational task environment successfully.")
	StartComputeTask      = NewEventType("0309009", "Starts the computation task.")
	CancelComputeTask     = NewEventType("0309010", "Cancel the execution of a computation task.")
	ExecuteComputeSucceed = NewEventType("0309011", "The computation task executed successfully.")
	ExecuteComputeFailed  = NewEventType("0309012", "The computation task failed to execute.")
	ReportComputeResult   = NewEventType("0309013", "Report of calculation results.")
	ReportTaskUsage       = NewEventType("0309014", "Resource usage for reporting tasks.")
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

func MakeComputerServiceEventInfo(event *event.TaskEvent) (*event.TaskEvent, error) {
	if _, ok := ComputerServiceEvent[event.Type]; ok {
		event.Content = ComputerServiceEvent[event.Type]
		return event, nil
	}

	return nil, IncEventType
}
