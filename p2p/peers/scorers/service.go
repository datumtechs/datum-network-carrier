package scorers

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/p2p/peers/peerdata"
	"github.com/libp2p/go-libp2p-core/peer"
	"math"
	"time"
)

//var _ Scorer = (*Service)(nil)

// ScoreRoundingFactor defines how many digits to keep in decimal part.
// This parameter is used in math.Round(score*ScoreRoundingFactor) / ScoreRoundingFactor.
const ScoreRoundingFactor = 10000

// BadPeerScore defines score that is returned for a bad peer (all other metrics are ignored).
const BadPeerScore = -1.00

// Scorer defines minimum set of methods every peer scorer must expose.
type Scorer interface {
	Score(pid peer.ID) float64
	IsBadPeer(pid peer.ID) bool
	BadPeers() []peer.ID
}

// Service manages peer scorers that are used to calculate overall peer score.
type Service struct {
	store   *peerdata.Store
	scorers struct {
		badResponsesScorer *BadResponsesScorer
		peerStatusScorer   *PeerStatusScorer
		gossipScorer       *GossipScorer
	}
	weights     map[Scorer]float64
	totalWeight float64
}

// Config holds configuration parameters for scoring service.
type Config struct {
	BadResponsesScorerConfig *BadResponsesScorerConfig
	PeerStatusScorerConfig   *PeerStatusScorerConfig
	GossipScorerConfig       *GossipScorerConfig
}

// NewService provides fully initialized peer scoring service.
func NewService(ctx context.Context, store *peerdata.Store, config *Config) *Service {
	s := &Service{
		store:   store,
		weights: make(map[Scorer]float64),
	}
	// Register scorers.
	s.scorers.badResponsesScorer = newBadResponsesScorer(store, config.BadResponsesScorerConfig)
	s.setScorerWeight(s.scorers.badResponsesScorer, 1.0)
	s.scorers.peerStatusScorer = newPeerStatusScorer(store, config.PeerStatusScorerConfig)
	s.setScorerWeight(s.scorers.peerStatusScorer, 0.0)
	s.scorers.gossipScorer = newGossipScorer(store, config.GossipScorerConfig)
	s.setScorerWeight(s.scorers.gossipScorer, 0.0)

	// start background tasks.
	go s.loop(ctx)

	return s
}

// BadResponsesScorer exposes bad responses scoring service.
func (s *Service) BadResponsesScorer() *BadResponsesScorer {
	return s.scorers.badResponsesScorer
}

// PeerStatusScorer exposes peer chain status scoring service.
func (s *Service) PeerStatusScorer() *PeerStatusScorer {
	return s.scorers.peerStatusScorer
}

// GossipScorer exposes the peer's gossip scoring service.
func (s *Service) GossipScorer() *GossipScorer {
	return s.scorers.gossipScorer
}

// ActiveScorersCount returns number of scorers that can affect score (have non-zero weight).
func (s *Service) ActiveScorersCount() int {
	cnt := 0
	for _, w := range s.weights {
		if w > 0 {
			cnt++
		}
	}
	return cnt
}

// Score returns calculated peer score across all tracked metrics.
func (s *Service) Score(pid peer.ID) float64 {
	s.store.RLock()
	defer s.store.RUnlock()

	score := float64(0)
	if _, ok := s.store.PeerData(pid); !ok {
		return 0
	}
	score += s.scorers.badResponsesScorer.score(pid) * s.scorerWeight(s.scorers.badResponsesScorer)
	score += s.scorers.peerStatusScorer.score(pid) * s.scorerWeight(s.scorers.peerStatusScorer)
	score += s.scorers.gossipScorer.score(pid) * s.scorerWeight(s.scorers.gossipScorer)
	return math.Round(score*ScoreRoundingFactor) / ScoreRoundingFactor
}

// IsBadPeer traverses all the scorers to see if any of them classifies peer as bad.
func (s *Service) IsBadPeer(pid peer.ID) bool {
	s.store.RLock()
	defer s.store.RUnlock()
	return s.isBadPeer(pid)
}

// isBadPeer is a lock-free version of isBadPeer.
func (s *Service) isBadPeer(pid peer.ID) bool {
	if s.scorers.badResponsesScorer.isBadPeer(pid) {
		return true
	}
	if s.scorers.peerStatusScorer.isBadPeer(pid) {
		return true
	}
	// TODO(#6043): Hook in gossip scorer's relevant
	// method to check if peer has a bad gossip score.
	return false
}

// BadPeers returns the peers that are considered bad by any of registered scorers.
func (s *Service) BadPeers() []peer.ID {
	s.store.RLock()
	defer s.store.RUnlock()

	badPeers := make([]peer.ID, 0)
	for pid := range s.store.Peers() {
		if s.isBadPeer(pid) {
			badPeers = append(badPeers, pid)
		}
	}
	return badPeers
}

// ValidationError returns peer data validation error, which potentially provides more information
// why peer is considered bad.
func (s *Service) ValidationError(pid peer.ID) error {
	s.store.RLock()
	defer s.store.RUnlock()

	peerData, ok := s.store.PeerData(pid)
	if !ok {
		return nil
	}
	return peerData.ChainStateValidationError
}

// loop handles background tasks.
func (s *Service) loop(ctx context.Context) {
	decayBadResponsesStats := time.NewTicker(s.scorers.badResponsesScorer.Params().DecayInterval)
	defer decayBadResponsesStats.Stop()
	for {
		select {
		case <-decayBadResponsesStats.C:
			s.scorers.badResponsesScorer.Decay()
		case <-ctx.Done():
			return
		}
	}
}

// setScorerWeight adds scorer to map of known scorers.
func (s *Service) setScorerWeight(scorer Scorer, weight float64) {
	s.weights[scorer] = weight
	s.totalWeight += s.weights[scorer]
}

// scorerWeight calculates contribution percentage of a given scorer in total score.
func (s *Service) scorerWeight(scorer Scorer) float64 {
	return s.weights[scorer] / s.totalWeight
}
