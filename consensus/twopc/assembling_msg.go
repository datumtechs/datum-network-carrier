package twopc

import (
	"bytes"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func makePrepareMsg(
	proposalId common.Hash,
	senderRole, receiverRole libtypes.TaskRole,
	senderPartyId, receiverPartyId string,
	nonConsTaks *types.NeedConsensusTask,
	startTime uint64,
) (*twopcpb.PrepareMsg, error) {

	// region receivers come from task.Receivers
	taskBytes := new(bytes.Buffer)
	err := nonConsTaks.GetTask().EncodePb(taskBytes)
	if err != nil {
		return nil, err
	}

	return &twopcpb.PrepareMsg{
		MsgOption: types.MakeMsgOption(proposalId, senderRole, receiverRole, senderPartyId, receiverPartyId, nonConsTaks.GetTask().GetTaskSender()),
		TaskInfo:  taskBytes.Bytes(),
		Evidence:  []byte(nonConsTaks.GetEvidence()),
		CreateAt:  startTime,
		Sign:      nil,
	}, nil
}

func makePrepareVote(
	proposalId common.Hash,
	senderRole, receiverRole libtypes.TaskRole,
	senderPartyId, receiverPartyId string,
	owner *libtypes.TaskOrganization,
	voteOption types.VoteOption,
	peerInfo *types.PrepareVoteResource,
	startTime uint64,
) *twopcpb.PrepareVote {

	return &twopcpb.PrepareVote{
		MsgOption:  types.MakeMsgOption(proposalId, senderRole, receiverRole, senderPartyId, receiverPartyId, owner),
		VoteOption: voteOption.Bytes(),
		PeerInfo:   types.ConvertTaskPeerInfo(peerInfo),
		CreateAt:   startTime,
		Sign:       nil,
	}
}

func makeConfirmMsg(
	proposalId common.Hash,
	senderRole, receiverRole libtypes.TaskRole,
	senderPartyId, receiverPartyId string,
	owner *libtypes.TaskOrganization,
	peers *twopcpb.ConfirmTaskPeerInfo,
	option types.TwopcMsgOption,
	startTime uint64,
) *twopcpb.ConfirmMsg {

	msg := &twopcpb.ConfirmMsg{
		MsgOption:     types.MakeMsgOption(proposalId, senderRole, receiverRole, senderPartyId, receiverPartyId, owner),
		ConfirmOption: option.Bytes(),
		Peers:         peers,
		CreateAt:      startTime,
		Sign:          nil,
	}

	return msg
}

func makeConfirmVote(
	proposalId common.Hash,
	senderRole, receiverRole libtypes.TaskRole,
	senderPartyId, receiverPartyId string,
	owner *libtypes.TaskOrganization,
	voteOption types.VoteOption,
	startTime uint64,
) *twopcpb.ConfirmVote {

	return &twopcpb.ConfirmVote{
		MsgOption:  types.MakeMsgOption(proposalId, senderRole, receiverRole, senderPartyId, receiverPartyId, owner),
		VoteOption: voteOption.Bytes(),
		CreateAt:   startTime,
		Sign:       nil,
	}
}

func makeCommitMsg(
	proposalId common.Hash,
	senderRole, receiverRole libtypes.TaskRole,
	senderPartyId, receiverPartyId string,
	owner *libtypes.TaskOrganization,
	option types.TwopcMsgOption,
	startTime uint64,
) *twopcpb.CommitMsg {

	msg := &twopcpb.CommitMsg{
		MsgOption:    types.MakeMsgOption(proposalId, senderRole, receiverRole, senderPartyId, receiverPartyId, owner),
		CommitOption: option.Bytes(),
		CreateAt:     startTime,
		Sign:         nil,
	}
	return msg
}

func fetchPrepareMsg(msg *types.PrepareMsgWrap) (*types.PrepareMsg, error) {

	if nil == msg || nil == msg.GetTaskInfo() {
		return nil, fmt.Errorf("receive nil prepareMsg or nil taskInfo")
	}

	task := types.NewTask(&libtypes.TaskPB{})
	err := task.DecodePb(msg.GetTaskInfo())
	if err != nil {
		return nil, fmt.Errorf("decode task info failed from prepareMsg, %s", err)
	}

	return &types.PrepareMsg{
			MsgOption: types.FetchMsgOption(msg.GetMsgOption()),
			TaskInfo:  task,
			Evidence:     string(msg.GetEvidence()),
			CreateAt:  msg.GetCreateAt(),
			Sign:      msg.GetSign(),
		},
		nil
}

func fetchProposalFromPrepareMsg(msg *types.PrepareMsg) *ctypes.ProposalTask {
	return &ctypes.ProposalTask{
		ProposalId: msg.GetMsgOption().GetProposalId(),
		TaskId:     msg.GetTask().GetTaskId(),
		CreateAt:   msg.GetCreateAt(),
	}
}

func fetchPrepareVote(vote *types.PrepareVoteWrap) *types.PrepareVote {
	return &types.PrepareVote{
		MsgOption:  types.FetchMsgOption(vote.GetMsgOption()),
		VoteOption: types.VoteOptionFromBytes(vote.GetVoteOption()),
		PeerInfo:   types.FetchTaskPeerInfo(vote.GetPeerInfo()),
		CreateAt:   vote.GetCreateAt(),
		Sign:       vote.GetSign(),
	}
}

func fetchConfirmMsg(msg *types.ConfirmMsgWrap) *types.ConfirmMsg {
	return &types.ConfirmMsg{
		MsgOption:     types.FetchMsgOption(msg.GetMsgOption()),
		ConfirmOption: types.TwopcMsgOptionFromBytes(msg.GetConfirmOption()),
		Peers:         msg.GetPeers(),
		CreateAt:      msg.GetCreateAt(),
		Sign:          msg.GetSign(),
	}
}

func fetchConfirmVote(vote *types.ConfirmVoteWrap) *types.ConfirmVote {
	return &types.ConfirmVote{
		MsgOption:  types.FetchMsgOption(vote.GetMsgOption()),
		VoteOption: types.VoteOptionFromBytes(vote.GetVoteOption()),
		CreateAt:   vote.GetCreateAt(),
		Sign:       vote.GetSign(),
	}
}

func fetchCommitMsg(msg *types.CommitMsgWrap) *types.CommitMsg {
	return &types.CommitMsg{
		MsgOption:    types.FetchMsgOption(msg.GetMsgOption()),
		CommitOption: types.TwopcMsgOptionFromBytes(msg.GetCommitOption()),
		CreateAt:     msg.GetCreateAt(),
		Sign:         msg.GetSign(),
	}
}
