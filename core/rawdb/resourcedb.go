package rawdb

import (
	"errors"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	ErrNotFound = errors.New("rawdb: not found")
)

//func storeNodeResource(db rawdb.DatabaseWriter, resource *types.LocalResourceTable) error {
//	val, err := rlp.EncodeToBytes(resource)
//	if nil != err {
//		return err
//	}
//	if err := db.Put(GetNodeResourceKey(resource.GetNodeId()), val); nil != err {
//		return err
//	}
//	return nil
//}
//func queryNodeResource(db rawdb.DatabaseReader, nodeId string) (*types.LocalResourceTable, error) {
//	has, err := db.Has(GetNodeResourceKey(nodeId))
//	if nil != err {
//		return nil, err
//	}
//	if !has {
//		return nil, ErrNotFound
//	}
//	b, err := db.Get(GetNodeResourceKey(nodeId))
//	if nil != err {
//		return nil, err
//	}
//	var resource types.LocalResourceTable
//	if err := rlp.DecodeBytes(b, &resource); nil != err {
//		return nil, err
//	}
//	return &resource, nil
//}

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
	if nil != err {
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
	if nil != err {
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

func QueryNodeResources (db DatabaseReader) ([]*types.LocalResourceTable, error) {
	has, err := db.Has(GetNodeResourceIdListKey())
	if nil != err {
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


		var resource types.LocalResourceTable

		key := GetNodeResourceKey(id)
		if err := rlp.DecodeBytes(key, &resource); nil != err {
			return nil, err
		}
		arr[i] = &resource
	}

	return arr, nil
}

// TODO 写到这里了 ...

func StoreOrgResources(db DatabaseWriter, resources []*types.RemoteResourceTable) error {
	val, err := rlp.EncodeToBytes(resources)
	if nil != err {
		return err
	}
	if err := db.Put(GetOrgResourceIdListKey(), val); nil != err {
		return err
	}
	return nil
}

func QueryOrgResources (db DatabaseReader) ([]*types.RemoteResourceTable, error) {
	has, err := db.Has(GetOrgResourceIdListKey())
	if nil != err {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	b, err := db.Get(GetOrgResourceIdListKey())
	if nil != err {
		return nil, err
	}
	var arr []*types.RemoteResourceTable
	if err := rlp.DecodeBytes(b, &arr); nil != err {
		return nil, err
	}
	return arr, nil
}

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

func QueryNodeResourceSlotUnit(db DatabaseReader) (*types.Slot, error) {
	has, err := db.Has(GetNodeResourceSlotUnitKey())
	if nil != err {
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