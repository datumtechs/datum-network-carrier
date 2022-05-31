package handler

import (
	"context"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/common/hashutil"
	"github.com/datumtechs/datum-network-carrier/common/traceutil"
	carriernetmsgtaskmngpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/taskmng"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"go.opencensus.io/trace"
)

func (s *Service) validateTaskTerminateMessagePubSub(ctx context.Context, pid peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
	// Validation runs on publish (not just subscriptions), so we should approve any message from ourselves.
	if pid == s.cfg.P2P.PeerID() {
		return pubsub.ValidationAccept
	}

	ctx, span := trace.StartSpanWithRemoteParent(ctx, "handler.validateTaskTerminateMessagePubSub",
		traceutil.GenerateParentSpanWithTaskTerminateMessage(pid, msg))
	defer span.End()

	m, err := s.decodePubsubMessage(msg)
	if err != nil {
		log.WithError(err).Debug("Could not decode message")
		traceutil.AnnotateError(span, err)
		return pubsub.ValidationReject
	}

	message, ok := m.(*carriernetmsgtaskmngpb.TaskTerminateMsg)
	if !ok {
		log.Errorf("Invalid message type in the validateTaskTerminateMessagePubSub, typ: %T", m)
		return pubsub.ValidationReject
	}

	if s.hasSeenTaskTerminateMsg(msg) {
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
func (s *Service) hasSeenTaskTerminateMsg(msg proto.Message) bool {
	s.seenTaskTerminateMsgLock.RLock()
	defer s.seenTaskTerminateMsgLock.RUnlock()
	v := hashutil.Hash([]byte(msg.String()))
	_, seen := s.seenTaskTerminateMsgCache.Get(common.Bytes2Hex(v[0:]))
	return seen
}

// Set proposalId in seen exit request cache.
func (s *Service) setTaskTerminateMsgSeen(msg proto.Message) {
	s.seenTaskTerminateMsgLock.Lock()
	defer s.seenTaskTerminateMsgLock.Unlock()
	v := hashutil.Hash([]byte(msg.String()))
	s.seenTaskTerminateMsgCache.Add(common.Bytes2Hex(v[0:]), true)
}