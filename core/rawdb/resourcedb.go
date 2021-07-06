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

func StoreNodeResources(db DatabaseWriter, resources []*types.LocalResourceTable) error {
	val, err := rlp.EncodeToBytes(resources)
	if nil != err {
		return err
	}
	if err := db.Put(GetNodeResourceIdListKey(), val); nil != err {
		return err
	}
	return nil
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
	var arr []*types.LocalResourceTable
	if err := rlp.DecodeBytes(b, &arr); nil != err {
		return nil, err
	}
	return arr, nil
}

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