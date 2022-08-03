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
		grpc.WithTimeout(2 * time.Second), // todo By default, grpc client connects to internal resources and gives 2S timeout first
	}
	c, err := grpc.DialContext(ctx, grpcurl, opts...)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// DialContext creates a client connection to the given target. By default, it's
// a non-blocking dial (the function won't wait for connections to be
// established, and connecting happens in the background). To make it a blocking
// dial, use WithBlock() dial option.
//
// In the non-blocking case, the ctx does not act against the connection. It
// only controls the setup steps.
//
// In the blocking case, ctx can be used to cancel or expire the pending
// connection. Once this function returns, the cancellation and expiration of
// ctx will be noop. Users should call ClientConn.Close to terminate all the
// pending operations after this function returns.
// grpcurl = ip:port
func DialContext(ctx context.Context, grpcurl string, withBlock bool) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithReturnConnectionError(),
		grpc.WithDefaultCallOptions(RPCMaxCallRecvMsgSize),
		grpc.WithTimeout(2 * time.Second), // todo By default, grpc client connects to internal resources and gives 2S timeout first
	}
	if withBlock {
		opts = append(opts, grpc.WithBlock())
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
		grpc.WithTimeout(2 * time.Second), // todo By default, grpc client connects to internal resources and gives 2S timeout first
	}
	c, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, err
	}
	return c, nil
}
