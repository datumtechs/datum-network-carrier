package twopc

import (
	"bytes"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	msgcommonpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/common"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func makeMsgOption(proposalId common.Hash,
	senderRole, receiverRole apicommonpb.TaskRole,
	senderPartyId, receiverPartyId string,
	sender *apicommonpb.TaskOrganization,
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

func makePrepareMsg(
	proposalId common.Hash,
	senderRole, receiverRole apicommonpb.TaskRole,
	senderPartyId, receiverPartyId string,
	task *types.Task,
	startTime uint64,
) (*twopcpb.PrepareMsg, error) {

	// region receivers come from task.Receivers
	bys := new(bytes.Buffer)
	err := task.EncodePb(bys)
	if err != nil {
		return nil, err
	}
	return &twopcpb.PrepareMsg{
		MsgOption: makeMsgOption(proposalId, senderRole, receiverRole, senderPartyId, receiverPartyId, task.GetTaskSender()),
		TaskInfo:  bys.Bytes(),
		CreateAt:  startTime,
		Sign:      nil,
	}, nil
}

func makePrepareVote(
	proposalId common.Hash,
	senderRole, receiverRole apicommonpb.TaskRole,
	senderPartyId, receiverPartyId string,
	task *types.Task,
	voteOption types.VoteOption,
	peerInfo *types.PrepareVoteResource,
	startTime uint64,
) *twopcpb.PrepareVote {

	return &twopcpb.PrepareVote{
		MsgOption:  makeMsgOption(proposalId, senderRole, receiverRole, senderPartyId, receiverPartyId, task.GetTaskSender()),
		VoteOption: voteOption.Bytes(),
		PeerInfo:   types.ConvertTaskPeerInfo(peerInfo),
		CreateAt:   startTime,
		Sign:       nil,
	}
}

func makeConfirmMsg(
	proposalId common.Hash,
	senderRole, receiverRole apicommonpb.TaskRole,
	senderPartyId, receiverPartyId string,
	task *types.Task,
	peers *twopcpb.ConfirmTaskPeerInfo,
	startTime uint64,
) *twopcpb.ConfirmMsg {

	msg := &twopcpb.ConfirmMsg{
		MsgOption: makeMsgOption(proposalId, senderRole, receiverRole, senderPartyId, receiverPartyId, task.GetTaskSender()),
		Peers:     peers,
		CreateAt:  startTime,
		Sign:      nil,
	}

	return msg
}

func makeConfirmVote(
	proposalId common.Hash,
	senderRole, receiverRole apicommonpb.TaskRole,
	senderPartyId, receiverPartyId string,
	task *types.Task,
	voteOption types.VoteOption,
	startTime uint64,
) *twopcpb.ConfirmVote {

	return &twopcpb.ConfirmVote{
		MsgOption:  makeMsgOption(proposalId, senderRole, receiverRole, senderPartyId, receiverPartyId, task.GetTaskSender()),
		VoteOption: voteOption.Bytes(),
		CreateAt:   startTime,
		Sign:       nil,
	}
}

func makeCommitMsg(
	proposalId common.Hash,
	senderRole, receiverRole apicommonpb.TaskRole,
	senderPartyId, receiverPartyId string,
	task *types.Task,
	startTime uint64,
) *twopcpb.CommitMsg {

	msg := &twopcpb.CommitMsg{
		MsgOption: makeMsgOption(proposalId, senderRole, receiverRole, senderPartyId, receiverPartyId, task.GetTaskSender()),
		CreateAt:  startTime,
		Sign:      nil,
	}
	return msg
}



func fetchPrepareMsg(msg *types.PrepareMsgWrap) (*types.PrepareMsg, error) {

	if nil == msg || nil == msg.TaskInfo {
		return nil, fmt.Errorf("receive nil prepareMsg or nil taskInfo")
	}

	task := types.NewTask(&libtypes.TaskPB{})
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

