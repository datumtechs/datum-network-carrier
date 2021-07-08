package grpclient

import (
	"context"
	"github.com/libp2p/go-libp2p-core/peer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
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

func (c *DataNodeClient) Close() {
	if c.cancel != nil {
		c.cancel()
	}
	c.conn.Close()
}

func (c *DataNodeClient) GetClientConn() *grpc.ClientConn {
	return c.conn
}


func (c *DataNodeClient) ConnStatus() connectivity.State {
	return c.conn.GetState()
}

func (c *DataNodeClient) IsConnected() bool {
	switch c.ConnStatus() {
	case connectivity.Ready, connectivity.Idle:
		return true
	default:
		return false
	}
}

func (c *DataNodeClient) Reconnect() error {
	if !c.IsConnected() {
		conn, err := dialContext(c.ctx, c.addr)
		if err != nil {
			return err
		}
		c.conn = conn
	}
	return nil
}