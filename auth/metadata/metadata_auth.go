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
	metadataAuth, err := ma.GetMetadataAuthority(audit.GetMetadataAuthId())
	if nil != err {
		log.Errorf("Failed to query old metadataAuth on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}, err: {%s}",
			audit.GetMetadataAuthId(), err)
		return apicommonpb.AuditMetadataOption_Audit_Pending, err
	}

	if metadataAuth.GetData().State != apicommonpb.MetadataAuthorityState_MAState_Released {
		log.Errorf("the old metadataAuth state is not release on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}",
			audit.GetMetadataAuthId())
		return metadataAuth.GetData().GetAuditOption(), fmt.Errorf("the old metadataAuth state is not release")
	}

	if metadataAuth.GetData().GetAuditOption() != apicommonpb.AuditMetadataOption_Audit_Pending {
		log.Errorf("the old metadataAuth has already audited on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}",
			audit.GetMetadataAuthId())
		return metadataAuth.GetData().GetAuditOption(), fmt.Errorf("the old metadataAuth has already audited")
	}

	metadataAuth.GetData().AuditOption = audit.GetAuditOption()
	metadataAuth.GetData().AuditSuggestion = audit.GetAuditSuggestion()

	err = ma.dataCenter.InsertMetadataAuthority(metadataAuth)
	if nil != err {
		log.Errorf("Failed to store metadataAuth after audit on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}, audit option:{%s}, err: {%s}",
			audit.GetMetadataAuthId(), audit.GetAuditOption().String(), err)
	}

	return metadataAuth.GetData().GetAuditOption(), err
}

func (ma *MetadataAuthority) ConsumeMetadataAuthority(metadataAuthId string) error {

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

	switch usageRule.UsageType {
	case apicommonpb.MetadataUsageType_Usage_Period:
		usedQuo.UsageType = apicommonpb.MetadataUsageType_Usage_Period
		if uint64(timeutils.UnixMsec()) >= usageRule.GetEndAt() {
			usedQuo.Expire = true
			metadataAuth.GetData().State = apicommonpb.MetadataAuthorityState_MAState_Invalid
		} else {
			usedQuo.Expire = false
		}
	case apicommonpb.MetadataUsageType_Usage_Times:
		usedQuo.UsageType = apicommonpb.MetadataUsageType_Usage_Times
		if usedQuo.GetUsedTimes() < usageRule.GetTimes() {
			usedQuo.UsedTimes += 1
		} else {
			metadataAuth.GetData().State = apicommonpb.MetadataAuthorityState_MAState_Invalid
		}
	default:
		log.Errorf("usageRule state of the old metadataAuth is invalid on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}",
			metadataAuthId)
		return fmt.Errorf("usageRule state of the old metadataAuth is invalid")
	}

	metadataAuth.GetData().UsedQuo = usedQuo

	err = ma.dataCenter.InsertMetadataAuthority(metadataAuth)
	if nil != err {
		log.Errorf("Failed to update metadataAuth after consume on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}, err: {%s}",
			metadataAuthId, err)
		return err
	}
	return nil
}

func (ma *MetadataAuthority) GetMetadataAuthority(metadataAuthId string) (*types.MetadataAuthority, error) {
	return ma.dataCenter.GetMetadataAuthority(metadataAuthId)
}

func (ma *MetadataAuthority) GetMetadataAuthorityList() (types.MetadataAuthArray, error) {
	identityId, err := ma.dataCenter.GetIdentityId()
	if nil != err {
		return nil, err
	}
	return ma.dataCenter.GetMetadataAuthorityListByIdentityId(identityId, uint64(timeutils.BeforeYearUnixMsec()))
}

func (ma *MetadataAuthority) GetMetadataAuthorityListByIds(metadataAuthIds []string) (types.MetadataAuthArray, error) {
	return ma.dataCenter.GetMetadataAuthorityListByIds(metadataAuthIds)
}

func (ma *MetadataAuthority) GetMetadataAuthorityListByUser(userType apicommonpb.UserType, user string) (types.MetadataAuthArray, error) {
	return ma.dataCenter.GetMetadataAuthorityListByUser(userType, user, uint64(timeutils.BeforeYearUnixMsec()))
}

