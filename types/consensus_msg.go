package types

import "github.com/RosettaFlow/Carrier-Go/common"

type ConsensusEngineType string
func (t ConsensusEngineType) String() string {return string(t)}
const (
	ChainconsTyp ConsensusEngineType = "ChainconsType"
	TwopcTyp ConsensusEngineType = "TwopcType"
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

