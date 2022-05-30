package handler

import (
	"context"
	"github.com/datumtechs/datum-network-carrier/common/bytesutil"
	pb "github.com/datumtechs/datum-network-carrier/pb/carrier/p2p/v1"
	"github.com/datumtechs/datum-network-carrier/p2p/peers"
	p2ptest "github.com/datumtechs/datum-network-carrier/p2p/testing"
	p2ptypes "github.com/datumtechs/datum-network-carrier/p2p/types"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/gogo/protobuf/proto"
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

func TestStatusRPCHandler_Disconnects_OnForkVersionMismatch(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)
	p1.Connect(p2)
	assert.Equal(t, 1, len(p1.BHost.Network().Peers()), "Expected peers to be connected")
	//root := [32]byte{'C'}

	r := &Service{
		cfg: &Config{
			P2P: p1,
		},
		rateLimiter: newRateLimiter(p1),
	}
	pcl := protocol.ID("/testing")
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)

	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		expectSuccess(t, stream)
		out := &pb.Status{}
		require.NoError(t, r.cfg.P2P.Encoding().DecodeWithMaxLength(stream, out))
		//assert.DeepEqual(t, root[:], out.FinalizedRoot)
		require.NoError(t, stream.Close())
	})

	pcl2 := protocol.ID("/rosettanet/carrier_chain/req/goodbye/1/ssz_snappy")
	topic = string(pcl2)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)
	var wg2 sync.WaitGroup
	wg2.Add(1)
	p2.BHost.SetStreamHandler(pcl2, func(stream network.Stream) {
		defer wg2.Done()
		msg := new(types.SSZUint64)
		require.NoError(t, r.cfg.P2P.Encoding().DecodeWithMaxLength(stream, msg))
		assert.Equal(t, p2ptypes.GoodbyeCodeWrongNetwork, *msg)
		require.NoError(t, stream.Close())
	})

	stream1, err := p1.BHost.NewStream(context.Background(), p2.BHost.ID(), pcl)
	require.NoError(t, err)
	require.NoError(t, r.statusRPCHandler(context.Background(), &pb.Status{ForkDigest: bytesutil.PadTo([]byte("f"), 4), HeadRoot: make([]byte, 32), FinalizedRoot: make([]byte, 32)}, stream1))

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}
	//if WaitTimeout(&wg2, 1*time.Second) {
	//	t.Fatal("Did not receive stream within 1 sec")
	//}

	//assert.Equal(t, 0, len(p1.BHost.Network().Peers()), "handler did not disconnect peer")
}

func TestStatusRPCHandler_ConnectsOnGenesis(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)
	p1.Connect(p2)
	assert.Equal(t, 1, len(p1.BHost.Network().Peers()), "Expected peers to be connected")
	root := [32]byte{}

	r := &Service{
		cfg: &Config{
			P2P: p1,
		},
		rateLimiter: newRateLimiter(p1),
	}
	pcl := protocol.ID("/testing")
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)

	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		expectSuccess(t, stream)
		out := &pb.Status{}
		require.NoError(t, r.cfg.P2P.Encoding().DecodeWithMaxLength(stream, out))
		assert.DeepEqual(t, root[:], out.FinalizedRoot)
	})

	stream1, err := p1.BHost.NewStream(context.Background(), p2.BHost.ID(), pcl)
	require.NoError(t, err)
	digest, err := r.forkDigest()
	require.NoError(t, err)

	err = r.statusRPCHandler(context.Background(), &pb.Status{ForkDigest: digest[:], FinalizedRoot: make([]byte, 32)}, stream1)
	require.NoError(t, err)

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}

	assert.Equal(t, 1, len(p1.BHost.Network().Peers()), "Handler disconnected with peer")
}

