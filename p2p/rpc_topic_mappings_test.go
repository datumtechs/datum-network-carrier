package p2p

import (
	carrierp2ppbv1 "github.com/datumtechs/datum-network-carrier/pb/carrier/p2p/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifyRPCMappings(t *testing.T) {
	assert.NoError(t, VerifyTopicMapping(RPCStatusTopic, &carrierp2ppbv1.Status{}), "Failed to verify status rpc topic")
	assert.NotNil(t, VerifyTopicMapping(RPCStatusTopic, new([]byte)), "Incorrect message type verified for status rpc topic")

	assert.NoError(t, VerifyTopicMapping(RPCMetaDataTopic, new(interface{})), "Failed to verify metadata rpc topic")
	assert.NotNil(t, VerifyTopicMapping(RPCStatusTopic, new([]byte)), "Incorrect message type verified for metadata rpc topic")
}

