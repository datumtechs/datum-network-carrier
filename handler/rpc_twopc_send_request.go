package handler

import (
	"context"
	libp2ptypes "github.com/RosettaFlow/Carrier-Go/lib/p2p/v1"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/libp2p/go-libp2p-core/peer"
)



// SendCarrierBlocksByRangeRequest sends CarrierBlocksByRange and returns fetched blocks, if any.
func SendTwoPcPrepareMsg(
	ctx context.Context, p2pProvider p2p.P2P, pid peer.ID,
	req *libp2ptypes.CarrierBlocksByRangeRequest, blockProcessor CarrierBlockProcessor,
) ([]*libtypes.BlockData, error) {

	// send request on the special topic.
	stream, err := p2pProvider.Send(ctx, req, p2p.RPCBlocksByRangeTopic, pid)
	if err != nil {
		return nil, err
	}
	defer closeStream(stream, log)

	// Augment block processing function, if non-nil block processor is provided.
	blocks := make([]*libtypes.BlockData, 0, req.Count)
	process := func(blk *libtypes.BlockData) error {
		blocks = append(blocks, blk)
		if blockProcessor != nil {
			return blockProcessor(blk)
		}
		return nil
	}
	_ = process
	//TODO: ....
	return blocks, nil
}