package handler

import (
	"context"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

func (s *Service) carrierBlockSubscriber(ctx context.Context, msg proto.Message) error {
	block, ok := msg.(*carriertypespb.BlockData)
	if !ok {
		return errors.New("message is not type *carriertypespb.BlockData")
	}

	if block == nil {
		return errors.New("nil block")
	}
	//TODO: need to coding...
	// cache for limit.

	// call chain, e.g. s.cfg.Chain.ReceiveBlock(ctx, signed, root)
	return nil
}
