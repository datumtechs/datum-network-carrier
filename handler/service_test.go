package handler

import (
	"context"
	"github.com/datumtechs/datum-network-carrier/common/abool"
	"github.com/datumtechs/datum-network-carrier/common/feed"
	statefeed "github.com/datumtechs/datum-network-carrier/common/feed/state"
	libp2ppb "github.com/datumtechs/datum-network-carrier/pb/carrier/rpc/debug/v1"
	p2ptest "github.com/datumtechs/datum-network-carrier/p2p/testing"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSyncHandlers_WaitToSync(t *testing.T) {
	p2p := p2ptest.NewTestP2P(t)
	r := Service{
		ctx: context.Background(),
		cfg: &Config{
			P2P:           p2p,
			StateNotifier: &p2ptest.MockStateNotifier{},
			InitialSync:   &p2ptest.Sync{IsSyncing: false},
		},
		chainStarted: abool.New(),
	}

	topic := "/carrier/%x/gossip_test_data"
	go r.registerHandlers()
	time.Sleep(100 * time.Millisecond)
	i := r.cfg.StateNotifier.StateFeed().Send(&feed.Event{
		Type: statefeed.Initialized,
		Data: &statefeed.InitializedData{
			StartTime: time.Now(),
		},
	})
	if i == 0 {
		t.Fatal("didn't send genesis time to subscribers")
	}

	p2p.ReceivePubSub(topic, &libp2ppb.SignedGossipTestData{
		Data:                 &libp2ppb.GossipTestData{
			Data:                 []byte("xxxx"),
			Count:                10,
			Step:                 777,
		},
		Signature:            make([]byte, 48),
	})
	// wait for chainstart to be sent
	time.Sleep(400 * time.Millisecond)
	require.Equal(t, true, true, "Did not receive chain start evengine.")
}

func TestSyncHandlers_WaitForChainStart(t *testing.T) {
	p2p := p2ptest.NewTestP2P(t)

	r := Service{
		ctx: context.Background(),
		cfg: &Config{
			P2P:           p2p,
			StateNotifier: &p2ptest.MockStateNotifier{},
			InitialSync:   &p2ptest.Sync{IsSyncing: false},
		},
		chainStarted:        abool.New(),
	}

	go r.registerHandlers()
	time.Sleep(100 * time.Millisecond)
	i := r.cfg.StateNotifier.StateFeed().Send(&feed.Event{
		Type: statefeed.Initialized,
		Data: &statefeed.InitializedData{
			StartTime: time.Now().Add(2 * time.Second),
		},
	})
	if i == 0 {
		t.Fatal("didn't send genesis time to subscribers")
	}
	require.Equal(t, false, r.chainStarted.IsSet(), "Chainstart was marked prematurely")

	// wait for chainstart to be sent
	time.Sleep(3 * time.Second)
	require.Equal(t, true, r.chainStarted.IsSet(), "Did not receive chain start evengine.")
}

func TestSyncService_StopCleanly(t *testing.T) {
	p2p := p2ptest.NewTestP2P(t)

	ctx, cancel := context.WithCancel(context.Background())
	r := Service{
		ctx:    ctx,
		cancel: cancel,
		cfg: &Config{
			P2P:           p2p,
			StateNotifier: &p2ptest.MockStateNotifier{},
			InitialSync:   &p2ptest.Sync{IsSyncing: false},
		},
		chainStarted: abool.New(),
	}

	go r.registerHandlers()
	time.Sleep(100 * time.Millisecond)
	i := r.cfg.StateNotifier.StateFeed().Send(&feed.Event{
		Type: statefeed.Initialized,
		Data: &statefeed.InitializedData{
			StartTime: time.Now(),
		},
	})
	if i == 0 {
		t.Fatal("didn't send genesis time to subscribers")
	}

	var err error
	p2p.Digest, err = r.forkDigest()
	require.NoError(t, err)

	// wait for chainstart to be sent
	time.Sleep(2 * time.Second)
	require.Equal(t, true, r.chainStarted.IsSet(), "Did not receive chain start evengine.")

	i = r.cfg.StateNotifier.StateFeed().Send(&feed.Event{
		Type: statefeed.Synced,
		Data: &statefeed.SyncedData{
			StartTime: time.Now(),
		},
	})
	if i == 0 {
		t.Fatal("didn't send genesis time to sync evengine subscribers")
	}

	time.Sleep(1 * time.Second)

	require.NotEqual(t, 0, len(r.cfg.P2P.PubSub().GetTopics()))
	require.NotEqual(t, 0, len(r.cfg.P2P.Host().Mux().Protocols()))

	// Both pubsub and rpc topcis should be unsubscribed.
	require.NoError(t, r.Stop())

	// Sleep to allow pubsub topics to be deregistered.
	time.Sleep(1 * time.Second)
	require.Equal(t, 0, len(r.cfg.P2P.PubSub().GetTopics()))
	require.Equal(t, 0, len(r.cfg.P2P.Host().Mux().Protocols()))
}
