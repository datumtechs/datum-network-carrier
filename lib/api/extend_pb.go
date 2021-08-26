package api

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	PrefixTypeYarnNode RegisteredNodeType = "yarnNode"
	PrefixTypeJobNode  RegisteredNodeType = "jobNode"
	PrefixTypeDataNode RegisteredNodeType = "dataNode"

	PrefixSeedNodeId = "seed:"
	PrefixJobNodeId  = "jobNode:"
	PrefixDataNodeId = "dataNode:"
)

type RegisteredNodeType string

func (typ RegisteredNodeType) String() string { return string(typ) }

func (seed *SeedPeer) SeedNodeId() string {
	if "" != seed.Id {
		return seed.Id
	}
	seed.Id = PrefixSeedNodeId + seed.hashByCreateTime().Hex()
	return seed.Id
}

func (seed *SeedPeer) hash() (h common.Hash) {
	hw := sha3.NewKeccak256()
	d := &struct {
		InternalIp   string
		InternalPort string
	}{
		InternalIp:   seed.InternalIp,
		InternalPort: seed.InternalPort,
	}
	rlp.Encode(hw, d)
	hw.Sum(h[:0])
	return h
}

func (seed *SeedPeer) hashByCreateTime() (h common.Hash) {
	hw := sha3.NewKeccak256()
	d := &struct {
		InternalIp   string
		InternalPort string
		CreateTime   uint64
	}{
		InternalIp:   seed.InternalIp,
		InternalPort: seed.InternalPort,
		CreateTime:   uint64(timeutils.UnixMsec()),
	}
	rlp.Encode(hw, d)
	hw.Sum(h[:0])
	return h
}

func (node *YarnRegisteredPeerDetail) SetJobNodeId() string {
	if "" != node.Id {
		return node.Id
	}
	node.Id = PrefixJobNodeId + node.hashByCreateTime(PrefixTypeJobNode).Hex()
	return node.Id
}

func (node *YarnRegisteredPeerDetail) SetDataNodeId() string {
	if "" != node.Id {
		return node.Id
	}
	node.Id = PrefixDataNodeId + node.hashByCreateTime(PrefixTypeDataNode).Hex()
	return node.Id
}

func (node *YarnRegisteredPeerDetail) hash(typ RegisteredNodeType) (h common.Hash) {
	hw := sha3.NewKeccak256()
	d := &struct {
		Type         RegisteredNodeType
		InternalIp   string
		InternalPort string
		ExternalIp   string
		ExternalPort string
	}{
		Type:         typ,
		InternalIp:   node.InternalIp,
		InternalPort: node.InternalPort,
		ExternalIp:   node.ExternalIp,
		ExternalPort: node.ExternalPort,
	}
	rlp.Encode(hw, d)
	hw.Sum(h[:0])
	return h
}

func (node *YarnRegisteredPeerDetail) hashByCreateTime(typ RegisteredNodeType) (h common.Hash) {
	hw := sha3.NewKeccak256()
	d := &struct {
		Type         RegisteredNodeType
		InternalIp   string
		InternalPort string
		ExternalIp   string
		ExternalPort string
		CreateAt     uint64
	}{
		Type:         typ,
		InternalIp:   node.InternalIp,
		InternalPort: node.InternalPort,
		ExternalIp:   node.ExternalIp,
		ExternalPort: node.ExternalPort,
		CreateAt:     uint64(timeutils.UnixMsec()),
	}
	rlp.Encode(hw, d)
	hw.Sum(h[:0])
	return h
}