package yarn

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (svr *Server) GetNodeInfo(ctx context.Context, req *emptypb.Empty) (*pb.GetNodeInfoResponse, error) {
	node, err := svr.B.GetNodeInfo()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetNodeInfo failed")
		return nil, ErrGetNodeInfo
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

func (svr *Server) GetRegisteredPeers(ctx context.Context, req *pb.GetRegisteredPeersRequest) (*pb.GetRegisteredPeersResponse, error) {

	if req.GetNodeType() == pb.NodeType_NodeType_YarnNode ||
		req.GetNodeType() == pb.NodeType_NodeType_SeedNode {
		log.Errorf("RPC-API:GetRegisteredPeers failed, invalid nodeType from req, nodeType: {%s}", req.GetNodeType().String())
		return nil, ErrGetRegisteredPeers
	}

	registerNodes, err := svr.B.GetRegisteredPeers(req.GetNodeType())
	if nil != err {
		log.WithError(err).Error("RPC-API:GetRegisteredPeers failed")
		return nil, ErrGetRegisteredPeers
	}
	log.Debugf("RPC-API:GetRegisteredPeers succeed, node len: {%d}", len(registerNodes))
	return &pb.GetRegisteredPeersResponse{
		Status: 0,
		Msg:    backend.OK,
		Nodes:  registerNodes,
	}, nil
}

func (svr *Server) SetSeedNode(ctx context.Context, req *pb.SetSeedNodeRequest) (*pb.SetSeedNodeResponse, error) {

	if "" == strings.Trim(req.GetAddr(), "") {
		return nil, backend.NewRpcBizErr(ErrSetSeedNodeInfo.Code, "require addr of seedNode")
	}

	seedNode := &pb.SeedPeer{
		Addr:      req.GetAddr(),
		IsDefault: false,
		ConnState: pb.ConnState_ConnState_UnConnected,
	}
	status, err := svr.B.SetSeedNode(seedNode)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:SetSeedNode failed, addr: {%s}", seedNode.GetAddr())
		errMsg := fmt.Sprintf(ErrSetSeedNodeInfo.Msg, "SetSeedNode", seedNode.GetAddr())
		return nil, backend.NewRpcBizErr(ErrSetSeedNodeInfo.Code, errMsg)
	}
	log.Debugf("RPC-API:SetSeedNode succeed, addr: {%s}, connStatus: {%d}", req.GetAddr(), status)
	seedNode.ConnState = status
	return &pb.SetSeedNodeResponse{
		Status: 0,
		Msg:    backend.OK,
		Node: seedNode,
	}, nil
}

func (svr *Server) DeleteSeedNode(ctx context.Context, req *pb.DeleteSeedNodeRequest) (*apicommonpb.SimpleResponse, error) {

	if "" == strings.Trim(req.GetAddr(), "") {
		return nil, backend.NewRpcBizErr(ErrDeleteSeedNodeInfo.Code, "require addr of seedNode")
	}

	err := svr.B.DeleteSeedNode(req.GetAddr())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RemoveSeedNode failed, addr: {%s}", req.GetAddr())
		errMsg := fmt.Sprintf(ErrDeleteSeedNodeInfo.Msg, req.GetAddr())
		return nil, backend.NewRpcBizErr(ErrDeleteSeedNodeInfo.Code, errMsg)
	}
	log.Debugf("RPC-API:RemoveSeedNode succeed, addr: {%s}", req.GetAddr())
	return &apicommonpb.SimpleResponse{Status: 0, Msg: backend.OK}, nil
}

func (svr *Server) GetSeedNodeList(ctx context.Context, req *emptypb.Empty) (*pb.GetSeedNodeListResponse, error) {
	list, err := svr.B.GetSeedNodeList()
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Error("RPC-API: call svr.B.GetSeedNodeList() failed")
		return nil, ErrGetSeedNodeList
	}
	return &pb.GetSeedNodeListResponse{
		Status: 0,
		Msg:    backend.OK,
		Nodes:  list,
	}, nil
}

