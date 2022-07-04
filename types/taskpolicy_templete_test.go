package types

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTaskMetadataPolicy(t *testing.T) {
	type TaskMetadataPolicyBINARY struct {
		PartyId      string `json:"partyId"`
		MetadataId   string `json:"metadataId"`
		MetadataName string `json:"metadataName"`
		InputType    uint32 `json:"inputType"`
	}

	type TaskMetadataPolicyBINARYByPay struct {
		*TaskMetadataPolicyBINARY
		PayType       uint32   `json:"payType"`
		PayOption     string   `json:"payOption"`
	}

	aty := &TaskMetadataPolicyBINARY{
		PartyId: "p0",
		MetadataId: "metadata:1",
		MetadataName: "metadata1",
		InputType: 1,
	}

	atyb, err := json.Marshal(aty)

	assert.Nil(t, err, "json failed")
	t.Logf("aty: %s \n", string(atyb))


	bty := &TaskMetadataPolicyBINARYByPay{
		TaskMetadataPolicyBINARY: &TaskMetadataPolicyBINARY{
			PartyId: "p1",
			MetadataId: "metadata:2",
			MetadataName: "metadata2",
			InputType: 1,
		},
		PayType: 1,
		PayOption: "NFT",
	}

	btyb, err := json.Marshal(bty)

	assert.Nil(t, err, "json failed")
	t.Logf("bty: %s \n", string(btyb))
}
