package grpclient

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/lib/center/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

const (
	DefaultGrpcRequestTimeout          = 2 * time.Second // Default timeout time for a grpc request.
	DefaultGrpcDialTimeout             = 2 * time.Second // Default timeout time for Grpc's dial-up.
	TweentySecondGrpcRequestTimeout    = 20 * time.Second
	SixHundredSecondGrpcRequestTimeout = 600 * time.Second
)

// Client defines typed wrapper for the CenterData GRPC API.
type GrpcClient struct {
	c *grpc.ClientConn

	// grpc service
	metadataService api.MetadataServiceClient
	resourceService api.ResourceServiceClient
	identityService api.IdentityServiceClient
	taskService     api.TaskServiceClient
}

// NewClient creates a client that uses the given GRPC client.
func NewGrpcClient(ctx context.Context, addr string) (*GrpcClient, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcDialTimeout)
	_ = cancel
	conn, err := dialContext(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &GrpcClient{
		c:               conn,
		metadataService: api.NewMetadataServiceClient(conn),
		resourceService: api.NewResourceServiceClient(conn),
		identityService: api.NewIdentityServiceClient(conn),
		taskService:     api.NewTaskServiceClient(conn),
	}, nil
}

func (gc *GrpcClient) Close() {
	if nil != gc {
		gc.c.Close()
	}
}

func (gc *GrpcClient) GetClientConn() *grpc.ClientConn {
	if nil == gc {
		return nil
	}
	return gc.c
}

// MetadataSave saves new metadata to database.
func (gc *GrpcClient) SaveMetadata(ctx context.Context, request *api.SaveMetadataRequest) (*apicommonpb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.metadataService.SaveMetadata(ctx, request)
}

func (gc *GrpcClient) GetMetadataSummaryList(ctx context.Context, request *api.ListMetadataSummaryRequest) (*api.ListMetadataSummaryResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.metadataService.ListMetadataSummary(ctx, request, RPCMaxCallRecvMsgSize)
}

func (gc *GrpcClient) GetMetadataList(ctx context.Context, request *api.ListMetadataRequest) (*api.ListMetadataResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.metadataService.ListMetadata(ctx, request, RPCMaxCallRecvMsgSize)
}

func (gc *GrpcClient) GetMetadataListByIdentityId(ctx context.Context, request *api.ListMetadataByIdentityIdRequest) (*api.ListMetadataResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.metadataService.ListMetadataByIdentityId(ctx, request, RPCMaxCallRecvMsgSize)
}

func (gc *GrpcClient) GetMetadataById(ctx context.Context, request *api.FindMetadataByIdRequest) (*api.FindMetadataByIdResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.metadataService.FindMetadataById(ctx, request)
}

func (gc *GrpcClient) RevokeMetadata(ctx context.Context, request *api.RevokeMetadataRequest) (*apicommonpb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.metadataService.RevokeMetadata(ctx, request)
}

// ************************************** Resource module *******************************************************

func (gc *GrpcClient) SaveResource(ctx context.Context, request *api.PublishPowerRequest) (*apicommonpb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.resourceService.PublishPower(ctx, request)
}

func (gc *GrpcClient) SyncPower(ctx context.Context, request *api.SyncPowerRequest) (*apicommonpb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.resourceService.SyncPower(ctx, request)
}

func (gc *GrpcClient) RevokeResource(ctx context.Context, request *api.RevokePowerRequest) (*apicommonpb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.resourceService.RevokePower(ctx, request)
}

func (gc *GrpcClient) GetPowerSummaryByIdentityId(ctx context.Context, request *api.GetPowerSummaryByIdentityRequest) (*api.PowerSummaryResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.resourceService.GetPowerSummaryByIdentityId(ctx, request, RPCMaxCallRecvMsgSize)
}

func (gc *GrpcClient) GetPowerGlobalSummaryList(ctx context.Context) (*api.ListPowerSummaryResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.resourceService.ListPowerSummary(ctx, &emptypb.Empty{}, RPCMaxCallRecvMsgSize)
}

func (gc *GrpcClient) GetPowerList(ctx context.Context, request *api.ListPowerRequest) (*api.ListPowerResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.resourceService.ListPower(ctx, request, RPCMaxCallRecvMsgSize)
}

// ************************************** Identity module *******************************************************

func (gc *GrpcClient) SaveIdentity(ctx context.Context, request *api.SaveIdentityRequest) (*apicommonpb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.identityService.SaveIdentity(ctx, request)
}

func (gc *GrpcClient) RevokeIdentityJoin(ctx context.Context, request *api.RevokeIdentityRequest) (*apicommonpb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.identityService.RevokeIdentity(ctx, request)
}

func (gc *GrpcClient) GetIdentityList(ctx context.Context, request *api.ListIdentityRequest) (*api.ListIdentityResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.identityService.ListIdentity(ctx, request)
}

// Store metadata authentication application records
func (gc *GrpcClient) SaveMetadataAuthority(ctx context.Context, request *api.MetadataAuthorityRequest) (*apicommonpb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.identityService.SaveMetadataAuthority(ctx, request)
}

// Data authorization audit, rules:
// 1. After authorization, the approval result can be bound to the original application record
func (gc *GrpcClient) UpdateMetadataAuthority(ctx context.Context, request *api.MetadataAuthorityRequest) (*apicommonpb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.identityService.UpdateMetadataAuthority(ctx, request)
}

// Obtain data authorization application list
// Rule: obtained according to the condition when the parameter exists; returned in full when the parameter does not exist
func (gc *GrpcClient) GetMetadataAuthorityList(ctx context.Context, request *api.ListMetadataAuthorityRequest) (*api.ListMetadataAuthorityResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.identityService.ListMetadataAuthority(ctx, request, RPCMaxCallRecvMsgSize)
}

func (gc *GrpcClient) FindMetadataAuthority(ctx context.Context, request *api.FindMetadataAuthorityRequest) (*api.FindMetadataAuthorityResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.identityService.FindMetadataAuthority(ctx, request)
}

// ************************************** GetTask module *******************************************************

func (gc *GrpcClient) SaveTask(ctx context.Context, request *api.SaveTaskRequest) (*apicommonpb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.taskService.SaveTask(ctx, request)
}

func (gc *GrpcClient) GetDetailTask(ctx context.Context, request *api.GetTaskDetailRequest) (*api.GetTaskDetailResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.taskService.GetTaskDetail(ctx, request)
}

func (gc *GrpcClient) ListTask(ctx context.Context, request *api.ListTaskRequest) (*api.ListTaskResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.taskService.ListTask(ctx, request, RPCMaxCallRecvMsgSize)
}

func (gc *GrpcClient) ListTaskByIdentity(ctx context.Context, request *api.ListTaskByIdentityRequest) (*api.ListTaskResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.taskService.ListTaskByIdentity(ctx, request, RPCMaxCallRecvMsgSize)
}

func (gc *GrpcClient) ListTaskByTaskIds(ctx context.Context, request *api.ListTaskByTaskIdsRequest) (*api.ListTaskResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.taskService.ListTaskByTaskIds(ctx, request, RPCMaxCallRecvMsgSize)
}

func (gc *GrpcClient) ListTaskEvent(ctx context.Context, request *api.ListTaskEventRequest) (*api.ListTaskEventResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.taskService.ListTaskEvent(ctx, request, RPCMaxCallRecvMsgSize)
}
