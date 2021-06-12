package grpclient

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/core/types"
	"github.com/RosettaFlow/Carrier-Go/lib/center/api"
	"google.golang.org/grpc"
	"time"
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

// ************************************** MetaData module *******************************************************

// MetaDataSave saves new metadata to database.
func (gc *GrpcClient) MetaDataSave(ctx context.Context, request *api.MetaDataSaveRequest) (*api.MetaDataSaveResponse, error) {
	cctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.metadataService.MetaDataSave(cctx, request)
}

func (gc *GrpcClient) GetMetaDataSummaryList(ctx context.Context) (*api.MetaDataSummaryListResponse, error) {
	return nil, nil
}

func (gc *GrpcClient) GetMetaDataSummaryByState(ctx context.Context, request *api.MetaDataSummaryByStateRequest) (*api.MetaDataSummaryByStateResponse, error) {
	return nil, nil
}

func (gc *GrpcClient) RevokeMetaData(ctx context.Context, request *api.RevokeMetaDataRequest) (*api.SimpleResponse, error) {
	return nil, nil
}

// ************************************** Resource module *******************************************************

func (gc *GrpcClient) SaveResource(ctx context.Context, request *api.PublishPowerRequest) (*api.PublishPowerResponse, error) {
	return nil, nil
}

func (gc *GrpcClient) RevokeResource(ctx context.Context, request *api.RevokePowerRequest) (*api.SimpleResponse, error) {
	return nil, nil
}

func (gc *GrpcClient) GetPowerTotalSummaryList(ctx context.Context) (*api.PowerTotalSummaryListResponse, error) {
	return nil, nil
}

// ************************************** Identity module *******************************************************

func (gc *GrpcClient) SaveIdentity(ctx context.Context, request *api.SaveIdentityRequest) (*api.SimpleResponse, error) {
	return nil, nil
}

func (gc *GrpcClient) RevokeIdentityJoin(ctx context.Context, request *api.RevokeIdentityJoinRequest) (*api.SimpleResponse, error) {
	return nil, nil
}

// ************************************** Task module *******************************************************

func (gc *GrpcClient) SaveTask(ctx context.Context, request *api.TaskDetail) error {
	return nil
}

func (gc *GrpcClient) GetDetailTask(ctx context.Context, request *api.DetailTaskRequest) (*api.TaskDetail, error) {
	return nil, nil
}

func (gc *GrpcClient) ListTask(ctx context.Context, request *api.TaskListRequest) (*api.TaskListResponse, error) {
	return nil, nil
}

func (gc *GrpcClient) ListTaskEvent(ctx context.Context, request *api.TaskEventRequest) (*api.TaskEventResponse, error) {
	return nil, nil
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