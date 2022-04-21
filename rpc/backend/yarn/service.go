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

func (svr *Server) QueryAvailableDataNode(ctx context.Context, req *pb.QueryAvailableDataNodeRequest) (*pb.QueryAvailableDataNodeResponse, error) {

	if req.GetFileType() == libtypes.OrigindataType_OrigindataType_Unknown {
		return &pb.QueryAvailableDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown fileType"}, nil
	}

	if req.GetFileSize() == 0 {
		return &pb.QueryAvailableDataNodeResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require fileSize"}, nil
	}

	dataResourceTables, err := svr.B.QueryDataResourceTables()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryAvailableDataNode-QueryDataResourceTables failed, fileType: {%s}, fileSize: {%d}",
			req.GetFileType(), req.GetFileSize())

		errMsg := fmt.Sprintf("%s, call QueryDataResourceTables() failed, %s, %d", backend.ErrQueryAvailableDataNode.Error(), req.GetFileType(), req.GetFileSize())
		return &pb.QueryAvailableDataNodeResponse{Status: backend.ErrQueryAvailableDataNode.ErrCode(), Msg: errMsg}, nil
	}

	var nodeId string
	for _, resource := range dataResourceTables {
		log.Debugf("QueryAvailableDataNode, dataResourceTable: %s, need disk: %d, remain disk: %d", resource.String(), req.GetFileSize(), resource.RemainDisk())
		if req.GetFileSize() < resource.RemainDisk() {
			nodeId = resource.GetNodeId()
			break
		}
	}

	if "" == strings.Trim(nodeId, "") {
		log.Errorf("RPC-API:QueryAvailableDataNode, not found available dataNodeId of dataNode, fileType: {%s}, fileSize: {%d}, dataNodeId: {%s}",
			req.GetFileType(), req.GetFileSize(), nodeId)

		errMsg := fmt.Sprintf("%s, not found available dataNodeId of dataNode, %s, %d, %s", backend.ErrQueryAvailableDataNode.Error(),
			req.GetFileType(), req.GetFileSize(), nodeId)
		return &pb.QueryAvailableDataNodeResponse{Status: backend.ErrQueryAvailableDataNode.ErrCode(), Msg: errMsg}, nil
	}

	dataNode, err := svr.B.GetRegisterNode(pb.PrefixTypeDataNode, nodeId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryAvailableDataNode-QueryRegisterNode failed, fileType: {%s}, fileSize: {%d}, dataNodeId: {%s}",
			req.GetFileType(), req.GetFileSize(), nodeId)

		errMsg := fmt.Sprintf("%s, call QueryRegisterNode() failed, %s, %d, %s", backend.ErrQueryAvailableDataNode.Error(),
			req.GetFileType(), req.GetFileSize(), nodeId)
		return &pb.QueryAvailableDataNodeResponse{Status: backend.ErrQueryAvailableDataNode.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:QueryAvailableDataNode succeed, fileType: {%s}, fileSize: {%d}, return dataNodeId: {%s}, dataNodeIp: {%s}, dataNodePort: {%s}",
		req.GetFileType(), req.GetFileSize(), dataNode.GetId(), dataNode.GetInternalIp(), dataNode.GetInternalPort())

	return &pb.QueryAvailableDataNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Information: &pb.QueryAvailableDataNode{
			Ip:   dataNode.GetInternalIp(),
			Port: dataNode.GetInternalPort(),
		},
	}, nil
}
func (svr *Server) QueryFilePosition(ctx context.Context, req *pb.QueryFilePositionRequest) (*pb.QueryFilePositionResponse, error) {

	if "" == req.GetOriginId() {
		return &pb.QueryFilePositionResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require originId"}, nil
	}

	dataResourceFileUpload, err := svr.B.QueryDataResourceFileUpload(req.GetOriginId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryFilePosition-QueryDataResourceFileUpload failed, originId: {%s}", req.GetOriginId())

		errMsg := fmt.Sprintf("%s, call QueryDataResourceFileUpload() failed, %s", backend.ErrQueryFilePosition.Error(), req.GetOriginId())
		return &pb.QueryFilePositionResponse{Status: backend.ErrQueryFilePosition.ErrCode(), Msg: errMsg}, nil
	}
	dataNode, err := svr.B.GetRegisterNode(pb.PrefixTypeDataNode, dataResourceFileUpload.GetNodeId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryFilePosition-QueryRegisterNode failed, originId: {%s}, dataNodeId: {%s}",
			req.GetOriginId(), dataResourceFileUpload.GetNodeId())

		errMsg := fmt.Sprintf("%s, call QueryRegisterNode() failed, %s, %s", backend.ErrQueryFilePosition.Error(),
			req.GetOriginId(), dataResourceFileUpload.GetNodeId())
		return &pb.QueryFilePositionResponse{Status: backend.ErrQueryFilePosition.ErrCode(), Msg: errMsg}, nil
	}

	log.Debugf("RPC-API:QueryFilePosition Succeed, originId: {%s}, return dataNodeIp: {%s}, dataNodePort: {%s}, filePath: {%s}",
		req.GetOriginId(), dataNode.GetInternalIp(), dataNode.GetInternalPort(), dataResourceFileUpload.GetDataPath())

	return &pb.QueryFilePositionResponse{
		Status: 0,
		Msg:    backend.OK,
		Information: &pb.QueryFilePosition{
			Ip:       dataNode.GetInternalIp(),
			Port:     dataNode.GetInternalPort(),
			DataPath: dataResourceFileUpload.GetDataPath(),
		},
	}, nil
}

