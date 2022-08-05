package grpclient

import (
	"context"
	"fmt"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	datacenterapipb "github.com/datumtechs/datum-network-carrier/pb/datacenter/api"
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
	metadataService     datacenterapipb.MetadataServiceClient
	resourceService     datacenterapipb.ResourceServiceClient
	identityService     datacenterapipb.IdentityServiceClient
	metadataAuthService datacenterapipb.MetadataAuthServiceClient
	taskService         datacenterapipb.TaskServiceClient
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
		c:                   conn,
		metadataService:     datacenterapipb.NewMetadataServiceClient(conn),
		resourceService:     datacenterapipb.NewResourceServiceClient(conn),
		identityService:     datacenterapipb.NewIdentityServiceClient(conn),
		metadataAuthService: datacenterapipb.NewMetadataAuthServiceClient(conn),
		taskService:         datacenterapipb.NewTaskServiceClient(conn),
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

// ************************************** Metadata module *******************************************************

// MetadataSave saves new metadata to database.
func (gc *GrpcClient) SaveMetadata(ctx context.Context, request *datacenterapipb.SaveMetadataRequest) (*carriertypespb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.metadataService.SaveMetadata(ctx, request)
}

func (gc *GrpcClient) GetMetadataSummaryList(ctx context.Context, request *datacenterapipb.ListMetadataSummaryRequest) (*datacenterapipb.ListMetadataSummaryResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.metadataService.ListMetadataSummary(ctx, request, RPCMaxCallRecvMsgSize)
}

func (gc *GrpcClient) GetMetadataList(ctx context.Context, request *datacenterapipb.ListMetadataRequest) (*datacenterapipb.ListMetadataResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.metadataService.ListMetadata(ctx, request, RPCMaxCallRecvMsgSize)
}

func (gc *GrpcClient) GetMetadataListByIdentityId(ctx context.Context, request *datacenterapipb.ListMetadataByIdentityIdRequest) (*datacenterapipb.ListMetadataResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.metadataService.ListMetadataByIdentityId(ctx, request, RPCMaxCallRecvMsgSize)
}

func (gc *GrpcClient) GetMetadataById(ctx context.Context, request *datacenterapipb.FindMetadataByIdRequest) (*datacenterapipb.FindMetadataByIdResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.metadataService.FindMetadataById(ctx, request)
}

func (gc *GrpcClient) GetMetadataByIds(ctx context.Context, request *datacenterapipb.FindMetadataByIdsRequest) (*datacenterapipb.ListMetadataResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.metadataService.FindMetadataByIds(ctx, request)
}

func (gc *GrpcClient) RevokeMetadata(ctx context.Context, request *datacenterapipb.RevokeMetadataRequest) (*carriertypespb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.metadataService.RevokeMetadata(ctx, request)
}

// add by v 0.4.0
func (gc *GrpcClient) UpdateMetadata(ctx context.Context, request *datacenterapipb.UpdateMetadataRequest) (*carriertypespb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.metadataService.UpdateMetadata(ctx, request)
}

// ************************************** Resource (power) module *******************************************************

func (gc *GrpcClient) SaveResource(ctx context.Context, request *datacenterapipb.PublishPowerRequest) (*carriertypespb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.resourceService.PublishPower(ctx, request)
}

func (gc *GrpcClient) SyncPower(ctx context.Context, request *datacenterapipb.SyncPowerRequest) (*carriertypespb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.resourceService.SyncPower(ctx, request)
}

func (gc *GrpcClient) RevokeResource(ctx context.Context, request *datacenterapipb.RevokePowerRequest) (*carriertypespb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.resourceService.RevokePower(ctx, request)
}

func (gc *GrpcClient) GetPowerSummaryByIdentityId(ctx context.Context, request *datacenterapipb.GetPowerSummaryByIdentityRequest) (*datacenterapipb.PowerSummaryResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.resourceService.GetPowerSummaryByIdentityId(ctx, request, RPCMaxCallRecvMsgSize)
}

func (gc *GrpcClient) GetPowerGlobalSummaryList(ctx context.Context) (*datacenterapipb.ListPowerSummaryResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.resourceService.ListPowerSummary(ctx, &emptypb.Empty{}, RPCMaxCallRecvMsgSize)
}

func (gc *GrpcClient) GetPowerList(ctx context.Context, request *datacenterapipb.ListPowerRequest) (*datacenterapipb.ListPowerResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.resourceService.ListPower(ctx, request, RPCMaxCallRecvMsgSize)
}

// ************************************** Identity module *******************************************************

func (gc *GrpcClient) SaveIdentity(ctx context.Context, request *datacenterapipb.SaveIdentityRequest) (*carriertypespb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.identityService.SaveIdentity(ctx, request)
}

func (gc *GrpcClient) RevokeIdentityJoin(ctx context.Context, request *datacenterapipb.RevokeIdentityRequest) (*carriertypespb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.identityService.RevokeIdentity(ctx, request)
}

func (gc *GrpcClient) GetIdentityById(ctx context.Context, request *datacenterapipb.FindIdentityRequest) (*datacenterapipb.FindIdentityResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.identityService.FindIdentity(ctx, request)
}

func (gc *GrpcClient) GetIdentityList(ctx context.Context, request *datacenterapipb.ListIdentityRequest) (*datacenterapipb.ListIdentityResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.identityService.ListIdentity(ctx, request)
}

func (gc *GrpcClient) UpdateIdentityCredential(ctx context.Context, request *datacenterapipb.UpdateIdentityCredentialRequest) (*carriertypespb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.identityService.UpdateIdentityCredential(ctx, request)
}

// ************************************** MetadataAuth module *******************************************************

// Store metadata authentication application records
func (gc *GrpcClient) SaveMetadataAuthority(ctx context.Context, request *datacenterapipb.MetadataAuthorityRequest) (*carriertypespb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.metadataAuthService.SaveMetadataAuthority(ctx, request)
}

// Data authorization audit, rules:
// 1. After authorization, the approval result can be bound to the original application record
func (gc *GrpcClient) UpdateMetadataAuthority(ctx context.Context, request *datacenterapipb.MetadataAuthorityRequest) (*carriertypespb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.metadataAuthService.UpdateMetadataAuthority(ctx, request)
}

// Obtain data authorization application list
// Rule: obtained according to the condition when the parameter exists; returned in full when the parameter does not exist
func (gc *GrpcClient) GetMetadataAuthorityList(ctx context.Context, request *datacenterapipb.ListMetadataAuthorityRequest) (*datacenterapipb.ListMetadataAuthorityResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.metadataAuthService.ListMetadataAuthority(ctx, request, RPCMaxCallRecvMsgSize)
}

func (gc *GrpcClient) FindMetadataAuthority(ctx context.Context, request *datacenterapipb.FindMetadataAuthorityRequest) (*datacenterapipb.FindMetadataAuthorityResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.metadataAuthService.FindMetadataAuthority(ctx, request)
}

// ************************************** Task module *******************************************************

func (gc *GrpcClient) SaveTask(ctx context.Context, request *datacenterapipb.SaveTaskRequest) (*carriertypespb.SimpleResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.taskService.SaveTask(ctx, request)
}

func (gc *GrpcClient) GetDetailTask(ctx context.Context, request *datacenterapipb.GetTaskDetailRequest) (*datacenterapipb.GetTaskDetailResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.taskService.GetTaskDetail(ctx, request)
}

func (gc *GrpcClient) ListTask(ctx context.Context, request *datacenterapipb.ListTaskRequest) (*datacenterapipb.ListTaskResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, DefaultGrpcRequestTimeout)
	defer cancel()
	return gc.taskService.ListTask(ctx, request, RPCMaxCallRecvMsgSize)
}

func (gc *GrpcClient) ListTaskByIdentity(ctx context.Context, request *datacenterapipb.ListTaskByIdentityRequest) (*datacenterapipb.ListTaskResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.taskService.ListTaskByIdentity(ctx, request, RPCMaxCallRecvMsgSize)
}

func (gc *GrpcClient) ListTaskByTaskIds(ctx context.Context, request *datacenterapipb.ListTaskByTaskIdsRequest) (*datacenterapipb.ListTaskResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.taskService.ListTaskByTaskIds(ctx, request, RPCMaxCallRecvMsgSize)
}

func (gc *GrpcClient) ListTaskEvent(ctx context.Context, request *datacenterapipb.ListTaskEventRequest) (*datacenterapipb.ListTaskEventResponse, error) {
	if nil == gc {
		return nil, fmt.Errorf("datacenter rpc client is nil")
	}
	// TODO: Requests take too long, consider stream processing
	ctx, cancel := context.WithTimeout(ctx, TweentySecondGrpcRequestTimeout)
	defer cancel()
	return gc.taskService.ListTaskEvent(ctx, request, RPCMaxCallRecvMsgSize)
}
