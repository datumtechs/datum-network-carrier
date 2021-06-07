package types

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/lib/types"
	"testing"
)

var metadata = NewMetadata(&types.MetaData{
	Identity:             "Identity",
	NodeId:               "nodeId",
	DataId:               "dataId",
	DataStatus:           "D",
	FilePath:             "/a/a",
	Desc:                 "desc",
	Rows:                 1,
	Columns:              2,
	Size_:                3,
	FileType:             "csv",
	State:                "create",
	HasTitleRow:          false,
	ColumnMetaList:       []*types.ColumnMeta{
		{
			Cindex:               2,
			Cname:                "cname",
			Ctype:                "ctype",
			Csize:                10,
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