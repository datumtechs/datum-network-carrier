package types

import (
	"fmt"
	taskmngpb "github.com/Metisnetwork/Metis-Carrier/lib/netmsg/taskmng"
	libtypes "github.com/Metisnetwork/Metis-Carrier/lib/types"
)

type TaskResultMsg struct {
	MsgOption     *MsgOption
	TaskEventList []*libtypes.TaskEvent
	CreateAt      uint64
	Sign          []byte
}

func (msg *TaskResultMsg) String() string {
	return fmt.Sprintf(`{"msgOption": %s, "createAt": %d, "sign": %v}`,
		msg.GetMsgOption().String(), msg.GetCreateAt(), msg.GetSign())
}

func (msg *TaskResultMsg) GetMsgOption() *MsgOption                { return msg.MsgOption }
func (msg *TaskResultMsg) GetTaskEventList() []*libtypes.TaskEvent { return msg.TaskEventList }
func (msg *TaskResultMsg) GetCreateAt() uint64                     { return msg.CreateAt }
func (msg *TaskResultMsg) GetSign() []byte                         { return msg.Sign }

func ConvertTaskResultMsg(msg *TaskResultMsg) *taskmngpb.TaskResultMsg {
	return &taskmngpb.TaskResultMsg{
		MsgOption:     ConvertMsgOption(msg.GetMsgOption()),
		TaskEventList: ConvertTaskEventArr(msg.GetTaskEventList()),
		CreateAt:      msg.GetCreateAt(),
		Sign:          msg.GetSign(),
	}
}

func FetchTaskResultMsg(msg *taskmngpb.TaskResultMsg) *TaskResultMsg {
	return &TaskResultMsg{
		MsgOption:     FetchMsgOption(msg.GetMsgOption()),
		TaskEventList: FetchTaskEventArr(msg.GetTaskEventList()),
		CreateAt:      msg.GetCreateAt(),
		Sign:          msg.GetSign(),
	}
}

type TaskResourceUsageMsg struct {
	MsgOption *MsgOption
	Usage     *TaskResuorceUsage
	CreateAt  uint64
	Sign      []byte
}

func (msg *TaskResourceUsageMsg) String() string {
	return fmt.Sprintf(`{"msgOption": %s, "usage": %s, "createAt": %d, "sign": %v}`,
		msg.GetMsgOption().String(), msg.GetUsage().String(), msg.GetCreateAt(), msg.GetSign())
}

func (msg *TaskResourceUsageMsg) GetMsgOption() *MsgOption     { return msg.MsgOption }
func (msg *TaskResourceUsageMsg) GetUsage() *TaskResuorceUsage { return msg.Usage }
func (msg *TaskResourceUsageMsg) GetCreateAt() uint64          { return msg.CreateAt }
func (msg *TaskResourceUsageMsg) GetSign() []byte              { return msg.Sign }

func FetchTaskResourceUsageMsg(msg *taskmngpb.TaskResourceUsageMsg) *TaskResourceUsageMsg {
	return &TaskResourceUsageMsg{
		MsgOption: FetchMsgOption(msg.GetMsgOption()),
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
		CreateAt: msg.GetCreateAt(),
		Sign:     msg.GetSign(),
	}
}

type TaskTerminateTaskMngMsg struct {
	MsgOption *MsgOption
	TaskId    string
	CreateAt  uint64
	Sign      []byte
}

func (msg *TaskTerminateTaskMngMsg) String() string {
	return fmt.Sprintf(`{"msgOption": %s, "taskId": %s, "createAt": %d, "sign": %v}`,
		msg.GetMsgOption().String(), msg.GetTaskId(), msg.GetCreateAt(), msg.GetSign())
}

func (msg *TaskTerminateTaskMngMsg) GetMsgOption() *MsgOption { return msg.MsgOption }
func (msg *TaskTerminateTaskMngMsg) GetTaskId() string        { return msg.TaskId }
func (msg *TaskTerminateTaskMngMsg) GetCreateAt() uint64      { return msg.CreateAt }
func (msg *TaskTerminateTaskMngMsg) GetSign() []byte          { return msg.Sign }

func FetchTaskTerminateTaskMngMsg(msg *taskmngpb.TaskTerminateMsg) *TaskTerminateTaskMngMsg {
	return &TaskTerminateTaskMngMsg{
		MsgOption: FetchMsgOption(msg.GetMsgOption()),
		TaskId:    string(msg.GetTaskId()),
		CreateAt:  msg.GetCreateAt(),
		Sign:      msg.GetSign(),
	}
}
