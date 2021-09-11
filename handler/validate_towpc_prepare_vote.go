package handler

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/traceutil"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"go.opencensus.io/trace"
)

// Clients who receive a prepare message on this topic MUST validate the conditions
// within process_prepare_vote before forwarding it across the network.
func (s *Service) validatePrepareVotePubSub(ctx context.Context, pid peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
	// Validation runs on publish (not just subscriptions), so we should approve any message from ourselves.
	if pid == s.cfg.P2P.PeerID() {
		return pubsub.ValidationAccept
	}

	ctx, span := trace.StartSpan(ctx, "sync.prepareVote")
	defer span.End()

	m, err := s.decodePubsubMessage(msg)
	if err != nil {
		log.WithError(err).Debug("Could not decode message")
		traceutil.AnnotateError(span, err)
		return pubsub.ValidationReject
	}

	message, ok := m.(*twopcpb.PrepareVote)
	if !ok {
		log.Errorf("Invalid message type in the validatePrepareVotePubSub, typ: %T", m)
		return pubsub.ValidationReject
	}

	if s.hasSeenPrepareVote(message.MsgOption.ProposalId) {
		return pubsub.ValidationIgnore
	}

	// validate prepareVote
	if err := s.validatePrepareVote(pid, message); err != nil {
		log.WithError(err).Errorf("Failed to call `validatePrepareVote`, proposalId: {%s}", common.BytesToHash(message.MsgOption.ProposalId).String())
		return pubsub.ValidationIgnore
	}

	msg.ValidatorData = message // Used in downstream subscriber

	return pubsub.ValidationAccept
}

// Returns true if the node has already received a prepare message request for the validator with index `proposalId`.
func (s *Service) hasSeenPrepareVote(proposalId []byte) bool {
	s.seenPrepareVoteLock.RLock()
	defer s.seenPrepareVoteLock.RUnlock()
	_, seen := s.seenPrepareVoteCache.Get(string(proposalId))
	return seen
}

// Set proposalId in seen exit request cache.
func (s *Service) setPrepareVoteSeen(proposalId []byte) {
	s.seenPrepareMsgLock.Lock()
	defer s.seenPrepareMsgLock.Unlock()
	s.seenPrepareVoteCache.Add(string(proposalId), true)
}