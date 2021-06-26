package p2p

import (
	libp2ppb "github.com/RosettaFlow/Carrier-Go/lib/p2p/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifyRPCMappings(t *testing.T) {
	assert.NoError(t, VerifyTopicMapping(RPCStatusTopic, &libp2ppb.Status{}), "Failed to verify status rpc topic")
	assert.NotNil(t, VerifyTopicMapping(RPCStatusTopic, new([]byte)), "Incorrect message type verified for status rpc topic")

	assert.NoError(t, VerifyTopicMapping(RPCMetaDataTopic, new(interface{})), "Failed to verify metadata rpc topic")
	assert.NotNil(t, VerifyTopicMapping(RPCStatusTopic, new([]byte)), "Incorrect message type verified for metadata rpc topic")
}

