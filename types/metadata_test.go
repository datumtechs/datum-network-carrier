package types

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/lib/types"
	"reflect"
	"testing"
)

var metadata = NewMetadata(&types.MetadataPB{
	IdentityId: "Identity",
	NodeId:     "nodeId",
	DataId:     "dataId",
	DataStatus: "D",
	FilePath:   "/a/a",
	Desc:       "desc",
	Rows:       1,
	Columns:    2,
	Size_:      3,
	FileType:   "csv",
	State:      "create",
	HasTitle:   false,
	MetadataColumnList: []*types.MetadataColumn {
		{
			CIndex: 2,
			CName:  "cname",
			CType:  "ctype",
			CSize:  10,
		},
	},
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

	if !bytes.Equal(buffer.Bytes(), dBuffer.Bytes()) {
		t.Fatalf("encode protobuf mismatch, got %x, want %x", common.Bytes2Hex(dBuffer.Bytes()), common.Bytes2Hex(buffer.Bytes()))
	}
}

func TestMetadata(t *testing.T) {
	byts := common.Hex2Bytes("0a084964656e7469747912066e6f646549641a066461746149642201442a042f612f6132046465736338014002480352036373765a066372656174656a1208021205636e616d651a056374797065200a")
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
