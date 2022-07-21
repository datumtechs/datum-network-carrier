package metadata

import (
	"encoding/json"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func mockForTestCheckConsumeOptionsParams(metadataOption string) (*carriertypespb.MetadataPB, *carrierapipb.UpdateMetadataRequest) {
	oldMetadata := &carriertypespb.MetadataPB{
		MetadataType:   1,
		DataType:       1, // this is a csv file
		DataHash:       "Hash001",
		State:          2,
		LocationType:   1,
		User:           "Tom",
		UserType:       1,
		MetadataOption: metadataOption,
	}
	req := &carrierapipb.UpdateMetadataRequest{
		Information: &carriertypespb.MetadataSummary{
			MetadataType:   1,
			DataType:       1, // this is a csv file
			DataHash:       "Hash001",
			State:          2,
			LocationType:   1,
			User:           "Tom",
			UserType:       1,
			MetadataOption: metadataOption,
		},
	}
	return oldMetadata, req
}
func TestCheckConsumeOptionsParams(t *testing.T) {
	{
		// test metadataOption Completely correct
		metadataOption := `{"originId": "originId001", "dataPath": "dataPath001", "rows": 12, "columns": 7, "size": 56, "hasTitle": true, "metadataColumns": [], "consumeTypes": [1, 2, 3], "consumeOptions": ["[]", "[{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 2}, {\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 4}, {\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF06\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 11}, {\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 7}]", "[\"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF11\", \"0x79e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\"]"]}`
		var test *types.MetadataOptionCSV
		err := json.Unmarshal([]byte(metadataOption), &test)
		assert.Nil(t, err)
		oldMetadata, req := mockForTestCheckConsumeOptionsParams(metadataOption)
		_, err = checkCanUpdateMetadataFieldIsLegal(oldMetadata, req)
		assert.Nil(t, err)
		//var info []map[string]interface{}
		//err = json.Unmarshal([]byte(test.GetConsumeOptions()[1]), &info)
		//for _, v := range info {
		//	fmt.Println(v["contract"].(string))
		//}
		//assert.Nil(t, err)
	}
	{
		// test metadataOption have repeat address
		metadataOption := `{"originId": "originId001", "dataPath": "dataPath001", "rows": 12, "columns": 7, "size": 56, "hasTitle": true, "metadataColumns": [], "consumeTypes": [1, 2, 3], "consumeOptions": ["[]", "[{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 2}, {\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 4}, {\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF06\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 11}, {\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 7}]", "[\"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"0x79e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\"]"]}`
		oldMetadata, req := mockForTestCheckConsumeOptionsParams(metadataOption)
		_, err := checkCanUpdateMetadataFieldIsLegal(oldMetadata, req)
		assert.NotNil(t, err)
	}
	{
		// test metadataOption not include "contract" filed
		metadataOption := `{"originId": "originId001", "dataPath": "dataPath001", "rows": 12, "columns": 7, "size": 56, "hasTitle": true, "metadataColumns": [], "consumeTypes": [1, 2, 3], "consumeOptions": ["[]", "[{\"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 2}, {\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 4}, {\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF06\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 11}, {\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 7}]", "[\"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF11\", \"0x79e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\"]"]}`
		oldMetadata, req := mockForTestCheckConsumeOptionsParams(metadataOption)
		_, err := checkCanUpdateMetadataFieldIsLegal(oldMetadata, req)
		assert.NotNil(t, err)
	}
	{
		// test metadataOption not json string
		metadataOption := `{{"originId": "originId001", "dataPath": "dataPath001", "rows": 12, "columns": 7, "size": 56, "hasTitle": true, "metadataColumns": [], "consumeTypes": [1, 2, 3], "consumeOptions": ["[]", "[{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 2}, {\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 4}, {\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF06\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 11}, {\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 7}]", "[\"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF11\", \"0x79e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\"]"]}`
		oldMetadata, req := mockForTestCheckConsumeOptionsParams(metadataOption)
		_, err := checkCanUpdateMetadataFieldIsLegal(oldMetadata, req)
		assert.NotNil(t, err)
	}
}
