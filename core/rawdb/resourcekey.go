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

	// key -> SlotUnit
	nodeResourceSlotUnitKey = []byte("nodeResourceSlotUnitKey")

	// prefix + taskId + partyId -> LocalTaskPowerUsed
	localTaskPowerUsedKeyPrefix = []byte("localTaskPowerUsedKeyPrefix:")
	// prefix + jobNodeId + taskId -> [partyId, ..., partyId]
	jobNodeTaskPartyIdsKeyPrefix = []byte("jobNodeTaskPartyIdsKeyPrefix:")
	// prefix + jobNodeId -> history task count
	jobNodeHistoryTaskCountKeyPrefix = []byte("jobNodeHistoryTaskCountKeyPrefix")




	// prefix + dataNodeId -> DataResourceTable{dataNodeId, totalDisk, usedDisk}
	dataResourceTableKeyPrefix = []byte("dataResourceTableKeyPrefix:")
	// key -> [dataNodeId, dataNodeId, ..., dataNodeId]
	dataResourceTableIdListKey = []byte("dataResourceTableIdListKey")

	// prefix + originId -> DataResourceFileUpload{originId, dataNodeId, metaDataId, filePath}
	dataResourceFileUploadKeyPrefix = []byte("dataResourceDataUsedKeyPrefix:")
	// key -> [originId, originId, ..., originId]
	dataResourceFileUploadIdListKey = []byte("dataResourceFileUploadIdListKey")

	// prefix + powerId -> jobNodeId
	powerIdJobNodeIdMapingKeyPrefix = []byte("powerIdJobNodeIdMapingKeyPrefix:")
	// prefix + metaDataId -> DataResourceDiskUsed{metaDataId, dataNodeId, diskUsed}
	dataResourceDiskUsedKeyPrefix = []byte("DataResourceDiskUsedKeyPrefix:")
	// prefix + taskId + partyId -> executeStatus
	localTaskExecuteStatusKeyPrefix = []byte("localTaskExecuteStatusKeyPrefix:")
	localTaskExecuteStatusValCons   = []byte("cons") // start consensus
	localTaskExecuteStatusValExec   = []byte("exec") // start execute

	// prefix + userType + user -> n
	userMetadataAuthUsedCountKey = []byte("userMetadataAuthUsedCountKey")
	// prefix + userType + user + n -> metadataAuthId
	userMetadataAuthUsedKeyPrefix = []byte("userMetadataAuthUsedKeyPrefix:")
	// prefix + userType + user + metadataId -> metadataAuthId (only one)
	userMetadataAuthByMetadataIdKeyPrefix = []byte("userMetadataAuthByMetadataIdKeyPrefix:")

	// metadataId -> taskCount (n)
	metadataUsedTaskIdCountKey = []byte("metadataUsedTaskIdCountKey")
	// metadataId + n -> taskId
	metadataUsedTaskIdKeyPrefix = []byte("metadataUsedTaskIdKeyPrefix:")

	// prefix + taskId -> n todo maybe support multi-result file about a taskId
	// prefix + taskId + n -> resultfile summary (auto build metadataId)

	// prefix + taskId -> resultfile summary (auto build metadataId)
	taskResultFileMetadataIdKeyPrefix = []byte("taskResultFileMetadataIdKeyPrefix:")
	//// prefix + taskId + partyId -> resourceUsed (totalProccesor, usedProccesor, totalMemory, usedMemory, totalBandwidth, usedBandwidth, totalDisk, usedDisk)
	//taskResuorceUsageKeyPrefix = []byte("taskResuorceUsageKeyPrefix:")
	// prefix + taskId -> powerPartyIds
	taskPowerPartyIdsKeyPrefix = []byte("taskPowerPartyIdsKeyPrefix:")
	// prefix + taskId -> [partyId, ..., partyId]  for task sender
	taskPartnerPartyIdsKeyPrefix = []byte("taskPartnerPartyIdsKeyPrefix:")

	needExecuteTaskKeyPrefix = []byte("needExecuteTaskKeyPrefix:")

	// ---------  for message_handler  ---------
	powerMsgKeyPrefix        = []byte("powerMsgKeyPrefix:")
	metadataMsgKeyPrefix     = []byte("metadataMsgKeyPrefix:")
	metadataAuthMsgKeyPrefix = []byte("metadataAuthMsgKeyPrefix:")
	taskMsgKeyPrefix         = []byte("taskMsgKeyPrefix:")

	// ---------- for scheduler (task bullet) ----------
	taskBulletKeyPrefix = []byte("taskBulletKeyPrefix:")
)

