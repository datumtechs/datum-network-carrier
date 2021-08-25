package twopc

import (
	"bytes"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func makePrepareMsgWithoutTaskRole(proposalId common.Hash, task *types.Task, startTime uint64) (*pb.PrepareMsg, error) {
	// region receivers come from task.Receivers
	bys := new(bytes.Buffer)
	err := task.EncodePb(bys)
	if err != nil {
		return nil, err
	}
	return &pb.PrepareMsg{
		ProposalId: proposalId.Bytes(),
		Owner: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(task.TaskData().NodeName),
			NodeId:     []byte(task.TaskData().NodeId),
			IdentityId: []byte(task.TaskData().IdentityId),
			PartyId:    []byte(task.TaskData().PartyId),
		},
		TaskInfo: bys.Bytes(),
		CreateAt: startTime,
	}, nil
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
			IdentityId: []byte(task.TaskData().IdentityId),
			PartyId:    []byte(task.TaskData().PartyId),
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
			IdentityId: []byte(task.TaskData().IdentityId),
			PartyId:    []byte(task.TaskData().PartyId),
		},
		CreateAt: startTime,
		Sign:     nil,
	}
	return msg
}

func makeTaskResultMsg(startTime uint64) *pb.TaskResultMsg {

	// 组装  TaskResultMsg
	return nil
}

func fetchPrepareMsg(prepareMsg *types.PrepareMsgWrap) (*types.PrepareMsg, error) {

	if nil == prepareMsg || nil == prepareMsg.TaskInfo {
		return nil, fmt.Errorf("receive nil prepareMsg or nil taskInfo")
	}

	task := types.NewTask(&libtypes.TaskData{})
	err := task.DecodePb(prepareMsg.TaskInfo)
	if err != nil {
		return nil, err
	}
	return &types.PrepareMsg{
			ProposalId:  common.BytesToHash(prepareMsg.ProposalId),
			TaskRole:    types.TaskRoleFromBytes(prepareMsg.TaskRole),
			TaskPartyId: string(prepareMsg.TaskPartyId),
			Owner: &types.TaskNodeAlias{
				Name:       string(prepareMsg.Owner.Name),
				NodeId:     string(prepareMsg.Owner.NodeId),
				IdentityId: string(prepareMsg.Owner.IdentityId),
				PartyId:    string(prepareMsg.Owner.PartyId),
			},
			TaskInfo: task,
			CreateAt: prepareMsg.CreateAt,
			Sign:     prepareMsg.Sign,
		},
		nil
}
func fetchProposalFromPrepareMsg(prepareMsg *types.PrepareMsg) *types.ProposalTask {
	return &types.ProposalTask{
		ProposalId: prepareMsg.ProposalId,
		Task:       prepareMsg.TaskInfo,
		CreateAt:   prepareMsg.CreateAt,
	}
}

func fetchPrepareVote(prepareVote *types.PrepareVoteWrap) (*types.PrepareVote, error) {
	//log.Debugf("=============================== partyId Len: %d", len(prepareVote.PeerInfo.PartyId))
	msg := &types.PrepareVote{
		ProposalId: common.BytesToHash(prepareVote.ProposalId),
		TaskRole:   types.TaskRoleFromBytes(prepareVote.TaskRole),
		Owner: &types.TaskNodeAlias{
			Name:       string(prepareVote.Owner.Name),
			NodeId:     string(prepareVote.Owner.NodeId),
			IdentityId: string(prepareVote.Owner.IdentityId),
			PartyId:    string(prepareVote.Owner.PartyId),
		},
		VoteOption: types.VoteOptionFromBytes(prepareVote.VoteOption),
		PeerInfo: &types.PrepareVoteResource{
			Id:      "",
			Ip:      string(prepareVote.PeerInfo.Ip),
			Port:    string(prepareVote.PeerInfo.Port),
			PartyId: string(prepareVote.PeerInfo.PartyId),
		},
		CreateAt: prepareVote.CreateAt,
		Sign:     prepareVote.Sign,
	}

	//if msg.VoteOption == types.Yes {
	//	msg.PeerInfo = &types.PrepareVoteResource{
	//		Id:      "",
	//		Ip:      string(prepareVote.PeerInfo.Ip),
	//		Port:    string(prepareVote.PeerInfo.Port),
	//		PartyId: string(prepareVote.PeerInfo.PartyId),
	//	}
	//}

	return msg, nil
}

