package types

import (
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)

type ReportTaskEvent struct {
	 PartyId   string
	 Event     *libtypes.TaskEvent
}

func NewReportTaskEvent(partyId string, event *libtypes.TaskEvent) *ReportTaskEvent {
	return &ReportTaskEvent{
		PartyId: partyId,
		Event: event,
	}
}
func (rte *ReportTaskEvent) GetPartyId() string            { return rte.PartyId }
func (rte *ReportTaskEvent) GetEvent() *libtypes.TaskEvent { return rte.Event }