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

// NewClient creates a client that uses the given GRPC client.
func NewGrpcClient(ctx context.Context, addr string) (*GrpcClient, error) {
	ctx, cancel := context.WithTimeout(ctx, 2 * time.Second)  //TODO 默认连接 dataCenter 的 grpc client 超时时间设为 2 s
	_ = cancel
	conn, err := dialContext(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &GrpcClient{
		c: conn,
		metadataService: api.NewMetaDataServiceClient(conn),
		resourceService: api.NewResourceServiceClient(conn),
		identityService: api.NewIdentityServiceClient(conn),
		taskService:     api.NewTaskServiceClient(conn),
	}, nil
}

func (gc *GrpcClient) Close() {
	gc.c.Close()
}

func (gc *GrpcClient) GetClientConn() *grpc.ClientConn {
	return gc.c
}

// MetaDataSave saves new metadata to database.
func (gc *GrpcClient) SaveMetaData(ctx context.Context, request *api.MetaDataSaveRequest) (*api.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return gc.metadataService.MetaDataSave(ctx, request)
}

func (gc *GrpcClient) GetMetaDataSummaryList(ctx context.Context) (*api.MetaDataSummaryListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return gc.metadataService.GetMetaDataSummaryList(ctx, &emptypb.Empty{})
}

func (gc *GrpcClient) GetMetadataList(ctx context.Context, request *api.MetadataListRequest) (*api.MetadataListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 600 * time.Second)
	defer cancel()
	return gc.metadataService.GetMetadataList(ctx, request)
}

func (gc *GrpcClient) GetMetadataById(ctx context.Context, request *api.MetadataByIdRequest) (*api.MetadataByIdResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return gc.metadataService.GetMetadataById(ctx, request)
}

func (gc *GrpcClient) RevokeMetaData(ctx context.Context, request *api.RevokeMetaDataRequest) (*api.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return gc.metadataService.RevokeMetaData(ctx, request)
}

// ************************************** Resource module *******************************************************

func (gc *GrpcClient) SaveResource(ctx context.Context, request *api.PublishPowerRequest) (*api.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return gc.resourceService.PublishPower(ctx, request)
}

func (gc *GrpcClient) SyncPower(ctx context.Context, request *api.SyncPowerRequest) (*api.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return gc.resourceService.SyncPower(ctx, request)
}

func (gc *GrpcClient) RevokeResource(ctx context.Context, request *api.RevokePowerRequest) (*api.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return gc.resourceService.RevokePower(ctx, request)
}

func (gc *GrpcClient) GetPowerSummaryByIdentityId(ctx context.Context, request *api.PowerSummaryByIdentityRequest) (*api.PowerTotalSummaryResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return gc.resourceService.GetPowerSummaryByIdentityId(ctx, request)
}

func (gc *GrpcClient) GetPowerTotalSummaryList(ctx context.Context) (*api.PowerTotalSummaryListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	return gc.resourceService.GetPowerTotalSummaryList(ctx, &emptypb.Empty{})
}

func (gc *GrpcClient) GetPowerList(ctx context.Context, request *api.PowerListRequest) (*api.PowerListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return gc.resourceService.GetPowerList(ctx, request)
}

// ************************************** Identity module *******************************************************

func (gc *GrpcClient) SaveIdentity(ctx context.Context, request *api.SaveIdentityRequest) (*api.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return gc.identityService.SaveIdentity(ctx, request)
}

func (gc *GrpcClient) RevokeIdentityJoin(ctx context.Context, request *api.RevokeIdentityJoinRequest) (*api.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return gc.identityService.RevokeIdentityJoin(ctx, request)
}

func (gc *GrpcClient) GetIdentityList(ctx context.Context, request *api.IdentityListRequest) (*api.IdentityListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return gc.identityService.GetIdentityList(ctx, request)
}

// ************************************** Task module *******************************************************

func (gc *GrpcClient) SaveTask(ctx context.Context, request *api.TaskDetail) (*api.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return gc.taskService.SaveTask(ctx, request)
}

func (gc *GrpcClient) GetDetailTask(ctx context.Context, request *api.DetailTaskRequest) (*api.TaskDetail, error) {
	ctx, cancel := context.WithTimeout(ctx,2*time.Second)
	defer cancel()
	return gc.taskService.GetDetailTask(ctx, request)
}

func (gc *GrpcClient) ListTask(ctx context.Context, request *api.TaskListRequest) (*api.TaskListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return gc.taskService.ListTask(ctx, request)
}

func (gc *GrpcClient) ListTaskByIdentity(ctx context.Context, request *api.TaskListByIdentityRequest) (*api.TaskListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return gc.taskService.ListTaskByIdentity(ctx, request)
}

func (gc *GrpcClient) ListTaskEvent(ctx context.Context, request *api.TaskEventRequest) (*api.TaskEventResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return gc.taskService.ListTaskEvent(ctx, request)
}