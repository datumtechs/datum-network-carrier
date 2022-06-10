package yarn

import (
	"context"
	"fmt"
	rawdb "github.com/datumtechs/datum-network-carrier/carrierdb/rawdb"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/rpc/backend"
	"github.com/datumtechs/datum-network-carrier/service/discovery"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (svr *Server) GetNodeInfo(ctx context.Context, req *emptypb.Empty) (*carrierapipb.GetNodeInfoResponse, error) {
	node, err := svr.B.GetNodeInfo()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetNodeInfo failed")
		return &carrierapipb.GetNodeInfoResponse{Status: backend.ErrQueryNodeInfo.ErrCode(), Msg: backend.ErrQueryNodeInfo.Error()}, nil
	}
	// setting gRPC server ip and port to internalIp and internalPort
	node.InternalIp = svr.RpcSvrIp
	node.InternalPort = svr.RpcSvrPort

	return &carrierapipb.GetNodeInfoResponse{
		Status:      0,
		Msg:         backend.OK,
		Information: node,
	}, nil
}

func (svr *Server) SetSeedNode(ctx context.Context, req *carrierapipb.SetSeedNodeRequest) (*carrierapipb.SetSeedNodeResponse, error) {

	if "" == strings.Trim(req.GetAddr(), "") {
		return &carrierapipb.SetSeedNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require addr of seedNode"}, nil
	}

	list, err := svr.B.GetSeedNodeList()
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Error("RPC-API: call svr.B.GetSeedNodeList() failed on SetSeedNode()")
		return &carrierapipb.SetSeedNodeResponse{Status: backend.ErrSetSeedNode.ErrCode(), Msg: fmt.Sprintf("query seedNode list failed, %s", err)}, nil
	}

	for _, seed := range list {
		if seed.GetAddr() == strings.Trim(req.GetAddr(), "") {
			log.WithError(err).Error("RPC-API: call svr.B.GetSeedNodeList() failed on SetSeedNode()")
			return &carrierapipb.SetSeedNodeResponse{Status: backend.ErrSetSeedNode.ErrCode(), Msg: "the addr has already exists"}, nil
		}
	}

	seedNode := &carrierapipb.SeedPeer{
		Addr:      strings.Trim(req.GetAddr(), ""),
		IsDefault: false,
		ConnState: carrierapipb.ConnState_ConnState_UnConnected,
	}
	status, err := svr.B.SetSeedNode(seedNode)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:SetSeedNode failed, addr: {%s}", seedNode.GetAddr())
		errMsg := fmt.Sprintf("%s, %s, %s", backend.ErrSetSeedNode.Error(), "SetSeedNode", seedNode.GetAddr())
		return &carrierapipb.SetSeedNodeResponse{Status: backend.ErrSetSeedNode.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:SetSeedNode succeed, addr: {%s}, connStatus: {%d}", req.GetAddr(), status)
	seedNode.ConnState = status
	return &carrierapipb.SetSeedNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Node:   seedNode,
	}, nil
}

func (svr *Server) DeleteSeedNode(ctx context.Context, req *carrierapipb.DeleteSeedNodeRequest) (*carriertypespb.SimpleResponse, error) {

	if "" == strings.Trim(req.GetAddr(), "") {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require addr of seedNode"}, nil
	}

	err := svr.B.DeleteSeedNode(strings.Trim(req.GetAddr(), ""))
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RemoveSeedNode failed, addr: {%s}", req.GetAddr())
		errMsg := fmt.Sprintf("%s, %s", backend.ErrDeleteSeedNode.Error(), req.GetAddr())
		return &carriertypespb.SimpleResponse{Status: backend.ErrDeleteSeedNode.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:RemoveSeedNode succeed, addr: {%s}", req.GetAddr())
	return &carriertypespb.SimpleResponse{Status: 0, Msg: backend.OK}, nil
}

func (svr *Server) GetSeedNodeList(ctx context.Context, req *emptypb.Empty) (*carrierapipb.GetSeedNodeListResponse, error) {
	list, err := svr.B.GetSeedNodeList()
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Error("RPC-API: call svr.B.GetSeedNodeList() failed")
		return &carrierapipb.GetSeedNodeListResponse{Status: backend.ErrQuerySeedNodeList.ErrCode(), Msg: backend.ErrQuerySeedNodeList.Error()}, nil
	}
	return &carrierapipb.GetSeedNodeListResponse{
		Status: 0,
		Msg:    backend.OK,
		Nodes:  list,
	}, nil
}

