package backend

import (
	"context"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func (svr *YarnServiceServer) GetNodeInfo(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetNodeInfoResponse, error) {
	node, err := svr.B.GetNodeInfo()
	if nil != err {
		return nil, NewRpcBizErr(ErrGetNodeInfoStr)
	}
	peers := make([]*pb.YarnRegisteredPeer, len(node.Peers))
	for i, v := range node.Peers {
		n := &pb.YarnRegisteredPeer{
			NodeType: v.NodeType,
			NodeDetail: &pb.YarnRegisteredPeerDetail{
				Id:           v.RegisteredNodeInfo.Id,
				InternalIp:   v.InternalIp,
				ExternalIp:   v.ExternalIp,
				InternalPort: v.InternalPort,
				ExternalPort: v.ExternalPort,
				ConnState:    v.ConnState.Int32(),
			},
		}
		peers[i] = n
	}
	seeds := make([]*pb.SeedPeer, len(node.SeedPeers))
	for i, v := range node.SeedPeers {
		n := &pb.SeedPeer{
			Id:           v.Id,
			InternalIp:   v.InternalIp,
			InternalPort: v.InternalPort,
			ConnState:    v.ConnState.Int32(),
		}
		seeds[i] = n
	}

	return &pb.GetNodeInfoResponse{
		Status: 0,
		Msg:    OK,
		Information: &pb.YarnNodeInfo{
			NodeType:     node.NodeType,
			NodeId:       node.NodeId,
			InternalIp:   node.InternalIp,
			ExternalIp:   node.ExternalIp,
			InternalPort: node.InternalPort,
			ExternalPort: node.ExternalPort,
			IdentityType: node.IdentityType,
			IdentityId:   node.IdentityId,
			State:        node.State,
			Name:         node.Name,
			SeedPeers:    seeds,
			Peers:        peers,
		},
	}, nil
}

func (svr *YarnServiceServer) GetRegisteredPeers(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetRegisteredPeersResponse, error) {
	registerNodes, err := svr.B.GetRegisteredPeers()
	if nil != err {
		return nil, NewRpcBizErr(ErrGetRegisteredPeersStr)
	}
	jobNodes := make([]*pb.YarnRegisteredJobNode, len(registerNodes.JobNodes))
	for i, v := range registerNodes.JobNodes {
		node := &pb.YarnRegisteredJobNode{
			Id:           v.Id,
			InternalIp:   v.InternalIp,
			ExternalIp:   v.ExternalIp,
			InternalPort: v.InternalPort,
			ExternalPort: v.ExternalPort,
			Information: &pb.ResourceUsedDetailShow{
				TotalMem:       v.ResourceUsage.TotalMem,
				UsedMem:        v.ResourceUsage.UsedMem,
				TotalProcessor: v.ResourceUsage.TotalProcessor,
				UsedProcessor:  v.ResourceUsage.UsedProcessor,
				TotalBandwidth: v.ResourceUsage.TotalBandwidth,
				UsedBandwidth:  v.ResourceUsage.UsedBandwidth,
			},
			Duration: v.Duration,
			Task: &pb.YarnRegisteredJobNodeTaskIds{
				Count:   v.Task.Count,
				TaskIds: v.Task.TaskIds,
			},
		}
		jobNodes[i] = node
	}

	dataNodes := make([]*pb.YarnRegisteredDataNode, len(registerNodes.DataNodes))
	for i, v := range registerNodes.DataNodes {
		node := &pb.YarnRegisteredDataNode{
			Id:           v.Id,
			InternalIp:   v.InternalIp,
			ExternalIp:   v.ExternalIp,
			InternalPort: v.InternalPort,
			ExternalPort: v.ExternalPort,
			Information: &pb.ResourceUsedDetailShow{
				TotalMem:       v.ResourceUsage.TotalMem,
				UsedMem:        v.ResourceUsage.UsedMem,
				TotalProcessor: v.ResourceUsage.TotalProcessor,
				UsedProcessor:  v.ResourceUsage.UsedProcessor,
				TotalBandwidth: v.ResourceUsage.TotalBandwidth,
				UsedBandwidth:  v.ResourceUsage.UsedBandwidth,
			},
			Duration: v.Duration,
			Delta: &pb.YarnRegisteredDataNodeDelta{
				FileCount:     v.Delta.FileCount,
				FileTotalSize: v.Delta.FileTotalSize,
			},
		}
		dataNodes[i] = node
	}
	return &pb.GetRegisteredPeersResponse{
		Status:    0,
		Msg:       OK,
		JobNodes:  jobNodes,
		DataNodes: dataNodes,
	}, nil
}

func (svr *YarnServiceServer) SetSeedNode(ctx context.Context, req *pb.SetSeedNodeRequest) (*pb.SetSeedNodeResponse, error) {

	seedNode := &types.SeedNodeInfo{
		InternalIp:   req.InternalIp,
		InternalPort: req.InternalPort,
		ConnState:    types.NONCONNECTED,
	}
	seedNode.SeedNodeId()
	_, err := svr.B.SetSeedNode(seedNode) // TODO 未完成 ...
	if nil != err {
		return nil, NewRpcBizErr(ErrSetSeedNodeInfoStr)
	}
	return &pb.SetSeedNodeResponse{
		Status: 0,
		Msg:    OK,
		SeedPeer: &pb.SeedPeer{
			Id:           seedNode.Id,
			InternalIp:   seedNode.InternalIp,
			InternalPort: seedNode.InternalPort,
			ConnState:    seedNode.ConnState.Int32(),
		},
	}, nil
}

func (svr *YarnServiceServer) UpdateSeedNode(ctx context.Context, req *pb.UpdateSeedNodeRequest) (*pb.SetSeedNodeResponse, error) {

	seedNode := &types.SeedNodeInfo{
		Id:           req.Id,
		InternalIp:   req.InternalIp,
		InternalPort: req.InternalPort,
		ConnState:    types.NONCONNECTED,
	}
	_, err := svr.B.SetSeedNode(seedNode) // TODO 未完成 ...
	if nil != err {
		return nil, NewRpcBizErr(ErrSetSeedNodeInfoStr)
	}
	return &pb.SetSeedNodeResponse{
		Status: 0,
		Msg:    OK,
		SeedPeer: &pb.SeedPeer{
			Id:           seedNode.Id,
			InternalIp:   seedNode.InternalIp,
			InternalPort: seedNode.InternalPort,
			ConnState:    seedNode.ConnState.Int32(),
		},
	}, nil

	return nil, nil
}

func (svr *YarnServiceServer) DeleteSeedNode(ctx context.Context, req *pb.DeleteRegisteredNodeRequest) (*pb.SimpleResponseCode, error) {
	err := svr.B.DeleteSeedNode(req.Id)
	if nil != err {
		return nil, NewRpcBizErr(ErrDeleteSeedNodeInfoStr)
	}
	return &pb.SimpleResponseCode{Status: 0, Msg: OK}, nil
}

func (svr *YarnServiceServer) GetSeedNodeList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetSeedNodeListResponse, error) {
	list, err := svr.B.GetSeedNodeList()
	if nil != err {
		return nil, NewRpcBizErr(ErrGetSeedNodeListStr)
	}
	seeds := make([]*pb.SeedPeer, len(list))
	for i, v := range list {
		s := &pb.SeedPeer{
			Id:           v.Id,
			InternalIp:   v.InternalIp,
			InternalPort: v.InternalPort,
			ConnState:    v.ConnState.Int32(),
		}
		seeds[i] = s
	}
	return &pb.GetSeedNodeListResponse{
		Status:    0,
		Msg:       OK,
		SeedPeers: seeds,
	}, nil
}

