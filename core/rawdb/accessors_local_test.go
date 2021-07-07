package rawdb

import (
	"github.com/RosettaFlow/Carrier-Go/db"
	"github.com/RosettaFlow/Carrier-Go/types"
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