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
	Usage 		  *TaskResuorceUsage
	CreateAt      uint64
	Sign          []byte
}


func (msg *TaskResourceUsageMsg) String() string {
	return fmt.Sprintf(`{"msgOption": %s, "usage": %s, "createAt": %d, "sign": %v}`,
		msg.MsgOption.String(), msg.Usage.String(), msg.CreateAt, msg.Sign)
}

func (msg *TaskResourceUsageMsg) GetMsgOption ()  *MsgOption { return msg.MsgOption }
func (msg *TaskResourceUsageMsg) GetUsage ()  *TaskResuorceUsage { return msg.Usage }
func (msg *TaskResourceUsageMsg) GetCreateAt ()  uint64 { return msg.CreateAt }
func (msg *TaskResourceUsageMsg) GetSign()  []byte { return msg.Sign }

func FetchTaskResourceUsageMsg (msg *taskmngpb.TaskResourceUsageMsg) *TaskResourceUsageMsg {
	return &TaskResourceUsageMsg{
		MsgOption:  FetchMsgOption(msg.GetMsgOption()),
		Usage: NewTaskResuorceUsage(
			string(msg.GetTaskId()),
			string(msg.GetMsgOption().GetSenderPartyId()),
			msg.GetUsage().GetTotalMem(),
			msg.GetUsage().GetTotalBandwidth(),
			msg.GetUsage().GetTotalDisk(),
			msg.GetUsage().GetUsedMem(),
			msg.GetUsage().GetUsedBandwidth(),
			msg.GetUsage().GetUsedDisk(),
			uint32(msg.GetUsage().GetTotalProcessor()),
			uint32(msg.GetUsage().GetUsedProcessor()),
			),
		CreateAt:      msg.CreateAt,
		Sign:          msg.Sign,
	}
}

