package types

import (
	"bytes"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"testing"
)

var testTaskdata = NewTask(&types.TaskPB{
	DataId:     "",
	DataStatus: carriertypespb.DataStatus_DataStatus_Unknown,
	TaskId:     "",
	State:      carriertypespb.TaskState_TaskState_Unknown,
	Reason:     "",
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
