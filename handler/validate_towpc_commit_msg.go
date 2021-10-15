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
// within process_commit_message before forwarding it across the network.
func (s *Service) validateCommitMessagePubSub(ctx context.Context, pid peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
	// Validation runs on publish (not just subscriptions), so we should approve any message from ourselves.
	if pid == s.cfg.P2P.PeerID() {
		return pubsub.ValidationAccept
	}

	ctx, span := trace.StartSpan(ctx, "sync.commitMessage")
	defer span.End()

	m, err := s.decodePubsubMessage(msg)
	if err != nil {
		log.WithError(err).Debug("Could not decode message")
		traceutil.AnnotateError(span, err)
		return pubsub.ValidationReject
	}

	message, ok := m.(*twopcpb.CommitMsg)
	if !ok {
		log.Errorf("Invalid message type in the validateCommitMessagePubSub, typ: %T", m)
		return pubsub.ValidationReject
	}

	if s.hasSeenCommitMsg(message.MsgOption.ProposalId, message.MsgOption.SenderPartyId, message.MsgOption.ReceiverPartyId) {
		return pubsub.ValidationIgnore
	}

	//// validate CommitMsg
	//if err := s.validateCommitMsg(pid, message); err != nil {
	//	log.WithError(err).Errorf("Failed to call `validateCommitMsg`, proposalId: {%s}", common.BytesToHash(message.MsgOption.ProposalId).String())
	//	return pubsub.ValidationIgnore
	//}

	msg.ValidatorData = message // Used in downstream subscriber

	return pubsub.ValidationAccept
}

// Returns true if the node has already received a prepare message request for the validator with index `proposalId`.
func (s *Service) hasSeenCommitMsg(proposalId []byte, senderPartId []byte, receivePartId []byte) bool {
	s.seenCommitMsgLock.RLock()
	defer s.seenCommitMsgLock.RUnlock()
	v := append(proposalId, senderPartId...)
	v = append(v, receivePartId...)
	_, seen := s.seenCommitMsgCache.Get(string(v))
	return seen
}

// Set proposalId in seen exit request cache.
func (s *Service) setCommitMsgSeen(proposalId []byte, senderPartId []byte, receivePartId []byte) {
	s.seenCommitMsgLock.Lock()
	defer s.seenCommitMsgLock.Unlock()
	v := append(proposalId, senderPartId...)
	v = append(v, receivePartId...)
	s.seenCommitMsgCache.Add(string(v), true)
}