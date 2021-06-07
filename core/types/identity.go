package types

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/common"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"io"
	"sync/atomic"
)


type Identity struct {
	data *libTypes.IdentityData

	// caches
	hash atomic.Value
	size atomic.Value
}

func NewIdentity(data *libTypes.IdentityData) *Identity {
	return &Identity{data: data}
}

func (m *Identity) EncodePb(w io.Writer) error {
	data, err := m.data.Marshal()
	if err != nil {
		w.Write(data)
	}
	return err
}

func (m *Identity) DecodePb(data []byte) error {
	m.size.Store(common.StorageSize(len(data)))
	return m.data.Unmarshal(data)
}

func (m *Identity) Hash() common.Hash {
	if hash := m.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	buffer := new(bytes.Buffer)
	m.EncodePb(buffer)
	v := protoBufHash(buffer.Bytes())
	m.hash.Store(v)
	return v
}

// IdentityArray is a Transaction slice type for basic sorting.
type IdentityArray []*Identity

// Len returns the length of s.
func (s IdentityArray) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s IdentityArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s IdentityArray) GetPb(i int) []byte {
	buffer := new(bytes.Buffer)
	s[i].EncodePb(buffer)
	return buffer.Bytes()
}

func NewIdentityArray(metaData []*libTypes.IdentityData) IdentityArray {
	var s IdentityArray
	for _, v := range metaData {
		s = append(s, NewIdentity(v))
	}
	return s
}

func (s IdentityArray) To() []*libTypes.IdentityData {
	arr := make([]*libTypes.IdentityData, 0, s.Len())
	for _, v := range s {
		arr = append(arr, v.data)
	}
	return arr
}
