package yarn

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (svr *Server) GetNodeInfo(ctx context.Context, req *emptypb.Empty) (*pb.GetNodeInfoResponse, error) {
	node, err := svr.B.GetNodeInfo()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetNodeInfo failed")
		return nil, ErrGetNodeInfo
	}
	peers := make([]*pb.YarnRegisteredPeer, len(node.Peers))
	for i, v := range node.Peers {
		/*n := &pb.YarnRegisteredPeer{
			NodeType: v.NodeType,
			NodeDetail: &pb.YarnRegisteredPeerDetail{
				Id:           v.RegisteredNodeInfo.Id,
				InternalIp:   v.InternalIp,
				InternalPort: v.InternalPort,
				ExternalIp:   v.ExternalIp,
				ExternalPort: v.ExternalPort,
				ConnState:    v.ConnState.Int32(),
			},
		}*/
		peers[i] = v
	}
	seeds := make([]*pb.SeedPeer, len(node.SeedPeers))
	for i, v := range node.SeedPeers {
		n := &pb.SeedPeer{
			Id:           v.Id,
			InternalIp:   v.InternalIp,
			InternalPort: v.InternalPort,
			ConnState:    v.ConnState,
		}
		seeds[i] = n
	}

	return &pb.GetNodeInfoResponse{
		Status: 0,
		Msg:    backend.OK,
		Information: &pb.YarnNodeInfo{
			NodeType:     node.NodeType,
			NodeId:       node.NodeId,
			InternalIp:   node.InternalIp,
			InternalPort: node.InternalPort,
			ExternalIp:   node.ExternalIp,
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

func (svr *Server) GetRegisteredPeers(ctx context.Context, req *pb.GetRegisteredPeersRequest) (*pb.GetRegisteredPeersResponse, error) {
	registerNodes, err := svr.B.GetRegisteredPeers()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetRegisteredPeers failed")
		return nil, ErrGetRegisteredPeers
	}

	//jobNodes := make([]*pb.YarnRegisteredPeerDetail, len(registerNodes.JobNodes))
	//for i, v := range registerNodes.JobNodes {
	//	node := &pb.YarnRegisteredPeerDetail{
	//		Id:           v.Id,
	//		InternalIp:   v.InternalIp,
	//		ExternalIp:   v.ExternalIp,
	//		InternalPort: v.InternalPort,
	//		ExternalPort: v.ExternalPort,
	//		//TODO: 数据结构调整导致必须注释的
	//		/*organization: &pb.ResourceUsedDetailShow{
	//			TotalMem:       v.ResourceUsage.TotalMem,
	//			UsedMem:        v.ResourceUsage.UsedMem,
	//			TotalProcessor: v.ResourceUsage.TotalProcessor,
	//			UsedProcessor:  v.ResourceUsage.UsedProcessor,
	//			TotalBandwidth: v.ResourceUsage.TotalBandwidth,
	//			UsedBandwidth:  v.ResourceUsage.UsedBandwidth,
	//		},
	//		Duration: v.Duration,
	//		GetTask: &pb.YarnRegisteredJobNodeTaskIds{
	//			Count:   v.GetTask.Count,
	//			TaskIds: v.GetTask.TaskIds,
	//		},*/
	//	}
	//	jobNodes[i] = node
	//}

	//TODO: 数据结构调整导致必须注释的
	/*dataNodes := make([]*pb.YarnRegisteredDataNode, len(registerNodes.DataNodes))
	for i, v := range registerNodes.DataNodes {
		node := &pb.YarnRegisteredDataNode{
			Id:           v.Id,
			InternalIp:   v.InternalIp,
			ExternalIp:   v.ExternalIp,
			InternalPort: v.InternalPort,
			ExternalPort: v.ExternalPort,
			organization: &pb.ResourceUsedDetailShow{
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
	}*/
	log.Debugf("RPC-API:GetRegisteredPeers succeed, node len: {%d}", len(registerNodes))
	return &pb.GetRegisteredPeersResponse{
		Status:    0,
		Msg:       backend.OK,
		Nodes: registerNodes,
	}, nil
}

func (svr *Server) SetSeedNode(ctx context.Context, req *pb.SetSeedNodeRequest) (*pb.SetSeedNodeResponse, error) {
	seedNode := &pb.SeedPeer{
		InternalIp:   req.InternalIp,
		InternalPort: req.InternalPort,
		ConnState:    pb.ConnState_ConnState_UnConnected,
	}
	seedNode.SeedNodeId()
	status, err := svr.B.SetSeedNode(seedNode)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:SetSeedNode failed, seedNodeId: {%s}, internalIp: {%s}, internalPort: {%s}",
			seedNode.Id, req.InternalIp, req.InternalPort)
		return nil, ErrSetSeedNodeInfo
	}
	log.Debugf("RPC-API:SetSeedNode succeed, seedNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, connStatus: {%d}",
		seedNode.Id, req.InternalIp, req.InternalPort, status)
	return &pb.SetSeedNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Node: &pb.SeedPeer{
			Id:           seedNode.Id,
			InternalIp:   seedNode.InternalIp,
			InternalPort: seedNode.InternalPort,
			ConnState:    status,
		},
	}, nil
}

