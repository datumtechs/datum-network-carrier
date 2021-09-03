// Copyright (C) 2021 The RosettaNet Authors.

package rawdb

import (
	"bytes"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	dbtype "github.com/RosettaFlow/Carrier-Go/lib/db"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/sirupsen/logrus"
	"strings"
	"sync/atomic"
)

const seedNodeToKeep = 50
const registeredNodeToKeep = 50

// ReadLocalIdentity retrieves the identity of local.
func ReadLocalIdentity(db DatabaseReader) (*apipb.Organization, error) {
	var blob apipb.Organization
	enc, err := db.Get(localIdentityKey)
	if nil != err {
		return nil, err
	}
	if err := blob.Unmarshal(enc); nil != err {
		return nil, err
	}
	return &apipb.Organization{
		NodeName:       blob.GetNodeName(),
		NodeId:     blob.GetNodeId(),
		IdentityId: blob.GetIdentityId(),
	}, nil
}

// WriteLocalIdentity stores the local identity.
func WriteLocalIdentity(db DatabaseWriter, localIdentity *apipb.Organization) {
	pb := &apipb.TaskOrganization{
		PartyId:              "",
		IdentityId:           localIdentity.GetIdentityId(),
		NodeId:               localIdentity.GetNodeId(),
		NodeName:             localIdentity.GetNodeName(),
	}
	enc, _ := pb.Marshal()
	if err := db.Put(localIdentityKey, enc); err != nil {
		log.WithError(err).Fatal("Failed to store local identity")
	}
}

// DeleteLocalIdentity deletes the local identity
func DeleteLocalIdentity(db DatabaseDeleter) {
	if err := db.Delete(localIdentityKey); err != nil {
		log.WithError(err).Fatal("Failed to delete local identity")
	}
}

// ReadSeedNode retrieves the seed node with the corresponding nodeId.
func ReadRunningTaskIDList(db DatabaseReader, jobNodeId string) ([]string, error) {
	blob, _ := db.Get(runningTaskIDListKey(jobNodeId))
	var array dbtype.StringArrayPB
	if len(blob) > 0 {
		if err := array.Unmarshal(blob); err != nil {
			log.WithError(err).Fatal("Failed to decode old RunningTaskIdList")
		}
		return array.GetArray(), nil
	}
	return nil, ErrNotFound
}

func WriteRunningTaskIDList(db KeyValueStore, jobNodeId, taskId string) {
	blob, err := db.Get(runningTaskIDListKey(jobNodeId))
	if err != nil {
		log.WithError(err).Warn("Failed to load old RunningTaskIdList")
	}
	var array dbtype.StringArrayPB
	if len(blob) > 0 {
		if err := array.Unmarshal(blob); err != nil {
			log.WithError(err).Fatal("Failed to decode old RunningTaskIdList")
		}
	}
	for _, s := range array.GetArray() {
		if strings.EqualFold(s, taskId) {
			log.WithFields(logrus.Fields{ "id": s }).Info("Skip duplicated running task id")
			return
		}
	}
	array.Array = append(array.Array, taskId)
	data, err := array.Marshal()
	if err != nil {
		log.WithError(err).Fatal("Failed to encode RunningTaskIdList")
	}
	if err := db.Put(runningTaskIDListKey(jobNodeId), data); err != nil {
		log.WithError(err).Fatal("Failed to write RunningTaskIdList")
	}
}

// DeleteRunningTaskIDList deletes the running taskID list of jobNode from the database with a special id
func DeleteRunningTaskIDList(db KeyValueStore, jobNodeId, taskId string) {
	blob, err := db.Get(runningTaskIDListKey(jobNodeId))
	if err != nil {
		log.WithError(err).Fatal("Failed to load old RunningTaskIdList")
	}
	var array dbtype.StringArrayPB
	if len(blob) > 0 {
		if err := array.Unmarshal(blob); err != nil {
			log.WithError(err).Fatal("Failed to decode old RunningTaskIdList")
		}
	}
	for idx, s := range array.GetArray() {
		if strings.EqualFold(s, taskId) {
			array.Array = append(array.Array[:idx], array.Array[idx+1:]...)
			break
		}
	}
	data, err := array.Marshal()
	if err != nil {
		log.WithError(err).Fatal("Failed to encode RunningTaskIdList")
	}
	if err := db.Put(runningTaskIDListKey(jobNodeId), data); err != nil {
		log.WithError(err).Fatal("Failed to write RunningTaskIdList")
	}
}

