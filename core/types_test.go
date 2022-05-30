package core

import (
	"github.com/datumtechs/datum-network-carrier/common/traceutil"
	carriernetmsgcommonpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/common"
	carriernetmsgtaskmngpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/taskmng"
	"testing"
)

func TestTraceIDRules(t *testing.T) {
	pbmsg := &carriernetmsgtaskmngpb.TaskResultMsg{
		MsgOption:            &carriernetmsgcommonpb.MsgOption{
		},
		TaskEventList:        make([]*carriernetmsgtaskmngpb.TaskEvent, 0),
		CreateAt:             10000,
		Sign:                 nil,
	}
	t.Log(traceutil.GenerateTraceID(pbmsg))


}
