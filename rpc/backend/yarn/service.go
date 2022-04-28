package yarn

import (
	"context"
	"fmt"
	"github.com/Metisnetwork/Metis-Carrier/core/rawdb"
	pb "github.com/Metisnetwork/Metis-Carrier/lib/api"
	libtypes "github.com/Metisnetwork/Metis-Carrier/lib/types"
	"github.com/Metisnetwork/Metis-Carrier/rpc/backend"
	"github.com/Metisnetwork/Metis-Carrier/service/discovery"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (svr *Server) GetNodeInfo(ctx context.Context, req *emptypb.Empty) (*pb.GetNodeInfoResponse, error) {
	node, err := svr.B.GetNodeInfo()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetNodeInfo failed")
		return &pb.GetNodeInfoResponse{Status: backend.ErrQueryNodeInfo.ErrCode(), Msg: backend.ErrQueryNodeInfo.Error()}, nil
	}
	// setting gRPC server ip and port to internalIp and internalPort
	node.InternalIp = svr.RpcSvrIp
	node.InternalPort = svr.RpcSvrPort

	return &pb.GetNodeInfoResponse{
		Status:      0,
		Msg:         backend.OK,
		Information: node,
	}, nil
}

func (svr *Server) SetSeedNode(ctx context.Context, req *pb.SetSeedNodeRequest) (*pb.SetSeedNodeResponse, error) {

	if "" == strings.Trim(req.GetAddr(), "") {
		return &pb.SetSeedNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require addr of seedNode"}, nil
	}

	list, err := svr.B.GetSeedNodeList()
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Error("RPC-API: call svr.B.GetSeedNodeList() failed on SetSeedNode()")
		return &pb.SetSeedNodeResponse{Status: backend.ErrSetSeedNode.ErrCode(), Msg: fmt.Sprintf("query seedNode list failed, %s", err)}, nil
	}

	for _, seed := range list {
		if seed.GetAddr() == strings.Trim(req.GetAddr(), "") {
			log.WithError(err).Error("RPC-API: call svr.B.GetSeedNodeList() failed on SetSeedNode()")
			return &pb.SetSeedNodeResponse{Status: backend.ErrSetSeedNode.ErrCode(), Msg: "the addr has already exists"}, nil
		}
	}

	seedNode := &pb.SeedPeer{
		Addr:      strings.Trim(req.GetAddr(), ""),
		IsDefault: false,
		ConnState: pb.ConnState_ConnState_UnConnected,
	}
	status, err := svr.B.SetSeedNode(seedNode)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:SetSeedNode failed, addr: {%s}", seedNode.GetAddr())
		errMsg := fmt.Sprintf("%s, %s, %s", backend.ErrSetSeedNode.Error(), "SetSeedNode", seedNode.GetAddr())
		return &pb.SetSeedNodeResponse{Status: backend.ErrSetSeedNode.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:SetSeedNode succeed, addr: {%s}, connStatus: {%d}", req.GetAddr(), status)
	seedNode.ConnState = status
	return &pb.SetSeedNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Node:   seedNode,
	}, nil
}

func (svr *Server) DeleteSeedNode(ctx context.Context, req *pb.DeleteSeedNodeRequest) (*libtypes.SimpleResponse, error) {

	if "" == strings.Trim(req.GetAddr(), "") {
		return &libtypes.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require addr of seedNode"}, nil
	}

	err := svr.B.DeleteSeedNode(strings.Trim(req.GetAddr(), ""))
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RemoveSeedNode failed, addr: {%s}", req.GetAddr())
		errMsg := fmt.Sprintf("%s, %s", backend.ErrDeleteSeedNode.Error(), req.GetAddr())
		return &libtypes.SimpleResponse{Status: backend.ErrDeleteSeedNode.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:RemoveSeedNode succeed, addr: {%s}", req.GetAddr())
	return &libtypes.SimpleResponse{Status: 0, Msg: backend.OK}, nil
}

func (svr *Server) GetSeedNodeList(ctx context.Context, req *emptypb.Empty) (*pb.GetSeedNodeListResponse, error) {
	list, err := svr.B.GetSeedNodeList()
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Error("RPC-API: call svr.B.GetSeedNodeList() failed")
		return &pb.GetSeedNodeListResponse{Status: backend.ErrQuerySeedNodeList.ErrCode(), Msg: backend.ErrQuerySeedNodeList.Error()}, nil
	}
	return &pb.GetSeedNodeListResponse{
		Status: 0,
		Msg:    backend.OK,
		Nodes:  list,
	}, nil
}

