package types

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
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
	Id      string
	Ip      string
	Port    string
	PartyId string
}

func NewPrepareVoteResource(id, ip, port, partyId string) *PrepareVoteResource {
	return &PrepareVoteResource{
		Id: id,
		Ip: ip,
		Port: port,
		PartyId: partyId,
	}
}
func (resource *PrepareVoteResource) String() string {
	return fmt.Sprintf(`{"id": %s, "ip": %s, "port": %s, "partyId": %s}`, resource.Id, resource.Ip, resource.Port, resource.PartyId)
}
func ConvertTaskPeerInfo(peerInfo *PrepareVoteResource) *pb.TaskPeerInfo {
	if nil == peerInfo {
		return &pb.TaskPeerInfo{}
	}
	return &pb.TaskPeerInfo{
		Ip:      []byte(peerInfo.Ip),
		Port:    []byte(peerInfo.Port),
		PartyId: []byte(peerInfo.PartyId),
	}
}
func FetchTaskPeerInfo(peerInfo *pb.TaskPeerInfo) *PrepareVoteResource {
	if nil == peerInfo {
		return &PrepareVoteResource{}
	}
	return &PrepareVoteResource{
		Ip:      string(peerInfo.Ip),
		Port:    string(peerInfo.Port),
		PartyId: string(peerInfo.PartyId),
	}
}
func ConvertTaskPeerInfoArr(resourceArr []*PrepareVoteResource) []*pb.TaskPeerInfo {

	arr := make([]*pb.TaskPeerInfo, len(resourceArr))

	for i, resource := range resourceArr {
		var peerInfo *pb.TaskPeerInfo
		if nil == resource {
			peerInfo = &pb.TaskPeerInfo{}
		} else {
			peerInfo = &pb.TaskPeerInfo{
				Ip:      []byte(resource.Ip),
				Port:    []byte(resource.Port),
				PartyId: []byte(resource.PartyId),
			}
		}

		arr[i] = peerInfo
	}
	return arr
}
func FetchTaskPeerInfoArr(peerInfoArr []*pb.TaskPeerInfo) []*PrepareVoteResource {

	arr := make([]*PrepareVoteResource, len(peerInfoArr))

	for i, peerInfo := range peerInfoArr {
		var resource *PrepareVoteResource
		if nil == peerInfo {
			resource = &PrepareVoteResource{}
		} else {
			resource = &PrepareVoteResource{
				Ip:      string(peerInfo.Ip),
				Port:    string(peerInfo.Port),
				PartyId: string(peerInfo.PartyId),
			}
		}
		arr[i] = resource
	}
	return arr
}

type PrepareMsg struct {
	ProposalId  common.Hash
	TaskRole    TaskRole
	TaskPartyId string
	Owner       *apipb.TaskOrganization
	TaskInfo    *Task
	CreateAt    uint64
	Sign        []byte
}

func (msg *PrepareMsg) String() string {
	return fmt.Sprintf(`{"proposalId": %s, "taskRole": %s, "taskPartyId": %s, "owner": %s, "createAt": %d, "sign": %v}`,
		msg.ProposalId.String(), msg.TaskRole.String(), msg.TaskPartyId, msg.Owner.String(), msg.CreateAt, msg.Sign)
}

type PrepareVote struct {
	ProposalId common.Hash
	TaskRole   TaskRole
	Owner      *apipb.TaskOrganization
	VoteOption VoteOption
	PeerInfo   *PrepareVoteResource
	CreateAt   uint64
	Sign       []byte
}

func (vote *PrepareVote) String() string {
	return fmt.Sprintf(`{"proposalId": %s, "taskRole": %s, "owner": %s, "voteOption": %s, "peerInfo": %s, "createAt": %d, "sign": %v}`,
		vote.ProposalId.String(), vote.TaskRole.String(), vote.Owner.String(), vote.VoteOption.String(), vote.PeerInfo.String(), vote.CreateAt, vote.Sign)
}

