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
	// ==================================================================== metadata policy option ====================================================================
	/**
	[
		{
	        "partyId": "p0",
	        "metadataId": "metadata:0xf7396b9a6be9ca6cd2099f0d5e67a748643898f2735573c0da44fc54880cf52d",
			"metadataName" "aaa",
	        "keyColumn": 1,
	        "selectedColumns": [1, 2, 3]
	    },
		{
	        "partyId": "p1",
	        "metadataId": "metadata:0x2843e8103c065268991943c56d970570bedc8b31066e05ef7d62ea20b6c537c7",
			"metadataName" "bbb",
	        "keyColumn": 1,
	        "selectedColumns": [1, 2, 3]
	    }
	]
	*/
	TASK_METADATA_POLICY_ROW_COLUMN = 1

	// ==================================================================== power policy option ====================================================================
	/**
	["q0", "q1", "q2"]
	*/
	TASK_POWER_POLICY_ASSIGNMENT_SYMBOL_RANDOM_ELECTION_POWER = 1
	/**
	["p0", "p1"]
	*/
	TASK_POWER_POLICY_DATANODE_PROVIDE_POWER = 2

	// ==================================================================== dataFlow policy option ====================================================================

)

// ==================================================================== metadata policy option ====================================================================

/**
enum:  TASK_METADATA_POLICY_ROW_COLUMN
value: 1
example:

	{
        "partyId": "p0",
        "metadataId": "metadata:0xf7396b9a6be9ca6cd2099f0d5e67a748643898f2735573c0da44fc54880cf52d",
		"metadataName" "xxx",
        "keyColumn": 1,
        "selectedColumns": [1, 2, 3]
    }


*/
type TaskMetadataPolicyRowAndColumn struct {
	PartyId         string   `json:"partyId"`
	MetadataId      string   `json:"metadataId"`
	MetadataName    string   `json:"metadataName"`
	KeyColumn       uint32   `json:"keyColumn"`
	SelectedColumns []uint32 `json:"selectedColumns"`
}

func (p *TaskMetadataPolicyRowAndColumn) GetPartyId() string {
	return p.PartyId
}
func (p *TaskMetadataPolicyRowAndColumn) GetMetadataId() string {
	return p.MetadataId
}
func (p *TaskMetadataPolicyRowAndColumn) GetMetadataName() string {
	return p.MetadataName
}
func (p *TaskMetadataPolicyRowAndColumn) QueryKeyColumn() uint32 {
	return p.KeyColumn
}
func (p *TaskMetadataPolicyRowAndColumn) QuerySelectedColumns() []uint32 {
	return p.SelectedColumns
}

// ==================================================================== power policy option ====================================================================

//type AssignmentProvidePower struct {
//	IdentityId string `json:"identityId"`
//	PartyId    string `json:"partyId"`
//}
//
//func (p *AssignmentProvidePower) GetIdentityId() string {
//	return p.IdentityId
//}
//func (p *AssignmentProvidePower) GetPartyId() string {
//	return p.PartyId
//}

// ==================================================================== dataFlow policy option ====================================================================
