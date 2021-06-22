package scores

import (
	"github.com/RosettaFlow/Carrier-Go/p2p/peers/peerdata"
	"github.com/libp2p/go-libp2p-core/peer"
)

var _ Scorer = (*Service)(nil)

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
		//badResponsesScorer  *BadResponsesScorer
		//blockProviderScorer *BlockProviderScorer
		//peerStatusScorer    *PeerStatusScorer
		//gossipScorer        *GossipScorer
	}
	weights     map[Scorer]float64
	totalWeight float64
}

func (s Service) Score(pid peer.ID) float64 {
	panic("implement me")
}

func (s Service) IsBadPeer(pid peer.ID) bool {
	panic("implement me")
}

func (s Service) BadPeers() []peer.ID {
	panic("implement me")
}
