package metadata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/rpc/backend"
	"github.com/datumtechs/datum-network-carrier/signsuite"
	"github.com/datumtechs/datum-network-carrier/types"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (svr *Server) GetGlobalMetadataDetailList(ctx context.Context, req *carrierapipb.GetGlobalMetadataDetailListRequest) (*carrierapipb.GetGlobalMetadataDetailListResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = backend.DefaultPageSize
	}
	metadataList, err := svr.B.GetGlobalMetadataDetailList(req.GetLastUpdated(), pageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetGlobalMetadataDetailList failed")
		return &carrierapipb.GetGlobalMetadataDetailListResponse{Status: backend.ErrQueryMetadataDetailList.ErrCode(), Msg: backend.ErrQueryMetadataDetailList.Error()}, nil
	}
	log.Debugf("Query all org's metadata list, len: {%d}", len(metadataList))
	return &carrierapipb.GetGlobalMetadataDetailListResponse{
		Status:    0,
		Msg:       backend.OK,
		Metadatas: metadataList,
	}, nil
}

func (svr *Server) GetLocalMetadataDetailList(ctx context.Context, req *carrierapipb.GetLocalMetadataDetailListRequest) (*carrierapipb.GetLocalMetadataDetailListResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = backend.DefaultPageSize
	}
	metadataList, err := svr.B.GetLocalMetadataDetailList(req.GetLastUpdated(), pageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetLocalMetadataDetailList failed")
		return &carrierapipb.GetLocalMetadataDetailListResponse{Status: backend.ErrQueryMetadataDetailList.ErrCode(), Msg: backend.ErrQueryMetadataDetailList.Error()}, nil
	}
	log.Debugf("Query current org's global metadata list, len: {%d}", len(metadataList))
	return &carrierapipb.GetLocalMetadataDetailListResponse{
		Status:    0,
		Msg:       backend.OK,
		Metadatas: metadataList,
	}, nil
}

func (svr *Server) GetLocalInternalMetadataDetailList(ctx context.Context, req *emptypb.Empty) (*carrierapipb.GetLocalMetadataDetailListResponse, error) {

	metadataList, err := svr.B.GetLocalInternalMetadataDetailList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetLocalInternalMetadataDetailList failed")
		return &carrierapipb.GetLocalMetadataDetailListResponse{Status: backend.ErrQueryMetadataDetailList.ErrCode(), Msg: backend.ErrQueryMetadataDetailList.Error()}, nil
	}
	log.Debugf("Query current org's internal metadata list, len: {%d}", len(metadataList))
	return &carrierapipb.GetLocalMetadataDetailListResponse{
		Status:    0,
		Msg:       backend.OK,
		Metadatas: metadataList,
	}, nil
}

func (svr *Server) PublishMetadata(ctx context.Context, req *carrierapipb.PublishMetadataRequest) (*carrierapipb.PublishMetadataResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishMetadata failed, query local identity failed, can not publish metadata")
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if nil == req.GetInformation() {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "the metadata infomation is empty"}, nil
	}
	if nil == req.GetInformation() {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "the metadata summary is empty"}, nil
	}

	metadataMsg := types.NewMetadataMessageFromRequest(req)
	metadataMsg.GenMetadataId()

	// check from
	from, _, err := signsuite.Sender(metadataMsg.GetMetadataSummary().GetUserType(), metadataMsg.Hash(), metadataMsg.GetMetadataSummary().GetSign())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishMetadata failed, cannot fetch sender from sign, userType: {%s}, user: {%s}",
			req.GetInformation().GetUserType().String(), req.GetInformation().GetUser())
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "cannot fetch sender from sign"}, nil
	}
	if from != req.GetInformation().GetUser() {
		log.WithError(err).Errorf("RPC-API:PublishMetadata failed, sender from sign and user is not sameone, userType: {%s}, user: {%s}, sender of sign: {%s}",
			req.GetInformation().GetUserType().String(), req.GetInformation().GetUser(), from)
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "the user sign is invalid"}, nil
	}

	if err := svr.B.SendMsg(metadataMsg); nil != err {
		log.WithError(err).Error("RPC-API:PublishMetadata failed")

		errMsg := fmt.Sprintf("%s, metadataId: {%s}", backend.ErrPublishMetadataMsg.Error(), metadataMsg.GetMetadataId())
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrPublishMetadataMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:PublishMetadata succeed, return metadataId: {%s}", metadataMsg.GetMetadataId())
	return &carrierapipb.PublishMetadataResponse{
		Status:     0,
		Msg:        backend.OK,
		MetadataId: metadataMsg.GetMetadataId(),
	}, nil
}

