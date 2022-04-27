package types

import (
	"fmt"
)

var (
	// about data policy
	NotFoundMetadataPolicy = fmt.Errorf("not found metadata policy")
	// about power policy
	NotFoundPowerPolicy = fmt.Errorf("not found power policy")
)

const (
	// ==================================================================== metadata policy option ====================================================================\

	// ==================================================================== power policy option ====================================================================
	/**
	"q0"
	*/
	TASK_POWER_POLICY_ASSIGNMENT_SYMBOL_RANDOM_ELECTION_POWER = 1
	/**
	"p0"
	*/
	TASK_POWER_POLICY_DATANODE_PROVIDE_POWER = 2

	// ==================================================================== dataFlow policy option ====================================================================

)

// ==================================================================== metadata policy option ====================================================================

/**
OrigindataType_Unknown
value: 0
example:

			{
				"partyId": "p0",
				"metadataId": "metadata:0x2843e8103c...6c537c7",
				"metadataName": "bbb",
				"inputType": 3 // 输入数据的类型，0:unknown, 1:origin_data, 2:psi_output, 3:model
			}


*/
type TaskMetadataPolicyUnknown struct {
	PartyId      string `json:"partyId"`
	MetadataId   string `json:"metadataId"`
	MetadataName string `json:"metadataName"`
	InputType    uint32 `json:"inputType"`
}

func (p *TaskMetadataPolicyUnknown) GetPartyId() string {
	return p.PartyId
}
func (p *TaskMetadataPolicyUnknown) GetMetadataId() string {
	return p.MetadataId
}
func (p *TaskMetadataPolicyUnknown) GetMetadataName() string {
	return p.MetadataName
}
func (p *TaskMetadataPolicyUnknown) QueryInputType() uint32 {
	return p.InputType
}


/**
OrigindataType_CSV
value: 1
example:

			{
				"partyId": "p0",
				"metadataId": "metadata:0xf7396b9a6be9c20...c54880c2d",
				"metadataName": "aaa",
				"inputType": 1, // 输入数据的类型，0:unknown, 1:origin_data, 2:psi_output, 3:model
				"keyColumn": 1,
				"selectedColumns": [1, 2, 3]
			}


*/
type TaskMetadataPolicyCSV struct {
	PartyId         string   `json:"partyId"`
	MetadataId      string   `json:"metadataId"`
	MetadataName    string   `json:"metadataName"`
	InputType       uint32   `json:"inputType"`
	KeyColumn       uint32   `json:"keyColumn"`
	SelectedColumns []uint32 `json:"selectedColumns"`
}

func (p *TaskMetadataPolicyCSV) GetPartyId() string {
	return p.PartyId
}
func (p *TaskMetadataPolicyCSV) GetMetadataId() string {
	return p.MetadataId
}
func (p *TaskMetadataPolicyCSV) GetMetadataName() string {
	return p.MetadataName
}
func (p *TaskMetadataPolicyCSV) QueryInputType() uint32 {
	return p.InputType
}
func (p *TaskMetadataPolicyCSV) QueryKeyColumn() uint32 {
	return p.KeyColumn
}
func (p *TaskMetadataPolicyCSV) QuerySelectedColumns() []uint32 {
	return p.SelectedColumns
}


/**
OrigindataType_BINARY
value: 0
example:

			{
				"partyId": "p0",
				"metadataId": "metadata:0x2843e8103c...6c537c7",
				"metadataName": "bbb",
				"inputType": 3 // 输入数据的类型，0:unknown, 1:origin_data, 2:psi_output, 3:model
			}


*/
type TaskMetadataPolicyBINARY struct {
	PartyId      string `json:"partyId"`
	MetadataId   string `json:"metadataId"`
	MetadataName string `json:"metadataName"`
	InputType    uint32 `json:"inputType"`
}

func (p *TaskMetadataPolicyBINARY) GetPartyId() string {
	return p.PartyId
}
func (p *TaskMetadataPolicyBINARY) GetMetadataId() string {
	return p.MetadataId
}
func (p *TaskMetadataPolicyBINARY) GetMetadataName() string {
	return p.MetadataName
}
func (p *TaskMetadataPolicyBINARY) QueryInputType() uint32 {
	return p.InputType
}


// ==================================================================== power policy option ====================================================================

// ==================================================================== dataFlow policy option ====================================================================
