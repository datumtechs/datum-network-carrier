package rawdb

import (
	"encoding/json"
	"github.com/datumtechs/datum-network-carrier/db"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gogo/protobuf/proto"
	assertPlus "github.com/stretchr/testify/assert"
	"gotest.tools/assert"
	"testing"
	"time"
)

var database = db.NewMemoryDatabase()

func TestSaveSendToTaskManager(t *testing.T) {
	taskId := "task:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e65"
	workflowId := "workflow:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e65"
	if err := SaveSendToTaskManager(database, taskId, workflowId); err != nil {
		t.Fatal("SaveSendToTaskManager fail,", err)
	}

	result, err := database.Get(GetSendToTaskManagerCacheKeyPrefix(taskId))
	assert.NilError(t, err)
	var getWorkflowId string
	if err := rlp.DecodeBytes(result, &getWorkflowId); err != nil {
		t.Fatal("DecodeBytes fail getWorkflowId")
	}
	assert.Equal(t, getWorkflowId, workflowId)

	if err := RemoveSendToTaskManager(database, taskId); err != nil {
		t.Fatal("RemoveSendToTaskManager fail", err)
	}
	result, err = database.Get(GetSendToTaskManagerCacheKeyPrefix(taskId))
	assertPlus.NotNil(t, err)
}

func TestSaveWorkflowCache(t *testing.T) {
	wfId := "workflow:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e65"
	workflowSave := &carriertypespb.Workflow{
		WorkflowId:   wfId,
		Desc:         "test_WorkflowCache",
		WorkflowName: "test_WorkflowName",
		PolicyType:   1,
		Policy:       `[{"origin": "aTask", "reference": [{"target": "bTask"}, {"target": "cTask"}, {"target": "eTask"}]}, {"origin": "bTask", "reference": [{"target": "dTask"}]}, {"origin": "cTask", "reference": [{"target": "dTask"}]}, {"origin": "dTask", "reference": []}, {"origin": "eTask", "reference": [{"target": "cTask"}, {"target": "dTask"},{"target": "aTask"}]}]`,
		User:         "User1",
		UserType:     1,
		Sign:         []byte("signTest"),
		Tasks: []*carriertypespb.TaskMsg{
			{
				Data: &carriertypespb.TaskPB{
					TaskName: "this is a test task1",
					TaskId:   "task:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e65",
					UserType: commonconstantpb.UserType_User_1,
				},
			},
			{
				Data: &carriertypespb.TaskPB{
					TaskName: "this is a test task2",
					TaskId:   "task:0xb7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e65",
					UserType: commonconstantpb.UserType_User_1,
				},
			},
		},
	}
	if err := SaveWorkflowCache(database, workflowSave); err != nil {
		t.Fatal("SaveSendToTaskManager fail,", err)
	}

	result, err := database.Get(GetWorkflowsCacheKeyPrefix(workflowSave.GetWorkflowId()))
	assert.NilError(t, err)
	var workflow carriertypespb.Workflow
	if err := proto.Unmarshal(result, &workflow); err != nil {
		t.Fatal("DecodeBytes workflow fail")
	}
	assert.Equal(t, wfId, workflow.GetWorkflowId())
	for _, tk := range workflow.Tasks {
		t.Logf("taskId=>%s", tk.Data.GetTaskId())
	}
}

func TestSaveWorkflowStatusCache(t *testing.T) {
	wfId := "workflow:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e65"
	workflowStatus := &types.WorkflowStatus{
		Status:   commonconstantpb.WorkFlowState_WorkFlowState_Running,
		UpdateAt: time.Now().Unix(),
	}
	err := SaveWorkflowStatusCache(database, wfId, workflowStatus)
	result, err := database.Get(GetWorkflowStatusCacheKeyPrefix(wfId))
	var workflowState *types.WorkflowStatus
	if err := json.Unmarshal(result, &workflowState); err != nil {
		t.Fatal("json.Unmarshal workflowState fail")
	} else {
		t.Logf("workflowState is %v", workflowStatus)
	}
	assertPlus.Nil(t, err)
}

func TestSaveWorkflowTaskStatusCache(t *testing.T) {
	wfId := "workflow:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e65"
	taskName := "this is a test name"
	if err := SaveWorkflowTaskStatusCache(database, wfId, &carrierapipb.WorkFlowTaskStatus{
		TaskId:   "task:0xb7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e65",
		Status:   2,
		TaskName: taskName,
	}); err != nil {
		t.Fatal("SaveSendToTaskManager fail,", err)
	}

	//workflowIdTaskName := fmt.Sprintf("%s%s", wfId, taskName)
	result, err := database.Get(GetWorkflowTaskStatusCacheKeyPrefix(wfId, taskName))
	assertPlus.Nil(t, err)
	var workflowStatus carrierapipb.WorkFlowTaskStatus
	if err := proto.Unmarshal(result, &workflowStatus); err != nil {
		t.Fatal("DecodeBytes workflow fail")
	}
	assert.Equal(t, taskName, workflowStatus.TaskName)
	testKey := GetWorkflowTaskStatusCacheKeyPrefix(wfId, taskName)

	workflowIdTaskNameK := testKey[len(workflowTaskStatusCacheKeyPrefix):]
	assert.Equal(t, wfId, string(workflowIdTaskNameK[:len(wfId)]))
	assert.Equal(t, taskName, string(workflowIdTaskNameK[len(wfId):]))
}