func (svr *Server) SetDataNode(ctx context.Context, req *pb.SetDataNodeRequest) (*pb.SetDataNodeResponse, error) {

	if "" == strings.Trim(req.GetInternalIp(), "") {
		return &pb.SetDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require internal Ip"}, nil
	}

	if "" == strings.Trim(req.GetInternalPort(), "") {
		return &pb.SetDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require internal port"}, nil
	}

	if "" == strings.Trim(req.GetExternalIp(), "") {
		return &pb.SetDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require external Ip"}, nil
	}

	if "" == strings.Trim(req.GetExternalPort(), "") {
		return &pb.SetDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require external port"}, nil
	}

	node := &pb.YarnRegisteredPeerDetail{
		InternalIp:   req.GetInternalIp(),
		InternalPort: req.GetInternalPort(),
		ExternalIp:   req.GetExternalIp(),
		ExternalPort: req.GetExternalPort(),
		ConnState:    pb.ConnState_ConnState_UnConnected,
	}
	node.Id = strings.Join([]string{discovery.DataNodeConsulServiceIdPrefix, req.GetInternalIp(), req.GetInternalPort()}, discovery.ConsulServiceIdSeparator)

	status, err := svr.B.SetRegisterNode(pb.PrefixTypeDataNode, node)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:SetDataNode failed, dataNodeId:{%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}",
			node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())

		errMsg := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s", backend.ErrSetDataNode.Error(), "SetDataNode",
			node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())
		return &pb.SetDataNodeResponse{Status: backend.ErrSetDataNode.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:SetDataNode succeed, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}, connStatus: {%d} return dataNodeId: {%s}",
		req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort(), status, node.GetId())
	return &pb.SetDataNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Node: &pb.YarnRegisteredPeerDetail{
			Id:           node.GetId(),
			InternalIp:   node.GetInternalIp(),
			InternalPort: node.GetInternalPort(),
			ExternalIp:   node.GetExternalIp(),
			ExternalPort: node.GetExternalPort(),
			ConnState:    status,
		},
	}, nil
}

func (svr *Server) UpdateDataNode(ctx context.Context, req *pb.UpdateDataNodeRequest) (*pb.SetDataNodeResponse, error) {

	if "" == strings.Trim(req.GetId(), "") {
		return &pb.SetDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require id of data node"}, nil
	}

	if "" == strings.Trim(req.GetInternalIp(), "") {
		return &pb.SetDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require internal Ip"}, nil
	}

	if "" == strings.Trim(req.GetInternalPort(), "") {
		return &pb.SetDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require internal port"}, nil
	}

	if "" == strings.Trim(req.GetExternalIp(), "") {
		return &pb.SetDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require external Ip"}, nil
	}

	if "" == strings.Trim(req.GetExternalPort(), "") {
		return &pb.SetDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require external port"}, nil
	}

	node := &pb.YarnRegisteredPeerDetail{
		Id:           req.GetId(),
		InternalIp:   req.GetInternalIp(),
		InternalPort: req.GetInternalPort(),
		ExternalIp:   req.GetExternalIp(),
		ExternalPort: req.GetExternalPort(),
		ConnState:    pb.ConnState_ConnState_UnConnected,
	}
	status, err := svr.B.UpdateRegisterNode(pb.PrefixTypeDataNode, node)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:UpdateDataNode failed, dataNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}",
			node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())

		errMsg := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s", backend.ErrSetDataNode.Error(), "UpdateDataNode",
			node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())
		return &pb.SetDataNodeResponse{Status: backend.ErrSetDataNode.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:UpdateDataNode succeed, dataNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}, connStatus: {%d}",
		node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort(), status)
	return &pb.SetDataNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Node: &pb.YarnRegisteredPeerDetail{
			Id:           node.GetId(),
			InternalIp:   node.GetInternalIp(),
			InternalPort: node.GetInternalPort(),
			ExternalIp:   node.GetExternalIp(),
			ExternalPort: node.GetExternalPort(),
			ConnState:    status,
		},
	}, nil
}

