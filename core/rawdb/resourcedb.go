package rawdb

import (
	"github.com/RosettaFlow/Carrier-Go/common/bytesutil"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/ethereum/go-ethereum/rlp"
	leveldberr "github.com/syndtr/goleveldb/leveldb/errors"
)

var (
	//ErrNotFound = errors.New("rawdb: not found")
	//ErrLeveldbNotFound = leveldberr.ErrNotFound

	ErrNotFound = leveldberr.ErrNotFound
)

func IsNoDBNotFoundErr(err error) bool {
	return nil != err && err != ErrNotFound
}
func IsDBNotFoundErr(err error) bool {
	return nil != err && err == ErrNotFound
}

// 操作 本组织 计算服务的资源
func StoreNodeResource(db KeyValueStore, resource *types.LocalResourceTable) error {

	key := GetNodeResourceKey(resource.GetNodeId())
	val, err := rlp.EncodeToBytes(resource)
	if nil != err {
		return err
	}

	if err := db.Put(key, val); nil != err {
		return err
	}

	has, err := db.Has(GetNodeResourceIdListKey())
	if IsNoDBNotFoundErr(err) {
		return err
	}

	var resourceIds []string
	if !has {
		resourceIds = []string{resource.GetNodeId()}
	} else {
		idsByte, err := db.Get(GetNodeResourceIdListKey())
		if nil != err {
			return err
		}
		if err := rlp.DecodeBytes(idsByte, &resourceIds); nil != err {
			return err
		}

		var include bool

		for _, id := range resourceIds {
			if id == resource.GetNodeId() {
				include = true
				break
			}
		}
		if !include {
			resourceIds = append(resourceIds, resource.GetNodeId())
		}
	}

	index, err := rlp.EncodeToBytes(resourceIds)
	if nil != err {
		return err
	}

	return db.Put(GetNodeResourceIdListKey(), index)
}

func StoreNodeResources(db KeyValueStore, resources []*types.LocalResourceTable) error {

	has, err := db.Has(GetNodeResourceIdListKey())
	if IsNoDBNotFoundErr(err) {
		return err
	}

	inputIds := make([]string, len(resources))
	for i, re := range resources {
		inputIds[i] =  re.GetNodeId()
		key := GetNodeResourceKey(re.GetNodeId())
		val, err := rlp.EncodeToBytes(re)
		if nil != err {
			return err
		}

		if err := db.Put(key, val); nil != err {
			return err
		}
	}

	var resourceIds []string
	if !has {
		resourceIds = inputIds
	} else {
		idsByte, err := db.Get(GetNodeResourceIdListKey())
		if nil != err {
			return err
		}
		if err := rlp.DecodeBytes(idsByte, &resourceIds); nil != err {
			return err
		}

		tmp := make(map[string]struct{})

		for _, id := range resourceIds {
			tmp[id] = struct{}{}
		}
		for _, id := range inputIds {
			if _, ok := tmp[id]; !ok {
				resourceIds = append(resourceIds, id)
			}
		}

	}

	index, err := rlp.EncodeToBytes(resourceIds)
	if nil != err {
		return err
	}

	return db.Put(GetNodeResourceIdListKey(), index)
}

