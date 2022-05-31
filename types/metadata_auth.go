package types

import (
	"bytes"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"io"
	"sync/atomic"
)

type MetadataAuthority struct {
	data *carriertypespb.MetadataAuthorityPB

	// caches
	hash atomic.Value
	size atomic.Value
}

func NewMetadataAuthority(data *carriertypespb.MetadataAuthorityPB) *MetadataAuthority {
	return &MetadataAuthority{data: data}
}

func (m *MetadataAuthority) EncodePB(w io.Writer) error {
	data, err := m.data.Marshal()
	if err == nil {
		w.Write(data)
	}
	return err
}

func (m *MetadataAuthority) DecodePB(data []byte) error {
	if m.data == nil {
		m.data = new(carriertypespb.MetadataAuthorityPB)
	}
	m.size.Store(common.StorageSize(len(data)))
	return m.data.Unmarshal(data)
}

func (m *MetadataAuthority) Hash() common.Hash {
	if hash := m.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	buffer := new(bytes.Buffer)
	m.EncodePB(buffer)
	v := protoBufHash(buffer.Bytes())
	m.hash.Store(v)
	return v
}

func (m *MetadataAuthority) GetUser() string                        { return m.data.GetUser() }
func (m *MetadataAuthority) GetUserType() commonconstantpb.UserType      { return m.data.GetUserType() }
func (m *MetadataAuthority) GetData() *carriertypespb.MetadataAuthorityPB { return m.data }

type MetadataAuthArray []*MetadataAuthority

// Len returns the length of s.
func (s MetadataAuthArray) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s MetadataAuthArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s MetadataAuthArray) GetPb(i int) []byte {
	buffer := new(bytes.Buffer)
	s[i].EncodePB(buffer)
	return buffer.Bytes()
}

func NewMetadataAuthArray(metaData []*carriertypespb.MetadataPB) MetadataArray {
	var s MetadataArray
	for _, v := range metaData {
		s = append(s, NewMetadata(v))
	}
	return s
}

func (s MetadataAuthArray) ToArray() []*carriertypespb.MetadataAuthorityPB {
	arr := make([]*carriertypespb.MetadataAuthorityPB, 0, s.Len())
	for _, v := range s {
		arr = append(arr, v.data)
	}
	return arr
}


type MetadataAuthAudit struct {
	MetadataAuthId  string
	AuditOption     commonconstantpb.AuditMetadataOption
	AuditSuggestion string
}

func NewMetadataAuthAudit(metadataAuthId, suggestion string, option commonconstantpb.AuditMetadataOption) *MetadataAuthAudit {
	return &MetadataAuthAudit{
		MetadataAuthId:  metadataAuthId,
		AuditOption:     option,
		AuditSuggestion: suggestion,
	}
}

func (maa *MetadataAuthAudit) GetMetadataAuthId() string                       { return maa.MetadataAuthId }
func (maa *MetadataAuthAudit) GetAuditOption() commonconstantpb.AuditMetadataOption { return maa.AuditOption }
func (maa *MetadataAuthAudit) GetAuditSuggestion() string                      { return maa.AuditSuggestion }

func (maa *MetadataAuthAudit) String() string {
	return fmt.Sprintf(`{"metadataAuthId": %s, "option": %s, "suggestion": %s}`,
		maa.GetMetadataAuthId(), maa.GetAuditOption().String(), maa.GetAuditSuggestion())
}