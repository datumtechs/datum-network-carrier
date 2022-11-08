// Copyright (C) 2021 The datum-network Authors.

package rawdb

import (
	"bytes"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common/bytesutil"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carrierdbpb "github.com/datumtechs/datum-network-carrier/pb/carrier/db"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
	"strings"
)

const seedNodeToKeep = 50
const registeredNodeToKeep = 50

// QueryLocalIdentity retrieves the identity of local.
func QueryLocalIdentity(db DatabaseReader) (*carriertypespb.Organization, error) {
	var blob carriertypespb.Organization
	enc, err := db.Get(localIdentityKey)
	if nil != err {
		return nil, err
	}
	if err := blob.Unmarshal(enc); nil != err {
		return nil, err
	}
	return &blob, nil
}

// StoreLocalIdentity stores the local identity.
func StoreLocalIdentity(db DatabaseWriter, identity *carriertypespb.Organization) error {
	enc, err := identity.Marshal()
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
	var yarnName carrierdbpb.StringPB
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
	pb := carrierdbpb.StringPB{
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

// QueryAllSeedNodes retrieves all the seed nodes in the database.
// All returned seed nodes are sorted in reverse.
func QueryAllSeedNodes(db DatabaseReader) ([]*carrierapipb.SeedPeer, error) {
	blob, err := db.Get(seedNodeKey)
	if nil != err {
		return nil, err
	}
	var seedNodes carrierdbpb.SeedPeerListPB
	if err := seedNodes.Unmarshal(blob); nil != err {
		return nil, err
	}
	var nodes []*carrierapipb.SeedPeer
	for _, seed := range seedNodes.SeedPeerList {
		nodes = append(nodes, &carrierapipb.SeedPeer{
			Addr:      seed.GetAddr(),
			IsDefault: false,
			ConnState: carrierapipb.ConnState_ConnState_UnConnected,
		})
	}
	return nodes, nil
}

// StoreSeedNode serializes the seed node into the database.
func StoreSeedNode(db KeyValueStore, seedNode *carrierapipb.SeedPeer) error {
	blob, err := db.Get(seedNodeKey)
	if IsNoDBNotFoundErr(err) {
		log.WithError(err).Error("Failed to load old seed nodes")
		return err
	}
	var seedNodes carrierdbpb.SeedPeerListPB
	if len(blob) > 0 {
		if err := seedNodes.Unmarshal(blob); nil != err {
			log.WithError(err).Error("Failed to decode old seed nodes")
			return err
		}

	}
	for _, s := range seedNodes.GetSeedPeerList() {
		if strings.EqualFold(s.GetAddr(), seedNode.GetAddr()) {
			log.WithFields(logrus.Fields{"addr": s.GetAddr()}).Info("Skip duplicated seed node")
			return nil
		}
	}

	// max limit for store seed node.
	if len(seedNodes.GetSeedPeerList()) >= seedNodeToKeep {
		return fmt.Errorf("The number of seed nodes overflows the maximum storage limit, has %d seed nodes", len(seedNodes.GetSeedPeerList()))
	}

	// append new seed node into arr
	seedNodes.SeedPeerList = append(seedNodes.SeedPeerList, &carrierdbpb.SeedPeerPB{Addr: seedNode.GetAddr()})

	data, err := seedNodes.Marshal()
	if nil != err {
		log.WithError(err).Error("Failed to encode seed nodes")
		return err
	}
	return db.Put(seedNodeKey, data)
}

// RemoveSeedNode deletes the seed nodes from the database with a special id
func RemoveSeedNode(db KeyValueStore, addr string) error {
	blob, err := db.Get(seedNodeKey)
	switch {
	case IsNoDBNotFoundErr(err):
		log.WithError(err).Error("Failed to load old seed nodes")
		return err
	case IsDBNotFoundErr(err), nil == err && len(blob) == 0:
		return nil
	}
	var seedNodes carrierdbpb.SeedPeerListPB
	if len(blob) > 0 {
		if err := seedNodes.Unmarshal(blob); nil != err {
			log.WithError(err).Error("Failed to decode old seed nodes")
			return err
		}
	}
	// removed the seed node.
	for idx, s := range seedNodes.GetSeedPeerList() {
		if strings.EqualFold(s.GetAddr(), addr) {
			seedNodes.SeedPeerList = append(seedNodes.SeedPeerList[:idx], seedNodes.SeedPeerList[idx+1:]...)
			break
		}
	}

	// removed emtpy seed nodes structure
	if len(seedNodes.GetSeedPeerList()) == 0 {
		return db.Delete(seedNodeKey)
	}

	data, err := seedNodes.Marshal()
	if nil != err {
		log.WithError(err).Fatal("Failed to encode seed nodes")
	}
	return db.Put(seedNodeKey, data)
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

func registryNodeKeyPrefix(nodeType carrierapipb.RegisteredNodeType) []byte {
	var key []byte
	if nodeType == carrierapipb.PrefixTypeJobNode {
		key = getJobNodeKeyPrefix()
	}
	if nodeType == carrierapipb.PrefixTypeDataNode {
		key = getDataNodeKeyPrefix()
	}
	return key
}

func registryNodeKey(nodeType carrierapipb.RegisteredNodeType, nodeId string) []byte {
	var key []byte
	if nodeType == carrierapipb.PrefixTypeJobNode {
		key = getJobNodeKey(nodeId)
	}
	if nodeType == carrierapipb.PrefixTypeDataNode {
		key = getDataNodeKey(nodeId)
	}
	return key
}

// QueryRegisterNode retrieves the register node with the corresponding nodeId.
func QueryRegisterNode(db DatabaseReader, nodeType carrierapipb.RegisteredNodeType, nodeId string) (*carrierapipb.YarnRegisteredPeerDetail, error) {
	blob, err := db.Get(registryNodeKey(nodeType, nodeId))
	if nil != err {
		return nil, err
	}
	registeredNode := &carrierdbpb.RegisteredNodePB{}
	if err := registeredNode.Unmarshal(blob); nil != err {
		log.WithError(err).Errorf("registeredNode decode failed")
		return nil, err
	}
	return &carrierapipb.YarnRegisteredPeerDetail{
		Id:           registeredNode.Id,
		InternalIp:   registeredNode.InternalIp,
		InternalPort: registeredNode.InternalPort,
		ExternalIp:   registeredNode.ExternalIp,
		ExternalPort: registeredNode.ExternalPort,
		ConnState:    carrierapipb.ConnState_ConnState_UnConnected,
	}, nil
}

// QueryAllRegisterNodes retrieves all the registered nodes in the database.
// All returned registered nodes are sorted in reverse.
func QueryAllRegisterNodes(db KeyValueStore, nodeType carrierapipb.RegisteredNodeType) ([]*carrierapipb.YarnRegisteredPeerDetail, error) {

	it := db.NewIteratorWithPrefixAndStart(registryNodeKeyPrefix(nodeType), nil)
	defer it.Release()

	arr := make([]*carrierapipb.YarnRegisteredPeerDetail, 0)
	for it.Next() {
		if blob := it.Value(); len(blob) != 0 {
			registeredNode := &carrierdbpb.RegisteredNodePB{}
			if err := registeredNode.Unmarshal(blob); nil != err {
				log.WithError(err).Warnf("Warning registeredNode decode failed")
				continue
			}
			arr = append(arr, &carrierapipb.YarnRegisteredPeerDetail{
				Id:           registeredNode.Id,
				InternalIp:   registeredNode.InternalIp,
				InternalPort: registeredNode.InternalPort,
				ExternalIp:   registeredNode.ExternalIp,
				ExternalPort: registeredNode.ExternalPort,
				ConnState:    carrierapipb.ConnState_ConnState_UnConnected,
			})
		}
	}

	if len(arr) == 0 {
		return nil, ErrNotFound
	}
	return arr, nil
}

// StoreRegisterNode serializes the registered node into the database. If the cumulated
// registered node exceeds the limitation, the oldest will be dropped.
func StoreRegisterNode(db DatabaseWriter, nodeType carrierapipb.RegisteredNodeType, registeredNode *carrierapipb.YarnRegisteredPeerDetail) error {

	key := registryNodeKey(nodeType, registeredNode.GetId())
	val := &carrierdbpb.RegisteredNodePB{
		Id:           registeredNode.Id,
		InternalIp:   registeredNode.InternalIp,
		InternalPort: registeredNode.InternalPort,
		ExternalIp:   registeredNode.ExternalIp,
		ExternalPort: registeredNode.ExternalPort,
	}
	data, err := val.Marshal()
	if nil != err {
		log.WithError(err).Error("Failed to encode registered node")
		return err
	}
	return db.Put(key, data)
}

func RemoveRegisterNode(db KeyValueStore, nodeType carrierapipb.RegisteredNodeType, id string) error {
	key := registryNodeKey(nodeType, id)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

// RemoveRegisterNodes deletes all the registered nodes from the database.
func RemoveRegisterNodes(db KeyValueStore, nodeType carrierapipb.RegisteredNodeType) error {

	it := db.NewIteratorWithPrefixAndStart(registryNodeKeyPrefix(nodeType), nil)
	defer it.Release()

	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			db.Delete(key)
		}
	}
	return nil
}

// QueryLocalResource retrieves the resource of local with the corresponding jobNodeId.
func QueryLocalResource(db DatabaseReader, jobNodeId string) (*types.LocalResource, error) {
	blob, err := db.Get(localResourceKey(jobNodeId))
	if nil != err {
		log.WithError(err).Errorf("Failed to read local resource")
		return nil, err
	}
	localResource := new(carriertypespb.LocalResourcePB)
	if err := localResource.Unmarshal(blob); nil != err {
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
			localResource := new(carriertypespb.LocalResourcePB)
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

// StoreLocalTask serializes the local task into the database.
func StoreLocalTask(db KeyValueStore, task *types.Task) error {

	key := getLocalTaskKey(task.GetTaskId())

	data, err := task.GetTaskData().Marshal()
	if nil != err {
		log.WithError(err).Error("Failed to encode local task node")
		return err
	}
	return db.Put(key, data)
}

// RemoveLocalTask deletes the local task from the database with a special taskId
func RemoveLocalTask(db KeyValueStore, taskId string) error {
	key := getLocalTaskKey(taskId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

// RemoveLocalAllTask deletes all the local task from the database.
func RemoveLocalAllTask(db KeyValueStore) error {
	it := db.NewIteratorWithPrefixAndStart(getLocalTaskKeyPrefix(), nil)
	defer it.Release()

	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			db.Delete(key)
		}
	}
	return nil
}

// QueryLocalTask retrieves the local task with the corresponding taskId.
func QueryLocalTask(db DatabaseReader, taskId string) (*types.Task, error) {
	blob, err := db.Get(getLocalTaskKey(taskId))
	if nil != err {
		return nil, err
	}
	task := &carriertypespb.TaskPB{}
	if err := task.Unmarshal(blob); nil != err {
		log.WithError(err).Errorf("local task decode failed")
		return nil, err
	}
	return types.NewTask(task), nil
}

// QueryLocalTaskByIds retrieves the local tasks with the corresponding taskIds.
func QueryLocalTaskByIds(db KeyValueStore, taskIds []string) (types.TaskDataArray, error) {

	arr := make(types.TaskDataArray, 0)
	for _, taskId := range taskIds {
		blob, err := db.Get(getLocalTaskKey(taskId))
		if nil != err {
			log.WithError(err).Warnf("Warning load local task failed")
			continue
		}
		task := &carriertypespb.TaskPB{}
		if err := task.Unmarshal(blob); nil != err {
			log.WithError(err).Warnf("Warning local task decode failed")
			continue
		}
		arr = append(arr, types.NewTask(task))
	}

	if len(arr) == 0 {
		return nil, ErrNotFound
	}
	return arr, nil
}

// QueryAllLocalTasks retrieves all the local task in the database.
func QueryAllLocalTasks(db KeyValueStore) (types.TaskDataArray, error) {
	it := db.NewIteratorWithPrefixAndStart(getLocalTaskKeyPrefix(), nil)
	defer it.Release()

	arr := make(types.TaskDataArray, 0)
	for it.Next() {
		if blob := it.Value(); len(blob) != 0 {
			task := &carriertypespb.TaskPB{}
			if err := task.Unmarshal(blob); nil != err {
				log.WithError(err).Warnf("Warning local task decode failed")
				continue
			}
			arr = append(arr, types.NewTask(task))
		}
	}

	if len(arr) == 0 {
		return nil, ErrNotFound
	}
	return arr, nil
}

// QueryTaskEvent retrieves the evengine of task with the corresponding taskId for all partyIds.
func QueryTaskEvent(db KeyValueStore, taskId string) ([]*carriertypespb.TaskEvent, error) {

	it := db.NewIteratorWithPrefixAndStart(append(taskEventPrefix, []byte(taskId)...), nil)
	defer it.Release()

	result := make([]*carriertypespb.TaskEvent, 0)

	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			blob, err := db.Get(key)
			if nil != err {
				continue
			}
			var events carrierdbpb.TaskEventArrayPB
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
func QueryTaskEventByPartyId(db DatabaseReader, taskId, partyId string) ([]*carriertypespb.TaskEvent, error) {

	key := taskEventKey(taskId, partyId)

	val, err := db.Get(key)
	if nil != err {
		return nil, err
	}
	var events carrierdbpb.TaskEventArrayPB
	if err := events.Unmarshal(val); nil != err {
		log.WithError(err).Errorf("Failed to encode task events")
		return nil, err
	}
	return events.GetTaskEventList(), nil
}

// QueryAllTaskEvents retrieves the task event with all (all taskIds and all partyIds).
func QueryAllTaskEvents(db KeyValueStore) ([]*carriertypespb.TaskEvent, error) {

	it := db.NewIteratorWithPrefixAndStart(taskEventPrefix, nil)
	defer it.Release()

	result := make([]*carriertypespb.TaskEvent, 0)

	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			blob, err := db.Get(key)
			if nil != err {
				continue
			}
			var events carrierdbpb.TaskEventArrayPB
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
func StoreTaskEvent(db KeyValueStore, taskEvent *carriertypespb.TaskEvent) error {
	key := taskEventKey(taskEvent.GetTaskId(), taskEvent.GetPartyId())
	val, err := db.Get(key)
	if IsNoDBNotFoundErr(err) {
		log.WithError(err).Error("Failed to load old task events")
		return err
	}
	var array carrierdbpb.TaskEventArrayPB
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

func StoreNeedExecuteTask(db KeyValueStore, task *types.NeedExecuteTask) error {
	key := GetNeedExecuteTaskKey(task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())

	var errStr string
	if nil != task.GetErr() {
		errStr = task.GetErr().Error()
	}

	conSumeTypes := make([]uint32, len(task.GetConsumeSpec().GetConsumeTypes()))
	consumeOptions := make([]string, len(task.GetConsumeSpec().GetConsumeTypes()))
	for i, _ := range task.GetConsumeSpec().GetConsumeTypes() {
		conSumeTypes[i] = uint32(task.GetConsumeSpec().GetConsumeTypes()[i])
		consumeOptions[i] = task.GetConsumeSpec().GetConsumeOptions()[i]
	}
	val, err := proto.Marshal(&carriertypespb.NeedExecuteTask{
		RemotePid:              task.GetRemotePID().String(),
		LocalTaskRole:          task.GetLocalTaskRole(),
		LocalTaskOrganization:  task.GetLocalTaskOrganization(),
		RemoteTaskRole:         task.GetRemoteTaskRole(),
		RemoteTaskOrganization: task.GetRemoteTaskOrganization(),
		TaskId:                 task.GetTaskId(),
		ConsStatus:             bytesutil.Uint16ToBytes(uint16(task.GetConsStatus())),
		LocalResource: &carriertypespb.PrepareVoteResource{
			Id:      task.GetLocalResource().GetId(),
			Ip:      task.GetLocalResource().GetIp(),
			Port:    task.GetLocalResource().GetPort(),
			PartyId: task.GetLocalResource().GetPartyId(),
		},
		Resources: task.GetResources(),
		ErrStr:    errStr,
		ConsumeSpec: &carriertypespb.TaskConsumeSpec{
			ConsumeTypes:   conSumeTypes,
			ConsumeOptions: consumeOptions,
		},
	})
	if nil != err {
		return fmt.Errorf("marshal needExecuteTask failed, %s", err)
	}
	return db.Put(key, val)
}

func RemoveNeedExecuteTaskByPartyId(db KeyValueStore, taskId, partyId string) error {
	key := GetNeedExecuteTaskKey(taskId, partyId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func RemoveNeedExecuteTask(db KeyValueStore, taskId string) error {
	prefix := append(needExecuteTaskKeyPrefix, []byte(taskId)...)
	it := db.NewIteratorWithPrefixAndStart(prefix, nil)
	defer it.Release()
	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			db.Delete(key)
		}
	}
	return nil
}

func ForEachNeedExecuteTaskWwithPrefix(db KeyValueStore, prefix []byte, f func(key, value []byte) error) error {
	it := db.NewIteratorWithPrefixAndStart(prefix, nil)
	defer it.Release()
	for it.Next() {
		if err := f(it.Key(), it.Value()); nil != err {
			return err
		}
	}
	return nil
}

func ForEachNeedExecuteTask(db KeyValueStore, f func(key, value []byte) error) error {
	it := db.NewIteratorWithPrefixAndStart(GetNeedExecuteTaskKeyPrefix(), nil)
	defer it.Release()
	for it.Next() {
		if err := f(it.Key(), it.Value()); nil != err {
			return err
		}
	}
	return nil
}

// StoreTaskBullet save scheduled tasks.
func StoreTaskBullet(db KeyValueStore, bullet *types.TaskBullet) error {

	key := GetTaskBulletKey(bullet.GetTaskId())

	val, err := rlp.EncodeToBytes(bullet)
	if nil != err {
		log.WithError(err).Errorf("encode taskBullet failed")
		return err
	}
	return db.Put(key, val)
}

func RemoveTaskBullet(db KeyValueStore, taskId string) error {
	key := GetTaskBulletKey(taskId)

	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func ForEachTaskBullets(db KeyValueStore, f func(key, value []byte) error) error {
	it := db.NewIteratorWithPrefixAndStart(GetTaskBulletKeyPrefix(), nil)
	defer it.Release()
	for it.Next() {
		if err := f(it.Key(), it.Value()); nil != err {
			return err
		}
	}
	return nil
}

// ================================= Local Metadata ==========================================
// QueryLocalMetadata retrieves the metadata of local with the corresponding metadataId.
func QueryLocalMetadata(db DatabaseReader, metadataId string) (*types.Metadata, error) {
	blob, err := db.Get(localMetadataKey(metadataId))
	if IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed to query local metadata")
		return nil, err
	}
	if IsDBNotFoundErr(err) {
		log.WithError(err).Warnf("Warning query local metadata not found")
		return nil, err
	}
	localMetadata := new(carriertypespb.MetadataPB)
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
			localMetadata := new(carriertypespb.MetadataPB)
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
