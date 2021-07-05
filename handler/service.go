package handler

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/abool"
	"github.com/RosettaFlow/Carrier-Go/common/feed"
	statefeed "github.com/RosettaFlow/Carrier-Go/common/feed/state"
	"github.com/RosettaFlow/Carrier-Go/common/runutil"
	timeutils "github.com/RosettaFlow/Carrier-Go/common/timeutil"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/RosettaFlow/Carrier-Go/types"
	lru "github.com/hashicorp/golang-lru"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"sync"
	"time"
)

var _ common.Service = (*Service)(nil)

const rangeLimit = 1024
const seenBlockSize = 1000
const seenAttSize = 10000
const seenExitSize = 100
const seenGossipTestDataSize = 100
const seenProposerSlashingSize = 100
const badBlockSize = 1000

const syncMetricsInterval = 10 * time.Second

var pendingBlockExpTime = time.Duration(params.CarrierChainConfig().SlotsPerEpoch.Mul(params.CarrierChainConfig().SecondsPerSlot)) * time.Second // Seconds in one epoch.

// Config to set up the regular sync service.
type Config struct {
	P2P           p2p.P2P
	Chain         blockchainService
	InitialSync   Checker
	StateNotifier statefeed.Notifier
	Engines       map[types.ConsensusEngineType]Engine
}

// Service is responsible for handling all run time p2p related operations as the
// main entry point for network messages.
type Service struct {
	cfg                 *Config
	ctx                 context.Context
	cancel              context.CancelFunc
	rateLimiter         *limiter
	seenGossipDataLock  sync.RWMutex
	seenGossipDataCache *lru.Cache
	badBlockCache       *lru.Cache
	badBlockLock        sync.RWMutex
	chainStarted        *abool.AtomicBool
}

// NewService initializes new regular sync service.
func NewService(ctx context.Context, cfg *Config) *Service {
	rLimiter := newRateLimiter(cfg.P2P)
	ctx, cancel := context.WithCancel(ctx)
	r := &Service{
		cfg:          cfg,
		ctx:          ctx,
		cancel:       cancel,
		rateLimiter:  rLimiter,
		chainStarted: abool.New(),
	}
	log.Info("Init handler service...")
	go r.registerHandlers()
	return r
}

// Start the regular sync service.
func (s *Service) Start() error {
	if err := s.initCaches(); err != nil {
		panic(err)
	}
	s.cfg.P2P.AddConnectionHandler(s.reValidatePeer, s.sendGoodbye)
	s.cfg.P2P.AddDisconnectionHandler(func(_ context.Context, _ peer.ID) error {
		// no-op
		return nil
	})
	s.cfg.P2P.AddPingMethod(s.sendPingRequest)
	s.processPendingBlocksQueue()
	//s.maintainPeerStatuses()
	// Update sync metrics.
	runutil.RunEvery(s.ctx, syncMetricsInterval, s.updateMetrics)

	// for testing
	/*runutil.RunEvery(s.ctx, 10*time.Second, func() {
		peerId := s.cfg.P2P.Peers().Active()
		//target, _ := peer.Decode("16Uiu2HAm7pq7heDZwmrmWt9rXV8C1t5ENmyPKtToZAV6pioc1CrW")
		for _, id := range peerId {
			// 16Uiu2HAm7pq7heDZwmrmWt9rXV8C1t5ENmyPKtToZAV6pioc1CrW
			log.Infof("send to %s", id.String())
			SendGossipTestDataByRangeRequest(s.ctx, s.cfg.P2P,
				id,
				&libp2ppb.GossipTestData{
					Data:                 []byte("test data data"),
					Count:                10,
					Step:                 5,
				})
		}
	})*/

	log.Info("Starting handler service")
	return nil
}

// Stop the regular sync service.
func (s *Service) Stop() error {
	defer func() {
		if s.rateLimiter != nil {
			s.rateLimiter.free()
		}
	}()
	// Removing RPC Stream handlers.
	for _, p := range s.cfg.P2P.Host().Mux().Protocols() {
		s.cfg.P2P.Host().RemoveStreamHandler(protocol.ID(p))
	}
	// Deregister Topic Subscribers.
	for _, t := range s.cfg.P2P.PubSub().GetTopics() {
		if err := s.cfg.P2P.PubSub().UnregisterTopicValidator(t); err != nil {
			log.Errorf("Could not successfully unregister for topic %s: %v", t, err)
		}
	}
	defer s.cancel()
	return nil
}

// Status of the currently running regular sync service.
func (s *Service) Status() error {
	return nil
}

// This initializes the caches to update seen beacon objects coming in from the wire
// and prevent DoS.
func (s *Service) initCaches() error {
	gossipCache, err := lru.New(seenGossipTestDataSize)
	if err != nil {
		return err
	}
	s.seenGossipDataCache = gossipCache
	return nil
}

func (s *Service) registerHandlers() {
	// Wait until chain start.
	stateChannel := make(chan *feed.Event, 1)
	stateSub := s.cfg.StateNotifier.StateFeed().Subscribe(stateChannel)
	defer stateSub.Unsubscribe()
	for {
		select {
		case event := <-stateChannel:
			switch event.Type {
			case statefeed.Initialized:
				data, ok := event.Data.(*statefeed.InitializedData)
				if !ok {
					log.Error("Event feed data is not type *statefeed.InitializedData")
					return
				}
				startTime := data.StartTime
				log.WithField("starttime", startTime).Debug("Received state initialized event")

				// Register respective rpc handlers at state initialized event.
				s.registerRPCHandlers()
				// Wait for chainstart in separate routine.
				go func() {
					if startTime.After(timeutils.Now()) {
						time.Sleep(timeutils.Until(startTime))
					}
					log.WithField("starttime", startTime).Debug("Chain started in sync service")
					s.markForChainStart()
				}()
			case statefeed.Synced:
				_, ok := event.Data.(*statefeed.SyncedData)
				if !ok {
					log.Error("Event feed data is not type *statefeed.SyncedData")
					return
				}
				// Register respective pubsub handlers at state synced event.
				s.registerSubscribers()
				return
			}
		case <-s.ctx.Done():
			log.Debug("Context closed, exiting goroutine")
			return
		case err := <-stateSub.Err():
			log.WithError(err).Error("Could not subscribe to state notifier")
			return
		}
	}
}

// marks the chain as having started.
func (s *Service) markForChainStart() {
	s.chainStarted.Set()
}
