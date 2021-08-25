package metadata

import (
	"context"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
)

func (svr *MetaDataServiceServer) GetMetaDataDetailListByOwner(ctx context.Context, req *pb.GetMetaDataDetailListByOwnerRequest) (*pb.GetMetaDataDetailListResponse, error) {
	metadataList, err := svr.B.GetMetaDataDetailListByOwner(req.IdentityId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetMetaDataDetailListByOwner failed, identityId: {%s}", req.IdentityId)
		return nil, ErrGetMetaDataDetailList
	}
	respList := make([]*pb.GetMetaDataDetailResponse, len(metadataList))
	for i, metadata := range metadataList {
		respList[i] = metadata
	}
	log.Debugf("RPC-API:GetMetaDataDetailListByOwner succeed, identityId: {%s}, metadataList len: {%d}", req.IdentityId, len(respList))
	return &pb.GetMetaDataDetailListResponse{
		Status:       0,
		Msg:          backend.OK,
		MetaDataList: respList,
	}, nil
}