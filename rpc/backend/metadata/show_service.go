package metadata

import (
	"context"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
)

func (svr *Server) GetMetadataDetailListByOwner(ctx context.Context, req *pb.GetMetadataDetailListByOwnerRequest) (*pb.GetMetadataDetailListResponse, error) {
	metadataList, err := svr.B.GetMetadataDetailListByOwner(req.IdentityId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetMetadataDetailListByOwner failed, identityId: {%s}", req.IdentityId)
		return nil, ErrGetMetadataDetailList
	}
	respList := make([]*pb.GetMetadataDetailResponse, len(metadataList))
	for i, metadata := range metadataList {
		respList[i] = metadata
	}
	log.Debugf("RPC-API:GetMetadataDetailListByOwner succeed, identityId: {%s}, metadataList len: {%d}", req.IdentityId, len(respList))
	return &pb.GetMetadataDetailListResponse{
		Status:       0,
		Msg:          backend.OK,
		MetadataList: respList,
	}, nil
}