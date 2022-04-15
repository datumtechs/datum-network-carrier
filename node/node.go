package node

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/carrier"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/debug"
	"github.com/RosettaFlow/Carrier-Go/common/feed"
	statefeed "github.com/RosettaFlow/Carrier-Go/common/feed/state"
	"github.com/RosettaFlow/Carrier-Go/common/flags"
	"github.com/RosettaFlow/Carrier-Go/common/sliceutil"
	"github.com/RosettaFlow/Carrier-Go/core"
	"github.com/RosettaFlow/Carrier-Go/db"
	"github.com/RosettaFlow/Carrier-Go/event"
	"github.com/RosettaFlow/Carrier-Go/gateway"
	"github.com/RosettaFlow/Carrier-Go/handler"
	"github.com/RosettaFlow/Carrier-Go/node/registration"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/rpc"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

// CarrierNode defines a struct that handles the services running a random rosetta net.
// It handles the lifecycle of the entire system and registers
// services to a service registry.
type CarrierNode struct {
	ctx       context.Context
	cancel    context.CancelFunc
	cliCtx    *cli.Context
	config    *Config
	services  *common.ServiceRegistry

	db        core.CarrierDB
	stateFeed *event.Feed
	lock      sync.RWMutex
	stop      chan struct{} // Channel to wait for termination notifications.
}

// New creates a new node instance, sets up configuration options, and registers
// every required service to the node.
func New(cliCtx *cli.Context) (*CarrierNode, error) {
	if err := configureTracing(cliCtx); err != nil {
		return nil, err
	}
	// todo: to init config
	config := makeConfig(cliCtx)
	configureNetwork(cliCtx)

	// Copy config and resolve the datadir so future changes to the current
	// working directory don't affect the node.
	confCopy := config.Node
	conf := &confCopy
	if conf.DataDir != "" {
		absdatadir, err := filepath.Abs(conf.DataDir)
		if err != nil {
			return nil, err
		}
		conf.DataDir = absdatadir
	}
	// init service register to accept some service.
	registry := common.NewServiceRegistry()

	ctx, cancel := context.WithCancel(cliCtx.Context)
	node := &CarrierNode{
		cliCtx:    cliCtx,
		ctx:       ctx,
		config:    conf,
		cancel:    cancel,
		services:  registry,
		stateFeed: new(event.Feed),
		stop:      make(chan struct{}),
	}

	// start db
	err := node.startDB(cliCtx, config.Carrier)
	if err != nil {
		log.WithError(err).Error("Failed to start DB")
		return nil, err
	}

	// register P2P service.
	if err := node.registerP2P(cliCtx); err != nil {
		return nil, err
	}

	// register backend service.
	if err := node.registerBackendService(config.Carrier, config.MockIdentityIdsFile,config.ConsensusStateFile); err != nil {
		return nil, err
	}

	// register network handler service.
	if err := node.registerHandlerService(); err != nil {
		return nil, err
	}

	// register rpc service.
	if err := node.registerRPCService(); err != nil {
		return nil, err
	}

	// register grpc gateway service.
	if err := node.registerGRPCGateway(); err != nil {
		return nil, err
	}

	// todo: some logic to be added here...
	return node, nil
}

func (node *CarrierNode) startDB (cliCtx *cli.Context, config *carrier.Config) error {
	dbPath := filepath.Join(node.config.DataDir, "datachain")
	log.WithField("database-path", dbPath).Info("Checking DB")
	config.DefaultConsensusWal = node.config.DataDir
	db, err := node.OpenDatabase(dbPath, config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return err
	}
	// setting database
	node.db = core.NewDataCenter(node.ctx, db)
	return nil
}

