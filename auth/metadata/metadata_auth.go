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

	// verify
	//
	// query metadataAuth with metadataAuthId from dataCenter
	metadataAuth, err := ma.GetMetadataAuthority(audit.GetMetadataAuthId())
	if nil != err {
		log.Errorf("Failed to query old metadataAuth on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}, err: {%s}",
			audit.GetMetadataAuthId(), err)
		return apicommonpb.AuditMetadataOption_Audit_Pending, err
	}

	// find metadataAuthId second (from local db)
	metadataAuthId, err := ma.dataCenter.QueryUserMetadataAuthIdByMetadataId(metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetMetadataId())
	if rawdb.IsNoDBNotFoundErr(err) {
		log.Errorf("Failed to query user metadataAuthId by metadataId on MetadataAuthority.AuditMetadataAuthority(), userType: {%s}, user: {%s}, metadataId: {%s}, err: {%s}",
			metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetMetadataId(), err)
		return apicommonpb.AuditMetadataOption_Audit_Pending, err
	}

	if audit.GetMetadataAuthId() == metadataAuthId {
		log.Errorf("Failed to verify metadataAuthId. this metadataAuth have been already audit on MetadataAuthority.AuditMetadataAuthority(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}",
			metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuthId)
		return apicommonpb.AuditMetadataOption_Audit_Pending, fmt.Errorf("This metadataAuth have been already audit")
	}

	if metadataAuth.GetData().GetState() != apicommonpb.MetadataAuthorityState_MAState_Released {
		log.Errorf("the old metadataAuth state is not release on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}",
			audit.GetMetadataAuthId())
		return metadataAuth.GetData().GetAuditOption(), fmt.Errorf("the old metadataAuth state is not release")
	}

	if metadataAuth.GetData().GetAuditOption() != apicommonpb.AuditMetadataOption_Audit_Pending {
		log.Errorf("the old metadataAuth has already audited on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}",
			audit.GetMetadataAuthId())
		return metadataAuth.GetData().GetAuditOption(), fmt.Errorf("the old metadataAuth has already audited")
	}

	auditOption := audit.GetAuditOption()
	auditSuggestion := audit.GetAuditSuggestion()

	// check usageType/endTime once again before store and pushlish
	switch metadataAuth.GetData().GetAuth().GetUsageRule().GetUsageType() {
	case apicommonpb.MetadataUsageType_Usage_Period:
		if timeutils.UnixMsecUint64() >= metadataAuth.GetData().GetAuth().GetUsageRule().GetEndAt() {
			metadataAuth.GetData().GetUsedQuo().Expire = true // update state, maybe state has invalid.
			metadataAuth.GetData().State = apicommonpb.MetadataAuthorityState_MAState_Invalid
			//
			auditOption = apicommonpb.AuditMetadataOption_Audit_Refused
			auditSuggestion = "metadataAuth has expired, refused it"
		}
	case apicommonpb.MetadataUsageType_Usage_Times:
		// do nothing
	default:
		log.Errorf("unknown usageType of the metadataAuth MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}",
			audit.GetMetadataAuthId())
		return metadataAuth.GetData().GetAuditOption(), fmt.Errorf("unknown usageType of the metadataAuth")
	}

	// update audit things.
	metadataAuth.GetData().AuditOption = auditOption
	metadataAuth.GetData().AuditSuggestion = auditSuggestion
	metadataAuth.GetData().AuditAt = timeutils.UnixMsecUint64()

	err = ma.dataCenter.UpdateMetadataAuthority(metadataAuth)
	if nil != err {
		log.Errorf("Failed to store metadataAuth after audit on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}, audit option:{%s}, err: {%s}",
			audit.GetMetadataAuthId(), audit.GetAuditOption().String(), err)
		return metadataAuth.GetData().GetAuditOption(), fmt.Errorf("update metadataAuth failed")
	}

	// prefix + userType + user + metadataId -> metadataAuthId (only one)
	err = ma.dataCenter.StoreUserMetadataAuthIdByMetadataId(metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetData().GetMetadataAuthId())
	if nil != err {
		log.Errorf("Failed to store metadataId and metadataAuthId mapping after audit on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}, err: {%s}",
			metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetUserType(), metadataAuth.GetUser(), err)
		return metadataAuth.GetData().GetAuditOption(), fmt.Errorf("store metadataId and last metadataAuthId mapping failed")
	}

	return metadataAuth.GetData().GetAuditOption(), nil
}

