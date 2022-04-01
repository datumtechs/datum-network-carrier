package types

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/lib/types"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

var metadata = NewMetadata(&types.MetadataPB{
	MetadataId: "xxxxx",
	Owner: &libtypes.Organization{
		IdentityId: "identity",
		NodeId:     "nodeId",
		NodeName:   "nodeName",
	},
	DataId:         "dataId",
	DataStatus:     libtypes.DataStatus_DataStatus_Valid,
	MetadataName:   "this is a metadata",
	MetadataType:   1,
	Desc:           "desc",
	Industry:       "1",
	State:          libtypes.MetadataState_MetadataState_Created,
	PublishAt:      timeutils.UnixMsecUint64(),
	UpdateAt:       timeutils.UnixMsecUint64(),
	Nonce:          1,
	MetadataOption: "option",
	TokenAddress:   "",
})

func TestMetadataEncode(t *testing.T) {
	buffer := new(bytes.Buffer)
	err := metadata.EncodePb(buffer)
	if err != nil {
		t.Fatal("encode protobuf failed, err: ", err)
	}

	dmetadata := new(Metadata)
	err = dmetadata.DecodePb(buffer.Bytes())
	if err != nil {
		t.Fatal("decode protobuf failed, err: ", err)
	}
	dBuffer := new(bytes.Buffer)
	dmetadata.EncodePb(dBuffer)
	t.Log(common.Bytes2Hex(buffer.Bytes()))
	if !bytes.Equal(buffer.Bytes(), dBuffer.Bytes()) {
		t.Fatalf("encode protobuf mismatch, got %x, want %x", common.Bytes2Hex(dBuffer.Bytes()), common.Bytes2Hex(buffer.Bytes()))
	}
}

func TestMetadata(t *testing.T) {
	byts := common.Hex2Bytes("12084964656e746974791a066e6f646549642a0664617461496430014a042f612f61520464657363580160026803700178018a011208021205636e616d651a056374797065200a")
	dmetadata := new(Metadata)
	err := dmetadata.DecodePb(byts)
	if err != nil {
		t.Fatal("decode protobuf failed, err: ", err)
	}
	check := func(f string, got, want interface{}) {
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s mismatch: got %v, want %v", f, got, want)
		}
	}
	check("hash", dmetadata.Hash(), metadata.Hash())
}

func TestMetadataArrLen(t *testing.T) {
	var metadataArray MetadataArray = nil

	t.Log(metadataArray.Len())
	// Len returns the length of s.
	assert.Equal(t, metadataArray.Len(), 0)
}
