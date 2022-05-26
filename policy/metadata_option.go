package policy

import (
	"encoding/json"
	"fmt"
	libtypes "github.com/datumtechs/datum-network-carrier/lib/types"
	"github.com/datumtechs/datum-network-carrier/types"
)

func FetchOriginId (dataType libtypes.OrigindataType, metadataOption string) (string, error) {

	switch dataType {
	case libtypes.OrigindataType_OrigindataType_CSV:
		var option *types.MetadataOptionCSV
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return "", fmt.Errorf("unmashal metadataOption to csv failed, %s", err)
		}
		return option.GetOriginId(), nil
	case libtypes.OrigindataType_OrigindataType_DIR:
		var option *types.MetadataOptionDIR
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return "", fmt.Errorf("unmashal metadataOption to dir failed, %s", err)
		}
		return option.GetOriginId(), nil
	case libtypes.OrigindataType_OrigindataType_BINARY:
		var option *types.MetadataOptionBINARY
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return "", fmt.Errorf("unmashal metadataOption to binary failed, %s", err)
		}
		return option.GetOriginId(), nil
	default:
		return "", types.CannotMatchMetadataOption
	}
}

func FetchDataPath (dataType libtypes.OrigindataType, metadataOption string) (string, error) {

	switch dataType {
	case libtypes.OrigindataType_OrigindataType_CSV:
		var option *types.MetadataOptionCSV
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return "", fmt.Errorf("unmashal metadataOption to csv failed, %s", err)
		}
		return option.GetDataPath(), nil
	case libtypes.OrigindataType_OrigindataType_DIR:
		var option *types.MetadataOptionDIR
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return "", fmt.Errorf("unmashal metadataOption to dir failed, %s", err)
		}
		return option.GetDirPath(), nil
	case libtypes.OrigindataType_OrigindataType_BINARY:
		var option *types.MetadataOptionBINARY
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return "", fmt.Errorf("unmashal metadataOption to binary failed, %s", err)
		}
		return option.GetDataPath(), nil
	default:
		return "", types.CannotMatchMetadataOption
	}
}