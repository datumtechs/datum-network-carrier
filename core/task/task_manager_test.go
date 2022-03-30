package task

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/types"
	"gotest.tools/assert"
	"math"
	"sync/atomic"

	//"sync"
	"testing"
	"time"
)

func TestExecuteTaskMonitor (t *testing.T) {

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


	m  := NewTaskManager(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	queue := m.syncExecuteTaskMonitors

	t.Log("now time ", now.Format("2006-01-02 15:04:05"), "timestamp", now.UnixNano()/1e6)

	timer := m.syncExecuteTaskMonitors.Timer()
	timer.Reset(time.Duration(math.MaxInt32) * time.Millisecond)

	ctx, cancelFn := context.WithCancel(context.Background())

	go func(cancelFn context.CancelFunc, queue *types.SyncExecuteTaskMonitorQueue) {

		t.Log("Start handle executeTask monitor queue")

		for {
			select {

			case <-timer.C:

				future := m.checkNeedExecuteTaskMonitors(timeutils.UnixMsec())
				timer.Reset(time.Duration(future-timeutils.UnixMsec()) * time.Millisecond)

				if  m.syncExecuteTaskMonitors.Len() == 0 {
					cancelFn()
					return
				}
			}
		}

	}(cancelFn, queue)

	var count uint32

	go func(queue *types.SyncExecuteTaskMonitorQueue) {
		t.Log("Start add new one member into executeTask monitor queue")
		for _, tm := range arr {
			queue.AddMonitor(types.NewExecuteTaskMonitor("taskId:"+fmt.Sprint(tm.UnixNano()), "partyId:" + fmt.Sprint(tm.UnixNano()), tm.UnixNano()/1e6, func() {
				atomic.AddUint32(&count, 1)
			}))
		}
	}(queue)

	<-ctx.Done()
	assert.Equal(t, int(count), len(arr), fmt.Sprintf("the number of monitors expected to be executed is %d, but the actual number is %d", len(arr), count))
}
