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
	DeleteSeedNodes(database)

	seedNodes = ReadAllSeedNodes(database)
	assert.Assert(t, len(seedNodes) == 0)
}
