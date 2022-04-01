package rawdb

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/db"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
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
	data01 := &libtypes.TaskPB{
		IdentityId: "identity",
		NodeId:     "nodeid",
		NodeName:   "nodename",
		DataId:     "taskId",
		DataStatus: libtypes.DataStatus_DataStatus_Deleted,
		TaskId:     "taskID",
		TaskName:   "taskName",
		State:      libtypes.TaskState_TaskState_Failed,
		Reason:     "reason",
		EventCount: 4,
		Desc:       "desc",
		CreateAt:   timeutils.UnixMsecUint64(),
		EndAt:      timeutils.UnixMsecUint64(),
	}
	StoreLocalTask(database, types.NewTask(data01))

	res, _ := QueryLocalTask(database, data01.TaskId)
	assert.Assert(t, strings.EqualFold(data01.TaskId, res.GetTaskId()))

	// test update state
	res.GetTaskData().State = libtypes.TaskState_TaskState_Failed
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
	seedNodeInfo := &pb.SeedPeer{
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
	registered := &pb.YarnRegisteredPeerDetail{
		Id:           "id",
		InternalIp:   "internalIp",
		InternalPort: "9999",
		ExternalIp:   "externalIp",
		ExternalPort: "999",
		ConnState:    1,
	}
	StoreRegisterNode(database, pb.PrefixTypeJobNode, registered)

	// get seed
	r, _ := QueryRegisterNode(database, pb.PrefixTypeJobNode, "id")
	t.Logf("registered info : %v", r)
	assert.Assert(t, strings.EqualFold("id", r.Id))

	// read all
	registeredNodes, _ := QueryAllRegisterNodes(database, pb.PrefixTypeJobNode)
	assert.Assert(t, len(registeredNodes) == 1)

	// delete
	RemoveRegisterNode(database, pb.PrefixTypeJobNode, "id")

	registeredNodes, _ = QueryAllRegisterNodes(database, pb.PrefixTypeJobNode)
	assert.Assert(t, len(registeredNodes) == 0)

	// delete
	RemoveRegisterNodes(database, pb.PrefixTypeJobNode)

	registeredNodes, _ = QueryAllRegisterNodes(database, pb.PrefixTypeJobNode)
	assert.Assert(t, len(registeredNodes) == 0)
}

func TestTaskEvent(t *testing.T) {
	database := db.NewMemoryDatabase()
	taskEvent := &libtypes.TaskEvent{
		Type:       "taskEventType",
		IdentityId: "taskEventIdentity",
		TaskId:     "taskEventTaskId",
		Content:    "taskEventContent",
		CreateAt:   timeutils.UnixMsecUint64(),
	}
	StoreTaskEvent(database, taskEvent)

	taskEvent2 := &libtypes.TaskEvent{
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
	nodeAlias := &libtypes.Organization{
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
	localResource01 := &libtypes.LocalResourcePB{
		IdentityId:     "01-identity",
		NodeId:         "01-nodeId",
		NodeName:       "01-nodename",
		JobNodeId:      "01",
		DataId:         "01-dataId",
		DataStatus:     libtypes.DataStatus_DataStatus_Deleted,
		State:          libtypes.PowerState_PowerState_Created,
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

	localResource02 := &libtypes.LocalResourcePB{
		IdentityId:     "01-identity",
		NodeId:         "01-nodeId",
		NodeName:       "01-nodename",
		JobNodeId:      "02",
		DataId:         "01-dataId",
		DataStatus:     libtypes.DataStatus_DataStatus_Normal,
		State:         libtypes.PowerState_PowerState_Created,
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
	localMetadata01 := &libtypes.MetadataPB{
		MetadataId:           "metadataId",
		IdentityId:           "identityId",
		NodeId:               "nodeId",
		NodeName:             "nodeName",
		DataId:               "dataId",
		DataStatus:           0,
		OriginId:             "originId",
		TableName:            "tableName",
		FilePath:             "filePath",
		Desc:                 "desc",
		Rows:                 0,
		Columns:              0,
		Size_:                0,
		FileType:             0,
		State:                0,
		HasTitle:             false,
		MetadataColumns:      nil,
		Industry:             "",
	}
	b, _ := localMetadata01.Marshal()
	_ = b
	StoreLocalMetadata(database, types.NewMetadata(localMetadata01))

	localMetadata02 := &libtypes.MetadataPB{
		MetadataId:           "metadataId-02",
		IdentityId:           "identityId-02",
		NodeId:               "nodeId-02",
		NodeName:             "nodeName",
		DataId:               "dataId",
		DataStatus:           0,
		OriginId:             "originId",
		TableName:            "tableName",
		FilePath:             "filePath",
		Desc:                 "desc",
		Rows:                 0,
		Columns:              0,
		Size_:                0,
		FileType:             0,
		State:                0,
		HasTitle:             false,
		MetadataColumns:      nil,
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