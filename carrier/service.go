package carrier

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/auth"
	"github.com/RosettaFlow/Carrier-Go/common/flags"
	"github.com/RosettaFlow/Carrier-Go/consensus/chaincons"
	"github.com/RosettaFlow/Carrier-Go/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/core"
	"github.com/RosettaFlow/Carrier-Go/core/election"
	"github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/message"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/core/schedule"
	"github.com/RosettaFlow/Carrier-Go/core/task"
	"github.com/RosettaFlow/Carrier-Go/db"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	"github.com/RosettaFlow/Carrier-Go/handler"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/service/discovery"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/urfave/cli/v2"
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
	TaskManager     handler.TaskManager
	authManager     *auth.AuthorityManager
	scheduler       schedule.Scheduler
	consulManager   *discovery.ConnectConsul
	runError        error

	// internal resource node set (Fighter node grpc client set)
	resourceClientSet *grpclient.InternalResourceClientSet
}

// NewService creates a new CarrierServer object (including the
// initialisation of the common Carrier object)
func NewService(ctx context.Context, cliCtx *cli.Context, config *Config, mockIdentityIdsFile, consensusStateFile string) (*Service, error) {
	ctx, cancel := context.WithCancel(ctx)
	_ = cancel // govet fix for lost cancel. Cancel is handled in service.Stop()

	nodeIdStr := config.P2P.NodeId()
	// read config from p2p config.
	nodeId, _ := p2p.HexID(nodeIdStr)

	pool := message.NewMempool(&message.MempoolConfig{NodeId: nodeIdStr})
	eventEngine := evengine.NewEventEngine(config.CarrierDB)

	// TODO The size of these ch is currently written dead ...
	localTaskMsgCh, needReplayScheduleTaskCh, needExecuteTaskCh, taskConsResultCh :=
		make(chan types.TaskDataArray, 100),
		make(chan *types.NeedReplayScheduleTask, 600),
		make(chan *types.NeedExecuteTask, 600),
		make(chan *types.TaskConsResult, 600)

	resourceClientSet := grpclient.NewInternalResourceNodeSet()
	resourceMng := resource.NewResourceManager(config.CarrierDB, mockIdentityIdsFile)
	authManager := auth.NewAuthorityManager(config.CarrierDB)
	scheduler := schedule.NewSchedulerStarveFIFO(election.NewVrfElector(config.P2P.PirKey(), resourceClientSet, resourceMng), eventEngine, resourceMng, authManager)
	twopcEngine := twopc.New(
		&twopc.Config{
			Option: &twopc.OptionConfig{
				NodePriKey: config.P2P.PirKey(),
				NodeID:     nodeId,
			},
			PeerMsgQueueSize:    1024,
			ConsensusStateFile:  consensusStateFile,
			DefaultConsensusWal: config.DefaultConsensusWal,
			DatabaseCache:       config.DatabaseCache,
			DatabaseHandles:     config.DatabaseHandles,
		},
		resourceMng,
		config.P2P,
		needReplayScheduleTaskCh,
		needExecuteTaskCh,
		taskConsResultCh,
	)
	taskManager := task.NewTaskManager(
		config.P2P,
		scheduler,
		twopcEngine,
		eventEngine,
		resourceMng,
		authManager,
		resourceClientSet,
		localTaskMsgCh,
		needReplayScheduleTaskCh,
		needExecuteTaskCh,
		taskConsResultCh,
	)

	s := &Service{
		ctx:               ctx,
		cancel:            cancel,
		config:            config,
		carrierDB:         config.CarrierDB,
		mempool:           pool,
		resourceManager:   resourceMng,
		messageManager:    message.NewHandler(pool, config.CarrierDB, taskManager, authManager, resourceClientSet),
		TaskManager:       taskManager,
		authManager:       authManager,
		scheduler:         scheduler,
		resourceClientSet: resourceClientSet,
		consulManager: discovery.New(&discovery.ConsulService{
			ServiceIP:   cliCtx.String(flags.RPCHost.Name),
			ServicePort: cliCtx.String(flags.RPCPort.Name),
			Tags:        cliCtx.StringSlice(flags.DiscoveryServerTags.Name),
			Name:        cliCtx.String(flags.DiscoveryServiceName.Name),
			Id:          cliCtx.String(flags.DiscoveryServiceId.Name),
			Interval:    cliCtx.Int(flags.DiscoveryServiceHealthCheckInterval.Name),
			Deregister:  cliCtx.Int(flags.DiscoveryServiceHealthCheckDeregister.Name),
		},
			cliCtx.String(flags.DiscoveryServerIP.Name),
			cliCtx.String(flags.DiscoveryServerPort.Name),
		),
	}

	//s.APIBackend = &CarrierAPIBackend{carrier: s}
	s.APIBackend = NewCarrierAPIBackend(s)
	s.Engines = make(map[types.ConsensusEngineType]handler.Engine, 0)
	s.Engines[types.TwopcTyp] = twopcEngine
	s.Engines[types.ChainconsTyp] = chaincons.New()

	// load stored jobNode and dataNode
	jobNodeList, err := s.carrierDB.QueryRegisterNodeList(pb.PrefixTypeJobNode)
	if err == nil {
		for _, node := range jobNodeList {
			client, err := grpclient.NewJobNodeClient(ctx, fmt.Sprintf("%s:%s", node.GetInternalIp(), node.GetInternalPort()), node.GetId())
			if err == nil {
				s.resourceClientSet.StoreJobNodeClient(node.GetId(), client)
			}
		}
	}
	dataNodeList, err := s.carrierDB.QueryRegisterNodeList(pb.PrefixTypeDataNode)
	if err == nil {
		for _, node := range dataNodeList {
			client, err := grpclient.NewDataNodeClient(ctx, fmt.Sprintf("%s:%s", node.GetInternalIp(), node.GetInternalPort()), node.GetId())
			if err == nil {
				s.resourceClientSet.StoreDataNodeClient(node.GetId(), client)
			}
		}
	}
	return s, nil
}

