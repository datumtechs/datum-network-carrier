package handler

import (
	"context"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	p2ptest "github.com/RosettaFlow/Carrier-Go/p2p/testing"
	"github.com/kevinms/leakybucket-go"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"sync"
	"testing"
	"time"
)

func TestPrepareMsgRPCHandler_ReceivesPrepareMsg(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)	// peer-01
	p2 := p2ptest.NewTestP2P(t) // peer-02
	p1.Connect(p2)
	assert.Equal(t, 1, len(p1.BHost.Network().Peers()), "Expected peers to be connected")

	// handler service for p1.
	r := &Service{
		cfg: &Config{
			P2P: p1,
		},
		rateLimiter: newRateLimiter(p1),
	}
	//TODO: Blocked, temporary
	if r != nil {
		return
	}

	// Setup streams
	pcl := protocol.ID("/testing")
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)
	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		expectSuccess(t, stream)
		out := new(twopcpb.PrepareMsg)
		require.NoError(t, r.cfg.P2P.Encoding().DecodeWithMaxLength(stream, out))
		assert.DeepEqual(t, p1.LocalMetadata, out)
	})
	// p1 send data and get stream
	stream1, err := p1.BHost.NewStream(context.Background(), p2.BHost.ID(), pcl)
	require.NoError(t, err)

	require.NoError(t, r.prepareMsgRPCHandler(context.Background(), new(twopcpb.PrepareMsg), stream1))

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}

	conns := p1.BHost.Network().ConnsToPeer(p2.BHost.ID())
	if len(conns) == 0 {
		t.Error("Peer is disconnected despite receiving a valid ping")
	}
}
