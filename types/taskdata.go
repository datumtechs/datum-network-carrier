package types

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/common"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"io"
	"sync/atomic"
)

type Task struct {
	data *libTypes.TaskData

	// caches
	hash atomic.Value
	size atomic.Value
}

func NewTask(data *libTypes.TaskData) *Task {
	return &Task{data: data}
}

func (m *Task) EncodePb(w io.Writer) error {
	if m.data == nil {
		m.data = new(libTypes.TaskData)
	}
	data, err := m.data.Marshal()
	if err == nil {
		w.Write(data)
	}
	return err
}

func (m *Task) DecodePb(data []byte) error {
	m.size.Store(common.StorageSize(len(data)))
	return m.data.Unmarshal(data)
}

func (m *Task) Hash() common.Hash {
	if hash := m.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	buffer := new(bytes.Buffer)
	m.EncodePb(buffer)
	v := protoBufHash(buffer.Bytes())
	m.hash.Store(v)
	return v
}

func (m *Task) TaskId() string {
	return m.data.TaskId
}

func (m *Task) TaskData() *libTypes.TaskData {
	return m.data
}

// TaskDataArray is a Transaction slice type for basic sorting.
type TaskDataArray []*Task

// Len returns the length of s.
func (s TaskDataArray) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s TaskDataArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s TaskDataArray) GetPb(i int) []byte {
	buffer := new(bytes.Buffer)
	s[i].EncodePb(buffer)
	return buffer.Bytes()
}

func NewTaskDataArray(metaData []*libTypes.TaskData) TaskDataArray {
	var s TaskDataArray
	for _, v := range metaData {
		s = append(s, NewTask(v))
	}
	return s
}

func (s TaskDataArray) To() []*libTypes.TaskData {
	arr := make([]*libTypes.TaskData, 0, s.Len())
	for _, v := range s {
		arr = append(arr, v.data)
	}
	return arr
}

type TaskDetailShow struct {
	TaskId        string                   `json:"taskId"`
	TaskName      string                   `json:"taskName"`
	Owner         *NodeAlias               `json:"owner"`
	AlgoSupplier  *NodeAlias               `json:"algoSupplier"`
	DataSupplier  []*TaskDataSupplierShow  `json:"dataSupplier"`
	PowerSupplier []*TaskPowerSupplierShow `json:"powerSupplier"`
	Receivers     []*NodeAlias             `json:"receivers"`
	CreateAt      string                   `json:"createat"`
	EndAt         string                   `json:"endAt"`
	State         string                   `json:"state"`
	OperationCost *TaskOperationCost       `json:"operationCost"`
}

type TaskDataSupplierShow struct {
	MemberInfo *NodeAlias `json:"memberInfo"`
	MetaId     string     `json:"metaId"`
	MetaName   string     `json:"metaName"`
}
type TaskPowerSupplierShow struct {
	MemberInfo    *NodeAlias     `json:"memberInfo"`
	ResourceUsage *ResourceUsage `json:"resourceUsage"`
}
