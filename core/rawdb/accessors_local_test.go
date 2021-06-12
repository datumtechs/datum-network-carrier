package rawdb

import (
	"github.com/RosettaFlow/Carrier-Go/db"
	"github.com/RosettaFlow/Carrier-Go/types"
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
}
