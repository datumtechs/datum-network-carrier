package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	rpcpb "github.com/RosettaFlow/Carrier-Go/lib/rpc/v1"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"
)

func (s *Service) gossipTestDataSubscriber(ctx context.Context, pid peer.ID, msg proto.Message) error {
	ve, ok := msg.(*rpcpb.SignedGossipTestData)
	if !ok {
		return fmt.Errorf("wrong type, expected: *ethpb.SignedVoluntaryExit got: %T", msg)
	}

	if ve.Data == nil {
		return errors.New("data can't be nil")
	}
	log.WithField("peer", pid).Debug("Receive gossip message")
	return nil
}

func (s *Service) prepareMessageSubscriber(ctx context.Context, pid peer.ID, msg proto.Message) error {
	message, ok := msg.(*pb.PrepareMsg)
	if !ok {
		return fmt.Errorf("wrong type, expected: *pb.PrepareMsg got: %T", msg)
	}

	s.setPrepareMsgSeen(message.ProposalId, message.TaskPartyId)

	// handle prepareMsg
	if err := s.onPrepareMsg(pid, message); err != nil {
		log.WithError(err).Errorf("Failed to call `onPrepareMsg`, proposalId: {%s}", common.BytesToHash(message.ProposalId).String())
		return err
	}
	return nil
}

func (s *Service) prepareVoteSubscriber(ctx context.Context, pid peer.ID, msg proto.Message) error {
	m, ok := msg.(*pb.PrepareVote)
	if !ok {
		return fmt.Errorf("wrong type, expected: *pb.PrepareVote got: %T", msg)
	}

	s.setPrepareVoteSeen(m.ProposalId)

	// handle prepareVote
	if err := s.onPrepareVote(pid, m); err != nil {
		log.WithError(err).Errorf("Failed to call `onPrepareVote`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}
	return nil
}

func (s *Service) confirmMessageSubscriber(ctx context.Context, pid peer.ID, msg proto.Message) error {
	m, ok := msg.(*pb.ConfirmMsg)
	if !ok {
		return fmt.Errorf("wrong type, expected: *pb.ConfirmMsg got: %T", msg)
	}

	s.setCommitMsgSeen(m.ProposalId, m.TaskPartyId)

	// handle ConfirmMsg
	if err := s.onConfirmMsg(pid, m); err != nil {
		log.WithError(err).Errorf("Failed to call `onConfirmMsg`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}
	return nil
}

func (s *Service) confirmVoteSubscriber(ctx context.Context, pid peer.ID, msg proto.Message) error {
	m, ok := msg.(*pb.ConfirmVote)
	if !ok {
		return fmt.Errorf("wrong type, expected: *pb.ConfirmVote got: %T", msg)
	}

	s.setConfirmVoteSeen(m.ProposalId)

	// handle ConfirmVote
	if err := s.onConfirmVote(pid, m); err != nil {
		log.WithError(err).Errorf("Failed to call `onConfirmVote`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}
	return nil
}

func (s *Service) commitMessageSubscriber(ctx context.Context, pid peer.ID, msg proto.Message) error {
	m, ok := msg.(*pb.CommitMsg)
	if !ok {
		return fmt.Errorf("wrong type, expected: *pb.CommitMsg got: %T", msg)
	}

	s.setCommitMsgSeen(m.ProposalId, m.TaskPartyId)

	// handle CommitMsg
	if err := s.onCommitMsg(pid, m); err != nil {
		log.WithError(err).Errorf("Failed to call `onCommitMsg`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}
	return nil
}

func (s *Service) taskResultMessageSubscriber(ctx context.Context, pid peer.ID, msg proto.Message) error {
	m, ok := msg.(*pb.TaskResultMsg)
	if !ok {
		return fmt.Errorf("wrong type, expected: *pb.TaskResultMsg got: %T", msg)
	}

	s.setTaskResultMsgSeen(m.ProposalId)

	// handle TaskResultMsg
	if err := s.onTaskResultMsg(pid, m); err != nil {
		log.WithError(err).Warnf("Warning to call `onTaskResultMsg`, proposalId: {%s}, taskId: {%s}", common.BytesToHash(m.ProposalId).String(), string(m.TaskId))
		return err
	}
	return nil
}