func (svr *Server) RevokeMetadata(ctx context.Context, req *carrierapipb.RevokeMetadataRequest) (*carriertypespb.SimpleResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokeMetadata failed, query local identity failed, can not revoke metadata")
		return &carriertypespb.SimpleResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if "" == strings.Trim(req.GetMetadataId(), "") {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataId"}, nil
	}

	metadataRevokeMsg := types.NewMetadataRevokeMessageFromRequest(req)

	if err := svr.B.SendMsg(metadataRevokeMsg); nil != err {
		log.WithError(err).Error("RPC-API:RevokeMetadata failed")

		errMsg := fmt.Sprintf("%s, metadataId: {%s}", backend.ErrRevokeMetadataMsg.Error(), req.GetMetadataId())
		return &carriertypespb.SimpleResponse{Status: backend.ErrRevokeMetadataMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:RevokeMetadata succeed, metadataId: {%s}", req.GetMetadataId())
	return &carriertypespb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) GetMetadataUsedTaskIdList(ctx context.Context, req *carrierapipb.GetMetadataUsedTaskIdListRequest) (*carrierapipb.GetMetadataUsedTaskIdListResponse, error) {

	if "" == req.GetMetadataId() {
		return &carrierapipb.GetMetadataUsedTaskIdListResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataId"}, nil
	}
	taskIds, err := svr.B.GetMetadataUsedTaskIdList(req.GetIdentityId(), req.GetMetadataId())
	if nil != err {
		errMsg := fmt.Sprintf("%s, IdentityId:{%s}, MetadataId:{%s}", backend.ErrQueryMetadataUsedTaskIdList.Error(), req.GetIdentityId(), req.GetMetadataId())
		return &carrierapipb.GetMetadataUsedTaskIdListResponse{Status: backend.ErrQueryMetadataUsedTaskIdList.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:GetMetadataUsedTaskIdList succeed, taskIds len: {%d}", len(taskIds))
	return &carrierapipb.GetMetadataUsedTaskIdListResponse{
		Status:  0,
		Msg:     backend.OK,
		TaskIds: taskIds,
	}, nil
}

func (svr *Server) PublishMetadataByInteranlMetadata(ctx context.Context, req *carrierapipb.PublishMetadataByInteranlMetadataRequest) (*carrierapipb.PublishMetadataResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishMetadataByInteranlMetadata failed, query local identity failed, can not publish metadata")
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if nil == req.GetInformation() {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataInfomation"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetMetadataId(), "") {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataId"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetMetadataName(), "") {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataName"}, nil
	}
	if req.GetInformation().GetMetadataType() == commonconstantpb.MetadataType_MetadataType_Unknown {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown metadataType"}, nil
	}
	// DataHash
	if "" == strings.Trim(req.GetInformation().GetDesc(), "") {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require desc"}, nil
	}
	if req.GetInformation().GetLocationType() == commonconstantpb.DataLocationType_DataLocationType_Unknown {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown locationType"}, nil
	}
	if req.GetInformation().GetDataType() == commonconstantpb.OrigindataType_OrigindataType_Unknown {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown dataType"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetIndustry(), "") {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require industry"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetMetadataOption(), "") {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataOption"}, nil
	}

	metadataId := strings.Trim(req.GetInformation().GetMetadataId(), "")

	metadata, err := svr.B.GetInternalMetadataDetail(metadataId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishMetadataByInteranlMetadata failed, query internal metadata failed, metadataId {%s}", metadataId)
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrPublishMetadataMsg.ErrCode(), Msg: fmt.Sprintf("%s, query internal metadata failed", backend.ErrPublishMetadataMsg.Error())}, nil
	}
	if nil == metadata {
		log.Errorf("RPC-API:PublishMetadataByInteranlMetadata failed, not found internal metadata, metadataId {%s}", metadataId)
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrPublishMetadataMsg.ErrCode(), Msg: fmt.Sprintf("%s, not found internal metadata", backend.ErrPublishMetadataMsg.Error())}, nil
	}

	// build metadata msg
	metadataMsg := &types.MetadataMsg{
		MetadataSummary: &carriertypespb.MetadataSummary{
			/**
			MetadataId     string
			MetadataName   string
			MetadataType   constant.MetadataType
			DataHash       string
			Desc           string
			LocationType   constant.DataLocationTy
			DataType       constant.OrigindataType
			Industry       string
			State          constant.MetadataState
			PublishAt      uint64
			UpdateAt       uint64
			Nonce          uint64
			MetadataOption string
			// add by v0.5.0
			User                 string
			UserType             constant.UserType
			Sign                 []byte
			*/
			MetadataId:     req.GetInformation().GetMetadataId(),
			MetadataName:   req.GetInformation().GetMetadataName(),
			MetadataType:   req.GetInformation().GetMetadataType(),
			DataHash:       metadata.GetData().GetDataHash(),
			Desc:           req.GetInformation().GetDesc(),
			LocationType:   req.GetInformation().GetLocationType(),
			DataType:       req.GetInformation().GetDataType(),
			Industry:       req.GetInformation().GetIndustry(),
			State:          req.GetInformation().GetState(),
			PublishAt:      req.GetInformation().GetPublishAt(),
			UpdateAt:       req.GetInformation().GetUpdateAt(),
			Nonce:          req.GetInformation().GetNonce(),
			MetadataOption: req.GetInformation().GetMetadataOption(),
			UserType:       req.GetInformation().GetUserType(),
			User:           req.GetInformation().GetUser(),
			Sign:           req.GetInformation().GetSign(),
		},
		CreateAt: timeutils.UnixMsecUint64(),
	}
	// check from
	from, _, err := signsuite.Sender(metadataMsg.GetMetadataSummary().GetUserType(), metadataMsg.Hash(), metadataMsg.GetMetadataSummary().GetSign())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishMetadataByInteranlMetadata failed, cannot fetch sender from sign, userType: {%s}, user: {%s}",
			req.GetInformation().GetUserType().String(), req.GetInformation().GetUser())
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "cannot fetch sender from sign"}, nil
	}
	if from != req.GetInformation().GetUser() {
		log.WithError(err).Errorf("RPC-API:PublishMetadataByInteranlMetadata failed, sender from sign and user is not sameone, userType: {%s}, user: {%s}, sender of sign: {%s}",
			req.GetInformation().GetUserType().String(), req.GetInformation().GetUser(), from)
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "the user sign is invalid"}, nil
	}
	if err := svr.B.SendMsg(metadataMsg); nil != err {
		log.WithError(err).Error("RPC-API:PublishMetadataByInteranlMetadata failed")

		errMsg := fmt.Sprintf("%s, metadataId: {%s}", backend.ErrPublishMetadataMsg.Error(), metadataMsg.GetMetadataId())
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrPublishMetadataMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:PublishMetadataByInteranlMetadata succeed, return metadataId: {%s}", metadataMsg.GetMetadataId())
	return &carrierapipb.PublishMetadataResponse{
		Status:     0,
		Msg:        backend.OK,
		MetadataId: metadataMsg.GetMetadataId(),
	}, nil
}

