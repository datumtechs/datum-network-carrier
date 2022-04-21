package handler

import (
	"context"
	p2ptest "github.com/Metisnetwork/Metis-Carrier/p2p/testing"
	p2ptypes "github.com/Metisnetwork/Metis-Carrier/p2p/types"
	"github.com/kevinms/leakybucket-go"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	types "github.com/prysmaticlabs/eth2-types"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"sync"
	"testing"
	"time"
)

func TestGoodByeRPCHandler_Disconnects_With_Peer(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)
	p1.Connect(p2)
	assert.Equal(t, 1, len(p1.BHost.Network().Peers()), "Expected peers to be connected")

	// Set up a head state in the database with data we expect.
	r := &Service{
		cfg: &Config{
			P2P: p1,
		},
		rateLimiter: newRateLimiter(p1),
	}

	// Setup streams
	pcl := protocol.ID("/testing")
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)
	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		expectResetStream(t, stream)
	})
	stream1, err := p1.BHost.NewStream(context.Background(), p2.BHost.ID(), pcl)
	require.NoError(t, err)
	failureCode := p2ptypes.GoodbyeCodeClientShutdown

	require.NoError(t, r.goodbyeRPCHandler(context.Background(), &failureCode, stream1))

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}

	conns := p1.BHost.Network().ConnsToPeer(p2.BHost.ID())
	if len(conns) > 0 {
		t.Error("Peer is still not disconnected despite sending a goodbye message")
	}
}

func TestGoodByeRPCHandler_BackOffPeer(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)
	p3 := p2ptest.NewTestP2P(t)

	p1.Connect(p2)
	p1.Connect(p3)
	assert.Equal(t, 2, len(p1.BHost.Network().Peers()), "Expected peers to be connected")

	// Set up a head state in the database with data we expect.
	r := &Service{
		cfg: &Config{
			P2P: p1,
		},
		rateLimiter: newRateLimiter(p1),
	}

	// Setup streams
	pcl := protocol.ID("/testing")
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)
	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		expectResetStream(t, stream)
	})
	stream1, err := p1.BHost.NewStream(context.Background(), p2.BHost.ID(), pcl)
	require.NoError(t, err)
	failureCode := p2ptypes.GoodbyeCodeClientShutdown

	require.NoError(t, r.goodbyeRPCHandler(context.Background(), &failureCode, stream1))

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}

	conns := p1.BHost.Network().ConnsToPeer(p2.BHost.ID())
	if len(conns) > 0 {
		t.Error("Peer is still not disconnected despite sending a goodbye message")
	}
	valTime, err := p1.Peers().NextValidTime(p2.BHost.ID())
	require.NoError(t, err)
	expectedTime := time.Now().Add(backOffTime[failureCode])
	diff := expectedTime.Sub(valTime)
	// Add a little bit of allowance
	require.Equal(t, true, diff.Seconds() <= 1)

	wg.Add(1)
	p3.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		expectResetStream(t, stream)
	})

	stream2, err := p1.BHost.NewStream(context.Background(), p3.BHost.ID(), pcl)
	require.NoError(t, err)
	failureCode = p2ptypes.GoodbyeCodeBanned

	require.NoError(t, r.goodbyeRPCHandler(context.Background(), &failureCode, stream2))

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}

	conns = p1.BHost.Network().ConnsToPeer(p3.BHost.ID())
	if len(conns) > 0 {
		t.Error("Peer is still not disconnected despite sending a goodbye message")
	}
	valTime, err = p1.Peers().NextValidTime(p3.BHost.ID())
	require.NoError(t, err)
	expectedTime = time.Now().Add(backOffTime[failureCode])
	diff = expectedTime.Sub(valTime)
	// Add a little bit of allowance
	require.Equal(t, true, diff.Seconds() <= 1)
}

func TestSendGoodbye_SendsMessage(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)
	p1.Connect(p2)
	assert.Equal(t, 1, len(p1.BHost.Network().Peers()), "Expected peers to be connected")

	// Set up a head state in the database with data we expect.
	r := &Service{
		cfg: &Config{
			P2P: p1,
		},
		rateLimiter: newRateLimiter(p1),
	}
	failureCode := p2ptypes.GoodbyeCodeClientShutdown

	// Setup streams
	pcl := protocol.ID("/rosettanet/carrier_chain/req/goodbye/1/ssz_snappy")
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)
	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		out := new(types.SSZUint64)
		require.NoError(t, r.cfg.P2P.Encoding().DecodeWithMaxLength(stream, out))
		assert.Equal(t, failureCode, *out)
		require.NoError(t, stream.Close())
	})

	err := r.sendGoodByeMessage(context.Background(), failureCode, p2.BHost.ID())
	require.NoError(t, err)

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}

	conns := p1.BHost.Network().ConnsToPeer(p1.BHost.ID())
	if len(conns) > 0 {
		t.Error("Peer is still not disconnected despite sending a goodbye message")
	}
}

func TestSendGoodbye_DisconnectWithPeer(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)
	p1.Connect(p2)
	assert.Equal(t, 1, len(p1.BHost.Network().Peers()), "Expected peers to be connected")

	// Set up a head state in the database with data we expect.
	r := &Service{
		cfg: &Config{
			P2P: p1,
		},
		rateLimiter: newRateLimiter(p1),
	}
	failureCode := p2ptypes.GoodbyeCodeClientShutdown

	// Setup streams
	pcl := protocol.ID("/rosettanet/carrier_chain/req/goodbye/1/ssz_snappy")
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)
	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		out := new(types.SSZUint64)
		require.NoError(t, r.cfg.P2P.Encoding().DecodeWithMaxLength(stream, out))
		assert.Equal(t, failureCode, *out)
		require.NoError(t, stream.Close())
	})

	require.NoError(t, r.sendGoodByeAndDisconnect(context.Background(), failureCode, p2.BHost.ID()))
	conns := p1.BHost.Network().ConnsToPeer(p2.BHost.ID())
	if len(conns) > 0 {
		t.Error("Peer is still not disconnected despite sending a goodbye message")
	}

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}

}

