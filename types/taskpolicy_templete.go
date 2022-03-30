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

//type TaskMetadataPolicy interface {
//	GetPartyId() string
//	GetMetadataId() string
//	GetMetadataName() string
//}

const (
	// data policy
	TASK_METADATA_POLICY_ROW_COLUMN = 1

	// power policy
	TASK_POWER_POLICY_ASSIGNMENT_SYMBOL_RANDOM_ELECTION_POWER = 1
	TASK_POWER_POLICY_DATANODE_PROVIDE_POWER                  = 2

	// dataFlow policy

)

// ==================================================================== metadata policy option ====================================================================

/**
TASK_METADATA_POLICY_ROW_COLUMN
*/
type TaskMetadataPolicyRowAndColumn struct {
	PartyId         string
	MetadataId      string
	MetadataName    string
	KeyColumn       uint32
	SelectedColumns []uint32
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

type AssignmentProvidePower struct {
	IdentityId string
	PartyId    string
}
func (p *AssignmentProvidePower) GetIdentityId() string {
	return p.IdentityId
}
func (p *AssignmentProvidePower) GetPartyId() string {
	return p.PartyId
}

// ==================================================================== dataFlow policy option ====================================================================
