package metadata

import (
	"encoding/json"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/carrierdb"
	"github.com/datumtechs/datum-network-carrier/carrierdb/rawdb"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	"github.com/datumtechs/datum-network-carrier/core/policy"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/rpc/backend"
	"github.com/datumtechs/datum-network-carrier/types"
)

var (
	ErrMetadataAuthHasAudited        = fmt.Errorf("the metadataAuth state was audited")
	ErrMetadataAuthHasNotReleased    = fmt.Errorf("the metadataAuth state was not released")
	ErrMetadataAuthHasExpired        = fmt.Errorf("the metadataAuth had been expired")
	ErrMetadataAuthHasNotEnoughTimes = fmt.Errorf("the metadataAuth had been not enough times")
	ErrMetadataAuthUnknoenUsageType  = fmt.Errorf("unknown usageType of the old metadataAuth")
)

type MetadataAuthority struct {
	dataCenter   carrierdb.CarrierDB
	policyEngine *policy.PolicyEngine
}

func NewMetadataAuthority(dataCenter carrierdb.CarrierDB, policyEngine *policy.PolicyEngine) *MetadataAuthority {
	return &MetadataAuthority{
		dataCenter:   dataCenter,
		policyEngine: policyEngine,
	}
}

func (ma *MetadataAuthority) ApplyMetadataAuthority(metadataAuth *types.MetadataAuthority) error {
	return ma.dataCenter.InsertMetadataAuthority(metadataAuth)
}

