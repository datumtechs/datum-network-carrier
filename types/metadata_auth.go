package types

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/common"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"io"
	"sync/atomic"
)

type MetadataAuth struct {
	data *libTypes.AuthRecordPB

	// caches
	hash atomic.Value
	size atomic.Value
}

func NewMedataAuth(data *libTypes.AuthRecordPB) *MetadataAuth {
	return &MetadataAuth{data: data}
}

func (m *MetadataAuth) EncodePB(w io.Writer) error {
	data, err := m.data.Marshal()
	if err == nil {
		w.Write(data)
	}
	return err
}

func (m *MetadataAuth) DecodePB(data []byte) error {
	if m.data == nil {
		m.data = new(libTypes.AuthRecordPB)
	}
	m.size.Store(common.StorageSize(len(data)))
	return m.data.Unmarshal(data)
}

func (m *MetadataAuth) Hash() common.Hash {
	if hash := m.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	buffer := new(bytes.Buffer)
	m.EncodePB(buffer)
	v := protoBufHash(buffer.Bytes())
	m.hash.Store(v)
	return v
}

func (m *MetadataAuth) User() string                   { return m.data.User }
func (m *MetadataAuth) UserType() apipb.UserType       { return m.data.GetUserType() }
func (m *MetadataAuth) Data() *libTypes.AuthRecordPB { return m.data }

type MetadataAuthArray []*MetadataAuth

// Len returns the length of s.
func (s MetadataAuthArray) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s MetadataAuthArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s MetadataAuthArray) GetPb(i int) []byte {
	buffer := new(bytes.Buffer)
	s[i].EncodePB(buffer)
	return buffer.Bytes()
}

func NewMetadataAuthArray(metaData []*libTypes.MetadataPB) MetadataArray {
	var s MetadataArray
	for _, v := range metaData {
		s = append(s, NewMetadata(v))
	}
	return s
}

func (s MetadataAuthArray) ToArray() []*libTypes.AuthRecordPB {
	arr := make([]*libTypes.AuthRecordPB, 0, s.Len())
	for _, v := range s {
		arr = append(arr, v.data)
	}
	return arr
}
