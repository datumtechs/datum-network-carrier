// Copyright (C) 2021 The RosettaNet Authors.

package rawdb

import (
	dbtype "github.com/RosettaFlow/Carrier-Go/lib/db"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/sirupsen/logrus"
	"strings"
)

const seedNodeToKeep = 50

// ReadSeedNode retrieves the seed node with the corresponding nodeId.
func ReadSeedNode(db DatabaseReader, nodeId string) *types.SeedNodeInfo {
	blob, err := db.Get(seedNodeKey)
	if err != nil {
		return nil
	}
	var seedNodes dbtype.SeedNodeListPB
	if err := seedNodes.Unmarshal(blob); err != nil {
		return nil
	}
	for _, seed := range seedNodes.GetSeedNodeList() {
		if strings.EqualFold(seed.Id, nodeId) {
			return &types.SeedNodeInfo{
				Id:           seed.Id,
				InternalIp:   seed.InternalIp,
				InternalPort: seed.InternalPort,
				ConnState:    types.NodeConnStatus(seed.ConnState),
			}
		}
	}
	return nil
}

// ReadAllSeedNodes retrieves all the seed nodes in the database.
// All returned seed nodes are sorted in reverse.
func ReadAllSeedNodes(db DatabaseReader) []*types.SeedNodeInfo {
	blob, err := db.Get(seedNodeKey)
	if err != nil {
		return nil
	}
	var seedNodes dbtype.SeedNodeListPB
	if err := seedNodes.Unmarshal(blob); err != nil {
		return nil
	}
	var nodes []*types.SeedNodeInfo
	for _, seed := range seedNodes.SeedNodeList {
		nodes = append(nodes, &types.SeedNodeInfo{
			Id:           seed.Id,
			InternalIp:   seed.InternalIp,
			InternalPort: seed.InternalPort,
			ConnState:    types.NodeConnStatus(seed.ConnState),
		})
	}
	return nodes
}

// WriteSeedNodes serializes the seed node into the database. If the cumulated
// seed node exceeds the limitation, the oldest will be dropped.
func WriteSeedNodes(db ethdb.KeyValueStore, seedNode *types.SeedNodeInfo) {
	blob, err := db.Get(seedNodeKey)
	if err != nil {
		log.Warn("Failed to load old seed nodes", "error", err)
	}
	var seedNodes dbtype.SeedNodeListPB
	if len(blob) > 0 {
		if err := seedNodes.Unmarshal(blob); err != nil {
			log.WithError(err).Error("Failed to decode old seed nodes")
		}

	}
	for _, s := range seedNodes.GetSeedNodeList() {
		if strings.EqualFold(s.Id, seedNode.Id) {
			log.WithFields(logrus.Fields{ "id": s.Id }).Info("Skip duplicated seed node")
			return
		}
	}
	seedNodes.SeedNodeList = append(seedNodes.SeedNodeList, &dbtype.SeedNodePB{
		Id:                   seedNode.Id,
		InternalIp:           seedNode.InternalIp,
		InternalPort:         seedNode.InternalPort,
		ConnState:            int32(seedNode.ConnState),
	})
	// max limit for store seed node.
	if len(seedNodes.SeedNodeList) > seedNodeToKeep {
		seedNodes.SeedNodeList = seedNodes.SeedNodeList[:seedNodeToKeep]
	}
	data, err := seedNodes.Marshal()
	if err != nil {
		log.WithError(err).Error("Failed to encode seed node")
	}
	if err := db.Put(seedNodeKey, data); err != nil {
		log.WithError(err).Error("Failed to write bad blocks")
	}
}

// DeleteSeedNodes deletes all the seed nodes from the database
func DeleteSeedNodes(db DatabaseDeleter) {
	if err := db.Delete(seedNodeKey); err != nil {
		log.WithError(err).Error("Failed to delete seed node")
	}
}