func (s *Service) Start() error {

	if nil != s.authManager {
		if err := s.authManager.Start(); nil != err {
			log.WithError(err).Errorf("Failed to start the authManager")
		}
	}

	for typ, engine := range s.Engines {
		if err := engine.Start(); nil != err {
			log.WithError(err).Errorf("Cound not start the consensus engine: %s", typ.String())
		}
	}
	if nil != s.resourceManager {
		if err := s.resourceManager.Start(); nil != err {
			log.WithError(err).Errorf("Failed to start the resourceManager")
		}
	}
	if nil != s.messageManager {
		if err := s.messageManager.Start(); nil != err {
			log.WithError(err).Errorf("Failed to start the messageManager")
		}
	}
	if nil != s.TaskManager {
		if err := s.TaskManager.Start(); nil != err {
			log.WithError(err).Errorf("Failed to start the TaskManager")
		}
	}
	if nil != s.scheduler {
		if err := s.scheduler.Start(); nil != err {
			log.WithError(err).Errorf("Failed to start the schedule")
		}
	}
	if nil != s.consulManager {
		err := s.consulManager.RegisterDiscoveryService()
		if err != nil {
			log.WithError(err).Fatal("Failed to register discovery service to discovery center")
		}
	}
	return nil
}

func (s *Service) Stop() error {
	if s.cancel != nil {
		defer s.cancel()
	}
	s.carrierDB.Stop()

	for typ, engine := range s.Engines {
		if err := engine.Stop(); nil != err {
			log.WithError(err).Errorf("Cound not close the consensus engine: %s, err: %v", typ.String(), err)
		}
	}
	if nil != s.resourceManager {
		if err := s.resourceManager.Stop(); nil != err {
			log.WithError(err).Errorf("Failed to stop the resourceManager, err: %v", err)
		}
	}
	if nil != s.messageManager {
		if err := s.messageManager.Stop(); nil != err {
			log.WithError(err).Errorf("Failed to stop the messageManager, err: %v", err)
		}
	}
	if nil != s.TaskManager {
		if err := s.TaskManager.Stop(); nil != err {
			log.WithError(err).Errorf("Failed to stop the TaskManager, err: %v", err)
		}
	}
	if nil != s.scheduler {
		if err := s.scheduler.Stop(); nil != err {
			log.WithError(err).Errorf("Failed to stop the schedule, err: %v", err)
		}
	}

	if nil != s.authManager {
		if err := s.authManager.Stop(); nil != err {
			log.WithError(err).Errorf("Failed to stop the authManager, err: %v", err)
		}
	}
	if nil != s.consulManager {
		if err := s.consulManager.Stop(); nil != err {
			log.WithError(err).Errorf("ServiceDeregister fail,serviceId is %v", err)
		}
	}
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