func (svr *Server) SetDataNode(ctx context.Context, req *carrierapipb.SetDataNodeRequest) (*carrierapipb.SetDataNodeResponse, error) {

	if "" == strings.Trim(req.GetInternalIp(), "") {
		return &carrierapipb.SetDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require internal Ip"}, nil
	}

	if "" == strings.Trim(req.GetInternalPort(), "") {
		return &carrierapipb.SetDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require internal port"}, nil
	}

	if "" == strings.Trim(req.GetExternalIp(), "") {
		return &carrierapipb.SetDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require external Ip"}, nil
	}

	if "" == strings.Trim(req.GetExternalPort(), "") {
		return &carrierapipb.SetDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require external port"}, nil
	}

	node := &carrierapipb.YarnRegisteredPeerDetail{
		InternalIp:   req.GetInternalIp(),
		InternalPort: req.GetInternalPort(),
		ExternalIp:   req.GetExternalIp(),
		ExternalPort: req.GetExternalPort(),
		ConnState:    carrierapipb.ConnState_ConnState_UnConnected,
	}
	node.Id = strings.Join([]string{discovery.DataNodeConsulServiceIdPrefix, req.GetInternalIp(), req.GetInternalPort()}, discovery.ConsulServiceIdSeparator)

	status, err := svr.B.SetRegisterNode(carrierapipb.PrefixTypeDataNode, node)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:SetDataNode failed, dataNodeId:{%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}",
			node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())

		errMsg := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s", backend.ErrSetDataNode.Error(), "SetDataNode",
			node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())
		return &carrierapipb.SetDataNodeResponse{Status: backend.ErrSetDataNode.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:SetDataNode succeed, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}, connStatus: {%d} return dataNodeId: {%s}",
		req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort(), status, node.GetId())
	return &carrierapipb.SetDataNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Node: &carrierapipb.YarnRegisteredPeerDetail{
			Id:           node.GetId(),
			InternalIp:   node.GetInternalIp(),
			InternalPort: node.GetInternalPort(),
			ExternalIp:   node.GetExternalIp(),
			ExternalPort: node.GetExternalPort(),
			ConnState:    status,
		},
	}, nil
}

func (svr *Server) UpdateDataNode(ctx context.Context, req *carrierapipb.UpdateDataNodeRequest) (*carrierapipb.SetDataNodeResponse, error) {

	if "" == strings.Trim(req.GetId(), "") {
		return &carrierapipb.SetDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require id of data node"}, nil
	}

	if "" == strings.Trim(req.GetInternalIp(), "") {
		return &carrierapipb.SetDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require internal Ip"}, nil
	}

	if "" == strings.Trim(req.GetInternalPort(), "") {
		return &carrierapipb.SetDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require internal port"}, nil
	}

	if "" == strings.Trim(req.GetExternalIp(), "") {
		return &carrierapipb.SetDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require external Ip"}, nil
	}

	if "" == strings.Trim(req.GetExternalPort(), "") {
		return &carrierapipb.SetDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require external port"}, nil
	}

	node := &carrierapipb.YarnRegisteredPeerDetail{
		Id:           req.GetId(),
		InternalIp:   req.GetInternalIp(),
		InternalPort: req.GetInternalPort(),
		ExternalIp:   req.GetExternalIp(),
		ExternalPort: req.GetExternalPort(),
		ConnState:    carrierapipb.ConnState_ConnState_UnConnected,
	}
	status, err := svr.B.UpdateRegisterNode(carrierapipb.PrefixTypeDataNode, node)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:UpdateDataNode failed, dataNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}",
			node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())

		errMsg := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s", backend.ErrSetDataNode.Error(), "UpdateDataNode",
			node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())
		return &carrierapipb.SetDataNodeResponse{Status: backend.ErrSetDataNode.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:UpdateDataNode succeed, dataNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}, connStatus: {%d}",
		node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort(), status)
	return &carrierapipb.SetDataNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Node: &carrierapipb.YarnRegisteredPeerDetail{
			Id:           node.GetId(),
			InternalIp:   node.GetInternalIp(),
			InternalPort: node.GetInternalPort(),
			ExternalIp:   node.GetExternalIp(),
			ExternalPort: node.GetExternalPort(),
			ConnState:    status,
		},
	}, nil
}

