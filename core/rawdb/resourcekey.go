package rawdb

import (
	"github.com/RosettaFlow/Carrier-Go/common/bytesutil"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
)

var (
	// prefix + jobNodeId -> LocalResourceTable
	nodeResourceKeyPrefix = []byte("NodeResourceKeyPrefix:")
	// key -> [jobNodeId, jobNodeId, ..., jobNodeId]
	nodeResourceIdListKey = []byte("nodeResourceIdListKey")
	// prefix + identityId -> RemoteResourceTable
	orgResourceKeyPrefix = []byte("OrgResourceKeyPrefix:")
	// key -> [identityId, identityId, ..., identityId]
	orgResourceIdListKey = []byte("OrgResourceIdListKey")
	// key -> SlotUnit
	nodeResourceSlotUnitKey = []byte("nodeResourceSlotUnitKey")
	// prefix + taskId -> LocalTaskPowerUsed
	localTaskPowerUsedKeyPrefix = []byte("localTaskPowerUsedKeyPrefix:")
	// key -> [taskId, taskId, ..., taskId]
	localTaskPowerUsedIdListKey = []byte("localTaskPowerUsedIdListKey")
	// prefix + dataNodeId -> DataResourceTable{dataNodeId, totalDisk, usedDisk}
	dataResourceTableKeyPrefix = []byte("dataResourceTableKeyPrefix:")
	// key -> [dataNodeId, dataNodeId, ..., dataNodeId]
	dataResourceTableIdListKey = []byte("dataResourceTableIdListKey")
	// prefix + originId -> DataResourceFileUpload{originId, dataNodeId, metaDataId, filePath}
	dataResourceFileUploadKeyPrefix = []byte("dataResourceDataUsedKeyPrefix:")
	// key -> [originId, originId, ..., originId]
	dataResourceFileUploadIdListKey = []byte("dataResourceFileUploadIdListKey")
	// prefix + jonNodeId -> [taskId, taskId, ..., taskId]
	resourceTaskIdsKeyPrefix = []byte("resourceTaskIdsKeyPrefix:")
	// prefix + powerId -> jobNodeId
	resourcePowerIdMapingKeyPrefix = []byte("resourcePowerIdMapingKeyPrefix:")
	// prefix + metaDataId -> DataResourceDiskUsed{metaDataId, dataNodeId, diskUsed}
	dataResourceDiskUsedKeyPrefix = []byte("DataResourceDiskUsedKeyPrefix:")
	// prefix + taskId -> executeStatus
	localTaskExecuteStatusKeyPrefix = []byte("localTaskExecuteStatusKeyPrefix:")

	// prefix + userType + user -> n
	userMetadataAuthUsedCountKey = []byte("userMetadataAuthUsedCountKey")
	// prefix + userType + user + n -> metadataId
	userMetadataAuthUsedKeyPrefix = []byte("userMetadataAuthUsedKeyPrefix:")
	// prefix + userType + user + metadataId -> metadataAuthId
	userMetadataAuthByMetadataIdKeyPrefix = []byte("userMetadataAuthByMetadataIdKeyPrefix:")

	// metadataId -> taskCount (n)
	metadataUsedTaskIdCountKey = []byte("metadataUsedTaskIdCountKey")
	// metadataId + n -> taskId
	metadataUsedTaskIdKeyPrefix = []byte("metadataUsedTaskIdKeyPrefix:")

	// prefix + taskId -> n todo maybe support multi-result file about a taskId
	// prefix + taskId + n -> resultfile summary (auto build metadataId)

	// taskId -> resultfile summary (auto build metadataId)
	taskResultFileMetadataIdKeyPrefix = []byte("taskResultFileMetadataIdKeyPrefix:")
)

// nodeResourceKey = NodeResourceKeyPrefix + jobNodeId
func GetNodeResourceKey(jobNodeId string) []byte {
	return append(nodeResourceKeyPrefix, []byte(jobNodeId)...)
}
func GetNodeResourceIdListKey() []byte {
	return nodeResourceIdListKey
}
func GetOrgResourceKey(identityId string) []byte {
	return append(orgResourceKeyPrefix, []byte(identityId)...)
}
func GetOrgResourceIdListKey() []byte {
	return orgResourceIdListKey
}
func GetNodeResourceSlotUnitKey() []byte {
	return nodeResourceSlotUnitKey
}

func GetLocalTaskPowerUsedKey(taskId string) []byte {
	return append(localTaskPowerUsedKeyPrefix, []byte(taskId)...)
}
func GetLocalTaskPowerUsedIdListKey() []byte {
	return localTaskPowerUsedIdListKey
}

