package types

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
	"sync/atomic"
)

const (
	PREFIX_POWER_ID = "power:"
	PREFIX_METADATA_ID = "metadata:"
	PREFIX_TASK_ID = "task:"
)

type Msg interface {
	Marshal() ([]byte, error)
	Unmarshal(b []byte) error
	String() string
	MsgType() string
}

// ------------------- SeedNode -------------------
type SeedNodeMsg struct {
	InternalIp   string                    `json:"internalIp"`
	InternalPort string                    `json:"internalPort"`
}
// ------------------- JobNode AND DataNode -------------------
type RegisteredNodeMsg struct {
	InternalIp   string                    `json:"internalIp"`
	InternalPort string                    `json:"internalPort"`
	ExternalIp   string                    `json:"externalIp"`
	ExternalPort string                    `json:"externalPort"`
}



// ------------------- identity -------------------
type IdentityMsg struct {
	*NodeAlias
	CreateAt uint64 `json:"createAt"`
}

// ------------------- power -------------------

type PowerMsg struct {
	// This is only used when marshaling to JSON.
	PowerId string `json:"powerId"`
	Data *powerData
	// caches
	hash atomic.Value
}
type powerData struct {
	*NodeAlias
	JobNodeId   string `json:"jobNodeId"`
	Information struct {
		Mem       uint64 `json:"mem,omitempty"`
		Processor uint64 `json:"processor,omitempty"`
		Bandwidth uint64 `json:"bandwidth,omitempty"`
	} `json:"information"`
	CreateAt uint64 `json:"createAt"`
}
type PowerMsgs []*PowerMsg


func (msg *PowerMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *PowerMsg) Unmarshal(b []byte) error { return nil }
func (msg *PowerMsg) String() string           { return "" }
func (msg *PowerMsg) MsgType() string          { return "" }

func (msg *PowerMsg) Onwer () *NodeAlias {
	return &NodeAlias{
		Name: msg.Data.Name,
		NodeId: msg.Data.NodeId,
		IdentityId: msg.Data.IdentityId,
	}
}
func (msg *PowerMsg) JobNodeId() string {
	return msg.Data.JobNodeId
}
func (msg *PowerMsg) Memory() string {
	return msg.Data.Information.Mem
}
func (msg)

func (msg *PowerMsg) GetPowerId () string {
	if "" != msg.PowerId {
		return msg.PowerId
	}
	msg.PowerId = PREFIX_POWER_ID + msg.Hash().Hex()
	return msg.PowerId
}
func (msg *PowerMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlpHash(msg.Data)
	msg.hash.Store(v)
	return v
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

// Len returns the length of s.
func (s PowerMsgs) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s PowerMsgs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// ------------------- metaData -------------------

type MetaDataMsg struct {
	MetaDataId string `json:"metaDataId"`
	Data *metadataData
	// caches
	hash atomic.Value
}
type metadataData struct {
	*NodeAlias
	Information struct {
		MetaSummary struct {
			OriginId   string `json:"originId,omitempty"`
			TableName  string `json:"tableName,omitempty"`
			Desc       string `json:"desc,omitempty"`
			FilePath   string `json:"filePath,omitempty"`
			Rows       uint32 `json:"rows,omitempty"`
			Columns    uint32 `json:"columns,omitempty"`
			Size_      string `json:"size,omitempty"`
			FileType   string `json:"fileType,omitempty"`
			HasTitle   bool   `json:"hasTitle,omitempty"`
			State      string `json:"state,omitempty"`
		} `json:"metaSummary"`
		ColumnMeta []*struct {
			Cindex   uint64 `json:"cindex,omitempty"`
			Cname    string `json:"cname,omitempty"`
			Ctype    string `json:"ctype,omitempty"`
			Csize    uint64 `json:"csize,omitempty"`
			Ccomment string `json:"ccomment,omitempty"`
		} `json:"columnMeta"`
	} `json:"information"`
	CreateAt uint64 `json:"createAt"`
}

type MetaDataMsgs []*MetaDataMsg

func (msg *MetaDataMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *MetaDataMsg) Unmarshal(b []byte) error { return nil }
func (msg *MetaDataMsg) String() string           { return "" }
func (msg *MetaDataMsg) MsgType() string          { return "" }

func (msg *MetaDataMsg) GetMetaDataId () string {
	if "" != msg.MetaDataId {
		return msg.MetaDataId
	}
	msg.MetaDataId = PREFIX_METADATA_ID + msg.Hash().Hex()
	return msg.MetaDataId
}
func (msg *MetaDataMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlpHash(msg.Data)
	msg.hash.Store(v)
	return v
}

// Len returns the length of s.
func (s MetaDataMsgs) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s MetaDataMsgs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// ------------------- task -------------------

type TaskMsg struct {
	TaskId                string                `json:"taskId"`
	Data *taskdata
	// caches
	hash atomic.Value
}
type taskdata struct {
	TaskName              string                `json:"taskName"`
	Owner                 *TaskSupplier         `json:"owner"`
	Partners              []*TaskSupplier       `json:"partners"`
	Receivers             []*TaskResultReceiver `json:"receivers"`
	CalculateContractCode string                `json:"calculateContractCode"`
	DataSplitContractCode string                `json:"dataSplitContractCode"`
	OperationCost         *TaskOperationCost    `json:"spend"`
	CreateAt              uint64                `json:"createAt"`
}

type TaskMsgs []*TaskMsg

func (msg *TaskMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *TaskMsg) Unmarshal(b []byte) error { return nil }
func (msg *TaskMsg) String() string           { return "" }
func (msg *TaskMsg) MsgType() string          { return "" }

func (msg *TaskMsg) GetTaskId () string {
	if "" != msg.TaskId {
		return msg.TaskId
	}
	msg.TaskId = PREFIX_TASK_ID + msg.Hash().Hex()
	return msg.TaskId
}
func (msg *TaskMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlpHash(msg.Data)
	msg.hash.Store(v)
	return v
}

// Len returns the length of s.
func (s TaskMsgs) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s TaskMsgs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type TaskSupplier struct {
	*NodeAlias
	MetaData *SupplierMetaData `json:"metaData"`
}

type SupplierMetaData struct {
	MetaId          string   `json:"metaId"`
	ColumnIndexList []uint32 `json:"columnIndexList"`
}

type TaskResultReceiver struct {
	*NodeAlias
	Providers []*NodeAlias `json:"providers"`
}

type TaskOperationCost struct {
	Processor uint64 `json:"processor"`
	Mem       uint64 `json:"mem"`
	Bandwidth uint64 `json:"bandwidth"`
	Duration  uint64 `json:"duration"`
}

// ------------------- data using authorize -------------------

type DataAuthorizationApply struct {
	PoposalHash string     `json:"poposalHash"`
	Proposer    *NodeAlias `json:"proposer"`
	Approver    *NodeAlias `json:"approver"`
	Apply       struct {
		MetaId       string `json:"metaId"`
		UseCount     uint64 `json:"useCount"`
		UseStartTime string `json:"useStartTime"`
		UseEndTime   string `json:"useEndTime"`
	} `json:"apply"`
}

type DataAuthorizationConfirm struct {
	PoposalHash string     `json:"poposalHash"`
	Proposer    *NodeAlias `json:"proposer"`
	Approver    *NodeAlias `json:"approver"`
	Approve     struct {
		Vote   uint16 `json:"vote"`
		Reason string `json:"reason"`
	} `json:"approve"`
}

// ------------------- common -------------------

type NodeAlias struct {
	Name       string `json:"name"`
	NodeId     string `json:"nodeId"`
	IdentityId string `json:"identityId"`
}
