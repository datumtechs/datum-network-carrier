package grcpclient

import (
	"context"
	"google.golang.org/grpc"
)

// Client defines typed wrapper for the CenterData GRPC API.
type GrpcClient struct {
	c *grpc.ClientConn
}

// Dial connects a client to the given URL.
func Dial(grpcurl string) (*GrpcClient, error) {
	return DialContext(context.Background(), grpcurl)
}

func DialContext(ctx context.Context, grpcurl string) (*GrpcClient, error) {
	c, err := grpc.Dial(grpcurl, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return NewGrpcClient(c), nil
}

// NewClient creates a client that uses the given GRPC client.
func NewGrpcClient(c *grpc.ClientConn) *GrpcClient {
	return &GrpcClient{c: c}
}

func (gc *GrpcClient) Close() {
	gc.c.Close()
}
