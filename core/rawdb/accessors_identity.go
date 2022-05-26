// Copyright (C) 2021 The RosettaNet Authors.

package rawdb

import "github.com/datumtechs/datum-network-carrier/common"

// ReadIdentityHash retrieves the hash assigned to a canonical number and index.
func ReadIdentityHash(db DatabaseReader, number uint64, index uint64) common.Hash {
	data, _ := db.Get(identityHashKey(number, index))
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteIdentityHash stores the hash assigned to a canonical number and index.
func WriteIdentityHash(db DatabaseWriter, number uint64, index uint64, hash common.Hash)  {
	if err := db.Put(identityHashKey(number, index), hash.Bytes()); err != nil {
		log.WithError(err).Fatal("Failed to store number-index to hash mapping for identity")
	}
}

// DeleteIdentityHash removes the number-index to hash canonical mapping.
func DeleteIdentityHash(db DatabaseDeleter, number uint64, index uint64) {
	if err := db.Delete(identityHashKey(number, index)); err != nil {
		log.WithError(err).Fatal("Failed to delete number-index to hash mapping for identity")
	}
}

// ReadIdentityId retrieves the dataId assigned to a canonical nodeId and hash.
func ReadIdentityId(db DatabaseReader, nodeId string, hash common.Hash) []byte {
	data, _ := db.Get(identityDataIdKey(common.Hex2Bytes(nodeId), hash))
	if len(data) == 0 {
		return []byte{}
	}
	return data
}

// WriteIdentityId stores the dataId assigned to a canonical nodeId and hash.
func WriteIdentityId(db DatabaseWriter, nodeId string, hash common.Hash, dataId string)  {
	if err := db.Put(identityDataIdKey(common.Hex2Bytes(nodeId), hash), common.Hex2Bytes(dataId)); err != nil {
		log.WithError(err).Fatal("Failed to store nodeId-hash to dataId mapping for identity")
	}
}

// DeleteIdentityId removes the nodeId-hash to dataId canonical mapping.
func DeleteIdentityId(db DatabaseDeleter, nodeId string, hash common.Hash,) {
	if err := db.Delete(identityDataIdKey(common.Hex2Bytes(nodeId), hash)); err != nil {
		log.WithError(err).Fatal("Failed to delete number-index to hash mapping for identity")
	}
}

// ReadIdentityTypeHash retrieves the hash assigned to a canonical type and dataId.
func ReadIdentityTypeHash(db DatabaseReader, dataId string, typ string) common.Hash {
	data, _ := db.Get(identityDataTypeHashKey(common.Hex2Bytes(dataId), common.Hex2Bytes(typ)))
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteIdentityTypeHash stores the hash assigned to a canonical type and dataId.
func WriteIdentityTypeHash(db DatabaseWriter, dataId string, typ string, hash common.Hash)  {
	if err := db.Put(identityDataTypeHashKey(common.Hex2Bytes(dataId), common.Hex2Bytes(typ)), hash.Bytes()); err != nil {
		log.WithError(err).Fatal("Failed to store type-dataId to hash mapping for identity")
	}
}

// DeleteIdentityTypeHash removes the type-dataId to hash canonical mapping.
func DeleteIdentityTypeHash(db DatabaseDeleter, dataId string, typ string) {
	if err := db.Delete(identityDataTypeHashKey(common.Hex2Bytes(dataId), common.Hex2Bytes(typ))); err != nil {
		log.WithError(err).Fatal("Failed to delete type-dataId to hash mapping for identity")
	}
}