func (svr *Server) SetDataNode(ctx context.Context, req *pb.SetDataNodeRequest) (*pb.SetDataNodeResponse, error) {

	if "" == strings.Trim(req.GetInternalIp(), "") {
		return nil, backend.NewRpcBizErr(ErrSetDataNodeInfo.Code, "require internal Ip")
	}

	if "" == strings.Trim(req.GetInternalPort(), "") {
		return nil, backend.NewRpcBizErr(ErrSetDataNodeInfo.Code, "require internal port")
	}

	if "" == strings.Trim(req.GetExternalIp(), "") {
		return nil, backend.NewRpcBizErr(ErrSetDataNodeInfo.Code, "require external Ip")
	}

	if "" == strings.Trim(req.GetExternalPort(), "") {
		return nil, backend.NewRpcBizErr(ErrSetDataNodeInfo.Code, "require external port")
	}

	node := &pb.YarnRegisteredPeerDetail{
		InternalIp:   req.GetInternalIp(),
		InternalPort: req.GetInternalPort(),
		ExternalIp:   req.GetExternalIp(),
		ExternalPort: req.GetExternalPort(),
		ConnState:    pb.ConnState_ConnState_UnConnected,
	}
	node.GenDataNodeId()
	status, err := svr.B.SetRegisterNode(pb.PrefixTypeDataNode, node)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:SetDataNode failed, dataNodeId:{%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}",
			node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())

		errMsg := fmt.Sprintf(ErrSetDataNodeInfo.Msg, "SetDataNode",
			node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())
		return nil, backend.NewRpcBizErr(ErrSetDataNodeInfo.Code, errMsg)
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
		return nil, backend.NewRpcBizErr(ErrSetDataNodeInfo.Code, "require id of data node")
	}

	if "" == strings.Trim(req.GetInternalIp(), "") {
		return nil, backend.NewRpcBizErr(ErrSetDataNodeInfo.Code, "require internal Ip")
	}

	if "" == strings.Trim(req.GetInternalPort(), "") {
		return nil, backend.NewRpcBizErr(ErrSetDataNodeInfo.Code, "require internal port")
	}

	if "" == strings.Trim(req.GetExternalIp(), "") {
		return nil, backend.NewRpcBizErr(ErrSetDataNodeInfo.Code, "require external Ip")
	}

	if "" == strings.Trim(req.GetExternalPort(), "") {
		return nil, backend.NewRpcBizErr(ErrSetDataNodeInfo.Code, "require external port")
	}

	node := &pb.YarnRegisteredPeerDetail{
		Id:           req.GetId(),
		InternalIp:   req.GetInternalIp(),
		InternalPort: req.GetInternalPort(),
		ExternalIp:   req.GetExternalIp(),
		ExternalPort: req.GetExternalPort(),
		ConnState:    pb.ConnState_ConnState_UnConnected,
	}
	// delete and insert.
	//svr.B.RemoveRegisterNode(types.PrefixTypeDataNode, node.Id)
	//status, err := svr.B.SetRegisterNode(types.PrefixTypeDataNode, node)
	status, err := svr.B.UpdateRegisterNode(pb.PrefixTypeDataNode, node)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:UpdateDataNode failed, dataNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}",
			node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())

		errMsg := fmt.Sprintf(ErrSetDataNodeInfo.Msg, "UpdateDataNode",
			node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())
		return nil, backend.NewRpcBizErr(ErrSetDataNodeInfo.Code, errMsg)
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