func (ma *MetadataAuthority) AuditMetadataAuthority(audit *types.MetadataAuthAudit) (commonconstantpb.AuditMetadataOption, error) {

	log.Debugf("Start audit metadataAuth on MetadataAuthority.AuditMetadataAuthority(), the audit: %s", audit.String())

	// verify
	//
	// query metadataAuth with metadataAuthId from dataCenter
	metadataAuth, err := ma.GetMetadataAuthority(audit.GetMetadataAuthId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query old metadataAuth on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}",
			audit.GetMetadataAuthId())
		return commonconstantpb.AuditMetadataOption_Audit_Pending, fmt.Errorf("query metadataAuth failed, %s", err)
	}

	pass, err := ma.VerifyMetadataAuthInfo(metadataAuth)
	if nil != err {
		log.WithError(err).Errorf("Failed to verify old metadataAuth on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}",
			audit.GetMetadataAuthId())
		return commonconstantpb.AuditMetadataOption_Audit_Pending, fmt.Errorf("verify metadataAuth failed, %s", err)
	}
	if !pass {
		log.Errorf("Invalid metadataAuth on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}",
			audit.GetMetadataAuthId())
		return commonconstantpb.AuditMetadataOption_Audit_Pending, fmt.Errorf("invalid metadataAuth")
	}

	identity, err := ma.dataCenter.QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to query local identity on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}",
			audit.GetMetadataAuthId())
		return commonconstantpb.AuditMetadataOption_Audit_Pending, fmt.Errorf("query local identity failed, %s", err)
	}

	if identity.GetIdentityId() != metadataAuth.GetData().GetAuth().GetOwner().GetIdentityId() {
		log.Errorf("Failed to verify local identity and identity of metadataAuth owner, they is not same identity on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}",
			audit.GetMetadataAuthId())
		return commonconstantpb.AuditMetadataOption_Audit_Pending, fmt.Errorf("metadataAuth did not current own, %s", err)
	}

	// Check whether metadataAuth has been audited
	has, err := ma.dataCenter.HasUserMetadataAuthIdByMetadataId(metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetData().GetMetadataAuthId())
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed to check whether metadataAuth has been audited on MetadataAuthority.AuditMetadataAuthority(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}",
			metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetData().GetMetadataAuthId())
		return commonconstantpb.AuditMetadataOption_Audit_Pending, fmt.Errorf("check whether metadataAuth has been audited failed, %s", err)
	}
	if has {
		log.Errorf("The metadataAuth has been audited on MetadataAuthority.AuditMetadataAuthority(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}",
			metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetData().GetMetadataAuthId())
		return commonconstantpb.AuditMetadataOption_Audit_Refused, fmt.Errorf("the metadataAuth has been audited, %s", err)
	}

	auditOption := audit.GetAuditOption()
	auditSuggestion := audit.GetAuditSuggestion()

	_, err = ma.VerifyMetadataAuthInfo(metadataAuth)
	switch err {
	case ErrMetadataAuthHasAudited:
		log.Errorf("the old metadataAuth has already audited on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}, audit option: {%s}",
			audit.GetMetadataAuthId(), audit.GetAuditOption().String())
		return metadataAuth.GetData().GetAuditOption(), fmt.Errorf("the old metadataAuth has already audited, %s", audit.GetAuditOption().String())
	case ErrMetadataAuthHasNotReleased:
		log.Errorf("the old metadataAuth state is not release on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}",
			audit.GetMetadataAuthId())
		return metadataAuth.GetData().GetAuditOption(), fmt.Errorf("the old metadataAuth state is %s", metadataAuth.GetData().GetState().String())
	case ErrMetadataAuthHasExpired:
		metadataAuth.GetData().GetUsedQuo().Expire = true // update state, maybe state has invalid.
		metadataAuth.GetData().State = commonconstantpb.MetadataAuthorityState_MAState_Invalid
		//
		auditOption = commonconstantpb.AuditMetadataOption_Audit_Refused
		auditSuggestion = "metadataAuth has expired, refused it"
	case ErrMetadataAuthHasNotEnoughTimes:
		metadataAuth.GetData().State = commonconstantpb.MetadataAuthorityState_MAState_Invalid
		//
		auditOption = commonconstantpb.AuditMetadataOption_Audit_Refused
		auditSuggestion = "metadataAuth has no enough remain times, refused it"
	case ErrMetadataAuthUnknoenUsageType:
		log.Errorf("unknown usageType of the metadataAuth MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}",
			audit.GetMetadataAuthId())
		return metadataAuth.GetData().GetAuditOption(), fmt.Errorf("unknown usageType of the metadataAuth")
	}

	// update audit things.
	metadataAuth.GetData().AuditOption = auditOption
	metadataAuth.GetData().AuditSuggestion = auditSuggestion
	metadataAuth.GetData().AuditAt = timeutils.UnixMsecUint64()

	if err := ma.dataCenter.UpdateMetadataAuthority(metadataAuth); nil != err {
		log.WithError(err).Errorf("Failed to update metadataAuth after audit on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}, audit option:{%s}",
			audit.GetMetadataAuthId(), audit.GetAuditOption().String())
		return metadataAuth.GetData().GetAuditOption(), fmt.Errorf("update metadataAuth failed, %s", err)
	}

	log.Debugf("metadataAuth audit succeed and call succeed `UpdateMetadataAuthority()` on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s} with audit option: {%s}, suggestion: {%s}, auditAt: {%d}",
		metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuditOption(), metadataAuth.GetData().GetAuditSuggestion(), metadataAuth.GetData().GetAuditAt())

	if err := ma.dataCenter.StoreValidUserMetadataAuthStatusByMetadataId(metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetMergeStatus()); nil != err {
		log.WithError(err).Errorf("Failed to update metadataAuth status AND audit option into local db while metadataAuth has invalid on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}, status: {%s}, auditOption: {%s}",
			metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetState().String(), metadataAuth.GetData().GetAuditOption().String())
		return metadataAuth.GetData().GetAuditOption(), fmt.Errorf("remove metadataId and invalid metadataAuthId mapping failed, %s", err)
	}

	return metadataAuth.GetData().GetAuditOption(), nil
}

