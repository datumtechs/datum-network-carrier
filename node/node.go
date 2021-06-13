package node

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/carrier"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/db"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

// CarrierNode defines a struct that handles the services running a random rosetta net.
// It handles the lifecycle of the entire system and registers
// services to a service registry.
type CarrierNode struct {
	cliCtx          *cli.Context
	config 			*Config
	ctx             context.Context
	cancel          context.CancelFunc
	services        *common.ServiceRegistry

	db 				db.Database

	lock            sync.RWMutex
	stop            chan struct{} // Channel to wait for termination notifications.
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

	registry := common.NewServiceRegistry()

	ctx, cancel := context.WithCancel(cliCtx.Context)
	node := &CarrierNode{
		cliCtx:          cliCtx,
		ctx:             ctx,
		config: 		 conf,
		cancel:          cancel,
		services:        registry,
		stop:            make(chan struct{}),
	}

	// start db
	node.startDB(cliCtx, &cfg.Carrier)

	// register P2P service
	if err := node.registerP2P(cliCtx); err != nil {
		return nil, err
	}
	// register core backend service
	if err := node.registerBackendService(&cfg.Carrier); err != nil {
		return nil, err
	}

	if err := node.registerSyncService(); err != nil {
		return nil, err
	}

	if err := node.registerRPCService(); err != nil {
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
	node.db = db
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
	return nil
}

func (b *CarrierNode) registerBackendService(carrierConfig *carrier.Config) error {
	backendService, err := carrier.NewService(b.ctx, carrierConfig, &params.DataCenterConfig{})
	if err != nil {
		return errors.Wrap(err, "could not register backend service")
	}
	return b.services.RegisterService(backendService)
}

func (b *CarrierNode) registerSyncService() error {
	return nil
}

func (b *CarrierNode) registerRPCService() error {
	return nil
}