func ConvertPrepareVote(vote *PrepareVote) *pb.PrepareVote {
	return &pb.PrepareVote{
		ProposalId: vote.ProposalId.Bytes(),
		TaskRole:   vote.TaskRole.Bytes(),
		Owner: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(vote.Owner.NodeName),
			NodeId:     []byte(vote.Owner.NodeId),
			IdentityId: []byte(vote.Owner.IdentityId),
			PartyId:    []byte(vote.Owner.PartyId),
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
		Owner: &apipb.TaskOrganization{
			NodeName:   string(vote.Owner.Name),
			NodeId:     string(vote.Owner.NodeId),
			IdentityId: string(vote.Owner.IdentityId),
			PartyId:    string(vote.Owner.PartyId),
		},
		VoteOption: VoteOptionFromBytes(vote.VoteOption),
		PeerInfo:   FetchTaskPeerInfo(vote.PeerInfo),
		CreateAt:   vote.CreateAt,
		Sign:       vote.Sign,
	}
}

type ConfirmMsg struct {
	ProposalId  common.Hash
	TaskRole    TaskRole
	TaskPartyId string
	Owner       *apipb.TaskOrganization
	PeerDesc    *pb.ConfirmTaskPeerInfo
	CreateAt    uint64
	Sign        []byte
}

func (msg *ConfirmMsg) String() string {
	return fmt.Sprintf(`{"proposalId": %s, "taskRole": %s, "taskPartyId": %s, "owner": %s, "peerDesc": %s, "createAt": %d, "sign": %v}`,
		msg.ProposalId.String(), msg.TaskRole.String(), msg.TaskPartyId, msg.Owner.String(), msg.PeerDesc.String(), msg.CreateAt, msg.Sign)
}

func ConvertConfirmMsg(msg *ConfirmMsg) *pb.ConfirmMsg {
	return &pb.ConfirmMsg{
		ProposalId:  msg.ProposalId.Bytes(),
		TaskRole:    msg.TaskRole.Bytes(),
		TaskPartyId: []byte(msg.TaskPartyId),
		Owner: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(msg.Owner.NodeName),
			NodeId:     []byte(msg.Owner.NodeId),
			IdentityId: []byte(msg.Owner.IdentityId),
			PartyId:    []byte(msg.Owner.PartyId),
		},
		PeerDesc: msg.PeerDesc,
		CreateAt: msg.CreateAt,
		Sign:     msg.Sign,
	}
}
func FetchConfirmMsg(msg *pb.ConfirmMsg) *ConfirmMsg {
	return &ConfirmMsg{
		ProposalId:  common.BytesToHash(msg.ProposalId),
		TaskRole:    TaskRoleFromBytes(msg.TaskRole),
		TaskPartyId: string(msg.TaskPartyId),
		Owner: &apipb.TaskOrganization{
			NodeName:   string(msg.Owner.Name),
			NodeId:     string(msg.Owner.NodeId),
			IdentityId: string(msg.Owner.IdentityId),
			PartyId:    string(msg.Owner.PartyId),
		},
		PeerDesc: msg.PeerDesc,
		CreateAt: msg.CreateAt,
		Sign:     msg.Sign,
	}
}

type ConfirmVote struct {
	ProposalId common.Hash
	TaskRole   TaskRole
	Owner      *apipb.TaskOrganization
	VoteOption VoteOption
	CreateAt   uint64
	Sign       []byte
}

func (vote *ConfirmVote) String() string {
	return fmt.Sprintf(`{"proposalId": %s, "taskRole": %s, "owner": %s, "voteOption": %s, "createAt": %d, "sign": %v}`,
		vote.ProposalId.String(), vote.TaskRole.String(), vote.Owner.String(), vote.VoteOption.String(), vote.CreateAt, vote.Sign)
}

func ConvertConfirmVote(vote *ConfirmVote) *pb.ConfirmVote {
	return &pb.ConfirmVote{
		ProposalId: vote.ProposalId.Bytes(),
		TaskRole:   vote.TaskRole.Bytes(),
		Owner: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(vote.Owner.NodeName),
			NodeId:     []byte(vote.Owner.NodeId),
			IdentityId: []byte(vote.Owner.IdentityId),
			PartyId:    []byte(vote.Owner.PartyId),
		},
		VoteOption: vote.VoteOption.Bytes(),
		CreateAt:   vote.CreateAt,
		Sign:       vote.Sign,
	}
}
func FetchConfirmVote(vote *pb.ConfirmVote) *ConfirmVote {
	return &ConfirmVote{
		ProposalId: common.BytesToHash(vote.ProposalId),
		TaskRole:   TaskRoleFromBytes(vote.TaskRole),
		Owner: &apipb.TaskOrganization{
			NodeName:   string(vote.Owner.Name),
			NodeId:     string(vote.Owner.NodeId),
			IdentityId: string(vote.Owner.IdentityId),
			PartyId:    string(vote.Owner.PartyId),
		},
		VoteOption: VoteOptionFromBytes(vote.VoteOption),
		CreateAt:   vote.CreateAt,
		Sign:       vote.Sign,
	}
}

