package twopc

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"gotest.tools/assert"
	"math"
	"sync/atomic"
	"testing"
	"time"
)

func TestProposalStateMonitor(t *testing.T) {

	now := time.Now()

	arr := []time.Time{
		now.Add(time.Duration(6) * time.Second),
		now.Add(time.Duration(1) * time.Second),
		now.Add(time.Duration(8) * time.Second),
		now.Add(time.Duration(4) * time.Second),
		now.Add(time.Duration(2) * time.Second),
		now.Add(time.Duration(3) * time.Second),
		now.Add(time.Duration(1) * time.Second),
	}

	consensus := &Twopc{
		state: &state{
			syncProposalStateMonitors: ctypes.NewSyncProposalStateMonitorQueue(0),
		},
	}

	queue := consensus.state.syncProposalStateMonitors

	t.Log("now time ", now.Format("2006-01-02 15:04:05"), "timestamp", now.UnixNano()/1e6)

	timer := consensus.proposalStateMonitorTimer()
	timer.Reset(time.Duration(math.MaxInt32) * time.Millisecond)

	ctx, cancelFn := context.WithCancel(context.Background())

	go func(cancelFn context.CancelFunc, queue *ctypes.SyncProposalStateMonitorQueue) {

		t.Log("Start handle 2pc consensus proposalState monitor queue")

		for {
			select {

			case <-timer.C:

				future := consensus.checkProposalStateMonitors(timeutils.UnixMsec())
				timer.Reset(time.Duration(future-timeutils.UnixMsec()) * time.Millisecond)

				if consensus.proposalStateMonitorsLen() == 0 {
					cancelFn()
					return
				}
			}
		}

	}(cancelFn, queue)

	var count uint32

	go func(queue *ctypes.SyncProposalStateMonitorQueue) {
		t.Log("Start add new one member into 2pc consensus proposalState monitor queue")
		for _, tm := range arr {
			queue.AddMonitor(ctypes.NewProposalStateMonitor(nil, tm.UnixNano()/1e6, tm.UnixNano()/1e6+1000,
				func(orgState *ctypes.OrgProposalState) {
					atomic.AddUint32(&count, 1)
				}))
		}
	}(queue)

	<-ctx.Done()
	assert.Equal(t, int(count), len(arr)*2, fmt.Sprintf("the number of monitors expected to be executed is %d, but the actual number is %d", len(arr)*2, count))
}
