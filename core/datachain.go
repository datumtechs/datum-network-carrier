package core

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/db"
	"github.com/RosettaFlow/Carrier-Go/event"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/RosettaFlow/Carrier-Go/types"
	lru "github.com/hashicorp/golang-lru"
	"sync"
	"sync/atomic"
)

const (
	blockCacheLimit = 256
)

type DataChain struct {
	chainConfig *params.DataChainConfig // network configuration

	db        db.Database // Low level persistent database to store final content.
	chainFeed event.Feed
	mu        sync.RWMutex // global mutex for locking chain operations.
	chainmu   sync.RWMutex // datachain insertion lock
	procmu    sync.RWMutex // data processor lock

	currentBlock atomic.Value // Current head of the data chain
	processor    Processor

	blockCache  *lru.Cache
	bodyCache   *lru.Cache
	bodyPbCache *lru.Cache
	quit        chan struct{} // datachain quit channel
	running     int32         // running must be called atomically

	procInterrupt int32          // interrupt signaler for block processing
	wg            sync.WaitGroup // chain processing wait group for shutting down
}

// NewDataChain returns a fully initialised data chain using information available in the database.
func NewDataChain(db db.Database, chainConfig *params.DataChainConfig) (*DataChain, error) {
	blockCache, _ := lru.New(blockCacheLimit)
	bodyCache, _ := lru.New(blockCacheLimit)
	bodyPbCache, _ := lru.New(blockCacheLimit)
	dc := &DataChain{
		chainConfig: chainConfig,
		db:          db,
		quit:        make(chan struct{}),
		blockCache:  blockCache,
		bodyCache:   bodyCache,
		bodyPbCache: bodyPbCache,
	}
	return dc, nil
}

