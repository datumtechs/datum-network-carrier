package types

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"io"
	"sync/atomic"
)

type Task struct {
	data *libtypes.TaskPB

	// caches
	hash atomic.Value
	size atomic.Value
}

func NewTask(data *libtypes.TaskPB) *Task {
	return &Task{data: data}
}

func (m *Task) EncodePb(w io.Writer) error {
	if m.data == nil {
		m.data = new(libtypes.TaskPB)
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

func (m *Task) GetTaskId() string {
	return m.data.TaskId
}

func (m *Task) GetTaskData() *libtypes.TaskPB {
	return m.data
}

func (m *Task) GetTaskSender() *libtypes.TaskOrganization {
	return m.data.GetSender()
}

func (m *Task) SetEventList(eventList []*libtypes.TaskEvent) {
	m.data.TaskEvents = eventList
}
func (m *Task) SetPowerSuppliers(arr []*libtypes.TaskOrganization) {
	m.data.PowerSuppliers = arr
}
func (m *Task) SetPowerResources(arr []*libtypes.TaskPowerResourceOption) {
	m.data.PowerResourceOptions = arr
}
func (m *Task) RemovePowerSuppliers() {
	m.data.PowerSuppliers = make([]*libtypes.TaskOrganization, 0)
}
func (m *Task) SetReceivers(arr []*libtypes.TaskOrganization) {
	m.data.Receivers = arr
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

func NewTaskDataArray(metaData []*libtypes.TaskPB) TaskDataArray {
	var s TaskDataArray
	for _, v := range metaData {
		s = append(s, NewTask(v))
	}
	return s
}

func (s TaskDataArray) To() []*libtypes.TaskPB {
	arr := make([]*libtypes.TaskPB, 0, s.Len())
	for _, v := range s {
		arr = append(arr, v.data)
	}
	return arr
}

type TaskResultFileSummary struct {
	TaskId       string
	MetadataId   string
	OriginId     string
	MetadataName string
	DataPath     string
	NodeId       string
	Extra        string
}

func NewTaskResultFileSummary(taskId, metadataId, originId, metadataName, dataPath, id, extra string) *TaskResultFileSummary {
	return &TaskResultFileSummary{
		TaskId:       taskId,
		MetadataId:   metadataId,
		OriginId:     originId,
		MetadataName: metadataName,
		DataPath:     dataPath,
		NodeId:       id,
		Extra:        extra,
	}
}
func (trfs *TaskResultFileSummary) GetTaskId() string       { return trfs.TaskId }
func (trfs *TaskResultFileSummary) GetMetadataName() string { return trfs.MetadataName }
func (trfs *TaskResultFileSummary) GetMetadataId() string   { return trfs.MetadataId }
func (trfs *TaskResultFileSummary) GetOriginId() string     { return trfs.OriginId }
func (trfs *TaskResultFileSummary) GetDataPath() string     { return trfs.DataPath }
func (trfs *TaskResultFileSummary) GetNodeId() string       { return trfs.NodeId }
func (trfs *TaskResultFileSummary) GetExtra() string        { return trfs.Extra }

type TaskResultFileSummaryArr []*TaskResultFileSummary
