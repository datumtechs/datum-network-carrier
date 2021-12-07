package core

import (
	"github.com/RosettaFlow/Carrier-Go/common/traceutil"
	"github.com/RosettaFlow/Carrier-Go/lib/netmsg/common"
	taskmngpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/taskmng"
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
