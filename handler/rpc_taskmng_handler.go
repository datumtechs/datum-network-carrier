package handler

import (
	"context"
	"errors"
	taskmngpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/taskmng"
	libp2pcore "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/peer"
)

// Deprecate: See taskResultMessageSubscriber in the subscriber_handlers.go.
func (s *Service) taskResultMsgRPCHandler(ctx context.Context, msg interface{}, stream libp2pcore.Stream) error {

	SetRPCStreamDeadlines(stream)

	m, ok := msg.(*taskmngpb.TaskResultMsg)
	if !ok {
		return errors.New("message is not type *taskmngpb.TaskResultMsg")
	}

	// validate TaskResultMsg
	if err := s.validateTaskResultMsg(stream.Conn().RemotePeer(), m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		return err
	}

	// handle TaskResultMsg
	if err := s.onTaskResultMsg(stream.Conn().RemotePeer(), m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		return err
	}

	// response code
	if _, err := stream.Write([]byte{responseCodeSuccess}); err != nil {
		return err
	}

	closeStream(stream, log)
	return nil
}

// Deprecate: See taskResourceUsageMessageSubscriber in the subscriber_handlers.go.
func (s *Service) taskResourceUsageMsgRPCHandler(ctx context.Context, msg interface{}, stream libp2pcore.Stream) error {

	SetRPCStreamDeadlines(stream)

	m, ok := msg.(*taskmngpb.TaskResourceUsageMsg)
	if !ok {
		return errors.New("message is not type *taskmngpb.TaskResourceUsageMsg")
	}

	// validate TaskResourceUsageMsg
	if err := s.validateTaskResourceUsageMsg(stream.Conn().RemotePeer(), m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		return err
	}

	// handle TaskResourceUsageMsg
	if err := s.onTaskResourceUsageMsg(stream.Conn().RemotePeer(), m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
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

// Deprecate: See taskTerminateMessageSubscriber in the subscriber_handlers.go.
func (s *Service) taskTerminateMsgRPCHandler(ctx context.Context, msg interface{}, stream libp2pcore.Stream) error {

	SetRPCStreamDeadlines(stream)

	m, ok := msg.(*taskmngpb.TaskTerminateMsg)
	if !ok {
		return errors.New("message is not type *taskmngpb.TaskTerminateMsg")
	}

	// validate TaskTerminateMsg
	if err := s.validateTaskTerminateMsg(stream.Conn().RemotePeer(), m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		return err
	}

	// handle TaskTerminateMsg
	if err := s.onTaskTerminateMsg(stream.Conn().RemotePeer(), m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
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


// --------------------------------------- validate fn ---------------------------------------
func (s *Service) validateTaskResultMsg(pid peer.ID, r *taskmngpb.TaskResultMsg) error {
	return s.cfg.TaskManager.ValidateTaskResultMsg(pid, r)
}

func (s *Service) validateTaskResourceUsageMsg(pid peer.ID, r *taskmngpb.TaskResourceUsageMsg) error {
	return s.cfg.TaskManager.ValidateTaskResourceUsageMsg(pid, r)
}

func (s *Service) validateTaskTerminateMsg(pid peer.ID, r *taskmngpb.TaskTerminateMsg) error {
	return s.cfg.TaskManager.ValidateTaskTerminateMsg(pid, r)
}




// --------------------------------------- handler fn ---------------------------------------
func (s *Service) onTaskResultMsg(pid peer.ID, r *taskmngpb.TaskResultMsg) error {
	return s.cfg.TaskManager.OnTaskResultMsg(pid, r)
}

func (s *Service) onTaskResourceUsageMsg(pid peer.ID, r *taskmngpb.TaskResourceUsageMsg) error {
	return s.cfg.TaskManager.OnTaskResourceUsageMsg(pid, r)
}

func (s *Service) onTaskTerminateMsg(pid peer.ID, r *taskmngpb.TaskTerminateMsg) error {
	return s.cfg.TaskManager.OnTaskTerminateMsg(pid, r)
}
