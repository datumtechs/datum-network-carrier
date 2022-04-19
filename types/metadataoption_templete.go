package types

import (
	"fmt"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)


var (
	CannotMatchMetadataOption = fmt.Errorf("cannot match metadata option")
)

func IsNotRowAndColumnData(fileType libtypes.OrigindataType) bool { return !IsRowAndColumnData(fileType) }
func IsRowAndColumnData(fileType libtypes.OrigindataType) bool {
	if fileType == libtypes.OrigindataType_OrigindataType_CSV {
		return true
	}
	return false
}

/**
{
    "originId": "d9b41e7138544c63f9fe25f6aa4983819793e5b46f14652a1ff1b51f99f71783",
    "dataPath": "/home/user1/data/data_root/bank_predict_partyA_20220218-090241.csv",
    "rows": 100,
    "columns": 27,
    "size": 12,
    "hasTitle": true,
    "metadataColumns": [
        {
            "index": 1,
            "name": "CLIENT_ID",
            "type": "string",
            "size": 0,
            "comment": ""
        }
    ],
}
*/
// libtypes.OriginFileType_FileType_CSV |
type MetadataOptionRowAndColumn struct {
	OriginId        string
	DataPath        string
	Rows            uint64
	Columns         uint64
	Size            uint64
	HasTitle        bool
	MetadataColumns []*MetadataColumn
}

func (option *MetadataOptionRowAndColumn) GetOriginId() string { return option.OriginId }
func (option *MetadataOptionRowAndColumn) GetDataPath() string { return option.DataPath }
func (option *MetadataOptionRowAndColumn) GetRows() uint64     { return option.Rows }
func (option *MetadataOptionRowAndColumn) GetColumns() uint64  { return option.Columns }
func (option *MetadataOptionRowAndColumn) GetSize() uint64     { return option.Size }
func (option *MetadataOptionRowAndColumn) GetHasTitle() bool   { return option.HasTitle }
func (option *MetadataOptionRowAndColumn) GetMetadataColumns() []*MetadataColumn {
	return option.MetadataColumns
}

type MetadataColumn struct {
	Index   uint32
	Name    string
	Type    string
	Comment string
	Size    uint64
}

func (mc *MetadataColumn) GetIndex() uint32   { return mc.Index }
func (mc *MetadataColumn) GetName() string    { return mc.Name }
func (mc *MetadataColumn) GetType() string    { return mc.Type }
func (mc *MetadataColumn) GetComment() string { return mc.Comment }
func (mc *MetadataColumn) GetSize() uint64    { return mc.Size }
