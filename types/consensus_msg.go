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
		Id:      id,
		Ip:      ip,
		Port:    port,
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

type MsgOption struct {
	ProposalId      common.Hash
	SenderRole      apipb.TaskRole
	SenderPartyId   string
	ReceiverRole    apipb.TaskRole
	ReceiverPartyId string
	Owner           *apipb.TaskOrganization
}

func (option *MsgOption) String() string {
	return fmt.Sprintf(`{"ProposalId": %s, "senderRole": %s, "senderPartyId": %s, "receiverRole": %s, "receiverPartyId": %s, "owner": %s}`,
		option.ProposalId.String(), option.SenderRole.String(), option.SenderPartyId, option.ReceiverRole.String(), option.ReceiverPartyId, option.Owner.String())
}

func ConvertMsgOption(option *MsgOption) *pb.MsgOption {
	return &pb.MsgOption{
		ProposalId:      option.ProposalId.Bytes(),
		SenderRole:      uint64(option.SenderRole),
		SenderPartyId:   []byte(option.SenderPartyId),
		ReceiverRole:    uint64(option.ReceiverRole),
		ReceiverPartyId: []byte(option.ReceiverPartyId),
		MsgOwner: &pb.TaskOrganizationIdentityInfo{
			Name:       []byte(option.Owner.GetNodeName()),
			NodeId:     []byte(option.Owner.GetNodeId()),
			IdentityId: []byte(option.Owner.GetIdentityId()),
			PartyId:    []byte(option.Owner.GetPartyId()),
		},
	}
}

func FetchMsgOption(option *pb.MsgOption) *MsgOption {
	return &MsgOption{
		ProposalId:      common.BytesToHash(option.ProposalId),
		SenderRole:      apipb.TaskRole(option.GetSenderRole()),
		SenderPartyId:   string(option.SenderPartyId),
		ReceiverRole:    apipb.TaskRole(option.GetReceiverRole()),
		ReceiverPartyId: string(option.ReceiverPartyId),
		Owner: &apipb.TaskOrganization{
			NodeName:   string(option.GetMsgOwner().GetName()),
			NodeId:     string(option.GetMsgOwner().GetNodeId()),
			IdentityId: string(option.GetMsgOwner().GetIdentityId()),
			PartyId:    string(option.GetMsgOwner().GetPartyId()),
		},
	}
}

type PrepareMsg struct {
	MsgOption *MsgOption
	TaskInfo  *Task
	CreateAt  uint64
	Sign      []byte
}

func (msg *PrepareMsg) String() string {
	return fmt.Sprintf(`{"msgOption": %s, "createAt": %d, "sign": %v}`,
		msg.MsgOption.String(), msg.CreateAt, msg.Sign)
}
func (msg *PrepareMsg) StringWithTask() string {
	return fmt.Sprintf(`{"msgOption": %s, "taskInfo": %s, "createAt": %d, "sign": %v}`,
		msg.MsgOption.String(), msg.TaskInfo.GetTaskData().String(), msg.CreateAt, msg.Sign)
}

type PrepareVote struct {
	MsgOption  *MsgOption
	VoteOption VoteOption
	PeerInfo   *PrepareVoteResource
	CreateAt   uint64
	Sign       []byte
}

func (vote *PrepareVote) String() string {
	return fmt.Sprintf(`{"msgOption": %s, "voteOption": %s, "peerInfo": %s, "createAt": %d, "sign": %v}`,
		vote.MsgOption.String(), vote.VoteOption.String(), vote.PeerInfo.String(), vote.CreateAt, vote.Sign)
}

func ConvertPrepareVote(vote *PrepareVote) *pb.PrepareVote {
	return &pb.PrepareVote{
		MsgOption:  ConvertMsgOption(vote.MsgOption),
		VoteOption: vote.VoteOption.Bytes(),
		PeerInfo:   ConvertTaskPeerInfo(vote.PeerInfo),
		CreateAt:   vote.CreateAt,
		Sign:       vote.Sign,
	}
}
func FetchPrepareVote(vote *pb.PrepareVote) *PrepareVote {
	return &PrepareVote{
		MsgOption:  FetchMsgOption(vote.GetMsgOption()),
		VoteOption: VoteOptionFromBytes(vote.VoteOption),
		PeerInfo:   FetchTaskPeerInfo(vote.PeerInfo),
		CreateAt:   vote.CreateAt,
		Sign:       vote.Sign,
	}
}

