package db

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
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
