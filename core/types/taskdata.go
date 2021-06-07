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
	data, err := m.data.Marshal()
	if err != nil {
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

// TaskDataArray is a Transaction slice type for basic sorting.
type TaskDataArray []*libTypes.TaskData

// Len returns the length of s.
func (s TaskDataArray) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s TaskDataArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s TaskDataArray) GetPb(i int) []byte {
	enc, _ := s[i].Marshal()
	return enc
}

