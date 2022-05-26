package handler

import (
	"context"
	pb "github.com/datumtechs/datum-network-carrier/lib/p2p/v1"
	libp2pcore "github.com/libp2p/go-libp2p-core"
	"github.com/pkg/errors"
	types "github.com/prysmaticlabs/eth2-types"
	"time"
)

// carrierBlocksByRangeRPCHandler looks up the request blocks from the database from a given start block.
func (s *Service) carrierBlocksByRangeRPCHandler(ctx context.Context, msg interface{}, stream libp2pcore.Stream) error {
	ctx, cancel := context.WithTimeout(ctx, respTimeout)
	defer cancel()
	SetRPCStreamDeadlines(stream)

	// Ticker to stagger out large requests.
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	m, ok := msg.(*pb.CarrierBlocksByRangeRequest)
	if !ok {
		return errors.New("message is not type *pb.CarrierBlockByRangeRequest")
	}
	if err := s.validateRangeRequest(m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		return err
	}
	//TODO: ....
	closeStream(stream, log)
	return nil
}

func (s *Service) writeBlockRangeToStream(ctx context.Context, startSlot, endSlot types.Slot,
	step uint64, stream libp2pcore.Stream) error {
	//TODO: ....
	return nil
}

func (s *Service) validateRangeRequest(r *pb.CarrierBlocksByRangeRequest) error {
	//TODO: ...
	return nil
}