func (ma *MetadataAuthority) ConsumeMetadataAuthority(metadataAuthId string) error {

	log.Debugf("Start consume metadataAuth, metadataAuthId: {%s}", metadataAuthId)

	// verify
	metadataAuth, err := ma.GetMetadataAuthority(metadataAuthId)
	if nil != err {
		log.Errorf("Failed to query old metadataAuth on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}, err: {%s}",
			metadataAuthId, err)
		return err
	}

	if metadataAuth.GetData().GetState() != apicommonpb.MetadataAuthorityState_MAState_Released {
		log.Errorf("the old metadataAuth state is not release on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}, state: {%s}",
			metadataAuthId, metadataAuth.GetData().GetState().String())
		return fmt.Errorf("the old metadataAuth state is not release")
	}

	if metadataAuth.GetData().GetAuditOption() == apicommonpb.AuditMetadataOption_Audit_Refused {
		log.Errorf("the old metadataAuth has already audited on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}",
			metadataAuthId)
		return fmt.Errorf("the old metadataAuth was refused")
	}

	usageRule := metadataAuth.GetData().GetAuth().GetUsageRule()
	usedQuo := metadataAuth.GetData().GetUsedQuo()

	switch usageRule.GetUsageType() {
	case apicommonpb.MetadataUsageType_Usage_Period:
		if timeutils.UnixMsecUint64() >= usageRule.GetEndAt() {
			usedQuo.Expire = true
			metadataAuth.GetData().State = apicommonpb.MetadataAuthorityState_MAState_Invalid
		} else {
			usedQuo.Expire = false
		}
	case apicommonpb.MetadataUsageType_Usage_Times:
		if usedQuo.GetUsedTimes() < usageRule.GetTimes() {
			usedQuo.UsedTimes += 1
		} else {
			metadataAuth.GetData().State = apicommonpb.MetadataAuthorityState_MAState_Invalid
		}
	default:
		log.Errorf("unknown usageType of the old metadataAuth on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}",
			metadataAuthId)
		return fmt.Errorf("unknown usageType of the old metadataAuth")
	}

	metadataAuth.GetData().UsedQuo = usedQuo
	err = ma.dataCenter.UpdateMetadataAuthority(metadataAuth)
	if nil != err {
		log.Errorf("Failed to update metadataAuth after consume on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}, err: {%s}",
			metadataAuthId, err)
		return err
	}

	// remove
	if metadataAuth.GetData().State == apicommonpb.MetadataAuthorityState_MAState_Invalid {
		// remove the invaid metadataAuthId from local db
		ma.dataCenter.RemoveUserMetadataAuthIdByMetadataId(metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetMetadataId())
	}
	return nil
}

func (ma *MetadataAuthority) GetMetadataAuthority(metadataAuthId string) (*types.MetadataAuthority, error) {
	return ma.dataCenter.QueryMetadataAuthority(metadataAuthId)
}

func (ma *MetadataAuthority) GetLocalMetadataAuthorityList() (types.MetadataAuthArray, error) {
	identityId, err := ma.dataCenter.QueryIdentityId()
	if nil != err {
		return nil, err
	}
	return ma.dataCenter.QueryMetadataAuthorityListByIdentityId(identityId, timeutils.BeforeYearUnixMsecUint64())
}

func (ma *MetadataAuthority) GetGlobalMetadataAuthorityList() (types.MetadataAuthArray, error) {
	return ma.dataCenter.QueryMetadataAuthorityList(timeutils.BeforeYearUnixMsecUint64())
}