// nodeResourceKey = NodeResourceKeyPrefix + jobNodeId
func GetNodeResourceKey(jobNodeId string) []byte {
	return append(nodeResourceKeyPrefix, []byte(jobNodeId)...)
}
func GetNodeResourceIdListKey() []byte {
	return nodeResourceIdListKey
}

func GetNodeResourceSlotUnitKey() []byte {
	return nodeResourceSlotUnitKey
}

func GetLocalTaskPowerUsedKey(taskId, partyId string) []byte {
	return append(append(localTaskPowerUsedKeyPrefix, []byte(taskId)...), []byte(partyId)...)
}

func GetLocalTaskPowerUsedKeyPrefix() []byte {
	return localTaskPowerUsedKeyPrefix
}

func GetLocalTaskPowerUsedKeyPrefixByTaskId(taskId string) []byte {
	return append(localTaskPowerUsedKeyPrefix, []byte(taskId)...)
}

// prefix + jobNodeId + taskId -> [partyId, ..., partyId]
func GetJobNodeTaskPartyIdsKey (jobNodeId, taskId string) []byte {
	return append(append(jobNodeTaskPartyIdsKeyPrefix, []byte(jobNodeId)...), []byte(taskId)...)
}

// prefix + jobNodeId + taskId -> [partyId, ..., partyId]
func GetJobNodeTaskPartyIdsKeyPrefixByJobNodeId (jobNodeId string) []byte {
	return append(jobNodeTaskPartyIdsKeyPrefix, []byte(jobNodeId)...)
}

// prefix + jobNodeId -> history task count
func GetJobNodeHistoryTaskCountKey (jobNodeId string) []byte {
	return append(jobNodeHistoryTaskCountKeyPrefix, []byte(jobNodeId)...)
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
// prefix + powerId -> jobNodeId
func GetPowerIdJobNodeIdMapingKey(powerId string) []byte {
	return append(powerIdJobNodeIdMapingKeyPrefix, []byte(powerId)...)
}

func GetDataResourceDiskUsedKey(metaDataId string) []byte {
	return append(dataResourceDiskUsedKeyPrefix, []byte(metaDataId)...)
}

func GetLocalTaskExecuteStatus(taskId, partyId string) []byte {
	return append(append(localTaskExecuteStatusKeyPrefix, []byte(taskId)...), []byte(partyId)...)
}
func GetLocalTaskExecuteStatusValCons() []byte {
	return localTaskExecuteStatusValCons
}
func GetLocalTaskExecuteStatusValExec() []byte {
	return localTaskExecuteStatusValExec
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

func GetTaskResultFileMetadataIdKeyPrefix() []byte {
	return taskResultFileMetadataIdKeyPrefix
}

func GetTaskPowerPartyIdsKey(taskId string) []byte {
	return append(taskPowerPartyIdsKeyPrefix, []byte(taskId)...)
}

func GetTaskPartnerPartyIdsKey(taskId string) []byte {
	return append(taskPartnerPartyIdsKeyPrefix, []byte(taskId)...)
}

func GetNeedExecuteTaskKeyPrefix() []byte {
	return needExecuteTaskKeyPrefix
}

func GetNeedExecuteTaskKey(taskId, partyId string) []byte {
	return append(append(needExecuteTaskKeyPrefix, []byte(taskId)...), []byte(partyId)...)
}

func GetPowerMsgKeyPrefix() []byte {
	return powerMsgKeyPrefix
}

func GetMetadataMsgKeyPrefix() []byte {
	return metadataMsgKeyPrefix
}

func GetMetadataAuthMsgKeyPrefix() []byte {
	return metadataAuthMsgKeyPrefix
}

func GetTaskMsgKeyPrefix() []byte {
	return taskMsgKeyPrefix
}

func GetPowerMsgKey(powerId string) []byte {
	return append(powerMsgKeyPrefix, []byte(powerId)...)
}

func GetMetadataMsgKey(metadataId string) []byte {
	return append(metadataMsgKeyPrefix, []byte(metadataId)...)
}

func GetMetadataAuthMsgKey(metadataAuthId string) []byte {
	return append(metadataAuthMsgKeyPrefix, []byte(metadataAuthId)...)
}

func GetTaskMsgKey(taskId string) []byte {
	return append(taskMsgKeyPrefix, []byte(taskId)...)
}


func GetTaskBulletKeyPrefix() []byte {
	return taskBulletKeyPrefix
}

// GetTaskBulletKey = taskBulletKeyPrefix + taskId
func GetTaskBulletKey(taskId string) []byte {
	return append(taskBulletKeyPrefix, []byte(taskId)...)
}