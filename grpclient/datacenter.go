package grpclient

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/lib/center/api"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

const DefaultGrpcRequestTimeout = 2 * time.Second

// Client defines typed wrapper for the CenterData GRPC API.
type GrpcClient struct {
	c *grpc.ClientConn

	// grpc service
	metadataService api.MetadataServiceClient
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
		metadataService: api.NewMetadataServiceClient(conn),
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

// MetadataSave saves new metadata to database.
func (gc *GrpcClient) SaveMetadata(ctx context.Context, request *api.MetadataSaveRequest) (*apipb.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return gc.metadataService.MetadataSave(ctx, request)
}

func (gc *GrpcClient) GetMetadataSummaryList(ctx context.Context) (*api.MetadataSummaryListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return gc.metadataService.GetMetadataSummaryList(ctx, &emptypb.Empty{})
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

func (gc *GrpcClient) RevokeMetadata(ctx context.Context, request *api.RevokeMetadataRequest) (*apipb.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return gc.metadataService.RevokeMetadata(ctx, request)
}

// ************************************** Resource module *******************************************************

func (gc *GrpcClient) SaveResource(ctx context.Context, request *api.PublishPowerRequest) (*apipb.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return gc.resourceService.PublishPower(ctx, request)
}

func (gc *GrpcClient) SyncPower(ctx context.Context, request *api.SyncPowerRequest) (*apipb.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return gc.resourceService.SyncPower(ctx, request)
}

func (gc *GrpcClient) RevokeResource(ctx context.Context, request *api.RevokePowerRequest) (*apipb.SimpleResponse, error) {
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

func (gc *GrpcClient) SaveIdentity(ctx context.Context, request *api.SaveIdentityRequest) (*apipb.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.identityService.SaveIdentity(ctx, request)
}

func (gc *GrpcClient) RevokeIdentityJoin(ctx context.Context, request *api.RevokeIdentityJoinRequest) (*apipb.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.identityService.RevokeIdentityJoin(ctx, request)
}

func (gc *GrpcClient) GetIdentityList(ctx context.Context, request *api.IdentityListRequest) (*api.IdentityListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return gc.identityService.GetIdentityList(ctx, request)
}

// 存储元数据鉴权申请记录
func (gc *GrpcClient) SaveMetadataAuthority(ctx context.Context, request *api.SaveMetadataAuthorityRequest) (*apipb.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.identityService.SaveMetadataAuthority(ctx, request)
}

// 数据授权审核，规则：
// 1、授权后，可以将审核结果绑定到原有申请记录之上
func (gc *GrpcClient) AuditMetadataAuthority(ctx context.Context, request *api.AuditMetadataAuthorityRequest) (*apipb.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.identityService.AuditMetadataAuthority(ctx, request)
}

// 获取数据授权申请列表
// 规则：参数存在时根据条件获取，参数不存在时全量返回
func (gc *GrpcClient) GetMetadataAuthorityList(ctx context.Context, request *api.MetadataAuthorityListRequest) (*api.MetadataAuthorityListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.identityService.GetMetadataAuthorityList(ctx, request)
}

// ************************************** GetTask module *******************************************************

func (gc *GrpcClient) SaveTask(ctx context.Context, request *libtypes.TaskDetail) (*apipb.SimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	return gc.taskService.SaveTask(ctx, request)
}

func (gc *GrpcClient) GetDetailTask(ctx context.Context, request *api.DetailTaskRequest) (*libtypes.TaskDetail, error) {
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