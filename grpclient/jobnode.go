package grpclient

import (
	"context"
	"github.com/libp2p/go-libp2p-core/peer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
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

func (c *JobNodeClient) Close() {
	if c.cancel != nil {
		c.cancel()
	}
	c.conn.Close()
}

func (c *JobNodeClient) GetClientConn() *grpc.ClientConn {
	return c.conn
}

func (c *JobNodeClient) ConnStatus() connectivity.State {
	return c.conn.GetState()
}

func (c *JobNodeClient) IsConnected() bool {
	switch c.ConnStatus() {
	case connectivity.Ready, connectivity.Idle:
		return true
	default:
		return false
	}
}

func (c *JobNodeClient) Reconnect() error {
	if !c.IsConnected() {
		conn, err := dialContext(c.ctx, c.addr)
		if err != nil {
			return err
		}
		c.conn = conn
	}
	return nil
}