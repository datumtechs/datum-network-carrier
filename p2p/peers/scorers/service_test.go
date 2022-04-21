package scorers_test

import (
	"context"
	"github.com/Metisnetwork/Metis-Carrier/p2p/peers"
	"github.com/Metisnetwork/Metis-Carrier/p2p/peers/scorers"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"gotest.tools/assert"
	"math"
	"testing"
	"time"
)

func TestScorers_Service_Init(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//batchSize := uint64(flags.Get().BlockBatchLimit)

	t.Run("default config", func(t *testing.T) {
		peerStatuses := peers.NewStatus(ctx, &peers.StatusConfig{
			PeerLimit:    30,
			ScorerParams: &scorers.Config{},
		})

		t.Run("bad responses scorer", func(t *testing.T) {
			params := peerStatuses.Scorers().BadResponsesScorer().Params()
			assert.Equal(t, scorers.DefaultBadResponsesThreshold, params.Threshold, "Unexpected threshold value")
			assert.Equal(t, scorers.DefaultBadResponsesDecayInterval,
				params.DecayInterval, "Unexpected decay interval value")
		})
	})

	t.Run("explicit config", func(t *testing.T) {
		peerStatuses := peers.NewStatus(ctx, &peers.StatusConfig{
			PeerLimit: 30,
			ScorerParams: &scorers.Config{
				BadResponsesScorerConfig: &scorers.BadResponsesScorerConfig{
					Threshold:     2,
					DecayInterval: 1 * time.Minute,
				},
			},
		})

		t.Run("bad responses scorer", func(t *testing.T) {
			params := peerStatuses.Scorers().BadResponsesScorer().Params()
			assert.Equal(t, 2, params.Threshold, "Unexpected threshold value")
			assert.Equal(t, 1*time.Minute, params.DecayInterval, "Unexpected decay interval value")
		})
	})
}

func TestScorers_Service_Score(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	peerScores := func(s *scorers.Service, pids []peer.ID) map[string]float64 {
		scores := make(map[string]float64, len(pids))
		for _, pid := range pids {
			scores[string(pid)] = s.Score(pid)
		}
		return scores
	}
	_ = peerScores

	pack := func(scorer *scorers.Service, s1, s2, s3 float64) map[string]float64 {
		return map[string]float64{
			"peer1": roundScore(s1),
			"peer2": roundScore(s2),
			"peer3": roundScore(s3),
		}
	}
	_ = pack
	setupScorer := func() (*scorers.Service, []peer.ID) {
		peerStatuses := peers.NewStatus(ctx, &peers.StatusConfig{
			PeerLimit: 30,
			ScorerParams: &scorers.Config{
				BadResponsesScorerConfig: &scorers.BadResponsesScorerConfig{
					Threshold: 5,
				},
			},
		})
		s := peerStatuses.Scorers()
		pids := []peer.ID{"peer1", "peer2", "peer3"}
		for _, pid := range pids {
			peerStatuses.Add(nil, pid, nil, network.DirUnknown)
			// Not yet used peer gets boosted score.
			//startScore := s.BlockProviderScorer().MaxScore()
			//assert.Equal(t, startScore/float64(s.ActiveScorersCount()), s.Score(pid), "Unexpected score for not yet used peer")
		}
		return s, pids
	}
	_ = setupScorer

	t.Run("no peer registered", func(t *testing.T) {
		peerStatuses := peers.NewStatus(ctx, &peers.StatusConfig{
			ScorerParams: &scorers.Config{},
		})
		s := peerStatuses.Scorers()
		assert.Equal(t, 0.0, s.BadResponsesScorer().Score("peer1"))
		//assert.Equal(t, s.BlockProviderScorer().MaxScore(), s.BlockProviderScorer().Score("peer1"))
		assert.Equal(t, 0.0, s.Score("peer1"))
	})

	t.Run("bad responses score", func(t *testing.T) {
		//s, pids := setupScorer()
		// Peers start with boosted start score (new peers are boosted by block provider).
		//penalty := (-1 / float64(s.BadResponsesScorer().Params().Threshold)) / float64(s.ActiveScorersCount())

		// Update peers' stats and test the effect on peer order.
		//s.BadResponsesScorer().Increment("peer2")
		//assert.DeepEqual(t, pack(s, startScore, startScore+penalty, startScore), peerScores(s, pids))
		//s.BadResponsesScorer().Increment("peer1")
		//s.BadResponsesScorer().Increment("peer1")
		//assert.DeepEqual(t, pack(s, startScore+2*penalty, startScore+penalty, startScore), peerScores(s, pids))
		//
		//// See how decaying affects order of peers.
		//s.BadResponsesScorer().Decay()
		//assert.DeepEqual(t, pack(s, startScore+penalty, startScore, startScore), peerScores(s, pids))
		//s.BadResponsesScorer().Decay()
		//assert.DeepEqual(t, pack(s, startScore, startScore, startScore), peerScores(s, pids))
	})
}

