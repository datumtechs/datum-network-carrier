package handler

import (
	"context"
	libp2ptypes "github.com/datumtechs/datum-network-carrier/lib/p2p/v1"
	libp2ppb "github.com/datumtechs/datum-network-carrier/lib/rpc/debug/v1"
	libtypes "github.com/datumtechs/datum-network-carrier/lib/types"
	"github.com/datumtechs/datum-network-carrier/p2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	"io"
)

// ErrInvalidFetchedData is thrown if stream fails to provide requested blocks.
var ErrInvalidFetchedData = errors.New("invalid data returned from peer")

// CarrierBlockProcessor defines a block processing function, which allows to start utilizing
// blocks even before all blocks are ready.
type CarrierBlockProcessor func(block *libtypes.BlockData) error

// SendCarrierBlocksByRangeRequest sends CarrierBlocksByRange and returns fetched blocks, if any.
func SendCarrierBlocksByRangeRequest(
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

// SendGossipTestDataByRangeRequest for testing
func SendGossipTestDataByRangeRequest(ctx context.Context, p2pProvider p2p.P2P, pid peer.ID,
	req *libp2ppb.GossipTestData) ([]*libp2ppb.SignedGossipTestData, error) {

	stream, err := p2pProvider.Send(ctx, req, p2p.RPCGossipTestDataByRangeTopic, pid)
	if err != nil {
		return nil, err
	}
	defer closeStream(stream, log)

	// Augment block processing function, if non-nil block processor is provided.
	datas := make([]*libp2ppb.SignedGossipTestData, 0, req.Count)
	process := func(blk *libp2ppb.SignedGossipTestData) error {
		log.Infof("Send done and response info, count: %d, step: %d", blk.Data.Count, blk.Data.Step)
		datas = append(datas, blk)
		return nil
	}
	for i := uint64(0); ; i++ {
		isFirstChunk := i == 0
		blk, err := ReadChunkedGossipTestData(stream, p2pProvider, isFirstChunk)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if err := process(blk); err != nil {
			return nil, err
		}
	}
	return datas, nil
}