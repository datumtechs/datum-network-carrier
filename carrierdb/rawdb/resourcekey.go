package rawdb

import (
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
)

var (

	// --------------- for resource ---------------
	//
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
	// prefix + originId -> DataResourceDataUpload{originId, dataNodeId, metaDataId, filePath}
	dataResourceDataUploadKeyPrefix = []byte("dataResourceDataUsedKeyPrefix:")
	// prefix + powerId -> jobNodeId
	powerIdJobNodeIdMapingKeyPrefix = []byte("powerIdJobNodeIdMapingKeyPrefix:")
	// prefix + metaDataId -> DataResourceDiskUsed{metaDataId, dataNodeId, diskUsed}
	dataResourceDiskUsedKeyPrefix = []byte("DataResourceDiskUsedKeyPrefix:")

	// --------------- for metadataAuth ---------------
	//
	// prefix + userType + user + metadataId -> metadataAuth totalCount
	userMetadataAuthTotalCountByMetadataIdPrefix = []byte("userMetadataAuthTotalCountByMetadataIdPrefix:")
	// prefix + userType + user + metadataId -> metadataAuth validCount
	userMetadataAuthValidCountByMetadataIdPrefix = []byte("userMetadataAuthValidCountByMetadataIdPrefix:")
	// prefix + uesrType + user + metadataId + metadataAuthId -> status
	// 数据授权信息的状态 (0: 未知; 1: 还未发布的数据授权; 2: 已发布的数据授权; 3: 已撤销的数据授权 <失效前主动撤回的>; 4: 已经失效的数据授权 <过期or达到使用上限的or被拒绝的>;)
	//
	// Use a value of uint16 to divide it into high-order uint8 and low-order uint8,
	// where the low-order uint8 represents the state of meadataAuth
	// and the high-order uint8 represents the audit option of metadataAauth,
	// the low-order uint8 => 0: unknown, 1: created, 2: released, 3: revoked, 4: invalid,
	// the high-order uint8 => 0: pending, 1: passed, 2: refused.
	userMetadataAuthStatusByMetadataIdPrefix = []byte("userMetadataAuthStatusByMetadataIdPrefix")

	// --------------- for task ---------------
	//
	// prefix + taskId + partyId -> executeStatus (uint64)
	localTaskExecuteStatusKeyPrefix = []byte("localTaskExecuteStatusKeyPrefix:")
	// prefix + metadataId -> history task count
	metadataHistoryTaskCountKeyPrefix = []byte("metadataHistoryTaskCountKeyPrefix:")
	// prefix + metadataId + taskId -> index
	metadataHistoryTaskKeyPrefix = []byte("metadataHistoryTaskKeyPrefix:")
	// prefix + taskId -> resultData summary (auto build metadataId)
	taskResultDataMetadataIdKeyPrefix = []byte("taskResultDataMetadataIdKeyPrefix:")
	// prefix + taskId -> [partyId, ..., partyId]  for task sender
	taskPartnerPartyIdsKeyPrefix = []byte("taskPartnerPartyIdsKeyPrefix:")
	// prefix + taskId -> needExecuteTask
	needExecuteTaskKeyPrefix = []byte("needExecuteTaskKeyPrefix:")

	// ---------------  for message_handler  ---------------
	powerMsgKeyPrefix        = []byte("powerMsgKeyPrefix:")
	metadataMsgKeyPrefix     = []byte("metadataMsgKeyPrefix:")
	metadataUpdateMsgPrefix  = []byte("metadataUpdateMsgPrefix:")
	metadataAuthMsgKeyPrefix = []byte("metadataAuthMsgKeyPrefix:")
	taskMsgKeyPrefix         = []byte("taskMsgKeyPrefix:")

	// --------------- for scheduler (task bullet) ---------------
	taskBulletKeyPrefix = []byte("taskBulletKeyPrefix:")

	// --------------- for organization built-in private key ---------------
	orgPriKeyPrefix = []byte("orgPriKeyPrefix:")

	// --------------- for msg nonce ---------------
	// identityMsgNonceKey -> nonce
	identityMsgNonceKey = []byte("identityMsgNonceKey")
	// metadataMsgNonceKey -> nonce
	metadataMsgNonceKey = []byte("metadataMsgNonceKey")
	// metadataAuthMsgNonceKey -> nonce
	metadataAuthMsgNonceKey = []byte("metadataAuthMsgNonceKey")
	// powerMsgNonceKey -> nonce
	powerMsgNonceKey = []byte("powerMsgNonceKey")
	// taskMsgNonceKey -> nonce
	taskMsgNonceKey = []byte("taskMsgNonceKey")
)

