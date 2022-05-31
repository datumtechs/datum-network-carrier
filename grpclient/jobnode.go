package grpclient

import (
	"context"
	"github.com/datumtechs/datum-network-carrier/common/runutil"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	fighterapicomputepb "github.com/datumtechs/datum-network-carrier/pb/fighter/api/compute"
	fightertypespb "github.com/datumtechs/datum-network-carrier/pb/fighter/types"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
	"time"
)

type JobNodeClient struct {
	ctx         context.Context
	cancel      context.CancelFunc
	conn        *grpc.ClientConn
	addr        string
	nodeId      string
	connMu      sync.RWMutex
	connStartAt int64

	//TODO: define some client...
	computeProviderClient fighterapicomputepb.ComputeProviderClient
}

func NewJobNodeClient(ctx context.Context, addr string, nodeId string) (*JobNodeClient, error) {
	ctx, cancel := context.WithCancel(ctx)
	conn, err := dialContext(ctx, addr)
	if err != nil {
		return nil, err
	}
	client := &JobNodeClient{
		ctx:                   ctx,
		cancel:                cancel,
		addr:                  addr,
		nodeId:                nodeId,
		conn:                  conn,
		connStartAt:           timeutils.UnixMsec(),
		computeProviderClient: fighterapicomputepb.NewComputeProviderClient(conn),
	}
	log.Debugf("New a jobNode Client, jobNodeId: {%s}, timestamp: {%d}", nodeId, client.connStartAt)
	// try to connect grpc server.
	runutil.RunEvery(client.ctx, 10*time.Second, func() {
		client.connecting()
	})
	return client, nil
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
	c.computeProviderClient = fighterapicomputepb.NewComputeProviderClient(conn)
	c.conn = conn
	c.connStartAt = timeutils.UnixMsec()
	log.Debugf("Update a jobNode Client, jobNodeId: {%s}, timestamp: {%d}", c.nodeId, c.connStartAt)
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

func (c *JobNodeClient) GetAddress() string {
	return c.addr
}

func (c *JobNodeClient) GetConnStartAt() int64 {
	return c.connStartAt
}

func (c *JobNodeClient) IsConnected() bool {
	switch c.ConnStatus() {
	case connectivity.Ready, connectivity.Idle:
		return true
	default:
		return false
	}
}
func (c *JobNodeClient) IsNotConnected() bool { return !c.IsConnected() }

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
		log.Warnf("Call jobNodeClient.RunningDuration(), the jobNode Client was not connected, jobNodeId: {%s}, timestamp: {%d}, connStatus: {%s}",
			c.nodeId, c.connStartAt, c.ConnStatus().String())
		return 0
	}
	log.Debugf("Call jobNodeClient.RunningDuration(), jobNodeId: {%s}, timestamp: {%d}, now: {%d}", c.nodeId, c.connStartAt, timeutils.UnixMsec())
	return timeutils.UnixMsec() - c.connStartAt
}

func (c *JobNodeClient) GetStatus() (*fighterapicomputepb.GetStatusReply, error) {
	ctx, cancel := context.WithTimeout(c.ctx, 20*defaultRequestTime)
	defer cancel()
	in := new(emptypb.Empty)
	return c.computeProviderClient.GetStatus(ctx, in)
}

func (c *JobNodeClient) GetTaskDetails(ctx context.Context, taskIds []string) (*fighterapicomputepb.GetTaskDetailsReply, error) {
	return nil, errors.New("method GetTaskDetails not implemented")
}

func (c *JobNodeClient) UploadShard(ctx context.Context) (fighterapicomputepb.ComputeProvider_UploadShardClient, error) {
	return nil, errors.New("method UploadShard not implemented")
}

func (c *JobNodeClient) HandleTaskReadyGo(req *fightertypespb.TaskReadyGoReq) (*fightertypespb.TaskReadyGoReply, error) {
	ctx, cancel := context.WithTimeout(c.ctx, 20*defaultRequestTime)
	defer cancel()
	return c.computeProviderClient.HandleTaskReadyGo(ctx, req)
}

func (c *JobNodeClient) HandleCancelTask(req *fightertypespb.TaskCancelReq) (*fightertypespb.TaskCancelReply, error) {
	ctx, cancel := context.WithTimeout(c.ctx, 20*defaultRequestTime)
	defer cancel()
	return c.computeProviderClient.HandleCancelTask(ctx, req)
}
