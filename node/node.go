package node

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/carrier"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/feed"
	statefeed "github.com/RosettaFlow/Carrier-Go/common/feed/state"
	"github.com/RosettaFlow/Carrier-Go/common/flags"
	"github.com/RosettaFlow/Carrier-Go/consensus/chaincons"
	"github.com/RosettaFlow/Carrier-Go/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/core"
	"github.com/RosettaFlow/Carrier-Go/db"
	"github.com/RosettaFlow/Carrier-Go/event"
	"github.com/RosettaFlow/Carrier-Go/handler"
	"github.com/RosettaFlow/Carrier-Go/node/registration"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/RosettaFlow/Carrier-Go/rpc"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/shared/sliceutil"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"path/filepath"
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
	Engines   map[types.ConsensusEngineType]handler.Engine
	db        core.CarrierDB
	stateFeed *event.Feed
	lock      sync.RWMutex
	stop      chan struct{} // Channel to wait for termination notifications.
}

// New creates a new node instance, sets up configuration options, and registers
// every required service to the node.
func New(cliCtx *cli.Context) (*CarrierNode, error) {
	// todo: to init config
	cfg := makeConfig(cliCtx)

	// Copy config and resolve the datadir so future changes to the current
	// working directory don't affect the node.
	confCopy := cfg.Node
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

	if err := node.registerConsensusEngine(); nil != err {
		log.Error("Failed to registerConsensusEngine", "err", err)
		return nil, err
	}

	// start db
	err := node.startDB(cliCtx, &cfg.Carrier)
	if err != nil {
		log.WithError(err).Error("Failed to start DB")
		return nil, err
	}
	// register P2P service
	if err := node.registerP2P(cliCtx); err != nil {
		return nil, err
	}
	if err := node.registerBackendService(&cfg.Carrier); err != nil {
		return nil, err
	}
	//
	if err := node.registerHandlerService(); err != nil {
		return nil, err
	}

	//TODO: need to config rpc
	if err := node.registerRPCService(&rpc.RpcConfig{
		Ip:       "0.0.0.0",
		Port:     "8888",
	}); err != nil {
		return nil, err
	}
	// todo: some logic to be added here...
	return node, nil
}

func (node *CarrierNode) startDB(cliCtx *cli.Context, config *carrier.Config) error {
	dbPath := filepath.Join(node.config.DataDir, "datachain")
	log.WithField("database-path", dbPath).Info("Checking DB")
	db, err := node.OpenDatabase(dbPath, config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return err
	}

	// setting database
	carrierDB, err := core.NewDataCenter(node.ctx, db, &params.DataCenterConfig{
		GrpcUrl: "192.168.112.32",
		Port:    9099,
	})
	if err != nil {
		return err
	}
	node.db = carrierDB
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
func (b *CarrierNode) Start() {
	b.lock.Lock()

	log.Info("Starting rosetta node")

	b.services.StartAll()

	// -------------------------------------------------------------
	//TODO: mock, Temporarily set the initial success of the system
	b.stateFeed.Send(&feed.Event{
		Type: statefeed.Initialized,
		Data: &statefeed.InitializedData{
			StartTime: time.Now(),
		},
	})
	// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

	stop := b.stop
	b.lock.Unlock()

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigc)
		<-sigc
		log.Info("Got interrupt, shutting down...")
		//debug.Exit(b.cliCtx) // Ensure trace and CPU profile data are flushed.
		go b.Close()
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
func (b *CarrierNode) Close() {
	b.lock.Lock()
	defer b.lock.Unlock()

	log.Info("Stopping rosetta node")
	b.services.StopAll()
	b.cancel()
	close(b.stop)
}

func (b *CarrierNode) registerP2P(cliCtx *cli.Context) error {
	bootstrapNodeAddrs, dataDir, err := registration.P2PPreregistration(cliCtx)
	if err != nil {
		return err
	}
	staticNodeAddrs, err := registration.P2PStaticNodeAddrs(cliCtx, dataDir)
	if err != nil {
		return err
	}
	svc, err := p2p.NewService(b.ctx, &p2p.Config{
		NoDiscovery:       cliCtx.Bool(flags.NoDiscovery.Name),
		StaticPeers:       staticNodeAddrs,
		BootstrapNodeAddr: bootstrapNodeAddrs,
		RelayNodeAddr:     cliCtx.String(flags.RelayNode.Name),
		DataDir:           dataDir,
		LocalIP:           cliCtx.String(flags.P2PIP.Name),
		HostAddress:       cliCtx.String(flags.P2PHost.Name),
		HostDNS:           cliCtx.String(flags.P2PHostDNS.Name),
		PrivateKey:        cliCtx.String(flags.P2PPrivKey.Name),
		MetaDataDir:       cliCtx.String(flags.P2PMetadata.Name),
		TCPPort:           cliCtx.Uint(flags.P2PTCPPort.Name),
		UDPPort:           cliCtx.Uint(flags.P2PUDPPort.Name),
		MaxPeers:          cliCtx.Uint(flags.P2PMaxPeers.Name),
		AllowListCIDR:     cliCtx.String(flags.P2PAllowList.Name),
		DenyListCIDR:      sliceutil.SplitCommaSeparated(cliCtx.StringSlice(flags.P2PDenyList.Name)),
		EnableUPnP:        cliCtx.Bool(flags.EnableUPnPFlag.Name),
		DisableDiscv5:     cliCtx.Bool(flags.DisableDiscv5.Name),
		StateNotifier:     b,
	})
	if err != nil {
		return err
	}
	return b.services.RegisterService(svc)
}

func (b *CarrierNode) registerBackendService(carrierConfig *carrier.Config) error {
	backendService, err := carrier.NewService(b.ctx, carrierConfig)
	if err != nil {
		return errors.Wrap(err, "could not register backend service")
	}
	return b.services.RegisterService(backendService)
}

func (b *CarrierNode) registerHandlerService() error {
	// use ` b.services.FetchService` to check whether the dependent service is registered.
	rs := handler.NewService(b.ctx, &handler.Config{
		P2P:           b.fetchP2P(),
		StateNotifier: b,
		Engines:       b.Engines,
	})
	return b.services.RegisterService(rs)
}

func (b *CarrierNode) registerRPCService(config *rpc.RpcConfig) error {
	rpcService := rpc.New(config, b.fetchBackend())
	return b.services.RegisterService(rpcService)
}

func (b *CarrierNode) fetchBackend() rpc.Backend {
	var s *carrier.Service
	if err := b.services.FetchService(&s); err != nil {
		panic(err)
	}
	return s.APIBackend
}

func (b *CarrierNode) fetchP2P() p2p.P2P {
	var p *p2p.Service
	if err := b.services.FetchService(&p); err != nil {
		panic(err)
	}
	return p
}

// StateFeed implements statefeed.Notifier.
func (b *CarrierNode) StateFeed() *event.Feed {
	return b.stateFeed
}

func (b *CarrierNode) registerConsensusEngine() error {

	b.Engines = make(map[types.ConsensusEngineType]handler.Engine, 0)
	b.Engines[types.TwopcTyp] = twopc.New(&twopc.Config{})
	b.Engines[types.ChainconsTyp] = chaincons.New()

	return nil
}