func (ma *MetadataAuthority) ConsumeMetadataAuthority(metadataAuthId string) error {

	log.Debugf("Start consume metadataAuth on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}", metadataAuthId)

	// checking start.
	metadataAuth, err := ma.GetMetadataAuthority(metadataAuthId)
	if nil != err {
		log.WithError(err).Errorf("Failed to query old metadataAuth on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}",
			metadataAuthId)
		return fmt.Errorf("query metadataAuth failed, %s", err)
	}

	// check audit option first
	if metadataAuth.GetData().GetAuditOption() != commonconstantpb.AuditMetadataOption_Audit_Passed {
		log.Errorf("the old metadataAuth audit option is not passed on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}, audit option: {%s}",
			metadataAuthId, metadataAuth.GetData().GetAuditOption().String())
		return fmt.Errorf("the old metadataAuth audit option is not passed")
	}

	// check state second
	if metadataAuth.GetData().GetState() != commonconstantpb.MetadataAuthorityState_MAState_Released {
		log.Errorf("the old metadataAuth state is not released on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}, state: {%s}",
			metadataAuthId, metadataAuth.GetData().GetState().String())
		return fmt.Errorf("the old metadataAuth state is not released")
	}

	usageRule := metadataAuth.GetData().GetAuth().GetUsageRule()
	usedQuo := metadataAuth.GetData().GetUsedQuo()

	var needUpdate bool

	// check anything else next...
	switch usageRule.GetUsageType() {
	case commonconstantpb.MetadataUsageType_Usage_Period:
		if timeutils.UnixMsecUint64() >= usageRule.GetEndAt() {
			usedQuo.Expire = true
			metadataAuth.GetData().State = commonconstantpb.MetadataAuthorityState_MAState_Invalid
			needUpdate = true
		}
	case commonconstantpb.MetadataUsageType_Usage_Times:

		usedQuo.UsedTimes += 1
		if usedQuo.GetUsedTimes() >= usageRule.GetTimes() {
			metadataAuth.GetData().State = commonconstantpb.MetadataAuthorityState_MAState_Invalid
		}
		needUpdate = true
	default:
		log.Errorf("unknown usageType of the old metadataAuth on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}",
			metadataAuthId)
		return fmt.Errorf("unknown usageType of the old metadataAuth")
	}

	if needUpdate {
		metadataAuth.GetData().UsedQuo = usedQuo
		if err := ma.dataCenter.UpdateMetadataAuthority(metadataAuth); nil != err {
			log.WithError(err).Errorf("Failed to update metadataAuth after consume on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}",
				metadataAuthId)
			return fmt.Errorf("update metadataAuth failed, %s", err)
		}

		log.Debugf("metadataAuth consume succeed and call succeed `UpdateMetadataAuthority()` on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s} with usageRule: %s usedQuo: %s, state: {%s}",
			metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetUsageRule().String(), metadataAuth.GetData().GetUsedQuo().String(), metadataAuth.GetData().GetState())
	}

	// remove
	if metadataAuth.GetData().GetState() == commonconstantpb.MetadataAuthorityState_MAState_Invalid {

		log.Debugf("the metadataAuth was invalid need to call `RemoveUserMetadataAuthIdByMetadataId()` on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}",
			metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetUserType(), metadataAuth.GetUser())

		if err := ma.dataCenter.StoreValidUserMetadataAuthStatusByMetadataId(metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetMergeStatus()); nil != err {
			log.WithError(err).Errorf("Failed to update metadataAuth status AND audit option into local db while metadataAuth has invalid on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}, status: {%s}, auditOption: {%s}",
				metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetState().String(), metadataAuth.GetData().GetAuditOption().String())
		}
	}

	return nil
}

func (ma *MetadataAuthority) GetMetadataAuthority(metadataAuthId string) (*types.MetadataAuthority, error) {
	return ma.dataCenter.QueryMetadataAuthority(metadataAuthId)
}

func (ma *MetadataAuthority) GetLocalMetadataAuthorityList(lastUpdate, pageSize uint64) (types.MetadataAuthArray, error) {
	identityId, err := ma.dataCenter.QueryIdentityId()
	if nil != err {
		return nil, err
	}
	return ma.dataCenter.QueryMetadataAuthorityListByIdentityId(identityId, lastUpdate, pageSize)
}