func (svr *Server) UpdateSeedNode(ctx context.Context, req *pb.UpdateSeedNodeRequest) (*pb.SetSeedNodeResponse, error) {
	seedNode := &pb.SeedPeer{
		Id:           req.Id,
		InternalIp:   req.InternalIp,
		InternalPort: req.InternalPort,
		ConnState:    pb.ConnState_ConnState_UnConnected,
	}
	svr.B.DeleteSeedNode(seedNode.Id)
	status, err := svr.B.SetSeedNode(seedNode)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:UpdateSeedNode failed, seedNodeId: {%s}, internalIp: {%s}, internalPort: {%s}",
			req.Id, req.InternalIp, req.InternalPort)
		return nil, ErrSetSeedNodeInfo
	}
	log.Debugf("RPC-API:UpdateSeedNode succeed, seedNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, connStatus: {%d}",
		req.Id, req.InternalIp, req.InternalPort, status)
	return &pb.SetSeedNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Node: &pb.SeedPeer{
			Id:           seedNode.Id,
			InternalIp:   seedNode.InternalIp,
			InternalPort: seedNode.InternalPort,
			ConnState:    status,
		},
	}, nil
}

func (svr *Server) DeleteSeedNode(ctx context.Context, req *pb.DeleteRegisteredNodeRequest) (*apicommonpb.SimpleResponse, error) {
	err := svr.B.DeleteSeedNode(req.Id)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:DeleteSeedNode failed, seedNodeId: {%s}", req.Id)
		return nil, ErrDeleteSeedNodeInfo
	}
	log.Debugf("RPC-API:DeleteSeedNode succeed, seedNodeId: {%s}", req.Id)
	return &apicommonpb.SimpleResponse{Status: 0, Msg: backend.OK}, nil
}

func (svr *Server) GetSeedNodeList(ctx context.Context, req *emptypb.Empty) (*pb.GetSeedNodeListResponse, error) {
	list, err := svr.B.GetSeedNodeList()
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Error("RPC-API:GetSeedNodeList failed")
		return nil, ErrGetSeedNodeList
	}
	seeds := make([]*pb.SeedPeer, len(list))
	for i, v := range list {
		s := &pb.SeedPeer{
			Id:           v.Id,
			InternalIp:   v.InternalIp,
			InternalPort: v.InternalPort,
			ConnState:    v.ConnState,
		}
		seeds[i] = s
	}
	return &pb.GetSeedNodeListResponse{
		Status:    0,
		Msg:       backend.OK,
		Nodes: 	   seeds,
	}, nil
}