func (svr *Server) DeleteDataNode(ctx context.Context, req *pb.DeleteRegisteredNodeRequest) (*libtypes.SimpleResponse, error) {

	if "" == strings.Trim(req.GetId(), "") {
		return &libtypes.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require id of data node"}, nil
	}

	if err := svr.B.DeleteRegisterNode(pb.PrefixTypeDataNode, req.GetId()); nil != err {
		log.WithError(err).Errorf("RPC-API:DeleteDataNode failed, dataNodeId: {%s}", req.GetId())

		errMsg := fmt.Sprintf("%s, %s", backend.ErrDeleteDataNode.Error(), req.GetId())
		return &libtypes.SimpleResponse{Status: backend.ErrDeleteDataNode.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:DeleteDataNode succeed, dataNodeId: {%s}", req.GetId())
	return &libtypes.SimpleResponse{Status: 0, Msg: backend.OK}, nil
}

func (svr *Server) GetDataNodeList(ctx context.Context, req *emptypb.Empty) (*pb.GetRegisteredNodeListResponse, error) {

	list, err := svr.B.GetRegisterNodeList(pb.PrefixTypeDataNode)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Error("RPC-API:GetDataNodeList failed")
		return &pb.GetRegisteredNodeListResponse{Status: backend.ErrQueryDataNodeList.ErrCode(), Msg: backend.ErrQueryDataNodeList.Error()}, nil
	}
	datas := make([]*pb.YarnRegisteredPeer, len(list))
	for i, v := range list {
		d := &pb.YarnRegisteredPeer{
			NodeType: pb.NodeType_NodeType_DataNode,
			NodeDetail: &pb.YarnRegisteredPeerDetail{
				Id:            v.GetId(),
				InternalIp:    v.GetInternalIp(),
				InternalPort:  v.GetInternalPort(),
				ExternalIp:    v.GetExternalIp(),
				ExternalPort:  v.GetExternalPort(),
				ConnState:     v.GetConnState(),
				Duration:      v.GetDuration(),
				TaskCount:     v.GetTaskCount(),
				TaskIdList:    v.GetTaskIdList(),
				FileCount:     v.GetFileCount(),
				FileTotalSize: v.GetFileTotalSize(),
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

	if "" == strings.Trim(req.GetInternalIp(), "") {
		return &pb.SetJobNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require internal Ip"}, nil
	}

	if "" == strings.Trim(req.GetInternalPort(), "") {
		return &pb.SetJobNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require internal port"}, nil
	}

	if "" == strings.Trim(req.GetExternalIp(), "") {
		return &pb.SetJobNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require external Ip"}, nil
	}

	if "" == strings.Trim(req.GetExternalPort(), "") {
		return &pb.SetJobNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require external port"}, nil
	}

	node := &pb.YarnRegisteredPeerDetail{
		InternalIp:   req.GetInternalIp(),
		InternalPort: req.GetInternalPort(),
		ExternalIp:   req.GetExternalIp(),
		ExternalPort: req.GetExternalPort(),
		ConnState:    pb.ConnState_ConnState_UnConnected,
	}
	node.Id = strings.Join([]string{discovery.JobNodeConsulServiceIdPrefix, req.GetInternalIp(), req.GetInternalPort()}, discovery.ConsulServiceIdSeparator)

	status, err := svr.B.SetRegisterNode(pb.PrefixTypeJobNode, node)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:SetJobNode failed, jobNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}",
			node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())

		errMsg := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s", backend.ErrSetJobNode.Error(), "SetJobNode",
			node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())
		return &pb.SetJobNodeResponse{Status: backend.ErrSetJobNode.ErrCode(), Msg: errMsg}, nil
	}

	log.Debugf("RPC-API:SetJobNode succeed, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}, connStats: {%d} return jobNodeId: {%s}",
		req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort(), status, node.GetId())
	return &pb.SetJobNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Node: &pb.YarnRegisteredPeerDetail{
			Id:           node.GetId(),
			InternalIp:   node.GetInternalIp(),
			InternalPort: node.GetInternalPort(),
			ExternalIp:   node.GetExternalIp(),
			ExternalPort: node.GetExternalPort(),
			ConnState:    status,
		},
	}, nil
}

