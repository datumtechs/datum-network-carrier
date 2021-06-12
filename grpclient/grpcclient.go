package grpclient

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/core/types"
	"github.com/RosettaFlow/Carrier-Go/lib/center/api"
	"google.golang.org/grpc"
)

// Client defines typed wrapper for the CenterData GRPC API.
type GrpcClient struct {
	c *grpc.ClientConn

	// grpc service
	metadataService api.MetaDataServiceClient
	resourceService api.ResourceServiceClient
	identityService api.IdentityServiceClient
	taskService 	api.TaskServiceClient
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
	return &GrpcClient{
		c: c,
		metadataService: api.NewMetaDataServiceClient(c),
		resourceService: api.NewResourceServiceClient(c),
		identityService: api.NewIdentityServiceClient(c),
		taskService:     api.NewTaskServiceClient(c),
	}
}

func (gc *GrpcClient) Close() {
	gc.c.Close()
}

func (gc *GrpcClient) GetClientConn() *grpc.ClientConn {
	return gc.c
}

// DataChain Access
func (gc *GrpcClient) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	/*conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())*/
	return nil, nil
}