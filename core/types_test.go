package core

import (
	"github.com/datumtechs/datum-network-carrier/common/traceutil"
	"github.com/datumtechs/datum-network-carrier/lib/netmsg/common"
	taskmngpb "github.com/datumtechs/datum-network-carrier/lib/netmsg/taskmng"
	"testing"
)

func TestTraceIDRules(t *testing.T) {
	pbmsg := &taskmngpb.TaskResultMsg{
		MsgOption:            &common.MsgOption{
		},
		TaskEventList:        make([]*taskmngpb.TaskEvent, 0),
		CreateAt:             10000,
		Sign:                 nil,
	}
	t.Log(traceutil.GenerateTraceID(pbmsg))


}
