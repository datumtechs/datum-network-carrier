package grpclient

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"
)

// Dial connects a client to the given URL.
func dial(grpcurl string) (*grpc.ClientConn, error) {
	return dialContext(context.Background(), grpcurl)
}

func dialContext(ctx context.Context, grpcurl string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithReturnConnectionError(),
		grpc.WithTimeout(1 * time.Second),
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
		grpc.WithTimeout(1 * time.Second),
	}
	c, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, err
	}
	return c, nil
}