func (svr *Server) DeleteDataNode(ctx context.Context, req *pb.DeleteRegisteredNodeRequest) (*apicommonpb.SimpleResponse, error) {

	if "" == strings.Trim(req.GetId(), "") {
		return nil, backend.NewRpcBizErr(ErrDeleteDataNodeInfo.Code, "require id of data node")
	}

	if err := svr.B.DeleteRegisterNode(pb.PrefixTypeDataNode, req.GetId()); nil != err {
		log.WithError(err).Errorf("RPC-API:DeleteDataNode failed, dataNodeId: {%s}", req.GetId())

		errMsg := fmt.Sprintf(ErrDeleteDataNodeInfo.Msg, req.GetId())
		return nil, backend.NewRpcBizErr(ErrDeleteDataNodeInfo.Code, errMsg)
	}
	log.Debugf("RPC-API:DeleteDataNode succeed, dataNodeId: {%s}", req.GetId())
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
		return nil, backend.NewRpcBizErr(ErrSetJobNodeInfo.Code, "require internal Ip")
	}

	if "" == strings.Trim(req.GetInternalPort(), "") {
		return nil, backend.NewRpcBizErr(ErrSetJobNodeInfo.Code, "require internal port")
	}

	if "" == strings.Trim(req.GetExternalIp(), "") {
		return nil, backend.NewRpcBizErr(ErrSetJobNodeInfo.Code, "require external Ip")
	}

	if "" == strings.Trim(req.GetExternalPort(), "") {
		return nil, backend.NewRpcBizErr(ErrSetJobNodeInfo.Code, "require external port")
	}

	node := &pb.YarnRegisteredPeerDetail{
		InternalIp:   req.GetInternalIp(),
		InternalPort: req.GetInternalPort(),
		ExternalIp:   req.GetExternalIp(),
		ExternalPort: req.GetExternalPort(),
		ConnState:    pb.ConnState_ConnState_UnConnected,
	}
	node.GenJobNodeId()
	status, err := svr.B.SetRegisterNode(pb.PrefixTypeJobNode, node)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:SetJobNode failed, jobNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}",
			node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())

		errMsg := fmt.Sprintf(ErrSetJobNodeInfo.Msg, "SetJobNode",
			node.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())
		return nil, backend.NewRpcBizErr(ErrSetJobNodeInfo.Code, errMsg)
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
		return nil, backend.NewRpcBizErr(ErrSetJobNodeInfo.Code, "require id of job node")
	}

	if "" == strings.Trim(req.GetInternalIp(), "") {
		return nil, backend.NewRpcBizErr(ErrSetJobNodeInfo.Code, "require internal Ip")
	}

	if "" == strings.Trim(req.GetInternalPort(), "") {
		return nil, backend.NewRpcBizErr(ErrSetJobNodeInfo.Code, "require internal port")
	}

	if "" == strings.Trim(req.GetExternalIp(), "") {
		return nil, backend.NewRpcBizErr(ErrSetJobNodeInfo.Code, "require external Ip")
	}

	if "" == strings.Trim(req.GetExternalPort(), "") {
		return nil, backend.NewRpcBizErr(ErrSetJobNodeInfo.Code, "require external port")
	}

	node := &pb.YarnRegisteredPeerDetail{
		Id:           req.GetId(),
		InternalIp:   req.GetInternalIp(),
		InternalPort: req.GetInternalPort(),
		ExternalIp:   req.GetExternalIp(),
		ExternalPort: req.GetExternalPort(),
		ConnState:    pb.ConnState_ConnState_UnConnected,
	}
	//svr.B.RemoveRegisterNode(types.PrefixTypeJobNode, node.Id)
	//status, err := svr.B.SetRegisterNode(types.PrefixTypeJobNode, node)
	status, err := svr.B.UpdateRegisterNode(pb.PrefixTypeJobNode, node)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:UpdateJobNode failed, jobNodeId: {%s}, internalIp: {%s}, internalPort: {%s}, externalIp: {%s}, externalPort: {%s}",
			req.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())

		errMsg := fmt.Sprintf(ErrSetJobNodeInfo.Msg, "UpdateJobNode",
			req.GetId(), req.GetInternalIp(), req.GetInternalPort(), req.GetExternalIp(), req.GetExternalPort())
		return nil, backend.NewRpcBizErr(ErrSetJobNodeInfo.Code, errMsg)
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

