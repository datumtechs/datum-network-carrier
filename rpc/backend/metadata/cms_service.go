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

func (svr *MetaDataServiceServer) GetMetaDataDetail(ctx context.Context, req *pb.GetMetaDataDetailRequest) (*pb.GetMetaDataDetailResponse, error) {
	if req.IdentityId == "" {
		return nil, errors.New("required identity")
	}
	if req.MetaDataId == "" {
		return nil, errors.New("required metadataId")
	}
	metaDataDetail, err := svr.B.GetMetaDataDetail(req.IdentityId, req.MetaDataId)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetMetaDataDetail failed")
		return nil, ErrGetMetaDataDetail
	}

	columns := make([]*libtypes.MetadataColumn, len(metaDataDetail.MetaData.ColumnMetas))
	for i, colv := range metaDataDetail.MetaData.ColumnMetas {
		column := &libtypes.MetadataColumn{
			CIndex:   colv.CIndex,
			CName:    colv.CName,
			CType:    colv.CType,
			CSize:    colv.CSize,
			CComment: colv.CComment,
		}
		columns[i] = column
	}

	return &pb.GetMetaDataDetailResponse{
		Owner:       types.ConvertNodeAliasToPB(metaDataDetail.Owner),
		Information: types.ConvertMetaDataInfoToPB(metaDataDetail.MetaData),
	}, nil
}

func (svr *MetaDataServiceServer) GetMetaDataDetailList(ctx context.Context, req *emptypb.Empty) (*pb.GetMetaDataDetailListResponse, error) {
	metaDataList, err := svr.B.GetMetaDataDetailList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetMetaDataDetailList failed")
		return nil, ErrGetMetaDataDetailList
	}
	respList := make([]*pb.GetMetaDataDetailResponse, len(metaDataList))
	for i, metaDataDetail := range metaDataList {
		resp := &pb.GetMetaDataDetailResponse{
			Owner:       types.ConvertNodeAliasToPB(metaDataDetail.Owner),
			Information: types.ConvertMetaDataInfoToPB(metaDataDetail.MetaData),
		}
		respList[i] = resp
	}
	log.Debugf("Query all org's metaData list, len: {%d}", len(respList))
	return &pb.GetMetaDataDetailListResponse{
		Status:       0,
		Msg:          backend.OK,
		MetaDataList: respList,
	}, nil
}

func (svr *MetaDataServiceServer) PublishMetaData(ctx context.Context, req *pb.PublishMetaDataRequest) (*pb.PublishMetaDataResponse, error) {
	if req.Information == nil {
		return nil, errors.New("required information")
	}
	if req.Information.MetaDataSummary == nil {
		return nil, errors.New("required metadata summary")
	}
	if req.Information.MetadataColumnList == nil {
		return nil, errors.New("required columnMeta of information")
	}

	metaDataMsg := types.NewMetaDataMessageFromRequest(req)
	//metaDataMsg.Data.CreateAt = uint64(timeutils.UnixMsec())

	ColumnMetas := make([]*libtypes.MetadataColumn, len(req.Information.MetadataColumnList))
	for i, v := range req.Information.MetadataColumnList {
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
	metaDataId := metaDataMsg.SetMetaDataId()

	err := svr.B.SendMsg(metaDataMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:PublishMetaData failed")
		return nil, ErrSendMetaDataMsg
	}
	log.Debugf("RPC-API:PublishMetaData succeed, originId: {%s}, return metadataId: {%s}", req.Information.MetaDataSummary.OriginId, metaDataId)
	return &pb.PublishMetaDataResponse{
		Status:     0,
		Msg:        backend.OK,
		MetaDataId: metaDataId,
	}, nil
}

func (svr *MetaDataServiceServer) RevokeMetaData(ctx context.Context, req *pb.RevokeMetaDataRequest) (*apipb.SimpleResponse, error) {
	metaDataRevokeMsg := types.NewMetadataRevokeMessageFromRequest(req)
	//metaDataRevokeMsg.CreateAt = uint64(timeutils.UnixMsec())

	err := svr.B.SendMsg(metaDataRevokeMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:RevokeMetaData failed")
		return nil, ErrSendMetaDataRevokeMsg
	}
	log.Debugf("RPC-API:RevokeMetaData succeed, metadataId: {%s}", req.MetaDataId)
	return &apipb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}
