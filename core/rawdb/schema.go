package rawdb


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
