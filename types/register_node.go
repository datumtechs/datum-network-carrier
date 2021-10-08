package types

//type NodeConnStatus int32
////type RegisteredNodeType string
//
//func (status NodeConnStatus) Int32() int32    { return int32(status) }
//
//const (
//	Connected    NodeConnStatus = 0 // 连接上就是未启用算力
//	NonConnected NodeConnStatus = -1
//	EnablePower  NodeConnStatus = 1 // 启用算力
//	BusyPower    NodeConnStatus = 2 // 算力被占用(有任务在执行 ...)
//)


//type SeedNodeInfo struct {
//	Id           string         `json:"id"`
//	InternalIp   string         `json:"internalIp"`
//	InternalPort string         `json:"internalPort"`
//	ConnState    NodeConnStatus `json:"connState"`
//}

//type RegisteredNodeInfo struct {
//	Id           string         `json:"id"`
//	InternalIp   string         `json:"internalIp"`
//	InternalPort string         `json:"internalPort"`
//	ExternalIp   string         `json:"externalIp"`
//	ExternalPort string         `json:"externalPort"`
//	ConnState    NodeConnStatus `json:"connState"`
//}
//func (n *RegisteredNodeInfo) String() string {
//	return fmt.Sprintf(`{"id": %s, "internalIp": %s, "internalPort": %s, "externalIp": %s, "externalPort": %s, "connState": %d}`,
//		n.Id, n.InternalIp, n.InternalPort, n.ExternalIp, n.ExternalPort, n.ConnState.Int32())
//}
//type RegisteredNodeDetail struct {
//	NodeType string `json:"nodeType"`
//	*RegisteredNodeInfo
//}

//func (seed *SeedNodeInfo) GenSeedNodeId() string {
//	if "" != seed.Id {
//		return seed.Id
//	}
//	seed.Id = PREFIX_SEEDNODE_ID + seed.hashByCreateTime().Hex()
//	return seed.Id
//}
//
//func (seed *SeedNodeInfo) hash() (h common.Hash) {
//	hw := sha3.NewKeccak256()
//	d := &struct {
//		InternalIp   string
//		InternalPort string
//	}{
//		InternalIp:   seed.InternalIp,
//		InternalPort: seed.InternalPort,
//	}
//	rlp.Encode(hw, d)
//	hw.Sum(h[:0])
//	return h
//}
//func (seed *SeedNodeInfo) hashByCreateTime() (h common.Hash) {
//	hw := sha3.NewKeccak256()
//	d := &struct {
//		InternalIp   string
//		InternalPort string
//		CreateTime   uint64
//	}{
//		InternalIp:   seed.InternalIp,
//		InternalPort: seed.InternalPort,
//		CreateTime:   uint64(timeutils.UnixMsec()),
//	}
//	rlp.Encode(hw, d)
//	hw.Sum(h[:0])
//	return h
//}

//func (node *RegisteredNodeInfo) GenJobNodeId() string {
//	if "" != node.Id {
//		return node.Id
//	}
//	node.Id = PREFIX_JOBNODE_ID + node.hashByCreateTime(PREFIX_TYPE_JOBNODE).Hex()
//	return node.Id
//}
//func (node *RegisteredNodeInfo) GenDataNodeId() string {
//	if "" != node.Id {
//		return node.Id
//	}
//	node.Id = PREFIX_DATANODE_ID + node.hashByCreateTime(PREFIX_TYPE_DATANODE).Hex()
//	return node.Id
//}
//func (node *RegisteredNodeInfo) hash(typ RegisteredNodeType) (h common.Hash) {
//	hw := sha3.NewKeccak256()
//	d := &struct {
//		Type         RegisteredNodeType
//		InternalIp   string
//		InternalPort string
//		ExternalIp   string
//		ExternalPort string
//	}{
//		Type:         typ,
//		InternalIp:   node.InternalIp,
//		InternalPort: node.InternalPort,
//		ExternalIp:   node.ExternalIp,
//		ExternalPort: node.ExternalPort,
//	}
//	rlp.Encode(hw, d)
//	hw.Sum(h[:0])
//	return h
//}
//
//func (node *RegisteredNodeInfo) hashByCreateTime(typ RegisteredNodeType) (h common.Hash) {
//	hw := sha3.NewKeccak256()
//	d := &struct {
//		Type         RegisteredNodeType
//		InternalIp   string
//		InternalPort string
//		ExternalIp   string
//		ExternalPort string
//		GetCreateAt     uint64
//	}{
//		Type:         typ,
//		InternalIp:   node.InternalIp,
//		InternalPort: node.InternalPort,
//		ExternalIp:   node.ExternalIp,
//		ExternalPort: node.ExternalPort,
//		GetCreateAt:     uint64(timeutils.UnixMsec()),
//	}
//	rlp.Encode(hw, d)
//	hw.Sum(h[:0])
//	return h
//}