func ReadRunningTaskCountForJobNode(db DatabaseReader, jobNodeId string) uint32 {
	var v dbtype.Uint32PB
	enc, _ := db.Get(runningTaskCountForJobNodeKey(jobNodeId))
	v.Unmarshal(enc)
	return v.GetV()
}

func IncreaseRunningTaskCountForJobNode(db KeyValueStore, jobNodeId string) uint32 {
	has, _ := db.Has(runningTaskCountForJobNodeKey(jobNodeId))
	if !has {
		WriteRunningTaskCountForJobNode(db, jobNodeId, 1)
	}
	count := ReadRunningTaskCountForJobNode(db, jobNodeId)
	newCount := atomic.AddUint32(&count, 1)
	WriteRunningTaskCountForJobNode(db, jobNodeId, newCount)
	return newCount
}

func WriteRunningTaskCountForJobNode(db DatabaseWriter, jobNodeId string, count uint32) {
	pb := dbtype.Uint32PB{
		V:                    count,
	}
	enc, _ := pb.Marshal()
	if err := db.Put(runningTaskCountForJobNodeKey(jobNodeId), enc); err != nil {
		log.WithError(err).Fatal("Failed to store RunningTaskCountForJobNode")
	}
}

func DeleteRunningTaskCountForJobNode(db DatabaseDeleter, jobNodeId string) {
	if err := db.Delete(runningTaskCountForJobNodeKey(jobNodeId)); err != nil {
		log.WithError(err).Fatal("Failed to delete RunningTaskCountForJobNode")
	}
}

func ReadRunningTaskCountForOrg(db DatabaseReader) uint32 {
	var v dbtype.Uint32PB
	enc, _ := db.Get(runningTaskCountForOrgKey)
	v.Unmarshal(enc)
	return v.GetV()
}

func IncreaseRunningTaskCountForOrg(db KeyValueStore) uint32 {
	has, _ := db.Has(runningTaskCountForOrgKey)
	if !has {
		WriteRunningTaskCountForOrg(db, 1)
	}
	count := ReadRunningTaskCountForOrg(db)
	newCount := atomic.AddUint32(&count, 1)
	WriteRunningTaskCountForOrg(db, newCount)
	return newCount
}

func WriteRunningTaskCountForOrg(db DatabaseWriter, count uint32) {
	pb := dbtype.Uint32PB{
		V:                    count,
	}
	enc, _ := pb.Marshal()
	if err := db.Put(runningTaskCountForOrgKey, enc); err != nil {
		log.WithError(err).Fatal("Failed to store RunningTaskCountForOrg")
	}
}

// ReadIdentityStr retrieves the identity string.
func ReadIdentityStr(db DatabaseReader) string {
	var str dbtype.StringPB
	enc, _ := db.Get(identityKey)
	str.Unmarshal(enc)
	return str.GetV()
}

// WriteIdentityStr stores the identity.
func WriteIdentityStr(db DatabaseWriter, identity string) {
	pb := dbtype.StringPB{
		V:                    identity,
	}
	enc, _ := pb.Marshal()
	if err := db.Put(identityKey, enc); err != nil {
		log.WithError(err).Fatal("Failed to store identity")
	}
}

// DeleteIdentityStr deletes the identity.
func DeleteIdentityStr(db DatabaseDeleter) {
	if err := db.Delete(identityKey); err != nil {
		log.WithError(err).Fatal("Failed to delete identity")
	}
}

// ReadYarnName retrieves the name of yarn.
func ReadYarnName(db DatabaseReader) (string, error) {
	var yarnName dbtype.StringPB
	enc, err := db.Get(yarnNameKey)
	if nil != err {
		return "", err
	}
	if err := yarnName.Unmarshal(enc); nil != err {
		return "", err
	}
	return yarnName.GetV(), nil
}

// WriteYarnName stores the name of yarn.
func WriteYarnName(db DatabaseWriter, yarnName string) {
	pb := dbtype.StringPB{
		V:                    yarnName,
	}
	enc, _ := pb.Marshal()
	if err := db.Put(yarnNameKey, enc); err != nil {
		log.WithError(err).Fatal("Failed to store yarn name")
	}
}

