package handler

import (
	"context"
	"errors"
	taskmngpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/taskmng"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/libp2p/go-libp2p-core/peer"
)

// SendTaskResultMsg sends taskResult to other peer, if the task has finished.
func SendTaskResultMsg(ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *taskmngpb.TaskResultMsg) error {
	ctx, cancel := context.WithTimeout(ctx, respTimeout)
	defer cancel()

	// send request on the special topic.
	stream, err := p2pProvider.Send(ctx, req, p2p.RPCTaskResultMsgTopic, pid)
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
