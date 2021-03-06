package scorers

import (
	"github.com/RosettaFlow/Carrier-Go/p2p/peers/peerdata"
	"github.com/libp2p/go-libp2p-core/peer"
	"time"
)

var _ Scorer = (*BadResponsesScorer)(nil)

const (
	// DefaultBadResponsesThreshold defines how many bad responses to tolerate before peer is deemed bad.
	DefaultBadResponsesThreshold = 10

	// DefaultBadResponsesDecayInterval defines how often to decay previous statistics.
	// Every interval bad responses counter will be decremented by 1.
	DefaultBadResponsesDecayInterval = time.Hour
)

// BadResponsesScorer represents bad responses scoring service.
type BadResponsesScorer struct {
	config *BadResponsesScorerConfig
	store  *peerdata.Store
}

type BadResponsesScorerConfig struct {
	// Threshold specifies number of bad responses tolerated, before peer is banned.
	Threshold     int
	DecayInterval time.Duration
}

// newBadResponsesScorer creates new bad responses scoring service.
func newBadResponsesScorer(store *peerdata.Store, config *BadResponsesScorerConfig) *BadResponsesScorer {
	if config == nil {
		config = &BadResponsesScorerConfig{}
	}
	scorer := &BadResponsesScorer{
		config: config,
		store:  store,
	}
	if scorer.config.Threshold == 0 {
		scorer.config.Threshold = DefaultBadResponsesThreshold
	}
	if scorer.config.DecayInterval == 0 {
		scorer.config.DecayInterval = DefaultBadResponsesDecayInterval
	}
	return scorer
}

// Score returns score (penalty) of bad responses peer produced.
func (s *BadResponsesScorer) Score(pid peer.ID) float64 {
	s.store.RLock()
	defer s.store.RUnlock()
	return s.score(pid)
}

// score is a lock-free version of Score
func (s *BadResponsesScorer) score(pid peer.ID) float64 {
	if s.isBadPeer(pid) {
		return BadPeerScore
	}
	score := float64(0)
	peerData, ok := s.store.PeerData(pid)
	if !ok {
		return score
	}
	if peerData.BadResponses > 0 {
		score = float64(peerData.BadResponses) / float64(s.config.Threshold)
		// Since score represents a penalty, negate it.
		score *= -1
	}
	return score
}

// Params exposes scorer's parameters.
func (s *BadResponsesScorer) Params() *BadResponsesScorerConfig {
	return s.config
}

// Count obtains the number of bad responses we have received from the given remote peer.
func (s *BadResponsesScorer) Count(pid peer.ID) (int, error) {
	s.store.RLock()
	defer s.store.RUnlock()
	return s.count(pid)
}

func (s *BadResponsesScorer) count(pid peer.ID) (int, error) {
	if peerData, ok := s.store.PeerData(pid); ok {
		return peerData.BadResponses, nil
	}
	return -1, peerdata.ErrPeerUnknown
}

// Increment increments the number of bad responses we have received from the given remote peer.
func (s *BadResponsesScorer) Increment(pid peer.ID) {
	s.store.Lock()
	defer s.store.Unlock()

	peerData, ok := s.store.PeerData(pid)
	if !ok {
		s.store.SetPeerData(pid, &peerdata.PeerData{
			BadResponses: 1,
		})
		return
	}
	peerData.BadResponses++
}

// IsBadPeer states if the peer is to be considered bad.
func (s *BadResponsesScorer) IsBadPeer(pid peer.ID) bool {
	s.store.RLock()
	defer s.store.RUnlock()
	return s.isBadPeer(pid)
}

// isBadPeer is lock-free version of IsBadPeer.
func (s *BadResponsesScorer) isBadPeer(pid peer.ID) bool {
	//TODO: need to remove....
	/*if peerData, ok := s.store.PeerData(pid); ok {
		return peerData.BadResponses >= s.config.Threshold
	}*/
	return false
}

// BadPeers returns the peers that are considered bad.
func (s *BadResponsesScorer) BadPeers() []peer.ID {
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

// Decay reduces the bad responses of all peers, giving reformed peers a chance to join the network.
// This can be run periodically, although note that each time it runs it does give all bad peers another chance as well
// to clog up the network with bad responses, so should not be run frequently; once an hour would be reasonable.
func (s *BadResponsesScorer) Decay() {
	s.store.Lock()
	defer s.store.Unlock()

	for _, peerData := range s.store.Peers() {
		if peerData.BadResponses > 0 {
			peerData.BadResponses--
		}
	}
}
