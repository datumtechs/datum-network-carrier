package types

import (
	"fmt"
	taskmngpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/taskmng"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)

type TaskResultMsg struct {
	MsgOption     *MsgOption
	TaskEventList []*libtypes.TaskEvent
	CreateAt      uint64
	Sign          []byte
}

func (msg *TaskResultMsg) String() string {
	return fmt.Sprintf(`{"msgOption": %s, "createAt": %d, "sign": %v}`,
		msg.MsgOption.String(), msg.CreateAt, msg.Sign)
}

func ConvertTaskResultMsg(msg *TaskResultMsg) *taskmngpb.TaskResultMsg {
	return &taskmngpb.TaskResultMsg{
		MsgOption: ConvertMsgOption(msg.MsgOption),
		TaskEventList: ConvertTaskEventArr(msg.TaskEventList),
		CreateAt:      msg.CreateAt,
		Sign:          msg.Sign,
	}
}

func FetchTaskResultMsg(msg *taskmngpb.TaskResultMsg) *TaskResultMsg {
	return &TaskResultMsg{
		MsgOption:  FetchMsgOption(msg.GetMsgOption()),
		TaskEventList: FetchTaskEventArr(msg.TaskEventList),
		CreateAt:      msg.CreateAt,
		Sign:          msg.Sign,
	}
}

type TaskResourceUsageMsg struct {
	MsgOption     *MsgOption
	Usage 		  TaskResuorceUsage
	CreateAt      uint64
	Sign          []byte
}


