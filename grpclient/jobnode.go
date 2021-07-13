package grpclient

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/common/runutil"
	"github.com/RosettaFlow/Carrier-Go/lib/fighter/common"
	"github.com/RosettaFlow/Carrier-Go/lib/fighter/computesvc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"sync"
	"time"
)

type JobNodeClient struct {
	ctx    context.Context
	cancel context.CancelFunc
	conn   *grpc.ClientConn
	addr   string
	nodeId string
	connMu sync.RWMutex

	//TODO: define some client...
	computeProviderClient computesvc.ComputeProviderClient
}

func NewJobNodeClient(ctx context.Context, addr string, nodeId string) (*JobNodeClient, error) {
	ctx, cancel := context.WithCancel(ctx)
	client := &JobNodeClient{
		ctx:    ctx,
		cancel: cancel,
		addr:   addr,
		nodeId: nodeId,
	}
	// try to connect grpc server.
	runutil.RunEvery(client.ctx, 2*time.Second, func() {
		client.connecting()
	})
	return client, nil
}

func NewJobNodeClientWithConn(ctx context.Context, addr string, nodeId string) (*JobNodeClient, error) {
	ctx, cancel := context.WithCancel(ctx)
	conn, err := dialContext(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &JobNodeClient{
		ctx:                   ctx,
		cancel:                cancel,
		conn:                  conn,
		addr:                  addr,
		nodeId:                nodeId,
		computeProviderClient: computesvc.NewComputeProviderClient(conn),
	}, nil
}

func (c *JobNodeClient) Close() {
	if c.cancel != nil {
		c.cancel()
	}
	c.conn.Close()
}

func (c *JobNodeClient) connecting() {
	if c.IsConnected() {
		return
	}
	c.connMu.Lock()
	conn, err := dialContext(c.ctx, c.addr)
	c.connMu.Unlock()
	// set client with conn for computeProviderClient
	c.computeProviderClient = computesvc.NewComputeProviderClient(conn)
	if err != nil {
		log.WithError(err).WithField("id", c.nodeId).Error("Connect GRPC server(for datanode) failed")
	}
	c.conn = conn
}

func (c *JobNodeClient) GetClientConn() *grpc.ClientConn {
	return c.conn
}

func (c *JobNodeClient) ConnStatus() connectivity.State {
	if c.conn == nil {
		return connectivity.Shutdown
	}
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
	c.connecting()
	return nil
}

func (c *JobNodeClient) GetStatus() (*computesvc.GetStatusReply, error) {
	return nil, errors.New("method GetStatus not implemented")
}

func (c *JobNodeClient) GetTaskDetails(ctx context.Context, taskIds []string) (*computesvc.GetTaskDetailsReply, error) {
	return nil, errors.New("method GetTaskDetails not implemented")
}

func (c *JobNodeClient) UploadShard(ctx context.Context) (computesvc.ComputeProvider_UploadShardClient, error){
	return nil, errors.New("method UploadShard not implemented")
}

func (c *JobNodeClient) HandleTaskReadyGo(req *common.TaskReadyGoReq) (*common.TaskReadyGoReply, error) {
	ctx, cancel := context.WithTimeout(c.ctx, defaultRequestTime)
	defer cancel()
	return c.computeProviderClient.HandleTaskReadyGo(ctx, req)
}
