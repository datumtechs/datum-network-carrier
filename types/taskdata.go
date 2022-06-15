package types

import (
	"bytes"
	"github.com/datumtechs/datum-network-carrier/common"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"io"
	"sync/atomic"
)

type Task struct {
	data *carriertypespb.TaskPB

	// caches
	hash atomic.Value
	size atomic.Value
}

func NewTask(data *carriertypespb.TaskPB) *Task {
	return &Task{data: data}
}

func (m *Task) EncodePb(w io.Writer) error {
	if m.data == nil {
		m.data = new(carriertypespb.TaskPB)
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
	return m.data.GetTaskId()
}

func (m *Task) GetTaskData() *carriertypespb.TaskPB {
	return m.data
}

func (m *Task) GetTaskSender() *carriertypespb.TaskOrganization {
	return m.data.GetSender()
}

func (m *Task) SetEventList(eventList []*carriertypespb.TaskEvent) {
	m.data.TaskEvents = eventList
}
func (m *Task) SetPowerSuppliers(arr []*carriertypespb.TaskOrganization) {
	m.data.PowerSuppliers = arr
}
func (m *Task) SetPowerResources(arr []*carriertypespb.TaskPowerResourceOption) {
	m.data.PowerResourceOptions = arr
}
func (m *Task) RemovePowerSuppliers() {
	m.data.PowerSuppliers = make([]*carriertypespb.TaskOrganization, 0)
}
func (m *Task) RemovePowerResources() {
	m.data.PowerResourceOptions = make([]*carriertypespb.TaskPowerResourceOption, 0)
}
func (m *Task) SetReceivers(arr []*carriertypespb.TaskOrganization) {
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

func NewTaskDataArray(metaData []*carriertypespb.TaskPB) TaskDataArray {
	var s TaskDataArray
	for _, v := range metaData {
		s = append(s, NewTask(v))
	}
	return s
}

func (s TaskDataArray) To() []*carriertypespb.TaskPB {
	arr := make([]*carriertypespb.TaskPB, 0, s.Len())
	for _, v := range s {
		arr = append(arr, v.data)
	}
	return arr
}

type TaskResultDataSummary struct {
	TaskId         string
	MetadataId     string
	OriginId       string
	MetadataName   string
	NodeId         string
	Extra          string
	DataHash       string
	MetadataOption string
	DataType       uint32
}

func NewTaskResultDataSummary(taskId, metadataId, originId, metadataName, dataHash, metadataOption, nodeId, extra string, dataType uint32) *TaskResultDataSummary {
	return &TaskResultDataSummary{
		TaskId:         taskId,
		MetadataId:     metadataId,
		OriginId:       originId,
		MetadataName:   metadataName,
		NodeId:         nodeId,
		Extra:          extra,
		DataHash:       dataHash,
		MetadataOption: metadataOption,
		DataType:       dataType,
	}
}
func (trfs *TaskResultDataSummary) GetTaskId() string         { return trfs.TaskId }
func (trfs *TaskResultDataSummary) GetMetadataName() string   { return trfs.MetadataName }
func (trfs *TaskResultDataSummary) GetMetadataId() string     { return trfs.MetadataId }
func (trfs *TaskResultDataSummary) GetOriginId() string       { return trfs.OriginId }
func (trfs *TaskResultDataSummary) GetNodeId() string         { return trfs.NodeId }
func (trfs *TaskResultDataSummary) GetExtra() string          { return trfs.Extra }
func (trfs *TaskResultDataSummary) GetDataHash() string       { return trfs.DataHash }
func (trfs *TaskResultDataSummary) GetMetadataOption() string { return trfs.MetadataOption }
func (trfs *TaskResultDataSummary) GetDataType() uint32       { return trfs.DataType }

type TaskResultDataSummaryArr []*TaskResultDataSummary
