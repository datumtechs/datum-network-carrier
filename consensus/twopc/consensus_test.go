package twopc

import (
	"context"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	ctypes "github.com/datumtechs/datum-network-carrier/consensus/twopc/types"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/types"
	"gotest.tools/assert"
	"math"
	"strconv"
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
				fmt.Sprintf("taskId:%d", i), carriertypespb.TaskRole_TaskRole_Unknown,
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
func mockTestData() *state {
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

	yesVotes := make(map[carriertypespb.TaskRole]uint32, 0)
	voteStatus := make(map[carriertypespb.TaskRole]uint32, 0)
	for i := 0; i < 5; i++ {
		yesVotes[carriertypespb.TaskRole(i)] = uint32(i + 1)
		voteStatus[carriertypespb.TaskRole(i)] = uint32(i + 2)
	}

	prepareVotes := make(map[common.Hash]*prepareVoteState, 0)
	votesP := make(map[string]*types.PrepareVote, 0)
	for i := 0; i < 2; i++ {
		for p := 0; p < 3; p++ {
			partyId := fmt.Sprintf("%s,%d", "p", i)
			votesP[partyId] = &types.PrepareVote{
				MsgOption: &types.MsgOption{
					ProposalId:      proposalIds[i],
					SenderRole:      carriertypespb.TaskRole(12),
					SenderPartyId:   partyId,
					ReceiverRole:    carriertypespb.TaskRole(23),
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
					SenderRole:      carriertypespb.TaskRole(12),
					SenderPartyId:   partyId,
					ReceiverRole:    carriertypespb.TaskRole(23),
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
	return &state{
		proposalTaskCache: proposalTaskCache,
		proposalSet:       proposalSet,
		prepareVotes:      prepareVotes,
		confirmVotes:      confirmVotes,
	}
}
func TestTwopc_Get2PcProposalStateByTaskId(t *testing.T) {
	twoPc := &Twopc{
		state: mockTestData(),
	}
	result1, _ := twoPc.Get2PcProposalStateByTaskId("task_00,0")
	assert.Equal(t, "0x35af63cf9e8f90dcc8a8e024dc78acbb268caa711a8a7339a9492f5ef2f8833e", result1.ProposalId)
	result2, _ := twoPc.Get2PcProposalStateByTaskId("task_00,1")
	assert.Equal(t, "0x3ff6fea93531aa400789b3f8ea0d28409790499fb16391f0814e033eef1a2ccf", result2.ProposalId)
}
func TestTwopc_Get2PcProposalStateByProposalId(t *testing.T) {
	twoPc := &Twopc{
		state: mockTestData(),
	}
	result1, _ := twoPc.Get2PcProposalStateByProposalId("0x35af63cf9e8f90dcc8a8e024dc78acbb268caa711a8a7339a9492f5ef2f8833e")
	assert.Equal(t, "task_00,0", result1.State["p,0"].TaskId)
	result2, _ := twoPc.Get2PcProposalStateByProposalId("0x3ff6fea93531aa400789b3f8ea0d28409790499fb16391f0814e033eef1a2ccf")
	assert.Equal(t, "task_00,0", result2.State["p,0"].TaskId)
}
func TestTwopc_Get2PcProposalPrepare(t *testing.T) {
	twoPc := &Twopc{
		state: mockTestData(),
	}
	_, err := twoPc.Get2PcProposalPrepare("0x35af63cf9e8f90dcc8a8e024dc78acbb268caa711a8a7339a9492f5ef2f8833e")
	assert.NilError(t, err)
}
func TestTwopc_Get2PcProposalConfirm(t *testing.T) {
	twoPc := &Twopc{
		state: mockTestData(),
	}
	_, err := twoPc.Get2PcProposalConfirm("0x35af63cf9e8f90dcc8a8e024dc78acbb268caa711a8a7339a9492f5ef2f8833e")
	assert.NilError(t, err)
}
