package auth

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

// for organization identity

func (svr *Server) ApplyIdentityJoin(ctx context.Context, req *pb.ApplyIdentityJoinRequest) (*libtypes.SimpleResponse, error) {

	identity, err := svr.B.GetNodeIdentity()
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("RPC-API:ApplyIdentityJoin failed, query local identity failed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}, imgUrl: {%s}, details: {%s}",
			req.GetInformation().GetIdentityId(), req.GetInformation().GetNodeId(), req.GetInformation().GetNodeName(), req.GetInformation().GetImageUrl(), req.GetInformation().GetDetails())

		errMsg := fmt.Sprintf("query local identity failed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}, imgUrl: {%s}, details: {%s}",
			req.GetInformation().GetIdentityId(), req.GetInformation().GetNodeId(), req.GetInformation().GetNodeName(), req.GetInformation().GetImageUrl(), req.GetInformation().GetDetails())
		return &libtypes.SimpleResponse{Status: backend.ErrApplyIdentityMsg.ErrCode(), Msg: errMsg}, nil
	}

	if nil != identity {
		log.Errorf("RPC-API:ApplyIdentityJoin failed, identity was already exist, old identityId: {%s}, old nodeId: {%s}, old nodeName: {%s}, old imgUrl: {%s}, old details: {%s}",
			identity.GetIdentityId(), identity.GetNodeId(), identity.GetName(), identity.GetImageUrl(), identity.GetDetails())

		errMsg := fmt.Sprintf("identity was already exist, identityId: {%s}, nodeId: {%s}, nodeName: {%s}, imgUrl: {%s}, details: {%s}",
			identity.GetIdentityId(), identity.GetNodeId(), identity.GetName(), identity.GetImageUrl(), identity.GetDetails())
		return &libtypes.SimpleResponse{Status: backend.ErrApplyIdentityMsg.ErrCode(), Msg: errMsg}, nil
	}

	if nil == req.GetInformation() {
		return &libtypes.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: backend.ErrRequireParams.Error()}, nil
	}

	if "" == strings.Trim(req.GetInformation().GetIdentityId(), "") ||
		"" == strings.Trim(req.GetInformation().GetNodeName(), "") {
		return &libtypes.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: backend.ErrRequireParams.Error()}, nil
	}

	identityMsg := types.NewIdentityMessageFromRequest(req)
	if err := identityMsg.CheckLength(); nil != err {
		errMsg := fmt.Sprintf("check fields len failed, %s , identityId: {%s}, nodeId: {%s}, nodeName: {%s}, imgUrl: {%s}, details: {%s}",
			err, identity.GetIdentityId(), identity.GetNodeId(), identity.GetName(), identity.GetImageUrl(), identity.GetDetails())
		return &libtypes.SimpleResponse{Status: backend.ErrApplyIdentityMsg.ErrCode(), Msg: errMsg}, nil
	}

	if err = svr.B.SendMsg(identityMsg); nil != err {
		log.WithError(err).Errorf("RPC-API:ApplyIdentityJoin failed, identityId: {%s}, nodeId: {%s}, nodeName: {%s}, imgUrl: {%s}, details: {%s}",
			req.GetInformation().GetIdentityId(), req.GetInformation().GetNodeId(), req.GetInformation().GetNodeName(), req.GetInformation().GetImageUrl(), req.GetInformation().GetDetails())

		errMsg := fmt.Sprintf("%s, identityId: {%s}, nodeId: {%s}, nodeName: {%s}", backend.ErrApplyIdentityMsg.Error(),
			req.GetInformation().GetIdentityId(), req.GetInformation().GetNodeId(), req.GetInformation().GetNodeName())
		return &libtypes.SimpleResponse{Status: backend.ErrApplyIdentityMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:ApplyIdentityJoin succeed SendMsg, identityId: {%s}, nodeId: {%s}, nodeName: {%s}, imgUrl: {%s}, details: {%s}",
		req.GetInformation().GetIdentityId(), req.GetInformation().GetNodeId(), req.GetInformation().GetNodeName(), req.GetInformation().GetImageUrl(), req.GetInformation().GetDetails())
	return &libtypes.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) RevokeIdentityJoin(ctx context.Context, req *emptypb.Empty) (*libtypes.SimpleResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if rawdb.IsDBNotFoundErr(err) {
		log.WithError(err).Errorf("RPC-API:RevokeIdentityJoin failed, the identity was not exist, can not revoke identity")
		errMsg := fmt.Sprintf("%s, the identity was not exist, can not revoke identity", backend.ErrRevokeIdentityMsg.Error())
		return &libtypes.SimpleResponse{Status: backend.ErrRevokeIdentityMsg.ErrCode(), Msg: errMsg}, nil
	}

	// what if local task we can not revoke identity
	has, err := svr.B.HasLocalTask()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokeIdentityJoin failed, can not check has local task")

		errMsg := fmt.Sprintf("%s, can not check has local task", backend.ErrRevokeIdentityMsg.Error())
		return &libtypes.SimpleResponse{Status: backend.ErrRevokeIdentityMsg.ErrCode(), Msg: errMsg}, nil
	}

	if has {
		log.WithError(err).Errorf("RPC-API:RevokeIdentityJoin failed, don't revoke identity when has local task")

		errMsg := fmt.Sprintf("%s, don't revoke identity when has local task", backend.ErrRevokeIdentityMsg.Error())
		return &libtypes.SimpleResponse{Status: backend.ErrRevokeIdentityMsg.ErrCode(), Msg: errMsg}, nil
	}

	identityRevokeMsg := types.NewIdentityRevokeMessage()

	if err = svr.B.SendMsg(identityRevokeMsg); nil != err {
		log.WithError(err).Error("RPC-API:RevokeIdentityJoin failed")

		errMsg := fmt.Sprintf("%s, send identity revoke msg failed", backend.ErrRevokeIdentityMsg.Error())
		return &libtypes.SimpleResponse{Status: backend.ErrRevokeIdentityMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debug("RPC-API:RevokeIdentityJoin succeed SendMsg")
	return &libtypes.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) GetNodeIdentity(ctx context.Context, req *emptypb.Empty) (*pb.GetNodeIdentityResponse, error) {
	identity, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetNodeIdentity failed")
		return &pb.GetNodeIdentityResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}
	return &pb.GetNodeIdentityResponse{
		Status: 0,
		Msg:    backend.OK,
		Information: &libtypes.Organization{
			NodeName:   identity.GetName(),
			NodeId:     identity.GetNodeId(),
			IdentityId: identity.GetIdentityId(),
			ImageUrl:   identity.GetImageUrl(),
			Details:    identity.GetDetails(),
			UpdateAt:   identity.GetUpdateAt(),
		},
	}, nil
}