func TestStatusRPCHandler_ReturnsHelloMessage(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)
	p1.Connect(p2)
	assert.Equal(t, 1, len(p1.BHost.Network().Peers()), "Expected peers to be connected")

	// Set up a head state with data we expect.
	r := &Service{
		cfg: &Config{
			P2P: p1,
		},
		rateLimiter: newRateLimiter(p1),
	}
	digest, err := r.forkDigest()
	require.NoError(t, err)

	// Setup streams
	pcl := protocol.ID("/testing")
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)
	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		expectSuccess(t, stream)
		out := &pb.Status{}
		require.NoError(t, r.cfg.P2P.Encoding().DecodeWithMaxLength(stream, out))
		expected := &pb.Status{
			ForkDigest:     digest[:],
			HeadSlot:       0,
			HeadRoot:       make([]byte, 32),
			FinalizedEpoch: 0,
			FinalizedRoot:  make([]byte, 32),
		}
		if !proto.Equal(out, expected) {
			t.Errorf("Did not receive expected message. Got %+v wanted %+v", out, expected)
		}
	})
	stream1, err := p1.BHost.NewStream(context.Background(), p2.BHost.ID(), pcl)
	require.NoError(t, err)

	err = r.statusRPCHandler(context.Background(), &pb.Status{
		ForkDigest:     digest[:],
		FinalizedRoot:  make([]byte, 32),
		FinalizedEpoch: 3,
	}, stream1)
	require.NoError(t, err)

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}
}

func TestHandshakeHandlers_Roundtrip(t *testing.T) {
	// Scenario is that p1 and p2 connect, exchange handshakes.
	// p2 disconnects and p1 should forget the handshake status.
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)

	p1.LocalMetadata = &pb.MetaData{
		SeqNumber: 2,
		Attnets:   bytesutil.PadTo([]byte{'A', 'B'}, 8),
	}

	p2.LocalMetadata = &pb.MetaData{
		SeqNumber: 2,
		Attnets:   bytesutil.PadTo([]byte{'C', 'D'}, 8),
	}
	r := &Service{
		cfg: &Config{
			P2P: p1,
		},
		ctx:         context.Background(),
		rateLimiter: newRateLimiter(p1),
	}
	digest, err := r.forkDigest()
	p1.Digest = digest
	require.NoError(t, err)

	r2 := &Service{
		cfg: &Config{
			P2P: p2,
		},
		rateLimiter: newRateLimiter(p2),
	}
	p2.Digest, err = r.forkDigest()
	require.NoError(t, err)

	r.Start()

	// Setup streams
	pcl := protocol.ID("/rosettanet/carrier_chain/req/status/1/ssz_snappy")
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)
	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		out := &pb.Status{}
		require.NoError(t, r.cfg.P2P.Encoding().DecodeWithMaxLength(stream, out))
		log.WithField("status", out).Warn("received status")
		resp := &pb.Status{
			HeadSlot:       0,
			HeadRoot:       make([]byte, 32),
			ForkDigest:     make([]byte, 4),
			FinalizedRoot:  make([]byte, 32),
			FinalizedEpoch: 0}
		_, err := stream.Write([]byte{responseCodeSuccess})
		require.NoError(t, err)
		_, err = r.cfg.P2P.Encoding().EncodeWithMaxLength(stream, resp)
		require.NoError(t, err)
		log.WithField("status", out).Warn("sending status")
		if err := stream.Close(); err != nil {
			t.Log(err)
		}
	})

	pcl = "/rosettanet/carrier_chain/req/ping/1/ssz_snappy"
	topic = string(pcl)
	r2.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)
	var wg2 sync.WaitGroup
	wg2.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg2.Done()
		out := new(types.SSZUint64)
		require.NoError(t, r.cfg.P2P.Encoding().DecodeWithMaxLength(stream, out))
		assert.Equal(t, uint64(2), uint64(*out))
		require.NoError(t, r2.pingHandler(context.Background(), out, stream))
		require.NoError(t, stream.Close())
	})

	numInactive1 := len(p1.Peers().Inactive())
	numActive1 := len(p1.Peers().Active())

	p1.Connect(p2)

	p1.Peers().Add(new(enr.Record), p2.BHost.ID(), p2.BHost.Addrs()[0], network.DirUnknown)
	p1.Peers().SetMetadata(p2.BHost.ID(), p2.LocalMetadata)

	p2.Peers().Add(new(enr.Record), p1.BHost.ID(), p1.BHost.Addrs()[0], network.DirUnknown)
	p2.Peers().SetMetadata(p1.BHost.ID(), p1.LocalMetadata)

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}
	if WaitTimeout(&wg2, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}

	// Wait for stream buffer to be read.
	time.Sleep(200 * time.Millisecond)

	numInactive2 := len(p1.Peers().Inactive())
	numActive2 := len(p1.Peers().Active())

	assert.Equal(t, numInactive1, numInactive1, "Number of inactive peers changed unexpectedly")
	assert.Equal(t, numActive1+1, numActive2, "Number of active peers unexpected")

	require.NoError(t, p2.Disconnect(p1.PeerID()))
	p1.Peers().SetConnectionState(p2.PeerID(), peers.PeerDisconnected)

	// Wait for disconnect evengine to trigger.
	time.Sleep(200 * time.Millisecond)

	numInactive3 := len(p1.Peers().Inactive())
	numActive3 := len(p1.Peers().Active())
	assert.Equal(t, numInactive2+1, numInactive3, "Number of inactive peers unexpected")
	assert.Equal(t, numActive2-1, numActive3, "Number of active peers unexpected")
}

