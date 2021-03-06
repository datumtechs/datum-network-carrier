package metadata

import (
	"context"
	"errors"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
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

	identity, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishMetaData failed, query local identity failed, can not publish metadata")
		return nil, ErrSendMetaDataMsg
	}

	if identity.IdentityId() != req.Owner.IdentityId {
		return nil, errors.New("invalid identityId of req")
	}
	if identity.NodeId() != req.Owner.NodeId {
		return nil, errors.New("invalid nodeId of req")
	}
	if identity.Name() != req.Owner.Name {
		return nil, errors.New("invalid nodeName of req")
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

	identity, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokeMetaData failed, query local identity failed, can not revoke metadata")
		return nil, ErrSendMetaDataRevokeMsg
	}

	if identity.IdentityId() != req.Owner.IdentityId {
		return nil, errors.New("invalid identityId of req")
	}
	if identity.NodeId() != req.Owner.NodeId {
		return nil, errors.New("invalid nodeId of req")
	}
	if identity.Name() != req.Owner.Name {
		return nil, errors.New("invalid nodeName of req")
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
