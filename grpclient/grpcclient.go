package grpclient

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"
)

var (
	RPCMaxCallRecvMsgSize = grpc.MaxCallRecvMsgSize(1 << 23) // 1 << 23 == 1024*1024*8 == 8388608 == 8mb
	RPCMaxCallSendMsgSize = grpc.MaxCallSendMsgSize(1 << 23)
)

// Dial connects a client to the given URL.
func dial(grpcurl string) (*grpc.ClientConn, error) {
	return dialContext(context.Background(), grpcurl)
}

func dialContext(ctx context.Context, grpcurl string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithReturnConnectionError(),
		grpc.WithDefaultCallOptions(RPCMaxCallRecvMsgSize),
		grpc.WithTimeout(2 * time.Second), // todo 默认 grpc client 连接内部资源, 先给 2s 超时
	}
	c, err := grpc.DialContext(ctx, grpcurl, opts...)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func dialContextWithSecure(ctx context.Context, addr string, remoteCert string) (*grpc.ClientConn, error) {
	creds, err := credentials.NewClientTLSFromFile(remoteCert, "")
	if err != nil {
		return nil, err
	}
	security := grpc.WithTransportCredentials(creds)
	opts := []grpc.DialOption{
		security,
		grpc.WithReturnConnectionError(),
		grpc.WithTimeout(2 * time.Second), // todo 默认 grpc client 连接内部资源, 先给 2s 超时
	}
	c, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, err
	}
	return c, nil
}