func (ma *MetadataAuthority) HasValidLastMetadataAuth(userType apicommonpb.UserType, user, metadataId string) (bool, error) {
	metadataAuthId, err := ma.dataCenter.QueryUserMetadataAuthIdByMetadataId(userType, user, metadataId)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed to query valid user metadataAuthId used on MetadataAuthority.InvalidLastMetadataAuth(), userType: {%s}, user:{%s}, metadataId: {%s}",
			userType.String(), user, metadataId)
		return false, err
	}

	if rawdb.IsDBNotFoundErr(err) {
		return false, nil
	}

	metadataAuth, err := ma.GetMetadataAuthority(metadataAuthId)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed to query last user metadataAuth info on MetadataAuthority.InvalidLastMetadataAuth(), userType: {%s}, user:{%s}, metadataId: {%s}, metadataAuthId: {%s}",
			userType.String(), user, metadataId, metadataAuthId)
		return false, err
	}

	if rawdb.IsDBNotFoundErr(err) {
		return false, nil
	}

	//if metadataAuth.GetData().GetAuth().GetMetadataId() != metadataId {
	//	log.Errorf("the metadataId of metadataAuth and current metadataId is not same on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}",
	//		userType.String(), user, metadataId, metadataAuthId)
	//	return false
	//}
	//
	//if metadataAuth.GetData().GetUserType() != userType || metadataAuth.GetData().GetUser() != user {
	//	log.Errorf("the userType or user of metadataAuth and current userType or user is not same on MetadataAuthority.VerifyMetadataAuth(), auth userType: {%s},auth user: {%s}, userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}",
	//		metadataAuth.GetData().GetUserType().String(), metadataAuth.GetData().GetUser(), userType.String(), user, metadataId, metadataAuthId)
	//	return false
	//}

	if metadataAuth.GetData().GetState() != apicommonpb.MetadataAuthorityState_MAState_Released {
		log.Debugf("the old metadataAuth state is not release on MetadataAuthority.InvalidLastMetadataAuth(), userType: {%s}, user:{%s}, metadataId: {%s}, metadataAuthId: {%s}, state: {%s}",
			userType.String(), user, metadataId, metadataAuthId, metadataAuth.GetData().GetState().String())
		return false, nil
	}

	usageRule := metadataAuth.GetData().GetAuth().GetUsageRule()
	usedQuo := metadataAuth.GetData().GetUsedQuo()

	switch usageRule.UsageType {
	case apicommonpb.MetadataUsageType_Usage_Period:
		usedQuo.UsageType = apicommonpb.MetadataUsageType_Usage_Period
		if uint64(timeutils.UnixMsec()) >= usageRule.GetEndAt() {
			usedQuo.Expire = true
			metadataAuth.GetData().State = apicommonpb.MetadataAuthorityState_MAState_Invalid
		} else {
			usedQuo.Expire = false
		}
	case apicommonpb.MetadataUsageType_Usage_Times:
		// do nothing
	default:
		log.Errorf("usageRule state of the old metadataAuth is invalid on MetadataAuthority.InvalidLastMetadataAuth(), userType: {%s}, user:{%s}, metadataId: {%s}, metadataAuthId: {%s}",
			userType.String(), user, metadataId, metadataAuthId)
		return false, fmt.Errorf("usageRule state of the old metadataAuth is invalid")
	}

	if usedQuo.Expire == true {

		log.Debugf("the old metadataAuth was expire on MetadataAuthority.InvalidLastMetadataAuth(), userType: {%s}, user:{%s}, metadataId: {%s}, metadataAuthId: {%s}, state: {%s}",
			userType.String(), user, metadataId, metadataAuthId, metadataAuth.GetData().GetState().String())

		metadataAuth.GetData().UsedQuo = usedQuo
		if err = ma.dataCenter.InsertMetadataAuthority(metadataAuth); nil != err {
			log.Errorf("Failed to update metadataAuth after consume on MetadataAuthority.InvalidLastMetadataAuth(), userType: {%s}, user:{%s}, metadataId: {%s}, metadataAuthId: {%s}",
				userType.String(), user, metadataId, metadataAuthId)
			return false, err
		}

		return false, nil
	}
	return true, nil
}

