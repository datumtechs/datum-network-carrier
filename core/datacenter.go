package core

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/core/types"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	"github.com/RosettaFlow/Carrier-Go/params"
	lru "github.com/hashicorp/golang-lru"
	"sync"
	"sync/atomic"
)

// DataCenter is mainly responsible for communicating with the data center service
type DataCenter struct {
	config 			*params.DataCenterConfig
	client 			*grpclient.GrpcClient
	mu     			sync.RWMutex 	// global mutex for locking data center operations.
	procmu 			sync.RWMutex 	// data processor lock

	blockCache  	*lru.Cache
	bodyCache   	*lru.Cache
	bodyPbCache 	*lru.Cache

	processor 		Processor // block processor interface

	quit        	chan struct{} 	// datachain quit channel
	running     	int32         	// running must be called atomically
	procInterrupt 	int32          	// interrupt signaler for block processing
	wg            	sync.WaitGroup 	// chain processing wait group for shutting down
}

// NewDataCenter returns a fully initialised data center using information available in the database.
func NewDataCenter(config *params.DataCenterConfig) (*DataCenter, error) {
	blockCache, _ := lru.New(blockCacheLimit)
	bodyCache, _ := lru.New(blockCacheLimit)
	bodyPbCache, _ := lru.New(blockCacheLimit)
	// todo: When to call Close??
	client, err := grpclient.Dial(fmt.Sprintf("%v:%v", config.GrpcUrl, config.Port))
	if err != nil {
		log.WithError(err).Error("dial grpc server failed")
		return nil, err
	}
	dc := &DataCenter{
		config: 		config,
		client: 		client,
		quit:          	make(chan struct{}),
		blockCache:    	blockCache,
		bodyCache:     	bodyCache,
		bodyPbCache:   	bodyPbCache,
	}
	return dc, nil
}

func (dc *DataCenter) getProcInterrupt() bool {
	return atomic.LoadInt32(&dc.procInterrupt) == 1
}

func (dc *DataCenter) SetProcessor(processor Processor) {
	dc.procmu.Lock()
	defer dc.procmu.Unlock()
	// do setting...
	dc.processor = processor
}

// interface:
// 1、插入元数据、资源、身份、任务（一个接口接入，传多个参数[列表]）
// 2、修改元数据、资源、身份、任务（同样支持批量多少数据）
// 3、删除元数据、资源、身份、任务（支持批量多数据）
// 4、查询元数据、资源、身份、任务列表
// 5、根据节点ID查询元数据、资源、身份、任务
// 6、根据节点ID查询参与的任务列表

func (dc *DataCenter) Insert(metas []*types.Metadata, resources []*types.Resource, identities []*types.Identity, tasks []*types.Task) error {
	return nil
}

func (dc *DataCenter) Update(metas []*types.Metadata, resources []*types.Resource, identities []*types.Identity, tasks []*types.Task) error {
	return nil
}

func (dc *DataCenter) Delete(metas []*types.Metadata, resources []*types.Resource, identities []*types.Identity, tasks []*types.Task) error {
	return nil
}

func (dc *DataCenter) MetadataList(nodeId string) (types.MetadataArray, error) {
	return nil, nil
}

func (dc *DataCenter) ResourceList(nodeId string) (types.ResourceArray, error) {
	return nil, nil
}

func (dc *DataCenter) IdentityList(nodeId string) (types.IdentityArray, error) {
	return nil, nil
}

func (dc *DataCenter) TaskDataList(nodeId string) (types.TaskDataArray, error) {
	return nil, nil
}