// DeleteYarnName deletes the name of yarn.
func DeleteYarnName(db DatabaseDeleter) {
	if err := db.Delete(yarnNameKey); err != nil {
		log.WithError(err).Fatal("Failed to delete yarn name")
	}
}

// ReadSeedNode retrieves the seed node with the corresponding nodeId.
func ReadSeedNode(db DatabaseReader, nodeId string) (*pb.SeedPeer, error) {
	blob, err := db.Get(seedNodeKey)
	if err != nil {
		return nil, err
	}
	var seedNodes dbtype.SeedNodeListPB
	if err := seedNodes.Unmarshal(blob); err != nil {
		return nil, err
	}
	for _, seed := range seedNodes.GetSeedNodeList() {
		if strings.EqualFold(seed.Id, nodeId) {
			return &pb.SeedPeer{
				Id:           seed.Id,
				InternalIp:   seed.InternalIp,
				InternalPort: seed.InternalPort,
				ConnState:    pb.ConnState(seed.ConnState),
			}, nil
		}
	}
	return nil, ErrNotFound
}

// ReadAllSeedNodes retrieves all the seed nodes in the database.
// All returned seed nodes are sorted in reverse.
func ReadAllSeedNodes(db DatabaseReader) ([]*pb.SeedPeer, error) {
	blob, err := db.Get(seedNodeKey)
	if err != nil {
		return nil, err
	}
	var seedNodes dbtype.SeedNodeListPB
	if err := seedNodes.Unmarshal(blob); err != nil {
		return nil, err
	}
	var nodes []*pb.SeedPeer
	for _, seed := range seedNodes.SeedNodeList {
		nodes = append(nodes, &pb.SeedPeer{
			Id:           seed.Id,
			InternalIp:   seed.InternalIp,
			InternalPort: seed.InternalPort,
			ConnState:    pb.ConnState(seed.ConnState),
		})
	}
	return nodes, nil
}

// WriteSeedNodes serializes the seed node into the database. If the cumulated
// seed node exceeds the limitation, the oldest will be dropped.
func WriteSeedNodes(db KeyValueStore, seedNode *pb.SeedPeer) {
	blob, err := db.Get(seedNodeKey)
	if err != nil {
		log.Warn("Failed to load old seed nodes", "error", err)
	}
	var seedNodes dbtype.SeedNodeListPB
	if len(blob) > 0 {
		if err := seedNodes.Unmarshal(blob); err != nil {
			log.WithError(err).Fatal("Failed to decode old seed nodes")
		}

	}
	for _, s := range seedNodes.GetSeedNodeList() {
		if strings.EqualFold(s.Id, seedNode.Id) {
			log.WithFields(logrus.Fields{ "id": s.Id }).Info("Skip duplicated seed node")
			return
		}
	}
	seedNodes.SeedNodeList = append(seedNodes.SeedNodeList, &dbtype.SeedNodePB{
		Id:                   seedNode.Id,
		InternalIp:           seedNode.InternalIp,
		InternalPort:         seedNode.InternalPort,
		ConnState:            int32(seedNode.ConnState),
	})
	// max limit for store seed node.
	if len(seedNodes.SeedNodeList) > seedNodeToKeep {
		seedNodes.SeedNodeList = seedNodes.SeedNodeList[:seedNodeToKeep]
	}
	data, err := seedNodes.Marshal()
	if err != nil {
		log.WithError(err).Fatal("Failed to encode seed node")
	}
	if err := db.Put(seedNodeKey, data); err != nil {
		log.WithError(err).Fatal("Failed to write seed node")
	}
}

// DeleteSeedNode deletes the seed nodes from the database with a special id
func DeleteSeedNode(db KeyValueStore, id string) {
	blob, err := db.Get(seedNodeKey)
	if err != nil {
		log.Warn("Failed to load old seed nodes", "error", err)
	}
	var seedNodes dbtype.SeedNodeListPB
	if len(blob) > 0 {
		if err := seedNodes.Unmarshal(blob); err != nil {
			log.WithError(err).Fatal("Failed to decode old seed nodes")
		}
	}
	// need to test.
	for idx, s := range seedNodes.GetSeedNodeList() {
		if strings.EqualFold(s.Id, id) {
			seedNodes.SeedNodeList = append(seedNodes.SeedNodeList[:idx], seedNodes.SeedNodeList[idx+1:]...)
			break
		}
	}
	data, err := seedNodes.Marshal()
	if err != nil {
		log.WithError(err).Fatal("Failed to encode seed nodes")
	}
	if err := db.Put(seedNodeKey, data); err != nil {
		log.WithError(err).Fatal("Failed to write seed nodes")
	}
}

