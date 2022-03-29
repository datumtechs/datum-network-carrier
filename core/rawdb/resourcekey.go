package rawdb

import (
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)

var (
	// prefix + jobNodeId -> LocalResourceTable
	nodeResourceKeyPrefix = []byte("NodeResourceKeyPrefix:")

	// key -> SlotUnit
	//nodeResourceSlotUnitKey = []byte("nodeResourceSlotUnitKey")

	// prefix + taskId + partyId -> LocalTaskPowerUsed
	localTaskPowerUsedKeyPrefix = []byte("localTaskPowerUsedKeyPrefix:")
	// prefix + jobNodeId + taskId -> [partyId, ..., partyId]
	jobNodeTaskPartyIdsKeyPrefix = []byte("jobNodeTaskPartyIdsKeyPrefix:")

	// prefix + jobNodeId -> history task count
	jobNodeHistoryTaskCountKeyPrefix = []byte("jobNodeHistoryTaskCountKeyPrefix:")
	// prefix + jobNodeId + taskId -> index
	jobNodeHistoryTaskKeyPrefix = []byte("jobNodeHistoryTaskKeyPrefix:")

	// prefix + dataNodeId -> DataResourceTable{dataNodeId, totalDisk, usedDisk}
	dataResourceTableKeyPrefix = []byte("dataResourceTableKeyPrefix:")

	// prefix + originId -> DataResourceFileUpload{originId, dataNodeId, metaDataId, filePath}
	dataResourceFileUploadKeyPrefix = []byte("dataResourceDataUsedKeyPrefix:")

	// prefix + powerId -> jobNodeId
	powerIdJobNodeIdMapingKeyPrefix = []byte("powerIdJobNodeIdMapingKeyPrefix:")
	// prefix + metaDataId -> DataResourceDiskUsed{metaDataId, dataNodeId, diskUsed}
	dataResourceDiskUsedKeyPrefix = []byte("DataResourceDiskUsedKeyPrefix:")
	// prefix + taskId + partyId -> executeStatus (uint64)
	localTaskExecuteStatusKeyPrefix = []byte("localTaskExecuteStatusKeyPrefix:")

	// prefix + userType + user + metadataId -> metadataAuthId (only one)
	userMetadataAuthByMetadataIdKeyPrefix = []byte("userMetadataAuthByMetadataIdKeyPrefix:")

	// prefix + metadataId -> history task count
	metadataHistoryTaskCountKeyPrefix = []byte("metadataHistoryTaskCountKeyPrefix:")
	// prefix + metadataId + taskId -> index
	metadataHistoryTaskKeyPrefix = []byte("metadataHistoryTaskKeyPrefix:")

	// prefix + taskId -> resultfile summary (auto build metadataId)
	taskResultFileMetadataIdKeyPrefix = []byte("taskResultFileMetadataIdKeyPrefix:")
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


const (
	/**
	######   ######   ######   ######   ######
	#   THE LOCAL NEEDEXECUTE TASK STATUS    #
	######   ######   ######   ######   ######
	*/
	OnConsensusExecuteTaskStatus   LocalTaskExecuteStatus = 1 << iota 	// 0001: the execute task is on consensus period now.
	OnRunningExecuteStatus                               				// 0010: the execute task is running now.
	OnTerminingExecuteStatus                               				// 0010: the execute task is termining now.
	UnKnownExecuteTaskStatus       = 0        							// 0000: the execute task status is unknown.
)

type LocalTaskExecuteStatus uint32

func (s LocalTaskExecuteStatus) Uint32() uint32 { return uint32(s) }

/// -- keys ... --

//
func GetNodeResourceKeyPrefix() []byte {
	return nodeResourceKeyPrefix
}

// nodeResourceKey = NodeResourceKeyPrefix + jobNodeId
func GetNodeResourceKey(jobNodeId string) []byte {
	return append(nodeResourceKeyPrefix, []byte(jobNodeId)...)
}

//func GetNodeResourceSlotUnitKey() []byte {
//	return nodeResourceSlotUnitKey
//}

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
// prefix + jobNodeId + taskId -> index
func GetJobNodeHistoryTaskKey (jobNodeId, taskId string) []byte {
	return append(append(jobNodeHistoryTaskKeyPrefix, []byte(jobNodeId)...), []byte(taskId)...)
}

// prefix + dataNodeId -> DataResourceTable{dataNodeId, totalDisk, usedDisk}
func GetDataResourceTableKeyPrefix() []byte {
	return dataResourceTableKeyPrefix
}
// prefix + dataNodeId -> DataResourceTable{dataNodeId, totalDisk, usedDisk}
func GetDataResourceTableKey(dataNodeId string) []byte {
	return append(dataResourceTableKeyPrefix, []byte(dataNodeId)...)
}

// prefix + originId -> DataResourceFileUpload{originId, dataNodeId, metaDataId, filePath}
func GetDataResourceFileUploadKeyPrefix() []byte {
	return dataResourceFileUploadKeyPrefix
}
// prefix + originId -> DataResourceFileUpload{originId, dataNodeId, metaDataId, filePath}
func GetDataResourceFileUploadKey(originId string) []byte {
	return append(dataResourceFileUploadKeyPrefix, []byte(originId)...)
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

func GetUserMetadataAuthByMetadataIdKey(userType libtypes.UserType, user, metadataId string) []byte {

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

// prefix + metadataId -> history task count
func GetMetadataHistoryTaskCountKey (metadataId string) []byte {
	return append(metadataHistoryTaskCountKeyPrefix, []byte(metadataId)...)
}
// prefix + metadataId + taskId -> index
func GetMetadataHistoryTaskKeyPrefixByMetadataId (metadataId string) []byte {
	return append(metadataHistoryTaskKeyPrefix, []byte(metadataId)...)
}
// prefix + metadataId + taskId -> index
func GetMetadataHistoryTaskKey (metadataId, taskId string) []byte {
	return append(append(metadataHistoryTaskKeyPrefix, []byte(metadataId)...), []byte(taskId)...)
}

func GetTaskResultFileMetadataIdKey(taskId string) []byte {
	return append(taskResultFileMetadataIdKeyPrefix, []byte(taskId)...)
}

func GetTaskResultFileMetadataIdKeyPrefix() []byte {
	return taskResultFileMetadataIdKeyPrefix
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