package types

import (
	"encoding/json"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/bytesutil"
	"github.com/RosettaFlow/Carrier-Go/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/consensus/twopc/utils"
	"github.com/RosettaFlow/Carrier-Go/crypto/sha3"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/ethereum/go-ethereum/rlp"
	"sync/atomic"
)

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

type taskOption struct {
	Role                  twopc.TaskRole           `json:"role"` // The role information of the current recipient of the task
	TaskId                string                   `json:"taskId"`
	TaskName              string                   `json:"taskName"`
	Owner                 *types.NodeAlias         `json:"owner"`
	AlgoSupplier          *types.NodeAlias         `json:"algoSupplier"`
	DataSupplier          []*dataSupplierOption    `json:"dataSupplier"`
	PowerSupplier         []*powerSupplierOption   `json:"powerSupplier"`
	Receivers             []*receiverOption        `json:"receivers"`
	OperationCost         *types.TaskOperationCost `json:"operationCost"`
	CalculateContractCode string                   `json:"calculateContractCode"`
	DataSplitContractCode string                   `json:"dataSplitContractCode"`
	CreateAt              uint64                   `json:"createat"`
}

func (t *taskOption) Hash() common.Hash {
	return rlpHash(t)
}

type dataSupplierOption struct {
	MemberInfo      *types.NodeAlias `json:"memberInfo"`
	MetaDataId      string           `json:"metaDataId"`
	ColumnIndexList []uint64         `json:"columnIndexList"`
}
type powerSupplierOption struct {
	MemberInfo *types.NodeAlias `json:"memberInfo"`
}
type receiverOption struct {
	MemberInfo *types.NodeAlias   `json:"memberInfo"`
	Providers  []*types.NodeAlias `json:"providers"`
}

type taskPeerInfo struct {
	// Used to connect when running task, internal network resuorce of org.
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

type PrepareMsg struct {
	ProposalID  common.Hash   `json:"proposalId"`
	TaskOption  *taskOption   `json:"taskOption"`
	CreateAt    uint64        `json:"createAt"`
	Sign        twopc.MsgSign `json:"sign"`
	messageHash atomic.Value  `rlp:"-"`
}

func (msg *PrepareMsg) String() string {
	b, _ := json.Marshal(msg)
	return string(b)
}
func (msg *PrepareMsg) MsgHash() common.Hash {
	if mhash := msg.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(twopc.PrepareProposalMsg, utils.MergeBytes(msg.ProposalID.Bytes(),
		msg.TaskOption.Hash().Bytes(), msg.Sign.Bytes(), bytesutil.Uint64ToBytes(msg.CreateAt)))
	msg.messageHash.Store(v)
	return v
}

type PrepareVote struct {
	ProposalID  common.Hash      `json:"proposalId"`
	Role        twopc.TaskRole   `json:"role"` // The role information of the current recipient of the task
	Owner       *types.NodeAlias `json:"owner"`
	VoteOption  twopc.VoteOption `json:"voteOption"`
	PeerInfo    *taskPeerInfo    `json:"peerInfo"`
	CreateAt    uint64           `json:"createAt"`
	Sign        twopc.MsgSign    `json:"sign"`
	messageHash atomic.Value     `rlp:"-"`
}

func (msg *PrepareVote) String() string {
	b, _ := json.Marshal(msg)
	return string(b)
}

func (msg *PrepareVote) MsgHash() common.Hash {
	if mhash := msg.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(twopc.PrepareVoteMsg, utils.MergeBytes(msg.ProposalID.Bytes(), /*msg.VoteNodeID.Bytes(), */ // TODO 编码 NodeAlias
		[]byte{msg.VoteOption.Byte()}, msg.Sign.Bytes(), bytesutil.Uint64ToBytes(msg.CreateAt)))
	msg.messageHash.Store(v)
	return v
}