func (svr *Server) GetIdentityList(ctx context.Context, req *pb.GetIdentityListRequest) (*pb.GetIdentityListResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = backend.DefaultPageSize
	}
	identityList, err := svr.B.GetIdentityList(req.GetLastUpdated(), pageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:QueryIdentityList failed")
		return &pb.GetIdentityListResponse{Status: backend.ErrQueryIdentityList.ErrCode(), Msg: backend.ErrQueryIdentityList.Error()}, nil
	}
	arr := make([]*libtypes.Organization, len(identityList))
	for i, identity := range identityList {
		iden := &libtypes.Organization{
			NodeName:   identity.GetName(),
			NodeId:     identity.GetNodeId(),
			IdentityId: identity.GetIdentityId(),
			ImageUrl:   identity.GetImageUrl(),
			Details:    identity.GetDetails(),
			UpdateAt:   identity.GetUpdateAt(),
			Status:     identity.GetDataStatus(),
		}
		arr[i] = iden
	}
	log.Debugf("Query all org's identity list, len: {%d}", len(identityList))
	return &pb.GetIdentityListResponse{
		Status:    0,
		Msg:       backend.OK,
		Identitys: arr,
	}, nil
}

// for metadata authority apply

func (svr *Server) ApplyMetadataAuthority(ctx context.Context, req *pb.ApplyMetadataAuthorityRequest) (*pb.ApplyMetadataAuthorityResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ApplyMetadataAuthority failed, query local identity failed")
		return &pb.ApplyMetadataAuthorityResponse{ Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if req.GetUser() == "" {
		return &pb.ApplyMetadataAuthorityResponse{ Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "the user is empty"}, nil
	}
	if !verifyUserType(req.GetUserType()) {
		return &pb.ApplyMetadataAuthorityResponse{ Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "the user type is wrong"}, nil
	}
	if req.GetAuth() == nil {
		return &pb.ApplyMetadataAuthorityResponse{ Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "the metadata auth is nil"}, nil
	}
	if len(req.GetSign()) == 0 {
		return &pb.ApplyMetadataAuthorityResponse{ Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "the user signature is nil"}, nil
	}
	signer := etypes.NewEIP155Signer(svr.B.GetCarrierChainConfig().BlockChainIdCache[req.GetUserType()])
	from, err := signer.Sender(tx)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ApplyMetadataAuthority failed, cannot fetch sender from sign, userType: {%s}, user: {%s}",
			req.GetUserType().String(), req.GetUser())
		return &pb.ApplyMetadataAuthorityResponse{ Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "cannot fetch sender from sign"}, nil
	}
	if from.Hex() != req.GetUser() {
		log.WithError(err).Errorf("RPC-API:ApplyMetadataAuthority failed, sender from sign and user is not sameone, userType: {%s}, user: {%s}, sender of sign: {%s}",
			req.GetUserType().String(), req.GetUser(), from.Hex())
		return &pb.ApplyMetadataAuthorityResponse{ Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "the user sign is invalid"}, nil
	}

	now := timeutils.UnixMsecUint64()
	switch req.GetAuth().GetUsageRule().GetUsageType() {
	case libtypes.MetadataUsageType_Usage_Period:
		if now >= req.GetAuth().GetUsageRule().GetEndAt() {
			log.Errorf("RPC-API:ApplyMetadataAuthority failed, usaageRule endTime of metadataAuth has expire, userType: {%s}, user: {%s}, metadataId: {%s}, usageType: {%s}, usageEndTime: {%d}, now: {%d}",
				req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId(), req.GetAuth().GetUsageRule().GetUsageType().String(), req.GetAuth().GetUsageRule().GetEndAt(), now)
			return &pb.ApplyMetadataAuthorityResponse{ Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "usaageRule endTime of metadataAuth has expire"}, nil
		}
	case libtypes.MetadataUsageType_Usage_Times:
		if req.GetAuth().GetUsageRule().GetTimes() == 0 {
			log.Errorf("RPC-API:ApplyMetadataAuthority failed, usaageRule times of metadataAuth must be greater than zero, userType: {%s}, user: {%s}, metadataId: {%s}, usageType: {%s}, usageEndTime: {%d}, now: {%d}",
				req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId(), req.GetAuth().GetUsageRule().GetUsageType().String(), req.GetAuth().GetUsageRule().GetEndAt(), now)
			return &pb.ApplyMetadataAuthorityResponse{ Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "usaageRule times of metadataAuth must be greater than zero"}, nil
		}
	default:
		log.Errorf("RPC-API:ApplyMetadataAuthority failed, unknown usageType of the metadataAuth, userType: {%s}, user: {%s}, metadataId: {%s}, usageType: {%s}",
			req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId(), req.GetAuth().GetUsageRule().GetUsageType().String())
		return &pb.ApplyMetadataAuthorityResponse{ Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "unknown usageType of the metadataAuth"}, nil
	}

	// ############################################
	// ############################################
	// NOTE:
	// 		check metadataId and identity of authority
	// ############################################
	// ############################################
	ideneityList, err := svr.B.GetIdentityList(timeutils.BeforeYearUnixMsecUint64(), backend.DefaultMaxPageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:ApplyMetadataAuthority failed, query global identity list failed")
		return &pb.ApplyMetadataAuthorityResponse{ Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "query global identity list failed"}, nil
	}
	var valid bool // false
	for _, identity := range ideneityList {
		if identity.GetIdentityId() == req.GetAuth().GetOwner().GetIdentityId() {
			valid = true
			break
		}
	}
	if !valid {
		log.WithError(err).Errorf("RPC-API:ApplyMetadataAuthority failed, not found identity with identityId of auth, identityId: {%s}",
			req.GetAuth().GetOwner().GetIdentityId())
		return &pb.ApplyMetadataAuthorityResponse{ Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "not found identity with identityId of auth"}, nil
	}
	// reset val
	valid = false

	metadataList, err := svr.B.GetGlobalMetadataDetailListByIdentityId(req.GetAuth().GetOwner().GetIdentityId(), timeutils.BeforeYearUnixMsecUint64(), backend.DefaultMaxPageSize)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ApplyMetadataAuthority failed, query global metadata list by identityId failed, identityId: {%s}", req.GetAuth().GetOwner().GetIdentityId())
		return &pb.ApplyMetadataAuthorityResponse{ Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "query global metadata list failed"}, nil
	}
	for _, metadata := range metadataList {
		if metadata.GetInformation().GetMetadataSummary().GetMetadataId() == req.GetAuth().GetMetadataId() {
			valid = true
			break
		}
	}
	if !valid {
		log.WithError(err).Errorf("RPC-API:ApplyMetadataAuthority failed, not found metadata with metadataId of auth, metadataId: {%s}",
			req.GetAuth().GetMetadataId())
		return &pb.ApplyMetadataAuthorityResponse{ Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: "not found metadata with metadataId of auth"}, nil
	}

	// check the metadataId whether has valid metadataAuth with current userType and user.
	has, err := svr.B.HasValidMetadataAuth(req.GetUserType(), req.GetUser(), req.GetAuth().GetOwner().GetIdentityId(), req.GetAuth().GetMetadataId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:ApplyMetadataAuthority failed, query valid user metadataAuth failed, userType: {%s}, user: {%s}, metadataId: {%s}",
			req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId())

		errMsg := fmt.Sprintf(backend.ErrApplyMetadataAuthority.Error(), "query valid user metadataAuth failed",
			req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId())
		return &pb.ApplyMetadataAuthorityResponse{ Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: errMsg}, nil
	}

	if has {
		log.Errorf("RPC-API:ApplyMetadataAuthority failed, has valid metadataAuth exists, userType: {%s}, user: {%s}, metadataId: {%s}",
			req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId())

		errMsg := fmt.Sprintf("%s, userType: {%s}, user: {%s}, metadataId: {%s}", backend.ErrApplyMetadataAuthority.Error(),
			req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId())
		return &pb.ApplyMetadataAuthorityResponse{ Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: errMsg}, nil
	}

	metadataAuthorityMsg := types.NewMetadataAuthorityMessageFromRequest(req)
	metadataAuthId := metadataAuthorityMsg.GetMetadataAuthId()

	if err = svr.B.SendMsg(metadataAuthorityMsg); nil != err {
		log.WithError(err).Error("RPC-API:ApplyMetadataAuthority failed")

		errMsg := fmt.Sprintf(backend.ErrApplyMetadataAuthority.Error(), "send metadata authority msg failed",
			req.GetUserType().String(), req.GetUser(), req.GetAuth().GetMetadataId())
		return &pb.ApplyMetadataAuthorityResponse{ Status: backend.ErrApplyMetadataAuthority.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:ApplyMetadataAuthority succeed, userType: {%s}, user: {%s}, metadataOwner: {%s}, metadataId: {%s}, usageRule: {%s},  return metadataAuthId: {%s}",
		req.GetUserType().String(), req.GetUser(), req.GetAuth().GetOwner().String(), req.GetAuth().GetMetadataId(), req.GetAuth().GetUsageRule().String(), metadataAuthId)
	return &pb.ApplyMetadataAuthorityResponse{
		Status:         0,
		Msg:            backend.OK,
		MetadataAuthId: metadataAuthId,
	}, nil
}