func (dc *DataChain) getProcInterrupt() bool {
	return atomic.LoadInt32(&dc.procInterrupt) == 1
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

// Processor returns the current processor.
func (dc *DataChain) Processor() Processor {
	dc.procmu.RLock()
	defer dc.procmu.RUnlock()
	return dc.processor
}

// HasBlock checks if a block is fully present in the database or not.
func (dc *DataChain) HasBlock(hash common.Hash, number uint64) bool {
	if dc.blockCache.Contains(hash) {
		return true
	}
	return rawdb.HasBody(dc.db, hash, number)
}

func (dc *DataChain) CurrentBlock() *types.Block {
	// convert type
	return dc.currentBlock.Load().(*types.Block)
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

// GetBlock retrieves a block from the database by hash and number,
// caching it if found.
func (dc *DataChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	// Short circuit if the block's already in the cache, retrieve otherwise
	if block, ok := dc.blockCache.Get(hash); ok {
		return block.(*types.Block)
	}
	block := rawdb.ReadBlock(dc.db, hash, number)
	if block == nil {
		return nil
	}
	// Cache the found block for next time and return
	dc.blockCache.Add(block.Hash(), block)
	return block
}

// GetBlockByHash retrieves a block from the database by hash, caching it if found.
func (dc *DataChain) GetBlockByHash(hash common.Hash) *types.Block {
	number := rawdb.ReadHeaderNumber(dc.db, hash)
	if number == nil {
		return nil
	}
	return dc.GetBlock(hash, *number)
}

// GetBlockByNumber retrieves a block from the database by number, caching it
// (associated with its hash) if found.
func (dc *DataChain) GetBlockByNumber(number uint64) *types.Block {
	hash := rawdb.ReadDataHash(dc.db, number)
	if hash == (common.Hash{}) {
		return nil
	}
	return dc.GetBlock(hash, number)
}

// InsertChain saves the data of block to the database.
func (dc *DataChain) InsertData(blocks types.Blocks) (int, error) {
	// metadata/resource/task...
	// todo: updateData()/revokeData()
	return 0, nil
}

// TODO 本地存储当前调度服务的 name
func (dc *DataChain) StoreYarnName(name string) error {
	rawdb.WriteYarnName(dc.db, name)
	return nil
}

// TODO 本地存储当前调度服务自身的  identity
func (dc *DataChain) StoreIdentity(identity string) error {
	rawdb.WriteIdentityStr(dc.db, identity)
	return nil
}

func (dc *DataChain) GetYarnName() (string, error) {
	return rawdb.ReadYarnName(dc.db), nil
}

func (dc *DataChain) GetIdentity() (string, error) {
	return rawdb.ReadIdentityStr(dc.db), nil
}

func (dc *DataChain) GetMetadataByHash(hash common.Hash) (*types.Metadata, error) {
	return nil, nil
}

func (dc *DataChain) GetMetadataByDataId(dataId string) (*types.Metadata, error) {
	return nil, nil
}

func (dc *DataChain) GetMetadataListByNodeId(nodeId string) (types.MetadataArray, error) {
	return nil, nil
}

func (dc *DataChain) GetResourceByHash(hash common.Hash) (*types.Resource, error) {
	return nil, nil
}

func (dc *DataChain) GetResourceByDataId(dataId string) (*types.Resource, error) {
	return nil, nil
}

func (dc *DataChain) GetResourceListByNodeId(nodeId string) (types.ResourceArray, error) {
	return nil, nil
}

func (dc *DataChain) GetIdentityByHash(hash common.Hash) (*types.Identity, error) {
	return nil, nil
}

func (dc *DataChain) GetIdentityByDataId(nodeId string) (*types.Identity, error) {
	return nil, nil
}

func (dc *DataChain) GetIdentityListByNodeId(nodeId string) (types.IdentityArray, error) {
	return nil, nil
}

func (dc *DataChain) GetTaskDataByHash(hash common.Hash) (*types.Task, error) {
	return nil, nil
}

func (dc *DataChain) GetTaskDataByTaskId(taskId string) (*types.Task, error) {
	return nil, nil
}

func (dc *DataChain) GetTaskDataListByNodeId(nodeId string) (types.TaskDataArray, error) {
	return nil, nil
}

func (dc *DataChain) SetSeedNode(seed *types.SeedNodeInfo) (types.NodeConnStatus, error) {
	return types.NONCONNECTED, nil
}

func (dc *DataChain) DeleteSeedNode(id string) error {
	rawdb.DeleteSeedNode(dc.db, id)
	return nil
}

func (dc *DataChain) GetSeedNode(id string) (*types.SeedNodeInfo, error) {
	return rawdb.ReadSeedNode(dc.db, id), nil
}

func (dc *DataChain) GetSeedNodeList() ([]*types.SeedNodeInfo, error) {
	return rawdb.ReadAllSeedNodes(dc.db), nil
}

func (dc *DataChain) SetRegisterNode(typ types.RegisteredNodeType, node *types.RegisteredNodeInfo) (types.NodeConnStatus, error) {
	rawdb.WriteRegisterNodes(dc.db, typ, node)
	// todo: need to establish conn to registered node. heartbeat detection
	return types.NONCONNECTED, nil
}

func (dc *DataChain) DeleteRegisterNode(typ types.RegisteredNodeType, id string) error {
	rawdb.DeleteRegisterNode(dc.db, typ, id)
	return nil
}

func (dc *DataChain) GetRegisterNode(typ types.RegisteredNodeType, id string) (*types.RegisteredNodeInfo, error) {
	return rawdb.ReadRegisterNode(dc.db, typ, id), nil
}

func (dc *DataChain) GetRegisterNodeList(typ types.RegisteredNodeType) ([]*types.RegisteredNodeInfo, error) {
	return rawdb.ReadAllRegisterNodes(dc.db, typ), nil
}

// TODO 存储当前组织正在参与运行的任务详情 (正在运行的, 管理台需要看 我的任务列表, 存储使用的 pb和数据中心一样, 最后需要两边都查回来本地合并列表展示)
func (dc *DataChain) StoreRunningTask(task *types.Task) error {

	return nil
}

// TODO 存储当前某个 计算服务正在执行的任务Id
func (dc *DataChain) StoreJobNodeRunningTaskId(jobNodeId, taskId string) {

}

//// TODO 存储当前某个 数据服务正在执行的任务Id  (先不要这个)
//func (dc *DataChain) StoreDataNodeRunningTaskId(dataNodeId, taskId string) {
//
//}

// TODO 存储当前组织 正在运行的任务总数 (递增)
func (dc *DataChain) IncreaseRunningTaskCountOnOrg() uint32 {

	return 0
}

// TODO 存储当前计算服务 正在运行的任务总数 (递增)
func (dc *DataChain) IncreaseRunningTaskCountOnJobNode(jobNodeId string) uint32 {
	return rawdb.IncreaseRunningTaskCountForOrg(dc.db)
}

// TODO 查询当前组织 正在运行的任务总数
func (dc *DataChain) GetRunningTaskCountOnOrg () uint32 {
	return rawdb.ReadRunningTaskCountForOrg(dc.db)
}

// TODO 查询当前计算服务 正在运行的任务总数
func (dc *DataChain) GetRunningTaskCountOnJobNode (jobNodeId string) uint32 {
	return 0
}

// TODO 查询当前组织计算服务正在执行的任务Id列表 (正在运行的还未结束的任务)
func (dc *DataChain) GetJobNodeRunningTaskIdList (jobNodeId string) []string {
	return nil
}

//// TODO 查询当前组织数据服务正在执行的任务Id列表 (正在运行的还未结束的任务) (先不要这个)
//func (dc *DataChain) GetDataNodeRunningTaskIdList (dataNodeId string) []string {
//
//}

// Config retrieves the datachain's chain configuration.
func (dc *DataChain) Config() *params.DataChainConfig { return dc.chainConfig }

// Stop stops the DataChain service. If any imports are currently in progress
// it will abort them using the procInterrupt.
func (dc *DataChain) Stop() {
	if !atomic.CompareAndSwapInt32(&dc.running, 0, 1) {
		return
	}
	// todo: add logic...
	log.Info("Datachain manager stopped")
}