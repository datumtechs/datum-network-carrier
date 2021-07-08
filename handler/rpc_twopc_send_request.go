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
func SendTwoPcPrepareVote (ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *pb.PrepareVote) error {
	// send request on the special topic.
	stream, err := p2pProvider.Send(ctx, req, p2p.RPCTwoPcPrepareVoteTopic, pid)
	if err != nil {
		return err
	}
	defer closeStream(stream, log)
	return nil
}


// SendTwoPcConfirmMsg sends 2pc ConfirmMsg to other peer.
func SendTwoPcConfirmMsg (ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *pb.ConfirmMsg) error {
	// send request on the special topic.
	stream, err := p2pProvider.Send(ctx, req, p2p.RPCTwoPcConfirmMsgTopic, pid)
	if err != nil {
		return err
	}
	defer closeStream(stream, log)
	return nil
}


// SendTwoPcConfirmVote sends 2pc ConfirmVote to other peer.
func SendTwoPcConfirmVote (ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *pb.ConfirmVote) error {
	// send request on the special topic.
	stream, err := p2pProvider.Send(ctx, req, p2p.RPCTwoPcConfirmVoteTopic, pid)
	if err != nil {
		return err
	}
	defer closeStream(stream, log)
	return nil
}


// SendTwoPcCmmitMsg sends 2pc CommitMsg to other peer.
func SendTwoPcCmmitMsg (ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *pb.CommitMsg) error {
	// send request on the special topic.
	stream, err := p2pProvider.Send(ctx, req, p2p.RPCTwoPcCommitMsgTopic, pid)
	if err != nil {
		return err
	}
	defer closeStream(stream, log)
	return nil
}