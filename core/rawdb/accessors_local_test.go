package rawdb

import (
	"github.com/RosettaFlow/Carrier-Go/db"
	"github.com/RosettaFlow/Carrier-Go/types"
	"gotest.tools/assert"
	"strings"
	"testing"
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
