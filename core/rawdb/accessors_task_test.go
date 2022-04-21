package rawdb

import (
	"github.com/Metisnetwork/Metis-Carrier/db"
	libtypes "github.com/Metisnetwork/Metis-Carrier/lib/types"
	"github.com/Metisnetwork/Metis-Carrier/types"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func TestRunningTask(t *testing.T) {
	database := db.NewMemoryDatabase()
	task := types.NewTask(&libtypes.TaskPB{
		Sender: &libtypes.TaskOrganization{
			PartyId:    "p0",
			IdentityId: "identity-task",
			NodeId:     "nodeId-task",
			NodeName:   "nodeName",
		},
		DataId:     "",
		DataStatus: libtypes.DataStatus_DataStatus_Valid,
		TaskId:     "taskID-01",
		TaskName:   "taskName-01",
		State:      libtypes.TaskState_TaskState_Succeed,
		Reason:     "",
		Desc:       "",
		CreateAt:   0,
		EndAt:      0,
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
