package types

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/common"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
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

func (m *Task) GetTaskSender() *apicommonpb.TaskOrganization {
	return m.data.GetSender()
}

func (m *Task) SetEventList(eventList []*libtypes.TaskEvent) {
	m.data.TaskEvents = eventList
}
func (m *Task) SetPowerSuppliers(arr []*apicommonpb.TaskOrganization) {
	m.data.PowerSuppliers = arr
}
func (m *Task) RemovePowerSuppliers() {
	m.data.PowerSuppliers = make([]*apicommonpb.TaskOrganization, 0)
}
func (m *Task) SetReceivers(arr []*apicommonpb.TaskOrganization) {
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

//type TaskEventShowAndRole struct {
//	Roles []apicommonpb.TaskRole
//	Data  *pb.TaskDetailShow
//}
//
//func (t *TaskEventShowAndRole) GetRoles() []apicommonpb.TaskRole { return t.Roles }
//func (t *TaskEventShowAndRole) GetData() *pb.TaskDetailShow   { return t.Data }

type TaskResultFileSummary struct {
	TaskId     string
	TableName  string
	MetadataId string
	OriginId   string
	FilePath   string
	NodeId     string
}

func NewTaskResultFileSummary(taskId, tableName, metadataId, originId, filePath, id string) *TaskResultFileSummary {
	return &TaskResultFileSummary{
		TaskId:     taskId,
		TableName:  tableName,
		MetadataId: metadataId,
		OriginId:   originId,
		FilePath:   filePath,
		NodeId:     id,
	}
}
func (trfs *TaskResultFileSummary) GetTaskId() string     { return trfs.TaskId }
func (trfs *TaskResultFileSummary) GetTableName() string  { return trfs.TableName }
func (trfs *TaskResultFileSummary) GetMetadataId() string { return trfs.MetadataId }
func (trfs *TaskResultFileSummary) GetOriginId() string   { return trfs.OriginId }
func (trfs *TaskResultFileSummary) GetFilePath() string   { return trfs.FilePath }
func (trfs *TaskResultFileSummary) GetNodeId() string     { return trfs.NodeId }

type TaskResultFileSummaryArr []*TaskResultFileSummary
