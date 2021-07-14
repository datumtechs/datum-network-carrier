package twopc

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/common"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func makePrepareMsg(proposalId common.Hash, task *types.ScheduleTask, startTime uint64) *pb.PrepareMsg {
	// region receivers come from task.Receivers
	receivers := make([]*pb.ReceiverOption, 0)
	for index, value := range task.Receivers {
		providers := make([]*pb.TaskOrganizationIdentityInfo, 0)
		for i, v := range value.Providers {
			providers[i] = &pb.TaskOrganizationIdentityInfo{
				Name:       []byte(v.Name),
				NodeId:     []byte(v.NodeId),
				IdentityId: []byte(v.IdentityId),
			}
		}
		receivers[index] = &pb.ReceiverOption{
			MemberInfo: &pb.TaskOrganizationIdentityInfo{
				Name:       []byte(value.Name),
				NodeId:     []byte(value.NodeId),
				IdentityId: []byte(value.IdentityId),
			},
			Providers: providers,
		}
	}
	// endregion
	// region powerSupplier come from PowerSuppliers
	powerSupplier := make([]*pb.PowerSupplierOption, 0)
	for index, value := range task.PowerSuppliers {
		powerSupplier[index] = &pb.PowerSupplierOption{MemberInfo: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(value.Name),
			NodeId:     []byte(value.NodeId),
			IdentityId: []byte(value.IdentityId),
		}}
	}
	// endregion
	// region dataSupplier come form Partners
	dataSupplier := make([]*pb.DataSupplierOption, 0)
	for index, value := range task.Partners {
		dataSupplier[index] = &pb.DataSupplierOption{
			MemberInfo: &pb.TaskOrganizationIdentityInfo{
				Name:       []byte(value.Name),
				NodeId:     []byte(value.NodeId),
				IdentityId: []byte(value.IdentityId),
			},
			MetaDataId:      []byte(value.MetaData.MetaDataId),
			ColumnIndexList: value.MetaData.ColumnIndexList,
		}
	}
	//endregion
	msg := &pb.PrepareMsg{
		CreateAt:   startTime,
		ProposalId: proposalId.Bytes(),
		TaskOption: &pb.TaskOption{
			TaskRole: nil,
			TaskId:   []byte(task.TaskId),
			TaskName: []byte(task.TaskName),
			Owner: &pb.TaskOrganizationIdentityInfo{
				Name:       []byte(task.Owner.Name),
				NodeId:     []byte(task.Owner.NodeId),
				IdentityId: []byte(task.Owner.IdentityId),
			},
			AlgoSupplier: &pb.TaskOrganizationIdentityInfo{
				Name:       []byte(task.Owner.Name),
				NodeId:     []byte(task.Owner.NodeId),
				IdentityId: []byte(task.Owner.IdentityId),
			},
			DataSupplier:  dataSupplier,
			PowerSupplier: powerSupplier,
			Receivers:     receivers,
			OperationCost: &pb.TaskOperationCost{
				CostMem:       0,
				CostProcessor: task.OperationCost.Processor,
				CostBandwidth: task.OperationCost.Bandwidth,
				Duration:      task.OperationCost.Duration},
			CalculateContractCode: []byte(task.CalculateContractCode),
			DatasplitContractCode: []byte(task.DataSplitContractCode),
			CreateAt:              task.CreateAt,
		},
		Sign: nil,
	}
	return msg
}

func makePrepareVote() *pb.PrepareVote {

	// TODO 组装  PrepareVote
	return nil
}

