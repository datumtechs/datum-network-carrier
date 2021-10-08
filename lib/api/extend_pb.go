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

func (seed *SeedPeer) GenSeedNodeId() string {
	if "" != seed.GetId() {
		return seed.GetId()
	}
	seed.Id = PrefixSeedNodeId + seed.hashByCreateTime().Hex()
	return seed.GetId()
}

func (seed *SeedPeer) hash() (h common.Hash) {
	hw := sha3.NewKeccak256()
	d := &struct {
		InternalIp   string
		InternalPort string
	}{
		InternalIp:   seed.GetInternalIp(),
		InternalPort: seed.GetInternalPort(),
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
		InternalIp:   seed.GetInternalIp(),
		InternalPort: seed.GetInternalPort(),
		CreateTime:   timeutils.UnixMsecUint64(),
	}
	rlp.Encode(hw, d)
	hw.Sum(h[:0])
	return h
}

func (node *YarnRegisteredPeerDetail) GenJobNodeId() string {
	if "" != node.GetId() {
		return node.GetId()
	}
	node.Id = PrefixJobNodeId + node.hashByCreateTime(PrefixTypeJobNode).Hex()
	return node.GetId()
}

func (node *YarnRegisteredPeerDetail) GenDataNodeId() string {
	if "" != node.GetId() {
		return node.GetId()
	}
	node.Id = PrefixDataNodeId + node.hashByCreateTime(PrefixTypeDataNode).Hex()
	return node.GetId()
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
		InternalIp:   node.GetInternalIp(),
		InternalPort: node.GetInternalPort(),
		ExternalIp:   node.GetExternalIp(),
		ExternalPort: node.GetExternalPort(),
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
		InternalIp:   node.GetInternalIp(),
		InternalPort: node.GetInternalPort(),
		ExternalIp:   node.GetExternalIp(),
		ExternalPort: node.GetExternalPort(),
		CreateAt:     timeutils.UnixMsecUint64(),
	}
	rlp.Encode(hw, d)
	hw.Sum(h[:0])
	return h
}