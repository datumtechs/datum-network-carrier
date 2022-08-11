package rawdb

import (
	"encoding/json"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common/bytesutil"
	"github.com/datumtechs/datum-network-carrier/db"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gogo/protobuf/proto"
	leveldberr "github.com/syndtr/goleveldb/leveldb/errors"
	"strings"
)

var (
	ErrNotFound = leveldberr.ErrNotFound
)

func IsNoDBNotFoundErr(err error) bool {
	return nil != err && err != ErrNotFound
}
func IsDBNotFoundErr(err error) bool {
	return nil != err && err == ErrNotFound
}

// Resources that operate the organization's jobNode services
func StoreNodeResource(db KeyValueStore, resource *types.LocalResourceTable) error {

	item_key := GetNodeResourceKey(resource.GetNodeId())
	val, err := rlp.EncodeToBytes(resource)
	if nil != err {
		return err
	}
	return db.Put(item_key, val)
}

func StoreNodeResources(db KeyValueStore, resources []*types.LocalResourceTable) error {

	for _, resource := range resources {
		key := GetNodeResourceKey(resource.GetNodeId())
		val, err := rlp.EncodeToBytes(resource)
		if nil != err {
			return err
		}

		if err := db.Put(key, val); nil != err {
			return err
		}
	}
	return nil
}

func RemoveNodeResource(db KeyValueStore, resourceId string) error {
	key := GetNodeResourceKey(resourceId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func QueryNodeResource(db DatabaseReader, resourceId string) (*types.LocalResourceTable, error) {
	key := GetNodeResourceKey(resourceId)
	vb, err := db.Get(key)
	if nil != err {
		return nil, err
	}

	var resource types.LocalResourceTable
	if err := rlp.DecodeBytes(vb, &resource); nil != err {
		return nil, err
	}
	return &resource, nil
}

func QueryNodeResources(db KeyValueStore) ([]*types.LocalResourceTable, error) {

	prefix := GetNodeResourceKeyPrefix()
	it := db.NewIteratorWithPrefixAndStart(prefix, nil)
	defer it.Release()

	arr := make([]*types.LocalResourceTable, 0)
	for it.Next() {
		if len(it.Key()) != 0 && len(it.Value()) != 0 {
			// prefix + jobNodeId -> LocalResourceTable
			var resource types.LocalResourceTable
			if err := rlp.DecodeBytes(it.Value(), &resource); nil != err {
				return nil, err
			}
			arr = append(arr, &resource)
		}
	}

	if len(arr) == 0 {
		return nil, ErrNotFound
	}

	return arr, nil
}

//// Operation resource slot unit definition
//func StoreNodeResourceSlotUnit(db DatabaseWriter, slot *types.Slot) error {
//	val, err := rlp.EncodeToBytes(slot)
//	if nil != err {
//		return err
//	}
//	return db.Put(GetNodeResourceSlotUnitKey(), val)
//}
//
//func RemoveNodeResourceSlotUnit(db KeyValueStore) error {
//	key := GetNodeResourceSlotUnitKey()
//	has, err := db.Has(key)
//	switch {
//	case IsNoDBNotFoundErr(err):
//		return err
//	case IsDBNotFoundErr(err), nil == err && !has:
//		return nil
//	}
//	return db.Delete(key)
//}
//
//func QueryNodeResourceSlotUnit(db DatabaseReader) (*types.Slot, error) {
//	has, err := db.Has(GetNodeResourceSlotUnitKey())
//	if IsNoDBNotFoundErr(err) {
//		return nil, err
//	}
//	if !has {
//		return nil, ErrNotFound
//	}
//	b, err := db.Get(GetNodeResourceSlotUnitKey())
//	if nil != err {
//		return nil, err
//	}
//	var slot *types.Slot
//	if err := rlp.DecodeBytes(b, &slot); nil != err {
//		return nil, err
//	}
//	return slot, nil
//}

// Operate the information of the jobNode resources being used by the local task
func StoreLocalTaskPowerUsed(db KeyValueStore, taskPowerUsed *types.LocalTaskPowerUsed) error {
	// prefix + taskId + partyId -> LocalTaskPowerUsed
	key := GetLocalTaskPowerUsedKey(taskPowerUsed.GetTaskId(), taskPowerUsed.GetPartyId())
	val, err := rlp.EncodeToBytes(taskPowerUsed)
	if nil != err {
		return err
	}
	log.Debugf("Call StoreLocalTaskPowerUsed, taskId: {%s}, partyId: {%s}, used: {%s}", taskPowerUsed.GetTaskId(), taskPowerUsed.GetPartyId(), taskPowerUsed.String())
	return db.Put(key, val)
}

func StoreLocalTaskPowerUseds(db KeyValueStore, taskPowerUseds []*types.LocalTaskPowerUsed) error {
	for _, used := range taskPowerUseds {
		key := GetLocalTaskPowerUsedKey(used.GetTaskId(), used.GetPartyId())
		val, err := rlp.EncodeToBytes(used)
		if nil != err {
			return err
		}

		if err := db.Put(key, val); nil != err {
			return err
		}
	}
	return nil
}

func HasLocalTaskPowerUsed(db DatabaseReader, taskId, partyId string) (bool, error) {

	has, err := db.Has(GetLocalTaskPowerUsedKey(taskId, partyId))

	switch {
	case IsNoDBNotFoundErr(err):
		return false, err
	case IsDBNotFoundErr(err), !has:
		return false, nil
	}
	return true, nil
}

func RemoveLocalTaskPowerUsed(db KeyValueStore, taskId, partyId string) error {
	key := GetLocalTaskPowerUsedKey(taskId, partyId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func RemoveLocalTaskPowerUsedByTaskId(db KeyValueStore, taskId string) error {
	it := db.NewIteratorWithPrefixAndStart(GetLocalTaskPowerUsedKeyPrefixByTaskId(taskId), nil)
	defer it.Release()

	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			db.Delete(key)
		}
	}

	return nil
}

func QueryLocalTaskPowerUsed(db DatabaseReader, taskId, partyId string) (*types.LocalTaskPowerUsed, error) {
	// prefix + taskId + partyId -> LocalTaskPowerUsed
	key := GetLocalTaskPowerUsedKey(taskId, partyId)
	vb, err := db.Get(key)
	if nil != err {
		return nil, err
	}

	var taskPowerUsed types.LocalTaskPowerUsed

	if err := rlp.DecodeBytes(vb, &taskPowerUsed); nil != err {
		return nil, err
	}
	used := &taskPowerUsed
	log.Debugf("Call QueryLocalTaskPowerUsed, taskId: {%s}, partyId: {%s}, used: {%s}", taskId, partyId, used.String())
	return &taskPowerUsed, nil
}

func QueryLocalTaskPowerUsedsByTaskId(db KeyValueStore, taskId string) ([]*types.LocalTaskPowerUsed, error) {
	// prefix + taskId + partyId -> LocalTaskPowerUsed
	it := db.NewIteratorWithPrefixAndStart(GetLocalTaskPowerUsedKeyPrefixByTaskId(taskId), nil)
	defer it.Release()

	arr := make([]*types.LocalTaskPowerUsed, 0)
	for it.Next() {
		if value := it.Value(); len(value) != 0 {
			var taskPowerUsed types.LocalTaskPowerUsed
			if err := rlp.DecodeBytes(value, &taskPowerUsed); nil != err {
				return nil, err
			}
			arr = append(arr, &taskPowerUsed)
		}
	}

	if len(arr) == 0 {
		return nil, ErrNotFound
	}
	return arr, nil
}

func QueryLocalTaskPowerUseds(db KeyValueStore) ([]*types.LocalTaskPowerUsed, error) {
	// prefix + taskId + partyId -> LocalTaskPowerUsed
	it := db.NewIteratorWithPrefixAndStart(GetLocalTaskPowerUsedKeyPrefix(), nil)
	defer it.Release()

	arr := make([]*types.LocalTaskPowerUsed, 0)
	for it.Next() {
		if value := it.Value(); len(value) != 0 {
			var taskPowerUsed types.LocalTaskPowerUsed
			if err := rlp.DecodeBytes(value, &taskPowerUsed); nil != err {
				return nil, err
			}
			arr = append(arr, &taskPowerUsed)
		}
	}

	if len(arr) == 0 {
		return nil, ErrNotFound
	}
	return arr, nil
}

func StoreJobNodeTaskPartyId(db KeyValueStore, jobNodeId, taskId, partyId string) error {
	// prefix + jobNodeId + taskId -> [partyId, ..., partyId]
	key := GetJobNodeTaskPartyIdsKey(jobNodeId, taskId)
	has, err := db.Has(key)
	if IsNoDBNotFoundErr(err) {
		return err
	}
	var partyIdArr []string
	if !has {
		partyIdArr = []string{partyId}
	} else {

		val, err := db.Get(key)
		if nil != err {
			return err
		}
		if err := rlp.DecodeBytes(val, &partyIdArr); nil != err {
			return err
		}

		var find bool
		for _, id := range partyIdArr {
			if id == partyId {
				find = true
				break
			}
		}
		if !find {
			partyIdArr = append(partyIdArr, partyId)
		}
	}
	val, err := rlp.EncodeToBytes(partyIdArr)
	if nil != err {
		return err
	}
	log.Debugf("Call StoreJobNodeTaskPartyId, jobNodeId: {%s}, taskId: {%s}, partyId: {%s}, partyIds: %s", jobNodeId, taskId, partyId, partyIdArr)
	return db.Put(key, val)
}

func RemoveJobNodeTaskPartyId(db KeyValueStore, jobNodeId, taskId, partyId string) error {
	// prefix + jobNodeId + taskId -> [partyId, ..., partyId]
	key := GetJobNodeTaskPartyIdsKey(jobNodeId, taskId)
	val, err := db.Get(key)

	var partyIdArr []string
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err):
		return nil
	case nil == err && len(val) != 0:
		if err := rlp.DecodeBytes(val, &partyIdArr); nil != err {
			return err
		}
	}

	//for i := 0; i < len(partyIdArr); i++ {
	//
	//	id := partyIdArr[i]
	//	if id == partyId {
	//		partyIdArr = append(partyIdArr[:i], partyIdArr[i+1:]...)
	//		i--
	//	}
	//}

	for i, id := range partyIdArr {
		if id == partyId {
			partyIdArr = append(partyIdArr[:i], partyIdArr[i+1:]...)
			break
		}
	}

	if len(partyIdArr) == 0 {
		log.Debugf("Call RemoveJobNodeTaskPartyId [clean all partyIds], jobNodeId: {%s}, taskId: {%s}, partyId: {%s}", jobNodeId, taskId, partyId)
		return db.Delete(key)
	}
	val, err = rlp.EncodeToBytes(partyIdArr)
	if nil != err {
		return err
	}
	log.Debugf("Call RemoveJobNodeTaskPartyId, jobNodeId: {%s}, taskId: {%s}, partyId: {%s}, partyIds: %s", jobNodeId, taskId, partyId, partyIdArr)
	return db.Put(key, val)
}

