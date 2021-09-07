// Copyright (C) 2021 The RosettaNet Authors.

package rawdb

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	dbtype "github.com/RosettaFlow/Carrier-Go/lib/db"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/sirupsen/logrus"
	"strings"
)

// ReadTaskDataHash retrieves the hash assigned to a canonical number and index.
func ReadTaskDataHash(db DatabaseReader, number uint64, index uint64) common.Hash {
	data, _ := db.Get(identityHashKey(number, index))
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteTaskDataHash stores the hash assigned to a canonical number and index.
func WriteTaskDataHash(db DatabaseWriter, number uint64, index uint64, hash common.Hash)  {
	if err := db.Put(identityHashKey(number, index), hash.Bytes()); err != nil {
		log.WithError(err).Fatal("Failed to store number-index to hash mapping for identity")
	}
}

// DeleteTaskDataHash removes the number-index to hash canonical mapping.
func DeleteTaskDataHash(db DatabaseDeleter, number uint64, index uint64) {
	if err := db.Delete(identityHashKey(number, index)); err != nil {
		log.WithError(err).Fatal("Failed to delete number-index to hash mapping for identity")
	}
}

// ReadTaskDataId retrieves the dataId assigned to a canonical nodeId and hash.
func ReadTaskDataId(db DatabaseReader, nodeId string, hash common.Hash) []byte {
	data, _ := db.Get(identityDataIdKey(common.Hex2Bytes(nodeId), hash))
	if len(data) == 0 {
		return []byte{}
	}
	return data
}

// WriteTaskDataId stores the dataId assigned to a canonical nodeId and hash.
func WriteTaskDataId(db DatabaseWriter, nodeId string, hash common.Hash, dataId string)  {
	if err := db.Put(identityDataIdKey(common.Hex2Bytes(nodeId), hash), common.Hex2Bytes(dataId)); err != nil {
		log.WithError(err).Fatal("Failed to store nodeId-hash to dataId mapping for identity")
	}
}

// DeleteTaskDataId removes the nodeId-hash to dataId canonical mapping.
func DeleteTaskDataId(db DatabaseDeleter, nodeId string, hash common.Hash,) {
	if err := db.Delete(identityDataIdKey(common.Hex2Bytes(nodeId), hash)); err != nil {
		log.WithError(err).Fatal("Failed to delete number-index to hash mapping for identity")
	}
}

// ReadTaskDataTypeHash retrieves the hash assigned to a canonical type and dataId.
func ReadTaskDataTypeHash(db DatabaseReader, dataId string, typ string) common.Hash {
	data, _ := db.Get(identityDataTypeHashKey(common.Hex2Bytes(dataId), common.Hex2Bytes(typ)))
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteTaskDataTypeHash stores the hash assigned to a canonical type and dataId.
func WriteTaskDataTypeHash(db DatabaseWriter, dataId string, typ string, hash common.Hash)  {
	if err := db.Put(identityDataTypeHashKey(common.Hex2Bytes(dataId), common.Hex2Bytes(typ)), hash.Bytes()); err != nil {
		log.WithError(err).Fatal("Failed to store type-dataId to hash mapping for identity")
	}
}

// DeleteTaskDataTypeHash removes the type-dataId to hash canonical mapping.
func DeleteTaskDataTypeHash(db DatabaseDeleter, dataId string, typ string) {
	if err := db.Delete(identityDataTypeHashKey(common.Hex2Bytes(dataId), common.Hex2Bytes(typ))); err != nil {
		log.WithError(err).Fatal("Failed to delete type-dataId to hash mapping for identity")
	}
}

// **************************************** current local operation **************************************************
// ReadRunningTask retrieves the task with the corresponding taskId.
func ReadRunningTask(db DatabaseReader, taskId string) *types.Task {
	blob, err := db.Get(runningTaskKey)
	if err != nil {
		return nil
	}
	var taskArray dbtype.TaskArrayPB
	if err := taskArray.Unmarshal(blob); err != nil {
		return nil
	}
	for _, task := range taskArray.GetTaskList() {
		if strings.EqualFold(task.GetTaskId(), taskId) {
			return types.NewTask(task)
		}
	}
	return nil
}

// ReadAllRunningTask retrieves all the running task in the database.
func ReadAllRunningTask(db DatabaseReader) []*types.Task {
	blob, err := db.Get(runningTaskKey)
	if err != nil {
		return nil
	}
	var array dbtype.TaskArrayPB
	if err := array.Unmarshal(blob); err != nil {
		return nil
	}
	var taskList []*types.Task
	for _, task := range array.GetTaskList() {
		taskList = append(taskList, types.NewTask(task))
	}
	return taskList
}

// WriteRunningTask serializes the running task into the database.
func WriteRunningTask(db KeyValueStore, task *types.Task) {
	blob, err := db.Get(runningTaskKey)
	if err != nil {
		log.Warn("Failed to load old running task", "error", err)
	}
	var array dbtype.TaskArrayPB
	if len(blob) > 0 {
		if err := array.Unmarshal(blob); err != nil {
			log.WithError(err).Fatal("Failed to decode old running task")
		}

	}
	for _, s := range array.GetTaskList() {
		if strings.EqualFold(s.TaskId, task.GetTaskId()) {
			log.WithFields(logrus.Fields{ "id": s.TaskId }).Info("Skip duplicated task")
			return
		}
	}
	array.TaskList = append(array.TaskList, task.GetTaskData())
	data, err := array.Marshal()
	if err != nil {
		log.WithError(err).Fatal("Failed to encode running task")
	}
	if err := db.Put(runningTaskKey, data); err != nil {
		log.WithError(err).Fatal("Failed to write running task")
	}
}

// DeleteRunningTask deletes the running task from the database with a special taskId
func DeleteRunningTask(db KeyValueStore, taskId string) {
	blob, err := db.Get(runningTaskKey)
	if err != nil {
		log.Warn("Failed to load old running task", "error", err)
	}
	var array dbtype.TaskArrayPB
	if len(blob) > 0 {
		if err := array.Unmarshal(blob); err != nil {
			log.WithError(err).Fatal("Failed to decode old running task")
		}
	}
	// need to test.
	for idx, s := range array.GetTaskList() {
		if strings.EqualFold(s.TaskId, taskId) {
			array.TaskList = append(array.TaskList[:idx], array.TaskList[idx+1:]...)
			break
		}
	}
	data, err := array.Marshal()
	if err != nil {
		log.WithError(err).Fatal("Failed to encode running task")
	}
	if err := db.Put(runningTaskKey, data); err != nil {
		log.WithError(err).Fatal("Failed to write running task")
	}
}