package rawdb

import (
	"encoding/binary"
	"github.com/RosettaFlow/Carrier-Go/common"
)

// The fields below define the low level database schema prefixing.
var (
	// databaseVersionKey tracks the current database version
	databaseVersionKey = []byte("DatabaseVersion")

	// headHeaderKey tracks the latest know header's hash
	headHeaderKey = []byte("LastHeader")

	// headBlockKey tracks the latest know full block's hash.
	headBlockKey = []byte("LastBlock")

	// Data item prefixes
	headerPrefix 		= []byte("h")	// headerPrefix + num (uint64 big endian) + hash -> header
	headerHashSuffix 	= []byte("n")	// headerPrefix + num (uint64 big endian) + headerHashSuffix -> hash
	headerNumberPrefix 	= []byte("H")	// headerNumberPrefix + hash -> num (uint64 big endian), headerNumberPrefix + nodeId + hash -> num (uint64 big endian)

	blockBodyPrefix		= []byte("b")	// blockBodyPrefix + num(uint64 big endian) + hash -> block body

	// data item prefixes for metadata
	metadataIdPrefix			= []byte("mh")
	metadataIdSuffix			= []byte("n")	// metadataIdPrefix + num + index + metadataIdSuffix -> metadata dataId
	metadataHashPrefix			= []byte("mn")	// metadataHashPrefix + nodeId + type -> metadata hash
	metadataTypeHashPrefix		= []byte("md")	// metadataTypeHashPrefix + type + dataId -> metadata hash
	metadataApplyHashPrefix 	= []byte("ma")	// metadataApplyHashPrefix + type + dataId -> apply metadata hash

	resourceDataIdPrefix 		= []byte("rh")
	resourceDataIdSuffix 		= []byte("n")	// resourceDataIdPrefix + num + index + resourceDataIdSuffix -> resource dataId
	resourceDataHashPrefix		= []byte("rn")	// resourceDataHashPrefix + nodeId + type -> resource hash
	resourceDataTypeHashPrefix 	= []byte("rd")	// resourceDataTypeHashPrefix + type + dataId -> resource hash

	identityDataIdPrefix 		= []byte("ch")
	identityDataIdSuffix 		= []byte("n")	// identityDataIdPrefix + num + index + identityDataIdSuffix -> identity dataId
	identityDataHashPrefix		= []byte("cn")	// identityDataHashPrefix + nodeId + type -> identity hash
	identityDataTypeHashPrefix 	= []byte("cb")	// identityDataTypeHashPrefix + type + dataId -> identity hash

	taskDataIdPrefix 			= []byte("th")
	taskDataIdSuffix 			= []byte("n")	// taskDataIdPrefix + num + index + taskDataIdSuffix -> task dataId
	taskDataHashPrefix			= []byte("tn")	// taskDataHashPrefix + nodeId + type -> task hash
	taskDataTypeHashPrefix 		= []byte("tb")	// taskDataTypeHashPrefix + type + dataId -> task hash
)

// encodeNumber encodes a number as big endian uint64
func encodeNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// headerKey = headerPrefix + num (uint64 big endian) + hash
func headerKey(number uint64, hash common.Hash) []byte {
	return append(append(headerPrefix, encodeNumber(number)...), hash.Bytes()...)
}

// headerHashKey = headerPrefix + num + headerHashSuffix
func headerHashKey(number uint64) []byte {
	return append(append(headerPrefix, encodeNumber(number)...), headerHashSuffix...)
}

// headerNumberKey = headerNumberPrefix + hash
func headerNumberKey(hash common.Hash) []byte {
	return append(headerNumberPrefix, hash.Bytes()...)
}

// blockBodyKey = blockBodyPrefix + num + hash
func blockBodyKey(number uint64, hash common.Hash) []byte {
	return append(append(blockBodyPrefix, encodeNumber(number)...), hash.Bytes()...)
}

// metadataIdKey = metadataIdPrefix + num + index + metadataIdSuffix
func metadataIdKey(number uint64, index uint64) []byte {
	return append(append(append(metadataIdPrefix, encodeNumber(number)...), encodeNumber(index)...), metadataIdSuffix...)
}

// metadataHashKey = metadataHashPrefix + nodeId + type
func metadataHashKey(nodeId []byte, typ []byte) []byte {
	return append(append(metadataHashPrefix, nodeId...), typ...)
}

// metadataTypeHashKey = metadataTypeHashPrefix + type + dataId
func metadataTypeHashKey(dataId []byte, typ []byte) []byte  {
	return append(append(metadataTypeHashPrefix, typ...), dataId...)
}

// metadataApplyHashKey = metadataApplyHashPrefix + type + dataId
func metadataApplyHashKey(tpy []byte, dataId []byte) []byte {
	return append(append(metadataApplyHashPrefix, tpy...), dataId...)
}

// resourceDataIdKey = resourceDataIdPrefix + num + index + resourceDataIdSuffix
func resourceDataIdKey(number uint64, index uint64) []byte {
	return append(append(append(resourceDataIdPrefix, encodeNumber(number)...), encodeNumber(index)...), resourceDataIdSuffix...)
}

// resourceDataHashKey = resourceDataHashPrefix + nodeId + type
func resourceDataHashKey(nodeId []byte, typ []byte) []byte {
	return append(append(resourceDataHashPrefix, nodeId...), typ...)
}

// resourceDataTypeHashKey = resourceDataTypeHashPrefix + type + dataId
func resourceDataTypeHashKey(dataId []byte, typ []byte) []byte  {
	return append(append(resourceDataTypeHashPrefix, typ...), dataId...)
}

// identityDataIdKey = identityDataIdPrefix + num + index + identityDataIdSuffix
func identityDataIdKey(number uint64, index uint64) []byte {
	return append(append(append(identityDataIdPrefix, encodeNumber(number)...), encodeNumber(index)...), identityDataIdSuffix...)
}

// identityDataHashKey = identityDataHashPrefix + nodeId + type
func identityDataHashKey(nodeId []byte, typ []byte) []byte {
	return append(append(identityDataHashPrefix, nodeId...), typ...)
}

// identityDataTypeHashKey = identityDataTypeHashPrefix + type + dataId
func identityDataTypeHashKey(dataId []byte, typ []byte) []byte  {
	return append(append(identityDataTypeHashPrefix, typ...), dataId...)
}

// taskDataIdKey = taskDataIdPrefix + num + index + taskDataIdSuffix
func taskDataIdKey(number uint64, index uint64) []byte {
	return append(append(append(taskDataIdPrefix, encodeNumber(number)...), encodeNumber(index)...), taskDataIdSuffix...)
}

// taskDataHashKey = taskDataHashPrefix + nodeId + type
func taskDataHashKey(nodeId []byte, typ []byte) []byte {
	return append(append(taskDataHashPrefix, nodeId...), typ...)
}

// taskDataTypeHashKey = taskDataTypeHashPrefix + type + dataId
func taskDataTypeHashKey(dataId []byte, typ []byte) []byte  {
	return append(append(taskDataTypeHashPrefix, typ...), dataId...)
}