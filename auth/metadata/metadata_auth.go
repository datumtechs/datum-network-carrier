package metadata

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/core"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type MetadataAuthority struct {
	dataCenter core.CarrierDB
}

func NewMetadataAuthority(dataCenter core.CarrierDB) *MetadataAuthority {
	return &MetadataAuthority{
		dataCenter: dataCenter,
	}
}

func (ma *MetadataAuthority) ApplyMetadataAuthority(metadataAuth *types.MetadataAuthority) error {
	return ma.dataCenter.InsertMetadataAuthority(metadataAuth)
}

func (ma *MetadataAuthority) AuditMetadataAuthority(audit *types.MetadataAuthAudit) (apicommonpb.AuditMetadataOption, error) {

	log.Debugf("Start audit metadataAuth on MetadataAuthority.AuditMetadataAuthority(), the audit: %s", audit.String())

	// verify
	//
	// query metadataAuth with metadataAuthId from dataCenter
	metadataAuth, err := ma.GetMetadataAuthority(audit.GetMetadataAuthId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query old metadataAuth on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}",
			audit.GetMetadataAuthId())
		return apicommonpb.AuditMetadataOption_Audit_Pending, fmt.Errorf("query metadataAuth failed, %s", err)
	}

	identity, err := ma.dataCenter.QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to query local identity on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}",
			audit.GetMetadataAuthId())
		return apicommonpb.AuditMetadataOption_Audit_Pending, fmt.Errorf("query local identity failed, %s", err)
	}

	if identity.GetIdentityId() != metadataAuth.GetData().GetAuth().GetOwner().GetIdentityId() {
		log.Errorf("Failed to verify local identity and identity of metadataAuth owner, they is not same identity on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}",
			audit.GetMetadataAuthId())
		return apicommonpb.AuditMetadataOption_Audit_Pending, fmt.Errorf("metadataAuth did not current own, %s", err)
	}

	// find metadataAuthId second (from local db)
	metadataAuthId, err := ma.dataCenter.QueryUserMetadataAuthIdByMetadataId(metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetMetadataId())
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed to query user metadataAuthId by metadataId on MetadataAuthority.AuditMetadataAuthority(), userType: {%s}, user: {%s}, metadataId: {%s}",
			metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetMetadataId())
		return apicommonpb.AuditMetadataOption_Audit_Pending, fmt.Errorf("query local only valid metadataAuthId by metadataId failed, %s", err)
	}

	if audit.GetMetadataAuthId() == metadataAuthId {
		log.Errorf("Failed to verify metadataAuthId. this metadataAuth have been already audit on MetadataAuthority.AuditMetadataAuthority(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}",
			metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuthId)
		return apicommonpb.AuditMetadataOption_Audit_Pending, fmt.Errorf("This metadataAuth have been already audited")
	}

	if metadataAuth.GetData().GetAuditOption() != apicommonpb.AuditMetadataOption_Audit_Pending {
		log.Errorf("the old metadataAuth has already audited on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}, audit option: {%s}",
			audit.GetMetadataAuthId(), audit.GetAuditOption().String())
		return metadataAuth.GetData().GetAuditOption(), fmt.Errorf("the old metadataAuth has already audited, %s", audit.GetAuditOption().String())
	}

	if metadataAuth.GetData().GetState() != apicommonpb.MetadataAuthorityState_MAState_Released {
		log.Errorf("the old metadataAuth state is not release on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}",
			audit.GetMetadataAuthId())
		return metadataAuth.GetData().GetAuditOption(), fmt.Errorf("the old metadataAuth state is %s", metadataAuth.GetData().GetState().String())
	}

	auditOption := audit.GetAuditOption()
	auditSuggestion := audit.GetAuditSuggestion()

	var invalid bool

	// check usageType/endTime once again before store and pushlish
	switch metadataAuth.GetData().GetAuth().GetUsageRule().GetUsageType() {
	case apicommonpb.MetadataUsageType_Usage_Period:
		if timeutils.UnixMsecUint64() >= metadataAuth.GetData().GetAuth().GetUsageRule().GetEndAt() {
			metadataAuth.GetData().GetUsedQuo().Expire = true // update state, maybe state has invalid.
			metadataAuth.GetData().State = apicommonpb.MetadataAuthorityState_MAState_Invalid
			//
			auditOption = apicommonpb.AuditMetadataOption_Audit_Refused
			auditSuggestion = "metadataAuth has expired, refused it"
			invalid = true
		}
	case apicommonpb.MetadataUsageType_Usage_Times:
		if  metadataAuth.GetData().GetUsedQuo().GetUsedTimes() >= metadataAuth.GetData().GetAuth().GetUsageRule().GetTimes() {
			metadataAuth.GetData().State = apicommonpb.MetadataAuthorityState_MAState_Invalid
			//
			auditOption = apicommonpb.AuditMetadataOption_Audit_Refused
			auditSuggestion = "metadataAuth has no enough remain times, refused it"
			invalid = true
		}
	default:
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

	if invalid {

		log.Debugf("the metadataAuth was invalid need to call `RemoveUserMetadataAuthIdByMetadataId()` on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}",
			metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetUserType(), metadataAuth.GetUser())

		// remove the invaid metadataAuthId from local db
		if err := ma.dataCenter.RemoveUserMetadataAuthIdByMetadataId(metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetMetadataId()); nil != err {
			log.WithError(err).Errorf("Failed to remove metadataId and metadataAuthId mapping while metadataAuth has invalid on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}",
				metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetUserType(), metadataAuth.GetUser())
			return metadataAuth.GetData().GetAuditOption(), fmt.Errorf("remove metadataId and invalid metadataAuthId mapping failed, %s", err)
		}
	} else {

		log.Debugf("the metadataAuth audit succeed need to call `StoreUserMetadataAuthIdByMetadataId()` on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}",
			metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetUserType(), metadataAuth.GetUser())

		// prefix + userType + user + metadataId -> metadataAuthId (only one)
		if err := ma.dataCenter.StoreUserMetadataAuthIdByMetadataId(metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetData().GetMetadataAuthId()); nil != err {
			log.WithError(err).Errorf("Failed to store metadataId and metadataAuthId mapping while metadataAuth has invalid on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}",
				metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetUserType(), metadataAuth.GetUser())
			return metadataAuth.GetData().GetAuditOption(), fmt.Errorf("store metadataId and valid metadataAuthId mapping failed, %s", err)
		}
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
	if metadataAuth.GetData().GetAuditOption() != apicommonpb.AuditMetadataOption_Audit_Passed {
		log.Errorf("the old metadataAuth audit option is not passed on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}, audit option: {%s}",
			metadataAuthId, metadataAuth.GetData().GetAuditOption().String())
		return fmt.Errorf("the old metadataAuth audit option is not passed")
	}

	// check state second
	if metadataAuth.GetData().GetState() != apicommonpb.MetadataAuthorityState_MAState_Released {
		log.Errorf("the old metadataAuth state is not released on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}, state: {%s}",
			metadataAuthId, metadataAuth.GetData().GetState().String())
		return fmt.Errorf("the old metadataAuth state is not released")
	}

	usageRule := metadataAuth.GetData().GetAuth().GetUsageRule()
	usedQuo := metadataAuth.GetData().GetUsedQuo()

	var needUpdate bool

	// check anything else next...
	switch usageRule.GetUsageType() {
	case apicommonpb.MetadataUsageType_Usage_Period:
		if timeutils.UnixMsecUint64() >= usageRule.GetEndAt() {
			usedQuo.Expire = true
			metadataAuth.GetData().State = apicommonpb.MetadataAuthorityState_MAState_Invalid
			needUpdate = true
		}
	case apicommonpb.MetadataUsageType_Usage_Times:

		usedQuo.UsedTimes += 1
		if usedQuo.GetUsedTimes() >= usageRule.GetTimes() {
			metadataAuth.GetData().State = apicommonpb.MetadataAuthorityState_MAState_Invalid
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
	if metadataAuth.GetData().GetState() == apicommonpb.MetadataAuthorityState_MAState_Invalid {

		log.Debugf("the metadataAuth was invalid need to call `RemoveUserMetadataAuthIdByMetadataId()` on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}",
			metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetUserType(), metadataAuth.GetUser())

		// remove the invaid metadataAuthId from local db
		if err := ma.dataCenter.RemoveUserMetadataAuthIdByMetadataId(metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetMetadataId()); nil != err {
			log.WithError(err).Errorf("Failed to removed metadataId and metadataAuthId mapping while metadataAuth has invalid on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}",
				metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetUserType(), metadataAuth.GetUser())
		}
	}

	return nil
}

func (ma *MetadataAuthority) GetMetadataAuthority (metadataAuthId string) (*types.MetadataAuthority, error) {
	return ma.dataCenter.QueryMetadataAuthority(metadataAuthId)
}

func (ma *MetadataAuthority) GetLocalMetadataAuthorityList() (types.MetadataAuthArray, error) {
	identityId, err := ma.dataCenter.QueryIdentityId()
	if nil != err {
		return nil, err
	}
	return ma.dataCenter.QueryMetadataAuthorityListByIdentityId(identityId, timeutils.BeforeYearUnixMsecUint64())
}

func (ma *MetadataAuthority) GetGlobalMetadataAuthorityList(lastUpdate uint64, pageSize uint64) (types.MetadataAuthArray, error) {
	return  ma.dataCenter.QueryMetadataAuthorityList(lastUpdate, pageSize)
}

func (ma *MetadataAuthority) GetMetadataAuthorityListByIds(metadataAuthIds []string) (types.MetadataAuthArray, error) {
	return ma.dataCenter.QueryMetadataAuthorityListByIds(metadataAuthIds)
}

func (ma *MetadataAuthority) HasValidMetadataAuth(userType apicommonpb.UserType, user, identityId, metadataId string) (bool, error) {


	// query metadataAuthList with target identityId (metadataId of target org)
	metadataAuthList, err := ma.dataCenter.QueryMetadataAuthorityListByIdentityId(identityId, timeutils.BeforeYearUnixMsecUint64())
	if nil != err {
		log.WithError(err).Errorf("Failed to QueryMetadataAuthorityListByIdentityId() on MetadataAuthority.HasValidMetadataAuth(), userType: {%s}, user:{%s}, identityId: {%s}, metadataId: {%s}",
			userType.String(), user, identityId,  metadataId)
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
			auth.GetData().GetState() == apicommonpb.MetadataAuthorityState_MAState_Released {

			// filter one by one

			if auth.GetData().GetAuditOption() == apicommonpb.AuditMetadataOption_Audit_Refused {
				continue
			}

			usageRule := auth.GetData().GetAuth().GetUsageRule()
			usedQuo := auth.GetData().GetUsedQuo()


			switch usageRule.GetUsageType() {
			case apicommonpb.MetadataUsageType_Usage_Period:
				if timeutils.UnixMsecUint64() >= usageRule.GetEndAt() {
					log.Warnf("the metadataAuth was expired on MetadataAuthority.HasValidMetadataAuth(), userType: {%s}, user:{%s}, identityId: {%s}, metadataId: {%s}",
						userType.String(), user,identityId,  metadataId)
					continue
				}
			case apicommonpb.MetadataUsageType_Usage_Times:
				if usedQuo.GetUsedTimes() >= usageRule.GetTimes() {
					log.Warnf("the metadataAuth was used times exceed the limit on MetadataAuthority.HasValidMetadataAuth(), userType: {%s}, user:{%s}, identityId: {%s}, metadataId: {%s}",
						userType.String(), user,identityId,  metadataId)
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

func (ma *MetadataAuthority) HasNotValidMetadataAuth(userType apicommonpb.UserType, user, ideneityId, metadataId string) (bool, error) {
	has, err := ma.HasValidMetadataAuth(userType, user, ideneityId, metadataId)
	if nil != err {
		return false, err
	}
	if has {
		return false, nil
	}
	return true, nil
}

func (ma *MetadataAuthority) VerifyMetadataAuth(userType apicommonpb.UserType, user, metadataId string) error {

	log.Debugf("Start verify metadataAuth, userType: {%s}, user: {%s}, metadataId: {%s}", userType.String(), user, metadataId)

	// If the metadata is internal metadata, no verify metadataAuth required
	flag, err := ma.dataCenter.IsInternalMetadataByDataId(metadataId)
	if nil != err {
		log.WithError(err).Errorf("Failed to check internal metadata by metadataId on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}",
			userType.String(), user, metadataId)
		return fmt.Errorf("check is internal metadata failed, %s", err)
	}
	// The internal metadata does not need to verify the authorization information.
	if flag {
		log.Debugf("The internal metadata verify the authorization information default audit `passed`, userType: {%s}, user: {%s}, metadataId: {%s}", userType.String(), user, metadataId)
		return nil
	}

	// query last metadataAuthId of metadataId with userType and user
	metadataAuthId, err := ma.dataCenter.QueryUserMetadataAuthIdByMetadataId(userType, user, metadataId)
	if nil != err {
		log.WithError(err).Errorf("Failed to query user metadataAuthId by metadataId on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}",
			userType.String(), user, metadataId)
		return fmt.Errorf("query metadataAuthId by metadataId failed, %s", err)
	}

	// checking
	//
	// query metadataAuth by metadataAuthId from dataCenter
	metadataAuth, err := ma.GetMetadataAuthority(metadataAuthId)
	if nil != err {
		log.WithError(err).Errorf("Failed to QueryMetadataAuthority on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}",
			userType.String(), user, metadataId, metadataAuthId)
		return fmt.Errorf("query metadataAuth info failed, %s", err)
	}

	if metadataAuth.GetData().GetAuth().GetMetadataId() != metadataId {
		log.Errorf("the metadataId of metadataAuth and current metadataId is not same on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}",
			userType.String(), user, metadataId, metadataAuthId)
		return fmt.Errorf("metadataId of metadataAuth and input params is defferent")
	}

	if metadataAuth.GetData().GetUserType() != userType || metadataAuth.GetData().GetUser() != user {
		log.Errorf("the userType or user of metadataAuth and current userType or user is not same on MetadataAuthority.VerifyMetadataAuth(), auth userType: {%s},auth user: {%s}, userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}",
			metadataAuth.GetData().GetUserType().String(), metadataAuth.GetData().GetUser(), userType.String(), user, metadataId, metadataAuthId)
		return fmt.Errorf("user information of metadataAuth and input params is defferent")
	}

	if metadataAuth.GetData().GetAuditOption() != apicommonpb.AuditMetadataOption_Audit_Passed {
		log.Errorf("the old metadataAuth audit option is not passed on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}, audit option: {%s}",
			userType.String(), user, metadataId, metadataAuthId, metadataAuth.GetData().GetAuditOption().String())
		return fmt.Errorf("the metadataAuth audit option is not passed")
	}

	if metadataAuth.GetData().GetState() != apicommonpb.MetadataAuthorityState_MAState_Released {
		log.Errorf("the old metadataAuth state is not release on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}, state: {%s}",
			userType.String(), user, metadataId, metadataAuthId, metadataAuth.GetData().GetState().String())
		return fmt.Errorf("the metadataAuth state is invalid")
	}

	usageRule := metadataAuth.GetData().GetAuth().GetUsageRule()
	usedQuo := metadataAuth.GetData().GetUsedQuo()

	var invalid bool

	switch usageRule.UsageType {
	case apicommonpb.MetadataUsageType_Usage_Period:
		usedQuo.UsageType = apicommonpb.MetadataUsageType_Usage_Period
		if timeutils.UnixMsecUint64() >= usageRule.GetEndAt() {
			usedQuo.Expire = true
			metadataAuth.GetData().State = apicommonpb.MetadataAuthorityState_MAState_Invalid
			invalid = true
		}
	case apicommonpb.MetadataUsageType_Usage_Times:
		if usedQuo.GetUsedTimes() >= usageRule.GetTimes() {
			metadataAuth.GetData().State = apicommonpb.MetadataAuthorityState_MAState_Invalid
			invalid = true
		}
	default:
		log.Errorf("unknown usageType of the metadataAuth on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}",
			userType.String(), user, metadataId, metadataAuthId)
		return fmt.Errorf("unknown usageType of the metadataAuth")
	}

	if invalid {

		log.Debugf("the old metadataAuth was invalid on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user:{%s}, metadataId: {%s}, metadataAuthId: {%s}, state: {%s}",
			userType.String(), user, metadataId, metadataAuthId, metadataAuth.GetData().GetState().String())

		metadataAuth.GetData().UsedQuo = usedQuo
		if err := ma.dataCenter.UpdateMetadataAuthority(metadataAuth); nil != err {
			log.WithError(err).Errorf("Failed to update metadataAuth after verify expire auth on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}",
				userType.String(), user, metadataId, metadataAuthId)
			return fmt.Errorf("update metadataAuth after verify expire auth failed, %s", err)
		}

		log.Debugf("the old metadataAuth was invalid and call succeed `UpdateMetadataAuthority()` on MetadataAuthority.VerifyMetadataAuth(), metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s} with usageRule: %s usedQuo: %s, state: {%s}",
			metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetUsageRule().String(), metadataAuth.GetData().GetUsedQuo().String(), metadataAuth.GetData().GetState())

		// remove the invaid metadataAuthId from local db
		if err := ma.dataCenter.RemoveUserMetadataAuthIdByMetadataId(metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetMetadataId()); nil != err {
			log.WithError(err).Errorf("Failed to remove metadataId and metadataAuthId mapping while metadataAuth has invalid on MetadataAuthority.VerifyMetadataAuth(), metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}",
				metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetUserType(), metadataAuth.GetUser())
		}
		return fmt.Errorf("the metadataAuth was invalid")
	}
	return nil
}

func (ma *MetadataAuthority) QueryMetadataAuthIdByMetadataId(userType apicommonpb.UserType, user, metadataId string) (string, error) {
	return ma.dataCenter.QueryUserMetadataAuthIdByMetadataId(userType, user, metadataId)
}

func (ma *MetadataAuthority) UpdateMetadataAuthority (metadataAuth *types.MetadataAuthority) error {
	return ma.dataCenter.UpdateMetadataAuthority(metadataAuth)
}

func (ma *MetadataAuthority) RemoveUserMetadataAuthIdByMetadataId(userType apicommonpb.UserType, user, metadataId string) error {
	return ma.dataCenter.RemoveUserMetadataAuthIdByMetadataId(userType, user, metadataId)
}