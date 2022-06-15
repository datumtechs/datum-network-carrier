package policy

import (
	"encoding/json"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/carrierdb/rawdb"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/types"
)


func (pe *PolicyEngine) FetchMetedataIdByPartyIdFromDataPolicy(partyId string, policyTypes []uint32, policyOptions []string) (string, error) {

	if len(policyTypes) != len(policyOptions) {
		return "", fmt.Errorf("type and option count is not same, types %d: policys: %d", len(policyTypes), len(policyOptions))
	}

	for i, policy := range policyOptions {
		metadataId, err := pe.FetchMetedataIdByPartyIdAndOptionFromDataPolicy(partyId, policyTypes[i], policy)
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

func (pe *PolicyEngine) FetchMetedataIdByPartyIdAndOptionFromDataPolicy(partyId string, policyType uint32, policyOption string) (string, error) {
	switch commonconstantpb.OrigindataType(policyType) {
	case types.TASK_DATA_POLICY_CSV:
		var policy *types.TaskMetadataPolicyCSV
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}
		if policy.GetPartyId() == partyId {
			return policy.GetMetadataId(), nil
		}
	case types.TASK_DATA_POLICY_DIR:
		var policy *types.TaskMetadataPolicyDIR
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}
		if policy.GetPartyId() == partyId {
			return policy.GetMetadataId(), nil
		}
	case types.TASK_DATA_POLICY_BINARY:
		var policy *types.TaskMetadataPolicyBINARY
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}
		if policy.GetPartyId() == partyId {
			return policy.GetMetadataId(), nil
		}
	case types.TASK_DATA_POLICY_CSV_WITH_TASKRESULTDATA:
		var policy *types.TaskMetadataPolicyCSVWithTaskResultData
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}
		if policy.GetPartyId() == partyId {
			resultData, err := pe.db.QueryTaskUpResulData(policy.GetTaskId())
			if rawdb.IsNoDBNotFoundErr(err) {
				return "", err
			}
			if rawdb.IsDBNotFoundErr(err) {
				return IgnoreMetadataId, nil
			}
			return resultData.GetMetadataId(), nil
		}
	}
	return "", types.NotFoundMetadataPolicy
}

func (pe *PolicyEngine) FetchAllMetedataIdsFromDataPolicy(policyTypes []uint32, policyOptions []string) ([]string, error) {

	if len(policyTypes) != len(policyOptions) {
		return nil, fmt.Errorf("type and option count is not same, types %d: policys: %d", len(policyTypes), len(policyOptions))
	}

	metadataIds := make([]string, len(policyOptions))

	for i, policyOption := range policyOptions {

		metadataId, err := pe.FetchAllMetedataIdByOptionFromDataPolicy(policyTypes[i], policyOption)
		if nil != err {
			return nil, err
		}
		metadataIds[i] = metadataId
	}

	return metadataIds, nil
}

func (pe *PolicyEngine) FetchAllMetedataIdByOptionFromDataPolicy(policyType uint32, policyOption string) (string, error) {
	switch commonconstantpb.OrigindataType(policyType) {
	case types.TASK_DATA_POLICY_CSV:
		var policy *types.TaskMetadataPolicyCSV
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}

		return policy.GetMetadataId(), nil
	case types.TASK_DATA_POLICY_DIR:
		var policy *types.TaskMetadataPolicyDIR
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}

		return policy.GetMetadataId(), nil
	case types.TASK_DATA_POLICY_BINARY:
		var policy *types.TaskMetadataPolicyBINARY
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}

		return policy.GetMetadataId(), nil

	case types.TASK_DATA_POLICY_CSV_WITH_TASKRESULTDATA:
		var policy *types.TaskMetadataPolicyCSVWithTaskResultData
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}
		resultData, err := pe.db.QueryTaskUpResulData(policy.GetTaskId())
		if rawdb.IsNoDBNotFoundErr(err) {
			return "", err
		}
		if rawdb.IsDBNotFoundErr(err) {
			return IgnoreMetadataId, nil
		}
		return resultData.GetMetadataId(), nil
	}
	return "", types.NotFoundMetadataPolicy
}

