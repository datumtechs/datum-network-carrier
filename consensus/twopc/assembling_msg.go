package twopc

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func makePrepareMsgWithoutTaskRole(proposalId common.Hash, task *types.Task, startTime uint64) *pb.PrepareMsg {
	// region receivers come from task.Receivers
	return &pb.PrepareMsg{
		ProposalId: proposalId.Bytes(),
		Owner: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(task.TaskData().NodeName),
			NodeId:     []byte(task.TaskData().NodeId),
			IdentityId: []byte(task.TaskData().Identity),
		},
		TaskInfo: task.TaskData(),
		CreateAt: startTime,
	}
}

func makePrepareVote() *pb.PrepareVote {

	// TODO 组装  PrepareVote
	return nil
}

func makeConfirmMsg(proposalId common.Hash, task *types.Task, startTime uint64) *pb.ConfirmMsg {
	msg := &pb.ConfirmMsg{
		ProposalId: proposalId.Bytes(),
		TaskRole:   nil,
		Owner: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(task.TaskData().NodeName),
			NodeId:     []byte(task.TaskData().NodeId),
			IdentityId: []byte(task.TaskData().Identity),
		},
		PeerDesc: nil,
		CreateAt: startTime,
		Sign:     nil,
	}

	return msg
}

func makeConfirmVote() *pb.ConfirmVote {

	// TODO 组装  ConfirmVote
	return nil
}

func makeCommitMsg(proposalId common.Hash, task *types.Task, startTime uint64) *pb.CommitMsg {
	msg := &pb.CommitMsg{
		ProposalId: proposalId.Bytes(),
		TaskRole:   nil,
		Owner: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(task.TaskData().NodeName),
			NodeId:     []byte(task.TaskData().NodeId),
			IdentityId: []byte(task.TaskData().Identity),
		},
		CreateAt: startTime,
		Sign:     nil,
	}
	return msg
}

func makeTaskResultMsg(startTime uint64) *pb.TaskResultMsg {

	// TODO 组装  TaskResultMsg
	return nil
}

func fetchPrepareMsg(prepareMsg *types.PrepareMsgWrap) (*types.ProposalTask, error) {

	return &types.ProposalTask{
		ProposalId: common.BytesToHash(prepareMsg.ProposalId),
		Task: types.NewTask(prepareMsg.TaskInfo),
		CreateAt: prepareMsg.CreateAt,
	}, nil
}
func fetchPrepareVote(prepareVote *types.PrepareVoteWrap) (*types.PrepareVote, error) {
	msg := &types.PrepareVote{
		ProposalId: common.BytesToHash(prepareVote.ProposalId),
		TaskRole:   types.TaskRoleFromBytes(prepareVote.TaskRole),
		Owner: &types.NodeAlias{
			Name:       string(prepareVote.Owner.Name),
			NodeId:     string(prepareVote.Owner.NodeId),
			IdentityId: string(prepareVote.Owner.IdentityId),
		},
		VoteOption: types.VoteOptionFromBytes(prepareVote.VoteOption),
		PeerInfo: &types.PrepareVoteResource{
			Id:   "",
			Ip:   string(prepareVote.PeerInfo.Ip),
			Port: string(prepareVote.PeerInfo.Port),
		},
		CreateAt: prepareVote.CreateAt,
		Sign:     prepareVote.Sign,
	}
	return msg, nil
}

func fetchConfirmMsg(confirmMsg *types.ConfirmMsgWrap) (*types.ConfirmMsg, error) {
	msg := &types.ConfirmMsg{
		ProposalId: common.BytesToHash(confirmMsg.ProposalId),
		TaskRole:   types.TaskRoleFromBytes(confirmMsg.TaskRole),
		Owner: &types.NodeAlias{
			Name:       string(confirmMsg.Owner.Name),
			NodeId:     string(confirmMsg.Owner.NodeId),
			IdentityId: string(confirmMsg.Owner.IdentityId),
		},
		CreateAt: confirmMsg.CreateAt,
		Sign:     confirmMsg.Sign,
	}
	return msg, nil
}

func fetchConfirmVote(confirmVote *types.ConfirmVoteWrap) (*types.ConfirmVote, error) {
	msg := &types.ConfirmVote{
		ProposalId: common.BytesToHash(confirmVote.ProposalId),
		TaskRole:   types.TaskRoleFromBytes(confirmVote.TaskRole),
		Owner: &types.NodeAlias{
			Name:       string(confirmVote.Owner.Name),
			NodeId:     string(confirmVote.Owner.NodeId),
			IdentityId: string(confirmVote.Owner.IdentityId),
		},
		VoteOption: types.VoteOptionFromBytes(confirmVote.VoteOption),
		CreateAt:   confirmVote.CreateAt,
		Sign:       confirmVote.Sign,
	}
	return msg, nil
}

func fetchCommitMsg(commitMsg *types.CommitMsgWrap) (*types.CommitMsg, error) {
	msg := &types.CommitMsg{
		ProposalId: common.BytesToHash(commitMsg.ProposalId),
		TaskRole:   types.TaskRoleFromBytes(commitMsg.TaskRole),
		Owner: &types.NodeAlias{
			Name:       string(commitMsg.Owner.Name),
			NodeId:     string(commitMsg.Owner.NodeId),
			IdentityId: string(commitMsg.Owner.IdentityId),
		},
		CreateAt: commitMsg.CreateAt,
		Sign:     commitMsg.Sign,
	}
	return msg, nil
}

func fetchTaskResultMsg(commitMsg *types.TaskResultMsgWrap) (*types.TaskResultMsg, error) {
	taskEventList := make([]*types.TaskEventInfo, 0)
	for index, value := range commitMsg.TaskEventList {
		taskEventList[index] = &types.TaskEventInfo{
			Type:       string(value.Type),
			TaskId:     string(value.TaskId),
			Identity:   string(value.IdentityId),
			Content:    string(value.Content),
			CreateTime: value.CreateAt,
		}
	}
	msg := &types.TaskResultMsg{
		ProposalId:    common.BytesToHash(commitMsg.ProposalId),
		TaskRole:      types.TaskRoleFromBytes(commitMsg.TaskRole),
		TaskId:        string(commitMsg.TaskId),
		TaskEventList: taskEventList,
		CreateAt:      commitMsg.CreateAt,
		Sign:          commitMsg.Sign,
	}
	return msg, nil
}
