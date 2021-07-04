package handler

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/common/runutil"
	timeutils "github.com/RosettaFlow/Carrier-Go/common/timeutil"
	pb "github.com/RosettaFlow/Carrier-Go/lib/p2p/v1"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/p2p/peers"
	p2ptypes "github.com/RosettaFlow/Carrier-Go/p2p/types"
	libp2pcore "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

// maintainPeerStatuses by infrequently polling peers for their latest status.
func (s *Service) maintainPeerStatuses() {
	//TODO: need to ensure the interval....
	interval := 2 * time.Minute
	runutil.RunEvery(s.ctx, interval, func() {
		wg := new(sync.WaitGroup)
		for _, pid := range s.cfg.P2P.Peers().Connected() {
			wg.Add(1)
			go func(id peer.ID) {
				defer wg.Done()
				// If our peer status has not been updated correctly we disconnect over here
				// and set the connection state over here instead.
				if s.cfg.P2P.Host().Network().Connectedness(id) != network.Connected {
					s.cfg.P2P.Peers().SetConnectionState(id, peers.PeerDisconnecting)
					if err := s.cfg.P2P.Disconnect(id); err != nil {
						log.Debugf("Error when disconnecting with peer: %v", err)
					}
					s.cfg.P2P.Peers().SetConnectionState(id, peers.PeerDisconnected)
					return
				}
				// Disconnect from peers that are considered bad by any of the registered scorers.
				if s.cfg.P2P.Peers().IsBad(id) {
					s.disconnectBadPeer(s.ctx, id)
					return
				}
				// If the status hasn't been updated in the recent interval time.
				lastUpdated, err := s.cfg.P2P.Peers().ChainStateLastUpdated(id)
				if err != nil {
					// Peer has vanished; nothing to do.
					return
				}
				if timeutils.Now().After(lastUpdated.Add(interval)) {
					if err := s.reValidatePeer(s.ctx, id); err != nil {
						log.WithField("peer", id).WithError(err).Debug("Could not revalidate peer")
						s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(id)
					}
				}
			}(pid)
		}
		// Wait for all status checks to finish and then proceed onwards to
		// pruning excess peers.
		wg.Wait()
		peerIds := s.cfg.P2P.Peers().PeersToPrune()
		//peerIds = s.filterNeededPeers(peerIds)
		for _, id := range peerIds {
			if err := s.sendGoodByeAndDisconnect(s.ctx, p2ptypes.GoodbyeCodeTooManyPeers, id); err != nil {
				log.WithField("peer", id).WithError(err).Debug("Could not disconnect with peer")
			}
		}
	})
}

// resyncIfBehind checks periodically to see if we are in normal sync but have fallen behind our peers
// by more than an epoch, in which case we attempt a resync using the initial sync method to catch up.
func (s *Service) resyncIfBehind() {
}

// shouldReSync returns true if the node is not syncing and falls behind two epochs.
func (s *Service) shouldReSync() bool {
	return false
}

// sendRPCStatusRequest for a given topic with an expected protobuf message type.
func (s *Service) sendRPCStatusRequest(ctx context.Context, id peer.ID) error {
	ctx, cancel := context.WithTimeout(ctx, respTimeout)
	defer cancel()

	forkDigest, err := s.forkDigest()
	if err != nil {
		return err
	}
	//TODO: need to change....
	resp := &pb.Status{
		ForkDigest:     forkDigest[:],
		FinalizedRoot:  make([]byte, 32),
		FinalizedEpoch: 0,
		HeadRoot:       make([]byte, 32),
		HeadSlot:       0,
	}
	stream, err := s.cfg.P2P.Send(ctx, resp, p2p.RPCStatusTopic, id)
	if err != nil {
		return err
	}
	defer closeStream(stream, log)

	code, errMsg, err := ReadStatusCode(stream, s.cfg.P2P.Encoding())
	if err != nil {
		return err
	}

	if code != 0 {
		s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(id)
		return errors.New(errMsg)
	}

	msg := &pb.Status{}
	if err := s.cfg.P2P.Encoding().DecodeWithMaxLength(stream, msg); err != nil {
		return err
	}

	// If validation fails, validation error is logged, and peer status scorer will mark peer as bad.
	err = s.validateStatusMessage(ctx, msg)
	s.cfg.P2P.Peers().Scorers().PeerStatusScorer().SetPeerStatus(id, msg, err)
	if s.cfg.P2P.Peers().IsBad(id) {
		s.disconnectBadPeer(s.ctx, id)
	}
	return err
}