type CommitMsg struct {
	ProposalId  common.Hash
	TaskRole    TaskRole
	TaskPartyId string
	Owner       *apipb.TaskOrganization
	CreateAt    uint64
	Sign        []byte
}

func (msg *CommitMsg) String() string {
	return fmt.Sprintf(`{"proposalId": %s, "taskRole": %s, "taskPartyId": %s, "owner": %s, "createAt": %d, "sign": %v}`,
		msg.ProposalId.String(), msg.TaskRole.String(), msg.TaskPartyId, msg.Owner.String(), msg.CreateAt, msg.Sign)
}

func ConvertCommitMsg(msg *CommitMsg) *pb.CommitMsg {
	return &pb.CommitMsg{
		ProposalId:  msg.ProposalId.Bytes(),
		TaskRole:    msg.TaskRole.Bytes(),
		TaskPartyId: []byte(msg.TaskPartyId),
		Owner: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(msg.Owner.NodeName),
			NodeId:     []byte(msg.Owner.NodeId),
			IdentityId: []byte(msg.Owner.IdentityId),
			PartyId:    []byte(msg.Owner.PartyId),
		},
		CreateAt: msg.CreateAt,
		Sign:     msg.Sign,
	}
}
func FetchCommitMsg(msg *pb.CommitMsg) *CommitMsg {
	return &CommitMsg{
		ProposalId:  common.BytesToHash(msg.ProposalId),
		TaskRole:    TaskRoleFromBytes(msg.TaskRole),
		TaskPartyId: string(msg.TaskPartyId),
		Owner: &apipb.TaskOrganization{
			NodeName:   string(msg.Owner.Name),
			NodeId:     string(msg.Owner.NodeId),
			IdentityId: string(msg.Owner.IdentityId),
			PartyId:    string(msg.Owner.PartyId),
		},
		CreateAt: msg.CreateAt,
		Sign:     msg.Sign,
	}
}

type TaskResultMsg struct {
	ProposalId    common.Hash
	TaskRole      TaskRole
	Owner         *apipb.TaskOrganization
	TaskId        string
	TaskEventList []*libTypes.TaskEvent
	CreateAt      uint64
	Sign          []byte
}

func (msg *TaskResultMsg) String() string {
	return fmt.Sprintf(`{"proposalId": %s, "taskRole": %s, "owner": %s, "taskId": %s, "createAt": %d, "sign": %v}`,
		msg.ProposalId.String(), msg.TaskRole.String(), msg.Owner.String(), msg.TaskId, msg.CreateAt, msg.Sign)
}

func ConvertTaskResultMsg(msg *TaskResultMsg) *pb.TaskResultMsg {
	return &pb.TaskResultMsg{
		ProposalId: msg.ProposalId.Bytes(),
		TaskRole:   msg.TaskRole.Bytes(),
		Owner: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(msg.Owner.NodeName),
			NodeId:     []byte(msg.Owner.NodeId),
			IdentityId: []byte(msg.Owner.IdentityId),
			PartyId:    []byte(msg.Owner.PartyId),
		},
		TaskId:        []byte(msg.TaskId),
		TaskEventList: ConvertTaskEventArr(msg.TaskEventList),
		CreateAt:      msg.CreateAt,
		Sign:          msg.Sign,
	}
}

func FetchTaskResultMsg(msg *pb.TaskResultMsg) *TaskResultMsg {
	return &TaskResultMsg{
		ProposalId: common.BytesToHash(msg.ProposalId),
		TaskRole:   TaskRoleFromBytes(msg.TaskRole),
		Owner: &apipb.TaskOrganization{
			NodeName:   string(msg.Owner.Name),
			NodeId:     string(msg.Owner.NodeId),
			IdentityId: string(msg.Owner.IdentityId),
			PartyId:    string(msg.Owner.PartyId),
		},
		TaskId:        string(msg.TaskId),
		TaskEventList: FetchTaskEventArr(msg.TaskEventList),
		CreateAt:      msg.CreateAt,
		Sign:          msg.Sign,
	}
}
