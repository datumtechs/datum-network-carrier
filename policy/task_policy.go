package policy

import (
	"encoding/json"
	"fmt"
	libtypes "github.com/Metisnetwork/Metis-Carrier/lib/types"
	"github.com/Metisnetwork/Metis-Carrier/types"
)

func FetchMetedataIdByPartyIdFromDataPolicy(partyId string, policyTypes []uint32, policyOptions []string) (string, error) {

	if len(policyTypes) != len(policyOptions) {
		return "", fmt.Errorf("type and option count is not same, types %d: policys: %d", len(policyTypes), len(policyOptions))
	}

	for i, policy := range policyOptions {
		metadataId, err := FetchMetedataIdByPartyIdAndOptionFromDataPolicy(partyId, policyTypes[i], policy)
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

func FetchMetedataIdByPartyIdAndOptionFromDataPolicy(partyId string, policyType uint32, policyOption string) (string, error) {
	switch libtypes.OrigindataType(policyType) {
	case libtypes.OrigindataType_OrigindataType_CSV:
		var policy *types.TaskMetadataPolicyCSV
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}
		if policy.GetPartyId() == partyId {
			return policy.GetMetadataId(), nil
		}
	case libtypes.OrigindataType_OrigindataType_DIR:
		var policy *types.TaskMetadataPolicyDIR
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}
		if policy.GetPartyId() == partyId {
			return policy.GetMetadataId(), nil
		}
	case libtypes.OrigindataType_OrigindataType_BINARY:
		var policy *types.TaskMetadataPolicyBINARY
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}
		if policy.GetPartyId() == partyId {
			return policy.GetMetadataId(), nil
		}
	}
	return "", types.NotFoundMetadataPolicy
}

func FetchAllMetedataIdsFromDataPolicy(policyTypes []uint32, policyOptions []string) ([]string, error) {

	if len(policyTypes) != len(policyOptions) {
		return nil, fmt.Errorf("type and option count is not same, types %d: policys: %d", len(policyTypes), len(policyOptions))
	}

	metadataIds := make([]string, len(policyOptions))

	for i, policyOption := range policyOptions {

		metadataId, err := FetchAllMetedataIdByOptionFromDataPolicy(policyTypes[i], policyOption)
		if nil != err {
			return nil, err
		}
		metadataIds[i] = metadataId
	}

	return metadataIds, nil
}

func FetchAllMetedataIdByOptionFromDataPolicy(policyType uint32, policyOption string) (string, error) {
	switch libtypes.OrigindataType(policyType) {
	case libtypes.OrigindataType_OrigindataType_CSV:
		var policy *types.TaskMetadataPolicyCSV
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}

		return policy.GetMetadataId(), nil
	case libtypes.OrigindataType_OrigindataType_DIR:
		var policy *types.TaskMetadataPolicyDIR
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}

		return policy.GetMetadataId(), nil
	case libtypes.OrigindataType_OrigindataType_BINARY:
		var policy *types.TaskMetadataPolicyBINARY
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}

		return policy.GetMetadataId(), nil
	}
	return "", types.NotFoundMetadataPolicy
}

func FetchMetedataNameByPartyIdFromDataPolicy(partyId string, policyTypes []uint32, policyOptions []string) (string, error) {

	if len(policyTypes) != len(policyOptions) {
		return "", fmt.Errorf("type and option count is not same, types %d: policys: %d", len(policyTypes), len(policyOptions))
	}

	for i, policy := range policyOptions {
		metadataName, err := FetchMetedataNameByPartyIdAndOptionFromDataPolicy(partyId, policyTypes[i], policy)
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

func FetchMetedataNameByPartyIdAndOptionFromDataPolicy(partyId string, policyType uint32, policyOption string) (string, error) {
	switch libtypes.OrigindataType(policyType) {
	case libtypes.OrigindataType_OrigindataType_CSV:
		var policy *types.TaskMetadataPolicyCSV
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}
		if policy.GetPartyId() == partyId {
			return policy.GetMetadataName(), nil
		}
	case libtypes.OrigindataType_OrigindataType_DIR:
		var policy *types.TaskMetadataPolicyDIR
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}
		if policy.GetPartyId() == partyId {
			return policy.GetMetadataName(), nil
		}
	case libtypes.OrigindataType_OrigindataType_BINARY:
		var policy *types.TaskMetadataPolicyBINARY
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}
		if policy.GetPartyId() == partyId {
			return policy.GetMetadataName(), nil
		}
	}
	return "", types.NotFoundMetadataPolicy
}

func FetchPowerPartyIdsFromPowerPolicy(policyTypes []uint32, policyOptions []string) ([]string, error) {

	if len(policyTypes) != len(policyOptions) {
		return nil, fmt.Errorf("type and option count is not same, types %d: policys: %d", len(policyTypes), len(policyOptions))
	}

	partyIds := make([]string, len(policyOptions))

	for i, policyOption := range policyOptions {

		partyId, err := FetchPowerPartyIdByOptionFromDataPolicy(policyTypes[i], policyOption)
		if nil != err {
			return nil, err
		}
		partyIds[i] = partyId
	}

	return partyIds, nil
}

func FetchPowerPartyIdByOptionFromDataPolicy(policyType uint32, policyOption string) (string, error) {
	switch policyType {
	case types.TASK_POWER_POLICY_ASSIGNMENT_SYMBOL_RANDOM_ELECTION:
		return policyOption, nil
	case types.TASK_POWER_POLICY_DATANODE_PROVIDE:
		return policyOption, nil
	default:
		return "", types.NotFoundPowerPolicy
	}
}