func (ma *MetadataAuthority) GetGlobalMetadataAuthorityList(lastUpdate uint64, pageSize uint64) (types.MetadataAuthArray, error) {
	return ma.dataCenter.QueryMetadataAuthorityList(lastUpdate, pageSize)
}

func (ma *MetadataAuthority) GetMetadataAuthorityListByIds(metadataAuthIds []string) (types.MetadataAuthArray, error) {
	return ma.dataCenter.QueryMetadataAuthorityListByIds(metadataAuthIds)
}

// check metadataAuth whether exists from datacenter ...
func (ma *MetadataAuthority) HasValidMetadataAuth(userType commonconstantpb.UserType, user, identityId, metadataId string) (bool, error) {

	// query metadataAuthList with target identityId (metadataId of target org)
	metadataAuthList, err := ma.dataCenter.QueryMetadataAuthorityListByIdentityId(identityId, timeutils.BeforeYearUnixMsecUint64(), backend.DefaultMaxPageSize)
	if nil != err {
		log.WithError(err).Errorf("Failed to QueryMetadataAuthorityListByIdentityId() on MetadataAuthority.HasValidMetadataAuth(), userType: {%s}, user:{%s}, identityId: {%s}, metadataId: {%s}",
			userType.String(), user, identityId, metadataId)
		return false, fmt.Errorf("query all valid metadataAuth list failed, %s", err)
	}

	if len(metadataAuthList) == 0 {
		log.Warnf("Not found metadataAuth list on MetadataAuthority.HasValidMetadataAuth()")
		return false, nil
	}

	var find bool

	// find valid metadataAauth only, and filter one by one.
	for _, auth := range metadataAuthList {
		if auth.GetUserType() == userType &&
			auth.GetUser() == user &&
			auth.GetData().GetAuth().GetMetadataId() == metadataId &&
			auth.GetData().GetState() == commonconstantpb.MetadataAuthorityState_MAState_Released {

			// filter one by one

			if auth.GetData().GetAuditOption() == commonconstantpb.AuditMetadataOption_Audit_Refused {
				continue
			}

			usageRule := auth.GetData().GetAuth().GetUsageRule()
			usedQuo := auth.GetData().GetUsedQuo()

			switch usageRule.GetUsageType() {
			case commonconstantpb.MetadataUsageType_Usage_Period:
				if timeutils.UnixMsecUint64() >= usageRule.GetEndAt() {
					log.Warnf("the metadataAuth was expired on MetadataAuthority.HasValidMetadataAuth(), userType: {%s}, user:{%s}, identityId: {%s}, metadataId: {%s}",
						userType.String(), user, identityId, metadataId)
					continue
				}
			case commonconstantpb.MetadataUsageType_Usage_Times:
				if usedQuo.GetUsedTimes() >= usageRule.GetTimes() {
					log.Warnf("the metadataAuth was used times exceed the limit on MetadataAuthority.HasValidMetadataAuth(), userType: {%s}, user:{%s}, identityId: {%s}, metadataId: {%s}",
						userType.String(), user, identityId, metadataId)
					continue
				}
			default:
				log.Errorf("unknown usageType of the old metadataAuth on MetadataAuthority.HasValidMetadataAuth(), userType: {%s}, user:{%s}, metadataId: {%s}, metadataAuthId: {%s}",
					userType.String(), user, metadataId, auth.GetData().GetMetadataAuthId())
				continue
			}

			// then find it
			find = true
			break
		}
	}
	return find, nil
}

func (ma *MetadataAuthority) HasNotValidMetadataAuth(userType commonconstantpb.UserType, user, ideneityId, metadataId string) (bool, error) {
	has, err := ma.HasValidMetadataAuth(userType, user, ideneityId, metadataId)
	if nil != err {
		return false, err
	}
	if has {
		return false, nil
	}
	return true, nil
}

