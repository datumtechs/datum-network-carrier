package rawdb

import (
	"github.com/RosettaFlow/Carrier-Go/db"
	pbcommon "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func TestRunningTask(t *testing.T) {
	database := db.NewMemoryDatabase()
	task := types.NewTask(&libtypes.TaskPB{
		IdentityId:           "identity-task",
		NodeId:               "nodeId-task",
		NodeName:             "nodeName",
		DataId:               "",
		DataStatus:           pbcommon.DataStatus_DataStatus_Normal,
		TaskId:               "taskID-01",
		TaskName:             "taskName-01",
		State:                pbcommon.TaskState_TaskState_Succeed,
		Reason:               "",
		EventCount:           0,
		Desc:                 "",
		CreateAt:             0,
		EndAt:                0,
	})
	WriteRunningTask(database, task)

	rtask := ReadRunningTask(database, "taskID-01")
	t.Logf("running task info : %v", rtask)
	assert.Assert(t, strings.EqualFold("taskID-01", rtask.GetTaskId()))

	// read all
	taskList := ReadAllRunningTask(database)
	assert.Assert(t, len(taskList) == 1)

	// delete
	DeleteRunningTask(database, "taskID-01")

	taskList = ReadAllRunningTask(database)
	assert.Assert(t, len(taskList) == 0)

}