const (
	/**
	######   ######   ######   ######   ######
	#   THE LOCAL NEEDEXECUTE TASK STATUS    #
	######   ######   ######   ######   ######
	*/
	OnConsensusExecuteTaskStatus LocalTaskExecuteStatus = 1 << iota // 0001: the execute task is on consensus period now.
	OnRunningExecuteStatus                                          // 0010: the execute task is running now.
	OnTerminingExecuteStatus                                        // 0010: the execute task is termining now.
	UnKnownExecuteTaskStatus     = 0                                // 0000: the execute task status is unknown.
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
func GetJobNodeTaskPartyIdsKey(jobNodeId, taskId string) []byte {
	return append(append(jobNodeTaskPartyIdsKeyPrefix, []byte(jobNodeId)...), []byte(taskId)...)
}

// prefix + jobNodeId + taskId -> [partyId, ..., partyId]
func GetJobNodeTaskPartyIdsKeyPrefixByJobNodeId(jobNodeId string) []byte {
	return append(jobNodeTaskPartyIdsKeyPrefix, []byte(jobNodeId)...)
}

// prefix + jobNodeId -> history task count
func GetJobNodeHistoryTaskCountKey(jobNodeId string) []byte {
	return append(jobNodeHistoryTaskCountKeyPrefix, []byte(jobNodeId)...)
}

// prefix + jobNodeId + taskId -> index
func GetJobNodeHistoryTaskKey(jobNodeId, taskId string) []byte {
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

// prefix + originId -> DataResourceDataUpload{originId, dataNodeId, metaDataId, filePath}
func GetDataResourceDataUploadKeyPrefix() []byte {
	return dataResourceDataUploadKeyPrefix
}

// prefix + originId -> DataResourceDataUpload{originId, dataNodeId, metaDataId, filePath}
func GetDataResourceDataUploadKey(originId string) []byte {
	return append(dataResourceDataUploadKeyPrefix, []byte(originId)...)
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

// prefix + userType + user + metadataId ->  metadataAuth totalCount
func GetTotalUserMetadataAuthCountByMetadataIdKey(userType commonconstantpb.UserType, user, metadataId string) []byte {

	// key: prefix + userType + user + metadataId ->  value: metadataAuth totalCount

	userTypeBytes := []byte(userType.String())
	userBytes := []byte(user)
	metadataIdBytes := []byte(metadataId)

	// some index of pivots
	prefixIndex := len(userMetadataAuthTotalCountByMetadataIdPrefix)
	userTypeIndex := prefixIndex + len(userTypeBytes)
	userIndex := userTypeIndex + len(userBytes)
	size := userIndex + len(metadataIdBytes)

	// construct key
	key := make([]byte, size)
	copy(key[:prefixIndex], userMetadataAuthTotalCountByMetadataIdPrefix)
	copy(key[prefixIndex:userTypeIndex], userTypeBytes)
	copy(key[userTypeIndex:userIndex], userBytes)
	copy(key[userIndex:], metadataIdBytes)

	return key
}

// prefix + userType + user + metadataId -> metadataAuth validCount
func GetValidUserMetadataAuthCountByMetadataIdKey(userType commonconstantpb.UserType, user, metadataId string) []byte {

	// key: prefix + userType + user + metadataId ->  value: metadataAuth validCount

	userTypeBytes := []byte(userType.String())
	userBytes := []byte(user)
	metadataIdBytes := []byte(metadataId)

	// some index of pivots
	prefixIndex := len(userMetadataAuthValidCountByMetadataIdPrefix)
	userTypeIndex := prefixIndex + len(userTypeBytes)
	userIndex := userTypeIndex + len(userBytes)
	size := userIndex + len(metadataIdBytes)

	// construct key
	key := make([]byte, size)
	copy(key[:prefixIndex], userMetadataAuthValidCountByMetadataIdPrefix)
	copy(key[prefixIndex:userTypeIndex], userTypeBytes)
	copy(key[userTypeIndex:userIndex], userBytes)
	copy(key[userIndex:], metadataIdBytes)

	return key
}

// prefix + uesrType + user + metadataId + metadataAuthId -> status <0: invalid, 1: valid>
func GetUserMetadataAuthStatusByMetadataIdKey(userType commonconstantpb.UserType, user, metadataId, metadataAuthId string) []byte {

	// key: prefix + uesrType + user + metadataId + metadataAuthId -> value status <0: invalid, 1: valid>

	userTypeBytes := []byte(userType.String())
	userBytes := []byte(user)
	metadataIdBytes := []byte(metadataId)
	metadataAuthIdBytes := []byte(metadataAuthId)

	// some index of pivots
	prefixIndex := len(userMetadataAuthStatusByMetadataIdPrefix)
	userTypeIndex := prefixIndex + len(userTypeBytes)
	userIndex := userTypeIndex + len(userBytes)
	metadataIdIndex := userIndex + len(metadataIdBytes)
	size := metadataIdIndex + len(metadataAuthIdBytes)

	// construct key
	key := make([]byte, size)
	copy(key[:prefixIndex], userMetadataAuthStatusByMetadataIdPrefix)
	copy(key[prefixIndex:userTypeIndex], userTypeBytes)
	copy(key[userTypeIndex:userIndex], userBytes)
	copy(key[userIndex:metadataIdIndex], metadataIdBytes)
	copy(key[metadataIdIndex:], metadataAuthIdBytes)

	return key
}

// prefix + uesrType + user + metadataId + metadataAuthId -> status <0: invalid, 1: valid>
func GetUserMetadataAuthStatusKeyPrefixByMetadataId(userType commonconstantpb.UserType, user, metadataId string) []byte {
	// key: prefix + uesrType + user + metadataId + metadataAuthId -> value status <0: invalid, 1: valid>

	userTypeBytes := []byte(userType.String())
	userBytes := []byte(user)
	metadataIdBytes := []byte(metadataId)

	// some index of pivots
	prefixIndex := len(userMetadataAuthStatusByMetadataIdPrefix)
	userTypeIndex := prefixIndex + len(userTypeBytes)
	userIndex := userTypeIndex + len(userBytes)
	size := userIndex + len(metadataIdBytes)

	// construct key
	key := make([]byte, size)
	copy(key[:prefixIndex], userMetadataAuthStatusByMetadataIdPrefix)
	copy(key[prefixIndex:userTypeIndex], userTypeBytes)
	copy(key[userTypeIndex:userIndex], userBytes)
	copy(key[userIndex:], metadataIdBytes)

	return key
}

// prefix + metadataId -> history task count
func GetMetadataHistoryTaskCountKey(metadataId string) []byte {
	return append(metadataHistoryTaskCountKeyPrefix, []byte(metadataId)...)
}

// prefix + metadataId + taskId -> index
func GetMetadataHistoryTaskKeyPrefixByMetadataId(metadataId string) []byte {
	return append(metadataHistoryTaskKeyPrefix, []byte(metadataId)...)
}

// prefix + metadataId + taskId -> index
func GetMetadataHistoryTaskKey(metadataId, taskId string) []byte {
	return append(append(metadataHistoryTaskKeyPrefix, []byte(metadataId)...), []byte(taskId)...)
}

func GetTaskResultDataMetadataIdKey(taskId string) []byte {
	return append(taskResultDataMetadataIdKeyPrefix, []byte(taskId)...)
}

func GetTaskResultDataMetadataIdKeyPrefix() []byte {
	return taskResultDataMetadataIdKeyPrefix
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

func GetMetadataUpdateMsgKeyPrefix() []byte {
	return metadataUpdateMsgPrefix
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

func GetMetadataUpdateMsgKey(metadataId string) []byte {
	return append(metadataUpdateMsgPrefix, []byte(metadataId)...)
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

func GetOrgPriKeyPrefix() []byte {
	return orgPriKeyPrefix
}

func GetIdentityMsgNonceKey() []byte {
	return identityMsgNonceKey
}

func GetMetadataMsgNonceKey() []byte {
	return metadataMsgNonceKey
}

func GetMetadataAuthMsgNonceKey() []byte {
	return metadataAuthMsgNonceKey
}

func GetPowerMsgNonceKey() []byte {
	return powerMsgNonceKey
}

func GetTaskMsgNonceKey() []byte {
	return taskMsgNonceKey
}