func (ma *MetadataAuthority) VerifyMetadataAuthWithMetadataOption(auth *types.MetadataAuthority) (bool, error) {

	identity, err := ma.dataCenter.QueryIdentityById(auth.GetData().GetAuth().GetOwner().GetIdentityId())
	if nil != err {
		return false, fmt.Errorf("can not query global identity list, %s", err)
	}
	if nil == identity {
		return false, fmt.Errorf("not found identity with identityId of auth, %s, identityId: {%s}",
			err, auth.GetData().GetAuth().GetOwner().GetIdentityId())
	}

	metadata, err := ma.dataCenter.QueryMetadataById(auth.GetData().GetAuth().GetMetadataId())
	if rawdb.IsNoDBNotFoundErr(err) {
		return false, fmt.Errorf("found global metadata arr by metadataId failed, %s, metadataId: {%s}", err, auth.GetData().GetAuth())
	}
	if nil == metadata {
		return false, fmt.Errorf("not found metadata with metadataId from metadataAuth, %s, metadataId: {%s}, metadataAuthId: {%s}",
			err, auth.GetData().GetAuth().GetMetadataId(), auth.GetData().GetMetadataAuthId())
	}
	if metadata.GetData().GetOwner().GetIdentityId() != auth.GetData().GetAuth().GetOwner().GetIdentityId() {
		return false, fmt.Errorf("the owner identityId of metadataAuth and the owner identityId metadata is not same, %s, metadataId: {%s}, metadataAuthId: {%s}, identityId of metadata: {%s}, identityId of metadataAuth: {%s}",
			err, auth.GetData().GetAuth().GetMetadataId(), auth.GetData().GetMetadataAuthId(), metadata.GetData().GetOwner().GetIdentityId(), auth.GetData().GetAuth().GetOwner().GetIdentityId())
	}

	consumeTypes, err := ma.policyEngine.FetchConsumeTypes(metadata.GetData().GetDataType(), metadata.GetData().GetMetadataOption())
	if nil != err {
		return false, fmt.Errorf("not fetch consumeTypes from metadataOption, %s, metadataId: {%s}", err, auth.GetData().GetAuth().GetMetadataId())
	}

	var index = -1
	for i, typ := range consumeTypes {
		if typ == types.ConsumeMetadataAuth {
			index = i
			break
		}
	}
	if -1 == index {
		return false, fmt.Errorf("not found metadataAuthConsumeType from metadataOption, %s, metadataId: {%s}, consumeTypes: %v",
			err, auth.GetData().GetAuth().GetMetadataId(), consumeTypes)
	}
	consumeOptions, err := ma.policyEngine.FetchConsumeOptions(metadata.GetData().GetDataType(), metadata.GetData().GetMetadataOption())
	if nil != err {
		return false, fmt.Errorf("not fetch consumeOptions from metadataOption, %s, metadataId: {%s}",
			err, auth.GetData().GetAuth().GetMetadataId())
	}
	if len(consumeTypes) != len(consumeOptions) {
		return false, fmt.Errorf("consumeTypesLen and consumeOptionsLen is not same fron metadataOption, %s, metadataId: {%s}",
			err, auth.GetData().GetAuth().GetMetadataId())
	}
	var option *types.MetadataConsumeOptionMetadataAuth
	if err := json.Unmarshal([]byte(consumeOptions[index]), &option); nil != err {
		return false, fmt.Errorf("can not unmashal consumeOptions to metadataConsumeOptionMetadataAuth, %s, metadataId: {%s}",
			err, auth.GetData().GetAuth().GetMetadataId())
	}

	now := timeutils.UnixMsecUint64()
	switch auth.GetData().GetAuth().GetUsageRule().GetUsageType() {
	case commonconstantpb.MetadataUsageType_Usage_Period:
		// 1、check type
		if option.GetStatus()&types.McomaStatusPeriodConsumeKind != types.McomaStatusPeriodConsumeKind {
			return false, fmt.Errorf("metadataAuth consume kind not support period kind, metadataId: {%s}, usageType: {%s}, usageEndTime: {%d}, consume kind status: {%d}",
				auth.GetData().GetAuth().GetMetadataId(), auth.GetData().GetAuth().GetUsageRule().GetUsageType().String(), auth.GetData().GetAuth().GetUsageRule().GetEndAt(), option.GetStatus())
		}

		// 2、check time
		if now >= auth.GetData().GetAuth().GetUsageRule().GetEndAt() {
			return false, fmt.Errorf("usaageRule endTime of metadataAuth has expire, metadataId: {%s}, usageType: {%s}, usageEndTime: {%d}, now: {%d}",
				auth.GetData().GetAuth().GetMetadataId(), auth.GetData().GetAuth().GetUsageRule().GetUsageType().String(), auth.GetData().GetAuth().GetUsageRule().GetEndAt(), now)
		}

	case commonconstantpb.MetadataUsageType_Usage_Times:

		// 1、check type
		if option.GetStatus()&types.McomaStatusTimesConsumeKind != types.McomaStatusTimesConsumeKind {
			return false, fmt.Errorf("metadataAuth consume kind not support times kind, metadataId: {%s}, usageType: {%s}, usageEndTime: {%d}, consume kind status: {%d}",
				auth.GetData().GetAuth().GetMetadataId(), auth.GetData().GetAuth().GetUsageRule().GetUsageType().String(), auth.GetData().GetAuth().GetUsageRule().GetEndAt(), option.GetStatus())
		}

		// 2、check count
		if auth.GetData().GetAuth().GetUsageRule().GetTimes() == 0 {
			return false, fmt.Errorf("usaageRule times of metadataAuth must be greater than zero, metadataId: {%s}, usageType: {%s}, usageEndTime: {%d}, now: {%d}",
				auth.GetData().GetAuth().GetMetadataId(), auth.GetData().GetAuth().GetUsageRule().GetUsageType().String(), auth.GetData().GetAuth().GetUsageRule().GetEndAt(), now)
		}
	default:
		return false, fmt.Errorf("unknown usageType of the metadataAuth, metadataId: {%s}, usageType: {%s}",
			auth.GetData().GetAuth().GetMetadataId(), auth.GetData().GetAuth().GetUsageRule().GetUsageType().String())
	}

	checkUsageTypeFn := func(metadataConsumeStatus uint64, authConsumeStatus commonconstantpb.MetadataUsageType) (bool, error) {
		switch {

		// there are options of consumption by times and consumption by time period.
		case metadataConsumeStatus&types.McomaStatusTimesConsumeKind == types.McomaStatusTimesConsumeKind &&
			metadataConsumeStatus&types.McomaStatusPeriodConsumeKind == types.McomaStatusPeriodConsumeKind:

			if authConsumeStatus != commonconstantpb.MetadataUsageType_Usage_Times &&
				authConsumeStatus != commonconstantpb.MetadataUsageType_Usage_Period {
				return false, fmt.Errorf("unknown usageType of metadataAuth, metadataAuthId: {%s}, usageType: {%s}",
					auth.GetData().GetMetadataAuthId(), authConsumeStatus.String())
			}

		// there is option of consumption by times only.
		case metadataConsumeStatus&types.McomaStatusTimesConsumeKind == types.McomaStatusTimesConsumeKind &&
			metadataConsumeStatus&types.McomaStatusPeriodConsumeKind != types.McomaStatusPeriodConsumeKind:

			if authConsumeStatus != commonconstantpb.MetadataUsageType_Usage_Times {
				return false, fmt.Errorf("usageType of metadataAuth not by times, metadataAuthId: {%s}, usageType: {%s}",
					auth.GetData().GetMetadataAuthId(), authConsumeStatus.String())
			}

		// there is option of consumption by time period only.
		case metadataConsumeStatus&types.McomaStatusTimesConsumeKind != types.McomaStatusTimesConsumeKind &&
			metadataConsumeStatus&types.McomaStatusPeriodConsumeKind == types.McomaStatusPeriodConsumeKind:

			if authConsumeStatus != commonconstantpb.MetadataUsageType_Usage_Period {
				return false, fmt.Errorf("usageType of metadataAuth not by time period, metadataAuthId: {%s}, usageType: {%s}",
					auth.GetData().GetMetadataAuthId(), authConsumeStatus.String())
			}
		// unknown consumption type.
		default:
			return false, fmt.Errorf("unknown consumption type on metadataAuth consume option of metadata, metadataId: {%s}, metadata consume status: {%d}",
				auth.GetData().GetAuth().GetMetadataId(), metadataConsumeStatus)
		}
		return true, nil
	}

	// check only one valid metadataAuth related one metadata? or multi metadataAuths related one metadata?
	//
	// only one valid related
	if option.GetStatus()&types.McomaStatusAuthMulti != types.McomaStatusAuthMulti {

		has, err := ma.HasValidMetadataAuth(auth.GetData().GetUserType(), auth.GetData().GetUser(), auth.GetData().GetAuth().GetOwner().GetIdentityId(), auth.GetData().GetAuth().GetMetadataId())
		if nil != err {
			return false, fmt.Errorf("cannot check if there is only one valid metadataAuth, %s, userType: {%s}, user: {%s}, metadata onwer identityId: {%s}, metadataId: {%s}",
				err, auth.GetData().GetUserType(), auth.GetData().GetUser(), auth.GetData().GetAuth().GetOwner().GetIdentityId(), auth.GetData().GetAuth().GetMetadataId())
		}
		if has {
			return false, fmt.Errorf("only one valid metadataAuth already exists, userType: {%s}, user: {%s}, metadata onwer identityId: {%s}, metadataId: {%s}",
				auth.GetData().GetUserType(), auth.GetData().GetUser(), auth.GetData().GetAuth().GetOwner().GetIdentityId(), auth.GetData().GetAuth().GetMetadataId())
		}

		//return checkUsageTypeFn(option.GetStatus(), auth.GetData().GetAuth().GetUsageRule().GetUsageType())
	}

	//// multi valid related
	//if option.GetStatus()&types.McomaStatusAuthMulti == types.McomaStatusAuthMulti {
	//	return checkUsageTypeFn(option.GetStatus(), auth.GetData().GetAuth().GetUsageRule().GetUsageType())
	//}

	return checkUsageTypeFn(option.GetStatus(), auth.GetData().GetAuth().GetUsageRule().GetUsageType())
}

