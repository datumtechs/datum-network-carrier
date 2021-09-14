package types

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/lib/types"
	"testing"
)

var resource_test = NewResource(&types.ResourceData{
	Identity:             "identity",
	NodeId:               "nodeId",
	NodeName:             "nodeName",
	DataId:               "dataId",
	DataStatus:           "D",
	State:                "1",
	TotalMem:             1,
	UsedMem:              2,
	TotalProcessor:       0,
	TotalBandWidth:       1,
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
