package types

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/common"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"io"
	"sync/atomic"
)

type MetadataAuthority struct {
	data *libtypes.MetadataAuthorityPB

	// caches
	hash atomic.Value
	size atomic.Value
}

func NewMetadataAuthority(data *libtypes.MetadataAuthorityPB) *MetadataAuthority {
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
		m.data = new(libtypes.MetadataAuthorityPB)
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
func (m *MetadataAuthority) GetUserType() apicommonpb.UserType      { return m.data.GetUserType() }
func (m *MetadataAuthority) GetData() *libtypes.MetadataAuthorityPB { return m.data }

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

func NewMetadataAuthArray(metaData []*libtypes.MetadataPB) MetadataArray {
	var s MetadataArray
	for _, v := range metaData {
		s = append(s, NewMetadata(v))
	}
	return s
}

func (s MetadataAuthArray) ToArray() []*libtypes.MetadataAuthorityPB {
	arr := make([]*libtypes.MetadataAuthorityPB, 0, s.Len())
	for _, v := range s {
		arr = append(arr, v.data)
	}
	return arr
}

//type MetadataAuthApply struct {
//	MetadataAuthId string
//	User           string
//	UserType       apicommonpb.UserType
//	Auth           *libtypes.MetadataAuthority
//}
//
//func NewMetadataAuthApply(metadataAuthId, user string, userType apicommonpb.UserType, apply *libtypes.MetadataAuthority) *MetadataAuthApply {
//	return &MetadataAuthApply{
//		MetadataAuthId: metadataAuthId,
//		User:           user,
//		UserType:       userType,
//		Auth:           apply,
//	}
//}
//
//func (maa *MetadataAuthApply) GetMetadataAuthId() string            { return maa.MetadataAuthId }
//func (maa *MetadataAuthApply) GetUser() string                      { return maa.User }
//func (maa *MetadataAuthApply) GetUserType() apicommonpb.UserType    { return maa.UserType }
//func (maa *MetadataAuthApply) GetAuth() *libtypes.MetadataAuthority { return maa.Auth }

type MetadataAuthAudit struct {
	MetadataAuthId  string
	AuditOption     apicommonpb.AuditMetadataOption
	AuditSuggestion string
}

func NewMetadataAuthAudit(metadataAuthId, suggestion string, option apicommonpb.AuditMetadataOption) *MetadataAuthAudit {
	return &MetadataAuthAudit{
		MetadataAuthId:  metadataAuthId,
		AuditOption:     option,
		AuditSuggestion: suggestion,
	}
}

func (maa *MetadataAuthAudit) GetMetadataAuthId() string                       { return maa.MetadataAuthId }
func (maa *MetadataAuthAudit) GetAuditOption() apicommonpb.AuditMetadataOption { return maa.AuditOption }
func (maa *MetadataAuthAudit) GetAuditSuggestion() string                      { return maa.AuditSuggestion }
