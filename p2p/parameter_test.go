package p2p

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"gotest.tools/assert"
	"testing"
)

func TestOverlayParameters(t *testing.T) {
	setPubSubParameters()
	assert.Equal(t, gossipSubD, pubsub.GossipSubD, "gossipSubD")
	assert.Equal(t, gossipSubDlo, pubsub.GossipSubDlo, "gossipSubDlo")
	assert.Equal(t, gossipSubDhi, pubsub.GossipSubDhi, "gossipSubDhi")
}
