package handler

import (
	"context"
	carrierrpcdebugpbv1 "github.com/datumtechs/datum-network-carrier/pb/carrier/rpc/debug/v1"
	p2ptypes "github.com/datumtechs/datum-network-carrier/p2p/types"
	libp2pcore "github.com/libp2p/go-libp2p-core"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/shared/traceutil"
	"go.opencensus.io/trace"
	"math/rand"
	"time"
)

func (s *Service) gossipTestDataByRangeRPCHandler(ctx context.Context, msg interface{}, stream libp2pcore.Stream) error {
	log.WithField("peer", stream.Conn().RemotePeer()).Debug("Receive gossipTestData message")
	ctx, span := trace.StartSpan(ctx, "handler.GossipTestDataByRangeRPCHandler")
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, respTimeout)
	defer cancel()
	SetRPCStreamDeadlines(stream)

	// Ticker to stagger out large requests.
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	m, ok := msg.(*carrierrpcdebugpbv1.GossipTestData) // as request.
	if !ok {
		return errors.New("message is not type *p2ppb.GossipTestData")
	}
	if err := s.validateGossipRangeRequest(m); err != nil {
		s.writeErrorResponseToStream(responseCodeInvalidRequest, err.Error(), stream)
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		traceutil.AnnotateError(span, err)
		return err
	}

	err := s.writeGossipTestDataRangeToStream(ctx, stream)
	if err != nil && !errors.Is(err, p2ptypes.ErrInvalidParent) {
		return err
	}
	closeStream(stream, log)
	return nil
}

func (s *Service) writeGossipTestDataRangeToStream(ctx context.Context, stream libp2pcore.Stream) error {
	ctx, span := trace.StartSpan(ctx, "sync.WriteGossipTestDataRangeToStream")
	defer span.End()

	blks, err := generateTestData()
	if err != nil {
		log.WithError(err).Debug("Could not get gossip test data")
		s.writeErrorResponseToStream(responseCodeServerError, p2ptypes.ErrGeneric.Error(), stream)
		traceutil.AnnotateError(span, err)
		return err
	}
	for _, b := range blks {
		if b == nil || b.Data == nil {
			continue
		}
		if chunkErr := s.chunkWriter(stream, b); chunkErr != nil {
			log.WithError(chunkErr).Debug("Could not send a chunked response")
			s.writeErrorResponseToStream(responseCodeServerError, p2ptypes.ErrGeneric.Error(), stream)
			traceutil.AnnotateError(span, chunkErr)
			return chunkErr
		}

	}
	// Return error in the evengine we have an invalid parent.
	return err
}

func (s *Service) validateGossipRangeRequest(r *carrierrpcdebugpbv1.GossipTestData) error {
	count := r.Count
	step := r.Step
	if count < 10 {
		//return p2ptypes.ErrInvalidRequest
	}
	if step < 1 {
		//return p2ptypes.ErrInvalidRequest
	}
	return nil
}

func generateTestData() ([]*carrierrpcdebugpbv1.SignedGossipTestData, error) {
	return []*carrierrpcdebugpbv1.SignedGossipTestData{
		{
			Data:                 &carrierrpcdebugpbv1.GossipTestData{
				Data:                 []byte("data01"),
				Count:                uint64(rand.Int63n(100)),
				Step:                 uint64(rand.Int63n(100)),
			},
			Signature:            make([]byte, 48),
		},
		{
			Data:                 &carrierrpcdebugpbv1.GossipTestData{
				Data:                 []byte("data02"),
				Count:                uint64(rand.Int63n(100)),
				Step:                 uint64(rand.Int63n(100)),
			},
			Signature:            make([]byte, 48),
		},
		{
			Data:                 &carrierrpcdebugpbv1.GossipTestData{
				Data:                 []byte("data03"),
				Count:                uint64(rand.Int63n(100)),
				Step:                 uint64(rand.Int63n(100)),
			},
			Signature:            make([]byte, 48),
		},
	}, nil
}
