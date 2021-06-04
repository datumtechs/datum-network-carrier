package rawdb

import (
	"encoding/binary"
	"github.com/RosettaFlow/Carrier-Go/common"
)

// ReadDataHash retrieves the hash assigned to a canonical block number
func ReadDataHash(db DatabaseReader, number uint64) common.Hash {
	data, _ := db.Get(headerHashKey(number))
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteDataHash stores the hash assigned to a canonical block number
func WriteDataHash(db DatabaseWriter, hash common.Hash, number uint64)  {
	if err := db.Put(headerHashKey(number), hash.Bytes()); err != nil {
		log.WithError(err).Error("Failed to store number to hash mapping")
	}
}

// DeleteDataHash removes the number to hash canonical mapping.
func DeleteDataHash(db DatabaseDeleter, number uint64) {
	if err := db.Delete(headerHashKey(number)); err != nil {
		log.WithError(err).Error("Failed to delete number to hash mapping")
	}
}

// ReadHeaderNumber returns the header number assigned to a hash.
func ReadHeaderNumber(db DatabaseReader, hash common.Hash) *uint64 {
	data, _ := db.Get(headerNumberKey(hash))
	if len(data) != 8 {
		return nil
	}
	number := binary.BigEndian.Uint64(data)
	return &number
}

// ReadHeadHeaderHash retrieves the hash of the current canonical head header.
func ReadHeadHeaderHash(db DatabaseReader) common.Hash {
	data, _ := db.Get(headHeaderKey)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteHeadHeaderHash stores the hash of the current canonical head header.
func WriteHeadHeaderHash(db DatabaseWriter, hash common.Hash) {
	if err := db.Put(headHeaderKey, hash.Bytes()); err != nil {
		log.WithError(err).Error("Failed to store last header's hash")
	}
}

// ReadHeadBlockHash retrieves the hash of the current canonical head block.
func ReadHeadBlockHash(db DatabaseReader) common.Hash {
	data, _ := db.Get(headBlockKey)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteHeadBlockHash stores the head block's hash.
func WriteHeadBlockHash(db DatabaseWriter, hash common.Hash) {
	if err := db.Put(headBlockKey, hash.Bytes()); err != nil {
		log.WithError(err).Error("Failed to store last block's hash")
	}
}