func (ma *MetadataAuthority) GetMetadataAuthorityListByIds(metadataAuthIds []string) (types.MetadataAuthArray, error) {
	return ma.dataCenter.QueryMetadataAuthorityListByIds(metadataAuthIds)
}

func (ma *MetadataAuthority) HasValidMetadataAuth(userType apicommonpb.UserType, user, identityId, metadataId string) (bool, error) {

	var (
		metadataAuth *types.MetadataAuthority
	)

	// query metadataAuthList with target identityId (metadataId of target org)
	metadataAuthList, err := ma.dataCenter.QueryMetadataAuthorityListByIdentityId(identityId, timeutils.BeforeYearUnixMsecUint64())
	if nil != err {
		log.WithError(err).Errorf("Failed to QueryMetadataAuthorityListByIdentityId() on MetadataAuthority.HasValidMetadataAuth(), userType: {%s}, user:{%s}, identityId: {%s}, metadataId: {%s}",
			userType.String(), user,identityId,  metadataId)
		return false, err
	}

	if len(metadataAuthList) == 0 {
		return false, nil
	}

	// find valid metadataAauth only
	for _, auth := range metadataAuthList {
		if auth.GetUserType() == userType &&
			auth.GetUser() == user &&
			auth.GetData().GetAuth().GetMetadataId() == metadataId &&
			auth.GetData().GetState() == apicommonpb.MetadataAuthorityState_MAState_Released {
			// then find it
			metadataAuth = auth
			break
		}
	}

	// what if we can not find it short circuit finally
	if nil == metadataAuth {
		return false, nil
	}

	usageRule := metadataAuth.GetData().GetAuth().GetUsageRule()
	usedQuo := metadataAuth.GetData().GetUsedQuo()

	switch usageRule.GetUsageType() {
	case apicommonpb.MetadataUsageType_Usage_Period:
		if timeutils.UnixMsecUint64() >= usageRule.GetEndAt() {
			usedQuo.Expire = true
			metadataAuth.GetData().State = apicommonpb.MetadataAuthorityState_MAState_Invalid
		} else {
			usedQuo.Expire = false
		}
	case apicommonpb.MetadataUsageType_Usage_Times:
		// do nothing
	default:
		log.Errorf("unknown usageType of the old metadataAuth on MetadataAuthority.HasValidMetadataAuth(), userType: {%s}, user:{%s}, metadataId: {%s}, metadataAuthId: {%s}",
			userType.String(), user, metadataId, metadataAuth.GetData().GetMetadataAuthId())
		return false, fmt.Errorf("unknown usageType of the old metadataAuth")
	}

	if usedQuo.Expire == true {

		log.Debugf("the old metadataAuth was expire on MetadataAuthority.HasValidMetadataAuth(), userType: {%s}, user:{%s}, metadataId: {%s}, metadataAuthId: {%s}, state: {%s}",
			userType.String(), user, metadataId, metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetData().GetState().String())

		// update the expired metadataAuth into datacenter
		metadataAuth.GetData().UsedQuo = usedQuo
		if err = ma.dataCenter.UpdateMetadataAuthority(metadataAuth); nil != err {
			log.Errorf("Failed to update metadataAuth after consume on MetadataAuthority.HasValidMetadataAuth(), userType: {%s}, user:{%s}, metadataId: {%s}, metadataAuthId: {%s}",
				userType.String(), user, metadataId, metadataAuth.GetData().GetMetadataAuthId())
			return false, err
		}
		// remove the expired metadataAuthId from local db
		ma.dataCenter.RemoveUserMetadataAuthIdByMetadataId(userType, user, metadataId)
		return false, nil
	}

	return true, nil
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

func (ma *MetadataAuthority) VerifyMetadataAuth(userType apicommonpb.UserType, user, metadataId string) bool {

	log.Debugf("Start verify metadataAuth, userType: {%s}, user: {%s}, metadataId: {%s}", userType.String(), user, metadataId)

	// If the metadata is internal metadata, no verify metadataAuth required
	flag, err := ma.dataCenter.IsInternalMetadataByDataId(metadataId)
	if nil != err {
		log.WithError(err).Errorf("Failed to check internal metadata by metadataId on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}",
			userType.String(), user, metadataId)
		return false
	}
	// The internal metadata does not need to verify the authorization information.
	if flag {
		log.Debugf("The internal metadata verify the authorization information default `pass`, userType: {%s}, user: {%s}, metadataId: {%s}", userType.String(), user, metadataId)
		return true
	}

	// query last metadataAuthId of metadataId with userType and user
	metadataAuthId, err := ma.dataCenter.QueryUserMetadataAuthIdByMetadataId(userType, user, metadataId)
	if nil != err {
		log.WithError(err).Errorf("Failed to query user metadataAuthId by metadataId on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}",
			userType.String(), user, metadataId)
		return false
	}

	// verify
	//
	// query metadataAuth by metadataAuthId from dataCenter
	metadataAuth, err := ma.GetMetadataAuthority(metadataAuthId)
	if nil != err {
		log.WithError(err).Errorf("Failed to QueryMetadataAuthority on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}",
			userType.String(), user, metadataId, metadataAuthId)
		return false
	}

	if metadataAuth.GetData().GetAuth().GetMetadataId() != metadataId {
		log.Errorf("the metadataId of metadataAuth and current metadataId is not same on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}",
			userType.String(), user, metadataId, metadataAuthId)
		return false
	}

	if metadataAuth.GetData().GetUserType() != userType || metadataAuth.GetData().GetUser() != user {
		log.Errorf("the userType or user of metadataAuth and current userType or user is not same on MetadataAuthority.VerifyMetadataAuth(), auth userType: {%s},auth user: {%s}, userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}",
			metadataAuth.GetData().GetUserType().String(), metadataAuth.GetData().GetUser(), userType.String(), user, metadataId, metadataAuthId)
		return false
	}

	if metadataAuth.GetData().GetState() != apicommonpb.MetadataAuthorityState_MAState_Released {
		log.Errorf("the old metadataAuth state is not release on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}, state: {%s}",
			userType.String(), user, metadataId, metadataAuthId, metadataAuth.GetData().GetState().String())
		return false
	}

	usageRule := metadataAuth.GetData().GetAuth().GetUsageRule()
	usedQuo := metadataAuth.GetData().GetUsedQuo()

	switch usageRule.UsageType {
	case apicommonpb.MetadataUsageType_Usage_Period:
		usedQuo.UsageType = apicommonpb.MetadataUsageType_Usage_Period
		if timeutils.UnixMsecUint64() >= usageRule.GetEndAt() {
			usedQuo.Expire = true
			metadataAuth.GetData().State = apicommonpb.MetadataAuthorityState_MAState_Invalid
		} else {
			usedQuo.Expire = false
		}
	case apicommonpb.MetadataUsageType_Usage_Times:
		// do nothing
	default:
		log.Errorf("unknown usageType of the metadataAuth on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}",
			userType.String(), user, metadataId, metadataAuthId)
		return false
	}

	if usedQuo.Expire == true {

		log.Debugf("the old metadataAuth was expire on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user:{%s}, metadataId: {%s}, metadataAuthId: {%s}, state: {%s}",
			userType.String(), user, metadataId, metadataAuthId, metadataAuth.GetData().GetState().String())

		metadataAuth.GetData().UsedQuo = usedQuo
		if err = ma.dataCenter.UpdateMetadataAuthority(metadataAuth); nil != err {
			log.Errorf("Failed to update metadataAuth after verify expire auth on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}",
				userType.String(), user, metadataId, metadataAuthId)
			return false
		}
		// remove the invaid metadataAuthId from local db
		ma.dataCenter.RemoveUserMetadataAuthIdByMetadataId(userType, user, metadataId)
		return false
	}
	return true
}

func (ma *MetadataAuthority) QueryMetadataAuthIdByMetadataId(userType apicommonpb.UserType, user, metadataId string) (string, error) {
	return ma.dataCenter.QueryUserMetadataAuthIdByMetadataId(userType, user, metadataId)
}

