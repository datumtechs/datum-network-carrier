package resource

import "github.com/RosettaFlow/Carrier-Go/common"

var (
	nodeResourceKeyPrefix   = []byte("NodeResourceKey:")
	nodeResourceIdListKey   = []byte("nodeResourceIdListKey:")
	orgResourceKeyPrefix    = []byte("OrgResourceKey:")
	orgResourceIdListKey    = []byte("OrgResourceIdListKey:")
	nodeResourceSlotUnitKey = []byte("nodeResourceSlotKey")
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
