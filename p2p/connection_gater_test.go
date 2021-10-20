package p2p

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/p2p/peers"
	"github.com/RosettaFlow/Carrier-Go/p2p/peers/scorers"
	"github.com/kevinms/leakybucket-go"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPeer_AtMaxLimit(t *testing.T) {
	// init: create host and remote peer
	ipAddr, pkey := createAddrAndPrivKey(t)
	ipAddr2, pkey2 := createAddrAndPrivKey(t)

	listen, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", ipAddr, 2000))
	require.NoError(t, err, "Failed to p2p listen")
	s := &Service{
		ipLimiter: leakybucket.NewCollector(ipLimit, ipBurst, false),
	}
	s.peers = peers.NewStatus(context.Background(), &peers.StatusConfig{
		PeerLimit: 0,
		ScorerParams: &scorers.Config{
			BadResponsesScorerConfig: &scorers.BadResponsesScorerConfig{
				Threshold: 3,
			},
		},
	})
	s.cfg = &Config{ MaxPeers: 0 }
	s.addrFilter, err = configureFilter(&Config{})
	require.NoError(t, err)
	h1, err := libp2p.New(context.Background(), []libp2p.Option{privKeyOption(pkey), libp2p.ListenAddrs(listen), libp2p.ConnectionGater(s)}...)
	require.NoError(t, err)
	s.host = h1
	defer func() {
		err := h1.Close()
		require.NoError(t, err)
	}()

	for i := 0; i < highWatermarkBuffer; i++ {
		addPeer(t, s.peers, peers.PeerConnected)
	}

	// create alternate host
	listen, err = ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", ipAddr2, 3000))
	require.NoError(t, err, "Failed to p2p listen")
	h2, err := libp2p.New(context.Background(), []libp2p.Option{privKeyOption(pkey2), libp2p.ListenAddrs(listen)}...)
	require.NoError(t, err)
	defer func() {
		err := h2.Close()
		require.NoError(t, err)
	}()
	multiAddress, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d/p2p/%s", ipAddr, 2000, h1.ID()))
	require.NoError(t, err)
	addrInfo, err := peer.AddrInfoFromP2pAddr(multiAddress)
	require.NoError(t, err)
	err = h2.Connect(context.Background(), *addrInfo)
	require.NotNil(t, err, "Wanted connection to fail with max peer")
}

func TestService_InterceptBannedIP(t *testing.T) {
	s := &Service{
		ipLimiter: leakybucket.NewCollector(ipLimit, ipBurst, false),
		peers: peers.NewStatus(context.Background(), &peers.StatusConfig{
			PeerLimit:    20,
			ScorerParams: &scorers.Config{},
		}),
	}
	var err error
	s.addrFilter, err = configureFilter(&Config{})
	require.NoError(t, err)
	ip := "192.168.1.111"
	multiAddress, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", ip, 3000))
	require.NoError(t, err)

	for i := 0; i < ipBurst; i++ {
		valid := s.validateDial(multiAddress)
		if !valid {
			t.Errorf("Expected multiaddress with ip %s to not be rejected", ip)
		}
	}
	valid := s.validateDial(multiAddress)
	if valid {
		t.Errorf("Expected multiaddress with ip %s to be rejected as it exceeds the burst limit", ip)
	}
}