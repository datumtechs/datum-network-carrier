package task

import (
	"context"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/types"
	teassert "github.com/stretchr/testify/assert"
	"gotest.tools/assert"
	"math"
	"sync/atomic"

	//"sync"
	"testing"
	"time"
)

func TestExecuteTaskMonitor(t *testing.T) {

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

	m, _ := NewTaskManager(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
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

				future := m.checkNeedExecuteTaskMonitors(timeutils.UnixMsec(), true)
				now := timeutils.UnixMsec()
				if future > now {
					timer.Reset(time.Duration(future-now) * time.Millisecond)
				} else if future < now {
					timer.Reset(time.Duration(now) * time.Millisecond)
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
			queue.AddMonitor(types.NewExecuteTaskMonitor("taskId:"+fmt.Sprint(tm.UnixNano()), "partyId:"+fmt.Sprint(tm.UnixNano()), tm.UnixNano()/1e6, func() {
				atomic.AddUint32(&count, 1)
			}))
		}
	}(queue)

	<-ctx.Done()
	assert.Equal(t, int(count), len(arr), fmt.Sprintf("the number of monitors expected to be executed is %d, but the actual number is %d", len(arr), count))
}

func TestCheckConsumeOptionsParams(t *testing.T) {
	{
		// correct
		tm := Manager{}
		task := types.NewTask(
			&carriertypespb.TaskPB{
				MetaAlgorithmId: "plaintext",
				DataPolicyTypes: []uint32{40001, 40001}, // csv
				DataPolicyOptions: []string{
					`{"partyId": "p1", "metadataId": "MetadataId001", "metadataName": "metadataName01", "inputType": 1, "keyColumn": 12, "selectedColumns": [1, 2, 3, 4], "consumeTypes": [2, 2, 3], "consumeOptions": ["{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"balance\": 10}", "{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"balance\": 2222}", "{\"contract\": \"0x79e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"takenId\": \"111222\"}"]}`,
					`{"partyId": "p2", "metadataId": "MetadataId002", "metadataName": "metadataName02", "inputType": 1, "keyColumn": 12, "selectedColumns": [1, 2, 3, 4], "consumeTypes": [2, 2, 3], "consumeOptions": ["{\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"balance\": 10}", "{\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"balance\": 2222}", "{\"contract\": \"0x87e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"takenId\": \"111222\"}"]}`,
				},
			},
		)
		err, dataConsumePolicy, _ := tm.checkConsumeOptionsParams(task, true)
		teassert.Nil(t, err)
		teassert.Equal(t, len(dataConsumePolicy), 2)
		//for consumeType, ConsumePolicyArray := range dataConsumePolicy {
		//	switch consumeType {
		//	case types.ConsumeMetadataAuth:
		//		//todo 等待实现
		//	case types.ConsumeTk20, types.ConsumeTk721:
		//		for _, consumePolicy := range ConsumePolicyArray {
		//			switch consumePolicy.(type) {
		//			case *types.Tk20Consume:
		//				fmt.Println((consumePolicy.(*types.Tk20Consume)).Address())
		//			case *types.Tk721Consume:
		//				fmt.Println((consumePolicy.(*types.Tk721Consume)).Address())
		//			}
		//		}
		//	}
		//}
	}
	{
		// test ciphertext balance less than cryptoAlgoConsumeUnit
		tm := Manager{}
		task := types.NewTask(
			&carriertypespb.TaskPB{
				MetaAlgorithmId: "ciphertext",
				DataPolicyTypes: []uint32{40001, 40001}, // csv
				DataPolicyOptions: []string{
					`{"partyId": "p1", "metadataId": "MetadataId001", "metadataName": "metadataName01", "inputType": 1, "keyColumn": 12, "selectedColumns": [1, 2, 3, 4], "consumeTypes": [2, 2, 3], "consumeOptions": ["{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"balance\": 2222}", "{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"balance\": 2222}", "{\"contract\": \"0x79e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"takenId\": \"111222\"}"]}`,
					`{"partyId": "p2", "metadataId": "MetadataId002", "metadataName": "metadataName02", "inputType": 1, "keyColumn": 12, "selectedColumns": [1, 2, 3, 4], "consumeTypes": [2, 2, 3], "consumeOptions": ["{\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"balance\": 2222}", "{\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"balance\": 2222}", "{\"contract\": \"0x87e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"takenId\": \"111222\"}"]}`,
				},
			},
		)
		err, _, _ := tm.checkConsumeOptionsParams(task, true)
		teassert.NotNil(t, err)
	}
	{
		// test plaintext balance less than cryptoAlgoConsumeUnit
		tm := Manager{}
		task := types.NewTask(
			&carriertypespb.TaskPB{
				MetaAlgorithmId: "plaintext",
				DataPolicyTypes: []uint32{40001, 40001}, // csv
				DataPolicyOptions: []string{
					`{"partyId": "p1", "metadataId": "MetadataId001", "metadataName": "metadataName01", "inputType": 1, "keyColumn": 12, "selectedColumns": [1, 2, 3, 4], "consumeTypes": [2, 2, 3], "consumeOptions": ["{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"balance\": 1}", "{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"balance\": 2222}", "{\"contract\": \"0x79e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"takenId\": \"111222\"}"]}`,
					`{"partyId": "p2", "metadataId": "MetadataId002", "metadataName": "metadataName02", "inputType": 1, "keyColumn": 12, "selectedColumns": [1, 2, 3, 4], "consumeTypes": [2, 2, 3], "consumeOptions": ["{\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"balance\": 1}", "{\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"balance\": 2222}", "{\"contract\": \"0x87e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"takenId\": \"111222\"}"]}`,
				},
			},
		)
		err, _, _ := tm.checkConsumeOptionsParams(task, true)
		teassert.NotNil(t, err)
	}
	{
		// test Check whether the contract address exists in the datacenter and the corresponding metadata
		// 0xb7e4b947F015f3f7C06E5173C2CfF41F2DDBAF04 not in dataCenter metadata
		tm := Manager{}
		task := types.NewTask(
			&carriertypespb.TaskPB{
				MetaAlgorithmId: "plaintext",
				DataPolicyTypes: []uint32{40001, 40001}, // csv
				DataPolicyOptions: []string{
					`{"partyId": "p1", "metadataId": "MetadataId001", "metadataName": "metadataName01", "inputType": 1, "keyColumn": 12, "selectedColumns": [1, 2, 3, 4], "consumeTypes": [2, 2, 3], "consumeOptions": ["{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"balance\": 9}", "{\"contract\": \"0xb7e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"balance\": 2222}", "{\"contract\": \"0x79e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"takenId\": \"111222\"}"]}`,
					`{"partyId": "p2", "metadataId": "MetadataId002", "metadataName": "metadataName02", "inputType": 1, "keyColumn": 12, "selectedColumns": [1, 2, 3, 4], "consumeTypes": [2, 2, 3], "consumeOptions": ["{\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"balance\": 9}", "{\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"balance\": 2222}", "{\"contract\": \"0x87e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"takenId\": \"111222\"}"]}`,
				},
			},
		)
		err, _, _ := tm.checkConsumeOptionsParams(task, true)
		teassert.NotNil(t, err)
	}
	{
		// test save metadataId include the contract address of the two same
		tm := Manager{}
		task := types.NewTask(
			&carriertypespb.TaskPB{
				MetaAlgorithmId: "plaintext",
				DataPolicyTypes: []uint32{40001, 40001}, // csv
				DataPolicyOptions: []string{
					`{"partyId": "p1", "metadataId": "MetadataId001", "metadataName": "metadataName01", "inputType": 1, "keyColumn": 12, "selectedColumns": [1, 2, 3, 4], "consumeTypes": [2, 2, 3], "consumeOptions": ["{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"balance\": 9}", "{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"balance\": 2222}", "{\"contract\": \"0x79e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"takenId\": \"111222\"}"]}`,
					`{"partyId": "p2", "metadataId": "MetadataId002", "metadataName": "metadataName02", "inputType": 1, "keyColumn": 12, "selectedColumns": [1, 2, 3, 4], "consumeTypes": [2, 2, 3], "consumeOptions": ["{\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"balance\": 9}", "{\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"balance\": 2222}", "{\"contract\": \"0x87e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"takenId\": \"111222\"}"]}`,
				},
			},
		)
		err, _, _ := tm.checkConsumeOptionsParams(task, true)
		teassert.NotNil(t, err)
	}
	{
		// test consumeTypes len not equal consumeOptions len
		tm := Manager{}
		task := types.NewTask(
			&carriertypespb.TaskPB{
				MetaAlgorithmId: "plaintext",
				DataPolicyTypes: []uint32{40001, 40001}, // csv
				DataPolicyOptions: []string{
					`{"partyId": "p1", "metadataId": "MetadataId001", "metadataName": "metadataName01", "inputType": 1, "keyColumn": 12, "selectedColumns": [1, 2, 3, 4], "consumeTypes": [2, 2], "consumeOptions": ["{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"balance\": 10}", "{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"balance\": 2222}", "{\"contract\": \"0x79e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"takenId\": \"111222\"}"]}`,
					`{"partyId": "p2", "metadataId": "MetadataId002", "metadataName": "metadataName02", "inputType": 1, "keyColumn": 12, "selectedColumns": [1, 2, 3, 4], "consumeTypes": [2, 2, 3], "consumeOptions": ["{\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"balance\": 10}", "{\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"balance\": 2222}", "{\"contract\": \"0x87e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"takenId\": \"111222\"}"]}`,
				},
			},
		)
		err, _, _ := tm.checkConsumeOptionsParams(task, true)
		teassert.NotNil(t, err)
	}
	{
		// isBeginConsume is false
		tm := Manager{}
		task := types.NewTask(
			&carriertypespb.TaskPB{
				MetaAlgorithmId: "plaintext",
				DataPolicyTypes: []uint32{40001, 40001}, // csv
				DataPolicyOptions: []string{
					`{"partyId": "p1", "metadataId": "MetadataId001", "metadataName": "metadataName01", "inputType": 1, "keyColumn": 12, "selectedColumns": [1, 2, 3, 4], "consumeTypes": [2, 2, 3], "consumeOptions": ["{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"balance\": 10}", "{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"balance\": 2222}", "{\"contract\": \"0x79e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"takenId\": \"111222\"}"]}`,
					`{"partyId": "p2", "metadataId": "MetadataId002", "metadataName": "metadataName02", "inputType": 1, "keyColumn": 12, "selectedColumns": [1, 2, 3, 4], "consumeTypes": [2, 2, 3], "consumeOptions": ["{\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"balance\": 10}", "{\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"balance\": 2222}", "{\"contract\": \"0x87e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"takenId\": \"111222\"}"]}`,
				},
			},
		)
		err, _, _ := tm.checkConsumeOptionsParams(task, false)
		teassert.Nil(t, err)
	}
}
