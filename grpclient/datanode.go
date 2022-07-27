package grpclient

import (
	"context"
	"errors"
	"github.com/datumtechs/datum-network-carrier/common/runutil"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	fighterapipb "github.com/datumtechs/datum-network-carrier/pb/fighter/api/data"
	fightertypespb "github.com/datumtechs/datum-network-carrier/pb/fighter/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
	"time"
)

const defaultRequestTime = 1 * time.Second

type DataNodeClient struct {
	ctx         context.Context
	cancel      context.CancelFunc
	conn        *grpc.ClientConn
	addr        string
	nodeId      string
	connMu      sync.RWMutex
	connStartAt int64

	//TODO: define some client...
	dataProviderClient fighterapipb.DataProviderClient
}

func NewDataNodeClient(ctx context.Context, addr string, nodeId string) (*DataNodeClient, error) {
	ctx, cancel := context.WithCancel(ctx)
	conn, err := dialContext(ctx, addr)
	if err != nil {
		return nil, err
	}
	client := &DataNodeClient{
		ctx:                ctx,
		cancel:             cancel,
		addr:               addr,
		nodeId:             nodeId,
		conn:               conn,
		connStartAt:        timeutils.UnixMsec(),
		dataProviderClient: fighterapipb.NewDataProviderClient(conn),
	}
	log.Debugf("New a dataNode Client, dataNodeId: {%s}, timestamp: {%d}", nodeId, client.connStartAt)
	// try to connect grpc server.
	runutil.RunEvery(client.ctx, 10*time.Second, func() {
		client.connecting()
	})
	return client, nil
}

func (c *DataNodeClient) Close() {
	c.connStartAt = 0
	if c.cancel != nil {
		c.cancel()
	}
	if nil != c.conn {
		c.conn.Close()
	}
}

func (c *DataNodeClient) connecting() error {
	if c.IsConnected() {
		return nil
	}
	c.connMu.Lock()
	defer c.connMu.Unlock()
	conn, err := dialContext(c.ctx, c.addr)
	if err != nil {
		log.WithError(err).WithField("id", c.nodeId).WithField("adrr", c.addr).Error("Connect GRPC server(for datanode) failed")
		return err
	}
	c.dataProviderClient = fighterapipb.NewDataProviderClient(conn)
	c.conn = conn
	c.connStartAt = timeutils.UnixMsec()
	log.Debugf("Update a dataNode Client, dataNodeId: {%s}, timestamp: {%d}", c.nodeId, c.connStartAt)
	return nil
}

func (c *DataNodeClient) GetClientConn() *grpc.ClientConn {
	return c.conn
}

func (c *DataNodeClient) ConnStatus() connectivity.State {
	c.connMu.RLock()
	defer c.connMu.RUnlock()
	if c.conn == nil {
		return connectivity.Shutdown
	}
	return c.conn.GetState()
}

func (c *DataNodeClient) GetAddress() string {
	return c.addr
}

func (c *DataNodeClient) GetConnStartAt() int64 {
	return c.connStartAt
}

func (c *DataNodeClient) IsConnected() bool {
	switch c.ConnStatus() {
	case connectivity.Ready, connectivity.Idle:
		return true
	default:
		return false
	}
}

func (c *DataNodeClient) IsNotConnected() bool { return !c.IsConnected() }

func (c *DataNodeClient) Reconnect() error {
	err := c.connecting()
	if err != nil {
		log.WithError(err).WithField("dataNodeId", c.nodeId).Debug("Reconnect to dataNode failed")
		return err
	}
	return nil
}

func (c *DataNodeClient) RunningDuration() int64 {
	if c.IsNotConnected() {
		c.connStartAt = 0
		log.Warnf("Call dataNodeClient.RunningDuration(), the dataNode Client was not connected, dataNodeId: {%s}, timestamp: {%d}, connStatus: {%s}",
			c.nodeId, c.connStartAt, c.ConnStatus().String())
		return 0
	}
	log.Debugf("Call dataNodeClient.RunningDuration(), dataNodeId: {%s}, timestamp: {%d}, now: {%d}", c.nodeId, c.connStartAt, timeutils.UnixMsec())
	return timeutils.UnixMsec() - c.connStartAt
}

func (c *DataNodeClient) GetStatus() (*fighterapipb.GetStatusReply, error) {
	ctx, cancel := context.WithTimeout(c.ctx, 20*defaultRequestTime)
	defer cancel()
	in := new(emptypb.Empty)
	return c.dataProviderClient.GetStatus(ctx, in)
}

func (c *DataNodeClient) ListData() (*fighterapipb.ListDataReply, error) {
	return nil, errors.New("method ListData not implemented")
}

func (c *DataNodeClient) UploadData(content []byte, metadata *fighterapipb.UploadRequest) (*fighterapipb.UploadReply, error) {
	return nil, errors.New("method UploadData not implemented")
}

func (c *DataNodeClient) BatchUpload() (fighterapipb.DataProvider_BatchUploadClient, error) {
	return nil, errors.New("method BatchUpload not implemented")
}

func (c *DataNodeClient) DownloadData(req *fighterapipb.DownloadRequest) (fighterapipb.DataProvider_DownloadDataClient, error) {
	ctx, cancel := context.WithTimeout(c.ctx, 20*defaultRequestTime)
	defer cancel()
	return c.dataProviderClient.DownloadData(ctx, req)
}

func (c *DataNodeClient) DeleteData(in *fighterapipb.DownloadRequest) (*fighterapipb.UploadReply, error) {
	return nil, errors.New("method DeleteData not implemented")
}

func (c *DataNodeClient) HandleTaskReadyGo(req *fightertypespb.TaskReadyGoReq) (*fightertypespb.TaskReadyGoReply, error) {
	ctx, cancel := context.WithTimeout(c.ctx, 20*defaultRequestTime)
	defer cancel()
	return c.dataProviderClient.HandleTaskReadyGo(ctx, req)
}

func (c *DataNodeClient) HandleCancelTask(req *fightertypespb.TaskCancelReq) (*fightertypespb.TaskCancelReply, error) {
	ctx, cancel := context.WithTimeout(c.ctx, 20*defaultRequestTime)
	defer cancel()
	return c.dataProviderClient.HandleCancelTask(ctx, req)
}
