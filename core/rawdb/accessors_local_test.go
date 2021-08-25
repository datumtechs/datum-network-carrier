package rawdb

import (
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/db"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func TestLocalTask(t *testing.T) {
	database := db.NewMemoryDatabase()
	data01 := &libTypes.TaskData{
		IdentityId:             "identity",
		NodeId:               "nodeid",
		NodeName:             "nodename",
		DataId:               "taskId",
		DataStatus:           "Y",
		TaskId:               "taskID",
		TaskName:             "taskName",
		State:                "SUCCESS",
		Reason:               "reason",
		EventCount:           4,
		Desc:                 "desc",
		CreateAt:             uint64(timeutils.UnixMsec()),
		EndAt:                uint64(timeutils.UnixMsec()),
	}
	WriteLocalTask(database, types.NewTask(data01))

	res, _ := ReadLocalTask(database, data01.TaskId)
	assert.Assert(t, strings.EqualFold(data01.TaskId, res.TaskId()))

	// test update state
	res.TaskData().State = "failed"
	DeleteLocalTask(database, data01.TaskId)
	WriteLocalTask(database, res)

	res, _ = ReadLocalTask(database, data01.TaskId)
	assert.Equal(t, "failed", res.TaskData().State)

	taskList, _ := ReadAllLocalTasks(database)
	assert.Assert(t, len(taskList) == 1)

	// delete
	DeleteLocalTask(database, data01.TaskId)

	taskList, _ = ReadAllLocalTasks(database)
	assert.Assert(t, len(taskList) == 0)

	// delete
	DeleteSeedNodes(database)

	taskList, _ = ReadAllLocalTasks(database)
	assert.Assert(t, len(taskList) == 0)
}

func TestSeedNode(t *testing.T) {
	// write seed
	database := db.NewMemoryDatabase()
	seedNodeInfo := &pb.SeedPeer{
		Id:           "id",
		InternalIp:   "internalIp",
		InternalPort: "9999",
		ConnState:    1,
	}
	WriteSeedNodes(database, seedNodeInfo)

	// get seed
	rseed, _ := ReadSeedNode(database, "id")
	t.Logf("seed info : %v", rseed)
	assert.Assert(t, strings.EqualFold("id", rseed.Id))

	// read all
	seedNodes, _ := ReadAllSeedNodes(database)
	assert.Assert(t, len(seedNodes) == 1)

	// delete
	DeleteSeedNode(database, "id")

	seedNodes, _ = ReadAllSeedNodes(database)
	assert.Assert(t, len(seedNodes) == 0)

	// delete
	DeleteSeedNodes(database)

	seedNodes, _ = ReadAllSeedNodes(database)
	assert.Assert(t, len(seedNodes) == 0)
}

func TestRegisteredNode(t *testing.T) {
	// write seed
	database := db.NewMemoryDatabase()
	registered := &pb.YarnRegisteredPeerDetail{
		Id:           "id",
		InternalIp:   "internalIp",
		InternalPort: "9999",
		ExternalIp:   "externalIp",
		ExternalPort: "999",
		ConnState:    1,
	}
	WriteRegisterNodes(database, pb.PrefixTypeJobNode, registered)

	// get seed
	r, _ := ReadRegisterNode(database, pb.PrefixTypeJobNode, "id")
	t.Logf("registered info : %v", r)
	assert.Assert(t, strings.EqualFold("id", r.Id))

	// read all
	registeredNodes, _ := ReadAllRegisterNodes(database, pb.PrefixTypeJobNode)
	assert.Assert(t, len(registeredNodes) == 1)

	// delete
	DeleteRegisterNode(database, pb.PrefixTypeJobNode, "id")

	registeredNodes, _ = ReadAllRegisterNodes(database, pb.PrefixTypeJobNode)
	assert.Assert(t, len(registeredNodes) == 0)

	// delete
	DeleteRegisterNodes(database, pb.PrefixTypeJobNode)

	registeredNodes, _ = ReadAllRegisterNodes(database, pb.PrefixTypeJobNode)
	assert.Assert(t, len(registeredNodes) == 0)
}

