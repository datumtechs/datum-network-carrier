package handler

import (
	"context"
	"errors"
	carriernetmsgtaskmngpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/taskmng"
	"github.com/datumtechs/datum-network-carrier/p2p"
	"github.com/libp2p/go-libp2p-core/peer"
)

// SendTaskResultMsg sends taskResult to other peer, if the task has finished.
func SendTaskResultMsg (ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *carriernetmsgtaskmngpb.TaskResultMsg) error {
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

func SendTaskResourceUsageMsg(ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *carriernetmsgtaskmngpb.TaskResourceUsageMsg) error {
	ctx, cancel := context.WithTimeout(ctx, respTimeout)
	defer cancel()

	// send request on the special topic.
	stream, err := p2pProvider.Send(ctx, req, p2p.RPCTaskResourceUsageMsgTopic, pid)
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

func SendTaskTerminateMsg(ctx context.Context, p2pProvider p2p.P2P, pid peer.ID, req *carriernetmsgtaskmngpb.TaskTerminateMsg) error {
	ctx, cancel := context.WithTimeout(ctx, respTimeout)
	defer cancel()

	// send request on the special topic.
	stream, err := p2pProvider.Send(ctx, req, p2p.RPCTaskTerminateMsgTopic, pid)
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