func (svr *Server) PublishMetadataByTaskResultFile(ctx context.Context, req *carrierapipb.PublishMetadataByTaskResultFileRequest) (*carrierapipb.PublishMetadataResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishMetadataByTaskResultFile failed, query local identity failed, can not publish metadata")
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if "" == strings.Trim(req.GetTaskId(), "") {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require taskId"}, nil
	}
	if nil == req.GetInformation() {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataInfomation"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetMetadataId(), "") {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataId"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetMetadataName(), "") {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataName"}, nil
	}
	if req.GetInformation().GetMetadataType() == commonconstantpb.MetadataType_MetadataType_Unknown {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown metadataType"}, nil
	}
	// DataHash
	if "" == strings.Trim(req.GetInformation().GetDesc(), "") {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require desc"}, nil
	}
	if req.GetInformation().GetLocationType() == commonconstantpb.DataLocationType_DataLocationType_Unknown {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown locationType"}, nil
	}
	if req.GetInformation().GetDataType() == commonconstantpb.OrigindataType_OrigindataType_Unknown {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown dataType"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetIndustry(), "") {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require industry"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetMetadataOption(), "") {
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataOption"}, nil
	}
	taskResultDataSummary, err := svr.B.QueryTaskResultDataSummary(req.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishMetadataByTaskResultFile-QueryTaskResultDataSummary failed, taskId: {%s}", req.GetTaskId())

		errMsg := fmt.Sprintf("%s, call QueryTaskResultDataSummary() failed, %s", backend.ErrPublishMetadataMsg.Error(), req.GetTaskId())
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrPublishMetadataMsg.ErrCode(), Msg: errMsg}, nil
	}

	metadataId := taskResultDataSummary.GetMetadataId()

	metadata, err := svr.B.GetInternalMetadataDetail(metadataId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishMetadataByTaskResultFile failed, query internal metadata failed, metadataId {%s}", metadataId)
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrPublishMetadataMsg.ErrCode(), Msg: fmt.Sprintf("%s, query internal metadata failed", backend.ErrPublishMetadataMsg.Error())}, nil
	}
	if nil == metadata {
		log.Errorf("RPC-API:PublishMetadataByTaskResultFile failed, not found internal metadata, metadataId {%s}", metadataId)
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrPublishMetadataMsg.ErrCode(), Msg: fmt.Sprintf("%s, not found internal metadata", backend.ErrPublishMetadataMsg.Error())}, nil
	}

	// build metadata msg
	metadataMsg := &types.MetadataMsg{
		MetadataSummary: &carriertypespb.MetadataSummary{
			/**
			MetadataId     string
			MetadataName   string
			MetadataType   constant.MetadataType
			DataHash       string
			Desc           string
			LocationType   constant.DataLocationTy
			DataType       constant.OrigindataType
			Industry       string
			State          constant.MetadataState
			PublishAt      uint64
			UpdateAt       uint64
			Nonce          uint64
			MetadataOption string
			// add by v0.5.0
			User                 string
			UserType             constant.UserType
			Sign                 []byte
			*/
			MetadataId:     req.GetInformation().GetMetadataId(),
			MetadataName:   req.GetInformation().GetMetadataName(),
			MetadataType:   req.GetInformation().GetMetadataType(),
			DataHash:       metadata.GetData().GetDataHash(),
			Desc:           req.GetInformation().GetDesc(),
			LocationType:   req.GetInformation().GetLocationType(),
			DataType:       req.GetInformation().GetDataType(),
			Industry:       req.GetInformation().GetIndustry(),
			State:          req.GetInformation().GetState(),
			PublishAt:      req.GetInformation().GetPublishAt(),
			UpdateAt:       req.GetInformation().GetUpdateAt(),
			Nonce:          req.GetInformation().GetNonce(),
			MetadataOption: req.GetInformation().GetMetadataOption(),
			UserType:       req.GetInformation().GetUserType(),
			User:           req.GetInformation().GetUser(),
			Sign:           req.GetInformation().GetSign(),
		},
		CreateAt: timeutils.UnixMsecUint64(),
	}
	// check from
	from, _, err := signsuite.Sender(metadataMsg.GetMetadataSummary().GetUserType(), metadataMsg.Hash(), metadataMsg.GetMetadataSummary().GetSign())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishMetadataByTaskResultFile failed, cannot fetch sender from sign, userType: {%s}, user: {%s}",
			req.GetInformation().GetUserType().String(), req.GetInformation().GetUser())
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "cannot fetch sender from sign"}, nil
	}
	if from != req.GetInformation().GetUser() {
		log.WithError(err).Errorf("RPC-API:PublishMetadataByTaskResultFile failed, sender from sign and user is not sameone, userType: {%s}, user: {%s}, sender of sign: {%s}",
			req.GetInformation().GetUserType().String(), req.GetInformation().GetUser(), from)
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "the user sign is invalid"}, nil
	}
	if err := svr.B.SendMsg(metadataMsg); nil != err {
		log.WithError(err).Error("RPC-API:PublishMetadataByTaskResultFile failed")

		errMsg := fmt.Sprintf("%s, metadataId: {%s}", backend.ErrPublishMetadataMsg.Error(), metadataMsg.GetMetadataId())
		return &carrierapipb.PublishMetadataResponse{Status: backend.ErrPublishMetadataMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:PublishMetadataByTaskResultFile succeed, return metadataId: {%s}", metadataMsg.GetMetadataId())
	return &carrierapipb.PublishMetadataResponse{
		Status:     0,
		Msg:        backend.OK,
		MetadataId: metadataMsg.GetMetadataId(),
	}, nil
}

func (svr *Server) UpdateMetadata(ctx context.Context, req *carrierapipb.UpdateMetadataRequest) (*carriertypespb.SimpleResponse, error) {
	log.Debugf("RPC-API:UpdateMetadata req is:%s", req.String())
	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:UpdateMetadata failed, query local identity failed, can not publish metadata")
		return &carriertypespb.SimpleResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if nil == req.GetInformation() {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataInfomation"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetMetadataId(), "") {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataId"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetMetadataName(), "") {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataName"}, nil
	}
	if req.GetInformation().GetMetadataType() == commonconstantpb.MetadataType_MetadataType_Unknown {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown metadataType"}, nil
	}
	// DataHash
	if "" == strings.Trim(req.GetInformation().GetDesc(), "") {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require desc"}, nil
	}
	if req.GetInformation().GetLocationType() == commonconstantpb.DataLocationType_DataLocationType_Unknown {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown locationType"}, nil
	}
	if req.GetInformation().GetDataType() == commonconstantpb.OrigindataType_OrigindataType_Unknown {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown dataType"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetIndustry(), "") {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require industry"}, nil
	}
	if "" == strings.Trim(req.GetInformation().GetMetadataOption(), "") {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require metadataOption"}, nil
	}

	// get old metadata info by dataCenter
	var oldMetadataInfo *types.Metadata
	oldMetadataInfo, err = svr.B.GetMetadataDetail(req.GetInformation().GetMetadataId())
	if err != nil {
		return &carriertypespb.SimpleResponse{Status: backend.ErrQueryMetadataDetailById.ErrCode(), Msg: "call GetMetadataDetail fail"}, err
	}
	oldMetadata := oldMetadataInfo.GetData()
	if result, err := checkCanUpdateMetadataFieldIsLegal(oldMetadata, req); err != nil {
		return result, err
	}

	metadataUpdateMsg := &types.MetadataUpdateMsg{
		/*
			MetadataId     string
			MetadataName   string
			MetadataType   constant.MetadataType
			DataHash       string
			Desc           string
			LocationType   constant.DataLocationTy
			DataType       constant.OrigindataType
			Industry       string
			State          constant.MetadataState
			PublishAt      uint64
			UpdateAt       uint64
			Nonce          uint64
			MetadataOption string
			// add by v0.5.0
			User                 string
			UserType             constant.UserType
			Sign                 []byte
		*/
		MetadataSummary: &carriertypespb.MetadataSummary{
			MetadataId:     req.GetInformation().GetMetadataId(),
			MetadataName:   req.GetInformation().GetMetadataName(),
			MetadataType:   req.GetInformation().GetMetadataType(),
			DataHash:       req.GetInformation().GetDataHash(),
			Desc:           req.GetInformation().GetDesc(),
			LocationType:   req.GetInformation().GetLocationType(),
			DataType:       req.GetInformation().GetDataType(),
			Industry:       req.GetInformation().GetIndustry(),
			State:          req.GetInformation().GetState(),
			PublishAt:      req.GetInformation().GetPublishAt(),
			UpdateAt:       req.GetInformation().GetUpdateAt(),
			Nonce:          req.GetInformation().GetNonce(),
			MetadataOption: req.GetInformation().GetMetadataOption(),
			User:           req.GetInformation().GetUser(),
			UserType:       req.GetInformation().GetUserType(),
			Sign:           req.GetInformation().GetSign(),
		},
		CreateAt: timeutils.UnixMsecUint64(),
	}
	// check from
	from, _, err := signsuite.Sender(metadataUpdateMsg.GetMetadataSummary().GetUserType(), metadataUpdateMsg.Hash(), metadataUpdateMsg.GetMetadataSummary().GetSign())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:UpdateMetadata failed, cannot fetch sender from sign, userType: {%s}, user: {%s}",
			req.GetInformation().GetUserType().String(), req.GetInformation().GetUser())
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "cannot fetch sender from sign"}, nil
	}
	if from != req.GetInformation().GetUser() {
		log.WithError(err).Errorf("RPC-API:UpdateMetadata failed, sender from sign and user is not sameone, userType: {%s}, user: {%s}, sender of sign: {%s}",
			req.GetInformation().GetUserType().String(), req.GetInformation().GetUser(), from)
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "the user sign is invalid"}, nil
	}
	if err := svr.B.SendMsg(metadataUpdateMsg); nil != err {
		log.WithError(err).Error("RPC-API:UpdateMetadata failed")

		errMsg := fmt.Sprintf("%s, metadataId: {%s}", backend.ErrPublishMetadataMsg.Error(), req.GetInformation().GetMetadataId())
		return &carriertypespb.SimpleResponse{Status: backend.ErrPublishMetadataMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:UpdateMetadata succeed, metadataId: {%s}", req.Information.GetMetadataId())
	return &carriertypespb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func checkCanUpdateMetadataFieldIsLegal(oldMetadata *carriertypespb.MetadataPB, req *carrierapipb.UpdateMetadataRequest) (*carriertypespb.SimpleResponse, error) {
	responseMsg := ""
	if oldMetadata.MetadataType != req.GetInformation().MetadataType {
		responseMsg = "update MetadataType not equal to old MetadataType"
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: responseMsg}, errors.New(responseMsg)
	}
	if oldMetadata.DataType != req.GetInformation().DataType {
		responseMsg = "update DataType not equal to old DataType"
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: responseMsg}, errors.New(responseMsg)
	}
	if oldMetadata.DataHash != req.GetInformation().DataHash {
		responseMsg = "update DataHash not equal to old DataHash"
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: responseMsg}, errors.New(responseMsg)
	}
	if oldMetadata.State != req.GetInformation().State {
		responseMsg = "update State not equal to old State"
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: responseMsg}, errors.New(responseMsg)
	}
	if oldMetadata.LocationType != req.GetInformation().LocationType {
		responseMsg = "update LocationType not equal to old LocationType"
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: responseMsg}, errors.New(responseMsg)
	}
	if oldMetadata.User != req.GetInformation().User {
		responseMsg = "update User not equal to old User"
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: responseMsg}, errors.New(responseMsg)
	}
	if oldMetadata.UserType != req.GetInformation().UserType {
		responseMsg = "update UserType not equal to old UserType"
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: responseMsg}, errors.New(responseMsg)
	}

	metadataOption := req.GetInformation().GetMetadataOption()
	var (
		consumeOptions   []string
		consumeTypes     []uint8
		duplicateAddress []string
	)
	if req.GetInformation().GetDataType() == commonconstantpb.OrigindataType_OrigindataType_CSV {
		var option *types.MetadataOptionCSV
		err := json.Unmarshal([]byte(metadataOption), &option)
		if err != nil {
			responseMsg = "MetadataOption is not a valid json string"
			return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: responseMsg}, errors.New(responseMsg)
		}
		consumeOptions, consumeTypes = option.GetConsumeOptions(), option.GetConsumeTypes()
		// Return directly without consumption, no need to do follow-up checks
		if len(consumeTypes) == 0 {
			return nil, nil
		}
	} else {
		//todo other dataType
		return nil, nil
	}

	if len(consumeOptions) != len(consumeTypes) {
		responseMsg = "consumeOptions len not equal to consumeTypes len"
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: responseMsg}, errors.New(responseMsg)
	}

	for index, consumeOption := range consumeOptions {
		log.Debugf("checkCanUpdateMetadataFieldIsLegal consumeOption str is %s", consumeOption)
		if consumeTypes[index] == types.ConsumeMetadataAuth {
			continue
		} else if consumeTypes[index] == types.ConsumeTk20 {
			var info []*types.MetadataConsumeOptionTK20
			if err := json.Unmarshal([]byte(consumeOption), &info); err != nil {
				responseMsg = fmt.Sprintf("json Unmarshal consumeOption %s field", consumeOption)
				return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: responseMsg}, errors.New(responseMsg)
			} else {
				set := make(map[string]struct{}, 0)
				for _, consumeMap := range info {
					contractAddress := consumeMap.GetContract()
					if contractAddress == "" {
						responseMsg = "tk20 consumeOption no contract field"
						return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: responseMsg}, errors.New(responseMsg)
					}
					if _, ok := set[contractAddress]; !ok {
						set[contractAddress] = struct{}{}
					} else {
						duplicateAddress = append(duplicateAddress, contractAddress)
					}
				}
			}
		} else if consumeTypes[index] == types.ConsumeTk721 {
			var addressArr []types.MetadataConsumeOptionTK721
			if err := json.Unmarshal([]byte(consumeOption), &addressArr); err != nil {
				responseMsg = fmt.Sprintf("tk721 json Unmarshal consumeOption %s field", consumeOption)
				return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: responseMsg}, errors.New(responseMsg)
			} else {
				set := make(map[string]struct{}, 0)
				for _, address := range addressArr {
					contractAddress := address.GetContract()
					if _, ok := set[contractAddress]; !ok {
						set[contractAddress] = struct{}{}
					} else {
						duplicateAddress = append(duplicateAddress, contractAddress)
					}
				}
			}
		} else {
			responseMsg = fmt.Sprintf("consumeType %s unknown,please check!", string(consumeTypes[index]))
			return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: responseMsg}, errors.New(responseMsg)
		}
	}
	if len(duplicateAddress) != 0 {
		jsonStr, _ := json.Marshal(duplicateAddress)
		responseMsg = fmt.Sprintf("exits duplicate contract Address,it's %s", jsonStr)
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: responseMsg}, errors.New(responseMsg)
	}
	return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: responseMsg}, nil
}
