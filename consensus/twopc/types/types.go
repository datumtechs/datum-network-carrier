package types

import (
	"encoding/json"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/consensus/twopc/utils"
	"github.com/RosettaFlow/Carrier-Go/crypto/sha3"
	"github.com/RosettaFlow/Carrier-Go/p2p"
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
	Role          twopc.TaskRole           `json:"role"`
	TaskId        string                   `json:"taskId"`
	TaskName      string                   `json:"taskName"`
	Owner         *types.NodeAlias         `json:"owner"`
	AlgoSupplier  *types.NodeAlias         `json:"algoSupplier"`
	DataSupplier  []*dataSupplierOption    `json:"dataSupplier"`
	PowerSupplier []*powerSupplierOption   `json:"powerSupplier"`
	Receivers     []*receiverOption        `json:"receivers"`
	CreateAt      uint64                   `json:"createat"`
	OperationCost *types.TaskOperationCost `json:"operationCost"`
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
	*types.NodeAlias
	Providers []*types.NodeAlias `json:"providers"`
}

type PrepareMsg struct {
	ProposalID  common.Hash  `json:"proposalID"`
	TaskOption  *taskOption  `json:"taskOption"`
	messageHash atomic.Value `rlp:"-"`
}

func (msg *PrepareMsg) String() string {
	b, _ := json.Marshal(msg)
	return string(b)
}
func (msg *PrepareMsg) MsgHash() common.Hash {
	if mhash := msg.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(twopc.PrepareProposalMsg, utils.MergeBytes(msg.ProposalID.Bytes(), msg.TaskOption.Hash().Bytes()))
	msg.messageHash.Store(v)
	return v
}

type PrepareVote struct {
	ProposalID  common.Hash      `json:"proposalID"`
	VoteNodeID  p2p.NodeID       `json:"voteNodeID"`
	VoteOption  twopc.VoteOption `json:"voteOption"`
	VoteSign    twopc.MsgSign    `json:"voteSign"`
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
	v := utils.BuildHash(twopc.PrepareVoteMsg, utils.MergeBytes(msg.ProposalID.Bytes(), msg.VoteNodeID.Bytes(),
		[]byte{msg.VoteOption.Byte()}, msg.VoteSign.Bytes()))
	msg.messageHash.Store(v)
	return v
}
