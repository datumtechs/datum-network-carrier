package handler

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/common/traceutil"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"go.opencensus.io/trace"
)

// Clients who receive a prepare message on this topic MUST validate the conditions
// within process_confirm_vote before forwarding it across the network.
func (s *Service) validateConfirmVotePubSub(ctx context.Context, pid peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
	// Validation runs on publish (not just subscriptions), so we should approve any message from ourselves.
	if pid == s.cfg.P2P.PeerID() {
		return pubsub.ValidationAccept
	}

	ctx, span := trace.StartSpan(ctx, "sync.confirmVote")
	defer span.End()

	m, err := s.decodePubsubMessage(msg)
	if err != nil {
		log.WithError(err).Debug("Could not decode message")
		traceutil.AnnotateError(span, err)
		return pubsub.ValidationReject
	}

	message, ok := m.(*twopcpb.ConfirmVote)
	if !ok {
		log.Errorf("Invalid message type in the validateConfirmVotePubSub, typ: %T", m)
		return pubsub.ValidationReject
	}

	if s.hasSeenConfirmVote(message.MsgOption.ProposalId) {
		return pubsub.ValidationIgnore
	}

	//// validate ConfirmVote
	//if err := s.validateConfirmVote(pid, message); err != nil {
	//	log.WithError(err).Errorf("Failed to call `validateConfirmVote`, proposalId: {%s}", common.BytesToHash(message.MsgOption.ProposalId).String())
	//	return pubsub.ValidationIgnore
	//}

	msg.ValidatorData = message // Used in downstream subscriber
	return pubsub.ValidationAccept
}

// Returns true if the node has already received a prepare message request for the validator with index `proposalId`.
func (s *Service) hasSeenConfirmVote(proposalId []byte) bool {
	s.seenConfirmVoteLock.RLock()
	defer s.seenConfirmVoteLock.RUnlock()
	_, seen := s.seenConfirmVoteCache.Get(string(proposalId))
	return seen
}

// Set proposalId in seen exit request cache.
func (s *Service) setConfirmVoteSeen(proposalId []byte) {
	s.seenConfirmVoteLock.Lock()
	defer s.seenConfirmVoteLock.Unlock()
	s.seenConfirmVoteCache.Add(string(proposalId), true)
}