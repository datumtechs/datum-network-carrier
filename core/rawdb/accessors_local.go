// Copyright (C) 2021 The RosettaNet Authors.

package rawdb

import (
	"bytes"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	dbtype "github.com/RosettaFlow/Carrier-Go/lib/db"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/sirupsen/logrus"
	"strings"
)

const seedNodeToKeep = 50
const registeredNodeToKeep = 50

// QueryLocalIdentity retrieves the identity of local.
func QueryLocalIdentity(db DatabaseReader) (*apicommonpb.Organization, error) {
	var blob apicommonpb.Organization
	enc, err := db.Get(localIdentityKey)
	if nil != err {
		return nil, err
	}
	if err := blob.Unmarshal(enc); nil != err {
		return nil, err
	}
	return &apicommonpb.Organization{
		NodeName:   blob.GetNodeName(),
		NodeId:     blob.GetNodeId(),
		IdentityId: blob.GetIdentityId(),
	}, nil
}

// StoreLocalIdentity stores the local identity.
func StoreLocalIdentity(db DatabaseWriter, localIdentity *apicommonpb.Organization) error {
	pb := &apicommonpb.Organization{
		IdentityId: localIdentity.GetIdentityId(),
		NodeId:     localIdentity.GetNodeId(),
		NodeName:   localIdentity.GetNodeName(),
	}
	enc, err := pb.Marshal()
	if nil != err {
		return err
	}
	if err := db.Put(localIdentityKey, enc); nil != err {
		log.WithError(err).Error("Failed to store local identity")
		return err
	}
	return nil
}

// RemoveLocalIdentity deletes the local identity
func RemoveLocalIdentity(db KeyValueStore) error {

	has, err := db.Has(localIdentityKey)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(localIdentityKey)
}