func (svr *Server) SetDataNode(ctx context.Context, req *pb.SetDataNodeRequest) (*pb.SetDataNodeResponse, error) {
	node := &pb.YarnRegisteredPeerDetail{
		InternalIp:   req.InternalIp,
		InternalPort: req.InternalPort,
		ExternalIp:   req.ExternalIp,
		ExternalPort: req.ExternalPort,
		ConnState:    pb.ConnState_ConnState_UnConnected,
	}
	node.SetDataNodeId()
	status, err := svr.B.SetRegisterNode(pb.PrefixTypeDataNode, node)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:SetDataNode failed, dataNodeId:{%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}",
			node.Id, req.InternalIp, req.InternalPort, req.ExternalIp, req.ExternalPort)
		return nil, ErrSetDataNodeInfo
	}
	log.Debugf("RPC-API:SetDataNode succeed, dataNodeId:{%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}, connStatus: {%d}",
		node.Id, req.InternalIp, req.InternalPort, req.ExternalIp, req.ExternalPort, status)
	return &pb.SetDataNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Node: &pb.YarnRegisteredPeerDetail{
			Id:           node.Id,
			InternalIp:   node.InternalIp,
			InternalPort: node.InternalPort,
			ExternalIp:   node.ExternalIp,
			ExternalPort: node.ExternalPort,
			ConnState:    status,
		},
	}, nil
}

func (svr *Server) UpdateDataNode(ctx context.Context, req *pb.UpdateDataNodeRequest) (*pb.SetDataNodeResponse, error) {
	node := &pb.YarnRegisteredPeerDetail{
		Id:           req.Id,
		InternalIp:   req.InternalIp,
		InternalPort: req.InternalPort,
		ExternalIp:   req.ExternalIp,
		ExternalPort: req.ExternalPort,
		ConnState:    pb.ConnState_ConnState_UnConnected,
	}
	// delete and insert.
	//svr.B.DeleteRegisterNode(types.PrefixTypeDataNode, node.Id)
	//status, err := svr.B.SetRegisterNode(types.PrefixTypeDataNode, node)
	status, err := svr.B.UpdateRegisterNode(pb.PrefixTypeDataNode, node)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:UpdateDataNode failed, dataNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}",
			req.Id, req.InternalIp, req.InternalPort, req.ExternalIp, req.ExternalPort)
		return nil, ErrSetDataNodeInfo
	}
	log.Debugf("RPC-API:UpdateDataNode succeed, dataNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}, connStatus: {%d}",
		req.Id, req.InternalIp, req.InternalPort, req.ExternalIp, req.ExternalPort, status)
	return &pb.SetDataNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Node: &pb.YarnRegisteredPeerDetail{
			Id:           node.Id,
			InternalIp:   node.InternalIp,
			InternalPort: node.InternalPort,
			ExternalIp:   node.ExternalIp,
			ExternalPort: node.ExternalPort,
			ConnState:    status,
		},
	}, nil
}

func (svr *Server) DeleteDataNode(ctx context.Context, req *pb.DeleteRegisteredNodeRequest) (*apicommonpb.SimpleResponse, error) {
	if err := svr.B.DeleteRegisterNode(pb.PrefixTypeDataNode, req.Id); nil != err {
		log.WithError(err).Errorf("RPC-API:DeleteDataNode failed, dataNodeId: {%s}", req.Id)
		return nil, ErrDeleteDataNodeInfo
	}
	log.Debugf("RPC-API:DeleteDataNode succeed, dataNodeId: {%s}", req.Id)
	return &apicommonpb.SimpleResponse{Status: 0, Msg: backend.OK}, nil
}

func (svr *Server) GetDataNodeList(ctx context.Context, req *emptypb.Empty) (*pb.GetRegisteredNodeListResponse, error) {

	list, err := svr.B.GetRegisterNodeList(pb.PrefixTypeDataNode)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Error("RPC-API:GetDataNodeList failed")
		return nil, ErrGetDataNodeList
	}
	datas := make([]*pb.YarnRegisteredPeer, len(list))
	for i, v := range list {
		d := &pb.YarnRegisteredPeer{
			NodeType: pb.NodeType(pb.NodeType_value[pb.PrefixTypeDataNode.String()]) ,
			NodeDetail: &pb.YarnRegisteredPeerDetail{
				Id:           v.Id,
				InternalIp:   v.InternalIp,
				InternalPort: v.InternalPort,
				ExternalIp:   v.ExternalIp,
				ExternalPort: v.ExternalPort,
				ConnState:    v.ConnState,
			},
		}
		datas[i] = d
	}
	return &pb.GetRegisteredNodeListResponse{
		Status: 0,
		Msg:    backend.OK,
		Nodes:  datas,
	}, nil
}

