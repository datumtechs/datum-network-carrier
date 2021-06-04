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

func headerKey(number uint64, hash common.Hash)  {
	
}

/*// TxLookupEntry is a positional metadata to help looking up the data content of

// headerKey = headerPrefix + num (uint64 big endian) + hash
func headerKey(number uint64, hash common.Hash) []byte {
	return append(append(headerPrefix, encodeBlockNumber(number)...), hash.Bytes()...)
}

// headerHashKey = headerPrefix + num (uint64 big endian) + headerHashSuffix
func headerHashKey(number uint64) []byte {
	return append(append(headerPrefix, encodeBlockNumber(number)...), headerHashSuffix...)
}

// headerNumberKey = headerNumberPrefix + hash
func headerNumberKey(hash common.Hash) []byte {
	return append(headerNumberPrefix, hash.Bytes()...)
}

// blockBodyKey = blockBodyPrefix + num (uint64 big endian) + hash
func blockBodyKey(number uint64, hash common.Hash) []byte {
	return append(append(blockBodyPrefix, encodeBlockNumber(number)...), hash.Bytes()...)
}

// blockReceiptsKey = blockReceiptsPrefix + num (uint64 big endian) + hash
func blockReceiptsKey(number uint64, hash common.Hash) []byte {
	return append(append(blockReceiptsPrefix, encodeBlockNumber(number)...), hash.Bytes()...)
}

// blockConfirmSignsKey = blockConfirmSignsPrefix + num (uint64 big endian) + hash
func blockConfirmSignsKey(number uint64, hash common.Hash) []byte {
	return append(append(blockConfirmSignsPrefix, encodeBlockNumber(number)...), hash.Bytes()...)
}

// txLookupKey = txLookupPrefix + hash
func txLookupKey(hash common.Hash) []byte {
	return append(txLookupPrefix, hash.Bytes()...)
}

// bloomBitsKey = bloomBitsPrefix + bit (uint16 big endian) + section (uint64 big endian) + hash
func bloomBitsKey(bit uint, section uint64, hash common.Hash) []byte {
	key := append(append(bloomBitsPrefix, make([]byte, 10)...), hash.Bytes()...)

	binary.BigEndian.PutUint16(key[1:], uint16(bit))
	binary.BigEndian.PutUint64(key[3:], section)

	return key
}

// preimageKey = preimagePrefix + hash
func preimageKey(hash common.Hash) []byte {
	return append(preimagePrefix, hash.Bytes()...)
}

// configKey = configPrefix + hash
func configKey(hash common.Hash) []byte {
	return append(configPrefix, hash.Bytes()...)
}

// economicModelKey = economicModelPrefix + hash
func economicModelKey(hash common.Hash) []byte {
	return append(economicModelPrefix, hash.Bytes()...)
}*/
