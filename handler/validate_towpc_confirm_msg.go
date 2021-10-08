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
// within process_confirm_message before forwarding it across the network.
func (s *Service) validateConfirmMessagePubSub(ctx context.Context, pid peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
	// Validation runs on publish (not just subscriptions), so we should approve any message from ourselves.
	if pid == s.cfg.P2P.PeerID() {
		return pubsub.ValidationAccept
	}

	ctx, span := trace.StartSpan(ctx, "sync.confirmMessage")
	defer span.End()

	m, err := s.decodePubsubMessage(msg)
	if err != nil {
		log.WithError(err).Debug("Could not decode message")
		traceutil.AnnotateError(span, err)
		return pubsub.ValidationReject
	}

	message, ok := m.(*twopcpb.ConfirmMsg)
	if !ok {
		log.Errorf("Invalid message type in the validateConfirmMessagePubSub, typ: %T", m)
		return pubsub.ValidationReject
	}

	if s.hasSeenConfirmMsg(message.MsgOption.ProposalId, message.MsgOption.SenderPartyId) {
		return pubsub.ValidationIgnore
	}

	//// validate ConfirmMsg
	//if err := s.validateConfirmMsg(pid, message); err != nil {
	//	log.WithError(err).Errorf("Failed to call `validateConfirmMsg`, proposalId: {%s}", common.BytesToHash(message.MsgOption.ProposalId).String())
	//	return pubsub.ValidationIgnore
	//}

	msg.ValidatorData = message // Used in downstream subscriber
	return pubsub.ValidationAccept
}

// Returns true if the node has already received a prepare message request for the validator with index `proposalId`.
func (s *Service) hasSeenConfirmMsg(proposalId []byte, taskPartyId []byte) bool {
	s.seenConfirmMsgLock.RLock()
	defer s.seenConfirmMsgLock.RUnlock()
	v := append(proposalId, taskPartyId...)
	_, seen := s.seenConfirmMsgCache.Get(string(v))
	return seen
}

// Set proposalId in seen exit request cache.
func (s *Service) setConfirmMsgSeen(proposalId []byte, taskPartyId []byte) {
	s.seenConfirmMsgLock.Lock()
	defer s.seenConfirmMsgLock.Unlock()
	v := append(proposalId, taskPartyId...)
	s.seenConfirmMsgCache.Add(string(v), true)
}