func (svr *YarnServiceServer) SetDataNode(ctx context.Context, req *pb.SetDataNodeRequest) (*pb.SetDataNodeResponse, error) {

	node := &types.RegisteredNodeInfo{
		InternalIp:   req.InternalIp,
		InternalPort: req.InternalPort,
		ExternalIp:   req.ExternalIp,
		ExternalPort: req.InternalPort,
		ConnState:    types.NONCONNECTED,
	}
	node.DataNodeId()
	_, err := svr.B.SetRegisterNode(types.PREFIX_TYPE_DATANODE, node) // TODO 未完成 ...
	if nil != err {
		return nil, NewRpcBizErr(ErrSetDataNodeInfoStr)
	}
	return &pb.SetDataNodeResponse{
		Status: 0,
		Msg:    OK,
		DataNode: &pb.YarnRegisteredPeerDetail{
			Id:           node.Id,
			InternalIp:   node.InternalIp,
			InternalPort: node.InternalPort,
			ExternalIp:   node.ExternalIp,
			ExternalPort: node.InternalPort,
			ConnState:    node.ConnState.Int32(),
		},
	}, nil
}

func (svr *YarnServiceServer) UpdateDataNode(ctx context.Context, req *pb.UpdateDataNodeRequest) (*pb.SetDataNodeResponse, error) {
	node := &types.RegisteredNodeInfo{
		Id:           req.Id,
		InternalIp:   req.InternalIp,
		InternalPort: req.InternalPort,
		ExternalIp:   req.ExternalIp,
		ExternalPort: req.InternalPort,
		ConnState:    types.NONCONNECTED,
	}

	_, err := svr.B.SetRegisterNode(types.PREFIX_TYPE_DATANODE, node) // TODO 未完成 ...
	if nil != err {
		return nil, NewRpcBizErr(ErrSetDataNodeInfoStr)
	}
	return &pb.SetDataNodeResponse{
		Status: 0,
		Msg:    OK,
		DataNode: &pb.YarnRegisteredPeerDetail{
			Id:           node.Id,
			InternalIp:   node.InternalIp,
			InternalPort: node.InternalPort,
			ExternalIp:   node.ExternalIp,
			ExternalPort: node.InternalPort,
			ConnState:    node.ConnState.Int32(),
		},
	}, nil
}

