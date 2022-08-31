package task

import (
	"context"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	"github.com/datumtechs/datum-network-carrier/core"
	"github.com/datumtechs/datum-network-carrier/core/resource"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"

	//"github.com/datumtechs/datum-network-carrier/core"
	//"github.com/datumtechs/datum-network-carrier/core/resource"
	//carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/types"
	//teassert "github.com/stretchr/testify/assert"
	"gotest.tools/assert"
	"math"
	"sync/atomic"

	//"sync"
	"testing"
	"time"
)

func TestExecuteTaskMonitor(t *testing.T) {

	start := time.Now()

	arr := []time.Time{
		start.Add(time.Duration(6) * time.Second),
		start.Add(time.Duration(1) * time.Second),
		start.Add(time.Duration(8) * time.Second),
		time.Unix((start.UnixNano()/1e6-9999)/1000, 0), // pass + 9999 ms == now
		start.Add(time.Duration(4) * time.Second),
		start.Add(time.Duration(2) * time.Second),
		start.Add(time.Duration(3) * time.Second),
		time.Unix((start.UnixNano()/1e6-3000)/1000, 0), // pass + 3000 ms == now
		start.Add(time.Duration(1) * time.Second),
	}

	m, _ := NewTaskManager(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	queue := m.syncExecuteTaskMonitors

	t.Log("now time ", start.Format("2006-01-02 15:04:05"), "timestamp", start.UnixNano()/1e6)

	timer := m.syncExecuteTaskMonitors.Timer()
	timer.Reset(time.Duration(math.MaxInt32) * time.Millisecond)

	ctx, cancelFn := context.WithCancel(context.Background())

	go func(cancelFn context.CancelFunc, queue *types.SyncExecuteTaskMonitorQueue) {

		t.Log("Start handle executeTask monitor queue")

		for {
			select {

			case <-timer.C:

				future := m.checkNeedExecuteTaskMonitors(timeutils.UnixMsec(), true)
				now := timeutils.UnixMsec()
				if future > now {
					timer.Reset(time.Duration(future-now) * time.Millisecond)
					t.Logf("match future time %d now %d duration %d \n", future, now, future-now)
				} else if future < now {
					timer.Reset(time.Duration(0) * time.Millisecond)
					t.Logf("match pass time %d now %d duration %d \n", future, now, future-now)
				}
				// when future value is 0, we do nothing

				if m.syncExecuteTaskMonitors.Len() == 0 {
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
			taskId := "taskId:" + fmt.Sprint(tm.UnixNano())
			queue.AddMonitor(types.NewExecuteTaskMonitor(taskId, "partyId:"+fmt.Sprint(tm.UnixNano()), tm.UnixNano()/1e6, func() {
				atomic.AddUint32(&count, 1)
				//t.Logf("execute taskId %s", taskId)
			}))
		}
	}(queue)

	<-ctx.Done()
	assert.Equal(t, int(count), len(arr), fmt.Sprintf("the number of monitors expected to be executed is %d, but the actual number is %d", len(arr), count))
}

func TestFetchConsumeOption(t *testing.T) {
	tm := Manager{
		resourceMng: resource.NewResourceManager(core.MockDataCenter{}, nil, ""),
	}
	{
		localTask := types.NewTask(
			&carriertypespb.TaskPB{
				MetaAlgorithmId: "plaintext",
				DataPolicyTypes: []uint32{40001, 40001}, // csv
				DataPolicyOptions: []string{
					`{"partyId": "p1", "metadataId": "MetadataId001", "metadataName": "metadataName01", "inputType": 1, "keyColumn": 12, "selectedColumns": [1, 2, 3, 4], "consumeTypes": [1, 2, 2, 3], "consumeOptions": ["metadataAuth:0x1c131be4afaccc849608214a9f0594cdb991dbab0887d091b6ca7b460ac9a816","{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"balance\": 10}", "{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"balance\": 2222}", "{\"contract\": \"0x79e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"takenId\": \"111222\"}"]}`,
					`{"partyId": "p2", "metadataId": "MetadataId002", "metadataName": "metadataName02", "inputType": 1, "keyColumn": 12, "selectedColumns": [1, 2, 3, 4], "consumeTypes": [1, 2, 2, 3], "consumeOptions": ["metadataAuth:0xf968d3ecc764d144bbedf6b65759852cb4f8aa0a47b33acf89d483f6aa53028f","{\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"balance\": 10}", "{\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"balance\": 2222}", "{\"contract\": \"0x87e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"takenId\": \"111222\"}"]}`,
				},
			},
		)
		dataConsumeOptionsCache, err := tm.fetchConsumeOption(localTask)
		assert.NilError(t, err)
		for consumeType, consumeOptions := range dataConsumeOptionsCache {
			switch consumeType {
			case types.ConsumeMetadataAuth:
				var authIds = map[string]struct{}{"metadataAuth:0x1c131be4afaccc849608214a9f0594cdb991dbab0887d091b6ca7b460ac9a816": {}}
				for _, c := range consumeOptions {
					_, ok := authIds[c.GetDataConsumeOption().Address()]
					assert.Equal(t, true, ok)
				}
			case types.ConsumeTk20:
				var addressCaches = map[string]struct{}{"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF03": {}, "0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF03": {}, "0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF04": {}, "0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF04": {}, "metadataAuth:0xf968d3ecc764d144bbedf6b65759852cb4f8aa0a47b33acf89d483f6aa53028f": {}}
				for _, c := range consumeOptions {
					_, ok := addressCaches[c.GetDataConsumeOption().Address()]
					assert.Equal(t, true, ok)
				}
			case types.ConsumeTk721:
				var addressCaches = map[string]struct{}{"0x79e4b947F015f3f7C06E5173C2CfF41F2DDBAF07": {}, "0x87e4b947F015f3f7C06E5173C2CfF41F2DDBAF07": {}}
				for _, c := range consumeOptions {
					_, ok := addressCaches[c.GetDataConsumeOption().Address()]
					assert.Equal(t, true, ok)
				}
			}
		}
	}
}
