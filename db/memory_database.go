package db

import (
	"errors"
	"github.com/RosettaFlow/Carrier-Go/common"
	"sync"
)

type MemoryDatabase struct {
	db 		map[string][]byte
	lock 	sync.RWMutex
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		db:  make(map[string][]byte),
	}
}

func NewMemoryDatabaseWithCap(size int) *MemoryDatabase {
	return &MemoryDatabase{
		db:  make(map[string][]byte, size),
	}
}

func (db *MemoryDatabase) Put(key []byte, value []byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	db.db[string(key)] = common.CopyBytes(value)
	return nil
}

func (db *MemoryDatabase) Has(key []byte) (bool, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	_, ok := db.db[string(key)]
	return ok, nil
}

func (db *MemoryDatabase) Get(key []byte) ([]byte, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if entry, ok := db.db[string(key)]; ok {
		return common.CopyBytes(entry), nil
	}
	return nil, errors.New("not found by special key")
}

func (db *MemoryDatabase) Keys() [][]byte {
	db.lock.RLock()
	defer db.lock.RUnlock()

	keys := [][]byte{}
	for key := range db.db {
		keys = append(keys, []byte(key))
	}
	return keys
}

func (db *MemoryDatabase) Delete(key []byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	delete(db.db, string(key))
	return nil
}

func (db *MemoryDatabase) Close() {
}

func (db *MemoryDatabase) NewBatch() Batch {
	return &memBatch{
		db: db,
	}
}

func (db *MemoryDatabase) Len() int {
	return len(db.db)
}

type kv struct {
	k []byte
	v []byte
	del bool
}

type memBatch struct {
	db 		*MemoryDatabase
	writes 	[]kv
	size 	int
}

func (b *memBatch) Put(key, value []byte) error {
	b.writes = append(b.writes, kv{common.CopyBytes(key), common.CopyBytes(value), false})
	b.size += len(value)
	return nil
}

func (b *memBatch) Delete(key []byte) error {
	b.writes = append(b.writes, kv{common.CopyBytes(key), nil, true})
	b.size += 1
	return nil
}

func (b *memBatch) Write() error {
	b.db.lock.Lock()
	defer b.db.lock.Unlock()

	for _, kv := range b.writes {
		if kv.del {
			delete(b.db.db, string(kv.k))
			continue
		}
		b.db.db[string(kv.k)] = kv.v
	}
	return nil
}

func (b *memBatch) ValueSize() int {
	return b.size
}

func (b *memBatch) Reset() {
	b.writes = b.writes[:0]
	b.size = 0
}