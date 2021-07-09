package rawdb

import (
	"github.com/RosettaFlow/Carrier-Go/db"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"strings"
	"testing"
	"time"
)

func TestSeedNode(t *testing.T) {
	// write seed
	database := db.NewMemoryDatabase()
	seedNodeInfo := &types.SeedNodeInfo{
		Id:           "id",
		InternalIp:   "internalIp",
		InternalPort: "9999",
		ConnState:    1,
	}
	WriteSeedNodes(database, seedNodeInfo)

	// get seed
	rseed := ReadSeedNode(database, "id")
	t.Logf("seed info : %v", rseed)
	assert.Assert(t, strings.EqualFold("id", rseed.Id))

	// read all
	seedNodes := ReadAllSeedNodes(database)
	assert.Assert(t, len(seedNodes) == 1)

	// delete
	DeleteSeedNode(database, "id")

	seedNodes = ReadAllSeedNodes(database)
	assert.Assert(t, len(seedNodes) == 0)

	// delete
	DeleteSeedNodes(database)

	seedNodes = ReadAllSeedNodes(database)
	assert.Assert(t, len(seedNodes) == 0)
}

func TestRegisteredNode(t *testing.T) {
	// write seed
	database := db.NewMemoryDatabase()
	registered := &types.RegisteredNodeInfo{
		Id:           "id",
		InternalIp:   "internalIp",
		InternalPort: "9999",
		ExternalIp:   "externalIp",
		ExternalPort: "999",
		ConnState:    1,
	}
	WriteRegisterNodes(database, types.PREFIX_TYPE_JOBNODE, registered)

	// get seed
	r := ReadRegisterNode(database, types.PREFIX_TYPE_JOBNODE, "id")
	t.Logf("registered info : %v", r)
	assert.Assert(t, strings.EqualFold("id", r.Id))

	// read all
	registeredNodes := ReadAllRegisterNodes(database, types.PREFIX_TYPE_JOBNODE)
	assert.Assert(t, len(registeredNodes) == 1)

	// delete
	DeleteRegisterNode(database, types.PREFIX_TYPE_JOBNODE, "id")

	registeredNodes = ReadAllRegisterNodes(database, types.PREFIX_TYPE_JOBNODE)
	assert.Assert(t, len(registeredNodes) == 0)

	// delete
	DeleteRegisterNodes(database, types.PREFIX_TYPE_JOBNODE)

	registeredNodes = ReadAllRegisterNodes(database, types.PREFIX_TYPE_JOBNODE)
	assert.Assert(t, len(registeredNodes) == 0)
}

func TestTaskEvent(t *testing.T) {
	database := db.NewMemoryDatabase()
	taskEvent := &types.TaskEventInfo{
		Type:       "taskEventType",
		Identity:   "taskEventIdentity",
		TaskId:     "taskEventTaskId",
		Content:    "taskEventContent",
		CreateTime: uint64(time.Now().Second()),
	}
	WriteTaskEvent(database, taskEvent)

	taskEvent2 := &types.TaskEventInfo{
		Type:       "taskEventType-02",
		Identity:   "taskEventIdentity",
		TaskId:     "taskEventTaskId",
		Content:    "taskEventContent-02",
		CreateTime: uint64(time.Now().Second()),
	}
	WriteTaskEvent(database, taskEvent2)

	revent := ReadTaskEvent(database, "taskEventTaskId")
	t.Logf("task evengine info : %v", len(revent))
	assert.Assert(t, strings.EqualFold("taskEventIdentity", revent[0].Identity))

	// read all
	taskEvents := ReadAllTaskEvents(database)
	assert.Assert(t, len(taskEvents) == 2)

	// delete
	DeleteTaskEvent(database, "taskEventTaskId")

	taskEvents = ReadAllTaskEvents(database)
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

	rnode := ReadLocalIdentity(database)
	assert.Equal(t, rnode.IdentityId, nodeAlias.IdentityId)
	assert.Equal(t, rnode.NodeId, nodeAlias.NodeId)
	assert.Equal(t, rnode.Name, nodeAlias.Name)

	DeleteLocalIdentity(database)
	rnode = ReadLocalIdentity(database)
	assert.Equal(t, rnode.IdentityId, "")
	assert.Equal(t, rnode.NodeId, "")
	assert.Equal(t, rnode.Name, "")
}

func TestLocalResource(t *testing.T) {
	database := db.NewMemoryDatabase()
	localResource01 := &libtypes.LocalResourceData{
		Identity:             "01-identity",
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

	localResource02 := &libtypes.LocalResourceData{
		Identity:             "01-identity",
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

	r1 := ReadLocalResource(database, localResource01.JobNodeId)
	assert.Equal(t, types.NewLocalResource(localResource01).Hash(), r1.Hash())

	array := ReadAllLocalResource(database)
	require.True(t, array.Len() == 2)

	DeleteLocalResource(database, localResource01.JobNodeId)
	array = ReadAllLocalResource(database)
	require.True(t, array.Len() == 1)
}