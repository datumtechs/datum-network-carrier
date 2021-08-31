package p2p

import (
	twopb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	librpcpb "github.com/RosettaFlow/Carrier-Go/lib/rpc/v1"
	"github.com/gogo/protobuf/proto"
	"reflect"
)

// GossipTopicMappings represent the protocol ID to protobuf message type map for easy lookup.
var GossipTopicMappings = map[string]proto.Message{
	GossipTestDataTopicFormat:     &librpcpb.SignedGossipTestData{},
	TwoPcPrepareMsgTopicFormat:    &twopb.PrepareMsg{},
	TwoPcPrepareVoteTopicFormat:   &twopb.PrepareVote{},
	TwoPcConfirmMsgTopicFormat:    &twopb.ConfirmMsg{},
	TwoPcConfirmVoteTopicFormat:   &twopb.ConfirmVote{},
	TwoPcCommitMsgTopicFormat:     &twopb.CommitMsg{},
	TwoPcTaskResultMsgTopicFormat: &twopb.TaskResultMsg{},
}

// GossipTypeMapping is the inverse of GossipTopicMappings so that an arbitrary protobuf message
// can be mapped to a protocol ID string.
var GossipTypeMapping = make(map[reflect.Type]string, len(GossipTopicMappings))

func init() {
	for k, v := range GossipTopicMappings {
		GossipTypeMapping[reflect.TypeOf(v)] = k
	}
}
