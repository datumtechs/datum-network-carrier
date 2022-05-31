package handler

import (
	"context"
	"github.com/datumtechs/datum-network-carrier/p2p"
	p2ptest "github.com/datumtechs/datum-network-carrier/p2p/testing"
	p2ptypes "github.com/datumtechs/datum-network-carrier/p2p/types"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"sync"
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	rlimiter := newRateLimiter(p2ptest.NewTestP2P(t))
	assert.Equal(t, len(rlimiter.limiterMap), 6, "correct number of topics not registered")
}

func TestNewRateLimiter_FreeCorrectly(t *testing.T) {
	rlimiter := newRateLimiter(p2ptest.NewTestP2P(t))
	rlimiter.free()
	assert.Equal(t, len(rlimiter.limiterMap), 0, "rate limiter not freed correctly")

}

func TestRateLimiter_ExceedCapacity(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)
	p1.Connect(p2)
	rlimiter := newRateLimiter(p1)

	// BlockByRange
	topic := p2p.RPCBlocksByRangeTopic + p1.Encoding().ProtocolSuffix()

	wg := sync.WaitGroup{}
	p2.BHost.SetStreamHandler(protocol.ID(topic), func(stream network.Stream) {
		defer wg.Done()
		code, errMsg, err := readStatusCodeNoDeadline(stream, p2.Encoding())
		require.NoError(t, err, "could not read incoming stream")
		assert.Equal(t, responseCodeInvalidRequest, code, "not equal response codes")
		assert.Equal(t, p2ptypes.ErrRateLimited.Error(), errMsg, "not equal errors")
	})
	wg.Add(1)
	stream, err := p1.BHost.NewStream(context.Background(), p2.PeerID(), protocol.ID(topic))
	require.NoError(t, err, "could not create stream")

	err = rlimiter.validateRequest(stream, 64)
	require.NoError(t, err, "could not validate incoming request")

	// Attempt to create an error, rate limit and lead to disconnect
	err = rlimiter.validateRequest(stream, 1000)
	require.NotNil(t, err, "could not get error from leaky bucket")

	require.NoError(t, stream.Close(), "could not close stream")

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}
}

func TestRateLimiter_ExceedRawCapacity(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)
	p1.Connect(p2)
	p1.Peers().Add(nil, p2.PeerID(), p2.BHost.Addrs()[0], network.DirOutbound)

	rlimiter := newRateLimiter(p1)

	// BlockByRange
	topic := p2p.RPCBlocksByRangeTopic + p1.Encoding().ProtocolSuffix()

	wg := sync.WaitGroup{}
	p2.BHost.SetStreamHandler(protocol.ID(topic), func(stream network.Stream) {
		defer wg.Done()
		code, errMsg, err := readStatusCodeNoDeadline(stream, p2.Encoding())
		require.NoError(t, err, "could not read incoming stream")
		assert.Equal(t, responseCodeInvalidRequest, code, "not equal response codes")
		assert.Equal(t, p2ptypes.ErrRateLimited.Error(), errMsg, "not equal errors")
	})
	wg.Add(1)
	stream, err := p1.BHost.NewStream(context.Background(), p2.PeerID(), protocol.ID(topic))
	require.NoError(t, err, "could not create stream")

	for i := 0; i < 2*defaultBurstLimit; i++ {
		err = rlimiter.validateRawRpcRequest(stream)
		rlimiter.addRawStream(stream)
		require.NoError(t, err, "could not validate incoming request")
	}
	// Triggers rate limit error on burst.
	assert.ErrorContains(t, rlimiter.validateRawRpcRequest(stream), p2ptypes.ErrRateLimited.Error())

	// Make Peer bad.
	for i := 0; i < defaultBurstLimit; i++ {
		assert.ErrorContains(t, rlimiter.validateRawRpcRequest(stream), p2ptypes.ErrRateLimited.Error())
	}
	//assert.Equal(t, true, p1.Peers().IsBad(p2.PeerID()), "peer is not marked as a bad peer")
	require.NoError(t, stream.Close(), "could not close stream")

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}
}

func Test_limiter_retrieveCollector_requiresLock(t *testing.T) {
	l := limiter{}
	_, err := l.retrieveCollector("")
	assert.ErrorContains(t, err,"caller must hold read/write lock")
}

// WaitTimeout will wait for a WaitGroup to resolve within a timeout interval.
// Returns true if the waitgroup exceeded the timeout.
func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	select {
	case <-ch:
		return false
	case <-time.After(timeout):
		return true
	}
}