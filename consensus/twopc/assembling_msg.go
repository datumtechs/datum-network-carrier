package twopc

import (
	"bytes"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func makeMsgOption(proposalId common.Hash,
	senderRole, receiverRole apipb.TaskRole,
	senderPartyId, receiverPartyId string,
	sender *apipb.TaskOrganization,
) *pb.MsgOption {
	return &pb.MsgOption{
		ProposalId:      proposalId.Bytes(),
		SenderRole:      uint64(senderRole),
		SenderPartyId:   []byte(senderPartyId),
		ReceiverRole:    uint64(receiverRole),
		ReceiverPartyId: []byte(receiverPartyId),
		MsgOwner: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(sender.GetNodeName()),
			NodeId:     []byte(sender.GetNodeId()),
			IdentityId: []byte(sender.GetIdentityId()),
			PartyId:    []byte(sender.GetPartyId()),
		},
	}
}

func makePrepareMsg(
	proposalId common.Hash,
	senderRole, receiverRole apipb.TaskRole,
	senderPartyId, receiverPartyId string,
	task *types.Task,
	startTime uint64,
) (*pb.PrepareMsg, error) {

	// region receivers come from task.Receivers
	bys := new(bytes.Buffer)
	err := task.EncodePb(bys)
	if err != nil {
		return nil, err
	}
	return &pb.PrepareMsg{
		MsgOption: makeMsgOption(proposalId, senderRole, receiverRole, senderPartyId, receiverPartyId, task.GetTaskSender()),
		TaskInfo:  bys.Bytes(),
		CreateAt:  startTime,
		Sign:      nil,
	}, nil
}

func makePrepareVote(
	proposalId common.Hash,
	senderRole, receiverRole apipb.TaskRole,
	senderPartyId, receiverPartyId string,
	task *types.Task,
	voteOption types.VoteOption,
	peerInfo *types.PrepareVoteResource,
	startTime uint64,
) *pb.PrepareVote {

	return &pb.PrepareVote{
		MsgOption:  makeMsgOption(proposalId, senderRole, receiverRole, senderPartyId, receiverPartyId, task.GetTaskSender()),
		VoteOption: voteOption.Bytes(),
		PeerInfo:   types.ConvertTaskPeerInfo(peerInfo),
		CreateAt:   startTime,
		Sign:       nil,
	}
}

func makeConfirmMsg(
	proposalId common.Hash,
	senderRole, receiverRole apipb.TaskRole,
	senderPartyId, receiverPartyId string,
	task *types.Task,
	peers *pb.ConfirmTaskPeerInfo,
	startTime uint64,
) *pb.ConfirmMsg {

	msg := &pb.ConfirmMsg{
		MsgOption: makeMsgOption(proposalId, senderRole, receiverRole, senderPartyId, receiverPartyId, task.GetTaskSender()),
		Peers:     peers,
		CreateAt:  startTime,
		Sign:      nil,
	}

	return msg
}

func makeConfirmVote(
	proposalId common.Hash,
	senderRole, receiverRole apipb.TaskRole,
	senderPartyId, receiverPartyId string,
	task *types.Task,
	voteOption types.VoteOption,
	startTime uint64,
) *pb.ConfirmVote {

	return &pb.ConfirmVote{
		MsgOption:  makeMsgOption(proposalId, senderRole, receiverRole, senderPartyId, receiverPartyId, task.GetTaskSender()),
		VoteOption: voteOption.Bytes(),
		CreateAt:   startTime,
		Sign:       nil,
	}
}

func makeCommitMsg(
	proposalId common.Hash,
	senderRole, receiverRole apipb.TaskRole,
	senderPartyId, receiverPartyId string,
	task *types.Task,
	startTime uint64,
) *pb.CommitMsg {

	msg := &pb.CommitMsg{
		MsgOption: makeMsgOption(proposalId, senderRole, receiverRole, senderPartyId, receiverPartyId, task.GetTaskSender()),
		CreateAt:  startTime,
		Sign:      nil,
	}
	return msg
}

func makeTaskResultMsg(
	proposalId common.Hash,
	senderRole, receiverRole apipb.TaskRole,
	senderPartyId, receiverPartyId string,
	task *types.Task,
	events []*libTypes.TaskEvent,
	startTime uint64,
) *pb.TaskResultMsg {
	return &pb.TaskResultMsg{
		MsgOption:     makeMsgOption(proposalId, senderRole, receiverRole, senderPartyId, receiverPartyId, task.GetTaskSender()),
		TaskEventList: types.ConvertTaskEventArr(events),
		CreateAt:      startTime,
		Sign:          nil,
	}
}

func fetchPrepareMsg(msg *types.PrepareMsgWrap) (*types.PrepareMsg, error) {

	if nil == msg || nil == msg.TaskInfo {
		return nil, fmt.Errorf("receive nil prepareMsg or nil taskInfo")
	}

	task := types.NewTask(&libTypes.TaskPB{})
	err := task.DecodePb(msg.TaskInfo)
	if err != nil {
		return nil, err
	}
	return &types.PrepareMsg{
			MsgOption: types.FetchMsgOption(msg.MsgOption),
			TaskInfo:  task,
			CreateAt:  msg.CreateAt,
			Sign:      msg.Sign,
		},
		nil
}
func fetchProposalFromPrepareMsg(msg *types.PrepareMsg) *types.ProposalTask {
	return &types.ProposalTask{
		ProposalId: msg.MsgOption.ProposalId,
		Task:       msg.TaskInfo,
		CreateAt:   msg.CreateAt,
	}
}

func fetchPrepareVote(vote *types.PrepareVoteWrap) *types.PrepareVote {
	return &types.PrepareVote{
		MsgOption:  types.FetchMsgOption(vote.MsgOption),
		VoteOption: types.VoteOptionFromBytes(vote.VoteOption),
		PeerInfo:   types.FetchTaskPeerInfo(vote.PeerInfo),
		CreateAt:   vote.CreateAt,
		Sign:       vote.Sign,
	}
}

func fetchConfirmMsg(msg *types.ConfirmMsgWrap) *types.ConfirmMsg {
	return &types.ConfirmMsg{
		MsgOption: types.FetchMsgOption(msg.MsgOption),
		Peers:     msg.Peers,
		CreateAt:  msg.CreateAt,
		Sign:      msg.Sign,
	}
}

func fetchConfirmVote(vote *types.ConfirmVoteWrap) *types.ConfirmVote {
	return &types.ConfirmVote{
		MsgOption:  types.FetchMsgOption(vote.MsgOption),
		VoteOption: types.VoteOptionFromBytes(vote.VoteOption),
		CreateAt:   vote.CreateAt,
		Sign:       vote.Sign,
	}
}

func fetchCommitMsg(msg *types.CommitMsgWrap) *types.CommitMsg {
	return &types.CommitMsg{
		MsgOption: types.FetchMsgOption(msg.MsgOption),
		CreateAt:  msg.CreateAt,
		Sign:      msg.Sign,
	}
}

func fetchTaskResultMsg(msg *types.TaskResultMsgWrap) *types.TaskResultMsg {
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
