package evengine

import (
	"github.com/Metisnetwork/Metis-Carrier/common/timeutils"
	"github.com/Metisnetwork/Metis-Carrier/core/iface"
	libtypes "github.com/Metisnetwork/Metis-Carrier/lib/types"
)

type EventEngine struct {
	dataCenter iface.TaskCarrierDB
}

func NewEventEngine(dataCenter iface.TaskCarrierDB) *EventEngine {
	return &EventEngine{
		dataCenter: dataCenter,
	}
}

func (e *EventEngine) GenerateEvent(typ, taskId, identityId, partyId, extra string) *libtypes.TaskEvent {
	return &libtypes.TaskEvent{
		Type:       typ,
		TaskId:     taskId,
		IdentityId: identityId,
		PartyId:    partyId,
		Content:    extra,
		CreateAt:   timeutils.UnixMsecUint64(),
	}
}

func (e *EventEngine) StoreEvent(event *libtypes.TaskEvent) {
	if err := e.dataCenter.StoreTaskEvent(event); nil != err {
		log.WithError(err).Errorf("Failed to Store task event, taskId: {%s}, evnet: {%s}", event.TaskId, event.String())
	}
	return
}


