package types

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)

type RegisterNodeConnectStatus int32
type RegisteredNodeType string

const (
	CONNECTED             RegisterNodeConnectStatus = 0
	NONCONNECTED          RegisterNodeConnectStatus = -1
	ENABLECOMPUTERESOURCE RegisterNodeConnectStatus = 1
)

const (
	PREFIX_TYPE_JOBNODE RegisteredNodeType = "jobNode"
	PREFIX_TYPE_DATANODE RegisteredNodeType = "dataNode"


	PREFIX_SEEDNODE_ID = "seed:"
	PREFIX_JOBNODE_ID = "jobNode:"
	PREFIX_DATANODE_ID = "dataNode:"
)

type SeedNodeInfo struct {
	Id           string                    `json:"id"`
	InternalIp   string                    `json:"internalIp"`
	InternalPort string                    `json:"internalPort"`
	ConnState    RegisterNodeConnectStatus `json:"connState"`
}

type RegisteredNodeInfo struct {
	Id           string                    `json:"id"`
	InternalIp   string                    `json:"internalIp"`
	InternalPort string                    `json:"internalPort"`
	ExternalIp   string                    `json:"externalIp"`
	ExternalPort string                    `json:"externalPort"`
	ConnState    RegisterNodeConnectStatus `json:"connState"`
}
func (seed *SeedNodeInfo) SeedNodeId() string {
	if "" != seed.Id {
		return seed.Id
	}
	seed.Id = PREFIX_SEEDNODE_ID + seed.hash().Hex()
	return seed.Id
}
func (seed *SeedNodeInfo) hash() (h common.Hash) {
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
func (node *RegisteredNodeInfo) JobNodeId() string {
	if "" != node.Id {
		return node.Id
	}
	node.Id = PREFIX_JOBNODE_ID + node.hash(PREFIX_TYPE_JOBNODE).Hex()
	return node.Id
}
func (node *RegisteredNodeInfo) DataNodeId() string {
	if "" != node.Id {
		return node.Id
	}
	node.Id = PREFIX_DATANODE_ID + node.hash(PREFIX_TYPE_DATANODE).Hex()
	return node.Id
}
func (node *RegisteredNodeInfo) hash(typ RegisteredNodeType) (h common.Hash) {
	hw := sha3.NewKeccak256()
	d := &struct {
		Type 		RegisteredNodeType
		InternalIp   string
		InternalPort string
		ExternalIp   string
		ExternalPort string
	}{
		Type: typ,
		InternalIp:   node.InternalIp,
		InternalPort: node.InternalPort,
		ExternalIp: node.ExternalIp,
		ExternalPort: node.ExternalPort,
	}
	rlp.Encode(hw, d)
	hw.Sum(h[:0])
	return h
}