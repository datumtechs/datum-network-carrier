package types

import (
	"fmt"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"math/big"
)

var (
	CannotMatchMetadataOption = fmt.Errorf("cannot match metadata option")
)

const (
	ConsumeMetadataAuth = iota + 1
	ConsumeTk20
	ConsumeTk721
)

func IsNotCSVdata(dataType commonconstantpb.OrigindataType) bool { return !IsCSVdata(dataType) }
func IsCSVdata(dataType commonconstantpb.OrigindataType) bool {
	if dataType == commonconstantpb.OrigindataType_OrigindataType_CSV {
		return true
	}
	return false
}

func IsNotDIRdata(dataType commonconstantpb.OrigindataType) bool { return !IsDIRdata(dataType) }
func IsDIRdata(dataType commonconstantpb.OrigindataType) bool {
	if dataType == commonconstantpb.OrigindataType_OrigindataType_DIR {
		return true
	}
	return false
}

func IsNotBINARYdata(dataType commonconstantpb.OrigindataType) bool { return !IsBINARYdata(dataType) }
func IsBINARYdata(dataType commonconstantpb.OrigindataType) bool {
	if dataType == commonconstantpb.OrigindataType_OrigindataType_BINARY {
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
	"condition": 3,
    "metadataColumns": [
        {
            "index": 1,
            "name": "CLIENT_ID",
            "type": "string",
            "size": 0,
            "comment": ""
        }
    ],
	"consumeTypes":[1,2,3],
    "consumeOptions": [
            "[
                {
                    "status": 3
                }
            ]",
            "[
                {
                    "contract": "0xbbb...eee",
                    "cryptoAlgoConsumeUnit": 1000000,
                    "plainAlgoConsumeUnit": 1
                }
            ]",
            "["0xaaa...fff", ..., "0xbbb...eee"]"
    ]
}
*/
// carriertypespb.OrigindataType_CSV
type MetadataOptionCSV struct {
	OriginId        string            `json:"originId"`
	DataPath        string            `json:"dataPath"`
	Rows            uint64            `json:"rows"`
	Columns         uint64            `json:"columns"`
	Size            uint64            `json:"size"`
	HasTitle        bool              `json:"hasTitle"`
	Condition       uint64            `json:"condition"`
	MetadataColumns []*MetadataColumn `json:"metadataColumns"`
	ConsumeTypes    []uint8           `json:"consumeTypes"`
	ConsumeOptions  []string          `json:"consumeOptions"`
}

func (option *MetadataOptionCSV) GetOriginId() string  { return option.OriginId }
func (option *MetadataOptionCSV) GetDataPath() string  { return option.DataPath }
func (option *MetadataOptionCSV) GetRows() uint64      { return option.Rows }
func (option *MetadataOptionCSV) GetColumns() uint64   { return option.Columns }
func (option *MetadataOptionCSV) GetSize() uint64      { return option.Size }
func (option *MetadataOptionCSV) GetHasTitle() bool    { return option.HasTitle }
func (option *MetadataOptionCSV) GetCondition() uint64 { return option.Condition }
func (option *MetadataOptionCSV) GetMetadataColumns() []*MetadataColumn {
	return option.MetadataColumns
}
func (option *MetadataOptionCSV) GetConsumeTypes() []uint8    { return option.ConsumeTypes }
func (option *MetadataOptionCSV) GetConsumeOptions() []string { return option.ConsumeOptions }

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
	"condition": 3,
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
	"filePaths": ["/home/user1/data/data_root/bank_predict_partyA_20220218-090241.csv"],
	"consumeTypes":[1,2,3],
    "consumeOptions": [
            "[
                {
                    "status": 3
                }
            ]",
            "[
                {
                    "contract": "0xbbb...eee",
                    "cryptoAlgoConsumeUnit": 1000000,
                    "plainAlgoConsumeUnit": 1
                }
            ]",
            "["0xaaa...fff", ..., "0xbbb...eee"]"
    ]
}
*/
// carriertypespb.OrigindataType_DIR |
type MetadataOptionDIR struct {
	OriginId       string               `json:"originId"`
	DirPath        string               `json:"dirPath"`
	Condition      uint64               `json:"condition"`
	Childs         []*MetadataOptionDIR `json:"childs"`
	Last           bool                 `json:"last"`
	FilePaths      []string             `json:"filePaths"`
	ConsumeTypes   []uint8              `json:"consumeTypes"`
	ConsumeOptions []string             `json:"consumeOptions"`
}

