package carrier

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/core"
	"github.com/RosettaFlow/Carrier-Go/core/task"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/RosettaFlow/Carrier-Go/types"
	"sync"
)

type Service struct {
	isRunning               bool
	processingLock          sync.RWMutex
	config 					*params.CarrierConfig
	proxy 					*core.DataCenter
	ctx                     context.Context
	cancel                  context.CancelFunc
	mempool 				*core.Mempool
	TaskMng 				*task.TaskManager
	runError                error
}

func NewService(ctx context.Context, config *params.CarrierConfig, dataCenterConfig *params.DataCenterConfig) (*Service, error) {
	ctx, cancel := context.WithCancel(ctx)
	_ = cancel // govet fix for lost cancel. Cancel is handled in service.Stop()
	proxy, err := core.NewDataCenter(dataCenterConfig)
	if err != nil {
		cancel()
		return nil, err
	}
	s := &Service{
		ctx:            ctx,
		cancel:         cancel,
		config:			config,
		proxy: 		 	proxy,
		mempool:        core.NewMempool(nil), // todo need  set mempool cfg
	}
	// todo: some logic could be added...
	return s, nil
}

func (s *Service) Start() {

}

func (s *Service) Stop() error {
	if s.cancel != nil {
		defer s.cancel()
	}
	// todo: could add some logic for here
	return nil
}

// Status is service health checks. Return nil or error.
func (s *Service) Status() error {
	// Service don't start
	if !s.isRunning {
		return nil
	}
	// get error from run function
	if s.runError != nil {
		return s.runError
	}
	return nil
}

func (s *Service) SendMsg (msg types.Msg) error {
	return s.mempool.Add(msg)
}
func (s *Service) SetSeedNode (seed *types.SeedNodeInfo) (types.NodeConnStatus,error) {
	return types.NONCONNECTED, nil
}
func (s *Service) DeleteSeedNode(id string) error {
	return nil
}
func (s *Service) GetSeedNode (id string) (*types.SeedNodeInfo, error) {
	return nil, nil
}
func (s *Service) GetSeedNodeList () ([]*types.SeedNodeInfo, error) {
	return nil, nil
}
func (s *Service) SetRegisterNode (typ types.RegisteredNodeType, node *types.RegisteredNodeInfo) (types.NodeConnStatus,error) {
	return types.NONCONNECTED, nil
}
func (s *Service) DeleteRegisterNode (typ types.RegisteredNodeType, id string) error {
	return nil
}
func (s *Service) GetRegisterNode (typ types.RegisteredNodeType, id string) (*types.RegisteredNodeInfo, error) {
	return nil, nil
}
func (s *Service) GetRegisterNodeList (typ types.RegisteredNodeType) ([]*types.RegisteredNodeInfo, error) {
	return nil, nil
}