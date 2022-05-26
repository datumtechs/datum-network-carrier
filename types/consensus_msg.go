package types

import (
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common"
	libtypes "github.com/datumtechs/datum-network-carrier/lib/types"
	msgcommonpb "github.com/datumtechs/datum-network-carrier/lib/netmsg/common"
	twopcpb "github.com/datumtechs/datum-network-carrier/lib/netmsg/consensus/twopc"
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
	return fmt.Sprintf(`{"id": %s, "ip": %s, "port": %s, "partyId": %s}`, resource.GetId(), resource.GetIp(), resource.GetPort(), resource.GetPartyId())
}

func (resource *PrepareVoteResource) GetId() string      { return resource.Id }
func (resource *PrepareVoteResource) GetIp() string      { return resource.Ip }
func (resource *PrepareVoteResource) GetPort() string    { return resource.Port }
func (resource *PrepareVoteResource) GetPartyId() string { return resource.PartyId }

func ConvertTaskPeerInfo(peerInfo *PrepareVoteResource) *twopcpb.TaskPeerInfo {
	if nil == peerInfo {
		return &twopcpb.TaskPeerInfo{}
	}
	return &twopcpb.TaskPeerInfo{
		Ip:      []byte(peerInfo.GetIp()),
		Port:    []byte(peerInfo.GetPort()),
		PartyId: []byte(peerInfo.GetPartyId()),
	}
}
func FetchTaskPeerInfo(peerInfo *twopcpb.TaskPeerInfo) *PrepareVoteResource {
	if nil == peerInfo {
		return &PrepareVoteResource{}
	}
	return &PrepareVoteResource{
		Ip:      string(peerInfo.GetIp()),
		Port:    string(peerInfo.GetPort()),
		PartyId: string(peerInfo.GetPartyId()),
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
				Ip:      []byte(resource.GetIp()),
				Port:    []byte(resource.GetPort()),
				PartyId: []byte(resource.GetPartyId()),
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
				Ip:      string(peerInfo.GetIp()),
				Port:    string(peerInfo.GetPort()),
				PartyId: string(peerInfo.GetPartyId()),
			}
		}
		arr[i] = resource
	}
	return arr
}

type MsgOption struct {
	ProposalId      common.Hash
	SenderRole      libtypes.TaskRole
	SenderPartyId   string
	ReceiverRole    libtypes.TaskRole
	ReceiverPartyId string
	Owner           *libtypes.TaskOrganization
}

