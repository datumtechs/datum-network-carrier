package grpclient

import (
	"context"
	"errors"
	"github.com/RosettaFlow/Carrier-Go/common/runutil"
	"github.com/RosettaFlow/Carrier-Go/lib/fighter/datasvc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"sync"
	"time"
)

const defaultRequestTime = 1 * time.Second

type DataNodeClient struct {
	ctx    context.Context
	cancel context.CancelFunc
	conn   *grpc.ClientConn
	addr   string
	nodeId string
	connMu sync.RWMutex

	//TODO: define some client...
	dataProviderClient datasvc.DataProviderClient
}

func NewDataNodeClient(ctx context.Context, addr string, nodeId string) (*DataNodeClient, error) {
	ctx, cancel := context.WithCancel(ctx)
	client := &DataNodeClient{
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

func NewDataNodeClientWithConn(ctx context.Context, addr string, nodeId string) (*DataNodeClient, error) {
	ctx, cancel := context.WithCancel(ctx)
	conn, err := dialContext(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &DataNodeClient{
		ctx:                ctx,
		cancel:             cancel,
		conn:               conn,
		addr:               addr,
		nodeId:             nodeId,
		dataProviderClient: datasvc.NewDataProviderClient(conn),
	}, nil
}

func (c *DataNodeClient) Close() {
	if c.cancel != nil {
		c.cancel()
	}
	c.conn.Close()
}

func (c *DataNodeClient) connecting() {
	if c.IsConnected() {
		return
	}
	c.connMu.Lock()
	conn, err := dialContext(c.ctx, c.addr)
	c.connMu.Unlock()
	c.dataProviderClient = datasvc.NewDataProviderClient(conn)
	if err != nil {
		log.WithError(err).WithField("id", c.nodeId).Error("Connect GRPC server(for datanode) failed")
	}
	c.conn = conn
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

func (c *DataNodeClient) GetStatus() (*datasvc.GetStatusReply, error) {
	return nil, errors.New("method GetStatus not implemented")
}

func (c *DataNodeClient) ListData() (*datasvc.ListDataReply, error) {
	return nil, errors.New("method ListData not implemented")
}

func (c *DataNodeClient) UploadData(content []byte, metadata *datasvc.FileInfo) (*datasvc.UploadReply, error) {
	return nil, errors.New("method UploadData not implemented")
}

func (c *DataNodeClient) BatchUpload() (datasvc.DataProvider_BatchUploadClient, error) {
	return nil, errors.New("method BatchUpload not implemented")
}

func (c *DataNodeClient) DownloadData(in *datasvc.DownloadRequest) (datasvc.DataProvider_DownloadDataClient, error) {
	return nil, errors.New("method DownloadData not implemented")
}

func (c *DataNodeClient) DeleteData(in *datasvc.DownloadRequest) (*datasvc.UploadReply, error) {
	return nil, errors.New("method DeleteData not implemented")
}

func (c *DataNodeClient) SendSharesData(in *datasvc.SendSharesDataRequest) (*datasvc.SendSharesDataReply, error) {
	return nil, errors.New("method SendSharesData not implemented")
}
