package node

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/carrier"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// CarrierNode defines a struct that handles the services running a random rosetta net.
// It handles the lifecycle of the entire system and registers
// services to a service registry.
type CarrierNode struct {
	cliCtx          *cli.Context
	ctx             context.Context
	cancel          context.CancelFunc
	services        *common.ServiceRegistry
	lock            sync.RWMutex
	stop            chan struct{} // Channel to wait for termination notifications.
}

// New creates a new node instance, sets up configuration options, and registers
// every required service to the node.
func New(cliCtx *cli.Context) (*CarrierNode, error) {

	// todo: to init config

	registry := common.NewServiceRegistry()

	ctx, cancel := context.WithCancel(cliCtx.Context)
	carrier := &CarrierNode{
		cliCtx:          cliCtx,
		ctx:             ctx,
		cancel:          cancel,
		services:        registry,
		stop:            make(chan struct{}),
	}
	// register P2P service
	if err := carrier.registerP2P(cliCtx); err != nil {
		return nil, err
	}
	// register core backend service
	if err := carrier.registerBackendService(); err != nil {
		return nil, err
	}

	if err := carrier.registerSyncService(); err != nil {
		return nil, err
	}

	if err := carrier.registerRPCService(); err != nil {
		return nil, err
	}
	// todo: some logic to be added here...
	return carrier, nil
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

func (b *CarrierNode) registerBackendService() error {
	backendService, err := carrier.NewService(b.ctx, &params.CarrierConfig{}, &params.DataCenterConfig{})
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
