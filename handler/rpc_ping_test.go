package handler

import (
	"context"
	pb "github.com/RosettaFlow/Carrier-Go/lib/p2p/v1"
	p2ptest "github.com/RosettaFlow/Carrier-Go/p2p/testing"
	p2ptypes "github.com/RosettaFlow/Carrier-Go/p2p/types"
	"github.com/ethereum/go-ethereum/p2p/enr"
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

func TestPingRPCHandler_ReceivesPing(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)
	p1.Connect(p2)
	assert.Equal(t, 1, len(p1.BHost.Network().Peers()), "Expected peers to be connected")
	p1.LocalMetadata = &pb.MetaData{
		SeqNumber: 2,
		Attnets:   []byte{'A', 'B'},
	}

	p2.LocalMetadata = &pb.MetaData{
		SeqNumber: 2,
		Attnets:   []byte{'C', 'D'},
	}

	// Set up a head state in the database with data we expect.
	r := &Service{
		cfg: &Config{
			P2P: p1,
		},
		rateLimiter: newRateLimiter(p1),
	}

	p1.Peers().Add(new(enr.Record), p2.BHost.ID(), p2.BHost.Addrs()[0], network.DirUnknown)
	p1.Peers().SetMetadata(p2.BHost.ID(), p2.LocalMetadata)

	// Setup streams
	pcl := protocol.ID("/testing")
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)
	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		expectSuccess(t, stream)
		out := new(types.SSZUint64)
		require.NoError(t, r.cfg.P2P.Encoding().DecodeWithMaxLength(stream, out))
		assert.Equal(t, uint64(2), uint64(*out))
	})
	stream1, err := p1.BHost.NewStream(context.Background(), p2.BHost.ID(), pcl)
	require.NoError(t, err)
	seqNumber := types.SSZUint64(2)

	require.NoError(t, r.pingHandler(context.Background(), &seqNumber, stream1))

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}

	conns := p1.BHost.Network().ConnsToPeer(p2.BHost.ID())
	if len(conns) == 0 {
		t.Error("Peer is disconnected despite receiving a valid ping")
	}
}

func TestPingRPCHandler_SendsPing(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)
	p1.Connect(p2)
	assert.Equal(t, 1, len(p1.BHost.Network().Peers()), "Expected peers to be connected")
	p1.LocalMetadata = &pb.MetaData{
		SeqNumber: 2,
		Attnets:   []byte{'A', 'B'},
	}

	p2.LocalMetadata = &pb.MetaData{
		SeqNumber: 2,
		Attnets:   []byte{'C', 'D'},
	}

	// Set up a head state in the database with data we expect.
	r := &Service{
		cfg: &Config{
			P2P: p1,
		},
		rateLimiter: newRateLimiter(p1),
	}

	p1.Peers().Add(new(enr.Record), p2.BHost.ID(), p2.BHost.Addrs()[0], network.DirUnknown)
	p1.Peers().SetMetadata(p2.BHost.ID(), p2.LocalMetadata)

	p2.Peers().Add(new(enr.Record), p1.BHost.ID(), p1.BHost.Addrs()[0], network.DirUnknown)
	p2.Peers().SetMetadata(p1.BHost.ID(), p1.LocalMetadata)

	r2 := &Service{
		cfg: &Config{
			P2P: p2,
		},
		rateLimiter: newRateLimiter(p2),
	}
	// Setup streams
	pcl := protocol.ID("/rosetta/carrier_chain/req/ping/1/ssz_snappy")
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)

	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		out := new(types.SSZUint64)
		require.NoError(t, r2.cfg.P2P.Encoding().DecodeWithMaxLength(stream, out))
		assert.Equal(t, uint64(2), uint64(*out))
		require.NoError(t, r2.pingHandler(context.Background(), out, stream))
	})

	require.NoError(t, r.sendPingRequest(context.Background(), p2.BHost.ID()))

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}

	conns := p1.BHost.Network().ConnsToPeer(p2.BHost.ID())
	if len(conns) == 0 {
		t.Error("Peer is disconnected despite receiving a valid ping")
	}
}

func TestPingRPCHandler_BadSequenceNumber(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)
	p1.Connect(p2)
	assert.Equal(t, 1, len(p1.BHost.Network().Peers()), "Expected peers to be connected")
	p1.LocalMetadata = &pb.MetaData{
		SeqNumber: 2,
		Attnets:   []byte{'A', 'B'},
	}

	p2.LocalMetadata = &pb.MetaData{
		SeqNumber: 2,
		Attnets:   []byte{'C', 'D'},
	}

	// Set up a head state in the database with data we expect.
	r := &Service{
		cfg: &Config{
			P2P: p1,
		},
		rateLimiter: newRateLimiter(p1),
	}

	badMetadata := &pb.MetaData{
		SeqNumber: 3,
		Attnets:   []byte{'E', 'F'},
	}

	p1.Peers().Add(new(enr.Record), p2.BHost.ID(), p2.BHost.Addrs()[0], network.DirUnknown)
	p1.Peers().SetMetadata(p2.BHost.ID(), badMetadata)

	// Setup streams
	pcl := protocol.ID("/testing")
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)
	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		expectFailure(t, responseCodeInvalidRequest, p2ptypes.ErrInvalidSequenceNum.Error(), stream)

	})
	stream1, err := p1.BHost.NewStream(context.Background(), p2.BHost.ID(), pcl)
	require.NoError(t, err)

	wantedSeq := types.SSZUint64(p2.LocalMetadata.SeqNumber)
	err = r.pingHandler(context.Background(), &wantedSeq, stream1)
	assert.ErrorContains(t, err, p2ptypes.ErrInvalidSequenceNum.Error())

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}

	res, err := p1.Peers().Scorers().BadResponsesScorer().Count(p2.BHost.ID())
	require.NoError(t, err)
	assert.Equal(t, 1, res, "Peer wasn't penalised for providing a bad sequence number")
}