func (svr *Server) RevokeMetadataAuthority(ctx context.Context, req *pb.RevokeMetadataAuthorityRequest) (*libtypes.SimpleResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, query local identity failed")
		return &libtypes.SimpleResponse{ Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if req.GetUser() == "" {
		return &libtypes.SimpleResponse{ Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "the user is empty"}, nil
	}
	if !verifyUserType(req.GetUserType()) {
		return &libtypes.SimpleResponse{ Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "the user type is wrong"}, nil
	}
	if req.GetMetadataAuthId() == "" {
		return &libtypes.SimpleResponse{ Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "the metadata auth is nil"}, nil
	}
	if len(req.GetSign()) == 0 {
		return &libtypes.SimpleResponse{ Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "the user signature is nil"}, nil
	}

	authorityList, err := svr.B.GetLocalMetadataAuthorityList(timeutils.BeforeYearUnixMsecUint64(), backend.DefaultMaxPageSize)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, query local metadataAuth list failed, %s", err)
		return &libtypes.SimpleResponse{ Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "query local metadataAuth list failed"}, nil
	}

	for _, auth := range authorityList {

		if auth.GetUserType() == req.GetUserType() &&
			auth.GetUser() == req.GetUser() &&
			auth.GetData().GetMetadataAuthId() == req.GetMetadataAuthId() {

			// The data authorization application information that has been audited and cannot be revoked
			if auth.GetData().GetAuditOption() != libtypes.AuditMetadataOption_Audit_Pending {
				log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, the metadataAuth was audited")
				return &libtypes.SimpleResponse{ Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "the metadataAuth state was audited"}, nil
			}

			if auth.GetData().GetState() != libtypes.MetadataAuthorityState_MAState_Released {
				log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, the metadataAuth state was not released")
				return &libtypes.SimpleResponse{ Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "the metadataAuth state was not released"}, nil
			}

			switch auth.GetData().GetAuth().GetUsageRule().GetUsageType() {
			case libtypes.MetadataUsageType_Usage_Period:
				if timeutils.UnixMsecUint64() >= auth.GetData().GetAuth().GetUsageRule().GetEndAt() {
					log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, the metadataAuth had been expire")
					return &libtypes.SimpleResponse{ Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "the metadataAuth had been expire"}, nil
				}
			case libtypes.MetadataUsageType_Usage_Times:
				if auth.GetData().GetUsedQuo().GetUsedTimes() >= auth.GetData().GetAuth().GetUsageRule().GetTimes() {
					log.WithError(err).Errorf("RPC-API:RevokeMetadataAuthority failed, the metadataAuth had been not enough times")
					return &libtypes.SimpleResponse{ Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "the metadataAuth had been not enough times"}, nil
				}
			default:
				log.Errorf("unknown usageType of the old metadataAuth on AuthorityManager.filterMetadataAuth(), metadataAuthId: {%s}", auth.GetData().GetMetadataAuthId())
				return &libtypes.SimpleResponse{ Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: "unknown usageType of the old metadataAuth"}, nil
			}
		}
	}

	metadataAuthorityRevokeMsg := types.NewMetadataAuthorityRevokeMessageFromRequest(req)
	metadataAuthId := metadataAuthorityRevokeMsg.GetMetadataAuthId()

	if err := svr.B.SendMsg(metadataAuthorityRevokeMsg); nil != err {
		log.WithError(err).Error("RPC-API:RevokeMetadataAuthority failed")

		errMsg := fmt.Sprintf("%s, userType: {%s}, user: {%s}, MetadataAuthId: {%s}", backend.ErrRevokeMetadataAuthority.Error(),
			req.GetUserType().String(), req.GetUser(), req.GetMetadataAuthId())
		return &libtypes.SimpleResponse{ Status: backend.ErrRevokeMetadataAuthority.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:RevokeMetadataAuthority succeed, userType: {%s}, user: {%s}, metadataAuthId: {%s}",
		req.GetUserType().String(), req.GetUser(), metadataAuthId)
	return &libtypes.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) AuditMetadataAuthority(ctx context.Context, req *pb.AuditMetadataAuthorityRequest) (*pb.AuditMetadataAuthorityResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:AuditMetadataAuthority failed, query local identity failed")
		return &pb.AuditMetadataAuthorityResponse{ Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if "" == req.GetMetadataAuthId() {
		return &pb.AuditMetadataAuthorityResponse{ Status: backend.ErrAuditMetadataAuth.ErrCode(), Msg: "the metadataAuth Id is empty"}, nil
	}

	if req.GetAudit() == libtypes.AuditMetadataOption_Audit_Pending {
		return &pb.AuditMetadataAuthorityResponse{ Status: backend.ErrAuditMetadataAuth.ErrCode(), Msg: "the valid audit metadata option must cannot pending"}, nil
	}

	option, err := svr.B.AuditMetadataAuthority(types.NewMetadataAuthAudit(req.GetMetadataAuthId(), req.GetSuggestion(), req.GetAudit()))
	if nil != err {
		log.WithError(err).Error("RPC-API:AuditMetadataAuthority failed")

		errMsg := fmt.Sprintf("%s, metadataAuthId: {%s}, audit option: {%s}, audit suggestion: {%s}", backend.ErrAuditMetadataAuth.Error(),
			req.GetMetadataAuthId(), req.GetAudit().String(), req.GetSuggestion())
		return &pb.AuditMetadataAuthorityResponse{ Status: backend.ErrAuditMetadataAuth.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:AuditMetadataAuthority succeed, metadataAuthId: {%s}, audit option: {%s}, audit suggestion: {%s}",
		req.GetMetadataAuthId(), req.GetAudit().String(), req.GetSuggestion())

	return &pb.AuditMetadataAuthorityResponse{
		Status: 0,
		Msg:    backend.OK,
		Audit:  option,
	}, nil
}

func (svr *Server) GetLocalMetadataAuthorityList(ctx context.Context, req *pb.GetMetadataAuthorityListRequest) (*pb.GetMetadataAuthorityListResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = backend.DefaultPageSize
	}
	authorityList, err := svr.B.GetLocalMetadataAuthorityList(req.GetLastUpdated(), pageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetLocalMetadataAuthorityList failed")
		return &pb.GetMetadataAuthorityListResponse{ Status: backend.ErrQueryAuthorityList.ErrCode(), Msg: backend.ErrQueryAuthorityList.Error()}, nil
	}
	arr := make([]*libtypes.MetadataAuthorityDetail, len(authorityList))
	for i, auth := range authorityList {
		arr[i] = &libtypes.MetadataAuthorityDetail{
			MetadataAuthId:  auth.GetData().GetMetadataAuthId(),
			User:            auth.GetData().GetUser(),
			UserType:        auth.GetData().GetUserType(),
			Auth:            auth.GetData().GetAuth(),
			AuditOption:     auth.GetData().GetAuditOption(),
			AuditSuggestion: auth.GetData().GetAuditSuggestion(),
			UsedQuo:         auth.GetData().GetUsedQuo(),
			ApplyAt:         auth.GetData().GetApplyAt(),
			AuditAt:         auth.GetData().GetAuditAt(),
			State:           auth.GetData().GetState(),
			UpdateAt:        auth.GetData().GetUpdateAt(),
		}
	}
	log.Debugf("RPC-API:GetLocalMetadataAuthorityList succeed, metadata authority list, len: {%d}", len(authorityList))
	return &pb.GetMetadataAuthorityListResponse{
		Status: 0,
		Msg:    backend.OK,
		MetadataAuths:   arr,
	}, nil
}

func (svr *Server) GetGlobalMetadataAuthorityList(ctx context.Context, req *pb.GetMetadataAuthorityListRequest) (*pb.GetMetadataAuthorityListResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = backend.DefaultPageSize
	}
	authorityList, err := svr.B.GetGlobalMetadataAuthorityList(req.GetLastUpdated(), pageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetGlobalMetadataAuthorityList failed")
		return &pb.GetMetadataAuthorityListResponse{ Status: backend.ErrQueryAuthorityList.ErrCode(), Msg: backend.ErrQueryAuthorityList.Error()}, nil
	}
	arr := make([]*libtypes.MetadataAuthorityDetail, len(authorityList))
	for i, auth := range authorityList {
		arr[i] = &libtypes.MetadataAuthorityDetail{
			MetadataAuthId:  auth.GetData().GetMetadataAuthId(),
			User:            auth.GetData().GetUser(),
			UserType:        auth.GetData().GetUserType(),
			Auth:            auth.GetData().GetAuth(),
			AuditOption:     auth.GetData().GetAuditOption(),
			AuditSuggestion: auth.GetData().GetAuditSuggestion(),
			UsedQuo:         auth.GetData().GetUsedQuo(),
			ApplyAt:         auth.GetData().GetApplyAt(),
			AuditAt:         auth.GetData().GetAuditAt(),
			State:           auth.GetData().GetState(),
			UpdateAt:        auth.GetData().GetUpdateAt(),
		}
	}
	log.Debugf("RPC-API:GetGlobalMetadataAuthorityList succeed, metadata authority list, len: {%d}", len(authorityList))
	return &pb.GetMetadataAuthorityListResponse{
		Status: 0,
		Msg:    backend.OK,
		MetadataAuths:   arr,
	}, nil
}

func verifyUserType(userType libtypes.UserType) bool {
	switch userType {
	case libtypes.UserType_User_1:
		return true
	case libtypes.UserType_User_2:
		return true
	case libtypes.UserType_User_3:
		return true
	default:
		return false
	}
}
