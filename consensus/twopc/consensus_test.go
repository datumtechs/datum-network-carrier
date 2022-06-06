package twopc

import (
	"context"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	ctypes "github.com/datumtechs/datum-network-carrier/consensus/twopc/types"
	carriertwopcpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/consensus/twopc"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/types"
	"gotest.tools/assert"
	"math"
	"strconv"
	"sync"
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

				future := consensus.checkProposalStateMonitors(timeutils.UnixMsec(), true)
				now := timeutils.UnixMsec()
				if future > now {
					timer.Reset(time.Duration(future-now) * time.Millisecond)
				} else if future < now {
					timer.Reset(time.Duration(now) * time.Millisecond)
				}
				// when future value is 0, we do nothing

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
		for i, tm := range arr {
			orgState := ctypes.NewOrgProposalState(common.Hash{byte(uint8(i))},
				fmt.Sprintf("taskId:%d", i), commonconstantpb.TaskRole_TaskRole_Unknown,
				&carriertypespb.TaskOrganization{
					PartyId:    fmt.Sprintf("senderPartyId:%d", i),
					NodeName:   fmt.Sprintf("senderNodeName:%d", i),
					NodeId:     fmt.Sprintf("senderNodeId:%d", i),
					IdentityId: fmt.Sprintf("senderIdentityId:%d", i),
				},
				&carriertypespb.TaskOrganization{
					PartyId:    fmt.Sprintf("partyId:%d", i),
					NodeName:   fmt.Sprintf("nodeName:%d", i),
					NodeId:     fmt.Sprintf("nodeId:%d", i),
					IdentityId: fmt.Sprintf("identityId:%d", i),
				}, timeutils.UnixMsecUint64())
			queue.AddMonitor(ctypes.NewProposalStateMonitor(orgState, tm.UnixNano()/1e6, tm.UnixNano()/1e6+1000,
				func(orgState *ctypes.OrgProposalState) {
					atomic.AddUint32(&count, 1)
				}))
		}
	}(queue)

	<-ctx.Done()
	assert.Equal(t, int(count), len(arr)*2, fmt.Sprintf("the number of monitors expected to be executed is %d, but the actual number is %d", len(arr)*2, count))
}
func mockTestData(t *testing.T) *state {
	proposalIds := []common.Hash{
		common.HexToHash("0x35af63cf9e8f90dcc8a8e024dc78acbb268caa711a8a7339a9492f5ef2f8833e"),
		common.HexToHash("0x3ff6fea93531aa400789b3f8ea0d28409790499fb16391f0814e033eef1a2ccf"),
	}
	proposalTaskCache := make(map[string]map[string]*ctypes.ProposalTask, 0)
	for i := 0; i < 2; i++ {
		cache := make(map[string]*ctypes.ProposalTask, 0)
		taskId := fmt.Sprintf("%s,%d", "task_00", i)
		for p := 0; p < 3; p++ {
			partyId := fmt.Sprintf("%s,%d", "p", i)
			cache[partyId] = &ctypes.ProposalTask{
				ProposalId: proposalIds[i],
				TaskId:     taskId,
				CreateAt:   timeutils.UnixMsecUint64(),
			}
		}
		proposalTaskCache[taskId] = cache
	}

	proposalSet := make(map[common.Hash]map[string]*ctypes.OrgProposalState, 0)
	orgProposalState := make(map[string]*ctypes.OrgProposalState, 0)
	for i := 0; i < 2; i++ {
		taskId := fmt.Sprintf("%s,%d", "task_00", i)
		for p := 0; p < 3; p++ {
			partyId := fmt.Sprintf("%s,%d", "p", i)
			taskSender := &carriertypespb.TaskOrganization{
				PartyId:    partyId,
				NodeName:   fmt.Sprintf("%s,%d", "NodeName_", i),
				NodeId:     "",
				IdentityId: fmt.Sprintf("%s,%d", "IdentityId_", i),
			}
			orgProposalState[partyId] = ctypes.NewOrgProposalState(proposalIds[i], taskId, 2, taskSender, nil, 7443)
			proposalSet[proposalIds[i]] = orgProposalState
		}
	}

	yesVotes := make(map[commonconstantpb.TaskRole]uint32, 0)
	voteStatus := make(map[commonconstantpb.TaskRole]uint32, 0)
	for i := 0; i < 5; i++ {
		yesVotes[commonconstantpb.TaskRole(i)] = uint32(i + 1)
		voteStatus[commonconstantpb.TaskRole(i)] = uint32(i + 2)
	}

	prepareVotes := make(map[common.Hash]*prepareVoteState, 0)
	votesP := make(map[string]*types.PrepareVote, 0)
	for i := 0; i < 2; i++ {
		for p := 0; p < 3; p++ {
			partyId := fmt.Sprintf("%s,%d", "p", i)
			votesP[partyId] = &types.PrepareVote{
				MsgOption: &types.MsgOption{
					ProposalId:      proposalIds[i],
					SenderRole:      commonconstantpb.TaskRole(12),
					SenderPartyId:   partyId,
					ReceiverRole:    commonconstantpb.TaskRole(23),
					ReceiverPartyId: partyId,
					Owner: &carriertypespb.TaskOrganization{
						PartyId:    partyId,
						NodeName:   "NodeName_" + strconv.Itoa(i),
						NodeId:     "",
						IdentityId: "IdentityId_" + strconv.Itoa(i),
					},
				},
				VoteOption: 12,
				CreateAt:   7777,
				Sign:       []byte("TestSignPrepareVote"),
			}
			prepareVotes[proposalIds[i]] = &prepareVoteState{
				votes: votesP,
			}
		}
	}

	confirmVotes := make(map[common.Hash]*confirmVoteState, 0)
	votesC := make(map[string]*types.ConfirmVote, 0)
	for i := 0; i < 2; i++ {
		for p := 0; p < 3; p++ {
			partyId := fmt.Sprintf("%s,%d", "p", i)
			votesC[partyId] = &types.ConfirmVote{
				MsgOption: &types.MsgOption{
					ProposalId:      proposalIds[i],
					SenderRole:      commonconstantpb.TaskRole(12),
					SenderPartyId:   partyId,
					ReceiverRole:    commonconstantpb.TaskRole(23),
					ReceiverPartyId: partyId,
					Owner: &carriertypespb.TaskOrganization{
						PartyId:    partyId,
						NodeName:   "NodeName_" + strconv.Itoa(i),
						NodeId:     "",
						IdentityId: "IdentityId_" + strconv.Itoa(i),
					},
				},
				VoteOption: 12,
				CreateAt:   7777,
				Sign:       []byte("TestSignConfirmVote"),
			}
			confirmVotes[proposalIds[i]] = &confirmVoteState{
				votes: votesC,
			}
		}
	}
	cache, err := NewTwopcMsgCache(default2pcMsgCacheSize)
	if nil != err {
		t.Fatalf("cannot make twopcMsgCache, %s", err)
	}

	return &state{
		proposalTaskCache: proposalTaskCache,
		proposalSet:       proposalSet,
		prepareVotes:      prepareVotes,
		confirmVotes:      confirmVotes,
		msgCache:          cache,
	}
}
func Testtwopc_Get2PcProposalStateByTaskId(t *testing.T) {
	twopc := &Twopc{
		state: mockTestData(t),
	}
	result1, _ := twopc.Get2PcProposalStateByTaskId("task_00,0")
	assert.Equal(t, "0x35af63cf9e8f90dcc8a8e024dc78acbb268caa711a8a7339a9492f5ef2f8833e", result1.ProposalId)
	result2, _ := twopc.Get2PcProposalStateByTaskId("task_00,1")
	assert.Equal(t, "0x3ff6fea93531aa400789b3f8ea0d28409790499fb16391f0814e033eef1a2ccf", result2.ProposalId)
}
func Testtwopc_Get2PcProposalStateByProposalId(t *testing.T) {
	twopc := &Twopc{
		state: mockTestData(t),
	}
	result1, _ := twopc.Get2PcProposalStateByProposalId("0x35af63cf9e8f90dcc8a8e024dc78acbb268caa711a8a7339a9492f5ef2f8833e")
	assert.Equal(t, "task_00,0", result1.State["p,0"].TaskId)
	result2, _ := twopc.Get2PcProposalStateByProposalId("0x3ff6fea93531aa400789b3f8ea0d28409790499fb16391f0814e033eef1a2ccf")
	assert.Equal(t, "task_00,0", result2.State["p,0"].TaskId)
}
func Testtwopc_Get2PcProposalPrepare(t *testing.T) {
	twopc := &Twopc{
		state: mockTestData(t),
	}
	_, err := twopc.Get2PcProposalPrepare("0x35af63cf9e8f90dcc8a8e024dc78acbb268caa711a8a7339a9492f5ef2f8833e")
	assert.NilError(t, err)
}
func Testtwopc_Get2PcProposalConfirm(t *testing.T) {
	twopc := &Twopc{
		state: mockTestData(t),
	}
	_, err := twopc.Get2PcProposalConfirm("0x35af63cf9e8f90dcc8a8e024dc78acbb268caa711a8a7339a9492f5ef2f8833e")
	assert.NilError(t, err)
}

