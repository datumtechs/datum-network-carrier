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
// within process_prepare_message before forwarding it across the network.
func (s *Service) validatePrepareMessagePubSub(ctx context.Context, pid peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
	// Validation runs on publish (not just subscriptions), so we should approve any message from ourselves.
	if pid == s.cfg.P2P.PeerID() {
		return pubsub.ValidationAccept
	}

	// The head state will be too far away to validate any voluntary exit.
	/*if s.cfg.InitialSync.Syncing() {
		return pubsub.ValidationIgnore
	}*/

	ctx, span := trace.StartSpanWithRemoteParent(ctx, "handler.validatePrepareMessagePubSub",
		traceutil.GenerateParentSpanWithPrepareMsg(pid, msg))
	defer span.End()

	m, err := s.decodePubsubMessage(msg)
	if err != nil {
		log.WithError(err).Debug("Could not decode message")
		traceutil.AnnotateError(span, err)
		return pubsub.ValidationReject
	}

	message, ok := m.(*twopcpb.PrepareMsg)
	if !ok {
		log.Errorf("Invalid message type in the validatePrepareMessagePubSub, typ: %T", m)
		return pubsub.ValidationReject
	}

	if s.hasSeenPrepareMsg(message.MsgOption.ProposalId, message.MsgOption.SenderPartyId, message.MsgOption.ReceiverPartyId) {
		return pubsub.ValidationIgnore
	}

	//// validate prepareMsg
	//if err := s.validatePrepareMsg(pid, message); err != nil {
	//	log.WithError(err).Errorf("Failed to call `validatePrepareMsg`, proposalId: {%s}", common.BytesToHash(message.MsgOption.ProposalId).String())
	//	return pubsub.ValidationIgnore
	//}

	msg.ValidatorData = message // Used in downstream subscriber

	return pubsub.ValidationAccept
}

// Returns true if the node has already received a prepare message request for the validator with index `proposalId`.
func (s *Service) hasSeenPrepareMsg(proposalId []byte, taskPartyId []byte, receivePartId []byte) bool {
	s.seenPrepareMsgLock.RLock()
	defer s.seenPrepareMsgLock.RUnlock()
	v := append(proposalId, taskPartyId...)
	v = append(v, receivePartId...)
	_, seen := s.seenPrepareMsgCache.Get(string(v))
	return seen
}

// Set proposalId in seen exit request cache.
func (s *Service) setPrepareMsgSeen(proposalId []byte, taskPartyId []byte, receivePartId []byte) {
	s.seenPrepareMsgLock.Lock()
	defer s.seenPrepareMsgLock.Unlock()
	v := append(proposalId, taskPartyId...)
	v = append(v, receivePartId...)
	s.seenPrepareMsgCache.Add(string(v), true)
}