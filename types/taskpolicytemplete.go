package types

import (
	"encoding/json"
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

func FetchMetedataIdByPartyId (partyId string, policyType uint32, policyOption string) (string, error) {
	switch policyType {
	case TASK_METADATA_POLICY_ROW_COLUMN:
		var policys []*TaskMetadataPolicyRowAndColumn
		if err := json.Unmarshal([]byte(policyOption), &policys); nil != err {
			return "", err
		}
		for _, policy := range policys {
			if policy.QueryPartyId() == partyId {
				return policy.QueryMetadataId(), nil
			}
		}
	}
	return "", NotFoundMetadataPolicy
}
func FetchMetedataNameByPartyId (partyId string, policyType uint32, policyOption string) (string, error) {
	switch policyType {
	case TASK_METADATA_POLICY_ROW_COLUMN:
		var policys []*TaskMetadataPolicyRowAndColumn
		if err := json.Unmarshal([]byte(policyOption), &policys); nil != err {
			return "", err
		}
		for _, policy := range policys {
			if policy.QueryPartyId() == partyId {
				return policy.QueryMetadataName(), nil
			}
		}
	}
	return "", NotFoundMetadataPolicy
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