func (svr *Server) SetJobNode(ctx context.Context, req *pb.SetJobNodeRequest) (*pb.SetJobNodeResponse, error) {
	node := &pb.YarnRegisteredPeerDetail{
		InternalIp:   req.InternalIp,
		InternalPort: req.InternalPort,
		ExternalIp:   req.ExternalIp,
		ExternalPort: req.ExternalPort,
		ConnState:    pb.ConnState_ConnState_UnConnected,
	}
	node.SetJobNodeId()
	status, err := svr.B.SetRegisterNode(pb.PrefixTypeJobNode, node)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:SetJobNode failed, jobNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}",
			node.Id, req.InternalIp, req.InternalPort, req.ExternalIp, req.ExternalPort)
		return nil, ErrSetJobNodeInfo
	}

	log.Debugf("RPC-API:SetJobNode succeed, jobNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}, connStats: {%d}",
		node.Id, req.InternalIp, req.InternalPort, req.ExternalIp, req.ExternalPort, status)
	return &pb.SetJobNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Node: &pb.YarnRegisteredPeerDetail{
			Id:           node.Id,
			InternalIp:   node.InternalIp,
			InternalPort: node.InternalPort,
			ExternalIp:   node.ExternalIp,
			ExternalPort: node.ExternalPort,
			ConnState:    status,
		},
	}, nil
}

func (svr *Server) UpdateJobNode(ctx context.Context, req *pb.UpdateJobNodeRequest) (*pb.SetJobNodeResponse, error) {
	node := &pb.YarnRegisteredPeerDetail{
		Id:           req.Id,
		InternalIp:   req.InternalIp,
		InternalPort: req.InternalPort,
		ExternalIp:   req.ExternalIp,
		ExternalPort: req.ExternalPort,
		ConnState:    pb.ConnState_ConnState_UnConnected,
	}
	//svr.B.DeleteRegisterNode(types.PrefixTypeJobNode, node.Id)
	//status, err := svr.B.SetRegisterNode(types.PrefixTypeJobNode, node)
	status, err := svr.B.UpdateRegisterNode(pb.PrefixTypeJobNode, node)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:UpdateJobNode failed, jobNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}",
			req.Id, req.InternalIp, req.InternalPort, req.ExternalIp, req.ExternalPort)
		return nil, ErrSetJobNodeInfo
	}

	log.Debugf("RPC-API:UpdateJobNode succeed, jobNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}, connStats: {%d}",
		req.Id, req.InternalIp, req.InternalPort, req.ExternalIp, req.ExternalPort, status)
	return &pb.SetJobNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Node: &pb.YarnRegisteredPeerDetail{
			Id:           node.Id,
			InternalIp:   node.InternalIp,
			InternalPort: node.InternalPort,
			ExternalIp:   node.ExternalIp,
			ExternalPort: node.ExternalPort,
			ConnState:    status,
		},
	}, nil
}

func (svr *Server) DeleteJobNode(ctx context.Context, req *pb.DeleteRegisteredNodeRequest) (*apicommonpb.SimpleResponse, error) {
	if err := svr.B.DeleteRegisterNode(pb.PrefixTypeJobNode, req.Id); nil != err {
		log.WithError(err).Errorf("RPC-API:DeleteJobNode failed, jobNodeId: {%s}", req.Id)
		return nil, ErrDeleteJobNodeInfo
	}
	log.Debugf("RPC-API:DeleteJobNode succeed, jobNodeId: {%s}", req.Id)
	return &apicommonpb.SimpleResponse{Status: 0, Msg: backend.OK}, nil
}

