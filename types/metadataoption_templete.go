package types

import (
	"fmt"
	libtypes "github.com/Metisnetwork/Metis-Carrier/lib/types"
)

var (
	CannotMatchMetadataOption = fmt.Errorf("cannot match metadata option")
)

func IsNotCSVdata(dataType libtypes.OrigindataType) bool { return !IsCSVdata(dataType) }
func IsCSVdata(dataType libtypes.OrigindataType) bool {
	if dataType == libtypes.OrigindataType_OrigindataType_CSV {
		return true
	}
	return false
}

func IsNotDIRdata(dataType libtypes.OrigindataType) bool { return !IsCSVdata(dataType) }
func IsDIRdata(dataType libtypes.OrigindataType) bool {
	if dataType == libtypes.OrigindataType_OrigindataType_DIR {
		return true
	}
	return false
}

func IsNotBINARYdata(dataType libtypes.OrigindataType) bool { return !IsCSVdata(dataType) }
func IsBINARYdata(dataType libtypes.OrigindataType) bool {
	if dataType == libtypes.OrigindataType_OrigindataType_BINARY {
		return true
	}
	return false
}

// ======================================================================================================

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
// libtypes.OrigindataType_CSV
type MetadataOptionCSV struct {
	OriginId        string            `json:"originId"`
	DataPath        string            `json:"dataPath"`
	Rows            uint64            `json:"rows"`
	Columns         uint64            `json:"columns"`
	Size            uint64            `json:"size"`
	HasTitle        bool              `json:"hasTitle"`
	MetadataColumns []*MetadataColumn `json:"metadataColumns"`
}

func (option *MetadataOptionCSV) GetOriginId() string { return option.OriginId }
func (option *MetadataOptionCSV) GetDataPath() string { return option.DataPath }
func (option *MetadataOptionCSV) GetRows() uint64     { return option.Rows }
func (option *MetadataOptionCSV) GetColumns() uint64  { return option.Columns }
func (option *MetadataOptionCSV) GetSize() uint64     { return option.Size }
func (option *MetadataOptionCSV) GetHasTitle() bool   { return option.HasTitle }
func (option *MetadataOptionCSV) GetMetadataColumns() []*MetadataColumn {
	return option.MetadataColumns
}

type MetadataColumn struct {
	Index   uint32 `json:"index"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Comment string `json:"comment"`
	Size    uint64 `json:"size"`
}

func (mc *MetadataColumn) GetIndex() uint32   { return mc.Index }
func (mc *MetadataColumn) GetName() string    { return mc.Name }
func (mc *MetadataColumn) GetType() string    { return mc.Type }
func (mc *MetadataColumn) GetComment() string { return mc.Comment }
func (mc *MetadataColumn) GetSize() uint64    { return mc.Size }

/**
{
    "originId": "d9b41e7138544c63f9fe25f6aa4983819793e5b46f14652a1ff1b51f99f71783",
    "dirPath": "/home/user1/data/data_root/",
	"childs": [
		{
    		"originId": "eefff343533377...4433dfaa",
    		"dirPath": "/home/user1/data/data_root/result_file/",
			"childs": [],
			"last": true,
			"filePaths": ["/home/user1/data/data_root/result_file/task_20220218_result.csv"]
		}
	],
	"last": false,
	"filePaths": ["/home/user1/data/data_root/bank_predict_partyA_20220218-090241.csv"]
}
*/
// libtypes.OrigindataType_DIR |
type MetadataOptionDIR struct {
	OriginId  string               `json:"originId"`
	DirPath   string               `json:"dirPath"`
	Childs    []*MetadataOptionDIR `json:"childs"`
	Last      bool                 `json:"last"`
	FilePaths []string             `json:"filePaths"`
}

func (option *MetadataOptionDIR) GetOriginId() string             { return option.OriginId }
func (option *MetadataOptionDIR) GetDirPath() string              { return option.DirPath }
func (option *MetadataOptionDIR) GetChilds() []*MetadataOptionDIR { return option.Childs }
func (option *MetadataOptionDIR) GetLast() bool                   { return option.Last }
func (option *MetadataOptionDIR) GetFilePaths() []string          { return option.FilePaths }

/**
{
    "originId": "d9b41e7138544c63f9fe25f6aa4983819793e5b46f14652a1ff1b51f99f71783",
    "dataPath": "/home/user1/data/data_root/bank_predict_partyA_20220218-090241.csv",
    "size": 12,
}
*/
// libtypes.OrigindataType_BINARY |
type MetadataOptionBINARY struct {
	OriginId string `json:"originId"`
	DataPath string `json:"dataPath"`
	Size     uint64 `json:"size"`
}

func (option *MetadataOptionBINARY) GetOriginId() string { return option.OriginId }
func (option *MetadataOptionBINARY) GetDataPath() string { return option.DataPath }
func (option *MetadataOptionBINARY) GetSize() uint64     { return option.Size }
