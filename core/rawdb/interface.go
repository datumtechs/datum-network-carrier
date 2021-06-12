// Copyright (C) 2021 The RosettaNet Authors.

package rawdb

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

type KeyValueStore interface {
	DatabaseReader
	DatabaseWriter
	DatabaseDeleter
}