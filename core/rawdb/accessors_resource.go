// Copyright (C) 2021 The RosettaNet Authors.

package rawdb


import "github.com/Metisnetwork/Metis-Carrier/common"

// ReadResourceHash retrieves the hash assigned to a canonical number and index.
func ReadResourceHash(db DatabaseReader, number uint64, index uint64) common.Hash {
	data, _ := db.Get(resourceHashKey(number, index))
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteResourceHash stores the hash assigned to a canonical number and index.
func WriteResourceHash(db DatabaseWriter, number uint64, index uint64, hash common.Hash)  {
	if err := db.Put(resourceHashKey(number, index), hash.Bytes()); err != nil {
		log.WithError(err).Fatal("Failed to store number-index to hash mapping for resource")
	}
}

// DeleteResourceHash removes the number-index to hash canonical mapping.
func DeleteResourceHash(db DatabaseDeleter, number uint64, index uint64) {
	if err := db.Delete(resourceHashKey(number, index)); err != nil {
		log.WithError(err).Fatal("Failed to delete number-index to hash mapping for resource")
	}
}

// ReadResourceId retrieves the dataId assigned to a canonical nodeId and hash.
func ReadResourceId(db DatabaseReader, nodeId string, hash common.Hash) []byte {
	data, _ := db.Get(resourceDataIdKey(common.Hex2Bytes(nodeId), hash))
	if len(data) == 0 {
		return []byte{}
	}
	return data
}

// WriteResourceId stores the dataId assigned to a canonical nodeId and hash.
func WriteResourceId(db DatabaseWriter, nodeId string, hash common.Hash, dataId string)  {
	if err := db.Put(resourceDataIdKey(common.Hex2Bytes(nodeId), hash), common.Hex2Bytes(dataId)); err != nil {
		log.WithError(err).Fatal("Failed to store nodeId-hash to dataId mapping for resource")
	}
}

// DeleteResourceId removes the nodeId-hash to dataId canonical mapping.
func DeleteResourceId(db DatabaseDeleter, nodeId string, hash common.Hash,) {
	if err := db.Delete(resourceDataIdKey(common.Hex2Bytes(nodeId), hash)); err != nil {
		log.WithError(err).Fatal("Failed to delete number-index to hash mapping for resource")
	}
}

// ReadResourceTypeHash retrieves the hash assigned to a canonical type and dataId.
func ReadResourceTypeHash(db DatabaseReader, dataId string, typ string) common.Hash {
	data, _ := db.Get(resourceDataTypeHashKey(common.Hex2Bytes(dataId), common.Hex2Bytes(typ)))
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteResourceTypeHash stores the hash assigned to a canonical type and dataId.
func WriteResourceTypeHash(db DatabaseWriter, dataId string, typ string, hash common.Hash)  {
	if err := db.Put(resourceDataTypeHashKey(common.Hex2Bytes(dataId), common.Hex2Bytes(typ)), hash.Bytes()); err != nil {
		log.WithError(err).Fatal("Failed to store type-dataId to hash mapping for resource")
	}
}

// DeleteResourceTypeHash removes the type-dataId to hash canonical mapping.
func DeleteResourceTypeHash(db DatabaseDeleter, dataId string, typ string) {
	if err := db.Delete(resourceDataTypeHashKey(common.Hex2Bytes(dataId), common.Hex2Bytes(typ))); err != nil {
		log.WithError(err).Fatal("Failed to delete type-dataId to hash mapping for resource")
	}
}