func (svr *Server) DeleteJobNode(ctx context.Context, req *pb.DeleteRegisteredNodeRequest) (*apicommonpb.SimpleResponse, error) {

	if "" == strings.Trim(req.GetId(), "") {
		return nil, backend.NewRpcBizErr(ErrDeleteJobNodeInfo.Code, "require id of job node")
	}

	if err := svr.B.DeleteRegisterNode(pb.PrefixTypeJobNode, req.GetId()); nil != err {
		log.WithError(err).Errorf("RPC-API:DeleteJobNode failed, jobNodeId: {%s}", req.GetId())

		errMsg := fmt.Sprintf(ErrDeleteJobNodeInfo.Msg, req.GetId())
		return nil, backend.NewRpcBizErr(ErrDeleteJobNodeInfo.Code, errMsg)
	}
	log.Debugf("RPC-API:DeleteJobNode succeed, jobNodeId: {%s}", req.GetId())
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
			NodeType: pb.NodeType_NodeType_JobNode,
			NodeDetail: &pb.YarnRegisteredPeerDetail{
				Id:           v.GetId(),
				InternalIp:   v.GetInternalIp(),
				InternalPort: v.GetInternalPort(),
				ExternalIp:   v.GetExternalIp(),
				ExternalPort: v.GetExternalPort(),
				ConnState:    v.GetConnState(),
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

	if req.GetFileType() == apicommonpb.OriginFileType_FileType_Unknown {
		return nil, fmt.Errorf("invalid fileType")
	}

	if req.GetFileSize() == 0 {
		return nil, fmt.Errorf("require fileSize")
	}

	dataResourceTables, err := svr.B.QueryDataResourceTables()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryAvailableDataNode-QueryDataResourceTables failed, fileType: {%s}, fileSize: {%d}",
			req.GetFileType(), req.GetFileSize())

		errMsg := fmt.Sprintf(ErrQueryDataResourceTableList.Msg, req.GetFileType(), req.GetFileSize())
		return nil, backend.NewRpcBizErr(ErrQueryDataResourceTableList.Code, errMsg)
	}

	var nodeId string
	for _, resource := range dataResourceTables {
		if req.GetFileSize() < resource.RemainDisk() {
			nodeId = resource.GetNodeId()
			break
		}
	}

	dataNode, err := svr.B.GetRegisterNode(pb.PrefixTypeDataNode, nodeId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryAvailableDataNode-QueryRegisterNode failed, fileType: {%s}, fileSize: {%d}, dataNodeId: {%s}",
			req.GetFileType(), req.GetFileSize(), nodeId)

		errMsg := fmt.Sprintf(ErrGetDataNodeInfoForQueryAvailableDataNode.Msg,
			req.GetFileType(), req.GetFileSize(), nodeId)
		return nil, backend.NewRpcBizErr(ErrGetDataNodeInfoForQueryAvailableDataNode.Code, errMsg)
	}
	log.Debugf("RPC-API:QueryAvailableDataNode succeed, fileType: {%s}, fileSize: {%d}, return dataNodeId: {%s}, dataNodeIp: {%s}, dataNodePort: {%s}",
		req.GetFileType(), req.GetFileSize(), dataNode.GetId(), dataNode.GetInternalIp(), dataNode.GetInternalPort())

	return &pb.QueryAvailableDataNodeResponse{
		Ip:   dataNode.GetInternalIp(),
		Port: dataNode.GetInternalPort(),
	}, nil
}
func (svr *Server) QueryFilePosition(ctx context.Context, req *pb.QueryFilePositionRequest) (*pb.QueryFilePositionResponse, error) {

	if "" == req.GetOriginId() {
		return nil, ErrReqOriginIdForQueryFilePosition
	}

	dataResourceFileUpload, err := svr.B.QueryDataResourceFileUpload(req.GetOriginId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryFilePosition-QueryDataResourceFileUpload failed, originId: {%s}", req.GetOriginId())

		errMsg := fmt.Sprintf(ErrQueryDataResourceDataUsed.Msg, req.GetOriginId())
		return nil, backend.NewRpcBizErr(ErrQueryDataResourceDataUsed.Code, errMsg)
	}
	dataNode, err := svr.B.GetRegisterNode(pb.PrefixTypeDataNode, dataResourceFileUpload.GetNodeId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryFilePosition-QueryRegisterNode failed, originId: {%s}, dataNodeId: {%s}",
			req.GetOriginId(), dataResourceFileUpload.GetNodeId())

		errMsg := fmt.Sprintf(ErrGetDataNodeInfoForQueryFilePosition.Msg,
			req.GetOriginId(), dataResourceFileUpload.GetNodeId())
		return nil, backend.NewRpcBizErr(ErrGetDataNodeInfoForQueryFilePosition.Code, errMsg)
	}

	log.Debugf("RPC-API:QueryFilePosition Succeed, originId: {%s}, return dataNodeIp: {%s}, dataNodePort: {%s}, filePath: {%s}",
		req.GetOriginId(), dataNode.GetInternalIp(), dataNode.GetInternalPort(), dataResourceFileUpload.GetFilePath())

	return &pb.QueryFilePositionResponse{
		Ip:       dataNode.GetInternalIp(),
		Port:     dataNode.GetInternalPort(),
		FilePath: dataResourceFileUpload.GetFilePath(),
	}, nil
}

