package core

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/core/types"
	"github.com/RosettaFlow/Carrier-Go/db"
	"github.com/RosettaFlow/Carrier-Go/event"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/params"
	lru "github.com/hashicorp/golang-lru"
	"sync"
	"sync/atomic"
)

const (
	blockCacheLimit = 256
)

type DataChain struct {
	rosettaConfig *params.CarrierConfig // network configuration

	db        db.Database // Low level persistent database to store final content.
	chainFeed event.Feed
	mu        sync.RWMutex // global mutex for locking chain operations.
	chainmu   sync.RWMutex // datachain insertion lock
	procmu    sync.RWMutex // data processor lock

	currentBlock atomic.Value // Current head of the data chain

	blockCache  *lru.Cache
	bodyCache   *lru.Cache
	bodyPbCache *lru.Cache
	quit        chan struct{} // datachain quit channel
	running     int32         // running must be called atomically

	procInterrupt int32          // interrupt signaler for block processing
	wg            sync.WaitGroup // chain processing wait group for shutting down
}

// NewDataChain returns a fully initialised data chain using information available in the database.
func NewDataChain(db db.Database, rosettaConfig *params.CarrierConfig) (*DataChain, error) {
	blockCache, _ := lru.New(blockCacheLimit)
	bodyCache, _ := lru.New(blockCacheLimit)
	bodyPbCache, _ := lru.New(blockCacheLimit)
	dc := &DataChain{
		rosettaConfig: rosettaConfig,
		db:            db,
		quit:          make(chan struct{}),
		blockCache:    blockCache,
		bodyCache:     bodyCache,
		bodyPbCache:   bodyPbCache,
	}
	return dc, nil
}

func (dc *DataChain) getProcInterrupt() bool {
	return atomic.LoadInt32(&dc.procInterrupt) == 1
}

// loadLastState loads the last known data chain state from the database.
func (dc *DataChain) loadLastState() error {
	head := rawdb.ReadHeadBlockHash(dc.db)
	if head == (common.Hash{}) {
		log.Warn("Empty database, resetting chain")
		return nil
	}
	return nil
}

func (dc *DataChain) CurrentBlock() *types.Block {
	// convert type
	return dc.currentBlock.Load().(*types.Block)
}

func (dc *DataChain) SetProcessor() {
	dc.procmu.Lock()
	defer dc.procmu.Unlock()
	// do setting...
}

func (dc *DataChain) SetValidator() {
	dc.procmu.Lock()
	defer dc.procmu.Unlock()
	// do setting...
}

func (dc *DataChain) insert(block *types.Block) {

}

// GetBody retrieves a block body (metadata/resource/identity/task) from database by hash, caching it if found.
func (dc *DataChain) GetBody(hash common.Hash) *libTypes.BodyData {
	if cached, ok := dc.bodyCache.Get(hash); ok {
		body := cached.(*libTypes.BodyData)
		return body
	}
	number := rawdb.ReadHeaderNumber(dc.db, hash)
	if number == nil {
		return nil
	}
	body := rawdb.ReadBody(dc.db, hash, *number)
	if body == nil {
		return nil
	}
	// Cache the found body for next time and return
	dc.bodyCache.Add(hash, body)
	return body
}

// GetBodyPb retrieves a block body in PB encoding from the database by hash, caching it if found
func (dc *DataChain) GetBodyPb(hash common.Hash) []byte {
	if cached, ok := dc.bodyPbCache.Get(hash); ok {
		return cached.([]byte)
	}
	number := rawdb.ReadHeaderNumber(dc.db, hash)
	if number == nil {
		return nil
	}
	body := rawdb.ReadBodyPB(dc.db, hash, *number)
	if len(body) == 0 {
		return nil
	}
	dc.bodyPbCache.Add(hash, body)
	return body
}

// HasBlock checks if a block is fully present in the database or not.
func (dc *DataChain) HasBlock(hash common.Hash, number uint64) bool {
	if dc.blockCache.Contains(hash) {
		return true
	}
	return rawdb.HasBody(dc.db, hash, number)
}
