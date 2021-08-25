package evengine

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/core/iface"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)


type EventEngine struct {
	dataCenter   iface.TaskCarrierDB
}


func NewEventEngine(dataCenter iface.TaskCarrierDB) *EventEngine {
	return &EventEngine{
		dataCenter: dataCenter,
	}
}

func (e *EventEngine) GenerateEvent(typ, taskId, identityId, extra string) *libTypes.TaskEvent {
	return &libTypes.TaskEvent{
		Type: typ,
		TaskId: taskId,
		IdentityId: identityId,
		Content: fmt.Sprintf("%s, reason: {%s}", ScheduleEvent[typ], extra),
		CreateAt: uint64(timeutils.UnixMsec()),
	}
}
func  (e *EventEngine) StoreEvent(event *libTypes.TaskEvent) {
	if err := e.dataCenter.StoreTaskEvent(event); nil != err {
		log.Errorf("Failed to Store task event, taskId: {%s}, evnet: {%s}, err: {%s}", event.TaskId, event.String(), err)
	}
	return
}
func  (e *EventEngine) GetTaskEventList(taskId string) ([]*libTypes.TaskEvent, error) {
	return e.dataCenter.GetTaskEventList(taskId)
}
func  (e *EventEngine)  RemoveTaskEventList(taskId string) error {
	return e.dataCenter.RemoveTaskEventList(taskId)
}