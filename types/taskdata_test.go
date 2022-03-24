package types

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/common"
	libcommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/lib/types"
	"testing"
)

var testTaskdata = NewTask(&types.TaskPB{
	IdentityId: "",
	NodeId:     "",
	NodeName:   "",
	DataId:     "",
	DataStatus: libcommonpb.DataStatus_DataStatus_Unknown,
	TaskId:     "",
	State:      libcommonpb.TaskState_TaskState_Unknown,
	Reason:     "",
	EventCount: 0,
	Desc:       "",
	TaskEvents: nil,
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