func (svr *YarnServiceServer) DeleteDataNode(ctx context.Context, req *pb.DeleteRegisteredNodeRequest) (*pb.SimpleResponseCode, error) {

	if err := svr.B.DeleteRegisterNode(types.PREFIX_TYPE_DATANODE, req.Id); nil != err {
		return nil, NewRpcBizErr(ErrDeleteDataNodeInfoStr)
	}
	return &pb.SimpleResponseCode{Status: 0, Msg: OK}, nil
}

func (svr *YarnServiceServer) GetDataNodeList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetRegisteredNodeListResponse, error) {

	list, err := svr.B.GetRegisterNodeList(types.PREFIX_TYPE_DATANODE)
	if nil != err {
		return nil, NewRpcBizErr(ErrGetDataNodeListStr)
	}
	datas := make([]*pb.YarnRegisteredPeer, len(list))
	for i, v := range list {
		d := &pb.YarnRegisteredPeer{
			NodeType: types.PREFIX_TYPE_DATANODE.String(),
			NodeDetail: &pb.YarnRegisteredPeerDetail{
				Id:           v.Id,
				InternalIp:   v.InternalIp,
				InternalPort: v.InternalPort,
				ExternalIp:   v.ExternalIp,
				ExternalPort: v.InternalPort,
				ConnState:    v.ConnState.Int32(),
			},
		}
		datas[i] = d
	}
	return &pb.GetRegisteredNodeListResponse{
		Status: 0,
		Msg:    OK,
		Nodes:  datas,
	}, nil
}

