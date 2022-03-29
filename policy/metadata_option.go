package policy

import (
	"encoding/json"
	"fmt"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func FetchOriginId (fileType libtypes.OrigindataType, metadataOption string) (string, error) {
	if fileType == libtypes.OrigindataType_OrigindataType_CSV {
		var option *types.MetadataOptionRowAndColumn
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return "", fmt.Errorf("unmashal metadataOption failed, %s", err)
		}
		return option.GetOriginId(), nil
	}
	return "", types.CannotMatchMetadataOption
}

func FetchFilePath (fileType libtypes.OrigindataType, metadataOption string) (string, error) {
	if fileType == libtypes.OrigindataType_OrigindataType_CSV {
		var option *types.MetadataOptionRowAndColumn
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return "", fmt.Errorf("unmashal metadataOption failed, %s", err)
		}
		return option.GetFilePath(), nil
	}
	return "", types.CannotMatchMetadataOption
}