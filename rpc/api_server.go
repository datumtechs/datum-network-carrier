package rpc

import (
	"context"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
)

type yarnServiceServer struct {
	pb.UnimplementedYarnServiceServer
	b 			Backend
}

func(svr *yarnServiceServer) GetNodeInfo (ctx context.Context, req *pb.EmptyGetParams) (*pb.GetNodeInfoResponse, error) {

	return nil, nil
}
func (svr *yarnServiceServer) GetRegisteredPeers(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetRegisteredPeersResponse, error) {

	return nil, nil
}
func (svr *yarnServiceServer) SetSeedNode(ctx context.Context, req *pb.SetSeedNodeRequest) (*pb.SetSeedNodeResponse, error) {
	return nil, nil
}
func (svr *yarnServiceServer) UpdateSeedNode(ctx context.Context, req *pb.UpdateSeedNodeRequest) (*pb.SetSeedNodeResponse, error) {
	return nil, nil
}
func (svr *yarnServiceServer) DeleteSeedNode(ctx context.Context, req *pb.DeleteRegisteredNodeRequest) (*pb.SimpleResponseCode, error) {
	return nil, nil
}
func (svr *yarnServiceServer) GetSeedNodeList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetSeedNodeListResponse, error) {
	return nil, nil
}
func (svr *yarnServiceServer) SetDataNode(ctx context.Context, req *pb.SetDataNodeRequest) (*pb.SetDataNodeResponse, error) {
	return nil, nil
}
func (svr *yarnServiceServer) UpdateDataNode(ctx context.Context, req *pb.UpdateDataNodeRequest) (*pb.SetDataNodeResponse, error) {
	return nil, nil
}
func (svr *yarnServiceServer) DeleteDataNode(ctx context.Context, req *pb.DeleteRegisteredNodeRequest) (*pb.SimpleResponseCode, error) {
	return nil, nil
}
func (svr *yarnServiceServer) GetDataNodeList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetRegisteredNodeListResponse, error) {
	return nil, nil
}
func (svr *yarnServiceServer) SetJobNode(ctx context.Context, req *pb.SetJobNodeRequest) (*pb.SetJobNodeResponse, error) {
	return nil, nil
}
func (svr *yarnServiceServer) UpdateJobNode(ctx context.Context, req *pb.UpdateJobNodeRequest) (*pb.SetJobNodeResponse, error) {
	return nil, nil
}
func (svr *yarnServiceServer) DeleteJobNode(ctx context.Context, req *pb.DeleteRegisteredNodeRequest) (*pb.SimpleResponseCode, error) {
	return nil, nil
}
func (svr *yarnServiceServer) GetJobNodeList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetRegisteredNodeListResponse, error) {
	return nil, nil
}
func (svr *yarnServiceServer) ReportTaskEvent(ctx context.Context, req *pb.ReportTaskEventRequest) (*pb.SimpleResponseCode, error) {
	return nil, nil
}
func (svr *yarnServiceServer) ReportTaskResourceExpense(ctx context.Context, req *pb.ReportTaskResourceExpenseRequest) (*pb.SimpleResponseCode, error) {
	return nil, nil
}


type metaDataServiceServer struct {
	pb.UnimplementedMetaDataServiceServer
	b 			Backend
}

func (svr *metaDataServiceServer) GetMetaDataSummaryList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetMetaDataSummaryListResponse, error) {
	return nil, nil
}
func (svr *metaDataServiceServer) GetMetaDataSummaryByState(ctx context.Context, req *pb.GetMetaDataSummaryByStateRequest) (*pb.GetMetaDataSummaryListResponse, error) {
	return nil, nil
}
func (svr *metaDataServiceServer) GetMetaDataSummaryByOwner(ctx context.Context, req *pb.GetMetaDataSummaryByOwnerRequest) (*pb.GetMetaDataSummaryListResponse, error) {
	return nil, nil
}
func (svr *metaDataServiceServer) GetMetaDataDetail(ctx context.Context, req *pb.GetMetaDataDetailRequest) (*pb.GetMetaDataDetailResponse, error) {
	return nil, nil
}
func (svr *metaDataServiceServer) PublishMetaData(ctx context.Context, req *pb.PublishMetaDataRequest) (*pb.PublishMetaDataResponse, error) {
	return nil, nil
}
func (svr *metaDataServiceServer) RevokeMetaData(ctx context.Context, req *pb.RevokeMetaDataRequest) (*pb.SimpleResponseCode, error) {
	return nil, nil
}


