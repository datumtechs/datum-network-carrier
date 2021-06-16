package core

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/db"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
)

// DataCenter is mainly responsible for communicating with the data center service
type DataCenter struct {
	ctx       context.Context
	config    *params.DataCenterConfig
	client    *grpclient.GrpcClient
	mu        sync.RWMutex // global mutex for locking data center operations.
	serviceMu sync.RWMutex // data processor lock

	db db.Database // Low level persistent database to store final content.

	processor     Processor      // block processor interface
	running       int32          // running must be called atomically
	procInterrupt int32          // interrupt signaler for block processing
	wg            sync.WaitGroup // chain processing wait group for shutting down
}

// NewDataCenter returns a fully initialised data center using information available in the database.
func NewDataCenter(config *params.DataCenterConfig) (*DataCenter, error) {
	// todo: When to call Close??
	client, err := grpclient.Dial(fmt.Sprintf("%v:%v", config.GrpcUrl, config.Port))
	if err != nil {
		log.WithError(err).Error("dial grpc server failed")
		return nil, err
	}
	dc := &DataCenter{
		config: config,
		client: client,
	}
	return dc, nil
}

func (dc *DataCenter) getProcInterrupt() bool {
	return atomic.LoadInt32(&dc.procInterrupt) == 1
}

func (dc *DataCenter) SetProcessor(processor Processor) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	// do setting...
	dc.processor = processor
}

// ************************************* public api (datachain) *******************************************

// InsertChain saves the data of block to the database.
func (dc *DataCenter) InsertData(blocks types.Blocks) (int, error) {
	if len(blocks) == 0 {
		return 0, nil
	}
	// check.
	for i := 0; i < len(blocks); i++ {
		if blocks[i].NumberU64() != blocks[i-1].NumberU64()+1 || blocks[i].ParentHash() != blocks[i-1].Hash() {
			log.WithFields(logrus.Fields{
				"number": blocks[i].NumberU64(),
				"hash":   blocks[i].Hash(),
			}).Error("Non contiguous block insert")
			return 0, fmt.Errorf("non contiguous insert: item %d is #%d", i-1, blocks[i-1].NumberU64())
		}
	}
	// pre-checks passed, start the full block imports
	dc.wg.Add(1)
	defer dc.wg.Done()

	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()

	headers := make([]*types.Header, len(blocks))
	seals := make([]bool, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
		seals[i] = true
	}
	for _, block := range blocks {
		if atomic.LoadInt32(&dc.procInterrupt) == 1 {
			log.Debug("Premature abort during blocks processing")
			break
		}
		err := dc.processor.Process(block, dc.config)
		if err != nil {
			// for err, how to deal with????
		}
	}
	return 0, nil
}

func (dc *DataCenter) GetMetadataByHash(hash common.Hash) (*types.Metadata, error) {
	return nil, nil
}

func (dc *DataCenter) GetMetadataByDataId(dataId string) (*types.Metadata, error) {
	return nil, nil
}

func (dc *DataCenter) GetMetadataListByNodeId(nodeId string) (types.MetadataArray, error) {
	return nil, nil
}

func (dc *DataCenter) GetResourceByHash(hash common.Hash) (*types.Resource, error) {
	return nil, nil
}

func (dc *DataCenter) GetResourceByDataId(dataId string) (*types.Resource, error) {
	return nil, nil
}

func (dc *DataCenter) GetResourceListByNodeId(nodeId string) (types.ResourceArray, error) {
	return nil, nil
}

func (dc *DataCenter) GetIdentityByHash(hash common.Hash) (*types.Identity, error) {
	return nil, nil
}

func (dc *DataCenter) GetIdentityByDataId(nodeId string) (*types.Identity, error) {
	return nil, nil
}

func (dc *DataCenter) GetIdentityListByNodeId(nodeId string) (types.IdentityArray, error) {
	return nil, nil
}

func (dc *DataCenter) GetTaskDataByHash(hash common.Hash) (*types.Task, error) {
	return nil, nil
}

func (dc *DataCenter) GetTaskDataByTaskId(taskId string) (*types.Task, error) {
	return nil, nil
}

func (dc *DataCenter) GetTaskDataListByNodeId(nodeId string) (types.TaskDataArray, error) {
	return nil, nil
}

func (dc *DataCenter) SetSeedNode(seed *types.SeedNodeInfo) (types.NodeConnStatus, error) {
	return types.NONCONNECTED, nil
}

func (dc *DataCenter) DeleteSeedNode(id string) error {
	rawdb.DeleteSeedNode(dc.db, id)
	return nil
}

func (dc *DataCenter) GetSeedNode(id string) (*types.SeedNodeInfo, error) {
	return rawdb.ReadSeedNode(dc.db, id), nil
}

func (dc *DataCenter) GetSeedNodeList() ([]*types.SeedNodeInfo, error) {
	return rawdb.ReadAllSeedNodes(dc.db), nil
}

func (dc *DataCenter) SetRegisterNode(typ types.RegisteredNodeType, node *types.RegisteredNodeInfo) (types.NodeConnStatus, error) {
	rawdb.WriteRegisterNodes(dc.db, typ, node)
	// todo: need to establish conn to registered node. heartbeat detection
	return types.NONCONNECTED, nil
}

func (dc *DataCenter) DeleteRegisterNode(typ types.RegisteredNodeType, id string) error {
	rawdb.DeleteRegisterNode(dc.db, typ, id)
	return nil
}

func (dc *DataCenter) GetRegisterNode(typ types.RegisteredNodeType, id string) (*types.RegisteredNodeInfo, error) {
	return rawdb.ReadRegisterNode(dc.db, typ, id), nil
}

func (dc *DataCenter) GetRegisterNodeList(typ types.RegisteredNodeType) ([]*types.RegisteredNodeInfo, error) {
	return rawdb.ReadAllRegisterNodes(dc.db, typ), nil
}

// ****************************************************************************************************************

func (dc *DataCenter) Stop() {
	if !atomic.CompareAndSwapInt32(&dc.running, 0, 1) {
		return
	}
	atomic.StoreInt32(&dc.procInterrupt, 1)
	dc.wg.Wait()
	dc.client.Close()
	log.Info("Datacenter manager stopped")
}
