package twopc

import (
	"bytes"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/bytesutil"
	"github.com/RosettaFlow/Carrier-Go/common/rlputil"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/gogo/protobuf/proto"
	"gotest.tools/assert"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

func RandStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func generateProposalId() common.Hash {
	now := timeutils.UnixMsecUint64()

	var buf bytes.Buffer
	buf.Write([]byte(RandStr(32)))
	buf.Write(bytesutil.Uint64ToBytes(now))
	proposalId := rlputil.RlpHash(buf.Bytes())
	return proposalId
}

func generateWalDB() *walDB {
	config := &Config{
		PeerMsgQueueSize:   12,
		ConsensusStateFile: "D:\\project\\src\\github.com\\RosettaFlow\\test.json",
	}
	return newWal(config)
}

func TestKeySplitProposalIdAndPartyId(t *testing.T) {
	db := generateWalDB()
	prefixLength := len([]byte("proposalSet:"))
	beforeProposalId := generateProposalId()
	result := db.GetProposalSetKey(beforeProposalId, "p9")
	prefix := string(result[:prefixLength])
	partyId := string(result[prefixLength+32:])
	proposalId := common.BytesToHash(result[prefixLength : prefixLength+32])
	assert.Equal(t, "p9", partyId)
	assert.Equal(t, beforeProposalId, proposalId)
	assert.Equal(t, "proposalSet:", prefix)
}

func TestUpdateOrgProposalState(t *testing.T) {
	db := generateWalDB()
	partyIds := []string{"p1", "p2", "p3"}
	for _, value := range partyIds {
		proposalId := generateProposalId()
		sender := &apicommonpb.TaskOrganization{
			PartyId:    value + "",
			NodeName:   value + "NodeName",
			NodeId:     value + "NodeId",
			IdentityId: value + "IdentityId",
		}
		orgState := &ctypes.OrgProposalState{
			PrePeriodStartTime: 1111,
			PeriodStartTime:    2222,
			DeadlineDuration:   3333,
			CreateAt:           4444,
			TaskId:             "TASK001",
			TaskRole:           3,
			TaskOrg: &apicommonpb.TaskOrganization{
				PartyId:    value + "P2",
				NodeName:   value + "NodeName",
				NodeId:     value + "NodeId",
				IdentityId: value + "IdentityId",
			},
			PeriodNum: 2,
		}
		db.StoreOrgProposalState(proposalId, sender, orgState)
	}
}

func TestUpdateConfirmTaskPeerInfo(t *testing.T) {
	db := generateWalDB()
	proposalId := generateProposalId()
	peerDesc := &twopcpb.ConfirmTaskPeerInfo{
		DataSupplierPeerInfos: []*twopcpb.TaskPeerInfo{
			{
				Ip:      []byte("192.157.222.112"),
				Port:    []byte("8890"),
				PartyId: []byte("P1"),
			},
		},
		PowerSupplierPeerInfos: []*twopcpb.TaskPeerInfo{
			{
				Ip:      []byte("192.157.222.113"),
				Port:    []byte("8889"),
				PartyId: []byte("P0"),
			},
		},
		ResultReceiverPeerInfos: []*twopcpb.TaskPeerInfo{
			{
				Ip:      []byte("192.157.222.114"),
				Port:    []byte("8888"),
				PartyId: []byte("P3"),
			},
		},
	}
	db.StoreConfirmTaskPeerInfo(proposalId, peerDesc)
}
func TestUpdatePrepareVotes(t *testing.T) {
	db := generateWalDB()
	partyIds := []string{"p1", "p2", "p3"}
	for _, value := range partyIds {
		proposalId := generateProposalId()
		vote := &types.PrepareVote{
			MsgOption: &types.MsgOption{
				ProposalId:      proposalId,
				SenderRole:      2,
				SenderPartyId:   value + "",
				ReceiverRole:    2,
				ReceiverPartyId: "P2",
				Owner: &apicommonpb.TaskOrganization{
					PartyId:    "P1",
					NodeName:   "NodeName",
					NodeId:     "NodeId",
					IdentityId: "IdentityId",
				},
			},
			VoteOption: 1,
			PeerInfo: &types.PrepareVoteResource{
				Id:      "ID",
				Ip:      "192.22.222.211",
				Port:    "9988",
				PartyId: "P8",
			},
			CreateAt: 2121,
			Sign:     []byte("this is a test"),
		}
		db.StorePrepareVote(vote)
	}
}
func TestUpdateConfirmVotes(t *testing.T) {
	db := generateWalDB()
	partyIds := []string{"p1", "p2", "p3"}
	for _, value := range partyIds {
		proposalId := generateProposalId()
		vote := &types.ConfirmVote{
			MsgOption: &types.MsgOption{
				ProposalId:      proposalId,
				SenderRole:      2,
				SenderPartyId:   value + "",
				ReceiverRole:    2,
				ReceiverPartyId: "P2",
				Owner: &apicommonpb.TaskOrganization{
					PartyId:    "P1",
					NodeName:   "NodeName",
					NodeId:     "NodeId",
					IdentityId: "IdentityId",
				},
			},
			VoteOption: 1,
			CreateAt:   2121,
			Sign:       []byte("this is a test"),
		}
		db.StoreConfirmVote(vote)
	}
}
func TestStoreProposalTask(t *testing.T) {
	db := generateWalDB()
	partyIds := []string{"p1", "p2", "p3"}
	taskIds := []string{"task:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e65",
		"task:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e66",
		"task:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e67"}
	for index, _ := range partyIds {
		proposalTask := &types.ProposalTask{
			ProposalId: generateProposalId(),
			TaskId: taskIds[index],
		}
		db.StoreProposalTask(partyIds[index], proposalTask)
	}
}
func TestDeleteState(t *testing.T) {
	db := generateWalDB()
	proposalId := "0x126e3fc23ace8c7351f2d7db7462ecc47812782509650f90851b49f99c064b79"
	db.DeleteState(db.GetProposalPeerInfoCacheKey(common.HexToHash(proposalId)))
}
func TestRecoveryState(t *testing.T) {
	db := generateWalDB()

	errCh := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	// recovery proposalSet (proposalId -> partyId -> orgState),StoreOrgProposalState
	go func(wg *sync.WaitGroup, errCh chan<- error) {

		defer wg.Done()

		prefixLength := len(proposalSetPrefix)
		proposalSet := make(map[common.Hash]*ctypes.ProposalState, 0)
		if err := db.ForEachKVWithPrefix(proposalSetPrefix, func(key, value []byte) error {

			if len(key) != 0 && len(value) != 0 {
				proposalId := common.BytesToHash(key[prefixLength : prefixLength+32])

				libOrgProposalState := &libtypes.OrgProposalState{}
				if err := proto.Unmarshal(value, libOrgProposalState); err != nil {
					return fmt.Errorf("unmarshal org proposalState failed, %s", err)
				}
				//proposalState, ok := t.state.proposalSet[proposalId]
				proposalState, ok := proposalSet[proposalId]
				if !ok {
					proposalState = ctypes.NewProposalState(proposalId, libOrgProposalState.GetTaskId(), libOrgProposalState.GetTaskSender())
				}
				proposalState.StoreOrgProposalStateUnSafe(&ctypes.OrgProposalState{
					PrePeriodStartTime: libOrgProposalState.GetPrePeriodStartTime(),
					PeriodStartTime:    libOrgProposalState.GetPeriodStartTime(),
					DeadlineDuration:   libOrgProposalState.GetDeadlineDuration(),
					CreateAt:           libOrgProposalState.GetCreateAt(),
					TaskId:             libOrgProposalState.GetTaskId(),
					TaskRole:           libOrgProposalState.GetTaskRole(),
					TaskOrg:            libOrgProposalState.GetTaskOrg(),
					PeriodNum:          ctypes.ProposalStatePeriod(libOrgProposalState.GetPeriodNum()),
				})
				//t.state.proposalSet[proposalId] = proposalState
				proposalSet[proposalId] = proposalState
			}
			return nil
		}); nil != err {
			errCh <- err
			return
		}
	}(&wg, errCh)

	// recovery proposalPeerInfoCache (proposalId -> ConfirmTaskPeerInfo)
	go func(wg *sync.WaitGroup, errCh chan<- error) {
		proposalPeerInfoCache := make(map[common.Hash]*twopcpb.ConfirmTaskPeerInfo, 0)
		defer wg.Done()

		prefixLength := len(proposalPeerInfoCachePrefix)
		if err := db.ForEachKVWithPrefix(proposalPeerInfoCachePrefix, func(key, value []byte) error {

			if len(key) != 0 && len(value) != 0 {
				proposalId := common.BytesToHash(key[prefixLength:])
				confirmTaskPeerInfo := &twopcpb.ConfirmTaskPeerInfo{}
				if err := proto.Unmarshal(value, confirmTaskPeerInfo); err != nil {
					return fmt.Errorf("unmarshal confirmTaskPeerInfo failed, %s", err)
				}
				proposalPeerInfoCache[proposalId] = confirmTaskPeerInfo
			}
			return nil
		}); nil != err {
			errCh <- err
			return
		}
	}(&wg, errCh)

	wg.Wait()
	close(errCh)

	errStrs := make([]string, 0)

	for err := range errCh {
		if nil != err {
			errStrs = append(errStrs, err.Error())
		}
	}
	if len(errStrs) != 0 {
		log.Fatalf(
			"recover consensus state failed: \n%s",
			strings.Join(errStrs, "\n"))
	}
}
func TestUnmarshal(t *testing.T) {
	db := generateWalDB()
	db.UnmarshalTest()
}