func RemoveJobNodeTaskIdAllPartyIds(db KeyValueStore, jobNodeId, taskId string) error {
	// prefix + jobNodeId + taskId -> [partyId, ..., partyId]
	key := GetJobNodeTaskPartyIdsKey(jobNodeId, taskId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func QueryJobNodeRunningTaskIds(db KeyValueStore, jobNodeId string) ([]string, error) {
	// prefix + jobNodeId + taskId -> [partyId, ..., partyId]
	prefixAndJobNodeId := GetJobNodeTaskPartyIdsKeyPrefixByJobNodeId(jobNodeId)
	it := db.NewIteratorWithPrefixAndStart(prefixAndJobNodeId, nil)
	defer it.Release()

	arr := make([]string, 0)
	tmp := make(map[string]struct{}, 0)
	for it.Next() {
		if len(it.Key()) != 0 && len(it.Value()) != 0 {
			// key len == len(prefix) + len([]byte(jobNodeId)) + len([]byte(taskId))
			taskId := string(it.Key()[len(prefixAndJobNodeId):])
			if _, ok := tmp[taskId]; !ok {
				tmp[taskId] = struct{}{}
				arr = append(arr, taskId)
			}
		}
	}

	if len(arr) == 0 {
		return nil, ErrNotFound
	}
	log.Debugf("Call QueryJobNodeRunningTaskIds, jobNodeId: {%s}, taskIds: %s", jobNodeId, "["+strings.Join(arr, ",")+"]")
	return arr, nil
}

func QueryJobNodeRunningTaskIdCount(db KeyValueStore, jobNodeId string) (uint32, error) {
	// prefix + jobNodeId + taskId -> [partyId, ..., partyId]
	prefixAndJobNodeId := GetJobNodeTaskPartyIdsKeyPrefixByJobNodeId(jobNodeId)
	it := db.NewIteratorWithPrefixAndStart(prefixAndJobNodeId, nil)
	defer it.Release()

	var count uint32
	tmp := make(map[string]struct{}, 0)
	for it.Next() {
		if len(it.Key()) != 0 && len(it.Value()) != 0 {
			// key len == len(prefix) + len([]byte(jobNodeId)) + len([]byte(taskId))
			taskId := string(it.Key()[len(prefixAndJobNodeId):])
			if _, ok := tmp[taskId]; !ok {
				tmp[taskId] = struct{}{}
				count++
			}
		}
	}
	log.Debugf("Call QueryJobNodeRunningTaskCount, jobNodeId: {%s}, taskIds count: %d", jobNodeId, count)
	return count, nil
}

func QueryJobNodeRunningTaskIdsAndPartyIdsPairs(db KeyValueStore, jobNodeId string) (map[string][]string, error) {
	// prefix + jobNodeId + taskId -> [partyId, ..., partyId]
	prefixAndJobNodeId := GetJobNodeTaskPartyIdsKeyPrefixByJobNodeId(jobNodeId)
	it := db.NewIteratorWithPrefixAndStart(prefixAndJobNodeId, nil)
	defer it.Release()

	res := make(map[string][]string, 0)
	tmp := make(map[string]struct{}, 0)
	for it.Next() {
		if len(it.Key()) != 0 && len(it.Value()) != 0 {
			// key len == len(prefix) + len([]byte(jobNodeId)) + len([]byte(taskId))
			taskId := string(it.Key()[len(prefixAndJobNodeId):])
			if _, ok := tmp[taskId]; !ok {

				tmp[taskId] = struct{}{}

				var partyIdArr []string
				if err := rlp.DecodeBytes(it.Value(), &partyIdArr); nil != err {
					return nil, err
				}
				res[taskId] = partyIdArr
			}
		}
	}

	if len(res) == 0 {
		return nil, ErrNotFound
	}

	return res, nil
}

func QueryJobNodeTaskAllPartyIds(db KeyValueStore, jobNodeId, taskId string) ([]string, error) {
	// prefix + jobNodeId + taskId -> [partyId, ..., partyId]
	key := GetJobNodeTaskPartyIdsKey(jobNodeId, taskId)
	val, err := db.Get(key)
	var partyIdArr []string
	switch {
	case IsNoDBNotFoundErr(err):
		return nil, err
	case IsDBNotFoundErr(err):
		return nil, ErrNotFound
	case nil == err && len(val) != 0:
		if err := rlp.DecodeBytes(val, &partyIdArr); nil != err {
			return nil, err
		}
	}
	return partyIdArr, nil
}

func HasJobNodeRunningTaskId(db DatabaseReader, jobNodeId, taskId string) (bool, error) {
	// prefix + jobNodeId + taskId -> [partyId, ..., partyId]
	key := GetJobNodeTaskPartyIdsKey(jobNodeId, taskId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return false, err
	case IsDBNotFoundErr(err):
		return false, nil
	case nil == err && !has:
		return false, nil
	}
	return true, nil
}

func HasJobNodeTaskPartyId(db DatabaseReader, jobNodeId, taskId, partyId string) (bool, error) {
	// prefix + jobNodeId + taskId -> [partyId, ..., partyId]
	key := GetJobNodeTaskPartyIdsKey(jobNodeId, taskId)
	val, err := db.Get(key)
	var partyIdArr []string
	switch {
	case IsNoDBNotFoundErr(err):
		return false, err
	case IsDBNotFoundErr(err):
		return false, nil
	case nil == err && len(val) != 0:
		if err := rlp.DecodeBytes(val, &partyIdArr); nil != err {
			return false, err
		}
	}

	for _, id := range partyIdArr {
		if id == partyId {
			return true, nil
		}
	}
	return false, nil
}

func QueryJobNodeTaskPartyIdCount(db DatabaseReader, jobNodeId, taskId string) (uint32, error) {
	// prefix + jobNodeId + taskId -> [partyId, ..., partyId]
	key := GetJobNodeTaskPartyIdsKey(jobNodeId, taskId)
	val, err := db.Get(key)

	var partyIdArr []string
	switch {
	case IsNoDBNotFoundErr(err):
		return 0, err
	case IsDBNotFoundErr(err):
		return 0, nil
	case nil == err && len(val) != 0:
		if err := rlp.DecodeBytes(val, &partyIdArr); nil != err {
			return 0, err
		}
	}
	return uint32(len(partyIdArr)), nil
}

// about jobNode history task
func StoreJobNodeHistoryTaskId(db KeyValueStore, jobNodeId, taskId string) error {

	// prefix + jobNodeId + taskId -> index
	item_key := GetJobNodeHistoryTaskKey(jobNodeId, taskId)
	has, err := db.Has(item_key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case nil == err && has:
		return nil // It have been exists, don't inscrease count
	}

	// When taskId have not on jobNode, inscrease jobNode taskId count
	// and put taskId on jobNodeId mapping.
	//
	// prefix + jobNodeId -> history task count
	count_key := GetJobNodeHistoryTaskCountKey(jobNodeId)
	count_val, err := db.Get(count_key)

	var count uint32

	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err):
		// do nothing
	case nil == err && len(count_val) != 0:
		count = bytesutil.BytesToUint32(count_val)
	}
	count++

	count_val = bytesutil.Uint32ToBytes(count)

	// First: put taskId on jobNode mapping.
	if err := db.Put(item_key, count_val); nil != err {
		return err
	}
	log.Debugf("InscreaseJobNodeHistoryTaskCount, jobNodeId: {%s}, taskId: {%s}, count: {%d}", jobNodeId, taskId, count)
	// Second: inscease taskId count on jobNode.
	return db.Put(count_key, count_val)
}

