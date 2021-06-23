package p2p

import (
	"reflect"

	pb "github.com/RosettaFlow/Carrier-Go/lib/p2p/v1"
	"github.com/pkg/errors"
	types "github.com/prysmaticlabs/eth2-types"
)

// Current schema version for our rpc protocol ID.
const schemaVersionV1 = "/1"

const (
	// RPCStatusTopic defines the topic for the status rpc method.
	RPCStatusTopic = "/rosetta/carrier_chain/req/status" + schemaVersionV1
	// RPCGoodByeTopic defines the topic for the goodbye rpc method.
	RPCGoodByeTopic = "/rosetta/carrier_chain/req/goodbye" + schemaVersionV1
	// RPCBlocksByRangeTopic defines the topic for the blocks by range rpc method.
	RPCBlocksByRangeTopic = "/rosetta/carrier_chain/req/carrier_blocks_by_range" + schemaVersionV1
	// RPCBlocksByRootTopic defines the topic for the blocks by root rpc method.
	RPCBlocksByRootTopic = "/rosetta/carrier_chain/req/carrier_blocks_by_root" + schemaVersionV1
	// RPCPingTopic defines the topic for the ping rpc method.
	RPCPingTopic = "/rosetta/carrier_chain/req/ping" + schemaVersionV1
	// RPCMetaDataTopic defines the topic for the metadata rpc method.
	RPCMetaDataTopic = "/rosetta/carrier_chain/req/metadata" + schemaVersionV1
)

// RPCTopicMappings map the base message type to the rpc request.
var RPCTopicMappings = map[string]interface{}{
	RPCStatusTopic:        new(pb.Status),
	RPCGoodByeTopic:       new(types.SSZUint64),
	RPCBlocksByRangeTopic: new(pb.BeaconBlocksByRangeRequest),
	//RPCBlocksByRootTopic:  new(p2ptypes.BeaconBlockByRootsReq),
	RPCPingTopic:          new(types.SSZUint64),
	RPCMetaDataTopic:      new(interface{}),
}

// VerifyTopicMapping verifies that the topic and its accompanying
// message type is correct.
func VerifyTopicMapping(topic string, msg interface{}) error {
	msgType, ok := RPCTopicMappings[topic]
	if !ok {
		return errors.New("rpc topic is not registered currently")
	}
	receivedType := reflect.TypeOf(msg)
	registeredType := reflect.TypeOf(msgType)
	typeMatches := registeredType.AssignableTo(receivedType)

	if !typeMatches {
		return errors.Errorf("accompanying message type is incorrect for topic: wanted %v  but got %v",
			registeredType.String(), receivedType.String())
	}
	return nil
}
