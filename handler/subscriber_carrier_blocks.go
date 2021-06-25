package handler

import (
	"context"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

func (s *Service) carrierBlockSubscriber(ctx context.Context, msg proto.Message) error {
	block, ok := msg.(*libtypes.BlockData)
	if !ok {
		return errors.New("message is not type *libtypes.BlockData")
	}

	if block == nil {
		return errors.New("nil block")
	}
	//TODO: need to coding...
	// cache for limit.

	// call chain, e.g. s.cfg.Chain.ReceiveBlock(ctx, signed, root)
	return nil
}