// OpenDatabase opens an existing database with the given name (or creates one
// if no previous can be found) from within the node's data directory. If the
// node is an ephemeral one, a memory database is returned.
func (node *CarrierNode) OpenDatabase(dbpath string, cache int, handles int) (db.Database, error) {
	if dbpath == "" {
		return db.NewMemoryDatabase(), nil
	}
	db, err := db.NewLDBDatabase(dbpath, cache, handles)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Start the CarrierNode and kicks off every registered service.
func (node *CarrierNode) Start() {
	node.lock.Lock()

	log.Info("Starting Carrier Node ...")

	node.services.StartAll()

	// -------------------------------------------------------------
	//TODO: mock, Temporarily set the initial success of the system
	go func() {
		node.stateFeed.Send(&feed.Event{
			Type: statefeed.Initialized,
			Data: &statefeed.InitializedData{
				StartTime: time.Now(),
			},
		})
	}()
	// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

	stop := node.stop
	node.lock.Unlock()

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigc)
		<-sigc
		log.Info("Got interrupt, shutting down...")
		debug.Exit(node.cliCtx) // Ensure trace and CPU profile data are flushed.
		go node.Close()
		for i := 10; i > 0; i-- {
			<-sigc
			if i > 1 {
				log.WithField("times", i-1).Info("Already shutting down, interrupt more to panic")
			}
		}
		panic("Panic closing the rosetta node")
	}()

	// Wait for stop channel to be closed.
	<-stop
}

// Close handles graceful shutdown of the system.
func (node *CarrierNode) Close() {
	node.lock.Lock()
	defer node.lock.Unlock()

	log.Info("Stopping carrier node")
	node.services.StopAll()
	node.cancel()
	close(node.stop)
}

func (node *CarrierNode) registerP2P(cliCtx *cli.Context) error {
	bootstrapNodeAddrs, dataDir, err := registration.P2PPreregistration(cliCtx)
	if err != nil {
		return err
	}
	staticNodeAddrs, err := registration.P2PStaticNodeAddrs(cliCtx, dataDir)
	if err != nil {
		return err
	}

	// load local seed node.
	seedNodeList, err := node.db.QuerySeedNodeList()
	localBootstrapAddr := make([]string, 0)
	if err == nil || len(seedNodeList) != 0 {
		for _, seedNode := range seedNodeList {
			localBootstrapAddr = append(localBootstrapAddr, seedNode.GetAddr())
		}
	}

	svc, err := p2p.NewService(node.ctx, &p2p.Config{
		EnableFakeNetwork:  cliCtx.Bool(flags.EnableFakeNetwork.Name),
		NoDiscovery:        cliCtx.Bool(flags.NoDiscovery.Name),
		StaticPeers:        staticNodeAddrs,
		BootstrapNodeAddr:  bootstrapNodeAddrs,
		LocalBootstrapAddr: localBootstrapAddr,
		RelayNodeAddr:      cliCtx.String(flags.RelayNode.Name),
		DataDir:            dataDir,
		LocalIP:            cliCtx.String(flags.P2PIP.Name),
		HostAddress:        cliCtx.String(flags.P2PHost.Name),
		HostDNS:            cliCtx.String(flags.P2PHostDNS.Name),
		PrivateKey:         cliCtx.String(flags.P2PPrivKey.Name),
		MetaDataDir:        cliCtx.String(flags.P2PMetadata.Name),
		TCPPort:            cliCtx.Uint(flags.P2PTCPPort.Name),
		UDPPort:            cliCtx.Uint(flags.P2PUDPPort.Name),
		MaxPeers:           cliCtx.Uint(flags.P2PMaxPeers.Name),
		AllowListCIDR:      cliCtx.String(flags.P2PAllowList.Name),
		DenyListCIDR:       sliceutil.SplitCommaSeparated(cliCtx.StringSlice(flags.P2PDenyList.Name)),
		EnableUPnP:         cliCtx.Bool(flags.EnableUPnPFlag.Name),
		DisableDiscv5:      cliCtx.Bool(flags.DisableDiscv5.Name),
		StateNotifier:      node,
	})
	if err != nil {
		return err
	}
	return node.services.RegisterService(svc)
}

func (node *CarrierNode) registerBackendService(carrierConfig *carrier.Config, mockIdentityIdsFile ,consensusStateFile string) error {
	carrierConfig.CarrierDB = node.db
	carrierConfig.P2P = node.fetchP2P()
	backendService, err := carrier.NewService(node.ctx,node.cliCtx, carrierConfig, mockIdentityIdsFile,consensusStateFile)
	if err != nil {
		return errors.Wrap(err, "could not register backend service")
	}
	return node.services.RegisterService(backendService)
}