func (svr *YarnServiceServer) SetJobNode(ctx context.Context, req *pb.SetJobNodeRequest) (*pb.SetJobNodeResponse, error) {
	node := &types.RegisteredNodeInfo{
		InternalIp:   req.InternalIp,
		InternalPort: req.InternalPort,
		ExternalIp:   req.ExternalIp,
		ExternalPort: req.InternalPort,
		ConnState:    types.NONCONNECTED,
	}
	node.DataNodeId()
	_, err := svr.B.SetRegisterNode(types.PREFIX_TYPE_JOBNODE, node) // TODO 未完成 ...
	if nil != err {
		return nil, NewRpcBizErr(ErrSetJobNodeInfoStr)
	}
	return &pb.SetJobNodeResponse{
		Status: 0,
		Msg:    OK,
		JobNode: &pb.YarnRegisteredPeerDetail{
			Id:           node.Id,
			InternalIp:   node.InternalIp,
			InternalPort: node.InternalPort,
			ExternalIp:   node.ExternalIp,
			ExternalPort: node.InternalPort,
			ConnState:    node.ConnState.Int32(),
		},
	}, nil
}

func (svr *YarnServiceServer) UpdateJobNode(ctx context.Context, req *pb.UpdateJobNodeRequest) (*pb.SetJobNodeResponse, error) {
	node := &types.RegisteredNodeInfo{
		Id:           req.Id,
		InternalIp:   req.InternalIp,
		InternalPort: req.InternalPort,
		ExternalIp:   req.ExternalIp,
		ExternalPort: req.InternalPort,
		ConnState:    types.NONCONNECTED,
	}
	_, err := svr.B.SetRegisterNode(types.PREFIX_TYPE_JOBNODE, node) // TODO 未完成 ...
	if nil != err {
		return nil, NewRpcBizErr(ErrSetJobNodeInfoStr)
	}
	return &pb.SetJobNodeResponse{
		Status: 0,
		Msg:    OK,
		JobNode: &pb.YarnRegisteredPeerDetail{
			Id:           node.Id,
			InternalIp:   node.InternalIp,
			InternalPort: node.InternalPort,
			ExternalIp:   node.ExternalIp,
			ExternalPort: node.InternalPort,
			ConnState:    node.ConnState.Int32(),
		},
	}, nil
}

func (svr *YarnServiceServer) DeleteJobNode(ctx context.Context, req *pb.DeleteRegisteredNodeRequest) (*pb.SimpleResponseCode, error) {
	if err := svr.B.DeleteRegisterNode(types.PREFIX_TYPE_JOBNODE, req.Id); nil != err {
		return nil, NewRpcBizErr(ErrDeleteJobNodeInfoStr)
	}
	return &pb.SimpleResponseCode{Status: 0, Msg: OK}, nil
}

func (svr *YarnServiceServer) GetJobNodeList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetRegisteredNodeListResponse, error) {
	list, err := svr.B.GetRegisterNodeList(types.PREFIX_TYPE_JOBNODE)
	if nil != err {
		return nil, NewRpcBizErr(ErrGetDataNodeListStr)
	}
	jobs := make([]*pb.YarnRegisteredPeer, len(list))
	for i, v := range list {
		d := &pb.YarnRegisteredPeer{
			NodeType: types.PREFIX_TYPE_JOBNODE.String(),
			NodeDetail: &pb.YarnRegisteredPeerDetail{
				Id:           v.Id,
				InternalIp:   v.InternalIp,
				InternalPort: v.InternalPort,
				ExternalIp:   v.ExternalIp,
				ExternalPort: v.InternalPort,
				ConnState:    v.ConnState.Int32(),
			},
		}
		jobs[i] = d
	}
	return &pb.GetRegisteredNodeListResponse{
		Status: 0,
		Msg:    OK,
		Nodes:  jobs,
	}, nil
}

func (svr *YarnServiceServer) ReportTaskEvent(ctx context.Context, req *pb.ReportTaskEventRequest) (*pb.SimpleResponseCode, error) {
	var err error
	go func() {
		err = svr.B.SendTaskEvent(&types.TaskEventInfo{
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
	return &pb.SimpleResponseCode{Status: 0, Msg: OK}, nil
}

func (svr *YarnServiceServer) ReportTaskResourceExpense(ctx context.Context, req *pb.ReportTaskResourceExpenseRequest) (*pb.SimpleResponseCode, error) {
	return nil, nil
}