func TestNewtwopcMsgCacheAddMsg(t *testing.T) {
	state := mockTestData(t)
	assert.Equal(t, false, state.AddMsg(12), "expect return false, but true")
	assert.Equal(t, true, state.AddMsg(&carriertwopcpb.PrepareMsg{}), "expect return true, but false")
}

func TestNewtwopcMsgCacheContainsOrAddMsg(t *testing.T) {
	state := mockTestData(t)
	state.ContainsOrAddMsg(&carriertwopcpb.PrepareMsg{})
	if err := state.ContainsOrAddMsg(&carriertwopcpb.PrepareMsg{}); nil == err {
		t.Fatalf("expect return err, but nil")
	}
	err := state.ContainsOrAddMsg(&carriertwopcpb.ConfirmMsg{})
	assert.NilError(t, err, fmt.Sprintf("expect return nil, but err: %s", err))
}

func TestConsensusasyncCallCh(t *testing.T) {

	twopc := &Twopc{
		asyncCallCh: make(chan func(), 1),
	}

	var (
		errOne = fmt.Errorf("I am err one")
		errTwo = fmt.Errorf("I am err two")
		count = uint32(0)
		wg sync.WaitGroup
	)

	wg.Add(3)

	ctx, cancal := context.WithCancel(context.Background())

	trigger := func(cancel context.CancelFunc, count *uint32) {
		atomic.AddUint32(count, 1)
		if c := atomic.LoadUint32(count); c >= 2 {
			cancel()
		}
	}

	go func(cancel context.CancelFunc, count *uint32) {

		defer wg.Done()

		f := func() error {

			errCh := make(chan error, 1)

			twopc.asyncCallCh <- func() {

				defer close(errCh)

				errCh <- errOne
				trigger(cancal, count)
				t.Log("finished func1")

			}

			return <-errCh
		}

		err := f()
		assert.Equal(t, errOne, err, "expect: %s, but: %s", errOne, err)

	}(cancal, &count)

	go func(cancel context.CancelFunc, count *uint32) {

		defer wg.Done()

		f := func() error {

			errCh := make(chan error, 1)

			twopc.asyncCallCh <- func() {

				defer close(errCh)

				errCh <- errTwo
				trigger(cancal, count)
				t.Log("finished func2")

			}

			return <-errCh
		}

		err := f()
		assert.Equal(t, errTwo, err, "expect: %s, but: %s", errTwo, err)

	}(cancal, &count)


	go func() {

		defer wg.Done()

		for {
			select {
			case fn := <- twopc.asyncCallCh:
				fn()
			case <-ctx.Done():
				t.Log("count: ", count)
				return
			}
		}

	}()

	wg.Wait()

}

func TestContex (t *testing.T)  {

	childCtx := context.WithValue(context.Background(), "key1", "value1")

	v := childCtx.Value("key1").(string)
	fmt.Println("first fetch value:", v)

	childCtx = context.WithValue(childCtx, "key1", "value2")
	childCtx = context.WithValue(childCtx, "key1", "value3")
	childCtx = context.WithValue(childCtx, "key1", "value4")

	v = childCtx.Value("key1").(string)
	fmt.Println("second fetch value:", v)

}