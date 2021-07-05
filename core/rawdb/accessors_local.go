// Copyright (C) 2021 The RosettaNet Authors.

package rawdb

import (
	"github.com/RosettaFlow/Carrier-Go/event"
	dbtype "github.com/RosettaFlow/Carrier-Go/lib/db"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/sirupsen/logrus"
	"strings"
	"sync/atomic"
)

const seedNodeToKeep = 50
const registeredNodeToKeep = 50

// ReadLocalIdentity retrieves the identity of local.
func ReadLocalIdentity(db DatabaseReader) *types.NodeAlias {
	var blob libtypes.OrganizationData
	enc, _ := db.Get(localIdentityKey)
	blob.Unmarshal(enc)
	return &types.NodeAlias{
		Name:       blob.GetNodeName(),
		NodeId:     blob.GetNodeId(),
		IdentityId: blob.GetIdentity(),
	}
}

// WriteLocalIdentity stores the local identity.
func WriteLocalIdentity(db DatabaseWriter, localIdentity *types.NodeAlias) {
	pb := &libtypes.OrganizationData{
		Alias:                localIdentity.GetNodeName(),
		Identity:             localIdentity.GetNodeIdentityId(),
		NodeId:               localIdentity.GetNodeIdStr(),
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
func ReadRunningTaskIDList(db DatabaseReader, jobNodeId string) []string {
	blob, _ := db.Get(runningTaskIDListKey(jobNodeId))
	var array dbtype.StringArrayPB
	if len(blob) > 0 {
		if err := array.Unmarshal(blob); err != nil {
			log.WithError(err).Fatal("Failed to decode old RunningTaskIdList")
		}
		return array.GetArray()
	}
	return nil
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
func ReadYarnName(db DatabaseReader) string {
	var yarnName dbtype.StringPB
	enc, _ := db.Get(yarnNameKey)
	yarnName.Unmarshal(enc)
	return yarnName.GetV()
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
func ReadSeedNode(db DatabaseReader, nodeId string) *types.SeedNodeInfo {
	blob, err := db.Get(seedNodeKey)
	if err != nil {
		return nil
	}
	var seedNodes dbtype.SeedNodeListPB
	if err := seedNodes.Unmarshal(blob); err != nil {
		return nil
	}
	for _, seed := range seedNodes.GetSeedNodeList() {
		if strings.EqualFold(seed.Id, nodeId) {
			return &types.SeedNodeInfo{
				Id:           seed.Id,
				InternalIp:   seed.InternalIp,
				InternalPort: seed.InternalPort,
				ConnState:    types.NodeConnStatus(seed.ConnState),
			}
		}
	}
	return nil
}

// ReadAllSeedNodes retrieves all the seed nodes in the database.
// All returned seed nodes are sorted in reverse.
func ReadAllSeedNodes(db DatabaseReader) []*types.SeedNodeInfo {
	blob, err := db.Get(seedNodeKey)
	if err != nil {
		return nil
	}
	var seedNodes dbtype.SeedNodeListPB
	if err := seedNodes.Unmarshal(blob); err != nil {
		return nil
	}
	var nodes []*types.SeedNodeInfo
	for _, seed := range seedNodes.SeedNodeList {
		nodes = append(nodes, &types.SeedNodeInfo{
			Id:           seed.Id,
			InternalIp:   seed.InternalIp,
			InternalPort: seed.InternalPort,
			ConnState:    types.NodeConnStatus(seed.ConnState),
		})
	}
	return nodes
}

// WriteSeedNodes serializes the seed node into the database. If the cumulated
// seed node exceeds the limitation, the oldest will be dropped.
func WriteSeedNodes(db KeyValueStore, seedNode *types.SeedNodeInfo) {
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

func registryNodeKey(nodeType types.RegisteredNodeType) []byte {
	var key []byte
	if nodeType == types.PREFIX_TYPE_JOBNODE {
		key = calcNodeKey
	}
	if nodeType == types.PREFIX_TYPE_DATANODE {
		key = dataNodeKey
	}
	return key
}

// ReadRegisterNode retrieves the register node with the corresponding nodeId.
func ReadRegisterNode(db DatabaseReader, nodeType types.RegisteredNodeType, nodeId string) *types.RegisteredNodeInfo {
	blob, err := db.Get(registryNodeKey(nodeType))
	if err != nil {
		return nil
	}
	var registeredNodes dbtype.RegisteredNodeListPB
	if err := registeredNodes.Unmarshal(blob); err != nil {
		return nil
	}
	for _, registered := range registeredNodes.GetRegisteredNodeList() {
		if strings.EqualFold(registered.Id, nodeId) {
			return &types.RegisteredNodeInfo{
				Id:           registered.Id,
				InternalIp:   registered.InternalIp,
				InternalPort: registered.InternalPort,
				ExternalIp:   registered.ExternalIp,
				ExternalPort: registered.ExternalPort,
				ConnState:    types.NodeConnStatus(registered.ConnState),
			}
		}
	}
	return nil
}

// ReadAllRegisterNodes retrieves all the registered nodes in the database.
// All returned registered nodes are sorted in reverse.
func ReadAllRegisterNodes(db DatabaseReader, nodeType types.RegisteredNodeType) []*types.RegisteredNodeInfo {
	blob, err := db.Get(registryNodeKey(nodeType))
	if err != nil {
		return nil
	}
	var registeredNodes dbtype.RegisteredNodeListPB
	if err := registeredNodes.Unmarshal(blob); err != nil {
		return nil
	}
	var nodes []*types.RegisteredNodeInfo
	for _, registered := range registeredNodes.GetRegisteredNodeList() {
		nodes = append(nodes, &types.RegisteredNodeInfo{
			Id:           registered.Id,
			InternalIp:   registered.InternalIp,
			InternalPort: registered.InternalPort,
			ExternalIp:   registered.ExternalIp,
			ExternalPort: registered.ExternalPort,
			ConnState:    types.NodeConnStatus(registered.ConnState),
		})
	}
	return nodes
}

// WriteRegisterNodes serializes the registered node into the database. If the cumulated
// registered node exceeds the limitation, the oldest will be dropped.
func WriteRegisterNodes(db KeyValueStore, nodeType types.RegisteredNodeType, registeredNode *types.RegisteredNodeInfo) {
	blob, err := db.Get(registryNodeKey(nodeType))
	if err != nil {
		log.Warn("Failed to load old seed nodes", "error", err)
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

func DeleteRegisterNode(db KeyValueStore, nodeType types.RegisteredNodeType, id string) {
	blob, err := db.Get(registryNodeKey(nodeType))
	if err != nil {
		log.Warn("Failed to load old registered nodes", "error", err)
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
func DeleteRegisterNodes(db DatabaseDeleter, nodeType types.RegisteredNodeType) {
	if err := db.Delete(registryNodeKey(nodeType)); err != nil {
		log.WithError(err).Fatal("Failed to delete registered node")
	}
}

// ReadTaskEvent retrieves the event of task with the corresponding taskId.
func ReadTaskEvent(db DatabaseReader, taskId string) []*event.TaskEvent {
	blob, err := db.Get(taskEventKey)
	if err != nil {
		return nil
	}
	var events dbtype.TaskEventArrayPB
	if err := events.Unmarshal(blob); err != nil {
		return nil
	}
	resEvent := make([]*event.TaskEvent, 0)
	for _, e := range events.GetTaskEventList() {
		if strings.EqualFold(e.GetTaskId(), taskId) {
			resEvent = append(resEvent, &event.TaskEvent{
				Type:       e.GetEventType(),
				Identity:   e.GetIdentity(),
				TaskId:     e.GetTaskId(),
				Content:    e.GetEventContent(),
				CreateTime: e.GetEventAt(),
			})
		}
	}
	return resEvent
}

// ReadAllTaskEvents retrieves all the task event in the database.
// All returned task events are sorted in reverse.
func ReadAllTaskEvents(db DatabaseReader) []*event.TaskEvent {
	blob, err := db.Get(taskEventKey)
	if err != nil {
		return nil
	}
	var taskEvents dbtype.TaskEventArrayPB
	if err := taskEvents.Unmarshal(blob); err != nil {
		return nil
	}
	var events []*event.TaskEvent
	for _, e := range taskEvents.GetTaskEventList() {
		events = append(events, &event.TaskEvent{
			Type:       e.GetEventType(),
			Identity:   e.GetIdentity(),
			TaskId:     e.GetTaskId(),
			Content:    e.GetEventContent(),
			CreateTime: e.GetEventAt(),
		})
	}
	return events
}

// WriteTaskEvent serializes the task event into the database.
func WriteTaskEvent(db KeyValueStore, taskEvent *event.TaskEvent) {
	blob, err := db.Get(taskEventKey)
	if err != nil {
		log.Warn("Failed to load old task event", "error", err)
	}
	var array dbtype.TaskEventArrayPB
	if len(blob) > 0 {
		if err := array.Unmarshal(blob); err != nil {
			log.WithError(err).Fatal("Failed to decode old task event")
		}

	}
	for _, s := range array.GetTaskEventList() {
		if strings.EqualFold(s.GetTaskId(), taskEvent.TaskId) && strings.EqualFold(s.GetIdentity(), taskEvent.Identity) {
			log.WithFields(logrus.Fields{ "identity": s.Identity }).Info("Skip duplicated task event")
			return
		}
	}
	array.TaskEventList = append(array.TaskEventList, &libtypes.EventData{
		TaskId:               taskEvent.TaskId,
		EventType:            taskEvent.Type,
		EventAt:              taskEvent.CreateTime,
		EventContent:         taskEvent.Content,
		Identity:             taskEvent.Identity,
	})

	data, err := array.Marshal()
	if err != nil {
		log.WithError(err).Fatal("Failed to encode task event")
	}
	if err := db.Put(taskEventKey, data); err != nil {
		log.WithError(err).Fatal("Failed to write task event")
	}
}

// DeleteTaskEvent deletes the task event from the database with a special taskId
func DeleteTaskEvent(db KeyValueStore, taskId string) {
	blob, err := db.Get(taskEventKey)
	if err != nil {
		log.Warn("Failed to load old task event", "error", err)
	}
	var array dbtype.TaskEventArrayPB
	if len(blob) > 0 {
		if err := array.Unmarshal(blob); err != nil {
			log.WithError(err).Fatal("Failed to decode old task event")
		}
	}
	for idx, s := range array.GetTaskEventList() {
		if strings.EqualFold(s.TaskId, taskId) {
			array.TaskEventList = append(array.TaskEventList[:idx], array.TaskEventList[idx+1:]...)
		}
	}
	data, err := array.Marshal()
	if err != nil {
		log.WithError(err).Fatal("Failed to encode task events")
	}
	if err := db.Put(seedNodeKey, data); err != nil {
		log.WithError(err).Fatal("Failed to write task events")
	}
}
