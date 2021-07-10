package grpclient

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/common/runutil"
	"github.com/RosettaFlow/Carrier-Go/lib/fighter/datasvc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"sync"
	"time"
)

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

func (c *DataNodeClient) GetStatus(ctx context.Context) (*datasvc.GetStatusReply, error) {
	return nil, nil
}

func (c *DataNodeClient) ListData(ctx context.Context) (*datasvc.ListDataReply, error) {
	return nil, nil
}

func (c *DataNodeClient) UploadData(ctx context.Context, content []byte, metadata *datasvc.FileInfo) (*datasvc.UploadReply, error) {
	stream, err := c.dataProviderClient.UploadData(ctx)
	if err != nil {
		log.WithError(err).Error("Invoke UploadData failed")
		return nil, err
	}
	// send metadata
	err = stream.Send(&datasvc.UploadRequest{
		Data: &datasvc.UploadRequest_Meta{
			Meta: metadata,
		},
	})
	if err != nil {
		log.WithError(err).Error("stream request(metadata) for uploadData failed")
		return nil, err
	}
	// send content
	err = stream.Send(&datasvc.UploadRequest{
		Data: &datasvc.UploadRequest_Content{
			Content: content,
		},
	})
	if err != nil {
		log.WithError(err).Error("stream request(content) for uploadData failed")
		return nil, err
	}
	// receive response
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.WithError(err).Error("UploadData get response failed")
		return nil, err
	}
	return res, nil
}

func (c *DataNodeClient) BatchUpload(ctx context.Context) (datasvc.DataProvider_BatchUploadClient, error) {
	return nil, nil
}

func (c *DataNodeClient) DownloadData(ctx context.Context, in *datasvc.DownloadRequest) (datasvc.DataProvider_DownloadDataClient, error) {
	return nil, nil
}

func (c *DataNodeClient) DeleteData(ctx context.Context, in *datasvc.DownloadRequest) (*datasvc.UploadReply, error) {
	return nil, nil
}

func (c *DataNodeClient) SendSharesData(ctx context.Context, in *datasvc.SendSharesDataRequest) (*datasvc.SendSharesDataReply, error) {
	return nil, nil
}