func (svr *Server) UpdateJobNode(ctx context.Context, req *pb.UpdateJobNodeRequest) (*pb.SetJobNodeResponse, error) {

	if "" == strings.Trim(req.GetId(), "") {
		return &pb.SetJobNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require id of job node"}, nil
	}

	if "" == strings.Trim(req.GetInternalIp(), "") {
		return &pb.SetJobNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require internal Ip"}, nil
	}

	if "" == strings.Trim(req.GetInternalPort(), "") {
		return &pb.SetJobNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require internal port"}, nil
	}

	if "" == strings.Trim(req.GetExternalIp(), "") {
		return &pb.SetJobNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require external Ip"}, nil
	}

	if "" == strings.Trim(req.GetExternalPort(), "") {
		return &pb.SetJobNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require external port"}, nil
	}

	node := &pb.YarnRegisteredPeerDetail{
		Id:           req.GetId(),
		InternalIp:   req.GetInternalIp(),
		InternalPort: req.GetInternalPort(),
		ExternalIp:   req.GetExternalIp(),
		ExternalPort: req.GetExternalPort(),
		ConnState:    pb.ConnState_ConnState_UnConnected,
	}
	status, err := svr.B.UpdateRegisterNode(pb.PrefixTypeJobNode, node)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:UpdateJobNode failed, jobNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}",
			req.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())

		errMsg := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s", backend.ErrSetJobNode.Error(), "UpdateJobNode",
			req.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())

		return &pb.SetJobNodeResponse{Status: backend.ErrSetJobNode.ErrCode(), Msg: errMsg}, nil
	}

	log.Debugf("RPC-API:UpdateJobNode succeed, jobNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}, connStats: {%d}",
		req.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort(), status)
	return &pb.SetJobNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Node: &pb.YarnRegisteredPeerDetail{
			Id:           node.GetId(),
			InternalIp:   node.GetInternalIp(),
			InternalPort: node.GetInternalPort(),
			ExternalIp:   node.GetExternalIp(),
			ExternalPort: node.GetExternalPort(),
			ConnState:    status,
		},
	}, nil
}

func (svr *Server) DeleteJobNode(ctx context.Context, req *pb.DeleteRegisteredNodeRequest) (*libtypes.SimpleResponse, error) {

	if "" == strings.Trim(req.GetId(), "") {
		return &libtypes.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require id of job node"}, nil
	}

	if err := svr.B.DeleteRegisterNode(pb.PrefixTypeJobNode, req.GetId()); nil != err {
		log.WithError(err).Errorf("RPC-API:DeleteJobNode failed, jobNodeId: {%s}", req.GetId())

		errMsg := fmt.Sprintf("%s, %s", backend.ErrDeleteJobNode.Error(), req.GetId())
		return &libtypes.SimpleResponse{Status: backend.ErrDeleteJobNode.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:DeleteJobNode succeed, jobNodeId: {%s}", req.GetId())
	return &libtypes.SimpleResponse{Status: 0, Msg: backend.OK}, nil
}

func (svr *Server) GetJobNodeList(ctx context.Context, req *emptypb.Empty) (*pb.GetRegisteredNodeListResponse, error) {
	list, err := svr.B.GetRegisterNodeList(pb.PrefixTypeJobNode)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Error("RPC-API:GetJobNodeList failed")
		return &pb.GetRegisteredNodeListResponse{Status: backend.ErrQueryJobNodeList.ErrCode(), Msg: backend.ErrQueryJobNodeList.Error()}, nil
	}
	jobs := make([]*pb.YarnRegisteredPeer, len(list))
	for i, v := range list {
		d := &pb.YarnRegisteredPeer{
			NodeType: pb.NodeType_NodeType_JobNode,
			NodeDetail: &pb.YarnRegisteredPeerDetail{
				Id:            v.GetId(),
				InternalIp:    v.GetInternalIp(),
				InternalPort:  v.GetInternalPort(),
				ExternalIp:    v.GetExternalIp(),
				ExternalPort:  v.GetExternalPort(),
				ConnState:     v.GetConnState(),
				Duration:      v.GetDuration(),
				TaskCount:     v.GetTaskCount(),
				TaskIdList:    v.GetTaskIdList(),
				FileCount:     v.GetFileCount(),
				FileTotalSize: v.GetFileTotalSize(),
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

func (svr *Server) GenerateObServerProxyWalletAddress(ctx context.Context, req *emptypb.Empty) (*pb.GenerateObServerProxyWalletAddressResponse, error) {

	wallet, err := svr.B.GenerateObServerProxyWalletAddress()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GenerateObServerProxyWalletAddress failed")
		return &pb.GenerateObServerProxyWalletAddressResponse{Status: backend.ErrGenerateObserverProxyWallet.ErrCode(), Msg: backend.ErrGenerateObserverProxyWallet.Error()}, nil
	}
	log.Debugf("RPC-API:GenerateObServerProxyWalletAddress Succeed, wallet: {%s}", wallet)
	return &pb.GenerateObServerProxyWalletAddressResponse{
		Status:  0,
		Msg:     backend.OK,
		Address: wallet,
	}, nil
}
