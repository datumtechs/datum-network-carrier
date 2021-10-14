package types

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	msgcommonpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/common"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
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
func ConvertTaskPeerInfo(peerInfo *PrepareVoteResource) *twopcpb.TaskPeerInfo {
	if nil == peerInfo {
		return &twopcpb.TaskPeerInfo{}
	}
	return &twopcpb.TaskPeerInfo{
		Ip:      []byte(peerInfo.Ip),
		Port:    []byte(peerInfo.Port),
		PartyId: []byte(peerInfo.PartyId),
	}
}
func FetchTaskPeerInfo(peerInfo *twopcpb.TaskPeerInfo) *PrepareVoteResource {
	if nil == peerInfo {
		return &PrepareVoteResource{}
	}
	return &PrepareVoteResource{
		Ip:      string(peerInfo.Ip),
		Port:    string(peerInfo.Port),
		PartyId: string(peerInfo.PartyId),
	}
}
func ConvertTaskPeerInfoArr(resourceArr []*PrepareVoteResource) []*twopcpb.TaskPeerInfo {

	arr := make([]*twopcpb.TaskPeerInfo, len(resourceArr))

	for i, resource := range resourceArr {
		var peerInfo *twopcpb.TaskPeerInfo
		if nil == resource {
			peerInfo = &twopcpb.TaskPeerInfo{}
		} else {
			peerInfo = &twopcpb.TaskPeerInfo{
				Ip:      []byte(resource.Ip),
				Port:    []byte(resource.Port),
				PartyId: []byte(resource.PartyId),
			}
		}

		arr[i] = peerInfo
	}
	return arr
}
func FetchTaskPeerInfoArr(peerInfoArr []*twopcpb.TaskPeerInfo) []*PrepareVoteResource {

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
	SenderRole      apicommonpb.TaskRole
	SenderPartyId   string
	ReceiverRole    apicommonpb.TaskRole
	ReceiverPartyId string
	Owner           *apicommonpb.TaskOrganization
}

func (option *MsgOption) String() string {
	return fmt.Sprintf(`{"ProposalId": "%s", "senderRole": "%s", "senderPartyId": "%s", "receiverRole": "%s", "receiverPartyId": "%s", "owner": %s}`,
		option.ProposalId.String(), option.SenderRole.String(), option.SenderPartyId, option.ReceiverRole.String(), option.ReceiverPartyId, option.Owner.String())
}

func ConvertMsgOption(option *MsgOption) *msgcommonpb.MsgOption {
	return &msgcommonpb.MsgOption{
		ProposalId:      option.ProposalId.Bytes(),
		SenderRole:      uint64(option.SenderRole),
		SenderPartyId:   []byte(option.SenderPartyId),
		ReceiverRole:    uint64(option.ReceiverRole),
		ReceiverPartyId: []byte(option.ReceiverPartyId),
		MsgOwner: &msgcommonpb.TaskOrganizationIdentityInfo{
			Name:       []byte(option.Owner.GetNodeName()),
			NodeId:     []byte(option.Owner.GetNodeId()),
			IdentityId: []byte(option.Owner.GetIdentityId()),
			PartyId:    []byte(option.Owner.GetPartyId()),
		},
	}
}

func FetchMsgOption(option *msgcommonpb.MsgOption) *MsgOption {
	return &MsgOption{
		ProposalId:      common.BytesToHash(option.ProposalId),
		SenderRole:      apicommonpb.TaskRole(option.GetSenderRole()),
		SenderPartyId:   string(option.SenderPartyId),
		ReceiverRole:    apicommonpb.TaskRole(option.GetReceiverRole()),
		ReceiverPartyId: string(option.ReceiverPartyId),
		Owner: &apicommonpb.TaskOrganization{
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

func (vote *PrepareVote) PeerInfoEmpty () bool { return nil == vote.PeerInfo }
func (vote *PrepareVote) String() string {
	return fmt.Sprintf(`{"msgOption": %s, "voteOption": "%s", "peerInfo": %s, "createAt": %d, "sign": %v}`,
		vote.MsgOption.String(), vote.VoteOption.String(), vote.PeerInfo.String(), vote.CreateAt, vote.Sign)
}

func ConvertPrepareVote(vote *PrepareVote) *twopcpb.PrepareVote {
	return &twopcpb.PrepareVote{
		MsgOption:  ConvertMsgOption(vote.MsgOption),
		VoteOption: vote.VoteOption.Bytes(),
		PeerInfo:   ConvertTaskPeerInfo(vote.PeerInfo),
		CreateAt:   vote.CreateAt,
		Sign:       vote.Sign,
	}
}
func FetchPrepareVote(vote *twopcpb.PrepareVote) *PrepareVote {
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
	Peers     *twopcpb.ConfirmTaskPeerInfo
	CreateAt  uint64
	Sign      []byte
}

func (msg *ConfirmMsg) PeersEmpty() bool {
	if nil == msg.Peers {
		return true
	}
	if nil == msg.Peers.GetOwnerPeerInfo() && len(msg.Peers.GetDataSupplierPeerInfoList()) == 0 &&
		len(msg.Peers.GetPowerSupplierPeerInfoList()) == 0 && len(msg.Peers.GetResultReceiverPeerInfoList()) == 0 {
		return true
	}
	return false
}

func (msg *ConfirmMsg) String() string {

	var peers string
	if msg.PeersEmpty() {
		peers = "{}"
	} else {
		peers = msg.Peers.String()
	}

	return fmt.Sprintf(`{"msgOption": %s, "peers": %s, "createAt": %d, "sign": %v}`,
		msg.MsgOption.String(), peers, msg.CreateAt, msg.Sign)
}

func ConvertConfirmMsg(msg *ConfirmMsg) *twopcpb.ConfirmMsg {
	return &twopcpb.ConfirmMsg{
		MsgOption: ConvertMsgOption(msg.MsgOption),
		Peers:     msg.Peers,
		CreateAt:  msg.CreateAt,
		Sign:      msg.Sign,
	}
}
func FetchConfirmMsg(msg *twopcpb.ConfirmMsg) *ConfirmMsg {
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

func ConvertConfirmVote(vote *ConfirmVote) *twopcpb.ConfirmVote {
	return &twopcpb.ConfirmVote{
		MsgOption: ConvertMsgOption(vote.MsgOption),
		VoteOption: vote.VoteOption.Bytes(),
		CreateAt:   vote.CreateAt,
		Sign:       vote.Sign,
	}
}
func FetchConfirmVote(vote *twopcpb.ConfirmVote) *ConfirmVote {
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

func ConvertCommitMsg(msg *CommitMsg) *twopcpb.CommitMsg {
	return &twopcpb.CommitMsg{
		MsgOption: ConvertMsgOption(msg.MsgOption),
		CreateAt: msg.CreateAt,
		Sign:     msg.Sign,
	}
}
func FetchCommitMsg(msg *twopcpb.CommitMsg) *CommitMsg {
	return &CommitMsg{
		MsgOption:  FetchMsgOption(msg.GetMsgOption()),
		CreateAt: msg.CreateAt,
		Sign:     msg.Sign,
	}
}

