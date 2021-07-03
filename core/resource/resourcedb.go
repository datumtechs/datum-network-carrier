package resource

import (
	"errors"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	ErrNotFound = errors.New("rawdb: not found")
)

func storeNodeResource(db rawdb.DatabaseWriter, resource *types.ResourceTable) error {
	val, err := rlp.EncodeToBytes(resource)
	if nil != err {
		return err
	}
	if err := db.Put(GetNodeResourceKey(resource.GetNodeId()), val); nil != err {
		return err
	}
	return nil
}
func queryNodeResource(db rawdb.DatabaseReader, nodeId string) (*types.ResourceTable, error) {
	has, err := db.Has(GetNodeResourceKey(nodeId))
	if nil != err {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	b, err := db.Get(GetNodeResourceKey(nodeId))
	if nil != err {
		return nil, err
	}
	var resource types.ResourceTable
	if err := rlp.DecodeBytes(b, &resource); nil != err {
		return nil, err
	}
	return &resource, nil
}


func storeNodeResources(db rawdb.DatabaseWriter, resources []*types.ResourceTable) error {
	val, err := rlp.EncodeToBytes(resources)
	if nil != err {
		return err
	}
	if err := db.Put(GetNodeResourceIdListKey(), val); nil != err {
		return err
	}
	return nil
}

func queryNodeResources (db rawdb.DatabaseReader) ([]*types.ResourceTable, error) {
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
	var arr []*types.ResourceTable
	if err := rlp.DecodeBytes(b, &arr); nil != err {
		return nil, err
	}
	return arr, nil
}


func storeNodeResourceSlotUnit(db rawdb.DatabaseWriter, slot *types.Slot) error {
	val, err := rlp.EncodeToBytes(slot)
	if nil != err {
		return err
	}
	if err := db.Put(GetNodeResourceSlotUnitKey(), val); nil != err {
		return err
	}
	return nil
}

func queryNodeResourceSlotUnit(db rawdb.DatabaseReader) (*types.Slot, error) {
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