package types

import (
	"bytes"
	"encoding/json"
	"github.com/RosettaFlow/Carrier-Go/common"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"io"
	"strings"
	"sync/atomic"
)

type Identity struct {
	data *libtypes.IdentityPB

	// caches
	hash atomic.Value
	size atomic.Value
}

func NewIdentity(data *libtypes.IdentityPB) *Identity {
	return &Identity{data: data}
}

func (m *Identity) EncodePb(w io.Writer) error {
	data, err := m.data.Marshal()
	if err == nil {
		w.Write(data)
	}
	return err
}

func IdentityDataTojson(identity *Identity) string {
	result, err := json.Marshal(identity.data)
	if err != nil {
		panic("Convert To json fail")
	}
	return string(result)
}

func (m *Identity) DecodePb(data []byte) error {
	if m.data == nil {
		m.data = new(libtypes.IdentityPB)
	}
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

func (m *Identity) GetName() string { return m.data.GetNodeName() }
func (m *Identity) GetNodeId() string { return m.data.GetNodeId() }
func (m *Identity) GetIdentityId() string { return m.data.GetIdentityId() }
func (m *Identity) GetImageUrl() string { return m.data.GetImageUrl() }
func (m *Identity) GetDetails() string { return m.data.GetDetails() }
func (m *Identity) GetCredential() string { return m.data.GetCredential() }
func (m *Identity) GetUpdateAt() uint64 { return m.data.GetUpdateAt() }
func (m *Identity) GetStatus() apicommonpb.CommonStatus { return m.data.GetStatus() }
func (m *Identity) GetDataStatus() apicommonpb.DataStatus { return m.data.GetDataStatus() }
func (m *Identity) GetDataId() string { return m.data.GetDataId() }

func (m *Identity) String() string {
	//return fmt.Sprintf(`{"identity": %s, "nodeId": %s, "nodeName": %s, "dataId": %s, "dataStatus": %s, "status": %s}`,
	//	m.data.Identity, m.data.NodeId, m.data.NodeName, m.data.DataId, m.data.DataStatus, m.data.Status)
	return m.data.String()
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

func NewIdentityArray(metaData []*libtypes.IdentityPB) IdentityArray {
	var s IdentityArray
	for _, v := range metaData {
		s = append(s, NewIdentity(v))
	}
	return s
}

func (s IdentityArray) To() []*libtypes.IdentityPB {
	arr := make([]*libtypes.IdentityPB, 0, s.Len())
	for _, v := range s {
		arr = append(arr, v.data)
	}
	return arr
}

func (s IdentityArray) String() string {
	arr := make([]string, len(s))
	for i, iden := range s {
		arr[i] = iden.String()
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return ""
}