func (ma *MetadataAuthority) HasInValidLastMetadataAuth(userType apicommonpb.UserType, user, metadataId string) (bool, error) {
	has, err := ma.HasValidLastMetadataAuth(userType, user, metadataId)
	if nil != err {
		return false, err
	}
	if has {
		return false, nil
	}
	return true, nil
}

func (ma *MetadataAuthority) VerifyMetadataAuth(userType apicommonpb.UserType, user, metadataId string) bool {

	metadataAuthId, err := ma.dataCenter.QueryUserMetadataAuthIdByMetadataId(userType, user, metadataId)
	if nil != err {
		log.Errorf("Failed to query user metadataAuthId by metadataId on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}, err: {%s}",
			userType.String(), user, metadataId, err)
		return false
	}

	// verify
	metadataAuth, err := ma.GetMetadataAuthority(metadataAuthId)
	if nil != err {
		log.Errorf("Failed to GetMetadataAuthority on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}, err: {%s}",
			userType.String(), user, metadataId, metadataAuthId, err)
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
		log.Debugf("the old metadataAuth state is not release on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}, state: {%s}",
			userType.String(), user, metadataId, metadataAuthId, metadataAuth.GetData().GetState().String())
		return false
	}

	usageRule := metadataAuth.GetData().GetAuth().GetUsageRule()
	usedQuo := metadataAuth.GetData().GetUsedQuo()

	switch usageRule.UsageType {
	case apicommonpb.MetadataUsageType_Usage_Period:
		usedQuo.UsageType = apicommonpb.MetadataUsageType_Usage_Period
		if uint64(timeutils.UnixMsec()) >= usageRule.GetEndAt() {
			usedQuo.Expire = true
			metadataAuth.GetData().State = apicommonpb.MetadataAuthorityState_MAState_Invalid
		} else {
			usedQuo.Expire = false
		}
	case apicommonpb.MetadataUsageType_Usage_Times:
		// do nothing
	default:
		log.Errorf("usageRule state of the old metadataAuth is invalid on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}",
			userType.String(), user, metadataId, metadataAuthId)
		return false
	}

	if usedQuo.Expire == true {

		log.Debugf("the old metadataAuth was expire on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user:{%s}, metadataId: {%s}, metadataAuthId: {%s}, state: {%s}",
			userType.String(), user, metadataId, metadataAuthId, metadataAuth.GetData().GetState().String())

		metadataAuth.GetData().UsedQuo = usedQuo
		if err = ma.dataCenter.InsertMetadataAuthority(metadataAuth); nil != err {
			log.Errorf("Failed to update metadataAuth after consume on MetadataAuthority.VerifyMetadataAuth(), userType: {%s}, user: {%s}, metadataId: {%s}, metadataAuthId: {%s}",
				userType.String(), user, metadataId, metadataAuthId)
			return false
		}
		return false
	}
	return true
}

func (ma *MetadataAuthority) StoreUserMetadataAuthUsed (userType apicommonpb.UserType, user, metadataAuthId string)  error {
	return ma.dataCenter.StoreUserMetadataAuthUsed(userType, user, metadataAuthId)
}

func (ma *MetadataAuthority) StoreUserMetadataAuthIdByMetadataId (userType apicommonpb.UserType, user, metadataId, metadataAuthId string) error {
	return ma.dataCenter.StoreUserMetadataAuthIdByMetadataId(userType, user, metadataId, metadataAuthId)
}

func (ma *MetadataAuthority) QueryMetadataAuthIdByMetadataId(userType apicommonpb.UserType, user, metadataId string) (string, error) {
	return ma.dataCenter.QueryUserMetadataAuthIdByMetadataId(userType, user, metadataId)
}

func (ma *MetadataAuthority) RemoveUserMetadataAuthIdByMetadataId (userType apicommonpb.UserType, user, metadataId string) error {
	return ma.dataCenter.RemoveUserMetadataAuthIdByMetadataId(userType, user, metadataId)
}