func makeConfirmMsg(proposalId common.Hash, task *types.ScheduleTask, startTime uint64) *pb.ConfirmMsg {
	msg := &pb.ConfirmMsg{
		ProposalId: proposalId.Bytes(),
		TaskRole:   nil,
		Epoch:      2, //从哪里获取？
		Owner: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(task.Owner.Name),
			NodeId:     []byte(task.Owner.NodeId),
			IdentityId: []byte(task.Owner.IdentityId),
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

func makeCommitMsg(proposalId common.Hash, task *types.ScheduleTask, startTime uint64) *pb.CommitMsg {
	msg := &pb.CommitMsg{
		ProposalId: proposalId.Bytes(),
		TaskRole:   nil,
		Owner: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(task.Owner.Name),
			NodeId:     []byte(task.Owner.NodeId),
			IdentityId: []byte(task.Owner.IdentityId),
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
	//region partners
	var ownerMetaData *types.SupplierMetaData
	partners := make([]*types.ScheduleTaskDataSupplier, 0)
	for index, value := range prepareMsg.TaskOption.DataSupplier {
		partners[index] = &types.ScheduleTaskDataSupplier{
			TaskNodeAlias: &types.TaskNodeAlias{
				PartyId:    "",
				Name:       string(value.MemberInfo.Name),
				NodeId:     string(value.MemberInfo.NodeId),
				IdentityId: string(value.MemberInfo.IdentityId),
			},
			MetaData: &types.SupplierMetaData{
				MetaDataId:      string(value.MetaDataId),
				ColumnIndexList: value.ColumnIndexList,
			},
		}
		if bytes.Compare(prepareMsg.TaskOption.Owner.IdentityId, value.MemberInfo.IdentityId) == 0 {
			ownerMetaData = &types.SupplierMetaData{
				MetaDataId:      string(value.MetaDataId),
				ColumnIndexList: value.ColumnIndexList,
			}
		}
	}
	//endregion
	// region powerSuppliers
	powerSuppliers := make([]*types.ScheduleTaskPowerSupplier, 0)
	for index, value := range prepareMsg.TaskOption.PowerSupplier {
		powerSuppliers[index] = &types.ScheduleTaskPowerSupplier{
			TaskNodeAlias: &types.TaskNodeAlias{
				PartyId:    "",
				Name:       string(value.MemberInfo.Name),
				NodeId:     string(value.MemberInfo.NodeId),
				IdentityId: string(value.MemberInfo.IdentityId),
			},
		}
	}
	// endregion
	// region receivers
	receivers := make([]*types.ScheduleTaskResultReceiver, 0)
	for index, value := range prepareMsg.TaskOption.Receivers {
		providers := make([]*types.TaskNodeAlias, 0)
		for i, v := range providers {
			providers[i] = &types.TaskNodeAlias{
				PartyId:    "",
				Name:       v.Name,
				NodeId:     v.NodeId,
				IdentityId: v.IdentityId,
			}
		}
		receivers[index] = &types.ScheduleTaskResultReceiver{
			TaskNodeAlias: &types.TaskNodeAlias{
				PartyId:    "",
				Name:       string(value.MemberInfo.Name),
				NodeId:     string(value.MemberInfo.NodeId),
				IdentityId: string(value.MemberInfo.IdentityId),
			},
			Providers: providers,
		}
	}
	// endregion

	msg := &types.ProposalTask{
		ProposalId: common.BytesToHash(prepareMsg.ProposalId),
		ScheduleTask: &types.ScheduleTask{
			TaskId:   string(prepareMsg.TaskOption.TaskId),
			TaskName: string(prepareMsg.TaskOption.TaskName),
			Owner: &types.ScheduleTaskDataSupplier{
				TaskNodeAlias: &types.TaskNodeAlias{
					PartyId:    "",
					Name:       string(prepareMsg.TaskOption.Owner.Name),
					NodeId:     string(prepareMsg.TaskOption.Owner.NodeId),
					IdentityId: string(prepareMsg.TaskOption.Owner.IdentityId),
				},
				MetaData: ownerMetaData,
			},
			Partners:              partners,
			PowerSuppliers:        powerSuppliers,
			Receivers:             receivers,
			CalculateContractCode: string(prepareMsg.TaskOption.CalculateContractCode),
			DataSplitContractCode: string(prepareMsg.TaskOption.DatasplitContractCode),
			OperationCost: &types.TaskOperationCost{
				Processor: prepareMsg.TaskOption.OperationCost.CostProcessor,
				Mem:       prepareMsg.TaskOption.OperationCost.CostMem,
				Bandwidth: prepareMsg.TaskOption.OperationCost.CostBandwidth,
				Duration:  prepareMsg.TaskOption.OperationCost.Duration,
			},
			CreateAt: prepareMsg.TaskOption.CreateAt,
		},
		CreateAt: prepareMsg.CreateAt,
	}
	return msg, nil
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
		Epoch:      confirmVote.Epoch,
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
