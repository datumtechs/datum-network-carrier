package policy

import (
	"encoding/json"
	"fmt"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/types"
)

func FetchOriginId (dataType commonconstantpb.OrigindataType, metadataOption string) (string, error) {

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

func FetchDataPath (dataType commonconstantpb.OrigindataType, metadataOption string) (string, error) {

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