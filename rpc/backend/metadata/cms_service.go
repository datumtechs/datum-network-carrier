package metadata

import (
	"context"
	"errors"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
)

//func (svr *MetaDataServiceServer) GetMetaDataSummaryList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetMetaDataSummaryListResponse, error) {
//	return nil, nil
//}
//func (svr *MetaDataServiceServer) GetMetaDataSummaryByState(ctx context.Context, req *pb.GetMetaDataSummaryByStateRequest) (*pb.GetMetaDataSummaryListResponse, error) {
//	return nil, nil
//}
//func (svr *MetaDataServiceServer) GetMetaDataSummaryByOwner(ctx context.Context, req *pb.GetMetaDataSummaryByOwnerRequest) (*pb.GetMetaDataSummaryListResponse, error) {
//	return nil, nil
//}
//func (svr *MetaDataServiceServer) GetMetaDataDetail(ctx context.Context, req *pb.GetMetaDataDetailRequest) (*pb.GetMetaDataDetailResponse, error) {
//	return nil, nil
//}
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

	columns := make([]*pb.MetaDataColumnDetail, len(metaDataDetail.MetaData.ColumnMetas))
	for i, colv := range metaDataDetail.MetaData.ColumnMetas {
		column := &pb.MetaDataColumnDetail{
			Cindex:   colv.Cindex,
			Cname:    colv.Cname,
			Ctype:    colv.Ctype,
			Csize:    colv.Csize,
			Ccomment: colv.Ccomment,
		}
		columns[i] = column
	}

	return &pb.GetMetaDataDetailResponse{
		Owner:       types.ConvertNodeAliasToPB(metaDataDetail.Owner),
		Information: types.ConvertMetaDataInfoToPB(metaDataDetail.MetaData),
	}, nil
}

func (svr *MetaDataServiceServer) GetMetaDataDetailList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetMetaDataDetailListResponse, error) {
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
	if req == nil || req.Owner == nil {
		return nil, errors.New("required owner")
	}
	if req.Information == nil {
		return nil, errors.New("required information")
	}
	if req.Information.MetaDataSummary == nil {
		return nil, errors.New("required metadata summary")
	}
	if req.Information.ColumnMeta == nil {
		return nil, errors.New("required columnMeta of information")
	}

	_, err := svr.B.GetNodeIdentity()
	if rawdb.IsDBNotFoundErr(err) {
		log.Errorf("RPC-API:PublishMetaData failed, the identity was not exist, can not revoke identity")
		return nil, ErrSendMetaDataMsg
	}

	metaDataMsg := types.NewMetaDataMessageFromRequest(req)
	//metaDataMsg.Data.CreateAt = uint64(timeutils.UnixMsec())


	ColumnMetas := make([]*libtypes.ColumnMeta, len(req.Information.ColumnMeta))
	for i, v := range req.Information.ColumnMeta {
		ColumnMeta := &libtypes.ColumnMeta{
			Cindex:   v.Cindex,
			Cname:    v.Cname,
			Ctype:    v.Ctype,
			Csize:    v.Csize,
			Ccomment: v.Ccomment,
		}
		ColumnMetas[i] = ColumnMeta
	}
	metaDataMsg.Data.Information.ColumnMetas = ColumnMetas
	metaDataId := metaDataMsg.SetMetaDataId()

	//b, _ := json.Marshal(metaDataMsg)
	//log.Debugf("############ Input req: {%s}", string(b))

	err = svr.B.SendMsg(metaDataMsg)
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

func (svr *MetaDataServiceServer) RevokeMetaData(ctx context.Context, req *pb.RevokeMetaDataRequest) (*pb.SimpleResponseCode, error) {
	if req == nil || req.Owner == nil {
		return nil, errors.New("required owner")
	}

	_, err := svr.B.GetNodeIdentity()
	if rawdb.IsDBNotFoundErr(err) {
		log.Errorf("RPC-API:RevokeMetaData failed, the identity was not exist, can not revoke identity")
		return nil, ErrSendMetaDataRevokeMsg
	}

	metaDataRevokeMsg := types.NewMetadataRevokeMessageFromRequest(req)
	//metaDataRevokeMsg.CreateAt = uint64(timeutils.UnixMsec())
	err = svr.B.SendMsg(metaDataRevokeMsg)
	if nil != err {
		log.WithError(err).Error("RPC-API:RevokeMetaData failed")
		return nil, ErrSendMetaDataRevokeMsg
	}
	log.Debugf("RPC-API:RevokeMetaData succeed, metadataId: {%s}", req.MetaDataId)
	return &pb.SimpleResponseCode{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}
