package handler

import (
	"context"
	"errors"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/libp2p/go-libp2p-core/peer"
)

// SendTwoPcPrepareMsg sends 2pc prepareMsg to other peer.
func SendTwoPcPrepareMsg(ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *pb.PrepareMsg) error {
	ctx, cancel := context.WithTimeout(ctx, respTimeout)
	defer cancel()

	// send request on the special topic.
	stream, err := p2pProvider.Send(ctx, req, p2p.RPCTwoPcPrepareMsgTopic, pid)
	if err != nil {
		return err
	}
	defer closeStream(stream, log)
	code, errMsg, err := ReadStatusCode(stream, p2pProvider.Encoding())
	if err != nil {
		return err
	}
	if code != 0 {
		return errors.New(errMsg)
	}
	return nil
}

// SendTwoPcPrepareVote sends 2pc prepareVote to other peer.
func SendTwoPcPrepareVote (ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *pb.PrepareVote) error {
	ctx, cancel := context.WithTimeout(ctx, respTimeout)
	defer cancel()

	// send request on the special topic.
	stream, err := p2pProvider.Send(ctx, req, p2p.RPCTwoPcPrepareVoteTopic, pid)
	if err != nil {
		return err
	}
	defer closeStream(stream, log)
	code, errMsg, err := ReadStatusCode(stream, p2pProvider.Encoding())
	if err != nil {
		return err
	}
	if code != 0 {
		return errors.New(errMsg)
	}
	return nil
}

// SendTwoPcConfirmMsg sends 2pc ConfirmMsg to other peer.
func SendTwoPcConfirmMsg (ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *pb.ConfirmMsg) error {
	ctx, cancel := context.WithTimeout(ctx, respTimeout)
	defer cancel()

	// send request on the special topic.
	stream, err := p2pProvider.Send(ctx, req, p2p.RPCTwoPcConfirmMsgTopic, pid)
	if err != nil {
		return err
	}
	defer closeStream(stream, log)
	code, errMsg, err := ReadStatusCode(stream, p2pProvider.Encoding())
	if err != nil {
		return err
	}
	if code != 0 {
		return errors.New(errMsg)
	}
	return nil
}

// SendTwoPcConfirmVote sends 2pc ConfirmVote to other peer.
func SendTwoPcConfirmVote (ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *pb.ConfirmVote) error {
	ctx, cancel := context.WithTimeout(ctx, respTimeout)
	defer cancel()

	// send request on the special topic.
	stream, err := p2pProvider.Send(ctx, req, p2p.RPCTwoPcConfirmVoteTopic, pid)
	if err != nil {
		return err
	}
	defer closeStream(stream, log)
	code, errMsg, err := ReadStatusCode(stream, p2pProvider.Encoding())
	if err != nil {
		return err
	}
	if code != 0 {
		return errors.New(errMsg)
	}
	return nil
}

// SendTwoPcCommitMsg sends 2pc CommitMsg to other peer.
func SendTwoPcCommitMsg (ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *pb.CommitMsg) error {
	ctx, cancel := context.WithTimeout(ctx, respTimeout)
	defer cancel()

	// send request on the special topic.
	stream, err := p2pProvider.Send(ctx, req, p2p.RPCTwoPcCommitMsgTopic, pid)
	if err != nil {
		return err
	}
	defer closeStream(stream, log)
	code, errMsg, err := ReadStatusCode(stream, p2pProvider.Encoding())
	if err != nil {
		return err
	}
	if code != 0 {
		return errors.New(errMsg)
	}
	return nil
}