func (node *CarrierNode) registerHandlerService() error {
	// use ` node.services.FetchService` to check whether the dependent service is registered.
	rs := handler.NewService(node.ctx, &handler.Config{
		P2P:           node.fetchP2P(),
		StateNotifier: node,
		Engines:       node.fetchBackend().Engines,
		TaskManager:   node.fetchBackend().TaskManager,
	})
	return node.services.RegisterService(rs)
}

func (node *CarrierNode) registerRPCService() error {
	backend := node.fetchRPCBackend()
	host := node.cliCtx.String(flags.RPCHost.Name)
	port := node.cliCtx.String(flags.RPCPort.Name)
	cert := node.cliCtx.String(flags.CertFlag.Name)
	key := node.cliCtx.String(flags.KeyFlag.Name)
	enableDebugRPCEndpoints := node.cliCtx.Bool(flags.EnableDebugRPCEndpoints.Name)
	maxMsgSize := node.cliCtx.Int(flags.GrpcMaxCallRecvMsgSizeFlag.Name)
	maxSendMsgSize := node.cliCtx.Int(flags.GrpcMaxCallSendMsgSizeFlag.Name)

	p2pService := node.fetchP2P()

	rpcService := rpc.NewService(node.ctx, &rpc.Config{
		Host:                    host,
		Port:                    port,
		CertFlag:                cert,
		KeyFlag:                 key,
		EnableDebugRPCEndpoints: enableDebugRPCEndpoints,
		Broadcaster:             p2pService,
		PeersFetcher:            p2pService,
		PeerManager:             p2pService,
		MetadataProvider:        p2pService,
		StateNotifier:           node,
		BackendAPI:              backend,
		MaxMsgSize:              maxMsgSize,
		MaxSendMsgSize:          maxSendMsgSize,
	})
	return node.services.RegisterService(rpcService)
}

func (node *CarrierNode) registerGRPCGateway() error {
	if node.cliCtx.Bool(flags.DisableGRPCGateway.Name) {
		return nil
	}
	gatewayPort := node.cliCtx.Int(flags.GRPCGatewayPort.Name)
	gatewayHost := node.cliCtx.String(flags.GRPCGatewayHost.Name)
	rpcHost := node.cliCtx.String(flags.RPCHost.Name)
	selfAddress := fmt.Sprintf("%s:%d", rpcHost, node.cliCtx.Int(flags.RPCPort.Name))
	gatewayAddress := fmt.Sprintf("%s:%d", gatewayHost, gatewayPort)
	allowedOrigins := strings.Split(node.cliCtx.String(flags.GPRCGatewayCorsDomain.Name), ",")
	enableDebugRPCEndpoints := node.cliCtx.Bool(flags.EnableDebugRPCEndpoints.Name)
	selfCert := node.cliCtx.String(flags.CertFlag.Name)

	return node.services.RegisterService(
		gateway.New(
			node.ctx,
			selfAddress,
			selfCert,
			gatewayAddress,
			nil, /*optional mux*/
			allowedOrigins,
			enableDebugRPCEndpoints,
			node.cliCtx.Uint64(flags.GrpcMaxCallRecvMsgSizeFlag.Name),
		),
	)
}

func (node *CarrierNode) fetchRPCBackend() backend.Backend {
	var s *carrier.Service
	if err := node.services.FetchService(&s); err != nil {
		panic(err)
	}
	return s.APIBackend
}

func (node *CarrierNode) fetchBackend() *carrier.Service {
	var s *carrier.Service
	if err := node.services.FetchService(&s); err != nil {
		panic(err)
	}
	return s
}

func (node *CarrierNode) fetchP2P() p2p.P2P {
	var p *p2p.Service
	if err := node.services.FetchService(&p); err != nil {
		panic(err)
	}
	return p
}

// StateFeed implements statefeed.Notifier.
func (node *CarrierNode) StateFeed() *event.Feed {
	return node.stateFeed
}

