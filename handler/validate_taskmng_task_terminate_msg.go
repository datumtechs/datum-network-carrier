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

func (s *Service) validateTaskTerminateMessagePubSub(ctx context.Context, pid peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
	// Validation runs on publish (not just subscriptions), so we should approve any message from ourselves.
	if pid == s.cfg.P2P.PeerID() {
		return pubsub.ValidationAccept
	}

	ctx, span := trace.StartSpan(ctx, "sync.TaskTerminateMsg")
	defer span.End()

	m, err := s.decodePubsubMessage(msg)
	if err != nil {
		log.WithError(err).Debug("Could not decode message")
		traceutil.AnnotateError(span, err)
		return pubsub.ValidationReject
	}

	message, ok := m.(*taskmngcpb.TaskTerminateMsg)
	if !ok {
		log.Errorf("Invalid message type in the validateTaskTerminateMessagePubSub, typ: %T", m)
		return pubsub.ValidationReject
	}

	if s.hasSeenTaskTerminateMsg(message.MsgOption.ProposalId) {
		return pubsub.ValidationIgnore
	}

	// validate TaskTerminateMsg
	if err := s.validateTaskTerminateMsg(pid, message); err != nil {
		log.WithError(err).Errorf("Failed to call `validateTaskTerminateMessagePubSub`, proposalId: {%s}", common.BytesToHash(message.MsgOption.ProposalId).String())
		return pubsub.ValidationIgnore
	}

	msg.ValidatorData = message // Used in downstream subscriber
	return pubsub.ValidationAccept
}

// Returns true if the node has already received a prepare message request for the validator with index `proposalId`.
func (s *Service) hasSeenTaskTerminateMsg(proposalId []byte) bool {
	s.seenTaskTerminateMsgLock.RLock()
	defer s.seenTaskTerminateMsgLock.RUnlock()
	_, seen := s.seenTaskTerminateMsgCache.Get(string(proposalId))
	return seen
}

// Set proposalId in seen exit request cache.
func (s *Service) setTaskTerminateMsgSeen(proposalId []byte) {
	s.seenTaskTerminateMsgLock.Lock()
	defer s.seenTaskTerminateMsgLock.Unlock()
	s.seenTaskTerminateMsgCache.Add(string(proposalId), true)
}