// QueryYarnName retrieves the name of yarn.
func QueryYarnName(db DatabaseReader) (string, error) {
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

// StoreYarnName stores the name of yarn.
func StoreYarnName(db DatabaseWriter, yarnName string) error {
	pb := dbtype.StringPB{
		V: yarnName,
	}
	enc, err := pb.Marshal()
	if nil != err {
		return err
	}
	if err := db.Put(yarnNameKey, enc); nil != err {
		log.WithError(err).Error("Failed to store yarn name")
		return err
	}
	return err
}

// RemoveYarnName deletes the name of yarn.
func RemoveYarnName(db KeyValueStore) error {

	has, err := db.Has(yarnNameKey)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(yarnNameKey)
}

// QuerySeedNode retrieves the seed node with the corresponding nodeId.
func QuerySeedNode(db DatabaseReader, nodeId string) (*pb.SeedPeer, error) {
	blob, err := db.Get(seedNodeKey)
	if nil != err {
		return nil, err
	}
	var seedNodes dbtype.SeedPeerListPB
	if err := seedNodes.Unmarshal(blob); nil != err {
		return nil, err
	}
	for _, seed := range seedNodes.GetSeedPeerList() {
		if strings.EqualFold(seed.Id, nodeId) {
			return &pb.SeedPeer{
				Id:           seed.Id,
				NodeId:       seed.NodeId,
				InternalIp:   seed.InternalIp,
				InternalPort: seed.InternalPort,
				ConnState:    pb.ConnState(seed.ConnState),
			}, nil
		}
	}
	return nil, ErrNotFound
}

// QueryAllSeedNodes retrieves all the seed nodes in the database.
// All returned seed nodes are sorted in reverse.
func QueryAllSeedNodes(db DatabaseReader) ([]*pb.SeedPeer, error) {
	blob, err := db.Get(seedNodeKey)
	if nil != err {
		return nil, err
	}
	var seedNodes dbtype.SeedPeerListPB
	if err := seedNodes.Unmarshal(blob); nil != err {
		return nil, err
	}
	var nodes []*pb.SeedPeer
	for _, seed := range seedNodes.SeedPeerList {
		nodes = append(nodes, &pb.SeedPeer{
			Id:           seed.Id,
			NodeId:   	  seed.NodeId,
			InternalIp:   seed.InternalIp,
			InternalPort: seed.InternalPort,
			ConnState:    pb.ConnState(seed.ConnState),
		})
	}
	return nodes, nil
}

// StoreSeedNode serializes the seed node into the database. If the cumulated
// seed node exceeds the limitation, the oldest will be dropped.
func StoreSeedNode(db KeyValueStore, seedNode *pb.SeedPeer) error {
	blob, err := db.Get(seedNodeKey)
	if IsNoDBNotFoundErr(err) {
		log.WithError(err).Error("Failed to load old seed nodes")
		return err
	}
	var seedNodes dbtype.SeedPeerListPB
	if len(blob) > 0 {
		if err := seedNodes.Unmarshal(blob); nil != err {
			log.WithError(err).Error("Failed to decode old seed nodes")
			return err
		}

	}
	for _, s := range seedNodes.GetSeedPeerList() {
		if strings.EqualFold(s.Id, seedNode.Id) {
			log.WithFields(logrus.Fields{"id": s.Id}).Info("Skip duplicated seed node")
			return nil
		}
	}
	seedNodes.SeedPeerList = append(seedNodes.SeedPeerList, &dbtype.SeedPeerPB{
		Id:           seedNode.Id,
		NodeId:       seedNode.NodeId,
		InternalIp:   seedNode.InternalIp,
		InternalPort: seedNode.InternalPort,
		ConnState:    int32(seedNode.ConnState),
	})
	// max limit for store seed node.
	if len(seedNodes.SeedPeerList) > seedNodeToKeep {
		seedNodes.SeedPeerList = seedNodes.SeedPeerList[:seedNodeToKeep]
	}
	data, err := seedNodes.Marshal()
	if nil != err {
		log.WithError(err).Error("Failed to encode seed node")
		return err
	}
	if err := db.Put(seedNodeKey, data); nil != err {
		log.WithError(err).Error("Failed to write seed node")
		return err
	}
	return nil
}

// RemoveSeedNode deletes the seed nodes from the database with a special id
func RemoveSeedNode(db KeyValueStore, id string) error {
	blob, err := db.Get(seedNodeKey)
	switch {
	case IsNoDBNotFoundErr(err):
		log.WithError(err).Error("Failed to load old seed nodes")
		return err
	case IsDBNotFoundErr(err), nil == err && len(blob) == 0:
		return nil
	}
	var seedNodes dbtype.SeedPeerListPB
	if len(blob) > 0 {
		if err := seedNodes.Unmarshal(blob); nil != err {
			log.WithError(err).Error("Failed to decode old seed nodes")
			return err
		}
	}
	// need to test.
	for idx, s := range seedNodes.GetSeedPeerList() {
		if strings.EqualFold(s.Id, id) {
			seedNodes.SeedPeerList = append(seedNodes.SeedPeerList[:idx], seedNodes.SeedPeerList[idx+1:]...)
			break
		}
	}

	if len(seedNodes.GetSeedPeerList()) == 0 {
		return db.Delete(seedNodeKey)
	}

	data, err := seedNodes.Marshal()
	if nil != err {
		log.WithError(err).Fatal("Failed to encode seed nodes")
	}
	if err := db.Put(seedNodeKey, data); nil != err {
		log.WithError(err).Error("Failed to write seed nodes")
		return err
	}
	return nil
}

// RemoveSeedNodes deletes all the seed nodes from the database
func RemoveSeedNodes(db KeyValueStore) error {

	has, err := db.Has(seedNodeKey)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(seedNodeKey)
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

// QueryRegisterNode retrieves the register node with the corresponding nodeId.
func QueryRegisterNode(db DatabaseReader, nodeType pb.RegisteredNodeType, nodeId string) (*pb.YarnRegisteredPeerDetail, error) {
	blob, err := db.Get(registryNodeKey(nodeType))
	if nil != err {
		return nil, err
	}
	var registeredNodes dbtype.RegisteredNodeListPB
	if err := registeredNodes.Unmarshal(blob); nil != err {
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

// QueryAllRegisterNodes retrieves all the registered nodes in the database.
// All returned registered nodes are sorted in reverse.
func QueryAllRegisterNodes(db DatabaseReader, nodeType pb.RegisteredNodeType) ([]*pb.YarnRegisteredPeerDetail, error) {
	blob, err := db.Get(registryNodeKey(nodeType))
	if nil != err {
		return nil, err
	}
	var registeredNodes dbtype.RegisteredNodeListPB
	if err := registeredNodes.Unmarshal(blob); nil != err {
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

// StoreRegisterNode serializes the registered node into the database. If the cumulated
// registered node exceeds the limitation, the oldest will be dropped.
func StoreRegisterNode (db KeyValueStore, nodeType pb.RegisteredNodeType, registeredNode *pb.YarnRegisteredPeerDetail) error {

	key := registryNodeKey(nodeType)

	blob, err := db.Get(key)
	if IsNoDBNotFoundErr(err) {
		log.WithError(err).Error("Failed to load old registered nodes")
		return err
	}
	var registeredNodes dbtype.RegisteredNodeListPB
	if len(blob) > 0 {
		if err := registeredNodes.Unmarshal(blob); nil != err {
			log.WithError(err).Error("Failed to decode old registered nodes")
			return err
		}

	}
	for _, s := range registeredNodes.GetRegisteredNodeList() {
		if strings.EqualFold(s.Id, registeredNode.Id) {
			log.WithFields(logrus.Fields{"id": s.Id}).Info("Skip duplicated registered node")
			return nil
		}
	}
	registeredNodes.RegisteredNodeList = append(registeredNodes.RegisteredNodeList, &dbtype.RegisteredNodePB{
		Id:           registeredNode.Id,
		InternalIp:   registeredNode.InternalIp,
		InternalPort: registeredNode.InternalPort,
		ExternalIp:   registeredNode.ExternalIp,
		ExternalPort: registeredNode.ExternalPort,
		ConnState:    int32(registeredNode.ConnState),
	})
	// max limit for store seed node.
	if len(registeredNodes.RegisteredNodeList) > registeredNodeToKeep {
		registeredNodes.RegisteredNodeList = registeredNodes.RegisteredNodeList[:registeredNodeToKeep]
	}
	data, err := registeredNodes.Marshal()
	if nil != err {
		log.WithError(err).Error("Failed to encode registered node")
		return err
	}
	if err := db.Put(key, data); nil != err {
		log.WithError(err).Error("Failed to write registered node")
		return err
	}
	return nil
}

func RemoveRegisterNode(db KeyValueStore, nodeType pb.RegisteredNodeType, id string) error {

	key := registryNodeKey(nodeType)
	blob, err := db.Get(key)

	switch {
	case IsNoDBNotFoundErr(err):
		log.WithError(err).Error("Failed to load old registered nodes")
		return err
	case IsDBNotFoundErr(err), nil == err && len(blob) == 0:
		return nil
	}

	var registeredNodes dbtype.RegisteredNodeListPB
	if len(blob) > 0 {
		if err := registeredNodes.Unmarshal(blob); nil != err {
			log.WithError(err).Fatal("Failed to decode old registered nodes")
		}
	}
	for i, s := range registeredNodes.GetRegisteredNodeList() {
		if strings.EqualFold(s.Id, id) {
			registeredNodes.RegisteredNodeList = append(registeredNodes.RegisteredNodeList[:i], registeredNodes.RegisteredNodeList[i+1:]...)
			break
		}
	}

	if len(registeredNodes.GetRegisteredNodeList()) == 0 {
		return db.Delete(key)
	}

	data, err := registeredNodes.Marshal()
	if nil != err {
		log.WithError(err).Fatal("Failed to encode registered node")
	}
	if err := db.Put(key, data); nil != err {
		log.WithError(err).Fatal("Failed to write registered node")
	}
	return nil
}

// RemoveRegisterNodes deletes all the registered nodes from the database.
func RemoveRegisterNodes(db KeyValueStore, nodeType pb.RegisteredNodeType) error {

	key := registryNodeKey(nodeType)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

// QueryTaskEvent retrieves the evengine of task with the corresponding taskId for all partyIds.
func QueryTaskEvent(db KeyValueStore, taskId string) ([]*libtypes.TaskEvent, error) {

	it := db.NewIteratorWithPrefixAndStart(append(taskEventPrefix, []byte(taskId)...), nil)
	defer it.Release()

	result := make([]*libtypes.TaskEvent, 0)

	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			blob, err := db.Get(key)
			if nil != err {
				continue
			}
			var events dbtype.TaskEventArrayPB
			if err := events.Unmarshal(blob); nil != err {
				continue
			}
			if len(events.GetTaskEventList()) != 0 {
				result = append(result, events.GetTaskEventList()...)
			}
		}
	}
	return result, nil
}

// QueryTaskEventByPartyId retrieves the events of task with the corresponding taskId and partyId.
func QueryTaskEventByPartyId (db DatabaseReader, taskId, partyId string) ([]*libtypes.TaskEvent, error) {

	key := taskEventKey(taskId, partyId)

	val, err := db.Get(key)
	if nil != err {
		return nil, err
	}
	var events dbtype.TaskEventArrayPB
	if err := events.Unmarshal(val); nil != err {
		log.WithError(err).Errorf("Failed to encode task events")
		return nil, err
	}
	return events.GetTaskEventList(), nil
}

// QueryAllTaskEvents retrieves the task event with all (all taskIds and all partyIds).
func QueryAllTaskEvents (db KeyValueStore) ([]*libtypes.TaskEvent, error) {

	it := db.NewIteratorWithPrefixAndStart(taskEventPrefix, nil)
	defer it.Release()

	result := make([]*libtypes.TaskEvent, 0)

	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			blob, err := db.Get(key)
			if nil != err {
				continue
			}
			var events dbtype.TaskEventArrayPB
			if err := events.Unmarshal(blob); nil != err {
				continue
			}
			if len(events.GetTaskEventList()) != 0 {
				result = append(result, events.GetTaskEventList()...)
			}
		}
	}
	return result, nil
}

// StoreTaskEvent serializes the task evengine into the database.
func StoreTaskEvent(db KeyValueStore, taskEvent *libtypes.TaskEvent) error {
	key := taskEventKey(taskEvent.GetTaskId(), taskEvent.GetPartyId())
	val, err := db.Get(key)
	if IsNoDBNotFoundErr(err) {
		log.WithError(err).Error("Failed to load old task events")
		return err
	}
	var array dbtype.TaskEventArrayPB
	if len(val) > 0 {
		if err := array.Unmarshal(val); nil != err {
			log.WithError(err).Errorf("Failed to decode old task events")
			return err
		}
	}

	array.TaskEventList = append(array.TaskEventList, taskEvent)
	data, err := array.Marshal()
	if nil != err {
		log.WithError(err).Errorf("Failed to encode task events")
		return err
	}
	if err := db.Put(key, data); nil != err {
		log.WithError(err).Errorf("Failed to write task events")
		return err
	}
	return nil
}

// RemoveTaskEvent remove the task events from the database with a special taskId for all partyIds
func RemoveTaskEvent(db KeyValueStore, taskId string) error {

	it := db.NewIteratorWithPrefixAndStart(append(taskEventPrefix, []byte(taskId)...), nil)
	defer it.Release()
	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			has, err := db.Has(key)
			switch {
			case IsNoDBNotFoundErr(err):
				return err
			case IsDBNotFoundErr(err), nil == err && !has:
				continue
			}
			if err := db.Delete(key); nil != err {
				return err
			}
		}
	}
	return nil
}

// RemoveTaskEventByPartyId remove the task events from the database with a special taskId and partyId
func RemoveTaskEventByPartyId(db KeyValueStore, taskId, partyId string) error {
	key := taskEventKey(taskId, partyId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

// QueryLocalResource retrieves the resource of local with the corresponding jobNodeId.
func QueryLocalResource(db DatabaseReader, jobNodeId string) (*types.LocalResource, error) {
	blob, err := db.Get(localResourceKey(jobNodeId))
	if nil != err {
		log.WithError(err).Errorf("Failed to read local resource")
		return nil, err
	}
	localResource := new(libtypes.LocalResourcePB)
	if err := localResource.Unmarshal(blob);nil != err {
		log.WithError(err).Errorf("Failed to unmarshal local resource")
		return nil, err
	}
	return types.NewLocalResource(localResource), nil
}

// QueryAllLocalResource retrieves the local resource with all.
func QueryAllLocalResource(db KeyValueStore) (types.LocalResourceArray, error) {
	prefix := localResourcePrefix
	it := db.NewIteratorWithPrefixAndStart(prefix, nil)
	defer it.Release()
	array := make(types.LocalResourceArray, 0)
	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			blob, err := db.Get(key)
			if nil != err {
				continue
			}
			localResource := new(libtypes.LocalResourcePB)
			if err := localResource.Unmarshal(blob); nil != err {
				continue
			}
			array = append(array, types.NewLocalResource(localResource))
		}
	}
	return array, nil
}

// StoreLocalResource serializes the local resource into the database.
func StoreLocalResource(db KeyValueStore, localResource *types.LocalResource) error {
	buffer := new(bytes.Buffer)
	err := localResource.EncodePb(buffer)
	if nil != err {
		log.WithError(err).Errorf("Failed to encode local resource")
		return err
	}
	if err := db.Put(localResourceKey(localResource.GetJobNodeId()), buffer.Bytes()); nil != err {
		log.WithError(err).Errorf("Failed to write local resource")
		return err
	}
	return nil
}

// RemoveLocalResource deletes the local resource from the database with a special jobNodeId
func RemoveLocalResource(db KeyValueStore, jobNodeId string) error {
	key := localResourceKey(jobNodeId)

	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

// QueryLocalTask retrieves the local task with the corresponding taskId.
func QueryLocalTask(db DatabaseReader, taskId string) (*types.Task, error) {
	blob, err := db.Get(localTaskKey)
	if nil != err {
		return nil, err
	}
	var taskList dbtype.TaskArrayPB
	if err := taskList.Unmarshal(blob); nil != err {
		return nil, err
	}
	for _, task := range taskList.GetTaskList() {
		if strings.EqualFold(task.TaskId, taskId) {
			return types.NewTask(task), nil
		}
	}
	return nil, ErrNotFound
}

// QueryLocalTaskByIds retrieves the local tasks with the corresponding taskIds.
func QueryLocalTaskByIds(db DatabaseReader, taskIds []string) (types.TaskDataArray, error) {
	blob, err := db.Get(localTaskKey)
	if nil != err {
		return nil, err
	}
	var taskList dbtype.TaskArrayPB
	if err := taskList.Unmarshal(blob); nil != err {
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

// QueryAllLocalTasks retrieves all the local task in the database.
func QueryAllLocalTasks(db DatabaseReader) (types.TaskDataArray, error) {
	blob, err := db.Get(localTaskKey)
	if nil != err {
		return nil, err
	}
	var array dbtype.TaskArrayPB
	if err := array.Unmarshal(blob); nil != err {
		return nil, err
	}
	result := make(types.TaskDataArray, 0)
	for _, task := range array.GetTaskList() {
		result = append(result, types.NewTask(task))
	}
	return result, nil
}

// StoreScheduling save scheduled tasks.
func StoreScheduling(db KeyValueStore, bullet *types.TaskBullet) error {
	val, err := rlp.EncodeToBytes(bullet)
	if nil != err {
		log.WithError(err).Errorf("Failed to encode scheduling task.")
		return err
	}
	if err := db.Put(schedulingKey(bullet.GetTaskId()), val); nil != err {
		log.WithError(err).Error("Failed to write scheduling task.")
		return err
	}
	return nil
}

func RemoveScheduling(db KeyValueStore, bullet *types.TaskBullet) error {
	key := schedulingKey(bullet.TaskId)

	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func RecoveryScheduling(db KeyValueStore) (*types.TaskBullets, *types.TaskBullets, map[string]*types.TaskBullet) {
	it := db.NewIteratorWithPrefixAndStart(schedulingPrefix, nil)
	defer it.Release()
	queue := make(types.TaskBullets,0)
	starveQueue := make(types.TaskBullets,0)
	schedulings := make(map[string]*types.TaskBullet, 0)
	for it.Next() {
		var result types.TaskBullet
		key := it.Key()
		taskId := string(key[len(schedulingPrefix):])

		value := it.Value()
		if err := rlp.DecodeBytes(value, &result); nil != err {
			log.Warning("RecoveryScheduling DecodeBytes fail,error info:", err)
		}
		schedulings[taskId] = &result
		if result.Starve == true {
			starveQueue.Push(result)
		} else {
			queue.Push(result)
		}
	}
	return &queue, &starveQueue, schedulings
}

// StoreLocalTask serializes the local task into the database.
func StoreLocalTask(db KeyValueStore, task *types.Task) error {
	blob, err := db.Get(localTaskKey)
	if IsNoDBNotFoundErr(err) {
		log.WithError(err).Error("Failed to load old local task")
		return err
	}
	var array dbtype.TaskArrayPB
	if len(blob) > 0 {
		if err := array.Unmarshal(blob); nil != err {
			log.WithError(err).Fatal("Failed to decode old local task")
		}
	}

	var duplicated bool
	// Check whether there is duplicate data.
	for i, s := range array.TaskList {
		if strings.EqualFold(s.TaskId, task.GetTaskId()) {
			log.WithFields(logrus.Fields{"taskId": s.TaskId}).Info("update duplicated local task")
			//return
			array.TaskList[i] = task.GetTaskData()
			duplicated = true
			break
		}
	}

	// add new task
	if !duplicated {
		array.TaskList = append(array.TaskList, task.GetTaskData())
	}

	data, err := array.Marshal()
	if nil != err {
		log.WithError(err).Error("Failed to encode local task")
		return err
	}
	if err := db.Put(localTaskKey, data); nil != err {
		log.WithError(err).Error("Failed to write local task")
		return err
	}
	return nil
}

// RemoveLocalTask deletes the local task from the database with a special taskId
func RemoveLocalTask(db KeyValueStore, taskId string) error {
	blob, err := db.Get(localTaskKey)

	switch {
	case IsNoDBNotFoundErr(err):
		log.WithError(err).Error("Failed to load old local task")
		return err
	case IsDBNotFoundErr(err), nil == err && len(blob) == 0:
		return nil
	}

	var array dbtype.TaskArrayPB
	if len(blob) > 0 {
		if err := array.Unmarshal(blob); nil != err {
			log.WithError(err).Error("Failed to decode old local task")
			return err
		}
	}
	// need to test.
	for idx, s := range array.GetTaskList() {
		if strings.EqualFold(s.TaskId, taskId) {
			array.TaskList = append(array.TaskList[:idx], array.TaskList[idx+1:]...)
			break
		}
	}

	if len(array.GetTaskList()) == 0 {
		return db.Delete(localTaskKey)
	}

	data, err := array.Marshal()
	if nil != err {
		log.WithError(err).Error("Failed to encode local task")
		return err
	}

	if err := db.Put(localTaskKey, data); nil != err {
		log.WithError(err).Error("Failed to write local task")
		return err
	}
	return nil
}

// RemoveLocalAllTask deletes all the local task from the database.
func RemoveLocalAllTask(db KeyValueStore) error {

	has, err := db.Has(localTaskKey)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(localTaskKey)
}

// ================================= Local Metadata ==========================================
// QueryLocalMetadata retrieves the metadata of local with the corresponding metadataId.
func QueryLocalMetadata(db DatabaseReader, metadataId string) (*types.Metadata, error) {
	blob, err := db.Get(localMetadataKey(metadataId))
	if nil != err {
		log.WithError(err).Error("Failed to read local metadata")
		return nil, err
	}
	localMetadata := new(libtypes.MetadataPB)
	if err := localMetadata.Unmarshal(blob); nil != err {
		log.WithError(err).Error("Failed to unmarshal local metadata")
		return nil, err
	}
	return types.NewMetadata(localMetadata), nil
}

// QueryAllLocalMetadata retrieves the local metadata with all.
func QueryAllLocalMetadata(db KeyValueStore) (types.MetadataArray, error) {

	it := db.NewIteratorWithPrefixAndStart(localMetadataPrefix, nil)
	defer it.Release()
	array := make(types.MetadataArray, 0)
	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			blob, err := db.Get(key)
			if nil != err {
				continue
			}
			localMetadata := new(libtypes.MetadataPB)
			if err := localMetadata.Unmarshal(blob); nil != err {
				continue
			}
			array = append(array, types.NewMetadata(localMetadata))
		}
	}
	return array, nil
}

// StoreLocalMetadata serializes the local metadata into the database.
func StoreLocalMetadata(db KeyValueStore, localMetadata *types.Metadata) error {
	buffer := new(bytes.Buffer)
	err := localMetadata.EncodePb(buffer)
	if nil != err {
		log.WithError(err).Error("Failed to encode local metadata")
		return err
	}
	if err := db.Put(localMetadataKey(localMetadata.GetData().GetMetadataId()), buffer.Bytes()); nil != err {
		log.WithError(err).Error("Failed to write local metadata")
		return err
	}
	return nil
}

// RemoveLocalMetadata deletes the local metadata from the database with a special metadataId
func RemoveLocalMetadata(db KeyValueStore, metadataId string) error {
	key := localMetadataKey(metadataId)

	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}