func TestScorers_Service_loop(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	peerStatuses := peers.NewStatus(ctx, &peers.StatusConfig{
		PeerLimit: 30,
		ScorerParams: &scorers.Config{
			BadResponsesScorerConfig: &scorers.BadResponsesScorerConfig{
				Threshold:     5,
				DecayInterval: 50 * time.Millisecond,
			},
		},
	})
	s1 := peerStatuses.Scorers().BadResponsesScorer()

	pid1 := peer.ID("peer1")
	peerStatuses.Add(nil, pid1, nil, network.DirUnknown)
	for i := 0; i < s1.Params().Threshold+5; i++ {
		s1.Increment(pid1)
	}
	assert.Equal(t, true, s1.IsBadPeer(pid1), "Peer should be marked as bad")

	done := make(chan struct{}, 1)
	go func() {
		defer func() {
			done <- struct{}{}
		}()
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if s1.IsBadPeer(pid1) == false {
					return
				}
			case <-ctx.Done():
				t.Error("Timed out")
				return
			}
		}
	}()

	<-done
	assert.Equal(t, false, s1.IsBadPeer(pid1), "Peer should not be marked as bad")
}

func TestScorers_Service_IsBadPeer(t *testing.T) {
	peerStatuses := peers.NewStatus(context.Background(), &peers.StatusConfig{
		PeerLimit: 30,
		ScorerParams: &scorers.Config{
			BadResponsesScorerConfig: &scorers.BadResponsesScorerConfig{
				Threshold:     2,
				DecayInterval: 50 * time.Second,
			},
		},
	})

	assert.Equal(t, false, peerStatuses.Scorers().IsBadPeer("peer1"))
	peerStatuses.Scorers().BadResponsesScorer().Increment("peer1")
	peerStatuses.Scorers().BadResponsesScorer().Increment("peer1")
	assert.Equal(t, true, peerStatuses.Scorers().IsBadPeer("peer1"))
}

func TestScorers_Service_BadPeers(t *testing.T) {
	peerStatuses := peers.NewStatus(context.Background(), &peers.StatusConfig{
		PeerLimit: 30,
		ScorerParams: &scorers.Config{
			BadResponsesScorerConfig: &scorers.BadResponsesScorerConfig{
				Threshold:     2,
				DecayInterval: 50 * time.Second,
			},
		},
	})

	assert.Equal(t, false, peerStatuses.Scorers().IsBadPeer("peer1"))
	assert.Equal(t, false, peerStatuses.Scorers().IsBadPeer("peer2"))
	assert.Equal(t, false, peerStatuses.Scorers().IsBadPeer("peer3"))
	assert.Equal(t, 0, len(peerStatuses.Scorers().BadPeers()))
	for _, pid := range []peer.ID{"peer1", "peer3"} {
		peerStatuses.Scorers().BadResponsesScorer().Increment(pid)
		peerStatuses.Scorers().BadResponsesScorer().Increment(pid)
	}
	assert.Equal(t, true, peerStatuses.Scorers().IsBadPeer("peer1"))
	assert.Equal(t, false, peerStatuses.Scorers().IsBadPeer("peer2"))
	assert.Equal(t, true, peerStatuses.Scorers().IsBadPeer("peer3"))
	assert.Equal(t, 2, len(peerStatuses.Scorers().BadPeers()))
}

// roundScore returns score rounded in accordance with the score manager's rounding factor.
func roundScore(score float64) float64 {
	return math.Round(score*scorers.ScoreRoundingFactor) / scorers.ScoreRoundingFactor
}

