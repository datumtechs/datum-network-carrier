package types

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)

type NodeConnStatus int32
type RegisteredNodeType string

func (status NodeConnStatus) Int32() int32    { return int32(status) }
func (typ RegisteredNodeType) String() string { return string(typ) }

const (
	CONNECTED    NodeConnStatus = 0 // 连接上就是未启用算力
	NONCONNECTED NodeConnStatus = -1
	ENABLE_POWER NodeConnStatus = 1 // 启用算力
	BUSY_POWER   NodeConnStatus = 2 // 算力被占用(有任务在执行 ...)

)

const (
	PREFIX_TYPE_YARNNODE RegisteredNodeType = "yarnNode"
	PREFIX_TYPE_JOBNODE  RegisteredNodeType = "jobNode"
	PREFIX_TYPE_DATANODE RegisteredNodeType = "dataNode"

	PREFIX_SEEDNODE_ID = "seed:"
	PREFIX_JOBNODE_ID  = "jobNode:"
	PREFIX_DATANODE_ID = "dataNode:"
)

type SeedNodeInfo struct {
	Id           string         `json:"id"`
	InternalIp   string         `json:"internalIp"`
	InternalPort string         `json:"internalPort"`
	ConnState    NodeConnStatus `json:"connState"`
}

type RegisteredNodeInfo struct {
	Id           string         `json:"id"`
	InternalIp   string         `json:"internalIp"`
	InternalPort string         `json:"internalPort"`
	ExternalIp   string         `json:"externalIp"`
	ExternalPort string         `json:"externalPort"`
	ConnState    NodeConnStatus `json:"connState"`
}
func (n *RegisteredNodeInfo) String() string {
	return fmt.Sprintf(`{"id": %s, "internalIp": %s, "internalPort": %s, "externalIp": %s, "externalPort": %s, "connState": %d}`,
		n.Id, n.InternalIp, n.InternalPort, n.ExternalIp, n.ExternalPort, n.ConnState.Int32())
}
type RegisteredNodeDetail struct {
	NodeType string `json:"nodeType"`
	*RegisteredNodeInfo
}

func (seed *SeedNodeInfo) SeedNodeId() string {
	if "" != seed.Id {
		return seed.Id
	}
	seed.Id = PREFIX_SEEDNODE_ID + seed.hashByCreateTime().Hex()
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
func (seed *SeedNodeInfo) hashByCreateTime() (h common.Hash) {
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

func (node *RegisteredNodeInfo) SetJobNodeId() string {
	if "" != node.Id {
		return node.Id
	}
	node.Id = PREFIX_JOBNODE_ID + node.hashByCreateTime(PREFIX_TYPE_JOBNODE).Hex()
	return node.Id
}
func (node *RegisteredNodeInfo) SetDataNodeId() string {
	if "" != node.Id {
		return node.Id
	}
	node.Id = PREFIX_DATANODE_ID + node.hashByCreateTime(PREFIX_TYPE_DATANODE).Hex()
	return node.Id
}
func (node *RegisteredNodeInfo) hash(typ RegisteredNodeType) (h common.Hash) {
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

func (node *RegisteredNodeInfo) hashByCreateTime(typ RegisteredNodeType) (h common.Hash) {
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
