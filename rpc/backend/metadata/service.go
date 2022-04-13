package metadata

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
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
		return &pb.GetGlobalMetadataDetailListResponse{Status: backend.ErrQueryMetadataDetailList.ErrCode(), Msg: backend.ErrQueryMetadataDetailList.Error()}, nil
	}
	log.Debugf("Query all org's metadata list, len: {%d}", len(metadataList))
	return &pb.GetGlobalMetadataDetailListResponse{
		Status:    0,
		Msg:       backend.OK,
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
		return &pb.GetLocalMetadataDetailListResponse{Status: backend.ErrQueryMetadataDetailList.ErrCode(), Msg: backend.ErrQueryMetadataDetailList.Error()}, nil
	}
	log.Debugf("Query current org's global metadata list, len: {%d}", len(metadataList))
	return &pb.GetLocalMetadataDetailListResponse{
		Status:    0,
		Msg:       backend.OK,
		Metadatas: metadataList,
	}, nil
}

func (svr *Server) GetLocalInternalMetadataDetailList(ctx context.Context, req *emptypb.Empty) (*pb.GetLocalMetadataDetailListResponse, error) {

	metadataList, err := svr.B.GetLocalInternalMetadataDetailList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetLocalInternalMetadataDetailList failed")
		return &pb.GetLocalMetadataDetailListResponse{Status: backend.ErrQueryMetadataDetailList.ErrCode(), Msg: backend.ErrQueryMetadataDetailList.Error()}, nil
	}
	log.Debugf("Query current org's internal metadata list, len: {%d}", len(metadataList))
	return &pb.GetLocalMetadataDetailListResponse{
		Status:    0,
		Msg:       backend.OK,
		Metadatas: metadataList,
	}, nil
}

