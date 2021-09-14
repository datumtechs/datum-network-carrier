package types

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/lib/types"
	"testing"
)

var testTaskdata = NewTask(&types.TaskData{
	Identity:             "",
	NodeId:               "",
	NodeName:             "",
	DataId:               "",
	DataStatus:           "",
	TaskId:               "",
	State:                "",
	Reason:               "",
	EventCount:           0,
	Desc:                 "",
	PartnerList:          nil,
	EventDataList:        nil,
})

func TestTaskDataEncode(t *testing.T) {
	buffer := new(bytes.Buffer)
	err := testTaskdata.EncodePb(buffer)
	if err != nil {
		t.Fatal("task encode protobuf failed, err: ", err)
	}

	dtaskdata := new(Task)
	err = dtaskdata.DecodePb(buffer.Bytes())
	if err != nil {
		t.Fatal("task decode protobuf failed, err: ", err)
	}
	dBuffer := new(bytes.Buffer)
	dtaskdata.EncodePb(dBuffer)

	if !bytes.Equal(buffer.Bytes(), dBuffer.Bytes()) {
		t.Fatalf("task encode protobuf mismatch, got %x, want %x", common.Bytes2Hex(dBuffer.Bytes()), common.Bytes2Hex(buffer.Bytes()))
	}
}
