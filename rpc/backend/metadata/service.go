package metadata

import (
	"context"
	"fmt"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	libcommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (svr *Server) GetGlobalMetadataDetailList(ctx context.Context, req *pb.GetGlobalMetadataDetailListRequest) (*pb.GetGlobalMetadataDetailListResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = backend.DefaultPageSize
	}
	metadataList, err := svr.B.GetGlobalMetadataDetailList(req.GetLastUpdated(), pageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetGlobalMetadataDetailList failed")
		return &pb.GetGlobalMetadataDetailListResponse{ Status: backend.ErrQueryMetadataDetailList.ErrCode(), Msg: backend.ErrQueryMetadataDetailList.Error()}, nil
	}
	log.Debugf("Query all org's metadata list, len: {%d}", len(metadataList))
	return &pb.GetGlobalMetadataDetailListResponse{
		Status:       0,
		Msg:          backend.OK,
		Metadatas: metadataList,
	}, nil
}

func (svr *Server) GetLocalMetadataDetailList(ctx context.Context, req *pb.GetLocalMetadataDetailListRequest) (*pb.GetLocalMetadataDetailListResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = backend.DefaultPageSize
	}
	metadataList, err := svr.B.GetLocalMetadataDetailList(req.GetLastUpdated(), pageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetLocalMetadataDetailList failed")
		return &pb.GetLocalMetadataDetailListResponse{ Status: backend.ErrQueryMetadataDetailList.ErrCode(), Msg: backend.ErrQueryMetadataDetailList.Error()}, nil
	}
	log.Debugf("Query current org's global metadata list, len: {%d}", len(metadataList))
	return &pb.GetLocalMetadataDetailListResponse{
		Status:       0,
		Msg:          backend.OK,
		Metadatas: metadataList,
	}, nil
}

func (svr *Server) GetLocalInternalMetadataDetailList(ctx context.Context, req *emptypb.Empty) (*pb.GetLocalMetadataDetailListResponse, error) {

	metadataList, err := svr.B.GetLocalInternalMetadataDetailList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetLocalInternalMetadataDetailList failed")
		return &pb.GetLocalMetadataDetailListResponse{ Status: backend.ErrQueryMetadataDetailList.ErrCode(), Msg: backend.ErrQueryMetadataDetailList.Error()}, nil
	}
	log.Debugf("Query current org's internal metadata list, len: {%d}", len(metadataList))
	return &pb.GetLocalMetadataDetailListResponse{
		Status:       0,
		Msg:          backend.OK,
		Metadatas: metadataList,
	}, nil
}

func (svr *Server) PublishMetadata(ctx context.Context, req *pb.PublishMetadataRequest) (*pb.PublishMetadataResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishMetadata failed, query local identity failed, can not publish metadata")
		return &pb.PublishMetadataResponse{ Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if nil == req.GetInformation() {
		return &pb.PublishMetadataResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "the metadata infomation is empty"}, nil
	}
	if req.GetInformation().GetMetadataSummary() == nil {
		return &pb.PublishMetadataResponse{ Status: backend.ErrRequireParams.ErrCode(), Msg: "the metadata summary is empty"}, nil
	}

	metadataMsg := types.NewMetadataMessageFromRequest(req)

	if err := svr.B.SendMsg(metadataMsg); nil != err {
		log.WithError(err).Error("RPC-API:PublishMetadata failed")

		errMsg := fmt.Sprintf("%s, metadataId: {%s}", backend.ErrPublishMetadataMsg.Error(), metadataMsg.GetMetadataId())
		return &pb.PublishMetadataResponse{ Status: backend.ErrPublishMetadataMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:PublishMetadata succeed, return metadataId: {%s}", metadataMsg.GetMetadataId())
	return &pb.PublishMetadataResponse{
		Status:     0,
		Msg:        backend.OK,
		MetadataId: metadataMsg.GetMetadataId(),
	}, nil
}

func (svr *Server) RevokeMetadata(ctx context.Context, req *pb.RevokeMetadataRequest) (*libcommonpb.SimpleResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokeMetadata failed, query local identity failed, can not revoke metadata")
		return &libcommonpb.SimpleResponse { Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if "" == strings.Trim(req.GetMetadataId(), "") {
		return &libcommonpb.SimpleResponse { Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataId"}, nil
	}

	metadataRevokeMsg := types.NewMetadataRevokeMessageFromRequest(req)

	if err := svr.B.SendMsg(metadataRevokeMsg); nil != err {
		log.WithError(err).Error("RPC-API:RevokeMetadata failed")

		errMsg := fmt.Sprintf("%s, metadataId: {%s}", backend.ErrRevokeMetadataMsg.Error(), req.GetMetadataId())
		return &libcommonpb.SimpleResponse { Status: backend.ErrRevokeMetadataMsg.ErrCode(), Msg: errMsg }, nil
	}
	log.Debugf("RPC-API:RevokeMetadata succeed, metadataId: {%s}", req.GetMetadataId())
	return &libcommonpb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) GetMetadataUsedTaskIdList(ctx context.Context, req *pb.GetMetadataUsedTaskIdListRequest) (*pb.GetMetadataUsedTaskIdListResponse, error) {

	if "" == req.GetMetadataId() {
		return &pb.GetMetadataUsedTaskIdListResponse { Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataId"}, nil
	}
	taskIds, err := svr.B.GetMetadataUsedTaskIdList(req.GetIdentityId(), req.GetMetadataId())
	if nil != err {
		errMsg := fmt.Sprintf("%s, IdentityId:{%s}, MetadataId:{%s}", backend.ErrQueryMetadataUsedTaskIdList.Error(), req.GetIdentityId(), req.GetMetadataId())
		return &pb.GetMetadataUsedTaskIdListResponse { Status: backend.ErrQueryMetadataUsedTaskIdList.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:GetMetadataUsedTaskIdList succeed, taskIds len: {%d}", len(taskIds))
	return &pb.GetMetadataUsedTaskIdListResponse{
		Status:  0,
		Msg:     backend.OK,
		TaskIds: taskIds,
	}, nil
}
