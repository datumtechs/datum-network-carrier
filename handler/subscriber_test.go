package handler

import (
	"context"
	"github.com/datumtechs/datum-network-carrier/common/abool"
	libp2ppb "github.com/datumtechs/datum-network-carrier/pb/carrier/rpc/debug/v1"
	"github.com/datumtechs/datum-network-carrier/p2p"
	p2ptest "github.com/datumtechs/datum-network-carrier/p2p/testing"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestSubscribe_ReceivesValidMessage(t *testing.T) {
	p2pService := p2ptest.NewTestP2P(t)
	r := Service{
		ctx: context.Background(),
		cfg: &Config{
			P2P:         p2pService,
			InitialSync: &p2ptest.Sync{IsSyncing: false},
		},
	}
	var err error
	p2pService.Digest, err = r.forkDigest()
	require.NoError(t, err)
	topic := "/carrier/%x/gossip_test_data"
	var wg sync.WaitGroup
	wg.Add(1)

	r.subscribe(topic, r.noopValidator, func(_ context.Context, pid peer.ID, msg proto.Message) error {
		m, ok := msg.(*libp2ppb.GossipTestData)
		assert.Equal(t, true, ok, "Object is not of type *pb.GossipTestData")
		if m.Step == 0 || m.Step != 55 {
			t.Errorf("Unexpected incoming message: %+v", m)
		}
		wg.Done()
		return nil
	})
	//r.markForChainStart()

	p2pService.ReceivePubSub(topic, &libp2ppb.GossipTestData{Step: 55})

	if WaitTimeout(&wg, 10 * time.Second) {
		t.Fatal("Did not receive PubSub in 1 second")
	}
}

func TestSubscribe_HandlesPanic(t *testing.T) {
	p := p2ptest.NewTestP2P(t)
	r := Service{
		ctx: context.Background(),
		cfg: &Config{
			P2P: p,
		},
		chainStarted: abool.New(),
	}
	var err error
	p.Digest, err = r.forkDigest()
	require.NoError(t, err)

	topic := p2p.GossipTypeMapping[reflect.TypeOf(&libp2ppb.GossipTestData{})]
	var wg sync.WaitGroup
	wg.Add(1)

	r.subscribe(topic, r.noopValidator, func(_ context.Context, peerId peer.ID, msg proto.Message) error {
		defer wg.Done()
		panic("bad")
	})
	r.markForChainStart()
	p.ReceivePubSub(topic, &libp2ppb.GossipTestData{ Step: 55 })

	if WaitTimeout(&wg, time.Second) {
		t.Fatal("Did not receive PubSub in 1 second")
	}
}

// Create peer and register them to provided topics.
func createPeer(t *testing.T, topics ...string) *p2ptest.TestP2P {
	p := p2ptest.NewTestP2P(t)
	for _, tp := range topics {
		jTop, err := p.PubSub().Join(tp)
		if err != nil {
			t.Fatal(err)
		}
		_, err = jTop.Subscribe()
		if err != nil {
			t.Fatal(err)
		}
	}
	return p
}

