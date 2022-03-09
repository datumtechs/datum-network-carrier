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

type TaskMetadataPolicy interface {
	QueryPartyId() string
	QueryMetadataId() string
	QueryMetadataName() string
}


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

func (p *TaskMetadataPolicyRowAndColumn) QueryPartyId () string {
	return p.PartyId
}
func (p *TaskMetadataPolicyRowAndColumn) QueryMetadataId () string {
	return p.MetadataId
}
func (p *TaskMetadataPolicyRowAndColumn) QueryMetadataName () string {
	return p.MetadataName
}


