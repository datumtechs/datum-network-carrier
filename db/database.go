package db

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"sync"
)

// LDBDatabase LevelDB implementation of db interface.
type LDBDatabase struct {
	fileName	string
	db 			*leveldb.DB

	// todo:
	quitLock	sync.Mutex
	quitChan chan chan error
}

// NewLDBDatabase returns a levelDB wrapped object
func NewLDBDatabase(filename string, cache int, handles int) (*LDBDatabase, error) {
	if cache < 16 {
		cache = 16
	}
	if handles < 16 {
		handles = 16
	}
	log.WithField("cache", cache).WithField("handles", handles).Info("Allocated cache and file handles")

	// Open the db and recover any potential corruptions
	db, err := leveldb.OpenFile(filename, &opt.Options{
		OpenFilesCacheCapacity: handles,
		BlockCacheCapacity: 	cache / 2 * opt.MiB,
		WriteBuffer: 			cache /4 * opt.MiB,	// Two of these are used internally
		Filter: 				filter.NewBloomFilter(10),
	})
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(filename, nil)
	}
	//
	if err != nil {
		return nil, err
	}
	return &LDBDatabase{
		fileName: 	filename,
		db: 		db,
	}, nil
}

// Path returns the path to the database directory.
func (db *LDBDatabase) Path() string {
	return db.fileName
}

// Put puts the given key/value to the queue.
func (db *LDBDatabase) Put(key []byte, value []byte) error {
	return db.db.Put(key, value, nil)
}

func (db *LDBDatabase) Has(key []byte) (bool, error) {
	return db.db.Has(key, nil)
}

// Get returns the given key if it's present.
func (db *LDBDatabase) Get(key []byte) ([]byte, error) {
	dat, err := db.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

// Delete deletes the key from the queue and database.
func (db *LDBDatabase) Delete(key []byte) error {
	return db.db.Delete(key, nil)
}

func (db *LDBDatabase) NewIterator() iterator.Iterator {
	return db.db.NewIterator(nil, nil)
}


