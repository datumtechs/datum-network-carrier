package types

import (
	"bytes"
	"encoding/json"
	"github.com/RosettaFlow/Carrier-Go/common"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"io"
	"sync/atomic"
)

type Metadata struct {
	data *libTypes.MetaData

	// caches
	hash atomic.Value
	size atomic.Value
}

func NewMetadata(data *libTypes.MetaData) *Metadata {
	return &Metadata{data: data}
}

func NewMetadataWithoutParam() *Metadata {
	return &Metadata{data: new(libTypes.MetaData)}
}

func (m *Metadata) EncodePb(w io.Writer) error {
	data, err := m.data.Marshal()
	if err == nil {
		w.Write(data)
	}
	return err
}

func MetaDataToJson(meta *Metadata) string {
	result, err := json.Marshal(meta.data)
	if err != nil {
		panic("Convert To json fail")
	}
	return string(result)
}

func (m *Metadata) DecodePb(data []byte) error {
	if m.data == nil {
		m.data = new(libTypes.MetaData)
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

func (m *Metadata) MetadataData() *libTypes.MetaData { return m.data }
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

func NewMetadataArray(metaData []*libTypes.MetaData) MetadataArray {
	var s MetadataArray
	for _, v := range metaData {
		s = append(s, NewMetadata(v))
	}
	return s
}

func (s MetadataArray) To() []*libTypes.MetaData {
	arr := make([]*libTypes.MetaData, 0, s.Len())
	for _, v := range s {
		arr = append(arr, v.data)
	}
	return arr
}

// ----------------------------------------- metaData -----------------------------------------
//type OrgMetaDataInfo struct {
//	Owner    *NodeAlias    `json:"owner"`
//	MetaData *MetaDataInfo `json:"metaData"`
//}
//
//type MetaDataInfo struct {
//	MetaDataSummary *libTypes.MetaDataSummary 			  `json:"metaDataSummary"`
//	ColumnMetas     []*libTypes.MetadataColumn    `json:"columnMetas"`
//}

//func ConvertMetaDataInfoToPB (metadata *MetaDataInfo) *libTypes.MetadataDetail {
//	columns := make([]*libTypes.MetadataColumn, len(metadata.ColumnMetas))
//	for j, column := range metadata.ColumnMetas {
//		columns[j] = column
//	}
//	return &libTypes.MetadataDetail{
//		MetaDataSummary: &libTypes.MetaDataSummary{
//			MetaDataId: metadata.MetaDataSummary.MetaDataId,
//			OriginId: metadata.MetaDataSummary.OriginId,
//			TableName: metadata.MetaDataSummary.TableName,
//			Desc: metadata.MetaDataSummary.Desc,
//			FilePath: metadata.MetaDataSummary.FilePath,
//			Rows:metadata.MetaDataSummary.Rows,
//			Columns: metadata.MetaDataSummary.Columns,
//			Size_: uint32(metadata.MetaDataSummary.Size()),
//			FileType: metadata.MetaDataSummary.FileType,
//			HasTitle: metadata.MetaDataSummary.HasTitle,
//			State: metadata.MetaDataSummary.State,
//		},
//		MetadataColumnList: columns,
//	}
//}