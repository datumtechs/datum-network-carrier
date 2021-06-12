package rpc

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/event"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"github.com/RosettaFlow/Carrier-Go/types"
	"time"
)

const (
	OK = "ok"
)

var (
	ErrSetSeedNodeInfoStr = "Failed to set seed node info"
	ErrDeleteSeedNodeInfoStr = "Failed to delete seed node info"
	ErrGetSeedNodeListStr = "Failed to get seed nodes"
	ErrSetDataNodeInfoStr = "Failed to set data node info"
	ErrDeleteDataNodeInfoStr = "Failed to delete data node info"
	ErrGetDataNodeListStr = "Failed to get data nodes"
	ErrSetJobNodeInfoStr = "Failed to set job node info"
	ErrDeleteJobNodeInfoStr = "Failed to delete job node info"
	ErrSendPowerMsgStr = "Failed to send powerMsg"

	ErrReportTaskEventStr = "Failed to report taskEvent"
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

	seedNode := &types.SeedNodeInfo{
		InternalIp: req.InternalIp,
		InternalPort: req.InternalPort,
		ConnState: types.NONCONNECTED,
	}
	seedNode.SeedNodeId()
	_, err := svr.b.SetSeedNode(seedNode)
	if nil != err {
		return nil, NewRpcBizErr(ErrSetSeedNodeInfoStr)
	}
	return &pb.SetSeedNodeResponse{
		Status: 0,
		Msg: OK,
		SeedPeer: &pb.SeedPeer{
			Id: seedNode.Id,
			InternalIp:  seedNode.InternalIp,
			InternalPort: seedNode.InternalPort,
			ConnState: seedNode.ConnState.Int32(),
		},
	}, nil
}
func (svr *yarnServiceServer) UpdateSeedNode(ctx context.Context, req *pb.UpdateSeedNodeRequest) (*pb.SetSeedNodeResponse, error) {

	seedNode := &types.SeedNodeInfo{
		Id: req.Id,
		InternalIp: req.InternalIp,
		InternalPort: req.InternalPort,
		ConnState: types.NONCONNECTED,
	}
	_, err := svr.b.SetSeedNode(seedNode)
	if nil != err {
		return nil, NewRpcBizErr(ErrSetSeedNodeInfoStr)
	}
	return &pb.SetSeedNodeResponse{
		Status: 0,
		Msg: OK,
		SeedPeer: &pb.SeedPeer{
			Id: seedNode.Id,
			InternalIp:  seedNode.InternalIp,
			InternalPort: seedNode.InternalPort,
			ConnState: seedNode.ConnState.Int32(),
		},
	}, nil

	return nil, nil
}
func (svr *yarnServiceServer) DeleteSeedNode(ctx context.Context, req *pb.DeleteRegisteredNodeRequest) (*pb.SimpleResponseCode, error) {
	err := svr.b.DeleteSeedNode(req.Id)
	if nil != err {
		return nil, NewRpcBizErr(ErrDeleteSeedNodeInfoStr)
	}
	return &pb.SimpleResponseCode{Status: 0, Msg:  OK}, nil
}
func (svr *yarnServiceServer) GetSeedNodeList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetSeedNodeListResponse, error) {
	list, err := svr.b.GetSeedNodeList()
	if nil != err {
		return nil, NewRpcBizErr(ErrGetSeedNodeListStr)
	}
	seeds := make([]*pb.SeedPeer, len(list))
	for i, v := range list {
		s := &pb.SeedPeer{
			Id: v.Id,
			InternalIp:  v.InternalIp,
			InternalPort: v.InternalPort,
			ConnState: v.ConnState.Int32(),
		}
		seeds[i] = s
	}
	return &pb.GetSeedNodeListResponse{
		Status: 0, 
		Msg: OK,
		SeedPeers: seeds,
	}, nil
}
func (svr *yarnServiceServer) SetDataNode(ctx context.Context, req *pb.SetDataNodeRequest) (*pb.SetDataNodeResponse, error) {

	node := &types.RegisteredNodeInfo{
		InternalIp: req.InternalIp,
		InternalPort: req.InternalPort,
		ExternalIp: req.ExternalIp,
		ExternalPort: req.InternalPort,
		ConnState: types.NONCONNECTED,
	}
	node.DataNodeId()
	_, err := svr.b.SetRegisterNode(types.PREFIX_TYPE_DATANODE, node)
	if nil != err {
		return nil, NewRpcBizErr(ErrSetDataNodeInfoStr)
	}
	return &pb.SetDataNodeResponse{
		Status: 0,
		Msg: OK,
		DataNode: &pb.YarnRegisteredPeerDetail{
			Id: node.Id,
			InternalIp: node.InternalIp,
			InternalPort: node.InternalPort,
			ExternalIp: node.ExternalIp,
			ExternalPort: node.InternalPort,
			ConnState: node.ConnState.Int32(),
		},
	}, nil
}
func (svr *yarnServiceServer) UpdateDataNode(ctx context.Context, req *pb.UpdateDataNodeRequest) (*pb.SetDataNodeResponse, error) {
	node := &types.RegisteredNodeInfo{
		Id: req.Id,
		InternalIp: req.InternalIp,
		InternalPort: req.InternalPort,
		ExternalIp: req.ExternalIp,
		ExternalPort: req.InternalPort,
		ConnState: types.NONCONNECTED,
	}

	_, err := svr.b.SetRegisterNode(types.PREFIX_TYPE_DATANODE, node)
	if nil != err {
		return nil, NewRpcBizErr(ErrSetDataNodeInfoStr)
	}
	return &pb.SetDataNodeResponse{
		Status: 0,
		Msg: OK,
		DataNode: &pb.YarnRegisteredPeerDetail{
			Id: node.Id,
			InternalIp: node.InternalIp,
			InternalPort: node.InternalPort,
			ExternalIp: node.ExternalIp,
			ExternalPort: node.InternalPort,
			ConnState: node.ConnState.Int32(),
		},
	}, nil
}
func (svr *yarnServiceServer) DeleteDataNode(ctx context.Context, req *pb.DeleteRegisteredNodeRequest) (*pb.SimpleResponseCode, error) {

	if err := svr.b.DeleteRegisterNode(types.PREFIX_TYPE_DATANODE, req.Id); nil != err {
		return nil, NewRpcBizErr(ErrDeleteDataNodeInfoStr)
	}
	return &pb.SimpleResponseCode{Status: 0, Msg:  OK}, nil
}
func (svr *yarnServiceServer) GetDataNodeList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetRegisteredNodeListResponse, error) {

	list, err := svr.b.GetRegisterNodeList(types.PREFIX_TYPE_DATANODE)
	if nil != err {
		return nil, NewRpcBizErr(ErrGetDataNodeListStr)
	}
	datas := make([]*pb.YarnRegisteredPeer, len(list))
	for i, v := range list {
		d := &pb.YarnRegisteredPeer{
			NodeType: types.PREFIX_TYPE_DATANODE.String(),
			NodeDetail: &pb.YarnRegisteredPeerDetail{
				InternalIp: v.InternalIp,
				InternalPort: v.InternalPort,
				ExternalIp: v.ExternalIp,
				ExternalPort: v.InternalPort,
				ConnState: v.ConnState.Int32(),
			},

		}
		datas[i] = d
	}
	return &pb.GetRegisteredNodeListResponse{
		Status: 0,
		Msg: OK,
		Nodes: datas,
	}, nil
}
func (svr *yarnServiceServer) SetJobNode(ctx context.Context, req *pb.SetJobNodeRequest) (*pb.SetJobNodeResponse, error) {
	node := &types.RegisteredNodeInfo{
		InternalIp: req.InternalIp,
		InternalPort: req.InternalPort,
		ExternalIp: req.ExternalIp,
		ExternalPort: req.InternalPort,
		ConnState: types.NONCONNECTED,
	}
	node.DataNodeId()
	_, err := svr.b.SetRegisterNode(types.PREFIX_TYPE_JOBNODE, node)
	if nil != err {
		return nil, NewRpcBizErr(ErrSetJobNodeInfoStr)
	}
	return &pb.SetJobNodeResponse{
		Status: 0,
		Msg: OK,
		JobNode: &pb.YarnRegisteredPeerDetail{
			Id: node.Id,
			InternalIp: node.InternalIp,
			InternalPort: node.InternalPort,
			ExternalIp: node.ExternalIp,
			ExternalPort: node.InternalPort,
			ConnState: node.ConnState.Int32(),
		},
	}, nil
}
func (svr *yarnServiceServer) UpdateJobNode(ctx context.Context, req *pb.UpdateJobNodeRequest) (*pb.SetJobNodeResponse, error) {
	node := &types.RegisteredNodeInfo{
		Id: req.Id,
		InternalIp: req.InternalIp,
		InternalPort: req.InternalPort,
		ExternalIp: req.ExternalIp,
		ExternalPort: req.InternalPort,
		ConnState: types.NONCONNECTED,
	}
	_, err := svr.b.SetRegisterNode(types.PREFIX_TYPE_JOBNODE, node)
	if nil != err {
		return nil, NewRpcBizErr(ErrSetJobNodeInfoStr)
	}
	return &pb.SetJobNodeResponse{
		Status: 0,
		Msg: OK,
		JobNode: &pb.YarnRegisteredPeerDetail{
			Id: node.Id,
			InternalIp: node.InternalIp,
			InternalPort: node.InternalPort,
			ExternalIp: node.ExternalIp,
			ExternalPort: node.InternalPort,
			ConnState: node.ConnState.Int32(),
		},
	}, nil
}
func (svr *yarnServiceServer) DeleteJobNode(ctx context.Context, req *pb.DeleteRegisteredNodeRequest) (*pb.SimpleResponseCode, error) {
	if err := svr.b.DeleteRegisterNode(types.PREFIX_TYPE_JOBNODE, req.Id); nil != err {
		return nil, NewRpcBizErr(ErrDeleteJobNodeInfoStr)
	}
	return &pb.SimpleResponseCode{Status: 0, Msg:  OK}, nil
}
func (svr *yarnServiceServer) GetJobNodeList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetRegisteredNodeListResponse, error) {
	list, err := svr.b.GetRegisterNodeList(types.PREFIX_TYPE_JOBNODE)
	if nil != err {
		return nil, NewRpcBizErr(ErrGetDataNodeListStr)
	}
	jobs := make([]*pb.YarnRegisteredPeer, len(list))
	for i, v := range list {
		d := &pb.YarnRegisteredPeer{
			NodeType: types.PREFIX_TYPE_JOBNODE.String(),
			NodeDetail: &pb.YarnRegisteredPeerDetail{
				InternalIp: v.InternalIp,
				InternalPort: v.InternalPort,
				ExternalIp: v.ExternalIp,
				ExternalPort: v.InternalPort,
				ConnState: v.ConnState.Int32(),
			},

		}
		jobs[i] = d
	}
	return &pb.GetRegisteredNodeListResponse{
		Status: 0,
		Msg: OK,
		Nodes: jobs,
	}, nil
}
func (svr *yarnServiceServer) ReportTaskEvent(ctx context.Context, req *pb.ReportTaskEventRequest) (*pb.SimpleResponseCode, error) {
	var err error
	go func() {
		err = svr.b.SendTaskEvent(&event.TaskEvent{
			Type:       req.TaskEvent.Type,
			Identity:   req.TaskEvent.IdentityId,
			TaskId:     req.TaskEvent.TaskId,
			Content:    req.TaskEvent.Content,
			CreateTime: req.TaskEvent.CreateAt,
		})
	}()
	if nil != err {
		return nil, NewRpcBizErr(ErrReportTaskEventStr)
	}
	return &pb.SimpleResponseCode{Status: 0, Msg:  OK}, nil
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

	powerMsg := new(types.PowerMsg)
	powerMsg.Data.JobNodeId = req.JobNodeId
	powerMsg.Data.CreateAt = uint64(time.Now().UnixNano())
	powerMsg.Data.Name = req.Owner.Name
	powerMsg.Data.NodeId = req.Owner.NodeId
	powerMsg.Data.IdentityId = req.Owner.IdentityId
	powerMsg.Data.Information.Processor = req.Information.Processor
	powerMsg.Data.Information.Mem = req.Information.Mem
	powerMsg.Data.Information.Bandwidth = req.Information.Bandwidth
	powerId := powerMsg.GetPowerId()

	err := svr.b.SendMsg(powerMsg)
	if nil != err {
		return nil, NewRpcBizErr(ErrSendPowerMsgStr)
	}
	return &pb.PublishPowerResponse{
		Status: 0,
		Msg: OK,
		PowerId: powerId,
	}, nil
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