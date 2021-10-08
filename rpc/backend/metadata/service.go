package metadata

import (
	"context"
	"fmt"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

//func (svr *Server) GetMetadataDetail(ctx context.Context, req *pb.GetMetadataDetailRequest) (*pb.GetMetadataDetailResponse, error) {
//	if req.GetMetadataId() == "" {
//		return nil, errors.New("required metadataId")
//	}
//	metadataDetail, err := svr.B.GetMetadataDetail(req.QueryIdentityId(), req.GetMetadataId())
//	if nil != err {
//		log.WithError(err).Error("RPC-API:GetMetadataDetail failed")
//		return nil, ErrGetMetadataDetail
//	}
//	return metadataDetail, nil
//}

func (svr *Server) GetGlobalMetadataDetailList(ctx context.Context, req *emptypb.Empty) (*pb.GetGlobalMetadataDetailListResponse, error) {
	metadataList, err := svr.B.GetGlobalMetadataDetailList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetGlobalMetadataDetailList failed")
		return nil, ErrGetMetadataDetailList
	}
	respList := make([]*pb.GetGlobalMetadataDetailResponse, len(metadataList))
	for i, metadataDetail := range metadataList {
		respList[i] = metadataDetail
	}
	log.Debugf("Query all org's metadata list, len: {%d}", len(respList))
	return &pb.GetGlobalMetadataDetailListResponse{
		Status:       0,
		Msg:          backend.OK,
		MetadataList: respList,
	}, nil
}

func (svr *Server) GetLocalMetadataDetailList(ctx context.Context, req *emptypb.Empty) (*pb.GetLocalMetadataDetailListResponse, error) {
	metadataList, err := svr.B.GetLocalMetadataDetailList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetLocalMetadataDetailList failed")
		return nil, ErrGetMetadataDetailList
	}
	respList := make([]*pb.GetLocalMetadataDetailResponse, len(metadataList))
	for i, metadataDetail := range metadataList {
		respList[i] = metadataDetail
	}
	log.Debugf("Query current org's metadata list, len: {%d}", len(respList))
	return &pb.GetLocalMetadataDetailListResponse{
		Status:       0,
		Msg:          backend.OK,
		MetadataList: respList,
	}, nil
}

func (svr *Server) PublishMetadata(ctx context.Context, req *pb.PublishMetadataRequest) (*pb.PublishMetadataResponse, error) {
	if req.GetInformation() == nil {
		return nil, ErrReqInfoForPublishMetadata
	}
	if req.GetInformation().GetMetadataSummary() == nil {
		return nil, ErrReqMetaSummaryForPublishMetadata
	}
	if req.GetInformation().GetMetadataColumns() == nil {
		return nil, ErrReqMetaColumnsForPublishMetadata
	}

	metadataMsg := types.NewMetadataMessageFromRequest(req)

	err := svr.B.SendMsg(metadataMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:PublishMetadata failed")

		errMsg := fmt.Sprintf(ErrSendMetadataMsg.Msg,
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

		errMsg := fmt.Sprintf(ErrSendMetadataRevokeMsg.Msg, req.GetMetadataId())
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
		errMsg := fmt.Sprintf(ErrReqListForMetadataUsedTaskIdList.Msg, req.GetIdentityId(), req.GetMetadataId())
		return nil, backend.NewRpcBizErr(ErrReqListForMetadataUsedTaskIdList.Code, errMsg)
	}
	log.Debugf("RPC-API:GetMetadataUsedTaskIdList succeed, taskIds len: {%d}", len(taskIds))
	return &pb.GetMetadataUsedTaskIdListResponse{
		Status:  0,
		Msg:     backend.OK,
		TaskIds: taskIds,
	}, nil
}
