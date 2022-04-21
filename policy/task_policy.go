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
			if policy.GetPartyId() == partyId {
				return policy.GetMetadataId(), nil
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
			if policy.GetPartyId() == partyId {
				return policy.GetMetadataName(), nil
			}
		}
	}
	return "", types.NotFoundMetadataPolicy
}

func FetchPowerPartyIds(policyType uint32, policyOption string) ([]string, error) {
	switch policyType {
	case types.TASK_POWER_POLICY_ASSIGNMENT_SYMBOL_RANDOM_ELECTION_POWER:
		var policys []string
		if err := json.Unmarshal([]byte(policyOption), &policys); nil != err {
			return nil, err
		}
		return policys, nil
	case types.TASK_POWER_POLICY_DATANODE_PROVIDE_POWER:
		var policys []string
		if err := json.Unmarshal([]byte(policyOption), &policys); nil != err {
			return nil, err
		}
		return policys, nil
	default:
		return nil, types.NotFoundPowerPolicy
	}
}