// DeleteSeedNodes deletes all the seed nodes from the database
func DeleteSeedNodes(db DatabaseDeleter) {
	if err := db.Delete(seedNodeKey); err != nil {
		log.WithError(err).Fatal("Failed to delete seed node")
	}
}

func registryNodeKey(nodeType pb.RegisteredNodeType) []byte {
	var key []byte
	if nodeType == pb.PrefixTypeJobNode {
		key = calcNodeKey
	}
	if nodeType == pb.PrefixTypeDataNode {
		key = dataNodeKey
	}
	return key
}

// ReadRegisterNode retrieves the register node with the corresponding nodeId.
func ReadRegisterNode(db DatabaseReader, nodeType pb.RegisteredNodeType, nodeId string) (*pb.YarnRegisteredPeerDetail, error) {
	blob, err := db.Get(registryNodeKey(nodeType))
	if err != nil {
		return nil, err
	}
	var registeredNodes dbtype.RegisteredNodeListPB
	if err := registeredNodes.Unmarshal(blob); err != nil {
		return nil, err
	}
	lis := registeredNodes.GetRegisteredNodeList()
	for _, registered := range lis {
		if strings.EqualFold(registered.Id, nodeId) {
			return &pb.YarnRegisteredPeerDetail{
				Id:           registered.Id,
				InternalIp:   registered.InternalIp,
				InternalPort: registered.InternalPort,
				ExternalIp:   registered.ExternalIp,
				ExternalPort: registered.ExternalPort,
				ConnState:    pb.ConnState(registered.ConnState),
			}, nil
		}
	}
	return nil, ErrNotFound
}

// ReadAllRegisterNodes retrieves all the registered nodes in the database.
// All returned registered nodes are sorted in reverse.
func ReadAllRegisterNodes(db DatabaseReader, nodeType pb.RegisteredNodeType) ([]*pb.YarnRegisteredPeerDetail, error) {
	blob, err := db.Get(registryNodeKey(nodeType))
	if err != nil {
		return nil, err
	}
	var registeredNodes dbtype.RegisteredNodeListPB
	if err := registeredNodes.Unmarshal(blob); err != nil {
		return nil, err
	}
	var nodes []*pb.YarnRegisteredPeerDetail
	for _, registered := range registeredNodes.GetRegisteredNodeList() {
		nodes = append(nodes, &pb.YarnRegisteredPeerDetail{
			Id:           registered.Id,
			InternalIp:   registered.InternalIp,
			InternalPort: registered.InternalPort,
			ExternalIp:   registered.ExternalIp,
			ExternalPort: registered.ExternalPort,
			ConnState:    pb.ConnState(registered.ConnState),
		})
	}
	return nodes, nil
}

// WriteRegisterNodes serializes the registered node into the database. If the cumulated
// registered node exceeds the limitation, the oldest will be dropped.
func WriteRegisterNodes(db KeyValueStore, nodeType pb.RegisteredNodeType, registeredNode *pb.YarnRegisteredPeerDetail) {
	blob, err := db.Get(registryNodeKey(nodeType))
	if err != nil {
		log.Warnf("Failed to load old registered nodes, err: {%s}",  err)
	}
	var registeredNodes dbtype.RegisteredNodeListPB
	if len(blob) > 0 {
		if err := registeredNodes.Unmarshal(blob); err != nil {
			log.WithError(err).Fatal("Failed to decode old registered nodes")
		}

	}
	for _, s := range registeredNodes.GetRegisteredNodeList() {
		if strings.EqualFold(s.Id, registeredNode.Id) {
			log.WithFields(logrus.Fields{ "id": s.Id }).Info("Skip duplicated registered node")
			return
		}
	}
	registeredNodes.RegisteredNodeList = append(registeredNodes.RegisteredNodeList, &dbtype.RegisteredNodePB{
		Id:                   registeredNode.Id,
		InternalIp:           registeredNode.InternalIp,
		InternalPort:         registeredNode.InternalPort,
		ExternalIp: 		  registeredNode.ExternalIp,
		ExternalPort: 	      registeredNode.ExternalPort,
		ConnState:            int32(registeredNode.ConnState),
	})
	// max limit for store seed node.
	if len(registeredNodes.RegisteredNodeList) > registeredNodeToKeep {
		registeredNodes.RegisteredNodeList = registeredNodes.RegisteredNodeList[:registeredNodeToKeep]
	}
	data, err := registeredNodes.Marshal()
	if err != nil {
		log.WithError(err).Fatal("Failed to encode registered node")
	}
	if err := db.Put(registryNodeKey(nodeType), data); err != nil {
		log.WithError(err).Fatal("Failed to write registered node")
	}
}

