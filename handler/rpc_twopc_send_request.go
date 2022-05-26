package handler

import (
	"context"
	"errors"
	twopcpb "github.com/datumtechs/datum-network-carrier/lib/netmsg/consensus/twopc"
	"github.com/datumtechs/datum-network-carrier/p2p"
	"github.com/libp2p/go-libp2p-core/peer"
)

// SendTwoPcPrepareMsg sends 2pc prepareMsg to other peer.
//
// Deprecate: Use Broadcast to send message.
func SendTwoPcPrepareMsg(ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *twopcpb.PrepareMsg) error {
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
//
// Deprecate: Use Broadcast to send message.
func SendTwoPcPrepareVote (ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *twopcpb.PrepareVote) error {
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
//
// Deprecate: Use Broadcast to send message.
func SendTwoPcConfirmMsg (ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *twopcpb.ConfirmMsg) error {
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
//
// Deprecate: Use Broadcast to send message.
func SendTwoPcConfirmVote (ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *twopcpb.ConfirmVote) error {
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
//
// Deprecate: Use Broadcast to send message.
func SendTwoPcCommitMsg (ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *twopcpb.CommitMsg) error {
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

