package types

import (
	"bytes"
	"github.com/datumtechs/datum-network-carrier/common"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"testing"
)

var resource_test = NewResource(&carriertypespb.ResourcePB{
	DataId:               "dataId",
	DataStatus:           commonconstantpb.DataStatus_DataStatus_Invalid,
	State:                commonconstantpb.PowerState_PowerState_Created,
	TotalMem:             1,
	UsedMem:              2,
	TotalProcessor:       0,
	TotalBandwidth:       1,
})

func TestResourceEncode(t *testing.T) {
	buffer := new(bytes.Buffer)
	err := resource_test.EncodePb(buffer)
	if err != nil {
		t.Fatal("resource encode protobuf failed, err: ", err)
	}

	dresource := new(Resource)
	err = dresource.DecodePb(buffer.Bytes())
	if err != nil {
		t.Fatal("resource decode protobuf failed, err: ", err)
	}
	dBuffer := new(bytes.Buffer)
	dresource.EncodePb(dBuffer)

	if !bytes.Equal(buffer.Bytes(), dBuffer.Bytes()) {
		t.Fatalf("encode protobuf mismatch, got %x, want %x", common.Bytes2Hex(dBuffer.Bytes()), common.Bytes2Hex(buffer.Bytes()))
	}
}
