package handler

import (
	"context"
	"errors"
	"fmt"
	rpcpb "github.com/RosettaFlow/Carrier-Go/lib/rpc/v1"
	"github.com/gogo/protobuf/proto"
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