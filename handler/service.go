package handler

import (
	"context"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/common/abool"
	"github.com/datumtechs/datum-network-carrier/common/feed"
	statefeed "github.com/datumtechs/datum-network-carrier/common/feed/state"
	"github.com/datumtechs/datum-network-carrier/common/runutil"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	"github.com/datumtechs/datum-network-carrier/p2p"
	"github.com/datumtechs/datum-network-carrier/params"
	carrierrpcdebugpbv1 "github.com/datumtechs/datum-network-carrier/pb/carrier/rpc/debug/v1"
	"github.com/datumtechs/datum-network-carrier/types"
	lru "github.com/hashicorp/golang-lru"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"math/rand"
	"sync"
	"time"
)

var _ common.Service = (*Service)(nil)

const rangeLimit = 1024
const seenBlockSize = 1000
const seenAttSize = 10000
const seenPrepareMsgSize = 1000
const seenPrepareVoteSize = 1000
const seenConfirmMsgSize = 1000
const seenConfirmVoteSize = 1000
const seenCommitMsgSize = 1000
const seenTaskResultMsgSize = 1000
const seenTaskResourceUsageMsgSize = 1000
const seenTaskTerminateMsgSize = 1000
const seenGossipTestDataSize = 100
const seenProposerSlashingSize = 100
const badBlockSize = 1000

const syncMetricsInterval = 10 * time.Second

var pendingBlockExpTime = time.Duration(params.CarrierConfig().SlotsPerEpoch.Mul(params.CarrierConfig().SecondsPerSlot)) * time.Second // Seconds in one epoch.

// Config to set up the regular sync service.
type Config struct {
	P2P           p2p.P2P
	Chain         blockchainService
	InitialSync   Checker
	StateNotifier statefeed.Notifier
	Engines       map[types.ConsensusEngineType]Engine
	TaskManager   TaskManager
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
	// Consensus-related
	seenPrepareMsgLock            sync.RWMutex
	seenPrepareMsgCache           *lru.Cache
	seenPrepareVoteLock           sync.RWMutex
	seenPrepareVoteCache          *lru.Cache
	seenConfirmMsgLock            sync.RWMutex
	seenConfirmMsgCache           *lru.Cache
	seenConfirmVoteLock           sync.RWMutex
	seenConfirmVoteCache          *lru.Cache
	seenCommitMsgLock             sync.RWMutex
	seenCommitMsgCache            *lru.Cache
	seenTaskResultMsgLock         sync.RWMutex
	seenTaskResultMsgCache        *lru.Cache
	seenTaskResourceUsageMsgLock  sync.RWMutex
	seenTaskResourceUsageMsgCache *lru.Cache
	seenTaskTerminateMsgLock      sync.RWMutex
	seenTaskTerminateMsgCache     *lru.Cache
	chainStarted                  *abool.AtomicBool
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
	// TODO: processPendingBlocksQueue,not needed for now
	//s.processPendingBlocksQueue()
	// TODO: Enable at the right time.
	//s.maintainPeerStatuses()
	// TODO: Update sync metrics.
	//runutil.RunEvery(s.ctx, syncMetricsInterval, s.updateMetrics)

	// for testing
	runutil.RunEvery(s.ctx, 5*time.Second, func() {
		sendPeer, _ := peer.Decode("16Uiu2HAm7pq7heDZwmrmWt9rXV8C1t5ENmyPKtToZAV6pioc1CrW")
		if s.cfg.P2P.PeerID() == sendPeer {
			err := s.cfg.P2P.Broadcast(s.ctx, &carrierrpcdebugpbv1.GossipTestData{
				Data:  []byte("data"),
				Count: rand.Uint64(),
				Step:  rand.Uint64(),
			})
			if err != nil {
				log.WithError(err).Error("Broadcast message failed")
			}
		}
	})

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
	prepareMsgCache, err := lru.New(seenPrepareMsgSize)
	if err != nil {
		return err
	}

	prepareVoteCache, err := lru.New(seenPrepareVoteSize)
	if err != nil {
		return err
	}
	confirmMsgCache, err := lru.New(seenConfirmMsgSize)
	if err != nil {
		return err
	}
	confirmVoteCache, err := lru.New(seenConfirmVoteSize)
	if err != nil {
		return err
	}
	commitMsgCache, err := lru.New(seenCommitMsgSize)
	if err != nil {
		return err
	}
	taskResultMsgCache, err := lru.New(seenTaskResultMsgSize)
	if err != nil {
		return err
	}
	taskResourceUsageMsgCache, err := lru.New(seenTaskResourceUsageMsgSize)
	if err != nil {
		return err
	}
	seenTaskTerminateMsgCache, err := lru.New(seenTaskTerminateMsgSize)
	if err != nil {
		return err
	}

	s.seenGossipDataCache = gossipCache
	s.seenPrepareMsgCache = prepareMsgCache
	s.seenPrepareVoteCache = prepareVoteCache
	s.seenConfirmMsgCache = confirmMsgCache
	s.seenConfirmVoteCache = confirmVoteCache
	s.seenCommitMsgCache = commitMsgCache
	s.seenTaskResultMsgCache = taskResultMsgCache
	s.seenTaskResourceUsageMsgCache = taskResourceUsageMsgCache
	s.seenTaskTerminateMsgCache = seenTaskTerminateMsgCache
	return nil
}

func (s *Service) registerHandlers() {
	// Wait until chain start.

	// Register respective rpc handlers at state initialized evengine.
	s.registerRPCHandlers()
	s.registerSubscribers()

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
				log.WithField("startTime", startTime).Debug("Received state initialized engine")

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
				// Register respective pubsub handlers at state synced evengine.
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
