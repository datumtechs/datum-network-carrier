package types

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
)

type ConsensusEngineType string

func (t ConsensusEngineType) String() string { return string(t) }

const (
	ChainconsTyp ConsensusEngineType = "ChainconsType"
	TwopcTyp     ConsensusEngineType = "TwopcType"
)

type ConsensusMsg interface {
	//Unmarshal
	String() string
	SealHash() common.Hash
	Hash() common.Hash
	Signature() []byte
}

type PrepareVoteResource struct {
	Id   string
	Ip   string
	Port string
}

func ConvertTaskPeerInfo(peerInfo *PrepareVoteResource) *pb.TaskPeerInfo {
	return &pb.TaskPeerInfo{
		Ip:   []byte(peerInfo.Ip),
		Port: []byte(peerInfo.Port),
	}
}
func FetchTaskPeerInfo(peerInfo *pb.TaskPeerInfo) *PrepareVoteResource {
	return &PrepareVoteResource{
		Ip:   string(peerInfo.Ip),
		Port: string(peerInfo.Port),
	}
}

type PrepareVote struct {
	ProposalId common.Hash
	TaskRole   TaskRole
	Owner      *NodeAlias
	VoteOption VoteOption
	PeerInfo   *PrepareVoteResource
	CreateAt   uint64
	Sign       []byte
}

func ConvertPrepareVote(vote *PrepareVote) *pb.PrepareVote {
	return &pb.PrepareVote{
		ProposalId: vote.ProposalId.Bytes(),
		TaskRole:   vote.TaskRole.Bytes(),
		Owner: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(vote.Owner.Name),
			NodeId:     []byte(vote.Owner.NodeId),
			IdentityId: []byte(vote.Owner.IdentityId),
		},
		VoteOption: vote.VoteOption.Bytes(),
		PeerInfo:   ConvertTaskPeerInfo(vote.PeerInfo),
		CreateAt:   vote.CreateAt,
		Sign:       vote.Sign,
	}
}
func FetchPrepareVote(vote *pb.PrepareVote) *PrepareVote {
	return &PrepareVote{
		ProposalId: common.BytesToHash(vote.ProposalId),
		TaskRole:   TaskRoleFromBytes(vote.TaskRole),
		Owner: &NodeAlias{
			Name:       string(vote.Owner.Name),
			NodeId:     string(vote.Owner.NodeId),
			IdentityId: string(vote.Owner.IdentityId),
		},
		VoteOption: VoteOptionFromBytes(vote.VoteOption),
		PeerInfo:   FetchTaskPeerInfo(vote.PeerInfo),
		CreateAt:   vote.CreateAt,
		Sign:       vote.Sign,
	}
}

type ConfirmMsg struct {
	ProposalId common.Hash
	TaskRole   TaskRole
	Epoch      uint64
	Owner      *NodeAlias
	CreateAt   uint64
	Sign       []byte
}

func ConvertConfirmMsg(msg *ConfirmMsg) *pb.ConfirmMsg {
	return &pb.ConfirmMsg{
		ProposalId: msg.ProposalId.Bytes(),
		TaskRole:   msg.TaskRole.Bytes(),
		Epoch:      msg.Epoch,
		Owner: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(msg.Owner.Name),
			NodeId:     []byte(msg.Owner.NodeId),
			IdentityId: []byte(msg.Owner.IdentityId),
		},
		CreateAt: msg.CreateAt,
		Sign:     msg.Sign,
	}
}
func FetchConfirmMsg(msg *pb.ConfirmMsg) *ConfirmMsg {
	return &ConfirmMsg{
		ProposalId: common.BytesToHash(msg.ProposalId),
		TaskRole:   TaskRoleFromBytes(msg.TaskRole),
		Epoch:      msg.Epoch,
		Owner: &NodeAlias{
			Name:       string(msg.Owner.Name),
			NodeId:     string(msg.Owner.NodeId),
			IdentityId: string(msg.Owner.IdentityId),
		},
		CreateAt: msg.CreateAt,
		Sign:     msg.Sign,
	}
}

