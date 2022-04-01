package rawdb

import (
	"bytes"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/bytesutil"
	"github.com/RosettaFlow/Carrier-Go/common/rlputil"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/db"
	dbtype "github.com/RosettaFlow/Carrier-Go/lib/db"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"
	"gotest.tools/assert"
	"math/rand"
	"strings"
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
func NeedExecuteTask() (KeyValueStore, dbtype.TaskArrayPB) {
	database := db.NewMemoryDatabase()
	var taskList dbtype.TaskArrayPB
	for _, taskId := range taskIds {
		taskPB := &libtypes.TaskPB{
			TaskId: taskId,
		}
		task := types.NewTask(taskPB)
		taskList.TaskList = append(taskList.TaskList, taskPB)
		for _, partyId := range partyIds {
			remotepid := "remotepid"
			localTaskOrganization := &libtypes.TaskOrganization{
				PartyId:    partyId,
				NodeName:   "NodeName",
				NodeId:     "NodeId_0001",
				IdentityId: "IdentityId_0001",
			}
			remoteTaskOrganization := &libtypes.TaskOrganization{
				PartyId:    partyId,
				NodeName:   "NodeName",
				NodeId:     "NodeId_0002",
				IdentityId: "IdentityId_0002",
			}
			localResource := &types.PrepareVoteResource{
				Id:      "PrepareVoteResourceId",
				Ip:      "2.2.2.2",
				Port:    "5555",
				PartyId: partyId,
			}
			resources := &twopcpb.ConfirmTaskPeerInfo{}

			_task := types.NewNeedExecuteTask(peer.ID(remotepid),1, 2, localTaskOrganization, remoteTaskOrganization, task.GetTaskId(),
				3, localResource, resources, nil)
			if err := StoreNeedExecuteTask(database, _task); err != nil {
				fmt.Printf("StoreNeedExecuteTask fail,taskId %s\n", taskId)
			}
		}
	}
	return database, taskList
}
func TestStoreNeedExecuteTask(t *testing.T) {
	NeedExecuteTask()
}

func TestDeleteNeedExecuteTask(t *testing.T) {
	database, _ := NeedExecuteTask()
	taskId1 := "task:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e65"
	taskId2 := "task:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e66"
	partyId := "P2"
	RemoveNeedExecuteTask(database, taskId1)
	count := 0
	it := database.NewIteratorWithPrefixAndStart(needExecuteTaskKeyPrefix, nil)
	defer it.Release()
	for it.Next() {
		count++
	}
	assert.Equal(t, 12, count)
	RemoveNeedExecuteTaskByPartyId(database, taskId2, partyId)
	assert.Equal(t, 11, count-1)
}
func MockQueryLocalTask(taskList dbtype.TaskArrayPB, taskId string) (*types.Task, error) {
	for _, task := range taskList.GetTaskList() {
		if strings.EqualFold(task.TaskId, taskId) {
			return types.NewTask(task), nil
		}
	}
	return nil, ErrNotFound
}
func TestRecoveryNeedExecuteTask(t *testing.T) {
	prefix := needExecuteTaskKeyPrefix
	database, localTask := NeedExecuteTask()
	runningTaskCache := make(map[string]map[string]*types.NeedExecuteTask, 0)
	if err := ForEachNeedExecuteTask(database, func(key, value []byte) error {
		if len(key) != 0 && len(value) != 0 {

			// task:${taskId hex} == 5 + 2 + 64 == "taskId:" + "0x" + "e33...fe4"
			taskId := string(key[len(prefix) : len(prefix)+71])
			partyId := string(key[len(prefix)+71:])

			task, err := MockQueryLocalTask(localTask, taskId)
			if nil != err {
				return fmt.Errorf("query local task failed on recover needExecuteTask from db, %s, taskId: {%s}", err, taskId)
			}

			var res libtypes.NeedExecuteTask

			if err := proto.Unmarshal(value, &res); nil != err {
				return fmt.Errorf("Unmarshal needExecuteTask failed, %s", err)
			}

			cache, ok := runningTaskCache[taskId]
			if !ok {
				cache = make(map[string]*types.NeedExecuteTask, 0)
			}
			cache[partyId] = types.NewNeedExecuteTask(
				peer.ID(res.GetRemotePid()),
				res.GetLocalTaskRole(),
				res.GetRemoteTaskRole(),
				res.GetLocalTaskOrganization(),
				res.GetRemoteTaskOrganization(),
				task.GetTaskId(),
				types.TaskActionStatus(bytesutil.BytesToUint16(res.GetConsStatus())),
				types.NewPrepareVoteResource(
					res.GetLocalResource().GetId(),
					res.GetLocalResource().GetIp(),
					res.GetLocalResource().GetPort(),
					res.GetLocalResource().GetPartyId(),
				),
				res.GetResources(),
				nil,
			)
			runningTaskCache[taskId] = cache
		}
		return nil
	}); nil != err {
		log.WithError(err).Fatalf("recover needExecuteTask failed")
	}

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
	for taskId, value := range runningTaskCache {
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

func TestStoreMessageCache(t *testing.T) {
	var err error
	database := db.NewMemoryDatabase()
	err = StoreMessageCache(database, &types.PowerMsg{
		PowerId:   "PowerId_111111",
		JobNodeId: "JobNodeId_222222",
		CreateAt:  2233,
	})
	if err != nil {
		panic(err.Error())
	}

	err = StoreMessageCache(database, &types.MetadataMsg{
		MetadataId: "MetadataId",
		MetadataSummary: &libtypes.MetadataSummary{
			MetadataId: "MetadataId",
			Desc:       "",
			Industry:   "",
			State:      2,
			PublishAt:  3344,
			UpdateAt:   4455,
		},
	})
	if err != nil {
		panic(err.Error())
	}

	err = StoreMessageCache(database, &types.MetadataAuthorityMsg{
		MetadataAuthId: "MetadataAuthId",
		User:           "user1",
		UserType:       2,
		Auth:           &libtypes.MetadataAuthority{},
		Sign:           []byte("sign"),
		CreateAt:       9988,
	})
	if err != nil {
		panic(err.Error())
	}

	err = StoreMessageCache(database, &types.TaskMsg{
		Data: types.NewTask(&libtypes.TaskPB{
			TaskId: "task:0xe7bdb5af4de9d851351c680fb0a9bfdff72bdc4ea86da3c2006d6a7a7d335e65",
		}),
	})
	if err != nil {
		panic(err.Error())
	}
}

func TestQueryRemoveMetadataAuthorityMsgArr(t *testing.T) {
	var err error
	database := db.NewMemoryDatabase()
	err = StoreMessageCache(database, &types.MetadataAuthorityMsg{
		MetadataAuthId: "MetadataAuthId",
		User:           "user1",
		UserType:       2,
		Auth:           &libtypes.MetadataAuthority{},
		Sign:           []byte("sign"),
		CreateAt:       9988,
	})
	if err != nil {
		panic(err.Error())
	}
	result, _ := QueryMetadataAuthorityMsgArr(database)
	fmt.Println(result)

	//TestRemove
	RemoveMetadataAuthMsg(database, "MetadataAuthId")
	assert.Equal(t, 0, database.Len())
}

func TestDataResourceTable(t *testing.T) {
	database := db.NewMemoryDatabase()

	dataNodeId := "dataService_192.168.10.150_8700"
	err := StoreDataResourceTable(database, types.NewDataResourceTable(dataNodeId, 63141883904, 23408070656, true))
	assert.NilError(t, err, "Failed to call StoreDataResourceTable()")

	arr, err := QueryDataResourceTables(database)

	assert.NilError(t, err, "Failed to call QueryDataResourceTables()")
	assert.Equal(t, len(arr), 1, "Not found dataNodes")
	assert.Equal(t, dataNodeId, arr[0].GetNodeId(), fmt.Sprintf("executed %s, but actual %s", dataNodeId, arr[0].GetNodeId()))
	t.Log("dataNode", arr[0].String())
}

func TestLocalTaskExecuteStatus(t *testing.T) {
	v := []byte{}


	val := bytesutil.BytesToUint32(v)
	val |= OnConsensusExecuteTaskStatus.Uint32()

	assert.Equal(t, OnConsensusExecuteTaskStatus.Uint32(), val, "failed test")

	v = bytesutil.Uint32ToBytes(OnConsensusExecuteTaskStatus.Uint32())

	val = bytesutil.BytesToUint32(v)
	val |= OnRunningExecuteStatus.Uint32()

	assert.Equal(t, OnConsensusExecuteTaskStatus.Uint32()|OnRunningExecuteStatus.Uint32(), val, "failed test")
	assert.Equal(t, val&OnRunningExecuteStatus.Uint32(), OnRunningExecuteStatus.Uint32(), "failed test")
}