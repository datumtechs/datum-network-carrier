package p2p

import (
	carriertwopcpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/consensus/twopc"
	carriernetmsgtaskmngpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/taskmng"
	carrierrpcdebugpbv1 "github.com/datumtechs/datum-network-carrier/pb/carrier/rpc/debug/v1"
	"github.com/gogo/protobuf/proto"
	"reflect"
)

// GossipTopicMappings represent the protocol ID to protobuf message type map for easy lookup.
var GossipTopicMappings = map[string]proto.Message{
	GossipTestDataTopicFormat:       &carrierrpcdebugpbv1.GossipTestData{},
	TwoPcPrepareMsgTopicFormat:      &carriertwopcpb.PrepareMsg{},
	TwoPcPrepareVoteTopicFormat:     &carriertwopcpb.PrepareVote{},
	TwoPcConfirmMsgTopicFormat:      &carriertwopcpb.ConfirmMsg{},
	TwoPcConfirmVoteTopicFormat:     &carriertwopcpb.ConfirmVote{},
	TwoPcCommitMsgTopicFormat:       &carriertwopcpb.CommitMsg{},
	TaskResultMsgTopicFormat:        &carriernetmsgtaskmngpb.TaskResultMsg{},
	TaskResourceUsageMsgTopicFormat: &carriernetmsgtaskmngpb.TaskResourceUsageMsg{},
	TaskTerminateMsgTopicFormat:     &carriernetmsgtaskmngpb.TaskTerminateMsg{},
}

// GossipTypeMapping is the inverse of GossipTopicMappings so that an arbitrary protobuf message
// can be mapped to a protocol ID string.
var GossipTypeMapping = make(map[reflect.Type]string, len(GossipTopicMappings))

func init() {
	for k, v := range GossipTopicMappings {
		GossipTypeMapping[reflect.TypeOf(v)] = k
	}
}
