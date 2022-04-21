package core

import (
	"github.com/Metisnetwork/Metis-Carrier/common/traceutil"
	"github.com/Metisnetwork/Metis-Carrier/lib/netmsg/common"
	taskmngpb "github.com/Metisnetwork/Metis-Carrier/lib/netmsg/taskmng"
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
