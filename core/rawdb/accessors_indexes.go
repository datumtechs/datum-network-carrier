// Copyright (C) 2021 The RosettaNet Authors.

package rawdb

import (
	"github.com/datumtechs/datum-network-carrier/common"
	libtypes "github.com/datumtechs/datum-network-carrier/lib/types"
)

// ReadDataLookupEntry retrieves the positional metadata associated with a metadata/resource/identity/task
// dataId hash to allow retrieving the metadata/resource/identity/task by the hash of dataId.
func ReadDataLookupEntry(db DatabaseReader, hash common.Hash) (common.Hash, uint64, uint64, string, string) {
	data, _ := db.Get(dataLookupKey(hash))
	if len(data) == 0 {
		return common.Hash{}, 0, 0, "", ""
	}
	var entry libtypes.DataLookupEntry
	if err := entry.Unmarshal(data); err != nil {
		log.WithField("hash", hash).WithError(err).Fatal("Invalid lookup entry ProtoBuf")
		return common.Hash{}, 0, 0, "", ""
	}
	return common.BytesToHash(entry.BlockHash), entry.BlockIndex, entry.Index, entry.NodeId, entry.Type
}

// WriteDataLookupEntries stores a positional metadata for every metadata/resource/identity from
// a block, enabling dataId hash based metadata、resource、identity and task lookups.
//func WriteDataLookupEntries(db DatabaseWriter, block *libtypes.BlockData) {
//	header := (*types.Header)(block.Header)
//	// todo: what is the type means???
//	for i, metadata := range block.Metadata {
//		writeDataLookupEntries(db, header.Hash(),
//			block.Header.Version, uint64(i), metadata.NodeId, "", metadata.DataId)
//		// store extra data for metadata
//	}
//	for i, resource := range block.Resourcedata {
//		writeDataLookupEntries(db, header.Hash(),
//			block.Header.Version, uint64(i), resource.NodeId, "", resource.DataId)
//		// store extra data for resource
//	}
//	for i, identity := range block.Identitydata {
//		writeDataLookupEntries(db, header.Hash(),
//			block.Header.Version, uint64(i), identity.NodeId, "", identity.DataId)
//		// store extra data for identity
//	}
//	for i, task := range block.Taskdata {
//		writeDataLookupEntries(db, header.Hash(),
//			block.Header.Version, uint64(i), task.NodeId, "", task.DataId)
//		// store extra data for task
//	}
//}
//
//// writeMetadata save detailed data separately for metadata.
//func writeMetadata(db DatabaseWriter, number uint64, index uint64, metadata *libtypes.MetadataPB) {
//	// todo: need to update, the param is wrong.
//	WriteMetadataHash(db, number, index, common.Hash{})
//	WriteMetadataId(db, metadata.NodeId, common.Hash{}, metadata.DataId)
//	WriteMetadataTypeHash(db, metadata.DataId, "", common.Hash{})
//}

func writeDataLookupEntries(db DatabaseWriter, hash common.Hash, number uint64,
							index uint64, nodeId string, typ string, dataId string)  {
	entry := libtypes.DataLookupEntry{
		BlockHash: hash.Bytes(),
		BlockIndex: number,
		Index: index,
		NodeId: nodeId,
		Type: typ,
	}
	data, err := entry.Marshal()
	if err != nil {
		log.WithField("type", typ).WithError(err).Fatal("Failed to encode lookup entry")
	}
	if err := db.Put(dataLookupKey(common.HexToHash(dataId)), data); err != nil {
		log.WithField("type", typ).WithError(err).Fatal("Failed to store lookup entry")
	}
}

// DeleteDataLookupEntry removes all metadata data associated with a dataId.
func DeleteDataLookupEntry(db DatabaseDeleter, dataId string) {
	db.Delete(dataLookupKey(common.HexToHash(dataId)))
}


// ReadMetadata retrieves a specific metadata from the database, along with
// its added positional metadata.
func ReadMetadata(db DatabaseReader, dataId string) (*libtypes.MetadataPB, common.Hash, uint64, uint64, string, string) {
	blockHash, blockNumber, index, nodeId, typ := ReadDataLookupEntry(db, common.HexToHash(dataId))
	if blockHash == (common.Hash{}) {
		return nil, common.Hash{}, 0, 0, "", ""
	}
	body := ReadBody(db, blockHash, blockNumber)
	if body == nil || len(body.Metadata) <= int(index) {
		log.WithField("number", blockNumber).
			WithField("hash", blockHash).
			WithField("index", index).
			Fatal("Medata referenced missing")
		return nil, common.Hash{}, 0, 0, "", ""
	}
	return body.Metadata[index], blockHash, blockNumber, index, nodeId, typ
}

// ReadResource retrieves a specific resource from the database, along with
// its added positional metadata.
func ReadResource(db DatabaseReader, dataId string) (*libtypes.ResourcePB, common.Hash, uint64, uint64, string, string) {
	blockHash, blockNumber, index, nodeId, typ := ReadDataLookupEntry(db, common.HexToHash(dataId))
	if blockHash == (common.Hash{}) {
		return nil, common.Hash{}, 0, 0, "", ""
	}
	body := ReadBody(db, blockHash, blockNumber)
	if body == nil || len(body.Resourcedata) <= int(index) {
		log.WithField("number", blockNumber).
			WithField("hash", blockHash).
			WithField("index", index).
			Fatal("Resource referenced missing")
		return nil, common.Hash{}, 0, 0, "", ""
	}
	return body.Resourcedata[index], blockHash, blockNumber, index, nodeId, typ
}

// ReadIdentity retrieves a specific identity from the database, along with
// its added positional metadata.
func ReadIdentity(db DatabaseReader, dataId string) (*libtypes.IdentityPB, common.Hash, uint64, uint64, string, string) {
	blockHash, blockNumber, index, nodeId, typ := ReadDataLookupEntry(db, common.HexToHash(dataId))
	if blockHash == (common.Hash{}) {
		return nil, common.Hash{}, 0, 0, "", ""
	}
	body := ReadBody(db, blockHash, blockNumber)
	if body == nil || len(body.Identitydata) <= int(index) {
		log.WithField("number", blockNumber).
			WithField("hash", blockHash).
			WithField("index", index).
			Fatal("Identity referenced missing")
		return nil, common.Hash{}, 0, 0, "", ""
	}
	return body.Identitydata[index], blockHash, blockNumber, index, nodeId, typ
}

// ReadIdentity retrieves a specific taskData from the database, along with
// its added positional metadata.
func ReadTask(db DatabaseReader, dataId string) (*libtypes.TaskPB, common.Hash, uint64, uint64, string, string) {
	blockHash, blockNumber, index, nodeId, typ := ReadDataLookupEntry(db, common.HexToHash(dataId))
	if blockHash == (common.Hash{}) {
		return nil, common.Hash{}, 0, 0, "", ""
	}
	body := ReadBody(db, blockHash, blockNumber)
	if body == nil || len(body.Taskdata) <= int(index) {
		log.WithField("number", blockNumber).
			WithField("hash", blockHash).
			WithField("index", index).
			Fatal("GetTask referenced missing")
		return nil, common.Hash{}, 0, 0, "", ""
	}
	return body.Taskdata[index], blockHash, blockNumber, index, nodeId, typ
}
