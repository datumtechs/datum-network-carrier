package evengine

import (
	"github.com/datumtechs/datum-network-carrier/carrierdb/iface"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
)

type EventEngine struct {
	dataCenter iface.TaskCarrierDB
}

func NewEventEngine(dataCenter iface.TaskCarrierDB) *EventEngine {
	return &EventEngine{
		dataCenter: dataCenter,
	}
}

func (e *EventEngine) GenerateEvent(typ, taskId, identityId, partyId, extra string) *carriertypespb.TaskEvent {
	return &carriertypespb.TaskEvent{
		Type:       typ,
		TaskId:     taskId,
		IdentityId: identityId,
		PartyId:    partyId,
		Content:    extra,
		CreateAt:   timeutils.UnixMsecUint64(),
	}
}

func (e *EventEngine) StoreEvent(event *carriertypespb.TaskEvent) {
	if err := e.dataCenter.StoreTaskEvent(event); nil != err {
		log.WithError(err).Errorf("Failed to Store task event, taskId: {%s}, evnet: {%s}", event.TaskId, event.String())
	}
	return
}


