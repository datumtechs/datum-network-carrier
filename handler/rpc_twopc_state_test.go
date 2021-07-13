package handler

import (
	"context"
	timeutils "github.com/RosettaFlow/Carrier-Go/common/timeutil"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/p2p"
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

func TestPrepareMsgRPCHandler_SendsPrepareMsg(t *testing.T) {
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

	r2 := &Service{
		cfg: &Config{
			P2P: p2,
		},
		rateLimiter: newRateLimiter(p2),
	}

	// Setup streams
	pcl := protocol.ID(p2p.RPCTwoPcPrepareMsgTopic + r.cfg.P2P.Encoding().ProtocolSuffix())
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)
	r2.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)

	prepareMsg := &twopcpb.PrepareMsg{
		ProposalId:           []byte("proposalId"),
		TaskOption:           &twopcpb.TaskOption{},
		CreateAt:             uint64(timeutils.Now().Unix()),
		Sign:                 make([]byte, 64),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		out := new(twopcpb.PrepareMsg)
		require.NoError(t, r.cfg.P2P.Encoding().DecodeWithMaxLength(stream, out))
		require.Equal(t, out.ProposalId, prepareMsg.ProposalId)
		if _, err := stream.Write([]byte{responseCodeSuccess}); err != nil {
			log.WithError(err).Error("Could not write to stream for response")
		}
		//require.NoError(t, r2.prepareMsgRPCHandler(context.Background(), prepareMsg, stream))
		require.NoError(t, stream.Close())
	})

	err := SendTwoPcPrepareMsg(context.Background(), r.cfg.P2P, p2.BHost.ID(), prepareMsg)
	require.NoError(t, err)

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}

	conns := p1.BHost.Network().ConnsToPeer(p2.BHost.ID())
	if len(conns) == 0 {
		t.Error("Peer is disconnected despite receiving a valid ping")
	}
}

func TestPrepareVoteRPCHandler_SendsPrepareVoteMsg(t *testing.T) {
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

	r2 := &Service{
		cfg: &Config{
			P2P: p2,
		},
		rateLimiter: newRateLimiter(p2),
	}

	// Setup streams
	pcl := protocol.ID(p2p.RPCTwoPcPrepareVoteTopic + r.cfg.P2P.Encoding().ProtocolSuffix())
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)
	r2.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)

	prepareVote := &twopcpb.PrepareVote{
		ProposalId:           []byte("proposalId"),
		CreateAt:             uint64(timeutils.Now().Unix()),
		Sign:                 make([]byte, 64),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		out := new(twopcpb.PrepareVote)
		require.NoError(t, r.cfg.P2P.Encoding().DecodeWithMaxLength(stream, out))
		require.Equal(t, out.ProposalId, prepareVote.ProposalId)
		if _, err := stream.Write([]byte{responseCodeSuccess}); err != nil {
			log.WithError(err).Error("Could not write to stream for response")
		}
		require.NoError(t, stream.Close())
	})

	err := SendTwoPcPrepareVote(context.Background(), r.cfg.P2P, p2.BHost.ID(), prepareVote)
	require.NoError(t, err)

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}

	conns := p1.BHost.Network().ConnsToPeer(p2.BHost.ID())
	if len(conns) == 0 {
		t.Error("Peer is disconnected despite receiving a valid ping")
	}
}

func TestConfirmMsgRPCHandler_SendsConfirmMsg(t *testing.T) {
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

	r2 := &Service{
		cfg: &Config{
			P2P: p2,
		},
		rateLimiter: newRateLimiter(p2),
	}

	// Setup streams
	pcl := protocol.ID(p2p.RPCTwoPcConfirmMsgTopic + r.cfg.P2P.Encoding().ProtocolSuffix())
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)
	r2.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)

	confirmMsg := &twopcpb.ConfirmMsg{
		ProposalId:           []byte("proposalId"),
		CreateAt:             uint64(timeutils.Now().Unix()),
		Sign:                 make([]byte, 64),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		out := new(twopcpb.ConfirmMsg)
		require.NoError(t, r.cfg.P2P.Encoding().DecodeWithMaxLength(stream, out))
		require.Equal(t, out.ProposalId, confirmMsg.ProposalId)
		if _, err := stream.Write([]byte{responseCodeSuccess}); err != nil {
			log.WithError(err).Error("Could not write to stream for response")
		}
		require.NoError(t, stream.Close())
	})

	err := SendTwoPcConfirmMsg(context.Background(), r.cfg.P2P, p2.BHost.ID(), confirmMsg)
	require.NoError(t, err)

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}

	conns := p1.BHost.Network().ConnsToPeer(p2.BHost.ID())
	if len(conns) == 0 {
		t.Error("Peer is disconnected despite receiving a valid ping")
	}
}

func TestConfirmVoteRPCHandler_SendsConfirmVote(t *testing.T) {
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

	r2 := &Service{
		cfg: &Config{
			P2P: p2,
		},
		rateLimiter: newRateLimiter(p2),
	}

	// Setup streams
	pcl := protocol.ID(p2p.RPCTwoPcConfirmVoteTopic + r.cfg.P2P.Encoding().ProtocolSuffix())
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)
	r2.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1, 1, false)

	confirmVote := &twopcpb.ConfirmVote{
		ProposalId:           []byte("proposalId"),
		CreateAt:             uint64(timeutils.Now().Unix()),
		Sign:                 make([]byte, 64),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		out := new(twopcpb.ConfirmVote)
		require.NoError(t, r.cfg.P2P.Encoding().DecodeWithMaxLength(stream, out))
		require.Equal(t, out.ProposalId, confirmVote.ProposalId)
		if _, err := stream.Write([]byte{responseCodeSuccess}); err != nil {
			log.WithError(err).Error("Could not write to stream for response")
		}
		require.NoError(t, stream.Close())
	})

	err := SendTwoPcConfirmVote(context.Background(), r.cfg.P2P, p2.BHost.ID(), confirmVote)
	require.NoError(t, err)

	if WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}

	conns := p1.BHost.Network().ConnsToPeer(p2.BHost.ID())
	if len(conns) == 0 {
		t.Error("Peer is disconnected despite receiving a valid ping")
	}
}