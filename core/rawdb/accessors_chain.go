// Copyright (C) 2021 The RosettaNet Authors.

package rawdb

import (
	"encoding/binary"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/datumtechs/datum-network-carrier/db"
	libtypes "github.com/datumtechs/datum-network-carrier/lib/types"
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
		log.WithError(err).Fatal("Failed to store number to hash mapping")
	}
}

// DeleteDataHash removes the number to hash canonical mapping.
func DeleteDataHash(db DatabaseDeleter, number uint64) {
	if err := db.Delete(headerHashKey(number)); err != nil {
		log.WithError(err).Fatal("Failed to delete number to hash mapping")
	}
}

// ReadAllHashes retrieves all the hashes assigned to blocks at a certain heights
func ReadAllHashes(db db.Iteratee, number uint64) []common.Hash {
	prefix := headerKeyPrefix(number)
	hashes := make([]common.Hash, 0, 1)
	it := db.NewIteratorWithPrefixAndStart(prefix, nil)
	defer it.Release()

	for it.Next() {
		if key := it.Key(); len(key) == len(prefix) + 32 {
			hashes = append(hashes, common.BytesToHash(key[len(key)-32:]))
		}
	}
	return hashes
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
		log.WithError(err).Fatal("Failed to store last header's hash")
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
		log.WithError(err).Fatal("Failed to store last block's hash")
	}
}

// ReadHeaderPB retrieves a block header in its raw RLP database encoding.
func ReadHeaderPB(db DatabaseReader, hash common.Hash, number uint64) []byte {
	data, _ := db.Get(headerKey(number, hash))
	return data
}

// HasHeader verifies the existence of a block header corresponding to the hash.
func HasHeader(db DatabaseReader, hash common.Hash, number uint64) bool {
	if has, err := db.Has(headerKey(number, hash)); !has || err != nil {
		return false
	}
	return true
}

// ReadHeader retrieves the block header corresponding to the hash.
func ReadHeader(db DatabaseReader, hash common.Hash, number uint64) *types.Header {
	data := ReadHeaderPB(db, hash, number)
	if len(data) == 0 {
		return nil
	}
	header := new(libtypes.HeaderPb)
	if err := header.Unmarshal(data); err != nil {
		log.WithField("hash", hash).WithError(err).Fatal("Invalid block header ProtoBuf")
		return nil
	}
	// todo: To be continued..
	return nil
}

// WriteHeader stores a block header into the database and also stores the hash-to-number mapping.
func WriteHeader(db DatabaseWriter, header *types.Header) {
	var (
		hash 	= header.Hash()
		number 	= header.Version
		encoded = encodeNumber(number)
	)
	key := headerNumberKey(hash)
	if err := db.Put(key, encoded); err != nil {
		log.WithError(err).Fatal("Failed to store hash to number mapping")
	}
	// todo: To be continued..
	//data, err := header.Marshal()
	//if err != nil {
	//	log.WithError(err).Error("Failed to PB encode header")
	//}
	//key  = headerKey(number, hash)
	//if err := db.Put(key, data); err != nil {
	//	log.WithError(err).Error("Failed to store header")
	//}
}

// DeleteHeader remove all block header data associated with a hash.
func DeleteHeader(db DatabaseDeleter, hash common.Hash, number uint64) {
	if err := db.Delete(headerKey(number, hash)); err != nil {
		log.WithError(err).Fatal("Failed to delete header")
	}
	if err := db.Delete(headerNumberKey(hash)); err != nil {
		log.WithError(err).Fatal("Failed to delete hash to number mapping")
	}
}

// ReadBodyPB retrieves the block body (metadata/resource/identity/task) in Protobuf encoding.
func ReadBodyPB(db DatabaseReader, hash common.Hash, number uint64) []byte {
	data, _ := db.Get(blockBodyKey(number, hash))
	return data
}

// WriteBodyPB stores an Protobuf encoded block body into the database.
func WriteBodyPB(db DatabaseWriter, hash common.Hash, number uint64, pb []byte) {
	if err := db.Put(blockBodyKey(number, hash), pb); err != nil {
		log.WithError(err).Fatal("Failed to store block body")
	}
}

// HasBody verifies the existence of a block body corresponding to the hash.
func HasBody(db DatabaseReader, hash common.Hash, number uint64) bool {
	if has, err := db.Has(blockBodyKey(number, hash)); !has || err != nil {
		return false
	}
	return true
}

// ReadBody retrieves the block body corresponding to the hash.
func ReadBody(db DatabaseReader, hash common.Hash, number uint64) *libtypes.BodyData {
	data := ReadBodyPB(db, hash, number)
	if len(data) == 0 {
		return nil
	}
	body := new(libtypes.BodyData)
	if err := body.Unmarshal(data); err != nil {
		log.WithField("hash", hash).WithError(err).Error("Invalid block body Protobuf")
		return nil
	}
	return body
}

// WriteBOdy store a block body into the database.
func WriteBody(db DatabaseWriter, hash common.Hash, number uint64, body *libtypes.BodyData) {
	data, err := body.Marshal()
	if err != nil {
		log.WithError(err).Fatal("Failed to ProtoBuf encode body")
	}
	WriteBodyPB(db, hash, number, data)
	// todo: Data analysis is required when storing the body.
	// metadata/resource/identity/...
}

// DeleteBody removes all block body data associated with a hash.
func DeleteBody(db DatabaseDeleter, hash common.Hash, number uint64) {
	if err := db.Delete(blockBodyKey(number, hash)); err != nil {
		log.WithError(err).Fatal("Failed to delete block body")
	}
	// todo: remove metadata/resource/identity/task at the same time.
}

func ReadBlock(db DatabaseReader, hash common.Hash, number uint64) *types.Block {
	header := ReadHeader(db, hash , number)
	if header == nil {
		return nil
	}
	body := ReadBody(db, hash, number)
	if body == nil {
		return nil
	}
	// todo: generated block by header and body.
	// return types.NewBlockWithHeader(header).WithBody(body.Transactions, body.ExtraData)
	return nil
}

// WriteBlock serializes a block into the database, header and body separately.
func WriteBlock(db DatabaseWriter, block *libtypes.BlockData) {
	//WriteBody(db, block.Hash(), block.NumberU64(), block.Body())
	//WriteHeader(db, &block.Header)
}

// DeleteBlock removes all block data associated with a hash.
func DeleteBlock(db DatabaseDeleter, hash common.Hash, number uint64) {
	DeleteHeader(db, hash, number)
	DeleteBody(db, hash, number)
}