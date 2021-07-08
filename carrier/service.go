package carrier

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/consensus/chaincons"
	"github.com/RosettaFlow/Carrier-Go/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/core"
	"github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/message"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/core/scheduler"
	"github.com/RosettaFlow/Carrier-Go/core/task"
	"github.com/RosettaFlow/Carrier-Go/db"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	"github.com/RosettaFlow/Carrier-Go/handler"
	"github.com/RosettaFlow/Carrier-Go/types"
	"sync"
)

type Service struct {
	isRunning      bool
	processingLock sync.RWMutex
	config         *Config
	carrierDB      core.CarrierDB
	ctx            context.Context
	cancel         context.CancelFunc
	mempool        *message.Mempool
	Engines        map[types.ConsensusEngineType]handler.Engine

	// DB interfaces
	dataDb     db.Database
	APIBackend *CarrierAPIBackend

	resourceManager *resource.Manager
	messageManager  *message.MessageHandler
	taskManager     *task.Manager
	scheduler       core.Scheduler
	runError        error

	// GRPC Client
	jobNodes  map[string]*grpclient.JobNodeClient
	dataNodes map[string]*grpclient.DataNodeClient
}

// NewService creates a new CarrierServer object (including the
// initialisation of the common Carrier object)
func NewService(ctx context.Context, config *Config) (*Service, error) {
	ctx, cancel := context.WithCancel(ctx)
	_ = cancel // govet fix for lost cancel. Cancel is handled in service.Stop()

	taskCh := make(chan types.TaskMsgs, 0)
	pool := message.NewMempool(nil) // todo need  set mempool cfg

	eventEngine := evengine.NewEventEngine(config.CarrierDB)

	// TODO 这些 Ch 的大小目前都是写死的 ...
	localTaskCh, schedTaskCh, remoteTaskCh, sendTaskCh, recvSchedTaskCh :=
		make(chan types.TaskMsgs, 27),
		make(chan *types.ConsensusTaskWrap, 100),
		make(chan *types.ScheduleTaskWrap, 100),
		make(chan types.TaskMsgs, 10),
		make(chan *types.ConsensusScheduleTask, 10)
	resourceMng := resource.NewResourceManager(config.CarrierDB)

	s := &Service{
		ctx:             ctx,
		cancel:          cancel,
		config:          config,
		carrierDB:       config.CarrierDB,
		mempool:         pool,
		resourceManager: resourceMng,
		messageManager:  message.NewHandler(pool, nil, config.CarrierDB, taskCh),                                 // todo need set dataChain
		taskManager:     task.NewTaskManager(nil, eventEngine, resourceMng, taskCh, sendTaskCh, recvSchedTaskCh), // todo need set dataChain
		scheduler: scheduler.NewSchedulerStarveFIFO(localTaskCh, schedTaskCh, remoteTaskCh,
			config.CarrierDB, recvSchedTaskCh, resourceMng, eventEngine),
		jobNodes:  make(map[string]*grpclient.JobNodeClient),
		dataNodes: make(map[string]*grpclient.DataNodeClient),
	}
	// todo: some logic could be added...
	s.APIBackend = &CarrierAPIBackend{carrier: s}
	s.Engines = make(map[types.ConsensusEngineType]handler.Engine, 0)
	s.Engines[types.TwopcTyp] = twopc.New(&twopc.Config{}, s.carrierDB, s.config.p2p, schedTaskCh, remoteTaskCh) // todo the 2pc config will be setup
	s.Engines[types.ChainconsTyp] = chaincons.New()
	// todo: set datachain....
	return s, nil
}

func (s *Service) Start() error {
	for typ, engine := range s.Engines {
		if err := engine.Start(); nil != err {
			log.WithError(err).Errorf("Cound not start the consensus engine: %s, err: %v", typ.String(), err)
		}
	}
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
