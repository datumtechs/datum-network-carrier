package grpclient

import (
	"context"
	"github.com/libp2p/go-libp2p-core/peer"
	"google.golang.org/grpc"
)

type JobNodeClient struct {
	ctx    context.Context
	cancel context.CancelFunc
	conn   *grpc.ClientConn
	addr   string
	peerId peer.ID

	//TODO: define some client...
}

func NewJobNodeClient(ctx context.Context, addr string, peerId string) (*JobNodeClient, error) {
	ctx, cancel := context.WithCancel(ctx)
	conn, err := dialContext(ctx, addr)
	if err != nil {
		return nil, err
	}
	pid, _ := peer.Decode(peerId)
	return &JobNodeClient{
		ctx:    ctx,
		cancel: cancel,
		conn:   conn,
		addr:   addr,
		peerId: pid,
	}, nil
}

func (gc *JobNodeClient) Close() {
	if gc.cancel != nil {
		gc.cancel()
	}
	gc.conn.Close()
}

func (gc *JobNodeClient) GetClientConn() *grpc.ClientConn {
	return gc.conn
}