func (svr *Server) GetTaskResultFileSummary(ctx context.Context, req *pb.GetTaskResultFileSummaryRequest) (*pb.GetTaskResultFileSummaryResponse, error) {

	if "" == req.GetTaskId() {
		return nil, ErrReqTaskIdForGetTaskResultFileSummary
	}

	taskResultFileSummary, err := svr.B.QueryTaskResultFileSummary(req.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetTaskResultFileSummary-QueryTaskResultFileSummary failed, taskId: {%s}", req.GetTaskId())

		errMsg := fmt.Sprintf(ErrQueryTaskResultFileSummary.Msg, req.GetTaskId())
		return nil, backend.NewRpcBizErr(ErrQueryTaskResultFileSummary.Code, fmt.Sprintf("%s: %s", errMsg, err))
	}
	dataNode, err := svr.B.GetRegisterNode(pb.PrefixTypeDataNode, taskResultFileSummary.GetNodeId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetTaskResultFileSummary-QueryRegisterNode failed, taskId: {%s}, dataNodeId: {%s}",
			req.GetTaskId(), taskResultFileSummary.GetNodeId())

		errMsg := fmt.Sprintf(ErrGetDataNodeInfoForTaskResultFileSummary.Msg,
			req.GetTaskId(), taskResultFileSummary.GetNodeId())
		return nil, backend.NewRpcBizErr(ErrGetDataNodeInfoForTaskResultFileSummary.Code, errMsg)
	}

	log.Debugf("RPC-API:GetTaskResultFileSummary Succeed, taskId: {%s}, return dataNodeIp: {%s}, dataNodePort: {%s}, metadataId: {%s}, originId: {%s}, fileName: {%s}, filePath: {%s}",
		req.GetTaskId(), dataNode.GetInternalIp(), dataNode.GetInternalPort(), taskResultFileSummary.GetMetadataId(), taskResultFileSummary.GetOriginId(), taskResultFileSummary.GetFileName(), taskResultFileSummary.GetFilePath())

	return &pb.GetTaskResultFileSummaryResponse{
		TaskId:     taskResultFileSummary.GetTaskId(),
		FileName:   taskResultFileSummary.GetFileName(),
		MetadataId: taskResultFileSummary.GetMetadataId(),
		OriginId:   taskResultFileSummary.GetOriginId(),
		FilePath:   taskResultFileSummary.GetFilePath(),
		Ip:         dataNode.GetInternalIp(),
		Port:       dataNode.GetInternalPort(),
	}, nil
}

func (svr *Server) GetTaskResultFileSummaryList(ctx context.Context, empty *emptypb.Empty) (*pb.GetTaskResultFileSummaryListResponse, error) {
	taskResultFileSummaryArr, err := svr.B.QueryTaskResultFileSummaryList()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetTaskResultFileSummaryList-QueryTaskResultFileSummaryList")
		return nil, ErrQueryTaskResultFileSummary
	}

	arr := make([]*pb.GetTaskResultFileSummaryResponse, 0)
	for _, summary := range taskResultFileSummaryArr {
		dataNode, err := svr.B.GetRegisterNode(pb.PrefixTypeDataNode, summary.GetNodeId())
		if nil != err {
			log.WithError(err).Errorf("RPC-API:GetTaskResultFileSummaryList-QueryRegisterNode failed, taskId: {%s}, dataNodeId: {%s}",
				summary.GetTaskId(), summary.GetNodeId())
			continue
		}
		arr = append(arr, &pb.GetTaskResultFileSummaryResponse{
			TaskId:     summary.GetTaskId(),
			FileName:   summary.GetFileName(),
			MetadataId: summary.GetMetadataId(),
			OriginId:   summary.GetOriginId(),
			FilePath:   summary.GetFilePath(),
			Ip:         dataNode.GetInternalIp(),
			Port:       dataNode.GetInternalPort(),
		})
	}

	log.Debugf("RPC-API:GetTaskResultFileSummaryList Succeed, task result file summary list len: {%d}", len(arr))
	return &pb.GetTaskResultFileSummaryListResponse{
		Status:       0,
		Msg:          backend.OK,
		MetadataList: arr,
	}, nil
}