func (ma *MetadataAuthority) VerifyMetadataAuthInfo(auth *types.MetadataAuthority) (bool, error) {

	if auth.GetData().GetAuditOption() != commonconstantpb.AuditMetadataOption_Audit_Pending {
		return false, ErrMetadataAuthHasAudited
	}

	if auth.GetData().GetState() != commonconstantpb.MetadataAuthorityState_MAState_Released {
		return false, ErrMetadataAuthHasNotReleased
	}

	switch auth.GetData().GetAuth().GetUsageRule().GetUsageType() {
	case commonconstantpb.MetadataUsageType_Usage_Period:
		if timeutils.UnixMsecUint64() >= auth.GetData().GetAuth().GetUsageRule().GetEndAt() {
			return false, ErrMetadataAuthHasExpired
		}
	case commonconstantpb.MetadataUsageType_Usage_Times:
		if auth.GetData().GetUsedQuo().GetUsedTimes() >= auth.GetData().GetAuth().GetUsageRule().GetTimes() {
			return false, ErrMetadataAuthHasNotEnoughTimes
		}
	default:
		return false, ErrMetadataAuthUnknoenUsageType
	}

	return true, nil
}

func (ma *MetadataAuthority) QueryMetadataAuthIdsByMetadataId(userType commonconstantpb.UserType, user, metadataId string) ([]string, error) {
	return ma.dataCenter.QueryValidUserMetadataAuthIdsByMetadataId(userType, user, metadataId)
}

func (ma *MetadataAuthority) UpdateMetadataAuthority(metadataAuth *types.MetadataAuthority) error {
	return ma.dataCenter.UpdateMetadataAuthority(metadataAuth)
}
