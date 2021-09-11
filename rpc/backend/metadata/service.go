package metadata

import (
	"context"
	"errors"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (svr *Server) GetMetadataDetail(ctx context.Context, req *pb.GetMetadataDetailRequest) (*pb.GetMetadataDetailResponse, error) {
	if req.GetIdentityId() == "" {
		return nil, errors.New("required identity")
	}
	if req.GetMetadataId() == "" {
		return nil, errors.New("required metadataId")
	}
	metaDataDetail, err := svr.B.GetMetadataDetail(req.GetIdentityId(), req.GetMetadataId())
	if nil != err {
		log.WithError(err).Error("RPC-API:GetMetadataDetail failed")
		return nil, ErrGetMetadataDetail
	}
	return metaDataDetail, nil
}

func (svr *Server) GetMetadataDetailList(ctx context.Context, req *emptypb.Empty) (*pb.GetMetadataDetailListResponse, error) {
	metaDataList, err := svr.B.GetMetadataDetailList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetMetadataDetailList failed")
		return nil, ErrGetMetadataDetailList
	}
	respList := make([]*pb.GetMetadataDetailResponse, len(metaDataList))
	for i, metaDataDetail := range metaDataList {
		respList[i] = metaDataDetail
	}
	log.Debugf("Query all org's metaData list, len: {%d}", len(respList))
	return &pb.GetMetadataDetailListResponse{
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

	metaDataMsg := types.NewMetadataMessageFromRequest(req)

	err := svr.B.SendMsg(metaDataMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:PublishMetadata failed")
		return nil, ErrSendMetadataMsg
	}
	log.Debugf("RPC-API:PublishMetadata succeed, originId: {%s}, return metadataId: {%s}", req.GetInformation().GetMetadataSummary().GetOriginId(), metaDataMsg.GetMetadataId())
	return &pb.PublishMetadataResponse{
		Status:     0,
		Msg:        backend.OK,
		MetadataId: metaDataMsg.GetMetadataId(),
	}, nil
}

func (svr *Server) RevokeMetadata(ctx context.Context, req *pb.RevokeMetadataRequest) (*apipb.SimpleResponse, error) {

	metaDataRevokeMsg := types.NewMetadataRevokeMessageFromRequest(req)
	//metaDataRevokeMsg.GetCreateAt = uint64(timeutils.UnixMsec())

	err := svr.B.SendMsg(metaDataRevokeMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:RevokeMetadata failed")
		return nil, ErrSendMetadataRevokeMsg
	}
	log.Debugf("RPC-API:RevokeMetadata succeed, metadataId: {%s}", req.GetMetadataId())
	return &apipb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}
