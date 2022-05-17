package carrier

import (
	"context"
	"fmt"
	"github.com/Metisnetwork/Metis-Carrier/ach/auth"
	"github.com/Metisnetwork/Metis-Carrier/ach/metispay"
	"github.com/Metisnetwork/Metis-Carrier/ach/metispay/kms"
	"github.com/Metisnetwork/Metis-Carrier/common/flags"
	"github.com/Metisnetwork/Metis-Carrier/consensus/chaincons"
	"github.com/Metisnetwork/Metis-Carrier/consensus/twopc"
	"github.com/Metisnetwork/Metis-Carrier/core"
	"github.com/Metisnetwork/Metis-Carrier/core/election"
	"github.com/Metisnetwork/Metis-Carrier/core/evengine"
	"github.com/Metisnetwork/Metis-Carrier/core/message"
	"github.com/Metisnetwork/Metis-Carrier/core/resource"
	"github.com/Metisnetwork/Metis-Carrier/core/schedule"
	"github.com/Metisnetwork/Metis-Carrier/core/task"
	"github.com/Metisnetwork/Metis-Carrier/db"
	"github.com/Metisnetwork/Metis-Carrier/grpclient"
	"github.com/Metisnetwork/Metis-Carrier/handler"
	pb "github.com/Metisnetwork/Metis-Carrier/lib/api"
	"github.com/Metisnetwork/Metis-Carrier/p2p"
	"github.com/Metisnetwork/Metis-Carrier/service/discovery"
	"github.com/Metisnetwork/Metis-Carrier/types"
	"github.com/urfave/cli/v2"
	"strconv"
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
	DebugAPIBackend *CarrierDebugAPIBackend

	resourceManager *resource.Manager
	messageManager  *message.MessageHandler
	TaskManager     handler.TaskManager
	authManager     *auth.AuthorityManager
	scheduler       schedule.Scheduler
	consulManager   *discovery.ConnectConsul
	runError        error
	metisPayManager *metispay.MetisPayManager
	quit            chan struct{}
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

	needReplayScheduleTaskCh, needExecuteTaskCh, taskConsResultCh :=
		make(chan *types.NeedReplayScheduleTask, config.TaskManagerConfig.NeedReplayScheduleTaskChanSize),
		make(chan *types.NeedExecuteTask, config.TaskManagerConfig.NeedExecuteTaskChanSize),
		make(chan *types.TaskConsResult, config.TaskManagerConfig.TaskConsResultChanSize)

	log.Debugf("Get some chan size value from config when carrier NewService, NeedReplayScheduleTaskChanSize: %d, NeedExecuteTaskChanSize: %d, TaskConsResultChanSize: %d",
		config.TaskManagerConfig.NeedReplayScheduleTaskChanSize, config.TaskManagerConfig.NeedExecuteTaskChanSize, config.TaskManagerConfig.TaskConsResultChanSize)

	resourceClientSet := grpclient.NewInternalResourceNodeSet()
	resourceMng := resource.NewResourceManager(config.CarrierDB, resourceClientSet, mockIdentityIdsFile)
	authManager := auth.NewAuthorityManager(config.CarrierDB)
	scheduler := schedule.NewSchedulerStarveFIFO(election.NewVrfElector(config.P2P.PirKey(), resourceMng), eventEngine, resourceMng, authManager)
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

	var metisPayManager *metispay.MetisPayManager

	if cliCtx.IsSet(flags.BlockChain.Name) {
		var metispayConfig *metispay.Config
		metispayConfig = &metispay.Config{URL: cliCtx.String(flags.BlockChain.Name)}

		var kmsConfig *kms.Config
		if cliCtx.IsSet(flags.KMSKeyId.Name) && cliCtx.IsSet(flags.KMSRegionId.Name) && cliCtx.IsSet(flags.KMSAccessKeyId.Name) && cliCtx.IsSet(flags.KMSAccessKeySecret.Name) {
			kmsConfig = &kms.Config{
				KeyId:           cliCtx.String(flags.KMSKeyId.Name),
				RegionId:        cliCtx.String(flags.KMSRegionId.Name),
				AccessKeyId:     cliCtx.String(flags.KMSAccessKeyId.Name),
				AccessKeySecret: cliCtx.String(flags.KMSAccessKeySecret.Name),
			}
		}
		metisPayManager = metispay.NewMetisPayManager(config.CarrierDB, metispayConfig, kmsConfig)
	}

	taskManager := task.NewTaskManager(
		config.P2P,
		scheduler,
		twopcEngine,
		eventEngine,
		resourceMng,
		authManager,
		metisPayManager,
		needReplayScheduleTaskCh,
		needExecuteTaskCh,
		taskConsResultCh,
		config.TaskManagerConfig,
	)

	s := &Service{
		ctx:             ctx,
		cancel:          cancel,
		config:          config,
		carrierDB:       config.CarrierDB,
		mempool:         pool,
		resourceManager: resourceMng,
		messageManager:  message.NewHandler(pool, resourceMng, taskManager, authManager),
		TaskManager:     taskManager,
		authManager:     authManager,
		metisPayManager: metisPayManager,
		scheduler:       scheduler,
		consulManager: discovery.NewConsulClient(&discovery.ConsulService{
			ServiceIP:   p2p.IpAddr().String(),
			ServicePort: strconv.Itoa(cliCtx.Int(flags.RPCPort.Name)),
			Tags:        config.DiscoverServiceConfig.DiscoveryServerTags,
			Name:        config.DiscoverServiceConfig.DiscoveryServiceName,
			Id:          config.DiscoverServiceConfig.DiscoveryServiceId,
			Interval:    config.DiscoverServiceConfig.DiscoveryServiceHealthCheckInterval,
			Deregister:  config.DiscoverServiceConfig.DiscoveryServiceHealthCheckDeregister,
		},
			config.DiscoverServiceConfig.DiscoveryServerIP,
			config.DiscoverServiceConfig.DiscoveryServerPort,
		),
		quit: make(chan struct{}),
	}

	//s.APIBackend = &CarrierAPIBackend{carrier: s}
	s.APIBackend = NewCarrierAPIBackend(s)
	s.DebugAPIBackend = NewCarrierDebugAPIBackend(twopcEngine)
	s.Engines = make(map[types.ConsensusEngineType]handler.Engine, 0)
	s.Engines[types.TwopcTyp] = twopcEngine
	s.Engines[types.ChainconsTyp] = chaincons.New()

	// load stored jobNode and dataNode
	jobNodeList, err := s.carrierDB.QueryRegisterNodeList(pb.PrefixTypeJobNode)
	if err == nil {
		for _, node := range jobNodeList {
			client, err := grpclient.NewJobNodeClient(ctx, fmt.Sprintf("%s:%s", node.GetInternalIp(), node.GetInternalPort()), node.GetId())
			if err == nil {
				s.resourceManager.StoreJobNodeClient(node.GetId(), client)
			}
		}
	}
	dataNodeList, err := s.carrierDB.QueryRegisterNodeList(pb.PrefixTypeDataNode)
	if err == nil {
		for _, node := range dataNodeList {
			client, err := grpclient.NewDataNodeClient(ctx, fmt.Sprintf("%s:%s", node.GetInternalIp(), node.GetInternalPort()), node.GetId())
			if err == nil {
				s.resourceManager.StoreDataNodeClient(node.GetId(), client)
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
	if err := s.initServicesWithDiscoveryCenter(); nil != err {
		log.Fatal(err)
	}
	go s.loop()
	return nil
}

func (s *Service) Stop() error {
	if s.cancel != nil {
		defer s.cancel()
	}
	s.carrierDB.Stop()

	for typ, engine := range s.Engines {
		if err := engine.Stop(); nil != err {
			log.WithError(err).Errorf("Cound not close the consensus engine: %s", typ.String())
		}
	}
	if nil != s.resourceManager {
		if err := s.resourceManager.Stop(); nil != err {
			log.WithError(err).Errorf("Failed to stop the resourceManager")
		}
	}
	if nil != s.messageManager {
		if err := s.messageManager.Stop(); nil != err {
			log.WithError(err).Errorf("Failed to stop the messageManager")
		}
	}
	if nil != s.TaskManager {
		if err := s.TaskManager.Stop(); nil != err {
			log.WithError(err).Errorf("Failed to stop the TaskManager")
		}
	}
	if nil != s.scheduler {
		if err := s.scheduler.Stop(); nil != err {
			log.WithError(err).Errorf("Failed to stop the schedule")
		}
	}

	if nil != s.authManager {
		if err := s.authManager.Stop(); nil != err {
			log.WithError(err).Errorf("Failed to stop the authManager")
		}
	}
	if nil != s.consulManager {
		if err := s.consulManager.DeregisterService2DiscoveryCenter(); nil != err {
			log.WithError(err).Errorf("Failed to deregister discover service")
		}
	}

	// stop service loop gorutine
	close(s.quit)
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
