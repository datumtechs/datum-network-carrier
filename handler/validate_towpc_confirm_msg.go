package handler

import (
	"context"
	"github.com/datumtechs/datum-network-carrier/common/traceutil"
	twopcpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/consensus/twopc"
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

	ctx, span := trace.StartSpanWithRemoteParent(ctx, "handler.validateConfirmMessagePubSub",
		traceutil.GenerateParentSpanWithConfirmMsg(pid, msg))
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

	if s.hasSeenConfirmMsg(message.MsgOption.ProposalId, message.MsgOption.SenderPartyId, message.MsgOption.ReceiverPartyId) {
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
func (s *Service) hasSeenConfirmMsg(proposalId []byte, senderPartId []byte, receivePartId []byte) bool {
	s.seenConfirmMsgLock.RLock()
	defer s.seenConfirmMsgLock.RUnlock()
	v := append(proposalId, senderPartId...)
	v = append(v, receivePartId...)
	_, seen := s.seenConfirmMsgCache.Get(string(v))
	return seen
}

// Set proposalId in seen exit request cache.
func (s *Service) setConfirmMsgSeen(proposalId []byte, senderPartId []byte, receivePartId []byte) {
	s.seenConfirmMsgLock.Lock()
	defer s.seenConfirmMsgLock.Unlock()
	v := append(proposalId, senderPartId...)
	v = append(v, receivePartId...)
	s.seenConfirmMsgCache.Add(string(v), true)
}