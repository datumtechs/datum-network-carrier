package types

import (
	"bytes"
	"encoding/json"
	"github.com/datumtechs/datum-network-carrier/common"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"io"
	"sync/atomic"
)

type Metadata struct {
	data *carriertypespb.MetadataPB

	// caches
	hash atomic.Value
	size atomic.Value
}

func NewMetadata(data *carriertypespb.MetadataPB) *Metadata {
	return &Metadata{data: data}
}

func NewMetadataWithoutParam() *Metadata {
	return &Metadata{data: new(carriertypespb.MetadataPB)}
}

func (m *Metadata) EncodePb(w io.Writer) error {
	data, err := m.data.Marshal()
	if err == nil {
		w.Write(data)
	}
	return err
}

func MetadataToJson(meta *Metadata) string {
	result, err := json.Marshal(meta.data)
	if err != nil {
		panic("Convert To json fail")
	}
	return string(result)
}

func (m *Metadata) DecodePb(data []byte) error {
	if m.data == nil {
		m.data = new(carriertypespb.MetadataPB)
	}
	m.size.Store(common.StorageSize(len(data)))
	return m.data.Unmarshal(data)
}

func (m *Metadata) Hash() common.Hash {
	if hash := m.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	buffer := new(bytes.Buffer)
	m.EncodePb(buffer)
	v := protoBufHash(buffer.Bytes())
	m.hash.Store(v)
	return v
}

func (m *Metadata) GetData() *carriertypespb.MetadataPB { return m.data }
// MetadataArray is a Transaction slice type for basic sorting.
type MetadataArray []*Metadata

// Len returns the length of s.
func (s MetadataArray) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s MetadataArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s MetadataArray) GetPb(i int) []byte {
	buffer := new(bytes.Buffer)
	s[i].EncodePb(buffer)
	return buffer.Bytes()
}

func NewMetadataArray(metaData []*carriertypespb.MetadataPB) MetadataArray {
	var s MetadataArray
	for _, v := range metaData {
		s = append(s, NewMetadata(v))
	}
	return s
}

func (s MetadataArray) To() []*carriertypespb.MetadataPB {
	arr := make([]*carriertypespb.MetadataPB, 0, s.Len())
	for _, v := range s {
		arr = append(arr, v.data)
	}
	return arr
}