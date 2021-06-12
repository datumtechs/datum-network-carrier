package types

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/common"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"io"
	"sync/atomic"
)

type Resource struct {
	data *libTypes.ResourceData

	// caches
	hash atomic.Value
	size atomic.Value
}

func NewResource(data *libTypes.ResourceData) *Resource {
	return &Resource{data: data}
}

func (m *Resource) EncodePb(w io.Writer) error {
	data, err := m.data.Marshal()
	if err == nil {
		w.Write(data)
	}
	return err
}

func (m *Resource) DecodePb(data []byte) error {
	if m.data == nil {
		m.data = new(libTypes.ResourceData)
	}
	m.size.Store(common.StorageSize(len(data)))
	return m.data.Unmarshal(data)
}

func (m *Resource) Hash() common.Hash {
	if hash := m.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	buffer := new(bytes.Buffer)
	m.EncodePb(buffer)
	v := protoBufHash(buffer.Bytes())
	m.hash.Store(v)
	return v
}

// ResourceArray is a Transaction slice type for basic sorting.
type ResourceArray []*Resource

// Len returns the length of s.
func (s ResourceArray) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s ResourceArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ResourceArray) GetPb(i int) []byte {
	buffer := new(bytes.Buffer)
	s[i].EncodePb(buffer)
	return buffer.Bytes()
}

func NewResourceArray(metaData []*libTypes.ResourceData) ResourceArray {
	var s ResourceArray
	for _, v := range metaData {
		s = append(s, NewResource(v))
	}
	return s
}

func (s ResourceArray) To() []*libTypes.ResourceData {
	arr := make([]*libTypes.ResourceData, 0, s.Len())
	for _, v := range s {
		arr = append(arr, v.data)
	}
	return arr
}

