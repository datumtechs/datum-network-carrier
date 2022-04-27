package policy

import (
	"encoding/json"
	"fmt"
	libtypes "github.com/Metisnetwork/Metis-Carrier/lib/types"
	"github.com/Metisnetwork/Metis-Carrier/types"
)

func FetchMetedataIdByPartyId (partyId string, policyTypes []uint32, policyOptions []string) (string, error) {

	if len(policyTypes) != len(policyOptions) {
		return "", fmt.Errorf("type and option count is not same, types %d: policys: %d", len(policyTypes), len(policyOptions))
	}

	for i, policy := range policyOptions {
		metadataId, err := FetchMetedataIdByPartyIdAndType(partyId, policyTypes[i], policy)
		if nil != err && err == types.NotFoundMetadataPolicy {
			continue
		}

		if nil != err {
			return "", err
		}
		return metadataId, nil
	}
	return "", types.NotFoundMetadataPolicy
}

func FetchMetedataIdByPartyIdAndType (partyId string, policyType uint32, policyOption string) (string, error) {
	switch policyType {
	case uint32(libtypes.OrigindataType_OrigindataType_CSV):
		var policys []*types.TaskMetadataPolicyCSV
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

func FetchAllMetedataIds (policyTypes []uint32, policyOptions []string) ([]string, error) {

	if len(policyTypes) != len(policyOptions) {
		return nil, fmt.Errorf("type and option count is not same, types %d: policys: %d", len(policyTypes), len(policyOptions))
	}

	metadataIds := make([]string, len(policyOptions))

	for i, policyOption := range policyOptions {

		metadataId, err := FetchAllMetedataIdByType(policyTypes[i], policyOption)
		if nil != err {
			return nil, err
		}
		metadataIds[i] = metadataId
	}

	return metadataIds, nil
}

func FetchAllMetedataIdByType (policyType uint32, policyOption string) (string, error) {
	switch policyType {
	case uint32(libtypes.OrigindataType_OrigindataType_CSV):
		var policy *types.TaskMetadataPolicyCSV
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}

		return policy.GetMetadataId(), nil
	}
	return "", types.NotFoundMetadataPolicy
}

func FetchMetedataNameByPartyId (partyId string, policyTypes []uint32, policyOptions []string) (string, error) {

	if len(policyTypes) != len(policyOptions) {
		return "", fmt.Errorf("type and option count is not same, types %d: policys: %d", len(policyTypes), len(policyOptions))
	}

	for i, policy := range policyOptions {
		metadataName, err := FetchMetedataNameByPartyIdAndType(partyId, policyTypes[i], policy)
		if nil != err && err == types.NotFoundMetadataPolicy {
			continue
		}

		if nil != err {
			return "", err
		}
		return metadataName, nil
	}
	return "", types.NotFoundMetadataPolicy
}

func FetchMetedataNameByPartyIdAndType (partyId string, policyType uint32, policyOption string) (string, error) {
	switch policyType {
	case uint32(libtypes.OrigindataType_OrigindataType_CSV):
		var policy *types.TaskMetadataPolicyCSV
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}
		if policy.GetPartyId() == partyId {
			return policy.GetMetadataName(), nil
		}
	}
	return "", types.NotFoundMetadataPolicy
}

func FetchPowerPartyIds (policyTypes []uint32, policyOptions []string) ([]string, error) {

	if len(policyTypes) != len(policyOptions) {
		return nil, fmt.Errorf("type and option count is not same, types %d: policys: %d", len(policyTypes), len(policyOptions))
	}

	partyIds := make([]string, len(policyOptions))

	for i, policyOption := range policyOptions {

		partyId, err := FetchPowerPartyIdByType(policyTypes[i], policyOption)
		if nil != err {
			return nil, err
		}
		partyIds[i] = partyId
	}

	return partyIds, nil
}

func FetchPowerPartyIdByType (policyType uint32, policyOption string) (string, error) {
	switch policyType {
	case types.TASK_POWER_POLICY_ASSIGNMENT_SYMBOL_RANDOM_ELECTION_POWER:
		return policyOption, nil
	case types.TASK_POWER_POLICY_DATANODE_PROVIDE_POWER:
		return policyOption, nil
	default:
		return "", types.NotFoundPowerPolicy
	}
}