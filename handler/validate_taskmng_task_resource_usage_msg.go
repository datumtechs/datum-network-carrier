package handler

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/traceutil"
	taskmngcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/taskmng"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"go.opencensus.io/trace"
)

func (s *Service) validateTaskResourceUsageMessagePubSub(ctx context.Context, pid peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
	// Validation runs on publish (not just subscriptions), so we should approve any message from ourselves.
	if pid == s.cfg.P2P.PeerID() {
		return pubsub.ValidationAccept
	}

	ctx, span := trace.StartSpanWithRemoteParent(ctx, "handler.validateTaskResourceUsageMessagePubSub",
		traceutil.GenerateParentSpanWithTaskResourceUsageMessage(pid, msg))
	defer span.End()

	m, err := s.decodePubsubMessage(msg)
	if err != nil {
		log.WithError(err).Debug("Could not decode message for TaskResourceUsageMsg")
		traceutil.AnnotateError(span, err)
		return pubsub.ValidationReject
	}

	message, ok := m.(*taskmngcpb.TaskResourceUsageMsg)
	if !ok {
		log.Errorf("Invalid message type in the validateTaskResourceUsageMessagePubSub, typ: %T", m)
		return pubsub.ValidationReject
	}

	if s.hasSeenTaskResourceUsageMsg(message.MsgOption.ProposalId, message.MsgOption.SenderPartyId, message.MsgOption.ReceiverPartyId) {
		return pubsub.ValidationIgnore
	}

	// validate TaskResourceUsageMsg
	if err := s.validateTaskResourceUsageMsg(pid, message); err != nil {
		log.WithError(err).Errorf("Failed to call `validateTaskResourceUsageMsg`, proposalId: {%s}", common.BytesToHash(message.MsgOption.ProposalId).String())
		return pubsub.ValidationIgnore
	}

	msg.ValidatorData = message // Used in downstream subscriber
	return pubsub.ValidationAccept
}

// Returns true if the node has already received a prepare message request for the validator with index `proposalId`.
func (s *Service) hasSeenTaskResourceUsageMsg(proposalId []byte, senderPartId []byte, receivePartId []byte) bool {
	s.seenTaskResourceUsageMsgLock.RLock()
	defer s.seenTaskResourceUsageMsgLock.RUnlock()
	v := append(proposalId, senderPartId...)
	v = append(v, receivePartId...)
	_, seen := s.seenTaskResourceUsageMsgCache.Get(string(v))
	return seen
}

// Set proposalId in seen exit request cache.
func (s *Service) setTaskResourceUsageMsgSeen(proposalId []byte, senderPartId []byte, receivePartId []byte) {
	s.seenTaskResourceUsageMsgLock.Lock()
	defer s.seenTaskResourceUsageMsgLock.Unlock()
	v := append(proposalId, senderPartId...)
	v = append(v, receivePartId...)
	s.seenTaskResourceUsageMsgCache.Add(string(v), true)
}