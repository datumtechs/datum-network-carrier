package policy

import (
	"encoding/json"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func FetchMetedataIdByPartyId (partyId string, policyType uint32, policyOption string) (string, error) {
	switch policyType {
	case types.TASK_METADATA_POLICY_ROW_COLUMN:
		var policys []*types.TaskMetadataPolicyRowAndColumn
		if err := json.Unmarshal([]byte(policyOption), &policys); nil != err {
			return "", err
		}
		for _, policy := range policys {
			if policy.QueryPartyId() == partyId {
				return policy.QueryMetadataId(), nil
			}
		}
	}
	return "", types.NotFoundMetadataPolicy
}
func FetchMetedataNameByPartyId (partyId string, policyType uint32, policyOption string) (string, error) {
	switch policyType {
	case types.TASK_METADATA_POLICY_ROW_COLUMN:
		var policys []*types.TaskMetadataPolicyRowAndColumn
		if err := json.Unmarshal([]byte(policyOption), &policys); nil != err {
			return "", err
		}
		for _, policy := range policys {
			if policy.QueryPartyId() == partyId {
				return policy.QueryMetadataName(), nil
			}
		}
	}
	return "", types.NotFoundMetadataPolicy
}

func FetchPowerPartyIds(policyType uint32, policyOption string) ([]string, error) {
	switch policyType {
	case types.TASK_POWER_POLICY_ASSIGNMENT_LABEL:
		var policys []string
		if err := json.Unmarshal([]byte(policyOption), &policys); nil != err {
			return nil, err
		}
		return policys, nil
	}
	return nil, types.NotFoundPowerPolicy
}

