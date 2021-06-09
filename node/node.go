package node

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/RosettaFlow/Carrier-Go/service"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// RosettaNode defines a struct that handles the services running a random rosetta net.
// It handles the lifecycle of the entire system and registers
// services to a service registry.
type RosettaNode struct {
	cliCtx          *cli.Context
	ctx             context.Context
	cancel          context.CancelFunc
	services        *common.ServiceRegistry
	lock            sync.RWMutex
	stop            chan struct{} // Channel to wait for termination notifications.
}

// New creates a new node instance, sets up configuration options, and registers
// every required service to the node.
func New(cliCtx *cli.Context) (*RosettaNode, error) {

	// todo: to init config

	registry := common.NewServiceRegistry()

	ctx, cancel := context.WithCancel(cliCtx.Context)
	rosetta := &RosettaNode{
		cliCtx:          cliCtx,
		ctx:             ctx,
		cancel:          cancel,
		services:        registry,
		stop:            make(chan struct{}),
	}
	// register P2P service
	if err := rosetta.registerP2P(cliCtx); err != nil {
		return nil, err
	}
	// register core backend service
	if err := rosetta.registerBackendService(); err != nil {
		return nil, err
	}

	if err := rosetta.registerSyncService(); err != nil {
		return nil, err
	}

	if err := rosetta.registerRPCService(); err != nil {
		return nil, err
	}
	// todo: some logic to be added here...
	return rosetta, nil
}

// Start the RosettaNode and kicks off every registered service.
func (b *RosettaNode) Start() {
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
func (b *RosettaNode) Close() {
	b.lock.Lock()
	defer b.lock.Unlock()

	log.Info("Stopping rosetta node")
	b.services.StopAll()
	b.cancel()
	close(b.stop)
}

func (b *RosettaNode) registerP2P(cliCtx *cli.Context) error {
	return nil
}

func (b *RosettaNode) registerBackendService() error {
	backendService, err := service.NewService(b.ctx, &params.RosettaConfig{}, &params.DataCenterConfig{})
	if err != nil {
		return errors.Wrap(err, "could not register backend service")
	}
	return b.services.RegisterService(backendService)
}

func (b *RosettaNode) registerSyncService() error {
	return nil
}

func (b *RosettaNode) registerRPCService() error {
	return nil
}