func (svr *Server) DeleteDataNode(ctx context.Context, req *carrierapipb.DeleteRegisteredNodeRequest) (*carriertypespb.SimpleResponse, error) {

	if "" == strings.Trim(req.GetId(), "") {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require id of data node"}, nil
	}

	if err := svr.B.DeleteRegisterNode(carrierapipb.PrefixTypeDataNode, req.GetId()); nil != err {
		log.WithError(err).Errorf("RPC-API:DeleteDataNode failed, dataNodeId: {%s}", req.GetId())

		errMsg := fmt.Sprintf("%s, %s", backend.ErrDeleteDataNode.Error(), req.GetId())
		return &carriertypespb.SimpleResponse{Status: backend.ErrDeleteDataNode.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:DeleteDataNode succeed, dataNodeId: {%s}", req.GetId())
	return &carriertypespb.SimpleResponse{Status: 0, Msg: backend.OK}, nil
}

func (svr *Server) GetDataNodeList(ctx context.Context, req *emptypb.Empty) (*carrierapipb.GetRegisteredNodeListResponse, error) {

	list, err := svr.B.GetRegisterNodeList(carrierapipb.PrefixTypeDataNode)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Error("RPC-API:GetDataNodeList failed")
		return &carrierapipb.GetRegisteredNodeListResponse{Status: backend.ErrQueryDataNodeList.ErrCode(), Msg: backend.ErrQueryDataNodeList.Error()}, nil
	}
	datas := make([]*carrierapipb.YarnRegisteredPeer, len(list))
	for i, v := range list {
		d := &carrierapipb.YarnRegisteredPeer{
			NodeType: carrierapipb.NodeType_NodeType_DataNode,
			NodeDetail: &carrierapipb.YarnRegisteredPeerDetail{
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
	return &carrierapipb.GetRegisteredNodeListResponse{
		Status: 0,
		Msg:    backend.OK,
		Nodes:  datas,
	}, nil
}

func (svr *Server) SetJobNode(ctx context.Context, req *carrierapipb.SetJobNodeRequest) (*carrierapipb.SetJobNodeResponse, error) {

	if "" == strings.Trim(req.GetInternalIp(), "") {
		return &carrierapipb.SetJobNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require internal Ip"}, nil
	}

	if "" == strings.Trim(req.GetInternalPort(), "") {
		return &carrierapipb.SetJobNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require internal port"}, nil
	}

	if "" == strings.Trim(req.GetExternalIp(), "") {
		return &carrierapipb.SetJobNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require external Ip"}, nil
	}

	if "" == strings.Trim(req.GetExternalPort(), "") {
		return &carrierapipb.SetJobNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require external port"}, nil
	}

	node := &carrierapipb.YarnRegisteredPeerDetail{
		InternalIp:   req.GetInternalIp(),
		InternalPort: req.GetInternalPort(),
		ExternalIp:   req.GetExternalIp(),
		ExternalPort: req.GetExternalPort(),
		ConnState:    carrierapipb.ConnState_ConnState_UnConnected,
	}
	node.Id = strings.Join([]string{discovery.JobNodeConsulServiceIdPrefix, req.GetInternalIp(), req.GetInternalPort()}, discovery.ConsulServiceIdSeparator)

	status, err := svr.B.SetRegisterNode(carrierapipb.PrefixTypeJobNode, node)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:SetJobNode failed, jobNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}",
			node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())

		errMsg := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s", backend.ErrSetJobNode.Error(), "SetJobNode",
			node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())
		return &carrierapipb.SetJobNodeResponse{Status: backend.ErrSetJobNode.ErrCode(), Msg: errMsg}, nil
	}

	log.Debugf("RPC-API:SetJobNode succeed, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}, connStats: {%d} return jobNodeId: {%s}",
		req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort(), status, node.GetId())
	return &carrierapipb.SetJobNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Node: &carrierapipb.YarnRegisteredPeerDetail{
			Id:           node.GetId(),
			InternalIp:   node.GetInternalIp(),
			InternalPort: node.GetInternalPort(),
			ExternalIp:   node.GetExternalIp(),
			ExternalPort: node.GetExternalPort(),
			ConnState:    status,
		},
	}, nil
}

