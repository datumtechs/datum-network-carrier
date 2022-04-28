package handler

import (
	"context"
	"github.com/Metisnetwork/Metis-Carrier/common/traceutil"
	twopcpb "github.com/Metisnetwork/Metis-Carrier/lib/netmsg/consensus/twopc"
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

	ctx, span := trace.StartSpanWithRemoteParent(ctx, "handler.validateConfirmVotePubSub",
		traceutil.GenerateParentSpanWithConfirmVote(pid, msg))
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

	if s.hasSeenConfirmVote(message.MsgOption.ProposalId, message.MsgOption.SenderPartyId, message.MsgOption.ReceiverPartyId) {
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
func (s *Service) hasSeenConfirmVote(proposalId []byte, senderPartId []byte, receivePartId []byte) bool {
	s.seenConfirmVoteLock.RLock()
	defer s.seenConfirmVoteLock.RUnlock()
	v := append(proposalId, senderPartId...)
	v = append(v, receivePartId...)
	_, seen := s.seenConfirmVoteCache.Get(string(v))
	return seen
}

// Set proposalId in seen exit request cache.
func (s *Service) setConfirmVoteSeen(proposalId []byte, senderPartId []byte, receivePartId []byte) {
	s.seenConfirmVoteLock.Lock()
	defer s.seenConfirmVoteLock.Unlock()
	v := append(proposalId, senderPartId...)
	v = append(v, receivePartId...)
	s.seenConfirmVoteCache.Add(string(v), true)
}