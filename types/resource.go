package types

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/common"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"io"
	"strings"
	"sync/atomic"
)

type Resource struct {
	data *libtypes.ResourcePB

	// caches
	hash atomic.Value
	size atomic.Value
}

func NewResource(data *libtypes.ResourcePB) *Resource {
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
		m.data = new(libtypes.ResourcePB)
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

func (m *Resource) GetIdentityId() string { return m.data.IdentityId }
func (m *Resource) GetNodeId() string { return m.data.NodeId }
func (m *Resource) GetNodeName() string                   { return m.data.NodeName }
func (m *Resource) GetDataStatus() apicommonpb.DataStatus { return m.data.DataStatus }
func (m *Resource) GetState() apicommonpb.PowerState      { return m.data.State }
func (m *Resource) GetTotalMem() uint64                   { return m.data.TotalMem }
func (m *Resource) GetUsedMem() uint64 { return m.data.UsedMem }
func (m *Resource) GetTotalProcessor() uint64 { return uint64(m.data.TotalProcessor) }
func (m *Resource) GetUsedProcessor() uint64 { return uint64(m.data.UsedProcessor) }
func (m *Resource) GetTotalBandWidth() uint64 { return m.data.TotalBandwidth }
func (m *Resource) GetUsedBandWidth() uint64 { return m.data.UsedBandwidth }
func (m *Resource) String() string {
	//return fmt.Sprintf(`{"identity": %s, "nodeId": %s, "nodeName": %s, "dataId": %s, "dataStatus": %s, "state": %s, "totalMem": %d, "usedMem": %d, "totalProcessor": %d, "usedProcessor": %d, "totalBandWidth": %d, "usedBandWidth": %d}`)
	return m.data.String()
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


func NewResourceArray(metadataArr []*libtypes.ResourcePB) ResourceArray {
	var s ResourceArray
	for _, v := range metadataArr {
		s = append(s, NewResource(v))
	}
	return s
}

func (s ResourceArray) To() []*libtypes.ResourcePB {
	arr := make([]*libtypes.ResourcePB, 0, s.Len())
	for _, v := range s {
		arr = append(arr, v.data)
	}
	return arr
}

func  (s ResourceArray) String () string {
	arr := make([]string, len(s))
	for i, iden := range s {
		arr[i] = iden.String()
	}
	if len(arr) != 0 {
		return "[" +  strings.Join(arr, ",") + "]"
	}
	return ""
}

// 新增 local Resource
type LocalResource struct {
	data *libtypes.LocalResourcePB

	// caches
	hash atomic.Value
	size atomic.Value
}

func NewLocalResource(data *libtypes.LocalResourcePB) *LocalResource {
	return &LocalResource{data: data}
}
func (m *LocalResource) GetData() *libtypes.LocalResourcePB { return m.data }
func (m *LocalResource) GetIdentityId() string              { return m.data.IdentityId }
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
		m.data = new(libtypes.LocalResourcePB)
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

func NewLocalResourceArray(metadataArr []*libtypes.LocalResourcePB) LocalResourceArray {
	var s LocalResourceArray
	for _, v := range metadataArr {
		s = append(s, NewLocalResource(v))
	}
	return s
}

func (s LocalResourceArray) To() []*libtypes.LocalResourcePB {
	arr := make([]*libtypes.LocalResourcePB, 0, s.Len())
	for _, v := range s {
		arr = append(arr, v.data)
	}
	return arr
}
func (s LocalResourceArray) String() string {
	arr := make([]string, len(s))
	for i, r := range s {
		arr[i] = r.data.String()
	}
	if len(arr) != 0 {
		return "[" +  strings.Join(arr, ",") + "]"
	}
	return ""
}