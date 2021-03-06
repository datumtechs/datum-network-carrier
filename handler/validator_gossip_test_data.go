package handler

import (
	"context"
	librpcpb "github.com/RosettaFlow/Carrier-Go/lib/rpc/v1"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Clients who receive a gossip test data on this topic MUST validate the conditions
// within process_gossip_test_data before forwarding it across the network.
func (s *Service) validateGossipTestData(ctx context.Context, pid peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
	// Validation runs on publish (not just subscriptions),
	// so we should approve any message from ourselves.
	if pid == s.cfg.P2P.PeerID() {
		return pubsub.ValidationAccept
	}

	m, err := s.decodePubsubMessage(msg)
	if err != nil {
		log.WithError(err).Debug("Could not decode message")
		return pubsub.ValidationReject
	}

	gossip, ok := m.(*librpcpb.SignedGossipTestData)
	if !ok {
		return pubsub.ValidationReject
	}

	if gossip.Data == nil {
		return pubsub.ValidationReject
	}
	if s.hasSeenGossipTestData(string(gossip.Data.GetData())) {
		return pubsub.ValidationIgnore
	}

	//TODO: more checking conditions...
	return pubsub.ValidationAccept
}

// Returns true if the node has already received a valid GossipTestData request.
func (s *Service) hasSeenGossipTestData(data string) bool {
	s.seenGossipDataLock.RLock()
	defer s.seenGossipDataLock.RUnlock()
	_, seen := s.seenGossipDataCache.Get(data)
	return seen
}

// Set exit request data in seen exit request cache.
func (s *Service) setGossipTestDataSeen(data string) {
	s.seenGossipDataLock.Lock()
	defer s.seenGossipDataLock.Unlock()
	s.seenGossipDataCache.Add(data, true)
}

