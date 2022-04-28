package handler

import (
	"context"
	libp2ppb "github.com/Metisnetwork/Metis-Carrier/lib/rpc/debug/v1"
	carrrierP2P "github.com/Metisnetwork/Metis-Carrier/p2p"
	"github.com/Metisnetwork/Metis-Carrier/p2p/encoder"
	p2ptest "github.com/Metisnetwork/Metis-Carrier/p2p/testing"
	libp2pcore "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"sync"
	"testing"
	"time"
)

// expectSuccess status code from a stream in regular sync.
func expectSuccess(t *testing.T, stream network.Stream) {
	code, errMsg, err := ReadStatusCode(stream, &encoder.SszNetworkEncoder{})
	require.NoError(t, err)
	require.Equal(t, uint8(0), code, "Received non-zero response code")
	require.Equal(t, "", errMsg, "Received error message from stream")
}

// expectSuccess status code from a stream in regular sync.
func expectFailure(t *testing.T, expectedCode uint8, expectedErrorMsg string, stream network.Stream) {
	code, errMsg, err := ReadStatusCode(stream, &encoder.SszNetworkEncoder{})
	require.NoError(t, err)
	require.NotEqual(t, uint8(0), code, "Expected request to fail but got a 0 response code")
	require.Equal(t, expectedCode, code, "Received incorrect response code")
	require.Equal(t, expectedErrorMsg, errMsg)
}

// expectResetStream status code from a stream in regular sync.
func expectResetStream(t *testing.T, stream network.Stream) {
	expectedErr := "stream reset"
	_, _, err := ReadStatusCode(stream, &encoder.SszNetworkEncoder{})
	assert.ErrorContains(t, err, expectedErr)
}

func TestRegisterRPC_ReceivesValidMessage(t *testing.T) {
	p2p := p2ptest.NewTestP2P(t)
	r := &Service{
		ctx:         context.Background(),
		cfg:         &Config{P2P: p2p},
		rateLimiter: newRateLimiter(p2p),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	topic := "/testing/foobar/1"
	handler := func(ctx context.Context, msg interface{}, stream libp2pcore.Stream) error {
		m, ok := msg.(*libp2ppb.GossipTestData)
		if !ok {
			t.Error("Object is not of type *pb.TestSimpleMessage")
		}
		assert.DeepEqual(t, []byte("gossipData"), m.Data)
		wg.Done()

		return nil
	}
	carrrierP2P.RPCTopicMappings[topic] = new(libp2ppb.GossipTestData)
	// Cleanup Topic mappings
	defer func() {
		delete(carrrierP2P.RPCTopicMappings, topic)
	}()
	r.registerRPC(topic, handler)

	p2p.ReceiveRPC(topic, &libp2ppb.GossipTestData{Data: []byte("gossipData"), Step: 1, Count: 3})

	if WaitTimeout(&wg, time.Second) {
		t.Fatal("Did not receive RPC in 1 second")
	}
}