func DeleteRegisterNode(db KeyValueStore, nodeType pb.RegisteredNodeType, id string) {
	blob, err := db.Get(registryNodeKey(nodeType))
	if err != nil {
		log.Warnf("Failed to load old registered nodes, err: {%s}", err)
	}
	var registeredNodes dbtype.RegisteredNodeListPB
	if len(blob) > 0 {
		if err := registeredNodes.Unmarshal(blob); err != nil {
			log.WithError(err).Fatal("Failed to decode old registered nodes")
		}
	}
	for i, s := range registeredNodes.GetRegisteredNodeList() {
		if strings.EqualFold(s.Id, id) {
			registeredNodes.RegisteredNodeList = append(registeredNodes.RegisteredNodeList[:i], registeredNodes.RegisteredNodeList[i+1:]...)
			break
		}
	}
	data, err := registeredNodes.Marshal()
	if err != nil {
		log.WithError(err).Fatal("Failed to encode registered node")
	}
	if err := db.Put(registryNodeKey(nodeType), data); err != nil {
		log.WithError(err).Fatal("Failed to write registered node")
	}
}

// DeleteRegisterNodes deletes all the registered nodes from the database.
func DeleteRegisterNodes(db DatabaseDeleter, nodeType pb.RegisteredNodeType) {
	if err := db.Delete(registryNodeKey(nodeType)); err != nil {
		log.WithError(err).Fatal("Failed to delete registered node")
	}
}

// ReadTaskEvent retrieves the evengine of task with the corresponding taskId.
func ReadTaskEvent(db DatabaseReader, taskId string) ([]*libTypes.TaskEvent, error) {
	blob, err := db.Get(taskEventKey(taskId))
	if err != nil {
		return nil, err
	}
	var events dbtype.TaskEventArrayPB
	if err := events.Unmarshal(blob); err != nil {
		return nil, err
	}
	resEvent := make([]*libTypes.TaskEvent, 0)
	for _, e := range events.GetTaskEventList() {
		if strings.EqualFold(e.GetTaskId(), taskId) {
			resEvent = append(resEvent, &libTypes.TaskEvent{
				Type:       e.GetType(),
				IdentityId:   e.GetIdentityId(),
				TaskId:     e.GetTaskId(),
				Content:    e.GetContent(),
				CreateAt: e.GetCreateAt(),
			})
		}
	}
	return resEvent, nil
}

// ReadAllTaskEvent retrieves the task event with all.
func ReadAllTaskEvents(db KeyValueStore) ( []*libTypes.TaskEvent, error) {
	prefix := taskEventPrefix
	it := db.NewIteratorWithPrefixAndStart(prefix, nil)
	defer it.Release()
	result := make([]*libTypes.TaskEvent, 0)

	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			blob, err := db.Get(key)
			if err != nil {
				continue
			}
			var events dbtype.TaskEventArrayPB
			if err := events.Unmarshal(blob); err != nil {
				continue
			}
			for _, e := range events.GetTaskEventList() {
				result = append(result, e)
			}
		}
	}
	return result, nil
}

