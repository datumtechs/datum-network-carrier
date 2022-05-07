package types

import (
	"bytes"
	"github.com/Metisnetwork/Metis-Carrier/common"
	"github.com/Metisnetwork/Metis-Carrier/lib/types"
	libtypes "github.com/Metisnetwork/Metis-Carrier/lib/types"
	"gotest.tools/assert"
	"testing"
)

var testTaskdata = NewTask(&types.TaskPB{
	DataId:     "",
	DataStatus: libtypes.DataStatus_DataStatus_Unknown,
	TaskId:     "",
	State:      libtypes.TaskState_TaskState_Unknown,
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

func TestTaskHashEqual(t *testing.T) {

	task1 := NewTask(&types.TaskPB{
		TaskId: "0xwqrqrqwr",
		TaskName: "Gavin",
		UserType: libtypes.UserType_User_1,
		User: "0x34234242",
	})

	task2 := NewTask(&types.TaskPB{
		TaskId: "0xwqrqrqwr",
		TaskName: "Gavin",
		UserType: libtypes.UserType_User_1,
		User: "0x34234242",
	})
	assert.Equal(t, task1.Hash(), task2.Hash(), "hash is not same")
}