func (pe *PolicyEngine) FetchMetedataNameByPartyIdFromDataPolicy(partyId string, policyTypes []uint32, policyOptions []string) (string, error) {

	if len(policyTypes) != len(policyOptions) {
		return "", fmt.Errorf("type and option count is not same, types %d: policys: %d", len(policyTypes), len(policyOptions))
	}

	for i, policy := range policyOptions {
		metadataName, err := pe.FetchMetedataNameByPartyIdAndOptionFromDataPolicy(partyId, policyTypes[i], policy)
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

func (pe *PolicyEngine) FetchMetedataNameByPartyIdAndOptionFromDataPolicy(partyId string, policyType uint32, policyOption string) (string, error) {
	switch commonconstantpb.OrigindataType(policyType) {
	case types.TASK_DATA_POLICY_CSV:
		var policy *types.TaskMetadataPolicyCSV
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}
		if policy.GetPartyId() == partyId {
			return policy.GetMetadataName(), nil
		}
	case types.TASK_DATA_POLICY_DIR:
		var policy *types.TaskMetadataPolicyDIR
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}
		if policy.GetPartyId() == partyId {
			return policy.GetMetadataName(), nil
		}
	case types.TASK_DATA_POLICY_BINARY:
		var policy *types.TaskMetadataPolicyBINARY
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}
		if policy.GetPartyId() == partyId {
			return policy.GetMetadataName(), nil
		}
	case types.TASK_DATA_POLICY_CSV_WITH_TASKRESULTDATA:
		var policy *types.TaskMetadataPolicyCSVWithTaskResultData
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}
		if policy.GetPartyId() == partyId {
			resultData, err := pe.db.QueryTaskUpResulData(policy.GetTaskId())
			if rawdb.IsNoDBNotFoundErr(err) {
				return "", err
			}
			if rawdb.IsDBNotFoundErr(err) {
				return IgnoreMetadataId, nil
			}
			metadata, err := pe.db.QueryInternalMetadataById(resultData.GetMetadataId())
			if rawdb.IsNoDBNotFoundErr(err) {
				return "", err
			}
			if rawdb.IsDBNotFoundErr(err) {
				return IgnoreMetadataName, nil
			}
			return metadata.GetData().GetMetadataName(), nil
		}
	}
	return "", types.NotFoundMetadataPolicy
}


// ==============================================================  power policy ==============================================================

func (pe *PolicyEngine) FetchPowerPartyIdsFromPowerPolicy(policyTypes []uint32, policyOptions []string) ([]string, error) {

	if len(policyTypes) != len(policyOptions) {
		return nil, fmt.Errorf("type and option count is not same, types %d: policys: %d", len(policyTypes), len(policyOptions))
	}

	partyIds := make([]string, len(policyOptions))

	for i, policyOption := range policyOptions {

		partyId, err := pe.FetchPowerPartyIdByOptionFromPowerPolicy(policyTypes[i], policyOption)
		if nil != err {
			return nil, err
		}
		partyIds[i] = partyId
	}

	return partyIds, nil
}

func (pe *PolicyEngine) FetchPowerPartyIdByOptionFromPowerPolicy(policyType uint32, policyOption string) (string, error) {
	switch policyType {
	case types.TASK_POWER_POLICY_ASSIGNMENT_SYMBOL_RANDOM_ELECTION:
		return policyOption, nil
	case types.TASK_POWER_POLICY_DATANODE_PROVIDE:
		var policy *types.TaskPowerPolicyDataNodeProvide
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}
		return policy.GetPowerPartyId(), nil
	default:
		return "", types.NotFoundPowerPolicy
	}
}


// ==============================================================  receiver policy ==============================================================
func (pe *PolicyEngine) FetchReceiverPartyIdsFromReceiverPolicy(policyTypes []uint32, policyOptions []string) ([]string, error) {

	if len(policyTypes) != len(policyOptions) {
		return nil, fmt.Errorf("type and option count is not same, types %d: policys: %d", len(policyTypes), len(policyOptions))
	}

	partyIds := make([]string, len(policyOptions))

	for i, policyOption := range policyOptions {

		partyId, err := pe.FetchReceiverPartyIdByOptionFromReceiverPolicy(policyTypes[i], policyOption)
		if nil != err {
			return nil, err
		}
		partyIds[i] = partyId
	}

	return partyIds, nil
}

func (pe *PolicyEngine) FetchReceiverPartyIdByOptionFromReceiverPolicy(policyType uint32, policyOption string) (string, error) {
	switch policyType {
	case types.TASK_RECEIVER_POLICY_RANDOM_ELECTION:
		return policyOption, nil
	case types.TASK_RECEIVER_POLICY_DATANODE_PROVIDE:
		var policy *types.TaskReceiverPolicyDataNodeProvide
		if err := json.Unmarshal([]byte(policyOption), &policy); nil != err {
			return "", err
		}
		return policy.GetReceiverPartyId(), nil
	default:
		return "", types.NotFoundPowerPolicy
	}
}
