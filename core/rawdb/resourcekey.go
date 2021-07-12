package rawdb

import "github.com/RosettaFlow/Carrier-Go/common"

var (
	nodeResourceKeyPrefix   = []byte("NodeResourceKey:")
	nodeResourceIdListKey   = []byte("nodeResourceIdListKey")
	orgResourceKeyPrefix    = []byte("OrgResourceKey:")
	orgResourceIdListKey    = []byte("OrgResourceIdListKey")
	nodeResourceSlotUnitKey = []byte("nodeResourceSlotKey")
	localTaskPowerUsedKeyPrefix = []byte("localTaskPowerUsedKey:")
	localTaskPowerUsedIdListKey = []byte("localTaskPowerUsedIdListKey")
	dataRereouceTableKeyPrefix = []byte("dataRereouceTableKey:")
	dataRereouceTableIdListKey = []byte("dataRereouceTableIdListKey")
	dataResourceDataUsedKeyPrefix = []byte("dataResourceDataUsedKey:")
	dataResourceDataUsedIdListKey = []byte("dataResourceDataUsedIdListKey")

	resourceTaskIdsKeyPrefix = []byte("resourceTaskIdsKeyPrefix:")

	resourcePowerIdMapingKeyPrefix = []byte("resourcePowerIdMapingKeyPrefix:")
)

// nodeResourceKey = NodeResourceKeyPrefix + jobNodeId
func GetNodeResourceKey(jobNodeId string) []byte {
	return append(nodeResourceKeyPrefix, common.Hex2Bytes(jobNodeId)...)
}
func GetNodeResourceIdListKey() []byte {
	return nodeResourceIdListKey
}
func GetOrgResourceKey(jobNodeId string) []byte {
	return append(orgResourceKeyPrefix, common.Hex2Bytes(jobNodeId)...)
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


func GetDataRereouceTableKey(nodeId string) []byte {
	return append(dataRereouceTableKeyPrefix, []byte(nodeId)...)
}
func GetDataRereouceTableIdListKey() []byte {
	return dataRereouceTableIdListKey
}


func GetDataResourceDataUsedKey(originId string) []byte {
	return append(dataResourceDataUsedKeyPrefix, []byte(originId)...)
}
func GetDataResourceDataUsedIdListKey() []byte {
	return dataResourceDataUsedIdListKey
}

func GetResourceTaskIdsKey(originId string) []byte {
	return append(resourceTaskIdsKeyPrefix, []byte(originId)...)
}

func GetResourcePowerIdMapingKey(powerId string) []byte {
	return append(resourcePowerIdMapingKeyPrefix, []byte(powerId)...)
}