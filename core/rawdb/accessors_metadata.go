// Copyright (C) 2021 The RosettaNet Authors.

package rawdb

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/db"
)

// ReadAllMetadataHashes retrieves all the hashes of metadata.
func ReadAllMetadataHashes(db db.Iteratee) []common.Hash {
	prefix := metadataIdPrefix

	hashes := make([]common.Hash, 0, 1)
	it := db.NewIterator(prefix, nil)
	defer it.Release()

	for it.Next() {
		// todo: need to determine the length of nodeId.
		if key := it.Key(); len(key) > len(prefix)+32 {
			hashes = append(hashes, common.BytesToHash(key[len(key)-32:]))
		}
	}
	return hashes
}

// ReadAllMetadataHashes retrieves all the hashes of metadata by special nodeId.
func ReadAllMetadataHashesByNodeId(db db.Iteratee, nodeId string) []common.Hash {
	prefix := metadataIdKeyPrefix(nodeId)
	hashes := make([]common.Hash, 0, 1)
	it := db.NewIterator(prefix, nil)
	defer it.Release()

	for it.Next() {
		if key := it.Key(); len(key) == len(prefix)+32 {
			hashes = append(hashes, common.BytesToHash(key[len(key)-32:]))
		}
	}
	return hashes
}

// ReadMetadataHash retrieves the hash assigned to a canonical number and index.
func ReadMetadataHash(db DatabaseReader, number uint64, index uint64) common.Hash {
	data, _ := db.Get(metadataHashKey(number, index))
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteMetadataHash stores the hash assigned to a canonical number and index.
func WriteMetadataHash(db DatabaseWriter, number uint64, index uint64, hash common.Hash)  {
	if err := db.Put(metadataHashKey(number, index), hash.Bytes()); err != nil {
		log.WithError(err).Error("Failed to store number-index to hash mapping for metadata")
	}
}

// DeleteMetadataHash removes the number-index to hash canonical mapping.
func DeleteMetadataHash(db DatabaseDeleter, number uint64, index uint64) {
	if err := db.Delete(metadataHashKey(number, index)); err != nil {
		log.WithError(err).Error("Failed to delete number-index to hash mapping for metadata")
	}
}

// ReadMetadataId retrieves the dataId assigned to a canonical nodeId and hash.
func ReadMetadataId(db DatabaseReader, nodeId string, hash common.Hash) []byte {
	data, _ := db.Get(metadataIdKey(common.Hex2Bytes(nodeId), hash))
	if len(data) == 0 {
		return []byte{}
	}
	return data
}

// WriteMetadataId stores the dataId assigned to a canonical nodeId and hash.
func WriteMetadataId(db DatabaseWriter, nodeId string, hash common.Hash, dataId string)  {
	if err := db.Put(metadataIdKey(common.Hex2Bytes(nodeId), hash), common.Hex2Bytes(dataId)); err != nil {
		log.WithError(err).Error("Failed to store nodeId-hash to dataId mapping for metadata")
	}
}

// DeleteMetadataId removes the nodeId-hash to dataId canonical mapping.
func DeleteMetadataId(db DatabaseDeleter, nodeId string, hash common.Hash,) {
	if err := db.Delete(metadataIdKey(common.Hex2Bytes(nodeId), hash)); err != nil {
		log.WithError(err).Error("Failed to delete number-index to hash mapping for metadata")
	}
}

// ReadMetadataTypeHash retrieves the hash assigned to a canonical type and dataId.
func ReadMetadataTypeHash(db DatabaseReader, dataId string, typ string) common.Hash {
	data, _ := db.Get(metadataTypeHashKey(common.Hex2Bytes(dataId), common.Hex2Bytes(typ)))
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteMetadataTypeHash stores the hash assigned to a canonical type and dataId.
func WriteMetadataTypeHash(db DatabaseWriter, dataId string, typ string, hash common.Hash)  {
	if err := db.Put(metadataTypeHashKey(common.Hex2Bytes(dataId), common.Hex2Bytes(typ)), hash.Bytes()); err != nil {
		log.WithError(err).Error("Failed to store type-dataId to hash mapping for metadata")
	}
}

// DeleteMetadataTypeHash removes the type-dataId to hash canonical mapping.
func DeleteMetadataTypeHash(db DatabaseDeleter, dataId string, typ string) {
	if err := db.Delete(metadataTypeHashKey(common.Hex2Bytes(dataId), common.Hex2Bytes(typ))); err != nil {
		log.WithError(err).Error("Failed to delete type-dataId to hash mapping for metadata")
	}
}

// todo: apply metadata data...