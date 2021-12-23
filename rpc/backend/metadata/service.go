package metadata

import (
	"context"
	"fmt"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strings"
)

func (svr *Server) GetGlobalMetadataDetailList(ctx context.Context, req *pb.GetGlobalMetadataDetailListRequest) (*pb.GetGlobalMetadataDetailListResponse, error) {
	metadataList, err := svr.B.GetGlobalMetadataDetailList(req.GetLastUpdated(), backend.DefaultPageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetGlobalMetadataDetailList failed")
		return nil, ErrGetMetadataDetailList
	}
	log.Debugf("Query all org's metadata list, len: {%d}", len(metadataList))
	return &pb.GetGlobalMetadataDetailListResponse{
		Status:       0,
		Msg:          backend.OK,
		MetadataList: metadataList,
	}, nil
}

func (svr *Server) GetLocalMetadataDetailList(ctx context.Context, req *pb.GetLocalMetadataDetailListRequest) (*pb.GetLocalMetadataDetailListResponse, error) {
	metadataList, err := svr.B.GetLocalMetadataDetailList(req.GetLastUpdated(), backend.DefaultPageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetLocalMetadataDetailList failed")
		return nil, ErrGetMetadataDetailList
	}
	log.Debugf("Query current org's metadata list, len: {%d}", len(metadataList))
	return &pb.GetLocalMetadataDetailListResponse{
		Status:       0,
		Msg:          backend.OK,
		MetadataList: metadataList,
	}, nil
}

func (svr *Server) PublishMetadata(ctx context.Context, req *pb.PublishMetadataRequest) (*pb.PublishMetadataResponse, error) {
	if req.GetInformation() == nil {
		return nil, ErrReqInfoForPublishMetadata
	}
	if req.GetInformation().GetMetadataSummary() == nil {
		return nil, ErrReqMetaSummaryForPublishMetadata
	}
	if len(req.GetInformation().GetMetadataColumns()) == 0 {
		return nil, ErrReqMetaColumnsForPublishMetadata
	}

	metadataMsg := types.NewMetadataMessageFromRequest(req)

	err := svr.B.SendMsg(metadataMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:PublishMetadata failed")

		errMsg := fmt.Sprintf("%s, originId: {%s}, metadataId: {%s}", ErrSendMetadataMsg.Msg,
			req.GetInformation().GetMetadataSummary().GetOriginId(), metadataMsg.GetMetadataId())
		return nil, backend.NewRpcBizErr(ErrSendMetadataMsg.Code, errMsg)
	}
	log.Debugf("RPC-API:PublishMetadata succeed, originId: {%s}, return metadataId: {%s}",
		req.GetInformation().GetMetadataSummary().GetOriginId(), metadataMsg.GetMetadataId())
	return &pb.PublishMetadataResponse{
		Status:     0,
		Msg:        backend.OK,
		MetadataId: metadataMsg.GetMetadataId(),
	}, nil
}

func (svr *Server) RevokeMetadata(ctx context.Context, req *pb.RevokeMetadataRequest) (*apicommonpb.SimpleResponse, error) {

	if "" == strings.Trim(req.GetMetadataId(), "") {
		return nil, backend.NewRpcBizErr(ErrSendMetadataRevokeMsg.Code, "require metadataId")
	}

	metadataRevokeMsg := types.NewMetadataRevokeMessageFromRequest(req)

	err := svr.B.SendMsg(metadataRevokeMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:RevokeMetadata failed")

		errMsg := fmt.Sprintf("%s, metadataId: {%s}", ErrSendMetadataRevokeMsg.Msg, req.GetMetadataId())
		return nil, backend.NewRpcBizErr(ErrSendMetadataRevokeMsg.Code, errMsg)
	}
	log.Debugf("RPC-API:RevokeMetadata succeed, metadataId: {%s}", req.GetMetadataId())
	return &apicommonpb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) GetMetadataUsedTaskIdList(ctx context.Context, req *pb.GetMetadataUsedTaskIdListRequest) (*pb.GetMetadataUsedTaskIdListResponse, error) {

	if "" == req.GetMetadataId() {
		return nil, ErrReqMetaIdForMetadataUsedTaskIdList
	}
	taskIds, err := svr.B.GetMetadataUsedTaskIdList(req.GetIdentityId(), req.GetMetadataId())
	if nil != err {
		errMsg := fmt.Sprintf("%s, IdentityId:{%s}, MetadataId:{%s}", ErrReqListForMetadataUsedTaskIdList.Msg, req.GetIdentityId(), req.GetMetadataId())
		return nil, backend.NewRpcBizErr(ErrReqListForMetadataUsedTaskIdList.Code, errMsg)
	}
	log.Debugf("RPC-API:GetMetadataUsedTaskIdList succeed, taskIds len: {%d}", len(taskIds))
	return &pb.GetMetadataUsedTaskIdListResponse{
		Status:  0,
		Msg:     backend.OK,
		TaskIds: taskIds,
	}, nil
}
