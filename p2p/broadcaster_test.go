package p2p

import (
	"context"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common/bytesutil"
	libp2ppb "github.com/datumtechs/datum-network-carrier/lib/rpc/debug/v1"
	libtypes "github.com/datumtechs/datum-network-carrier/lib/types"
	p2ptest "github.com/datumtechs/datum-network-carrier/p2p/testing"
	"github.com/gogo/protobuf/proto"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestService_Broadcast(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)
	p1.Connect(p2)
	if len(p1.BHost.Network().Peers()) == 0 {
		t.Fatal("No peers")
	}

	p := &Service{
		host:                  p1.BHost,
		pubsub:                p1.PubSub(),
		joinedTopics:          map[string]*pubsub.Topic{},
		cfg:                   &Config{},
		genesisTime:   	        time.Now(),
	}

	msg := &libp2ppb.GossipTestData{
		Data:                 []byte{0x0,0x1},
		Count:                11,
		Step:                 23,
	}

	topic := "/carrier/%x/testing"
	// Set a test gossip mapping for testpb.TestSimpleMessage.
	GossipTypeMapping[reflect.TypeOf(msg)] = topic
	digest, err := p.forkDigest()
	require.NoError(t, err)
	topic = fmt.Sprintf(topic, digest)

	// External peer subscribes to the topic.
	topic += p.Encoding().ProtocolSuffix()
	sub, err := p2.SubscribeToTopic(topic)
	require.NoError(t, err)

	time.Sleep(50 * time.Millisecond) // libp2p fails without this delay...

	// Async listen for the pubsub, must be before the broadcast.
	var wg sync.WaitGroup
	wg.Add(1)
	go func(tt *testing.T) {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		incomingMessage, err := sub.Next(ctx)
		require.NoError(t, err)

		result := &libp2ppb.GossipTestData{}
		require.NoError(t, p.Encoding().DecodeGossip(incomingMessage.Data, result))
		if !proto.Equal(result, msg) {
			tt.Errorf("Did not receive expected message, got %+v, wanted %+v", result, msg)
		}
	}(t)

	// Broadcast to peers and wait.
	require.NoError(t, p.Broadcast(context.Background(), msg))
	if WaitTimeout(&wg, 1*time.Second) {
		t.Error("Failed to receive pubsub within 1s")
	}
}


func TestService_Broadcast_ReturnsErr_TopicNotMapped(t *testing.T) {
	p := Service{
		genesisTime:           time.Now(),
		genesisValidatorsRoot: bytesutil.PadTo([]byte{'A'}, 32),
	}
	assert.ErrorContains(t, p.Broadcast(context.Background(), &libtypes.BlockData{}), ErrMessageNotMapped.Error())
}