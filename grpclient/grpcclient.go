package grpclient

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/lib/center/api"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
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
func (gc *GrpcClient) SaveMetaData(ctx context.Context, request *api.MetaDataSaveRequest) (*api.MetaDataSaveResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.metadataService.MetaDataSave(ctx, request)
}

func (gc *GrpcClient) GetMetaDataSummaryList(ctx context.Context) (*api.MetaDataSummaryListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.metadataService.GetMetaDataSummaryList(ctx, &emptypb.Empty{})
}

func (gc *GrpcClient) GetMetadataList(ctx context.Context, request *api.MetadataListRequest) (*api.MetadataListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.metadataService.GetMetadataList(ctx, request)
}

func (gc *GrpcClient) GetMetadataById(ctx context.Context, request *api.MetadataByIdRequest) (*api.MetadataByIdResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.metadataService.GetMetadataById(ctx, request)
}

func (gc *GrpcClient) RevokeMetaData(ctx context.Context, request *api.RevokeMetaDataRequest) (*api.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.metadataService.RevokeMetaData(ctx, request)
}

// ************************************** Resource module *******************************************************

func (gc *GrpcClient) SaveResource(ctx context.Context, request *api.PublishPowerRequest) (*api.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.resourceService.PublishPower(ctx, request)
}

func (gc *GrpcClient) SyncPower(ctx context.Context, request *api.SyncPowerRequest) (*api.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.resourceService.SyncPower(ctx, request)
}

func (gc *GrpcClient) RevokeResource(ctx context.Context, request *api.RevokePowerRequest) (*api.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.resourceService.RevokePower(ctx, request)
}

func (gc *GrpcClient) GetPowerSummaryByNodeId(ctx context.Context, request *api.PowerSummaryByNodeIdRequest) (*api.PowerTotalSummaryResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.resourceService.GetPowerSummaryByNodeId(ctx, request)
}

func (gc *GrpcClient) GetPowerTotalSummaryList(ctx context.Context) (*api.PowerTotalSummaryListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.resourceService.GetPowerTotalSummaryList(ctx, &emptypb.Empty{})
}

func (gc *GrpcClient) GetPowerList(ctx context.Context, request *api.PowerListRequest) (*api.PowerListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.resourceService.GetPowerList(ctx, request)
}

// ************************************** Identity module *******************************************************

func (gc *GrpcClient) SaveIdentity(ctx context.Context, request *api.SaveIdentityRequest) (*api.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.identityService.SaveIdentity(ctx, request)
}

func (gc *GrpcClient) RevokeIdentityJoin(ctx context.Context, request *api.RevokeIdentityJoinRequest) (*api.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.identityService.RevokeIdentityJoin(ctx, request)
}

func (gc *GrpcClient) GetIdentityList(ctx context.Context, request *api.IdentityListRequest) (*api.IdentityListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.GetIdentityList(ctx, request)
}

// ************************************** Task module *******************************************************

func (gc *GrpcClient) SaveTask(ctx context.Context, request *api.TaskDetail) (*api.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.taskService.SaveTask(ctx, request)
}

func (gc *GrpcClient) GetDetailTask(ctx context.Context, request *api.DetailTaskRequest) (*api.TaskDetail, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.taskService.GetDetailTask(ctx, request)
}

func (gc *GrpcClient) ListTask(ctx context.Context, request *api.TaskListRequest) (*api.TaskListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.taskService.ListTask(ctx, request)
}

func (gc *GrpcClient) ListTaskEvent(ctx context.Context, request *api.TaskEventRequest) (*api.TaskEventResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return gc.taskService.ListTaskEvent(ctx, request)
}