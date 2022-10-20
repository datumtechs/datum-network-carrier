package workflow

import (
	"encoding/json"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	"github.com/datumtechs/datum-network-carrier/types"
	"gotest.tools/assert"
	"testing"
)

func TestDirectedAcyclicGraphTopologicalSort(t *testing.T) {
	{
		// not exits ring
		graph := make(map[string][]string, 0)
		graph["aTask"] = []string{"bTask", "cTask", "eTask"}
		graph["bTask"] = []string{"dTask"}
		graph["cTask"] = []string{"dTask"}
		graph["dTask"] = []string{}
		graph["eTask"] = []string{"cTask", "dTask"}
		result := directedAcyclicGraphTopologicalSort(graph)
		assert.DeepEqual(t, []string{"aTask", "eTask", "cTask", "bTask", "dTask"}, result)
	}
	{
		// exits ring
		graph := make(map[string][]string, 0)
		graph["aTask"] = []string{"bTask", "cTask", "eTask"}
		graph["bTask"] = []string{"dTask"}
		graph["cTask"] = []string{"dTask"}
		graph["dTask"] = []string{}
		graph["eTask"] = []string{"cTask", "dTask", "aTask"}
		result := directedAcyclicGraphTopologicalSort(graph)
		assert.Equal(t, 0, len(result))
	}
}

func TestCheckWorkflowTaskListReferTo(t *testing.T) {
	{
		testReq := &carrierapipb.PublishWorkFlowDeclareRequest{
			PolicyType: 1,
			Policy:     `[{"origin": "aTask", "reference": [{"target": "bTask"}, {"target": "cTask"}, {"target": "eTask"}]}, {"origin": "bTask", "reference": [{"target": "dTask"}]}, {"origin": "cTask", "reference": [{"target": "dTask"}]}, {"origin": "dTask", "reference": []}, {"origin": "eTask", "reference": [{"target": "cTask"}, {"target": "dTask"}]}]`,
			TaskList: []*carrierapipb.PublishTaskDeclareRequest{
				{
					TaskName: "aTask",
				},
				{
					TaskName: "eTask",
				},
				{
					TaskName: "dTask",
				},
				{
					TaskName: "cTask",
				},
				{
					TaskName: "bTask",
				},
			},
		}
		result := checkWorkflowTaskListReferTo(testReq)
		assert.Equal(t, false, result)
		t.Logf("taskList %v", testReq.TaskList)
	}
	{
		testReq := &carrierapipb.PublishWorkFlowDeclareRequest{
			PolicyType: 1,
			Policy:     `[{"origin": "aTask", "reference": [{"target": "bTask"}, {"target": "cTask"}, {"target": "eTask"}]}, {"origin": "bTask", "reference": [{"target": "dTask"}]}, {"origin": "cTask", "reference": [{"target": "dTask"}]}, {"origin": "dTask", "reference": []}, {"origin": "eTask", "reference": [{"target": "cTask"}, {"target": "dTask"},{"target": "aTask"}]}]`,
			TaskList: []*carrierapipb.PublishTaskDeclareRequest{
				{
					TaskName: "aTask",
				},
				{
					TaskName: "eTask",
				},
				{
					TaskName: "dTask",
				},
				{
					TaskName: "cTask",
				},
				{
					TaskName: "bTask",
				},
			},
		}
		result := checkWorkflowTaskListReferTo(testReq)
		assert.Equal(t, true, result)
		t.Logf("taskList %v", testReq.TaskList)
	}
	test := "{\"inputType\": 1, \"keyColumnName\": \"12\", \"selectedColumnNames\": []}"
	var p *types.TaskMetadataPolicyCSVWithTaskResultData
	_ = json.Unmarshal([]byte(test), &p)
	t.Logf("%v,%s", p, p.PartyId)
}