func (svr *Server) GetJobNodeList(ctx context.Context, req *emptypb.Empty) (*pb.GetRegisteredNodeListResponse, error) {
	list, err := svr.B.GetRegisterNodeList(pb.PrefixTypeJobNode)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Error("RPC-API:GetJobNodeList failed")
		return nil, ErrGetDataNodeList
	}
	jobs := make([]*pb.YarnRegisteredPeer, len(list))
	for i, v := range list {
		d := &pb.YarnRegisteredPeer{
			NodeType: pb.NodeType(pb.NodeType_value[pb.PrefixTypeJobNode.String()]),
			NodeDetail: &pb.YarnRegisteredPeerDetail{
				Id:           v.Id,
				InternalIp:   v.InternalIp,
				InternalPort: v.InternalPort,
				ExternalIp:   v.ExternalIp,
				ExternalPort: v.ExternalPort,
				ConnState:    v.ConnState,
			},
		}
		jobs[i] = d
	}
	return &pb.GetRegisteredNodeListResponse{
		Status: 0,
		Msg:    backend.OK,
		Nodes:  jobs,
	}, nil
}

func (svr *Server) QueryAvailableDataNode(ctx context.Context, req *pb.QueryAvailableDataNodeRequest) (*pb.QueryAvailableDataNodeResponse, error) {
	dataResourceTables, err := svr.B.QueryDataResourceTables()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryAvailableDataNode-QueryDataResourceTables failed, fileType: {%s}, fileSize: {%s}", req.FileType, req.FileSize)
		return nil, ErrQueryDataResourceTableList
	}

	var nodeId string
	for _, resource := range dataResourceTables {
		if req.FileSize < resource.RemainDisk() {
			nodeId = resource.GetNodeId()
			break
		}
	}

	dataNode, err := svr.B.GetRegisterNode(pb.PrefixTypeDataNode, nodeId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryAvailableDataNode-GetRegisterNode failed, fileType: {%s}, fileSize: {%s}, dataNodeId: {%s}",
			req.FileType, req.FileSize, nodeId)
		return nil, ErrGetDataNodeInfo
	}
	log.Debugf("RPC-API:QueryAvailableDataNode succeed, fileType: {%s}, fileSize: {%d}, return dataNodeId: {%s}, dataNodeIp: {%s}, dataNodePort: {%s}",
		req.FileType, req.FileSize, dataNode.Id, dataNode.InternalIp, dataNode.InternalPort)

	return &pb.QueryAvailableDataNodeResponse{
		Ip:   dataNode.InternalIp,
		Port: dataNode.InternalPort,
	}, nil
}
func (svr *Server) QueryFilePosition(ctx context.Context, req *pb.QueryFilePositionRequest) (*pb.QueryFilePositionResponse, error) {
	dataResourceFileUpload, err := svr.B.QueryDataResourceFileUpload(req.OriginId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryFilePosition-QueryDataResourceFileUpload failed, originId: {%s}", req.OriginId)
		return nil, ErrQueryDataResourceDataUsed
	}
	dataNode, err := svr.B.GetRegisterNode(pb.PrefixTypeDataNode, dataResourceFileUpload.GetNodeId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryFilePosition-GetRegisterNode failed, originId: {%s}", req.OriginId)
		return nil, ErrGetDataNodeInfo
	}

	log.Debugf("RPC-API:QueryFilePosition Succeed, originId: {%s}, return dataNodeIp: {%s}, dataNodePort: {%s}, filePath: {%s}",
		req.OriginId, dataNode.InternalIp, dataNode.InternalPort, dataResourceFileUpload.GetFilePath())

	return &pb.QueryFilePositionResponse{
		Ip:       dataNode.InternalIp,
		Port:     dataNode.InternalPort,
		FilePath: dataResourceFileUpload.GetFilePath(),
	}, nil
}

func (svr *Server) ReportTaskResultFileSummary(ctx context.Context, request *pb.ReportTaskResultFileSummaryRequest) (*apicommonpb.SimpleResponse, error) {
	panic("implement me")
}

func (svr *Server) GetTaskResultFileSummary(ctx context.Context, request *pb.GetTaskResultFileSummaryRequest) (*pb.GetTaskResultFileSummaryResponse, error) {
	panic("implement me")
}

func (svr *Server) GetTaskResultFileSummaryList(ctx context.Context, empty *emptypb.Empty) (*pb.GetTaskResultFileSummaryListResponse, error) {
	panic("implement me")
}