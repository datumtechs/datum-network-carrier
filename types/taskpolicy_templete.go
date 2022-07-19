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
	TASK_DATA_POLICY_UNKNOWN = 0

	TASK_DATA_POLICY_CSV = 1 // csv
	/**
	"{
		"partyId": "p0",
		"metadataId": "metadata:0xf7396b9a6be9c20...c54880c2d",
		"metadataName": "aaa",
		"inputType": 1, // 输入数据的类型，0:unknown, 1:origin_data, 2:psi_output, 3:model
		"keyColumn": 1,
		"selectedColumns": [1, 2, 3]
	}"
	*/
	TASK_DATA_POLICY_DIR = 2 // dir
	/**
	"{
		"partyId": "p0",
		"metadataId": "metadata:0xf7396b9a6be9c20...c54880c2d",
		"metadataName": "aaa",
		"inputType": 3 // 输入数据的类型，0:unknown, 1:origin_data, 2:psi_output, 3:model
	}"
	*/
	TASK_DATA_POLICY_BINARY = 3 // binary
	/**
	"{
		"partyId": "p0",
		"metadataId": "metadata:0xf7396b9a6be9c20...c54880c2d",
		"metadataName": "aaa",
		"inputType": 1 // 输入数据的类型，0:unknown, 1:origin_data, 2:psi_output, 3:model
	}"
	*/
	TASK_DATA_POLICY_XLS  = 4 // xls
	TASK_DATA_POLICY_XLSX = 5 // xlsx
	TASK_DATA_POLICY_TXT  = 6 // txt
	TASK_DATA_POLICY_JSON = 7 // json

	TASK_DATA_POLICY_CSV_WITH_TASKRESULTDATA = 30001 // csv for task result data
	/*
	"{
		"partyId": "p0",
		"taskId": "task:0x43b3d8c65b877adfd05a77dc6b3bb1ad27e4727edbccb3cc76ffd51f78794479",
		"inputType": 1, // 输入数据的类型，0:unknown, 1:origin_data, 2:psi_output, 3:model
		"keyColumnName": "id",
		"selectedColumnNames": ["name", "age", "point"]
	}"
	*/
	TASK_DATA_POLICY_IS_CSV_HAVE_CONSUME = 40001
	/*
		"{
			"partyId": "p0",
			"metadataId": "metadata:0xf7396b9a6be9c20...c54880c2d",
			"metadataName": "aaa",
			"inputType": 1, // 输入数据的类型，0:unknown, 1:origin_data, 2:psi_output, 3:model
			"keyColumn": 1,
			"selectedColumns": [1, 2, 3],
			"consumeTypes": [],  // 消费该元数据的方式类型说明, 0: unknown, 1: metadataAuth, 2: ERC20, 3: ERC721, ...
			"consumeOptions": ["具体看消费说明option定义的字符串"]
		}"
	*/
	// ==================================================================== power policy option ====================================================================

	TASK_POWER_POLICY_UNKNOWN = 0

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
	TASK_RECEIVER_POLICY_RANDOM_UNKNOWN = 0

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
	TASK_DATAFLOW_POLICY_GENERAL_FULL_UNKNOWN = 0

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
TASK_DATA_POLICY_UNKNOWN
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
TASK_DATA_POLICY_CSV
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
type TaskMetadataPolicyCsvConsume struct {
	TaskMetadataPolicyCSV
	ConsumeTypes   []uint8  `json:"consumeTypes"`
	ConsumeOptions []string `json:"consumeOptions"`
}

type DataConsumePolicy interface {
	Address() string
}
type MetadataAuthConsume string

type Tk20Consume struct {
	Contract string `json:"contract"`
	Balance  uint64 `json:"balance"`
}

type Tk721Consume struct {
	Contract string `json:"contract"`
	TokenId  string `json:"tokenId"`
}

func (tk *Tk721Consume) Address() string {
	return tk.Contract
}
func (tk *Tk20Consume) Address() string {
	return tk.Contract
}
func (tk MetadataAuthConsume) Address() string {
	return string(tk)
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
func (p *TaskMetadataPolicyCsvConsume) GetConsumeTypes() []uint8 {
	return p.ConsumeTypes
}
func (p *TaskMetadataPolicyCsvConsume) GetConsumeOptions() []string {
	return p.ConsumeOptions
}
/**
TASK_DATA_POLICY_DIR
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
TASK_DATA_POLICY_BINARY
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

/**
TASK_DATA_POLICY_CSV_WITH_TASKRESULTDATA
value: 1
example:

			{
				"partyId": "p0",
				"taskId": "task:0x43b3d8c65b877adfd05a77dc6b3bb1ad27e4727edbccb3cc76ffd51f78794479",
				"inputType": 1, // 输入数据的类型，0:unknown, 1:origin_data, 2:psi_output, 3:model
				"keyColumnName": "id",
				"selectedColumnNames": ["name", "age", "point"]
			}


*/
type TaskMetadataPolicyCSVWithTaskResultData struct {
	PartyId             string   `json:"partyId"`
	TaskId              string   `json:"taskId"`
	InputType           uint32   `json:"inputType"`
	KeyColumnName       string   `json:"keyColumnName"`
	SelectedColumnNames []string `json:"selectedColumnNames"`
}

func (p *TaskMetadataPolicyCSVWithTaskResultData) GetPartyId() string {
	return p.PartyId
}
func (p *TaskMetadataPolicyCSVWithTaskResultData) GetTaskId() string {
	return p.TaskId
}
func (p *TaskMetadataPolicyCSVWithTaskResultData) QueryInputType() uint32 {
	return p.InputType
}
func (p *TaskMetadataPolicyCSVWithTaskResultData) QueryKeyColumnName() string {
	return p.KeyColumnName
}
func (p *TaskMetadataPolicyCSVWithTaskResultData) QuerySelectedColumnNames() []string {
	return p.SelectedColumnNames
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
func (p *TaskPowerPolicyDataNodeProvide) GetPowerPartyId() string    { return p.PowerPartyId }

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
