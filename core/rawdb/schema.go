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

	// seedNodeKey tracks the seed node list.
	seedNodeKey = []byte("SeedNodeKey")

	// resourceNodeKey tracks the node for calc node.
	calcNodeKey = []byte("CalcNodeKey")

	dataNodeKey = []byte("DataNodeKey")

	// Data item prefixes
	headerPrefix 		= []byte("h")	// headerPrefix + num (uint64 big endian) + hash -> header
	headerHashSuffix 	= []byte("n")	// headerPrefix + num (uint64 big endian) + headerHashSuffix -> hash
	headerNumberPrefix 	= []byte("H")	// headerNumberPrefix + hash -> num (uint64 big endian), headerNumberPrefix + nodeId + hash -> num (uint64 big endian)

	blockBodyPrefix		= []byte("b")	// blockBodyPrefix + num(uint64 big endian) + hash -> block body

	// data item prefixes for metadata
	metadataHashPrefix			= []byte("mh")
	metadataHashSuffix			= []byte("n")	// metadataHashPrefix + num + index + metadataHashSuffix -> metadata hash
	metadataIdPrefix			= []byte("mn")	// metadataIdPrefix + nodeId + hash -> metadata dataId
	metadataTypeHashPrefix		= []byte("md")	// metadataTypeHashPrefix + type + dataId -> metadata hash
	metadataApplyHashPrefix 	= []byte("ma")	// metadataApplyHashPrefix + type + dataId -> apply metadata hash

	resourceHashPrefix 			= []byte("rh")
	resourceHashSuffix 			= []byte("n")	// resourceHashPrefix + num + index + resourceHashSuffix -> resource hash
	resourceDataIdPrefix		= []byte("rn")	// resourceDataHashPrefix + nodeId + hash -> resource dataId
	resourceDataTypeHashPrefix 	= []byte("rd")	// resourceDataTypeHashPrefix + type + dataId -> resource hash

	identityHashPrefix 			= []byte("ch")
	identityHashSuffix 			= []byte("n")	// identityHashPrefix + num + index + identityHashSuffix -> identity hash
	identityDataIdPrefix		= []byte("cn")	// identityDataIdPrefix + nodeId + hash -> identity dataId
	identityDataTypeHashPrefix 	= []byte("cb")	// identityDataTypeHashPrefix + type + dataId -> identity hash

	taskDataHashPrefix 			= []byte("th")
	taskDataHashSuffix 			= []byte("n")	// taskDataHashPrefix + num + index + taskDataHashSuffix -> task hash
	taskDataIdPrefix			= []byte("tn")	// taskDataIdPrefix + nodeId + hash -> task dataId
	taskDataTypeHashPrefix 		= []byte("tb")	// taskDataTypeHashPrefix + type + dataId -> task hash

	dataLookupPrefix  			= []byte("l") 	// dataLookupPrefix + dataId -> metadata lookup metadata
	//resourceLookupPrefix  		= []byte("r") 	// resourceLookupPrefix + dataId -> resource lookup metadata
	//identityLookupPrefix  		= []byte("i") 	// identityLookupPrefix + dataId -> identity lookup metadata
	//taskLookupPrefix  			= []byte("t") 	// taskLookupPrefix + dataId -> task lookup metadata
)

// headerKeyPrefix = headerPrefix + num (uint64 big endian)
func headerKeyPrefix(number uint64) []byte {
	return append(headerPrefix, encodeNumber(number)...)
}

// metadataIdKeyPrefix = metadataIdPrefix + nodeId
func metadataIdKeyPrefix(nodeId string) []byte {
	return append(metadataIdPrefix, common.Hex2Bytes(nodeId)...)
}

// encodeNumber encodes a number as big endian uint64
func encodeNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// dataLookupKey = dataLookupPrefix + hash
func dataLookupKey(hash common.Hash) []byte {
	return append(dataLookupPrefix, hash.Bytes()...)
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

// metadataHashKey = metadataHashPrefix + num + index + metadataHashSuffix
func metadataHashKey(number uint64, index uint64) []byte {
	return append(append(append(metadataHashPrefix, encodeNumber(number)...), encodeNumber(index)...), metadataHashSuffix...)
}

// metadataIdKey = metadataIdPrefix + nodeId + hash
func metadataIdKey(nodeId []byte, hash common.Hash) []byte {
	return append(append(metadataIdPrefix, nodeId...), hash.Bytes()...)
}

// metadataTypeHashKey = metadataTypeHashPrefix + type + dataId
func metadataTypeHashKey(dataId []byte, typ []byte) []byte  {
	return append(append(metadataTypeHashPrefix, typ...), dataId...)
}

// metadataApplyHashKey = metadataApplyHashPrefix + type + dataId
func metadataApplyHashKey(dataId []byte, tpy []byte) []byte {
	return append(append(metadataApplyHashPrefix, tpy...), dataId...)
}

// resourceHashKey = resourceHashPrefix + num + index + resourceHashSuffix
func resourceHashKey(number uint64, index uint64) []byte {
	return append(append(append(resourceHashPrefix, encodeNumber(number)...), encodeNumber(index)...), resourceHashSuffix...)
}

// resourceDataIdKey = resourceDataIdKey + nodeId + hash
func resourceDataIdKey(nodeId []byte, hash common.Hash) []byte {
	return append(append(resourceDataIdPrefix, nodeId...), hash.Bytes()...)
}

// resourceDataTypeHashKey = resourceDataTypeHashPrefix + type + dataId
func resourceDataTypeHashKey(dataId []byte, typ []byte) []byte  {
	return append(append(resourceDataTypeHashPrefix, typ...), dataId...)
}

// identityHashKey = identityHashPrefix + num + index + identityHashSuffix
func identityHashKey(number uint64, index uint64) []byte {
	return append(append(append(identityHashPrefix, encodeNumber(number)...), encodeNumber(index)...), identityHashSuffix...)
}

// identityDataIdKey = identityDataIdPrefix + nodeId + hash
func identityDataIdKey(nodeId []byte, hash common.Hash) []byte {
	return append(append(identityDataIdPrefix, nodeId...), hash.Bytes()...)
}

// identityDataTypeHashKey = identityDataTypeHashPrefix + type + dataId
func identityDataTypeHashKey(dataId []byte, typ []byte) []byte  {
	return append(append(identityDataTypeHashPrefix, typ...), dataId...)
}

// taskDataHashKey = taskDataHashPrefix + num + index + taskDataHashSuffix
func taskDataHashKey(number uint64, index uint64) []byte {
	return append(append(append(taskDataHashPrefix, encodeNumber(number)...), encodeNumber(index)...), taskDataHashSuffix...)
}

// taskDataIdKey = taskDataIdPrefix + nodeId + type
func taskDataIdKey(nodeId []byte, typ []byte) []byte {
	return append(append(taskDataIdPrefix, nodeId...), typ...)
}

// taskDataTypeHashKey = taskDataTypeHashPrefix + type + dataId
func taskDataTypeHashKey(dataId []byte, typ []byte) []byte  {
	return append(append(taskDataTypeHashPrefix, typ...), dataId...)
}