func TestStatusRPCRequest_RequestSent(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)

	r := &Service{
		cfg: &Config{
			P2P: p1,
		},
		ctx:         context.Background(),
		rateLimiter: newRateLimiter(p1),
	}

	// Setup streams
	pcl := protocol.ID("/rosettanet/carrier_chain/req/status/1/ssz_snappy")
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)
	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		out := &pb.Status{}
		require.NoError(t, r.cfg.P2P.Encoding().DecodeWithMaxLength(stream, out))
		digest, err := r.forkDigest()
		require.NoError(t, err)
		expected := &pb.Status{
			ForkDigest:     digest[:],
			HeadSlot:       0,
			HeadRoot:       make([]byte, 32),
			FinalizedEpoch: 0,
			FinalizedRoot:  make([]byte, 32),
		}
		if !proto.Equal(out, expected) {
			t.Errorf("Did not receive expected message. Got %+v wanted %+v", out, expected)
		}
	})

	p1.AddConnectionHandler(r.sendRPCStatusRequest, nil)
	p1.Connect(p2)

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}

	assert.Equal(t, 1, len(p1.BHost.Network().Peers()), "Expected peers to continue being connected")
}

func TestStatusRPCRequest_BadPeerHandshake(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)
	p1.LocalMetadata = &pb.MetaData{}
	p2.LocalMetadata = &pb.MetaData{
		SeqNumber:            1,
		Attnets:              []byte{},
	}
	r := &Service{
		cfg: &Config{
			P2P: p1,
		},
		ctx:         context.Background(),
		rateLimiter: newRateLimiter(p1),
	}
	r.Start()

	// Setup streams
	pcl := protocol.ID("/rosettanet/carrier_chain/req/status/1/ssz_snappy")
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)
	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		out := &pb.Status{}
		require.NoError(t, r.cfg.P2P.Encoding().DecodeWithMaxLength(stream, out))
		expected := &pb.Status{
			ForkDigest:     []byte{1, 1, 1, 1},
			HeadSlot:       0,
			HeadRoot:       make([]byte, 32),
			FinalizedEpoch: 5,
			FinalizedRoot:  make([]byte, 32),
		}
		if _, err := stream.Write([]byte{responseCodeSuccess}); err != nil {
			log.WithError(err).Debug("Could not write to stream")
		}
		_, err := r.cfg.P2P.Encoding().EncodeWithMaxLength(stream, expected)
		require.NoError(t, err)
	})

	assert.Equal(t, false, p1.Peers().Scorers().IsBadPeer(p2.PeerID()), "Peer is marked as bad")
	p1.Connect(p2)

	if WaitTimeout(&wg, time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}
	time.Sleep(100 * time.Millisecond)

	connectionState, err := p1.Peers().ConnectionState(p2.PeerID())
	require.NoError(t, err, "Could not obtain peer connection state")
	//assert.Equal(t, peers.PeerDisconnected, connectionState, "Expected peer to be disconnected")
	assert.Equal(t, peers.PeerConnected, connectionState, "Expected peer to be disconnected")

	//assert.Equal(t, true, p1.Peers().Scorers().IsBadPeer(p2.PeerID()), "Peer is not marked as bad")
}
