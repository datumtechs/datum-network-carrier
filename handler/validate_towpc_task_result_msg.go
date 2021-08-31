package handler

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/traceutil"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"go.opencensus.io/trace"
)

// Clients who receive a prepare message on this topic MUST validate the conditions
// within process_task_result_message before forwarding it across the network.
func (s *Service) validateTaskResultMessagePubSub(ctx context.Context, pid peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
	// Validation runs on publish (not just subscriptions), so we should approve any message from ourselves.
	if pid == s.cfg.P2P.PeerID() {
		return pubsub.ValidationAccept
	}

	ctx, span := trace.StartSpan(ctx, "sync.TaskResultMsg")
	defer span.End()

	m, err := s.decodePubsubMessage(msg)
	if err != nil {
		log.WithError(err).Debug("Could not decode message")
		traceutil.AnnotateError(span, err)
		return pubsub.ValidationReject
	}

	message, ok := m.(*pb.TaskResultMsg)
	if !ok {
		log.Errorf("Invalid message type in the validateTaskResultMessagePubSub, typ: %T", m)
		return pubsub.ValidationReject
	}

	if s.hasSeenTaskResultMsg(message.ProposalId) {
		return pubsub.ValidationIgnore
	}

	// validate TaskResultMsg
	if err := s.validateTaskResultMsg(pid, message); err != nil {
		log.WithError(err).Errorf("Failed to call `validateTaskResultMsg`, proposalId: {%s}, taskId: {%s}", common.BytesToHash(message.ProposalId).String(), string(message.TaskId))
		return pubsub.ValidationIgnore
	}

	msg.ValidatorData = message // Used in downstream subscriber
	return pubsub.ValidationAccept
}

// Returns true if the node has already received a prepare message request for the validator with index `proposalId`.
func (s *Service) hasSeenTaskResultMsg(proposalId []byte) bool {
	s.seenTaskResultMsgLock.RLock()
	defer s.seenTaskResultMsgLock.RUnlock()
	_, seen := s.seenTaskResultMsgCache.Get(string(proposalId))
	return seen
}

// Set proposalId in seen exit request cache.
func (s *Service) setTaskResultMsgSeen(proposalId []byte) {
	s.seenTaskResultMsgLock.Lock()
	defer s.seenTaskResultMsgLock.Unlock()
	s.seenTaskResultMsgCache.Add(string(proposalId), true)
}