type powerServiceServer struct {
	pb.UnimplementedPowerServiceServer
	b 			Backend
}
func (svr *powerServiceServer) GetPowerTotalSummaryList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetPowerTotalSummaryListResponse, error) {
	return nil, nil
}
func (svr *powerServiceServer) GetPowerSingleSummaryList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetPowerSingleSummaryListResponse, error) {
	return nil, nil
}
func (svr *powerServiceServer) GetPowerTotalSummaryByState(ctx context.Context, req *pb.GetPowerTotalSummaryByStateRequest) (*pb.GetPowerTotalSummaryListResponse, error) {
	return nil, nil
}
func (svr *powerServiceServer) GetPowerSingleSummaryByState(ctx context.Context, req *pb.GetPowerSingleSummaryByStateRequest) (*pb.GetPowerSingleSummaryListResponse, error) {
	return nil, nil
}
func (svr *powerServiceServer) GetPowerTotalSummaryByOwner(ctx context.Context, req *pb.GetPowerTotalSummaryByOwnerRequest) (*pb.GetPowerTotalSummaryResponse, error) {
	return nil, nil
}
func (svr *powerServiceServer) GetPowerSingleSummaryByOwner(ctx context.Context, req *pb.GetPowerSingleSummaryByOwnerRequest) (*pb.GetPowerSingleSummaryListResponse, error) {
	return nil, nil
}
func (svr *powerServiceServer) GetPowerSingleDetail(ctx context.Context, req *pb.GetPowerSingleDetailRequest) (*pb.GetPowerSingleDetailResponse, error) {
	return nil, nil
}
func (svr *powerServiceServer) PublishPower(ctx context.Context, req *pb.PublishPowerRequest) (*pb.PublishPowerResponse, error) {
	return nil, nil
}
func (svr *powerServiceServer) RevokePower(ctx context.Context, req *pb.RevokePowerRequest) (*pb.SimpleResponseCode, error) {
	return nil, nil
}


type authServiceServer struct {
	pb.UnimplementedAuthServiceServer
	b 			Backend
}
func (svr *authServiceServer) ApplyIdentityJoin(ctx context.Context, req *pb.ApplyIdentityJoinRequest) (*pb.SimpleResponseCode, error) {
	return nil, nil
}
func (svr *authServiceServer) RevokeIdentityJoin(ctx context.Context, req *pb.RevokeIdentityJoinRequest) (*pb.SimpleResponseCode, error) {
	return nil, nil
}



type taskServiceServer struct {
	pb.UnimplementedTaskServiceServer
	b 			Backend
}
func (svr *taskServiceServer) GetTaskSummaryList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetTaskSummaryListResponse, error) {
	return nil, nil
}
func (svr *taskServiceServer) GetTaskJoinSummaryList(ctx context.Context, req *pb.GetTaskJoinSummaryListRequest) (*pb.GetTaskJoinSummaryListResponse, error) {
	return nil, nil
}
func (svr *taskServiceServer) GetTaskDetail(ctx context.Context, req *pb.GetTaskDetailRequest) (*pb.GetTaskDetailResponse, error) {
	return nil, nil
}
func (svr *taskServiceServer) GetTaskEventList(ctx context.Context, req *pb.GetTaskEventListRequest) (*pb.GetTaskEventListResponse, error) {
	return nil, nil
}
func (svr *taskServiceServer) PublishTaskDeclare(ctx context.Context, req *pb.PublishTaskDeclareRequest) (*pb.PublishTaskDeclareResponse, error) {
	return nil, nil
}