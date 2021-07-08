package types

import "github.com/RosettaFlow/Carrier-Go/common"

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

type PrepareVote struct {
	ProposalId common.Hash
	TaskRole   TaskRole
	Owner      *NodeAlias
	VoteOption VoteOption
	PeerInfo   *PrepareVoteResource
	CreateAt   uint64
	Sign       []byte
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
