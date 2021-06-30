package handler

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/types"
	libp2pcore "github.com/libp2p/go-libp2p-core"
	"time"
)

// carrierBlocksByRangeRPCHandler looks up the request blocks from the database from a given start block.
func (s *Service) sendPrepareMsgRPCHandler(ctx context.Context, msg interface{}, stream libp2pcore.Stream) error {
	ctx, cancel := context.WithTimeout(ctx, respTimeout)
	defer cancel()
	SetRPCStreamDeadlines(stream)

	// Ticker to stagger out large requests.
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	m, ok := msg.(*pb.PrepareMsg)
	if !ok {
		return errors.New("message is not type *pb.PrepareMsg")
	}
	// validate prepareMsg
	if err := s.validatePrepareMsg(m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		return err
	}
	// hanlde prepareMsg
	if err := s.onPrepareMsg(m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		return err
	}

	closeStream(stream, log)
	return nil
}

func (s *Service) validatePrepareMsg(r *pb.PrepareMsg) error {
	engine, ok := s.cfg.Engines[types.TwopcTyp]
	if !ok {
		return fmt.Errorf("Failed to fecth 2pc engine instanse ...")
	}
	return engine.ValidateConsensusMsg(&types.PrepareMsgWrap{PrepareMsg: r})
}

func (s *Service) onPrepareMsg(r *pb.PrepareMsg) error {
	engine, ok := s.cfg.Engines[types.TwopcTyp]
	if !ok {
		return fmt.Errorf("Failed to fecth 2pc engine instanse ...")
	}
	return engine.OnConsensusMsg(&types.PrepareMsgWrap{PrepareMsg: r})
}