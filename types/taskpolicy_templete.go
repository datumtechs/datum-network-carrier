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
	TASK_POWER_POLICY_ASSIGNMENT_LABEL = 1

	// dataFlow policy


)

/**

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
func (p *TaskMetadataPolicyRowAndColumn) QueryKeyColumn () uint32 {
	return p.KeyColumn
}
func (p *TaskMetadataPolicyRowAndColumn) QuerySelectedColumns () []uint32 {
	return p.SelectedColumns
}


