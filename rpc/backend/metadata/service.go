package metadata

import (
	"context"
	"errors"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (svr *Server) GetMetadataDetail(ctx context.Context, req *pb.GetMetadataDetailRequest) (*pb.GetMetadataDetailResponse, error) {
	if req.IdentityId == "" {
		return nil, errors.New("required identity")
	}
	if req.MetadataId == "" {
		return nil, errors.New("required metadataId")
	}
	metaDataDetail, err := svr.B.GetMetadataDetail(req.IdentityId, req.MetadataId)
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
	if req.Information == nil {
		return nil, errors.New("required information")
	}
	if req.Information.MetadataSummary == nil {
		return nil, errors.New("required metadata summary")
	}
	if req.Information.MetadataColumns == nil {
		return nil, errors.New("required columnMeta of information")
	}

	metaDataMsg := types.NewMetadataMessageFromRequest(req)
	//metaDataMsg.Data.CreateAt = uint64(timeutils.UnixMsec())

	ColumnMetas := make([]*libtypes.MetadataColumn, len(req.Information.MetadataColumns))
	for i, v := range req.Information.MetadataColumns {
		ColumnMeta := &libtypes.MetadataColumn{
			CIndex:   v.CIndex,
			CName:    v.CName,
			CType:    v.CType,
			CSize:    v.CSize,
			CComment: v.CComment,
		}
		ColumnMetas[i] = ColumnMeta
	}
	metaDataMsg.Data.Information.ColumnMetas = ColumnMetas
	metaDataId := metaDataMsg.SetMetadataId()

	err := svr.B.SendMsg(metaDataMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:PublishMetadata failed")
		return nil, ErrSendMetadataMsg
	}
	log.Debugf("RPC-API:PublishMetadata succeed, originId: {%s}, return metadataId: {%s}", req.Information.MetadataSummary.OriginId, metaDataId)
	return &pb.PublishMetadataResponse{
		Status:     0,
		Msg:        backend.OK,
		MetadataId: metaDataId,
	}, nil
}

func (svr *Server) RevokeMetadata(ctx context.Context, req *pb.RevokeMetadataRequest) (*apipb.SimpleResponse, error) {
	metaDataRevokeMsg := types.NewMetadataRevokeMessageFromRequest(req)
	//metaDataRevokeMsg.CreateAt = uint64(timeutils.UnixMsec())

	err := svr.B.SendMsg(metaDataRevokeMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:RevokeMetadata failed")
		return nil, ErrSendMetadataRevokeMsg
	}
	log.Debugf("RPC-API:RevokeMetadata succeed, metadataId: {%s}", req.MetadataId)
	return &apipb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}
