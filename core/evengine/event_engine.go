package evengine

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/core/iface"
	"github.com/RosettaFlow/Carrier-Go/types"
	"time"
)


type EventEngine struct {
	dataCenter   iface.TaskCarrierDB
}


func NewEventEngine(dataCenter iface.TaskCarrierDB) *EventEngine {
	return &EventEngine{
		dataCenter: dataCenter,
	}
}

func (e *EventEngine) GenerateEvent(typ, taskId, identityId, extra string) *types.TaskEventInfo {
	return &types.TaskEventInfo{
		Type: typ,
		TaskId: taskId,
		Identity: identityId,
		Content: fmt.Sprintf("%s, reason: %s", ScheduleEvent[typ], extra),
		CreateTime: uint64(time.Now().UnixNano()),
	}
}
func  (e *EventEngine) StoreEvent(event *types.TaskEventInfo) error {
	return e.dataCenter.StoreTaskEvent(event)
}
func  (e *EventEngine) GetTaskEventList(taskId string) ([]*types.TaskEventInfo, error) {
	return e.dataCenter.GetTaskEventList(taskId)
}
func  (e *EventEngine)  CleanTaskEventList(taskId string) error {
	return e.dataCenter.CleanTaskEventList(taskId)
}