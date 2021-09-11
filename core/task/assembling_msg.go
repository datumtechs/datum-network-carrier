package task

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	msgcommonpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/common"
	taskmngpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/taskmng"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func makeMsgOption(proposalId common.Hash,
	senderRole, receiverRole apipb.TaskRole,
	senderPartyId, receiverPartyId string,
	sender *apipb.TaskOrganization,
) *msgcommonpb.MsgOption {
	return &msgcommonpb.MsgOption{
		ProposalId:      proposalId.Bytes(),
		SenderRole:      uint64(senderRole),
		SenderPartyId:   []byte(senderPartyId),
		ReceiverRole:    uint64(receiverRole),
		ReceiverPartyId: []byte(receiverPartyId),
		MsgOwner: &msgcommonpb.TaskOrganizationIdentityInfo{
			Name:       []byte(sender.GetNodeName()),
			NodeId:     []byte(sender.GetNodeId()),
			IdentityId: []byte(sender.GetIdentityId()),
			PartyId:    []byte(sender.GetPartyId()),
		},
	}
}


func makeTaskResultMsg(
	proposalId common.Hash,
	senderRole, receiverRole apipb.TaskRole,
	senderPartyId, receiverPartyId string,
	task *types.Task,
	events []*libTypes.TaskEvent,
	startTime uint64,
) *taskmngpb.TaskResultMsg {
	return &taskmngpb.TaskResultMsg{
		MsgOption:     makeMsgOption(proposalId, senderRole, receiverRole, senderPartyId, receiverPartyId, task.GetTaskSender()),
		TaskEventList: types.ConvertTaskEventArr(events),
		CreateAt:      startTime,
		Sign:          nil,
	}
}


func fetchTaskResultMsg(msg *taskmngpb.TaskResultMsg) *types.TaskResultMsg {
	taskEventList := make([]*libTypes.TaskEvent, len(msg.TaskEventList))
	for index, value := range msg.TaskEventList {
		taskEventList[index] = &libTypes.TaskEvent{
			Type:       string(value.Type),
			TaskId:     string(value.TaskId),
			IdentityId: string(value.IdentityId),
			Content:    string(value.Content),
			CreateAt:   value.CreateAt,
		}
	}
	return &types.TaskResultMsg{
		MsgOption:     types.FetchMsgOption(msg.MsgOption),
		TaskEventList: taskEventList,
		CreateAt:      msg.CreateAt,
		Sign:          msg.Sign,
	}
}