func (svr *Server) PublishMetadata(ctx context.Context, req *pb.PublishMetadataRequest) (*pb.PublishMetadataResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishMetadata failed, query local identity failed, can not publish metadata")
		return &pb.PublishMetadataResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if nil == req.GetInformation() {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "the metadata infomation is empty"}, nil
	}
	if nil == req.GetInformation() {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "the metadata summary is empty"}, nil
	}

	metadataMsg := types.NewMetadataMessageFromRequest(req)

	if err := svr.B.SendMsg(metadataMsg); nil != err {
		log.WithError(err).Error("RPC-API:PublishMetadata failed")

		errMsg := fmt.Sprintf("%s, metadataId: {%s}", backend.ErrPublishMetadataMsg.Error(), metadataMsg.GetMetadataId())
		return &pb.PublishMetadataResponse{Status: backend.ErrPublishMetadataMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:PublishMetadata succeed, return metadataId: {%s}", metadataMsg.GetMetadataId())
	return &pb.PublishMetadataResponse{
		Status:     0,
		Msg:        backend.OK,
		MetadataId: metadataMsg.GetMetadataId(),
	}, nil
}

func (svr *Server) RevokeMetadata(ctx context.Context, req *pb.RevokeMetadataRequest) (*libtypes.SimpleResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokeMetadata failed, query local identity failed, can not revoke metadata")
		return &libtypes.SimpleResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if "" == strings.Trim(req.GetMetadataId(), "") {
		return &libtypes.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataId"}, nil
	}

	metadataRevokeMsg := types.NewMetadataRevokeMessageFromRequest(req)

	if err := svr.B.SendMsg(metadataRevokeMsg); nil != err {
		log.WithError(err).Error("RPC-API:RevokeMetadata failed")

		errMsg := fmt.Sprintf("%s, metadataId: {%s}", backend.ErrRevokeMetadataMsg.Error(), req.GetMetadataId())
		return &libtypes.SimpleResponse{Status: backend.ErrRevokeMetadataMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:RevokeMetadata succeed, metadataId: {%s}", req.GetMetadataId())
	return &libtypes.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) GetMetadataUsedTaskIdList(ctx context.Context, req *pb.GetMetadataUsedTaskIdListRequest) (*pb.GetMetadataUsedTaskIdListResponse, error) {

	if "" == req.GetMetadataId() {
		return &pb.GetMetadataUsedTaskIdListResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataId"}, nil
	}
	taskIds, err := svr.B.GetMetadataUsedTaskIdList(req.GetIdentityId(), req.GetMetadataId())
	if nil != err {
		errMsg := fmt.Sprintf("%s, IdentityId:{%s}, MetadataId:{%s}", backend.ErrQueryMetadataUsedTaskIdList.Error(), req.GetIdentityId(), req.GetMetadataId())
		return &pb.GetMetadataUsedTaskIdListResponse{Status: backend.ErrQueryMetadataUsedTaskIdList.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:GetMetadataUsedTaskIdList succeed, taskIds len: {%d}", len(taskIds))
	return &pb.GetMetadataUsedTaskIdListResponse{
		Status:  0,
		Msg:     backend.OK,
		TaskIds: taskIds,
	}, nil
}

func (svr *Server) BindDataTokenAddress(ctx context.Context, req *pb.BindDataTokenAddressRequest) (*libtypes.SimpleResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:BindDataTokenAddress failed, query local identity failed, can not publish metadata")
		return &libtypes.SimpleResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if "" == strings.Trim(req.GetMetadataId(), "") {
		return &libtypes.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataId"}, nil
	}
	if "" == strings.Trim(req.GetTokenAddress(), "") {
		return &libtypes.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require dataToken address"}, nil
	}

	metadataId := strings.Trim(req.GetMetadataId(), "")

	var metadata *types.Metadata

	metadata, err = svr.B.GetInternalMetadataDetail(metadataId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:BindDataTokenAddress failed, check is internal metadata failed, metadataId: {%s}, dataTokenAddress: {%s}", req.GetMetadataId(), req.GetTokenAddress())
		return &libtypes.SimpleResponse{Status: backend.ErrBindDataTokenAddress.ErrCode(), Msg: fmt.Sprintf("%s, check is internal metadata failed", backend.ErrBindDataTokenAddress.Error())}, nil
	}
	if nil != metadata {
		log.Errorf("RPC-API:BindDataTokenAddress failed, internal metadata be not able to bind datatoken address, metadataId: {%s}, dataTokenAddress: {%s}", req.GetMetadataId(), req.GetTokenAddress())
		return &libtypes.SimpleResponse{Status: backend.ErrBindDataTokenAddress.ErrCode(), Msg: fmt.Sprintf("%s, internal metadata be not able to bind datatoken address", backend.ErrBindDataTokenAddress.Error())}, nil
	}

	metadata, err = svr.B.GetMetadataDetail(metadataId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:BindDataTokenAddress failed, query metadata failed, metadataId: {%s}, dataTokenAddress: {%s}", req.GetMetadataId(), req.GetTokenAddress())
		return &libtypes.SimpleResponse{Status: backend.ErrBindDataTokenAddress.ErrCode(), Msg: fmt.Sprintf("%s, query metadata failed", backend.ErrBindDataTokenAddress.Error())}, nil
	}
	if nil == metadata {
		log.Errorf("RPC-API:BindDataTokenAddress failed, not found metadata")
		return &libtypes.SimpleResponse{Status: backend.ErrBindDataTokenAddress.ErrCode(), Msg: fmt.Sprintf("%s, not found metadata", backend.ErrBindDataTokenAddress.Error())}, nil
	}
	if "" != metadata.GetData().GetTokenAddress() {
		log.Errorf("RPC-API:BindDataTokenAddress failed, the metadata had tokenAddress already")
		return &libtypes.SimpleResponse{Status: backend.ErrBindDataTokenAddress.ErrCode(), Msg: fmt.Sprintf("%s, the metadata had tokenAddress already", backend.ErrBindDataTokenAddress.Error())}, nil

	}

	metadata.GetData().TokenAddress = strings.Trim(req.GetTokenAddress(), "")
	if err := svr.B.UpdateGlobalMetadata(metadata); nil != err {
		log.WithError(err).Errorf("RPC-API:BindDataTokenAddress failed, update global metadata failed, metadataId: {%s}, dataTokenAddress: {%s}", req.GetMetadataId(), req.GetTokenAddress())
		return &libtypes.SimpleResponse{Status: backend.ErrBindDataTokenAddress.ErrCode(), Msg: fmt.Sprintf("%s, update global metadata failed", backend.ErrBindDataTokenAddress.Error())}, nil
	}
	log.Debugf("RPC-API:BindDataTokenAddress succeed, metadataId: {%s}, dataTokenAddress: {%s}", req.GetMetadataId(), req.GetTokenAddress())
	return &libtypes.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) PublishMetadataByInteranlMetadata(ctx context.Context, req *pb.PublishMetadataByInteranlMetadataRequest) (*pb.PublishMetadataResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishMetadataByInteranlMetadata failed, query local identity failed, can not publish metadata")
		return &pb.PublishMetadataResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if nil == req.GetInformation() {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataInfomation"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetMetadataId(), "") {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataId"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetMetadataName(), "") {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataName"}, nil
	}
	if req.GetInformation().GetMetadataType() == libtypes.MetadataType_MetadataType_Unknown {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown metadataType"}, nil
	}
	// DataHash
	if "" == strings.Trim(req.GetInformation().GetDesc(), "") {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require desc"}, nil
	}
	if req.GetInformation().GetDataType() == libtypes.OrigindataType_OrigindataType_Unknown {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown dataType"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetIndustry(), "") {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require industry"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetMetadataOption(), "") {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataOption"}, nil
	}
	// AllowExpose
	// TokenAddress

	metadataId := strings.Trim(req.GetInformation().GetMetadataId(), "")

	metadata, err := svr.B.GetInternalMetadataDetail(metadataId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishMetadataByInteranlMetadata failed, query internal metadata failed, metadataId {%s}", metadataId)
		return &pb.PublishMetadataResponse{Status: backend.ErrPublishMetadataMsg.ErrCode(), Msg: fmt.Sprintf("%s, query internal metadata failed", backend.ErrPublishMetadataMsg.Error())}, nil
	}
	if nil == metadata {
		log.Errorf("RPC-API:PublishMetadataByInteranlMetadata failed, not found internal metadata, metadataId {%s}", metadataId)
		return &pb.PublishMetadataResponse{Status: backend.ErrPublishMetadataMsg.ErrCode(), Msg: fmt.Sprintf("%s, not found internal metadata", backend.ErrPublishMetadataMsg.Error())}, nil
	}

	// build metadata msg
	metadataMsg := &types.MetadataMsg{
		MetadataSummary: &libtypes.MetadataSummary{
			MetadataId:     req.GetInformation().GetMetadataId(),
			MetadataName:   req.GetInformation().GetMetadataName(),
			MetadataType:   req.GetInformation().GetMetadataType(),
			DataHash:       metadata.GetData().GetDataHash(),
			Desc:           req.GetInformation().GetDesc(),
			DataType:       req.GetInformation().GetDataType(),
			Industry:       req.GetInformation().GetIndustry(),
			State:          req.GetInformation().GetState(),
			PublishAt:      req.GetInformation().GetPublishAt(),
			UpdateAt:       req.GetInformation().GetUpdateAt(),
			Nonce:          req.GetInformation().GetNonce(),
			MetadataOption: req.GetInformation().GetMetadataOption(),
			AllowExpose:    req.GetInformation().GetAllowExpose(),
			//TokenAddress:   req.GetInformation().GetTokenAddress(),
		},
		CreateAt: timeutils.UnixMsecUint64(),
	}

	if err := svr.B.SendMsg(metadataMsg); nil != err {
		log.WithError(err).Error("RPC-API:PublishMetadataByInteranlMetadata failed")

		errMsg := fmt.Sprintf("%s, metadataId: {%s}", backend.ErrPublishMetadataMsg.Error(), metadataMsg.GetMetadataId())
		return &pb.PublishMetadataResponse{Status: backend.ErrPublishMetadataMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:PublishMetadataByInteranlMetadata succeed, return metadataId: {%s}", metadataMsg.GetMetadataId())
	return &pb.PublishMetadataResponse{
		Status:     0,
		Msg:        backend.OK,
		MetadataId: metadataMsg.GetMetadataId(),
	}, nil
}

func (svr *Server) PublishMetadataByTaskResultFile(ctx context.Context, req *pb.PublishMetadataByTaskResultFileRequest) (*pb.PublishMetadataResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishMetadataByTaskResultFile failed, query local identity failed, can not publish metadata")
		return &pb.PublishMetadataResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if "" == strings.Trim(req.GetTaskId(), "") {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require taskId"}, nil
	}
	if nil == req.GetInformation() {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataInfomation"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetMetadataId(), "") {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataId"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetMetadataName(), "") {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataName"}, nil
	}
	if req.GetInformation().GetMetadataType() == libtypes.MetadataType_MetadataType_Unknown {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown metadataType"}, nil
	}
	// DataHash
	if "" == strings.Trim(req.GetInformation().GetDesc(), "") {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require desc"}, nil
	}
	if req.GetInformation().GetDataType() == libtypes.OrigindataType_OrigindataType_Unknown {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown dataType"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetIndustry(), "") {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require industry"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetMetadataOption(), "") {
		return &pb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataOption"}, nil
	}
	// AllowExpose
	// TokenAddress
	taskResultFileSummary, err := svr.B.QueryTaskResultFileSummary(req.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishMetadataByTaskResultFile-QueryTaskResultFileSummary failed, taskId: {%s}", req.GetTaskId())

		errMsg := fmt.Sprintf("%s, call QueryTaskResultFileSummary() failed, %s", backend.ErrPublishMetadataMsg.Error(), req.GetTaskId())
		return &pb.PublishMetadataResponse{Status: backend.ErrPublishMetadataMsg.ErrCode(), Msg: errMsg}, nil
	}

	metadataId := taskResultFileSummary.GetMetadataId()

	metadata, err := svr.B.GetInternalMetadataDetail(metadataId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishMetadataByTaskResultFile failed, query internal metadata failed, metadataId {%s}", metadataId)
		return &pb.PublishMetadataResponse{Status: backend.ErrPublishMetadataMsg.ErrCode(), Msg: fmt.Sprintf("%s, query internal metadata failed", backend.ErrPublishMetadataMsg.Error())}, nil
	}
	if nil == metadata {
		log.Errorf("RPC-API:PublishMetadataByTaskResultFile failed, not found internal metadata, metadataId {%s}", metadataId)
		return &pb.PublishMetadataResponse{Status: backend.ErrPublishMetadataMsg.ErrCode(), Msg: fmt.Sprintf("%s, not found internal metadata", backend.ErrPublishMetadataMsg.Error())}, nil
	}

	// build metadata msg
	metadataMsg := &types.MetadataMsg{
		MetadataSummary: &libtypes.MetadataSummary{
			MetadataId:     req.GetInformation().GetMetadataId(),
			MetadataName:   req.GetInformation().GetMetadataName(),
			MetadataType:   req.GetInformation().GetMetadataType(),
			DataHash:       metadata.GetData().GetDataHash(),
			Desc:           req.GetInformation().GetDesc(),
			DataType:       req.GetInformation().GetDataType(),
			Industry:       req.GetInformation().GetIndustry(),
			State:          req.GetInformation().GetState(),
			PublishAt:      req.GetInformation().GetPublishAt(),
			UpdateAt:       req.GetInformation().GetUpdateAt(),
			Nonce:          req.GetInformation().GetNonce(),
			MetadataOption: req.GetInformation().GetMetadataOption(),
			AllowExpose:    req.GetInformation().GetAllowExpose(),
			//TokenAddress:   req.GetInformation().GetTokenAddress(),
		},
		CreateAt: timeutils.UnixMsecUint64(),
	}

	if err := svr.B.SendMsg(metadataMsg); nil != err {
		log.WithError(err).Error("RPC-API:PublishMetadataByTaskResultFile failed")

		errMsg := fmt.Sprintf("%s, metadataId: {%s}", backend.ErrPublishMetadataMsg.Error(), metadataMsg.GetMetadataId())
		return &pb.PublishMetadataResponse{Status: backend.ErrPublishMetadataMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:PublishMetadataByTaskResultFile succeed, return metadataId: {%s}", metadataMsg.GetMetadataId())
	return &pb.PublishMetadataResponse{
		Status:     0,
		Msg:        backend.OK,
		MetadataId: metadataMsg.GetMetadataId(),
	}, nil
}