func MakeMsgOption(proposalId common.Hash,
	senderRole, receiverRole libtypes.TaskRole,
	senderPartyId, receiverPartyId string,
	sender *libtypes.TaskOrganization,
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

func (option *MsgOption) String() string {
	return fmt.Sprintf(`{"ProposalId": "%s", "senderRole": "%s", "senderPartyId": "%s", "receiverRole": "%s", "receiverPartyId": "%s", "owner": %s}`,
		option.GetProposalId().String(), option.GetSenderRole().String(), option.GetSenderPartyId(), option.GetReceiverRole().String(), option.GetReceiverPartyId(), option.GetOwner().String())
}

func (option *MsgOption) GetProposalId() common.Hash              { return option.ProposalId }
func (option *MsgOption) GetSenderRole() libtypes.TaskRole     { return option.SenderRole }
func (option *MsgOption) GetSenderPartyId() string                { return option.SenderPartyId }
func (option *MsgOption) GetReceiverRole() libtypes.TaskRole   { return option.ReceiverRole }
func (option *MsgOption) GetReceiverPartyId() string              { return option.ReceiverPartyId }
func (option *MsgOption) GetOwner() *libtypes.TaskOrganization { return option.Owner }

func ConvertMsgOption(option *MsgOption) *msgcommonpb.MsgOption {
	return &msgcommonpb.MsgOption{
		ProposalId:      option.GetProposalId().Bytes(),
		SenderRole:      uint64(option.GetSenderRole()),
		SenderPartyId:   []byte(option.GetSenderPartyId()),
		ReceiverRole:    uint64(option.GetReceiverRole()),
		ReceiverPartyId: []byte(option.GetReceiverPartyId()),
		MsgOwner: &msgcommonpb.TaskOrganizationIdentityInfo{
			Name:       []byte(option.GetOwner().GetNodeName()),
			NodeId:     []byte(option.GetOwner().GetNodeId()),
			IdentityId: []byte(option.GetOwner().GetIdentityId()),
			PartyId:    []byte(option.GetOwner().GetPartyId()),
		},
	}
}

func FetchMsgOption(option *msgcommonpb.MsgOption) *MsgOption {
	return &MsgOption{
		ProposalId:      common.BytesToHash(option.GetProposalId()),
		SenderRole:      libtypes.TaskRole(option.GetSenderRole()),
		SenderPartyId:   string(option.GetSenderPartyId()),
		ReceiverRole:    libtypes.TaskRole(option.GetReceiverRole()),
		ReceiverPartyId: string(option.GetReceiverPartyId()),
		Owner: &libtypes.TaskOrganization{
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
	Evidence  string
	CreateAt  uint64
	Sign      []byte
}

func (msg *PrepareMsg) String() string {
	return fmt.Sprintf(`{"msgOption": %s, "evidence": %s, "createAt": %d, "sign": %v}`,
		msg.GetMsgOption().String(), msg.GetEvidence(), msg.GetCreateAt(), msg.GetSign())
}
func (msg *PrepareMsg) StringWithTask() string {
	return fmt.Sprintf(`{"msgOption": %s, "taskInfo": %s, "evidence": %s, "createAt": %d, "sign": %v}`,
		msg.GetMsgOption().String(), msg.GetTask().GetTaskData().String(), msg.GetEvidence(), msg.GetCreateAt(), msg.GetSign())
}

func (msg *PrepareMsg) GetMsgOption() *MsgOption { return msg.MsgOption }
func (msg *PrepareMsg) GetTask() *Task           { return msg.TaskInfo }
func (msg *PrepareMsg) GetEvidence() string { return msg.Evidence }
func (msg *PrepareMsg) GetCreateAt() uint64 { return msg.CreateAt }
func (msg *PrepareMsg) GetSign() []byte     { return msg.Sign }

//type PrepareMsgExtra struct {
//	Nonce   []byte
//	Weights [][]byte
//}
//func (msg *PrepareMsgExtra) String() string {
//	nonceStr := "0x"
//	if len(msg.Nonce) != 0 {
//		nonceStr = common.BytesToHash(msg.Nonce).Hex()
//	}
//	weightsStr := "[]"
//	if len(msg.Weights) != 0 {
//		arr := make([]string, len(msg.Weights))
//		for i, weight := range msg.Weights {
//			arr[i] = new(big.Int).SetBytes(weight).String()
//		}
//		weightsStr = "[" + strings.Join(arr, ",") + "]"
//	}
//	return fmt.Sprintf(`{"nonce": %s, "Weights": %s}`,
//		nonceStr, weightsStr)
//}
//func (msg *PrepareMsgExtra) GetNonce() []byte { return msg.Nonce }
//func (msg *PrepareMsgExtra) GetWeights() [][]byte { return msg.Weights }
//func (msg *PrepareMsgExtra) GetNonceHex() string {
//	nonceStr := "0x"
//	if len(msg.Nonce) != 0 {
//		nonceStr = common.BytesToHash(msg.Nonce).Hex()
//	}
//	return nonceStr
//}
//
//func (msg *PrepareMsgExtra) GetWeightsBigInt() []*big.Int {
//	var weights []*big.Int
//	if len(msg.Weights) != 0 {
//		weights = make([]*big.Int, len(msg.Weights))
//		for i, weight := range msg.Weights {
//			weights[i] = new(big.Int).SetBytes(weight)
//		}
//	}
//	return weights
//}

type PrepareVote struct {
	MsgOption  *MsgOption
	VoteOption VoteOption
	PeerInfo   *PrepareVoteResource
	CreateAt   uint64
	Sign       []byte
}

func (vote *PrepareVote) PeerInfoEmpty() bool { return nil == vote.GetPeerInfo() }
func (vote *PrepareVote) String() string {
	return fmt.Sprintf(`{"msgOption": %s, "voteOption": "%s", "peerInfo": %s, "createAt": %d, "sign": %v}`,
		vote.GetMsgOption().String(), vote.GetVoteOption().String(), vote.GetPeerInfo().String(), vote.GetCreateAt(), vote.GetSign())
}

func (vote *PrepareVote) GetMsgOption() *MsgOption          { return vote.MsgOption }
func (vote *PrepareVote) GetVoteOption() VoteOption         { return vote.VoteOption }
func (vote *PrepareVote) GetPeerInfo() *PrepareVoteResource { return vote.PeerInfo }
func (vote *PrepareVote) GetCreateAt() uint64               { return vote.CreateAt }
func (vote *PrepareVote) GetSign() []byte                   { return vote.Sign }

func ConvertPrepareVote(vote *PrepareVote) *twopcpb.PrepareVote {
	return &twopcpb.PrepareVote{
		MsgOption:  ConvertMsgOption(vote.GetMsgOption()),
		VoteOption: vote.GetVoteOption().Bytes(),
		PeerInfo:   ConvertTaskPeerInfo(vote.GetPeerInfo()),
		CreateAt:   vote.GetCreateAt(),
		Sign:       vote.GetSign(),
	}
}
func FetchPrepareVote(vote *twopcpb.PrepareVote) *PrepareVote {
	return &PrepareVote{
		MsgOption:  FetchMsgOption(vote.GetMsgOption()),
		VoteOption: VoteOptionFromBytes(vote.GetVoteOption()),
		PeerInfo:   FetchTaskPeerInfo(vote.GetPeerInfo()),
		CreateAt:   vote.GetCreateAt(),
		Sign:       vote.GetSign(),
	}
}

type ConfirmMsg struct {
	MsgOption     *MsgOption
	ConfirmOption TwopcMsgOption
	Peers         *twopcpb.ConfirmTaskPeerInfo
	CreateAt      uint64
	Sign          []byte
}

func (msg *ConfirmMsg) PeersEmpty() bool {
	if nil == msg.GetPeers() {
		return true
	}
	if len(msg.GetPeers().GetDataSupplierPeerInfos()) == 0 &&
		len(msg.GetPeers().GetPowerSupplierPeerInfos()) == 0 && len(msg.GetPeers().GetResultReceiverPeerInfos()) == 0 {
		return true
	}
	return false
}

func (msg *ConfirmMsg) String() string {

	var peers string
	if msg.PeersEmpty() {
		peers = "{}"
	} else {
		peers = msg.GetPeers().String()
	}

	return fmt.Sprintf(`{"msgOption": %s, "confirmOption": "%s", "peers": %s, "createAt": %d, "sign": %v}`,
		msg.GetMsgOption().String(), msg.GetConfirmOption().String(), peers, msg.GetCreateAt(), msg.GetSign())
}

func (msg *ConfirmMsg) GetMsgOption() *MsgOption               { return msg.MsgOption }
func (msg *ConfirmMsg) GetConfirmOption() TwopcMsgOption       { return msg.ConfirmOption }
func (msg *ConfirmMsg) GetPeers() *twopcpb.ConfirmTaskPeerInfo { return msg.Peers }
func (msg *ConfirmMsg) GetCreateAt() uint64                    { return msg.CreateAt }
func (msg *ConfirmMsg) GetSign() []byte                        { return msg.Sign }

func ConvertConfirmMsg(msg *ConfirmMsg) *twopcpb.ConfirmMsg {
	return &twopcpb.ConfirmMsg{
		MsgOption:     ConvertMsgOption(msg.GetMsgOption()),
		ConfirmOption: msg.GetConfirmOption().Bytes(),
		Peers:         msg.GetPeers(),
		CreateAt:      msg.GetCreateAt(),
		Sign:          msg.GetSign(),
	}
}
func FetchConfirmMsg(msg *twopcpb.ConfirmMsg) *ConfirmMsg {
	return &ConfirmMsg{
		MsgOption:     FetchMsgOption(msg.GetMsgOption()),
		ConfirmOption: TwopcMsgOptionFromBytes(msg.GetConfirmOption()),
		Peers:         msg.GetPeers(),
		CreateAt:      msg.GetCreateAt(),
		Sign:          msg.GetSign(),
	}
}

type ConfirmVote struct {
	MsgOption  *MsgOption
	VoteOption VoteOption
	CreateAt   uint64
	Sign       []byte
}

func (vote *ConfirmVote) String() string {
	return fmt.Sprintf(`{"msgOption": %s, "voteOption": "%s", "createAt": %d, "sign": %v}`,
		vote.GetMsgOption().String(), vote.GetVoteOption().String(), vote.GetCreateAt(), vote.GetSign())
}

func (vote *ConfirmVote) GetMsgOption() *MsgOption  { return vote.MsgOption }
func (vote *ConfirmVote) GetVoteOption() VoteOption { return vote.VoteOption }
func (vote *ConfirmVote) GetCreateAt() uint64       { return vote.CreateAt }
func (vote *ConfirmVote) GetSign() []byte           { return vote.Sign }

func ConvertConfirmVote(vote *ConfirmVote) *twopcpb.ConfirmVote {
	return &twopcpb.ConfirmVote{
		MsgOption:  ConvertMsgOption(vote.GetMsgOption()),
		VoteOption: vote.GetVoteOption().Bytes(),
		CreateAt:   vote.GetCreateAt(),
		Sign:       vote.GetSign(),
	}
}
func FetchConfirmVote(vote *twopcpb.ConfirmVote) *ConfirmVote {
	return &ConfirmVote{
		MsgOption:  FetchMsgOption(vote.GetMsgOption()),
		VoteOption: VoteOptionFromBytes(vote.GetVoteOption()),
		CreateAt:   vote.GetCreateAt(),
		Sign:       vote.GetSign(),
	}
}

type CommitMsg struct {
	MsgOption    *MsgOption
	CommitOption TwopcMsgOption
	CreateAt     uint64
	Sign         []byte
}

func (msg *CommitMsg) String() string {
	return fmt.Sprintf(`{"msgOption": %s, "commitOption", "%s", "createAt": %d, "sign": %v}`,
		msg.GetMsgOption().String(), msg.GetCommitOption().String(), msg.GetCreateAt(), msg.GetSign())
}

func (msg *CommitMsg) GetMsgOption() *MsgOption        { return msg.MsgOption }
func (msg *CommitMsg) GetCommitOption() TwopcMsgOption { return msg.CommitOption }
func (msg *CommitMsg) GetCreateAt() uint64             { return msg.CreateAt }
func (msg *CommitMsg) GetSign() []byte                 { return msg.Sign }

func ConvertCommitMsg(msg *CommitMsg) *twopcpb.CommitMsg {
	return &twopcpb.CommitMsg{
		MsgOption:    ConvertMsgOption(msg.GetMsgOption()),
		CommitOption: msg.GetCommitOption().Bytes(),
		CreateAt:     msg.GetCreateAt(),
		Sign:         msg.GetSign(),
	}
}
func FetchCommitMsg(msg *twopcpb.CommitMsg) *CommitMsg {
	return &CommitMsg{
		MsgOption:    FetchMsgOption(msg.GetMsgOption()),
		CommitOption: TwopcMsgOptionFromBytes(msg.GetCommitOption()),
		CreateAt:     msg.GetCreateAt(),
		Sign:         msg.GetSign(),
	}
}