// WriteTaskEvent serializes the task evengine into the database.
func WriteTaskEvent(db KeyValueStore, taskEvent *libTypes.TaskEvent) {
	blob, err := db.Get(taskEventKey(taskEvent.TaskId))
	if err != nil {
		log.WithError(err).Warn("Failed to load old task events")
	}
	var array dbtype.TaskEventArrayPB
	if len(blob) > 0 {
		if err := array.Unmarshal(blob); err != nil {
			log.WithError(err).Fatal("Failed to decode old task events")
		}
	}
	//for _, s := range array.GetTaskEventList() {
	//	if strings.EqualFold(s.TaskId, taskEvent.TaskId) &&
	//		strings.EqualFold(s.Identity, taskEvent.Identity) &&
	//		strings.EqualFold(s.EventContent, taskEvent.Content) {
	//		log.WithFields(logrus.Fields{ "identity": s.Identity, "taskId": s.TaskId }).Info("Skip duplicated task evengine")
	//		return
	//	}
	//}
	array.TaskEventList = append(array.TaskEventList, taskEvent)

	data, err := array.Marshal()
	if err != nil {
		log.WithError(err).Fatal("Failed to encode task events")
	}
	if err := db.Put(taskEventKey(taskEvent.TaskId), data); err != nil {
		log.WithError(err).Fatal("Failed to write task events")
	}
}

// DeleteTaskEvent deletes the task evengine from the database with a special taskId
func DeleteTaskEvent(db KeyValueStore, taskId string) {
	blob, err := db.Get(taskEventKey(taskId))
	if err != nil {
		log.Warnf("Failed to load old task evengine, err: {%s}", err)
	}
	var array dbtype.TaskEventArrayPB
	if len(blob) > 0 {
		if err := array.Unmarshal(blob); err != nil {
			log.WithError(err).Fatal("Failed to decode old task events")
		}
	}
	finalArray := new(dbtype.TaskEventArrayPB)
	for _, s := range array.GetTaskEventList() {
		if !strings.EqualFold(s.TaskId, taskId) {
			finalArray.TaskEventList = append(finalArray.TaskEventList, s)
		}
	}
	data, err := finalArray.Marshal()
	if err != nil {
		log.WithError(err).Fatal("Failed to encode task events")
	}
	if err := db.Put(taskEventKey(taskId), data); err != nil {
		log.WithError(err).Fatal("Failed to write task events")
	}
}

// ReadLocalResource retrieves the resource of local with the corresponding jobNodeId.
func ReadLocalResource(db DatabaseReader, jobNodeId string) (*types.LocalResource, error) {
	blob, err := db.Get(localResourceKey(jobNodeId))
	if err != nil {
		log.WithError(err).Fatal("Failed to read local resource")
		return nil, err
	}
	localResource := new(libTypes.LocalResourcePB)
	if err := localResource.Unmarshal(blob); err != nil {
		log.WithError(err).Fatal("Failed to unmarshal local resource")
		return nil, err
	}
	return types.NewLocalResource(localResource), nil
}

// ReadAllLocalResource retrieves the local resource with all.
func ReadAllLocalResource(db KeyValueStore) (types.LocalResourceArray, error) {
	prefix := localResourcePrefix
	it := db.NewIteratorWithPrefixAndStart(prefix, nil)
	defer it.Release()
	array := make(types.LocalResourceArray, 0)
	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			blob, err := db.Get(key)
			if err != nil {
				continue
			}
			localResource := new(libTypes.LocalResourcePB)
			if err := localResource.Unmarshal(blob); err != nil {
				continue
			}
			array = append(array, types.NewLocalResource(localResource))
		}
	}
	return array, nil
}

// WriteLocalResource serializes the local resource into the database.
func WriteLocalResource(db KeyValueStore, localResource *types.LocalResource) {
	buffer := new(bytes.Buffer)
	err := localResource.EncodePb(buffer)
	if err != nil {
		log.WithError(err).Fatal("Failed to encode local resource")
	}
	if err := db.Put(localResourceKey(localResource.GetJobNodeId()), buffer.Bytes()); err != nil {
		log.WithError(err).Fatal("Failed to write local resource")
	}
}

// DeleteLocalResource deletes the local resource from the database with a special jobNodeId
func DeleteLocalResource(db KeyValueStore, jobNodeId string) {
	if err := db.Delete(localResourceKey(jobNodeId)); err != nil {
		log.WithError(err).Fatal("Failed to delete local resource")
	}
}

