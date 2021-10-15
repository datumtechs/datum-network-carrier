package twopc

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/bytesutil"
	"github.com/RosettaFlow/Carrier-Go/common/rlputil"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/types"
	"gotest.tools/assert"
	"math/rand"
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
	count := 0
	partyIds := [5]string{"p1", "p2", "p3"}
	for {
		count += 1
		proposalId := generateProposalId()
		for _, value := range partyIds {
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
			db.UpdateOrgProposalState(proposalId, sender, orgState)
		}
		if count == 3 {
			break
		}
	}
}

func TestUpdateConfirmTaskPeerInfo(t *testing.T) {
	db := generateWalDB()
	proposalId := generateProposalId()
	peerDesc := &twopcpb.ConfirmTaskPeerInfo{
		OwnerPeerInfo: &twopcpb.TaskPeerInfo{
			Ip:      []byte("192.157.222.111"),
			Port:    []byte("8899"),
			PartyId: []byte("P2"),
		},
		DataSupplierPeerInfoList: []*twopcpb.TaskPeerInfo{
			{
				Ip:      []byte("192.157.222.112"),
				Port:    []byte("8890"),
				PartyId: []byte("P1"),
			},
		},
		PowerSupplierPeerInfoList: []*twopcpb.TaskPeerInfo{
			{
				Ip:      []byte("192.157.222.113"),
				Port:    []byte("8889"),
				PartyId: []byte("P0"),
			},
		},
		ResultReceiverPeerInfoList: []*twopcpb.TaskPeerInfo{
			{
				Ip:      []byte("192.157.222.114"),
				Port:    []byte("8888"),
				PartyId: []byte("P3"),
			},
		},
	}
	db.UpdateConfirmTaskPeerInfo(proposalId, peerDesc)
}
func TestUpdatePrepareVotes(t *testing.T) {
	db := generateWalDB()
	count := 0
	partyIds := [5]string{"p1", "p2", "p3"}
	for {
		proposalId := generateProposalId()
		for _, value := range partyIds {
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
			db.UpdatePrepareVotes(vote)
		}
		count += 1
		if count == 3 {
			break
		}
	}
}
func TestUpdateConfirmVotes(t *testing.T) {
	db := generateWalDB()
	count := 0
	partyIds := [5]string{"p1", "p2", "p3"}
	for {
		proposalId := generateProposalId()
		for _, value := range partyIds {
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
			db.UpdateConfirmVotes(vote)
		}
		count += 1
		if count == 3 {
			break
		}
	}
}
func TestDeleteState(t *testing.T) {
	db := generateWalDB()
	proposalId := "0x126e3fc23ace8c7351f2d7db7462ecc47812782509650f90851b49f99c064b79"
	db.DeleteState(db.GetProposalPeerInfoCacheKey(common.HexToHash(proposalId)))
}
func TestUnmarshal(t *testing.T) {
	db := generateWalDB()
	db.UnmarshalTest()
}
