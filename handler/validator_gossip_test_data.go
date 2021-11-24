package handler

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/common/traceutil"
	librpcpb "github.com/RosettaFlow/Carrier-Go/lib/rpc/v1"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"go.opencensus.io/trace"
)

// Clients who receive a gossip test data on this topic MUST validate the conditions
// within process_gossip_test_data before forwarding it across the network.
func (s *Service) validateGossipTestData(ctx context.Context, pid peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
	// Validation runs on publish (not just subscriptions),
	// so we should approve any message from ourselves.
	if pid == s.cfg.P2P.PeerID() {
		return pubsub.ValidationAccept
	}

	ctx, span := trace.StartSpanWithRemoteParent(ctx, "handler.validateGossipTestData",
		traceutil.GenerateParentSpanWithGossipTestData(pid, msg))
	defer span.End()

	m, err := s.decodePubsubMessage(msg)
	if err != nil {
		log.WithError(err).Debug("Could not decode message")
		traceutil.AnnotateError(span, err)
		return pubsub.ValidationReject
	}

	gossip, ok := m.(*librpcpb.GossipTestData)
	if !ok {
		log.Errorf("Invalid message type in the validateGossipTestData, typ: %T", m)
		return pubsub.ValidationReject
	}

	if gossip.Data == nil {
		return pubsub.ValidationReject
	}
	if s.hasSeenGossipTestData(string(gossip.GetData())) {
		return pubsub.ValidationIgnore
	}

	msg.ValidatorData = gossip // Used in downstream subscriber

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