func (s *Service) reValidatePeer(ctx context.Context, id peer.ID) error {
	//s.cfg.P2P.Peers().Scorers().PeerStatusScorer().SetHeadSlot(s.cfg.Chain.HeadSlot())
	if err := s.sendRPCStatusRequest(ctx, id); err != nil {
		return err
	}
	// Do not return an error for ping requests.
	if err := s.sendPingRequest(ctx, id); err != nil {
		log.WithError(err).Debug("Could not ping peer")
	}
	return nil
}

// statusRPCHandler reads the incoming Status RPC from the peer and responds with our version of a status message.
// This handler will disconnect any peer that does not match our fork version.
func (s *Service) statusRPCHandler(ctx context.Context, msg interface{}, stream libp2pcore.Stream) error {
	ctx, cancel := context.WithTimeout(ctx, ttfbTimeout)
	defer cancel()
	SetRPCStreamDeadlines(stream)
	log := log.WithField("handler", "status")
	m, ok := msg.(*pb.Status)
	if !ok {
		return errors.New("message is not type *pb.Status")
	}
	if err := s.rateLimiter.validateRequest(stream, 1); err != nil {
		return err
	}
	s.rateLimiter.add(stream, 1)

	remotePeer := stream.Conn().RemotePeer()
	if err := s.validateStatusMessage(ctx, m); err != nil {
		log.WithFields(logrus.Fields{
			"peer":  remotePeer,
			"error": err,
		}).Debug("Invalid status message from peer")

		respCode := byte(0)
		switch err {
		case p2ptypes.ErrGeneric:
			respCode = responseCodeServerError
		case p2ptypes.ErrWrongForkDigestVersion:
			// Respond with our status and disconnect with the peer.
			s.cfg.P2P.Peers().SetChainState(remotePeer, m)
			if err := s.respondWithStatus(ctx, stream); err != nil {
				return err
			}
			// Close before disconnecting, and wait for the other end to ack our response.
			closeStreamAndWait(stream, log)
			if err := s.sendGoodByeAndDisconnect(ctx, p2ptypes.GoodbyeCodeWrongNetwork, remotePeer); err != nil {
				return err
			}
			return nil
		default:
			respCode = responseCodeInvalidRequest
			s.cfg.P2P.Peers().Scorers().BadResponsesScorer().Increment(remotePeer)
		}

		originalErr := err
		resp, err := s.generateErrorResponse(respCode, err.Error())
		if err != nil {
			log.WithError(err).Debug("Could not generate a response error")
		} else if _, err := stream.Write(resp); err != nil {
			// The peer may already be ignoring us, as we disagree on fork version, so log this as debug only.
			log.WithError(err).Debug("Could not write to stream")
		}
		closeStreamAndWait(stream, log)
		if err := s.sendGoodByeAndDisconnect(ctx, p2ptypes.GoodbyeCodeGenericError, remotePeer); err != nil {
			return err
		}
		return originalErr
	}
	s.cfg.P2P.Peers().SetChainState(remotePeer, m)

	if err := s.respondWithStatus(ctx, stream); err != nil {
		return err
	}
	closeStream(stream, log)
	return nil
}

func (s *Service) respondWithStatus(ctx context.Context, stream network.Stream) error {
	forkDigest, err := s.forkDigest()
	if err != nil {
		return err
	}
	resp := &pb.Status{
		ForkDigest:     forkDigest[:],
		FinalizedRoot:  make([]byte, 32),
		FinalizedEpoch: 0,
		HeadRoot:       make([]byte, 32),
		HeadSlot:       0,
	}

	if _, err := stream.Write([]byte{responseCodeSuccess}); err != nil {
		log.WithError(err).Debug("Could not write to stream")
	}
	_, err = s.cfg.P2P.Encoding().EncodeWithMaxLength(stream, resp)
	return err
}

func (s *Service) validateStatusMessage(ctx context.Context, msg *pb.Status) error {
	return nil
}