func fetchConfirmMsg(confirmMsg *types.ConfirmMsgWrap) (*types.ConfirmMsg, error) {
	msg := &types.ConfirmMsg{
		ProposalId:  common.BytesToHash(confirmMsg.ProposalId),
		TaskRole:    types.TaskRoleFromBytes(confirmMsg.TaskRole),
		TaskPartyId: string(confirmMsg.TaskPartyId),
		Owner: &types.TaskNodeAlias{
			Name:       string(confirmMsg.Owner.Name),
			NodeId:     string(confirmMsg.Owner.NodeId),
			IdentityId: string(confirmMsg.Owner.IdentityId),
			PartyId:    string(confirmMsg.Owner.PartyId),
		},
		PeerDesc: confirmMsg.PeerDesc,
		CreateAt: confirmMsg.CreateAt,
		Sign:     confirmMsg.Sign,
	}
	return msg, nil
}

func fetchConfirmVote(confirmVote *types.ConfirmVoteWrap) (*types.ConfirmVote, error) {
	msg := &types.ConfirmVote{
		ProposalId: common.BytesToHash(confirmVote.ProposalId),
		TaskRole:   types.TaskRoleFromBytes(confirmVote.TaskRole),
		Owner: &types.TaskNodeAlias{
			Name:       string(confirmVote.Owner.Name),
			NodeId:     string(confirmVote.Owner.NodeId),
			IdentityId: string(confirmVote.Owner.IdentityId),
			PartyId:    string(confirmVote.Owner.PartyId),
		},
		VoteOption: types.VoteOptionFromBytes(confirmVote.VoteOption),
		CreateAt:   confirmVote.CreateAt,
		Sign:       confirmVote.Sign,
	}
	return msg, nil
}

func fetchCommitMsg(commitMsg *types.CommitMsgWrap) (*types.CommitMsg, error) {
	msg := &types.CommitMsg{
		ProposalId:  common.BytesToHash(commitMsg.ProposalId),
		TaskRole:    types.TaskRoleFromBytes(commitMsg.TaskRole),
		TaskPartyId: string(commitMsg.TaskPartyId),
		Owner: &types.TaskNodeAlias{
			Name:       string(commitMsg.Owner.Name),
			NodeId:     string(commitMsg.Owner.NodeId),
			IdentityId: string(commitMsg.Owner.IdentityId),
			PartyId:    string(commitMsg.Owner.PartyId),
		},
		CreateAt: commitMsg.CreateAt,
		Sign:     commitMsg.Sign,
	}
	return msg, nil
}

func fetchTaskResultMsg(commitMsg *types.TaskResultMsgWrap) (*types.TaskResultMsg, error) {
	taskEventList := make([]*types.TaskEventInfo, len(commitMsg.TaskEventList))
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
		ProposalId: common.BytesToHash(commitMsg.ProposalId),
		TaskRole:   types.TaskRoleFromBytes(commitMsg.TaskRole),
		Owner: &types.TaskNodeAlias{
			Name:       string(commitMsg.Owner.Name),
			NodeId:     string(commitMsg.Owner.NodeId),
			IdentityId: string(commitMsg.Owner.IdentityId),
			PartyId:    string(commitMsg.Owner.PartyId),
		},
		TaskId:        string(commitMsg.TaskId),
		TaskEventList: taskEventList,
		CreateAt:      commitMsg.CreateAt,
		Sign:          commitMsg.Sign,
	}
	return msg, nil
}
