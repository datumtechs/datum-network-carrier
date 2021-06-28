package p2p

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/bytesutil"
	"github.com/RosettaFlow/Carrier-Go/common/feed"
	statefeed "github.com/RosettaFlow/Carrier-Go/common/feed/state"
	timeutils "github.com/RosettaFlow/Carrier-Go/common/timeutil"
	"github.com/RosettaFlow/Carrier-Go/event"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	noise "github.com/libp2p/go-libp2p-noise"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
	"time"
)


type mockListener struct {
	localNode *enode.LocalNode
}

func (m mockListener) Self() *enode.Node {
	return m.localNode.Node()
}

func (mockListener) Close() {
	// no-op
}

func (mockListener) Lookup(enode.ID) []*enode.Node {
	panic("implement me")
}

func (mockListener) ReadRandomNodes(_ []*enode.Node) int {
	panic("implement me")
}

func (mockListener) Resolve(*enode.Node) *enode.Node {
	panic("implement me")
}

func (mockListener) Ping(*enode.Node) error {
	panic("implement me")
}

func (mockListener) RequestENR(*enode.Node) (*enode.Node, error) {
	panic("implement me")
}

func (mockListener) LocalNode() *enode.LocalNode {
	panic("implement me")
}

func (mockListener) RandomNodes() enode.Iterator {
	panic("implement me")
}

// initializeStateWithForkDigest sets up the state feed initialized event and returns the fork
// digest associated with that genesis event.
func initializeStateWithForkDigest(ctx context.Context, t *testing.T, ef *event.Feed) [4]byte {
	gt := timeutils.Now()
	gvr := bytesutil.PadTo([]byte("genesis validator root"), 32)
	for n := 0; n == 0; {
		if ctx.Err() != nil {
			t.Fatal(ctx.Err())
		}
		n = ef.Send(&feed.Event{
			Type: statefeed.Initialized,
			Data: &statefeed.InitializedData{
				StartTime:             gt,
				GenesisValidatorsRoot: gvr,
			},
		})
	}

	fd := [4]byte{0x1, 0x1, 0x1, 0x1,}
	require.NoError(t, nil)

	time.Sleep(50 * time.Millisecond) // wait for pubsub filter to initialize.

	return fd
}

func createHost(t *testing.T, port int) (host.Host, *ecdsa.PrivateKey, net.IP) {
	_, pkey := createAddrAndPrivKey(t)
	ipAddr := net.ParseIP("127.0.0.1")
	listen, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", ipAddr, port))
	require.NoError(t, err, "Failed to p2p listen")
	h, err := libp2p.New(context.Background(), []libp2p.Option{privKeyOption(pkey), libp2p.ListenAddrs(listen), libp2p.Security(noise.ID, noise.New)}...)
	require.NoError(t, err)
	return h, pkey, ipAddr
}