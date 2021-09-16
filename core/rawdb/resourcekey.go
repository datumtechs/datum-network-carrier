package rawdb

import (
	"github.com/RosettaFlow/Carrier-Go/common/bytesutil"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
)

var (
	// prefix + jobNodeId -> LocalResourceTable
	nodeResourceKeyPrefix           = []byte("NodeResourceKeyPrefix:")
	// key -> [jobNodeId, jobNodeId, ..., jobNodeId]
	nodeResourceIdListKey           = []byte("nodeResourceIdListKey")
	// prefix + identityId -> RemoteResourceTable
	orgResourceKeyPrefix            = []byte("OrgResourceKeyPrefix:")
	// key -> [identityId, identityId, ..., identityId]
	orgResourceIdListKey            = []byte("OrgResourceIdListKey")
	// key -> SlotUnit
	nodeResourceSlotUnitKey         = []byte("nodeResourceSlotUnitKey")
	// prefix + taskId -> LocalTaskPowerUsed
	localTaskPowerUsedKeyPrefix     = []byte("localTaskPowerUsedKeyPrefix:")
	// key -> [taskId, taskId, ..., taskId]
	localTaskPowerUsedIdListKey     = []byte("localTaskPowerUsedIdListKey")
	// prefix + dataNodeId -> DataResourceTable{dataNodeId, totalDisk, usedDisk}
	dataResourceTableKeyPrefix      = []byte("dataResourceTableKeyPrefix:")
	// key -> [dataNodeId, dataNodeId, ..., dataNodeId]
	dataResourceTableIdListKey      = []byte("dataResourceTableIdListKey")
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
	userMetadataAuthUsedTotalKey = []byte("userMetadataAuthUsedTotalKey")
	// prefix + userType + user + n -> metadataId
	userMetadataAuthUsedKeyPrefix = []byte("userMetadataAuthUsedKeyPrefix:")


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

func GetUserMetadataAuthUsedTotalKey(userType apicommonpb.UserType, user string) []byte {
	return append(append(userMetadataAuthUsedTotalKey, []byte(userType.String())...), []byte(user)...)
}

func GetUserMetadataAuthUsedKey(userType apicommonpb.UserType, user string, n uint32) []byte {

	userTypeBytes := []byte(userType.String())
	userBytes := []byte(user)
	nBytes := bytesutil.Uint32ToBytes(n)

	prefixLen := len(localTaskExecuteStatusKeyPrefix)
	userTypeLen := len(userTypeBytes)
	userLen := len(userBytes)
	nLen := len(nBytes)

	key := make([]byte, prefixLen+userTypeLen+userLen+nLen)
	copy(key[:prefixLen], localTaskExecuteStatusKeyPrefix)
	copy(key[prefixLen:userTypeLen], userTypeBytes)
	copy(key[userTypeLen:userLen], userBytes)
	copy(key[userLen:nLen], nBytes)

	return key
}
