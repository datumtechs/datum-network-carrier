package p2p

import (
	"gotest.tools/assert"
	"testing"
)

func TestMakePeer_InvalidMultiaddress(t *testing.T) {
	_, err := MakePeer("/ip4")
	assert.ErrorContains(t, err, "failed to parse multiaddr \"/ip4\"",  "Expect error when invalid multiaddress was provided")
}