func RemoveNodeResource (db KeyValueStore, resourceId string) error {
	has, err := db.Has(GetNodeResourceIdListKey())
	if IsNoDBNotFoundErr(err) {
		return err
	}

	var resourceIds []string
	if !has {
		return nil
	} else {
		idsByte, err := db.Get(GetNodeResourceIdListKey())
		if nil != err {
			return err
		}
		if err := rlp.DecodeBytes(idsByte, &resourceIds); nil != err {
			return err
		}


		for i := 0; i < len(resourceIds); i++ {
			id := resourceIds[i]
			if id == resourceId {
				key := GetNodeResourceKey(resourceId)
				if err := db.Delete(key); nil != err {
					return err
				}
				resourceIds = append(resourceIds[:i], resourceIds[i+1:]...)
				i--
				break
			}
		}
	}

	index, err := rlp.EncodeToBytes(resourceIds)
	if nil != err {
		return err
	}

	return db.Put(GetNodeResourceIdListKey(), index)
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

func QueryNodeResources (db DatabaseReader) ([]*types.LocalResourceTable, error) {
	has, err := db.Has(GetNodeResourceIdListKey())
	if IsNoDBNotFoundErr(err) {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	b, err := db.Get(GetNodeResourceIdListKey())
	if nil != err {
		return nil, err
	}
	var ids []string
	if err := rlp.DecodeBytes(b, &ids); nil != err {
		return nil, err
	}

	arr := make([]*types.LocalResourceTable, len(ids))
	for i, id := range ids {

		key := GetNodeResourceKey(id)
		vb, err := db.Get(key)
		if nil != err {
			return nil, err
		}

		var resource types.LocalResourceTable
		if err := rlp.DecodeBytes(vb, &resource); nil != err {
			return nil, err
		}
		arr[i] = &resource
	}

	return arr, nil
}

// 操作 全网组织 算力资源
func StoreOrgResource(db KeyValueStore, resource *types.RemoteResourceTable) error {
	key := GetOrgResourceKey(resource.GetIdentityId())
	val, err := rlp.EncodeToBytes(resource)
	if nil != err {
		return err
	}

	if err := db.Put(key, val); nil != err {
		return err
	}

	has, err := db.Has(GetOrgResourceIdListKey())
	if IsNoDBNotFoundErr(err) {
		return err
	}

	var resourceIds []string
	if !has {
		resourceIds = []string{resource.GetIdentityId()}
	} else {
		idsByte, err := db.Get(GetOrgResourceIdListKey())
		if nil != err {
			return err
		}
		if err := rlp.DecodeBytes(idsByte, &resourceIds); nil != err {
			return err
		}

		var include bool

		for _, id := range resourceIds {
			if id == resource.GetIdentityId() {
				include = true
				break
			}
		}
		if !include {
			resourceIds = append(resourceIds, resource.GetIdentityId())
		}
	}

	index, err := rlp.EncodeToBytes(resourceIds)
	if nil != err {
		return err
	}

	return db.Put(GetOrgResourceIdListKey(), index)
}

func StoreOrgResources(db KeyValueStore, resources []*types.RemoteResourceTable) error {

	has, err := db.Has(GetOrgResourceIdListKey())
	if IsNoDBNotFoundErr(err) {
		return err
	}

	inputIds := make([]string, len(resources))
	for i, re := range resources {
		inputIds[i] =  re.GetIdentityId()
		key := GetOrgResourceKey(re.GetIdentityId())
		val, err := rlp.EncodeToBytes(re)
		if nil != err {
			return err
		}

		if err := db.Put(key, val); nil != err {
			return err
		}
	}

	var resourceIds []string
	if !has {
		resourceIds = inputIds
	} else {
		idsByte, err := db.Get(GetOrgResourceIdListKey())
		if nil != err {
			return err
		}
		if err := rlp.DecodeBytes(idsByte, &resourceIds); nil != err {
			return err
		}

		tmp := make(map[string]struct{})

		for _, id := range resourceIds {
			tmp[id] = struct{}{}
		}
		for _, id := range inputIds {
			if _, ok := tmp[id]; !ok {
				resourceIds = append(resourceIds, id)
			}
		}
	}

	index, err := rlp.EncodeToBytes(resourceIds)
	if nil != err {
		return err
	}

	return db.Put(GetOrgResourceIdListKey(), index)
}

func RemoveOrgResource (db KeyValueStore, identityId string) error {
	has, err := db.Has(GetOrgResourceIdListKey())
	if IsNoDBNotFoundErr(err) {
		return err
	}

	var identityIds []string
	if !has {
		return nil
	} else {
		idsByte, err := db.Get(GetOrgResourceIdListKey())
		if nil != err {
			return err
		}
		if err := rlp.DecodeBytes(idsByte, &identityIds); nil != err {
			return err
		}


		for i := 0; i < len(identityIds); i++ {
			id := identityIds[i]
			if id == identityId {
				key := GetOrgResourceKey(identityId)
				if err := db.Delete(key); nil != err {
					return err
				}
				identityIds = append(identityIds[:i], identityIds[i+1:]...)
				i--
				break
			}
		}
	}

	index, err := rlp.EncodeToBytes(identityIds)
	if nil != err {
		return err
	}

	return db.Put(GetOrgResourceIdListKey(), index)
}

func QueryOrgResource (db DatabaseReader, identityId string) (*types.RemoteResourceTable, error){
	key := GetOrgResourceKey(identityId)
	vb, err := db.Get(key)
	if nil != err {
		return nil, err
	}

	var resource types.RemoteResourceTable

	if err := rlp.DecodeBytes(vb, &resource); nil != err {
		return nil, err
	}
	return &resource, nil
}

func QueryOrgResources (db DatabaseReader) ([]*types.RemoteResourceTable, error) {
	has, err := db.Has(GetOrgResourceIdListKey())
	if IsNoDBNotFoundErr(err) {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	b, err := db.Get(GetOrgResourceIdListKey())
	if nil != err {
		return nil, err
	}
	var ids []string
	if err := rlp.DecodeBytes(b, &ids); nil != err {
		return nil, err
	}

	arr := make([]*types.RemoteResourceTable, len(ids))
	for i, id := range ids {

		key := GetOrgResourceKey(id)
		vb, err := db.Get(key)
		if nil != err {
			return nil, err
		}

		var resource types.RemoteResourceTable

		if err := rlp.DecodeBytes(vb, &resource); nil != err {
			return nil, err
		}
		arr[i] = &resource
	}

	return arr, nil
}


// 操作 资源单位定义
func StoreNodeResourceSlotUnit(db DatabaseWriter, slot *types.Slot) error {
	val, err := rlp.EncodeToBytes(slot)
	if nil != err {
		return err
	}
	if err := db.Put(GetNodeResourceSlotUnitKey(), val); nil != err {
		return err
	}
	return nil
}

func RemoveNodeResourceSlotUnit (db KeyValueStore) error {
	has, err := db.Has(GetNodeResourceSlotUnitKey())
	if IsNoDBNotFoundErr(err) {
		return err
	}
	if !has {
		return nil
	}
	if err := db.Delete(GetNodeResourceSlotUnitKey()); nil != err {
		return err
	}
	return nil
}

func QueryNodeResourceSlotUnit(db DatabaseReader) (*types.Slot, error) {
	has, err := db.Has(GetNodeResourceSlotUnitKey())
	if IsNoDBNotFoundErr(err) {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	b, err := db.Get(GetNodeResourceSlotUnitKey())
	if nil != err {
		return nil, err
	}
	var slot *types.Slot
	if err := rlp.DecodeBytes(b, &slot); nil != err {
		return nil, err
	}
	return slot, nil
}


// 操作 本地task 正在使用的 算力资源信息
func StoreLocalTaskPowerUsed(db KeyValueStore, taskPowerUsed *types.LocalTaskPowerUsed) error {

	key := GetLocalTaskPowerUsedKey(taskPowerUsed.GetTaskId(), taskPowerUsed.GetPartyId())
	val, err := rlp.EncodeToBytes(taskPowerUsed)
	if nil != err {
		return err
	}

	if err := db.Put(key, val); nil != err {
		return err
	}

	return nil
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

func RemoveLocalTaskPowerUsed(db KeyValueStore, taskId, partyId string) error {
	key := GetLocalTaskPowerUsedKey(taskId, partyId)
	if err := db.Delete(key); nil != err {
		return err
	}
	return nil
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

func QueryLocalTaskPowerUsed (db DatabaseReader, taskId, partyId string) (*types.LocalTaskPowerUsed, error) {
	key := GetLocalTaskPowerUsedKey(taskId, partyId)
	vb, err := db.Get(key)
	if nil != err {
		return nil, err
	}

	var taskPowerUsed types.LocalTaskPowerUsed

	if err := rlp.DecodeBytes(vb, &taskPowerUsed); nil != err {
		return nil, err
	}
	return &taskPowerUsed, nil
}

func QueryLocalTaskPowerUsedsByTaskId (db KeyValueStore, taskId string) ([]*types.LocalTaskPowerUsed, error) {
	it := db.NewIteratorWithPrefixAndStart(GetLocalTaskPowerUsedKeyPrefixByTaskId(taskId), nil)
	defer it.Release()

	arr := make([]*types.LocalTaskPowerUsed, 0)
	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			var taskPowerUsed types.LocalTaskPowerUsed
			if err := rlp.DecodeBytes(it.Value(), &taskPowerUsed); nil != err {
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

func QueryLocalTaskPowerUseds (db KeyValueStore) ([]*types.LocalTaskPowerUsed, error) {
	it := db.NewIteratorWithPrefixAndStart(GetLocalTaskPowerUsedKeyPrefix(), nil)
	defer it.Release()

	arr := make([]*types.LocalTaskPowerUsed, 0)
	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			var taskPowerUsed types.LocalTaskPowerUsed
			if err := rlp.DecodeBytes(it.Value(), &taskPowerUsed); nil != err {
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


// 操作 本地 数据服务 资源信息
func StoreDataResourceTable(db KeyValueStore, dataResourceTable *types.DataResourceTable) error {

	key := GetDataResourceTableKey(dataResourceTable.GetNodeId())
	val, err := rlp.EncodeToBytes(dataResourceTable)
	if nil != err {
		return err
	}

	if err := db.Put(key, val); nil != err {
		return err
	}

	has, err := db.Has(GetDataResourceTableIdListKey())
	if IsNoDBNotFoundErr(err) {
		return err
	}

	var nodeIds []string
	if !has {
		nodeIds = []string{dataResourceTable.GetNodeId()}
	} else {
		idsByte, err := db.Get(GetDataResourceTableIdListKey())
		if nil != err {
			return err
		}
		if err := rlp.DecodeBytes(idsByte, &nodeIds); nil != err {
			return err
		}

		var include bool

		for _, id := range nodeIds {
			if id == dataResourceTable.GetNodeId() {
				include = true
				break
			}
		}
		if !include {
			nodeIds = append(nodeIds, dataResourceTable.GetNodeId())
		}
	}

	index, err := rlp.EncodeToBytes(nodeIds)
	if nil != err {
		return err
	}

	return db.Put(GetDataResourceTableIdListKey(), index)
}

func StoreDataResourceTables(db KeyValueStore, dataResourceTables []*types.DataResourceTable) error {

	has, err := db.Has(GetDataResourceTableIdListKey())
	if IsNoDBNotFoundErr(err) {
		return err
	}

	inputIds := make([]string, len(dataResourceTables))
	for i, dataResourceTable := range dataResourceTables {
		inputIds[i] =  dataResourceTable.GetNodeId()
		key := GetDataResourceTableKey(dataResourceTable.GetNodeId())
		val, err := rlp.EncodeToBytes(dataResourceTable)
		if nil != err {
			return err
		}

		if err := db.Put(key, val); nil != err {
			return err
		}
	}

	var nodeIds []string
	if !has {
		nodeIds = inputIds
	} else {
		idsByte, err := db.Get(GetDataResourceTableIdListKey())
		if nil != err {
			return err
		}
		if err := rlp.DecodeBytes(idsByte, &nodeIds); nil != err {
			return err
		}

		tmp := make(map[string]struct{})

		for _, id := range nodeIds {
			tmp[id] = struct{}{}
		}
		for _, id := range inputIds {
			if _, ok := tmp[id]; !ok {
				nodeIds = append(nodeIds, id)
			}
		}

	}

	index, err := rlp.EncodeToBytes(nodeIds)
	if nil != err {
		return err
	}

	return db.Put(GetDataResourceTableIdListKey(), index)
}

func RemoveDataResourceTable(db KeyValueStore, nodeId string) error {
	has, err := db.Has(GetDataResourceTableIdListKey())
	if IsNoDBNotFoundErr(err) {
		return err
	}

	var nodeIds []string
	if !has {
		return nil
	} else {
		idsByte, err := db.Get(GetDataResourceTableIdListKey())
		if nil != err {
			return err
		}
		if err := rlp.DecodeBytes(idsByte, &nodeIds); nil != err {
			return err
		}


		for i := 0; i < len(nodeIds); i++ {
			id := nodeIds[i]
			if id == nodeId {
				key := GetDataResourceTableKey(nodeId)
				if err := db.Delete(key); nil != err {
					return err
				}
				nodeIds = append(nodeIds[:i], nodeIds[i+1:]...)
				i--
				break
			}
		}
	}

	index, err := rlp.EncodeToBytes(nodeIds)
	if nil != err {
		return err
	}

	return db.Put(GetDataResourceTableIdListKey(), index)
}

func QueryDataResourceTable (db DatabaseReader, nodeId string) (*types.DataResourceTable, error) {
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

func QueryDataResourceTables (db DatabaseReader) ([]*types.DataResourceTable, error) {
	has, err := db.Has(GetDataResourceTableIdListKey())
	if IsNoDBNotFoundErr(err) {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	b, err := db.Get(GetDataResourceTableIdListKey())
	if nil != err {
		return nil, err
	}
	var nodeIds []string
	if err := rlp.DecodeBytes(b, &nodeIds); nil != err {
		return nil, err
	}

	arr := make([]*types.DataResourceTable, len(nodeIds))
	for i, nodeId := range nodeIds {

		key := GetDataResourceTableKey(nodeId)
		vb, err := db.Get(key)
		if nil != err {
			return nil, err
		}

		var dataResourceTable types.DataResourceTable

		if err := rlp.DecodeBytes(vb, &dataResourceTable); nil != err {
			return nil, err
		}
		arr[i] = &dataResourceTable
	}

	return arr, nil
}


// 操作 原始文件Id 所在的 数据服务信息  (originId -> {nodeId/metaDataId/filePath}})
func StoreDataResourceFileUpload(db KeyValueStore, dataResourceFileUpload *types.DataResourceFileUpload) error {

	key := GetDataResourceFileUploadKey(dataResourceFileUpload.GetOriginId())
	val, err := rlp.EncodeToBytes(dataResourceFileUpload)
	if nil != err {
		return err
	}

	if err := db.Put(key, val); nil != err {
		return err
	}

	has, err := db.Has(GetDataResourceFileUploadIdListKey())
	if IsNoDBNotFoundErr(err) {
		return err
	}

	var originIds []string
	if !has {
		originIds = []string{dataResourceFileUpload.GetOriginId()}
	} else {
		idsByte, err := db.Get(GetDataResourceFileUploadIdListKey())
		if nil != err {
			return err
		}
		if err := rlp.DecodeBytes(idsByte, &originIds); nil != err {
			return err
		}

		var include bool

		for _, id := range originIds {
			if id == dataResourceFileUpload.GetOriginId() {
				include = true
				break
			}
		}
		if !include {
			originIds = append(originIds, dataResourceFileUpload.GetOriginId())
		}
	}

	index, err := rlp.EncodeToBytes(originIds)
	if nil != err {
		return err
	}

	return db.Put(GetDataResourceFileUploadIdListKey(), index)
}

func StoreDataResourceFileUploads(db KeyValueStore, dataResourceDataUseds []*types.DataResourceFileUpload) error {

	has, err := db.Has(GetDataResourceFileUploadIdListKey())
	if IsNoDBNotFoundErr(err) {
		return err
	}

	inputIds := make([]string, len(dataResourceDataUseds))
	for i, dataResourceDataUsed := range dataResourceDataUseds {
		inputIds[i] =  dataResourceDataUsed.GetOriginId()
		key := GetDataResourceFileUploadKey(dataResourceDataUsed.GetOriginId())
		val, err := rlp.EncodeToBytes(dataResourceDataUsed)
		if nil != err {
			return err
		}

		if err := db.Put(key, val); nil != err {
			return err
		}
	}

	var originIds []string
	if !has {
		originIds = inputIds
	} else {
		idsByte, err := db.Get(GetDataResourceFileUploadIdListKey())
		if nil != err {
			return err
		}
		if err := rlp.DecodeBytes(idsByte, &originIds); nil != err {
			return err
		}

		tmp := make(map[string]struct{})

		for _, id := range originIds {
			tmp[id] = struct{}{}
		}
		for _, id := range inputIds {
			if _, ok := tmp[id]; !ok {
				originIds = append(originIds, id)
			}
		}

	}

	index, err := rlp.EncodeToBytes(originIds)
	if nil != err {
		return err
	}

	return db.Put(GetDataResourceFileUploadIdListKey(), index)
}

func RemoveDataResourceFileUpload(db KeyValueStore, originId string) error {
	has, err := db.Has(GetDataResourceFileUploadIdListKey())
	if  IsNoDBNotFoundErr(err) {
		return err
	}

	var originIds []string
	if !has {
		return nil
	} else {
		idsByte, err := db.Get(GetDataResourceFileUploadIdListKey())
		if nil != err {
			return err
		}
		if err := rlp.DecodeBytes(idsByte, &originIds); nil != err {
			return err
		}


		for i := 0; i < len(originIds); i++ {
			id := originIds[i]
			if id == originId {
				key := GetDataResourceFileUploadKey(originId)
				if err := db.Delete(key); nil != err {
					return err
				}
				originIds = append(originIds[:i], originIds[i+1:]...)
				i--
				break
			}
		}
	}

	index, err := rlp.EncodeToBytes(originIds)
	if nil != err {
		return err
	}

	return db.Put(GetDataResourceFileUploadIdListKey(), index)
}

func QueryDataResourceFileUpload (db DatabaseReader, originId string) (*types.DataResourceFileUpload, error) {
	key := GetDataResourceFileUploadKey(originId)
	vb, err := db.Get(key)
	if nil != err {
		return nil, err
	}

	var dataResourceDataUsed types.DataResourceFileUpload

	if err := rlp.DecodeBytes(vb, &dataResourceDataUsed); nil != err {
		return nil, err
	}
	return &dataResourceDataUsed, nil
}

func QueryDataResourceFileUploads (db DatabaseReader) ([]*types.DataResourceFileUpload, error) {
	has, err := db.Has(GetDataResourceFileUploadIdListKey())
	if IsNoDBNotFoundErr(err) {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	b, err := db.Get(GetDataResourceFileUploadIdListKey())
	if nil != err {
		return nil, err
	}
	var originIds []string
	if err := rlp.DecodeBytes(b, &originIds); nil != err {
		return nil, err
	}

	arr := make([]*types.DataResourceFileUpload, len(originIds))
	for i, originId := range originIds {

		key := GetDataResourceFileUploadKey(originId)
		vb, err := db.Get(key)
		if nil != err {
			return nil, err
		}

		var dataResourceDataUsed types.DataResourceFileUpload

		if err := rlp.DecodeBytes(vb, &dataResourceDataUsed); nil != err {
			return nil, err
		}
		arr[i] = &dataResourceDataUsed
	}

	return arr, nil
}

func StoreResourceTaskId(db KeyValueStore, jobNodeId, taskId string) error {
	key := GetResourceTaskIdsKey(jobNodeId)
	has, err := db.Has(key)
	if IsNoDBNotFoundErr(err) {
		return err
	}
	var taskIds []string
	if !has {
		taskIds = []string{taskId}
	} else {

		idsByte, err := db.Get(key)
		if nil != err {
			return err
		}
		if err := rlp.DecodeBytes(idsByte, &taskIds); nil != err {
			return err
		}

		var find bool
		for _, id := range taskIds {
			if id == taskId {
				find = true
				break
			}
		}
		if !find {
			taskIds = append(taskIds, taskId)
		}
	}
	index, err := rlp.EncodeToBytes(taskIds)
	if nil != err {
		return err
	}
	return db.Put(key, index)
}

func RemoveResourceTaskId(db KeyValueStore, jobNodeId, taskId string) error {
	key := GetResourceTaskIdsKey(jobNodeId)
	has, err := db.Has(key)
	if IsNoDBNotFoundErr(err) {
		return err
	}
	var taskIds []string
	if !has {
		return nil
	} else {

		idsByte, err := db.Get(key)
		if nil != err {
			return err
		}
		if err := rlp.DecodeBytes(idsByte, &taskIds); nil != err {
			return err
		}
		for i := 0; i < len(taskIds); i++ {
			id := taskIds[i]
			if id == taskId {
				taskIds = append(taskIds[:i], taskIds[i+1:]...)
				i--
				break
			}
		}
	}
	index, err := rlp.EncodeToBytes(taskIds)
	if nil != err {
		return err
	}
	return db.Put(key, index)
}

func QueryResourceTaskIds(db KeyValueStore, jobNodeId string) ([]string, error) {
	key := GetResourceTaskIdsKey(jobNodeId)
	has, err := db.Has(key)
	if IsNoDBNotFoundErr(err) {
		return nil, err
	}
	var taskIds []string
	if !has {
		return nil, ErrNotFound
	} else {
		idsByte, err := db.Get(key)
		if nil != err {
			return nil, err
		}
		if err := rlp.DecodeBytes(idsByte, &taskIds); nil != err {
			return nil, err
		}
	}
	return taskIds, nil
}

func IncreaseResourceTaskPartyIdCount (db KeyValueStore, jobNodeId, taskId string) error {
	count_key := GetResourceTaskPartyIdCountKey(jobNodeId, taskId)
	val, err := db.Get(count_key)

	var count uint32
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case nil == err && len(val) != 0:
		count = bytesutil.BytesToUint32(val)
	}

	count++

	if err := db.Put(count_key, bytesutil.Uint32ToBytes(count)); nil != err {
		return err
	}
	return nil
}

func DecreaseResourceTaskPartyIdCount (db KeyValueStore, jobNodeId, taskId string) error {
	count_key := GetResourceTaskPartyIdCountKey(jobNodeId, taskId)
	val, err := db.Get(count_key)

	var count uint32
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err):
		return nil
	case nil == err && len(val) != 0:
		count = bytesutil.BytesToUint32(val)
	}

	if count == 0 {
		if err := db.Delete(count_key); nil != err {
			return err
		}
	} else {

		count++

		if err := db.Put(count_key, bytesutil.Uint32ToBytes(count)); nil != err {
			return err
		}
	}

	return nil
}

func QueryResourceTaskPartyIdCount (db DatabaseReader, jobNodeId, taskId string) (uint32, error) {
	count_key := GetResourceTaskPartyIdCountKey(jobNodeId, taskId)
	val, err := db.Get(count_key)

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

func StoreLocalResourceIdByPowerId(db DatabaseWriter, powerId, resourceId string) error {
	key := GetResourcePowerIdMapingKey(powerId)
	index, err := rlp.EncodeToBytes(resourceId)
	if nil != err {
		return err
	}
	return db.Put(key, index)
}

func RemoveLocalResourceIdByPowerId(db DatabaseDeleter, powerId string) error {
	key := GetResourcePowerIdMapingKey(powerId)
	return db.Delete(key)
}

func QueryLocalResourceIdByPowerId(db DatabaseReader, powerId string) (string, error) {
	key := GetResourcePowerIdMapingKey(powerId)
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
	var resourceId string
	if err := rlp.DecodeBytes(idsByte, &resourceId); nil != err {
		return "", err
	}
	return resourceId, nil
}

func StoreDataResourceDiskUsed(db DatabaseWriter, dataResourceDiskUsed *types.DataResourceDiskUsed) error {
	key := GetDataResourceDiskUsedKey(dataResourceDiskUsed.GetMetadataId())
	val, err := rlp.EncodeToBytes(dataResourceDiskUsed)
	if nil != err {
		return err
	}
	return db.Put(key, val)
}

func RemoveDataResourceDiskUsed(db DatabaseDeleter, metaDataId string) error {
	key := GetDataResourceDiskUsedKey(metaDataId)
	return db.Delete(key)
}

func QueryDataResourceDiskUsed(db DatabaseReader, metaDataId string) (*types.DataResourceDiskUsed, error) {
	key := GetDataResourceDiskUsedKey(metaDataId)
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

func StoreLocalTaskExecuteStatus(db DatabaseWriter, taskId string) error {
	key := GetLocalTaskExecuteStatus(taskId)
	val, err := rlp.EncodeToBytes("yes")
	if nil != err {
		return err
	}
	return db.Put(key, val)
}

func RemoveLocalTaskExecuteStatus(db DatabaseDeleter, taskId string) error {
	key := GetLocalTaskExecuteStatus(taskId)
	return db.Delete(key)
}

func HasLocalTaskExecute(db DatabaseReader, taskId string) (bool, error) {
	key := GetLocalTaskExecuteStatus(taskId)
	has, err := db.Has(key)
	if IsNoDBNotFoundErr(err) {
		return false, err
	}
	if !has {
		return false, nil
	}
	return true, nil
}


func StoreUserMetadataAauthUsed (db KeyValueStore, userType apicommonpb.UserType, user, metadataAuthId string) error {
	count_key := GetUserMetadataAuthUsedCountKey(userType, user)
	val, err := db.Get(count_key)

	var count uint32
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case nil == err && len(val) != 0:
		count = bytesutil.BytesToUint32(val)
	}

	count++

	if err := db.Put(count_key, bytesutil.Uint32ToBytes(count)); nil != err {
		return err
	}

	item_key := GetUserMetadataAuthUsedKey(userType, user, count)

	item, err := rlp.EncodeToBytes(metadataAuthId)
	if nil != err {
		return err
	}
	return db.Put(item_key, item)
}

func QueryUserMetadataAuthUsedCount (db DatabaseReader, userType apicommonpb.UserType, user string) (uint32, error) {
	count_key := GetUserMetadataAuthUsedCountKey(userType, user)
	val, err := db.Get(count_key)

	var count uint32
	switch {
	case IsNoDBNotFoundErr(err):
		return 0, err
	case nil == err && len(val) != 0:
		count = bytesutil.BytesToUint32(val)
	}
	return count, nil
}
func QueryUserMetadataAuthUsedByIndex (db DatabaseReader, userType apicommonpb.UserType, user string, index uint32) (string, error) {
	item_key := GetUserMetadataAuthUsedKey(userType, user, index)

	var metadataAuthId string
	item_val, err := db.Get(item_key)
	if nil != err {
		return "", err
	}

	if err = rlp.DecodeBytes(item_val, &metadataAuthId); nil != err {
		return "", err
	}
	if "" == metadataAuthId {
		return "", ErrNotFound
	}

	return metadataAuthId, nil
}


func RemoveUserMetadataAuthUsedCount (db KeyValueStore, userType apicommonpb.UserType, user string) error {
	count_key := GetUserMetadataAuthUsedCountKey(userType, user)

	has, err := db.Has(count_key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}

	return db.Delete(count_key)
}

func RemoveUserMetadataAuthUsedByIndex (db KeyValueStore, userType apicommonpb.UserType, user string, index uint32) error {
	item_key := GetUserMetadataAuthUsedKey(userType, user, index)

	has, err := db.Has(item_key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}

	return db.Delete(item_key)
}

func StoreUserMetadataAuthIdByMetadataId (db DatabaseWriter, userType apicommonpb.UserType, user, metadataId, metadataAuthId string) error {
	key := GetUserMetadataAuthByMetadataIdKey(userType, user, metadataId)
	val, err := rlp.EncodeToBytes(metadataAuthId)
	if nil != err {
		return err
	}
	return db.Put(key, val)
}

func QueryUserMetadataAuthIdByMetadataId (db DatabaseReader, userType apicommonpb.UserType, user, metadataId string) (string, error) {
	key := GetUserMetadataAuthByMetadataIdKey(userType, user, metadataId)

	var metadataAuthId string
	val, err := db.Get(key)
	if nil != err {
		return "", err
	}

	if err = rlp.DecodeBytes(val, &metadataAuthId); nil != err {
		return "", err
	}
	if "" == metadataAuthId {
		return "", ErrNotFound
	}
	return metadataAuthId, nil
}

func HasUserMetadataAuthIdByMetadataId (db DatabaseReader, userType apicommonpb.UserType, user, metadataId string) (bool, error) {
	key := GetUserMetadataAuthByMetadataIdKey(userType, user, metadataId)

	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return false, err
	case IsDBNotFoundErr(err), nil == err && !has:
		return false, nil
	}
	return true, nil
}


func RemoveUserMetadataAuthIdByMetadataId (db KeyValueStore, userType apicommonpb.UserType, user, metadataId string) error {
	key := GetUserMetadataAuthByMetadataIdKey(userType, user, metadataId)

	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}



func StoreMetadataUsedTaskId (db KeyValueStore, metadataId, taskId string) error {
	count_key := GetMetadataUsedTaskIdCountKey(metadataId)
	val, err := db.Get(count_key)

	var count uint32
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case nil == err && len(val) != 0:
		count = bytesutil.BytesToUint32(val)
	}

	count++

	if err := db.Put(count_key, bytesutil.Uint32ToBytes(count)); nil != err {
		return err
	}

	item_key := GetMetadataUsedTaskIdKey(metadataId, count)

	item, err := rlp.EncodeToBytes(taskId)
	if nil != err {
		return err
	}
	return db.Put(item_key, item)
}

func QueryMetadataUsedTaskIdCount (db DatabaseReader, metadataId string) (uint32, error) {
	count_key := GetMetadataUsedTaskIdCountKey(metadataId)
	val, err := db.Get(count_key)

	var count uint32
	switch {
	case IsNoDBNotFoundErr(err):
		return 0, err
	case nil == err && len(val) != 0:
		count = bytesutil.BytesToUint32(val)
	}
	return count, nil
}
func QueryMetadataUsedTaskIdByIndex (db DatabaseReader, metadataId string, index uint32) (string, error) {
	item_key := GetMetadataUsedTaskIdKey(metadataId, index)

	var taskId string
	item_val, err := db.Get(item_key)
	if nil != err {
		return "", err
	}

	if err = rlp.DecodeBytes(item_val, &taskId); nil != err {
		return "", err
	}
	if "" == taskId {
		return "", ErrNotFound
	}

	return taskId, nil
}


func RemoveMetadataUsedTaskIdCount (db KeyValueStore, metadataId string) error {
	count_key := GetMetadataUsedTaskIdCountKey(metadataId)

	has, err := db.Has(count_key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}

	return db.Delete(count_key)
}

func RemoveMetadataUsedTaskIdByIndex (db KeyValueStore, metadataId string, index uint32) error {
	item_key := GetMetadataUsedTaskIdKey(metadataId, index)

	has, err := db.Has(item_key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}

	return db.Delete(item_key)
}



func StoreTaskUpResultFile (db DatabaseWriter, turf *types.TaskUpResultFile)  error {
	key := GetTaskResultFileMetadataIdKey(turf.GetTaskId())
	val, err := rlp.EncodeToBytes(turf)
	if nil != err {
		return err
	}
	return db.Put(key, val)
}

func QueryTaskUpResultFile (db DatabaseReader, taskId string)  (*types.TaskUpResultFile, error) {
	key := GetTaskResultFileMetadataIdKey(taskId)
	vb, err := db.Get(key)
	if nil != err {
		return nil, err
	}
	var taskUpResultFile types.TaskUpResultFile
	if err = rlp.DecodeBytes(vb, &taskUpResultFile); nil != err {
		return nil, err
	}
	return &taskUpResultFile, nil
}

func QueryTaskUpResultFileList (db DatabaseIteratee) ([]*types.TaskUpResultFile, error) {

	it := db.NewIteratorWithPrefixAndStart(GetTaskResultFileMetadataIdKeyPrefix(), nil)
	defer it.Release()

	arr := make([]*types.TaskUpResultFile, 0)
	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			var taskUpResultFile types.TaskUpResultFile
			if err := rlp.DecodeBytes(it.Value(), &taskUpResultFile); nil != err {
				log.Errorf("Failed to call QueryAllTaskUpResultFile, decode db val failed, err: {%s}", err)
				continue
			}
			arr = append(arr, &taskUpResultFile)
		}
	}

	if len(arr) == 0 {
		return nil, ErrNotFound
	}

	return arr, nil
}

func RemoveTaskUpResultFile (db KeyValueStore, taskId string) error {
	key := GetTaskResultFileMetadataIdKey(taskId)

	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}



func StoreTaskResuorceUsage (db DatabaseWriter, usage *types.TaskResuorceUsage) error {
	key := GetTaskResuorceUsageKey(usage.GetTaskId(), usage.GetPartyId())
	val, err := rlp.EncodeToBytes(usage)
	if nil != err {
		return err
	}
	return db.Put(key, val)
}

func QueryTaskResuorceUsage (db DatabaseReader, taskId, partyId string) (*types.TaskResuorceUsage, error) {
	key := GetTaskResuorceUsageKey(taskId, partyId)

	vb, err := db.Get(key)
	if nil != err {
		return nil, err
	}
	var taskResuorceUsage types.TaskResuorceUsage
	if err = rlp.DecodeBytes(vb, &taskResuorceUsage); nil != err {
		return nil, err
	}
	return &taskResuorceUsage, nil
}

func RemoveTaskResuorceUsage (db KeyValueStore, taskId, partyId string) error {
	key := GetTaskResuorceUsageKey(taskId, partyId)

	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func RemoveTaskResuorceUsageByTaskId (db KeyValueStore, taskId string) error {

	it := db.NewIteratorWithPrefixAndStart(GetTaskResuorceUsageKeyPrefixByTaskId(taskId), nil)
	defer it.Release()

	for it.Next() {
		if key := it.Key(); len(key) != 0 {
			db.Delete(key)
		}
	}

	return nil
}


func StoreMessageCache(db KeyValueStore, value interface{}){
	byteArray, _ := rlp.EncodeToBytes(value)
	key:=""
	switch value.(type) {
	case types.PowerMsgArr:
		key="PowerMsgArr"
	case types.MetadataMsgArr:
		key="MetadataMsgArr"
	case types.MetadataAuthorityMsgArr:
		key="MetadataAuthorityMsgArr"
	case types.TaskMsgArr:
		key="TaskMsgArr"
	}
	if err := db.Put([]byte(key), byteArray); err != nil {
		log.Warning("StoreMessageCache fail,key is:",key)
	}
}

func QueryPowerMsgArr(db KeyValueStore) (types.PowerMsgArr, error) {
	key := []byte("PowerMsgArr")
	has, err := db.Has(key)
	if IsNoDBNotFoundErr(err) {
		return nil, err
	}
	var result types.PowerMsgArr
	if !has {
		return nil, ErrNotFound
	} else {
		byteArray, er := db.Get(key)
		if nil != er {
			return nil, er
		}
		if err := rlp.DecodeBytes(byteArray, &result); nil != err {
			return nil, err
		}
	}
	return result, nil
}

func QueryMetadataMsgArr(db KeyValueStore) (types.MetadataMsgArr, error) {
	key := []byte("MetadataMsgArr")
	has, err := db.Has(key)
	if IsNoDBNotFoundErr(err) {
		return nil, err
	}
	var result types.MetadataMsgArr
	if !has {
		return nil, ErrNotFound
	} else {
		byteArray, er := db.Get(key)
		if nil != er {
			return nil, er
		}
		if err := rlp.DecodeBytes(byteArray, &result); nil != err {
			return nil, err
		}
	}
	return result, nil
}

func QueryMetadataAuthorityMsgArr(db KeyValueStore) (types.MetadataAuthorityMsgArr, error)  {
	key := []byte("MetadataAuthorityMsgArr")
	has, err := db.Has(key)
	if IsNoDBNotFoundErr(err) {
		return nil, err
	}
	var result types.MetadataAuthorityMsgArr
	if !has {
		return nil, ErrNotFound
	} else {
		byteArray, er := db.Get(key)
		if nil != er {
			return nil, er
		}
		if err := rlp.DecodeBytes(byteArray, &result); nil != err {
			return nil, err
		}
	}
	return result, nil
}

func QueryTaskMsgArr(db KeyValueStore) (types.TaskMsgArr, error) {
	key := []byte("TaskMsgArr")
	has, err := db.Has(key)
	if IsNoDBNotFoundErr(err) {
		return nil, err
	}
	var result types.TaskMsgArr
	if !has {
		return nil, ErrNotFound
	} else {
		byteArray, er := db.Get(key)
		if nil != er {
			return nil, er
		}
		if err := rlp.DecodeBytes(byteArray, &result); nil != err {
			return nil, err
		}
	}
	return result, nil
}