type ConfirmMsg struct {
	MsgOption *MsgOption
	Peers     *pb.ConfirmTaskPeerInfo
	CreateAt  uint64
	Sign      []byte
}

func (msg *ConfirmMsg) String() string {
	return fmt.Sprintf(`{"msgOption": %s, "peerDesc": %s, "createAt": %d, "sign": %v}`,
		msg.MsgOption.String(), msg.Peers.String(), msg.CreateAt, msg.Sign)
}

func ConvertConfirmMsg(msg *ConfirmMsg) *pb.ConfirmMsg {
	return &pb.ConfirmMsg{
		MsgOption: ConvertMsgOption(msg.MsgOption),
		Peers:     msg.Peers,
		CreateAt:  msg.CreateAt,
		Sign:      msg.Sign,
	}
}
func FetchConfirmMsg(msg *pb.ConfirmMsg) *ConfirmMsg {
	return &ConfirmMsg{
		MsgOption:  FetchMsgOption(msg.GetMsgOption()),
		Peers: msg.Peers,
		CreateAt: msg.CreateAt,
		Sign:     msg.Sign,
	}
}

type ConfirmVote struct {
	MsgOption  *MsgOption
	VoteOption VoteOption
	CreateAt   uint64
	Sign       []byte
}

func (vote *ConfirmVote) String() string {
	return fmt.Sprintf(`{"msgOption": %s, "voteOption": %s, "createAt": %d, "sign": %v}`,
		vote.MsgOption.String(), vote.VoteOption.String(), vote.CreateAt, vote.Sign)
}

func ConvertConfirmVote(vote *ConfirmVote) *pb.ConfirmVote {
	return &pb.ConfirmVote{
		MsgOption: ConvertMsgOption(vote.MsgOption),
		VoteOption: vote.VoteOption.Bytes(),
		CreateAt:   vote.CreateAt,
		Sign:       vote.Sign,
	}
}
func FetchConfirmVote(vote *pb.ConfirmVote) *ConfirmVote {
	return &ConfirmVote{
		MsgOption:  FetchMsgOption(vote.GetMsgOption()),
		VoteOption: VoteOptionFromBytes(vote.VoteOption),
		CreateAt:   vote.CreateAt,
		Sign:       vote.Sign,
	}
}

type CommitMsg struct {
	MsgOption *MsgOption
	CreateAt  uint64
	Sign      []byte
}

func (msg *CommitMsg) String() string {
	return fmt.Sprintf(`{"msgOption": %s, "createAt": %d, "sign": %v}`,
		msg.MsgOption.String(), msg.CreateAt, msg.Sign)
}

func ConvertCommitMsg(msg *CommitMsg) *pb.CommitMsg {
	return &pb.CommitMsg{
		MsgOption: ConvertMsgOption(msg.MsgOption),
		CreateAt: msg.CreateAt,
		Sign:     msg.Sign,
	}
}
func FetchCommitMsg(msg *pb.CommitMsg) *CommitMsg {
	return &CommitMsg{
		MsgOption:  FetchMsgOption(msg.GetMsgOption()),
		CreateAt: msg.CreateAt,
		Sign:     msg.Sign,
	}
}

type TaskResultMsg struct {
	MsgOption     *MsgOption
	TaskEventList []*libTypes.TaskEvent
	CreateAt      uint64
	Sign          []byte
}

func (msg *TaskResultMsg) String() string {
	return fmt.Sprintf(`{"msgOption": %s, "createAt": %d, "sign": %v}`,
		msg.MsgOption.String(), msg.CreateAt, msg.Sign)
}

func ConvertTaskResultMsg(msg *TaskResultMsg) *pb.TaskResultMsg {
	return &pb.TaskResultMsg{
		MsgOption: ConvertMsgOption(msg.MsgOption),
		TaskEventList: ConvertTaskEventArr(msg.TaskEventList),
		CreateAt:      msg.CreateAt,
		Sign:          msg.Sign,
	}
}

func FetchTaskResultMsg(msg *pb.TaskResultMsg) *TaskResultMsg {
	return &TaskResultMsg{
		MsgOption:  FetchMsgOption(msg.GetMsgOption()),
		TaskEventList: FetchTaskEventArr(msg.TaskEventList),
		CreateAt:      msg.CreateAt,
		Sign:          msg.Sign,
	}
}
