package task

import (
	"errors"
	"fmt"
	"github.com/RosettaFlow/carrier-go/event"
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

func MakeScheduleEventInfo(event *event.TaskEvent) (string, error) {
	if _, ok := ScheduleEvent[event.Type]; ok {
		return fmt.Sprintf("TaskID:%s, Identity:%s, Content:%s, CreateTime:%s, Information:%s",
			event.TaskId, event.Identity, event.Content, event.CreateTime, ScheduleEvent[event.Type]), nil
	}

	return "", IncEventType
}

// 数据服务事件

