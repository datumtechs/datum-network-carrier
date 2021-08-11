package grpclient

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/common/runutil"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
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
	connStartAt   int64

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
	runutil.RunEvery(client.ctx, 10 * time.Second, func() {
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
	c.connStartAt = 0
	if c.cancel != nil {
		c.cancel()
	}
	if nil != c.conn {
		c.conn.Close()
	}
}

func (c *JobNodeClient) connecting() error {
	if c.IsConnected() {
		return nil
	}
	c.connMu.Lock()
	defer c.connMu.Unlock()
	conn, err := dialContext(c.ctx, c.addr)
	if err != nil {
		log.WithError(err).WithField("id", c.nodeId).WithField("addr", c.addr).Error("Connect GRPC server(for jobnode) failed")
		return err
	}
	// set client with conn for computeProviderClient
	c.computeProviderClient = computesvc.NewComputeProviderClient(conn)
	c.conn = conn
	c.connStartAt = timeutils.UnixMsec()
	return nil
}

func (c *JobNodeClient) GetClientConn() *grpc.ClientConn {
	return c.conn
}

func (c *JobNodeClient) ConnStatus() connectivity.State {
	c.connMu.RLock()
	defer c.connMu.RUnlock()
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
func (c *JobNodeClient) IsNotConnected() bool { return !c.IsConnected()}

func (c *JobNodeClient) Reconnect() error {
	err := c.connecting()
	if err != nil {
		log.WithError(err).WithField("jobNode", c.nodeId).Debug("Reconnect to jobNode failed")
		return err
	}
	return nil
}

func (c *JobNodeClient) RunningDuration() int64 {
	if c.IsNotConnected() {
		c.connStartAt = 0
		return 0
	}
	return timeutils.UnixMsec() - c.connStartAt
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