func HasJobNodeHistoryTaskId(db DatabaseReader, jobNodeId, taskId string) (bool, error) {

	item_key := GetJobNodeHistoryTaskKey(jobNodeId, taskId)
	has, err := db.Has(item_key)
	switch {
	case IsNoDBNotFoundErr(err):
		return false, err
	case IsDBNotFoundErr(err):
		return false, nil
	case nil == err && !has:
		return false, nil
	}
	return true, nil
}

func QueryJobNodeHistoryTaskCount(db KeyValueStore, jobNodeId string) (uint32, error) {
	// prefix + jobNodeId -> history task count
	key := GetJobNodeHistoryTaskCountKey(jobNodeId)
	val, err := db.Get(key)

	var count uint32

	switch {
	case IsNoDBNotFoundErr(err):
		return 0, err
	case IsDBNotFoundErr(err):
		return 0, nil
	case nil == err && len(val) != 0:
		count = bytesutil.BytesToUint32(val)
	}
	return count, nil
}

// Operation local dataNode resource information
func StoreDataResourceTable(db KeyValueStore, dataResourceTable *types.DataResourceTable) error {

	key := GetDataResourceTableKey(dataResourceTable.GetNodeId())
	val, err := rlp.EncodeToBytes(dataResourceTable)
	if nil != err {
		return err
	}
	return db.Put(key, val)
}

func StoreDataResourceTables(db KeyValueStore, dataResourceTables []*types.DataResourceTable) error {

	for _, dataResourceTable := range dataResourceTables {
		key := GetDataResourceTableKey(dataResourceTable.GetNodeId())
		val, err := rlp.EncodeToBytes(dataResourceTable)
		if nil != err {
			return err
		}
		if err := db.Put(key, val); nil != err {
			return err
		}
	}
	return nil
}