func GetDataResourceTableKey(dataNodeId string) []byte {
	return append(dataResourceTableKeyPrefix, []byte(dataNodeId)...)
}
func GetDataResourceTableIdListKey() []byte {
	return dataResourceTableIdListKey
}

func GetDataResourceFileUploadKey(originId string) []byte {
	return append(dataResourceFileUploadKeyPrefix, []byte(originId)...)
}
func GetDataResourceFileUploadIdListKey() []byte {
	return dataResourceFileUploadIdListKey
}

func GetResourceTaskIdsKey(jonNodeId string) []byte {
	return append(resourceTaskIdsKeyPrefix, []byte(jonNodeId)...)
}

func GetResourcePowerIdMapingKey(powerId string) []byte {
	return append(resourcePowerIdMapingKeyPrefix, []byte(powerId)...)
}

//func GetResourceMetadataIdMapingKey(powerId string) []byte {
//	return append(resourceMetadataIdMapingKeyPrefix, []byte(powerId)...)
//}

func GetDataResourceDiskUsedKey(metaDataId string) []byte {
	return append(dataResourceDiskUsedKeyPrefix, []byte(metaDataId)...)
}

func GetLocalTaskExecuteStatus(taskId string) []byte {
	return append(localTaskExecuteStatusKeyPrefix, []byte(taskId)...)
}

func GetUserMetadataAuthUsedCountKey(userType apicommonpb.UserType, user string) []byte {
	return append(append(userMetadataAuthUsedCountKey, []byte(userType.String())...), []byte(user)...)
}

func GetUserMetadataAuthUsedKey(userType apicommonpb.UserType, user string, n uint32) []byte {

	// key: prefix + userType + user + n

	userTypeBytes := []byte(userType.String())
	userBytes := []byte(user)
	nBytes := bytesutil.Uint32ToBytes(n)

	// some index of pivots
	prefixIndex := len(userMetadataAuthUsedKeyPrefix)
	userTypeIndex := prefixIndex + len(userTypeBytes)
	userIndex := userTypeIndex + len(userBytes)
	size := userIndex + len(nBytes)

	// construct key
	key := make([]byte, size)
	copy(key[:prefixIndex], userMetadataAuthUsedKeyPrefix)
	copy(key[prefixIndex:userTypeIndex], userTypeBytes)
	copy(key[userTypeIndex:userIndex], userBytes)
	copy(key[userIndex:], nBytes)

	return key
}

func GetUserMetadataAuthByMetadataIdKey(userType apicommonpb.UserType, user, metadataId string) []byte {

	// key: prefix + userType + user + metadataId

	userTypeBytes := []byte(userType.String())
	userBytes := []byte(user)
	metadataIdBytes := []byte(metadataId)

	// some index of pivots
	prefixIndex := len(userMetadataAuthByMetadataIdKeyPrefix)
	userTypeIndex := prefixIndex + len(userTypeBytes)
	userIndex := userTypeIndex + len(userBytes)
	size := userIndex + len(metadataIdBytes)

	// construct key
	key := make([]byte, size)
	copy(key[:prefixIndex], userMetadataAuthByMetadataIdKeyPrefix)
	copy(key[prefixIndex:userTypeIndex], userTypeBytes)
	copy(key[userTypeIndex:userIndex], userBytes)
	copy(key[userIndex:], metadataIdBytes)

	return key
}

func GetMetadataUsedTaskIdCountKey(metadataId string) []byte {
	return append(metadataUsedTaskIdCountKey, []byte(metadataId)...)
}

func GetMetadataUsedTaskIdKey(metadataId string, n uint32) []byte {

	// key: prefix + metadataId + n

	metadataIdBytes := []byte(metadataId)
	nBytes := bytesutil.Uint32ToBytes(n)

	// some index of pivots
	prefixIndex := len(metadataUsedTaskIdKeyPrefix)
	metadataIdIndex := prefixIndex + len(metadataIdBytes)
	size := metadataIdIndex + len(nBytes)

	// construct key
	key := make([]byte, size)
	copy(key[:prefixIndex], userMetadataAuthUsedKeyPrefix)
	copy(key[prefixIndex:metadataIdIndex], metadataIdBytes)
	copy(key[metadataIdIndex:], nBytes)

	return key
}

func GetTaskResultFileMetadataIdKey(taskId string) []byte {
	return append(taskResultFileMetadataIdKeyPrefix, []byte(taskId)...)
}
