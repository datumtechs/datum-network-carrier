package metadata

import (
	"context"
	"errors"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"google.golang.org/protobuf/types/known/emptypb"
)

//func (svr *Server) GetMetadataDetail(ctx context.Context, req *pb.GetMetadataDetailRequest) (*pb.GetMetadataDetailResponse, error) {
//	if req.GetMetadataId() == "" {
//		return nil, errors.New("required metadataId")
//	}
//	metadataDetail, err := svr.B.GetMetadataDetail(req.GetIdentityId(), req.GetMetadataId())
//	if nil != err {
//		log.WithError(err).Error("RPC-API:GetMetadataDetail failed")
//		return nil, ErrGetMetadataDetail
//	}
//	return metadataDetail, nil
//}

func (svr *Server) GetTotalMetadataDetailList(ctx context.Context, req *emptypb.Empty) (*pb.GetTotalMetadataDetailListResponse, error) {
	metadataList, err := svr.B.GetTotalMetadataDetailList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetTotalMetadataDetailList failed")
		return nil, ErrGetMetadataDetailList
	}
	respList := make([]*pb.GetTotalMetadataDetailResponse, len(metadataList))
	for i, metadataDetail := range metadataList {
		respList[i] = metadataDetail
	}
	log.Debugf("Query all org's metadata list, len: {%d}", len(respList))
	return &pb.GetTotalMetadataDetailListResponse{
		Status:       0,
		Msg:          backend.OK,
		MetadataList: respList,
	}, nil
}

func (svr *Server) GetSelfMetadataDetailList(ctx context.Context, req *emptypb.Empty) (*pb.GetSelfMetadataDetailListResponse, error) {
	metadataList, err := svr.B.GetSelfMetadataDetailList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetSelfMetadataDetailList failed")
		return nil, ErrGetMetadataDetailList
	}
	respList := make([]*pb.GetSelfMetadataDetailResponse, len(metadataList))
	for i, metadataDetail := range metadataList {
		respList[i] = metadataDetail
	}
	log.Debugf("Query current org's metadata list, len: {%d}", len(respList))
	return &pb.GetSelfMetadataDetailListResponse{
		Status:       0,
		Msg:          backend.OK,
		MetadataList: respList,
	}, nil
}

func (svr *Server) PublishMetadata(ctx context.Context, req *pb.PublishMetadataRequest) (*pb.PublishMetadataResponse, error) {
	if req.GetInformation() == nil {
		return nil, errors.New("required information")
	}
	if req.GetInformation().GetMetadataSummary()== nil {
		return nil, errors.New("required metadata summary")
	}
	if req.GetInformation().GetMetadataColumns() == nil {
		return nil, errors.New("required columnMeta of information")
	}

	metadataMsg := types.NewMetadataMessageFromRequest(req)

	err := svr.B.SendMsg(metadataMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:PublishMetadata failed")
		return nil, ErrSendMetadataMsg
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

	metadataRevokeMsg := types.NewMetadataRevokeMessageFromRequest(req)
	//metaDataRevokeMsg.GetCreateAt = uint64(timeutils.UnixMsec())

	err := svr.B.SendMsg(metadataRevokeMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:RevokeMetadata failed")
		return nil, ErrSendMetadataRevokeMsg
	}
	log.Debugf("RPC-API:RevokeMetadata succeed, metadataId: {%s}", req.GetMetadataId())
	return &apicommonpb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}


// GetMetadataUsedTaskIdList (GetMetadataUsedTaskIdListRequest) returns (GetMetadataUsedTaskIdListResponse) {
func (svr *Server) GetMetadataUsedTaskIdList (ctx context.Context, req *pb.GetMetadataUsedTaskIdListRequest) (*pb.GetMetadataUsedTaskIdListResponse, error) {

	return &pb.GetMetadataUsedTaskIdListResponse{
		Status:     0,
		Msg:        backend.OK,
		TaskIds:    []string{},
	}, nil
}