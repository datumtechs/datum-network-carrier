// Copyright (C) 2021 The RosettaNet Authors.

package rawdb

import "github.com/datumtechs/datum-network-carrier/db"

type DatabaseReader interface {
	Has(key []byte) (bool, error)
	Get(key []byte) ([]byte, error)
}

type DatabaseWriter interface {
	Put(key, value []byte) error
}

type DatabaseDeleter interface {
	Delete(key []byte) error
}

type DatabaseIteratee interface {
	NewIteratorWithPrefixAndStart(prefix []byte, start []byte) db.Iterator
}

type KeyValueStore interface {
	DatabaseReader
	DatabaseWriter
	DatabaseDeleter
	DatabaseIteratee
}