func (svr *Server) UpdateJobNode(ctx context.Context, req *carrierapipb.UpdateJobNodeRequest) (*carrierapipb.SetJobNodeResponse, error) {

	if "" == strings.Trim(req.GetId(), "") {
		return &carrierapipb.SetJobNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require id of job node"}, nil
	}

	if "" == strings.Trim(req.GetInternalIp(), "") {
		return &carrierapipb.SetJobNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require internal Ip"}, nil
	}

	if "" == strings.Trim(req.GetInternalPort(), "") {
		return &carrierapipb.SetJobNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require internal port"}, nil
	}

	if "" == strings.Trim(req.GetExternalIp(), "") {
		return &carrierapipb.SetJobNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require external Ip"}, nil
	}

	if "" == strings.Trim(req.GetExternalPort(), "") {
		return &carrierapipb.SetJobNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require external port"}, nil
	}

	node := &carrierapipb.YarnRegisteredPeerDetail{
		Id:           req.GetId(),
		InternalIp:   req.GetInternalIp(),
		InternalPort: req.GetInternalPort(),
		ExternalIp:   req.GetExternalIp(),
		ExternalPort: req.GetExternalPort(),
		ConnState:    carrierapipb.ConnState_ConnState_UnConnected,
	}
	status, err := svr.B.UpdateRegisterNode(carrierapipb.PrefixTypeJobNode, node)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:UpdateJobNode failed, jobNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}",
			req.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())

		errMsg := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s", backend.ErrSetJobNode.Error(), "UpdateJobNode",
			req.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())

		return &carrierapipb.SetJobNodeResponse{Status: backend.ErrSetJobNode.ErrCode(), Msg: errMsg}, nil
	}

	log.Debugf("RPC-API:UpdateJobNode succeed, jobNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}, connStats: {%d}",
		req.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort(), status)
	return &carrierapipb.SetJobNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Node: &carrierapipb.YarnRegisteredPeerDetail{
			Id:           node.GetId(),
			InternalIp:   node.GetInternalIp(),
			InternalPort: node.GetInternalPort(),
			ExternalIp:   node.GetExternalIp(),
			ExternalPort: node.GetExternalPort(),
			ConnState:    status,
		},
	}, nil
}

func (svr *Server) DeleteJobNode(ctx context.Context, req *carrierapipb.DeleteRegisteredNodeRequest) (*carriertypespb.SimpleResponse, error) {

	if "" == strings.Trim(req.GetId(), "") {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require id of job node"}, nil
	}

	if err := svr.B.DeleteRegisterNode(carrierapipb.PrefixTypeJobNode, req.GetId()); nil != err {
		log.WithError(err).Errorf("RPC-API:DeleteJobNode failed, jobNodeId: {%s}", req.GetId())

		errMsg := fmt.Sprintf("%s, %s", backend.ErrDeleteJobNode.Error(), req.GetId())
		return &carriertypespb.SimpleResponse{Status: backend.ErrDeleteJobNode.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:DeleteJobNode succeed, jobNodeId: {%s}", req.GetId())
	return &carriertypespb.SimpleResponse{Status: 0, Msg: backend.OK}, nil
}

func (svr *Server) GetJobNodeList(ctx context.Context, req *emptypb.Empty) (*carrierapipb.GetRegisteredNodeListResponse, error) {
	list, err := svr.B.GetRegisterNodeList(carrierapipb.PrefixTypeJobNode)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Error("RPC-API:GetJobNodeList failed")
		return &carrierapipb.GetRegisteredNodeListResponse{Status: backend.ErrQueryJobNodeList.ErrCode(), Msg: backend.ErrQueryJobNodeList.Error()}, nil
	}
	jobs := make([]*carrierapipb.YarnRegisteredPeer, len(list))
	for i, v := range list {
		d := &carrierapipb.YarnRegisteredPeer{
			NodeType: carrierapipb.NodeType_NodeType_JobNode,
			NodeDetail: &carrierapipb.YarnRegisteredPeerDetail{
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
	return &carrierapipb.GetRegisteredNodeListResponse{
		Status: 0,
		Msg:    backend.OK,
		Nodes:  jobs,
	}, nil
}

func (svr *Server) GenerateObServerProxyWalletAddress(ctx context.Context, req *emptypb.Empty) (*carrierapipb.GenerateObServerProxyWalletAddressResponse, error) {

	wallet, err := svr.B.GenerateObServerProxyWalletAddress()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GenerateObServerProxyWalletAddress failed")
		return &carrierapipb.GenerateObServerProxyWalletAddressResponse{Status: backend.ErrGenerateObserverProxyWallet.ErrCode(), Msg: backend.ErrGenerateObserverProxyWallet.Error()}, nil
	}
	log.Debugf("RPC-API:GenerateObServerProxyWalletAddress Succeed, wallet: {%s}", wallet)
	return &carrierapipb.GenerateObServerProxyWalletAddressResponse{
		Status:  0,
		Msg:     backend.OK,
		Address: wallet,
	}, nil
}
