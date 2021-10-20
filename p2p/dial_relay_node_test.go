package p2p

import (
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"testing"
)

func TestMakePeer_InvalidMultiaddress(t *testing.T) {
	_, err := MakePeer("/ip4")
	assert.ErrorContains(t, err, "failed to parse multiaddr \"/ip4\"",  "Expect error when invalid multiaddress was provided")
}

func TestMakePeer_OK(t *testing.T) {
	a, err := MakePeer("/ip4/127.0.0.1/tcp/3333/p2p/QmUn6ycS8Fu6L462uZvuEfDoSgYX6kqP4aSZWMa7z1tWAX")
	require.NoError(t, err, "Unexpected error when making a valid peer")
	assert.Equal(t, "QmUn6ycS8Fu6L462uZvuEfDoSgYX6kqP4aSZWMa7z1tWAX", a.ID.Pretty(), "Unexpected peer ID")
}