type ConfirmVote struct {
	ProposalId common.Hash
	Epoch      uint64
	TaskRole   TaskRole
	Owner      *NodeAlias
	VoteOption VoteOption
	CreateAt   uint64
	Sign       []byte
}

func ConvertConfirmVote(vote *ConfirmVote) *pb.ConfirmVote {
	return &pb.ConfirmVote{
		ProposalId: vote.ProposalId.Bytes(),
		Epoch:      vote.Epoch,
		TaskRole:   vote.TaskRole.Bytes(),
		Owner: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(vote.Owner.Name),
			NodeId:     []byte(vote.Owner.NodeId),
			IdentityId: []byte(vote.Owner.IdentityId),
		},
		VoteOption: vote.VoteOption.Bytes(),
		CreateAt:   vote.CreateAt,
		Sign:       vote.Sign,
	}
}
func FetchConfirmVote(vote *pb.ConfirmVote) *ConfirmVote {
	return &ConfirmVote{
		ProposalId: common.BytesToHash(vote.ProposalId),
		Epoch:      vote.Epoch,
		TaskRole:   TaskRoleFromBytes(vote.TaskRole),
		Owner: &NodeAlias{
			Name:       string(vote.Owner.Name),
			NodeId:     string(vote.Owner.NodeId),
			IdentityId: string(vote.Owner.IdentityId),
		},
		VoteOption: VoteOptionFromBytes(vote.VoteOption),
		CreateAt:   vote.CreateAt,
		Sign:       vote.Sign,
	}
}

type CommitMsg struct {
	ProposalId common.Hash
	TaskRole   TaskRole
	Owner      *NodeAlias
	CreateAt   uint64
	Sign       []byte
}

func ConvertCommitMsg(msg *CommitMsg) *pb.CommitMsg {
	return &pb.CommitMsg{
		ProposalId: msg.ProposalId.Bytes(),
		TaskRole:   msg.TaskRole.Bytes(),
		Owner: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(msg.Owner.Name),
			NodeId:     []byte(msg.Owner.NodeId),
			IdentityId: []byte(msg.Owner.IdentityId),
		},
		CreateAt: msg.CreateAt,
		Sign:     msg.Sign,
	}
}
func FetchCommitMsg(msg *pb.CommitMsg) *CommitMsg {
	return &CommitMsg{
		ProposalId: common.BytesToHash(msg.ProposalId),
		TaskRole:   TaskRoleFromBytes(msg.TaskRole),
		Owner: &NodeAlias{
			Name:       string(msg.Owner.Name),
			NodeId:     string(msg.Owner.NodeId),
			IdentityId: string(msg.Owner.IdentityId),
		},
		CreateAt: msg.CreateAt,
		Sign:     msg.Sign,
	}
}

type TaskResultMsg struct {
	ProposalId    common.Hash
	TaskRole      TaskRole
	TaskId        string
	TaskEventList []*TaskEventInfo
	CreateAt      uint64
	Sign          []byte
}

func ConvertTaskResultMsg(msg *TaskResultMsg) *pb.TaskResultMsg {
	return &pb.TaskResultMsg{
		ProposalId:    msg.ProposalId.Bytes(),
		TaskRole:      msg.TaskRole.Bytes(),
		TaskId:        []byte(msg.TaskId),
		TaskEventList: ConvertTaskEventArr(msg.TaskEventList),
		CreateAt:      msg.CreateAt,
		Sign:          msg.Sign,
	}
}

func FetchTaskResultMsg(msg *pb.TaskResultMsg) *TaskResultMsg {
	return &TaskResultMsg{
		ProposalId:    common.BytesToHash(msg.ProposalId),
		TaskRole:      TaskRoleFromBytes(msg.TaskRole),
		TaskId:        string(msg.TaskId),
		TaskEventList: FetchTaskEventArr(msg.TaskEventList),
		CreateAt:      msg.CreateAt,
		Sign:          msg.Sign,
	}
}
