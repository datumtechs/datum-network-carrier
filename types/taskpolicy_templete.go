package types

import (
	"fmt"
)

var (
	NotFoundMetadataPolicy = fmt.Errorf("not found metadata policy")
)

type TaskMetadataPolicy interface {
	QueryPartyId() string
	QueryMetadataId() string
	QueryMetadataName() string
}


const (
	TASK_METADATA_POLICY_ROW_COLUMN = 1
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
