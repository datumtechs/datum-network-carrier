package carrier

import (
	"context"
	"fmt"
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
	"github.com/RosettaFlow/Carrier-Go/p2p"
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

	// internal resource node set (Fighter node grpc client set)
	resourceClientSet *grpclient.InternalResourceClientSet
}

// NewService creates a new CarrierServer object (including the
// initialisation of the common Carrier object)
func NewService(ctx context.Context, config *Config) (*Service, error) {
	ctx, cancel := context.WithCancel(ctx)
	_ = cancel // govet fix for lost cancel. Cancel is handled in service.Stop()

	nodeId := config.P2P.NodeId()
	pool := message.NewMempool(&message.MempoolConfig{NodeId: nodeId})
	eventEngine := evengine.NewEventEngine(config.CarrierDB)

	// TODO 这些 Ch 的大小目前都是写死的 ...
	localTaskMsgCh, needConsensusTaskCh, replayScheduleTaskCh, doneScheduleTaskCh :=
		make(chan types.TaskMsgs, 27),
		make(chan *types.ConsensusTaskWrap, 100),
		make(chan *types.ReplayScheduleTaskWrap, 100),
		make(chan *types.DoneScheduleTaskChWrap, 10)

	resourceClientSet := grpclient.NewInternalResourceNodeSet()

	resourceMng := resource.NewResourceManager(config.CarrierDB)

	taskManager := task.NewTaskManager(
		config.CarrierDB,
		eventEngine,
		resourceMng,
		resourceClientSet,
		localTaskMsgCh,
		doneScheduleTaskCh,
	)

	s := &Service{
		ctx:             ctx,
		cancel:          cancel,
		config:          config,
		carrierDB:       config.CarrierDB,
		mempool:         pool,
		resourceManager: resourceMng,
		messageManager:  message.NewHandler(pool, config.CarrierDB, taskManager),
		taskManager:     taskManager,
		scheduler: scheduler.NewSchedulerStarveFIFO(
			eventEngine,
			resourceMng,
			config.CarrierDB,
			localTaskMsgCh,
			needConsensusTaskCh,
			replayScheduleTaskCh,
			doneScheduleTaskCh,
		),
		resourceClientSet: resourceClientSet,
	}
	// read config from p2p config.
	NodeId, _ := p2p.HexID(nodeId)

	s.APIBackend = &CarrierAPIBackend{carrier: s}
	s.Engines = make(map[types.ConsensusEngineType]handler.Engine, 0)
	s.Engines[types.TwopcTyp] = twopc.New(
		&twopc.Config{
			Option: &twopc.OptionConfig{
				NodePriKey: s.config.P2P.PirKey(),
				NodeID:     NodeId,
			},
			PeerMsgQueueSize: 1024,
		},
		s.carrierDB,
		resourceMng,
		s.config.P2P,
		needConsensusTaskCh,
		replayScheduleTaskCh,
		doneScheduleTaskCh,
	)
	s.Engines[types.ChainconsTyp] = chaincons.New()

	// load stored jobNode and dataNode
	jobNodeList, err := s.carrierDB.GetRegisterNodeList(types.PREFIX_TYPE_JOBNODE)
	if err == nil {
		for _, node := range jobNodeList {
			client, err := grpclient.NewJobNodeClient(ctx, fmt.Sprintf("%s:%s", node.InternalIp, node.InternalPort), node.Id)
			if err == nil {
				s.resourceClientSet.StoreJobNodeClient(node.Id, client)
			}
		}
	}
	dataNodeList, err := s.carrierDB.GetRegisterNodeList(types.PREFIX_TYPE_DATANODE)
	if err == nil {
		for _, node := range dataNodeList {
			client, err := grpclient.NewDataNodeClient(ctx, fmt.Sprintf("%s:%s", node.InternalIp, node.InternalPort), node.Id)
			if err == nil {
				s.resourceClientSet.StoreDataNodeClient(node.Id, client)
			}
		}
	}
	return s, nil
}

func (s *Service) Start() error {
	for typ, engine := range s.Engines {
		if err := engine.Start(); nil != err {
			log.WithError(err).Errorf("Cound not start the consensus engine: %s, err: %v", typ.String(), err)
		}
	}
	if nil != s.resourceManager {
		if err := s.resourceManager.Start(); nil != err {
			log.WithError(err).Errorf("Failed to start the resourceManager, err: %v", err)
		}
	}
	if nil != s.messageManager {
		if err := s.messageManager.Start(); nil != err {
			log.WithError(err).Errorf("Failed to start the messageManager, err: %v", err)
		}
	}
	if nil != s.taskManager {
		if err := s.taskManager.Start(); nil != err {
			log.WithError(err).Errorf("Failed to start the taskManager, err: %v", err)
		}
	}
	if nil != s.scheduler {
		if err := s.scheduler.Start(); nil != err {
			log.WithError(err).Errorf("Failed to start the scheduler, err: %v", err)
		}
	}

	return nil
}

func (s *Service) Stop() error {
	if s.cancel != nil {
		defer s.cancel()
	}
	// todo: could add some logic for here（some logic for stop.）
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
