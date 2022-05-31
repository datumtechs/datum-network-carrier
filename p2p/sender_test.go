package p2p

import (
	"context"
	carrierrpcdebugpbv1 "github.com/datumtechs/datum-network-carrier/pb/carrier/rpc/debug/v1"
	testp2p "github.com/datumtechs/datum-network-carrier/p2p/testing"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestService_Send(t *testing.T) {
	p1 := testp2p.NewTestP2P(t)
	p2 := testp2p.NewTestP2P(t)
	p1.Connect(p2)

	svc := &Service{
		host: p1.BHost,
		cfg:  &Config{},
	}

	msg := &carrierrpcdebugpbv1.GossipTestData{
		Data:                 nil,
		Count:                1,
		Step:                 3,
	}

	// Register external listener which will repeat the message back.
	var wg sync.WaitGroup
	wg.Add(1)
	topic := "/testing/1"
	RPCTopicMappings[topic] = new(carrierrpcdebugpbv1.GossipTestData)
	defer func() {
		delete(RPCTopicMappings, topic)
	}()
	p2.SetStreamHandler(topic+"/ssz_snappy", func(stream network.Stream) {
		rcvd := &carrierrpcdebugpbv1.GossipTestData{}
		require.NoError(t, svc.Encoding().DecodeWithMaxLength(stream, rcvd))
		_, err := svc.Encoding().EncodeWithMaxLength(stream, rcvd)
		require.NoError(t, err)
		assert.NoError(t, stream.Close())
		wg.Done()
	})

	stream, err := svc.Send(context.Background(), msg, "/testing/1", p2.BHost.ID())
	require.NoError(t, err)

	WaitTimeout(&wg, 1*time.Second)

	rcvd := &carrierrpcdebugpbv1.GossipTestData{}
	require.NoError(t, svc.Encoding().DecodeWithMaxLength(stream, rcvd))
	if !proto.Equal(rcvd, msg) {
		t.Errorf("Expected identical message to be received. got %v want %v", rcvd, msg)
	}
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

