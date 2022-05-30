package p2p

import (
	libp2ppb "github.com/datumtechs/datum-network-carrier/pb/carrier/rpc/debug/v1"
	"reflect"

	twopcpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/consensus/twopc"
	taskmngpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/taskmng"
	p2ppb "github.com/datumtechs/datum-network-carrier/pb/carrier/p2p/v1"
	"github.com/pkg/errors"
	types "github.com/prysmaticlabs/eth2-types"
)

// Current schema version for our rpc protocol ID.
const schemaVersionV1 = "/1"

const (
	// RPCStatusTopic defines the topic for the status rpc method.
	RPCStatusTopic = "/rosettanet/carrier_chain/req/status" + schemaVersionV1
	// RPCGoodByeTopic defines the topic for the goodbye rpc method.
	RPCGoodByeTopic = "/rosettanet/carrier_chain/req/goodbye" + schemaVersionV1
	// RPCPingTopic defines the topic for the ping rpc method.
	RPCPingTopic = "/rosettanet/carrier_chain/req/ping" + schemaVersionV1
	// RPCMetaDataTopic defines the topic for the metadata rpc method.
	RPCMetaDataTopic = "/rosettanet/carrier_chain/req/metadata" + schemaVersionV1

	// RPCBlocksByRangeTopic defines the topic for the blocks by range rpc method.
	RPCBlocksByRangeTopic = "/rosettanet/carrier_chain/req/carrier_blocks_by_range" + schemaVersionV1

	// for test communication.
	RPCGossipTestDataByRangeTopic = "/rosettanet/carrier_chain/req/gossip_test_data_by_range" + schemaVersionV1

	// for 2pc consensus
	RPCTwoPcPrepareMsgTopic  = "/rosettanet/consensus/twopc/send_prepare_msg" + schemaVersionV1
	RPCTwoPcPrepareVoteTopic = "/rosettanet/consensus/twopc/send_prepare_vote" + schemaVersionV1
	RPCTwoPcConfirmMsgTopic  = "/rosettanet/consensus/twopc/send_confirm_msg" + schemaVersionV1
	RPCTwoPcConfirmVoteTopic = "/rosettanet/consensus/twopc/send_confirm_vote" + schemaVersionV1
	RPCTwoPcCommitMsgTopic   = "/rosettanet/consensus/twopc/send_commit_msg" + schemaVersionV1

	// for task manager
	RPCTaskResultMsgTopic        = "/rosettanet/consensus/task/send_task_result_msg" + schemaVersionV1
	RPCTaskResourceUsageMsgTopic = "/rosettanet/consensus/task/send_task_resource_usage_msg" + schemaVersionV1
	RPCTaskTerminateMsgTopic     = "/rosettanet/consensus/task/send_task_terminate_msg" + schemaVersionV1
)

// RPCTopicMappings map the base message type to the rpc request.
var RPCTopicMappings = map[string]interface{}{
	RPCStatusTopic:                new(p2ppb.Status),
	RPCGoodByeTopic:               new(types.SSZUint64),
	RPCBlocksByRangeTopic:         new(p2ppb.CarrierBlocksByRangeRequest),
	RPCPingTopic:                  new(types.SSZUint64),
	RPCMetaDataTopic:              new(interface{}),
	RPCGossipTestDataByRangeTopic: new(libp2ppb.GossipTestData),
	RPCTwoPcPrepareMsgTopic:       new(twopcpb.PrepareMsg),
	RPCTwoPcPrepareVoteTopic:      new(twopcpb.PrepareVote),
	RPCTwoPcConfirmMsgTopic:       new(twopcpb.ConfirmMsg),
	RPCTwoPcConfirmVoteTopic:      new(twopcpb.ConfirmVote),
	RPCTwoPcCommitMsgTopic:        new(twopcpb.CommitMsg),
	RPCTaskResultMsgTopic:         new(taskmngpb.TaskResultMsg),
	RPCTaskResourceUsageMsgTopic:  new(taskmngpb.TaskResourceUsageMsg),
	RPCTaskTerminateMsgTopic:      new(taskmngpb.TaskTerminateMsg),
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
