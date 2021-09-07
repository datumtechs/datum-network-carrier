package rawdb


var (
	// jobNodeId -> LocalResourceTable
	nodeResourceKeyPrefix           = []byte("NodeResourceKey:")
	// key -> [jobNodeId, jobNodeId, ..., jobNodeId]
	nodeResourceIdListKey           = []byte("nodeResourceIdListKey")
	// identityId -> RemoteResourceTable
	orgResourceKeyPrefix            = []byte("OrgResourceKey:")
	// key -> [identityId, identityId, ..., identityId]
	orgResourceIdListKey            = []byte("OrgResourceIdListKey")
	// key -> SlotUnit
	nodeResourceSlotUnitKey         = []byte("nodeResourceSlotUnitKey")
	// taskId -> LocalTaskPowerUsed
	localTaskPowerUsedKeyPrefix     = []byte("localTaskPowerUsedKey:")
	// key -> [taskId, taskId, ..., taskId]
	localTaskPowerUsedIdListKey     = []byte("localTaskPowerUsedIdListKey")
	// dataNodeId -> DataResourceTable{dataNodeId, totalDisk, usedDisk}
	dataResourceTableKeyPrefix      = []byte("dataResourceTableKey:")
	// key -> [dataNodeId, dataNodeId, ..., dataNodeId]
	dataResourceTableIdListKey      = []byte("dataResourceTableIdListKey")
	// originId -> DataResourceFileUpload{originId, dataNodeId, metaDataId, filePath}
	dataResourceFileUploadKeyPrefix = []byte("dataResourceDataUsedKey:")
	// key -> [originId, originId, ..., originId]
	dataResourceFileUploadIdListKey = []byte("dataResourceFileUploadIdListKey")
	// jonNodeId -> [taskId, taskId, ..., taskId]
	resourceTaskIdsKeyPrefix = []byte("resourceTaskIdsKeyPrefix:")
	// powerId -> jobNodeId
	resourcePowerIdMapingKeyPrefix = []byte("resourcePowerIdMapingKeyPrefix:")
	// metaDataId -> DataResourceDiskUsed{metaDataId, dataNodeId, diskUsed}
	dataResourceDiskUsedKeyPrefix = []byte("DataResourceDiskUsedKeyPrefix:")
	// taskId -> executeStatus
	localTaskExecuteStatusPrefix = []byte("localTaskExecuteStatus")
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
	return append(localTaskExecuteStatusPrefix, []byte(taskId)...)
}
