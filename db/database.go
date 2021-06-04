package db

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"sync"
)

// LDBDatabase LevelDB implementation of db interface.
type LDBDatabase struct {
	fileName string
	db       *leveldb.DB

	// todo:
	quitLock sync.Mutex
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
		BlockCacheCapacity:     cache / 2 * opt.MiB,
		WriteBuffer:            cache / 4 * opt.MiB, // Two of these are used internally
		Filter:                 filter.NewBloomFilter(10),
	})
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(filename, nil)
	}
	//
	if err != nil {
		return nil, err
	}
	return &LDBDatabase{
		fileName: filename,
		db:       db,
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

// NewIteratorWithPrefix returns a iterator to iterate over subset of database content with a particular prefix.
func (db *LDBDatabase) NewIteratorWithPrefix(prefix []byte) iterator.Iterator {
	return db.db.NewIterator(util.BytesPrefix(prefix), nil)
}

func (db *LDBDatabase) Close() {
	db.quitLock.Lock()
	defer db.quitLock.Unlock()
	if db.quitChan != nil {
		errc := make(chan error)
		db.quitChan <- errc
		if err := <- errc; err != nil {
			log.WithError(err).Error("Metrics collection failed")
		}
		db.quitChan = nil
	}
	err := db.db.Close()
	if err == nil {
		log.Info("Database closed")
	} else {
		log.WithError(err).Error("Failed to close database")
	}
}

func (db *LDBDatabase) LDB() *leveldb.DB {
	return db.db
}

func (db *LDBDatabase) NewBatch() Batch {
	return &ldbBatch{
		db: db.db,
		b: new(leveldb.Batch),
	}
}

type ldbBatch struct {
	db		*leveldb.DB
	b 		*leveldb.Batch
	size 	int
}

func (b *ldbBatch) Put(key, value []byte) error {
	b.b.Put(key, value)
	b.size += len(value)
	return nil
}

func (b *ldbBatch) Delete(key []byte) error {
	b.b.Delete(key)
	b.size += 1
	return nil
}

func (b *ldbBatch) Write() error {
	return b.db.Write(b.b, nil)
}

func (b *ldbBatch) ValueSize() int {
	return b.size
}

func (b *ldbBatch) Reset()  {
	b.b.Reset()
	b.size = 0
}
