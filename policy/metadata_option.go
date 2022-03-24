package policy

import (
	"encoding/json"
	"fmt"
	libcommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func FetchOriginId (fileType libcommonpb.OriginFileType, metadataOption string) (string, error) {
	if fileType == libcommonpb.OriginFileType_FileType_CSV {
		var option *types.MetadataOptionRowAndColumn
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return "", fmt.Errorf("unmashal metadataOption failed, %s", err)
		}
		return option.GetOriginId(), nil
	}
	return "", types.CannotMatchMetadataOption
}

func FetchFilePath (fileType libcommonpb.OriginFileType, metadataOption string) (string, error) {
	if fileType == libcommonpb.OriginFileType_FileType_CSV {
		var option *types.MetadataOptionRowAndColumn
		if err := json.Unmarshal([]byte(metadataOption), &option); nil != err {
			return "", fmt.Errorf("unmashal metadataOption failed, %s", err)
		}
		return option.GetFilePath(), nil
	}
	return "", types.CannotMatchMetadataOption
}