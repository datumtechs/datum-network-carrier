package policy

import (
	"encoding/json"
	"fmt"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func FetchOriginId (fileType apicommonpb.OriginFileType, metadataOption string) (string, error) {
	if fileType == apicommonpb.OriginFileType_FileType_CSV {
		var option *types.MetadataOptionRowAndColumn
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return "", fmt.Errorf("unmashal metadataOption failed, %s", err)
		}
		return option.GetOriginId(), nil
	}
	return "", types.CannotMatchMetadataOption
}

func FetchFilePath (fileType apicommonpb.OriginFileType, metadataOption string) (string, error) {
	if fileType == apicommonpb.OriginFileType_FileType_CSV {
		var option *types.MetadataOptionRowAndColumn
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return "", fmt.Errorf("unmashal metadataOption failed, %s", err)
		}
		return option.GetFilePath(), nil
	}
	return "", types.CannotMatchMetadataOption
}