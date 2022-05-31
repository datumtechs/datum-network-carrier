package rawdb

import (
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	"github.com/datumtechs/datum-network-carrier/db"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestLocalTask(t *testing.T) {
	database := db.NewMemoryDatabase()
	data01 := &carriertypespb.TaskPB{
		DataId:     "taskId",
		DataStatus: commonconstantpb.DataStatus_DataStatus_Invalid,
		TaskId:     "taskID",
		TaskName:   "taskName",
		State:      commonconstantpb.TaskState_TaskState_Failed,
		Reason:     "reason",
		Desc:       "desc",
		CreateAt:   timeutils.UnixMsecUint64(),
		EndAt:      timeutils.UnixMsecUint64(),
	}
	StoreLocalTask(database, types.NewTask(data01))

	res, _ := QueryLocalTask(database, data01.TaskId)
	assert.Assert(t, strings.EqualFold(data01.TaskId, res.GetTaskId()))

	// test update state
	res.GetTaskData().State = commonconstantpb.TaskState_TaskState_Failed
	RemoveLocalTask(database, data01.TaskId)
	StoreLocalTask(database, res)

	res, _ = QueryLocalTask(database, data01.TaskId)
	assert.Equal(t, data01.State, res.GetTaskData().State)

	taskList, _ := QueryAllLocalTasks(database)
	assert.Assert(t, len(taskList) == 1)

	// delete
	RemoveLocalTask(database, data01.TaskId)

	taskList, _ = QueryAllLocalTasks(database)
	assert.Assert(t, len(taskList) == 0)

	// delete
	RemoveSeedNodes(database)

	taskList, _ = QueryAllLocalTasks(database)
	assert.Assert(t, len(taskList) == 0)
}

func TestSeedNode(t *testing.T) {
	// write seed
	database := db.NewMemoryDatabase()
	seedNodeInfo := &carrierapipb.SeedPeer{
		Addr:"addr1",
		IsDefault: false,
		ConnState:    1,
	}
	StoreSeedNode(database, seedNodeInfo)

	// read all
	seedNodes, _ := QueryAllSeedNodes(database)
	assert.Assert(t, len(seedNodes) == 1)

	// delete
	RemoveSeedNode(database, "id")

	seedNodes, _ = QueryAllSeedNodes(database)
	assert.Assert(t, len(seedNodes) == 0)

	// delete
	RemoveSeedNodes(database)

	seedNodes, _ = QueryAllSeedNodes(database)
	assert.Assert(t, len(seedNodes) == 0)
}

func TestRegisteredNode(t *testing.T) {
	// write seed
	database := db.NewMemoryDatabase()
	registered := &carrierapipb.YarnRegisteredPeerDetail{
		Id:           "id",
		InternalIp:   "internalIp",
		InternalPort: "9999",
		ExternalIp:   "externalIp",
		ExternalPort: "999",
		ConnState:    1,
	}
	StoreRegisterNode(database, carrierapipb.PrefixTypeJobNode, registered)

	// get seed
	r, _ := QueryRegisterNode(database, carrierapipb.PrefixTypeJobNode, "id")
	t.Logf("registered info : %v", r)
	assert.Assert(t, strings.EqualFold("id", r.Id))

	// read all
	registeredNodes, _ := QueryAllRegisterNodes(database, carrierapipb.PrefixTypeJobNode)
	assert.Assert(t, len(registeredNodes) == 1)

	// delete
	RemoveRegisterNode(database, carrierapipb.PrefixTypeJobNode, "id")

	registeredNodes, _ = QueryAllRegisterNodes(database, carrierapipb.PrefixTypeJobNode)
	assert.Assert(t, len(registeredNodes) == 0)

	// delete
	RemoveRegisterNodes(database, carrierapipb.PrefixTypeJobNode)

	registeredNodes, _ = QueryAllRegisterNodes(database, carrierapipb.PrefixTypeJobNode)
	assert.Assert(t, len(registeredNodes) == 0)
}

func TestTaskEvent(t *testing.T) {
	database := db.NewMemoryDatabase()
	taskEvent := &carriertypespb.TaskEvent{
		Type:       "taskEventType",
		IdentityId: "taskEventIdentity",
		TaskId:     "taskEventTaskId",
		Content:    "taskEventContent",
		CreateAt:   timeutils.UnixMsecUint64(),
	}
	StoreTaskEvent(database, taskEvent)

	taskEvent2 := &carriertypespb.TaskEvent{
		Type:       "taskEventType-02",
		IdentityId: "taskEventIdentity",
		TaskId:     "taskEventTaskId",
		Content:    "taskEventContent-02",
		CreateAt:   timeutils.UnixMsecUint64(),
	}
	StoreTaskEvent(database, taskEvent2)

	revent, _ := QueryTaskEvent(database, "taskEventTaskId")
	t.Logf("task evengine info : %v", len(revent))
	assert.Assert(t, strings.EqualFold("taskEventIdentity", revent[0].IdentityId))

	// read all
	taskEvents, _ := QueryAllTaskEvents(database)
	assert.Assert(t, len(taskEvents) == 2)

	// delete
	RemoveTaskEvent(database, "taskEventTaskId")

	taskEvents, _ = QueryAllTaskEvents(database)
	assert.Assert(t, len(taskEvents) == 0)
}

func TestLocalIdentity(t *testing.T) {
	database := db.NewMemoryDatabase()
	nodeAlias := &carriertypespb.Organization{
		NodeName:   "node-name",
		NodeId:     "node-nodeId",
		IdentityId: "node-identityId",
	}
	StoreLocalIdentity(database, nodeAlias)
	StoreLocalIdentity(database, nodeAlias)

	rnode, _ := QueryLocalIdentity(database)
	assert.Equal(t, rnode.IdentityId, nodeAlias.IdentityId)
	assert.Equal(t, rnode.NodeId, nodeAlias.NodeId)
	assert.Equal(t, rnode.NodeName, nodeAlias.NodeName)

	RemoveLocalIdentity(database)
	rnode, _ = QueryLocalIdentity(database)
	require.Nil(t, rnode)
}

func TestLocalResource(t *testing.T) {
	database := db.NewMemoryDatabase()
	localResource01 := &carriertypespb.LocalResourcePB{
		JobNodeId:      "01",
		DataId:         "01-dataId",
		DataStatus:     commonconstantpb.DataStatus_DataStatus_Invalid,
		State:          commonconstantpb.PowerState_PowerState_Created,
		TotalMem:       111,
		UsedMem:        222,
		TotalProcessor: 11,
		UsedProcessor:  33,
		TotalBandwidth: 33,
		UsedBandwidth:  44,
	}
	b, _ := localResource01.Marshal()
	_ = b
	StoreLocalResource(database, types.NewLocalResource(localResource01))

	localResource02 := &carriertypespb.LocalResourcePB{
		JobNodeId:      "02",
		DataId:         "01-dataId",
		DataStatus:     commonconstantpb.DataStatus_DataStatus_Valid,
		State:         commonconstantpb.PowerState_PowerState_Created,
		TotalMem:       111,
		UsedMem:        222,
		TotalProcessor: 11,
		UsedProcessor:  33,
		TotalBandwidth: 33,
		UsedBandwidth:  44,
	}
	StoreLocalResource(database, types.NewLocalResource(localResource02))

	r1, _ := QueryLocalResource(database, localResource01.JobNodeId)
	assert.Equal(t, types.NewLocalResource(localResource01).Hash(), r1.Hash())

	array, _ := QueryAllLocalResource(database)
	require.True(t, array.Len() == 2)

	RemoveLocalResource(database, localResource01.JobNodeId)
	array, _ = QueryAllLocalResource(database)
	require.True(t, array.Len() == 1)
}

func TestLocalMetadata(t *testing.T) {
	database := db.NewMemoryDatabase()
	localMetadata01 := &carriertypespb.MetadataPB{
		MetadataId:           "metadataId",
		DataId:               "dataId",
		DataStatus:           0,
		Desc:                 "desc",
		State:                0,
		Industry:             "",
	}
	b, _ := localMetadata01.Marshal()
	_ = b
	StoreLocalMetadata(database, types.NewMetadata(localMetadata01))

	localMetadata02 := &carriertypespb.MetadataPB{
		MetadataId:           "metadataId-02",
		DataId:               "dataId",
		DataStatus:           0,
		Desc:                 "desc",
		State:                0,
		Industry:             "",
	}
	StoreLocalMetadata(database, types.NewMetadata(localMetadata02))

	r1, _ := QueryLocalMetadata(database, localMetadata01.MetadataId)
	assert.Equal(t, types.NewMetadata(localMetadata01).Hash(), r1.Hash())

	array, _ := QueryAllLocalMetadata(database)
	require.True(t, array.Len() == 2)

	RemoveLocalMetadata(database, localMetadata01.MetadataId)
	array, _ = QueryAllLocalMetadata(database)
	require.True(t, array.Len() == 1)
}
func TestWriteScheduling(t *testing.T) {
	var taskBullt *types.TaskBullet
	database := db.NewMemoryDatabase()
	var count uint32
	var starve bool
	rand.Seed(time.Now().UnixNano())

	queue := make(types.TaskBullets,0)
	starveQueue := make(types.TaskBullets,0)

	for {
		if count > 20 {
			break
		}
		count++
		if a := rand.Intn(2); a == 1 {
			starve = true
		} else {
			starve = false
		}

		taskBullt = &types.TaskBullet{
			TaskId:  "taskId_" + strconv.FormatUint(uint64(count), 10),
			Starve:  starve,
			Term:    count,
			Resched: count,
		}
		StoreTaskBullet(database, taskBullt)
		if starve == true {
			starveQueue.Push(taskBullt)
		} else {
			queue.Push(taskBullt)
		}
	}
	fmt.Println("-----------starveQueue--------------")
	for _, value := range starveQueue {
		fmt.Println(value.IsStarve(),value.GetTaskId())
	}
	fmt.Println("-----------queue--------------------")
	for _,value:=range queue{
		fmt.Println(value.IsStarve(),value.GetTaskId())
	}
}