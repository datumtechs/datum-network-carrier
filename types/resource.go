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

func (m *Resource) GetIdentityId() string { return m.data.Identity }
func (m *Resource) GetNodeId() string { return m.data.NodeId }
func (m *Resource) GetNodeName() string { return m.data.NodeName }
func (m *Resource) GetDataStatus() string { return m.data.DataStatus }
func (m *Resource) GetState() string { return m.data.State }
func (m *Resource) GetTotalMem() uint64 { return m.data.TotalMem }
func (m *Resource) GetUsedMem() uint64 { return m.data.UsedMem }
func (m *Resource) GetTotalProcessor() uint64 { return m.data.TotalProcessor }
func (m *Resource) GetUsedProcessor() uint64 { return m.data.UsedProcessor }
func (m *Resource) GetTotalBandWidth() uint64 { return m.data.TotalBandWidth }
func (m *Resource) GetUsedBandWidth() uint64 { return m.data.UsedBandWidth }

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

// 新增 local Resource
type LocalResource struct {
	data *libTypes.LocalResourceData

	// caches
	hash atomic.Value
	size atomic.Value
}

func NewLocalResource(data *libTypes.LocalResourceData) *LocalResource {
	return &LocalResource{data: data}
}
func (m *LocalResource) GetData() *libTypes.LocalResourceData { return m.data }
func (m *LocalResource) GetIdentityId() string { return m.data.Identity }
func (m *LocalResource) GetJobNodeId() string { return m.data.JobNodeId }
func (m *LocalResource) EncodePb(w io.Writer) error {
	data, err := m.data.Marshal()
	if err == nil {
		w.Write(data)
	}
	return err
}

func (m *LocalResource) DecodePb(data []byte) error {
	if m.data == nil {
		m.data = new(libTypes.LocalResourceData)
	}
	m.size.Store(common.StorageSize(len(data)))
	return m.data.Unmarshal(data)
}

func (m *LocalResource) Hash() common.Hash {
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
type LocalResourceArray []*LocalResource

// Len returns the length of s.
func (s LocalResourceArray) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s LocalResourceArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s LocalResourceArray) GetPb(i int) []byte {
	buffer := new(bytes.Buffer)
	s[i].EncodePb(buffer)
	return buffer.Bytes()
}

func NewLocalResourceArray(metaData []*libTypes.LocalResourceData) LocalResourceArray {
	var s LocalResourceArray
	for _, v := range metaData {
		s = append(s, NewLocalResource(v))
	}
	return s
}

func (s LocalResourceArray) To() []*libTypes.LocalResourceData {
	arr := make([]*libTypes.LocalResourceData, 0, s.Len())
	for _, v := range s {
		arr = append(arr, v.data)
	}
	return arr
}