package backend

import (
	"context"
	"errors"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"time"
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
		return nil, NewRpcBizErr(ErrGetMetaDataDetailStr)
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
		Owner:       types.ConvertTaskNodeAliasToPB(metaDataDetail.Owner),
		Information: types.ConvertMetaDataInfoToPB(metaDataDetail.MetaData),
	}, nil
}

func (svr *MetaDataServiceServer) GetMetaDataDetailList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetMetaDataDetailListResponse, error) {
	metaDataList, err := svr.B.GetMetaDataDetailList()
	if nil != err {
		return nil, NewRpcBizErr(ErrGetMetaDataDetailListStr)
	}
	respList := make([]*pb.GetMetaDataDetailResponse, len(metaDataList))
	for i, metaDataDetail := range metaDataList {
		resp := &pb.GetMetaDataDetailResponse{
			Owner:       types.ConvertTaskNodeAliasToPB(metaDataDetail.Owner),
			Information: types.ConvertMetaDataInfoToPB(metaDataDetail.MetaData),
		}
		respList[i] = resp
	}

	return &pb.GetMetaDataDetailListResponse{
		Status:       0,
		Msg:          OK,
		MetaDataList: respList,
	}, nil
}

func (svr *MetaDataServiceServer) GetMetaDataDetailListByOwner(ctx context.Context, req *pb.GetMetaDataDetailListByOwnerRequest) (*pb.GetMetaDataDetailListResponse, error) {
	metaDataList, err := svr.B.GetMetaDataDetailListByOwner(req.IdentityId)
	if nil != err {
		return nil, NewRpcBizErr(ErrGetMetaDataDetailListStr)
	}
	respList := make([]*pb.GetMetaDataDetailResponse, len(metaDataList))
	for i, metaDataDetail := range metaDataList {
		resp := &pb.GetMetaDataDetailResponse{
			Owner:       types.ConvertTaskNodeAliasToPB(metaDataDetail.Owner),
			Information: types.ConvertMetaDataInfoToPB(metaDataDetail.MetaData),
		}
		respList[i] = resp
	}

	return &pb.GetMetaDataDetailListResponse{
		Status:       0,
		Msg:          OK,
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
	metaDataMsg := types.NewMetaDataMessageFromRequest(req)
	metaDataMsg.Data.CreateAt = uint64(time.Now().UnixNano())

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
	metaDataId := metaDataMsg.GetMetaDataId()

	err := svr.B.SendMsg(metaDataMsg)
	if nil != err {
		return nil, NewRpcBizErr(ErrSendMetaDataMsgStr)
	}
	return &pb.PublishMetaDataResponse{
		Status:     0,
		Msg:        OK,
		MetaDataId: metaDataId,
	}, nil
}

func (svr *MetaDataServiceServer) RevokeMetaData(ctx context.Context, req *pb.RevokeMetaDataRequest) (*pb.SimpleResponseCode, error) {
	if req == nil || req.Owner == nil {
		return nil, errors.New("required owner")
	}
	metaDataRevokeMsg := types.NewMetadataRevokeMessageFromRequest(req)
	metaDataRevokeMsg.CreateAt = uint64(time.Now().UnixNano())
	err := svr.B.SendMsg(metaDataRevokeMsg)
	if nil != err {
		return nil, NewRpcBizErr(ErrSendMetaDataRevokeMsgStr)
	}
	return &pb.SimpleResponseCode{
		Status: 0,
		Msg:    OK,
	}, nil
}
