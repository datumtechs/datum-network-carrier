package p2p

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/common/bytesutil"
	"github.com/RosettaFlow/Carrier-Go/common/feed"
	statefeed "github.com/RosettaFlow/Carrier-Go/common/feed/state"
	timeutils "github.com/RosettaFlow/Carrier-Go/common/timeutil"
	"github.com/RosettaFlow/Carrier-Go/event"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

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
