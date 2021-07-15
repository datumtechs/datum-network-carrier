package rawdb

import (
	"github.com/RosettaFlow/Carrier-Go/db"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"testing"
)

func TestResourceDataUsed(t *testing.T) {
	database := db.NewMemoryDatabase()
	dataUsed := types.NewDataResourceFileUpload("node_id", "origin_id_01", "metadata_id", "/a/b/c/d")

	// save
	err := StoreDataResourceDataUsed(database, dataUsed)
	require.Nil(t, err)

	// query
	queryUsed, err := QueryDataResourceDataUsed(database, dataUsed.GetOriginId())
	require.Nil(t, err)
	assert.Equal(t, dataUsed.GetOriginId(), queryUsed.GetOriginId())
}