func (svr *Server) GetTaskResultFileSummary(ctx context.Context, req *pb.GetTaskResultFileSummaryRequest) (*pb.GetTaskResultFileSummaryResponse, error) {

	if "" == req.GetTaskId() {
		return &pb.GetTaskResultFileSummaryResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require taskId"}, nil
	}

	summary, err := svr.B.QueryTaskResultFileSummary(req.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetTaskResultFileSummary-QueryTaskResultFileSummary failed, taskId: {%s}", req.GetTaskId())

		errMsg := fmt.Sprintf("%s, call QueryTaskResultFileSummary() failed, %s", backend.ErrQueryTaskResultFileSummary.Error(), req.GetTaskId())
		return &pb.GetTaskResultFileSummaryResponse{Status: backend.ErrQueryTaskResultFileSummary.ErrCode(), Msg: errMsg}, nil
	}
	dataNode, err := svr.B.GetRegisterNode(pb.PrefixTypeDataNode, summary.GetNodeId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetTaskResultFileSummary-QueryRegisterNode failed, taskId: {%s}, dataNodeId: {%s}",
			req.GetTaskId(), summary.GetNodeId())

		errMsg := fmt.Sprintf("%s, call QueryRegisterNode() failed, %s, %s", backend.ErrQueryTaskResultFileSummary.Error(),
			req.GetTaskId(), summary.GetNodeId())
		return &pb.GetTaskResultFileSummaryResponse{Status: backend.ErrQueryTaskResultFileSummary.ErrCode(), Msg: errMsg}, nil
	}

	log.Debugf("RPC-API:GetTaskResultFileSummary Succeed, taskId: {%s}, return dataNodeIp: {%s}, dataNodePort: {%s}, metadataId: {%s}, originId: {%s}, fileName: {%s}, filePath: {%s}",
		req.GetTaskId(), dataNode.GetInternalIp(), dataNode.GetInternalPort(), summary.GetMetadataId(), summary.GetOriginId(), summary.GetMetadataName(), summary.GetDataPath())

	return &pb.GetTaskResultFileSummaryResponse{
		Status: 0,
		Msg:    backend.OK,
		Information: &pb.GetTaskResultFileSummary{
			TaskId:       summary.GetTaskId(),
			MetadataName: summary.GetMetadataName(),
			MetadataId:   summary.GetMetadataId(),
			OriginId:     summary.GetOriginId(),
			DataPath:     summary.GetDataPath(),
			Ip:           dataNode.GetInternalIp(),
			Port:         dataNode.GetInternalPort(),
			Extra:        summary.GetExtra(),
		},
	}, nil
}

func (svr *Server) GetTaskResultFileSummaryList(ctx context.Context, empty *emptypb.Empty) (*pb.GetTaskResultFileSummaryListResponse, error) {
	taskResultFileSummaryArr, err := svr.B.QueryTaskResultFileSummaryList()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetTaskResultFileSummaryList-QueryTaskResultFileSummaryList failed")
		return &pb.GetTaskResultFileSummaryListResponse{Status: backend.ErrQueryTaskResultFileSummaryList.ErrCode(), Msg: backend.ErrQueryTaskResultFileSummaryList.Error()}, nil
	}

	arr := make([]*pb.GetTaskResultFileSummary, 0)
	for _, summary := range taskResultFileSummaryArr {
		dataNode, err := svr.B.GetRegisterNode(pb.PrefixTypeDataNode, summary.GetNodeId())
		if nil != err {
			log.WithError(err).Errorf("RPC-API:GetTaskResultFileSummaryList-QueryRegisterNode failed, taskId: {%s}, dataNodeId: {%s}",
				summary.GetTaskId(), summary.GetNodeId())
			continue
		}
		arr = append(arr, &pb.GetTaskResultFileSummary{
			TaskId:       summary.GetTaskId(),
			MetadataName: summary.GetMetadataName(),
			MetadataId:   summary.GetMetadataId(),
			OriginId:     summary.GetOriginId(),
			DataPath:     summary.GetDataPath(),
			Ip:           dataNode.GetInternalIp(),
			Port:         dataNode.GetInternalPort(),
			Extra:        summary.GetExtra(),
		})
	}

	log.Debugf("RPC-API:GetTaskResultFileSummaryList Succeed, task result file summary list len: {%d}", len(arr))
	return &pb.GetTaskResultFileSummaryListResponse{
		Status:          0,
		Msg:             backend.OK,
		TaskResultFiles: arr,
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
