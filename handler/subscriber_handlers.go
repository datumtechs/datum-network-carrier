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

func (s *Service) gossipTestDataSubscriber(ctx context.Context, msg proto.Message) error {
	ve, ok := msg.(*rpcpb.SignedGossipTestData)
	if !ok {
		return fmt.Errorf("wrong type, expected: *ethpb.SignedVoluntaryExit got: %T", msg)
	}

	if ve.Data == nil {
		return errors.New("data can't be nil")
	}
	//TODO: do cache, callback...
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