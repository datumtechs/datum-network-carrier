package metadata

import (
	"context"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func (svr *MetaDataServiceServer) GetMetaDataDetailListByOwner(ctx context.Context, req *pb.GetMetaDataDetailListByOwnerRequest) (*pb.GetMetaDataDetailListResponse, error) {
	metaDataList, err := svr.B.GetMetaDataDetailListByOwner(req.IdentityId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetMetaDataDetailListByOwner failed, identityId: {%s}", req.IdentityId)
		return nil, backend.NewRpcBizErr(ErrGetMetaDataDetailListStr)
	}
	respList := make([]*pb.GetMetaDataDetailResponse, len(metaDataList))
	for i, metaDataDetail := range metaDataList {
		resp := &pb.GetMetaDataDetailResponse{
			Owner:       types.ConvertNodeAliasToPB(metaDataDetail.Owner),
			Information: types.ConvertMetaDataInfoToPB(metaDataDetail.MetaData),
		}
		respList[i] = resp
	}
	log.Debugf("RPC-API:GetMetaDataDetailListByOwner succeed, identityId: {%s}, metaDataList len: {%d}", req.IdentityId, len(respList))
	return &pb.GetMetaDataDetailListResponse{
		Status:       0,
		Msg:          backend.OK,
		MetaDataList: respList,
	}, nil
}