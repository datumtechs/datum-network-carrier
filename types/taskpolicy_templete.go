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

	TASK_POWER_POLICY_ASSIGNMENT_SYMBOL_RANDOM_ELECTION = 1
	/**
	"q0"
	*/

	TASK_POWER_POLICY_DATANODE_PROVIDE = 2
	/**
	"{
		"providerPartyId": "p0",
		"powerPartyId": "y0"
	}"
	*/

	// ==================================================================== receiver policy option ====================================================================
	TASK_RECEIVER_POLICY_RANDOM_ELECTION = 1
	/**
	"q0"
	 */

	TASK_RECEIVER_POLICY_DATANODE_PROVIDE = 2
	/**
	{
		"providerPartyId": "p0",
		"receiverPartyId": "q0"
	}
	*/

	// ==================================================================== dataFlow policy option ====================================================================
	TASK_DATAFLOW_POLICY_GENERAL_FULL_CONNECT = 1
	/**
	{}
	 */

	TASK_DATAFLOW_POLICY_GENERAL_DIRECTIONAL_CONNECT = 2
	/**
	"{
		"p0": ["p1", "y0", "y1"],
		"p1": ["p0", "y0", "y1"],
		"y0": ["q0", "q1"],
		"y1": ["q0", "q1"]
	}"
	 */
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
OrigindataType_DIR
value: 2
example:

			{
				"partyId": "p0",
				"metadataId": "metadata:0x2843e8103c...6c537c7",
				"metadataName": "bbb",
				"inputType": 3 // 输入数据的类型，0:unknown, 1:origin_data, 2:psi_output, 3:model
			}


*/
type TaskMetadataPolicyDIR struct {
	PartyId      string `json:"partyId"`
	MetadataId   string `json:"metadataId"`
	MetadataName string `json:"metadataName"`
	InputType    uint32 `json:"inputType"`
}

func (p *TaskMetadataPolicyDIR) GetPartyId() string {
	return p.PartyId
}
func (p *TaskMetadataPolicyDIR) GetMetadataId() string {
	return p.MetadataId
}
func (p *TaskMetadataPolicyDIR) GetMetadataName() string {
	return p.MetadataName
}
func (p *TaskMetadataPolicyDIR) QueryInputType() uint32 {
	return p.InputType
}

/**
OrigindataType_BINARY
value: 3
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

/**
{
	"providerPartyId": "p0",
	"powerPartyId": "y0"
}
*/
type TaskPowerPolicyDataNodeProvide struct {
	ProviderPartyId string `json:"providerPartyId"`
	PowerPartyId    string `json:"powerPartyId"`
}

func (p *TaskPowerPolicyDataNodeProvide) GetProviderPartyId() string { return p.ProviderPartyId }
func (p *TaskPowerPolicyDataNodeProvide) GetPowerPartyId() string { return p.PowerPartyId }

// ==================================================================== receiver policy option ====================================================================

/**
{
	"providerPartyId": "p0",
	"receiverPartyId": "q0"
}
*/
type TaskReceiverPolicyDataNodeProvide struct {
	ProviderPartyId string `json:"providerPartyId"`
	ReceiverPartyId string `json:"receiverPartyId"`
}

func (p *TaskReceiverPolicyDataNodeProvide) GetProviderPartyId() string { return p.ProviderPartyId }
func (p *TaskReceiverPolicyDataNodeProvide) GetReceiverPartyId() string { return p.ReceiverPartyId }

// ==================================================================== dataFlow policy option ====================================================================
