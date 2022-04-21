package handler

import (
	"context"
	"github.com/Metisnetwork/Metis-Carrier/common"
	"github.com/Metisnetwork/Metis-Carrier/common/traceutil"
	taskmngcpb "github.com/Metisnetwork/Metis-Carrier/lib/netmsg/taskmng"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/prysmaticlabs/prysm/shared/hashutil"
	"go.opencensus.io/trace"
)

// Clients who receive a prepare message on this topic MUST validate the conditions
// within process_task_result_message before forwarding it across the network.
func (s *Service) validateTaskResultMessagePubSub(ctx context.Context, pid peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
	// Validation runs on publish (not just subscriptions), so we should approve any message from ourselves.
	if pid == s.cfg.P2P.PeerID() {
		return pubsub.ValidationAccept
	}

	ctx, span := trace.StartSpanWithRemoteParent(ctx, "handler.validateTaskResultMessagePubSub",
		traceutil.GenerateParentSpanWithTaskResultMessage(pid, msg))
	defer span.End()

	m, err := s.decodePubsubMessage(msg)
	if err != nil {
		log.WithError(err).Debug("Could not decode message")
		traceutil.AnnotateError(span, err)
		return pubsub.ValidationReject
	}

	message, ok := m.(*taskmngcpb.TaskResultMsg)
	if !ok {
		log.Errorf("Invalid message type in the validateTaskResultMessagePubSub, typ: %T", m)
		return pubsub.ValidationReject
	}

	if s.hasSeenTaskResultMsg(msg) {
		return pubsub.ValidationIgnore
	}

	// validate TaskResultMsg
	if err := s.validateTaskResultMsg(pid, message); err != nil {
		log.WithError(err).Errorf("Failed to call `validateTaskResultMsg`, proposalId: {%s}", common.BytesToHash(message.MsgOption.ProposalId).String())
		return pubsub.ValidationIgnore
	}

	msg.ValidatorData = message // Used in downstream subscriber
	return pubsub.ValidationAccept
}

// Returns true if the node has already received a prepare message request for the validator with index `proposalId`.
func (s *Service) hasSeenTaskResultMsg(msg proto.Message) bool {
	s.seenTaskResultMsgLock.RLock()
	defer s.seenTaskResultMsgLock.RUnlock()
	v := hashutil.Hash([]byte(msg.String()))
	_, seen := s.seenTaskResultMsgCache.Get(common.Bytes2Hex(v[0:]))
	return seen
}

// Set proposalId in seen exit request cache.
func (s *Service) setTaskResultMsgSeen(msg proto.Message) {
	s.seenTaskResultMsgLock.Lock()
	defer s.seenTaskResultMsgLock.Unlock()
	v := hashutil.Hash([]byte(msg.String()))
	s.seenTaskResultMsgCache.Add(common.Bytes2Hex(v[0:]), true)
}