// ReadLocalTask retrieves the local task with the corresponding taskId.
func ReadLocalTask(db DatabaseReader, taskId string) (*types.Task, error) {
	blob, err := db.Get(localTaskKey)
	if err != nil {
		return nil, err
	}
	var taskList dbtype.TaskArrayPB
	if err := taskList.Unmarshal(blob); err != nil {
		return nil, err
	}
	for _, task := range taskList.GetTaskList() {
		if strings.EqualFold(task.TaskId, taskId) {
			return types.NewTask(task), nil
		}
	}
	return nil, ErrNotFound
}

// ReadLocalTaskByIds retrieves the local tasks with the corresponding taskIds.
func ReadLocalTaskByIds(db DatabaseReader, taskIds []string) (types.TaskDataArray, error) {
	blob, err := db.Get(localTaskKey)
	if err != nil {
		return nil, err
	}
	var taskList dbtype.TaskArrayPB
	if err := taskList.Unmarshal(blob); err != nil {
		return nil, err
	}
	tmp := make(map[string]struct{}, 0)
	taskArr := make(types.TaskDataArray, 0)
	for _, taskId := range taskIds {
		tmp[taskId] = struct{}{}
	}
	for _, task := range taskList.GetTaskList() {
		if _, ok := tmp[task.TaskId]; ok {
			taskArr = append(taskArr, types.NewTask(task))
		}
	}
	return taskArr, nil
}

// ReadAllLocalTasks retrieves all the local task in the database.
func ReadAllLocalTasks(db DatabaseReader) (types.TaskDataArray, error) {
	blob, err := db.Get(localTaskKey)
	if err != nil {
		return nil, err
	}
	var array dbtype.TaskArrayPB
	if err := array.Unmarshal(blob); err != nil {
		return nil, err
	}
	result := make(types.TaskDataArray, 0)
	for _, task := range array.GetTaskList() {
		result = append(result, types.NewTask(task))
	}
	return result, nil
}

// WriteLocalTask serializes the local task into the database.
func WriteLocalTask(db KeyValueStore, task *types.Task) {
	blob, err := db.Get(localTaskKey)
	if err != nil {
		log.Warnf("Failed to load old local task, err: {%s}", err)
	}
	var array dbtype.TaskArrayPB
	if len(blob) > 0 {
		if err := array.Unmarshal(blob); err != nil {
			log.WithError(err).Fatal("Failed to decode old local task")
		}
	}

	var duplicated bool = false
	// Check whether there is duplicate data.
	for i, s := range array.TaskList {
		if strings.EqualFold(s.TaskId, task.TaskId()) {
			log.WithFields(logrus.Fields{ "taskId": s.TaskId }).Info("update duplicated local task")
			//return
			array.TaskList[i] = task.TaskData()
			duplicated = true
			break
		}
	}

	// add new task
	if !duplicated {
		array.TaskList = append(array.TaskList, task.TaskData())
	}

	data, err := array.Marshal()
	if err != nil {
		log.WithError(err).Fatal("Failed to encode local task")
	}
	if err := db.Put(localTaskKey, data); err != nil {
		log.WithError(err).Fatal("Failed to write local task")
	}
}

// DeleteLocalTask deletes the local task from the database with a special taskId
func DeleteLocalTask(db KeyValueStore, taskId string) {
	blob, err := db.Get(localTaskKey)
	if err != nil {
		log.Warnf("Failed to load old local task, err: {%s}", err)
	}
	var array dbtype.TaskArrayPB
	if len(blob) > 0 {
		if err := array.Unmarshal(blob); err != nil {
			log.WithError(err).Fatal("Failed to decode old local task")
		}
	}
	// need to test.
	for idx, s := range array.TaskList {
		if strings.EqualFold(s.TaskId, taskId) {
			array.TaskList = append(array.TaskList[:idx], array.TaskList[idx+1:]...)
			break
		}
	}
	data, err := array.Marshal()
	if err != nil {
		log.WithError(err).Fatal("Failed to encode local task")
	}
	if err := db.Put(localTaskKey, data); err != nil {
		log.WithError(err).Fatal("Failed to write local task")
	}
}

// DeleteLocalTask deletes all the local task from the database.
func DeleteLocalTasks(db DatabaseDeleter) {
	if err := db.Delete(localTaskKey); err != nil {
		log.WithError(err).Fatal("Failed to delete local task with all")
	}
}