func RemoveDataResourceTable(db KeyValueStore, nodeId string) error {
	key := GetDataResourceTableKey(nodeId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func QueryDataResourceTable(db DatabaseReader, nodeId string) (*types.DataResourceTable, error) {
	key := GetDataResourceTableKey(nodeId)
	vb, err := db.Get(key)
	if nil != err {
		return nil, err
	}

	var dataResourceTable types.DataResourceTable

	if err := rlp.DecodeBytes(vb, &dataResourceTable); nil != err {
		return nil, err
	}
	return &dataResourceTable, nil
}

func QueryDataResourceTables(db KeyValueStore) ([]*types.DataResourceTable, error) {

	prefix := GetDataResourceTableKeyPrefix()
	it := db.NewIteratorWithPrefixAndStart(prefix, nil)
	defer it.Release()

	arr := make([]*types.DataResourceTable, 0)
	for it.Next() {
		if len(it.Key()) != 0 && len(it.Value()) != 0 {
			// prefix + dataNodeId -> LocalResourceTable
			var dataResourceTable types.DataResourceTable
			if err := rlp.DecodeBytes(it.Value(), &dataResourceTable); nil != err {
				return nil, err
			}
			arr = append(arr, &dataResourceTable)
		}
	}
	if len(arr) == 0 {
		return nil, ErrNotFound
	}
	return arr, nil
}

// The dataNode service information where the operation original file ID is located (originid - > {nodeid / metadataid / filepath})
func StoreDataResourceDataUpload(db KeyValueStore, dataResourceDataUpload *types.DataResourceDataUpload) error {

	key := GetDataResourceDataUploadKey(dataResourceDataUpload.GetOriginId())
	val, err := rlp.EncodeToBytes(dataResourceDataUpload)
	if nil != err {
		return err
	}
	return db.Put(key, val)
}

func StoreDataResourceDataUploads(db KeyValueStore, dataResourceDataUseds []*types.DataResourceDataUpload) error {

	for _, dataResourceDataUsed := range dataResourceDataUseds {
		key := GetDataResourceDataUploadKey(dataResourceDataUsed.GetOriginId())
		val, err := rlp.EncodeToBytes(dataResourceDataUsed)
		if nil != err {
			return err
		}

		if err := db.Put(key, val); nil != err {
			return err
		}
	}
	return nil
}

func RemoveDataResourceDataUpload(db KeyValueStore, originId string) error {

	key := GetDataResourceDataUploadKey(originId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func QueryDataResourceDataUpload(db DatabaseReader, originId string) (*types.DataResourceDataUpload, error) {
	key := GetDataResourceDataUploadKey(originId)
	vb, err := db.Get(key)
	if nil != err {
		return nil, err
	}

	var dataResourceDataUsed types.DataResourceDataUpload

	if err := rlp.DecodeBytes(vb, &dataResourceDataUsed); nil != err {
		return nil, err
	}
	return &dataResourceDataUsed, nil
}

func QueryDataResourceDataUploads(db KeyValueStore) ([]*types.DataResourceDataUpload, error) {

	prefix := GetDataResourceDataUploadKeyPrefix()
	it := db.NewIteratorWithPrefixAndStart(prefix, nil)
	defer it.Release()

	arr := make([]*types.DataResourceDataUpload, 0)
	for it.Next() {
		if len(it.Key()) != 0 && len(it.Value()) != 0 {
			// prefix + originId -> DataResourceDataUpload{originId, dataNodeId, metaDataId, filePath}
			var dataResourceDataUsed types.DataResourceDataUpload
			if err := rlp.DecodeBytes(it.Value(), &dataResourceDataUsed); nil != err {
				return nil, err
			}
			arr = append(arr, &dataResourceDataUsed)
		}
	}
	if len(arr) == 0 {
		return nil, ErrNotFound
	}
	return arr, nil
}

func StoreJobNodeIdByPowerId(db DatabaseWriter, powerId, jobNodeId string) error {
	key := GetPowerIdJobNodeIdMapingKey(powerId)
	index, err := rlp.EncodeToBytes(jobNodeId)
	if nil != err {
		return err
	}
	return db.Put(key, index)
}

func RemoveJobNodeIdByPowerId(db KeyValueStore, powerId string) error {
	key := GetPowerIdJobNodeIdMapingKey(powerId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func QueryJobNodeIdByPowerId(db DatabaseReader, powerId string) (string, error) {
	key := GetPowerIdJobNodeIdMapingKey(powerId)
	has, err := db.Has(key)
	if IsNoDBNotFoundErr(err) {
		return "", err
	}

	if !has {
		return "", ErrNotFound
	}
	idsByte, err := db.Get(key)
	if nil != err {
		return "", err
	}
	var jobNodeId string
	if err := rlp.DecodeBytes(idsByte, &jobNodeId); nil != err {
		return "", err
	}
	return jobNodeId, nil
}

func StoreDataResourceDiskUsed(db DatabaseWriter, dataResourceDiskUsed *types.DataResourceDiskUsed) error {
	key := GetDataResourceDiskUsedKey(dataResourceDiskUsed.GetMetadataId())
	val, err := rlp.EncodeToBytes(dataResourceDiskUsed)
	if nil != err {
		return err
	}
	return db.Put(key, val)
}

func RemoveDataResourceDiskUsed(db KeyValueStore, metadataId string) error {
	key := GetDataResourceDiskUsedKey(metadataId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func QueryDataResourceDiskUsed(db DatabaseReader, metadataId string) (*types.DataResourceDiskUsed, error) {
	key := GetDataResourceDiskUsedKey(metadataId)
	has, err := db.Has(key)
	if IsNoDBNotFoundErr(err) {
		return nil, err
	}

	if !has {
		return nil, ErrNotFound
	}
	vb, err := db.Get(key)
	if nil != err {
		return nil, err
	}
	var dataResourceDiskUsed types.DataResourceDiskUsed
	if err := rlp.DecodeBytes(vb, &dataResourceDiskUsed); nil != err {
		return nil, err
	}
	return &dataResourceDiskUsed, nil
}

func StoreLocalTaskExecuteStatusValConsByPartyId(db KeyValueStore, taskId, partyId string) error {

	key := GetLocalTaskExecuteStatus(taskId, partyId)

	v, err := db.Get(key)
	if IsNoDBNotFoundErr(err) {
		return err
	}

	val := bytesutil.BytesToUint32(v)
	val |= OnConsensusExecuteTaskStatus.Uint32()

	return db.Put(GetLocalTaskExecuteStatus(taskId, partyId), bytesutil.Uint32ToBytes(val))
}
func StoreLocalTaskExecuteStatusValExecByPartyId(db KeyValueStore, taskId, partyId string) error {
	key := GetLocalTaskExecuteStatus(taskId, partyId)

	v, err := db.Get(key)
	if IsNoDBNotFoundErr(err) {
		return err
	}

	val := bytesutil.BytesToUint32(v)
	val |= OnRunningExecuteStatus.Uint32()

	return db.Put(GetLocalTaskExecuteStatus(taskId, partyId), bytesutil.Uint32ToBytes(val))
}

func StoreLocalTaskExecuteStatusValTerminateByPartyId(db KeyValueStore, taskId, partyId string) error {
	key := GetLocalTaskExecuteStatus(taskId, partyId)

	v, err := db.Get(key)
	if IsNoDBNotFoundErr(err) {
		return err
	}

	val := bytesutil.BytesToUint32(v)
	val |= OnTerminingExecuteStatus.Uint32()

	return db.Put(GetLocalTaskExecuteStatus(taskId, partyId), bytesutil.Uint32ToBytes(val))
}

func RemoveLocalTaskExecuteStatusByPartyId(db KeyValueStore, taskId, partyId string) error {
	key := GetLocalTaskExecuteStatus(taskId, partyId) // prefix + taskId + partyId -> executeStatus (uint64)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func HasLocalTaskExecuteStatusParty(db KeyValueStore, taskId string) (bool, error) {
	prefix := append(localTaskExecuteStatusKeyPrefix, []byte(taskId)...)
	it := db.NewIteratorWithPrefixAndStart(prefix, nil)
	defer it.Release()

	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			has, err := db.Has(key)
			if IsNoDBNotFoundErr(err) {
				return false, err
			}
			// As long as there is a K-V existence,
			// it is a existence about task party execStatus
			if has {
				return true, nil
			}
		}
	}
	return false, nil
}

func HasLocalTaskExecuteStatusByPartyId(db DatabaseReader, taskId, partyId string) (bool, error) {
	key := GetLocalTaskExecuteStatus(taskId, partyId)
	has, err := db.Has(key)
	if IsNoDBNotFoundErr(err) {
		return false, err
	}
	if !has {
		return false, nil
	}
	return true, nil
}

func HasLocalTaskExecuteStatusConsensusByPartyId(db DatabaseReader, taskId, partyId string) (bool, error) {
	key := GetLocalTaskExecuteStatus(taskId, partyId)
	has, err := db.Has(key)
	if IsNoDBNotFoundErr(err) {
		return false, err
	}
	if !has {
		return false, nil
	}

	vb, err := db.Get(key)
	if nil != err {
		return false, err
	}
	if bytesutil.BytesToUint32(vb)&OnConsensusExecuteTaskStatus.Uint32() != OnConsensusExecuteTaskStatus.Uint32() {
		return false, nil
	}
	return true, nil
}

func HasLocalTaskExecuteStatusRunningByPartyId(db DatabaseReader, taskId, partyId string) (bool, error) {
	key := GetLocalTaskExecuteStatus(taskId, partyId)
	has, err := db.Has(key)
	if IsNoDBNotFoundErr(err) {
		return false, err
	}
	if !has {
		return false, nil
	}

	vb, err := db.Get(key)
	if nil != err {
		return false, err
	}
	if bytesutil.BytesToUint32(vb)&OnRunningExecuteStatus.Uint32() != OnRunningExecuteStatus.Uint32() {
		return false, nil
	}
	return true, nil
}

func HasLocalTaskExecuteStatusTerminateByPartyId(db DatabaseReader, taskId, partyId string) (bool, error) {
	key := GetLocalTaskExecuteStatus(taskId, partyId)
	has, err := db.Has(key)
	if IsNoDBNotFoundErr(err) {
		return false, err
	}
	if !has {
		return false, nil
	}

	vb, err := db.Get(key)
	if nil != err {
		return false, err
	}
	if bytesutil.BytesToUint32(vb)&OnTerminingExecuteStatus.Uint32() != OnTerminingExecuteStatus.Uint32() {
		return false, nil
	}
	return true, nil
}

// for metadataAuth
func StoreValidUserMetadataAuthStatusByMetadataId(db KeyValueStore, userType commonconstantpb.UserType, user, metadataId, metadataAuthId string, status uint16) error {

	// key: prefix + uesrType + user + metadataId + metadataAuthId -> value status <0: invalid, 1: valid>
	item_key := GetUserMetadataAuthStatusByMetadataIdKey(userType, user, metadataId, metadataAuthId)
	item_val, err := db.Get(item_key)

	var oldStatus uint16
	if IsNoDBNotFoundErr(err) {
		return err
	} else {
		oldStatus = bytesutil.BytesToUint16(item_val)
	}

	if status == oldStatus {
		return nil // It have been exists, don't store status.
	}

	statusVal := oldStatus

	// metadataAuth state, the low-order uint8.
	oldMetadataStatus := oldStatus &^ (0xFF << 8)
	metadataStatus := status &^ (0xFF << 8)
	// audit option, the high-order uint8.
	oldAuditOption := (oldStatus &^ 0xFF) >> 8
	auditOption := (status &^ 0xFF) >> 8

	// increase OR decrease valid count value.
	// 0: do nothing, -1: decrease, 1: increase
	var increaseORdecrease int

	switch {
	// has `unknown` old status for metadataAuth OR
	// has `invalid` old status for metadataAuth OR
	// has `revoke` old status for metadataAuth
	case oldMetadataStatus == uint16(commonconstantpb.MetadataAuthorityState_MAState_Unknown),
		oldMetadataStatus == uint16(commonconstantpb.MetadataAuthorityState_MAState_Invalid),
		oldMetadataStatus == uint16(commonconstantpb.MetadataAuthorityState_MAState_Revoked):
		return fmt.Errorf("old metadataAuth status has not `created` or `released`")
	case oldMetadataStatus == uint16(commonconstantpb.MetadataAuthorityState_MAState_Released) &&
		metadataStatus == uint16(commonconstantpb.MetadataAuthorityState_MAState_Created):
		return fmt.Errorf("can not change metadataStatus to `created` status from `released` status")
	// will change to `unknown` status for metadataAuth
	case metadataStatus == uint16(commonconstantpb.MetadataAuthorityState_MAState_Unknown):
		return fmt.Errorf("can not change metadataStatus to `unknown` status")
	// has not `pending` audit option.
	case oldAuditOption != uint16(commonconstantpb.AuditMetadataOption_Audit_Pending):
		return fmt.Errorf("old audit option has not `pending`")
	// will change to `invalid` status or`revoke` status for metadataAuth
	case metadataStatus == uint16(commonconstantpb.MetadataAuthorityState_MAState_Invalid),
		metadataStatus == uint16(commonconstantpb.MetadataAuthorityState_MAState_Revoked):
		increaseORdecrease = -1
	case oldMetadataStatus == uint16(commonconstantpb.MetadataAuthorityState_MAState_Created) &&
		metadataStatus == uint16(commonconstantpb.MetadataAuthorityState_MAState_Released):
		// do nothing...
	default: // first set `created` OR `released` status into.
		increaseORdecrease = 1
	}

	// new value: high-order uint8 + low-order uint8
	statusVal = auditOption + metadataStatus

	if increaseORdecrease != 0 {

		var valid_count uint32
		// key: prefix + userType + user + metadataId ->  value: metadataAuth validCount
		valid_count_key := GetValidUserMetadataAuthCountByMetadataIdKey(userType, user, metadataId)
		valid_count_val, verr := db.Get(valid_count_key)
		switch {
		case IsNoDBNotFoundErr(verr):
			return verr
		case IsDBNotFoundErr(verr):
			// do nothing
		case nil == verr && len(valid_count_val) != 0:
			valid_count = bytesutil.BytesToUint32(valid_count_val)
		}

		if increaseORdecrease == 1 {
			// When metadataAuthId have not by metadataId, inscrease total count of metadataAuth and valid count of metadataAuth
			// and put metadataAuth status `valid`.
			//
			// key: prefix + userType + user + metadataId ->  value: metadataAuth totalCount
			total_count_key := GetTotalUserMetadataAuthCountByMetadataIdKey(userType, user, metadataId)
			total_count_val, terr := db.Get(total_count_key)

			var total_count uint32
			switch {
			case IsNoDBNotFoundErr(terr):
				return terr
			case IsDBNotFoundErr(terr):
				// do nothing
			case nil == terr && len(total_count_val) != 0:
				total_count = bytesutil.BytesToUint32(total_count_val)
			}
			total_count++
			total_count_val = bytesutil.Uint32ToBytes(total_count)

			// inscease total count of metadataAuth and valid count of metadataAuth.
			if err := db.Put(total_count_key, total_count_val); nil != err {
				return fmt.Errorf("update metadataAuth totalCount %s", err)
			}

			valid_count++
			valid_count_val = bytesutil.Uint32ToBytes(valid_count)

			if err := db.Put(valid_count_key, valid_count_val); nil != err {
				return fmt.Errorf("update metadataAuth validCount %s", err)
			}

		} else {
			if valid_count > 0 {
				valid_count--
			}
			// Second: descease or delete valid count of metadataAuth.
			if valid_count == 0 {
				if err := db.Delete(valid_count_key); nil != err {
					return fmt.Errorf("delete metadataAuth validCount %s", err)
				}
			} else {
				valid_count_val = bytesutil.Uint32ToBytes(valid_count)
				if err := db.Put(valid_count_key, valid_count_val); nil != err {
					return fmt.Errorf("update metadataAuth validCount %s", err)
				}
			}
		}
	}

	// update metadataStatus into metadataAauth mapping.
	if err := db.Put(item_key, bytesutil.Uint16ToBytes(statusVal)); nil != err {
		return fmt.Errorf("update metadataAuth status %s", err)
	}

	return nil
}

func HasValidUserMetadataAuthStatusByMetadataId(db DatabaseReader, userType commonconstantpb.UserType, user, metadataId string) (bool, error) {

	key := GetValidUserMetadataAuthCountByMetadataIdKey(userType, user, metadataId)
	val, err := db.Get(key)

	var count uint32

	switch {
	case IsNoDBNotFoundErr(err):
		return false, err
	case IsDBNotFoundErr(err):
		return false, nil
	case nil == err && len(val) != 0:
		count = bytesutil.BytesToUint32(val)
	}
	return count > 0, nil
}

func HasUserMetadataAuthStatusByMetadataId(db DatabaseReader, userType commonconstantpb.UserType, user, metadataId string) (bool, error) {

	key := GetTotalUserMetadataAuthCountByMetadataIdKey(userType, user, metadataId)
	val, err := db.Get(key)

	var count uint32

	switch {
	case IsNoDBNotFoundErr(err):
		return false, err
	case IsDBNotFoundErr(err):
		return false, nil
	case nil == err && len(val) != 0:
		count = bytesutil.BytesToUint32(val)
	}
	return count > 0, nil
}

func HasUserMetadataAuthIdByMetadataId(db DatabaseReader, userType commonconstantpb.UserType, user, metadataId, metadataAuthId string) (bool, error) {
	// key: prefix + uesrType + user + metadataId + metadataAuthId -> value status <0: invalid, 1: valid>
	item_key := GetUserMetadataAuthStatusByMetadataIdKey(userType, user, metadataId, metadataAuthId)
	has, err := db.Has(item_key)
	switch {
	case IsNoDBNotFoundErr(err):
		return false, err
	case IsDBNotFoundErr(err):
		return false, nil
	case nil == err && !has:
		return false, nil
	}
	return true, nil
}

func QueryTotalUserMetadataAuthCountByMetadataId(db DatabaseReader, userType commonconstantpb.UserType, user, metadataId string) (uint32, error) {

	key := GetTotalUserMetadataAuthCountByMetadataIdKey(userType, user, metadataId)
	val, err := db.Get(key)

	var count uint32

	switch {
	case IsNoDBNotFoundErr(err):
		return 0, err
	case IsDBNotFoundErr(err):
		return 0, nil
	case nil == err && len(val) != 0:
		count = bytesutil.BytesToUint32(val)
	}
	return count, nil
}

func QueryValidUserMetadataAuthCountByMetadataId(db DatabaseReader, userType commonconstantpb.UserType, user, metadataId string) (uint32, error) {

	key := GetValidUserMetadataAuthCountByMetadataIdKey(userType, user, metadataId)
	val, err := db.Get(key)

	var count uint32

	switch {
	case IsNoDBNotFoundErr(err):
		return 0, err
	case IsDBNotFoundErr(err):
		return 0, nil
	case nil == err && len(val) != 0:
		count = bytesutil.BytesToUint32(val)
	}
	return count, nil
}

func QueryValidUserMetadataAuthIdsByMetadataId(db KeyValueStore, userType commonconstantpb.UserType, user, metadataId string) ([]string, error) {

	// prefix + uesrType + user + metadataId + metadataAuthId -> status <0: invalid, 1: valid>
	prefixAndMetadataId := GetUserMetadataAuthStatusKeyPrefixByMetadataId(userType, user, metadataId)
	it := db.NewIteratorWithPrefixAndStart(prefixAndMetadataId, nil)
	defer it.Release()

	arr := make([]string, 0)
	tmp := make(map[string]struct{}, 0)
	for it.Next() {
		if len(it.Key()) != 0 && len(it.Value()) != 0 {
			// key len == len(prefix) + len([]byte(uesrType)) + len([]byte(user)) + len([]byte(metadataId)) + len([]byte(metadataAuthId))
			metadataAuthId := string(it.Key()[len(prefixAndMetadataId):])
			if _, ok := tmp[metadataAuthId]; !ok {
				tmp[metadataAuthId] = struct{}{}
				arr = append(arr, metadataAuthId)
			}
		}
	}
	if len(arr) == 0 {
		return nil, ErrNotFound
	}

	return arr, nil
}

func RemoveUserMetadataAuthStatusByMetadataId(db KeyValueStore, userType commonconstantpb.UserType, user, metadataId, metadataAuthId string) error {

	// key: prefix + uesrType + user + metadataId + metadataAuthId -> value status <0: invalid, 1: valid>
	item_key := GetUserMetadataAuthStatusByMetadataIdKey(userType, user, metadataId, metadataAuthId)

	// When metadataAuthId have not by metadataId, descrease total count of metadataAuth and valid count of metadataAuth
	// and put metadataAuth status `valid`.
	//
	// key: prefix + userType + user + metadataId ->  value: metadataAuth totalCount
	total_count_key := GetTotalUserMetadataAuthCountByMetadataIdKey(userType, user, metadataId)
	total_count_val, terr := db.Get(total_count_key)
	// key: prefix + userType + user + metadataId ->  value: metadataAuth validCount
	valid_count_key := GetValidUserMetadataAuthCountByMetadataIdKey(userType, user, metadataId)
	valid_count_val, verr := db.Get(valid_count_key)

	var (
		total_count uint32
		valid_count uint32
	)
	switch {
	case IsNoDBNotFoundErr(terr):
		return terr
	case IsDBNotFoundErr(terr):
		// do nothing
	case nil == terr && len(total_count_val) != 0:
		total_count = bytesutil.BytesToUint32(total_count_val)
	}
	switch {
	case IsNoDBNotFoundErr(verr):
		return verr
	case IsDBNotFoundErr(verr):
		// do nothing
	case nil == verr && len(valid_count_val) != 0:
		valid_count = bytesutil.BytesToUint32(valid_count_val)
	}

	if total_count > 0 {
		total_count--
	}
	if valid_count > 0 {
		valid_count--
	}

	// First: delete metadataStatus `invalid` on metadataAauth mapping.
	if err := db.Delete(item_key); nil != err {
		return err
	}

	// Second: decscease or delete total count of metadataAuth and valid count of metadataAuth.
	if total_count == 0 {
		if err := db.Delete(total_count_key); nil != err {
			return fmt.Errorf("delete metadataAuth totalCount %s", err)
		}
	} else {
		total_count_val = bytesutil.Uint32ToBytes(total_count)
		if err := db.Put(total_count_key, total_count_val); nil != err {
			return fmt.Errorf("update metadataAuth totalCount %s", err)
		}
	}
	if valid_count == 0 {
		if err := db.Delete(valid_count_key); nil != err {
			return fmt.Errorf("delete metadataAuth validCount %s", err)
		}
	} else {
		valid_count_val = bytesutil.Uint32ToBytes(valid_count)
		if err := db.Put(valid_count_key, valid_count_val); nil != err {
			return fmt.Errorf("update metadataAuth validCount %s", err)
		}
	}
	return nil
}

// about metadata history used task.
func StoreMetadataHistoryTaskId(db KeyValueStore, metadataId, taskId string) error {
	// prefix + metadataId + taskId -> index
	item_key := GetMetadataHistoryTaskKey(metadataId, taskId)
	has, err := db.Has(item_key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case nil == err && has:
		return nil // It have been exists, don't inscrease count
	}

	// When taskId have not by metadata, inscrease metadata used taskId count
	// and put taskId on metadataId mapping.
	//
	// prefix + metadataId -> history task count
	count_key := GetMetadataHistoryTaskCountKey(metadataId)
	count_val, err := db.Get(count_key)

	var count uint32

	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err):
		// do nothing
	case nil == err && len(count_val) != 0:
		count = bytesutil.BytesToUint32(count_val)
	}
	count++

	count_val = bytesutil.Uint32ToBytes(count)

	// First: put taskId on metadata mapping.
	if err := db.Put(item_key, count_val); nil != err {
		return err
	}
	//log.Debugf("InscreaseMetadataHistoryTaskCount, metadataId: {%s}, taskId: {%s}, count: {%d}", metadataId, taskId, count)
	// Second: inscease taskId count by metadata.
	return db.Put(count_key, count_val)
}

func HasMetadataHistoryTaskId(db DatabaseReader, metadataId, taskId string) (bool, error) {

	item_key := GetMetadataHistoryTaskKey(metadataId, taskId)
	has, err := db.Has(item_key)
	switch {
	case IsNoDBNotFoundErr(err):
		return false, err
	case IsDBNotFoundErr(err):
		return false, nil
	case nil == err && !has:
		return false, nil
	}
	return true, nil
}

func QueryMetadataHistoryTaskIdCount(db DatabaseReader, metadataId string) (uint32, error) {
	// prefix + metadataId -> history task count
	key := GetMetadataHistoryTaskCountKey(metadataId)
	val, err := db.Get(key)

	var count uint32

	switch {
	case IsNoDBNotFoundErr(err):
		return 0, err
	case IsDBNotFoundErr(err):
		return 0, nil
	case nil == err && len(val) != 0:
		count = bytesutil.BytesToUint32(val)
	}
	return count, nil
}

func QueryMetadataHistoryTaskIds(db KeyValueStore, metadataId string) ([]string, error) {
	// prefix + metadataId + taskId -> index
	prefixAndMetadataId := GetMetadataHistoryTaskKeyPrefixByMetadataId(metadataId)
	it := db.NewIteratorWithPrefixAndStart(prefixAndMetadataId, nil)
	defer it.Release()

	arr := make([]string, 0)
	tmp := make(map[string]struct{}, 0)
	for it.Next() {
		if len(it.Key()) != 0 && len(it.Value()) != 0 {
			// key len == len(prefix) + len([]byte(metadataId)) + len([]byte(taskId))
			taskId := string(it.Key()[len(prefixAndMetadataId):])
			if _, ok := tmp[taskId]; !ok {
				tmp[taskId] = struct{}{}
				arr = append(arr, taskId)
			}
		}
	}
	if len(arr) == 0 {
		return nil, ErrNotFound
	}

	return arr, nil
}

func StoreTaskUpResultData(db DatabaseWriter, turf *types.TaskUpResultData) error {
	key := GetTaskResultDataMetadataIdKey(turf.GetTaskId())
	val, err := rlp.EncodeToBytes(turf)
	if nil != err {
		return err
	}
	return db.Put(key, val)
}

func QueryTaskUpResultData(db DatabaseReader, taskId string) (*types.TaskUpResultData, error) {
	key := GetTaskResultDataMetadataIdKey(taskId)
	vb, err := db.Get(key)
	if nil != err {
		return nil, err
	}
	var taskUpResultData types.TaskUpResultData
	if err = rlp.DecodeBytes(vb, &taskUpResultData); nil != err {
		return nil, err
	}
	return &taskUpResultData, nil
}

func QueryTaskUpResultDataList(db DatabaseIteratee) ([]*types.TaskUpResultData, error) {

	it := db.NewIteratorWithPrefixAndStart(GetTaskResultDataMetadataIdKeyPrefix(), nil)
	defer it.Release()

	arr := make([]*types.TaskUpResultData, 0)
	for it.Next() {
		if value := it.Value(); len(value) != 0 {
			var taskUpResultData types.TaskUpResultData
			if err := rlp.DecodeBytes(value, &taskUpResultData); nil != err {
				log.WithError(err).Errorf("Failed to call QueryAllTaskUpResultData, decode db val failed")
				continue
			}
			arr = append(arr, &taskUpResultData)
		}
	}

	if len(arr) == 0 {
		return nil, ErrNotFound
	}

	return arr, nil
}

func RemoveTaskUpResultData(db KeyValueStore, taskId string) error {
	key := GetTaskResultDataMetadataIdKey(taskId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func StoreTaskPartnerPartyIds(db DatabaseWriter, taskId string, partyIds []string) error {
	key := GetTaskPartnerPartyIdsKey(taskId)
	val, err := rlp.EncodeToBytes(partyIds)
	if nil != err {
		return err
	}
	return db.Put(key, val)
}

func HasTaskPartnerPartyIds(db DatabaseReader, taskId string) (bool, error) {
	key := GetTaskPartnerPartyIdsKey(taskId)

	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return false, err
	case IsDBNotFoundErr(err), !has:
		return false, nil
	}
	return true, nil
}

func QueryTaskPartnerPartyIds(db DatabaseReader, taskId string) ([]string, error) {
	key := GetTaskPartnerPartyIdsKey(taskId)

	vb, err := db.Get(key)
	if nil != err {
		return nil, err
	}
	var partyIdArr []string
	if err = rlp.DecodeBytes(vb, &partyIdArr); nil != err {
		return nil, err
	}
	return partyIdArr, nil
}

func RemoveTaskPartnerPartyId(db KeyValueStore, taskId, partyId string) error {
	key := GetTaskPartnerPartyIdsKey(taskId)
	vb, err := db.Get(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && len(vb) == 0:
		return nil
	}

	var partyIdArr []string
	if err = rlp.DecodeBytes(vb, &partyIdArr); nil != err {
		return err
	}

	for i, id := range partyIdArr {
		if id == partyId {
			partyIdArr = append(partyIdArr[:i], partyIdArr[i+1:]...)
			break
		}
	}
	if len(partyIdArr) == 0 {
		return db.Delete(key)
	}
	vb, err = rlp.EncodeToBytes(partyIdArr)
	if nil != err {
		return err
	}
	return db.Put(key, vb)
}

func RemoveTaskPartnerPartyIds(db KeyValueStore, taskId string) error {
	key := GetTaskPartnerPartyIdsKey(taskId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func StoreMessageCache(db KeyValueStore, value interface{}) error {
	var (
		key []byte
		val []byte
		err error
	)
	switch v := value.(type) {
	case *types.PowerMsg:

		key = GetPowerMsgKey(v.GetPowerId())
		val, err = proto.Marshal(&carriertypespb.PowerMsg{
			PowerId:   v.GetPowerId(),
			JobNodeId: v.GetJobNodeId(),
			CreateAt:  v.GetCreateAt(),
		})
		if nil != err {
			return fmt.Errorf("marshal powerMsg failed, %s", err)
		}

	case *types.MetadataMsg:
		key = GetMetadataMsgKey(v.GetMetadataId())
		val, err = proto.Marshal(&carriertypespb.MetadataMsg{
			MetadataId:      v.GetMetadataId(),
			MetadataSummary: v.GetMetadataSummary(),
			CreateAt:        v.GetCreateAt(),
		})
		if nil != err {
			return fmt.Errorf("marshal metadataMsg failed, %s", err)
		}
	case *types.MetadataAuthorityMsg:
		key = GetMetadataAuthMsgKey(v.GetMetadataAuthId())
		val, err = proto.Marshal(&carriertypespb.MetadataAuthorityMsg{
			MetadataAuthId: v.GetMetadataAuthId(),
			User:           v.GetUser(),
			UserType:       v.GetUserType(),
			Auth:           v.GetMetadataAuthority(),
			Sign:           v.GetSign(),
			CreateAt:       v.GetCreateAt(),
		})
		if nil != err {
			return fmt.Errorf("marshal metadataAuthorityMsg failed, %s", err)
		}
	case *types.TaskMsg:
		key = GetTaskMsgKey(v.GetTaskId())
		val, err = proto.Marshal(&carriertypespb.TaskMsg{
			Data: v.GetTaskData(),
		})
		if nil != err {
			return fmt.Errorf("marshal taskMsg failed, %s", err)
		}
	case types.MetadataUpdateMsg:
		key = GetMetadataUpdateMsgKey(v.GetMetadataId())
		val, err = json.Marshal(&types.MetadataUpdateMsg{
			MetadataSummary: v.GetMetadataSummary(),
			CreateAt:        v.GetCreateAt(),
		})
		if nil != err {
			return fmt.Errorf("marshal metadataMsg failed, %s", err)
		}
	}
	return db.Put(key, val)
}

func RemovePowerMsg(db KeyValueStore, powerId string) error {
	key := GetPowerMsgKey(powerId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func RemoveAllPowerMsg(db KeyValueStore) error {
	it := db.NewIteratorWithPrefixAndStart(GetPowerMsgKeyPrefix(), nil)
	defer it.Release()

	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			db.Delete(key)
		}
	}
	return nil
}

func RemoveMetadataMsg(db KeyValueStore, metadataId string) error {
	key := GetMetadataMsgKey(metadataId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func RemoveMetadataUpdateMsg(db KeyValueStore, metadataId string) error {
	key := GetMetadataUpdateMsgKey(metadataId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func RemoveAllMetadataMsg(db KeyValueStore) error {
	it := db.NewIteratorWithPrefixAndStart(GetMetadataMsgKeyPrefix(), nil)
	defer it.Release()

	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			db.Delete(key)
		}
	}
	return nil
}

func RemoveMetadataAuthMsg(db KeyValueStore, metadataAuthId string) error {
	key := GetMetadataAuthMsgKey(metadataAuthId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func RemoveAllMetadataAuthMsg(db KeyValueStore) error {
	it := db.NewIteratorWithPrefixAndStart(GetMetadataAuthMsgKeyPrefix(), nil)
	defer it.Release()

	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			db.Delete(key)
		}
	}
	return nil
}

func RemoveTaskMsg(db KeyValueStore, taskId string) error {
	key := GetTaskMsgKey(taskId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func RemoveAllTaskMsg(db KeyValueStore) error {
	it := db.NewIteratorWithPrefixAndStart(GetTaskMsgKeyPrefix(), nil)
	defer it.Release()

	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			db.Delete(key)
		}
	}
	return nil
}

func QueryPowerMsgArr(db KeyValueStore) (types.PowerMsgArr, error) {

	it := db.NewIteratorWithPrefixAndStart(GetPowerMsgKeyPrefix(), nil)
	defer it.Release()

	arr := make(types.PowerMsgArr, 0)

	for it.Next() {
		if val := it.Value(); len(val) != 0 {
			var res carriertypespb.PowerMsg
			if err := proto.Unmarshal(val, &res); nil != err {
				continue
			}
			arr = append(arr, &types.PowerMsg{
				PowerId:   res.GetPowerId(),
				JobNodeId: res.GetJobNodeId(),
				CreateAt:  res.GetCreateAt(),
			})
		}
	}
	if len(arr) == 0 {
		return nil, ErrNotFound
	}
	return arr, nil
}

func QueryMetadataMsgArr(db KeyValueStore) (types.MetadataMsgArr, error) {
	it := db.NewIteratorWithPrefixAndStart(GetMetadataMsgKeyPrefix(), nil)
	defer it.Release()

	arr := make(types.MetadataMsgArr, 0)

	for it.Next() {
		if val := it.Value(); len(val) != 0 {
			var res carriertypespb.MetadataMsg
			if err := proto.Unmarshal(val, &res); nil != err {
				continue
			}
			arr = append(arr, &types.MetadataMsg{
				MetadataSummary: res.GetMetadataSummary(),
				CreateAt:        res.GetCreateAt(),
			})
		}
	}
	if len(arr) == 0 {
		return nil, ErrNotFound
	}
	return arr, nil
}

func QueryMetadataUpdateMsgArr(db KeyValueStore) (types.MetadataUpdateMsgArr, error) {
	it := db.NewIteratorWithPrefixAndStart(GetMetadataMsgKeyPrefix(), nil)
	defer it.Release()

	arr := make(types.MetadataUpdateMsgArr, 0)

	for it.Next() {
		if val := it.Value(); len(val) != 0 {
			var res types.MetadataUpdateMsg
			if err := json.Unmarshal(val, &res); nil != err {
				continue
			}
			arr = append(arr, &types.MetadataUpdateMsg{
				MetadataSummary: res.GetMetadataSummary(),
				CreateAt:        res.GetCreateAt(),
			})
		}
	}
	if len(arr) == 0 {
		return nil, ErrNotFound
	}
	return arr, nil
}

func QueryMetadataAuthorityMsgArr(db KeyValueStore) (types.MetadataAuthorityMsgArr, error) {
	it := db.NewIteratorWithPrefixAndStart(GetMetadataAuthMsgKeyPrefix(), nil)
	defer it.Release()

	arr := make(types.MetadataAuthorityMsgArr, 0)

	for it.Next() {
		if val := it.Value(); len(val) != 0 {
			var res carriertypespb.MetadataAuthorityMsg
			if err := proto.Unmarshal(val, &res); nil != err {
				continue
			}
			arr = append(arr, &types.MetadataAuthorityMsg{
				MetadataAuthId: res.GetMetadataAuthId(),
				User:           res.GetUser(),
				UserType:       res.GetUserType(),
				Auth:           res.GetAuth(),
				Sign:           res.GetSign(),
				CreateAt:       res.GetCreateAt(),
			})
		}
	}
	if len(arr) == 0 {
		return nil, ErrNotFound
	}
	return arr, nil
}

func QueryTaskMsgArr(db KeyValueStore) (types.TaskMsgArr, error) {
	it := db.NewIteratorWithPrefixAndStart(GetMetadataAuthMsgKeyPrefix(), nil)
	defer it.Release()

	arr := make(types.TaskMsgArr, 0)

	for it.Next() {
		if val := it.Value(); len(val) != 0 {
			var res carriertypespb.TaskMsg
			if err := proto.Unmarshal(val, &res); nil != err {
				continue
			}
			arr = append(arr, &types.TaskMsg{
				Data: types.NewTask(res.GetData()),
			})
		}
	}
	if len(arr) == 0 {
		return nil, ErrNotFound
	}
	return arr, nil
}

func SaveOrgPriKey(db db.Database, priKey string) error {
	key := GetOrgPriKeyPrefix()
	val, err := rlp.EncodeToBytes(priKey)
	if nil != err {
		return err
	}
	return db.Put(key, val)
}

// FindOrgPriKey does not return ErrNotFound if the organization private key not found.
func FindOrgPriKey(db DatabaseReader) (string, error) {
	key := GetOrgPriKeyPrefix()
	if has, err := db.Has(key); err != nil {
		return "", err
	} else if has {
		if val, err := db.Get(key); err != nil {
			return "", err
		} else {
			var priKey string
			if err := rlp.DecodeBytes(val, &priKey); err != nil {
				return "", err
			} else {
				return priKey, nil
			}
		}
	}
	return "", nil
}

func QueryIdentityMsgNonce(db DatabaseReader) (uint64, error) {

	val, err := db.Get(GetIdentityMsgNonceKey())

	var nonce uint64
	switch {
	case IsNoDBNotFoundErr(err):
		return 0, err
	case IsDBNotFoundErr(err):
		return 0, nil
	case nil == err && len(val) != 0:
		nonce = bytesutil.BytesToUint64(val)
	}
	return nonce, nil
}
func IncreaseIdentityMsgNonce(db KeyValueStore) (uint64, error) {

	key := GetIdentityMsgNonceKey()
	val, err := db.Get(key)

	var nonce uint64
	switch {
	case IsNoDBNotFoundErr(err):
		return 0, err
	case IsDBNotFoundErr(err):
		// do nothing
	case nil == err && len(val) != 0:
		nonce = bytesutil.BytesToUint64(val)
	}
	nonce++

	return nonce, db.Put(key, bytesutil.Uint64ToBytes(nonce))
}
func QueryMetadataMsgNonce(db DatabaseReader) (uint64, error) {

	val, err := db.Get(GetMetadataMsgNonceKey())

	var nonce uint64
	switch {
	case IsNoDBNotFoundErr(err):
		return 0, err
	case IsDBNotFoundErr(err):
		return 0, nil
	case nil == err && len(val) != 0:
		nonce = bytesutil.BytesToUint64(val)
	}
	return nonce, nil
}
func IncreaseMetadataMsgNonce(db KeyValueStore) (uint64, error) {

	key := GetMetadataMsgNonceKey()
	val, err := db.Get(key)

	var nonce uint64
	switch {
	case IsNoDBNotFoundErr(err):
		return 0, err
	case IsDBNotFoundErr(err):
		// do nothing
	case nil == err && len(val) != 0:
		nonce = bytesutil.BytesToUint64(val)
	}
	nonce++

	return nonce, db.Put(key, bytesutil.Uint64ToBytes(nonce))
}
func QueryMetadataAuthMsgNonce(db DatabaseReader) (uint64, error) {

	val, err := db.Get(GetMetadataAuthMsgNonceKey())

	var nonce uint64
	switch {
	case IsNoDBNotFoundErr(err):
		return 0, err
	case IsDBNotFoundErr(err):
		return 0, nil
	case nil == err && len(val) != 0:
		nonce = bytesutil.BytesToUint64(val)
	}
	return nonce, nil
}
func IncreaseMetadataAuthMsgNonce(db KeyValueStore) (uint64, error) {

	key := GetMetadataAuthMsgNonceKey()
	val, err := db.Get(key)

	var nonce uint64
	switch {
	case IsNoDBNotFoundErr(err):
		return 0, err
	case IsDBNotFoundErr(err):
		// do nothing
	case nil == err && len(val) != 0:
		nonce = bytesutil.BytesToUint64(val)
	}
	nonce++

	return nonce, db.Put(key, bytesutil.Uint64ToBytes(nonce))
}
func QueryPowerMsgNonce(db DatabaseReader) (uint64, error) {

	val, err := db.Get(GetPowerMsgNonceKey())

	var nonce uint64
	switch {
	case IsNoDBNotFoundErr(err):
		return 0, err
	case IsDBNotFoundErr(err):
		return 0, nil
	case nil == err && len(val) != 0:
		nonce = bytesutil.BytesToUint64(val)
	}
	return nonce, nil
}
func IncreasePowerMsgNonce(db KeyValueStore) (uint64, error) {

	key := GetPowerMsgNonceKey()
	val, err := db.Get(key)

	var nonce uint64
	switch {
	case IsNoDBNotFoundErr(err):
		return 0, err
	case IsDBNotFoundErr(err):
		// do nothing
	case nil == err && len(val) != 0:
		nonce = bytesutil.BytesToUint64(val)
	}
	nonce++

	return nonce, db.Put(key, bytesutil.Uint64ToBytes(nonce))
}
func QueryTaskMsgNonce(db DatabaseReader) (uint64, error) {

	val, err := db.Get(GetTaskMsgNonceKey())

	var nonce uint64
	switch {
	case IsNoDBNotFoundErr(err):
		return 0, err
	case IsDBNotFoundErr(err):
		return 0, nil
	case nil == err && len(val) != 0:
		nonce = bytesutil.BytesToUint64(val)
	}
	return nonce, nil
}
func IncreaseTaskMsgNonce(db KeyValueStore) (uint64, error) {

	key := GetTaskMsgNonceKey()
	val, err := db.Get(key)

	var nonce uint64
	switch {
	case IsNoDBNotFoundErr(err):
		return 0, err
	case IsDBNotFoundErr(err):
		// do nothing
	case nil == err && len(val) != 0:
		nonce = bytesutil.BytesToUint64(val)
	}
	nonce++

	return nonce, db.Put(key, bytesutil.Uint64ToBytes(nonce))
}