func (option *MetadataOptionDIR) GetOriginId() string             { return option.OriginId }
func (option *MetadataOptionDIR) GetDirPath() string              { return option.DirPath }
func (option *MetadataOptionDIR) GetCondition() uint64            { return option.Condition }
func (option *MetadataOptionDIR) GetChilds() []*MetadataOptionDIR { return option.Childs }
func (option *MetadataOptionDIR) GetLast() bool                   { return option.Last }
func (option *MetadataOptionDIR) GetFilePaths() []string          { return option.FilePaths }
func (option *MetadataOptionDIR) GetConsumeTypes() []uint8        { return option.ConsumeTypes }
func (option *MetadataOptionDIR) GetConsumeOptions() []string     { return option.ConsumeOptions }

/**
{
    "originId": "d9b41e7138544c63f9fe25f6aa4983819793e5b46f14652a1ff1b51f99f71783",
    "dataPath": "/home/user1/data/data_root/bank_predict_partyA_20220218-090241.csv",
    "size": 12,
	"condition": 3,
    "consumeTypes":[1,2,3],
    "consumeOptions": [
            "[
                {
                    "status": 3
                }
            ]",
            "[
                {
                    "contract": "0xbbb...eee",
                    "cryptoAlgoConsumeUnit": 1000000,
                    "plainAlgoConsumeUnit": 1
                }
            ]",
            "["0xaaa...fff", ..., "0xbbb...eee"]"
    ]
}
*/
// carriertypespb.OrigindataType_BINARY |
type MetadataOptionBINARY struct {
	OriginId       string   `json:"originId"`
	DataPath       string   `json:"dataPath"`
	Size           uint64   `json:"size"`
	Condition      uint64   `json:"condition"`
	ConsumeTypes   []uint8  `json:"consumeTypes"`
	ConsumeOptions []string `json:"consumeOptions"`
}

func (option *MetadataOptionBINARY) GetOriginId() string         { return option.OriginId }
func (option *MetadataOptionBINARY) GetDataPath() string         { return option.DataPath }
func (option *MetadataOptionBINARY) GetSize() uint64             { return option.Size }
func (option *MetadataOptionBINARY) GetCondition() uint64        { return option.Condition }
func (option *MetadataOptionBINARY) GetConsumeTypes() []uint8    { return option.ConsumeTypes }
func (option *MetadataOptionBINARY) GetConsumeOptions() []string { return option.ConsumeOptions }

// ## metadata consume option

/**

metadataAuth consume kind:

{
  "status": 3
}
*/
type MetadataConsumeOptionMetadataAuth struct {
	Status uint64 `json:"status"`
}

const (
	// `Mcoma` mean `MetadataConsumeOptionMetadataAuth`
	//
	// ######  status values of MetadataConsumeOptionMetadataAuth ######
	//
	// NOTE: multiple of these values can exist together.
	//
	// 0000...0000   (only one valid metadataAuth related one metadata, that is default)
	// 0000...0001   (multi metadataAuths related one metadata)
	// 0000...0010   (has count consume kind)
	// 0000...0100   (has time expire consume kind)

	McomaStatusAuthOnlyOne       uint64 = 0
	McomaStatusAuthMulti                = iota
	McomaStatusTimesConsumeKind         = McomaStatusAuthMulti << 1
	McomaStatusPeriodConsumeKind        = McomaStatusAuthMulti << 2
)

func (o *MetadataConsumeOptionMetadataAuth) GetStatus() uint64 { return o.Status }

/**

tk20 consume kind:

{
    "contract": "0xbbb...eee",
    "cryptoAlgoConsumeUnit": 1000000,
    "plainAlgoConsumeUnit": 1
}
*/
type MetadataConsumeOptionTK20 struct {
	Contract              string   `json:"contract"`              // the tk20 contract address
	CryptoAlgoConsumeUnit *big.Int `json:"cryptoAlgoConsumeUnit"` // Pricing unit for ciphertext algorithm (number of tks, minimum unit)
	PlainAlgoConsumeUnit  *big.Int `json:"plainAlgoConsumeUnit"`  // Pricing unit for plain algorithm (number of tks, minimum unit)
}

func (o *MetadataConsumeOptionTK20) GetContract() string { return o.Contract }
func (o *MetadataConsumeOptionTK20) GetCryptoAlgoConsumeUnit() *big.Int {
	return o.CryptoAlgoConsumeUnit
}
func (o *MetadataConsumeOptionTK20) GetCryptoAlgoConsumeUnitUint64() uint64 {
	return o.CryptoAlgoConsumeUnit.Uint64()
}
func (o *MetadataConsumeOptionTK20) GetPlainAlgoConsumeUnit() *big.Int {
	return o.PlainAlgoConsumeUnit
}
func (o *MetadataConsumeOptionTK20) GetPlainAlgoConsumeUnitUint64() uint64 {
	return o.PlainAlgoConsumeUnit.Uint64()
}

/**

tk721 consume kind:

0xaaa...fff
*/
type MetadataConsumeOptionTK721 string

func (o MetadataConsumeOptionTK721) GetContract() string { return string(o) }