func TestTaskEvent(t *testing.T) {
	database := db.NewMemoryDatabase()
	taskEvent := &libTypes.TaskEvent{
		Type:       "taskEventType",
		IdentityId:   "taskEventIdentity",
		TaskId:     "taskEventTaskId",
		Content:    "taskEventContent",
		CreateAt: uint64(timeutils.UnixMsec()),
	}
	WriteTaskEvent(database, taskEvent)

	taskEvent2 := &libTypes.TaskEvent{
		Type:       "taskEventType-02",
		IdentityId:   "taskEventIdentity",
		TaskId:     "taskEventTaskId",
		Content:    "taskEventContent-02",
		CreateAt: uint64(timeutils.UnixMsec()),
	}
	WriteTaskEvent(database, taskEvent2)

	revent, _ := ReadTaskEvent(database, "taskEventTaskId")
	t.Logf("task evengine info : %v", len(revent))
	assert.Assert(t, strings.EqualFold("taskEventIdentity", revent[0].IdentityId))

	// read all
	taskEvents, _ := ReadAllTaskEvents(database)
	assert.Assert(t, len(taskEvents) == 2)

	// delete
	DeleteTaskEvent(database, "taskEventTaskId")

	taskEvents, _ = ReadAllTaskEvents(database)
	assert.Assert(t, len(taskEvents) == 0)
}

func TestLocalIdentity(t *testing.T) {
	database := db.NewMemoryDatabase()
	nodeAlias := &types.NodeAlias{
		Name:       "node-name",
		NodeId:     "node-nodeId",
		IdentityId: "node-identityId",
	}
	WriteLocalIdentity(database, nodeAlias)
	WriteLocalIdentity(database, nodeAlias)

	rnode, _ := ReadLocalIdentity(database)
	assert.Equal(t, rnode.IdentityId, nodeAlias.IdentityId)
	assert.Equal(t, rnode.NodeId, nodeAlias.NodeId)
	assert.Equal(t, rnode.Name, nodeAlias.Name)

	DeleteLocalIdentity(database)
	rnode, _ = ReadLocalIdentity(database)
	assert.Equal(t, rnode.IdentityId, "")
	assert.Equal(t, rnode.NodeId, "")
	assert.Equal(t, rnode.Name, "")
}

func TestLocalResource(t *testing.T) {
	database := db.NewMemoryDatabase()
	localResource01 := &libTypes.LocalResourceData{
		IdentityId:             "01-identity",
		NodeId:               "01-nodeId",
		NodeName:             "01-nodename",
		JobNodeId:            "01",
		DataId:               "01-dataId",
		DataStatus:           "y",
		State:                "SUCCESS",
		TotalMem:             111,
		UsedMem:              222,
		TotalProcessor:       11,
		UsedProcessor:        33,
		TotalBandWidth:       33,
		UsedBandWidth:        44,
	}
	b, _ := localResource01.Marshal()
	_ = b
	WriteLocalResource(database, types.NewLocalResource(localResource01))

	localResource02 := &libTypes.LocalResourceData{
		IdentityId:             "01-identity",
		NodeId:               "01-nodeId",
		NodeName:             "01-nodename",
		JobNodeId:            "02",
		DataId:               "01-dataId",
		DataStatus:           "y",
		State:                "SUCCESS",
		TotalMem:             111,
		UsedMem:              222,
		TotalProcessor:       11,
		UsedProcessor:        33,
		TotalBandWidth:       33,
		UsedBandWidth:        44,
	}
	WriteLocalResource(database, types.NewLocalResource(localResource02))

	r1, _ := ReadLocalResource(database, localResource01.JobNodeId)
	assert.Equal(t, types.NewLocalResource(localResource01).Hash(), r1.Hash())

	array, _ := ReadAllLocalResource(database)
	require.True(t, array.Len() == 2)

	DeleteLocalResource(database, localResource01.JobNodeId)
	array, _ = ReadAllLocalResource(database)
	require.True(t, array.Len() == 1)
}