package rawdb

import (
	"bytes"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/bytesutil"
	"github.com/RosettaFlow/Carrier-Go/common/rlputil"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/db"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"gotest.tools/assert"
	"math/rand"
	"testing"
	"time"
)

func TestResourceDataUsed(t *testing.T) {
	//database := db.NewMemoryDatabase()
	//dataUsed := types.NewDataResourceFileUpload("node_id", "origin_id_01", "metadata_id", "/a/b/c/d")

	// save
	//err := StoreDataResourceDataUsed(database, dataUsed)
	//require.Nil(t, err)
	//
	//// query
	//queryUsed, err := QueryDataResourceDataUsed(database, dataUsed.GetOriginId())
	//require.Nil(t, err)
	//assert.Equal(t, dataUsed.GetOriginId(), queryUsed.GetOriginId())
}

var (
	taskIds = []string{"task:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e65",
		"task:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e66",
		"task:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e67",
		"task:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e68",
		"task:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e69"}
	partyIds = []string{"P0", "P1", "P2"}
)

func generateProposalId() common.Hash {
	now := timeutils.UnixMsecUint64()

	var buf bytes.Buffer
	buf.Write([]byte(RandStr(32)))
	buf.Write(bytesutil.Uint64ToBytes(now))
	proposalId := rlputil.RlpHash(buf.Bytes())
	return proposalId
}
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
func NeedExecuteTask() KeyValueStore {
	database := db.NewMemoryDatabase()
	for _, taskId := range taskIds {
		for _, partyId := range partyIds {
			remotepid := "remotepid"
			proposalId := generateProposalId()
			localTaskOrganization := &apicommonpb.TaskOrganization{
				PartyId:    partyId,
				NodeName:   "NodeName",
				NodeId:     "NodeId_0001",
				IdentityId: "IdentityId_0001",
			}
			remoteTaskOrganization := &apicommonpb.TaskOrganization{
				PartyId:    partyId,
				NodeName:   "NodeName",
				NodeId:     "NodeId_0002",
				IdentityId: "IdentityId_0002",
			}
			task := &types.Task{}
			localResource := &types.PrepareVoteResource{
				Id:      "PrepareVoteResourceId",
				Ip:      "2.2.2.2",
				Port:    "5555",
				PartyId: partyId,
			}
			resources := &twopcpb.ConfirmTaskPeerInfo{}

			_task := types.NewNeedExecuteTask(peer.ID(remotepid), proposalId, 1, localTaskOrganization, 2, remoteTaskOrganization, task,
				3, localResource, resources)
			if err := StoreNeedExecuteTask(database, _task, taskId, partyId); err != nil {
				fmt.Printf("StoreNeedExecuteTask fail,taskId %s\n", taskId)
			}
		}
	}
	return database
}
func TestStoreNeedExecuteTask(t *testing.T) {
	NeedExecuteTask()
}

func TestDeleteNeedExecuteTask(t *testing.T) {
	database := NeedExecuteTask()
	taskId1 := "task:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e65"
	taskId2 := "task:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e66"
	partyId := "P2"
	RemoveNeedExecuteTask(database, taskId1)
	count := 0
	iter := database.NewIteratorWithPrefixAndStart(needExecuteTaskPrefix, nil)
	for iter.Next() {
		count++
	}
	assert.Equal(t, 12, count)
	RemoveNeedExecuteTaskByPartyId(database, taskId2, partyId)
	assert.Equal(t, 11, count-1)
}

func TestRecoveryNeedExecuteTask(t *testing.T) {
	result := RecoveryNeedExecuteTask(NeedExecuteTask())
	count := 0
	checkRepet := func(partyId string, p []string) bool {
		for _, value := range p {
			if partyId == value {
				return false
			}
		}
		return true
	}
	taskIdsResult := make([]string, 0)
	partyIdsResult := make([]string, 0)
	for taskId, value := range result {
		taskIdsResult = append(taskIdsResult, taskId)
		for partyId, _ := range value {
			if true == checkRepet(partyId, partyIdsResult) {
				partyIdsResult = append(partyIdsResult, partyId)
			}
			count++
		}
	}
	assert.Equal(t, len(taskIds), len(taskIdsResult))
	assert.Equal(t, len(partyIds), len(partyIdsResult))
	assert.Equal(t, 15, count)
}
