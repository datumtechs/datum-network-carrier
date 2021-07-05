package carrier

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/consensus"
	"github.com/RosettaFlow/Carrier-Go/core"
	"github.com/RosettaFlow/Carrier-Go/core/message"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/core/task"
	"github.com/RosettaFlow/Carrier-Go/db"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/RosettaFlow/Carrier-Go/types"
	"sync"
)

type Service struct {
	isRunning      bool
	processingLock sync.RWMutex
	config         *Config
	dataCenter     *core.DataCenter
	datachain      *core.DataChain
	ctx            context.Context
	cancel         context.CancelFunc
	mempool        *message.Mempool

	// Consensuses
	cons 		  map[string]consensus.Engine

	// DB interfaces
	dataDb     db.Database
	APIBackend *CarrierAPIBackend

	resourceManager *resource.Manager
	messageManager  *message.MessageHandler
	taskManager     *task.Manager
	runError        error
}

// NewService creates a new CarrierServer object (including the
// initialisation of the common Carrier object)
func NewService(ctx context.Context, config *Config, dataCenterConfig *params.DataCenterConfig) (*Service, error) {
	ctx, cancel := context.WithCancel(ctx)
	_ = cancel // govet fix for lost cancel. Cancel is handled in service.Stop()

	proxy, err := core.NewDataCenter(ctx, dataCenterConfig)
	if err != nil {
		cancel()
		return nil, err
	}
	taskCh := make(chan types.TaskMsgs, 0)
	pool := message.NewMempool(nil) // todo need  set mempool cfg
	s := &Service{
		ctx:             ctx,
		cancel:          cancel,
		config:          config,
		dataCenter:      proxy,
		mempool:         pool,
		resourceManager: resource.NewResourceManager(),
		messageManager: message.NewHandler(pool, nil, proxy, taskCh, nil), // todo need set dataChain
		taskManager:     task.NewTaskManager(nil, taskCh, nil), // todo need set dataChain
	}
	// todo: some logic could be added...

	// todo: set datachain....
	return s, nil
}

func (s *Service) Start() error {
	return nil
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
