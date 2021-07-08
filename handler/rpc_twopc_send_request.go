package handler

import (
	"context"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/libp2p/go-libp2p-core/peer"
)



// SendTwoPcPrepareMsg sends 2pc prepareMsg to other peer.
func SendTwoPcPrepareMsg(ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *pb.PrepareMsg) error {
	// send request on the special topic.
	stream, err := p2pProvider.Send(ctx, req, p2p.RPCTwoPcPrepareMsgTopic, pid)
	if err != nil {
		return err
	}
	defer closeStream(stream, log)
	return nil
}

// SendTwoPcPrepareVote sends 2pc prepareVote to other peer.
func SendTwoPcPrepareVote (ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *pb.PrepareMsg) error {
	// send request on the special topic.
	stream, err := p2pProvider.Send(ctx, req, p2p.RPCTwoPcPrepareVoteTopic, pid)
	if err != nil {
		return err
	}
	defer closeStream(stream, log)
	return nil
}
