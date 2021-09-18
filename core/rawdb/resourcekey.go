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

	userTypeBytes := []byte(userType.String())
	userBytes := []byte(user)
	nBytes := bytesutil.Uint32ToBytes(n)

	prefixLen := len(userMetadataAuthUsedKeyPrefix)
	userTypeLen := len(userTypeBytes)
	userLen := len(userBytes)
	nLen := len(nBytes)

	key := make([]byte, prefixLen+userTypeLen+userLen+nLen)
	copy(key[:prefixLen], userMetadataAuthUsedKeyPrefix)
	copy(key[prefixLen:userTypeLen], userTypeBytes)
	copy(key[userTypeLen:userLen], userBytes)
	copy(key[userLen:nLen], nBytes)

	return key
}

func GetUserMetadataAuthByMetadataIdKey(userType apicommonpb.UserType, user, metadataId string) []byte {

	userTypeBytes := []byte(userType.String())
	userBytes := []byte(user)
	metadataIdBytes := []byte(metadataId)

	prefixLen := len(userMetadataAuthByMetadataIdKeyPrefix)
	userTypeLen := len(userTypeBytes)
	userLen := len(userBytes)
	metadataIdLen := len(metadataIdBytes)

	key := make([]byte, prefixLen+userTypeLen+userLen+metadataIdLen)
	copy(key[:prefixLen], userMetadataAuthByMetadataIdKeyPrefix)
	copy(key[prefixLen:userTypeLen], userTypeBytes)
	copy(key[userTypeLen:userLen], userBytes)
	copy(key[userLen:metadataIdLen], metadataIdBytes)

	return key
}

func GetMetadataUsedTaskIdCountKey(metadataId string) []byte {
	return append(metadataUsedTaskIdCountKey, []byte(metadataId)...)
}

func GetMetadataUsedTaskIdKey(metadataId string, n uint32) []byte {

	metadataIdBytes := []byte(metadataId)
	nBytes := bytesutil.Uint32ToBytes(n)

	prefixLen := len(metadataUsedTaskIdKeyPrefix)
	metadataIdLen := len(metadataIdBytes)
	nLen := len(nBytes)

	key := make([]byte, prefixLen+metadataIdLen+nLen)
	copy(key[:prefixLen], userMetadataAuthUsedKeyPrefix)
	copy(key[prefixLen:metadataIdLen], metadataIdBytes)
	copy(key[metadataIdLen:nLen], nBytes)

	return key
}

func GetTaskResultFileMetadataIdKey(taskId string) []byte {
	return append(taskResultFileMetadataIdKeyPrefix, []byte(taskId)...)
}
