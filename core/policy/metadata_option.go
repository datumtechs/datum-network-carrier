package policy

import (
	"encoding/json"
	"fmt"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/types"
)

func (pe *PolicyEngine) FetchOriginId(dataType commonconstantpb.OrigindataType, metadataOption string) (string, error) {

	switch dataType {
	case commonconstantpb.OrigindataType_OrigindataType_CSV:
		var option *types.MetadataOptionCSV
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return "", fmt.Errorf("unmashal metadataOption to csv failed, %s", err)
		}
		return option.GetOriginId(), nil
	case commonconstantpb.OrigindataType_OrigindataType_DIR:
		var option *types.MetadataOptionDIR
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return "", fmt.Errorf("unmashal metadataOption to dir failed, %s", err)
		}
		return option.GetOriginId(), nil
	case commonconstantpb.OrigindataType_OrigindataType_BINARY:
		var option *types.MetadataOptionBINARY
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return "", fmt.Errorf("unmashal metadataOption to binary failed, %s", err)
		}
		return option.GetOriginId(), nil
	default:
		return "", types.CannotMatchMetadataOption
	}
}

func (pe *PolicyEngine) FetchDataPath(dataType commonconstantpb.OrigindataType, metadataOption string) (string, error) {

	switch dataType {
	case commonconstantpb.OrigindataType_OrigindataType_CSV:
		var option *types.MetadataOptionCSV
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return "", fmt.Errorf("unmashal metadataOption to csv failed, %s", err)
		}
		return option.GetDataPath(), nil
	case commonconstantpb.OrigindataType_OrigindataType_DIR:
		var option *types.MetadataOptionDIR
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return "", fmt.Errorf("unmashal metadataOption to dir failed, %s", err)
		}
		return option.GetDirPath(), nil
	case commonconstantpb.OrigindataType_OrigindataType_BINARY:
		var option *types.MetadataOptionBINARY
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return "", fmt.Errorf("unmashal metadataOption to binary failed, %s", err)
		}
		return option.GetDataPath(), nil
	default:
		return "", types.CannotMatchMetadataOption
	}
}

func (pe *PolicyEngine) FetchConsumeTypes(dataType commonconstantpb.OrigindataType, metadataOption string) ([]uint8, error) {

	switch dataType {
	case commonconstantpb.OrigindataType_OrigindataType_CSV:
		var option *types.MetadataOptionCSV
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return nil, fmt.Errorf("unmashal metadataOption to csv failed, %s", err)
		}
		return option.GetConsumeTypes(), nil
	case commonconstantpb.OrigindataType_OrigindataType_DIR:
		var option *types.MetadataOptionDIR
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return nil, fmt.Errorf("unmashal metadataOption to dir failed, %s", err)
		}
		return option.GetConsumeTypes(), nil
	case commonconstantpb.OrigindataType_OrigindataType_BINARY:
		var option *types.MetadataOptionBINARY
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return nil, fmt.Errorf("unmashal metadataOption to binary failed, %s", err)
		}
		return option.GetConsumeTypes(), nil
	default:
		return nil, types.CannotMatchMetadataOption
	}
}

func (pe *PolicyEngine) FetchConsumeOptions(dataType commonconstantpb.OrigindataType, metadataOption string) ([]string, error) {

	switch dataType {
	case commonconstantpb.OrigindataType_OrigindataType_CSV:
		var option *types.MetadataOptionCSV
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return nil, fmt.Errorf("unmashal metadataOption to csv failed, %s", err)
		}
		return option.GetConsumeOptions(), nil
	case commonconstantpb.OrigindataType_OrigindataType_DIR:
		var option *types.MetadataOptionDIR
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return nil, fmt.Errorf("unmashal metadataOption to dir failed, %s", err)
		}
		return option.GetConsumeOptions(), nil
	case commonconstantpb.OrigindataType_OrigindataType_BINARY:
		var option *types.MetadataOptionBINARY
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return nil, fmt.Errorf("unmashal metadataOption to binary failed, %s", err)
		}
		return option.GetConsumeOptions(), nil
	default:
		return nil, types.CannotMatchMetadataOption
	}
}
