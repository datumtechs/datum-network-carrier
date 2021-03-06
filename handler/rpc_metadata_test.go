package handler

import (
	"context"
	pb "github.com/RosettaFlow/Carrier-Go/lib/p2p/v1"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	p2ptest "github.com/RosettaFlow/Carrier-Go/p2p/testing"
	"github.com/kevinms/leakybucket-go"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestMetaDataRPCHandler_ReceivesMetadata(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)
	p1.Connect(p2)
	assert.Equal(t, 1, len(p1.BHost.Network().Peers()), "Expected peers to be connected")
	bitfield := [8]byte{'A', 'B'}
	p1.LocalMetadata = &pb.MetaData{
		SeqNumber: 2,
		Attnets:   bitfield[:],
	}

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
		expectSuccess(t, stream)
		out := new(pb.MetaData)
		require.NoError(t, r.cfg.P2P.Encoding().DecodeWithMaxLength(stream, out))
		assert.DeepEqual(t, p1.LocalMetadata, out)
	})
	stream1, err := p1.BHost.NewStream(context.Background(), p2.BHost.ID(), pcl)
	require.NoError(t, err)

	require.NoError(t, r.metaDataHandler(context.Background(), new(interface{}), stream1))

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}

	conns := p1.BHost.Network().ConnsToPeer(p2.BHost.ID())
	if len(conns) == 0 {
		t.Error("Peer is disconnected despite receiving a valid ping")
	}
}

func TestMetadataRPCHandler_SendsMetadata(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)
	p1.Connect(p2)
	assert.Equal(t, 1, len(p1.BHost.Network().Peers()), "Expected peers to be connected")
	bitfield := [8]byte{'A', 'B'}
	p2.LocalMetadata = &pb.MetaData{
		SeqNumber: 2,
		Attnets:   bitfield[:],
	}

	// Set up a head state in the database with data we expect.
	r := &Service{
		cfg: &Config{
			P2P: p1,
		},
		rateLimiter: newRateLimiter(p1),
	}

	r2 := &Service{
		cfg: &Config{
			P2P: p2,
		},
		rateLimiter: newRateLimiter(p2),
	}

	// Setup streams
	pcl := protocol.ID(p2p.RPCMetaDataTopic + r.cfg.P2P.Encoding().ProtocolSuffix())
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)
	r2.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)

	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		require.NoError(t, r2.metaDataHandler(context.Background(), new(interface{}), stream))
	})

	metadata, err := r.sendMetaDataRequest(context.Background(), p2.BHost.ID())
	require.NoError(t, err)

	if !reflect.DeepEqual(metadata, p2.LocalMetadata) {
		t.Fatalf("Metadata unequal, received %v but wanted %v", metadata, p2.LocalMetadata)
	}

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}

	conns := p1.BHost.Network().ConnsToPeer(p2.BHost.ID())
	if len(conns) == 0 {
		t.Error("Peer is disconnected despite receiving a valid ping")
	}
}

