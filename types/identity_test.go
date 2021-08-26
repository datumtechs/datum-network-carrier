package types

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/lib/types"
	"testing"
)

var identities = NewIdentity(&types.IdentityData{
	IdentityId: "",
	NodeId:     "",
	NodeName:   "",
	DataId:     "",
	DataStatus: "",
	Status:     "",
	Credential: "",
})

func TestIdentitiesEncode(t *testing.T) {
	buffer := new(bytes.Buffer)
	err := identities.EncodePb(buffer)
	if err != nil {
		t.Fatal("identity encode protobuf failed, err: ", err)
	}

	didentities := new(Identity)
	err = didentities.DecodePb(buffer.Bytes())
	if err != nil {
		t.Fatal("identity decode protobuf failed, err: ", err)
	}
	dBuffer := new(bytes.Buffer)
	didentities.EncodePb(dBuffer)

	if !bytes.Equal(buffer.Bytes(), dBuffer.Bytes()) {
		t.Fatalf("identity encode protobuf mismatch, got %x, want %x", common.Bytes2Hex(dBuffer.Bytes()), common.Bytes2Hex(buffer.Bytes()))
	}
}
