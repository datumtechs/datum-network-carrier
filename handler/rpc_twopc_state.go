package handler

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/types"
	libp2pcore "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/peer"
)

func (s *Service) prepareMsgRPCHandler(ctx context.Context, msg interface{}, stream libp2pcore.Stream) error {

	SetRPCStreamDeadlines(stream)

	m, ok := msg.(*pb.PrepareMsg)
	if !ok {
		//log.Errorf("Failed to convert `PrepareMsg` from msg, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return errors.New("message is not type *pb.PrepareMsg")
	}

	//TODO: validate request by rateLimiter.

	// validate prepareMsg
	if err := s.validatePrepareMsg(stream.Conn().RemotePeer(), m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		//log.WithError(err).Errorf("Failed to call `validatePrepareMsg`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}

	// handle prepareMsg
	if err := s.onPrepareMsg(stream.Conn().RemotePeer(), m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		//log.WithError(err).Errorf("Failed to call `onPrepareMsg`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}

	// response code
	if _, err := stream.Write([]byte{responseCodeSuccess}); err != nil {
		//log.WithError(err).Errorf("Could not write to stream for response, after to call `onPrepareMsg`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}

	closeStream(stream, log)
	return nil
}

func (s *Service) prepareVoteRPCHandler(ctx context.Context, msg interface{}, stream libp2pcore.Stream) error {

	SetRPCStreamDeadlines(stream)

	m, ok := msg.(*pb.PrepareVote)
	if !ok {
		//log.Errorf("Failed to convert `PrepareVote` from msg, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return errors.New("message is not type *pb.PrepareVote")
	}

	// validate prepareVote
	if err := s.validatePrepareVote(stream.Conn().RemotePeer(), m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		//log.WithError(err).Errorf("Failed to call `validatePrepareVote`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}

	// handle prepareVote
	if err := s.onPrepareVote(stream.Conn().RemotePeer(), m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		//log.WithError(err).Errorf("Failed to call `onPrepareVote`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}

	// response code
	if _, err := stream.Write([]byte{responseCodeSuccess}); err != nil {
		//log.WithError(err).Errorf("Could not write to stream for response, after to call `onPrepareVote`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}

	closeStream(stream, log)
	return nil
}


func (s *Service) confirmMsgRPCHandler(ctx context.Context, msg interface{}, stream libp2pcore.Stream) error {

	SetRPCStreamDeadlines(stream)

	m, ok := msg.(*pb.ConfirmMsg)
	if !ok {
		//log.Errorf("Failed to convert `ConfirmMsg` from msg, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return errors.New("message is not type *pb.ConfirmMsg")
	}

	// validate ConfirmMsg
	if err := s.validateConfirmMsg(stream.Conn().RemotePeer(), m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		//log.WithError(err).Errorf("Failed to call `validateConfirmMsg`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}

	// handle ConfirmMsg
	if err := s.onConfirmMsg(stream.Conn().RemotePeer(), m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		//log.WithError(err).Errorf("Failed to call `onConfirmMsg`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}

	// response code
	if _, err := stream.Write([]byte{responseCodeSuccess}); err != nil {
		//log.WithError(err).Errorf("Could not write to stream for response, after to call `onConfirmMsg`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}

	closeStream(stream, log)
	return nil
}


func (s *Service) confirmVoteRPCHandler(ctx context.Context, msg interface{}, stream libp2pcore.Stream) error {

	SetRPCStreamDeadlines(stream)

	m, ok := msg.(*pb.ConfirmVote)
	if !ok {
		//log.Errorf("Failed to convert `ConfirmVote` from msg, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return errors.New("message is not type *pb.ConfirmVote")
	}

	// validate ConfirmVote
	if err := s.validateConfirmVote(stream.Conn().RemotePeer(), m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		//log.WithError(err).Errorf("Failed to call `validateConfirmVote`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}

	// handle ConfirmVote
	if err := s.onConfirmVote(stream.Conn().RemotePeer(), m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		//log.WithError(err).Errorf("Failed to call `onConfirmVote`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}

	// response code
	if _, err := stream.Write([]byte{responseCodeSuccess}); err != nil {
		//log.WithError(err).Errorf("Could not write to stream for response, after to call `onConfirmVote`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}

	closeStream(stream, log)
	return nil
}

func (s *Service) commitMsgRPCHandler(ctx context.Context, msg interface{}, stream libp2pcore.Stream) error {

	SetRPCStreamDeadlines(stream)

	m, ok := msg.(*pb.CommitMsg)
	if !ok {
		//log.Errorf("Failed to convert `CommitMsg` from msg, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return errors.New("message is not type *pb.CommitMsg")
	}

	// validate CommitMsg
	if err := s.validateCommitMsg(stream.Conn().RemotePeer(), m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		//log.WithError(err).Errorf("Failed to call `validateCommitMsg`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}

	// handle CommitMsg
	if err := s.onCommitMsg(stream.Conn().RemotePeer(), m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		//log.WithError(err).Errorf("Failed to call `onCommitMsg`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}

	// response code
	if _, err := stream.Write([]byte{responseCodeSuccess}); err != nil {
		//log.WithError(err).Errorf("Could not write to stream for response, after to call `onCommitMsg`, proposalId: {%s}", common.BytesToHash(m.ProposalId).String())
		return err
	}

	closeStream(stream, log)
	return nil
}

func (s *Service) taskResultMsgRPCHandler(ctx context.Context, msg interface{}, stream libp2pcore.Stream) error {

	SetRPCStreamDeadlines(stream)

	m, ok := msg.(*pb.TaskResultMsg)
	if !ok {
		//log.Errorf("Failed to convert `TaskResultMsg` from msg, proposalId: {%s}, taskId: {%s}", common.BytesToHash(m.ProposalId).String(), string(m.GetTaskId))
		return errors.New("message is not type *pb.TaskResultMsg")
	}

	// validate TaskResultMsg
	if err := s.validateTaskResultMsg(stream.Conn().RemotePeer(), m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		//log.WithError(err).Errorf("Failed to call `validateTaskResultMsg`, proposalId: {%s}, taskId: {%s}", common.BytesToHash(m.ProposalId).String(), string(m.GetTaskId))
		return err
	}

	// handle TaskResultMsg
	if err := s.onTaskResultMsg(stream.Conn().RemotePeer(), m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		//log.WithError(err).Warnf("Warning to call `onTaskResultMsg`, proposalId: {%s}, taskId: {%s}", common.BytesToHash(m.ProposalId).String(), string(m.GetTaskId))
		return err
	}

	// response code
	if _, err := stream.Write([]byte{responseCodeSuccess}); err != nil {
		//log.WithError(err).Errorf("Could not write to stream for response, after to call `onTaskResultMsg`, proposalId: {%s}, taskId: {%s}", common.BytesToHash(m.ProposalId).String(), string(m.GetTaskId))
		return err
	}

	closeStream(stream, log)
	return nil
}


// ------------------------------------  some validate Fn  ------------------------------------

func (s *Service) validatePrepareMsg(pid peer.ID, r *pb.PrepareMsg) error {
	engine, ok := s.cfg.Engines[types.TwopcTyp]
	if !ok {
		return fmt.Errorf("Failed to fecth 2pc engine instanse ...")
	}
	return engine.ValidateConsensusMsg(pid, &types.PrepareMsgWrap{PrepareMsg: r})
}

func (s *Service) validatePrepareVote(pid peer.ID, r *pb.PrepareVote) error {
	engine, ok := s.cfg.Engines[types.TwopcTyp]
	if !ok {
		return fmt.Errorf("Failed to fecth 2pc engine instanse ...")
	}
	return engine.ValidateConsensusMsg(pid, &types.PrepareVoteWrap{PrepareVote: r})
}

func (s *Service) validateConfirmMsg(pid peer.ID, r *pb.ConfirmMsg) error {
	engine, ok := s.cfg.Engines[types.TwopcTyp]
	if !ok {
		return fmt.Errorf("Failed to fecth 2pc engine instanse ...")
	}
	return engine.ValidateConsensusMsg(pid, &types.ConfirmMsgWrap{ConfirmMsg: r})
}

func (s *Service) validateConfirmVote(pid peer.ID, r *pb.ConfirmVote) error {
	engine, ok := s.cfg.Engines[types.TwopcTyp]
	if !ok {
		return fmt.Errorf("Failed to fecth 2pc engine instanse ...")
	}
	return engine.ValidateConsensusMsg(pid, &types.ConfirmVoteWrap{ConfirmVote: r})
}

func (s *Service) validateCommitMsg(pid peer.ID, r *pb.CommitMsg) error {
	engine, ok := s.cfg.Engines[types.TwopcTyp]
	if !ok {
		return fmt.Errorf("Failed to fecth 2pc engine instanse ...")
	}
	return engine.ValidateConsensusMsg(pid, &types.CommitMsgWrap{CommitMsg: r})
}



// ------------------------------------  some handle Fn  ------------------------------------

func (s *Service) onPrepareMsg(pid peer.ID, r *pb.PrepareMsg) error {
	engine, ok := s.cfg.Engines[types.TwopcTyp]
	if !ok {
		return fmt.Errorf("Failed to fecth 2pc engine instanse ...")
	}
	return engine.OnConsensusMsg(pid, &types.PrepareMsgWrap{PrepareMsg: r})
}

func (s *Service) onPrepareVote(pid peer.ID, r *pb.PrepareVote) error {
	engine, ok := s.cfg.Engines[types.TwopcTyp]
	if !ok {
		return fmt.Errorf("Failed to fecth 2pc engine instanse ...")
	}
	return engine.OnConsensusMsg(pid, &types.PrepareVoteWrap{PrepareVote: r})
}

func (s *Service) onConfirmMsg(pid peer.ID, r *pb.ConfirmMsg) error {
	engine, ok := s.cfg.Engines[types.TwopcTyp]
	if !ok {
		return fmt.Errorf("Failed to fecth 2pc engine instanse ...")
	}
	return engine.OnConsensusMsg(pid, &types.ConfirmMsgWrap{ConfirmMsg: r})
}

func (s *Service) onConfirmVote(pid peer.ID, r *pb.ConfirmVote) error {
	engine, ok := s.cfg.Engines[types.TwopcTyp]
	if !ok {
		return fmt.Errorf("Failed to fecth 2pc engine instanse ...")
	}
	return engine.OnConsensusMsg(pid, &types.ConfirmVoteWrap{ConfirmVote: r})
}

func (s *Service) onCommitMsg(pid peer.ID, r *pb.CommitMsg) error {
	engine, ok := s.cfg.Engines[types.TwopcTyp]
	if !ok {
		return fmt.Errorf("Failed to fecth 2pc engine instanse ...")
	}
	return engine.OnConsensusMsg(pid, &types.CommitMsgWrap{CommitMsg: r})
}

