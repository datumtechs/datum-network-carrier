package grpclient

import (
	"context"
	"github.com/libp2p/go-libp2p-core/peer"
	"google.golang.org/grpc"
)

type DataNodeClient struct {
	ctx    context.Context
	cancel context.CancelFunc
	conn   *grpc.ClientConn
	addr   string
	peerId peer.ID

	//TODO: define some client...
}

func NewDataNodeClient(ctx context.Context, addr string, peerId string) (*DataNodeClient, error) {
	ctx, cancel := context.WithCancel(ctx)
	conn, err := dialContext(ctx, addr)
	if err != nil {
		return nil, err
	}
	pid, _ := peer.Decode(peerId)
	return &DataNodeClient{
		ctx:    ctx,
		cancel: cancel,
		conn:   conn,
		addr:   addr,
		peerId: pid,
	}, nil
}

func (gc *DataNodeClient) Close() {
	if gc.cancel != nil {
		gc.cancel()
	}
	gc.conn.Close()
}

func (gc *DataNodeClient) GetClientConn() *grpc.ClientConn {
	return gc.conn
}
